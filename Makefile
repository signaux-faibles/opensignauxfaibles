.DEFAULT_GOAL := build

build: ## Build the sfdata binary
	# Note: don't forget to run `go generate ./...` before building.
	go build -o "sfdata" -ldflags "-X main.GitCommit=$(shell git rev-parse HEAD)"

build-prod: ## Build the sfdata binary for our production environment
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(MAKE) build

test: ## Run automated tests
	./test-all.sh

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build build-prod test help
