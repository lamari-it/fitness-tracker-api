# Database Configuration
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATE_PATH=./migrations

# Load environment variables from .env file if it exists
-include .env
export

# Default values if not set in environment
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= password
DB_NAME ?= fitflow
DB_SSL_MODE ?= disable

.PHONY: build run test clean deps help migrate-up migrate-down migrate-create migrate-force migrate-version migrate-drop

help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build commands
build: ## Build the application
	go build -o bin/fitflow-api main.go

build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/fitflow-api main.go

# Run commands
run: ## Run the application
	go run main.go

dev: ## Run with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
	air

# Test and quality commands
test: ## Run tests
	go test -v ./...

fmt: ## Format code
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

# Dependency management
deps: ## Install dependencies
	go mod download
	go mod tidy

clean: ## Clean build artifacts
	rm -rf bin/

# Database migration commands
migrate-up: check-migrate ## Run all up migrations
	@echo "Running database migrations..."
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down: check-migrate ## Run all down migrations (WARNING: This will drop all tables)
	@echo "WARNING: This will drop all tables. Are you sure? [y/N]" && read ans && [ $${ans:-N} = y ]
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

migrate-down-1: check-migrate ## Rollback last migration
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-create: check-migrate ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(NAME)

migrate-force: check-migrate ## Force migration to specific version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make migrate-force VERSION=version_number"; exit 1; fi
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(VERSION)

migrate-version: check-migrate ## Show current migration version
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" version

migrate-drop: check-migrate ## Drop all tables and migration history (WARNING: Destructive)
	@echo "WARNING: This will drop all tables and migration history. Are you sure? [y/N]" && read ans && [ $${ans:-N} = y ]
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" drop

# Development database commands
dev-reset: ## Reset database for development (drop, migrate up, seed)
	@echo "Resetting development database..."
	make migrate-drop || true
	make migrate-up
	@echo "Database reset complete"

dev-fresh: ## Fresh database setup (migrate up from clean state)
	make migrate-up

# Legacy database commands (kept for compatibility)
createdb: ## Create database
	createdb fitflow

dropdb: ## Drop database
	dropdb fitflow

# Tool installation and checks
install-migrate: ## Install golang-migrate tool
	@which migrate > /dev/null || (echo "Installing golang-migrate..." && \
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)

check-migrate: ## Check if migrate tool is available
	@which migrate > /dev/null || (echo "Error: migrate tool not found. Run 'make install-migrate' first." && exit 1)