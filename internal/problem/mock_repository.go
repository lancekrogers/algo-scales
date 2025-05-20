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
func (m *MockRepository) GetAll() ([]Problem, error) {
	return m.Problems, nil
}

// GetByID retrieves a specific problem by its ID
func (m *MockRepository) GetByID(id string) (*Problem, error) {
	if problem, ok := m.ProblemsByID[id]; ok {
		return problem, nil
	}
	return nil, ErrProblemNotFound
}

// GetByPattern returns problems matching a specific pattern
func (m *MockRepository) GetByPattern(pattern string) ([]Problem, error) {
	if pattern == "" {
		return m.Problems, nil
	}
	
	var filtered []Problem
	for _, p := range m.Problems {
		for _, patternName := range p.Patterns {
			if patternName == pattern {
				filtered = append(filtered, p)
				break
			}
		}
	}
	
	return filtered, nil
}

// GetPatterns returns all available algorithm patterns
func (m *MockRepository) GetPatterns() ([]string, error) {
	return m.AvailablePatterns, nil
}

// GetLanguages returns all available programming languages
func (m *MockRepository) GetLanguages() ([]string, error) {
	return m.Languages, nil
}

// GetByDifficulty returns problems with a specific difficulty level
func (m *MockRepository) GetByDifficulty(difficulty string) ([]Problem, error) {
	var filtered []Problem
	for _, p := range m.Problems {
		if p.Difficulty == difficulty {
			filtered = append(filtered, p)
		}
	}
	
	return filtered, nil
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