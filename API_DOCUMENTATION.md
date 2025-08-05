# FitFlow API Documentation

## Overview
This document provides comprehensive documentation for the FitFlow API, including all endpoints, request/response formats, and database schema.

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

## Testing

Use the provided Postman collection in the `/postman/` directory for comprehensive API testing.

## Support

For API issues or questions, please refer to the main README.md file or create an issue in the project repository.