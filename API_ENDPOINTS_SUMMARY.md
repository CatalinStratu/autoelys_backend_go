# Vehicle API Endpoints - Complete Summary

## Authentication & Authorization

All endpoints under `/api/user/vehicles` require **Bearer Token Authentication**.

### How to Get a Token

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "your_password"
  }'

# Response includes token
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...}
}
```

---

## Vehicle Endpoints Overview

| Method | Endpoint | Auth | Authorization | Description |
|--------|----------|------|---------------|-------------|
| **POST** | `/api/user/vehicles` | ✅ Required | **Authenticated users** | Create new vehicle |
| **GET** | `/api/user/vehicles` | ✅ Required | **Authenticated users** | Get all user's vehicles |
| **GET** | `/api/user/vehicles/:uuid` | ✅ Required | **Owner or Admin** | Get vehicle by UUID |
| **PUT** | `/api/user/vehicles/:uuid` | ✅ Required | **Owner or Admin** | Update vehicle by UUID |
| **GET** | `/api/vehicles/:id` | ❌ Public | **Anyone** | Get vehicle by ID |

---

## 1. Create Vehicle (POST /api/user/vehicles)

### 🔒 Authentication Required

**Only authenticated users can create vehicles.**

### Request

```bash
curl -X POST http://localhost:8080/api/user/vehicles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "title=BMW 320d 2019" \
  -F "category=autoturisme" \
  -F "description=Excellent condition" \
  -F "price=18500" \
  -F "currency=lei" \
  -F "negotiable=true" \
  -F "person_type=persoana_fizica" \
  -F "brand=BMW" \
  -F "model=320d" \
  -F "engine_capacity=1995" \
  -F "power_hp=190" \
  -F "fuel_type=motorina" \
  -F "body_type=sedan" \
  -F "kilometers=85000" \
  -F "color=Black" \
  -F "year=2019" \
  -F "number_of_keys=2" \
  -F "condition=utilizat" \
  -F "transmission=automata" \
  -F "steering=stanga" \
  -F "registered=true" \
  -F "city=București" \
  -F "contact_name=John Doe" \
  -F "email=john@example.com" \
  -F "phone=0721123456" \
  -F "images=@/path/to/image1.jpg" \
  -F "images=@/path/to/image2.jpg"
```

### Required Fields

| Field | Type | Validation | Options |
|-------|------|------------|---------|
| `title` | string | 5-255 chars | - |
| `category` | string | required | - |
| `price` | number | > 0 | - |
| `currency` | string | required | lei, euro, usd |
| `person_type` | string | required | persoana_fizica, firma |
| `brand` | string | required | - |
| `model` | string | required | - |
| `fuel_type` | string | required | benzina, motorina, electric, hibrid, gpl |
| `body_type` | string | required | sedan, suv, break, coupe, cabrio, hatchback, pickup, van, monovolum |
| `year` | integer | 1970-2030 | - |
| `condition` | string | required | utilizat, nou |
| `transmission` | string | required | manuala, automata |
| `steering` | string | required | stanga, dreapta |
| `city` | string | required | - |
| `contact_name` | string | required | - |
| `email` | string | valid email | - |

### Optional Fields

- `description`, `negotiable`, `engine_capacity`, `power_hp`, `kilometers`, `color`, `number_of_keys`, `registered`, `phone`, `images`

### Response (201 Created)

```json
{
  "status": "success",
  "message": "Vehicle added successfully",
  "vehicle_id": 11,
  "data": {
    "id": 11,
    "user_id": 3,
    "uuid": "new-uuid-here",
    "slug": "bmw-320d-2019",
    "title": "BMW 320d 2019",
    "category": "autoturisme",
    "price": 18500.00,
    "currency": "lei",
    "negotiable": true,
    "brand": "BMW",
    "model": "320d",
    "year": 2019,
    "images": [
      {
        "id": 1,
        "vehicle_id": 11,
        "image_url": "/uploads/vehicles/...",
        "created_at": "2025-11-23T10:00:00Z"
      }
    ],
    "created_at": "2025-11-23T10:00:00Z",
    "updated_at": "2025-11-23T10:00:00Z"
  }
}
```

### Features

✅ **Auto-assignment**: Vehicle automatically assigned to authenticated user
✅ **UUID Generation**: Unique UUID generated automatically
✅ **Slug Generation**: SEO-friendly slug from title (Romanian chars supported)
✅ **Image Upload**: Up to 8 images (JPEG/PNG)
✅ **Validation**: Comprehensive field validation

---

## 2. Get All User Vehicles (GET /api/user/vehicles)

### 🔒 Authentication Required

**Returns all vehicles created by the authenticated user.**

### Request

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/user/vehicles
```

### Response (200 OK)

```json
{
  "status": "success",
  "count": 10,
  "data": [
    {
      "id": 1,
      "user_id": 3,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "slug": "bmw-320d-2019-xdrive-impecabil",
      "title": "BMW 320d 2019 xDrive - Impecabil",
      "category": "autoturisme",
      "price": 18500.00,
      "currency": "lei",
      "brand": "BMW",
      "model": "320d",
      "year": 2019,
      "images": [...]
    },
    ...
  ]
}
```

### Features

