package template

import (
	"fmt"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// MockService provides a mock implementation of TemplateService for testing
type MockService struct {
	Templates          map[string]string
	TestHarnesses      map[string]string
	SupportedLanguages []string
	FunctionNames      map[string]string
}

// NewMockService creates a new mock template service
func NewMockService() *MockService {
	return &MockService{
		Templates:          make(map[string]string),
		TestHarnesses:      make(map[string]string),
		SupportedLanguages: []string{"go", "python", "javascript"},
		FunctionNames:      make(map[string]string),
	}
}

// GetTemplate returns a code template for a given problem and language
func (s *MockService) GetTemplate(prob *interfaces.Problem, language string) (string, error) {
	key := fmt.Sprintf("%s:%s", prob.ID, language)
	if template, ok := s.Templates[key]; ok {
		return template, nil
	}
	
	// Return a basic template if not found
	return fmt.Sprintf("Mock template for %s in %s", prob.ID, language), nil
}

// GetTestHarness generates a test harness for a language
func (s *MockService) GetTestHarness(prob *interfaces.Problem, solutionCode, language string) (string, error) {
	key := fmt.Sprintf("%s:%s", prob.ID, language)
	if harness, ok := s.TestHarnesses[key]; ok {
		return harness, nil
	}
	
	// Return a basic test harness if not found
	return fmt.Sprintf("Mock test harness for %s in %s", prob.ID, language), nil
}

// GetSupportedLanguages returns a list of supported languages
func (s *MockService) GetSupportedLanguages() []string {
	return s.SupportedLanguages
}

// GetFunctionName extracts the primary function name from the code
func (s *MockService) GetFunctionName(code, language string) (string, error) {
	key := fmt.Sprintf("%s:%s", code[:min(10, len(code))], language)
	if name, ok := s.FunctionNames[key]; ok {
		return name, nil
	}
	
	// Return a default function name if not found
	return "solution", nil
}

// SetTemplate sets a template for a problem and language
func (s *MockService) SetTemplate(probID, language, template string) *MockService {
	key := fmt.Sprintf("%s:%s", probID, language)
	s.Templates[key] = template
	return s
}

// SetTestHarness sets a test harness for a problem and language
func (s *MockService) SetTestHarness(probID, language, harness string) *MockService {
	key := fmt.Sprintf("%s:%s", probID, language)
	s.TestHarnesses[key] = harness
	return s
}

// SetFunctionName sets a function name for a code snippet and language
func (s *MockService) SetFunctionName(code, language, name string) *MockService {
	key := fmt.Sprintf("%s:%s", code[:min(10, len(code))], language)
	s.FunctionNames[key] = name
	return s
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}