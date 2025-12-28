# virsh-sandbox API Makefile

.PHONY: all build run clean fmt lint vet test check generate-openapi create-test-vm setup-lima help install-tools

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
BINARY_NAME=virsh-sandbox-api
CMD_PATH=./cmd/api

# Default target
all: check build

## Build targets

build: ## Build the API binary
	$(GOBUILD) -o bin/$(BINARY_NAME) $(CMD_PATH)

run: ## Run the API server
	$(GORUN) $(CMD_PATH)

clean: ## Clean build artifacts
	rm -rf bin/
	$(GOCMD) clean

## Code quality targets

fmt: ## Format code with gofumpt
	./scripts/fmt.sh

lint: ## Run golangci-lint
	./scripts/lint.sh

vet: ## Run go vet
	./scripts/vet.sh

test: ## Run tests
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

check: fmt vet lint ## Run all code quality checks (fmt, vet, lint)

## Dependency management

deps: ## Download dependencies
	$(GOMOD) download

tidy: ## Tidy and verify dependencies
	$(GOMOD) tidy
	$(GOMOD) verify

## Code generation

generate-openapi: ## Generate OpenAPI/Swagger documentation
	./scripts/generate-openapi.sh

generate: generate-openapi ## Run all code generation

## Development setup

install-tools: ## Install development tools (gofumpt, golangci-lint, swag)
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

setup-lima: ## Setup Lima libvirt environment
	./scripts/setup-lima-libvirt.sh

create-test-vm: ## Create a test VM for development
	./scripts/create-test-vm.sh

## Docker targets

docker-build: ## Build Docker image
	docker build -f api.Dockerfile -t virsh-sandbox-api .

## Help

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
