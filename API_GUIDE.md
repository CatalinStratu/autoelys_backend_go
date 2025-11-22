# AutoElys Backend API Guide

## Overview

This API provides authentication and user management functionality following Go best practices.

## Base URL
```
http://localhost:8080
```

## Authentication Endpoints

### 1. Register User
Create a new user account.

**Endpoint:** `POST /api/auth/register`

**Request:**
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

**Response (201):**
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

---

### 2. Login
Authenticate existing user.

**Endpoint:** `POST /api/auth/login`

**Request:**
```json
{
  "email": "john@example.com",
  "password": "Password123"
}
```

**Response (200):**
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

**Error Responses:**
- `401 Unauthorized`: Invalid credentials or inactive account
- `422 Unprocessable Entity`: Validation errors

---

### 3. Forgot Password
Request a password reset token.

**Endpoint:** `POST /api/auth/forgot-password`

**Request:**
```json
{
  "email": "john@example.com"
}
```

**Response (200):**
```json
{
  "message": "If the email exists, a password reset link has been sent."
}
```

**Notes:**
- Response is always 200 to prevent email enumeration
- Reset token expires in 1 hour
- Previous tokens are invalidated when new token is generated
- Email with reset link will be sent (currently logged to console)

---

### 4. Reset Password
Reset password using token from email.

**Endpoint:** `POST /api/auth/reset-password`

**Request:**
```json
{
  "token": "a1b2c3d4e5f6...",
  "password": "NewPassword123",
  "password_confirmation": "NewPassword123"
}
```

**Response (200):**
```json
{
  "message": "Password has been reset successfully. You can now login with your new password."
}
```

**Error Responses:**
- `400 Bad Request`: Invalid, expired, or already used token
- `422 Unprocessable Entity`: Validation errors (password requirements not met)

---

## Security Best Practices Implemented

### 1. Password Security
- **Bcrypt hashing** with default cost (currently 10)
- Minimum 8 characters
- Must contain letters and digits
- Passwords never returned in responses

### 2. JWT Tokens
- 24-hour expiration
- Contains user_id, email, and role_id
- Signed with secure secret key (base64 encoded, 64 bytes)
- HMAC-SHA256 algorithm

### 3. Rate Limiting
- 10 requests per minute per IP
- Burst allowance of 5 requests
- Automatic cleanup of old visitor records

### 4. Password Reset
- Cryptographically secure tokens (32 bytes, hex encoded)
- 1-hour expiration
- Single-use tokens (marked as used after reset)
- No email enumeration (same response for existing/non-existing emails)
- Old tokens invalidated when new one is requested

### 5. Input Validation
- go-playground/validator for struct validation
- Custom validators for phone (E.164) and password strength
- Detailed, user-friendly error messages
- Email normalization to lowercase

### 6. Database Security
- Prepared statements (SQL injection prevention)
- Foreign key constraints
- Proper indexes for performance
- Cascade deletes for related data

## Validation Rules

### Register
- `first_name`: required, min 2 characters
- `last_name`: required, min 2 characters
- `email`: required, valid email, unique
- `phone`: optional, E.164 format (+country code + number)
- `password`: required, min 8 characters, letters + digits
- `password_confirmation`: must match password
- `accepted_terms`: must be true

### Login
- `email`: required, valid email
- `password`: required

### Forgot Password
- `email`: required, valid email

### Reset Password
- `token`: required
- `password`: required, min 8 characters, letters + digits
- `password_confirmation`: must match password

## Error Responses

### Validation Error (422)
```json
{
  "errors": {
    "email": ["The email is already taken."],
    "password": ["Must be at least 8 characters long and contain both letters and digits."]
  }
}
```

### Unauthorized (401)
```json
{
  "error": "Invalid email or password"
}
```

### Too Many Requests (429)
```json
{
  "error": "Too many requests. Please try again later."
}
```

### Internal Server Error (500)
```json
{
  "error": "Internal server error"
}
```

## Testing Examples

### cURL Examples

**Register:**
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

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "Password123"
  }'
```

**Forgot Password:**
```bash
curl -X POST http://localhost:8080/api/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

**Reset Password:**
```bash
curl -X POST http://localhost:8080/api/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-reset-token-here",
    "password": "NewPassword123",
    "password_confirmation": "NewPassword123"
  }'
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT UNSIGNED NOT NULL DEFAULT 2,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    password_hash VARCHAR(255) NOT NULL,
    active TINYINT(1) DEFAULT 1,
    accepted_terms_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Password Reset Tokens Table
```sql
CREATE TABLE password_reset_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Roles Table
```sql
CREATE TABLE roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Default roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator with full access'),
    ('user', 'Regular user with standard access');
```

## Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=autoelys_backend

# JWT
JWT_SECRET=4HsOPRyKXp09rMogbEXP3r3a+TGsvzwj3cPqh7QlB6Jxeoytl8QIwKEmCgIeAXFUL3v7SF7xXLckllFQ0Q7yoA==

# Server
PORT=8080

# Email (optional - for production)
EMAIL_FROM=noreply@autoelys.com
APP_URL=http://localhost:8080
```

## Future Enhancements

- [ ] Email verification on registration
- [ ] Implement actual email sending (SMTP/SendGrid/AWS SES)
- [ ] Refresh tokens
- [ ] OAuth2 social login
- [ ] Two-factor authentication (2FA)
- [ ] Account lockout after failed login attempts
- [ ] Password history (prevent reusing recent passwords)
- [ ] Session management
- [ ] Audit logging
