# Vehicle UUID and Slug Implementation

## ✅ Complete: UUID and Slug Added to Vehicles

Successfully added UUID and slug fields to the vehicles table and integrated slug generation from vehicle titles.

---

## Database Changes

### Updated Table Structure

The `vehicles` table now includes:

```sql
CREATE TABLE vehicles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid CHAR(36) NOT NULL,                    -- NEW: Unique identifier
    slug VARCHAR(255) NOT NULL,                -- NEW: URL-friendly slug
    title VARCHAR(255) NOT NULL,
    ...
    UNIQUE INDEX idx_uuid (uuid),              -- NEW: Unique constraint
    UNIQUE INDEX idx_slug (slug)               -- NEW: Unique constraint
);
```

**Fields Added:**
- ✅ **uuid** (CHAR(36), NOT NULL, UNIQUE) - Universally unique identifier
- ✅ **slug** (VARCHAR(255), NOT NULL, UNIQUE) - SEO-friendly URL slug

---

## Slug Generation

### Algorithm

The slug is automatically generated from the vehicle title using the following process:

1. **Convert to lowercase**
2. **Remove Romanian diacritics** (ă→a, â→a, î→i, ș→s, ț→t)
3. **Replace spaces and special characters** with hyphens
4. **Remove consecutive hyphens**
5. **Trim leading/trailing hyphens**
6. **Limit to 200 characters**

### Examples

| Title | Generated Slug |
|-------|---------------|
| "BMW 320d 2019 - Impecabil" | "bmw-320d-2019-impecabil" |
| "Audi A4 în stare PERFECTĂ!" | "audi-a4-in-stare-perfecta" |
| "Mercedes-Benz E-Class 2020" | "mercedes-benz-e-class-2020" |
| "Mașină utilizată - Good Deal" | "masina-utilizata-good-deal" |

### Romanian Character Support

Special handling for Romanian diacritics:

| Character | Converted To |
|-----------|--------------|
| ă, Ă | a |
| â, Â | a |
| î, Î | i |
| ș, Ș | s |
| ț, Ț | t |

---

## Code Implementation

### 1. Slug Utility (`internal/utils/slug.go`)

```go
func GenerateSlug(text string) string {
    // Convert to lowercase
    slug := strings.ToLower(text)

    // Remove diacritics (ă, â, î, ș, ț, etc.)
    slug = removeDiacritics(slug)

    // Replace special characters with hyphens
    reg := regexp.MustCompile("[^a-z0-9]+")
    slug = reg.ReplaceAllString(slug, "-")

    // Clean up
    slug = strings.Trim(slug, "-")

    // Limit length to 200 characters
    if len(slug) > 200 {
        slug = slug[:200]
        slug = strings.TrimRight(slug, "-")
    }

    return slug
}
```

### 2. Updated Vehicle Model

```go
type Vehicle struct {
    ID             uint64    `json:"id"`
    UUID           string    `json:"uuid"`        // NEW
    Slug           string    `json:"slug"`        // NEW
    Title          string    `json:"title"`
    Category       string    `json:"category"`
    ...
}
```

### 3. Updated Repository

**Create Method:**
```go
func (r *VehicleRepository) Create(vehicle *models.Vehicle) (*models.Vehicle, error) {
    query := `INSERT INTO vehicles (
        uuid, slug, title, category, ...
    ) VALUES (?, ?, ?, ?, ...)`

    result, err := r.db.Exec(query,
        vehicle.UUID,
        vehicle.Slug,
        vehicle.Title,
        ...
    )
    ...
}
```

**GetByID Method:**
```go
query := `SELECT
    v.id, v.uuid, v.slug, v.title, ...
FROM vehicles v
...`
```

### 4. Updated Handler

**CreateVehicle Method:**
```go
func (h *VehicleHandler) CreateVehicle(c *gin.Context) {
    ...

    // Generate UUID and slug
    vehicleUUID := uuid.New().String()
    slug := utils.GenerateSlug(req.Title)

    // Create vehicle model
    vehicle := &models.Vehicle{
        UUID:  vehicleUUID,
        Slug:  slug,
        Title: req.Title,
        ...
    }

    // Save to database
    createdVehicle, err := h.vehicleRepo.Create(vehicle)
    ...
}
```

---

## API Response Format

### Create Vehicle Response

```json
{
  "status": "success",
  "message": "Vehicle added successfully",
  "vehicle_id": 1,
  "data": {
    "id": 1,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "slug": "bmw-320d-2019-impecabil",
    "title": "BMW 320d 2019 - Impecabil",
    "category": "autoturisme",
    "brand": "BMW",
    "model": "320d",
    "price": 18000,
    ...
  }
}
```

