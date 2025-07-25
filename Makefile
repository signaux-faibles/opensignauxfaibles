.DEFAULT_GOAL := build

build: ## Build the sfdata binary
	go build -o "sfdata" -ldflags "-X main.GitCommit=$(shell git rev-parse HEAD)"

build-prod: ## Build the sfdata binary for our production environment
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(MAKE) build

test: ## Run automated tests
	./test-all.sh

test-update: ## Run automated tests, update snaphots and golden files
	./test-all.sh --update

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

start-postgres: ## Run a postgres container for tests
	@docker run --name test-postgres \
    -e POSTGRES_DB=testdb \
    -e POSTGRES_USER=testuser \
    -e POSTGRES_PASSWORD=testpass \
    -p 5432:5432 \
    -d postgres:17

stop-postgres: ## Stop the postgrescontainer. All data will be lost
	@docker stop test-postgres
	@docker rm test-postgres

.PHONY: build build-prod test test-update help
