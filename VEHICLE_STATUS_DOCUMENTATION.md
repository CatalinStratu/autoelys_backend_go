# Vehicle Status Feature Documentation

## Overview

A status system has been added to vehicles to track their state (active, inactive, banned).

---

## Database Schema

### Status Column

```sql
ALTER TABLE vehicles
ADD COLUMN status TINYINT UNSIGNED NOT NULL DEFAULT 1
COMMENT '1=active, 2=inactive, 3=banned' AFTER user_id;
```

**Column Details:**
- **Name**: `status`
- **Type**: `TINYINT UNSIGNED`
- **Default**: `1` (active)
- **Nullable**: `NO`
- **Comment**: `1=active, 2=inactive, 3=banned`
- **Index**: `idx_status`
- **Position**: After `user_id`

---

## Status Values

| Value | Constant | Name | Description |
|-------|----------|------|-------------|
| 1 | `VehicleStatusActive` | active | Vehicle is publicly visible and searchable |
| 2 | `VehicleStatusInactive` | inactive | Vehicle is hidden from public (draft, sold, etc.) |
| 3 | `VehicleStatusBanned` | banned | Vehicle violates rules, permanently hidden |

---

## Code Implementation

### Constants (models/vehicle.go)

```go
const (
    VehicleStatusActive   uint8 = 1
    VehicleStatusInactive uint8 = 2
    VehicleStatusBanned   uint8 = 3
)
```

### Helper Function

```go
func GetStatusName(status uint8) string {
    switch status {
    case VehicleStatusActive:
        return "active"
    case VehicleStatusInactive:
        return "inactive"
    case VehicleStatusBanned:
        return "banned"
    default:
        return "unknown"
    }
}
```

### Vehicle Model Fields

```go
type Vehicle struct {
    // ...
    Status     uint8  `json:"status"`             // 1=active, 2=inactive, 3=banned
    StatusName string `json:"status_name,omitempty"` // Human-readable status
    // ...
}
```

---

## API Response

### Vehicle Response with Status

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "user_id": 3,
    "status": 1,
    "status_name": "active",
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "slug": "bmw-320d-2019",
    "title": "BMW 320d 2019",
    ...
  }
}
```

---

## Usage Examples

### Setting Status on Create

```go
vehicle := &models.Vehicle{
    UserID: userID,
    Status: models.VehicleStatusActive, // Optional, defaults to active
    Title:  "BMW 320d",
    // ... other fields
}

createdVehicle, err := vehicleRepo.Create(vehicle)
```

**Note**: If `Status` is not set (0), it automatically defaults to `VehicleStatusActive` (1).

### Checking Vehicle Status

```go
vehicle, err := vehicleRepo.GetByUUID(uuid)
if err != nil {
    return err
}

if vehicle.Status == models.VehicleStatusBanned {
    return errors.New("Vehicle is banned")
}

if vehicle.Status != models.VehicleStatusActive {
    return errors.New("Vehicle is not active")
}
```

### Updating Vehicle Status

```go
// Mark as inactive (e.g., sold)
vehicle.Status = models.VehicleStatusInactive
err := vehicleRepo.Update(uuid, vehicle)

// Ban vehicle
vehicle.Status = models.VehicleStatusBanned
err := vehicleRepo.Update(uuid, vehicle)

// Reactivate vehicle
vehicle.Status = models.VehicleStatusActive
err := vehicleRepo.Update(uuid, vehicle)
```

---

## Query Examples

### Get Only Active Vehicles

```sql
SELECT * FROM vehicles
WHERE status = 1
ORDER BY created_at DESC;
```

### Get User's Active Vehicles

```sql
SELECT * FROM vehicles
WHERE user_id = ? AND status = 1
ORDER BY created_at DESC;
```

### Count Vehicles by Status

```sql
SELECT
    status,
    CASE
        WHEN status = 1 THEN 'active'
        WHEN status = 2 THEN 'inactive'
        WHEN status = 3 THEN 'banned'
        ELSE 'unknown'
    END as status_name,
    COUNT(*) as count
FROM vehicles
GROUP BY status;
```

### Get Banned Vehicles (Admin)

```sql
SELECT * FROM vehicles
WHERE status = 3
ORDER BY updated_at DESC;
```

---

## Use Cases

### 1. **Active (status = 1)**
- Default state for new vehicles
- Publicly visible
- Searchable
- Appears in listings

**When to use:**
- Vehicle is available for sale
- All information is complete
- No violations

### 2. **Inactive (status = 2)**
- Temporarily hidden
- Not searchable
- Owner can reactivate

**When to use:**
- Vehicle is sold
- Draft listing (not ready to publish)
- User wants to pause listing
- Seasonal vehicles (store for winter)
- Under review/verification

### 3. **Banned (status = 3)**
- Permanently hidden (or until admin review)
- Not accessible to owner
- Cannot be reactivated by owner

**When to use:**
- Violates terms of service
- Fraudulent listing
- Spam/scam
- Duplicate listing
- Admin moderation required

---

## Frontend Integration

### Display Status Badge

```jsx
const StatusBadge = ({ status, statusName }) => {
  const getStatusColor = (status) => {
    switch (status) {
      case 1: return 'green';    // active
      case 2: return 'orange';   // inactive
      case 3: return 'red';      // banned
      default: return 'gray';
    }
  };

  return (
    <span
      className={`badge badge-${getStatusColor(status)}`}
    >
      {statusName || 'unknown'}
    </span>
  );
};

