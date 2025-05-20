package stats

// This file provides backward compatibility functions
// to maintain the same exported API while using the new service internally

import (
	"fmt"
	"path/filepath"
	"time"
)

// DefaultService is the default stats service instance
var DefaultService = NewService()

// RecordSession records a session's statistics
var RecordSession = func(stats SessionStats) error {
	return DefaultService.RecordSession(stats)
}

// GetSummary returns summary statistics
var GetSummary = func() (*Summary, error) {
	return DefaultService.GetSummary()
}

// GetByPattern returns statistics by pattern
var GetByPattern = func() (map[string]PatternStats, error) {
	return DefaultService.GetByPattern()
}

// GetTrends returns usage trends over time
var GetTrends = func() (*Trends, error) {
	return DefaultService.GetTrends()
}

// Reset resets all statistics
var Reset = func() error {
	return DefaultService.Reset()
}

// GetAllSessions returns all recorded sessions
var GetAllSessions = func() ([]SessionStats, error) {
	return DefaultService.GetAllSessions()
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