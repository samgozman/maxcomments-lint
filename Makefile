BINARY := ./bin/custom-gcl

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the custom golangci-lint binary (reads .custom-gcl.yml)
	golangci-lint custom

.PHONY: run
run: ## Lint this repo with the custom binary
	$(BINARY) run

.PHONY: test
test: ## Run the Go test suite
	go test ./... -v

.PHONY: tidy
tidy: ## Tidy module dependencies
	go mod tidy
