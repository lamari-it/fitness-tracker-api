# Unified Weight System Documentation

## Overview

FitFlow API implements a **unified global weight handling system** that provides:

- **Canonical storage** in kilograms for all weights
- **Preservation** of original user input (value + unit)
- **Flexible input** accepting both kg and lb
- **User-specific responses** converted to preferred unit
- **Consistent API** across all weight-related endpoints

## Core Principles

### 1. Canonical Storage (PostgreSQL)

All weights are stored canonically in **kilograms** using `DECIMAL(6,2)`:

```sql
-- Every weight uses THREE fields:
{prefix}_weight_kg              DECIMAL(6,2)  -- Canonical kg storage
original_{prefix}_weight_value  DECIMAL(6,2)  -- Original value as entered
original_{prefix}_weight_unit   VARCHAR(2)    -- Original unit ('kg' or 'lb')
```

**Examples:**
- `current_weight_kg`, `original_current_weight_value`, `original_current_weight_unit`
- `target_weight_kg`, `original_target_weight_value`, `original_target_weight_unit`
- `actual_weight_kg`, `original_actual_weight_value`, `original_actual_weight_unit`

### 2. User Preferences

Weight unit preferences are stored at the **user level** (not profile level):

```sql
-- users table
preferred_weight_unit    VARCHAR(2) DEFAULT 'kg'  -- 'kg' or 'lb'
preferred_height_unit    VARCHAR(5) DEFAULT 'cm'  -- 'cm', 'ft', or 'ft_in'
preferred_distance_unit  VARCHAR(2) DEFAULT 'km'  -- 'km' or 'mi'
```

### 3. Conversion Factor

Uses the exact conversion factor:

```
1 lb = 0.45359237 kg
1 kg = 2.20462262 lb
```

## API Request Format

### Weight Input Structure

All weight inputs use the `WeightInput` format:

```json
{
  "current_weight": {
    "weight_value": 185.5,
    "weight_unit": "lb"
  }
}
```

**Go Structure:**
```go
type WeightInput struct {
    WeightValue *float64 `json:"weight_value,omitempty"`
    WeightUnit  *string  `json:"weight_unit,omitempty" binding:"omitempty,oneof=kg lb"`
}
```

### Examples

**Create Fitness Profile (kg):**
```json
POST /api/v1/user/fitness-profile
{
  "date_of_birth": "1990-01-15",
  "gender": "male",
  "height_cm": 180,
  "current_weight": {
    "weight_value": 80.0,
    "weight_unit": "kg"
  },
  "target_weight": {
    "weight_value": 75.0,
    "weight_unit": "kg"
  },
  "fitness_goal_ids": ["uuid-here"]
}
```

**Log Weight (lb):**
```json
POST /api/v1/user/fitness-profile/log-weight
{
  "weight": {
    "weight_value": 185.5,
    "weight_unit": "lb"
  }
}
```

**Create Workout Prescription (lb):**
```json
POST /api/v1/workouts/{id}/prescriptions
{
  "exercise_id": "uuid-here",
  "sets": 3,
  "reps": 10,
  "target_weight": {
    "weight_value": 225.0,
    "weight_unit": "lb"
  }
}
```

**Update Session Set (kg):**
```json
PUT /api/v1/session-sets/{id}
{
  "actual_reps": 10,
  "actual_weight": {
    "weight_value": 100.0,
    "weight_unit": "kg"
  },
  "completed": true
}
```

## API Response Format

### Weight Output Structure

All weight outputs use the `WeightOutput` format, **converted to user's preferred unit**:

```json
{
  "current_weight": {
    "weight_value": 185.5,
    "weight_unit": "lb"
  }
}
```

**Go Structure:**
```go
type WeightOutput struct {
    WeightValue *float64 `json:"weight_value,omitempty"`
    WeightUnit  *string  `json:"weight_unit,omitempty"`
}
```

### Response Behavior

1. **System reads** user's `preferred_weight_unit` from `users` table
2. **System converts** canonical kg value to preferred unit
3. **System returns** weight in user's preference

**Example:**
- User has `preferred_weight_unit = "lb"`
- Database stores: `current_weight_kg = 80.0`
- API returns:
  ```json
  {
    "current_weight": {
      "weight_value": 176.37,
      "weight_unit": "lb"
    }
  }
  ```

## Affected Endpoints

### User Fitness Profile

| Endpoint | Method | Weight Fields |
|----------|--------|---------------|
| `/api/v1/user/fitness-profile` | POST | `current_weight`, `target_weight` |
| `/api/v1/user/fitness-profile` | GET | Returns both weights in preferred unit |
| `/api/v1/user/fitness-profile` | PUT | `current_weight`, `target_weight` |
| `/api/v1/user/fitness-profile/log-weight` | POST | `weight` |

### Workout Prescriptions

| Endpoint | Method | Weight Fields |
|----------|--------|---------------|
| `/api/v1/workouts/{id}/prescriptions` | POST | `target_weight` |
| `/api/v1/workouts/{id}` | GET | Returns `target_weight` in preferred unit |
| `/api/v1/prescriptions/{id}` | PUT | `target_weight` |

### Workout Sessions (Logging)

| Endpoint | Method | Weight Fields |
|----------|--------|---------------|
| `/api/v1/session-sets/{id}` | PUT | `actual_weight` |
| `/api/v1/session-sets/{id}/complete` | POST | `actual_weight` (optional) |
| `/api/v1/session-exercises/{id}/sets` | POST | `actual_weight` |
| `/api/v1/sessions/{id}` | GET | Returns all `actual_weight` in preferred unit |

## Database Schema Changes

### Migration: `20251125090902_unified_weight_system`

