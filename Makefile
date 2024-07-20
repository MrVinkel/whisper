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

.PHONY: dev
dev: ## Setup dev vault with default secrets, policies and users
	$(eval EXIST := $(shell [ "$$(docker ps -a | grep dev-vault)" ] && echo true || echo false))
	$(eval RUNNING := $(shell [ "$(EXIST)" = "true" ] && docker container inspect -f '{{.State.Running}}' 'dev-vault'))
	@if [ "$(RUNNING)" = "true" ]; then \
		echo "Dev vault is already running"; \
	elif [ "$(EXIST)" = "true" ]; then \
		echo "Starting dev vault"; \
		docker restart dev-vault; \
	else \
		echo "Creating dev vault"; \
		docker run -e 'VAULT_DEV_ROOT_TOKEN_ID=root' -e 'VAULT_TOKEN=root' -e 'VAULT_ADDR=http://127.0.0.1:8200' -v ${PWD}/testdata:/testdata:ro --cap-add=IPC_LOCK -p=8200:8200 -d --name=dev-vault hashicorp/vault; \
		echo "Waiting for vault to start..."; \
		until docker exec dev-vault vault status 2>/dev/null; do sleep 1; done; \
		echo "Vault is ready"; \
		docker exec dev-vault vault auth enable userpass; \
		docker exec dev-vault vault policy write writer /testdata/writer-policy.hcl; \
		docker exec dev-vault vault policy write reader /testdata/reader-policy.hcl; \
		docker exec dev-vault vault write auth/userpass/users/reader \password=reader \policies=reader; \
		docker exec dev-vault vault write auth/userpass/users/writer \password=writer \policies=writer; \
		docker exec dev-vault vault kv put -mount=secret mysecret foo=bar hello=world; \
	fi

.PHONY: dev-oidc
dev-oidc: dev ## Setup dev vault with oidc authentication
	@docker exec dev-vault vault auth enable oidc
	@docker exec dev-vault vault write auth/oidc/config oidc_discovery_url="$(OIDC_DOMAIN)" oidc_client_id="$(OIDC_CLIENT_ID)" oidc_client_secret="$(OIDC_CLIENT_SECRET)" default_role="reader"
	@docker exec dev-vault vault write auth/oidc/role/reader bound_audiences="$(OIDC_CLIENT_ID)" allowed_redirect_uris="http://localhost:8200/ui/vault/auth/oidc/oidc/callback" allowed_redirect_uris="http://localhost:8250/oidc/callback" user_claim="sub" token_policies="reader"

.PHONY: dev-azure
dev-azure: dev ## Setup dev vault with azure authentication
	@docker exec dev-vault vault auth enable azure
	@docker exec dev-vault vault write auth/azure/config tenant_id=$(AZURE_TENANT) resource=https://management.azure.com/ client_id=$(AZURE_CLIENT_ID) client_secret=$(AZURE_CLIENT_SECRET)

.PHONY: dev-clean
dev-clean: ## Stop and remove dev vault
	@docker stop dev-vault || true
	@docker rm dev-vault || true

.PHONY: help
help: ## Shows this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'