.PHONY: help build run test clean migrate-up migrate-down docker-up docker-down lint dev fmt vet

# Variables
APP_NAME=go-vuln
BUILD_DIR=./bin
MAIN_PATH=./cmd/server
MIGRATION_PATH=./migrations
DOCKER_COMPOSE=docker compose -f docker/docker-compose.yml

# Colors for terminal output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Show this help message
	@echo '$(GREEN)Available targets:$(NC)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build the application
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

run: ## Run the application
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	@go run $(MAIN_PATH)/main.go

dev: ## Run with hot reload (requires air)
	@echo "$(GREEN)Starting development server with hot reload...$(NC)"
	@air

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -rf tmp
	@echo "$(GREEN)Clean complete$(NC)"

migrate-up: ## Run database migrations
	@echo "$(GREEN)Running database migrations...$(NC)"
	@if [ ! -f .env ]; then echo "$(YELLOW)Warning: .env file not found$(NC)"; fi
	@PGPASSWORD=$${DB_PASSWORD:-vulnpilot} psql -h $${DB_HOST:-localhost} -U $${DB_USER:-vulnpilot} -d $${DB_NAME:-vulnpilot_db} -f $(MIGRATION_PATH)/000001_init.up.sql
	@echo "$(GREEN)Migrations complete$(NC)"

migrate-down: ## Rollback database migrations
	@echo "$(GREEN)Rolling back database migrations...$(NC)"
	@PGPASSWORD=$${DB_PASSWORD:-vulnpilot} psql -h $${DB_HOST:-localhost} -U $${DB_USER:-vulnpilot} -d $${DB_NAME:-vulnpilot_db} -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "$(GREEN)Rollback complete$(NC)"

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "$(YELLOW)Usage: make migrate-create NAME=migration_name$(NC)"; exit 1; fi
	@echo "$(GREEN)Creating migration: $(NAME)$(NC)"
	@touch $(MIGRATION_PATH)/$(shell date +%Y%m%d%H%M%S)_$(NAME).up.sql
	@touch $(MIGRATION_PATH)/$(shell date +%Y%m%d%H%M%S)_$(NAME).down.sql
	@echo "$(GREEN)Migration files created$(NC)"

docker-up: ## Start Docker containers (PostgreSQL, Redis)
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Containers started$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432$(NC)"
	@echo "$(YELLOW)Redis: localhost:6379$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Containers stopped$(NC)"

docker-logs: ## View Docker container logs
	@$(DOCKER_COMPOSE) logs -f

docker-clean: ## Remove Docker containers and volumes
	@echo "$(GREEN)Removing Docker containers and volumes...$(NC)"
	@$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)Cleanup complete$(NC)"

lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Format complete$(NC)"

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Vet complete$(NC)"

deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

setup: docker-up deps migrate-up ## Complete setup (Docker + deps + migrations)
	@echo "$(GREEN)Setup complete!$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Copy .env.example to .env and configure your settings"
	@echo "  2. Run 'make dev' to start the development server"

install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Tools installed$(NC)"

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "$(GREEN)All checks passed!$(NC)"

db-shell: ## Connect to PostgreSQL database
	@PGPASSWORD=$${DB_PASSWORD:-vulnpilot} psql -h $${DB_HOST:-localhost} -U $${DB_USER:-vulnpilot} -d $${DB_NAME:-vulnpilot_db}

redis-cli: ## Connect to Redis
	@redis-cli -h $${REDIS_HOST:-localhost} -p $${REDIS_PORT:-6379}

.DEFAULT_GOAL := help
