# LamariFit API Documentation

## Overview
This document provides comprehensive documentation for the LamariFit API, including all endpoints, request/response formats, and database schema.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

---

## Standard Response Format

All API responses follow a consistent structure to ensure predictable client-side handling.

### Success Response
```json
{
  "success": true,
  "message": "Resource fetched successfully.",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Workout Plan A",
    "exercises": [ ... ]
  },
  "errors": null
}
```

### Error Response
```json
{
  "success": false,
  "message": "Validation failed.",
  "data": null,
  "errors": {
    "name": ["The name field is required."],
    "email": ["Please provide a valid email address."]
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "Workouts fetched successfully.",
  "data": [ ... ],
  "errors": null,
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_pages": 5,
    "total_items": 47
  }
}
```

### Response Fields
- **success** (boolean): Indicates whether the request was successful
- **message** (string): Human-readable message describing the result
- **data** (any): The actual response data (null for errors)
- **errors** (object/null): Field-specific error messages for validation failures
- **meta** (object): Pagination metadata (only present for paginated responses)

### HTTP Status Codes
- **200 OK**: Successful GET, PUT requests
- **201 Created**: Successful POST requests that create resources
- **204 No Content**: Successful DELETE requests
- **400 Bad Request**: Validation errors or malformed requests
- **401 Unauthorized**: Missing or invalid authentication token
- **403 Forbidden**: Authenticated but lacks permission
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists
- **500 Internal Server Error**: Server-side errors

---

## Authentication Endpoints

### Register User
- **POST** `/auth/register`
- **Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "password_confirm": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```
- **Response:** User object with JWT token
- **Validation:** Password and password_confirm must match

### Login User
- **POST** `/auth/login`
- **Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
- **Response:** User object with JWT token

### Get User Profile
- **GET** `/auth/profile`
- **Headers:** Authorization: Bearer <token>
- **Response:** Current user profile

### Google OAuth
- **GET** `/auth/google`
- Redirects to Google OAuth consent screen

### Apple Sign-In
- **POST** `/auth/apple`
- **Body:**
```json
{
  "identity_token": "apple_identity_token",
  "auth_code": "apple_auth_code",
  "first_name": "John",
  "last_name": "Doe"
}
```

---

## Exercise Endpoints

### Create Exercise
- **POST** `/exercises`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "name": "Push-ups",
  "description": "A bodyweight exercise targeting chest, shoulders, and triceps",
  "is_bodyweight": true,
  "instructions": "Start in plank position, lower your body until chest nearly touches the floor.",
  "video_url": "https://example.com/pushup-video",
  "muscle_groups": [
    {
      "muscle_group_id": "uuid",
      "primary": true,
      "intensity": "high"
    }
  ]
}
```
- **Note:** The `slug` field is automatically generated from the exercise name (e.g., "Push-ups" becomes "push-ups")

### Get All Exercises
- **GET** `/exercises`
- **Headers:** Authorization: Bearer <token>
- **Query Parameters:**
  - `search` - Search by name
  - `muscle_group_id` - Filter by muscle group
  - `bodyweight` - Filter by bodyweight (true/false)
  - `primary_only` - Filter by primary muscle groups only

### Get Exercise by ID
- **GET** `/exercises/:id`
- **Headers:** Authorization: Bearer <token>

### Get Exercise by Slug
- **GET** `/exercises/by-slug/:slug`
- **Headers:** Authorization: Bearer <token>
- **Response:** Exercise details with muscle groups and equipment
- **Example:** `/exercises/by-slug/push-ups`

### Update Exercise
- **PUT** `/exercises/:id`
- **Headers:** Authorization: Bearer <token>
- **Body:** Same as create exercise

### Delete Exercise
- **DELETE** `/exercises/:id`
- **Headers:** Authorization: Bearer <token>

