.PHONY: run test vet lint build clean fmt help
.DEFAULT_GOAL := help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Start the server
	go run . server

test: ## Run tests
	go test -v -cover ./...

vet: ## go vet
	go vet ./...

fmt: ## go fmt
	go fmt ./...

lint: ## golangci-lint
	golangci-lint run --timeout=5m

build: ## Build the binary
	go build -o bin/SERVICENAME .

clean:
	go clean ./...
	rm -rf bin/
