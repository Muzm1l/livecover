package config

import (
	"errors"
	"fmt"
	"time"
)

type AgentConfig struct {
	NodeID        string `hcl:"node_id"`
	Environment   string `hcl:"environment"`
	AggregatorURL string `hcl:"aggregator_url"`

	Discovery struct {
		Mode             string   `hcl:"mode"`
		PollInterval     string   `hcl:"poll_interval"`
		IncludeServices  []string `hcl:"include_services,optional"`
		ExcludeBinaries  []string `hcl:"exclude_binaries,optional"`
	} `hcl:"discovery,block"`

	Instrumentation struct {
		MaxProbesPerProcess int  `hcl:"max_probes_per_process"`
		MaxEventsPerSecond  int  `hcl:"max_events_per_second"`
		AttachTimeout       string `hcl:"attach_timeout"`
		DropOnBackpressure  bool `hcl:"drop_on_backpressure"`
	} `hcl:"instrumentation,block"`

	Security struct {
		RequireAllowlist bool     `hcl:"require_allowlist"`
		AllowExecPaths   []string `hcl:"allow_exec_paths,optional"`
	} `hcl:"security,block"`
}

type AggregatorConfig struct {
	ListenGRPC string `hcl:"listen_grpc"`
	ListenHTTP string `hcl:"listen_http"`

	Storage struct {
		Driver       string `hcl:"driver"`
		DSN          string `hcl:"dsn"`
		FlushInterval string `hcl:"flush_interval"`
		RetentionDays int    `hcl:"retention_days"`
	} `hcl:"storage,block"`

	Realtime struct {
		WSPath              string `hcl:"ws_path"`
		SnapshotOnSubscribe bool   `hcl:"snapshot_on_subscribe"`
	} `hcl:"realtime,block"`

	Deadcode struct {
		Enabled             bool   `hcl:"enabled"`
		DefaultWindow       string `hcl:"default_window"`
		MinObservationHosts int    `hcl:"min_observation_hosts"`
		MinConfidence       string `hcl:"min_confidence"`
	} `hcl:"deadcode,block"`
}

func (c AgentConfig) Validate() error {
	if c.NodeID == "" || c.AggregatorURL == "" || c.Environment == "" {
		return errors.New("agent config requires node_id, environment, and aggregator_url")
	}
	if c.Discovery.Mode != "procfs" && c.Discovery.Mode != "k8s" {
		return fmt.Errorf("discovery.mode must be procfs or k8s, got %q", c.Discovery.Mode)
	}
	if _, err := time.ParseDuration(c.Discovery.PollInterval); err != nil {
		return fmt.Errorf("invalid discovery.poll_interval: %w", err)
	}
	if _, err := time.ParseDuration(c.Instrumentation.AttachTimeout); err != nil {
		return fmt.Errorf("invalid instrumentation.attach_timeout: %w", err)
	}
	if c.Instrumentation.MaxProbesPerProcess <= 0 || c.Instrumentation.MaxEventsPerSecond <= 0 {
		return errors.New("instrumentation limits must be positive")
	}
	if c.Security.RequireAllowlist && len(c.Security.AllowExecPaths) == 0 {
		return errors.New("security.allow_exec_paths is required when require_allowlist=true")
	}
	return nil
}

func (c AggregatorConfig) Validate() error {
	if c.ListenGRPC == "" || c.ListenHTTP == "" {
		return errors.New("aggregator config requires listen_grpc and listen_http")
	}
	if c.Storage.Driver == "" || c.Storage.DSN == "" {
		return errors.New("storage.driver and storage.dsn are required")
	}
	if _, err := time.ParseDuration(c.Storage.FlushInterval); err != nil {
		return fmt.Errorf("invalid storage.flush_interval: %w", err)
	}
	if c.Storage.RetentionDays < 1 {
		return errors.New("storage.retention_days must be >= 1")
	}
	if c.Realtime.WSPath == "" {
		return errors.New("realtime.ws_path is required")
	}
	if _, err := time.ParseDuration(c.Deadcode.DefaultWindow); err != nil {
		return fmt.Errorf("invalid deadcode.default_window: %w", err)
	}
	switch c.Deadcode.MinConfidence {
	case "high", "medium", "low":
	default:
		return fmt.Errorf("deadcode.min_confidence must be high, medium, or low (got %q)", c.Deadcode.MinConfidence)
	}
	return nil
}