### Assign Equipment to Exercise
- **POST** `/exercises/:exercise_id/equipment`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "equipment_id": "uuid",
  "optional": false,
  "notes": "Use standard barbell"
}
```

### Get Exercise Equipment
- **GET** `/exercises/:exercise_id/equipment`
- **Headers:** Authorization: Bearer <token>

### Remove Equipment from Exercise
- **DELETE** `/exercises/:exercise_id/equipment/:equipment_id`
- **Headers:** Authorization: Bearer <token>

---

## Equipment Endpoints

### Create Equipment
- **POST** `/equipment`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "name": "Barbell",
  "slug": "barbell",
  "description": "Standard olympic barbell",
  "category": "free_weight",
  "image_url": "https://example.com/barbell.jpg"
}
```
- **Categories:** machine, free_weight, cable, cardio, other
- **Note:** Slug must be unique across all equipment

### Get All Equipment
- **GET** `/equipment`
- **Headers:** Authorization: Bearer <token>
- **Query Parameters:**
  - `search` - Search by name
  - `category` - Filter by category
  - `page` - Page number (default: 1)
  - `limit` - Items per page (default: 10, max: 100)

### Get Equipment by ID
- **GET** `/equipment/:id`
- **Headers:** Authorization: Bearer <token>
- **Response:** Equipment details with exercise count

### Update Equipment
- **PUT** `/equipment/:id`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "name": "Olympic Barbell",
  "slug": "olympic_barbell",
  "description": "45lb olympic barbell",
  "category": "free_weight",
  "image_url": "https://example.com/olympic-barbell.jpg"
}
```

### Delete Equipment
- **DELETE** `/equipment/:id`
- **Headers:** Authorization: Bearer <token>
- **Note:** Cannot delete if equipment is assigned to exercises

---

## Fitness Level Endpoints

### Get All Fitness Levels
- **GET** `/fitness-levels`
- **Headers:** Authorization: Bearer <token>
- **Response:** List of all fitness levels sorted by sort_order

### Get Fitness Level by ID
- **GET** `/fitness-levels/:id`
- **Headers:** Authorization: Bearer <token>
- **Response:** Fitness level details with user count

### Create Fitness Level
- **POST** `/fitness-levels`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "name": "Expert",
  "description": "Professional athlete level",
  "sort_order": 4
}
```

### Update Fitness Level
- **PUT** `/fitness-levels/:id`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "name": "Expert",
  "description": "Elite athlete level",
  "sort_order": 4
}
```

### Delete Fitness Level
- **DELETE** `/fitness-levels/:id`
- **Headers:** Authorization: Bearer <token>
- **Note:** Cannot delete if users are assigned to this level

---

## Fitness Goal Endpoints

### Get All Fitness Goals
- **GET** `/fitness-goals`
- **Headers:** Authorization: Bearer <token>
- **Response:** List of all fitness goals

### Get Fitness Goal by ID
- **GET** `/fitness-goals/:id`
- **Headers:** Authorization: Bearer <token>
- **Response:** Fitness goal details with user count

### Create Fitness Goal
- **POST** `/fitness-goals`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "name": "Cardiovascular Health"
}
```

### Update Fitness Goal
- **PUT** `/fitness-goals/:id`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "name": "Heart Health"
}
```

### Delete Fitness Goal
- **DELETE** `/fitness-goals/:id`
- **Headers:** Authorization: Bearer <token>
- **Note:** Cannot delete if users have selected this goal

---

## User Fitness Settings Endpoints

### Get User Fitness Goals
- **GET** `/user/fitness/goals`
- **Headers:** Authorization: Bearer <token>
- **Response:** List of user's selected fitness goals with details

### Set User Fitness Goals
- **PUT** `/user/fitness/goals`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "goals": [
    {
      "fitness_goal_id": "uuid1",
      "priority": 1,
      "target_date": "2024-12-31T00:00:00Z",
      "notes": "Lose 20 pounds by end of year"
    },
    {
      "fitness_goal_id": "uuid2",
      "priority": 2,
      "target_date": null,
      "notes": "Build muscle mass"
    }
  ]
}
```
- **Note:** This replaces all existing goals with the provided list
- **Priority:** 1 = primary, 2 = secondary, etc.

