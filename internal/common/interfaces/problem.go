// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import "context"

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
	StarterCode map[string]string
}

// TestCase represents a problem test case
type TestCase struct {
	Input    string
	Expected string
}

// ProblemRepository defines the interface for accessing algorithm problems
type ProblemRepository interface {
	// GetAll returns all available problems
	GetAll(ctx context.Context) ([]Problem, error)
	
	// GetByID retrieves a specific problem by its ID
	GetByID(ctx context.Context, id string) (*Problem, error)
	
	// GetByPattern returns problems matching a specific pattern
	GetByPattern(ctx context.Context, pattern string) ([]Problem, error)
	
	// GetByDifficulty returns problems with a specific difficulty level
	GetByDifficulty(ctx context.Context, difficulty string) ([]Problem, error)
	
	// GetByTags returns problems matching any of the specified tags
	GetByTags(ctx context.Context, tags []string) ([]Problem, error)
	
	// GetRandom returns a random problem
	GetRandom(ctx context.Context) (*Problem, error)
	
	// GetRandomByPattern returns a random problem matching a pattern
	GetRandomByPattern(ctx context.Context, pattern string) (*Problem, error)
	
	// GetRandomByDifficulty returns a random problem with a difficulty
	GetRandomByDifficulty(ctx context.Context, difficulty string) (*Problem, error)
	
	// GetRandomByTags returns a random problem matching tags
	GetRandomByTags(ctx context.Context, tags []string) (*Problem, error)
	
	// GetPatterns returns all available algorithm patterns
	GetPatterns(ctx context.Context) ([]string, error)
	
	// GetLanguages returns all available programming languages
	GetLanguages(ctx context.Context) ([]string, error)
}