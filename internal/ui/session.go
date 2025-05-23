package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/common/config"
	"github.com/lancekrogers/algo-scales/internal/common/logging"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Update handles updates for the session screen
func (m Model) updateSession(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Initialize or update viewport
		if m.session.viewport.Width == 0 {
			m.session.viewport = viewport.New(msg.Width-4, msg.Height-10)
			m.session.viewport.SetContent(m.sessionContent())
		} else {
			m.session.viewport.Width = msg.Width - 4
			m.session.viewport.Height = msg.Height - 10
		}
		
	case sessionTickMsg:
		// Update timer
		if !m.session.timerPaused {
			m.session.duration = time.Since(m.session.startTime)
		}
		return m, sessionTick()
		
	case sessionStartedMsg:
		m.session.sessionID = msg.sessionID
		m.session.startTime = time.Now()
		return m, sessionTick()
		
	case testResultsMsg:
		m.session.testResults = msg.results
		m.session.viewport.SetContent(m.sessionContent())
		
	case editorFinishedMsg:
		m.session.message = "Editor closed. Press 't' to run tests."
		return m, nil
		
	case editorErrorMsg:
		m.session.message = fmt.Sprintf("Error opening editor: %v", msg.error)
		return m, nil
		
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			// Open editor
			return m, openEditor(m.session.sessionID, m.config.Language, m.session.problem)
		case "t":
			// Run tests
			return m, runTests(m.session.sessionID, m.config.Language)
		case "h":
			// Toggle hint
			m.session.showHint = !m.session.showHint
			m.session.viewport.SetContent(m.sessionContent())
		case "s":
			// Toggle solution
			m.session.showSolution = !m.session.showSolution
			m.session.viewport.SetContent(m.sessionContent())
		case "p":
			// Pause/unpause timer
			m.session.timerPaused = !m.session.timerPaused
		case "enter":
			// Submit solution
			return m.submitSolution()
		case "ctrl+c", "q":
			// Confirmation before quitting
			if m.session.confirmQuit {
				return m.navigate(StateHome), nil
			}
			m.session.confirmQuit = true
			return m, nil
		default:
			// Pass through to viewport
			m.session.viewport, cmd = m.session.viewport.Update(msg)
		}
	}
	
	return m, cmd
}

// View renders the session screen
func (m Model) viewSession() string {
	var b strings.Builder
	
	// Header with problem title and timer
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62"))
	
	timerStyle := lipgloss.NewStyle().
		Bold(true)
	
	if m.session.duration > 30*time.Minute {
		timerStyle = timerStyle.Foreground(lipgloss.Color("196")) // Red
	} else if m.session.duration > 20*time.Minute {
		timerStyle = timerStyle.Foreground(lipgloss.Color("214")) // Orange
	} else {
		timerStyle = timerStyle.Foreground(lipgloss.Color("46")) // Green
	}
	
	pauseIndicator := ""
	if m.session.timerPaused {
		pauseIndicator = " (PAUSED)"
	}
	
	header := headerStyle.Render(m.session.problem.Title)
	timer := timerStyle.Render(formatDuration(m.session.duration) + pauseIndicator)
	
	headerBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		header,
		strings.Repeat(" ", max(0, m.width-lipgloss.Width(header)-lipgloss.Width(timer))),
		timer,
	)
	
	b.WriteString(headerBar)
	b.WriteString("\n\n")
	
	// Viewport with session content
	b.WriteString(m.session.viewport.View())
	b.WriteString("\n\n")
	
	// Message or confirmation
	if m.session.confirmQuit {
		confirmStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))
		b.WriteString(confirmStyle.Render("Really quit? Press q or ctrl+c again to confirm, any other key to cancel."))
		b.WriteString("\n")
	} else if m.session.message != "" {
		msgStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
		b.WriteString(msgStyle.Render(m.session.message))
		b.WriteString("\n")
	}
	
	// Action bar
	actionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	
	actions := []string{
		"e: Edit Code",
		"t: Run Tests",
		"h: Toggle Hint",
		"s: Show Solution",
		"p: Pause Timer",
		"Enter: Submit",
		"Esc: Back",
	}
	
	b.WriteString(actionStyle.Render(strings.Join(actions, " • ")))
	
	return b.String()
}

