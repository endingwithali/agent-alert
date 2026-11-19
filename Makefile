.DEFAULT_GOAL := help

.PHONY: help run test build

help: ## Show the available commands
	@printf 'Agent Monitor\n\n'
	@printf 'Commands:\n'
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the session reader (loads .env automatically)
	go run ./cmd/session-reader

test: ## Run all tests
	go test ./...

build: ## Build bin/session-reader
	@mkdir -p bin
	go build -o bin/session-reader ./cmd/session-reader
