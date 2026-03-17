# ==============================================================================
# Makefile — common development tasks
#
# Teaching note: Makefiles are a great way to standardize commands across
# developer machines AND CI pipelines. Your GitHub Actions workflows can call
# `make test` instead of `go test ./...`, so the same command works everywhere.
# ==============================================================================

.PHONY: all build test lint fmt vet coverage run clean help

# Default target: build and test
all: build test

## build: Compile the server binary into ./bin/
build:
	@echo "→ Building..."
	go build -o bin/server ./cmd/server

## test: Run all tests
test:
	@echo "→ Running tests..."
	go test ./...

## test-verbose: Run all tests with verbose output
test-verbose:
	@echo "→ Running tests (verbose)..."
	go test -v ./...

## coverage: Run tests and generate a coverage report
coverage:
	@echo "→ Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "→ Coverage report written to coverage.html"

## lint: Run golangci-lint
lint:
	@echo "→ Linting..."
	golangci-lint run ./...

## fmt: Format all Go source files
fmt:
	@echo "→ Formatting..."
	gofmt -w .

## vet: Run go vet
vet:
	@echo "→ Running go vet..."
	go vet ./...

## run: Start the server locally on port 8080
run:
	@echo "→ Starting server on :8080..."
	go run ./cmd/server

## clean: Remove build artifacts
clean:
	@echo "→ Cleaning..."
	rm -rf bin/ coverage.out coverage.html

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'
