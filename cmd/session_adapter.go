// Session adapter for command-line tools
package cmd

import (
	"context"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/session"
)

// SessionAdapter adapts a session.Session to implement the SessionImpl interface methods needed by CLI
type SessionAdapter struct {
	*session.Session
	Implementation interfaces.Session
}

// ensureImplementation creates a SessionImpl if it doesn't exist
func (s *SessionAdapter) ensureImplementation() {
	if s.Implementation == nil {
		// Create a properly initialized SessionImpl
		impl := session.NewSessionImpl(convertOptions(s.Options), s.Problem)
		// Set additional fields that the constructor doesn't handle
		impl.Workspace = s.Workspace
		impl.CodeFile = s.CodeFile
		impl.ShowPattern = s.ShowPattern
		impl.StartTime = s.StartTime
		s.Implementation = impl
	}
}

// SetCode implements the SetCode method for CLI usage
func (s *SessionAdapter) SetCode(code string) error {
	s.ensureImplementation()
	return s.Implementation.SetCode(code)
}

// RunTests implements the RunTests method for CLI usage
func (s *SessionAdapter) RunTests(ctx context.Context) ([]interfaces.TestResult, bool, error) {
	s.ensureImplementation()
	return s.Implementation.RunTests(ctx)
}

// ShowHints implements the ShowHints method for CLI usage
func (s *SessionAdapter) ShowHints(show bool) {
	s.ShowPattern = show
}

// ShowSolution implements the ShowSolution method for CLI usage
func (s *SessionAdapter) ShowSolution(show bool) {
	// Nothing to do here for CLI mode
}

// FinishSession implements the session finish method
func (s *SessionAdapter) FinishSession(solved bool) error {
	return s.Session.FinishSession(solved)
}

// Helper function to convert between session option types
func convertOptions(opts session.Options) interfaces.SessionOptions {
	return interfaces.SessionOptions{
		Mode:       interfaces.SessionMode(opts.Mode),
		Language:   opts.Language,
		Timer:      opts.Timer,
		Pattern:    opts.Pattern,
		Difficulty: opts.Difficulty,
		ProblemID:  opts.ProblemID,
	}
}