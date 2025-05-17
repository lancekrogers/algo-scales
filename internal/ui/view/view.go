// Package view handles UI rendering and presentation
package view

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/common/highlight"
	"github.com/lancekrogers/algo-scales/internal/ui/model"
)

// View handles rendering the UI based on the model state
type View struct {
	// The UI model
	Model *model.UIModel

	// UI components
	spinner           spinner.Model
	syntaxHighlighter *highlight.SyntaxHighlighter
	patternViz        *PatternVisualization
}

// NewView creates a new view with the given model
func NewView(m *model.UIModel) *View {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498db"))

	return &View{
		Model:             m,
		spinner:           s,
		syntaxHighlighter: NewSyntaxHighlighter("monokai"),
		patternViz:        NewPatternVisualization(),
	}
}

// Render renders the UI based on the model state
func (v *View) Render() string {
	width, height := 80, 24 // Default terminal size

	// Render different screens based on app state
	switch v.Model.AppState {
	case model.StateInitial:
		return v.renderWelcomeScreen(width, height)

	case model.StateOnboarding:
		return v.renderOnboardingScreen(width, height)

	case model.StateModeSelection:
		return v.renderModeSelectionScreen(width, height)

	case model.StatePatternSelection:
		return v.renderPatternSelectionScreen(width, height)

	case model.StateProblemSelection:
		return v.renderProblemSelectionScreen(width, height)

	case model.StateSession:
		return v.renderSessionScreen(width, height)

	case model.StateStatistics:
		return v.renderStatisticsScreen(width, height)

	default:
		return "Unknown state"
	}
}

// renderWelcomeScreen renders the welcome screen
func (v *View) renderWelcomeScreen(width, height int) string {
	title := TitleStyle.Render("Algo Scales")
	subtitle := SubtitleStyle.Render("The Musician's Approach to Algorithm Mastery")

	// ASCII art logo
	logo := `
    ♩     ♫ 
  ♪ ───────── ♬
     SCALES
	`

	// Welcome message
	message := `
Welcome to Algo Scales, a terminal-based algorithm study tool
designed for efficient interview preparation.

Press ENTER to start your practice session.
Press ? for help at any time.
Press q to quit.
`

	// Combine all elements
	content := logo + "\n" + title + "\n" + subtitle + "\n" + message

	// Show loading spinner if loading
	if v.Model.Loading {
		content += "\n\n" + v.spinner.View() + " Loading..."
	}

	// Show error message if any
	if v.Model.ErrorMessage != "" {
		content += "\n\n" + ErrorStyle.Render(v.Model.ErrorMessage)
	}

	return content
}

// renderOnboardingScreen renders the interactive tutorial
func (v *View) renderOnboardingScreen(width, height int) string {
	title := TitleStyle.Render("Getting Started with Algo Scales")

	// Tutorial content
	content := `
Algo Scales organizes problems by algorithm patterns, helping you
recognize and master common interview techniques.

Each pattern is represented by a musical scale, creating a unique
learning metaphor:

- C Major (Sliding Window): The fundamental scale
- G Major (Two Pointers): Balanced and efficient
- D Major (Fast & Slow Pointers): The cycle detector
- And many more...

Press SPACE to continue or ESC to skip.
`

	return title + "\n" + content
}

// renderModeSelectionScreen renders the mode selection screen
func (v *View) renderModeSelectionScreen(width, height int) string {
	title := TitleStyle.Render("Select Learning Mode")

	// Mode options
	modes := []struct {
		name        string
		description string
	}{
		{"Learn Mode", "Focused learning with detailed explanations"},
		{"Practice Mode", "Test your skills with hints available"},
		{"Cram Mode", "Rapid-fire practice for interview prep"},
	}

	// Render each mode option
	var modeOptions strings.Builder
	for i, mode := range modes {
		option := "  "
		if i == v.Model.SelectedIndex {
			option = FocusedItemStyle.Render("▶ " + mode.name)
		} else {
			option = UnfocusedItemStyle.Render("  " + mode.name)
		}
		option += " // " + mode.description
		modeOptions.WriteString(option + "\n")
	}

	return title + "\n\n" + modeOptions.String() + "\n\nUse arrow keys to navigate, Enter to select"
}

