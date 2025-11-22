#!/bin/bash

echo "Fixing Go dependencies..."

# Remove go.sum to force regeneration
rm -f go.sum

# Download all dependencies
go mod download

# Tidy up and regenerate go.sum
go mod tidy

# Verify dependencies
go mod verify

echo "Dependencies fixed! You can now run: go run main.go"