### Get Vehicle Response

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "slug": "bmw-320d-2019-impecabil",
    "title": "BMW 320d 2019 - Impecabil",
    ...
  }
}
```

---

## Use Cases

### 1. SEO-Friendly URLs

Instead of:
```
https://autoelys.com/vehicles/123
```

Use:
```
https://autoelys.com/vehicles/bmw-320d-2019-impecabil
```

### 2. Public Sharing Links

UUID-based sharing (harder to guess, more secure):
```
https://autoelys.com/v/550e8400-e29b-41d4-a716-446655440000
```

### 3. Database Lookups

**By ID** (fastest):
```sql
SELECT * FROM vehicles WHERE id = 1
```

**By UUID** (public-facing):
```sql
SELECT * FROM vehicles WHERE uuid = '550e8400-e29b-41d4-a716-446655440000'
```

**By Slug** (SEO-friendly):
```sql
SELECT * FROM vehicles WHERE slug = 'bmw-320d-2019-impecabil'
```

---

## Benefits

### ✅ UUID Benefits

1. **Global Uniqueness**: Can merge databases without ID conflicts
2. **Security**: Harder to guess than sequential IDs
3. **Public Sharing**: Safe to expose in URLs without revealing business metrics
4. **Distributed Systems**: Can generate offline without coordination

### ✅ Slug Benefits

1. **SEO Optimization**: Search engines prefer descriptive URLs
2. **User-Friendly**: Readable and memorable URLs
3. **Click-Through Rate**: Users more likely to click descriptive links
4. **Social Sharing**: Better appearance when shared on social media

---

## Future Enhancements

### 1. Slug Uniqueness Handling

Currently, slugs must be unique. To handle duplicate titles:

```go
// Check if slug exists
existingCount := checkSlugExists(slug)

// Append counter if duplicate
if existingCount > 0 {
    slug = fmt.Sprintf("%s-%d", slug, existingCount)
}
```

### 2. Get Vehicle by Slug Endpoint

Add new endpoint:

```go
GET /api/vehicles/slug/:slug
```

**Example**:
```bash
curl http://localhost:8080/api/vehicles/slug/bmw-320d-2019-impecabil
```

### 3. Slug History/Redirects

Keep old slugs when title changes:

```sql
CREATE TABLE vehicle_slug_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    vehicle_id BIGINT UNSIGNED NOT NULL,
    old_slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE
);
```

### 4. Custom Slugs

Allow users to customize slugs:

```json
{
  "title": "BMW 320d 2019",
  "custom_slug": "my-awesome-bmw"
}
```

---

## Testing

### Create Test Vehicle

```bash
curl -X POST http://localhost:8080/api/vehicles \
  -F "title=Audi A4 în stare perfectă - 2020" \
  -F "category=autoturisme" \
  -F "price=20000" \
  -F "currency=lei" \
  -F "person_type=persoana_fizica" \
  -F "brand=Audi" \
  -F "model=A4" \
  -F "fuel_type=motorina" \
  -F "body_type=sedan" \
  -F "year=2020" \
  -F "condition=utilizat" \
  -F "transmission=automata" \
  -F "steering=stanga" \
  -F "city=Bucuresti" \
  -F "contact_name=Test User" \
  -F "email=test@example.com"
```

**Expected Response:**
```json
{
  "status": "success",
  "vehicle_id": 1,
  "data": {
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "slug": "audi-a4-in-stare-perfecta-2020",
    "title": "Audi A4 în stare perfectă - 2020"
  }
}
```

### Verify in Database

```bash
mysql -u root -p autoelys_backend -e "
SELECT id, uuid, slug, title
FROM vehicles
LIMIT 5
"
```

---

## Files Modified

1. ✅ **internal/models/vehicle.go**
   - Added UUID and Slug fields

2. ✅ **internal/utils/slug.go** (NEW)
   - Slug generation utility
   - Romanian diacritic removal

3. ✅ **internal/repository/vehicle_repository.go**
   - Updated Create() to include uuid and slug
   - Updated GetByID() to return uuid and slug

4. ✅ **internal/handlers/vehicle_handler.go**
   - Added UUID generation
   - Added slug generation from title

---

## Summary

**Status**: ✅ **Complete**

UUID and slug functionality has been successfully implemented for all vehicle listings.

**Key Features**:
- ✅ Automatic UUID generation (v4)
- ✅ Automatic slug generation from title
- ✅ Romanian character support
- ✅ Unique constraints on both fields
- ✅ SEO-friendly URLs ready
- ✅ Public sharing links ready

**Database State**:
- vehicles table has uuid and slug columns
- Both columns have unique indexes
- All new vehicles will have UUID and slug auto-generated

**API Changes**:
- Response now includes uuid and slug fields
- No changes required to request format
- Fully backward compatible

---

**Created**: 2025-11-22
**Version**: 1.0
**Dependencies**: google/uuid, golang.org/x/text
