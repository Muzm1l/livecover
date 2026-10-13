node_id        = "node-a1"
environment    = "prod"
aggregator_url = "grpc://coverage-agg:9090"

discovery {
  mode = "procfs" # procfs | k8s
  poll_interval = "5s"
  include_services = ["checkout", "payments"]
  exclude_binaries = ["/usr/bin/*", "/bin/*"]
}

instrumentation {
  max_probes_per_process = 5000
  max_events_per_second  = 200000
  attach_timeout         = "3s"
  drop_on_backpressure   = true
}

security {
  require_allowlist = true
  allow_exec_paths  = ["/srv/apps/*"]
}
