// Package controller handles user interactions and app logic
package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/common/highlight"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/utils"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/lancekrogers/algo-scales/internal/stats"
	"github.com/lancekrogers/algo-scales/internal/ui/model"
	"github.com/lancekrogers/algo-scales/internal/ui/view"
)

// Controller handles interactions between the model and view
type Controller struct {
	// The UI model
	Model *model.UIModel

	// Visualization components
	syntaxHighlighter *highlight.SyntaxHighlighter
	spinners          view.CustomSpinners
	patternViz        *view.PatternVisualization

	// Session management
	sessionManager interfaces.SessionManager
	activeSession  interfaces.Session
	
	// Stats service
	statsService   interfaces.StatsService
	
	// Problem repository
	problemRepo    interfaces.ProblemRepository
}

// NewController creates a new controller with the model and initializes components
func NewController(m *model.UIModel) *Controller {
	return &Controller{
		Model:             m,
		syntaxHighlighter: highlight.NewSyntaxHighlighter("monokai"),
		spinners:          view.NewCustomSpinners(),
		patternViz:        view.NewPatternVisualization(),
		sessionManager:    session.NewManager(),
		statsService:      stats.NewService(),
		problemRepo:       problem.NewRepository(),
	}
}

// Initialize loads initial data and sets up the application
func (c *Controller) Initialize() tea.Cmd {
	// Load initial problems
	return func() tea.Msg {
		// Load all problems using the problem repository
		problems, err := c.problemRepo.GetAll(context.Background())
		if err != nil {
			return model.ErrorMsg(fmt.Sprintf("Failed to load problems: %v", err))
		}

		// Load stats using StatsService
		summary, err := c.statsService.GetSummary(context.Background())
		if err != nil {
			log.Printf("Failed to load stats: %v", err)
			// Non-critical error, continue without stats
		} else {
			// Update model with stats  
			c.Model.Stats.ProblemsSolved = summary.TotalSolved
			c.Model.Stats.ProblemsAttempted = summary.TotalAttempted
			// SuccessRate can be calculated: summary.SuccessRate
		}
		
		// Load pattern stats
		patternStats, err := c.statsService.GetByPattern(context.Background())
		if err == nil {
			// Update model with pattern stats
			for pattern, stat := range patternStats {
				c.Model.Stats.PatternCounts[pattern] = stat.Solved
				if stat.Attempted > 0 {
					c.Model.Stats.PatternsProgress[pattern] = float64(stat.Solved) / float64(stat.Attempted)
				}
			}
		}

		// Convert interfaces.Problem to problem.Problem
		convertedProblems := make([]problem.Problem, len(problems))
		for i, p := range problems {
			convertedProblems[i] = c.convertInterfaceToLocalProblem(p)
		}

		return model.ProblemsLoadedMsg{Problems: convertedProblems}
	}
}

// Update handles messages and updates the model accordingly
func (c *Controller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case model.ProblemsLoadedMsg:
		// Update model with loaded problems
		c.Model.AvailableProblems = msg.Problems
		c.Model.Loading = false
		// No command needed

	case model.ErrorMsg:
		// Handle error messages
		c.Model.ErrorMessage = string(msg)
		c.Model.Loading = false
		// No command needed

	case model.TickMsg:
		// Handle timer ticks
		if c.activeSession != nil {
			// Update model with remaining time from session
			remaining := c.activeSession.GetTimeRemaining()
			c.Model.Session.TimeRemaining = remaining
			cmd = c.tickTimer()
		}

	case model.ProblemSelectedMsg:
		// Handle problem selection
		for _, p := range c.Model.AvailableProblems {
			if p.ID == msg.ProblemID {
				// Start a session with this problem
				cmd = c.startSession(p, msg.Mode)
				break
			}
		}

	case model.ShowHintsMsg:
		// Toggle hints visibility
		if c.activeSession != nil {
			c.activeSession.ShowHints(msg.Show)
			c.Model.Session.ShowHints = msg.Show
		}
		// No command needed

	case model.ShowSolutionMsg:
		// Toggle solution visibility
		if c.activeSession != nil {
			c.activeSession.ShowSolution(msg.Show)
			c.Model.Session.ShowSolution = msg.Show
		}
		// No command needed

	case model.CodeUpdatedMsg:
		// Code has been updated, refresh the view
		// No command needed

	case model.TestResultsMsg:
		// Handle test results
		c.Model.Session.TestResults = msg.Results
		if msg.AllPassed {
			// Problem solved!
			if c.activeSession.Finish(context.Background(), true) != nil {
				c.updateStatistics()
				cmd = c.checkAchievements()
			}
		}

	case model.SelectionMsg:
		// Handle selection changes based on app state
		cmd = c.handleSelection(msg.Index)

	case model.EditCodeMsg:
		// User wants to edit code in external editor
		cmd = func() tea.Msg {
			return c.openEditor()
		}

	case model.QuitMsg:
		// Quit the application
		return c.Model, tea.Quit
	}

	return c.Model, cmd
}

