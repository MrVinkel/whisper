.DEFAULT_GOAL := build

.PHONY: setup
setup: ## Install tools and download dependencies
	@go mod download
	@go install gotest.tools/gotestsum@latest
	@go install github.com/boumenot/gocover-cobertura@latest

.PHONY: build
build: ## Build gateway and forwarder
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/whisper ./cmd

.PHONY: test
test: ## Run tests
	@gotestsum --format pkgname  -- -coverprofile=./bin/cobertura-coverage.txt -covermode count ./...
	@gocover-cobertura < ./bin/cobertura-coverage.txt > ./bin/cobertura-coverage.xml

.PHONY: run
run: build ## Run whisper
	@./bin/whisper

.PHONY: lint
lint: build ## Lint code
	@golangci-lint run ./...

.PHONY: clean
clean: ## Clean all build files
	@rm -rf bin
	@go cache clean

.PHONY: help
help: ## Shows this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'