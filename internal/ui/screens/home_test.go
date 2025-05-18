package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/lancekrogers/algo-scales/internal/ui"
)

func TestHomeModel_Init(t *testing.T) {
	home := newHomeModel()
	home.width = 80
	home.height = 24
	
	cmd := home.Init()
	
	// Init should return nil for home screen
	assert.Nil(t, cmd)
}

func TestHomeModel_Update(t *testing.T) {
	home := newHomeModel()
	home.width = 80
	home.height = 24
	
	// Test arrow down
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, cmd := home.Update(msg)
	
	h, ok := updatedModel.(homeModel)
	require.True(t, ok)
	
	assert.Equal(t, 1, h.cursor)
	assert.Nil(t, cmd)
	
	// Test arrow up
	msg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, cmd = h.Update(msg)
	
	h, ok = updatedModel.(homeModel)
	require.True(t, ok)
	
	assert.Equal(t, 0, h.cursor)
	assert.Nil(t, cmd)
	
	// Test j key (vim down)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, cmd = h.Update(msg)
	
	h, ok = updatedModel.(homeModel)
	require.True(t, ok)
	
	assert.Equal(t, 1, h.cursor)
	assert.Nil(t, cmd)
	
	// Test k key (vim up)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedModel, cmd = h.Update(msg)
	
	h, ok = updatedModel.(homeModel)
	require.True(t, ok)
	
	assert.Equal(t, 0, h.cursor)
	assert.Nil(t, cmd)
	
	// Test enter key on "Start"
	h.cursor = 0
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd = h.Update(msg)
	
	// Should return a command to change state
	assert.NotNil(t, cmd)
	
	// Execute the command to see what message it produces
	resultMsg := cmd()
	selectionMsg, ok := resultMsg.(ui.SelectionChangedMsg)
	require.True(t, ok)
	assert.Equal(t, ui.StatePatternSelection, selectionMsg.GetState())
	
	// Test wrap around at bottom
	h.cursor = len(h.choices) - 1
	msg = tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = h.Update(msg)
	
	h, ok = updatedModel.(homeModel)
	require.True(t, ok)
	
	assert.Equal(t, 0, h.cursor) // Should wrap to top
	
	// Test wrap around at top
	h.cursor = 0
	msg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = h.Update(msg)
	
	h, ok = updatedModel.(homeModel)
	require.True(t, ok)
	
	assert.Equal(t, len(h.choices)-1, h.cursor) // Should wrap to bottom
}

func TestHomeModel_View(t *testing.T) {
	home := newHomeModel()
	home.cursor = 1
	home.width = 80
	home.height = 24
	
	view := home.View()
	
	// Check that all menu items are in the view
	assert.Contains(t, view, "Start")
	assert.Contains(t, view, "Daily Challenge")
	assert.Contains(t, view, "Stats")
	assert.Contains(t, view, "Settings")
	
	// Check for title
	assert.Contains(t, view, "AlgoScales")
	
	// Check for cursor indicator (assuming the selected item has a pointer)
	// The exact format depends on the implementation
	assert.Contains(t, view, ">") // Assuming ">" is used as cursor
}

func TestHomeModel_WindowSize(t *testing.T) {
	home := newHomeModel()
	home.width = 80
	home.height = 24
	
	// Test window size update
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, _ := home.Update(msg)
	
	h, ok := updatedModel.(homeModel)
	require.True(t, ok)
	
	assert.Equal(t, 120, h.width)
	assert.Equal(t, 40, h.height)
}