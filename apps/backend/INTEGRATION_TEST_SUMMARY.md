# Integration Testing Setup - Quick Reference

## What Was Created

A complete integration testing framework for the backend with Docker Compose isolation.

## Files Created

```
apps/backend/
├── docker-compose.test.yaml           # Isolated test database (port 5433)
├── configs.test.toml                  # Test configuration
├── INTEGRATION_TESTING.md             # Full documentation
├── tests/
│   └── integration/
│       ├── setup_test.go              # Test server initialization
│       └── auth_test.go               # Auth endpoint tests (5 test suites)
└── internal/
    └── testutil/
        ├── testutil.go                # HTTP helpers (MakeRequest, ParseResponse, AssertStatusCode)
        ├── database.go                # DB setup/teardown, seeding
        └── fixtures.go                # Test data creation (users, teams, products)
```

## Makefile Commands Added

```bash
make test-integration          # Run all integration tests
make test-integration-coverage # Run with coverage report
make test-setup               # Start test database only
make test-teardown            # Stop test database
make test-clean               # Clean up all test resources
```

## Test Suites Implemented

### Auth Routes (`auth_test.go`)
1. **TestAuthRegister** - User registration (5 test cases)
   - Successful registration
   - Duplicate username
   - Duplicate email
   - Missing required fields
   - Password too short

2. **TestAuthLogin** - User login (5 test cases)
   - Login with username
   - Login with email
   - Wrong password
   - Non-existent user
   - Missing credentials

3. **TestAuthRefreshToken** - Token refresh (3 test cases)
   - Successful refresh
   - Missing token
   - Invalid token

4. **TestAuthLogout** - Logout functionality (1 test case)
   - Successful logout with cookie clearing

5. **TestAuthProtectedRoute** - Protected endpoint access (3 test cases)
   - Valid token
   - No token
   - Invalid token

## Quick Start

```bash
# Run all tests
make test-integration

# Run with coverage
make test-integration-coverage

# View coverage in browser
open apps/backend/coverage.html
```

## Test Architecture

```
┌─────────────────────────────────────────────────┐
│  Integration Test                               │
│  ┌───────────────────────────────────────────┐  │
│  │  SetupTestServer()                        │  │
│  │  - Creates isolated DB connection        │  │
│  │  - Runs migrations                        │  │
│  │  - Seeds basic data                       │  │
│  │  - Initializes Gin router                 │  │
│  │  - Sets up controllers & middleware       │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  ┌───────────────────────────────────────────┐  │
│  │  Test Execution                           │  │
│  │  - Make HTTP request                      │  │
│  │  - Assert status code                     │  │
│  │  - Parse JSON response                    │  │
│  │  - Verify response data                   │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  ┌───────────────────────────────────────────┐  │
│  │  Teardown()                               │  │
│  │  - Clean database tables                  │  │
│  │  - Close connections                      │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
          ▼
┌─────────────────────────────────────────────────┐
│  Docker Compose Test Environment               │
│  - PostgreSQL on port 5433                     │
│  - Data stored in tmpfs (fast)                 │
│  - Isolated from dev database                  │
└─────────────────────────────────────────────────┘
```

## Key Features

1. **Isolated Environment**: Test database runs on separate port (5433)
2. **Fast Execution**: Uses tmpfs for in-memory storage
3. **Complete Integration**: Tests full HTTP → DB → Response cycle
4. **Easy Fixtures**: Helper functions to create test data
5. **Comprehensive Coverage**: Success and error cases
6. **Clean Separation**: Test code isolated from production code

## Example Usage

```go
func TestYourFeature(t *testing.T) {
    // Setup
    ts := SetupTestServer(t)
    defer ts.Teardown(t)

    // Create test data
    user := testutil.CreateTestUser(t, ts.DB, "test@example.com", "testuser", "pass123")

    // Make request
    w := testutil.MakeRequest(t, ts.Router, "GET", "/api/endpoint", nil, nil)

    // Assert response
    testutil.AssertStatusCode(t, w, http.StatusOK)

    var response map[string]interface{}
    testutil.ParseResponse(t, w, &response)
    assert.Contains(t, response, "data")
}
```

## Test Utilities

### HTTP Helpers
- `testutil.MakeRequest()` - Make HTTP requests with optional body and headers
- `testutil.ParseResponse()` - Parse JSON response into structs
- `testutil.AssertStatusCode()` - Assert HTTP status code

### Database Helpers
- `testutil.SetupTestDB()` - Initialize test database with migrations
- `testutil.TeardownTestDB()` - Clean up test database
- `testutil.SeedBasicData()` - Seed roles, team_roles, schema types

### Fixture Helpers
- `testutil.CreateTestUser()` - Create test user with hashed password
- `testutil.CreateTestTeam()` - Create team with owner membership
- `testutil.CreateTestProduct()` - Create product for a team
- `testutil.CreateTestProductCategory()` - Create product category

## Coverage

Run tests with coverage:
```bash
make test-integration-coverage
```

This generates:
- `apps/backend/coverage.out` - Coverage data
- `apps/backend/coverage.html` - HTML coverage report

## Next Steps

To add more integration tests:

1. Create new test file: `tests/integration/feature_test.go`
2. Use the same pattern as `auth_test.go`
3. Import testutil helpers
4. Follow table-driven test approach

See `INTEGRATION_TESTING.md` for full documentation and examples.
