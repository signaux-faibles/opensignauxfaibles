.DEFAULT_GOAL := help

test: ## Run some tests
	@npx typescript dbmongo/js/common/raison_sociale.ts
	@cd dbmongo/js/test/ && ./test_common.sh
	@echo "✅ Tests passed."

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test help
