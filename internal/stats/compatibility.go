package stats

// This file provides backward compatibility functions
// to maintain the same exported API while using the new service internally

import (
	"fmt"
	"path/filepath"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// DefaultService is the default stats service instance
var DefaultService = NewService()

// RecordSession records a session's statistics
var RecordSession = func(stats SessionStats) error {
	// Convert to interface type
	interfaceStats := interfaces.SessionStats{
		ProblemID:    stats.ProblemID,
		StartTime:    stats.StartTime,
		EndTime:      stats.EndTime,
		Duration:     stats.Duration,
		Solved:       stats.Solved,
		Mode:         stats.Mode,
		HintsUsed:    stats.HintsUsed,
		SolutionUsed: stats.SolutionUsed,
		Patterns:     stats.Patterns,
		Difficulty:   stats.Difficulty,
	}
	return DefaultService.RecordSession(interfaceStats)
}

// GetSummary returns summary statistics
var GetSummary = func() (*Summary, error) {
	interfaceSummary, err := DefaultService.GetSummary()
	if err != nil {
		return nil, err
	}
	
	// Convert to local type
	localSummary := &Summary{
		TotalAttempted: interfaceSummary.TotalAttempted,
		TotalSolved:    interfaceSummary.TotalSolved,
		AvgSolveTime:   interfaceSummary.AvgSolveTime,
		SuccessRate:    interfaceSummary.SuccessRate,
		FastestSolve:   interfaceSummary.FastestSolve,
		MostChallenging: interfaceSummary.MostChallenging,
	}
	return localSummary, nil
}

// GetByPattern returns statistics by pattern
var GetByPattern = func() (map[string]PatternStats, error) {
	interfaceStats, err := DefaultService.GetByPattern()
	if err != nil {
		return nil, err
	}
	
	// Convert to local type
	localStats := make(map[string]PatternStats)
	for pattern, stats := range interfaceStats {
		localStats[pattern] = PatternStats{
			Pattern:     stats.Pattern,
			Attempted:   stats.Attempted,
			Solved:      stats.Solved,
			SuccessRate: stats.SuccessRate,
			AvgTime:     stats.AvgTime,
		}
	}
	return localStats, nil
}

// GetTrends returns usage trends over time
var GetTrends = func() (*Trends, error) {
	interfaceTrends, err := DefaultService.GetTrends()
	if err != nil {
		return nil, err
	}
	
	// Convert to local type
	localTrends := &Trends{
		Daily:  make([]DailyTrend, len(interfaceTrends.Daily)),
		Weekly: make([]WeeklyTrend, len(interfaceTrends.Weekly)),
	}
	
	for i, daily := range interfaceTrends.Daily {
		localTrends.Daily[i] = DailyTrend{
			Date:    daily.Date,
			Solved:  daily.Solved,
			AvgTime: daily.AvgTime,
		}
	}
	
	for i, weekly := range interfaceTrends.Weekly {
		localTrends.Weekly[i] = WeeklyTrend{
			StartDate:   weekly.StartDate,
			EndDate:     weekly.EndDate,
			Solved:      weekly.Solved,
			SuccessRate: weekly.SuccessRate,
		}
	}
	
	return localTrends, nil
}

// Reset resets all statistics
var Reset = func() error {
	return DefaultService.Reset()
}

// GetAllSessions returns all recorded sessions
var GetAllSessions = func() ([]SessionStats, error) {
	interfaceSessions, err := DefaultService.GetAllSessions()
	if err != nil {
		return nil, err
	}
	
	// Convert to local type
	localSessions := make([]SessionStats, len(interfaceSessions))
	for i, s := range interfaceSessions {
		localSessions[i] = SessionStats{
			ProblemID:    s.ProblemID,
			StartTime:    s.StartTime,
			EndTime:      s.EndTime,
			Duration:     s.Duration,
			Solved:       s.Solved,
			Mode:         s.Mode,
			HintsUsed:    s.HintsUsed,
			SolutionUsed: s.SolutionUsed,
			Patterns:     s.Patterns,
			Difficulty:   s.Difficulty,
		}
	}
	return localSessions, nil
}

// Helper functions that remain as internal utilities

// getYearWeek returns a string representing the year and week
func getYearWeek(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

// startOfWeek returns the start of the week for a time
func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekday-- // Adjust to make Monday the first day

	return time.Date(t.Year(), t.Month(), t.Day()-weekday, 0, 0, 0, 0, t.Location())
}

// endOfWeek returns the end of the week for a time
func endOfWeek(t time.Time) time.Time {
	return startOfWeek(t).AddDate(0, 0, 6)
}

// isStatsFile checks if a filename is a stats file
func isStatsFile(filename string) bool {
	return len(filename) > 8 && filename[:8] == "session_" && filepath.Ext(filename) == ".json"
}