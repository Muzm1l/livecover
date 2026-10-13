package agentpipe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"livecover/internal/model"
)

type HTTPSender struct {
	endpoint string
	client   *http.Client
	logger   *slog.Logger
}

func NewHTTPSender(endpoint string, logger *slog.Logger) *HTTPSender {
	if logger == nil {
		logger = slog.Default()
	}
	return &HTTPSender{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger.With("component", "agentpipe"),
	}
}

func (s *HTTPSender) Send(ctx context.Context, batch model.Batch) error {
	start := time.Now()
	body := struct {
		Events []model.HitEvent `json:"events"`
	}{
		Events: batch.Hits,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal batch: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	s.logger.Debug("sending hit batch",
		"endpoint", s.endpoint,
		"node_id", batch.NodeID,
		"events", len(batch.Hits),
		"bytes", len(payload),
	)

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("send batch failed",
			"endpoint", s.endpoint,
			"node_id", batch.NodeID,
			"events", len(batch.Hits),
			"err", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return fmt.Errorf("send batch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		s.logger.Error("ingest endpoint rejected batch",
			"endpoint", s.endpoint,
			"status", resp.StatusCode,
			"events", len(batch.Hits),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return fmt.Errorf("ingest endpoint returned status %d", resp.StatusCode)
	}

	s.logger.Info("hit batch sent",
		"endpoint", s.endpoint,
		"node_id", batch.NodeID,
		"events", len(batch.Hits),
		"status", resp.StatusCode,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return nil
}
