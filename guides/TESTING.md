# Testing Guide for Algo Scales

This document describes the testing approach for the Algo Scales project.

## Testing Philosophy

Algo Scales follows a comprehensive testing strategy that includes:

- **Unit Tests**: Testing individual components in isolation
- **Integration Tests**: Testing interactions between components
- **End-to-End Tests**: Testing the application as a whole

All new code should include corresponding tests, and we aim for high test coverage (>80%).

## Running Tests

### Prerequisites

- Go 1.24 or higher
- Make (optional, for using Makefile commands)

### Basic Test Commands

```bash
# Run all tests
go test ./...

# Run all tests with verbose output
go test -v ./...

# Run tests with coverage reporting
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Using Make
make test
make test-coverage
```

### Testing Specific Packages

```bash
# Test a specific package
go test ./internal/problem

# Test a specific function
go test ./internal/problem -run TestGetByID
```

## Test Structure

We follow standard Go testing practices:

1. Test files are named with the `_test.go` suffix
2. Tests are prefixed with `Test`
3. We use table-driven tests where appropriate
4. We use subtests for related test cases

Example:

```go
func TestFeature(t *testing.T) {
    t.Run("Scenario1", func(t *testing.T) {
        // Test code
    })

    t.Run("Scenario2", func(t *testing.T) {
        // Test code
    })
}
```

## Mocking

We use a simple approach to mocking:

1. Export the original function to a variable
2. Replace it with a mock implementation for testing
3. Restore the original function after the test

Example:

```go
// Save original function
original := somePackage.SomeFunction
defer func() { somePackage.SomeFunction = original }()

// Replace with mock
somePackage.SomeFunction = func(args) {
    return mockResult
}
```

## Test Fixtures

Test fixtures are stored in the `testdata` directory and loaded during tests as needed.

## Continuous Integration

Tests are automatically run on GitHub Actions for:

- Every push to the main branch
- Every pull request

See the `.github/workflows/build.yml` file for details.

## Test Coverage

We track test coverage using:

1. Local coverage reports
2. CI integration with Codecov

To view the current coverage report, run:

```bash
make test-coverage
```

## Best Practices

1. **Keep tests independent**: Each test should run in isolation
2. **Use descriptive test names**: Names should indicate what's being tested
3. **Use test helpers**: Extract common setup/teardown code
4. **Test edge cases**: Include boundary conditions and error cases
5. **Keep tests fast**: Slow tests discourage regular testing
6. **Clean up after tests**: Tests should not leave behind files or resources
