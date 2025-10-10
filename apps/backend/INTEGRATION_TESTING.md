# Integration Testing Guide

This document describes the integration testing setup for the backend application.

## Overview

The integration testing framework allows you to test complete HTTP request/response cycles with a real PostgreSQL database running in Docker. This ensures your routes, controllers, services, and database interactions work correctly together.

## Architecture

### Components

1. **Docker Compose Test Environment** (`docker-compose.test.yaml`)
   - Isolated PostgreSQL test database on port 5433
   - Uses tmpfs for faster performance (data stored in memory)
   - Automatically seeded with basic data (roles, team_roles, schema types)

2. **Test Utilities** (`internal/testutil/`)
   - `testutil.go`: HTTP request helpers and assertions
   - `database.go`: Database setup, teardown, and seeding utilities
   - `fixtures.go`: Helper functions to create test data (users, teams, products, etc.)

3. **Test Setup** (`tests/integration/setup_test.go`)
   - `TestServer`: Encapsulates all test dependencies (router, DB, config, controllers)
   - `SetupTestServer()`: Initializes the complete test environment
   - Configures Gin router with middleware and routes

4. **Test Suites** (`tests/integration/*_test.go`)
   - Auth tests: Registration, login, token refresh, logout, protected routes
   - Additional test suites can be added for other features

## Directory Structure

```
apps/backend/
├── docker-compose.test.yaml          # Test database configuration
├── configs.test.toml                  # Test environment configuration
├── tests/
│   └── integration/
│       ├── setup_test.go              # Test server setup
│       ├── auth_test.go               # Auth integration tests
│       └── ...                        # Additional test files
└── internal/
    └── testutil/
        ├── testutil.go                # HTTP and assertion helpers
        ├── database.go                # Database setup utilities
        └── fixtures.go                # Test data creation helpers
```

## Running Tests

### Quick Start

```bash
# Run all integration tests
make test-integration

# Run with coverage report
make test-integration-coverage

# Clean up test resources
make test-clean
```

### Manual Setup

```bash
# 1. Start test database
make test-setup

# 2. Run tests
cd apps/backend
go test -v ./tests/integration/...

# 3. Stop test database
make test-teardown
```

### Run Specific Tests

```bash
# Run only auth tests
cd apps/backend
go test -v ./tests/integration/... -run TestAuth

# Run specific test case
go test -v ./tests/integration/... -run TestAuthRegister
```

## Writing Integration Tests

### Basic Structure

```go
func TestYourFeature(t *testing.T) {
    // Setup test server
    ts := SetupTestServer(t)
    defer ts.Teardown(t)

    // Create test data if needed
    testUser := testutil.CreateTestUser(t, ts.DB, "test@example.com", "testuser", "password123")

    // Define test cases
    tests := []struct {
        name           string
        payload        YourRequestDTO
        expectedStatus int
        checkResponse  func(t *testing.T, body map[string]interface{})
    }{
        {
            name: "successful request",
            payload: YourRequestDTO{
                Field: "value",
            },
            expectedStatus: http.StatusOK,
            checkResponse: func(t *testing.T, body map[string]interface{}) {
                assert.Contains(t, body, "data")
                // Add more assertions
            },
        },
        // More test cases...
    }

    // Run test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := testutil.MakeRequest(t, ts.Router, "POST", "/your/endpoint", tt.payload, nil)
            testutil.AssertStatusCode(t, w, tt.expectedStatus)

            var response map[string]interface{}
            testutil.ParseResponse(t, w, &response)

            if tt.checkResponse != nil {
                tt.checkResponse(t, response)
            }
        })
    }
}
```

### Using Test Utilities

#### Making HTTP Requests

```go
// Simple request
w := testutil.MakeRequest(t, router, "GET", "/endpoint", nil, nil)

// Request with JSON body
payload := dto.LoginDto{
    UserNameEmail: "user@example.com",
    Password:      "password123",
}
w := testutil.MakeRequest(t, router, "POST", "/auth/login", payload, nil)

// Request with headers (e.g., authorization)
headers := map[string]string{
    "Authorization": "Bearer " + token,
}
w := testutil.MakeRequest(t, router, "GET", "/protected", nil, headers)
```

#### Parsing Responses

```go
// Parse JSON response
var response map[string]interface{}
testutil.ParseResponse(t, w, &response)

// Access nested data
data := response["data"].(map[string]interface{})
token := data["token"].(string)
```

#### Asserting Status Codes

```go
testutil.AssertStatusCode(t, w, http.StatusOK)
testutil.AssertStatusCode(t, w, http.StatusBadRequest)
testutil.AssertStatusCode(t, w, http.StatusUnauthorized)
```

### Creating Test Data

#### Create Test User

```go
user := testutil.CreateTestUser(t, db, "user@example.com", "username", "password123")
// user.ID, user.Email, user.Username are now available
```

#### Create Test Team

```go
team := testutil.CreateTestTeam(t, db, "Team Name", "Description", ownerUserID)
// team.ID, team.Name are now available
// Owner is automatically added as a team member
```

#### Create Test Product

```go
product := testutil.CreateTestProduct(t, db, "Product Name", "Description", 99.99, teamID)
// product.ID, product.Name, product.Price are now available
```

#### Create Test Product Category

```go
category := testutil.CreateTestProductCategory(t, db, "Category Name", teamID)
// category.ID, category.Name are now available
```

### Database Seeding

Basic data (roles, team_roles, report_json_schema_types) is automatically seeded when you call `SetupTestServer()`. You can also manually seed data:

