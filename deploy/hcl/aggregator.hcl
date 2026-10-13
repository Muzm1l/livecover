listen_grpc = ":9090"
listen_http = ":8080"

storage {
  driver = "postgres"
  dsn = "postgres://coverage:change-me@db:5432/coverage?sslmode=disable"
  flush_interval = "30s"
  retention_days = 90
}

realtime {
  ws_path = "/v1/stream"
  snapshot_on_subscribe = true
}

deadcode {
  enabled = true
  default_window = "720h" # 30d
  min_observation_hosts = 5
  min_confidence = "medium"
}
