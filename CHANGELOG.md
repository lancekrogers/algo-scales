# Changelog

## Current Status

### Work in Progress
- Terminal UI (TUI) mode is under active development and not yet functional
- The `--tui` flag and split-screen mode are disabled pending completion
- CLI mode is the stable, recommended interface for all users

## CLI Flag Fix

### Fixed
- Fixed `--cli` flag handling to properly work with Cobra, enabling CLI mode without errors
- CLI commands can now be run directly with `algo-scales --cli <command>`

### Changed
- Improved flag handling by registering all mode flags with Cobra
- Moved UI mode selection logic from main.go to cmd/root.go for better organization
- Moved terminal detection functionality to cmd/root.go
- Simplified main.go to delegate flag handling to Cobra
- Improved help output by hiding aliased flags

### Added
- Better error handling for terminal detection
- Utility function for hiding flags in help output across all commands
- More descriptive flag documentation

## Benefits
- More consistent command-line interface
- Proper integration with Cobra command framework
- Better code organization and maintainability
- Improved user experience with cleaner help output
- Direct access to CLI commands with `--cli` flag