### Update User Fitness Level
- **PUT** `/user/fitness/level`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "fitness_level_id": "uuid"
}
```

---

## User Equipment Endpoints

### Get User Equipment
- **GET** `/user/equipment`
- **Headers:** Authorization: Bearer <token>
- **Query Parameters:**
  - `location_type` - Filter by location (home, gym)
- **Response:**
```json
{
  "equipment": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "equipment_id": "uuid",
      "location_type": "home",
      "gym_location": null,
      "notes": "Adjustable 5-50lbs",
      "equipment": {
        "id": "uuid",
        "name": "Dumbbells",
        "slug": "dumbbells",
        "description": "Adjustable or fixed weight dumbbells",
        "category": "free_weight",
        "image_url": null
      },
      "created_at": "2024-01-20T10:00:00Z",
      "updated_at": "2024-01-20T10:00:00Z"
    }
  ],
  "total": 1
}
```

### Add User Equipment
- **POST** `/user/equipment`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "equipment_id": "uuid",
  "location_type": "home",
  "gym_location": null,
  "notes": "20lb dumbbells"
}
```
- **Location Types:** home, gym
- **Note:** Cannot add duplicate equipment for the same location

### Bulk Add User Equipment
- **POST** `/user/equipment/bulk`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "equipment": [
    {
      "equipment_id": "uuid1",
      "location_type": "home",
      "notes": "Adjustable 5-50lbs"
    },
    {
      "equipment_id": "uuid2",
      "location_type": "gym",
      "gym_location": "Downtown Fitness",
      "notes": null
    }
  ]
}
```
- **Response:** List of created equipment entries
- **Note:** Skips duplicates without error

### Update User Equipment
- **PUT** `/user/equipment/:id`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "location_type": "gym",
  "gym_location": "Uptown Gym",
  "notes": "Recently upgraded to 30lb"
}
```
- **Note:** Can only update your own equipment

### Remove User Equipment
- **DELETE** `/user/equipment/:id`
- **Headers:** Authorization: Bearer <token>
- **Note:** Can only delete your own equipment

### Get User Equipment by Location
- **GET** `/user/equipment/location/:location`
- **Headers:** Authorization: Bearer <token>
- **Parameters:**
  - `location` - Either "home" or "gym"
- **Response:**
```json
{
  "location": "home",
  "equipment_by_category": {
    "free_weight": [
      {
        "id": "uuid",
        "equipment": {
          "name": "Dumbbells",
          "slug": "dumbbells",
          "category": "free_weight"
        },
        "notes": "Adjustable 5-50lbs"
      }
    ],
    "other": [
      {
        "id": "uuid",
        "equipment": {
          "name": "Resistance Bands",
          "slug": "resistance_bands",
          "category": "other"
        },
        "notes": "Light, medium, heavy bands"
      }
    ]
  },
  "total": 2
}
```

---

## Workout Plan Endpoints

### Create Workout Plan
- **POST** `/workout-plans`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "title": "Full Body Workout",
  "description": "A comprehensive full body workout plan",
  "visibility": "private"
}
```

### Get All Workout Plans
- **GET** `/workout-plans`
- **Headers:** Authorization: Bearer <token>
- **Query Parameters:**
  - `page` - Page number (default: 1)
  - `limit` - Items per page (default: 10)

### Get Workout Plan by ID
- **GET** `/workout-plans/:id`
- **Headers:** Authorization: Bearer <token>
- **Response:** Workout plan with associated workouts and exercises

### Update Workout Plan
- **PUT** `/workout-plans/:id`
- **Headers:** Authorization: Bearer <token>
- **Body:** Same as create workout plan

### Delete Workout Plan
- **DELETE** `/workout-plans/:id`
- **Headers:** Authorization: Bearer <token>

---

## Set Group Endpoints

### Create Set Group
- **POST** `/set-groups`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "workout_id": "uuid",
  "group_type": "superset",
  "name": "Chest & Back Superset",
  "notes": "Perform exercises back-to-back with minimal rest",
  "order_number": 1,
  "rest_between_sets": 60,
  "rounds": 3
}
```
- **Group Types:** `straight`, `superset`, `circuit`, `giant_set`, `drop_set`, `pyramid`, `rest_pause`

