package agg

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"translator/internal/model"
)

type bucketKey struct {
	service string
	buildID string
	funcID  uint64
	bucket  int64
}

type funcKey struct {
	service string
	buildID string
	funcID  uint64
}

type memoryEngine struct {
	mu        sync.RWMutex
	buckets   map[bucketKey]uint64
	lastSeen  map[funcKey]model.LastSeen
	subs      map[uint64]subscriber
	nextSubID uint64
	logger    *slog.Logger

	ingestedTotal atomic.Uint64
	droppedDeltas atomic.Uint64
}

type subscriber struct {
	filter model.QueryFilter
	ch     chan model.CoverageDelta
}

func NewMemoryEngine(logger *slog.Logger) Engine {
	if logger == nil {
		logger = slog.Default()
	}
	return &memoryEngine{
		buckets:  make(map[bucketKey]uint64),
		lastSeen: make(map[funcKey]model.LastSeen),
		subs:     make(map[uint64]subscriber),
		logger:   logger.With("component", "agg"),
	}
}

func (e *memoryEngine) IngestHits(_ context.Context, hits []model.HitEvent) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	deltas := make([]model.CoverageDelta, 0, len(hits))
	for _, hit := range hits {
		bucketUnix := int64(hit.TsUnixNano / uint64(time.Minute.Nanoseconds()))
		bk := bucketKey{
			service: hit.Service,
			buildID: hit.BuildID,
			funcID:  hit.FuncID,
			bucket:  bucketUnix,
		}
		e.buckets[bk]++
		deltas = append(deltas, model.CoverageDelta{
			Service:    hit.Service,
			BuildID:    hit.BuildID,
			FuncID:     hit.FuncID,
			BucketUnix: bucketUnix * 60,
			Hits:       1,
		})

		ts := time.Unix(0, int64(hit.TsUnixNano)).UTC()
		fk := funcKey{
			service: hit.Service,
			buildID: hit.BuildID,
			funcID:  hit.FuncID,
		}
		current, ok := e.lastSeen[fk]
		if !ok {
			e.lastSeen[fk] = model.LastSeen{
				Service:    hit.Service,
				BuildID:    hit.BuildID,
				FuncID:     hit.FuncID,
				LastSeenTs: ts,
				TotalHits:  1,
			}
			continue
		}
		current.TotalHits++
		if ts.After(current.LastSeenTs) {
			current.LastSeenTs = ts
		}
		e.lastSeen[fk] = current
	}
	e.ingestedTotal.Add(uint64(len(hits)))
	e.publishDeltasLocked(deltas)

	e.logger.Debug("ingested hits",
		"count", len(hits),
		"subscribers", len(e.subs),
		"ingested_total", e.ingestedTotal.Load(),
		"dropped_deltas_total", e.droppedDeltas.Load(),
	)

	return nil
}

func (e *memoryEngine) Snapshot(filter model.QueryFilter) model.Snapshot {
	cutoff := time.Time{}
	if filter.Window > 0 {
		cutoff = time.Now().Add(-filter.Window)
	}
	cutoffBucket := cutoff.Unix() / 60

	e.mu.RLock()
	defer e.mu.RUnlock()

	deltas := make([]model.CoverageDelta, 0, len(e.buckets))
	for key, hits := range e.buckets {
		if filter.Service != "" && key.service != filter.Service {
			continue
		}
		if filter.BuildID != "" && key.buildID != filter.BuildID {
			continue
		}
		if !cutoff.IsZero() && key.bucket < cutoffBucket {
			continue
		}
		deltas = append(deltas, model.CoverageDelta{
			Service:    key.service,
			BuildID:    key.buildID,
			FuncID:     key.funcID,
			BucketUnix: key.bucket * 60,
			Hits:       uint32(hits),
		})
	}

	lastSeen := make([]model.LastSeen, 0, len(e.lastSeen))
	for _, value := range e.lastSeen {
		if filter.Service != "" && value.Service != filter.Service {
			continue
		}
		if filter.BuildID != "" && value.BuildID != filter.BuildID {
			continue
		}
		if !cutoff.IsZero() && value.LastSeenTs.Before(cutoff) {
			continue
		}
		lastSeen = append(lastSeen, value)
	}

	sort.Slice(deltas, func(i, j int) bool {
		if deltas[i].Service != deltas[j].Service {
			return deltas[i].Service < deltas[j].Service
		}
		if deltas[i].BuildID != deltas[j].BuildID {
			return deltas[i].BuildID < deltas[j].BuildID
		}
		if deltas[i].FuncID != deltas[j].FuncID {
			return deltas[i].FuncID < deltas[j].FuncID
		}
		return deltas[i].BucketUnix < deltas[j].BucketUnix
	})

	sort.Slice(lastSeen, func(i, j int) bool {
		if lastSeen[i].Service != lastSeen[j].Service {
			return lastSeen[i].Service < lastSeen[j].Service
		}
		if lastSeen[i].BuildID != lastSeen[j].BuildID {
			return lastSeen[i].BuildID < lastSeen[j].BuildID
		}
		return lastSeen[i].FuncID < lastSeen[j].FuncID
	})

	e.logger.Debug("snapshot built",
		"service", filter.Service,
		"build_id", filter.BuildID,
		"window", filter.Window.String(),
		"deltas", len(deltas),
		"last_seen", len(lastSeen),
	)

	return model.Snapshot{
		GeneratedAt: time.Now().UTC(),
		Service:     filter.Service,
		BuildID:     filter.BuildID,
		Deltas:      deltas,
		LastSeen:    lastSeen,
	}
}

func (e *memoryEngine) Subscribe(filter model.QueryFilter) (<-chan model.CoverageDelta, func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	id := e.nextSubID
	e.nextSubID++
	ch := make(chan model.CoverageDelta, 256)
	e.subs[id] = subscriber{
		filter: filter,
		ch:     ch,
	}

	e.logger.Info("subscriber connected",
		"sub_id", id,
		"service", filter.Service,
		"build_id", filter.BuildID,
		"active_subscribers", len(e.subs),
	)

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			e.mu.Lock()
			defer e.mu.Unlock()
			sub, ok := e.subs[id]
			if ok {
				delete(e.subs, id)
				close(sub.ch)
				e.logger.Info("subscriber disconnected",
					"sub_id", id,
					"active_subscribers", len(e.subs),
				)
			}
		})
	}
	return ch, unsubscribe
}

func (e *memoryEngine) publishDeltasLocked(deltas []model.CoverageDelta) {
	var dropped uint64
	for _, sub := range e.subs {
		for _, delta := range deltas {
			if !matchesFilter(delta, sub.filter) {
				continue
			}
			select {
			case sub.ch <- delta:
			default:
				// Drop when subscriber is slow to preserve ingest path.
				dropped++
			}
		}
	}
	if dropped > 0 {
		e.droppedDeltas.Add(dropped)
		e.logger.Warn("dropped live deltas for slow subscribers",
			"dropped", dropped,
			"dropped_deltas_total", e.droppedDeltas.Load(),
			"subscribers", len(e.subs),
		)
	}
}

func matchesFilter(delta model.CoverageDelta, filter model.QueryFilter) bool {
	if filter.Service != "" && delta.Service != filter.Service {
		return false
	}
	if filter.BuildID != "" && delta.BuildID != filter.BuildID {
		return false
	}
	return true
}
