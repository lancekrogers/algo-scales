package problem

import (
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/utils"
)

// MockRepository provides a test double for the ProblemRepository interface
type MockRepository struct {
	Problems         []Problem
	ProblemsByID     map[string]*Problem
	AvailablePatterns []string
	Languages        []string
	fs              interfaces.FileSystem
}

// NewMockRepository creates a new mock problem repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		Problems:     []Problem{},
		ProblemsByID: make(map[string]*Problem),
		AvailablePatterns: []string{},
		Languages:    []string{},
		fs:          utils.NewMockFileSystem(),
	}
}

// WithFileSystem sets a custom filesystem for the mock repository
func (m *MockRepository) WithFileSystem(fs interfaces.FileSystem) *MockRepository {
	m.fs = fs
	return m
}

// AddProblem adds a problem to the mock repository
func (m *MockRepository) AddProblem(p Problem) *MockRepository {
	m.Problems = append(m.Problems, p)
	m.ProblemsByID[p.ID] = &p
	return m
}

// SetPatterns sets the available patterns
func (m *MockRepository) SetPatterns(patterns []string) *MockRepository {
	m.AvailablePatterns = patterns
	return m
}

// SetLanguages sets the available languages
func (m *MockRepository) SetLanguages(languages []string) *MockRepository {
	m.Languages = languages
	return m
}

// Ensure MockRepository implements ProblemRepository
var _ interfaces.ProblemRepository = (*MockRepository)(nil)

// GetAll returns all available problems
func (m *MockRepository) GetAll() ([]interfaces.Problem, error) {
	result := make([]interfaces.Problem, len(m.Problems))
	for i, p := range m.Problems {
		result[i] = m.convertToInterface(p)
	}
	return result, nil
}

// GetByID retrieves a specific problem by its ID
func (m *MockRepository) GetByID(id string) (*interfaces.Problem, error) {
	if problem, ok := m.ProblemsByID[id]; ok {
		converted := m.convertToInterface(*problem)
		return &converted, nil
	}
	return nil, ErrProblemNotFound
}

// GetByPattern returns problems matching a specific pattern
func (m *MockRepository) GetByPattern(pattern string) ([]interfaces.Problem, error) {
	var filtered []Problem
	if pattern == "" {
		filtered = m.Problems
	} else {
		for _, p := range m.Problems {
			for _, patternName := range p.Patterns {
				if patternName == pattern {
					filtered = append(filtered, p)
					break
				}
			}
		}
	}
	
	result := make([]interfaces.Problem, len(filtered))
	for i, p := range filtered {
		result[i] = m.convertToInterface(p)
	}
	return result, nil
}

// GetByDifficulty returns problems with a specific difficulty level
func (m *MockRepository) GetByDifficulty(difficulty string) ([]interfaces.Problem, error) {
	var filtered []Problem
	for _, p := range m.Problems {
		if p.Difficulty == difficulty {
			filtered = append(filtered, p)
		}
	}
	
	result := make([]interfaces.Problem, len(filtered))
	for i, p := range filtered {
		result[i] = m.convertToInterface(p)
	}
	return result, nil
}

// GetByTags returns problems matching any of the specified tags
func (m *MockRepository) GetByTags(tags []string) ([]interfaces.Problem, error) {
	var filtered []Problem
	for _, p := range m.Problems {
		// Check if any pattern matches any tag
		for _, pattern := range p.Patterns {
			for _, tag := range tags {
				if pattern == tag {
					filtered = append(filtered, p)
					goto next_problem
				}
			}
		}
		next_problem:
	}
	
	result := make([]interfaces.Problem, len(filtered))
	for i, p := range filtered {
		result[i] = m.convertToInterface(p)
	}
	return result, nil
}

// GetRandom returns a random problem
func (m *MockRepository) GetRandom() (*interfaces.Problem, error) {
	if len(m.Problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple selection for mock
	randomIndex := len(m.Problems) / 2
	converted := m.convertToInterface(m.Problems[randomIndex])
	return &converted, nil
}

// GetRandomByPattern returns a random problem matching a pattern
func (m *MockRepository) GetRandomByPattern(pattern string) (*interfaces.Problem, error) {
	problems, err := m.GetByPattern(pattern)
	if err != nil {
		return nil, err
	}
	
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple selection for mock
	randomIndex := len(problems) / 2
	return &problems[randomIndex], nil
}

// GetRandomByDifficulty returns a random problem with a difficulty
func (m *MockRepository) GetRandomByDifficulty(difficulty string) (*interfaces.Problem, error) {
	problems, err := m.GetByDifficulty(difficulty)
	if err != nil {
		return nil, err
	}
	
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple selection for mock
	randomIndex := len(problems) / 2
	return &problems[randomIndex], nil
}

// GetRandomByTags returns a random problem matching tags
func (m *MockRepository) GetRandomByTags(tags []string) (*interfaces.Problem, error) {
	problems, err := m.GetByTags(tags)
	if err != nil {
		return nil, err
	}
	
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple selection for mock
	randomIndex := len(problems) / 2
	return &problems[randomIndex], nil
}

// GetPatterns returns all available algorithm patterns
func (m *MockRepository) GetPatterns() ([]string, error) {
	return m.AvailablePatterns, nil
}

// GetLanguages returns all available programming languages
func (m *MockRepository) GetLanguages() ([]string, error) {
	return m.Languages, nil
}

// GetByCompany returns problems from a specific company
func (m *MockRepository) GetByCompany(company string) ([]Problem, error) {
	var filtered []Problem
	for _, p := range m.Problems {
		for _, c := range p.Companies {
			if c == company {
				filtered = append(filtered, p)
				break
			}
		}
	}
	
	return filtered, nil
}

// convertToInterface converts a local Problem to interfaces.Problem
func (m *MockRepository) convertToInterface(p Problem) interfaces.Problem {
	var pattern string
	if len(p.Patterns) > 0 {
		pattern = p.Patterns[0] // Use first pattern for simplicity
	}
	
	var languages []string
	for lang := range p.StarterCode {
		languages = append(languages, lang)
	}
	
	// Convert test cases
	testCases := make([]interfaces.TestCase, len(p.TestCases))
	for i, tc := range p.TestCases {
		testCases[i] = interfaces.TestCase{
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}
	
	return interfaces.Problem{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Pattern:     pattern,
		Difficulty:  p.Difficulty,
		Companies:   p.Companies,
		Tags:        p.Patterns, // Use patterns as tags
		TestCases:   testCases,
		Languages:   languages,
	}
}