### Get Set Groups for Workout
- **GET** `/workouts/:workout_id/set-groups`
- **Headers:** Authorization: Bearer <token>
- **Response:** List of set groups with associated exercises

### Add Exercise to Set Group
- **POST** `/workout-exercises`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "workout_id": "uuid",
  "set_group_id": "uuid",
  "exercise_id": "uuid",
  "order_number": 1,
  "target_sets": 3,
  "target_reps": 10,
  "target_weight": 135.5,
  "target_rest_sec": 60
}
```
- **Note:** All exercises must belong to a SetGroup. The set_group_id is required.

### Exercise Logging with Set Groups
When logging exercises during a workout session, the `set_group_id` can be included to track which set group configuration was being performed:
- The `set_group_id` in `exercise_logs` is nullable to support free-form workouts
- When present, it links the logged exercise to the specific set group (superset, circuit, etc.) that was being performed
- This allows tracking whether exercises were performed as supersets, circuits, or other groupings during the actual workout

### Weight Units in Set Logs
When logging individual sets, the weight unit can be specified:
- Default weight unit is `kg` (kilograms)
- Supported units: `kg`, `lbs`, `lb`
- Each set can have its own weight unit, allowing flexibility for different equipment or user preferences

---

## Friend System Endpoints

### Send Friend Request
- **POST** `/friends/request`
- **Headers:** Authorization: Bearer <token>
- **Body:**
```json
{
  "friend_email": "friend@example.com"
}
```

### Get Friend Requests
- **GET** `/friends/requests`
- **Headers:** Authorization: Bearer <token>
- **Response:** List of pending friend requests

### Respond to Friend Request
- **PUT** `/friends/requests/:id/:action`
- **Headers:** Authorization: Bearer <token>
- **URL Parameters:**
  - `:id` - Friend request ID
  - `:action` - "accept" or "decline"

### Get Friends List
- **GET** `/friends`
- **Headers:** Authorization: Bearer <token>
- **Response:** List of accepted friends

### Remove Friend
- **DELETE** `/friends/:id`
- **Headers:** Authorization: Bearer <token>
- **URL Parameters:**
  - `:id` - Friendship ID

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    provider VARCHAR DEFAULT 'local',
    google_id VARCHAR UNIQUE,
    apple_id VARCHAR UNIQUE,
    fitness_level_id UUID REFERENCES fitness_levels(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Trainer Profiles Table
```sql
CREATE TABLE trainer_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    specialties TEXT[],
    hourly_rate NUMERIC(10,2),
    location TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Exercises Table
```sql
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(255) NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    is_bodyweight BOOLEAN DEFAULT false,
    instructions TEXT,
    video_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Equipment Table
```sql
CREATE TABLE equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) UNIQUE,
    description TEXT,
    category VARCHAR(50),
    image_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Exercise Equipment Table
```sql
CREATE TABLE exercise_equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    optional BOOLEAN DEFAULT false,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(exercise_id, equipment_id)
);
```

### User Equipment Table
```sql
CREATE TABLE user_equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    location_type VARCHAR(10) NOT NULL CHECK (location_type IN ('home', 'gym')),
    gym_location TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, equipment_id, location_type)
);
```

### Fitness Levels Table
```sql
CREATE TABLE fitness_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Fitness Goals Table
```sql
CREATE TABLE fitness_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### User Fitness Goals Table
```sql
CREATE TABLE user_fitness_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fitness_goal_id UUID NOT NULL REFERENCES fitness_goals(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, fitness_goal_id)
);
```

