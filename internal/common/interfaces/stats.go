// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import (
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// StatsService defines the interface for accessing algorithm practice statistics
type StatsService interface {
	// RecordSession records a session's statistics
	RecordSession(sessionStats stats.SessionStats) error
	
	// GetSummary returns summary statistics
	GetSummary() (*stats.Summary, error)
	
	// GetByPattern returns statistics by pattern
	GetByPattern() (map[string]stats.PatternStats, error)
	
	// GetTrends returns usage trends over time
	GetTrends() (*stats.Trends, error)
	
	// Reset resets all statistics
	Reset() error
	
	// GetAllSessions returns all recorded sessions
	GetAllSessions() ([]stats.SessionStats, error)
}

// StatsStorage defines the interface for storing and retrieving statistics
type StatsStorage interface {
	// SaveSession saves a session's statistics
	SaveSession(session stats.SessionStats) error
	
	// LoadAllSessions loads all session statistics
	LoadAllSessions() ([]stats.SessionStats, error)
	
	// ClearAllSessions removes all session statistics
	ClearAllSessions() error
}