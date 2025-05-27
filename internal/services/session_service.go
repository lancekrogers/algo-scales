package services

import (
	"context"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/session"
)

// SessionService provides business logic for session management
type SessionService interface {
	// StartSession starts a new practice session
	StartSession(ctx context.Context, opts session.Options) error
	
	// StartSessionWithProblem starts a session with a specific problem
	StartSessionWithProblem(ctx context.Context, opts session.Options, problem *problem.Problem) error
}

// SessionServiceImpl implements SessionService
type SessionServiceImpl struct {
	problemService ProblemService
	sessionManager interfaces.SessionManager
}

// NewSessionService creates a new session service
func NewSessionService(problemService ProblemService, sessionManager interfaces.SessionManager) SessionService {
	return &SessionServiceImpl{
		problemService: problemService,
		sessionManager: sessionManager,
	}
}

// StartSession starts a new practice session
func (s *SessionServiceImpl) StartSession(ctx context.Context, opts session.Options) error {
	// Delegate to the session package for now to maintain compatibility
	return session.Start(opts)
}

// StartSessionWithProblem starts a session with a specific problem
func (s *SessionServiceImpl) StartSessionWithProblem(ctx context.Context, opts session.Options, problem *problem.Problem) error {
	// For now, delegate to session package while maintaining the interface
	// In the future, this would use the session manager to create proper sessions
	if problem != nil {
		opts.ProblemID = problem.ID
	}
	return session.Start(opts)
}