### Workout Plans Table
```sql
CREATE TABLE workout_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    visibility VARCHAR(20) DEFAULT 'private',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Workouts Table
```sql
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES workout_plans(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    day_number INTEGER NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Set Groups Table
```sql
CREATE TABLE set_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    group_type VARCHAR(20) NOT NULL DEFAULT 'straight',
    name VARCHAR(255),
    notes TEXT,
    order_number INTEGER NOT NULL,
    rest_between_sets INTEGER,
    rounds INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT check_group_type CHECK (group_type IN ('straight', 'superset', 'circuit', 'giant_set', 'drop_set', 'pyramid', 'rest_pause'))
);
```

### Workout Exercises Table
```sql
CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    set_group_id UUID NOT NULL REFERENCES set_groups(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    order_number INTEGER NOT NULL,
    target_sets INTEGER,
    target_reps INTEGER,
    target_weight NUMERIC(10,2),
    target_rest_sec INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Friendships Table
```sql
CREATE TABLE friendships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Workout Sessions Table
```sql
CREATE TABLE workout_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_id UUID REFERENCES workouts(id) ON DELETE SET NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Exercise Logs Table
```sql
CREATE TABLE exercise_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    set_group_id UUID REFERENCES set_groups(id) ON DELETE SET NULL,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    order_number INTEGER NOT NULL,
    notes TEXT,
    difficulty_rating INTEGER CHECK (difficulty_rating >= 1 AND difficulty_rating <= 10),
    difficulty_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Set Logs Table
```sql
CREATE TABLE set_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_log_id UUID NOT NULL REFERENCES exercise_logs(id) ON DELETE CASCADE,
    set_number INTEGER NOT NULL,
    weight NUMERIC(10,2),
    weight_unit VARCHAR(5) DEFAULT 'kg',
    reps INTEGER,
    rest_after_sec INTEGER,
    tempo VARCHAR(10),
    rpe NUMERIC(3,1) CHECK (rpe >= 1 AND rpe <= 10),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT check_weight_unit CHECK (weight_unit IN ('kg', 'lbs', 'lb'))
);
```

---

## Error Responses

All endpoints return consistent error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request data"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header is required"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 409 Conflict
```json
{
  "error": "Resource already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Rate Limiting

Currently, there are no rate limits implemented. Consider adding rate limiting for production use.

## Pagination

List endpoints support pagination with the following parameters:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)

Response format:
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

## Filtering and Search

Many endpoints support filtering and search:
- Use query parameters for filtering
- Use `search` parameter for text search
- Multiple filters can be combined

Example:
```
GET /exercises?search=push&muscle_group=chest&bodyweight=true
```

---

## Future Enhancements

### Planned Features
1. **Workout Logging** - Track actual workout sessions
2. **Social Features** - Share workouts and progress
3. **Progress Tracking** - Charts and analytics
4. **Nutrition Tracking** - Meal planning and logging
5. **Trainer Features** - Client management and programs
6. **Mobile App Support** - Enhanced mobile API features

### API Versioning
The API uses URL versioning (`/api/v1/`). Future versions will be available at `/api/v2/`, etc.

### WebSocket Support
Real-time features may be added using WebSocket connections for:
- Live workout tracking
- Real-time friend activity
- Instant messaging between trainers and clients

---

## Role-Based Access Control (RBAC)

### Overview
The API implements a comprehensive RBAC system with roles, permissions, and user-role associations. This replaces the previous `is_admin` flag with a more flexible permission-based system.

### Database Schema

#### Roles Table
- `id` (SERIAL PRIMARY KEY): Unique identifier
- `name` (VARCHAR(50) UNIQUE NOT NULL): Role name (e.g., 'admin', 'user', 'trainer')
- `description` (TEXT): Description of the role
- `created_at`, `updated_at`, `deleted_at`: Timestamps

#### Permissions Table
- `id` (SERIAL PRIMARY KEY): Unique identifier
- `name` (VARCHAR(100) UNIQUE NOT NULL): Permission name (e.g., 'exercises.create')
- `resource` (VARCHAR(50) NOT NULL): Resource being accessed (e.g., 'exercises')
- `action` (VARCHAR(50) NOT NULL): Action being performed (e.g., 'create', 'read', 'update', 'delete')
- `description` (TEXT): Description of the permission
- `created_at`, `updated_at`, `deleted_at`: Timestamps

#### Role_Permissions Table (Junction)
- `role_id` (INTEGER): Foreign key to roles table
- `permission_id` (INTEGER): Foreign key to permissions table
- `created_at`: Timestamp
- Primary Key: (role_id, permission_id)

#### User_Roles Table (Junction)
- `user_id` (UUID): Foreign key to users table
- `role_id` (INTEGER): Foreign key to roles table
- `created_at`: Timestamp
- Primary Key: (user_id, role_id)

### Default Roles

1. **Admin**: Full system access with all permissions
2. **Trainer**: Can manage workouts, clients, and view most resources
3. **User**: Basic access to personal data and public resources

### Permission Structure

Permissions follow the format: `resource.action`

Example permissions:
- `exercises.create`: Create new exercises
- `exercises.read`: View exercises
- `exercises.update`: Modify existing exercises
- `exercises.delete`: Remove exercises
- `workout_plans.create`: Create workout plans
- `translations.manage`: Manage translations (admin only)

### Migration Notes

- Existing users with `is_admin = true` are automatically assigned the 'admin' role
- Existing users with `is_admin = false` are assigned the 'user' role
- The `is_admin` field is maintained for backward compatibility but should be considered deprecated

### Seeding

Run the database seed to populate:
1. Default roles (admin, trainer, user)
2. All permissions based on existing endpoints
3. Role-permission mappings
4. User-role assignments for existing users

---

## Testing

Use the provided Postman collection in the `/postman/` directory for comprehensive API testing.

## Support

For API issues or questions, please refer to the main README.md file or create an issue in the project repository.# New API Endpoints Documentation

## Workouts

### Create Workout
Creates a new standalone workout for the authenticated user.

**Endpoint:** `POST /api/v1/workouts`

**Request Body:**
```json
{
  "title": "Morning Workout",
  "description": "Quick morning routine",
  "visibility": "private"  // Options: "private", "public", "friends"
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Workout created successfully.",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "456e7890-e89b-12d3-a456-426614174000",
    "title": "Morning Workout",
    "description": "Quick morning routine",
    "visibility": "private",
    "created_at": "2024-01-01T08:00:00Z",
    "updated_at": "2024-01-01T08:00:00Z"
  }
}
```

### Get User Workouts
Retrieves all workouts for the authenticated user with pagination.

**Endpoint:** `GET /api/v1/workouts`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Workouts fetched successfully.",
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Morning Workout",
      "description": "Quick morning routine",
      "visibility": "private",
      "set_groups": [],
      "exercises": []
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Get Workout by ID
Retrieves a specific workout with all its exercises and set groups.

**Endpoint:** `GET /api/v1/workouts/{workout_id}`

**Response:** `200 OK`

### Update Workout
Updates an existing workout's details.

**Endpoint:** `PUT /api/v1/workouts/{workout_id}`

**Request Body:**
```json
{
  "title": "Updated Morning Workout",
  "description": "Enhanced morning routine",
  "visibility": "friends"
}
```

**Response:** `200 OK`

### Delete Workout
Deletes a workout and all associated data.

**Endpoint:** `DELETE /api/v1/workouts/{workout_id}`

**Response:** `200 OK`

### Duplicate Workout
Creates a copy of an existing workout with all its exercises and set groups.

**Endpoint:** `POST /api/v1/workouts/{workout_id}/duplicate`

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Workout duplicated successfully.",
  "data": {
    "id": "789e4567-e89b-12d3-a456-426614174000",
    "title": "Morning Workout (Copy)",
    "description": "Quick morning routine",
    "visibility": "private"
  }
}
```

