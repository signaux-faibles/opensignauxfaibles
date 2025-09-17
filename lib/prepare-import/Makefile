.DEFAULT_GOAL := prepare-import

prepare-import: *.go prepareimport/*.go createfilter/*.go
	@make install
	@go build

install: ## Install dependencies (including development/testing dependencies)
	@go get -t ./...

test: ## Run tests, including end-to-end binary tests
	@make prepare-import
	@go test ./... -test.count=1 # prevent cache

test-update: ## Run tests and update snapshots / golden files
	@make prepare-import
	@go test main_test.go -update -test.count=1 # prevent cache and update golden file

format: ## Fix the formatting of .go files
	@go fmt

format-doc: ## Fix the formatting of md and yml files
	@npm install --save-dev --save-exact prettier
	@npx prettier --write .

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install test test-update format format-doc help
