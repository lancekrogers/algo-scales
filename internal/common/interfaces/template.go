// Package interfaces defines the core interfaces for Algo Scales
package interfaces


// TemplateService defines the interface for generating code templates
type TemplateService interface {
	// GetTemplate returns a code template for a given problem and language
	GetTemplate(prob *Problem, language string) (string, error)
	
	// GetTestHarness generates a test harness for a language
	GetTestHarness(prob *Problem, solutionCode, language string) (string, error)
	
	// GetSupportedLanguages returns a list of supported languages
	GetSupportedLanguages() []string
}