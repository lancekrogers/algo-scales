// Package interfaces defines the core interfaces for Algo Scales
package interfaces

// Problem represents an algorithm problem
type Problem struct {
	ID          string
	Title       string
	Description string
	Pattern     string
	Difficulty  string
	Companies   []string
	Tags        []string
	TestCases   []TestCase
	Languages   []string
}

// TestCase represents a problem test case
type TestCase struct {
	Input    interface{}
	Expected interface{}
}

// ProblemRepository defines the interface for accessing algorithm problems
type ProblemRepository interface {
	// GetAll returns all available problems
	GetAll() ([]Problem, error)
	
	// GetByID retrieves a specific problem by its ID
	GetByID(id string) (*Problem, error)
	
	// GetByPattern returns problems matching a specific pattern
	GetByPattern(pattern string) ([]Problem, error)
	
	// GetByDifficulty returns problems with a specific difficulty level
	GetByDifficulty(difficulty string) ([]Problem, error)
	
	// GetByTags returns problems matching any of the specified tags
	GetByTags(tags []string) ([]Problem, error)
	
	// GetRandom returns a random problem
	GetRandom() (*Problem, error)
	
	// GetRandomByPattern returns a random problem matching a pattern
	GetRandomByPattern(pattern string) (*Problem, error)
	
	// GetRandomByDifficulty returns a random problem with a difficulty
	GetRandomByDifficulty(difficulty string) (*Problem, error)
	
	// GetRandomByTags returns a random problem matching tags
	GetRandomByTags(tags []string) (*Problem, error)
}