### Add Exercise to Workout
Adds an exercise to a workout with prescription details.

**Endpoint:** `POST /api/v1/workouts/{workout_id}/exercises`

**Request Body:**
```json
{
  "exercise_id": "123e4567-e89b-12d3-a456-426614174000",
  "set_group_id": "456e7890-e89b-12d3-a456-426614174000",
  "order_number": 1,
  "target_sets": 3,
  "target_reps": 12,
  "target_weight": 50,
  "target_rest_sec": 60,
  "prescription": "reps",  // Options: "reps", "time"
  "target_duration_sec": 0  // Used when prescription is "time"
}
```

**Response:** `201 Created`

### Get Workout Exercises
Retrieves all exercises in a workout, ordered by their position.

**Endpoint:** `GET /api/v1/workouts/{workout_id}/exercises`

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Workout exercises fetched successfully.",
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "workout_id": "456e7890-e89b-12d3-a456-426614174000",
      "exercise_id": "789e4567-e89b-12d3-a456-426614174000",
      "set_group_id": "abc1234-e89b-12d3-a456-426614174000",
      "order_number": 1,
      "target_sets": 3,
      "target_reps": 12,
      "target_weight": 50,
      "target_rest_sec": 60,
      "prescription": "reps",
      "exercise": {
        "id": "789e4567-e89b-12d3-a456-426614174000",
        "name": "Bench Press",
        "slug": "bench-press"
      },
      "set_group": {
        "id": "abc1234-e89b-12d3-a456-426614174000",
        "group_type": "straight",
        "name": "Main Set"
      }
    }
  ]
}
```

### Remove Exercise from Workout
Removes an exercise from a workout.

**Endpoint:** `DELETE /api/v1/workouts/{workout_id}/exercises/{exercise_id}`

**Response:** `200 OK`

---

## Workout Plan Items

### Add Workout to Plan
Adds an existing workout to a workout plan.

**Endpoint:** `POST /api/v1/workout-plans/{plan_id}/workouts`

**Request Body:**
```json
{
  "workout_id": "123e4567-e89b-12d3-a456-426614174000",
  "week_index": 0  // For multi-week plans (0-based)
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Workout added to plan successfully.",
  "data": {
    "id": "999e4567-e89b-12d3-a456-426614174000",
    "plan_id": "111e4567-e89b-12d3-a456-426614174000",
    "workout_id": "123e4567-e89b-12d3-a456-426614174000",
    "week_index": 0,
    "workout": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Leg Day"
    }
  }
}
```

### Get Plan Workouts
Retrieves all workouts in a workout plan, ordered by week and creation date.

**Endpoint:** `GET /api/v1/workout-plans/{plan_id}/workouts`

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Plan workouts fetched successfully.",
  "data": [
    {
      "id": "999e4567-e89b-12d3-a456-426614174000",
      "plan_id": "111e4567-e89b-12d3-a456-426614174000",
      "workout_id": "123e4567-e89b-12d3-a456-426614174000",
      "week_index": 0,
      "workout": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "title": "Leg Day",
        "exercises": []
      }
    }
  ]
}
```

