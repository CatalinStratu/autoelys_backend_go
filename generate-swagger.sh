#!/bin/bash

echo "=== Generating Swagger Documentation ==="

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Swag CLI not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest

    # Add Go bin to PATH if not already there
    export PATH=$PATH:$(go env GOPATH)/bin

    # Check again after installation
    if ! command -v swag &> /dev/null; then
        echo "Error: Failed to install swag. Please install manually:"
        echo "  go install github.com/swaggo/swag/cmd/swag@latest"
        echo "  export PATH=\$PATH:\$(go env GOPATH)/bin"
        exit 1
    fi
fi

# Generate Swagger documentation
echo "Generating Swagger docs..."
swag init

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Swagger documentation generated successfully!"
    echo ""
    echo "Generated files:"
    echo "  - docs/docs.go"
    echo "  - docs/swagger.json"
    echo "  - docs/swagger.yaml"
    echo ""
    echo "Access Swagger UI at: http://localhost:8080/swagger/index.html"
else
    echo ""
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi
