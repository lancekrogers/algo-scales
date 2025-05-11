package model

import (
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Message types for communication between components

// ErrorMsg is sent when an error occurs
type ErrorMsg string

// ProblemsLoadedMsg is sent when the problem list is loaded
type ProblemsLoadedMsg struct {
	Problems []problem.Problem
}

// TickMsg is sent every second for timers
type TickMsg struct{}

// ProblemSelectedMsg is sent when a problem is selected
type ProblemSelectedMsg struct {
	ProblemID string
	Mode      string
}

// CodeUpdatedMsg is sent when the code is updated in the editor
type CodeUpdatedMsg struct{}

// TestResultsMsg is sent after running tests
type TestResultsMsg struct {
	Results   []TestResult
	AllPassed bool
}

// AchievementUnlockedMsg is sent when an achievement is unlocked
type AchievementUnlockedMsg struct {
	AchievementID string
}