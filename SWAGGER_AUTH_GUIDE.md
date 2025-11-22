# Swagger Authentication Guide

## How to Use Protected Endpoints in Swagger UI

### Step 1: Get a Token
First, use the `/api/auth/login` or `/api/auth/register` endpoint to get a JWT token.

Example response:
```json
{
  "message": "Login successful.",
  "user": { ... },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Step 2: Authorize in Swagger
1. Click the **"Authorize"** button (🔒 lock icon) at the top of the Swagger UI page
2. In the "Value" field, enter: **`Bearer `** followed by your token

   **IMPORTANT:** Include the word "Bearer" with a space after it!

   Example:
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJyb2xlX2lkIjoyLCJleHAiOjE3NjM5MjAxODEsIm5iZiI6MTc2MzgzMzc4MSwiaWF0IjoxNzYzODMzNzgxfQ.vnYntK3eCeqazSt9qjr0E7PwE3pe-9Rca_Q7Ne9O85c
   ```

3. Click **"Authorize"**
4. Click **"Close"**

### Step 3: Test Protected Endpoints
Now you can test protected endpoints like `/api/auth/me`

## Common Mistakes

❌ **Wrong:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

✅ **Correct:**
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Testing with cURL

If you prefer using cURL directly:

```bash
# 1. Login first
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"Password123"}'

# 2. Copy the token from the response

# 3. Use the token (with Bearer prefix)
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Token Expiration

- Tokens expire after **24 hours**
- If you get "Invalid or expired token" error, login again to get a new token
