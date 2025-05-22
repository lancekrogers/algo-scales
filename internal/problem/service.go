package problem

import (
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// Service provides operations on problems
type Service struct {
	repository interfaces.ProblemRepository
}

// NewService creates a new problem service
func NewService() *Service {
	return &Service{
		repository: NewRepository(),
	}
}

// WithRepository sets a custom repository for the service
func (s *Service) WithRepository(repo interfaces.ProblemRepository) *Service {
	s.repository = repo
	return s
}

// ListAll returns all available problems
func (s *Service) ListAll() ([]Problem, error) {
	interfaceProblems, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	
	// Convert to local types
	problems := make([]Problem, len(interfaceProblems))
	for i, p := range interfaceProblems {
		problems[i] = s.convertFromInterface(p)
	}
	return problems, nil
}

// GetByID retrieves a problem by its ID
func (s *Service) GetByID(id string) (*Problem, error) {
	interfaceProblem, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	localProblem := s.convertFromInterface(*interfaceProblem)
	return &localProblem, nil
}

// GetByPattern returns problems matching a specific pattern
func (s *Service) GetByPattern(pattern string) ([]Problem, error) {
	interfaceProblems, err := s.repository.GetByPattern(pattern)
	if err != nil {
		return nil, err
	}
	
	// Convert to local types
	problems := make([]Problem, len(interfaceProblems))
	for i, p := range interfaceProblems {
		problems[i] = s.convertFromInterface(p)
	}
	return problems, nil
}

// GetPatterns returns all available algorithm patterns
func (s *Service) GetPatterns() ([]string, error) {
	return s.repository.GetPatterns()
}

// GetLanguages returns all available programming languages
func (s *Service) GetLanguages() ([]string, error) {
	return s.repository.GetLanguages()
}

// TestSolution tests a user's solution against test cases
func (s *Service) TestSolution(problemID, code, language string) ([]struct {
	Input    string
	Expected string
	Actual   string
	Passed   bool
}, error) {
	// Get the problem
	p, err := s.repository.GetByID(problemID)
	if err != nil {
		return nil, err
	}

	// For now, we'll simulate testing by returning mock results
	results := make([]struct {
		Input    string
		Expected string
		Actual   string
		Passed   bool
	}, 0, len(p.TestCases))

	// Generate simulated test results
	for _, tc := range p.TestCases {
		// For demonstration purposes, we'll simulate most tests passing
		// In a real implementation, we would execute the code against test cases
		passed := true
		
		// Convert interface{} to string with type assertion
		inputStr := ""
		if str, ok := tc.Input.(string); ok {
			inputStr = str
		}
		
		expectedStr := ""
		if str, ok := tc.Expected.(string); ok {
			expectedStr = str
		}
		
		actual := expectedStr // For simulation

		results = append(results, struct {
			Input    string
			Expected string
			Actual   string
			Passed   bool
		}{
			Input:    inputStr,
			Expected: expectedStr,
			Actual:   actual,
			Passed:   passed,
		})
	}

	return results, nil
}

// convertFromInterface converts an interfaces.Problem to a local Problem
func (s *Service) convertFromInterface(p interfaces.Problem) Problem {
	// Convert test cases
	testCases := make([]TestCase, len(p.TestCases))
	for i, tc := range p.TestCases {
		inputStr := ""
		if str, ok := tc.Input.(string); ok {
			inputStr = str
		}
		
		expectedStr := ""
		if str, ok := tc.Expected.(string); ok {
			expectedStr = str
		}
		
		testCases[i] = TestCase{
			Input:    inputStr,
			Expected: expectedStr,
		}
	}
	
	// Create starter code map
	starterCode := make(map[string]string)
	for _, lang := range p.Languages {
		starterCode[lang] = "" // Empty starter code for now
	}
	
	return Problem{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Difficulty:  p.Difficulty,
		Patterns:    p.Tags, // Use tags as patterns
		Companies:   p.Companies,
		TestCases:   testCases,
		StarterCode: starterCode,
		Solutions:   make(map[string]string),
	}
}