package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Update handles updates for the problem list screen
func (m Model) updateProblemList(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case problemsLoadedMsg:
		m.problems.problems = msg.problems
		m.problems.loading = false
		return m, nil
		
	case problemsErrorMsg:
		m.problems.loading = false
		// Handle error (could set an error message)
		return m, nil
		
	case tea.KeyMsg:
		if m.problems.loading {
			return m, nil
		}
		
		switch msg.String() {
		case "up", "k":
			if m.problems.selectedIndex > 0 {
				m.problems.selectedIndex--
			}
		case "down", "j":
			if m.problems.selectedIndex < len(m.problems.problems)-1 {
				m.problems.selectedIndex++
			}
		case "enter", "right", "l":
			if len(m.problems.problems) > 0 {
				selectedProblem := m.problems.problems[m.problems.selectedIndex]
				m.problemDetail.problem = selectedProblem
				return m.navigate(StateProblemDetail), nil
			}
		}
	}
	return m, nil
}

// View renders the problem list screen
func (m Model) viewProblemList() string {
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62")).
		MarginBottom(2)
	
	b.WriteString(titleStyle.Render(fmt.Sprintf("%s Problems", m.problems.pattern)))
	b.WriteString("\n\n")
	
	if m.problems.loading {
		b.WriteString("Loading problems...")
		return b.String()
	}
	
	if len(m.problems.problems) == 0 {
		b.WriteString("No problems found for this pattern.")
		return b.String()
	}
	
	// Problem list
	for i, problem := range m.problems.problems {
		cursor := "  "
		if i == m.problems.selectedIndex {
			cursor = "> "
		}
		
		// Difficulty color
		diffColor := "243"
		switch problem.Difficulty {
		case "easy":
			diffColor = "46"  // Green
		case "medium":
			diffColor = "214" // Orange
		case "hard":
			diffColor = "196" // Red
		}
		
		diffStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(diffColor))
		
		line := fmt.Sprintf("%s%-30s %s", cursor, problem.Title, diffStyle.Render(problem.Difficulty))
		
		if i == m.problems.selectedIndex {
			line = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")).
				Render(problem.Title) + " " + diffStyle.Render(problem.Difficulty)
			line = cursor + line
		}
		
		b.WriteString(line + "\n")
	}
	
	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(2)
	
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Back"))
	
	return b.String()
}