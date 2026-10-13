# Runbooks

## No events arriving

1. Verify agent has loaded config successfully.
2. Check target process discovery output.
3. Confirm aggregator endpoint reachability.
4. Inspect `probe_attach_failures_total` and permission errors.

## High dropped events

1. Inspect `coverage_events_dropped_total`.
2. Raise batch frequency or reduce event volume.
3. Lower probes-per-process cap for hottest binaries.
4. Confirm aggregator ingest and storage flush latency.

## Low symbol confidence

1. Check binary build includes required symbol/debug metadata.
2. Validate build-id extraction and resolver fallback behavior.
3. Exclude low-confidence results from dead-code scoring.

## Post-upgrade regression

1. Compare kernel and eBPF support matrix.
2. Run attach smoke tests on canary host.
3. Revert to previous agent image if attach failures spike.
