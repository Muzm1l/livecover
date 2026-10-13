# Live Coverage Server (Go + eBPF)

This repository is a starter implementation scaffold for a live production code coverage system:

- `agent`: host daemon that discovers processes, resolves symbols, and emits function hit events.
- `aggregator`: ingests events, builds coverage rollups, and serves API/WebSocket streams.
- `web`: minimal HTTP endpoint placeholder for UI/static assets.
- `cli`: operations/debug command entrypoint.

The project is configuration-driven with HCL (`.hcl`) files and includes full architecture and operations documentation under `docs/`.

## Quick start

1. Create configs from examples in `deploy/hcl/`.
2. Build binaries:
   - `go build ./cmd/agent`
   - `go build ./cmd/aggregator`
   - `go build ./cmd/web`
   - `go build ./cmd/cli`
3. Run:
   - `./agent -config deploy/hcl/agent.hcl`
   - `./aggregator -config deploy/hcl/aggregator.hcl`

## MVP API available now

- `GET /v1/health`
- `POST /v1/ingest` for JSON hit batches
- `GET /v1/coverage/snapshot?service=<name>&build_id=<id>&window=24h`
- `GET /v1/stream` (WebSocket) for `coverage.snapshot` + `coverage.delta`

## End-to-end demo

1. Start aggregator:
   - `go run ./cmd/aggregator -config deploy/hcl/aggregator.hcl`
2. Send one synthetic event from agent:
   - `go run ./cmd/agent -config deploy/hcl/agent.hcl -demo-send -demo-ingest-url http://localhost:8080/v1/ingest`
3. Query snapshot:
   - `curl "http://localhost:8080/v1/coverage/snapshot?service=demo-service&build_id=demo-build&window=1h"`

## Logging

Use `-log-level` and `-log-format` on agent/aggregator:

```bash
go run ./cmd/aggregator -config deploy/hcl/aggregator.hcl -log-level=debug -log-format=json
```

## Current status

This is an MVP scaffold with a working in-memory ingest and snapshot path in the aggregator. eBPF probe attachment, symbol parsing, gRPC ingestion, WebSocket streaming, and persistent storage adapters are intentionally left behind interfaces for incremental implementation.
