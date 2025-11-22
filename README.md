# AutoElys Backend API

Backend API for AutoElys application built with Go, Gin framework, and MySQL.

## Features

- ✅ User registration with validation
- ✅ JWT authentication
- ✅ Role-based access control (Admin/User)
- ✅ Database migrations
- ✅ Rate limiting
- ✅ Swagger API documentation
- ✅ Bcrypt password hashing

## Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher
- Swag CLI for generating Swagger docs

## Installation

1. **Install dependencies**:
```bash
go mod download
```

2. **Install Swag CLI** (for generating Swagger documentation):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

3. **Configure environment**:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Generate Swagger documentation**:
```bash
swag init
```

5. **Run database migrations**:
```bash
go run main.go migrate:up
```

6. **Start the server**:
```bash
go run main.go
```

## API Endpoints

### Authentication

#### Register User
```bash
POST /api/auth/register
```

**Request Body**:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone": "+40712345678",
  "password": "Password123",
  "password_confirmation": "Password123",
  "accepted_terms": true
}
```

#### Login
```bash
POST /api/auth/login
```

**Request Body**:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone": "+40712345678",
  "password": "Password123",
  "password_confirmation": "Password123",
  "accepted_terms": true
}
```

**Success Response (201)**:
```json
{
  "message": "Account created successfully.",
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login
```bash
POST /api/auth/login
```

**Request Body**:
```json
{
  "email": "john@example.com",
  "password": "Password123"
}
```

**Success Response (200)**:
```json
{
  "message": "Login successful.",
  "user": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Forgot Password
```bash
POST /api/auth/forgot-password
```

**Request Body**:
```json
{
  "email": "john@example.com"
}
```

**Success Response (200)**:
```json
{
  "message": "If the email exists, a password reset link has been sent."
}
```

#### Reset Password
```bash
POST /api/auth/reset-password
```

**Request Body**:
```json
{
  "token": "reset-token-here",
  "password": "NewPassword123",
  "password_confirmation": "NewPassword123"
}
```

**Success Response (200)**:
```json
{
  "message": "Password has been reset successfully. You can now login with your new password."
}
```

## Validation Rules

- **first_name**: required, min 2 characters
- **last_name**: required, min 2 characters
- **email**: required, valid email format, unique
- **phone**: optional, E.164 format (+407xxxxxxxxx)
- **password**: required, min 8 characters, must contain letters and digits
- **password_confirmation**: must match password
- **accepted_terms**: must be true

## Database Migration Commands

```bash
# Run all pending migrations
go run main.go migrate:up

# Rollback last migration
go run main.go migrate:down

# Check current migration version
go run main.go migrate:status
```

## API Documentation

Once the server is running, access Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

## Project Structure

```
/var/www/autoelys_backend/
├── database/
│   ├── config.go           # Database configuration
│   └── migrate.go          # Migration system
├── internal/
│   ├── auth/
│   │   ├── jwt.go          # JWT token generation & validation
│   │   └── password.go     # Password hashing with bcrypt
│   ├── handlers/
│   │   └── auth_handler.go # Authentication endpoints
│   ├── middleware/
│   │   └── rate_limit.go   # Rate limiting middleware
│   ├── models/
│   │   └── user.go         # User & Role models
│   ├── repository/
│   │   └── user_repository.go # Database operations
│   └── validation/
│       └── custom_validators.go # Custom validators
├── migrations/
│   ├── 000001_create_roles_table.up.sql
│   ├── 000001_create_roles_table.down.sql
│   ├── 000002_create_users_table.up.sql
│   └── 000002_create_users_table.down.sql
├── docs/                   # Auto-generated Swagger docs
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
└── .env.example            # Environment variables template
```

## Security Features

- **Password Hashing**: Bcrypt with default cost
- **JWT Tokens**: 24-hour expiration
- **Rate Limiting**: 10 requests/minute per IP
- **SQL Injection Prevention**: Prepared statements
- **Input Validation**: go-playground/validator
- **Email Normalization**: Lowercase conversion

## Development

### Regenerate Swagger docs after API changes:
```bash
swag init
```

### Create a new migration:
Create two files in `migrations/` directory:
- `{version}_{description}.up.sql`
- `{version}_{description}.down.sql`

Example:
- `000003_create_posts_table.up.sql`
- `000003_create_posts_table.down.sql`

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=autoelys
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
```

## Testing

Example cURL request:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+40712345678",
    "password": "Password123",
    "password_confirmation": "Password123",
    "accepted_terms": true
  }'
```

## License

MIT