// handleSelection processes a selection based on current state
func (c *Controller) handleSelection(index int) tea.Cmd {
	return func() tea.Msg {
		c.Model.SelectedIndex = index

		switch c.Model.AppState {
		case model.StateInitial:
			// Initial menu selection
			switch index {
			case 0: // Start Practice
				c.Model.AppState = model.StateModeSelection
				c.Model.SelectedIndex = 0
				return nil

			case 1: // View Stats
				c.Model.AppState = model.StateStatistics
				return nil

			case 2: // Settings
				c.Model.AppState = model.StateSettings
				return nil

			case 3: // Quit
				return model.QuitMsg{}
			}
			// Default return for any other index
			return nil

		case model.StateModeSelection:
			// Practice mode selection
			mode := interfaces.PracticeMode // Default to practice mode

			switch index {
			case 0: // Learn mode
				mode = interfaces.LearnMode
			case 1: // Practice mode
				mode = interfaces.PracticeMode
			case 2: // Cram mode
				mode = interfaces.CramMode
			}

			c.Model.AppState = model.StatePatternSelection
			c.Model.SelectedIndex = 0
			c.Model.Session.Mode = string(mode)
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
				c.Model.AppState = model.StateProblemSelection
				c.Model.SelectedIndex = 0

				// Filter problems by pattern
				interfaceFiltered, err := c.problemRepo.GetByPattern(context.Background(), pattern)
				if err != nil {
					return model.ErrorMsg(fmt.Sprintf("Failed to filter problems: %v", err))
				}
				
				// Convert to local problem type
				filtered := make([]problem.Problem, len(interfaceFiltered))
				for i, p := range interfaceFiltered {
					filtered[i] = c.convertInterfaceToLocalProblem(p)
				}

				// Update available problems list
				c.Model.AvailableProblems = filtered
				return nil
			}
			// Default return if index is out of range
			return nil

		case model.StateProblemSelection:
			// Problem selection
			if index >= 0 && index < len(c.Model.AvailableProblems) {
				problem := c.Model.AvailableProblems[index]
				return model.ProblemSelectedMsg{
					ProblemID: problem.ID,
					Mode:      c.Model.Session.Mode,
				}
			}
			// Default return if no valid problem selected
			return nil

		case model.StateSession:
			// Session controls
			// Default no-op for session state
			return nil

		case model.StateStatistics:
			// Statistics view controls
			// Return to main menu on any selection
			c.Model.AppState = model.StateInitial
			c.Model.SelectedIndex = 0
			return nil

		case model.StateSettings:
			// Settings controls
			// Return to main menu on any selection
			c.Model.AppState = model.StateInitial
			c.Model.SelectedIndex = 0
			return nil
		}

		// Default return for any state not explicitly handled
		return nil
	}
}

