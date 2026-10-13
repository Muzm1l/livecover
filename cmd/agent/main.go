package main

import (
	"context"
	"flag"
	"os"
	"time"

	"livecover/internal/agentpipe"
	"livecover/internal/config"
	"livecover/internal/logging"
	"livecover/internal/model"
)

func main() {
	configPath := flag.String("config", "", "Path to agent HCL config file")
	demoSend := flag.Bool("demo-send", false, "Send one synthetic batch to an HTTP ingest endpoint")
	demoIngestURL := flag.String("demo-ingest-url", "http://localhost:8080/v1/ingest", "HTTP ingest endpoint for demo sender")
	logLevel := flag.String("log-level", "info", "Log level: debug|info|warn|error")
	logFormat := flag.String("log-format", "text", "Log format: text|json")
	flag.Parse()

	logger := logging.New(logging.Options{
		Level:     *logLevel,
		Format:    *logFormat,
		Component: "agent",
	})
	logging.SetDefault(logger)

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	cfg, err := config.LoadAgentConfig(*configPath)
	if err != nil {
		logger.Error("failed to load agent config", "path", *configPath, "err", err)
		os.Exit(1)
	}

	logger.Info("agent started",
		"node_id", cfg.NodeID,
		"environment", cfg.Environment,
		"discovery_mode", cfg.Discovery.Mode,
		"aggregator_url", cfg.AggregatorURL,
		"max_probes_per_process", cfg.Instrumentation.MaxProbesPerProcess,
		"max_events_per_second", cfg.Instrumentation.MaxEventsPerSecond,
	)

	if !*demoSend {
		logger.Info("agent idle; use -demo-send to post a synthetic hit batch for end-to-end validation")
		return
	}

	sender := agentpipe.NewHTTPSender(*demoIngestURL, logger)
	now := time.Now().UTC()
	batch := model.Batch{
		NodeID: cfg.NodeID,
		Hits: []model.HitEvent{
			{
				TsUnixNano: uint64(now.UnixNano()),
				PID:        uint32(os.Getpid()),
				FuncID:     1,
				BuildID:    "demo-build",
				Service:    "demo-service",
				CPU:        0,
				NodeID:     cfg.NodeID,
			},
		},
	}
	logger.Info("sending demo batch",
		"ingest_url", *demoIngestURL,
		"events", len(batch.Hits),
		"service", "demo-service",
		"build_id", "demo-build",
	)
	if err := sender.Send(context.Background(), batch); err != nil {
		logger.Error("demo send failed", "err", err)
		os.Exit(1)
	}
	logger.Info("demo send succeeded", "events", len(batch.Hits), "ingest_url", *demoIngestURL)
}
