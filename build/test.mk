COVERAGE_FILE   ?= coverage.out
COVERAGE_REPORT	?= coverage.html
COVERAGE_XML	?= coverage.xml
TEST_TIMEOUT  	= 300

.PHONY: test-init
test-init: ## download temporary resources for tests
	sh $(PWD)/scripts/setup-testdata.sh

.PHONY: test
test: ## run unit tests
	@echo "=== $(PROJECT_NAME) === [ test ]: running unit tests..."
	@$(GO) test -timeout $(TEST_TIMEOUT)s -race -v ./...

.PHONY: test-coverage
test-coverage: ## run unit tests with coverage
	@echo "=== $(PROJECT_NAME) === [ test-coverage ]: running unit tests with coverage..."
	@$(GO) test -timeout $(TEST_TIMEOUT)s -race -v -covermode=atomic -coverpkg=./... -coverprofile $(COVERAGE_FILE) ./...
	@$(GO) tool cover -func $(COVERAGE_FILE)
	@$(GO) tool cover -html $(COVERAGE_FILE) -o $(COVERAGE_REPORT)
	@gocov convert $(COVERAGE_FILE) | gocov-xml > $(COVERAGE_XML)


.PHONY: test-clean
test-clean: ## remove coverage files
	@rm -f $(COVERAGE_FILE) $(COVERAGE_REPORT)


