# Vim Mode Testing Guide

## Overview

This document describes the testing infrastructure for AlgoScales vim mode and Neovim plugin integration.

## Test Structure

### 1. Unit Tests
- **Location**: `cmd/vim_commands_test.go`
- **Coverage**: Individual vim mode commands
- **Run**: `go test ./cmd -run TestVimCommands`

### 2. Integration Tests
- **Location**: `cmd/vim_integration_test.go`
- **Coverage**: Full CLI integration with vim mode
- **Features**:
  - Builds actual binary
  - Tests all vim mode commands
  - Verifies JSON output format
  - Tests multi-level hints
- **Run**: `go test ./cmd -run TestVimModeIntegration`

### 3. Manual Tests
- **Location**: `tests/manual_vim_test.sh`
- **Purpose**: Quick manual verification
- **Run**: `./tests/manual_vim_test.sh`

### 4. Neovim Plugin Tests
- **Location**: `algo-scales-nvim/tests/algo-scales_spec.lua`
- **Framework**: Plenary (Neovim testing framework)
- **Coverage**:
  - Plugin setup and configuration
  - Window layout creation
  - Hint level tracking
  - Workspace management
- **Run**: `./algo-scales-nvim/tests/test_runner.sh`

## Make Targets

```bash
# Run all tests including vim mode
make test-all

# Run only vim mode tests
make test-vim

# Run integration tests (includes vim mode)
make test-integration

# Run standard test suite
make test-chart
```

## Test Output Examples

### Successful Vim Mode Test
```json
{
  "hint": "The two-pointers technique involves...",
  "level": 1,
  "walkthrough": ["Step 1...", "Step 2..."],
  "solution": "func pairWithTargetSum(...) { ... }"
}
```

### Test Results
- ✅ 270+ tests passing
- ✅ Vim mode commands validated
- ✅ Multi-level hints working
- ✅ Integration with Neovim tested

## Continuous Integration

The vim mode tests are automatically run when:
1. Running `make test-all`
2. Running `make test-integration`
3. Running `make test-vim`

## Adding New Tests

### For CLI Commands
1. Add test cases to `vim_commands_test.go`
2. Add integration scenarios to `vim_integration_test.go`

### For Neovim Plugin
1. Add test specs to `algo-scales_spec.lua`
2. Test new features with mock vim.fn.system calls

## Known Issues

1. **Hint Levels**: CLI commands don't persist hint levels between invocations (by design)
2. **Submit Test**: One test case in two_sum has issues with certain inputs
3. **Neovim Tests**: Require plenary.nvim to be installed

## Future Improvements

1. Add GitHub Actions workflow for automated testing
2. Add performance benchmarks for vim mode
3. Create integration tests for AI hint functionality
4. Add end-to-end tests with actual Neovim instance