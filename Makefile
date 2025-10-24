.PHONY: help build run test docker-up docker-down migrate lint

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make docker-up   - Start Docker containers"
	@echo "  make docker-down - Stop Docker containers"
	@echo "  make migrate     - Run database migrations"
	@echo "  make lint        - Run linters"

build:
	go build -o bin/nutrient-api cmd/api/main.go

run:
	go run cmd/api/main.go server --config=configs/config.dev.yaml

test:
	go test -v -race -coverprofile=coverage.out ./...

docker-up:
	docker-compose -f deployments/docker-compose.yml up -d

docker-down:
	docker-compose -f deployments/docker-compose.yml down

migrate:
	go run cmd/api/main.go migrate --config=configs/config.dev.yaml

lint:
	golangci-lint run
