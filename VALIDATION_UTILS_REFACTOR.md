# Validation Utilities Refactoring

## Issue

The codebase had duplicate `formatValidationErrors` function declarations in both `auth_handler.go` and `vehicle_handler.go`, causing a compilation error:

```
internal/handlers/vehicle_handler.go:295:6: formatValidationErrors redeclared in this block
    internal/handlers/auth_handler.go:523:6: other declaration of formatValidationErrors
```

## Solution

Created a centralized validation utility module to eliminate code duplication and provide reusable validation error formatting functions.

---

## Changes Made

### 1. Created New Utility File

**File**: `internal/utils/validation.go`

This file contains three exported functions:

#### `FormatValidationErrorsArray()`
- **Purpose**: Formats validation errors with support for multiple errors per field
- **Returns**: `map[string][]string`
- **Used by**: `auth_handler.go`
- **Features**:
  - Converts field names to snake_case
  - Supports auth-specific validation rules (strong_password, phone_e164, etc.)
  - Handles AcceptedTerms field specially
  - Returns array of error messages per field

#### `FormatValidationErrorsSimple()`
- **Purpose**: Formats validation errors with single error per field
- **Returns**: `map[string]string`
- **Used by**: `vehicle_handler.go`
- **Features**:
  - Keeps original field names
  - Simple one-error-per-field format
  - Handles vehicle-specific validations (oneof, gt, min, max)

#### `ToSnakeCase()`
- **Purpose**: Converts CamelCase to snake_case
- **Used by**: `FormatValidationErrorsArray()`
- **Example**: `ContactName` → `contact_name`

---

## Updated Files

### `internal/handlers/auth_handler.go`

**Changes**:
1. ✅ Added import: `"autoelys_backend/internal/utils"`
2. ✅ Removed import: `"strings"` (no longer needed)
3. ✅ Removed duplicate function: `formatValidationErrors()`
4. ✅ Removed duplicate function: `toSnakeCase()`
5. ✅ Updated 5 function calls:
   - `formatValidationErrors()` → `utils.FormatValidationErrorsArray()`

**Affected Functions**:
- `Register()` - line 144
- `Login()` - line 241
- `RequestPasswordReset()` - line 305
- `ResetPassword()` - line 362
- `UpdateProfile()` - line 473

### `internal/handlers/vehicle_handler.go`

**Changes**:
1. ✅ Import already existed: `"autoelys_backend/internal/utils"`
2. ✅ Removed duplicate function: `formatValidationErrors()`
3. ✅ Updated 1 function call:
   - `formatValidationErrors()` → `utils.FormatValidationErrorsSimple()`

**Affected Functions**:
- `CreateVehicle()` - line 112

---

## Function Signatures Comparison

### Before (Duplicated)

**auth_handler.go**:
```go
func formatValidationErrors(errs validator.ValidationErrors) map[string][]string {
    // Returns array of errors per field
    // Converts to snake_case
}

func toSnakeCase(s string) string {
    // Converts CamelCase to snake_case
}
```

**vehicle_handler.go**:
```go
func formatValidationErrors(errors validator.ValidationErrors) map[string]string {
    // Returns single error per field
    // Keeps original field names
}
```

### After (Centralized)

**utils/validation.go**:
```go
func FormatValidationErrorsArray(errs validator.ValidationErrors) map[string][]string {
    // For auth handler - array format
}

func FormatValidationErrorsSimple(errors validator.ValidationErrors) map[string]string {
    // For vehicle handler - simple format
}

func ToSnakeCase(s string) string {
    // Shared utility
}
```

---

## Benefits

### ✅ Code Reusability
- Single source of truth for validation error formatting
- No code duplication across handlers

### ✅ Maintainability
- Easy to update error messages in one place
- Consistent error formatting across the API

### ✅ Flexibility
- Two different formats available based on needs:
  - Array format for complex validation scenarios
  - Simple format for straightforward validation

### ✅ Scalability
- Future handlers can easily reuse these utilities
- Can add more validation utilities to the same file

---

## Validation Error Format Examples

### Array Format (Auth Handler)

**Request**:
```json
{
  "email": "invalid-email",
  "password": "weak"
}
```

**Response**:
```json
{
  "status": "error",
  "message": "Validation failed",
  "errors": {
    "email": ["Must be a valid email address."],
    "password": ["Must be at least 8 characters long and contain both letters and digits."]
  }
}
```

### Simple Format (Vehicle Handler)

**Request**:
```json
{
  "title": "Car",
  "price": -100
}
```

**Response**:
```json
{
  "status": "error",
  "message": "Validation failed",
  "errors": {
    "Title": "Title must be at least 5 characters",
    "Price": "Price must be greater than 0"
  }
}
```

---

## Testing

To verify the fix works:

1. **Test Auth Endpoint**:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid",
    "password": "weak"
  }'
```

Expected: Validation errors in array format with snake_case fields

2. **Test Vehicle Endpoint**:
```bash
curl -X POST http://localhost:8080/api/vehicles \
  -F "title=Car" \
  -F "price=-100"
```

Expected: Validation errors in simple format with original field names

---

## Files Summary

### Created
- ✅ `internal/utils/validation.go` (82 lines)

### Modified
- ✅ `internal/handlers/auth_handler.go`
  - Removed ~46 lines (duplicate functions)
  - Updated 5 function calls
  - Fixed imports

- ✅ `internal/handlers/vehicle_handler.go`
  - Removed ~24 lines (duplicate function)
  - Updated 1 function call
  - No import changes needed

### Net Result
- **Lines removed**: ~70 lines of duplicate code
- **Lines added**: 82 lines of reusable utilities
- **Compilation errors**: 0 ❌ → ✅

---

## Migration Notes

**No Breaking Changes**: The refactoring maintains the exact same API response format, so no frontend changes are required.

**Backward Compatible**: All existing validation error responses remain identical.

---

**Status**: ✅ **Complete**

**Fixed**: Duplicate function declaration compilation error
**Created**: Centralized validation utilities module
**Impact**: Zero breaking changes, improved code quality
