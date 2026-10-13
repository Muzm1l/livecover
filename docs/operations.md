# Operations

## Runtime requirements

- Linux kernel with eBPF and uprobes support.
- Agent privileges: CAP_BPF/CAP_PERFMON or platform-equivalent.
- Access to target process metadata and executable paths.

## Rollout strategy

1. Single-host validation in staging.
2. Canary rollout to low-risk production nodes.
3. Expand by service with per-service overhead checks.
4. Enable dead-code reporting only after stable observation windows.

## SLO targets

- Agent CPU overhead per node: under 2%.
- Agent memory footprint: under 150MB.
- New process attach latency: under 5s.
- Event loss: under 0.1% per 5-minute window.

## Metrics

- `coverage_events_received_total`
- `coverage_events_dropped_total`
- `probe_attach_failures_total`
- `probe_active_count`
- `flush_duration_ms`
- `stream_subscriber_count`

## Logging

Both `agent` and `aggregator` use structured logging via Go `slog`.

### Flags

- `-log-level=debug|info|warn|error` (default: `info`)
- `-log-format=text|json` (default: `text`)

### Examples

```bash
go run ./cmd/aggregator -config deploy/hcl/aggregator.hcl -log-level=debug -log-format=json
go run ./cmd/agent -config deploy/hcl/agent.hcl -demo-send -log-level=debug
```

### What is logged

- Startup config summary (addresses, ws path, instrumentation limits)
- HTTP request middleware: method, path, status, duration_ms, bytes, remote
- Ingest accept/reject with event counts and service list
- Snapshot serve counts (deltas / last_seen)
- WebSocket connect/disconnect and write failures
- Live delta drops for slow subscribers (warn)
- Agent batch send success/failure with duration