// startSession starts a new practice session with the given problem
func (c *Controller) startSession(p problem.Problem, mode string) tea.Cmd {
	// Convert mode string to enum
	sessionMode := interfaces.PracticeMode
	switch mode {
	case "learn":
		sessionMode = interfaces.LearnMode
	case "practice":
		sessionMode = interfaces.PracticeMode
	case "cram":
		sessionMode = interfaces.CramMode
	}

	// Prepare session options
	options := interfaces.SessionOptions{
		Mode:      sessionMode,
		ProblemID: p.ID,
		Language:  c.Model.Session.Language,
		Pattern:   c.Model.Session.CurrentPattern,
	}

	// Create session
	session, err := c.sessionManager.StartSession(context.Background(), options)
	if err != nil {
		return func() tea.Msg {
			return model.ErrorMsg(fmt.Sprintf("Failed to start session: %v", err))
		}
	}

	// Store the active session
	c.activeSession = session

	// Update the model with session info
	interfaceProblem := session.GetProblem()
	
	// Convert interface problem to local problem type
	localProblem := c.convertInterfaceToLocalProblem(*interfaceProblem)

	c.Model.Session = model.Session{
		Active:         true,
		Mode:           string(options.Mode),
		Problem:        &localProblem,
		StartTime:      session.GetStartTime(),
		TimeRemaining:  session.GetTimeRemaining(),
		ShowHints:      session.AreHintsShown(),
		ShowSolution:   session.IsSolutionShown(),
		Language:       session.GetLanguage(),
		Code:           session.GetCode(),
		CurrentPattern: options.Pattern,
	}

	// Update app state
	c.Model.AppState = model.StateSession

	// Start the session timer
	return c.tickTimer()
}

// tickTimer handles the session timer
func (c *Controller) tickTimer() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return model.TickMsg{}
	})
}

// openEditor opens the code in an external editor
func (c *Controller) openEditor() tea.Msg {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "algo-scales-*."+view.GetLanguageExtension(c.Model.Session.Language))
	if err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to create temp file: %v", err))
	}
	defer func() {
		if rmErr := os.Remove(tmpfile.Name()); rmErr != nil {
			log.Printf("tmpfile cleanup failed: %v", rmErr)
		}
	}()

	// Write the current code to the file
	if _, err := tmpfile.Write([]byte(c.Model.Session.Code)); err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to write to temp file: %v", err))
	}
	if err := tmpfile.Close(); err != nil {
		return model.ErrorMsg(fmt.Sprintf("Failed to close temp file: %v", err))
	}

	// Use our utility function to open the editor with exec.Command
	editor := utils.OpenEditor(tmpfile.Name())
	cmd := tea.ExecProcess(editor, func(err error) tea.Msg {
		if err != nil {
			return model.ErrorMsg(fmt.Sprintf("Editor exited with error: %v", err))
		}

		// Read the updated file
		content, err := utils.ReadFile(tmpfile.Name())
		if err != nil {
			return model.ErrorMsg(fmt.Sprintf("Failed to read temp file: %v", err))
		}

		// Update the code in the session and model
		if err := c.activeSession.SetCode(string(content)); err != nil {
			return model.ErrorMsg(fmt.Sprintf("failed updating session: %v", err))
		}
		c.Model.Session.Code = string(content)
		return model.CodeUpdatedMsg{}
	})

	return cmd
}

// submitSolution submits the current solution for testing
func (c *Controller) submitSolution() tea.Cmd {
	return func() tea.Msg {
		// Validate that we're in an active session
		if c.activeSession == nil {
			return model.ErrorMsg("No active session")
		}

		// Create a loading message
		c.Model.Loading = true

		// Execute the tests using our session
		results, allPassed, err := c.activeSession.RunTests(context.Background())
		if err != nil {
			return model.ErrorMsg(fmt.Sprintf("Failed to run tests: %v", err))
		}

		// Convert result types
		modelResults := make([]model.TestResult, len(results))
		for i, result := range results {
			modelResults[i] = model.TestResult{
				Input:    result.Input,
				Expected: result.Expected,
				Actual:   result.Actual,
				Passed:   result.Passed,
			}
		}

		// Update session with test results
		c.Model.Session.TestResults = modelResults
		return model.TestResultsMsg{Results: modelResults, AllPassed: allPassed}
	}
}

