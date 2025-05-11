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

// ShowHintsMsg is sent to toggle visibility of hints
type ShowHintsMsg struct {
	Show bool
}

// ShowSolutionMsg is sent to toggle visibility of solution
type ShowSolutionMsg struct {
	Show bool
}

// SelectionMsg is sent when a selection is made from a menu
type SelectionMsg struct {
	Index int
}

// EditCodeMsg is sent when the user wants to edit code
type EditCodeMsg struct{}

// QuitMsg is sent when the application should quit
type QuitMsg struct{}