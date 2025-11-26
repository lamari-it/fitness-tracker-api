# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

FitFlow API is a Go-based REST API for fitness tracking using Gin framework, PostgreSQL, and GORM ORM. Features JWT authentication, OAuth (Google/Apple), multi-language support, and comprehensive workout management.

## Essential Commands

### Development
```bash
make run              # Start the application
make dev              # Run with auto-reload (requires air)
make deps             # Install dependencies
make build            # Build binary to ./bin/fitflow-api
make test             # Run tests
make fmt              # Format code
make lint             # Run linter
```

### Database Operations
```bash
# Via Makefile
make migrate-up                        # Apply all migrations
make migrate-down-1                    # Rollback last migration
make migrate-create NAME=migration_name # Create new migration
make dev-reset                         # Reset database (migrate + seed)

# Via CLI tool (after make build)
./bin/fitflow-api db seed              # Run seeders
./bin/fitflow-api db seed:fresh        # Drop data and re-seed
./bin/fitflow-api db reset             # Full reset
./bin/fitflow-api db migrate           # Run migrations
./bin/fitflow-api db migrate:fresh     # Drop tables and migrate
```

### Docker
```bash
docker-compose up -d  # Start PostgreSQL (port 5467) + Adminer (port 8081)
```

## Architecture & Key Patterns

### Directory Structure
- `controllers/` - HTTP handlers, all return standardized JSON responses
- `models/` - GORM models with relationships and validations
- `middleware/` - Auth (JWT), i18n, admin/role checks
- `routes/` - Route definitions grouped by resource
- `utils/` - JWT handling, password hashing, response formatting
- `migrations/` - Database migrations (19 total)
- `cmd/fitflow/` - CLI commands for database operations

### Core Models & Relationships
```
User → WorkoutPlans → Workouts → WorkoutExercises → Exercises
User → UserEquipment → Equipment
Exercise → ExerciseMuscleGroups → MuscleGroups
Exercise → ExerciseEquipment → Equipment
WorkoutSession → ExerciseLogs → SetLogs
```

### API Response Format
All endpoints must return this structure:
```go
utils.RespondWithJSON(c, statusCode, utils.Response{
    Success: bool,
    Message: string,
    Data:    interface{},
    Errors:  map[string]string, // Optional
    Meta:    interface{},        // Optional, for pagination
})
```

### Authentication Flow
1. JWT tokens via `utils.GenerateJWT()` and `utils.ValidateJWT()`
2. Protected routes use `middleware.AuthMiddleware()`
3. Admin routes use `middleware.AdminMiddleware()` 
4. Role-based permissions via `roles` and `permissions` tables

### Database Conventions
- All models embed `gorm.Model` (ID, CreatedAt, UpdatedAt, DeletedAt)
- Use UUID for IDs via `google/uuid`
- Soft deletes enabled by default
- Foreign key constraints defined in models
- Indexes on frequently queried fields

### Key Features
- **Set Groups**: Support supersets, circuits, drop sets via `set_groups` table
- **Dual Prescriptions**: Exercises can be time-based or rep-based (`is_timed` flag)
- **Weight Units**: Flexible kg/lbs per set (`weight_unit` in `set_logs`)
- **Multi-language**: i18n via `X-Language` header, translations in `locales/`
- **OAuth**: Google and Apple Sign-In integrated

## Environment Configuration

Required `.env` variables:
```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
JWT_SECRET, JWT_EXPIRES_IN
GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URL
APPLE_CLIENT_ID, APPLE_TEAM_ID, APPLE_KEY_ID, APPLE_PRIVATE_KEY_PATH
USE_MIGRATIONS=false  # Set true for golang-migrate, false for GORM AutoMigrate
APP_ENV=development
```

## Testing

API testing via Postman collection in `postman/` directory with:
- Environment variables for base URL and auth tokens
- Auto-token management in collection scripts
- Comprehensive endpoint coverage

**IMPORTANT**: The Postman collection must be kept in sync with the codebase. When making changes:
- Adding endpoints: Add corresponding requests to the appropriate folder in the collection
- Modifying request/response schemas: Update the request body and example responses
- Removing endpoints: Remove the corresponding requests from the collection

## API Documentation

Swagger/OpenAPI documentation is in `swagger/` directory:
- `swagger.yaml` - Main OpenAPI 3.0.3 specification entry point
- `components/schemas/` - Data models organized by domain (auth, user, exercise, workout, session, trainer, fitness)
- `components/parameters.yaml` - Reusable path/query parameters
- `components/responses.yaml` - Standard response definitions
- `paths/` - API endpoint definitions grouped by resource

**IMPORTANT**: The Swagger documentation must be kept in sync with the codebase. When making changes:
- Adding/modifying endpoints: Update the corresponding `paths/*.yaml` file
- Adding/modifying request/response schemas: Update the corresponding `components/schemas/*.yaml` file
- Adding new resources: Create new path and schema files as needed
- Update `swagger.yaml` if adding new path or schema references

## Common Development Tasks

### Adding a New API Endpoint
1. Define model in `models/`
2. Create migration: `make migrate-create NAME=add_feature`
3. Add controller in `controllers/`
4. Define routes in `routes/`
5. Use `utils.RespondWithJSON()` for responses
6. Add middleware for auth/admin as needed
7. Add seeders in `database/seeds/...`
8. Add new endpoint to Postman collection
9. Add new endpoint to Swagger docs
10. Add new endpoint to README

### Working with Migrations
- Migrations in `migrations/` use golang-migrate format
- Naming: `YYYYMMDDHHMMSS_description.up.sql` and `.down.sql`
- Always include rollback logic in down migrations
- Test with `make migrate-up` and `make migrate-down-1`

### Handling Translations
IGNORE Translations for now
- Add translations to `locales/{lang}/messages.json`
- Use `middleware.TranslateMessage(c, "key")` in controllers
- Client sets language via `X-Language` header