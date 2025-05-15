// Package splitscreen implements a split-screen terminal UI for AlgoScales
package splitscreen

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Model represents the main application model for the split-screen UI
type Model struct {
	// Window dimensions
	windowWidth  int
	windowHeight int

	// Panel components
	problemView  viewport.Model  // Left panel: Problem description
	codeEditor   textarea.Model  // Right panel: Code editor
	terminal     viewport.Model  // Bottom panel: Command output
	terminalInput textinput.Model // Bottom panel: Command input
	
	// Application state
	focusedPanel    focusedPanel
	codeLanguage    string
	theme           ScaleTheme
	styles          map[string]lipgloss.Style
	elapsedTime     time.Duration
	startTime       time.Time
	runningCommand  bool
	showHelp        bool
	ready           bool
	
	// Current problem
	currentProblem *problem.Problem
	
	// Vim mode (for code editor)
	vimMode VimMode
}

// focusedPanel represents which panel currently has focus
type focusedPanel int

const (
	problemPanel focusedPanel = iota
	codePanel
	terminalPanel
)

// VimMode represents the current vim editing mode
type VimMode int

const (
	NormalMode VimMode = iota
	InsertMode
	VisualMode
)

// NewModel creates a new model with initialized components
func NewModel() Model {
	// Initialize with empty components
	// They will be properly set up once we know the terminal dimensions
	
	// Set default theme
	defaultTheme := MajorTheme
	
	return Model{
		// Default values
		focusedPanel: codePanel, // Start with focus on code editor
		codeLanguage: "go",      // Default language
		theme:        defaultTheme,
		styles:       ThemeStyles(defaultTheme),
		vimMode:      InsertMode,
		showHelp:     false,
		ready:        false,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	// Return a command that will be called after Init
	return tea.Batch(
		waitForActivity(time.Second),
	)
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Initialize components with proper dimensions
		m = m.updateWindowSize(msg.Width, msg.Height)
		m.ready = true
		return m, nil
		
	case tea.KeyMsg:
		// Handle global key presses first
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
			
		case "tab":
			// Cycle focus between panels
			m.focusedPanel = (m.focusedPanel + 1) % 3
			return m, nil
			
		case "shift+tab":
			// Reverse cycle focus between panels
			if m.focusedPanel == 0 {
				m.focusedPanel = 2
			} else {
				m.focusedPanel--
			}
			return m, nil
			
		case "ctrl+s":
			// Switch language
			switch m.codeLanguage {
			case "go":
				m.codeLanguage = "python"
			case "python":
				m.codeLanguage = "javascript"
			case "javascript":
				m.codeLanguage = "go"
			}
			return m, nil
			
		case "?":
			// Toggle help
			m.showHelp = !m.showHelp
			return m, nil
		}
		
		// Route key messages to the focused panel
		switch m.focusedPanel {
		case problemPanel:
			// Handle problem view navigation
			switch msg.String() {
			case "up", "k":
				m.problemView.LineUp(1)
			case "down", "j":
				m.problemView.LineDown(1)
			case "pgup":
				m.problemView.HalfViewUp()
			case "pgdown":
				m.problemView.HalfViewDown()
			}
			
		case codePanel:
			// Update code editor
			var cmd tea.Cmd
			m.codeEditor, cmd = m.codeEditor.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			
		case terminalPanel:
			// Handle terminal input
			switch msg.String() {
			case "enter":
				// Execute command from input
				command := m.terminalInput.Value()
				m.terminalInput.Reset()
				cmds = append(cmds, runCommand(command, m.codeEditor.Value()))
				m.runningCommand = true
			default:
				// Update terminal input
				var cmd tea.Cmd
				m.terminalInput, cmd = m.terminalInput.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
		
	case execResultMsg:
		// Process command execution results
		m.runningCommand = false
		m.terminal.SetContent(m.terminal.View() + "\n$ " + msg.command + "\n" + msg.output)
		m.terminal.GotoBottom()
		
	case statusTickMsg:
		// Update elapsed time
		m.elapsedTime = time.Since(m.startTime)
		cmds = append(cmds, waitForActivity(time.Second))
	}

	// Return the updated model and commands
	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Calculate panel dimensions
	leftPanelWidth := m.windowWidth / 2
	rightPanelWidth := m.windowWidth - leftPanelWidth
	topSectionHeight := m.windowHeight - 10 // Bottom panel is 10 rows high
	
	// Create styles based on focus state
	problemPanelStyle := lipgloss.NewStyle().
		Width(leftPanelWidth).
		Height(topSectionHeight).
		BorderStyle(lipgloss.RoundedBorder())

	codePanelStyle := lipgloss.NewStyle().
		Width(rightPanelWidth).
		Height(topSectionHeight).
		BorderStyle(lipgloss.RoundedBorder())

	bottomPanelStyle := lipgloss.NewStyle().
		Width(m.windowWidth).
		Height(9).
		BorderStyle(lipgloss.RoundedBorder())
	
	// Update border colors based on focus
	switch m.focusedPanel {
	case problemPanel:
		problemPanelStyle = problemPanelStyle.
			BorderForeground(lipgloss.Color(m.theme.BrightColor))
	case codePanel:
		codePanelStyle = codePanelStyle.
			BorderForeground(lipgloss.Color(m.theme.BrightColor))
	case terminalPanel:
		bottomPanelStyle = bottomPanelStyle.
			BorderForeground(lipgloss.Color(m.theme.BrightColor))
	}
	
	// Set default border colors
	if m.focusedPanel != problemPanel {
		problemPanelStyle = problemPanelStyle.
			BorderForeground(lipgloss.Color(m.theme.MutedColor))
	}
	
	if m.focusedPanel != codePanel {
		codePanelStyle = codePanelStyle.
			BorderForeground(lipgloss.Color(m.theme.MutedColor))
	}
	
	if m.focusedPanel != terminalPanel {
		bottomPanelStyle = bottomPanelStyle.
			BorderForeground(lipgloss.Color(m.theme.MutedColor))
	}

	// Panel titles
	codeTitle := " Code Editor (" + m.codeLanguage + ") "
	
	// Add vim mode indicator to code editor title if in code panel
	if m.focusedPanel == codePanel {
		var modeText string
		switch m.vimMode {
		case NormalMode:
			modeText = " [NORMAL] "
		case InsertMode:
			modeText = " [INSERT] "
		case VisualMode:
			modeText = " [VISUAL] "
		}
		codeTitle += modeText
	}
	
	// Apply title styles with border titles
	problemPanelStyle = problemPanelStyle.
		BorderTop(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.theme.MutedColor))
	if m.focusedPanel == problemPanel {
		problemPanelStyle = problemPanelStyle.BorderForeground(lipgloss.Color(m.theme.BrightColor))
	}
	
	codePanelStyle = codePanelStyle.
		BorderTop(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.theme.MutedColor))
	if m.focusedPanel == codePanel {
		codePanelStyle = codePanelStyle.BorderForeground(lipgloss.Color(m.theme.BrightColor))
	}
	
	bottomPanelStyle = bottomPanelStyle.
		BorderTop(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.theme.MutedColor))
	if m.focusedPanel == terminalPanel {
		bottomPanelStyle = bottomPanelStyle.BorderForeground(lipgloss.Color(m.theme.BrightColor))
	}

	// Render panel content
	leftPanelRendered := problemPanelStyle.Render(m.problemView.View())
	rightPanelRendered := codePanelStyle.Render(m.codeEditor.View())
	
	// Combine terminal viewport and input for bottom panel
	terminalContent := m.terminal.View() + "\n\n> " + m.terminalInput.View()
	bottomPanelRendered := bottomPanelStyle.Render(terminalContent)

	// Format status bar
	hours := int(m.elapsedTime.Hours())
	minutes := int(m.elapsedTime.Minutes()) % 60
	seconds := int(m.elapsedTime.Seconds()) % 60
	timeStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.theme.BrightColor)).
		Render(
			lipgloss.NewStyle().Bold(true).Render("Time:") + 
			lipgloss.NewStyle().Render(
				lipgloss.NewStyle().Foreground(lipgloss.Color("#f8e71c")).
				Render(
					lipgloss.NewStyle().Bold(true).
					Render(
						lipgloss.NewStyle().Italic(true).
						Render(
							lipgloss.NewStyle().Underline(true).
							Render(
								lipgloss.NewStyle().Faint(false).
								Render(
									lipgloss.NewStyle().Render(
										lipgloss.NewStyle().Render(
											lipgloss.NewStyle().Render(
												lipgloss.NewStyle().Render(
													lipgloss.NewStyle().Render(
														lipgloss.NewStyle().
														Render(
															lipgloss.NewStyle().
															Render(
																lipgloss.NewStyle().
																Render(
																	lipgloss.NewStyle().
																	Render(
																		lipgloss.NewStyle().
																		Render(
																			lipgloss.NewStyle().
																			Render(
																				lipgloss.NewStyle().
																				Render(
																					lipgloss.NewStyle().
																					Render(
																						lipgloss.NewStyle().
																						Render(
																							lipgloss.NewStyle().
																							Render(
																								lipgloss.NewStyle().
																								Render(
																									lipgloss.NewStyle().
																									Render(
																										lipgloss.NewStyle().
																										Render(
																											lipgloss.NewStyle().
																											Render(
																												lipgloss.NewStyle().
																												Render(
																													lipgloss.NewStyle().
																													Render(
																														lipgloss.NewStyle().
																														Render(
																															lipgloss.NewStyle().
																															Render(
																																lipgloss.NewStyle().
																																Render(
																																	lipgloss.NewStyle().
																																	Render(fmt.Sprintf(" %02d:%02d:%02d", hours, minutes, seconds)),
																																),
																															),
																														),
																													),
																												),
																											),
																										),
																									),
																								),
																							),
																						),
																					),
																				),
																			),
																		),
																	),
																),
															),
														),
													),
												),
											),
										),
									),
								),
							),
						),
					),
				),
			),
		)
	
	// Format language indicator
	languageStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.theme.AccentColor)).
		Render(
			lipgloss.NewStyle().Bold(true).Render("Language:") + " " + 
			lipgloss.NewStyle().Italic(true).Render(m.codeLanguage),
		)
	
	// Format key bindings
	keybindingsStr := "Tab: Switch Panel | Ctrl+S: Switch Language | ?: Toggle Help | Ctrl+C: Quit"
	if m.showHelp {
		keybindingsStr = "k/j: Scroll Up/Down | Ctrl+R: Run Code | Esc: Exit Help | Tab: Switch Panel"
	}
	
	helpStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.theme.MutedColor)).
		Render(keybindingsStr)
	
	// Create status bar
	statusBarStyle := lipgloss.NewStyle().
		Width(m.windowWidth).
		Padding(0, 1).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(m.theme.BaseColor))
	
	// Format status bar content with proper spacing
	leftStatus := timeStr
	rightStatus := languageStr + " | " + helpStr
	
	// Calculate padding needed between left and right status elements
	padding := m.windowWidth - lipgloss.Width(leftStatus) - lipgloss.Width(rightStatus) - 2
	if padding < 0 {
		padding = 0
	}
	
	statusContent := leftStatus + strings.Repeat(" ", padding) + rightStatus
	statusBar := statusBarStyle.Render(statusContent)

	// Join horizontal panels (left and right)
	topSection := lipgloss.JoinHorizontal(lipgloss.Top, leftPanelRendered, rightPanelRendered)
	
	// Join vertical sections (top, bottom, and status bar)
	return lipgloss.JoinVertical(lipgloss.Left, topSection, bottomPanelRendered, statusBar)
}

