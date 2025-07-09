.DEFAULT_GOAL := build

build: ## Build the sfdata binary
	# Note: don't forget to run `go generate ./...` before building.
	go build -o "sfdata" -ldflags "-X main.GitCommit=$(shell git rev-parse HEAD)"

build-prod: ## Build the sfdata binary for our production environment
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(MAKE) build

test: ## Run automated tests
	./test-all.sh

test-update: ## Run automated tests, update snaphots and golden files
	./test-all.sh --update

postgres-container: ## Spin up a Postgresql container for testing (container: test-postgres, pw: testpass, port: 5432, db: testdb, user: testuser)
	@ docker run --name test-postgres \
		-e POSTGRES_DB=testdb \
		-e POSTGRES_USER=testuser \
		-e POSTGRES_PASSWORD=testpass \
		-p 5432:5432 \
		-d postgres:17

postgres-migrations:
	@ for file in ./migrations/*; do psql postgres://testuser:testpass@localhost:5432/testdb  -f "$${file}"; done

postgres-cleanup: ## Clean up the Postgresql container
	@ docker stop test-postgres
	@ docker rm test-postgres



help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build build-prod test test-update help
