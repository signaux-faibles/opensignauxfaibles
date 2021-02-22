.DEFAULT_GOAL := build

build: ## Build the sfdata binary
	# Note: don't forget to run `go generate ./...` before building.
	GOARCH=amd64 go build -o "sfdata" -ldflags "-X main.GitCommit=$(shell git rev-parse HEAD)"
	@ # Other note: GOARCH=amd64 environment variable is necessary for the binary to work in production environment

test: ## Run automated tests
	./test-all.sh

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build test help