// updateWindowSize updates the window dimensions and adjusts all components accordingly
func (m Model) updateWindowSize(width, height int) Model {
	m.windowWidth = width
	m.windowHeight = height
	
	// Calculate panel dimensions
	leftPanelWidth := width / 2
	rightPanelWidth := width - leftPanelWidth
	topSectionHeight := height - 10 // Bottom panel is 10 rows high
	
	// Adjust problem view
	m.problemView = viewport.New(leftPanelWidth-4, topSectionHeight-2) // Adjust for border and padding
	m.problemView.SetContent("Loading problem description...")
	
	// Adjust code editor
	m.codeEditor = textarea.New()
	m.codeEditor.SetWidth(rightPanelWidth - 4)
	m.codeEditor.SetHeight(topSectionHeight - 2)
	m.codeEditor.ShowLineNumbers = true
	m.codeEditor.Placeholder = "// Write your code here"
	
	// Adjust terminal
	m.terminal = viewport.New(width-4, 6) // Adjust for border and padding
	m.terminal.SetContent("Welcome to AlgoScales Terminal\nType commands here and press Enter to execute.\n")
	
	// Adjust terminal input
	m.terminalInput = textinput.New()
	m.terminalInput.Width = width - 6
	m.terminalInput.Placeholder = "Type command here"
	
	// Only focus the input if terminal panel is focused
	if m.focusedPanel == terminalPanel {
		m.terminalInput.Focus()
	}
	
	return m
}

