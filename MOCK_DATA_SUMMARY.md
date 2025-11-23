# Mock Vehicle Data Summary

## Overview

Mock vehicle data has been successfully added to the database for testing purposes.

---

## Mock Data Statistics

- **Total Vehicles**: 10
- **Owner User ID**: 3 (john@example.com)
- **Vehicles with Images**: 3 (BMW, Audi, Mercedes)
- **Total Images**: 4

---

## Mock Vehicles List

| ID | Brand | Model | Year | Price | Category | Condition |
|----|-------|-------|------|-------|----------|-----------|
| 1 | BMW | 320d | 2019 | 18,500 lei | Autoturisme | Utilizat |
| 2 | Audi | A4 | 2020 | 22,000 lei | Autoturisme | Utilizat |
| 3 | Mercedes-Benz | C200 | 2021 | 28,000 lei | Autoturisme | Utilizat |
| 4 | Volkswagen | Golf GTI | 2022 | 26,500 lei | Autoturisme | Utilizat |
| 5 | Toyota | RAV4 | 2023 | 32,000 lei | Autoturisme | Utilizat |
| 6 | Ford | Mustang GT | 2020 | 35,000 lei | Autoturisme | Utilizat |
| 7 | Dacia | Duster | 2021 | 14,500 lei | Autoturisme | Utilizat |
| 8 | Porsche | 911 Carrera | 2022 | 95,000 lei | Autoturisme | Utilizat |
| 9 | Renault | Clio | 2023 | 13,500 lei | Autoturisme | NOU |
| 10 | Tesla | Model 3 | 2023 | 42,000 lei | Autoturisme | Utilizat |

---

## Vehicle Details

### 1. BMW 320d 2019 xDrive - Impecabil
- **UUID**: `550e8400-e29b-41d4-a716-446655440000`
- **Price**: 18,500 lei (Negotiable)
- **Fuel**: Diesel
- **Transmission**: Automatic
- **Body Type**: Sedan
- **Kilometers**: 85,000 km
- **Power**: 190 HP
- **Engine**: 1,995 cc
- **Color**: Black
- **Images**: 2
- **Location**: București

### 2. Audi A4 2020 Quattro S-Line
- **UUID**: `660e8400-e29b-41d4-a716-446655440001`
- **Price**: 22,000 lei (Fixed)
- **Fuel**: Diesel
- **Transmission**: Automatic
- **Body Type**: Sedan
- **Kilometers**: 62,000 km
- **Power**: 204 HP
- **Engine**: 1,984 cc
- **Color**: Grey
- **Images**: 1
- **Location**: Cluj-Napoca

### 3. Mercedes-Benz C200 2021 AMG
- **UUID**: `770e8400-e29b-41d4-a716-446655440002`
- **Price**: 28,000 lei (Negotiable)
- **Fuel**: Petrol
- **Transmission**: Automatic
- **Body Type**: Sedan
- **Kilometers**: 45,000 km
- **Power**: 204 HP
- **Engine**: 1,991 cc
- **Color**: White
- **Images**: 1
- **Location**: Timișoara
- **Seller Type**: Company

### 4. VW Golf 8 2022 GTI
- **UUID**: `880e8400-e29b-41d4-a716-446655440003`
- **Price**: 26,500 lei (Fixed)
- **Fuel**: Petrol
- **Transmission**: Automatic
- **Body Type**: Hatchback
- **Kilometers**: 28,000 km
- **Power**: 245 HP
- **Engine**: 1,984 cc
- **Color**: Red
- **Location**: Brașov

### 5. Toyota RAV4 2023 Hybrid AWD
- **UUID**: `990e8400-e29b-41d4-a716-446655440004`
- **Price**: 32,000 lei (Negotiable)
- **Fuel**: Hybrid
- **Transmission**: Automatic
- **Body Type**: SUV
- **Kilometers**: 15,000 km
- **Power**: 222 HP
- **Engine**: 2,487 cc
- **Color**: Silver
- **Location**: Constanța

### 6. Ford Mustang 2020 GT V8
- **UUID**: `aa0e8400-e29b-41d4-a716-446655440005`
- **Price**: 35,000 lei (Fixed)
- **Fuel**: Petrol
- **Transmission**: Automatic
- **Body Type**: Coupe
- **Kilometers**: 42,000 km
- **Power**: 450 HP
- **Engine**: 5,038 cc
- **Color**: Orange
- **Location**: București

