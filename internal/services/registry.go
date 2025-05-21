package services

import (
	"sync"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// ServiceRegistry holds all business logic services for the application
type ServiceRegistry struct {
	problemService     ProblemService
	sessionService     SessionService
	statsService       StatsCommandService
	mutex              sync.RWMutex
}

// Global default service registry
var DefaultRegistry = NewServiceRegistry()

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{}
}

// WithProblemService sets the problem service
func (r *ServiceRegistry) WithProblemService(service ProblemService) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.problemService = service
	return r
}

// WithSessionService sets the session service
func (r *ServiceRegistry) WithSessionService(service SessionService) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.sessionService = service
	return r
}

// WithStatsService sets the stats service
func (r *ServiceRegistry) WithStatsService(service StatsCommandService) *ServiceRegistry {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.statsService = service
	return r
}

// GetProblemService returns the problem service, creating a default if none exists
func (r *ServiceRegistry) GetProblemService() ProblemService {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.problemService == nil {
		// Create default service with nil repository (will cause fallback behavior)
		r.problemService = NewProblemService(nil)
	}
	
	return r.problemService
}

// GetSessionService returns the session service, creating a default if none exists
func (r *ServiceRegistry) GetSessionService() SessionService {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.sessionService == nil {
		// Create default service
		r.sessionService = NewSessionService(r.GetProblemService(), nil)
	}
	
	return r.sessionService
}

// GetStatsService returns the stats service, creating a default if none exists
func (r *ServiceRegistry) GetStatsService() StatsCommandService {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.statsService == nil {
		// Create default legacy stats service
		r.statsService = NewStatsCommandService(nil)
	}
	
	return r.statsService
}

// InitializeDefaults initializes the default registry with concrete implementations
func InitializeDefaults(
	problemRepo interfaces.ProblemRepository,
	sessionManager interfaces.SessionManager,
	statsService interfaces.StatsService,
) {
	DefaultRegistry.WithProblemService(NewProblemService(problemRepo))
	DefaultRegistry.WithSessionService(NewSessionService(DefaultRegistry.GetProblemService(), sessionManager));
	DefaultRegistry.WithStatsService(NewStatsCommandService(statsService));
}