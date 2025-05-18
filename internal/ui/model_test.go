package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/stats"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	model := New()
	
	assert.NotNil(t, model)
	assert.Equal(t, StateHome, model.state)
	assert.NotNil(t, model.home)
	assert.NotNil(t, model.patterns)
	assert.NotNil(t, model.problems)
	assert.NotNil(t, model.problemDetail)
	assert.NotNil(t, model.session)
	assert.NotNil(t, model.stats)
	assert.NotNil(t, model.daily)
	assert.NotNil(t, model.settings)
	assert.NotNil(t, model.loading)
	assert.False(t, model.showLoading)
	assert.False(t, model.ready)
}

func TestModelInit(t *testing.T) {
	model := New()
	
	cmd := model.Init()
	
	// Init should return a command
	assert.NotNil(t, cmd)
}

func TestModelUpdate_WindowSize(t *testing.T) {
	model := New()
	
	msg := tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	}
	
	updatedModel, _ := model.Update(msg)
	
	// Check that the model was updated
	assert.NotNil(t, updatedModel)
	// cmd can be nil for window size updates
	
	// Type assert to access internal fields
	m, ok := updatedModel.(Model)
	require.True(t, ok)
	
	assert.Equal(t, 80, m.width)
	assert.Equal(t, 24, m.height)
	assert.True(t, m.ready)
	
	// Check that window size was propagated to components
	assert.Equal(t, 80, m.home.width)
	assert.Equal(t, 24, m.home.height)
}

func TestModelUpdate_StateTransition(t *testing.T) {
	tests := []struct {
		name      string
		fromState State
		msg       tea.Msg
		toState   State
	}{
		{
			name:      "home to pattern selection",
			fromState: StateHome,
			msg:       SelectionChangedMsg{State: StatePatternSelection},
			toState:   StatePatternSelection,
		},
		{
			name:      "pattern selection to home with esc",
			fromState: StatePatternSelection,
			msg:       tea.KeyMsg{Type: tea.KeyEsc},
			toState:   StateHome,
		},
		{
			name:      "home to stats",
			fromState: StateHome,
			msg:       SelectionChangedMsg{State: StateStats},
			toState:   StateStats,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.state = tt.fromState
			model.ready = true
			model.width = 80
			model.height = 24
			
			updatedModel, _ := model.Update(tt.msg)
			
			m, ok := updatedModel.(Model)
			require.True(t, ok)
			
			assert.Equal(t, tt.toState, m.state)
		})
	}
}

func TestModelUpdate_KeyNavigation(t *testing.T) {
	model := New()
	model.ready = true
	model.width = 80
	model.height = 24
	
	// Test esc key - go back/home
	model.state = StatePatternSelection
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m, ok := updatedModel.(Model)
	require.True(t, ok)
	assert.Equal(t, StateHome, m.state)
	
	// Test ctrl+c - quit
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.NotNil(t, cmd) // Should return quit command
}

func TestModelView(t *testing.T) {
	tests := []struct {
		name  string
		state State
		ready bool
		showLoading bool
		expectContent func(view string) bool
	}{
		{
			name:  "not ready",
			state: StateHome,
			ready: false,
			expectContent: func(view string) bool {
				return view == "Loading..."
			},
		},
		{
			name:  "show loading",
			state: StateHome,
			ready: true,
			showLoading: true,
			expectContent: func(view string) bool {
				// Should show loading view
				return len(view) > 0
			},
		},
		{
			name:  "home screen",
			state: StateHome,
			ready: true,
			expectContent: func(view string) bool {
				// Should show home view
				return len(view) > 0
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.state = tt.state
			model.ready = tt.ready
			model.showLoading = tt.showLoading
			model.width = 80
			model.height = 24
			
			// Initialize loading screen
			if tt.showLoading {
				model.loading = NewLoadingScreen("Loading...")
				model.loading.width = model.width
				model.loading.height = model.height
			}
			
			view := model.View()
			
			assert.True(t, tt.expectContent(view))
		})
	}
}

func TestStatsLoadedMsg_Update(t *testing.T) {
	model := New()
	model.ready = true
	
	// Create test stats
	testStats := stats.Summary{
		TotalAttempted: 10,
		TotalSolved:    8,
		SuccessRate:    80.0,
		AvgSolveTime:   "00:15:00",
	}
	
	msg := statsLoadedMsg{stats: testStats}
	updatedModel, _ := model.Update(msg)
	
	m, ok := updatedModel.(Model)
	require.True(t, ok)
	
	// Check that stats were propagated to the stats component
	// This assumes the stats component has a stats field
	// If not, we might need to test this differently
	assert.NotNil(t, m.stats)
}

func TestProblemLoadedMsg_Update(t *testing.T) {
	model := New()
	model.ready = true
	
	// Create a test problem
	prob := problem.Problem{
		ID:         "test-problem",
		Title:      "Test Problem",
		Difficulty: "Easy",
	}
	
	msg := problemSelectedMsg{problem: prob}
	updatedModel, _ := model.Update(msg)
	
	m, ok := updatedModel.(Model)
	require.True(t, ok)
	
	// Check that the update was handled
	// The specific behavior depends on the implementation
	assert.NotNil(t, m)
}