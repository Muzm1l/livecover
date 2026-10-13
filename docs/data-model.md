# Data Model

## SQL schema (MVP)

```sql
create table builds (
  build_id text primary key,
  service text not null,
  version text,
  first_seen timestamptz not null default now()
);

create table functions (
  build_id text not null references builds(build_id) on delete cascade,
  func_id bigint not null,
  name text not null,
  package text,
  file text,
  line int,
  confidence text not null check (confidence in ('high','medium','low')),
  primary key (build_id, func_id)
);

create table coverage_1m (
  service text not null,
  build_id text not null,
  func_id bigint not null,
  bucket_ts timestamptz not null,
  hits bigint not null,
  hosts_seen int not null default 1,
  primary key (service, build_id, func_id, bucket_ts)
);

create table function_last_seen (
  service text not null,
  build_id text not null,
  func_id bigint not null,
  last_seen_ts timestamptz not null,
  total_hits bigint not null default 0,
  primary key (service, build_id, func_id)
);

create index idx_cov_service_bucket on coverage_1m(service, bucket_ts desc);
create index idx_cov_build_bucket on coverage_1m(build_id, bucket_ts desc);
create index idx_last_seen_service on function_last_seen(service, last_seen_ts desc);
```

## Retention

- Keep minute-level coverage rows for 90 days by default.
- Optionally downsample old windows to hourly aggregates.

## Compaction

- Merge duplicate bucket updates per flush cycle.
- Write last-seen updates idempotently with upsert semantics.
