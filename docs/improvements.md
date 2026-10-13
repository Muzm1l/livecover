# Improvement Backlog

## Done (this iteration)

- Structured logging (`slog`) with level/format flags
- HTTP request middleware with timing and status
- Ingest / snapshot / websocket lifecycle logs
- Slow-subscriber delta drop warnings
- Agent sender success/failure timing logs

## High priority next

1. **gRPC ingest** — replace HTTP demo transport for production backpressure.
2. **Prometheus metrics** — promote log counters (`ingested_total`, `dropped_deltas`) to `/v1/metrics`.
3. **Persistent store flush** — wire Postgres adapter + flush ticker from config.
4. **Request IDs** — propagate `X-Request-ID` across agent → aggregator → logs.
5. **Authn/Authz** — protect ingest and stream endpoints (mTLS or token).
6. **Symbol resolution MVP** — ELF symbol table for one Go binary.
7. **Process discovery** — procfs watcher with allowlist enforcement.
8. **Circuit breaker** — detach probes when host CPU overhead exceeds SLO.
9. **Dead-code scorer** — implement keep-list + host-coverage gates.
10. **UI heatmap** — consume `/v1/stream` and color by hit density.

## Reliability / correctness

- Coalesce per-bucket deltas before publish (reduce WS fanout).
- Persist last-seen across restarts.
- Build-id keyed function catalogs.
- Explicit event-loss accounting in ingest ACKs (not always `dropped: 0`).
