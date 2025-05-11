// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import (
	"time"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

// SessionMode represents a practice session mode
type SessionMode string

const (
	// LearnMode shows pattern explanations and walkthroughs
	LearnMode SessionMode = "learn"
	// PracticeMode hides solutions but allows hints
	PracticeMode SessionMode = "practice"
	// CramMode focuses on rapid-fire practice with timers
	CramMode SessionMode = "cram"
)

// SessionOptions represents configuration options for a session
type SessionOptions struct {
	Mode       SessionMode
	Language   string
	Timer      int
	Pattern    string
	Difficulty string
	ProblemID  string
}

// TestResult represents the result of a test case
type TestResult struct {
	Input    string
	Expected string
	Actual   string
	Passed   bool
}

// Session represents an active problem-solving session
type Session interface {
	// GetProblem returns the current problem
	GetProblem() *problem.Problem
	
	// GetOptions returns the session options
	GetOptions() SessionOptions
	
	// GetStartTime returns when the session started
	GetStartTime() time.Time
	
	// GetTimeRemaining returns the remaining session time
	GetTimeRemaining() time.Duration
	
	// GetLanguage returns the programming language
	GetLanguage() string
	
	// ShowHints toggles hint display
	ShowHints(show bool)
	
	// ShowSolution toggles solution display
	ShowSolution(show bool)
	
	// AreHintsShown returns if hints are visible
	AreHintsShown() bool
	
	// IsSolutionShown returns if solution is visible
	IsSolutionShown() bool
	
	// FormatDescription returns formatted problem description
	FormatDescription() string
	
	// GetCode returns the current solution code
	GetCode() string
	
	// SetCode updates the solution code
	SetCode(code string) error
	
	// RunTests executes tests on the current solution
	RunTests() ([]TestResult, bool, error)
	
	// Finish completes the session and records stats
	Finish(solved bool) error
}

// SessionManager creates and manages problem-solving sessions
type SessionManager interface {
	// StartSession begins a new practice session
	StartSession(opts SessionOptions) (Session, error)
	
	// GetSessionByID retrieves an active session
	GetSessionByID(id string) (Session, bool)
	
	// FinishSession completes a session
	FinishSession(sessionID string, solved bool) error
}