// renderPatternSelectionScreen renders the algorithm pattern selection screen
func (v *View) renderPatternSelectionScreen(width, height int) string {
	title := TitleStyle.Render("Select Algorithm Pattern")

	// Get all patterns
	var patternOptions strings.Builder
	patternOptions.WriteString("  All Patterns\n")
	
	// Convert map to slice for consistent ordering
	patterns := make([]string, 0, len(MusicScales))
	for pattern := range MusicScales {
		patterns = append(patterns, pattern)
	}

	// Render pattern options
	for i, pattern := range patterns {
		scale := MusicScales[pattern]
		option := "  "
		
		// Highlight selected pattern
		if i+1 == v.Model.SelectedIndex {
			option = FocusedItemStyle.Render("▶ " + scale.Name)
		} else {
			// Style based on the pattern's color theme
			patternStyle, _, _ := GetPatternStyle(pattern)
			option = patternStyle.Render("  " + scale.Name)
		}
		
		// Add description
		option += " // " + scale.Description
		patternOptions.WriteString(option + "\n")
	}

	// Show progress bars for each pattern if user has statistics
	var progressSection strings.Builder
	if len(v.Model.Stats.PatternsProgress) > 0 {
		progressSection.WriteString("\nYour Pattern Mastery:\n")
		
		for pattern, progress := range v.Model.Stats.PatternsProgress {
			if scale, ok := MusicScales[pattern]; ok {
				// Only show patterns with some progress
				if progress > 0 {
					progressBar := ProgressBar(20, progress, pattern)
					name := scale.Name
					if len(name) > 20 {
						name = name[:17] + "..."
					}
					progressSection.WriteString(
						fmt.Sprintf("%-20s %s\n", name, progressBar),
					)
				}
			}
		}
	}

	return title + "\n\n" + patternOptions.String() + "\n" + progressSection.String() + 
	       "\n\nUse arrow keys to navigate, Enter to select"
}

// renderProblemSelectionScreen renders the problem selection screen
func (v *View) renderProblemSelectionScreen(width, height int) string {
	pattern := v.Model.Session.CurrentPattern
	var title string
	
	if pattern != "" {
		scale, ok := MusicScales[pattern]
		if ok {
			title = TitleStyle.Render("Select Problem - " + scale.Name)
		} else {
			title = TitleStyle.Render("Select Problem")
		}
	} else {
		title = TitleStyle.Render("Select Problem")
	}

	// Show pattern visualization if a pattern is selected
	var patternVisual string
	if pattern != "" {
		patternVisual = v.patternViz.VisualizePattern(pattern, "", width) + "\n\n"
	}

	// Problem list
	var problemList strings.Builder
	if len(v.Model.AvailableProblems) == 0 {
		problemList.WriteString("No problems available for this pattern.")
	} else {
		for i, problem := range v.Model.AvailableProblems {
			option := "  "
			
			// Style based on difficulty
			difficultyStyle := lipgloss.NewStyle()
			switch problem.Difficulty {
			case "easy":
				difficultyStyle = difficultyStyle.Foreground(lipgloss.Color("#2ecc71"))
			case "medium":
				difficultyStyle = difficultyStyle.Foreground(lipgloss.Color("#f1c40f"))
			case "hard":
				difficultyStyle = difficultyStyle.Foreground(lipgloss.Color("#e74c3c"))
			}
			
			// Highlight selected problem
			if i == v.Model.SelectedIndex {
				option = FocusedItemStyle.Render("▶ " + problem.Title)
			} else {
				option = UnfocusedItemStyle.Render("  " + problem.Title)
			}
			
			// Add difficulty and estimated time
			difficulty := difficultyStyle.Render(strings.Title(problem.Difficulty))
			timeEstimate := fmt.Sprintf("%d min", problem.EstimatedTime)
			
			problemList.WriteString(fmt.Sprintf("%-40s [%s | %s]\n", option, difficulty, timeEstimate))
		}
	}

	// Show loading spinner if loading
	var loadingIndicator string
	if v.Model.Loading {
		loadingIndicator = "\n" + v.spinner.View() + " Loading problems..."
	}

	return title + "\n\n" + patternVisual + problemList.String() + loadingIndicator +
	       "\n\nUse arrow keys to navigate, Enter to select, Backspace to return"
}

