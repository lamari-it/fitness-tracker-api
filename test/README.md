# FitFlow API Tests

This directory contains the test suite for the FitFlow API using the [httpexpect](https://github.com/gavv/httpexpect) library.

## Setup

### Prerequisites

1. PostgreSQL database running (default on port 5467 for tests)
2. Test database created (fitflow_test)
3. Environment configuration

### Configuration

Tests use `.env.test` file for configuration. Copy the example configuration:

```bash
cp .env.test.example .env.test
```

Default test database configuration:
- Host: localhost
- Port: 5467
- Database: fitflow_test
- User: fitflow_user
- Password: fitflow_password

## Running Tests

### Run all tests
```bash
make test
```

### Run authentication tests only
```bash
make test-auth
```

### Run tests with coverage
```bash
make test-coverage
```

### Clean test cache
```bash
make test-clean
```

## Test Structure

- `test_setup.go` - Test configuration, database setup, and helper functions
- `auth_test.go` - Authentication tests (registration, login, protected routes)

## Test Database

The test suite automatically:
1. Connects to the test database
2. Runs migrations
3. Cleans all data before each test run
4. Provides helper functions for creating test data

## Writing New Tests

Example test structure:

```go
func TestFeature(t *testing.T) {
    // Setup test app
    e := SetupTestApp(t)
    
    // Create test data if needed
    userData := CreateTestUser(t)
    
    // Run subtests
    t.Run("SubTest", func(t *testing.T) {
        response := e.GET("/api/v1/endpoint").
            Expect().
            Status(200).
            JSON().
            Object()
            
        response.Value("success").Boolean().IsTrue()
    })
}
```

## Helper Functions

- `SetupTestApp(t)` - Creates test application with httpexpect
- `SetupTestDatabase(t)` - Initializes test database
- `CleanDatabase(t)` - Clears all test data
- `CreateTestUser(t)` - Creates a standard test user
- `GetAuthToken(e, email, password)` - Gets JWT token for testing protected routes

## Coverage Reports

After running `make test-coverage`, view the coverage report:
- Terminal: Check coverage.out
- Browser: Open coverage.html