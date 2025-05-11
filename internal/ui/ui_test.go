// Tests for UI components
package ui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test session
func createTestSession() *session.Session {
	return &session.Session{
		Options: session.Options{
			Mode:       session.LearnMode,
			Language:   "go",
			Timer:      30,
			Pattern:    "",
			Difficulty: "",
			ProblemID:  "test-problem",
		},
		Problem: &problem.Problem{
			ID:            "test-problem",
			Title:         "Test Problem",
			Difficulty:    "Easy",
			Patterns:      []string{"hash-map"},
			EstimatedTime: 15,
			Companies:     []string{"Test Company"},
			Description:   "This is a test problem description.",
			Examples: []problem.Example{
				{
					Input:       "nums = [1,2,3], target = 5",
					Output:      "[1,2]",
					Explanation: "Because nums[1] + nums[2] == 5",
				},
			},
			Constraints:         []string{"1 <= nums.length <= 100"},
			PatternExplanation:  "This is a test pattern explanation.",
			SolutionWalkthrough: []string{"Step 1", "Step 2"},
			StarterCode: map[string]string{
				"go":         "func solution() {}\n",
				"python":     "def solution():\n    pass\n",
				"javascript": "function solution() {}\n",
			},
			Solutions: map[string]string{
				"go":         "func solution() { return [1, 2] }\n",
				"python":     "def solution():\n    return [1, 2]\n",
				"javascript": "function solution() { return [1, 2]; }\n",
			},
			TestCases: []problem.TestCase{
				{
					Input:    "[1,2,3], 5",
					Expected: "[1,2]",
				},
			},
		},
		ShowHints:    true,
		ShowPattern:  true,
		ShowSolution: false,
		CodeFile:     "/tmp/test-solution.go",
	}
}

// Test the initialization of the UI model
func TestModelInit(t *testing.T) {
	// Create a test session
	s := createTestSession()

	// Create a model
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

	// Initialize the model
	cmd := m.Init()

	// Check that the command is not nil
	assert.NotNil(t, cmd)
}

// Test the update function for window sizing
func TestUpdateWindowSize(t *testing.T) {
	// Create a test session
	s := createTestSession()

	// Create a model
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

	// Create a window size message
	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	// Update the model
	updatedModel, _ := m.Update(msg)
	typedModel, ok := updatedModel.(Model)

	// Check that the update worked
	assert.True(t, ok)
	assert.Equal(t, 100, typedModel.width)
	assert.Equal(t, 50, typedModel.height)
	assert.True(t, typedModel.ready)
	assert.NotNil(t, typedModel.viewport)
}

// Test the key handling
func TestUpdateKeyPress(t *testing.T) {
	// Create a test session
	s := createTestSession()

	// Create a model
	m := Model{
		session:      s,
		keyMap:       DefaultKeyMap(),
		help:         help.New(),
		messageStyle: infoStyle,
		message:      "Press '?' for help, 'e' to open editor",
		ready:        true, // Pretend we've received window size
	}

	// Set up timer
	timerDuration := time.Duration(s.Options.Timer) * time.Minute
	m.timer = timer.NewWithInterval(timerDuration, time.Second)
	m.timeRemaining = timerDuration

	// Create a viewport for testing
	m.viewport = viewport.New(100, 40)

	// Test case for showing help
	t.Run("ToggleHelp", func(t *testing.T) {
		// Create a key message for '?'
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}

		// Update the model
		updatedModel, _ := m.Update(msg)
		typedModel, ok := updatedModel.(Model)

		// Check that help was toggled
		assert.True(t, ok)
		assert.True(t, typedModel.showHelp)

		// Toggle back
		updatedModel, _ = typedModel.Update(msg)
		typedModel, ok = updatedModel.(Model)

		// Check that help was toggled back
		assert.True(t, ok)
		assert.False(t, typedModel.showHelp)
	})

	// Test case for showing solution
	t.Run("ShowSolution", func(t *testing.T) {
		// Create a key message for 's'
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}

		// Update the model
		updatedModel, _ := m.Update(msg)
		typedModel, ok := updatedModel.(Model)

		// Check that solution was shown
		assert.True(t, ok)
		assert.True(t, typedModel.session.ShowSolution)
		assert.Contains(t, typedModel.message, "Solution shown")
	})

	// Test case for submitting solution
	t.Run("SubmitSolution", func(t *testing.T) {
		// Create a key message for 'enter'
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		// Update the model
		updatedModel, _ := m.Update(msg)
		typedModel, ok := updatedModel.(Model)

		// Check that solution was submitted
		assert.True(t, ok)
		assert.True(t, typedModel.problemCompleted)
		assert.Contains(t, typedModel.message, "submitted successfully")
	})

	// Test case for quit confirmation
	t.Run("QuitConfirmation", func(t *testing.T) {
		// Set up model for testing
		testModel := m
		testModel.session.Options.Mode = session.CramMode

		// Create a key message for 'q'
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

		// Update the model
		updatedModel, _ := testModel.Update(msg)
		typedModel, ok := updatedModel.(Model)

		// Check that confirmation was requested
		assert.True(t, ok)
		assert.True(t, typedModel.confirmQuit)
		assert.Contains(t, typedModel.message, "Quit without completing")

		// Confirm quit
		yesMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}

		// Update the model
		updatedModel, cmd := typedModel.Update(yesMsg)
		typedModel, ok = updatedModel.(Model)

		// Check that quit was confirmed
		assert.True(t, ok)
		assert.True(t, typedModel.shouldQuit)
		assert.NotNil(t, cmd)
	})
}

