// Package screens contains UI screens for different app states
package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/common/highlight"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/ui/view"
)

// SessionKeyMap defines key mappings for the session screen
type SessionKeyMap struct {
	EditCode     key.Binding
	ShowHints    key.Binding
	ShowSolution key.Binding
	RunTests     key.Binding
	Submit       key.Binding
	Skip         key.Binding
	Help         key.Binding
	Quit         key.Binding
}

// NewSessionKeyMap creates a new key map for the session
func NewSessionKeyMap() SessionKeyMap {
	return SessionKeyMap{
		EditCode: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit code"),
		),
		ShowHints: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "show hints"),
		),
		ShowSolution: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "show solution"),
		),
		RunTests: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "run tests"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit solution"),
		),
		Skip: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next problem"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// SessionModel represents the session screen model
type SessionModel struct {
	// Session state
	Problem          *problem.Problem
	Mode             string    // "learn", "practice", "cram"
	Language         string
	Code             string
	StartTime        time.Time
	ShowHints        bool
	ShowSolution     bool
	ProblemCompleted bool
	CurrentPattern   string

	// UI components
	ProblemViewport viewport.Model
	CodeViewport    viewport.Model
	CodeInput       textinput.Model
	Timer           timer.Model
	TimeRemaining   time.Duration
	Spinner         spinner.Model
	Help            help.Model
	KeyMap          SessionKeyMap

	// Screen state
	ShowHelp     bool
	Message      string
	MessageStyle lipgloss.Style
	Testing      bool
	TestResults  []TestResult
	AllPassed    bool
	Loading      bool
	ConfirmQuit  bool
	Width        int
	Height       int
	Ready        bool
	EditorOpened bool

	// Rendering components
	SyntaxHighlighter *highlight.SyntaxHighlighter
	PatternViz        *view.PatternVisualization
}

// TestResult represents the result of a test case
type TestResult struct {
	Input    string
	Expected string
	Actual   string
	Passed   bool
}

// NewSessionModel creates a new session model
func NewSessionModel(prob *problem.Problem, mode, language string, currentPattern string) SessionModel {
	// Create key map
	keyMap := NewSessionKeyMap()

	// Create help component
	help := help.New()

	// Create timer
	var timerDuration time.Duration
	switch mode {
	case "learn":
		timerDuration = 45 * time.Minute
	case "practice":
		timerDuration = 30 * time.Minute
	case "cram":
		timerDuration = 15 * time.Minute
	default:
		timerDuration = 30 * time.Minute
	}
	
	t := timer.NewWithInterval(timerDuration, time.Second)

	// Create spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	// Create syntax highlighter
	syntaxHighlighter := highlight.NewSyntaxHighlighter("monokai")

	// Create pattern visualization
	patternViz := view.NewPatternVisualization()

	// Get starter code for the selected language
	code := ""
	if prob != nil && prob.StarterCode != nil {
		if starter, ok := prob.StarterCode[language]; ok {
			code = starter
		}
	}

	return SessionModel{
		Problem:           prob,
		Mode:              mode,
		Language:          language,
		StartTime:         time.Now(),
		TimeRemaining:     timerDuration,
		CurrentPattern:    currentPattern,
		KeyMap:            keyMap,
		Help:              help,
		Timer:             t,
		Spinner:           s,
		Message:           "Press '?' for help, 'e' to open editor",
		MessageStyle:      view.InfoStyle,
		SyntaxHighlighter: syntaxHighlighter,
		PatternViz:        patternViz,
		Code:              code,
	}
}

// Init initializes the session model
func (m SessionModel) Init() tea.Cmd {
	return tea.Batch(
		m.Timer.Init(),
		spinner.Tick,
	)
}