### 7. Dacia Duster 2021 4x4 Prestige
- **UUID**: `bb0e8400-e29b-41d4-a716-446655440006`
- **Price**: 14,500 lei (Negotiable)
- **Fuel**: Diesel
- **Transmission**: Manual
- **Body Type**: SUV
- **Kilometers**: 68,000 km
- **Power**: 115 HP
- **Engine**: 1,461 cc
- **Color**: Orange
- **Location**: Iași

### 8. Porsche 911 Carrera 2022
- **UUID**: `cc0e8400-e29b-41d4-a716-446655440007`
- **Price**: 95,000 lei (Fixed)
- **Fuel**: Petrol
- **Transmission**: Automatic
- **Body Type**: Coupe
- **Kilometers**: 18,000 km
- **Power**: 385 HP
- **Engine**: 2,981 cc
- **Color**: Blue
- **Location**: București
- **Seller Type**: Company

### 9. Renault Clio 2023 NOU - Intens
- **UUID**: `dd0e8400-e29b-41d4-a716-446655440008`
- **Price**: 13,500 lei (Fixed)
- **Fuel**: Petrol
- **Transmission**: Manual
- **Body Type**: Hatchback
- **Kilometers**: 0 km (NEW)
- **Power**: 90 HP
- **Engine**: 999 cc
- **Color**: White
- **Condition**: NEW
- **Location**: Pitești
- **Seller Type**: Company

### 10. Tesla Model 3 2023 Long Range
- **UUID**: `ee0e8400-e29b-41d4-a716-446655440009`
- **Price**: 42,000 lei (Negotiable)
- **Fuel**: Electric
- **Transmission**: Automatic
- **Body Type**: Sedan
- **Kilometers**: 32,000 km
- **Power**: 283 HP
- **Engine**: 0 cc (Electric)
- **Color**: Black
- **Location**: Cluj-Napoca

---

## Testing API with Mock Data

### Get All User Vehicles
```bash
# First, login to get the token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "your_password_here"
  }'

# Use the token to get all vehicles
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/user/vehicles
```

### Get Specific Vehicle by UUID
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/user/vehicles/550e8400-e29b-41d4-a716-446655440000
```

### Get Vehicle by ID (Public)
```bash
curl http://localhost:8080/api/vehicles/1
```

---

## Price Range Summary

| Price Range | Count | Vehicles |
|-------------|-------|----------|
| Under 15,000 lei | 2 | Renault Clio, Dacia Duster |
| 15,000 - 25,000 lei | 2 | BMW 320d, Audi A4 |
| 25,000 - 35,000 lei | 3 | Mercedes C200, VW Golf GTI, Toyota RAV4 |
| 35,000 - 50,000 lei | 2 | Ford Mustang, Tesla Model 3 |
| Above 50,000 lei | 1 | Porsche 911 |

---

## Fuel Type Distribution

- **Petrol**: 4 vehicles (Porsche, Ford, Mercedes, VW Golf, Renault)
- **Diesel**: 3 vehicles (BMW, Audi, Dacia)
- **Electric**: 1 vehicle (Tesla)
- **Hybrid**: 1 vehicle (Toyota)

---

## Transmission Distribution

- **Automatic**: 8 vehicles
- **Manual**: 2 vehicles (Dacia Duster, Renault Clio)

---

## Body Type Distribution

- **Sedan**: 4 vehicles (BMW, Audi, Mercedes, Tesla)
- **SUV**: 2 vehicles (Toyota RAV4, Dacia Duster)
- **Coupe**: 2 vehicles (Porsche 911, Ford Mustang)
- **Hatchback**: 2 vehicles (VW Golf, Renault Clio)

---

## Re-seeding Data

To re-seed the mock data (remove old and add fresh):

```bash
# Remove existing mock vehicles
mysql -u root -p autoelys_backend -e "DELETE FROM vehicles WHERE user_id = 3;"

# Re-run the seed file
mysql -u root -p autoelys_backend < migrations/seeds/mock_vehicles.sql
```

---

**Created**: 2025-11-23
**User ID**: 3 (john@example.com)
**Total Mock Vehicles**: 10
**Seed File**: `/var/www/autoelys_backend/migrations/seeds/mock_vehicles.sql`
