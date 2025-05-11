// Terminal UI using Bubble Tea
package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/lancekrogers/algo-scales/internal/ui/model"
)

// Define styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0A84FF")).
			Padding(0, 1)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFCC00")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF453A")).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#30D158")).
			Padding(0, 1)

	timerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF9500")).
			Padding(0, 1).
			MarginLeft(1)

	timerWarningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#FF453A")).
				Padding(0, 1).
				MarginLeft(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			MarginTop(1).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A9DEDE"))

	paragraphStyle = lipgloss.NewStyle().
			MarginTop(1).
			MarginBottom(1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#DDDDDD")).
			Padding(0, 1).
			Width(100).
			Height(1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF9500")).
			Padding(0, 3).
			MarginRight(1)

	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 3).
				MarginRight(1)

	fieldStyle = lipgloss.NewStyle().
			Width(25).
			MarginRight(2)

	codeBlockStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1E1E1E")).
			Foreground(lipgloss.Color("#E2E2E2")).
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#404040"))
)

// KeyMap defines the key mappings
type KeyMap struct {
	Help         key.Binding
	Quit         key.Binding
	Editor       key.Binding
	ShowHints    key.Binding
	ShowSolution key.Binding
	Submit       key.Binding
	Skip         key.Binding
	Testing      key.Binding
}

// DefaultKeyMap returns the default key map
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Editor: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "open editor"),
		),
		ShowHints: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "show hints"),
		),
		ShowSolution: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "show solution"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit solution"),
		),
		Skip: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "skip to next problem"),
		),
		Testing: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "run tests"),
		),
	}
}

// Model represents the UI model
type Model struct {
	session          *session.Session
	viewport         viewport.Model
	timer            timer.Model
	keyMap           KeyMap
	help             help.Model
	showHelp         bool
	width            int
	height           int
	editorOpened     bool
	message          string
	messageStyle     lipgloss.Style
	timeRemaining    time.Duration
	problemCompleted bool
	ready            bool
	testing          bool
	testResults      string
	shouldQuit       bool
	confirmQuit      bool
}

