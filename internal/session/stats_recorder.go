package session

import (
	"context"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// SessionStatsRecorderImpl implements the SessionStatsRecorder interface
type SessionStatsRecorderImpl struct {
	statsService interfaces.StatsService
}

// NewSessionStatsRecorder creates a new session stats recorder
func NewSessionStatsRecorder(statsService interfaces.StatsService) interfaces.SessionStatsRecorder {
	return &SessionStatsRecorderImpl{
		statsService: statsService,
	}
}

// RecordSession records session statistics
func (sr *SessionStatsRecorderImpl) RecordSession(ctx context.Context, sessionStats interfaces.SessionStats) error {
	// Convert to stats package format if needed
	statsSession := stats.SessionStats{
		ProblemID:    sessionStats.ProblemID,
		StartTime:    sessionStats.StartTime,
		EndTime:      sessionStats.EndTime,
		Duration:     sessionStats.Duration,
		Solved:       sessionStats.Solved,
		Mode:         sessionStats.Mode,
		HintsUsed:    sessionStats.HintsUsed,
		SolutionUsed: sessionStats.SolutionUsed,
		Patterns:     sessionStats.Patterns,
		Difficulty:   sessionStats.Difficulty,
	}
	
	// Use the legacy function for now to maintain compatibility
	return stats.RecordSession(statsSession)
}