// renderSessionScreen renders the active session screen with problem and code editor
func (v *View) renderSessionScreen(width, height int) string {
	if !v.Model.Session.Active || v.Model.Session.Problem == nil {
		return "No active session"
	}

	// Get the current problem and session details
	problem := v.Model.Session.Problem
	
	// Get pattern colors for styling
	primaryStyle, secondaryStyle, _ := GetPatternStyle(v.Model.Session.CurrentPattern)
	
	// Format title with pattern information
	title := TitleStyle.Render(problem.Title)
	if v.Model.Session.CurrentPattern != "" {
		scale, ok := MusicScales[v.Model.Session.CurrentPattern]
		if ok {
			title = TitleStyle.Render(problem.Title + " - " + scale.Name)
		}
	}
	
	// Problem difficulty and metadata
	metadata := fmt.Sprintf(
		"Difficulty: %s | Estimated Time: %d min | Language: %s | Mode: %s",
		strings.Title(problem.Difficulty),
		problem.EstimatedTime,
		strings.Title(v.Model.Session.Language),
		strings.Title(v.Model.Session.Mode),
	)
	metadata = secondaryStyle.Render(metadata)
	
	// Format timer
	timer := fmt.Sprintf(
		"Time Remaining: %02d:%02d:%02d",
		int(v.Model.Session.TimeRemaining.Hours()),
		int(v.Model.Session.TimeRemaining.Minutes())%60,
		int(v.Model.Session.TimeRemaining.Seconds())%60,
	)
	timer = primaryStyle.Render(timer)
	
	// Format problem description
	description := problem.Description
	if v.Model.Session.ShowHints && problem.PatternExplanation != "" {
		description += "\n\n" + BorderedBoxStyle.Render(
			primaryStyle.Render("Pattern Explanation:") + "\n" + problem.PatternExplanation,
		)
	}
	
	// Add pattern visualization
	var patternViz string
	if v.Model.Session.CurrentPattern != "" {
		// Get example data from the problem
		var exampleData string
		if len(problem.Examples) > 0 {
			exampleData = problem.Examples[0].Input
		}
		
		patternViz = "\n" + BorderedBoxStyle.Render(
			primaryStyle.Render("Pattern Visualization:") + "\n" + 
			v.patternViz.VisualizePattern(v.Model.Session.CurrentPattern, exampleData, width/2),
		)
	}
	
	// Highlight code
	highlightedCode, err := v.syntaxHighlighter.Highlight(
		v.Model.Session.Code, 
		v.Model.Session.Language,
	)
	if err != nil {
		highlightedCode = v.Model.Session.Code
	}
	
	// Format test results
	var testResults string
	if len(v.Model.Session.TestResults) > 0 {
		var testOutput strings.Builder
		testOutput.WriteString(primaryStyle.Render("Test Results:") + "\n")
		
		allPassed := true
		for i, test := range v.Model.Session.TestResults {
			testOutput.WriteString(fmt.Sprintf("Test %d: ", i+1))
			if test.Passed {
				testOutput.WriteString(SuccessStyle.Render("✓ PASSED") + "\n")
			} else {
				testOutput.WriteString(ErrorStyle.Render("✗ FAILED") + "\n")
				testOutput.WriteString(fmt.Sprintf("  Input: %s\n", test.Input))
				testOutput.WriteString(fmt.Sprintf("  Expected: %s\n", test.Expected))
				testOutput.WriteString(fmt.Sprintf("  Got: %s\n", test.Actual))
				allPassed = false
			}
		}
		
		if allPassed {
			testOutput.WriteString("\n" + SuccessStyle.Render("All tests passed! 🎉"))
		}
		
		testResults = "\n" + BorderedBoxStyle.Render(testOutput.String())
	}
	
	// Show solution if requested
	var solution string
	if v.Model.Session.ShowSolution && problem.Solutions != nil {
		solutionCode, ok := problem.Solutions[v.Model.Session.Language]
		if !ok {
			// Try to find any solution if the current language isn't available
			for _, code := range problem.Solutions {
				solutionCode = code
				break
			}
		}
		
		// Highlight the solution code
		highlightedSolution, err := v.syntaxHighlighter.Highlight(
			solutionCode,
			v.Model.Session.Language,
		)
		if err != nil {
			highlightedSolution = solutionCode
		}
		
		solution = "\n" + BorderedBoxStyle.Render(
			primaryStyle.Render("Solution:") + "\n" + highlightedSolution,
		)
	}
	
	// Format keyboard shortcuts
	shortcuts := "e: Edit Code | h: Show Hints | s: Show Solution | Enter: Submit | q: Quit"
	shortcuts = HelpStyle.Render(shortcuts)
	
	// Layout the screen with split view
	content := title + "\n" + metadata + "\n\n" + 
		description + patternViz + "\n\n" +
		primaryStyle.Render("Your Solution:") + "\n" +
		highlightedCode + "\n" +
		testResults + solution + "\n\n" +
		timer + "\n" + shortcuts
	
	// Show loading spinner if running tests
	if v.Model.Loading {
		content += "\n\n" + v.spinner.View() + " Running tests..."
	}
	
	return content
}

