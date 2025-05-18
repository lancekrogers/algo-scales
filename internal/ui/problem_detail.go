package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Update handles updates for the problem detail screen
func (m Model) updateProblemDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Initialize or update viewport
		if m.problemDetail.viewport.Width == 0 {
			m.problemDetail.viewport = viewport.New(msg.Width-4, msg.Height-10)
			m.problemDetail.viewport.SetContent(m.problemDetailContent())
		} else {
			m.problemDetail.viewport.Width = msg.Width - 4
			m.problemDetail.viewport.Height = msg.Height - 10
		}
		
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Start session with selected problem
			m.session.problem = m.problemDetail.problem
			m.session.startTime = time.Now()
			return m.navigate(StateSession), startSession(m.problemDetail.problem)
		case "h":
			// Toggle hint
			m.problemDetail.showHint = !m.problemDetail.showHint
			m.problemDetail.viewport.SetContent(m.problemDetailContent())
		case "i":
			// Show additional info
			m.problemDetail.showInfo = !m.problemDetail.showInfo
			m.problemDetail.viewport.SetContent(m.problemDetailContent())
		default:
			// Pass through to viewport
			m.problemDetail.viewport, cmd = m.problemDetail.viewport.Update(msg)
		}
	}
	
	return m, cmd
}

// View renders the problem detail screen
func (m Model) viewProblemDetail() string {
	var b strings.Builder
	
	// Title bar
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62")).
		MarginBottom(1)
	
	difficultyColor := "243"
	switch m.problemDetail.problem.Difficulty {
	case "easy":
		difficultyColor = "46"
	case "medium":
		difficultyColor = "214"
	case "hard":
		difficultyColor = "196"
	}
	
	diffStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(difficultyColor)).
		Bold(true)
	
	title := fmt.Sprintf("%s %s", 
		m.problemDetail.problem.Title,
		diffStyle.Render(fmt.Sprintf("(%s)", m.problemDetail.problem.Difficulty)))
	
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")
	
	// Viewport with problem content
	b.WriteString(m.problemDetail.viewport.View())
	b.WriteString("\n\n")
	
	// Action bar
	actionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	
	actions := []string{
		"Enter: Start Session",
		"h: Toggle Hint",
		"i: More Info",
		"Esc: Back",
	}
	
	b.WriteString(actionStyle.Render(strings.Join(actions, " • ")))
	
	// Progress indicator
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Align(lipgloss.Right)
	
	progress := fmt.Sprintf("%d%% ", m.problemDetail.viewport.ScrollPercent())
	b.WriteString(progressStyle.Render(progress))
	
	return b.String()
}

// problemDetailContent generates the content for the problem detail viewport
func (m Model) problemDetailContent() string {
	var content strings.Builder
	p := m.problemDetail.problem
	
	// Description
	content.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		Render("Description"))
	content.WriteString("\n\n")
	content.WriteString(p.Description)
	content.WriteString("\n\n")
	
	// Examples
	if len(p.Examples) > 0 {
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Render("Examples"))
		content.WriteString("\n\n")
		
		for i, example := range p.Examples {
			content.WriteString(fmt.Sprintf("Example %d:\n", i+1))
			content.WriteString(fmt.Sprintf("Input: %s\n", example.Input))
			content.WriteString(fmt.Sprintf("Output: %s\n", example.Output))
			if example.Explanation != "" {
				content.WriteString(fmt.Sprintf("Explanation: %s\n", example.Explanation))
			}
			content.WriteString("\n")
		}
	}
	
	// Hint (if toggled)
	if m.problemDetail.showHint && p.PatternExplanation != "" {
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214")).
			Render("💡 Pattern Explanation"))
		content.WriteString("\n\n")
		content.WriteString(p.PatternExplanation)
		content.WriteString("\n\n")
	}
	
	// Additional Info (if toggled)
	if m.problemDetail.showInfo {
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			Render("ℹ️  Additional Information"))
		content.WriteString("\n\n")
		
		// Patterns
		if len(p.Patterns) > 0 {
			content.WriteString("**Patterns:** " + strings.Join(p.Patterns, ", "))
			content.WriteString("\n")
		}
		
		// Estimated Time
		if p.EstimatedTime > 0 {
			content.WriteString(fmt.Sprintf("**Estimated Time:** %d minutes", p.EstimatedTime))
			content.WriteString("\n")
		}
		
		// Companies
		if len(p.Companies) > 0 {
			content.WriteString("**Asked by:** " + strings.Join(p.Companies, ", "))
			content.WriteString("\n")
		}
		
		content.WriteString("\n")
	}
	
	return content.String()
}