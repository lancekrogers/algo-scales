// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import (
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// ProblemRepository defines the interface for accessing algorithm problems
type ProblemRepository interface {
	// GetAll returns all available problems
	GetAll() ([]problem.Problem, error)
	
	// GetByID retrieves a specific problem by its ID
	GetByID(id string) (*problem.Problem, error)
	
	// GetByPattern returns problems matching a specific pattern
	GetByPattern(pattern string) ([]problem.Problem, error)
	
	// GetPatterns returns all available algorithm patterns
	GetPatterns() ([]string, error)
	
	// GetLanguages returns all available programming languages
	GetLanguages() ([]string, error)
	
	// GetByDifficulty returns problems with a specific difficulty level
	GetByDifficulty(difficulty string) ([]problem.Problem, error)
	
	// GetByCompany returns problems from a specific company
	GetByCompany(company string) ([]problem.Problem, error)
}

// TestRunner defines the interface for testing problem solutions
type TestRunner interface {
	// TestSolution tests a solution against test cases
	TestSolution(problemID, code, language string) ([]TestResult, error)
}