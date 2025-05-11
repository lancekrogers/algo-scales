// Package controller handles user interactions and app logic
package controller

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/ui/model"
	"github.com/lancekrogers/algo-scales/internal/ui/view"
)

// Controller handles interactions between the model and view
type Controller struct {
	// The UI model
	Model *model.UIModel

	// Visualization components
	syntaxHighlighter *view.SyntaxHighlighter
	spinners          view.CustomSpinners
	patternViz        *view.PatternVisualization

	// Problem service
	problemService *problem.Service
}

// NewController creates a new controller with the given model
func NewController(m *model.UIModel) *Controller {
	return &Controller{
		Model:             m,
		syntaxHighlighter: view.NewSyntaxHighlighter("monokai"),
		spinners:          view.NewCustomSpinners(),
		patternViz:        view.NewPatternVisualization(),
		problemService:    problem.NewService(),
	}
}

// Init initializes the controller and loads initial data
func (c *Controller) Init() tea.Cmd {
	// Return a command that loads the problem list
	return c.loadProblems
}

// loadProblems loads the available problems
func (c *Controller) loadProblems() tea.Msg {
	// Create a loading message
	c.Model.Loading = true

	// Load problems from the problem service
	problems, err := c.problemService.ListAll()
	if err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to load problems: %v", err))
	}

	// Return a message with the loaded problems
	return model.ProblemsLoadedMsg{Problems: problems}
}

// Update handles model updates based on messages
func (c *Controller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key presses
		cmd = c.handleKeyPress(msg)

	case model.ProblemsLoadedMsg:
		// Update model with loaded problems
		c.Model.AvailableProblems = msg.Problems
		c.Model.Loading = false

	case model.ErrorMsg:
		// Handle error messages
		c.Model.ErrorMessage = string(msg)
		c.Model.Loading = false

	case model.TickMsg:
		// Handle timer ticks
		if c.Model.Session.Active && c.Model.Session.TimeRemaining > 0 {
			c.Model.Session.TimeRemaining -= time.Second
			return c.Model, c.tickTimer()
		}

	case model.ProblemSelectedMsg:
		// Handle problem selection
		cmd = c.startSession(msg.ProblemID, msg.Mode)

	case model.AchievementUnlockedMsg:
		// Handle unlocked achievements
		achievement, exists := c.Model.Achievements[msg.AchievementID]
		if exists && !achievement.Earned {
			achievement.Earned = true
			achievement.EarnedDate = time.Now()
			c.Model.Achievements[msg.AchievementID] = achievement
		}
	}

	return c.Model, cmd
}

// handleKeyPress processes keyboard input
func (c *Controller) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		// Exit the application
		return tea.Quit

	case "?":
		// Toggle help
		c.Model.ShowHelp = !c.Model.ShowHelp
		return nil

	case "h":
		// Show hints when in a session
		if c.Model.Session.Active {
			c.Model.Session.ShowHints = true
		}
		return nil

	case "s":
		// Show solution when in a session
		if c.Model.Session.Active {
			c.Model.Session.ShowSolution = true
		}
		return nil

	case "e":
		// Open code in editor
		if c.Model.Session.Active {
			return c.openEditor
		}

	case "enter":
		// Submit solution or select menu item
		if c.Model.Session.Active {
			return c.submitSolution
		} else {
			return c.handleSelection
		}

	case "up", "k":
		// Navigate menu up
		if c.Model.SelectedIndex > 0 {
			c.Model.SelectedIndex--
		}
		return nil

	case "down", "j":
		// Navigate menu down
		if c.Model.AppState == model.StateProblemSelection && 
		   c.Model.SelectedIndex < len(c.Model.AvailableProblems)-1 {
			c.Model.SelectedIndex++
		} else if c.Model.AppState == model.StateModeSelection && 
				 c.Model.SelectedIndex < 2 { // 3 modes
			c.Model.SelectedIndex++
		} else if c.Model.AppState == model.StatePatternSelection && 
				 c.Model.SelectedIndex < len(view.MusicScales) {
			c.Model.SelectedIndex++
		}
		return nil
	}

	return nil
}

