# Makefile for Go CRUD API

# Variables
BINARY_NAME=bin/api
MAIN_FILE=cmd/api/main.go
DOCKER_COMPOSE_FILE=docker-compose.yml

.PHONY: all build run test clean docker-up docker-down verify help

all: build

# Build the application
build:
	@echo "Building..."
	@mkdir -p bin
	@go build -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application locally (default .env)
run:
	@echo "Running with default .env..."
	@go run $(MAIN_FILE)

# Run with .env.dev
run-dev:
	@echo "Running with .env.dev..."
	@go run $(MAIN_FILE) -env=.env.dev

# Run with .env.local
run-local:
	@echo "Running with .env.local..."
	@go run $(MAIN_FILE) -env=.env.local

# Hot Reloading with Air
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed. Install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

watch-dev:
	@if command -v air > /dev/null; then \
	    air -- -env=.env.dev; \
	    echo "Watching with .env.dev...";\
	else \
	    echo "Air not installed. Run 'make watch' to install."; \
	fi

watch-local:
	@if command -v air > /dev/null; then \
	    air -- -env=.env.local; \
	    echo "Watching with .env.local...";\
	else \
	    echo "Air not installed. Run 'make watch' to install."; \
	fi

# Generate Swagger Docs
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go

# Run tests
test:
	@echo "Testing..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Run verification script
verify:
	@echo "Verifying..."
	@chmod +x verify.sh
	@./verify.sh

# Database Migrations
migrate:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go -env=.env

migrate-dev:
	@echo "Running migrations (dev)..."
	@go run cmd/migrate/main.go -env=.env.dev

migrate-reset:
	@echo "Resetting database..."
	@go run cmd/migrate/main.go -env=.env -action=reset

migrate-reset-dev:
	@echo "Resetting database (dev)..."
	@go run cmd/migrate/main.go -env=.env.dev -action=reset

# Show help
help:
	@echo "Available targets:"
	@echo "  make build             - Build the binary"
	@echo "  make run               - Run the application locally"
	@echo "  make test              - Run tests"
	@echo "  make clean             - Remove binary"
	@echo "  make docker-up         - Start database container"
	@echo "  make docker-down       - Stop database container"
	@echo "  make verify            - Run verification script"
	@echo "  make migrate           - Run database migrations"
	@echo "  make migrate-reset     - Reset database (drop & migrate)"
	@echo "  make migrate-reset-dev - Reset dev database"