// Test timer functionality
func TestTimer(t *testing.T) {
	// Create a test session
	s := createTestSession()

	// Create a model
	m := Model{
		session:      s,
		keyMap:       DefaultKeyMap(),
		help:         help.New(),
		messageStyle: infoStyle,
		message:      "Press '?' for help, 'e' to open editor",
		ready:        true, // Pretend we've received window size
	}

	// Set up timer with a very short duration for testing
	timerDuration := 100 * time.Millisecond
	m.timer = timer.NewWithInterval(timerDuration, 50*time.Millisecond)
	m.timeRemaining = timerDuration

	// Initialize the timer
	m.timer.Init()

	// Create a tick message
	tickMsg := timer.TickMsg{}

	// Update the model
	updatedModel, _ := m.Update(tickMsg)
	typedModel, ok := updatedModel.(Model)

	// Check that the timer was updated
	assert.True(t, ok)
	assert.True(t, typedModel.timeRemaining < timerDuration)

	// Create a timeout message
	timeoutMsg := timer.TimeoutMsg{}

	// Set to cram mode for testing
	typedModel.session.Options.Mode = session.CramMode

	// Update the model
	updatedModel, cmd := typedModel.Update(timeoutMsg)
	typedModel, ok = updatedModel.(Model)

	// Check that timeout was handled
	assert.True(t, ok)
	assert.Contains(t, typedModel.message, "Time's up")
	assert.NotNil(t, cmd)
}

// Test format duration
func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "00:30"},
		{2 * time.Minute, "02:00"},
		{2*time.Minute + 30*time.Second, "02:30"},
		{1*time.Hour + 15*time.Minute + 45*time.Second, "75:45"},
	}

	for _, tc := range testCases {
		t.Run("FormatDuration", func(t *testing.T) {
			result := formatDuration(tc.duration)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Test the view rendering
func TestView(t *testing.T) {
	// Create a test session
	s := createTestSession()

	// Create a model
	m := Model{
		session:      s,
		keyMap:       DefaultKeyMap(),
		help:         help.New(),
		messageStyle: infoStyle,
		message:      "Test message",
		ready:        true, // Pretend we've received window size
		width:        100,
		height:       50,
	}

	// Set up timer
	timerDuration := time.Duration(s.Options.Timer) * time.Minute
	m.timer = timer.NewWithInterval(timerDuration, time.Second)
	m.timeRemaining = timerDuration

	// Create a viewport for testing
	m.viewport = viewport.New(100, 40)
	m.viewport.SetContent("Test content")

	// Render the view
	view := m.View()

	// Check that the view contains the expected elements
	assert.Contains(t, view, s.Problem.Title)
	assert.Contains(t, view, s.Problem.Difficulty)
	assert.Contains(t, view, "Test message")
	assert.Contains(t, view, formatDuration(timerDuration))
	assert.Contains(t, view, "Test content")
}