### Update Plan Item
Updates a workout's position in the plan (e.g., move to different week).

**Endpoint:** `PUT /api/v1/workout-plans/{plan_id}/workouts/{item_id}`

**Request Body:**
```json
{
  "week_index": 1
}
```

**Response:** `200 OK`

### Remove Workout from Plan
Removes a workout from a workout plan.

**Endpoint:** `DELETE /api/v1/workout-plans/{plan_id}/workouts/{item_id}`

**Response:** `200 OK`

---

## Plan Enrollments

### Enroll in Plan
Enrolls the authenticated user in a workout plan with scheduling preferences.

**Endpoint:** `POST /api/v1/enrollments`

**Request Body:**
```json
{
  "plan_id": "123e4567-e89b-12d3-a456-426614174000",
  "start_date": "2024-01-01",
  "days_per_week": 3,
  "schedule_mode": "rolling",  // Options: "rolling", "calendar"
  "preferred_weekdays": [1, 3, 5]  // Only for calendar mode (0=Mon, 6=Sun)
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Successfully enrolled in workout plan.",
  "data": {
    "id": "555e4567-e89b-12d3-a456-426614174000",
    "plan_id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "789e4567-e89b-12d3-a456-426614174000",
    "start_date": "2024-01-01T00:00:00Z",
    "days_per_week": 3,
    "current_index": 0,
    "schedule_mode": "rolling",
    "preferred_weekdays": [],
    "status": "active",
    "created_at": "2024-01-01T08:00:00Z"
  }
}
```

