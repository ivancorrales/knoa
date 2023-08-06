ROJECT_NAME := github.com/ivancorrales/mapify
GO_VERSION := $(shell cat .go-version)
GO  = GOFLAGS=-mod=readonly go
GO_CMD          ?= go
UNAME := $(shell uname)

all: init fmt test lint

help: ## Display this help screen
	@echo "Makefile targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: git-hooks tools deps  ## setup hooks & install tools and deps

.PHONY: clean
clean: test-clean ## remove the temporary resources
	@echo "=== $(PROJECT_NAME) === [ clean ]: removing binaries and coverage file..."

include build/tools.mk
include build/deps.mk
include build/code.mk
include build/test.mk