// handleSelection processes selection in menus
func (c *Controller) handleSelection() tea.Msg {
	switch c.Model.AppState {
	case model.StateModeSelection:
		// Select mode
		mode := ""
		switch c.Model.SelectedIndex {
		case 0:
			mode = "learn"
		case 1:
			mode = "practice"
		case 2:
			mode = "cram"
		}

		// Move to pattern selection
		c.Model.AppState = model.StatePatternSelection
		c.Model.SelectedIndex = 0 // Reset selection
		c.Model.Session.Mode = mode
		return nil

	case model.StatePatternSelection:
		// Get all pattern keys
		patterns := make([]string, 0, len(view.MusicScales))
		for pattern := range view.MusicScales {
			patterns = append(patterns, pattern)
		}

		if c.Model.SelectedIndex == 0 {
			// All patterns selected, go to problem selection with all problems
			c.Model.AppState = model.StateProblemSelection
			c.Model.SelectedIndex = 0
			return nil
		} else if c.Model.SelectedIndex <= len(patterns) {
			// Specific pattern selected
			pattern := patterns[c.Model.SelectedIndex-1]
			c.Model.Session.CurrentPattern = pattern

			// Filter problems by pattern
			c.filterProblemsByPattern(pattern)
			
			// Move to problem selection
			c.Model.AppState = model.StateProblemSelection
			c.Model.SelectedIndex = 0
			return nil
		}

	case model.StateProblemSelection:
		// Select a problem
		if c.Model.SelectedIndex < len(c.Model.AvailableProblems) {
			problem := c.Model.AvailableProblems[c.Model.SelectedIndex]
			return model.ProblemSelectedMsg{
				ProblemID: problem.ID, 
				Mode:      c.Model.Session.Mode,
			}
		}
	}

	return nil
}

// filterProblemsByPattern filters the available problems by pattern
func (c *Controller) filterProblemsByPattern(pattern string) {
	// Filter the problems by pattern
	var filtered []problem.Problem
	for _, p := range c.Model.AvailableProblems {
		for _, patt := range p.Patterns {
			if patt == pattern {
				filtered = append(filtered, p)
				break
			}
		}
	}
	
	// Update the model with filtered problems
	if len(filtered) > 0 {
		c.Model.AvailableProblems = filtered
	}
}

// startSession begins a new problem session
func (c *Controller) startSession(problemID, mode string) tea.Cmd {
	// Load the problem
	p, err := c.problemService.GetByID(problemID)
	if err != nil {
		return tea.Batch(
			func() tea.Msg {
				return model.ErrorMsg(fmt.Sprintf("Failed to load problem: %v", err))
			},
		)
	}

	// Initialize session
	c.Model.Session = model.Session{
		Active:        true,
		Mode:          mode,
		Problem:       p,
		StartTime:     time.Now(),
		TimeRemaining: time.Duration(p.EstimatedTime) * time.Minute,
		ShowHints:     mode == "learn",
		ShowSolution:  false,
		Language:      c.Model.Session.Language,
	}

	if c.Model.Session.Language == "" {
		c.Model.Session.Language = "go" // Default language
	}

	// Get starter code for the language
	if starterCode, ok := p.StarterCode[c.Model.Session.Language]; ok {
		c.Model.Session.Code = starterCode
	} else {
		// Try to find any starter code if the preferred language isn't available
		for lang, code := range p.StarterCode {
			c.Model.Session.Language = lang
			c.Model.Session.Code = code
			break
		}
	}

	// Start the session timer
	return c.tickTimer
}

// tickTimer handles the session timer
func (c *Controller) tickTimer() tea.Msg {
	time.Sleep(time.Second)
	return model.TickMsg{}
}

// openEditor opens the code in an external editor
func (c *Controller) openEditor() tea.Msg {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "algo-scales-*."+view.GetLanguageExtension(c.Model.Session.Language))
	if err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to create temp file: %v", err))
	}
	defer os.Remove(tmpfile.Name())

	// Write the current code to the file
	if _, err := tmpfile.Write([]byte(c.Model.Session.Code)); err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to write to temp file: %v", err))
	}
	if err := tmpfile.Close(); err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to close temp file: %v", err))
	}

	// Determine which editor to use
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default to vim if EDITOR is not set
	}

	// Open the file in the editor
	cmd := tea.ExecProcess(tea.ShellCommand(editor + " " + tmpfile.Name()), func(err error) tea.Msg {
		if err != nil {
			return model.ErrorMsg(fmt.Sprintf("Editor exited with error: %v", err))
		}

		// Read the updated file
		content, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			return model.ErrorMsg(fmt.Sprintf("Failed to read temp file: %v", err))
		}

		// Update the code in the model
		c.Model.Session.Code = string(content)
		return model.CodeUpdatedMsg{}
	})

	return cmd
}

