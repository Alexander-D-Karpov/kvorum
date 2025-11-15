.PHONY: help up down logs migrate-up migrate-down clean build test

help:
	@echo "Available commands:"
	@echo "  make up           - Start all services"
	@echo "  make down         - Stop all services"
	@echo "  make logs         - Show logs"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback migrations"
	@echo "  make clean        - Clean all data"
	@echo "  make build        - Rebuild all images"
	@echo "  make test         - Run tests"

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

migrate-up:
	docker-compose run --rm migrator up

migrate-down:
	docker-compose run --rm migrator down

clean:
	docker-compose down -v
	docker system prune -f

build:
	docker-compose build

test:
	cd backend && go test -v ./...
	cd frontend && npm test

.DEFAULT_GOAL := help