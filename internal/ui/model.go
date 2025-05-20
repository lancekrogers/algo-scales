package ui

import (
	"fmt"
	"strings"
	"time"
	
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/common/config"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// New creates a new model instance
func New() Model {
	// Load user config with defaults
	cfg, err := config.LoadConfig()
	if err != nil {
		// Use default config if loading fails
		cfg = config.DefaultConfig()
	}
	
	return Model{
		state: StateHome,
		config: cfg,
		home: homeModel{
			selectedOption: 0,
			options: []string{
				"Start Practice Session", 
				"Daily Scales",
				"View Statistics",
				"Settings",
			},
		},
		patterns:      patternModel{},
		problems:      problemListModel{},
		problemDetail: problemDetailModel{},
		session:       sessionModel{},
		stats:         statsModel{},
		daily:         dailyModel{},
		settings:      settingsModel{},
		keys:          globalKeyMap{
			Quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
			Back: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "back"),
			),
		},
		animation:     Animation{Type: AnimationNone},
		loading:       LoadingScreen{},
	}
}

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
	
	// Animation state
	animation     Animation
	loading       LoadingScreen
	showLoading   bool
	spinnerTicks  int
	
	// Error state
	errorMessage  string
}

// homeModel represents the home screen state
type homeModel struct {
	selectedOption int
	options        []string
	width          int
	height         int
}

// Init initializes the home model
func (m homeModel) Init() tea.Cmd {
	return nil
}

// Update handles updates for the home model
func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
		
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedOption > 0 {
				m.selectedOption--
			}
		case "down", "j":
			if m.selectedOption < len(m.options)-1 {
				m.selectedOption++
			}
		case "enter", "right", "l":
			// Return appropriate state change message
			switch m.selectedOption {
			case 0: // Start Practice Session
				return m, func() tea.Msg { return SelectionChangedMsg{State: StatePatternSelection} }
			case 1: // Daily Scales
				return m, func() tea.Msg { return SelectionChangedMsg{State: StateDaily} }
			case 2: // View Statistics  
				return m, func() tea.Msg { return SelectionChangedMsg{State: StateStats} }
			case 3: // Settings
				return m, func() tea.Msg { return SelectionChangedMsg{State: StateSettings} }
			}
		}
	}
	return m, nil
}

// View renders the home model
func (m homeModel) View() string {
	var b strings.Builder
	
	// Title
	b.WriteString(titleStyle.Render("🎵 AlgoScales"))
	b.WriteString("\n\n")
	
	// Menu options
	for i, option := range m.options {
		cursor := "  "
		if i == m.selectedOption {
			cursor = cursorStyle.Render("> ")
			option = selectedItemStyle.Render(option)
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, option))
	}
	
	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • q: Quit"))
	
	return b.String()
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
	return tea.Batch(
		loadProblems(),
		loadConfig(),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		if m.showLoading {
			m.loading.width = msg.Width
			m.loading.height = msg.Height
		}
		// Propagate window size to all components
		m.home.width = msg.Width
		m.home.height = msg.Height
		// TODO: Propagate to other components when they need it
		return m, nil
		
	case animationTickMsg:
		m.animation.Update()
		if !m.animation.Complete {
			cmds = append(cmds, AnimationTick())
		}
		
	case spinnerTickMsg:
		m.spinnerTicks++
		if m.showLoading {
			m.loading.spinnerFrame = m.spinnerTicks
			cmds = append(cmds, tickSpinner())
		}
		
	case startLoadingMsg:
		m.showLoading = true
		m.loading = NewLoadingScreen(msg.message)
		m.loading.width = m.width
		m.loading.height = m.height
		cmds = append(cmds, tickSpinner())
		
	case stopLoadingMsg:
		m.showLoading = false
		
	case problemsLoadedMsg:
		m.allProblems = msg.problems
		
	case configLoadedMsg:
		m.config = msg.config
		
	case navigateBackMsg:
		m, cmd = m.handleBack()
		cmds = append(cmds, cmd)
		// Start slide animation
		m.animation = NewAnimation(AnimationSlideLeft, 300*time.Millisecond)
		cmds = append(cmds, AnimationTick())
		return m, tea.Batch(cmds...)
		
	case SelectionChangedMsg:
		m = m.navigate(msg.State)
		cmds = append(cmds, AnimationTick())
		return m, tea.Batch(cmds...)
		
	case tea.KeyMsg:
		// Handle global key bindings
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			m, cmd = m.handleBack()
			cmds = append(cmds, cmd)
			// Start slide animation
			m.animation = NewAnimation(AnimationSlideLeft, 300*time.Millisecond)
			cmds = append(cmds, AnimationTick())
			return m, tea.Batch(cmds...)
		}
	}
	
	// Handle loading screen updates
	if m.showLoading {
		m.loading, cmd = m.loading.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	
	// Route updates to current state
	m, cmd = m.routeUpdate(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

// View renders the current state
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}
	
	// Show loading screen if active
	if m.showLoading {
		return m.loading.View()
	}
	
	// Render current state
	var content string
	switch m.state {
	case StateHome:
		content = m.viewHome()
	case StatePatternSelection:
		content = m.viewPatterns()
	case StateProblemList:
		content = m.viewProblemList()
	case StateProblemDetail:
		content = m.viewProblemDetail()
	case StateSession:
		content = m.viewSession()
	case StateStats:
		content = m.viewStats()
	case StateDaily:
		content = m.viewDaily()
	case StateSettings:
		content = m.viewSettings()
	default:
		content = "Unknown state"
	}
	
	// Apply animation if active
	if !m.animation.Complete {
		content = m.animation.Apply(content, m.width, m.height)
	}
	
	return content
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
	
	// Start appropriate animation based on state transition
	if newState > m.previousState {
		// Moving forward
		m.animation = NewAnimation(AnimationSlideRight, 300*time.Millisecond)
	} else {
		// Moving backward
		m.animation = NewAnimation(AnimationSlideLeft, 300*time.Millisecond)
	}
	
	return m
}

// routeUpdate routes updates to the appropriate component
func (m Model) routeUpdate(msg tea.Msg) (Model, tea.Cmd) {
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