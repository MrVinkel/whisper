.DEFAULT_GOAL := build

ifdef CI_COMMIT_TAG
VERSION := $(CI_COMMIT_TAG)
else
VERSION := dev
endif
VERSION_PACKAGE := gitlab.com/mr_vinkel/whisper/cmd/whisper

.PHONY: setup
setup: ## Install tools and download dependencies
	@go mod download
	@go install gotest.tools/gotestsum@latest
	@go install github.com/boumenot/gocover-cobertura@latest

.PHONY: build
build: ## Build gateway and forwarder
	@$(eval VERSIONFLAGS=-X '$(VERSION_PACKAGE).Version=$(VERSION)')
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s $(VERSIONFLAGS)" -o ./bin/whisper ./cmd

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

.PHONY: dev
dev:
	@docker run -e 'VAULT_DEV_ROOT_TOKEN_ID=potato' --cap-add=IPC_LOCK -p=8200:8200 -d --name=dev-vault hashicorp/vault

.PHONY: dev-clean
dev-clean:
	@docker stop dev-vault
	@docker rm dev-vault

.PHONY: help
help: ## Shows this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'