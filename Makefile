# RimPay Makefile

.PHONY: help build test test-watch clean lint fmt vet

# Default target
help:
	@echo "Available commands:"
	@echo "  build      - Build the project"
	@echo "  test       - Run tests"
	@echo "  test-watch - Run tests in watch mode"
	@echo "  clean      - Clean build artifacts"
	@echo "  lint       - Run golangci-lint"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"

# Build the project
build:
	go build -v ./...

# Run tests
test:
	go test -v ./...

# Run tests in watch mode (requires entr or similar)
test-watch:
	find . -name "*.go" | entr -c go test -v ./...

# Clean build artifacts
clean:
	go clean
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run linter (if golangci-lint is installed)
lint:
	@which golangci-lint > /dev/null && golangci-lint run || echo "golangci-lint not installed, skipping..."