// Update handles updates to the session model
func (m SessionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		m.Width = msg.Width
		m.Height = msg.Height

		// Set up split view
		if !m.Ready {
			// Calculate dimensions for split view
			headerHeight := 5  // Title + Mode + Timer + separator
			footerHeight := 5  // Status bar + help + message
			contentHeight := m.Height - headerHeight - footerHeight

			// Split content area 40/60 for problem/code
			problemWidth := m.Width * 4 / 10
			codeWidth := m.Width - problemWidth - 2 // 2 for separator

			// Set up problem viewport
			m.ProblemViewport = viewport.New(problemWidth, contentHeight)
			m.ProblemViewport.SetContent(m.formatProblemContent())

			// Set up code viewport
			m.CodeViewport = viewport.New(codeWidth, contentHeight)
			m.CodeViewport.SetContent(m.formatCodeContent())

			// Set up help menu
			m.Help.Width = m.Width

			m.Ready = true
		} else {
			// Recalculate viewport sizes
			headerHeight := 5
			footerHeight := 5
			contentHeight := m.Height - headerHeight - footerHeight

			problemWidth := m.Width * 4 / 10
			codeWidth := m.Width - problemWidth - 2

			m.ProblemViewport.Width = problemWidth
			m.ProblemViewport.Height = contentHeight

			m.CodeViewport.Width = codeWidth
			m.CodeViewport.Height = contentHeight

			m.Help.Width = m.Width
		}

		// Refresh contents
		m.ProblemViewport.SetContent(m.formatProblemContent())
		m.CodeViewport.SetContent(m.formatCodeContent())

		return m, nil

	case tea.KeyMsg:
		// Don't process input when editor is opened or when confirming quit
		if m.EditorOpened {
			return m, nil
		}

		// Handle quit confirmation
		if m.ConfirmQuit {
			switch msg.String() {
			case "y", "Y":
				return m, tea.Quit
			case "n", "N", "esc":
				m.ConfirmQuit = false
				m.Message = "Quit cancelled"
				m.MessageStyle = view.InfoStyle
				return m, nil
			}
		}

		// Handle key press
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			if m.Mode == "cram" && !m.ProblemCompleted {
				// In Cram mode, ask for confirmation before quitting
				m.ConfirmQuit = true
				m.Message = "Quit without completing? (y/n)"
				m.MessageStyle = view.WarningStyle
				return m, nil
			}
			return m, tea.Quit

		case key.Matches(msg, m.KeyMap.Help):
			m.ShowHelp = !m.ShowHelp
			return m, nil

		case key.Matches(msg, m.KeyMap.EditCode):
			// Placeholder for opening editor
			m.EditorOpened = true
			m.Message = "Opening editor..."
			m.MessageStyle = view.InfoStyle
			return m, func() tea.Msg {
				// Simulate editing - in real implementation, this would open an external editor
				time.Sleep(500 * time.Millisecond)
				return editorFinishedMsg{code: m.Code}
			}

		case key.Matches(msg, m.KeyMap.ShowHints):
			m.ShowHints = true
			m.Message = "Hints shown"
			m.MessageStyle = view.InfoStyle
			// Update problem viewport with hints
			m.ProblemViewport.SetContent(m.formatProblemContent())
			return m, nil

		case key.Matches(msg, m.KeyMap.ShowSolution):
			m.ShowSolution = true
			m.Message = "Solution shown"
			m.MessageStyle = view.InfoStyle
			// Update problem viewport with solution
			m.ProblemViewport.SetContent(m.formatProblemContent())
			return m, nil

		case key.Matches(msg, m.KeyMap.RunTests):
			m.Testing = true
			m.Loading = true
			m.Message = "Running tests..."
			m.MessageStyle = view.InfoStyle
			// Simulate running tests
			return m, func() tea.Msg {
				// This would actually run tests in a real implementation
				time.Sleep(1 * time.Second)
				return testResultsMsg{
					Results: []TestResult{
						{
							Input:    "Input 1",
							Expected: "Expected 1",
							Actual:   "Actual 1",
							Passed:   true,
						},
						{
							Input:    "Input 2",
							Expected: "Expected 2",
							Actual:   "Not matching",
							Passed:   false,
						},
					},
					AllPassed: false,
				}
			}

		case key.Matches(msg, m.KeyMap.Submit):
			m.ProblemCompleted = true
			m.Message = "Solution submitted successfully!"
			m.MessageStyle = view.SuccessStyle
			return m, nil

		case key.Matches(msg, m.KeyMap.Skip):
			if m.Mode == "cram" && !m.ProblemCompleted {
				// In Cram mode, ask for confirmation before skipping
				m.ConfirmQuit = true
				m.Message = "Skip this problem? (y/n)"
				m.MessageStyle = view.WarningStyle
				return m, nil
			}
			return m, tea.Quit
		}

	case timer.TickMsg:
		// Update the timer
		var timerCmd tea.Cmd
		m.Timer, timerCmd = m.Timer.Update(msg)
		cmds = append(cmds, timerCmd)

		// Update time remaining
		m.TimeRemaining = m.Timer.Timeout

		// Change timer style if less than 5 minutes left
		if m.TimeRemaining < 5*time.Minute && m.TimeRemaining > 0 {
			if m.Message != "Less than 5 minutes remaining!" {
				m.Message = "Less than 5 minutes remaining!"
				m.MessageStyle = view.WarningStyle
			}
		}

	case timer.TimeoutMsg:
		// Timer has expired
		m.Message = "Time's up!"
		m.MessageStyle = view.ErrorStyle

		// In cram mode, move to next problem
		if m.Mode == "cram" {
			m.Message += " Moving to next problem..."
			return m, tea.Quit
		}

	case spinner.TickMsg:
		// Update the spinner
		var spinnerCmd tea.Cmd
		m.Spinner, spinnerCmd = m.Spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)

	case editorFinishedMsg:
		// Editor has finished
		m.EditorOpened = false
		m.Code = msg.code
		m.Message = "Code saved"
		m.MessageStyle = view.InfoStyle
		// Update code viewport
		m.CodeViewport.SetContent(m.formatCodeContent())

	case testResultsMsg:
		// Test results received
		m.Testing = false
		m.Loading = false
		m.TestResults = msg.Results
		m.AllPassed = msg.AllPassed

		// Update message based on test results
		if m.AllPassed {
			m.Message = "All tests passed!"
			m.MessageStyle = view.SuccessStyle
		} else {
			m.Message = "Some tests failed!"
			m.MessageStyle = view.ErrorStyle
		}

		// Update the code viewport to show test results
		m.CodeViewport.SetContent(m.formatCodeContent())
	}

	// Update viewports
	var problemCmd, codeCmd tea.Cmd
	m.ProblemViewport, problemCmd = m.ProblemViewport.Update(msg)
	m.CodeViewport, codeCmd = m.CodeViewport.Update(msg)
	cmds = append(cmds, problemCmd, codeCmd)

	return m, tea.Batch(cmds...)
}