// renderStatisticsScreen renders the statistics and achievements screen
func (v *View) renderStatisticsScreen(width, height int) string {
	title := TitleStyle.Render("Your Progress Statistics")
	
	// Basic stats
	stats := fmt.Sprintf(
		"Problems Solved: %d\nCurrent Streak: %d days\nLongest Streak: %d days\nTotal Practice Time: %s",
		v.Model.Stats.ProblemsSolved,
		v.Model.Stats.CurrentStreak,
		v.Model.Stats.LongestStreak,
		formatDuration(v.Model.Stats.TotalTime),
	)
	
	// Pattern progress
	var patternStats strings.Builder
	patternStats.WriteString("\n\nPattern Mastery:\n")
	
	for pattern, progress := range v.Model.Stats.PatternsProgress {
		if scale, ok := MusicScales[pattern]; ok {
			// Skip patterns with no progress
			if progress <= 0 {
				continue
			}
			
			count := v.Model.Stats.PatternCounts[pattern]
			bar := ProgressBar(20, progress, pattern)
			patternStats.WriteString(fmt.Sprintf(
				"%-20s %s %d/10\n",
				scale.Name,
				bar,
				count,
			))
		}
	}
	
	// Difficulty breakdown
	var difficultyStats strings.Builder
	difficultyStats.WriteString("\nDifficulty Distribution:\n")
	
	difficulties := []string{"easy", "medium", "hard"}
	for _, diff := range difficulties {
		count := v.Model.Stats.DifficultyCounts[diff]
		
		// Skip difficulties with no solved problems
		if count <= 0 {
			continue
		}
		
		// Select color based on difficulty
		var color lipgloss.Color
		switch diff {
		case "easy":
			color = lipgloss.Color("#2ecc71")
		case "medium":
			color = lipgloss.Color("#f1c40f")
		case "hard":
			color = lipgloss.Color("#e74c3c")
		}
		
		diffStyle := lipgloss.NewStyle().Foreground(color)
		difficultyStats.WriteString(fmt.Sprintf(
			"%s: %d problems\n",
			diffStyle.Render(strings.Title(diff)),
			count,
		))
	}
	
	// Achievements
	var achievementsSection strings.Builder
	achievementsSection.WriteString("\nAchievements:\n")
	
	for _, achievement := range v.Model.Achievements {
		status := "□"
		if achievement.Earned {
			status = "✓"
		}
		
		achievementsSection.WriteString(fmt.Sprintf(
			"%s %s: %s\n",
			status,
			achievement.Title,
			achievement.Description,
		))
	}
	
	// Combine all sections
	content := title + "\n\n" + stats + patternStats.String() + difficultyStats.String() + achievementsSection.String()
	content += "\n\nPress q to return to main menu"
	
	return content
}

// Helper function to format duration in a human-readable format
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}