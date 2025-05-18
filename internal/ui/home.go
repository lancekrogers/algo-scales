package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles updates for the home screen
func (m Model) updateHome(msg tea.Msg) (Model, tea.Cmd) {
	// Let the homeModel handle its own updates
	updatedHome, cmd := m.home.Update(msg)
	if h, ok := updatedHome.(homeModel); ok {
		m.home = h
	}
	
	// Handle selection change messages from homeModel
	if cmd != nil {
		if cmdMsg := cmd(); cmdMsg != nil {
			if selMsg, ok := cmdMsg.(SelectionChangedMsg); ok {
				m = m.navigate(selMsg.State)
				return m, nil
			}
		}
	}
	
	return m, cmd
}

// View renders the home screen
func (m Model) viewHome() string {
	return m.home.View()
}

