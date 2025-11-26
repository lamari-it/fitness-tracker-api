# LamariFit API

A Go REST API with PostgreSQL, GORM, JWT authentication, and OAuth integration for Google and Apple Sign-In.

## Features

- User registration and authentication
- JWT token-based authentication
- Google OAuth integration
- Apple Sign-In integration
- PostgreSQL database with GORM ORM
- Database migrations with golang-migrate
- Protected routes with middleware
- CORS support
- Multi-language support (i18n)

## Setup

### Prerequisites

- Go 1.19+
- PostgreSQL 12+
- Google OAuth credentials (for Google Sign-In)
- Apple Developer credentials (for Apple Sign-In)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd lamari-fit-api
```

2. Install dependencies:
```bash
go mod download
```

3. Create a PostgreSQL database:
```sql
CREATE DATABASE lamarifit;
```

4. Copy environment variables:
```bash
cp .env.example .env
```

5. Configure your `.env` file with your database and OAuth credentials.

6. Run the application:
```bash
go run main.go
```

The server will start on port 8080.

## Database Migrations

This project supports two database schema management approaches:

### 1. GORM AutoMigrate (Default)
The default approach uses GORM's AutoMigrate feature for development:
```bash
# Uses GORM AutoMigrate (default)
go run main.go
```

### 2. golang-migrate (Recommended for Production)
For production deployments, use proper database migrations:

#### Install migration tool:
```bash
make install-migrate
```

#### Migration commands:
```bash
# Run all migrations
make migrate-up

# Rollback last migration
make migrate-down-1

# Rollback all migrations (destructive)
make migrate-down

# Check migration status
make migrate-version

# Create new migration
make migrate-create NAME=add_new_table

# Development reset (drop all and recreate)
make dev-reset
```

#### Enable migrations in application:
Set `USE_MIGRATIONS=true` in your `.env` file to use golang-migrate instead of GORM AutoMigrate.

#### Migration files:
All migration files are located in the `migrations/` directory with the following naming convention:
- `000001_migration_name.up.sql` - Forward migration
- `000001_migration_name.down.sql` - Rollback migration

## API Endpoints

### Authentication

#### Register
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Google OAuth
```
GET /api/v1/auth/google
```
Redirects to Google OAuth consent screen.

```
GET /api/v1/auth/google/callback
```
Handles Google OAuth callback.

#### Apple Sign-In
```
POST /api/v1/auth/apple
Content-Type: application/json

{
  "identity_token": "apple_identity_token",
  "auth_code": "apple_auth_code",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Get Profile
```
GET /api/v1/auth/profile
Authorization: Bearer <jwt_token>
```

### Protected Routes

All protected routes require the `Authorization: Bearer <jwt_token>` header.

#### Dashboard
```
GET /api/v1/dashboard
Authorization: Bearer <jwt_token>
```

#### Workouts
```
GET /api/v1/workouts
Authorization: Bearer <jwt_token>
```

#### Nutrition
```
GET /api/v1/nutrition
Authorization: Bearer <jwt_token>
```

### Health Check
```
GET /health
```

## Database Schema

### Users Table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| email | VARCHAR | Unique email address |
| password | VARCHAR | Hashed password |
| first_name | VARCHAR | User's first name |
| last_name | VARCHAR | User's last name |
| provider | VARCHAR | Auth provider (local, google, apple) |
| google_id | VARCHAR | Google OAuth ID |
| apple_id | VARCHAR | Apple Sign-In ID |
| is_active | BOOLEAN | Account status |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | |
| DB_NAME | Database name | lamarifit |
| DB_SSLMODE | SSL mode | disable |
| USE_MIGRATIONS | Use golang-migrate instead of GORM AutoMigrate | false |
| JWT_SECRET | JWT signing secret | |
| JWT_EXPIRES_IN | JWT expiration time | 24h |
| GOOGLE_CLIENT_ID | Google OAuth client ID | |
| GOOGLE_CLIENT_SECRET | Google OAuth client secret | |
| GOOGLE_REDIRECT_URL | Google OAuth redirect URL | |
| APPLE_CLIENT_ID | Apple Sign-In client ID | |
| APPLE_TEAM_ID | Apple Developer team ID | |
| APPLE_KEY_ID | Apple Sign-In key ID | |
| APPLE_PRIVATE_KEY_PATH | Path to Apple private key | |
| APPLE_REDIRECT_URL | Apple Sign-In redirect URL | |

## Project Structure

```
lamari-fit-api/
├── config/         # Configuration management
├── controllers/    # HTTP handlers
├── database/       # Database connection and migrations
├── middleware/     # HTTP middleware
├── migrations/     # Database migration files
├── models/         # Data models
├── routes/         # Route definitions
├── utils/          # Utility functions
├── locales/        # Translation files
├── postman/        # Postman collection
├── main.go         # Application entry point
├── go.mod          # Go module dependencies
├── Makefile        # Build and migration commands
└── README.md       # This file
```

## OAuth Setup

### Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Create OAuth 2.0 credentials
5. Add your redirect URI: `http://localhost:8080/auth/google/callback`
6. Copy the client ID and secret to your `.env` file

### Apple Sign-In

1. Go to [Apple Developer Portal](https://developer.apple.com/)
2. Create a new App ID with Sign In with Apple capability
3. Create a new Service ID
4. Generate a new key for Sign In with Apple
5. Download the private key and update your `.env` file

## Security

- Passwords are hashed using bcrypt
- JWT tokens are signed with a secret key
- OAuth tokens are validated on the server side
- CORS is configured for cross-origin requests
- SQL injection protection via GORM

## Development

### Quick Start
```bash
# Install dependencies
make deps

# Run with GORM AutoMigrate (development)
make run

# Run with auto-reload
make dev
```

### Production Build
```bash
# Build for production
make build

# Build for production (optimized)
make build-prod
```

### Database Management
```bash
# Show all available commands
make help

# Install migration tool
make install-migrate

# Run migrations
make migrate-up

# Reset database (development)
make dev-reset
```

## License

This project is licensed under the MIT License.