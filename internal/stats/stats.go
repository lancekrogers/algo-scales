// Statistics tracking and analysis

package stats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

// RecordSession records a session's statistics
func RecordSession(stats SessionStats) error {
	// Get the stats directory
	statsDir := filepath.Join(getConfigDir(), "stats")
	if err := os.MkdirAll(statsDir, 0755); err != nil {
		return err
	}

	// Generate a unique filename
	filename := fmt.Sprintf("session_%s_%s.json", stats.ProblemID, stats.StartTime.Format("20060102_150405"))
	statsFile := filepath.Join(statsDir, filename)

	// Save stats to file
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(statsFile, data, 0644)
}

// GetSummary returns summary statistics
// Exported as variable for testing
var GetSummary = func() (*Summary, error) {
	// Load all session stats
	sessions, err := loadAllSessions()
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return &Summary{
			TotalAttempted: 0,
			TotalSolved:    0,
			SuccessRate:    0,
			AvgSolveTime:   "00:00",
			FastestSolve: struct {
				ProblemID string `json:"problem_id"`
				Time      string `json:"time"`
			}{
				ProblemID: "",
				Time:      "",
			},
			MostChallenging: struct {
				ProblemID string `json:"problem_id"`
				Attempts  int    `json:"attempts"`
			}{
				ProblemID: "",
				Attempts:  0,
			},
		}, nil
	}

	// Calculate summary stats
	summary := &Summary{}
	summary.TotalAttempted = len(sessions)

	var totalSolveTime time.Duration
	var solvedCount int

	// Track problem attempts
	problemAttempts := make(map[string]int)
	problemSolved := make(map[string]bool)

	// Track fastest solve
	var fastestTime time.Duration
	var fastestProblem string

	for _, session := range sessions {
		problemAttempts[session.ProblemID]++

		if session.Solved {
			solvedCount++
			totalSolveTime += session.Duration

			// Update fastest solve if this is faster
			if fastestProblem == "" || session.Duration < fastestTime {
				fastestTime = session.Duration
				fastestProblem = session.ProblemID
			}

			problemSolved[session.ProblemID] = true
		}
	}

	summary.TotalSolved = solvedCount
	if solvedCount > 0 {
		avgTime := totalSolveTime / time.Duration(solvedCount)
		summary.AvgSolveTime = formatDuration(avgTime)
	}

	if summary.TotalAttempted > 0 {
		summary.SuccessRate = float64(solvedCount) / float64(summary.TotalAttempted) * 100
	}

	// Find most challenging problem
	var maxAttempts int
	var mostChallengingProblem string

	for problem, attempts := range problemAttempts {
		if attempts > maxAttempts && !problemSolved[problem] {
			maxAttempts = attempts
			mostChallengingProblem = problem
		}
	}

	if fastestProblem != "" {
		summary.FastestSolve.ProblemID = fastestProblem
		summary.FastestSolve.Time = formatDuration(fastestTime)
	}

	if mostChallengingProblem != "" {
		summary.MostChallenging.ProblemID = mostChallengingProblem
		summary.MostChallenging.Attempts = maxAttempts
	}

	return summary, nil
}

// GetByPattern returns statistics by pattern
// Exported as variable for testing
var GetByPattern = func() (map[string]PatternStats, error) {
	// Load all session stats
	sessions, err := loadAllSessions()
	if err != nil {
		return nil, err
	}

	// Group by pattern
	patternStats := make(map[string]PatternStats)

	// Track pattern solve times
	patternTimes := make(map[string][]time.Duration)

	for _, session := range sessions {
		for _, pattern := range session.Patterns {
			// Initialize if not exists
			if _, ok := patternStats[pattern]; !ok {
				patternStats[pattern] = PatternStats{
					Pattern: pattern,
				}
			}

			// Update stats
			stats := patternStats[pattern]
			stats.Attempted++

			if session.Solved {
				stats.Solved++
				patternTimes[pattern] = append(patternTimes[pattern], session.Duration)
			}

			// Calculate success rate
			stats.SuccessRate = float64(stats.Solved) / float64(stats.Attempted) * 100

			// Update in map
			patternStats[pattern] = stats
		}
	}

	// Calculate average solve times
	for pattern, times := range patternTimes {
		if len(times) > 0 {
			var total time.Duration
			for _, t := range times {
				total += t
			}
			avgTime := total / time.Duration(len(times))

			stats := patternStats[pattern]
			stats.AvgTime = formatDuration(avgTime)
			patternStats[pattern] = stats
		}
	}

	return patternStats, nil
}

