TOOL_DIR     		?= tools

.PHONY: tools
tools: ## install required tools
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@cd $(TOOL_DIR) && $(GO) mod download
	@cd $(TOOL_DIR) && cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
	@cd $(TOOL_DIR) && $(GO) mod tidy