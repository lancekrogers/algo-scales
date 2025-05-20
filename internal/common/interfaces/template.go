// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import (
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// TemplateService defines the interface for generating code templates
type TemplateService interface {
	// GetTemplate returns a code template for a given problem and language
	GetTemplate(prob *problem.Problem, language string) (string, error)
	
	// GetTestHarness generates a test harness for a language
	GetTestHarness(prob *problem.Problem, solutionCode, language string) (string, error)
	
	// GetSupportedLanguages returns a list of supported languages
	GetSupportedLanguages() []string
}