// submitSolution submits the current solution for testing
func (c *Controller) submitSolution() tea.Msg {
	// Validate that we're in an active session
	if !c.Model.Session.Active || c.Model.Session.Problem == nil {
		return model.ErrorMsg("No active session")
	}

	// Create a loading message
	c.Model.Loading = true

	// Execute the tests
	results, err := c.problemService.TestSolution(
		c.Model.Session.Problem.ID,
		c.Model.Session.Code,
		c.Model.Session.Language,
	)

	// Update model with test results
	c.Model.Loading = false
	if err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to test solution: %v", err))
	}

	// Check if all tests passed
	allPassed := true
	for _, result := range results {
		if !result.Passed {
			allPassed = false
			break
		}
	}

	// Update stats if all tests passed
	if allPassed {
		c.updateStatistics()
		c.checkAchievements()
	}

	// Update session with test results
	c.Model.Session.TestResults = results
	return model.TestResultsMsg{Results: results, AllPassed: allPassed}
}

// updateStatistics updates user statistics after solving a problem
func (c *Controller) updateStatistics() {
	// Increment problems solved count
	c.Model.Stats.ProblemsSolved++

	// Update time spent
	elapsedTime := time.Since(c.Model.Session.StartTime)
	c.Model.Stats.TotalTime += elapsedTime

	// Update pattern counts
	for _, pattern := range c.Model.Session.Problem.Patterns {
		c.Model.Stats.PatternCounts[pattern]++
		
		// Calculate pattern progress (10 problems = 100%)
		progress := float64(c.Model.Stats.PatternCounts[pattern]) / 10.0
		if progress > 1.0 {
			progress = 1.0
		}
		c.Model.Stats.PatternsProgress[pattern] = progress
	}

	// Update difficulty counts
	difficulty := c.Model.Session.Problem.Difficulty
	c.Model.Stats.DifficultyCounts[difficulty]++

	// Update streak
	today := time.Now().Format("2006-01-02")
	lastPractice := c.Model.Stats.LastPracticeDate.Format("2006-01-02")
	
	if lastPractice == "" {
		// First practice
		c.Model.Stats.CurrentStreak = 1
		c.Model.Stats.LongestStreak = 1
	} else if lastPractice == today {
		// Already practiced today, streak unchanged
	} else if lastPractice == time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
		// Practiced yesterday, streak continues
		c.Model.Stats.CurrentStreak++
		if c.Model.Stats.CurrentStreak > c.Model.Stats.LongestStreak {
			c.Model.Stats.LongestStreak = c.Model.Stats.CurrentStreak
		}
	} else {
		// Streak broken
		c.Model.Stats.CurrentStreak = 1
	}
	
	c.Model.Stats.LastPracticeDate = time.Now()
}

// checkAchievements checks for any newly unlocked achievements
func (c *Controller) checkAchievements() {
	// Check pattern mastery achievements
	for pattern, count := range c.Model.Stats.PatternCounts {
		if count >= 10 {
			achievementID := "pattern-master-" + pattern
			if achievement, exists := c.Model.Achievements[achievementID]; exists && !achievement.Earned {
				return func() tea.Msg {
					return model.AchievementUnlockedMsg{AchievementID: achievementID}
				}
			}
		}
	}
	
	// Check streak achievement
	if c.Model.Stats.CurrentStreak >= 30 {
		if achievement, exists := c.Model.Achievements["streak-virtuoso"]; exists && !achievement.Earned {
			return func() tea.Msg {
				return model.AchievementUnlockedMsg{AchievementID: "streak-virtuoso"}
			}
		}
	}
	
	// Check performance ace achievement (solve hard problem in under 15 min)
	if c.Model.Session.Problem.Difficulty == "hard" {
		elapsedTime := time.Since(c.Model.Session.StartTime)
		if elapsedTime < 15*time.Minute {
			if achievement, exists := c.Model.Achievements["performance-ace"]; exists && !achievement.Earned {
				return func() tea.Msg {
					return model.AchievementUnlockedMsg{AchievementID: "performance-ace"}
				}
			}
		}
	}
	
	return nil
}