// sessionContent generates the content for the session viewport
func (m Model) sessionContent() string {
	var content strings.Builder
	p := m.session.problem
	
	// Problem description
	content.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		Render("Problem"))
	content.WriteString("\n\n")
	if p.Description != "" {
		content.WriteString(p.Description)
	} else {
		content.WriteString("No description available")
	}
	content.WriteString("\n\n")
	
	// Examples
	if len(p.Examples) > 0 {
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Render("Examples"))
		content.WriteString("\n\n")
		
		for i, example := range p.Examples {
			content.WriteString(fmt.Sprintf("Example %d:\n", i+1))
			codeStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Background(lipgloss.Color("235")).
				Padding(0, 1)
			
			content.WriteString("Input: ")
			content.WriteString(codeStyle.Render(example.Input))
			content.WriteString("\n")
			
			content.WriteString("Output: ")
			content.WriteString(codeStyle.Render(example.Output))
			content.WriteString("\n")
			
			if example.Explanation != "" {
				content.WriteString("Explanation: " + example.Explanation + "\n")
			}
			content.WriteString("\n")
		}
	}
	
	// Test results
	if m.session.testResults != "" {
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Render("Test Results"))
		content.WriteString("\n\n")
		content.WriteString(m.session.testResults)
		content.WriteString("\n\n")
	}
	
	// Pattern Explanation
	if m.session.showHint && p.PatternExplanation != "" {
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214")).
			Render("💡 Pattern Explanation"))
		content.WriteString("\n\n")
		content.WriteString(p.PatternExplanation)
		content.WriteString("\n\n")
	}
	
	// Solution
	if m.session.showSolution && len(p.SolutionWalkthrough) > 0 {
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			Render("✅ Solution Walkthrough"))
		content.WriteString("\n\n")
		for i, step := range p.SolutionWalkthrough {
			content.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
		}
		content.WriteString("\n")
	}
	
	return content.String()
}