### Get User Enrollments
Retrieves all enrollments for the authenticated user with optional filtering.

**Endpoint:** `GET /api/v1/enrollments`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `status` (optional): Filter by status ("active", "paused", "completed", "cancelled")

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Enrollments fetched successfully.",
  "data": [
    {
      "id": "555e4567-e89b-12d3-a456-426614174000",
      "plan_id": "123e4567-e89b-12d3-a456-426614174000",
      "start_date": "2024-01-01T00:00:00Z",
      "days_per_week": 3,
      "status": "active"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

### Get Enrollment by ID
Retrieves details of a specific enrollment.

**Endpoint:** `GET /api/v1/enrollments/{enrollment_id}`

**Response:** `200 OK`

### Update Enrollment
Updates enrollment settings such as days per week, schedule mode, or status.

**Endpoint:** `PUT /api/v1/enrollments/{enrollment_id}`

**Request Body:**
```json
{
  "days_per_week": 4,
  "schedule_mode": "calendar",
  "preferred_weekdays": [1, 2, 4, 5],
  "status": "paused"  // Options: "active", "paused", "completed"
}
```

**Response:** `200 OK`

### Cancel Enrollment
Cancels an active enrollment. This is a soft delete that changes status to "cancelled".

**Endpoint:** `DELETE /api/v1/enrollments/{enrollment_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Enrollment cancelled successfully.",
  "data": {
    "id": "555e4567-e89b-12d3-a456-426614174000",
    "status": "cancelled"
  }
}
```

---

## Database Schema Updates

### Workouts Table
```sql
CREATE TABLE workouts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  title TEXT NOT NULL,
  description TEXT,
  visibility VARCHAR(20) DEFAULT 'private',
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### Workout Plan Items Table
```sql
CREATE TABLE workout_plan_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_id UUID NOT NULL REFERENCES workout_plans(id),
  workout_id UUID NOT NULL REFERENCES workouts(id),
  week_index INT DEFAULT 0,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### Plan Enrollments Table
```sql
CREATE TABLE plan_enrollments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_id UUID NOT NULL REFERENCES workout_plans(id),
  user_id UUID NOT NULL REFERENCES users(id),
  start_date TIMESTAMP NOT NULL,
  days_per_week INT NOT NULL CHECK (days_per_week >= 1 AND days_per_week <= 7),
  current_index INT DEFAULT 0,
  schedule_mode VARCHAR(20) DEFAULT 'rolling',
  preferred_weekdays INT[] DEFAULT '{}',
  status VARCHAR(20) DEFAULT 'active',
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### Workout Exercises Table (Updated)
```sql
CREATE TABLE workout_exercises (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workout_id UUID NOT NULL REFERENCES workouts(id),
  set_group_id UUID NOT NULL REFERENCES set_groups(id),
  exercise_id UUID NOT NULL REFERENCES exercises(id),
  order_number INT NOT NULL,
  target_sets INT,
  target_reps INT,
  target_weight NUMERIC(10,2),
  target_rest_sec INT,
  prescription VARCHAR(20) DEFAULT 'reps',
  target_duration_sec INT DEFAULT 0,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

---

## Error Responses

### 400 Bad Request
Returned when trying to enroll in a plan you're already enrolled in, or when validation fails.

```json
{
  "success": false,
  "message": "You are already enrolled in this plan.",
  "data": null,
  "errors": null
}
```

### 404 Not Found
Returned when a requested resource doesn't exist or doesn't belong to the user.

```json
{
  "success": false,
  "message": "Workout not found.",
  "data": null,
  "errors": null
}
```

### 422 Unprocessable Entity
Returned for validation errors.

```json
{
  "success": false,
  "message": "Validation failed.",
  "data": null,
  "errors": {
    "days_per_week": ["Days per week must be between 1 and 7"],
    "schedule_mode": ["Schedule mode must be 'rolling' or 'calendar'"]
  }
}
```