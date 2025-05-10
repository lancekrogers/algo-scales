// Tests for stats command

package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/lancekrogers/algo-scales/internal/stats"
	"github.com/stretchr/testify/assert"
)

// Mock stats.GetSummary for testing
func mockGetSummary(summary *stats.Summary, err error) func() {
	original := stats.GetSummary
	stats.GetSummary = func() (*stats.Summary, error) {
		return summary, err
	}
	return func() {
		stats.GetSummary = original
	}
}

// Mock stats.GetByPattern for testing
func mockGetByPattern(patternStats map[string]stats.PatternStats, err error) func() {
	original := stats.GetByPattern
	stats.GetByPattern = func() (map[string]stats.PatternStats, error) {
		return patternStats, err
	}
	return func() {
		stats.GetByPattern = original
	}
}

// Mock stats.GetTrends for testing
func mockGetTrends(trends *stats.Trends, err error) func() {
	original := stats.GetTrends
	stats.GetTrends = func() (*stats.Trends, error) {
		return trends, err
	}
	return func() {
		stats.GetTrends = original
	}
}

// Mock stats.Reset for testing
func mockReset(err error) func() {
	original := stats.Reset
	stats.Reset = func() error {
		return err
	}
	return func() {
		stats.Reset = original
	}
}

func TestStatsCommand(t *testing.T) {
	t.Run("Summary", func(t *testing.T) {
		// Create a sample summary
		summary := &stats.Summary{
			TotalAttempted: 10,
			TotalSolved:    7,
			AvgSolveTime:   "15:30",
			SuccessRate:    70.0,
			FastestSolve: struct {
				ProblemID string `json:"problem_id"`
				Time      string `json:"time"`
			}{
				ProblemID: "test-problem",
				Time:      "05:45",
			},
			MostChallenging: struct {
				ProblemID string `json:"problem_id"`
				Attempts  int    `json:"attempts"`
			}{
				ProblemID: "hard-problem",
				Attempts:  3,
			},
		}

		// Mock GetSummary
		restore := mockGetSummary(summary, nil)
		defer restore()

		// Execute stats command
		output, err := executeCommand(rootCmd, "stats")
		assert.NoError(t, err)

		// Check output contains summary info
		assert.Contains(t, output, "10")           // TotalAttempted
		assert.Contains(t, output, "7")            // TotalSolved
		assert.Contains(t, output, "15:30")        // AvgSolveTime
		assert.Contains(t, output, "test-problem") // FastestSolve
		assert.Contains(t, output, "05:45")        // FastestSolve time
		assert.Contains(t, output, "hard-problem") // MostChallenging
		assert.Contains(t, output, "3")            // MostChallenging attempts
	})

	t.Run("PatternStats", func(t *testing.T) {
		// Create sample pattern stats
		patternStats := map[string]stats.PatternStats{
			"hash-map": {
				Pattern:     "hash-map",
				Attempted:   5,
				Solved:      4,
				SuccessRate: 80.0,
				AvgTime:     "10:15",
			},
			"two-pointers": {
				Pattern:     "two-pointers",
				Attempted:   3,
				Solved:      2,
				SuccessRate: 66.7,
				AvgTime:     "12:30",
			},
		}

		// Mock GetByPattern
		restore := mockGetByPattern(patternStats, nil)
		defer restore()

		// Execute patterns command
		output, err := executeCommand(rootCmd, "stats", "patterns")
		assert.NoError(t, err)

		// Check output contains pattern stats
		assert.Contains(t, output, "hash-map")
		assert.Contains(t, output, "two-pointers")
		assert.Contains(t, output, "5")     // hash-map attempted
		assert.Contains(t, output, "4")     // hash-map solved
		assert.Contains(t, output, "80")    // hash-map success rate
		assert.Contains(t, output, "10:15") // hash-map avg time
		assert.Contains(t, output, "3")     // two-pointers attempted
		assert.Contains(t, output, "2")     // two-pointers solved
	})

	t.Run("Trends", func(t *testing.T) {
		// Create sample trends
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		lastWeek := now.AddDate(0, 0, -7)

		trends := &stats.Trends{
			Daily: []stats.DailyTrend{
				{
					Date:    now.Format("2006-01-02"),
					Solved:  3,
					AvgTime: "11:30",
				},
				{
					Date:    yesterday.Format("2006-01-02"),
					Solved:  2,
					AvgTime: "15:45",
				},
			},
			Weekly: []stats.WeeklyTrend{
				{
					StartDate:   now.AddDate(0, 0, -6).Format("2006-01-02"),
					EndDate:     now.Format("2006-01-02"),
					Solved:      5,
					SuccessRate: 71.4,
				},
				{
					StartDate:   lastWeek.AddDate(0, 0, -6).Format("2006-01-02"),
					EndDate:     lastWeek.Format("2006-01-02"),
					Solved:      4,
					SuccessRate: 66.7,
				},
			},
		}

		// Mock GetTrends
		restore := mockGetTrends(trends, nil)
		defer restore()

		// Execute trends command
		output, err := executeCommand(rootCmd, "stats", "trends")
		assert.NoError(t, err)

		// Check output contains trend info
		assert.Contains(t, output, now.Format("2006-01-02"))
		assert.Contains(t, output, "3")     // Solved today
		assert.Contains(t, output, "11:30") // Avg time today
		assert.Contains(t, output, yesterday.Format("2006-01-02"))
		assert.Contains(t, output, "2")    // Solved yesterday
		assert.Contains(t, output, "5")    // Solved this week
		assert.Contains(t, output, "71.4") // Success rate this week
	})

	t.Run("Reset", func(t *testing.T) {
		// We can't easily test the interactive part of reset that requires user input
		// but we can test the command execution

		// Mock Reset to return no error
		restore := mockReset(nil)
		defer restore()

		// Execute reset command with automatic yes response
		// Note: In a real test, we would mock the prompt/response
		output, err := executeCommand(rootCmd, "stats", "reset")
		assert.NoError(t, err)

		// Since we can't easily mock the prompt response, we can just check
		// that the command ran without error
		assert.NotContains(t, output, "Error resetting stats")
	})

	t.Run("GetSummaryError", func(t *testing.T) {
		// Mock GetSummary to return an error
		restore := mockGetSummary(nil, errors.New("summary error"))
		defer restore()

		// Execute stats command
		output, err := executeCommand(rootCmd, "stats")
		assert.NoError(t, err)                               // Command itself should not error
		assert.Contains(t, output, "Error retrieving stats") // But output should contain error message
	})

	t.Run("GetByPatternError", func(t *testing.T) {
		// Mock GetByPattern to return an error
		restore := mockGetByPattern(nil, errors.New("pattern stats error"))
		defer restore()

		// Execute patterns command
		output, err := executeCommand(rootCmd, "stats", "patterns")
		assert.NoError(t, err)                                       // Command itself should not error
		assert.Contains(t, output, "Error retrieving pattern stats") // But output should contain error message
	})

	t.Run("GetTrendsError", func(t *testing.T) {
		// Mock GetTrends to return an error
		restore := mockGetTrends(nil, errors.New("trends error"))
		defer restore()

		// Execute trends command
		output, err := executeCommand(rootCmd, "stats", "trends")
		assert.NoError(t, err)                                     // Command itself should not error
		assert.Contains(t, output, "Error retrieving trend stats") // But output should contain error message
	})
}
