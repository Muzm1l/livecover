# HTTP API

## Health

- `GET /v1/health`
  - returns `200 OK` with body `ok` when service is live.

## Ingest

- `POST /v1/ingest`
  - accepts a JSON batch of hit events for MVP transport compatibility.

Request body:

```json
{
  "events": [
    {
      "tsUnixNano": 1715530800000000000,
      "pid": 1234,
      "funcID": 42,
      "buildID": "sha256:abcd",
      "service": "checkout",
      "cpu": 1,
      "nodeID": "node-a1"
    }
  ]
}
```

Response body:

```json
{
  "accepted": 1,
  "dropped": 0
}
```

## Snapshot

- `GET /v1/coverage/snapshot?service=<name>&build_id=<id>&window=24h`
  - returns aggregated in-memory bucket deltas and last-seen rows.

## Planned endpoints

- `GET /v1/services`
- `GET /v1/builds?service=<name>`
- `GET /v1/functions/<func_id>/history?service=<name>&window=24h`
- `GET /v1/dead-candidates?service=<name>&window=30d`

## Response conventions

- RFC3339 timestamps in JSON.
- Pagination for list endpoints via `limit` and `cursor`.
- Error shape:

```json
{
  "error": {
    "code": "invalid_argument",
    "message": "service is required"
  }
}
```
