package execution

import (
	"context"
	"fmt"
	"sync"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// RunnerRegistry implements the TestRunnerRegistry interface
type RunnerRegistry struct {
	runners map[string]interfaces.TestRunner
	mutex   sync.RWMutex
}

// NewRunnerRegistry creates a new test runner registry with default runners
func NewRunnerRegistry() *RunnerRegistry {
	registry := &RunnerRegistry{
		runners: make(map[string]interfaces.TestRunner),
	}
	
	// Register default runners
	registry.RegisterRunner(NewGoTestRunner())
	registry.RegisterRunner(NewPythonTestRunner())
	registry.RegisterRunner(NewJavaScriptTestRunner())
	
	return registry
}

// GetRunner returns a test runner for the specified language
func (r *RunnerRegistry) GetRunner(language string) (interfaces.TestRunner, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if runner, ok := r.runners[language]; ok {
		return runner, nil
	}
	
	return nil, fmt.Errorf("no test runner available for language: %s", language)
}

// RegisterRunner adds a test runner to the registry
func (r *RunnerRegistry) RegisterRunner(runner interfaces.TestRunner) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	language := runner.GetLanguage()
	if language == "" {
		return fmt.Errorf("test runner must specify a language")
	}
	
	r.runners[language] = runner
	return nil
}

// GetSupportedLanguages returns a list of supported languages
func (r *RunnerRegistry) GetSupportedLanguages() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	languages := make([]string, 0, len(r.runners))
	for lang := range r.runners {
		languages = append(languages, lang)
	}
	
	return languages
}

// DefaultRegistry is the default test runner registry instance
var DefaultRegistry = NewRunnerRegistry()

// ExecuteTests is a convenience function using the default registry
func ExecuteTests(ctx context.Context, prob *interfaces.Problem, code, language string, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
	runner, err := DefaultRegistry.GetRunner(language)
	if err != nil {
		return nil, false, err
	}
	
	return runner.ExecuteTests(ctx, prob, code, timeout)
}