package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"translator/internal/agg"
	"translator/internal/config"
	"translator/internal/logging"
	"translator/internal/model"
)

func main() {
	configPath := flag.String("config", "", "Path to aggregator HCL config file")
	logLevel := flag.String("log-level", "info", "Log level: debug|info|warn|error")
	logFormat := flag.String("log-format", "text", "Log format: text|json")
	flag.Parse()

	logger := logging.New(logging.Options{
		Level:     *logLevel,
		Format:    *logFormat,
		Component: "aggregator",
	})
	logging.SetDefault(logger)

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	cfg, err := config.LoadAggregatorConfig(*configPath)
	if err != nil {
		logger.Error("failed to load aggregator config", "path", *configPath, "err", err)
		os.Exit(1)
	}

	engine := agg.NewMemoryEngine(logger)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool { return true },
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/v1/ingest", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if r.Method != http.MethodPost {
			logger.Warn("ingest rejected: method not allowed", "method", r.Method, "remote", r.RemoteAddr)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Events []model.HitEvent `json:"events"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Warn("ingest rejected: invalid json", "err", err, "remote", r.RemoteAddr)
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		if len(req.Events) == 0 {
			logger.Warn("ingest rejected: empty events", "remote", r.RemoteAddr)
			http.Error(w, "events cannot be empty", http.StatusBadRequest)
			return
		}
		if err := engine.IngestHits(r.Context(), req.Events); err != nil {
			logger.Error("ingest failed", "err", err, "events", len(req.Events), "remote", r.RemoteAddr)
			http.Error(w, "failed to ingest events", http.StatusInternalServerError)
			return
		}

		services := uniqueServices(req.Events)
		logger.Info("ingest accepted",
			"events", len(req.Events),
			"services", services,
			"remote", r.RemoteAddr,
			"duration_ms", time.Since(start).Milliseconds(),
		)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"accepted": len(req.Events),
			"dropped":  0,
		})
	})
	mux.HandleFunc("/v1/coverage/snapshot", func(w http.ResponseWriter, r *http.Request) {
		filter, err := queryFilterFromRequest(r)
		if err != nil {
			logger.Warn("snapshot rejected: invalid window", "err", err, "remote", r.RemoteAddr)
			http.Error(w, "invalid window duration", http.StatusBadRequest)
			return
		}
		snapshot := engine.Snapshot(filter)
		logger.Info("snapshot served",
			"service", filter.Service,
			"build_id", filter.BuildID,
			"window", filter.Window.String(),
			"deltas", len(snapshot.Deltas),
			"last_seen", len(snapshot.LastSeen),
			"remote", r.RemoteAddr,
		)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(snapshot)
	})
	mux.HandleFunc(cfg.Realtime.WSPath, func(w http.ResponseWriter, r *http.Request) {
		filter, err := queryFilterFromRequest(r)
		if err != nil {
			logger.Warn("ws rejected: invalid window", "err", err, "remote", r.RemoteAddr)
			http.Error(w, "invalid window duration", http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("ws upgrade failed", "err", err, "remote", r.RemoteAddr)
			return
		}
		defer conn.Close()

		logger.Info("ws connected",
			"path", cfg.Realtime.WSPath,
			"service", filter.Service,
			"build_id", filter.BuildID,
			"remote", r.RemoteAddr,
		)
		defer logger.Info("ws disconnected",
			"service", filter.Service,
			"build_id", filter.BuildID,
			"remote", r.RemoteAddr,
		)

		sub, unsubscribe := engine.Subscribe(filter)
		defer unsubscribe()

		snapshot := engine.Snapshot(filter)
		snapshotMessage := map[string]any{
			"type":    "coverage.snapshot",
			"ts":      time.Now().UTC(),
			"payload": snapshot,
		}
		if err := conn.WriteJSON(snapshotMessage); err != nil {
			logger.Warn("ws snapshot write failed", "err", err, "remote", r.RemoteAddr)
			return
		}
		logger.Debug("ws snapshot sent",
			"deltas", len(snapshot.Deltas),
			"last_seen", len(snapshot.LastSeen),
			"remote", r.RemoteAddr,
		)

		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		conn.SetPongHandler(func(_ string) error {
			return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		})

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()

		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		var deltasSent int
		for {
			select {
			case <-done:
				logger.Debug("ws read loop ended", "deltas_sent", deltasSent, "remote", r.RemoteAddr)
				return
			case delta, ok := <-sub:
				if !ok {
					logger.Debug("ws subscription closed", "deltas_sent", deltasSent, "remote", r.RemoteAddr)
					return
				}
				msg := map[string]any{
					"type":    "coverage.delta",
					"ts":      time.Now().UTC(),
					"payload": delta,
				}
				if err := conn.WriteJSON(msg); err != nil {
					logger.Warn("ws delta write failed", "err", err, "remote", r.RemoteAddr)
					return
				}
				deltasSent++
				logger.Debug("ws delta sent",
					"service", delta.Service,
					"build_id", delta.BuildID,
					"func_id", delta.FuncID,
					"hits", delta.Hits,
					"remote", r.RemoteAddr,
				)
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(2*time.Second)); err != nil {
					logger.Debug("ws ping failed", "err", err, "remote", r.RemoteAddr)
					return
				}
			}
		}
	})

	handler := logging.Middleware(logger)(mux)

	logger.Info("aggregator started",
		"http", cfg.ListenHTTP,
		"grpc", cfg.ListenGRPC,
		"ws_path", cfg.Realtime.WSPath,
		"storage_driver", cfg.Storage.Driver,
		"deadcode_enabled", cfg.Deadcode.Enabled,
		"log_level", *logLevel,
		"log_format", *logFormat,
	)
	if err := http.ListenAndServe(cfg.ListenHTTP, handler); err != nil {
		logger.Error("aggregator http serve failed", "err", err)
		os.Exit(1)
	}
}

func queryFilterFromRequest(r *http.Request) (model.QueryFilter, error) {
	q := r.URL.Query()
	filter := model.QueryFilter{
		Service: q.Get("service"),
		BuildID: q.Get("build_id"),
	}
	if rawWindow := q.Get("window"); rawWindow != "" {
		window, err := time.ParseDuration(rawWindow)
		if err != nil {
			return model.QueryFilter{}, err
		}
		filter.Window = window
	}
	return filter, nil
}

func uniqueServices(events []model.HitEvent) []string {
	seen := make(map[string]struct{}, len(events))
	out := make([]string, 0, 4)
	for _, ev := range events {
		if ev.Service == "" {
			continue
		}
		if _, ok := seen[ev.Service]; ok {
			continue
		}
		seen[ev.Service] = struct{}{}
		out = append(out, ev.Service)
	}
	return out
}
