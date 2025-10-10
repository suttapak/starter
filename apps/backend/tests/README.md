# Backend Tests

This directory contains all tests for the backend application.

## Test Types

### Integration Tests (`integration/`)

Full end-to-end tests with real database and HTTP server.

**Run tests:**
```bash
make test-integration
```

**Features:**
- Tests complete HTTP request/response cycle
- Uses isolated PostgreSQL database in Docker
- Includes test utilities and fixtures
- Covers auth routes with 17+ test cases

**Documentation:** See `../INTEGRATION_TESTING.md` for full guide

### Unit Tests (throughout codebase)

Individual function/method tests located alongside the code they test.

**Run tests:**
```bash
cd apps/backend
go test ./...
```

**Examples:**
- `internal/service/auth.service_test.go`
- `internal/controller/auth.controller_test.go` (currently commented out)

## Quick Reference

```bash
# Run all tests (unit + integration)
cd apps/backend && go test -v ./...

# Run only integration tests
make test-integration

# Run with coverage
make test-integration-coverage

# Run specific test
go test -v ./tests/integration/... -run TestAuthRegister

# Setup test database only
make test-setup

# Clean up test resources
make test-clean
```

## File Structure

```
tests/
├── README.md                  # This file
└── integration/
    ├── setup_test.go         # Test server initialization
    ├── auth_test.go          # Auth endpoint integration tests
    └── ...                   # Add more test files here
```

## Writing Tests

### Integration Test Template

```go
package integration_test

import (
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/suttapak/starter/internal/testutil"
)

func TestYourFeature(t *testing.T) {
    ts := SetupTestServer(t)
    defer ts.Teardown(t)

    // Create test data
    user := testutil.CreateTestUser(t, ts.DB, "test@example.com", "user", "pass123")

    // Test your feature
    w := testutil.MakeRequest(t, ts.Router, "POST", "/api/endpoint", payload, nil)
    testutil.AssertStatusCode(t, w, http.StatusOK)

    var response map[string]interface{}
    testutil.ParseResponse(t, w, &response)

    assert.Contains(t, response, "data")
}
```

### Unit Test Template

```go
package service_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestYourFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "test case 1",
            input:    "input",
            expected: "expected",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := YourFunction(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Test Utilities

The `internal/testutil` package provides helpers:

- **HTTP**: `MakeRequest()`, `ParseResponse()`, `AssertStatusCode()`
- **Database**: `SetupTestDB()`, `TeardownTestDB()`, `SeedBasicData()`
- **Fixtures**: `CreateTestUser()`, `CreateTestTeam()`, `CreateTestProduct()`

## Best Practices

1. Each test should be independent
2. Use `defer ts.Teardown(t)` for cleanup
3. Create minimal test data
4. Test both success and error cases
5. Use descriptive test names
6. Follow table-driven test pattern
7. Check response structure, not just status codes

## CI/CD Integration

These tests can run in CI pipelines:

```yaml
# Example GitHub Actions
- name: Run tests
  run: make test-integration-coverage

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    files: ./apps/backend/coverage.out
```

## Documentation

- **Full Integration Testing Guide**: `../INTEGRATION_TESTING.md`
- **Quick Summary**: `../INTEGRATION_TEST_SUMMARY.md`
- **Examples**: Check `integration/auth_test.go` for reference

## Support

For questions or issues with tests:
1. Check the documentation files
2. Review existing test examples
3. Ensure test database is running (`make test-setup`)
