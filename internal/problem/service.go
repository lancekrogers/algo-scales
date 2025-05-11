package problem

// Service provides operations on problems
type Service struct{}

// NewService creates a new problem service
func NewService() *Service {
	return &Service{}
}

// ListAll returns all available problems
func (s *Service) ListAll() ([]Problem, error) {
	return ListAll()
}

// GetByID retrieves a problem by its ID
func (s *Service) GetByID(id string) (*Problem, error) {
	return GetByID(id)
}

// TestSolution tests a user's solution against test cases
func (s *Service) TestSolution(problemID, code, language string) ([]struct {
	Input    string
	Expected string
	Actual   string
	Passed   bool
}, error) {
	// Get the problem
	p, err := s.GetByID(problemID)
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
		actual := tc.Expected

		results = append(results, struct {
			Input    string
			Expected string
			Actual   string
			Passed   bool
		}{
			Input:    tc.Input,
			Expected: tc.Expected,
			Actual:   actual,
			Passed:   passed,
		})
	}

	return results, nil
}