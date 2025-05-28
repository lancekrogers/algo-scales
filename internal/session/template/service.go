// Package template provides code template generation functionality
package template

import (
	"fmt"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// Service implements the TemplateService interface
type Service struct {
	generators map[string]LanguageGenerator
}

// LanguageGenerator generates code templates for a specific language
type LanguageGenerator interface {
	// GetTemplate returns a code template for a problem
	GetTemplate(prob interfaces.Problem) string
	
	// GetTestHarness generates a test harness
	GetTestHarness(prob interfaces.Problem, solutionCode string) string
	
	// GetLanguage returns the language this generator supports
	GetLanguage() string
	
	// GetFunctionName extracts the primary function name from the code
	GetFunctionName(code string) string
}

// NewService creates a new template service
func NewService() *Service {
	service := &Service{
		generators: make(map[string]LanguageGenerator),
	}
	
	// Register default generators
	service.RegisterGenerator(NewGoGenerator())
	service.RegisterGenerator(NewPythonGenerator())
	service.RegisterGenerator(NewJavaScriptGenerator())
	
	return service
}

// RegisterGenerator adds a language generator to the service
func (s *Service) RegisterGenerator(generator LanguageGenerator) {
	s.generators[generator.GetLanguage()] = generator
}

// GetTemplate returns a code template for a given problem and language
func (s *Service) GetTemplate(prob *interfaces.Problem, language string) (string, error) {
	generator, ok := s.generators[language]
	if !ok {
		// Fallback to a generic template
		return s.getGenericTemplate(prob), nil
	}
	
	return generator.GetTemplate(*prob), nil
}

// GetTestHarness generates a test harness for a language
func (s *Service) GetTestHarness(prob *interfaces.Problem, solutionCode, language string) (string, error) {
	generator, ok := s.generators[language]
	if !ok {
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	
	return generator.GetTestHarness(*prob, solutionCode), nil
}

// GetSupportedLanguages returns a list of supported languages
func (s *Service) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(s.generators))
	for lang := range s.generators {
		languages = append(languages, lang)
	}
	return languages
}

// GetFunctionName extracts the primary function name from the code
func (s *Service) GetFunctionName(code, language string) (string, error) {
	generator, ok := s.generators[language]
	if !ok {
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	
	return generator.GetFunctionName(code), nil
}

// getGenericTemplate returns a generic template for any language
func (s *Service) getGenericTemplate(prob *interfaces.Problem) string {
	return fmt.Sprintf(`// %s
// %s

/*
 * Step 1: Understand the problem
 * - Read the problem description carefully
 * - Identify input/output requirements
 * - Consider edge cases
 * 
 * Step 2: Plan your approach
 * - What algorithm pattern applies here?
 * - What data structures do you need?
 * - What's the time/space complexity?
 * 
 * Step 3: Implement your solution
 * - Replace this with your actual implementation
 */
`, prob.Title, prob.Description)
}