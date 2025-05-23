package config

import (
	"time"
)

// TimeoutConfig contains configurable timeout settings
type TimeoutConfig struct {
	// TestExecution is the default timeout for test execution
	TestExecution time.Duration
	
	// FileOperations is the default timeout for file I/O operations
	FileOperations time.Duration
	
	// TemplateGeneration is the default timeout for template generation
	TemplateGeneration time.Duration
	
	// ProcessExecution is the default timeout for external process execution
	ProcessExecution time.Duration
	
	// DirectoryOperations is the default timeout for directory operations
	DirectoryOperations time.Duration
}

// DefaultTimeouts returns the default timeout configuration
func DefaultTimeouts() TimeoutConfig {
	return TimeoutConfig{
		TestExecution:       30 * time.Second,
		FileOperations:      5 * time.Second,
		TemplateGeneration:  10 * time.Second,
		ProcessExecution:    30 * time.Second,
		DirectoryOperations: 5 * time.Second,
	}
}

// DevelopmentTimeouts returns more generous timeouts for development
func DevelopmentTimeouts() TimeoutConfig {
	return TimeoutConfig{
		TestExecution:       60 * time.Second,
		FileOperations:      10 * time.Second,
		TemplateGeneration:  20 * time.Second,
		ProcessExecution:    60 * time.Second,
		DirectoryOperations: 10 * time.Second,
	}
}

// ProductionTimeouts returns stricter timeouts for production
func ProductionTimeouts() TimeoutConfig {
	return TimeoutConfig{
		TestExecution:       20 * time.Second,
		FileOperations:      3 * time.Second,
		TemplateGeneration:  5 * time.Second,
		ProcessExecution:    20 * time.Second,
		DirectoryOperations: 3 * time.Second,
	}
}

// Global default instance
var DefaultTimeoutConfig = DefaultTimeouts()