// View renders the session screen
func (m SessionModel) View() string {
	if !m.Ready {
		return "Loading..."
	}

	// Format the header with title and pattern information
	title := m.formatTitle()
	modeInfo := m.formatModeInfo()
	timerView := m.formatTimer()
	header := lipgloss.JoinHorizontal(lipgloss.Top, title, modeInfo, timerView)

	// Format the split view
	splitView := m.formatSplitView()

	// Format the footer with message and help
	message := m.formatMessage()
	helpView := m.formatHelp()
	footer := lipgloss.JoinVertical(lipgloss.Left, message, helpView)

	// Put it all together
	return lipgloss.JoinVertical(lipgloss.Left, header, splitView, footer)
}

// formatTitle formats the title section
func (m SessionModel) formatTitle() string {
	if m.Problem == nil {
		return view.TitleStyle.Render("No Problem Selected")
	}

	title := m.Problem.Title
	
	// Add pattern info if available
	if m.CurrentPattern != "" {
		if scale, ok := view.MusicScales[m.CurrentPattern]; ok {
			title += " — " + scale.Name
		}
	}
	
	return view.TitleStyle.Copy().
		Width(m.Width / 2).
		Render(title)
}

// formatModeInfo formats the mode information
func (m SessionModel) formatModeInfo() string {
	difficulty := ""
	if m.Problem != nil {
		difficulty = m.Problem.Difficulty
	}
	
	info := fmt.Sprintf("%s | %s", 
		strings.Title(m.Mode), 
		strings.Title(difficulty),
	)
	
	return view.StatusBarStyle.Copy().
		Width(m.Width / 4).
		Render(info)
}

// formatTimer formats the timer display
func (m SessionModel) formatTimer() string {
	hours := int(m.TimeRemaining.Hours())
	mins := int(m.TimeRemaining.Minutes()) % 60
	secs := int(m.TimeRemaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
	
	if m.TimeRemaining < 5*time.Minute {
		return view.TimerWarningStyle.Copy().
			Width(m.Width / 4).
			Render("Time: " + timeStr)
	}
	
	return view.TimerStyle.Copy().
		Width(m.Width / 4).
		Render("Time: " + timeStr)
}

// formatSplitView formats the split view with problem and code
func (m SessionModel) formatSplitView() string {
	// Add divider between problem and code sections
	divider := view.BorderedBoxStyle.
		Copy().
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Width(1).
		Height(m.ProblemViewport.Height).
		Render("")

	// Join the problem and code viewports with the divider
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.ProblemViewport.View(),
		divider,
		m.CodeViewport.View(),
	)
}

// formatMessage formats the message line
func (m SessionModel) formatMessage() string {
	if m.Loading {
		return m.MessageStyle.Copy().
			Width(m.Width).
			Render(m.Message + " " + m.Spinner.View())
	}
	
	return m.MessageStyle.Copy().
		Width(m.Width).
		Render(m.Message)
}

