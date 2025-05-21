package registry

import (
	"fmt"
	"sync"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/session/execution"
)

// ServiceRegistry manages all application services
type ServiceRegistry struct {
	problemRepo     interfaces.ProblemRepository
	fileSystem      interfaces.FileSystem
	testRunnerReg   *execution.RunnerRegistry
	statsService    interfaces.StatsService
	templateService interfaces.TemplateService
	mutex           sync.RWMutex
}

// NewServiceRegistry creates a new service registry with default implementations
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		testRunnerReg: execution.NewRunnerRegistry(),
	}
}

// WithProblemRepository sets the problem repository
func (r *ServiceRegistry) WithProblemRepository(repo interfaces.ProblemRepository) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.problemRepo = repo
	return r
}

// WithFileSystem sets the file system
func (r *ServiceRegistry) WithFileSystem(fs interfaces.FileSystem) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.fileSystem = fs
	return r
}

// WithTestRunnerRegistry sets the test runner registry
func (r *ServiceRegistry) WithTestRunnerRegistry(reg *execution.RunnerRegistry) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.testRunnerReg = reg
	return r
}

// WithStatsService sets the stats service
func (r *ServiceRegistry) WithStatsService(service interfaces.StatsService) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.statsService = service
	return r
}

// WithTemplateService sets the template service
func (r *ServiceRegistry) WithTemplateService(service interfaces.TemplateService) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.templateService = service
	return r
}

// GetProblemRepository returns the problem repository
func (r *ServiceRegistry) GetProblemRepository() (interfaces.ProblemRepository, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.problemRepo == nil {
		return nil, fmt.Errorf("problem repository not configured")
	}
	return r.problemRepo, nil
}

// GetFileSystem returns the file system
func (r *ServiceRegistry) GetFileSystem() (interfaces.FileSystem, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.fileSystem == nil {
		return nil, fmt.Errorf("file system not configured")
	}
	return r.fileSystem, nil
}

// GetTestRunnerRegistry returns the test runner registry
func (r *ServiceRegistry) GetTestRunnerRegistry() *execution.RunnerRegistry {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.testRunnerReg
}

// GetStatsService returns the stats service
func (r *ServiceRegistry) GetStatsService() (interfaces.StatsService, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.statsService == nil {
		return nil, fmt.Errorf("stats service not configured")
	}
	return r.statsService, nil
}

// GetTemplateService returns the template service
func (r *ServiceRegistry) GetTemplateService() (interfaces.TemplateService, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.templateService == nil {
		return nil, fmt.Errorf("template service not configured")
	}
	return r.templateService, nil
}

// DefaultRegistry is the default service registry instance
var DefaultRegistry = NewServiceRegistry()

// InitializeDefaults sets up the default registry with concrete implementations
// This function should be called from the main package to avoid import cycles
func InitializeDefaults(
	fs interfaces.FileSystem,
	repo interfaces.ProblemRepository,
	stats interfaces.StatsService,
	templates interfaces.TemplateService,
) {
	DefaultRegistry.WithFileSystem(fs)
	DefaultRegistry.WithProblemRepository(repo)
	DefaultRegistry.WithStatsService(stats)
	DefaultRegistry.WithTemplateService(templates)
}