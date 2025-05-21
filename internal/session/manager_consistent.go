package session

import (
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/utils"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/session/execution"
)

// ManagerOption defines a configuration option for Manager
type ManagerOption func(*Manager)

// WithFileSystem sets a custom file system
func WithFileSystem(fs interfaces.FileSystem) ManagerOption {
	return func(m *Manager) {
		m.fs = fs
	}
}

// WithProblemRepository sets a custom problem repository
func WithProblemRepository(repo interfaces.ProblemRepository) ManagerOption {
	return func(m *Manager) {
		m.problemRepo = repo
	}
}

// WithTestRegistry sets a custom test runner registry
func WithTestRegistry(registry interfaces.TestRunnerRegistry) ManagerOption {
	return func(m *Manager) {
		m.testRegistry = registry
	}
}

// NewConsistentManager creates a new session manager with standardized constructor pattern
func NewConsistentManager(opts ...ManagerOption) *Manager {
	m := &Manager{
		sessions: make(map[string]interfaces.Session),
	}
	
	// Apply provided options
	for _, opt := range opts {
		opt(m)
	}
	
	// Set defaults for any unspecified dependencies
	if m.fs == nil {
		m.fs = utils.NewFileSystem()
	}
	if m.problemRepo == nil {
		m.problemRepo = problem.NewRepository()
	}
	if m.testRegistry == nil {
		m.testRegistry = execution.DefaultRegistry
	}
	
	return m
}

// NewManagerWithDefaults creates a manager with all default dependencies
func NewManagerWithDefaults() *Manager {
	return NewConsistentManager()
}

// NewManagerForTesting creates a manager suitable for testing with mock dependencies
func NewManagerForTesting(
	fs interfaces.FileSystem,
	repo interfaces.ProblemRepository,
	registry interfaces.TestRunnerRegistry,
) *Manager {
	return NewConsistentManager(
		WithFileSystem(fs),
		WithProblemRepository(repo),
		WithTestRegistry(registry),
	)
}