✅ **User-specific**: Only shows logged-in user's vehicles
✅ **Ordered**: Newest first (by created_at DESC)
✅ **Complete data**: Includes all fields and images
✅ **Count**: Total number of vehicles

---

## 3. Get Vehicle by UUID (GET /api/user/vehicles/:uuid)

### 🔒 Authentication Required
### 🔐 Authorization: Owner or Admin Only

**Only the vehicle owner or admin can access vehicle details.**

### Request

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/user/vehicles/550e8400-e29b-41d4-a716-446655440000
```

### Response (200 OK)

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "user_id": 3,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "slug": "bmw-320d-2019-xdrive-impecabil",
    "title": "BMW 320d 2019 xDrive - Impecabil",
    "description": "BMW Seria 3 320d in stare excelenta...",
    "price": 18500.00,
    "currency": "lei",
    "negotiable": true,
    "person_type": "Persoană Fizică",
    "brand": "BMW",
    "model": "320d",
    "engine_capacity": 1995,
    "power_hp": 190,
    "fuel_type": "Motorină",
    "body_type": "Sedan",
    "kilometers": 85000,
    "color": "Negru",
    "year": 2019,
    "number_of_keys": 2,
    "condition": "Utilizat",
    "transmission": "Automată",
    "steering": "Stânga",
    "registered": true,
    "city": "București",
    "contact_name": "Ion Popescu",
    "email": "john@example.com",
    "phone": "0721123456",
    "images": [...],
    "created_at": "2025-11-23T10:00:00Z",
    "updated_at": "2025-11-23T10:00:00Z"
  }
}
```

### Error Response (403 Forbidden)

```json
{
  "status": "error",
  "message": "You don't have permission to view this vehicle"
}
```

### Authorization Rules

| User Type | Can Access? |
|-----------|-------------|
| Vehicle Owner | ✅ Yes |
| Admin (role_id = 1) | ✅ Yes |
| Other Users | ❌ No (403 Forbidden) |
| Not Authenticated | ❌ No (401 Unauthorized) |

---

## 4. Update Vehicle (PUT /api/user/vehicles/:uuid)

### 🔒 Authentication Required
### 🔐 Authorization: Owner or Admin Only

**Only the vehicle owner or admin can update vehicles.**

### Request

```bash
curl -X PUT http://localhost:8080/api/user/vehicles/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "price=19000" \
  -F "kilometers=90000" \
  -F "description=Updated description"
```

### Features

✅ **Partial Updates**: Only send fields you want to change
✅ **Slug Regeneration**: Slug auto-updates if title changes
✅ **Validation**: All fields validated
✅ **Owner/Admin Only**: Authorization check

### Response (200 OK)

```json
{
  "status": "success",
  "message": "Vehicle updated successfully",
  "data": {
    "id": 1,
    "user_id": 3,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "price": 19000.00,
    "kilometers": 90000,
    "description": "Updated description",
    ...
  }
}
```

### Error Response (403 Forbidden)

```json
{
  "status": "error",
  "message": "You don't have permission to update this vehicle"
}
```

---

## 5. Get Vehicle by ID (GET /api/vehicles/:id)

### 🔓 Public Endpoint (No Authentication)

**Anyone can access vehicle details by numeric ID.**

### Request

```bash
curl http://localhost:8080/api/vehicles/1
```

### Response (200 OK)

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "user_id": 3,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "slug": "bmw-320d-2019-xdrive-impecabil",
    "title": "BMW 320d 2019 xDrive - Impecabil",
    ...
  }
}
```

---

## Authorization Matrix

| Endpoint | Public | Authenticated | Owner | Admin |
|----------|--------|---------------|-------|-------|
| POST /api/user/vehicles | ❌ | ✅ | - | - |
| GET /api/user/vehicles | ❌ | ✅ | - | - |
| GET /api/user/vehicles/:uuid | ❌ | ❌ | ✅ | ✅ |
| PUT /api/user/vehicles/:uuid | ❌ | ❌ | ✅ | ✅ |
| GET /api/vehicles/:id | ✅ | ✅ | ✅ | ✅ |

---

## Error Codes

| Status Code | Meaning | When? |
|-------------|---------|-------|
| 200 | Success | Request succeeded |
| 201 | Created | Vehicle created successfully |
| 400 | Bad Request | Validation failed |
| 401 | Unauthorized | No/invalid token |
| 403 | Forbidden | Not owner/admin |
| 404 | Not Found | Vehicle doesn't exist |
| 500 | Server Error | Internal error |

---

## Testing with Mock Data

The database contains 10 mock vehicles (all owned by user ID 3).

### Sample UUIDs for Testing

```bash
# BMW 320d
550e8400-e29b-41d4-a716-446655440000

# Audi A4
660e8400-e29b-41d4-a716-446655440001

# Mercedes C200
770e8400-e29b-41d4-a716-446655440002

# Tesla Model 3
ee0e8400-e29b-41d4-a716-446655440009

# Porsche 911
cc0e8400-e29b-41d4-a716-446655440007
```

---

## Swagger Documentation

Interactive API documentation available at:

```
http://localhost:8080/swagger/index.html
```

All endpoints are documented with request/response examples, validation rules, and authorization requirements.

---

**Last Updated**: 2025-11-23
**API Version**: 1.0
**Authentication**: JWT Bearer Token
