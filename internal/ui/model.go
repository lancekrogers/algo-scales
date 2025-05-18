package ui

import (
	"time"
	
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/common/config"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// Model represents the main application model
type Model struct {
	// Current application state
	state State
	
	// Previous state for back navigation
	previousState State
	
	// Component models
	home          homeModel
	patterns      patternModel
	problems      problemListModel
	problemDetail problemDetailModel
	session       sessionModel
	stats         statsModel
	daily         dailyModel
	settings      settingsModel
	
	// Common data
	config    config.UserConfig
	allProblems  []problem.Problem
	width     int
	height    int
	ready     bool
	
	// Global key bindings
	keys globalKeyMap
}

// homeModel represents the home screen state
type homeModel struct {
	selectedOption int
	options        []string
}

// patternModel represents the pattern selection state
type patternModel struct {
	patterns       []string
	selectedIndex  int
	selectedPattern string
}

// problemListModel represents the problem list state
type problemListModel struct {
	problems      []problem.Problem
	selectedIndex int
	pattern       string
	loading       bool
}

// problemDetailModel represents the problem detail view state
type problemDetailModel struct {
	problem  problem.Problem
	showHint bool
	showInfo bool
	viewport viewport.Model
}

// sessionModel represents the active session state
type sessionModel struct {
	sessionID    string
	problem      problem.Problem
	showHint     bool
	showSolution bool
	timerPaused  bool
	startTime    time.Time
	duration     time.Duration
	viewport     viewport.Model
	testResults  string
	message      string
	confirmQuit  bool
}

// statsModel represents the statistics view state
type statsModel struct {
	loading bool
	summary stats.Summary
	viewport viewport.Model
}

// dailyModel represents the daily challenge state
type dailyModel struct {
	currentScale string
	progress     interface{} // Can be daily.ScaleProgress
	loading      bool
}

// settingsModel represents the settings view state
type settingsModel struct {
	selectedOption int
	editing        bool
	editingField   string
	editValue      string
	message        string
}

// globalKeyMap defines global keyboard shortcuts
type globalKeyMap struct {
	Quit key.Binding
	Back key.Binding
	Help key.Binding
}

// NewModel creates a new application model
func NewModel() Model {
	return Model{
		state: StateHome,
		home: homeModel{
			options: []string{
				"Start Practice Session",
				"Daily Scales",
				"View Statistics",
				"Settings",
			},
		},
		keys: globalKeyMap{
			Quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q/ctrl+c", "quit"),
			),
			Back: key.NewBinding(
				key.WithKeys("esc", "left"),
				key.WithHelp("esc/←", "back"),
			),
			Help: key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("?", "help"),
			),
		},
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Load initial data
	return loadProblems()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil
		
	case tea.KeyMsg:
		// Handle global key bindings
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m.handleBack()
		}
	}
	
	// Route updates to current state
	return m.routeUpdate(msg)
}

// View renders the current state
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}
	
	switch m.state {
	case StateHome:
		return m.viewHome()
	case StatePatternSelection:
		return m.viewPatterns()
	case StateProblemList:
		return m.viewProblemList()
	case StateProblemDetail:
		return m.viewProblemDetail()
	case StateSession:
		return m.viewSession()
	case StateStats:
		return m.viewStats()
	case StateDaily:
		return m.viewDaily()
	case StateSettings:
		return m.viewSettings()
	default:
		return "Unknown state"
	}
}

// Navigation methods
func (m Model) handleBack() (Model, tea.Cmd) {
	// Navigate back to previous state
	if m.previousState != m.state {
		m.state = m.previousState
	} else {
		// Default back navigation
		switch m.state {
		case StatePatternSelection, StateDaily, StateStats, StateSettings:
			m.state = StateHome
		case StateProblemList:
			m.state = StatePatternSelection
		case StateProblemDetail:
			m.state = StateProblemList
		case StateSession:
			m.state = StateProblemDetail
		default:
			m.state = StateHome
		}
	}
	return m, nil
}

func (m Model) navigate(newState State) Model {
	m.previousState = m.state
	m.state = newState
	return m
}

// routeUpdate routes updates to the appropriate component
func (m Model) routeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateHome:
		return m.updateHome(msg)
	case StatePatternSelection:
		return m.updatePatterns(msg)
	case StateProblemList:
		return m.updateProblemList(msg)
	case StateProblemDetail:
		return m.updateProblemDetail(msg)
	case StateSession:
		return m.updateSession(msg)
	case StateStats:
		return m.updateStats(msg)
	case StateDaily:
		return m.updateDaily(msg)
	case StateSettings:
		return m.updateSettings(msg)
	default:
		return m, nil
	}
}