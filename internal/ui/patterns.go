package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Pattern categories
var patterns = []string{
	"Two Pointers",
	"Sliding Window",
	"Fast & Slow Pointers",
	"Hash Maps",
	"Binary Search",
	"BFS",
	"DFS",
	"Dynamic Programming",
	"Greedy",
	"Heap/Priority Queue",
	"Union Find",
}

// Update handles updates for the pattern selection screen
func (m Model) updatePatterns(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.patterns.selectedIndex > 0 {
				m.patterns.selectedIndex--
			}
		case "down", "j":
			if m.patterns.selectedIndex < len(patterns)-1 {
				m.patterns.selectedIndex++
			}
		case "enter", "right", "l":
			m.patterns.selectedPattern = patterns[m.patterns.selectedIndex]
			m.problems.pattern = m.patterns.selectedPattern
			return m.navigate(StateProblemList), loadProblemsForPattern(m.patterns.selectedPattern)
		}
	}
	return m, nil
}

// View renders the pattern selection screen
func (m Model) viewPatterns() string {
	var b strings.Builder
	
	// Title
	
	b.WriteString(titleStyle.Render("Select a Pattern"))
	b.WriteString("\n\n")
	
	// Pattern list
	for i, pattern := range patterns {
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