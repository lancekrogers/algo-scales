package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Stub implementations for remaining screens

// Problem Detail
func (m Model) updateProblemDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Start session with selected problem
			return m.navigate(StateSession), startSession(m.problemDetail.problem)
		}
	}
	return m, nil
}

func (m Model) viewProblemDetail() string {
	p := m.problemDetail.problem
	return "Problem: " + p.Title + "\n\nPress Enter to start, Esc to go back"
}

// Session
func (m Model) updateSession(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case sessionTickMsg:
		// Update timer
		return m, sessionTick()
	}
	return m, nil
}

func (m Model) viewSession() string {
	return "Session in progress...\n\nPress Esc to go back"
}

// Stats
func (m Model) updateStats(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsLoadedMsg:
		m.stats.stats = msg.stats
		m.stats.loading = false
	}
	return m, nil
}

func (m Model) viewStats() string {
	if m.stats.loading {
		return "Loading statistics..."
	}
	return "Statistics\n\nPress Esc to go back"
}

// Daily
func (m Model) updateDaily(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dailyScaleLoadedMsg:
		m.daily.currentScale = msg.scale
		m.daily.progress = msg.progress
	}
	return m, nil
}

func (m Model) viewDaily() string {
	return "Daily Scale: " + m.daily.currentScale + "\n\nPress Esc to go back"
}

// Settings
func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) viewSettings() string {
	return "Settings\n\nPress Esc to go back"
}

// startSession creates a command to start a new session
func startSession(prob problem.Problem) tea.Cmd {
	return func() tea.Msg {
		// Create session (simplified for now)
		return sessionStartedMsg{
			sessionID: "session-" + prob.ID,
		}
	}
}