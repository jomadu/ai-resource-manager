.PHONY: build build-all install test lint fmt clean install-tools setup-hooks check setup

# Build variables
BINARY_NAME := arm
BIN_DIR := bin
DIST_DIR := dist
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.buildVersion=$(VERSION) -X main.buildCommit=$(COMMIT) -X main.buildTimestamp=$(BUILD_TIME) -X github.com/jomadu/ai-rules-manager/internal/v3/version.Version=$(VERSION) -X github.com/jomadu/ai-rules-manager/internal/v3/version.Commit=$(COMMIT) -X github.com/jomadu/ai-rules-manager/internal/v3/version.BuildTime=$(BUILD_TIME) -s -w"

# Build the internal/v4 package
build:
	go build ./internal/v4/...

# Build for all platforms (internal/v4)
build-all:
	@echo "Building internal/v4 for all platforms..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build ./internal/v4/...
	GOOS=linux GOARCH=arm64 go build ./internal/v4/...
	GOOS=darwin GOARCH=amd64 go build ./internal/v4/...
	GOOS=darwin GOARCH=arm64 go build ./internal/v4/...
	GOOS=windows GOARCH=amd64 go build ./internal/v4/...

# Install binary to system PATH (disabled - internal/v4 is a library)
install: build
	@echo "Note: internal/v4 is a library package, not a binary. Nothing to install."

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./internal/v4/...

# Format code
fmt:
	gofmt -w .
	$(shell go env GOPATH)/bin/goimports -w .

# Run linter
lint:
	$(shell go env GOPATH)/bin/golangci-lint run

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR) $(DIST_DIR) coverage.out .venv

# Install development tools
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Setup pre-commit hooks
setup-hooks:
	python3 -m venv .venv
	.venv/bin/pip install pre-commit
	.venv/bin/pre-commit install
	.venv/bin/pre-commit install --hook-type commit-msg

# Run all checks
check: fmt lint test

# Development setup
setup: install-tools setup-hooks
	go mod tidy
