# Dead Code Methodology

## Definition

A function is marked as a dead-code candidate only when no hit is observed in the configured window and observation quality gates pass.

## Quality gates

- Sufficient host coverage (`hosts_seen >= threshold`).
- Sufficient time window (for example, 30 days).
- Confidence threshold met (`high` or `medium` in MVP default).
- Function not in explicit keep-list (startup/DR/admin paths).

## Scoring

Example score:

`score = (not_seen_days * 2.0) + (host_coverage_ratio * 30) - risk_penalty`

Suggested interpretation:

- `>= 80`: high-confidence candidate
- `50-79`: medium-confidence candidate
- `< 50`: informational only

## Caveats

- Not seen in production does not guarantee safe deletion.
- Rare or feature-flagged paths require manual review.
- Disaster recovery and maintenance endpoints need exclusions.