// Usage
<StatusBadge status={vehicle.status} statusName={vehicle.status_name} />
```

### Filter Active Vehicles

```javascript
const activeVehicles = vehicles.filter(v => v.status === 1);
const inactiveVehicles = vehicles.filter(v => v.status === 2);
const bannedVehicles = vehicles.filter(v => v.status === 3);
```

### Status Toggle (Admin)

```javascript
const toggleVehicleStatus = async (uuid, newStatus) => {
  const response = await fetch(
    `http://localhost:8080/api/admin/vehicles/${uuid}/status`,
    {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${adminToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ status: newStatus })
    }
  );

  const data = await response.json();
  return data;
};

// Ban vehicle
await toggleVehicleStatus(uuid, 3);

// Reactivate
await toggleVehicleStatus(uuid, 1);
```

---

## Migration Files

### Up Migration
**File**: `migrations/000016_add_status_to_vehicles.up.sql`

```sql
ALTER TABLE vehicles
ADD COLUMN status TINYINT UNSIGNED NOT NULL DEFAULT 1
COMMENT '1=active, 2=inactive, 3=banned' AFTER user_id,
ADD INDEX idx_status (status);

UPDATE vehicles SET status = 1 WHERE status IS NULL OR status = 0;
```

### Down Migration
**File**: `migrations/000016_add_status_to_vehicles.down.sql`

```sql
ALTER TABLE vehicles
DROP INDEX idx_status,
DROP COLUMN status;
```

---

## Benefits

✅ **Moderation**: Admins can hide inappropriate listings
✅ **User Control**: Owners can pause listings
✅ **Sold Tracking**: Mark vehicles as sold without deletion
✅ **Draft Mode**: Users can save incomplete listings
✅ **Spam Prevention**: Ban fraudulent listings
✅ **Data Retention**: Keep history even when hidden
✅ **Performance**: Index on status for fast filtering

---

## Future Enhancements

### 1. Status History Table

Track status changes over time:

```sql
CREATE TABLE vehicle_status_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    vehicle_id BIGINT UNSIGNED NOT NULL,
    old_status TINYINT UNSIGNED,
    new_status TINYINT UNSIGNED NOT NULL,
    changed_by BIGINT UNSIGNED,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_vehicle_id (vehicle_id)
);
```

### 2. Additional Statuses

```go
const (
    VehicleStatusActive      uint8 = 1
    VehicleStatusInactive    uint8 = 2
    VehicleStatusBanned      uint8 = 3
    VehicleStatusPending     uint8 = 4  // Awaiting approval
    VehicleStatusSold        uint8 = 5  // Marked as sold
    VehicleStatusExpired     uint8 = 6  // Listing expired
    VehicleStatusUnderReview uint8 = 7  // Admin reviewing
)
```

### 3. Auto-Expiry

Automatically set status to expired after X days:

```sql
UPDATE vehicles
SET status = 6
WHERE status = 1
  AND created_at < DATE_SUB(NOW(), INTERVAL 90 DAY);
```

### 4. Admin Endpoints

```
PUT /api/admin/vehicles/:uuid/status
GET /api/admin/vehicles/banned
GET /api/admin/vehicles/pending
POST /api/admin/vehicles/:uuid/approve
POST /api/admin/vehicles/:uuid/ban
```

---

## Testing

### Check Status Column

```bash
mysql -u root -p autoelys_backend -e "SHOW FULL COLUMNS FROM vehicles WHERE Field='status';"
```

### Verify Default Values

```bash
mysql -u root -p autoelys_backend -e "SELECT id, title, status FROM vehicles LIMIT 5;"
```

### Test Status Assignment

```bash
mysql -u root -p autoelys_backend -e "UPDATE vehicles SET status = 2 WHERE id = 1;"
mysql -u root -p autoelys_backend -e "SELECT id, status FROM vehicles WHERE id = 1;"
```

---

## Notes

- ✅ Status defaults to `1` (active) for all new vehicles
- ✅ Existing vehicles automatically set to active
- ✅ Status is included in all vehicle API responses
- ✅ Status name is human-readable (active/inactive/banned)
- ✅ Database has comment explaining values
- ✅ Index added for performance
- ✅ NOT NULL constraint ensures all vehicles have a status

---

**Created**: 2025-11-23
**Migration**: 000016_add_status_to_vehicles
**Status Values**: 1=active, 2=inactive, 3=banned
**Default**: 1 (active)
