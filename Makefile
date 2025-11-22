.PHONY: setup swagger run migrate-up migrate-down migrate-status help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

setup: ## Install dependencies and generate Swagger docs
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Generating Swagger documentation..."
	swag init
	@echo "Creating .env file if it doesn't exist..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env - please update with your credentials"; fi
	@echo "Setup complete!"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	swag init
	@echo "Swagger docs generated successfully!"

run: ## Run the application
	go run main.go

migrate-up: ## Run all pending migrations
	go run main.go migrate:up

migrate-down: ## Rollback last migration
	go run main.go migrate:down

migrate-status: ## Check current migration version
	go run main.go migrate:status

build: ## Build the application
	go build -o bin/autoelys_backend main.go

clean: ## Clean build artifacts and generated docs
	rm -rf bin/
	rm -rf docs/

test: ## Run tests
	go test -v ./...

dev: swagger run ## Generate swagger and run the application
