# WebSocket Contract

## Endpoint

- `ws://<host>/v1/stream?service=<name>&build_id=<id>&window=24h`
- path is configurable with `realtime.ws_path` in `aggregator.hcl`.

## Envelope

```json
{
  "type": "coverage.delta",
  "ts": "2026-05-11T19:00:00Z",
  "payload": {}
}
```

## Message types

- `coverage.snapshot`: first full state after subscribe.
- `coverage.delta`: incremental hit updates.
- `process.state`: attach/detach and process lifecycle events.
- `deadcode.update`: dead-code candidate updates.
- `health.notice`: stream or source degradation notices.

## Delivery behavior

- At-least-once for deltas.
- Client should deduplicate by `(service, build_id, func_id, bucket_unix)`.
- Server may coalesce bursts to reduce network overhead.
- Server sends a `coverage.snapshot` message first, then streams `coverage.delta`.