```go
// Seed all basic data
testutil.SeedBasicData(t, db)

// Or seed individually
testutil.SeedRoles(t, db)
testutil.SeedTeamRoles(t, db)
testutil.SeedReportSchemaTypes(t, db)
```

## Test Examples

### Example 1: Testing User Registration

```go
func TestAuthRegister(t *testing.T) {
    ts := SetupTestServer(t)
    defer ts.Teardown(t)

    payload := service.UserRegisterDto{
        Username: "newuser",
        Email:    "new@example.com",
        Password: "password123",
        FullName: "New User",
    }

    w := testutil.MakeRequest(t, ts.Router, "POST", "/auth/register", payload, nil)
    testutil.AssertStatusCode(t, w, http.StatusOK)

    var response map[string]interface{}
    testutil.ParseResponse(t, w, &response)

    assert.Contains(t, response, "data")
    data := response["data"].(map[string]interface{})
    assert.NotEmpty(t, data["token"])
    assert.NotEmpty(t, data["refresh_token"])
}
```

### Example 2: Testing Protected Routes

```go
func TestProtectedRoute(t *testing.T) {
    ts := SetupTestServer(t)
    defer ts.Teardown(t)

    // Create user and login
    user := testutil.CreateTestUser(t, ts.DB, "test@example.com", "testuser", "password123")

    loginPayload := service.LoginDto{
        UserNameEmail: "testuser",
        Password:      "password123",
    }
    loginResp := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", loginPayload, nil)

    var loginResponse map[string]interface{}
    testutil.ParseResponse(t, loginResp, &loginResponse)
    token := loginResponse["data"].(map[string]interface{})["token"].(string)

    // Test protected endpoint
    headers := map[string]string{
        "Authorization": "Bearer " + token,
    }
    w := testutil.MakeRequest(t, ts.Router, "GET", "/protected-endpoint", nil, headers)
    testutil.AssertStatusCode(t, w, http.StatusOK)
}
```

### Example 3: Testing with Multiple Test Cases

```go
func TestLogin(t *testing.T) {
    ts := SetupTestServer(t)
    defer ts.Teardown(t)

    testUser := testutil.CreateTestUser(t, ts.DB, "test@example.com", "testuser", "password123")

    tests := []struct {
        name           string
        payload        service.LoginDto
        expectedStatus int
    }{
        {
            name: "successful login",
            payload: service.LoginDto{
                UserNameEmail: "testuser",
                Password:      "password123",
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "wrong password",
            payload: service.LoginDto{
                UserNameEmail: "testuser",
                Password:      "wrongpassword",
            },
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "non-existent user",
            payload: service.LoginDto{
                UserNameEmail: "nonexistent",
                Password:      "password123",
            },
            expectedStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := testutil.MakeRequest(t, ts.Router, "POST", "/auth/login", tt.payload, nil)
            testutil.AssertStatusCode(t, w, tt.expectedStatus)
        })
    }
}
```

## Best Practices

1. **Test Isolation**: Each test should be independent and not rely on other tests
2. **Database Cleanup**: Use `defer ts.Teardown(t)` to ensure cleanup happens even if tests fail
3. **Descriptive Test Names**: Use clear, descriptive names for test cases
4. **Test Both Success and Failure**: Cover happy paths and error cases
5. **Use Table-Driven Tests**: Group related test cases using structs
6. **Check Response Structure**: Verify not just status codes but also response data
7. **Create Minimal Test Data**: Only create the data you need for each test
8. **Use Subtests**: Organize related tests with `t.Run()`

## Configuration

### Test Database

The test database runs on port **5433** (not 5432) to avoid conflicts with development databases.

Connection string:
```
host=localhost user=test_user password=test_password dbname=test_db port=5433 sslmode=disable
```

### Test Config

Test-specific configuration is in `configs.test.toml`:
- Isolated JWT secrets
- Test database connection
- Test mail settings
- Other test-specific configurations

## Troubleshooting

### Tests Failing to Connect to Database

```bash
# Ensure test database is running
make test-setup

# Check database status
docker compose -f apps/backend/docker-compose.test.yaml ps

# Check logs
docker compose -f apps/backend/docker-compose.test.yaml logs
```

### Tests Hanging

The test database setup waits for the database to be ready. If tests hang:
- Increase the sleep duration in `make test-setup`
- Check Docker resources (CPU, memory)
- Verify port 5433 is not already in use

### Migration Errors

```bash
# Clean and restart test database
make test-clean
make test-setup
```

### Coverage Reports Not Generated

```bash
# Ensure you have write permissions
chmod +w apps/backend

# Run with coverage
make test-integration-coverage

# Check apps/backend/coverage.html
```

## Adding New Test Suites

1. Create a new test file in `tests/integration/` (e.g., `product_test.go`)
2. Follow the standard structure (setup, test cases, assertions)
3. Use existing testutil helpers or add new ones as needed
4. Document any special setup requirements

Example:

```go
package integration_test

import (
    "testing"
    // ... imports
)

func TestProductCreate(t *testing.T) {
    ts := SetupTestServer(t)
    defer ts.Teardown(t)

    // Your test implementation
}
```

## Continuous Integration

These integration tests can be run in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run integration tests
  run: make test-integration-coverage

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    files: ./apps/backend/coverage.out
```

## Additional Resources

- [Testing in Go](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Gin Testing](https://github.com/gin-gonic/gin#testing)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
