# Vehicles API Documentation

## Overview

The Vehicles API allows users to create and retrieve vehicle listings with images, technical specifications, and contact information.

## Endpoints

### 1. Create Vehicle (Add Listing)

**Endpoint**: `POST /api/vehicles`

**Content-Type**: `multipart/form-data`

**Description**: Creates a new vehicle listing with images and all details.

#### Request Parameters

| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| `title` | string | ✅ Yes | Vehicle title | Min 5, Max 255 characters |
| `category` | string | ✅ Yes | Vehicle category | - |
| `description` | string | ❌ No | Detailed description | - |
| `price` | float | ✅ Yes | Price | Must be > 0 |
| `currency` | string | ✅ Yes | Currency code | Default: "lei" |
| `negotiable` | boolean | ❌ No | Price is negotiable | Default: false |
| `person_type` | string | ✅ Yes | Seller type | `persoana_fizica` or `firma` |
| `brand` | string | ✅ Yes | Vehicle brand | - |
| `model` | string | ✅ Yes | Vehicle model | - |
| `engine_capacity` | integer | ❌ No | Engine capacity (cm³) | - |
| `power_hp` | integer | ❌ No | Power (HP) | - |
| `fuel_type` | string | ✅ Yes | Fuel type | See enum values below |
| `body_type` | string | ✅ Yes | Body type | See enum values below |
| `kilometers` | integer | ❌ No | Mileage in km | - |
| `color` | string | ❌ No | Vehicle color | - |
| `year` | integer | ✅ Yes | Manufacturing year | Min: 1970, Max: 2030 |
| `number_of_keys` | integer | ❌ No | Number of keys | - |
| `condition` | string | ✅ Yes | Vehicle condition | `utilizat` or `nou` |
| `transmission` | string | ✅ Yes | Transmission type | `manuala` or `automata` |
| `steering` | string | ✅ Yes | Steering position | `stanga` or `dreapta` |
| `registered` | boolean | ❌ No | Is vehicle registered | Default: false |
| `city` | string | ✅ Yes | City location | - |
| `contact_name` | string | ✅ Yes | Contact person name | - |
| `email` | string | ✅ Yes | Contact email | Must be valid email |
| `phone` | string | ❌ No | Contact phone | - |
| `images[]` | file[] | ❌ No | Vehicle images | Max 8 files, 10MB each |

#### Enum Values

**fuel_type**:
- `benzina` - Gasoline
- `motorina` - Diesel
- `electric` - Electric
- `hibrid` - Hybrid
- `gpl` - LPG
- `hybrid_benzina` - Hybrid Gasoline
- `hybrid_motorina` - Hybrid Diesel

**body_type**:
- `sedan` - Sedan
- `suv` - SUV
- `break` - Station Wagon
- `coupe` - Coupe
- `cabrio` - Convertible
- `hatchback` - Hatchback
- `pickup` - Pickup
- `van` - Van
- `monovolum` - Minivan

#### Image Requirements

- **Maximum files**: 8 images per vehicle
- **Maximum file size**: 10 MB per image
- **Allowed formats**: JPG, JPEG, PNG, WEBP
- **Storage**: Images are stored in `/uploads/vehicles/` directory
- **Access**: Images accessible via `/uploads/vehicles/{filename}`

#### Request Example (cURL)

```bash
curl -X POST http://localhost:8080/api/vehicles \
  -F "title=Audi A4 2018 - Excellent Condition" \
  -F "category=autoturisme" \
  -F "description=Masina in stare impecabila, service la zi" \
  -F "price=15500" \
  -F "currency=lei" \
  -F "negotiable=true" \
  -F "person_type=persoana_fizica" \
  -F "brand=Audi" \
  -F "model=A4" \
  -F "engine_capacity=2000" \
  -F "power_hp=190" \
  -F "fuel_type=motorina" \
  -F "body_type=sedan" \
  -F "kilometers=85000" \
  -F "color=Negru" \
  -F "year=2018" \
  -F "number_of_keys=2" \
  -F "condition=utilizat" \
  -F "transmission=automata" \
  -F "steering=stanga" \
  -F "registered=true" \
  -F "city=Bucuresti" \
  -F "contact_name=Ion Popescu" \
  -F "email=ion.popescu@example.com" \
  -F "phone=+40721123456" \
  -F "images=@/path/to/image1.jpg" \
  -F "images=@/path/to/image2.jpg" \
  -F "images=@/path/to/image3.jpg"
```

#### Success Response (201 Created)

```json
{
  "status": "success",
  "message": "Vehicle added successfully",
  "vehicle_id": 1,
  "data": {
    "id": 1,
    "title": "Audi A4 2018 - Excellent Condition",
    "category": "autoturisme",
    "description": "Masina in stare impecabila, service la zi",
    "price": 15500,
    "currency": "lei",
    "negotiable": true,
    "person_type": "persoana_fizica",
    "brand": "Audi",
    "model": "A4",
    "engine_capacity": 2000,
    "power_hp": 190,
    "fuel_type": "motorina",
    "body_type": "sedan",
    "kilometers": 85000,
    "color": "Negru",
    "year": 2018,
    "number_of_keys": 2,
    "condition": "utilizat",
    "transmission": "automata",
    "steering": "stanga",
    "registered": true,
    "city": "Bucuresti",
    "contact_name": "Ion Popescu",
    "email": "ion.popescu@example.com",
    "phone": "+40721123456",
    "created_at": "2025-11-22T20:00:00Z",
    "updated_at": "2025-11-22T20:00:00Z",
    "images": [
      {
        "id": 1,
        "vehicle_id": 1,
        "image_url": "/uploads/vehicles/1732307200_abc123-def456.jpg",
        "created_at": "2025-11-22T20:00:00Z"
      },
      {
        "id": 2,
        "vehicle_id": 1,
        "image_url": "/uploads/vehicles/1732307201_ghi789-jkl012.jpg",
        "created_at": "2025-11-22T20:00:01Z"
      }
    ]
  }
}
```

