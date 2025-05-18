package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/common/config"
	"github.com/lancekrogers/algo-scales/internal/daily"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// loadProblems loads all available problems
func loadProblems() tea.Cmd {
	return tea.Sequence(
		startLoading("Loading problems..."),
		func() tea.Msg {
			// Simulate a small delay for visual effect
			time.Sleep(200 * time.Millisecond)
			
			problems, err := problem.LoadLocalProblems()
			if err != nil {
				return problemsErrorMsg{err: err}
			}
			return problemsLoadedMsg{problems: problems}
		},
		stopLoading(),
	)
}

// loadConfig loads the user configuration
func loadConfig() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.LoadConfig()
		if err != nil {
			// Use default config on error
			cfg = config.DefaultConfig()
		}
		return configLoadedMsg{config: cfg}
	}
}

// loadProblemsForPattern loads problems for a specific pattern
func loadProblemsForPattern(pattern string) tea.Cmd {
	return func() tea.Msg {
		problems, err := problem.LoadLocalProblems()
		if err != nil {
			return problemsErrorMsg{err: err}
		}
		
		// Convert pattern name to directory/tag format for matching
		// e.g., "Two Pointers" -> "two-pointers"
		normalizedPattern := strings.ToLower(strings.ReplaceAll(pattern, " ", "-"))
		normalizedPattern = strings.ReplaceAll(normalizedPattern, "&", "")
		normalizedPattern = strings.ReplaceAll(normalizedPattern, "/", "-")
		
		// Filter problems by pattern
		filtered := make([]problem.Problem, 0)
		for _, p := range problems {
			for _, tag := range p.Patterns {
				// Also normalize the tag from the problem
				normalizedTag := strings.ToLower(tag)
				if normalizedTag == normalizedPattern || tag == pattern {
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

// Loading commands
func startLoading(message string) tea.Cmd {
	return func() tea.Msg {
		return startLoadingMsg{message: message}
	}
}

func stopLoading() tea.Cmd {
	return func() tea.Msg {
		return stopLoadingMsg{}
	}
}

// Loading messages
type startLoadingMsg struct {
	message string
}

type stopLoadingMsg struct{}