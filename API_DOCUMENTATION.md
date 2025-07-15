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
  "muscle_group": "chest",
  "equipment": "none",
  "is_bodyweight": true
}
```

### Get All Exercises
- **GET** `/exercises`
- **Headers:** Authorization: Bearer <token>
- **Query Parameters:**
  - `search` - Search by name
  - `muscle_group` - Filter by muscle group
  - `equipment` - Filter by equipment
  - `bodyweight` - Filter by bodyweight (true/false)

### Get Exercise by ID
- **GET** `/exercises/:id`
- **Headers:** Authorization: Bearer <token>

### Update Exercise
- **PUT** `/exercises/:id`
- **Headers:** Authorization: Bearer <token>
- **Body:** Same as create exercise

### Delete Exercise
- **DELETE** `/exercises/:id`
- **Headers:** Authorization: Bearer <token>

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
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    muscle_group TEXT,
    equipment TEXT,
    is_bodyweight BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
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

### Workout Exercises Table
```sql
CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
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
    reps INTEGER,
    rest_after_sec INTEGER,
    tempo VARCHAR(10),
    rpe NUMERIC(3,1) CHECK (rpe >= 1 AND rpe <= 10),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
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