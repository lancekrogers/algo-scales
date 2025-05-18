package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Update handles updates for the home screen
func (m Model) updateHome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.home.selectedOption > 0 {
				m.home.selectedOption--
			}
		case "down", "j":
			if m.home.selectedOption < len(m.home.options)-1 {
				m.home.selectedOption++
			}
		case "enter", "right", "l":
			return m.selectHomeOption()
		}
	}
	return m, nil
}

// View renders the home screen
func (m Model) viewHome() string {
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62")).
		MarginBottom(2)
	
	b.WriteString(titleStyle.Render("🎵 AlgoScales"))
	b.WriteString("\n\n")
	
	// Menu options
	for i, option := range m.home.options {
		cursor := "  "
		if i == m.home.selectedOption {
			cursor = "> "
			option = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")).
				Render(option)
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, option))
	}
	
	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(2)
	
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • q: Quit"))
	
	return b.String()
}

// selectHomeOption handles option selection on the home screen
func (m Model) selectHomeOption() (Model, tea.Cmd) {
	switch m.home.selectedOption {
	case 0: // Start Practice Session
		return m.navigate(StatePatternSelection), nil
	case 1: // Daily Scales
		return m.navigate(StateDaily), loadDailyScale()
	case 2: // View Statistics
		return m.navigate(StateStats), loadStats()
	case 3: // Settings
		return m.navigate(StateSettings), nil
	}
	return m, nil
}