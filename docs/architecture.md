# Architecture

## Goal

Collect live, production function coverage from Go binaries without recompilation, then expose real-time and historical coverage views.

## Components

- **Agent (`cmd/agent`)**: discovers target processes, resolves symbols, attaches eBPF uprobes, batches hit events.
- **Aggregator (`cmd/aggregator`)**: accepts event batches, maintains rollups, persists data, serves API/WebSocket.
- **Web (`cmd/web`)**: serves UI/static assets and optional API gateway.
- **Store (`internal/store`)**: persistence abstraction (Postgres in MVP).

## Data flow

1. Agent discovers a target process.
2. Agent resolves function symbols for the process build.
3. Agent attaches uprobes for function entry points.
4. eBPF emits hit events to user space.
5. Agent batches and sends events to aggregator.
6. Aggregator updates in-memory counters and flushes persistent buckets.
7. API and WebSocket endpoints serve snapshots and live deltas.

## Boundaries

- Kernel/eBPF specifics are isolated in `internal/bpf`.
- Symbol extraction logic is isolated in `internal/symbols`.
- Business logic for dead-code scoring is isolated in `internal/deadcode`.