**users table:**
```sql
ADD COLUMN preferred_weight_unit VARCHAR(2) DEFAULT 'kg'
ADD COLUMN preferred_height_unit VARCHAR(5) DEFAULT 'cm'
ADD COLUMN preferred_distance_unit VARCHAR(2) DEFAULT 'km'
```

**user_fitness_profiles table:**
```sql
DROP COLUMN preferred_weight_unit  -- Moved to users table

-- Current weight
ALTER COLUMN current_weight_kg TYPE DECIMAL(6,2)
ADD COLUMN original_current_weight_value DECIMAL(6,2)
ADD COLUMN original_current_weight_unit VARCHAR(2)

-- Target weight
ALTER COLUMN target_weight_kg TYPE DECIMAL(6,2)
ADD COLUMN original_target_weight_value DECIMAL(6,2)
ADD COLUMN original_target_weight_unit VARCHAR(2)
```

**workout_prescriptions table:**
```sql
DROP COLUMN weight_kg  -- Unused, removed

-- Target weight
ALTER COLUMN target_weight_kg TYPE DECIMAL(6,2)
ADD COLUMN original_target_weight_value DECIMAL(6,2)
ADD COLUMN original_target_weight_unit VARCHAR(2)
```

**session_sets table:**
```sql
ALTER COLUMN actual_weight_kg TYPE DECIMAL(6,2)
ADD COLUMN original_actual_weight_value DECIMAL(6,2)
ADD COLUMN original_actual_weight_unit VARCHAR(2)
```

## Code Examples

### Backend: Processing Weight Input

```go
import "fit-flow-api/utils"

// In controller
if req.CurrentWeight != nil {
    currentWeightKg, originalValue, originalUnit := utils.ProcessWeightInput(req.CurrentWeight)
    profile.CurrentWeightKg = currentWeightKg
    profile.OriginalCurrentWeightValue = originalValue
    profile.OriginalCurrentWeightUnit = originalUnit
}
```

### Backend: Converting for Response

```go
// Get user's preference
preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

// Convert to response
response := profile.ToResponse(preferredWeightUnit)
```

### Backend: Model ToResponse Methods

```go
func (p *UserFitnessProfile) ToResponse(preferredWeightUnit string) UserFitnessProfileResponse {
    response := UserFitnessProfileResponse{
        // ... other fields
    }

    // Convert current weight
    if p.CurrentWeightKg != nil {
        currentWeight := utils.ConvertWeightForResponse(p.CurrentWeightKg, preferredWeightUnit)
        response.CurrentWeight = currentWeight
    }

    return response
}
```

## Testing Guidelines

### Unit Tests

Update test requests to use `WeightInput` format:

```go
// OLD FORMAT (deprecated)
data := map[string]interface{}{
    "current_weight_kg": 80.0,
}

// NEW FORMAT
data := map[string]interface{}{
    "current_weight": map[string]interface{}{
        "weight_value": 80.0,
        "weight_unit":  "kg",
    },
}
```

Update test assertions for `WeightOutput` format:

```go
// OLD ASSERTION (deprecated)
data.Value("current_weight_kg").Number().IsEqual(80.0)
data.Value("current_weight_lbs").Number().Gt(0)

// NEW ASSERTION
currentWeight := data.Value("current_weight").Object()
currentWeight.Value("weight_value").Number().IsEqual(80.0)
currentWeight.Value("weight_unit").String().IsEqual("kg")
```

### Integration Tests

Test scenarios:
1. Create with kg, verify response in kg
2. Create with lb, verify response in lb
3. Update user preference, verify responses change accordingly
4. Mixed units across different fields
5. Nil/optional weight values

## Migration Guide

### For Existing Databases

If migrating from old schema:

1. **Backup database** before running migration
2. Run migration: `make migrate-up` or `USE_MIGRATIONS=true make migrate-up`
3. Existing weight data in kg will remain valid
4. Original value/unit fields will be NULL for old records
5. New records will populate all three fields

### For Fresh Installations

1. Models use GORM AutoMigrate by default
2. Run `make dev-reset` to create fresh schema
3. All weight fields will be created with correct structure

## Breaking Changes

⚠️ **API Breaking Changes:**

1. **Request format changed:**
   - Old: `"current_weight_kg": 80.0`
   - New: `"current_weight": {"weight_value": 80.0, "weight_unit": "kg"}`

2. **Response format changed:**
   - Old: `"current_weight_kg": 80.0, "current_weight_lbs": 176.37`
   - New: `"current_weight": {"weight_value": 80.0, "weight_unit": "kg"}`

3. **Preference location moved:**
   - Old: `user_fitness_profiles.preferred_weight_unit`
   - New: `users.preferred_weight_unit`

4. **Field removed:**
   - `workout_prescriptions.weight_kg` (unused, removed)

## Best Practices

### DO ✅

- Always use `WeightInput` for all weight inputs
- Always convert to user's preferred unit for responses
- Store canonical kg + original value/unit
- Validate weight units are 'kg' or 'lb' only
- Use `utils.ProcessWeightInput()` for processing
- Use model `ToResponse(preferredWeightUnit)` methods

### DON'T ❌

- Don't store weights in multiple units
- Don't expose canonical kg values directly in API responses
- Don't hardcode conversion factors (use utility functions)
- Don't forget to pass `preferredWeightUnit` to `ToResponse()`
- Don't use old `_kg` suffixed fields in requests/responses
- Don't store preferences in `user_fitness_profiles` table

## Support

For questions or issues:
- Check this documentation first
- Review code examples in `utils/weight_conversion.go`
- Check controller implementations for patterns
- See test examples for request/response formats

## Version History

- **v1.0** (2025-11-25): Initial unified weight system implementation
  - Canonical kg storage with original preservation
  - User-level preferences
  - Consistent WeightInput/WeightOutput DTOs
  - Migration for existing databases
