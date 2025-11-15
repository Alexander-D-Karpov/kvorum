.PHONY: help build run-api run-worker migrate-up migrate-down test clean docker-build docker-up docker-down

help:
	@echo "Available targets:"
	@echo "  build         - Build all binaries"
	@echo "  run-api       - Run API server"
	@echo "  run-worker    - Run background worker"
	@echo "  migrate-up    - Apply database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker images"
	@echo "  docker-up     - Start services with Docker Compose"
	@echo "  docker-down   - Stop services with Docker Compose"

build:
	@echo "Building binaries..."
	@go build -o bin/api ./cmd/api
	@go build -o bin/worker ./cmd/worker
	@go build -o bin/migrator ./cmd/migrator
	@echo "Done!"

run-api:
	@echo "Starting API server..."
	@go run ./cmd/api/main.go

run-worker:
	@echo "Starting worker..."
	@go run ./cmd/worker/main.go

migrate-up:
	@echo "Applying migrations..."
	@go run ./cmd/migrator/main.go up

migrate-down:
	@echo "Rolling back migrations..."
	@go run ./cmd/migrator/main.go down

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Done!"

docker-build:
	@echo "Building Docker images..."
	@docker-compose build

docker-up:
	@echo "Starting services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping services..."
	@docker-compose down

lint:
	@echo "Running linters..."
	@golangci-lint run ./...

.DEFAULT_GOAL := help