# Build Plan

## Milestone A: Local Agent Signal

- Process discovery for one target binary.
- Symbol resolution and function map extraction.
- eBPF uprobe attach for entry hits.
- Local event stream verification.

## Milestone B: Aggregation Pipeline

- Agent batch transport to aggregator.
- In-memory rollup by minute buckets.
- SQL flush for `coverage_1m` and `function_last_seen`.

## Milestone C: API + Realtime

- Snapshot API endpoints.
- WebSocket stream for deltas and process state.
- Basic heatmap client integration.

## Milestone D: Hardening

- Backpressure and drop accounting.
- Multi-host testing and reliability checks.
- Security and operations readiness.
