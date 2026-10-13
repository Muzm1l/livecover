# Security

## Principles

- Default deny for target binaries.
- Least privilege for all daemon processes.
- Strong service/environment isolation for event streams.

## Controls

- Executable allowlist in `agent.hcl` (`security.allow_exec_paths`).
- Reject instrumentation outside explicit allowlist.
- Audit log for instrumentation start/stop actions.
- Redaction option for symbol names if telemetry leaves trusted boundary.

## Risks

- Coverage telemetry can reveal internals (symbol names, package paths).
- Misconfigured allowlists may instrument unintended binaries.
- Kernel capability drift can break expected privilege boundaries.

## Recommended hardening

- Run agent in dedicated runtime profile.
- Rotate credentials used by transport/storage.
- Gate access to dead-code reports by role.
