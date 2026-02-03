.PHONY: build run test test-race test-coverage clean

# Binary name
BINARY_NAME=server

# Build the application
build:
	go build -o $(BINARY_NAME) cmd/server/main.go

# Run the application
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Run tests with race detector
test-race:
	go test -race ./...

# Run tests with coverage
test-coverage:
	go test -race -cover ./...

# Run tests with detailed coverage report
test-coverage-report:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter (if golangci-lint is installed)
lint:
	golangci-lint run

# Help
help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make run                - Build and run the application"
	@echo "  make test               - Run tests"
	@echo "  make test-race          - Run tests with race detector"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo "  make test-coverage-report - Generate HTML coverage report"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make deps               - Install dependencies"
	@echo "  make lint               - Run linter"
