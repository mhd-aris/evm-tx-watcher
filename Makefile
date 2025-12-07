.PHONY: build clean test run-api run-worker docker-build docker-up docker-down migrate-up migrate-down

# Build configuration
BINARY_API=bin/api
BINARY_WORKER=bin/worker
BUILD_DIR=bin

# Load environment variables from .env
include .env
export $(shell sed 's/=.*//' .env)


# Build binaries
build:
	@echo "Building binaries..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY_API) ./cmd/api
	go build -o $(BINARY_WORKER) ./cmd/worker
	@echo "Build completed: $(BINARY_API), $(BINARY_WORKER)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run API server
run-api: build
	@echo "Starting API server..."
	./$(BINARY_API)

# Run worker
run-worker: build
	@echo "Starting worker..."
	./$(BINARY_WORKER)

# Legacy commands for compatibility
run: run-api

lint:
	golangci-lint run ./...

# Database migration commands
migrate-create:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Usage: make migrate-create create_addresses_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(filter-out $@,$(MAKECMDGOALS))
%:
	@:


migrate-up:
	@echo "Running database migrations up..."
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Running database migrations down..."
	migrate -path db/migrations -database "$(DB_URL)" down

# Development setup
setup-dev:
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		if [ -f env.example ]; then \
			cp env.example .env; \
			echo "Created .env file from env.example"; \
		elif [ -f .env.example ]; then \
			cp .env.example .env; \
			echo "Created .env file from .env.example"; \
		else \
			echo "No example file found. Please create .env manually"; \
		fi; \
	else \
		echo ".env file already exists"; \
	fi
	@echo "Please configure your .env file with proper RPC URLs and database credentials"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Generate swagger docs
docs:
	@echo "Generating swagger documentation..."
	swag init -g cmd/api/main.go -o docs

# Quick development workflow
dev: clean build migrate-up
	@echo "Development environment ready!"
	@echo "Start API: make run-api"
	@echo "Start Worker: make run-worker"

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build API and Worker binaries"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  run-api      - Build and run API server"
	@echo "  run-worker   - Build and run worker"
	@echo "  setup-dev    - Setup development environment"
	@echo "  deps         - Install Go dependencies"
	@echo "  fmt          - Format Go code"
	@echo "  lint         - Run linter"
	@echo "  docs         - Generate swagger documentation"
	@echo "  migrate-up   - Run database migrations up"
	@echo "  migrate-down - Run database migrations down"
	@echo "  dev          - Quick development setup"
