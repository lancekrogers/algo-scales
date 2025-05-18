package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Pattern categories - will be loaded dynamically from problems
var defaultPatterns = []string{
	"Two Pointers",
	"Sliding Window",
	"Fast & Slow Pointers",
	"Hash Map",
	"Binary Search",
	"BFS",
	"DFS",
	"Dynamic Programming",
	"Greedy",
	"Heap",
	"Union Find",
}

// Update handles updates for the pattern selection screen
func (m Model) updatePatterns(msg tea.Msg) (Model, tea.Cmd) {
	// Initialize patterns from problems if not done yet
	if len(m.patterns.patterns) == 0 && len(m.allProblems) > 0 {
		m.patterns.patterns = GetAvailablePatterns(m.allProblems)
	}
	// Use default patterns if no problems loaded
	if len(m.patterns.patterns) == 0 {
		m.patterns.patterns = defaultPatterns
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.patterns.selectedIndex > 0 {
				m.patterns.selectedIndex--
			}
		case "down", "j":
			if m.patterns.selectedIndex < len(m.patterns.patterns)-1 {
				m.patterns.selectedIndex++
			}
		case "enter", "right", "l":
			m.patterns.selectedPattern = m.patterns.patterns[m.patterns.selectedIndex]
			m.problems.pattern = m.patterns.selectedPattern
			return m.navigate(StateProblemList), loadProblemsForPattern(m.patterns.selectedPattern)
		}
	}
	return m, nil
}

// View renders the pattern selection screen
func (m Model) viewPatterns() string {
	var b strings.Builder
	
	// Initialize patterns if needed
	patternList := m.patterns.patterns
	if len(patternList) == 0 && len(m.allProblems) > 0 {
		patternList = GetAvailablePatterns(m.allProblems)
	}
	if len(patternList) == 0 {
		patternList = defaultPatterns
	}
	
	// Title
	b.WriteString(titleStyle.Render("Select a Pattern"))
	b.WriteString("\n\n")
	
	// Pattern list
	for i, pattern := range patternList {
		cursor := "  "
		if i == m.patterns.selectedIndex {
			cursor = cursorStyle.Render("> ")
			pattern = selectedItemStyle.Render(pattern)
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, pattern))
	}
	
	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Back"))
	
	return b.String()
}

// GetAvailablePatterns extracts unique patterns from all problems
func GetAvailablePatterns(problems []problem.Problem) []string {
	return problem.GetPatterns(problems)
}