// updateStatistics updates user statistics after solving a problem
func (c *Controller) updateStatistics() {
	// Increment problems solved count
	c.Model.Stats.ProblemsSolved++

	// Update time spent
	c.Model.Stats.TotalTime += time.Since(c.activeSession.GetStartTime())

	// Update pattern stats if defined
	problem := c.activeSession.GetProblem()
	for _, pattern := range problem.Tags {
		c.Model.Stats.PatternCounts[pattern]++
		// Skip updating progress during stats update to avoid infinite recursion
		// Progress will be calculated during initial load
	}

	// Update difficulty stats
	c.Model.Stats.DifficultyCounts[problem.Difficulty]++

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
	
	// Store session stats
	sessionStats := interfaces.SessionStats{
		ProblemID:    problem.ID,
		StartTime:    c.activeSession.GetStartTime(),
		EndTime:      time.Now(),
		Duration:     time.Since(c.activeSession.GetStartTime()),
		Solved:       true,
		Mode:         c.Model.Session.Mode,
		HintsUsed:    c.activeSession.AreHintsShown(),
		SolutionUsed: c.activeSession.IsSolutionShown(),
		Patterns:     problem.Tags,
		Difficulty:   problem.Difficulty,
	}
	
	// Record stats using the stats service
	if err := c.statsService.RecordSession(context.Background(), sessionStats); err != nil {
		log.Printf("Failed to record session stats: %v", err)
	}
}

// checkAchievements checks for any newly unlocked achievements
func (c *Controller) checkAchievements() tea.Cmd {
	return func() tea.Msg {
		// Check pattern mastery achievements
		for pattern, count := range c.Model.Stats.PatternCounts {
			if count >= 10 {
				achievementID := "pattern-master-" + pattern
				if achievement, exists := c.Model.Achievements[achievementID]; exists && !achievement.Earned {
					return model.AchievementUnlockedMsg{AchievementID: achievementID}
				}
			}
		}

		// Check streak achievement
		if c.Model.Stats.CurrentStreak >= 30 {
			if achievement, exists := c.Model.Achievements["streak-virtuoso"]; exists && !achievement.Earned {
				return model.AchievementUnlockedMsg{AchievementID: "streak-virtuoso"}
			}
		}

		// Check performance ace achievement (solve hard problem in under 15 min)
		problem := c.activeSession.GetProblem()
		if problem.Difficulty == "hard" {
			elapsedTime := time.Since(c.activeSession.GetStartTime())
			if elapsedTime < 15*time.Minute {
				if achievement, exists := c.Model.Achievements["performance-ace"]; exists && !achievement.Earned {
					return model.AchievementUnlockedMsg{AchievementID: "performance-ace"}
				}
			}
		}

		return nil
	}
}
// convertInterfaceToLocalProblem converts an interfaces.Problem to a problem.Problem
func (c *Controller) convertInterfaceToLocalProblem(p interfaces.Problem) problem.Problem {
	// Convert test cases
	testCases := make([]problem.TestCase, len(p.TestCases))
	for i, tc := range p.TestCases {
		testCases[i] = problem.TestCase{
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}
	
	// Create starter code map
	starterCode := make(map[string]string)
	if p.StarterCode != nil {
		starterCode = p.StarterCode
	} else {
		for _, lang := range p.Languages {
			starterCode[lang] = ""
		}
	}
	
	return problem.Problem{
		ID:                  p.ID,
		Title:               p.Title,
		Description:         p.Description,
		Difficulty:          p.Difficulty,
		Patterns:            p.Tags, // Map Tags to Patterns
		Companies:           p.Companies,
		TestCases:           testCases,
		StarterCode:         starterCode,
		Solutions:           make(map[string]string),
		EstimatedTime:       30, // Default value
		Examples:            []problem.Example{}, // Empty for now
		Constraints:         []string{}, // Empty for now
		PatternExplanation:  "", // Empty for now
		SolutionWalkthrough: []string{}, // Empty for now
	}
}
