package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Update handles updates for the stats screen
func (m Model) updateStats(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Initialize or update viewport
		if m.stats.viewport.Width == 0 {
			m.stats.viewport = viewport.New(msg.Width-4, msg.Height-8)
			m.stats.viewport.SetContent(m.statsContent())
		} else {
			m.stats.viewport.Width = msg.Width - 4
			m.stats.viewport.Height = msg.Height - 8
		}
		
	case statsLoadedMsg:
		m.stats.summary = msg.stats
		m.stats.loading = false
		m.stats.viewport.SetContent(m.statsContent())
		
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// Refresh stats
			m.stats.loading = true
			return m, loadStats()
		default:
			// Pass through to viewport
			m.stats.viewport, cmd = m.stats.viewport.Update(msg)
		}
	}
	
	return m, cmd
}

// View renders the stats screen
func (m Model) viewStats() string {
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62")).
		MarginBottom(2)
	
	b.WriteString(titleStyle.Render("📊 Statistics"))
	b.WriteString("\n\n")
	
	if m.stats.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
		b.WriteString(loadingStyle.Render("Loading statistics..."))
		return b.String()
	}
	
	// Viewport with stats content
	b.WriteString(m.stats.viewport.View())
	b.WriteString("\n\n")
	
	// Action bar
	actionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	
	actions := []string{
		"r: Refresh",
		"Esc: Back",
	}
	
	b.WriteString(actionStyle.Render(strings.Join(actions, " • ")))
	
	return b.String()
}

// statsContent generates the content for the stats viewport
func (m Model) statsContent() string {
	var content strings.Builder
	s := m.stats.summary
	
	// Overview
	overviewStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		MarginBottom(1)
	
	content.WriteString(overviewStyle.Render("Overview"))
	content.WriteString("\n\n")
	
	statsBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(40)
	
	overviewContent := fmt.Sprintf(
		"Total Problems Attempted: %d\n"+
		"Total Problems Solved: %d\n"+
		"Success Rate: %.1f%%\n"+
		"Average Solve Time: %s",
		s.TotalAttempted,
		s.TotalSolved,
		s.SuccessRate*100,
		s.AvgSolveTime,
	)
	
	content.WriteString(statsBoxStyle.Render(overviewContent))
	content.WriteString("\n\n")
	
	// Best Performance
	if s.FastestSolve.ProblemID != "" {
		content.WriteString(overviewStyle.Render("🏆 Fastest Solve"))
		content.WriteString("\n\n")
		
		fastestContent := fmt.Sprintf(
			"Problem: %s\n"+
			"Time: %s",
			s.FastestSolve.ProblemID,
			s.FastestSolve.Time,
		)
		
		successStyle := statsBoxStyle.Copy().
			BorderForeground(lipgloss.Color("46"))
		
		content.WriteString(successStyle.Render(fastestContent))
		content.WriteString("\n\n")
	}
	
	// Most Challenging
	if s.MostChallenging.ProblemID != "" {
		content.WriteString(overviewStyle.Render("🔥 Most Challenging"))
		content.WriteString("\n\n")
		
		challengingContent := fmt.Sprintf(
			"Problem: %s\n"+
			"Attempts: %d",
			s.MostChallenging.ProblemID,
			s.MostChallenging.Attempts,
		)
		
		challengeStyle := statsBoxStyle.Copy().
			BorderForeground(lipgloss.Color("214"))
		
		content.WriteString(challengeStyle.Render(challengingContent))
		content.WriteString("\n\n")
	}
	
	// Pattern Breakdown (if we add PatternStats in the future)
	
	return content.String()
}