// openEditor opens the code file in the user's editor
func openEditor(sessionID, language string, problem problem.Problem) tea.Cmd {
	return func() tea.Msg {
		// Get the session directory
		sessionDir := fmt.Sprintf("/tmp/algo-scales/sessions/%s", sessionID)
		codeFile := fmt.Sprintf("%s/solution.%s", sessionDir, getFileExtension(language))
		
		// Create the file if it doesn't exist
		if _, err := os.Stat(codeFile); os.IsNotExist(err) {
			os.MkdirAll(sessionDir, 0755)
			// Write starter code
			starterCode := problem.StarterCode[language]
			if starterCode == "" {
				// Provide a basic template if no starter code
				starterCode = getDefaultTemplate(language, problem)
			}
			os.WriteFile(codeFile, []byte(starterCode), 0644)
		}
		
		// Get editor from config or environment
		cfg, _ := config.LoadConfig()
		editor := cfg.EditorCommand
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		if editor == "" {
			if runtime.GOOS == "windows" {
				editor = "notepad"
			} else {
				editor = "vim"
			}
		}
		
		// Create session state for error logging
		sessionState := &logging.SessionSnapshot{
			ProblemID:    problem.ID,
			Language:     language,
			Mode:         "editor_session",
			StartTime:    time.Now(),
			Patterns:     problem.Patterns,
			Difficulty:   problem.Difficulty,
			Workspace:    sessionDir,
			CodeFile:     codeFile,
			CustomFields: map[string]string{
				"editor": editor,
			},
		}
		
		// Open editor
		cmd := exec.Command(editor, codeFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		// Log editor operation start
		ctx := context.Background()
		ctx = logging.WithOperation(ctx, "open_editor")
		ctx = logging.WithComponent(ctx, "UI")
		logger := logging.NewLogger("EditorSession").WithContext(ctx)
		
		logger.Info("Opening editor: %s for file: %s", editor, codeFile)
		
		err := cmd.Run()
		if err != nil {
			// Log detailed editor error
			if logging.GlobalErrorLogger != nil {
				logging.GlobalErrorLogger.LogEditorError(ctx, err, editor, codeFile, sessionState)
			}
			logger.Error("Editor failed: %v", err)
			return editorErrorMsg{err}
		}
		
		logger.Info("Editor session completed successfully")
		return editorFinishedMsg{}
	}
}

// runTests runs tests on the current solution
func runTests(sessionID, language string) tea.Cmd {
	return func() tea.Msg {
		// Get the session directory and code file
		sessionDir := fmt.Sprintf("/tmp/algo-scales/sessions/%s", sessionID)
		codeFile := fmt.Sprintf("%s/solution.%s", sessionDir, getFileExtension(language))
		
		// Check if file exists
		if _, err := os.Stat(codeFile); os.IsNotExist(err) {
			return testResultsMsg{results: "Error: No solution file found. Press 'e' to edit your solution first."}
		}
		
		// Simulate test run for now
		time.Sleep(1 * time.Second)
		
		results := "Running tests...\n\n"
		results += "✅ Test 1: PASSED\n"
		results += "✅ Test 2: PASSED\n"
		results += "❌ Test 3: FAILED\n"
		results += "   Expected: [1, 2, 3]\n"
		results += "   Got: [1, 3, 2]\n\n"
		results += "2/3 tests passed"
		
		return testResultsMsg{results: results}
	}
}

// submitSolution handles solution submission
func (m Model) submitSolution() (Model, tea.Cmd) {
	// Save session stats
	duration := m.session.duration
	
	// Simple completion check
	completed := strings.Contains(m.session.testResults, "tests passed") &&
		!strings.Contains(m.session.testResults, "FAILED")
	
	// Create completion message
	msg := fmt.Sprintf("Session completed in %s", formatDuration(duration))
	if completed {
		msg += " - All tests passed! 🎉"
	} else {
		msg += " - Some tests failed"
	}
	
	m.session.message = msg
	
	// Return to problem list after a delay
	return m, tea.Sequence(
		tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return navigateBackMsg{}
		}),
	)
}

// Custom message types for session
type editorFinishedMsg struct{}
type editorErrorMsg struct{ error }
type testResultsMsg struct{ results string }

// Helper to get file extension
func getFileExtension(language string) string {
	switch language {
	case "python":
		return "py"
	case "javascript":
		return "js"
	case "typescript":
		return "ts"
	case "java":
		return "java"
	case "cpp":
		return "cpp"
	case "rust":
		return "rs"
	case "go":
		return "go"
	default:
		return "txt"
	}
}

// Helper to get max of two ints
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}


// getDefaultTemplate provides a basic template when no starter code is available
func getDefaultTemplate(language string, problem problem.Problem) string {
	switch language {
	case "python":
		return fmt.Sprintf(`# %s
# %s

def solution():
    # TODO: Implement your solution here
    pass

if __name__ == "__main__":
    # Test your solution here
    pass
`, problem.Title, problem.Description)
	case "javascript":
		return fmt.Sprintf(`// %s
// %s

function solution() {
    // TODO: Implement your solution here
}

// Test your solution here
`, problem.Title, problem.Description)
	case "go":
		return fmt.Sprintf(`// %s
// %s

package main

import "fmt"

func solution() {
    // TODO: Implement your solution here
}

func main() {
    // Test your solution here
    fmt.Println("Solution not implemented yet")
}
`, problem.Title, problem.Description)
	default:
		return fmt.Sprintf(`// %s
// %s

// TODO: Implement your solution here
`, problem.Title, problem.Description)
	}
}