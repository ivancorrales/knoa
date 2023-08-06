GOIMPORT		= goimport
GOFUMPT			= gofumpt
GOLINTER		?= golangci-lint
GIT_HOOKS_PATH 	?= .githooks

.PHONY: git-hooks
git-hooks: ## setup githooks for the local repository
	@echo "=== $(PROJECT_NAME) === [ git-hooks        ]: Configuring git hooks..."
	git config core.hooksPath $(GIT_HOOKS_PATH)

.PHONY: fmt
fmt: ## format the code
	@echo "=== $(PROJECT_NAME) === [ fmt ]: formatting the code with goimport..."
	@#$($(GOFUMPT)  $(find . -type f -name '*.go' -not -path "./vendor/*")
	@goimports -w . && gofumpt -w .

.PHONY: lint
lint:
	@echo "=== $(PROJECT_NAME) === [ lint ]: Validating source code running $(GOLINTER)..."
	$(GOLINTER) run ./...