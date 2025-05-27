# CLI Testing Guide

This guide explains how to test CLI commands in AlgoScales to prevent runtime errors and ensure commands work as expected.

## Test Categories

### 1. Unit Tests
Test individual command functions with mocked dependencies.

**Location**: `cmd/*_test.go`

**Example**:
```go
func TestStartCramCommand(t *testing.T) {
    // Mock session.Start
    restore := mockSessionStart(nil)
    defer restore()
    
    output, err := executeCommand(rootCmd, "start", "cram")
    assert.NoError(t, err)
    assert.NotContains(t, output, "Error starting session")
}
```

### 2. Integration Tests
Test end-to-end command execution without launching UI.

**Location**: `cmd/integration_test.go`

**Features**:
- Tests commands don't panic
- Verifies session creation
- Checks flag handling

### 3. Manual Testing
Test actual command execution with UI.

## Testing Environment Setup

The testing framework uses environment variables to control behavior:

- `TESTING=1` - Skips UI launches and setup processes
- `DEBUG=1` - Enables debug logging

## Key Testing Functions

### `executeCommand(root *cobra.Command, args ...string)`
- Captures command output
- Sets `TESTING=1` automatically
- Returns output and error

### `mockSessionStart(err error)`
- Mocks `session.Start` function
- Returns restore function
- Allows testing without actual session creation

## Manual Testing Checklist

Before committing changes to CLI commands, test these scenarios:

### Basic Commands
- [ ] `./bin/algo-scales` (help)
- [ ] `./bin/algo-scales --help`
- [ ] `./bin/algo-scales version`

### Start Commands
- [ ] `./bin/algo-scales start cram`
- [ ] `./bin/algo-scales start learn`
- [ ] `./bin/algo-scales start practice`
- [ ] `./bin/algo-scales start cram --language python`
- [ ] `./bin/algo-scales start learn --tui`
- [ ] `./bin/algo-scales start practice --split`

### CLI Commands
- [ ] `./bin/algo-scales solve`
- [ ] `./bin/algo-scales list`
- [ ] `./bin/algo-scales stats`

### Error Conditions
- [ ] Invalid flags
- [ ] Non-existent problems
- [ ] Network issues (if applicable)

## Common Issues and Solutions

### 1. Silent Command Failure
**Problem**: Command executes but nothing happens
**Cause**: Missing UI launch after session creation
**Solution**: Ensure `launchUI()` is called after `session.Start()`

### 2. Panic on Nil Pointer
**Problem**: Runtime panic due to uninitialized fields
**Cause**: Improper struct initialization
**Solution**: Use constructor functions with proper field initialization

### 3. Flag Parsing Issues
**Problem**: Flags not recognized or parsed incorrectly
**Cause**: Flag not properly registered or inherited
**Solution**: Check flag registration in `init()` functions

## Test Automation

### Running Tests
```bash
# Run all CLI tests
go test ./cmd -v

# Run specific test
go test ./cmd -run TestStartCommand -v

# Run integration tests
go test ./cmd -run TestIntegration -v
```

### Continuous Integration
Tests are run automatically on:
- Pull requests
- Main branch commits
- Release builds

## Best Practices

1. **Always test manually** after making CLI changes
2. **Mock external dependencies** in unit tests
3. **Use integration tests** for command flow verification
4. **Test error conditions** as well as success cases
5. **Document breaking changes** in CLI behavior

## Adding New Commands

When adding a new command:

1. Create the command file (`cmd/newcommand.go`)
2. Add unit tests (`cmd/newcommand_test.go`)
3. Add integration tests (update `cmd/integration_test.go`)
4. Update this guide with new test cases
5. Test manually before committing

## Debugging Tips

### Enable Debug Mode
```bash
./bin/algo-scales --debug start cram
```

### Check Logs
```bash
tail -f ~/.algo-scales/logs/algo-scales.log
```

### Verbose Testing
```bash
go test ./cmd -v -run TestSpecificFunction
```