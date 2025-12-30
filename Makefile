.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building backend..."
	@cd cmd/api && go build -o ../../bin/api

run: ## Run the application
	@echo "Running backend..."
	@cd cmd/api && go run main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path ./migrations -database "postgres://entorno35:entorno35@localhost:5432/entorno35?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@migrate -path ./migrations -database "postgres://entorno35:entorno35@localhost:5432/entorno35?sslmode=disable" down

setup: ## Initial project setup
	@echo "Setting up project..."
	@go mod download
	@go mod tidy

