package agg

import (
	"context"
	"testing"
	"time"

	"livecover/internal/logging"
	"livecover/internal/model"
)

func TestMemoryEngineIngestAndSnapshot(t *testing.T) {
	engine := NewMemoryEngine(logging.Discard())
	now := time.Now().UTC()

	hits := []model.HitEvent{
		{
			TsUnixNano: uint64(now.UnixNano()),
			FuncID:     7,
			BuildID:    "sha256:abc",
			Service:    "checkout",
		},
		{
			TsUnixNano: uint64(now.Add(10 * time.Second).UnixNano()),
			FuncID:     7,
			BuildID:    "sha256:abc",
			Service:    "checkout",
		},
	}

	if err := engine.IngestHits(context.Background(), hits); err != nil {
		t.Fatalf("IngestHits() error = %v", err)
	}

	snapshot := engine.Snapshot(model.QueryFilter{
		Service: "checkout",
		BuildID: "sha256:abc",
		Window:  time.Hour,
	})

	if len(snapshot.Deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(snapshot.Deltas))
	}
	if snapshot.Deltas[0].Hits != 2 {
		t.Fatalf("expected 2 hits, got %d", snapshot.Deltas[0].Hits)
	}
	if len(snapshot.LastSeen) != 1 {
		t.Fatalf("expected 1 last_seen row, got %d", len(snapshot.LastSeen))
	}
	if snapshot.LastSeen[0].TotalHits != 2 {
		t.Fatalf("expected total_hits=2, got %d", snapshot.LastSeen[0].TotalHits)
	}
}

func TestMemoryEngineSubscribeReceivesDelta(t *testing.T) {
	engine := NewMemoryEngine(logging.Discard())
	sub, unsubscribe := engine.Subscribe(model.QueryFilter{
		Service: "checkout",
		BuildID: "sha256:abc",
	})
	defer unsubscribe()

	now := time.Now().UTC()
	err := engine.IngestHits(context.Background(), []model.HitEvent{
		{
			TsUnixNano: uint64(now.UnixNano()),
			FuncID:     7,
			BuildID:    "sha256:abc",
			Service:    "checkout",
		},
	})
	if err != nil {
		t.Fatalf("IngestHits() error = %v", err)
	}

	select {
	case delta := <-sub:
		if delta.Service != "checkout" || delta.BuildID != "sha256:abc" || delta.FuncID != 7 {
			t.Fatalf("unexpected delta: %#v", delta)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for delta")
	}
}
