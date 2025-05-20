// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import (
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// TestRunner defines an interface for running code tests
type TestRunner interface {
	// ExecuteTests runs tests for a solution and returns the results
	ExecuteTests(prob *problem.Problem, code string, timeout time.Duration) ([]TestResult, bool, error)
	
	// GetLanguage returns the language this runner supports
	GetLanguage() string
	
	// GenerateTestCode creates test code for a given problem
	GenerateTestCode(prob *problem.Problem, solutionCode string) (string, error)
}

// TestRunnerRegistry provides access to language-specific test runners
type TestRunnerRegistry interface {
	// GetRunner returns a test runner for the specified language
	GetRunner(language string) (TestRunner, error)
	
	// RegisterRunner adds a test runner to the registry
	RegisterRunner(runner TestRunner) error
	
	// GetSupportedLanguages returns a list of supported languages
	GetSupportedLanguages() []string
}