// formatHelp formats the help view
func (m SessionModel) formatHelp() string {
	if m.ShowHelp {
		// Use a simple help format instead of the help component
		helpText := "e: Edit Code | h: Hints | s: Solution | t: Tests | Enter: Submit | q: Quit"
		return view.HelpStyle.Render(helpText)
	}
	
	return view.HelpStyle.
		Render("Press ? for help")
}

// formatProblemContent formats the problem description
func (m SessionModel) formatProblemContent() string {
	if m.Problem == nil {
		return "No problem selected"
	}

	// Start with the problem description
	content := fmt.Sprintf("%s\n\n", m.Problem.Description)

	// Add examples
	if len(m.Problem.Examples) > 0 {
		content += view.HeaderStyle.Render("Examples:") + "\n\n"
		for i, example := range m.Problem.Examples {
			content += fmt.Sprintf("Example %d:\n", i+1)
			content += fmt.Sprintf("Input: %s\n", example.Input)
			content += fmt.Sprintf("Output: %s\n", example.Output)
			if example.Explanation != "" {
				content += fmt.Sprintf("Explanation: %s\n", example.Explanation)
			}
			content += "\n"
		}
	}

	// Add constraints
	if len(m.Problem.Constraints) > 0 {
		content += view.HeaderStyle.Render("Constraints:") + "\n\n"
		for _, constraint := range m.Problem.Constraints {
			content += fmt.Sprintf("• %s\n", constraint)
		}
		content += "\n"
	}

	// Add pattern explanation if in learn mode or hints are shown
	if m.Mode == "learn" || m.ShowHints {
		if m.Problem.PatternExplanation != "" {
			content += view.HeaderStyle.Render("Pattern Explanation:") + "\n\n"
			content += m.Problem.PatternExplanation + "\n\n"
		}
	}

	// Add solution walkthrough if in learn mode or solution is shown
	if m.Mode == "learn" || m.ShowSolution {
		if len(m.Problem.SolutionWalkthrough) > 0 {
			content += view.HeaderStyle.Render("Solution Walkthrough:") + "\n\n"
			for i, step := range m.Problem.SolutionWalkthrough {
				content += fmt.Sprintf("%d. %s\n", i+1, step)
			}
			content += "\n"
		}

		// Add solution code
		if m.Problem.Solutions != nil {
			if solution, ok := m.Problem.Solutions[m.Language]; ok {
				content += view.HeaderStyle.Render("Solution Code:") + "\n\n"
				highlightedSolution, _ := m.SyntaxHighlighter.Highlight(solution, m.Language)
				content += highlightedSolution + "\n\n"
			}
		}
	}

	// Add pattern visualization if available
	if m.CurrentPattern != "" {
		content += view.HeaderStyle.Render("Pattern Visualization:") + "\n\n"
		vizWidth := m.ProblemViewport.Width - 4
		
		// Get example data from the problem
		var exampleData string
		if len(m.Problem.Examples) > 0 {
			exampleData = m.Problem.Examples[0].Input
		}
		
		viz := m.PatternViz.VisualizePattern(m.CurrentPattern, exampleData, vizWidth)
		content += viz + "\n\n"
	}

	return content
}

// formatCodeContent formats the code editor and test results
func (m SessionModel) formatCodeContent() string {
	// Start with the code section header
	content := view.HeaderStyle.Render("Your Solution:") + "\n\n"

	// Add highlighted code
	highlightedCode, _ := m.SyntaxHighlighter.Highlight(m.Code, m.Language)
	content += highlightedCode + "\n\n"

	// Add test results if available
	if len(m.TestResults) > 0 {
		content += view.HeaderStyle.Render("Test Results:") + "\n\n"
		
		for i, result := range m.TestResults {
			if result.Passed {
				content += view.SuccessStyle.Render(fmt.Sprintf("✓ Test %d: PASSED", i+1)) + "\n"
			} else {
				content += view.ErrorStyle.Render(fmt.Sprintf("✗ Test %d: FAILED", i+1)) + "\n"
				content += fmt.Sprintf("  Input: %s\n", result.Input)
				content += fmt.Sprintf("  Expected: %s\n", result.Expected)
				content += fmt.Sprintf("  Actual: %s\n", result.Actual)
			}
			content += "\n"
		}
		
		if m.AllPassed {
			content += view.SuccessStyle.Render("All tests passed! 🎉") + "\n"
		}
	}

	return content
}

// Custom message types
type (
	editorFinishedMsg struct {
		code string
	}
	
	testResultsMsg struct {
		Results   []TestResult
		AllPassed bool
	}
)