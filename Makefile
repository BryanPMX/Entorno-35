.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down setup

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
	@if [ -f .env ]; then \
		set -a && . ./.env && set +a && cd cmd/api && go run main.go; \
	else \
		cd cmd/api && go run main.go; \
	fi

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html

docker-up: ## Start Docker containers (uses docker-compose.prod.yml; set DB_PASSWORD etc. for local)
	@echo "Starting Docker containers..."
	@docker-compose -f docker-compose.prod.yml up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose -f docker-compose.prod.yml down

migrate-up: ## Run database migrations up (requires DB_URL env var)
	@echo "Running migrations up..."
	@if [ -z "$$DB_URL" ]; then \
		echo "Error: DB_URL environment variable is required."; \
		echo "Example: export DB_URL=postgres://user:pass@host:5432/dbname?sslmode=disable"; \
		exit 1; \
	fi
	@migrate -path ./migrations -database "$$DB_URL" up

migrate-down: ## Run database migrations down (requires DB_URL env var)
	@echo "Running migrations down..."
	@if [ -z "$$DB_URL" ]; then \
		echo "Error: DB_URL environment variable is required."; \
		echo "Example: export DB_URL=postgres://user:pass@host:5432/dbname?sslmode=disable"; \
		exit 1; \
	fi
	@migrate -path ./migrations -database "$$DB_URL" down

setup: ## Initial project setup
	@echo "Setting up project..."
	@go mod download
	@go mod tidy

seed: ## Seed database with NOM-035 questions (requires DB_* env vars; run after migrate)
	@echo "Seeding database with questions..."
	@go run cmd/seeder/main.go

