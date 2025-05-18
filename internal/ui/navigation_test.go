package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNavigationFlow tests the complete navigation flow
func TestNavigationFlow(t *testing.T) {
	// Start with home screen
	model := New()
	model.ready = true
	model.width = 80
	model.height = 24
	
	assert.Equal(t, StateHome, model.state)
	
	// Navigate to pattern selection via selection message
	msg := SelectionChangedMsg{State: StatePatternSelection}
	updatedModel, _ := model.Update(msg)
	m, ok := updatedModel.(Model)
	require.True(t, ok)
	assert.Equal(t, StatePatternSelection, m.state)
	assert.Equal(t, StateHome, m.previousState)
	
	// Navigate back to home using esc key
	msg2 := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ = m.Update(msg2)
	m, ok = updatedModel.(Model)
	require.True(t, ok)
	assert.Equal(t, StateHome, m.state)
	
	// Navigate to stats
	msg3 := SelectionChangedMsg{State: StateStats}
	updatedModel, _ = m.Update(msg3)
	m, ok = updatedModel.(Model)
	require.True(t, ok)
	assert.Equal(t, StateStats, m.state)
	assert.Equal(t, StateHome, m.previousState)
}

// TestBackNavigation tests the back navigation
func TestBackNavigation(t *testing.T) {
	tests := []struct {
		name              string
		currentState      State
		previousState     State
		expectedBackState State
	}{
		{
			name:              "pattern selection to home (with previous)",
			currentState:      StatePatternSelection,
			previousState:     StateHome,
			expectedBackState: StateHome,
		},
		{
			name:              "pattern selection to home (default)",
			currentState:      StatePatternSelection,
			previousState:     StatePatternSelection, // same as current, triggers default
			expectedBackState: StateHome,
		},
		{
			name:              "problem list to pattern selection",  
			currentState:      StateProblemList,
			previousState:     StateProblemList, // triggers default
			expectedBackState: StatePatternSelection,
		},
		{
			name:              "problem detail to problem list",
			currentState:      StateProblemDetail,
			previousState:     StateProblemDetail, // triggers default
			expectedBackState: StateProblemList,
		},
		{
			name:              "session to problem detail",
			currentState:      StateSession,
			previousState:     StateSession, // triggers default
			expectedBackState: StateProblemDetail,
		},
		{
			name:              "stats to home",
			currentState:      StateStats,
			previousState:     StateStats, // triggers default
			expectedBackState: StateHome,
		},
		{
			name:              "daily to home",
			currentState:      StateDaily,
			previousState:     StateDaily, // triggers default
			expectedBackState: StateHome,
		},
		{
			name:              "settings to home",
			currentState:      StateSettings,
			previousState:     StateSettings, // triggers default
			expectedBackState: StateHome,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.state = tt.currentState
			model.previousState = tt.previousState
			
			newModel, _ := model.handleBack()
			
			assert.Equal(t, tt.expectedBackState, newModel.state)
		})
	}
}

// TestNavigationAnimation tests that navigation triggers animations
func TestNavigationAnimation(t *testing.T) {
	model := New()
	model.ready = true
	
	// Initial state should have no animation
	assert.Equal(t, AnimationNone, model.animation.Type)
	
	// Navigate forward should trigger slide right animation
	msg := SelectionChangedMsg{State: StatePatternSelection}
	updatedModel, cmd := model.Update(msg)
	m, ok := updatedModel.(Model)
	require.True(t, ok)
	
	assert.Equal(t, AnimationSlideRight, m.animation.Type)
	assert.NotNil(t, cmd) // Should return animation tick command
	
	// Navigate back should trigger slide left animation
	msg2 := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd = m.Update(msg2)
	m, ok = updatedModel.(Model)
	require.True(t, ok)
	
	assert.Equal(t, AnimationSlideLeft, m.animation.Type)
	assert.NotNil(t, cmd) // Should return animation tick command
}

// TestNavigateMethod tests the navigate method directly
func TestNavigateMethod(t *testing.T) {
	model := New()
	model.state = StateHome
	
	// Navigate to new state
	newModel := model.navigate(StatePatternSelection)
	
	assert.Equal(t, StatePatternSelection, newModel.state)
	assert.Equal(t, StateHome, newModel.previousState)
	assert.Equal(t, AnimationSlideRight, newModel.animation.Type)
	
	// Navigate back
	newModel = newModel.navigate(StateHome)
	
	assert.Equal(t, StateHome, newModel.state)
	assert.Equal(t, StatePatternSelection, newModel.previousState)
	assert.Equal(t, AnimationSlideLeft, newModel.animation.Type)
}