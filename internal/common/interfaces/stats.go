// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import "context"

// Summary represents summary statistics
type Summary struct {
	TotalAttempted int     `json:"total_attempted"`
	TotalSolved    int     `json:"total_solved"`
	AvgSolveTime   string  `json:"avg_solve_time"`
	SuccessRate    float64 `json:"success_rate"`
	FastestSolve   struct {
		ProblemID string `json:"problem_id"`
		Time      string `json:"time"`
	} `json:"fastest_solve"`
	MostChallenging struct {
		ProblemID string `json:"problem_id"`
		Attempts  int    `json:"attempts"`
	} `json:"most_challenging"`
}

// PatternStats represents statistics for a pattern
type PatternStats struct {
	Pattern     string  `json:"pattern"`
	Attempted   int     `json:"attempted"`
	Solved      int     `json:"solved"`
	SuccessRate float64 `json:"success_rate"`
	AvgTime     string  `json:"avg_time"`
}

// Trends represents trends over time
type Trends struct {
	Daily  []DailyTrend  `json:"daily"`
	Weekly []WeeklyTrend `json:"weekly"`
}

// DailyTrend represents a daily trend
type DailyTrend struct {
	Date    string `json:"date"`
	Solved  int    `json:"solved"`
	AvgTime string `json:"avg_time"`
}

// WeeklyTrend represents a weekly trend
type WeeklyTrend struct {
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Solved      int     `json:"solved"`
	SuccessRate float64 `json:"success_rate"`
}

// OverallStats contains general overview statistics
type OverallStats struct {
	Summary *Summary
	Trends  *Trends
}

// StatsService defines the interface for accessing algorithm practice statistics
type StatsService interface {
	// RecordSession records a session's statistics
	RecordSession(ctx context.Context, sessionStats SessionStats) error
	
	// GetSummary returns summary statistics
	GetSummary(ctx context.Context) (*Summary, error)
	
	// GetByPattern returns statistics by pattern
	GetByPattern(ctx context.Context) (map[string]PatternStats, error)
	
	// GetTrends returns usage trends over time
	GetTrends(ctx context.Context) (*Trends, error)
	
	// Reset resets all statistics
	Reset(ctx context.Context) error
	
	// GetAllSessions returns all recorded sessions
	GetAllSessions(ctx context.Context) ([]SessionStats, error)
}

// DifficultyStats represents statistics by difficulty level
type DifficultyStats struct {
	Difficulty  string  `json:"difficulty"`
	Attempted   int     `json:"attempted"`
	Solved      int     `json:"solved"`
	SuccessRate float64 `json:"success_rate"`
	AvgTime     string  `json:"avg_time"`
}

// LanguageStats represents statistics by programming language
type LanguageStats struct {
	Language    string  `json:"language"`
	Attempted   int     `json:"attempted"`
	Solved      int     `json:"solved"`
	SuccessRate float64 `json:"success_rate"`
	AvgTime     string  `json:"avg_time"`
}

// DailyStats represents statistics for daily practice
type DailyStats struct {
	Date          string   `json:"date"`
	ProblemsToday int      `json:"problems_today"`
	StreakDays    int      `json:"streak_days"`
	PatternsToday []string `json:"patterns_today"`
	Complete      bool     `json:"complete"`
}

// StatsStorage defines the interface for storing and retrieving statistics
type StatsStorage interface {
	// SaveSession saves a session's statistics
	SaveSession(ctx context.Context, session SessionStats) error
	
	// LoadAllSessions loads all session statistics
	LoadAllSessions(ctx context.Context) ([]SessionStats, error)
	
	// ClearAllSessions removes all session statistics
	ClearAllSessions(ctx context.Context) error
}