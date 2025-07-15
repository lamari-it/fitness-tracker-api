# FitFlow API Postman Workspace

This directory contains a complete Postman workspace for testing the FitFlow API, including collection and environment files.

## Files

- `FitFlow-API.postman_collection.json` - Complete API collection with all endpoints
- `FitFlow-API.postman_environment.json` - Environment variables and configuration
- `README.md` - This documentation file

## Setup Instructions

### 1. Import Collection and Environment

1. Open Postman
2. Click "Import" in the top left
3. Import both files:
   - `FitFlow-API.postman_collection.json`
   - `FitFlow-API.postman_environment.json`

### 2. Select Environment

1. In Postman, click the environment dropdown (top right)
2. Select "FitFlow API Environment"

### 3. Configure Environment Variables

The environment is pre-configured with:
- `base_url`: `http://localhost:8080` (default API URL)
- `test_email`: `test@example.com` (for testing)
- `test_password`: `testpassword123` (for testing)

Auto-populated variables (set by tests):
- `auth_token`: JWT token after login/register
- `user_id`: Current user ID
- `user_email`: Current user email

## Collection Structure

### 1. Health Check
- **GET** `/health` - Check API status

### 2. Authentication
- **POST** `/api/v1/auth/register` - Register new user
- **POST** `/api/v1/auth/login` - Login user
- **GET** `/api/v1/auth/profile` - Get user profile (protected)

### 3. OAuth
- **GET** `/api/v1/auth/google` - Google OAuth login
- **POST** `/api/v1/auth/apple` - Apple Sign-In

### 4. Protected Routes
- **GET** `/api/v1/dashboard` - User dashboard
- **GET** `/api/v1/workouts` - User workouts
- **GET** `/api/v1/nutrition` - User nutrition data

### 5. Error Scenarios
- Invalid login credentials
- Unauthorized access attempts
- Invalid token tests
- Duplicate registration tests

## Testing Workflow

### Quick Start Testing
1. Run "Health Check" to verify API is running
2. Run "Register User" to create a test account
3. Run "Login User" to get authentication token
4. Run any protected route to test authentication

### Authentication Flow Testing
1. **Register**: Creates new user and auto-saves token
2. **Login**: Authenticates user and auto-saves token
3. **Profile**: Tests token validation
4. **Protected Routes**: Test various authenticated endpoints

### Error Testing
1. Run requests in "Error Scenarios" folder
2. Verify proper error responses and status codes

## Auto-Generated Variables

The collection includes JavaScript test scripts that automatically:

- Extract JWT tokens from login/register responses
- Save user ID and email to environment variables
- Validate response structure and content
- Check response times and content types

## Request Examples

### Register User
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

### Login User
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

### Apple Sign-In
```json
{
  "identity_token": "your_apple_identity_token_here",
  "auth_code": "your_apple_auth_code_here",
  "first_name": "John",
  "last_name": "Doe"
}
```

## Common Issues

### 1. API Not Running
- Error: Connection refused
- Solution: Start the API with `go run main.go`

### 2. Database Connection Issues
- Error: Database connection failed
- Solution: Ensure PostgreSQL is running and configured

### 3. Missing Environment Variables
- Error: Variables not found
- Solution: Select "FitFlow API Environment" in Postman

### 4. Token Expired
- Error: 401 Unauthorized
- Solution: Run login request again to get new token

## Advanced Usage

### Environment Variables
You can modify environment variables for different setups:
- `base_url`: Change for different API environments
- `test_email`/`test_password`: Use different test credentials

### Running Tests
1. Individual: Click "Send" on any request
2. Collection: Click "..." on collection → "Run collection"
3. Folder: Click "..." on folder → "Run folder"

### Custom Tests
Each request includes test scripts that:
- Validate response status codes
- Check response structure
- Verify data types
- Auto-extract tokens and IDs

## OAuth Testing Notes

### Google OAuth
- The Google OAuth endpoint redirects to Google's consent screen
- Use this in a browser, not directly in Postman
- The callback will return a token

### Apple Sign-In
- Requires valid Apple identity token from iOS/web app
- Use actual Apple Sign-In integration to get tokens
- Replace placeholder tokens with real ones

## Collection Features

- **Automatic token management**: Tokens are extracted and used automatically
- **Response validation**: Built-in tests validate all responses
- **Error handling**: Comprehensive error scenario testing
- **Environment flexibility**: Easy to switch between environments
- **Documentation**: All requests include detailed descriptions

## Support

For API issues, refer to:
- Main API documentation: `../README.md`
- API source code: `../`
- Postman documentation: https://learning.postman.com/