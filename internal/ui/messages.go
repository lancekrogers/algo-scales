package ui

import (
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/config"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// Navigation messages
type navigateToMsg struct {
	state State
}

// navigateBackMsg is used to navigate back to the previous screen
type navigateBackMsg struct{}

// Data loading messages
type problemsLoadedMsg struct {
	problems []problem.Problem
}

type problemsErrorMsg struct {
	err error
}

type statsLoadedMsg struct {
	stats stats.Summary
}

type statsErrorMsg struct {
	err error
}

// Session messages
type startSessionMsg struct {
	problem problem.Problem
}

type sessionStartedMsg struct {
	sessionID string
}

type sessionTickMsg struct{}

type sessionCompletedMsg struct {
	duration time.Duration
	solved   bool
}

type showHintMsg struct{}

type showSolutionMsg struct{}

// Pattern selection messages
type patternSelectedMsg struct {
	pattern string
}

// Problem selection messages
type problemSelectedMsg struct {
	problem problem.Problem
}

// Daily scale messages
type dailyScaleLoadedMsg struct {
	scale string
	progress interface{}
}

// Settings messages
type settingChangedMsg struct {
	key   string
	value interface{}
}

type saveSettingsMsg struct{}

// Config messages
type configLoadedMsg struct {
	config config.UserConfig
}

// SelectionChangedMsg is sent when the user makes a selection
type SelectionChangedMsg struct {
	State State
}

// GetState returns the state from the message
func (m SelectionChangedMsg) GetState() State {
	return m.State
}