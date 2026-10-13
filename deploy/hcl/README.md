# HCL Configuration

This project is configured with HashiCorp Configuration Language (`.hcl`).

## Files

- `agent.hcl`: host-level daemon configuration
- `aggregator.hcl`: central service configuration

## Validation rules

Validation is implemented in `internal/config/types.go` and enforced at startup for both binaries.

## Notes

- Duration values use Go duration syntax (`5s`, `30s`, `720h`).
- `deadcode.default_window` is stored as a duration string; prefer explicit hours for portability.
