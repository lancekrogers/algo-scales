package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/daily"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// loadProblems loads all available problems
func loadProblems() tea.Cmd {
	return func() tea.Msg {
		problems, err := problem.LoadLocalProblems()
		if err != nil {
			return problemsErrorMsg{err: err}
		}
		return problemsLoadedMsg{problems: problems}
	}
}

// loadProblemsForPattern loads problems for a specific pattern
func loadProblemsForPattern(pattern string) tea.Cmd {
	return func() tea.Msg {
		// Convert pattern name to directory format
		dirName := strings.ToLower(strings.ReplaceAll(pattern, " ", "-"))
		dirName = strings.ReplaceAll(dirName, "&", "")
		
		problems, err := problem.LoadLocalProblems()
		if err != nil {
			return problemsErrorMsg{err: err}
		}
		
		// Filter problems by pattern
		filtered := make([]problem.Problem, 0)
		for _, p := range problems {
			for _, tag := range p.Patterns {
				if tag == pattern {
					filtered = append(filtered, p)
					break
				}
			}
		}
		
		return problemsLoadedMsg{problems: filtered}
	}
}

// loadStats loads user statistics
func loadStats() tea.Cmd {
	return func() tea.Msg {
		summary, err := stats.GetSummary()
		if err != nil {
			return statsErrorMsg{err: err}
		}
		return statsLoadedMsg{stats: *summary}
	}
}

// loadDailyScale loads the daily scale challenge
func loadDailyScale() tea.Cmd {
	return func() tea.Msg {
		progress, err := daily.LoadProgress()
		if err != nil {
			// Handle error or create default progress
			progress = daily.ScaleProgress{
				Current:       0,
				LastPracticed: time.Time{},
				Completed:     []string{},
				Streak:        0,
				LongestStreak: 0,
			}
		}
		
		scale := daily.GetNextScale(progress.Completed)
		if scale == nil {
			// All scales completed, start over
			progress.Completed = []string{}
			scale = daily.GetNextScale(progress.Completed)
		}
		
		return dailyScaleLoadedMsg{
			scale:    scale.Pattern,
			progress: progress,
		}
	}
}

// sessionTick generates periodic ticks for the session timer
func sessionTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return sessionTickMsg{}
	})
}

// startSession creates a command to start a new session
func startSession(prob problem.Problem) tea.Cmd {
	return func() tea.Msg {
		// Create session ID
		sessionID := fmt.Sprintf("session-%s-%d", prob.ID, time.Now().Unix())
		return sessionStartedMsg{
			sessionID: sessionID,
		}
	}
}