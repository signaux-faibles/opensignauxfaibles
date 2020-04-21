.DEFAULT_GOAL := help

test: ## Run some tests
	@echo "Transpiling TypeScript files, and generating the jsFunctions.go bundle..."
	@cd dbmongo/lib/engine && go generate -x
	@echo "Running tests against the JS files (including the ones transpiled from TS)..."
	@cd dbmongo/js/test/ && ./test_common.sh
	@echo "✅ Tests passed."

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test help
