// Statistics tracking and analysis

package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SessionStats represents statistics for a single session
type SessionStats struct {
	ProblemID    string        `json:"problem_id"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	Solved       bool          `json:"solved"`
	Mode         string        `json:"mode"`
	HintsUsed    bool          `json:"hints_used"`
	SolutionUsed bool          `json:"solution_used"`
	Patterns     []string      `json:"patterns"`
	Difficulty   string        `json:"difficulty"`
}

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

// Helper functions

// loadAllSessions loads all session stats from files
func loadAllSessions() ([]SessionStats, error) {
	var sessions []SessionStats

	statsDir := filepath.Join(getConfigDir(), "stats")

	// Check if directory exists
	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		return sessions, nil
	}

	// Read all stats files
	files, err := os.ReadDir(statsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !isStatsFile(file.Name()) {
			continue
		}

		// Read file
		data, err := os.ReadFile(filepath.Join(statsDir, file.Name()))
		if err != nil {
			return nil, err
		}

		var session SessionStats
		if err := json.Unmarshal(data, &session); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// formatDuration formats a duration as a string (HH:MM:SS)
func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds / 60) % 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// getConfigDir returns the configuration directory
// Exported as variable for testing
var getConfigDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}