// GetTrends returns usage trends over time
// Exported as variable for testing
var GetTrends = func() (*Trends, error) {
	// Load all session stats
	sessions, err := loadAllSessions()
	if err != nil {
		return nil, err
	}

	// Sort sessions by start time
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.Before(sessions[j].StartTime)
	})

	// Group by day and week
	trends := &Trends{}
	dailyStats := make(map[string]struct {
		Solved    int
		TotalTime time.Duration
		Count     int
	})

	weeklyStats := make(map[string]struct {
		Solved    int
		Attempted int
		StartDate time.Time
		EndDate   time.Time
	})

	for _, session := range sessions {
		// Daily stats
		dateStr := session.StartTime.Format("2006-01-02")
		daily := dailyStats[dateStr]
		if session.Solved {
			daily.Solved++
			daily.TotalTime += session.Duration
			daily.Count++
		}
		dailyStats[dateStr] = daily

		// Weekly stats
		yearWeek := getYearWeek(session.StartTime)
		weekly := weeklyStats[yearWeek]
		if weekly.StartDate.IsZero() || session.StartTime.Before(weekly.StartDate) {
			weekly.StartDate = startOfWeek(session.StartTime)
		}
		if weekly.EndDate.IsZero() || session.StartTime.After(weekly.EndDate) {
			weekly.EndDate = endOfWeek(session.StartTime)
		}
		weekly.Attempted++
		if session.Solved {
			weekly.Solved++
		}
		weeklyStats[yearWeek] = weekly
	}

	// Convert maps to slices for the last 7 days
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		dailyTrend := DailyTrend{
			Date:   dateStr,
			Solved: 0,
		}

		// Add data if available
		if stats, ok := dailyStats[dateStr]; ok {
			dailyTrend.Solved = stats.Solved
			if stats.Count > 0 {
				avgTime := stats.TotalTime / time.Duration(stats.Count)
				dailyTrend.AvgTime = formatDuration(avgTime)
			}
		}

		trends.Daily = append(trends.Daily, dailyTrend)
	}

	// Sort daily trends by date
	sort.Slice(trends.Daily, func(i, j int) bool {
		return trends.Daily[i].Date < trends.Daily[j].Date
	})

	// Convert weekly stats map to slice
	for _, stats := range weeklyStats {
		weeklyTrend := WeeklyTrend{
			StartDate:   stats.StartDate.Format("2006-01-02"),
			EndDate:     stats.EndDate.Format("2006-01-02"),
			Solved:      stats.Solved,
			SuccessRate: 0,
		}

		if stats.Attempted > 0 {
			weeklyTrend.SuccessRate = float64(stats.Solved) / float64(stats.Attempted) * 100
		}

		trends.Weekly = append(trends.Weekly, weeklyTrend)
	}

	// Sort weekly trends by start date
	sort.Slice(trends.Weekly, func(i, j int) bool {
		return trends.Weekly[i].StartDate < trends.Weekly[j].StartDate
	})

	// Limit to last 4 weeks
	if len(trends.Weekly) > 4 {
		trends.Weekly = trends.Weekly[len(trends.Weekly)-4:]
	}

	return trends, nil
}

// Reset resets all statistics
// Exported as variable for testing
var Reset = func() error {
	statsDir := filepath.Join(getConfigDir(), "stats")

	// Create directory if it doesn't exist
	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		return nil // Nothing to reset
	}

	// Remove all files in the directory
	files, err := ioutil.ReadDir(statsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(filepath.Join(statsDir, file.Name())); err != nil {
			return err
		}
	}

	return nil
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
	files, err := ioutil.ReadDir(statsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !isStatsFile(file.Name()) {
			continue
		}

		// Read file
		data, err := ioutil.ReadFile(filepath.Join(statsDir, file.Name()))
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

// isStatsFile checks if a filename is a stats file
func isStatsFile(filename string) bool {
	return len(filename) > 8 && filename[:8] == "session_" && filepath.Ext(filename) == ".json"
}

// formatDuration formats a duration as a string (MM:SS)
func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

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

// getConfigDir returns the configuration directory
// Exported as variable for testing
var getConfigDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}
