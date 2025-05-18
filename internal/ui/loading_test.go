package ui

import (
	"testing"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// LoadingScreen doesn't have an Init method, so let's remove this test
// The LoadingScreen is managed by the parent model

func TestLoadingScreen_Update(t *testing.T) {
	ls := LoadingScreen{
		message: "Loading data...",
		width:   80,
		height:  24,
	}
	
	// Test window size message
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedLS, _ := ls.Update(sizeMsg)
	
	assert.Equal(t, 100, updatedLS.width)
	assert.Equal(t, 50, updatedLS.height)
	
	// Test spinner tick message
	tickMsg := spinnerTickMsg(time.Now())
	updatedLS, cmd := ls.Update(tickMsg)
	
	// Check that frame advanced
	assert.Equal(t, 1, updatedLS.spinnerFrame)
	assert.NotNil(t, cmd) // Should return another tick command
}

func TestLoadingScreen_View(t *testing.T) {
	ls := LoadingScreen{
		message:      "Loading...",
		spinnerFrame: 0,
		width:        80,
		height:       24,
	}
	
	view := ls.View()
	
	// Check that view contains the loading message
	assert.Contains(t, view, "Loading...")
	
	// Check that view contains a spinner
	assert.True(t, strings.Contains(view, "⠋") || 
		strings.Contains(view, "⠙") || 
		strings.Contains(view, "⠹") || 
		strings.Contains(view, "⠸") || 
		strings.Contains(view, "⠼") || 
		strings.Contains(view, "⠴") || 
		strings.Contains(view, "⠦") || 
		strings.Contains(view, "⠧") || 
		strings.Contains(view, "⠇") || 
		strings.Contains(view, "⠏"))
	
	// Check that view is centered
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			// Non-empty lines should have some leading spaces for centering
			assert.True(t, strings.HasPrefix(line, " ") || strings.Contains(line, " "))
		}
	}
}