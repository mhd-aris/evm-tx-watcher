# EVM Transaction Watcher Makefile

.PHONY: run build test lint clean help deps docker-up docker-down

# Default target
help:
	@echo "Available commands:"
	@echo "  run         - Run the application in development mode"
	@echo "  build       - Build the application binary"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linter"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download dependencies"
	@echo "  docker-up   - Start PostgreSQL with docker-compose"
	@echo "  docker-down - Stop docker containers"

# Run application
run:
	go run cmd/api/main.go

# Build application
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Start PostgreSQL with Docker
docker-up:
	docker run --name evm-tx-watcher-db \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=evm_tx_watcher \
		-p 5432:5432 \
		-d postgres:15

# Stop Docker containers
docker-down:
	docker stop evm-tx-watcher-db || true
	docker rm evm-tx-watcher-db || true

# Setup database (run migrations)
db-setup:
	psql -h localhost -U postgres -d evm_tx_watcher -f migrations/001_initial_schema.sql

# Development setup
setup: deps docker-up
	@echo "Waiting for database to be ready..."
	@sleep 5
	@make db-setup
	@echo "Setup complete! Run 'make run' to start the server."
