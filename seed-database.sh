#!/bin/bash

echo "=== Database Seeding Script ==="

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Database connection details
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME:-autoelys_backend}

# Check if mysql is available
if ! command -v mysql &> /dev/null; then
    echo "Error: mysql command not found. Please install MySQL client."
    exit 1
fi

echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo ""

# Function to run seed file
run_seed() {
    local seed_file=$1
    local table_name=$2

    if [ ! -f "$seed_file" ]; then
        echo "⚠️  Seed file not found: $seed_file"
        return 1
    fi

    echo "Seeding $table_name..."
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$seed_file"

    if [ $? -eq 0 ]; then
        # Count records
        count=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -D "$DB_NAME" -se "SELECT COUNT(*) FROM $table_name")
        echo "✅ $table_name seeded successfully! ($count records)"
    else
        echo "❌ Failed to seed $table_name"
        return 1
    fi
}

# Seed brands first (because automobiles depends on it)
echo "--- Seeding Brands ---"
run_seed "migrations/seeds/brands_seed.sql" "brands"

echo ""
echo "--- Seeding Automobiles ---"
run_seed "migrations/seeds/automobiles_seed.sql" "automobiles"

echo ""
echo "=== Seeding Complete ==="
