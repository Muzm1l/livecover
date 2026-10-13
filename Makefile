GO ?= go

.PHONY: tidy build test run-agent run-aggregator

tidy:
	$(GO) mod tidy

build:
	$(GO) build ./cmd/agent
	$(GO) build ./cmd/aggregator
	$(GO) build ./cmd/web
	$(GO) build ./cmd/cli

test:
	$(GO) test ./...

run-agent:
	$(GO) run ./cmd/agent -config deploy/hcl/agent.hcl

run-aggregator:
	$(GO) run ./cmd/aggregator -config deploy/hcl/aggregator.hcl
