# Makefile for tmux-agent

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod
BINARY_NAME=tmux-agent
BINARY_DIR=bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/server

# Build for production (static, stripped)
.PHONY: build-prod
build-prod:
	@echo "Building $(BINARY_NAME) for production..."
	@mkdir -p $(BINARY_DIR)
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/server

# Run the server
.PHONY: run
run: build
	./$(BINARY_DIR)/$(BINARY_NAME)

# Run with configuration file
.PHONY: run-config
run-config: build
	./$(BINARY_DIR)/$(BINARY_NAME) -config config.yaml

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Lint (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@rm -f audit.log audit.log.*

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Create example config if it doesn't exist
.PHONY: config
config:
	@if [ ! -f config.yaml ]; then \
		cp config.example.yaml config.yaml; \
		echo "Created config.yaml from config.example.yaml"; \
	else \
		echo "config.yaml already exists"; \
	fi

# Show version
.PHONY: version
version:
	@echo "Version: $(VERSION)"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t tmux-agent:$(VERSION) .

# Help
.PHONY: help
help:
	@echo "tmux-agent Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make              Build the binary"
	@echo "  make build        Build the binary"
	@echo "  make build-prod   Build for production (static, stripped)"
	@echo "  make run          Build and run the server"
	@echo "  make run-config   Build and run with config.yaml"
	@echo "  make test         Run tests"
	@echo "  make test-coverage Run tests with coverage report"
	@echo "  make vet          Run go vet"
	@echo "  make fmt          Format code"
	@echo "  make lint         Run golangci-lint"
	@echo "  make tidy         Tidy dependencies"
	@echo "  make deps         Download dependencies"
	@echo "  make clean        Clean build artifacts"
	@echo "  make install      Install binary to GOPATH/bin"
	@echo "  make config       Create config.yaml from example"
	@echo "  make version      Show version"
	@echo "  make docker-build Build Docker image"
	@echo "  make help         Show this help"
