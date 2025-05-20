// Session adapter for command-line tools
package cmd

import (
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/session"
)

// SessionAdapter adapts a session.Session to implement the SessionImpl interface methods needed by CLI
type SessionAdapter struct {
	*session.Session
	Implementation interfaces.Session
}

// SetCode implements the SetCode method for CLI usage
func (s *SessionAdapter) SetCode(code string) error {
	// Create an implementation if not already exists
	if s.Implementation == nil {
		// Create a simple SessionImpl adapter
		s.Implementation = &session.SessionImpl{
			Problem:  s.Problem,
			Options:  convertOptions(s.Options),
			StartTime: s.StartTime,
			Workspace: s.Workspace,
			CodeFile:  s.CodeFile,
			ShowPattern: s.ShowPattern,
		}
	}
	
	return s.Implementation.SetCode(code)
}

// RunTests implements the RunTests method for CLI usage
func (s *SessionAdapter) RunTests() ([]interfaces.TestResult, bool, error) {
	// Create an implementation if not already exists
	if s.Implementation == nil {
		// Create a simple SessionImpl adapter
		s.Implementation = &session.SessionImpl{
			Problem:  s.Problem,
			Options:  convertOptions(s.Options),
			StartTime: s.StartTime,
			Workspace: s.Workspace,
			CodeFile:  s.CodeFile,
			ShowPattern: s.ShowPattern,
		}
	}
	
	return s.Implementation.RunTests()
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