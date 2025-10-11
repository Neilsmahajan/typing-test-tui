SHELL := /bin/sh

APP_NAME ?= typing-test-tui
BIN_DIR ?= bin
BIN := $(BIN_DIR)/$(APP_NAME)
PKG := .
BUILD_FLAGS ?=
LDFLAGS ?=
ARGS ?=
GOFMT ?= gofmt

.DEFAULT_GOAL := help

.PHONY: help build run test lint fmt tidy clean

help: ## Show available Make targets
	@printf "\nUsage: make <target> [VARIABLE=value]\n\n"
	@awk 'BEGIN {FS = ":.*##"; printf "Targets:\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@printf "\nVariables:\n"
	@printf "  APP_NAME=%s\n  BIN_DIR=%s\n  BUILD_FLAGS=%s\n  LDFLAGS=%s\n  ARGS=%s\n\n" "$(APP_NAME)" "$(BIN_DIR)" "$(BUILD_FLAGS)" "$(LDFLAGS)" "$(ARGS)"

build: ## Compile the TUI binary into ./$(BIN_DIR)
	@mkdir -p $(BIN_DIR)
	go build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BIN) $(PKG)

run: ## Run the TUI locally (pass CLI args with ARGS='--mode quote')
	go run $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" $(PKG) $(ARGS)

test: ## Run all Go tests
	go test ./...

lint: ## Run go vet and optional golangci-lint (if installed)
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed; skipping"; \
	fi

fmt: ## Format Go source files in place
	$(GOFMT) -w $(shell find . -type f -name '*.go' -not -path './vendor/*')

tidy: ## Sync go.mod and go.sum
	go mod tidy

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
