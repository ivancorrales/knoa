VENDOR_CMD	= $(GO_CMD) mod vendor

.PHONY: deps
deps: ## install dependencies
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@$(VENDOR_CMD)