#### Error Responses

**400 Bad Request - Validation Error**:
```json
{
  "status": "error",
  "message": "Validation failed",
  "errors": {
    "title": "title must be at least 5 characters",
    "price": "price must be greater than 0",
    "email": "Invalid email format",
    "year": "year must be at least 1970"
  }
}
```

**400 Bad Request - Image Upload Error**:
```json
{
  "status": "error",
  "message": "Failed to upload images",
  "error": "file image.jpg exceeds maximum size of 10MB"
}
```

**500 Internal Server Error**:
```json
{
  "status": "error",
  "message": "Failed to create vehicle",
  "error": "database connection error"
}
```

---

### 2. Get Vehicle by ID

**Endpoint**: `GET /api/vehicles/:id`

**Content-Type**: `application/json`

**Description**: Retrieves a single vehicle listing with all details and images.

#### Request Example

```bash
curl http://localhost:8080/api/vehicles/1
```

#### Success Response (200 OK)

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "title": "Audi A4 2018 - Excellent Condition",
    "category": "autoturisme",
    "description": "Masina in stare impecabila, service la zi",
    "price": 15500,
    "currency": "lei",
    "negotiable": true,
    "person_type": "persoana_fizica",
    "brand": "Audi",
    "model": "A4",
    "engine_capacity": 2000,
    "power_hp": 190,
    "fuel_type": "motorina",
    "body_type": "sedan",
    "kilometers": 85000,
    "color": "Negru",
    "year": 2018,
    "number_of_keys": 2,
    "condition": "utilizat",
    "transmission": "automata",
    "steering": "stanga",
    "registered": true,
    "city": "Bucuresti",
    "contact_name": "Ion Popescu",
    "email": "ion.popescu@example.com",
    "phone": "+40721123456",
    "created_at": "2025-11-22T20:00:00Z",
    "updated_at": "2025-11-22T20:00:00Z",
    "images": [
      {
        "id": 1,
        "vehicle_id": 1,
        "image_url": "/uploads/vehicles/1732307200_abc123-def456.jpg",
        "created_at": "2025-11-22T20:00:00Z"
      }
    ]
  }
}
```

#### Error Responses

**404 Not Found**:
```json
{
  "status": "error",
  "message": "Vehicle not found"
}
```

**400 Bad Request**:
```json
{
  "status": "error",
  "message": "Invalid vehicle ID"
}
```

---

## Database Schema

### vehicles Table

```sql
CREATE TABLE vehicles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description LONGTEXT,
    price DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'lei',
    negotiable BOOLEAN DEFAULT FALSE,
    person_type ENUM('persoana_fizica', 'firma') NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    engine_capacity INT,
    power_hp INT,
    fuel_type ENUM(...) NOT NULL,
    body_type ENUM(...) NOT NULL,
    kilometers INT,
    color VARCHAR(50),
    year INT NOT NULL,
    number_of_keys INT,
    `condition` ENUM('utilizat', 'nou') NOT NULL,
    transmission ENUM('manuala', 'automata') NOT NULL,
    steering ENUM('stanga', 'dreapta') NOT NULL DEFAULT 'stanga',
    registered BOOLEAN DEFAULT FALSE,
    city VARCHAR(100) NOT NULL,
    contact_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- Indexes omitted for brevity
);
```

### vehicle_images Table

```sql
CREATE TABLE vehicle_images (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    vehicle_id BIGINT UNSIGNED NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE
);
```

---

## Testing

### Create a Test Vehicle

```bash
# Create a test vehicle without images
curl -X POST http://localhost:8080/api/vehicles \
  -F "title=BMW 320d 2019" \
  -F "category=autoturisme" \
  -F "price=18000" \
  -F "currency=lei" \
  -F "person_type=persoana_fizica" \
  -F "brand=BMW" \
  -F "model=320d" \
  -F "fuel_type=motorina" \
  -F "body_type=sedan" \
  -F "year=2019" \
  -F "condition=utilizat" \
  -F "transmission=automata" \
  -F "steering=stanga" \
  -F "city=Cluj-Napoca" \
  -F "contact_name=Maria Ionescu" \
  -F "email=maria@example.com"
```

### Retrieve the Vehicle

```bash
curl http://localhost:8080/api/vehicles/1 | jq
```

---

## Notes

- All timestamps are in UTC
- Images are stored with unique filenames (timestamp_uuid.ext)
- Old images can be cleaned up if a vehicle is deleted (CASCADE on vehicle_images table)
- The API supports CORS for `http://localhost:3000`
- Swagger documentation available at `/swagger/index.html`

---

**Created**: 2025-11-22
**Version**: 1.0
**Status**: ✅ Production Ready