// SetProblem sets the current problem and updates the problem view
func (m *Model) SetProblem(p *problem.Problem) {
	m.currentProblem = p
	
	// Format the problem description
	description := fmt.Sprintf("# %s\n\n", p.Title)
	description += fmt.Sprintf("**Difficulty**: %s\n\n", p.Difficulty)
	description += p.Description + "\n\n"
	
	// Add examples
	if len(p.Examples) > 0 {
		description += "## Examples\n\n"
		for i, example := range p.Examples {
			description += fmt.Sprintf("### Example %d:\n\n", i+1)
			description += fmt.Sprintf("**Input**: %s\n\n", example.Input)
			description += fmt.Sprintf("**Output**: %s\n\n", example.Output)
			if example.Explanation != "" {
				description += fmt.Sprintf("**Explanation**: %s\n\n", example.Explanation)
			}
		}
	}
	
	// Add constraints
	if len(p.Constraints) > 0 {
		description += "## Constraints\n\n"
		for _, constraint := range p.Constraints {
			description += "- " + constraint + "\n"
		}
	}
	
	// Update the problem view with the formatted description
	m.problemView.SetContent(description)
	m.problemView.GotoTop()
}

// SetFocus sets the focus to the specified panel
func (m *Model) SetFocus(panel focusedPanel) {
	m.focusedPanel = panel
	
	// Update focus state of relevant components
	switch panel {
	case terminalPanel:
		m.terminalInput.Focus()
	case codePanel:
		m.codeEditor.Focus()
	}
}

// Custom message types
type (
	// statusTickMsg is sent every second to update the timer
	statusTickMsg struct{}
	
	// execResultMsg is sent when a command execution is complete
	execResultMsg struct {
		command string
		output  string
		err     error
	}
)

// waitForActivity returns a command that sends a statusTickMsg after the specified duration
func waitForActivity(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return statusTickMsg{}
	})
}

// runCommand executes a command and returns the result
func runCommand(command string, input string) tea.Cmd {
	return func() tea.Msg {
		// Here we would implement actual command execution
		// For now, just echo the command and input
		output := "Command execution not implemented yet.\n"
		output += "Command: " + command + "\n"
		output += "With input from editor:\n"
		output += input
		
		return execResultMsg{
			command: command,
			output:  output,
			err:     nil,
		}
	}
}