// StartSession starts the UI session
func StartSession(s *session.Session) error {
	// Initialize the model
	m := Model{
		session:      s,
		keyMap:       DefaultKeyMap(),
		help:         help.New(),
		messageStyle: infoStyle,
		message:      "Press '?' for help, 'e' to open editor",
	}

	// Set up timer
	timerDuration := time.Duration(s.Options.Timer) * time.Minute
	m.timer = timer.NewWithInterval(timerDuration, time.Second)
	m.timeRemaining = timerDuration

	// Create the program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run the program
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running UI: %v", err)
	}

	return nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Start the timer
	return tea.Batch(
		m.timer.Init(),
		tea.EnterAltScreen,
	)
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't process input when editor is opened
		if m.editorOpened {
			return m, nil
		}

		// Ask for confirmation before quitting
		if m.confirmQuit {
			switch msg.String() {
			case "y", "Y":
				// Save session stats before quitting
				if !m.problemCompleted {
					_ = m.session.FinishSession(false)
				}
				m.shouldQuit = true
				return m, tea.Quit
			case "n", "N", "esc":
				m.confirmQuit = false
				m.message = "Quit cancelled"
				m.messageStyle = infoStyle
				return m, nil
			}
		}

		// Handle key press
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			if m.session.Options.Mode == session.CramMode && !m.problemCompleted {
				// In Cram mode, ask for confirmation before quitting
				m.confirmQuit = true
				m.message = "Quit without completing? (y/n)"
				m.messageStyle = warningStyle
				return m, nil
			}

			// Record the session as not solved if not completed
			if !m.problemCompleted {
				_ = m.session.FinishSession(false)
			}
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, m.keyMap.Editor):
			// Open the editor
			m.editorOpened = true
			m.message = "Opening editor..."
			m.messageStyle = infoStyle
			return m, m.openEditor

		case key.Matches(msg, m.keyMap.ShowHints):
			if !m.session.ShowHints {
				m.session.ShowHints = true
				m.message = "Hints shown"
				m.messageStyle = infoStyle

				// Update the problem description
				content := m.session.FormatProblemDescription()
				m.viewport.SetContent(content)
			}
			return m, nil

		case key.Matches(msg, m.keyMap.ShowSolution):
			if !m.session.ShowSolution {
				m.session.ShowSolution = true
				m.message = "Solution shown"
				m.messageStyle = infoStyle

				// Update the problem description
				content := m.session.FormatProblemDescription()
				m.viewport.SetContent(content)
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Submit):
			// This would verify the solution in a real implementation
			// For MVP, we'll just mark it as solved
			m.message = "Solution submitted successfully!"
			m.messageStyle = successStyle
			m.problemCompleted = true

			// Record the session as solved
			_ = m.session.FinishSession(true)

			// In cram mode, we might move to the next problem automatically
			if m.session.Options.Mode == session.CramMode {
				m.message += " Press 'n' for next problem"
			}

			return m, nil

		case key.Matches(msg, m.keyMap.Skip):
			if m.session.Options.Mode == session.CramMode && !m.problemCompleted {
				// In Cram mode, ask for confirmation before skipping
				m.confirmQuit = true
				m.message = "Skip this problem? (y/n)"
				m.messageStyle = warningStyle
				return m, nil
			}

			if !m.problemCompleted {
				// Record the session as not solved
				_ = m.session.FinishSession(false)
			}

			// In a real implementation, this would move to the next problem
			// For MVP, we'll just quit
			m.message = "Skipped to next problem"
			m.messageStyle = warningStyle
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.Testing):
			// This would run tests in a real implementation
			m.testing = true
			m.message = "Running tests..."
			m.messageStyle = infoStyle
			return m, m.runTests
		}

	case tea.WindowSizeMsg:
		// Handle window resize
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// This is the first time we're getting a window size
			headerHeight := 3 // Title + Timer
			footerHeight := 3 // Status bar + help
			contentHeight := m.height - headerHeight - footerHeight

			// Set up the viewport for the problem description
			m.viewport = viewport.New(msg.Width, contentHeight)
			content := m.session.FormatProblemDescription()
			m.viewport.SetContent(content)

			// Set up help menu
			m.help.Width = msg.Width

			m.ready = true
		} else {
			// Just adjust the viewport size
			headerHeight := 3 // Title + Timer
			footerHeight := 3 // Status bar + help
			contentHeight := m.height - headerHeight - footerHeight

			m.viewport.Width = msg.Width
			m.viewport.Height = contentHeight

			// Set help width
			m.help.Width = msg.Width
		}

	case timer.TickMsg:
		// Update timer
		var timerCmd tea.Cmd
		m.timer, timerCmd = m.timer.Update(msg)
		cmds = append(cmds, timerCmd)

		// Update time remaining
		m.timeRemaining = m.timer.Timeout

		// Change timer style if less than 10 minutes left
		if m.timeRemaining < 10*time.Minute && m.timeRemaining > 0 {
			// Only show the warning once
			if m.message != "Less than 10 minutes remaining!" {
				m.message = "Less than 10 minutes remaining!"
				m.messageStyle = warningStyle
			}
		}

	case timer.TimeoutMsg:
		// Timer has expired
		m.message = "Time's up!"
		m.messageStyle = errorStyle

		// In practice mode, we might want to allow continuing
		// In cram mode, we might want to move to the next problem
		if m.session.Options.Mode == session.CramMode {
			// Record the session as not solved
			_ = m.session.FinishSession(false)
			m.message += " Moving to next problem..."
			return m, tea.Quit
		}

	case editorFinishedMsg:
		// Editor has finished
		m.editorOpened = false
		m.message = "Code saved"
		m.messageStyle = infoStyle

	case uiTestResultsMsg:
		// Test results
		m.testing = false
		m.testResults = fmt.Sprintf("%v", msg)

		if strings.Contains(m.testResults, "PASS") {
			m.message = "Tests passed!"
			m.messageStyle = successStyle
		} else {
			m.message = "Tests failed!"
			m.messageStyle = errorStyle
		}
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.shouldQuit {
		return ""
	}

	// Format remaining time
	timeStr := formatDuration(m.timeRemaining)
	var timerView string
	if m.timeRemaining < 10*time.Minute {
		timerView = timerWarningStyle.Render(fmt.Sprintf("Time: %s", timeStr))
	} else {
		timerView = timerStyle.Render(fmt.Sprintf("Time: %s", timeStr))
	}

	// Build the header
	modeStr := string(m.session.Options.Mode)

	// Get musical scale info if available
	scaleInfo := ""
	if len(m.session.Problem.Patterns) > 0 {
		pattern := m.session.Problem.Patterns[0]
		switch pattern {
		case "sliding-window":
			scaleInfo = " ♪ C Major (Sliding Window) ♪"
		case "two-pointers":
			scaleInfo = " ♪ G Major (Two Pointers) ♪"
		case "fast-slow-pointers":
			scaleInfo = " ♪ D Major (Fast & Slow Pointers) ♪"
		case "hash-map":
			scaleInfo = " ♪ A Major (Hash Maps) ♪"
		case "binary-search":
			scaleInfo = " ♪ E Major (Binary Search) ♪"
		case "dfs":
			scaleInfo = " ♪ B Major (DFS) ♪"
		case "bfs":
			scaleInfo = " ♪ F# Major (BFS) ♪"
		case "dynamic-programming":
			scaleInfo = " ♪ Db Major (Dynamic Programming) ♪"
		case "greedy":
			scaleInfo = " ♪ Ab Major (Greedy) ♪"
		case "union-find":
			scaleInfo = " ♪ Eb Major (Union-Find) ♪"
		case "heap":
			scaleInfo = " ♪ Bb Major (Heap) ♪"
		}
	}

	// Build title with optional scale info
	title := fmt.Sprintf("%s (%s)%s", m.session.Problem.Title, m.session.Problem.Difficulty, scaleInfo)
	header := titleStyle.Copy().Width(m.width - len(modeStr) - len(timeStr) - 4).Render(title)

	// Add mode and timer to header
	modeView := infoStyle.Render(modeStr)
	headerBar := lipgloss.JoinHorizontal(lipgloss.Top, header, modeView, timerView)

	// Build the footer with message line
	footer := ""
	if m.message != "" {
		footer = m.messageStyle.Copy().Width(m.width).Render(m.message)
	}

	// Add test results if available
	testView := ""
	if m.testResults != "" {
		testView = "\n" + codeBlockStyle.Width(m.width-4).Render(m.testResults)
	}

	// Build the help line
	helpView := ""
	if m.showHelp {
		// Fixed help view with empty key bindings
		helpView = "\n" + "Press ? for help, q to quit"
	}

	// Show confirmation dialog
	if m.confirmQuit {
		return fmt.Sprintf("%s\n%s\n\n%s%s%s",
			headerBar,
			m.viewport.View(),
			footer,
			testView,
			helpView)
	}

	// Put it all together
	return fmt.Sprintf("%s\n%s\n\n%s%s%s",
		headerBar,
		m.viewport.View(),
		footer,
		testView,
		helpView)
}

