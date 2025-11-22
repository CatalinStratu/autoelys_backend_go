# Database Setup Guide

## Overview

The database includes the following tables:
1. **roles** - User roles (admin, user)
2. **users** - User accounts with authentication
3. **password_reset_tokens** - Password reset functionality
4. **brands** - Automobile brands (124 brands)
5. **automobiles** - Automobile models (7000+ vehicles)

## Setup Instructions

### 1. Configure Environment

Make sure your `.env` file has the correct database credentials:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=autoelys_backend
```

### 2. Run Migrations

Run all migrations to create the database schema:

```bash
go run main.go migrate:up
```

This will create the following tables in order:
- `000001_create_roles_table.up.sql` - Creates roles table and seeds admin/user roles
- `000002_create_users_table.up.sql` - Creates users table
- `000003_create_password_reset_tokens_table.up.sql` - Creates password reset tokens table
- `000004_add_uuid_to_users.up.sql` - Adds UUID column to users
- `000005_create_brands_table.up.sql` - Creates brands table
- `000006_create_automobiles_table.up.sql` - Creates automobiles table

### 3. Seed Data

After migrations are complete, seed the brands and automobiles data:

```bash
./seed-database.sh
```

This will:
- Load **124 automobile brands** (AC, Acura, Alfa Romeo, Audi, BMW, etc.)
- Load **7000+ automobile models** with their details

## Database Schema

### Brands Table

```sql
CREATE TABLE brands (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    url_hash VARCHAR(191) NOT NULL,
    url TEXT NOT NULL,
    name VARCHAR(191) NOT NULL,
    logo TEXT,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    INDEX idx_url_hash (url_hash),
    INDEX idx_name (name)
);
```

**Sample Brands:**
- AC
- ACURA
- ALFA ROMEO
- AUDI
- BMW
- FERRARI
- FORD
- MERCEDES BENZ
- TESLA
- TOYOTA
- And 114 more...

### Automobiles Table

```sql
CREATE TABLE automobiles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    url_hash VARCHAR(191) NOT NULL,
    url LONGTEXT NOT NULL,
    brand_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(191) NOT NULL,
    description LONGTEXT,
    press_release LONGTEXT,
    photos LONGTEXT,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    INDEX idx_url_hash (url_hash),
    INDEX idx_brand_id (brand_id),
    INDEX idx_name (name),
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
);
```

## Migration Commands

```bash
# Run all pending migrations
go run main.go migrate:up

# Rollback last migration
go run main.go migrate:down

# Check migration status
go run main.go migrate:status
```

## Seed Files Location

- **Brands seed**: `migrations/seeds/brands_seed.sql`
- **Automobiles seed**: `migrations/seeds/automobiles_seed.sql`

## Verify Data

After seeding, verify the data was loaded:

```bash
# Count brands
mysql -u root -p -e "SELECT COUNT(*) FROM brands" autoelys_backend

# Count automobiles
mysql -u root -p -e "SELECT COUNT(*) FROM automobiles" autoelys_backend

# View some brands
mysql -u root -p -e "SELECT id, name FROM brands LIMIT 10" autoelys_backend
```

Expected results:
- **Brands**: 124 records
- **Automobiles**: 7000+ records

## Troubleshooting

### Foreign Key Constraint Error

If you get a foreign key error when seeding automobiles, make sure brands are seeded first:

```bash
# Seed brands first
mysql -u root -p autoelys_backend < migrations/seeds/brands_seed.sql

# Then seed automobiles
mysql -u root -p autoelys_backend < migrations/seeds/automobiles_seed.sql
```

### Reset Database

To start fresh:

```bash
# Rollback all migrations
go run main.go migrate:down

# Run all migrations again
go run main.go migrate:up

# Seed data again
./seed-database.sh
```

## Next Steps

After setting up the database:
1. Create API endpoints for brands and automobiles
2. Implement search and filtering
3. Add pagination for large result sets
4. Create relationships between users and automobiles (favorites, etc.)
