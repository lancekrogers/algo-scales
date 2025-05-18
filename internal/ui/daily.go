package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/daily"
)

// Update handles updates for the daily scale screen
func (m Model) updateDaily(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dailyScaleLoadedMsg:
		m.daily.currentScale = msg.scale
		m.daily.progress = msg.progress
		m.daily.loading = false
		
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Find problems for this scale pattern
			pattern := scaleToPattern(m.daily.currentScale)
			return m.navigate(StateProblemList), loadProblemsForPattern(pattern)
		case "n":
			// Skip to next scale
			if p, ok := m.daily.progress.(daily.ScaleProgress); ok {
				p.Completed = append(p.Completed, m.daily.currentScale)
				return m, loadDailyScale()
			}
		case "r":
			// Reset progress
			m.daily.loading = true
			return m, resetDailyProgress()
		}
	}
	
	return m, nil
}

// View renders the daily scale screen
func (m Model) viewDaily() string {
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62")).
		MarginBottom(2)
	
	b.WriteString(titleStyle.Render("🎵 Daily Scales"))
	b.WriteString("\n\n")
	
	if m.daily.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
		b.WriteString(loadingStyle.Render("Loading daily scale..."))
		return b.String()
	}
	
	// Scale information
	scaleBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)
	
	scale := getScaleInfo(m.daily.currentScale)
	scaleContent := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		lipgloss.NewStyle().Bold(true).Render(scale.MusicalName),
		scale.Pattern,
		lipgloss.NewStyle().Italic(true).Render(scale.Description),
	)
	
	b.WriteString(scaleBoxStyle.Render(scaleContent))
	b.WriteString("\n\n")
	
	// Progress information
	if p, ok := m.daily.progress.(daily.ScaleProgress); ok {
		progressStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))
		
		// Streak information
		streakStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))
		
		if p.Streak > 0 {
			streakText := fmt.Sprintf("🔥 %d day streak!", p.Streak)
			b.WriteString(streakStyle.Render(streakText))
			b.WriteString("\n\n")
		}
		
		// Completion progress
		progressText := fmt.Sprintf("Completed: %d/12 scales today", len(p.Completed))
		b.WriteString(progressStyle.Render(progressText))
		
		// Last practice
		if !p.LastPracticed.IsZero() {
			lastPracticeText := fmt.Sprintf("\nLast practice: %s", 
				p.LastPracticed.Format("Jan 2, 3:04 PM"))
			b.WriteString(progressStyle.Render(lastPracticeText))
		}
		
		b.WriteString("\n\n")
	}
	
	// Action bar
	actionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	
	actions := []string{
		"Enter: Practice this scale",
		"n: Next scale",
		"r: Reset progress",
		"Esc: Back",
	}
	
	b.WriteString(actionStyle.Render(strings.Join(actions, " • ")))
	
	return b.String()
}

// scaleToPattern converts a scale pattern to a problem pattern
func scaleToPattern(scale string) string {
	// Map scale names to pattern names
	mapping := map[string]string{
		"sliding-window":      "Sliding Window",
		"two-pointers":        "Two Pointers",
		"fast-slow-pointers":  "Fast & Slow Pointers",
		"hash-map":            "Hash Maps",
		"binary-search":       "Binary Search",
		"dfs":                 "DFS",
		"bfs":                 "BFS",
		"dynamic-programming": "Dynamic Programming",
		"greedy":              "Greedy",
		"union-find":          "Union Find",
		"heap":                "Heap/Priority Queue",
	}
	
	if pattern, ok := mapping[scale]; ok {
		return pattern
	}
	return scale
}

// getScaleInfo returns scale information for display
func getScaleInfo(pattern string) daily.Scale {
	// Map of pattern to scale info
	scales := map[string]daily.Scale{
		"sliding-window": {
			Pattern:     "Sliding Window",
			MusicalName: "C Major",
			Description: "The fundamental scale, elegant and versatile",
		},
		"two-pointers": {
			Pattern:     "Two Pointers",
			MusicalName: "G Major",
			Description: "Balanced and efficient, the workhorse of array manipulation",
		},
		"fast-slow-pointers": {
			Pattern:     "Fast & Slow Pointers",
			MusicalName: "D Major",
			Description: "The cycle detector, bright and revealing",
		},
		"hash-map": {
			Pattern:     "Hash Maps",
			MusicalName: "A Major",
			Description: "The lookup accelerator, crisp and direct",
		},
		"binary-search": {
			Pattern:     "Binary Search",
			MusicalName: "E Major",
			Description: "The divide and conquer virtuoso",
		},
		"dfs": {
			Pattern:     "DFS",
			MusicalName: "B Major",
			Description: "The deep explorer, methodical and thorough",
		},
		"bfs": {
			Pattern:     "BFS",
			MusicalName: "F# Major",
			Description: "The level-wise traverser, systematic and complete",
		},
		"dynamic-programming": {
			Pattern:     "Dynamic Programming",
			MusicalName: "Db Major",
			Description: "The optimization maestro, building upon the past",
		},
		"greedy": {
			Pattern:     "Greedy",
			MusicalName: "Ab Major",
			Description: "The local optimizer, decisive and swift",
		},
		"union-find": {
			Pattern:     "Union Find",
			MusicalName: "Eb Major",
			Description: "The connection tracker, uniting disparate elements",
		},
		"heap": {
			Pattern:     "Heap/Priority Queue",
			MusicalName: "Bb Major",
			Description: "The priority manager, always serving the most important",
		},
	}
	
	if scale, ok := scales[pattern]; ok {
		return scale
	}
	
	return daily.Scale{
		Pattern:     pattern,
		MusicalName: "Unknown Scale",
		Description: "Mystery pattern",
	}
}

// resetDailyProgress resets the daily progress
func resetDailyProgress() tea.Cmd {
	return func() tea.Msg {
		// Reset progress
		progress := daily.ScaleProgress{
			Current:       0,
			LastPracticed: time.Now(),
			Completed:     []string{},
			Streak:        0,
			LongestStreak: 0,
		}
		
		// Save and reload
		daily.SaveProgress(progress)
		
		scale := daily.GetNextScale(progress.Completed)
		return dailyScaleLoadedMsg{
			scale:    scale.Pattern,
			progress: progress,
		}
	}
}