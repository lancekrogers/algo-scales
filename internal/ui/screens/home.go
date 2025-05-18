package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/ui"
)

// homeModel represents the home screen
type homeModel struct {
	choices  []string
	cursor   int
	width    int
	height   int
}

// newHomeModel creates a new home screen model
func newHomeModel() homeModel {
	return homeModel{
		choices: []string{"Start", "Daily Challenge", "Stats", "Settings"},
		cursor:  0,
	}
}

// Init initializes the home screen
func (m homeModel) Init() tea.Cmd {
	return nil
}

// Update handles updates to the home screen
func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
			
		case tea.KeyDown:
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
			
		case tea.KeyEnter:
			// Return selection change message based on cursor position
			switch m.cursor {
			case 0: // Start
				return m, func() tea.Msg {
					return ui.SelectionChangedMsg{State: ui.StatePatternSelection}
				}
			case 1: // Daily Challenge
				return m, func() tea.Msg {
					return ui.SelectionChangedMsg{State: ui.StateDaily}
				}
			case 2: // Stats
				return m, func() tea.Msg {
					return ui.SelectionChangedMsg{State: ui.StateStats}
				}
			case 3: // Settings
				return m, func() tea.Msg {
					return ui.SelectionChangedMsg{State: ui.StateSettings}
				}
			}
			
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "j": // vim down
				m.cursor++
				if m.cursor >= len(m.choices) {
					m.cursor = 0
				}
			case "k": // vim up
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.choices) - 1
				}
			}
		}
	}
	
	return m, nil
}

// View renders the home screen
func (m homeModel) View() string {
	title := "AlgoScales"
	subtitle := "Master algorithm patterns through musical scales"
	
	// Build menu items
	var menu strings.Builder
	for i, choice := range m.choices {
		cursor := "  " // two spaces
		if m.cursor == i {
			cursor = "> " // cursor indicator
		}
		
		menu.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
	}
	
	// Center the content
	lines := []string{
		title,
		subtitle,
		"",
		menu.String(),
		"",
		"Use arrow keys or j/k to navigate, Enter to select",
	}
	
	content := strings.Join(lines, "\n")
	
	// Simple centering - in production you'd use lipgloss.Place
	return content
}