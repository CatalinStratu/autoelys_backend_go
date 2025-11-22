#!/bin/bash

echo "=== AutoElys Backend Setup ==="

# Install dependencies
echo "Installing Go dependencies..."
go mod tidy

# Install swag if not installed
if ! command -v swag &> /dev/null; then
    echo "Installing Swag CLI..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger docs
echo "Generating Swagger documentation..."
swag init

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    echo "Please update .env with your database credentials"
fi

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Next steps:"
echo "1. Update .env with your database credentials"
echo "2. Run migrations: go run main.go migrate:up"
echo "3. Start server: go run main.go"
echo ""
echo "Swagger UI will be available at: http://localhost:8080/swagger/index.html"
