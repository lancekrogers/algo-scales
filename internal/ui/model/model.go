// Package model contains the data models for the UI
package model

import (
	"time"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

// UIModel represents the state of the UI
type UIModel struct {
	// Application state
	AppState AppState

	// Current session details
	Session Session

	// Available problems for selection
	AvailableProblems []problem.Problem

	// User statistics and achievements
	Stats     Statistics
	Achievements map[string]Achievement

	// UI state
	ActiveScreen  ScreenType
	SelectedIndex int
	InputField    string
	ShowHelp      bool
	Loading       bool
	ErrorMessage  string
}

// AppState represents the application state
type AppState int

const (
	StateInitial AppState = iota
	StateOnboarding
	StateProblemSelection
	StateSession
	StateStatistics
	StateSettings
)

// ScreenType represents the type of screen currently displayed
type ScreenType int

const (
	ScreenWelcome ScreenType = iota
	ScreenOnboarding
	ScreenModeSelection
	ScreenPatternSelection 
	ScreenProblemList
	ScreenProblem
	ScreenStats
	ScreenAchievements
	ScreenHelp
	ScreenSettings
)

// Session represents the current practice session
type Session struct {
	Active         bool
	Mode           string // "learn", "practice", "cram"
	Problem        *problem.Problem
	StartTime      time.Time
	TimeRemaining  time.Duration
	ShowHints      bool
	ShowSolution   bool
	Language       string
	Code           string
	TestResults    []TestResult
	CurrentPattern string
}

// TestResult represents the result of a test case
type TestResult struct {
	Input    string
	Expected string
	Actual   string
	Passed   bool
}

// Statistics represents user statistics
type Statistics struct {
	ProblemsAttempted  int
	ProblemsSolved     int
	TotalTime          time.Duration
	PatternCounts      map[string]int
	DifficultyCounts   map[string]int
	CurrentStreak      int
	LongestStreak      int
	LastPracticeDate   time.Time
	PatternsProgress   map[string]float64 // 0.0 to 1.0
}

// Achievement represents a user achievement
type Achievement struct {
	ID          string
	Title       string
	Description string
	Earned      bool
	EarnedDate  time.Time
	Icon        string
}

// Theme contains theming options for the UI
type Theme struct {
	Background  string
	Text        string
	Accent      string
	Error       string
	Success     string
	Border      string
	CurrentMode string
}

// NewModel creates a new UI model
func NewModel() UIModel {
	return UIModel{
		AppState:    StateInitial,
		ActiveScreen: ScreenWelcome,
		Stats: Statistics{
			PatternCounts:    make(map[string]int),
			DifficultyCounts: make(map[string]int),
			PatternsProgress: make(map[string]float64),
		},
		Achievements: make(map[string]Achievement),
		Session: Session{
			Active: false,
		},
	}
}