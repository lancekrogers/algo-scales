# Test Fixes for Algo Scales

This document describes how test failures were fixed in the Algo Scales project.

## Import Cycle Resolution

One of the main issues was a circular dependency between the `session` and `ui` packages:

```
session → ui → session
```

This was fixed by:

1. Creating a proper dependency direction:
   ```
   ui → session → problem
   ```

2. Creating a common `interfaces` package for shared types:
   ```
   session implements interfaces.SessionManager
   ui depends on interfaces.SessionManager, not directly on session
   ```

3. Removing the direct UI reference from the session package:
   ```diff
   - return ui.StartSession(session)
   + return nil
   ```

## Mockable Functions for Testing

Many tests were failing because they attempted to mock functions that were not designed to be mockable. We fixed this by converting private functions to function variables:

```diff
- func getConfigDir() string {
+ var getConfigDir = func() string {
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".algo-scales")
  }
```

This allows tests to temporarily replace these functions:

```go
origGetConfigDir := getConfigDir
defer func() { getConfigDir = origGetConfigDir }()

getConfigDir = func() string {
    return "/tmp/test-config"
}
```

We applied this pattern to several functions:

- `getConfigDir` in problem, license, api, and stats packages
- `verifySignature` in license package
- `ValidateLicense` in license package

## Architecture Improvements

1. **Separation of Concerns**: UI logic is now separate from core session logic

2. **Clear Interface Boundaries**: Using interfaces allows for better testing and flexibility

3. **Common Packages**: Shared functionality moved to common packages:
   - highlight
   - interfaces
   - utils

4. **MVC Pattern**: UI follows a Model-View-Controller pattern with clear separation:
   - Model: Data structures and state
   - View: Rendering and display
   - Controller: User input and logic

## Testing Approach

For effective testing:

1. Make functions mockable using function variables
2. Use interfaces to decouple components
3. Use dependency injection for better testing
4. Ensure all test cases cover both success and failure paths
5. Organize tests to match the structure of the code being tested

These fixes ensure a more testable and maintainable codebase, with proper separation of concerns.