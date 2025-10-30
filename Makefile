.PHONY: help build run test docker-up docker-down docker-dev docker-logs docker-rebuild migrate lint

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make docker-up     - Start Docker containers (production)"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-dev    - Start Docker containers with hot reload (Air)"
	@echo "  make docker-logs   - View Docker container logs"
	@echo "  make docker-rebuild - Rebuild Docker images"
	@echo "  make migrate       - Run database migrations"
	@echo "  make lint          - Run linters"

build:
	go build -o bin/nutrient-api cmd/api/main.go

run:
	go run cmd/api/main.go server --config=configs/config.dev.yaml

test:
	go test -v -race -coverprofile=coverage.out ./...

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

docker-dev:
	@echo "Starting Docker containers with hot reload (Air)..."
	docker compose -f deployments/docker-compose.yml up --build

docker-logs:
	docker compose -f deployments/docker-compose.yml logs -f api

docker-rebuild:
	docker compose -f deployments/docker-compose.yml build --no-cache api

migrate:
	go run cmd/api/main.go migrate --config=configs/config.dev.yaml

lint:
	golangci-lint run