// openEditor opens the code in the user's editor
func (m Model) openEditor() tea.Msg {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Default to vim or notepad based on OS
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vim"
		}
	}

	// Prepare the command
	cmd := exec.Command(editor, m.session.CodeFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the editor
	err := cmd.Run()
	if err != nil {
		return editorErrorMsg{err}
	}

	return editorFinishedMsg{}
}

// runTests runs the tests for the current problem
func (m Model) runTests() tea.Msg {
	// This would run the tests in a real implementation
	// For MVP, we'll just return a mock result

	language := m.session.Options.Language
	problemID := m.session.Problem.ID

	// In a real implementation, we would compile and run the code
	// For now, just return a mock result
	mockResults := fmt.Sprintf("Running tests for %s in %s...\n\n", problemID, language)

	// For demo purposes, generate a random result
	if time.Now().UnixNano()%2 == 0 {
		mockResults += "✓ PASS: Test case 1\n"
		mockResults += "✓ PASS: Test case 2\n"
		mockResults += "✓ PASS: Test case 3\n\n"
		mockResults += "All tests passed!"
	} else {
		mockResults += "✓ PASS: Test case 1\n"
		mockResults += "✗ FAIL: Test case 2\n"
		mockResults += "  Expected: [1, 2], Got: [1, 3]\n"
		mockResults += "✓ PASS: Test case 3\n\n"
		mockResults += "1 test failed"
	}

	// Wait a bit to simulate running tests
	time.Sleep(1 * time.Second)

	return uiTestResultsMsg{
		Results:   []model.TestResult{},
		AllPassed: strings.Contains(mockResults, "✓ PASS"),
	}
}

// Custom message types
type (
	editorFinishedMsg struct{}
	editorErrorMsg    struct{ error }
	// Using local UI package test results message
	uiTestResultsMsg   struct {
		Results   []model.TestResult
		AllPassed bool
	}
)

// formatDuration formats a duration as MM:SS
func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
