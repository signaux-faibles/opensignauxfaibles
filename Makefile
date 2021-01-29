.DEFAULT_GOAL := build

build: ## Build the sfdata binary
	# TODO: go generate ./...
	go build -o "sfdata"

test: ## Run automated tests
	./test-all.sh

format: ## Fix the formatting of .go files
	go fmt

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build test format help
