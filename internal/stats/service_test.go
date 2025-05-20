package stats

import (
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

func TestStatsService(t *testing.T) {
	// Create a mock storage for testing
	mockStorage := NewMockStorage()
	
	// Create a stats service with the mock storage
	service := NewService().WithStorage(mockStorage)
	
	// Verify the service implements the StatsService interface
	var _ interfaces.StatsService = service
	
	// Test GetSummary with no sessions
	t.Run("GetSummary_Empty", func(t *testing.T) {
		summary, err := service.GetSummary()
		assert.NoError(t, err)
		assert.Equal(t, 0, summary.TotalAttempted)
		assert.Equal(t, 0, summary.TotalSolved)
		assert.Equal(t, 0.0, summary.SuccessRate)
	})
	
	// Add some test sessions
	now := time.Now()
	
	// Session 1: Successfully solved
	session1 := SessionStats{
		ProblemID:    "problem1",
		StartTime:    now.Add(-30 * time.Minute),
		EndTime:      now.Add(-15 * time.Minute),
		Duration:     15 * time.Minute,
		Solved:       true,
		Mode:         "practice",
		HintsUsed:    false,
		SolutionUsed: false,
		Patterns:     []string{"two-pointers", "sliding-window"},
		Difficulty:   "medium",
	}
	mockStorage.AddSession(session1)
	
	// Session 2: Failed attempt
	session2 := SessionStats{
		ProblemID:    "problem2",
		StartTime:    now.Add(-60 * time.Minute),
		EndTime:      now.Add(-45 * time.Minute),
		Duration:     15 * time.Minute,
		Solved:       false,
		Mode:         "practice",
		HintsUsed:    true,
		SolutionUsed: false,
		Patterns:     []string{"dynamic-programming"},
		Difficulty:   "hard",
	}
	mockStorage.AddSession(session2)
	
	// Session 3: Multiple attempts on same problem
	session3 := SessionStats{
		ProblemID:    "problem2",
		StartTime:    now.Add(-40 * time.Minute),
		EndTime:      now.Add(-30 * time.Minute),
		Duration:     10 * time.Minute,
		Solved:       true,
		Mode:         "practice",
		HintsUsed:    true,
		SolutionUsed: false,
		Patterns:     []string{"dynamic-programming"},
		Difficulty:   "hard",
	}
	mockStorage.AddSession(session3)
	
	// Test GetSummary with sessions
	t.Run("GetSummary", func(t *testing.T) {
		summary, err := service.GetSummary()
		assert.NoError(t, err)
		assert.Equal(t, 3, summary.TotalAttempted)
		assert.Equal(t, 2, summary.TotalSolved)
		assert.InDelta(t, 66.67, summary.SuccessRate, 0.01)
		assert.Equal(t, "00:12:30", summary.AvgSolveTime) // Average of 15 min and 10 min
		assert.Equal(t, "problem2", summary.FastestSolve.ProblemID)
		assert.Equal(t, "00:10:00", summary.FastestSolve.Time)
	})
	
	// Test GetByPattern
	t.Run("GetByPattern", func(t *testing.T) {
		patternStats, err := service.GetByPattern()
		assert.NoError(t, err)
		
		// Check two-pointers pattern
		twoPointers, ok := patternStats["two-pointers"]
		assert.True(t, ok)
		assert.Equal(t, 1, twoPointers.Attempted)
		assert.Equal(t, 1, twoPointers.Solved)
		assert.InDelta(t, 100.0, twoPointers.SuccessRate, 0.01)
		assert.Equal(t, "00:15:00", twoPointers.AvgTime)
		
		// Check dynamic-programming pattern
		dp, ok := patternStats["dynamic-programming"]
		assert.True(t, ok)
		assert.Equal(t, 2, dp.Attempted)
		assert.Equal(t, 1, dp.Solved)
		assert.InDelta(t, 50.0, dp.SuccessRate, 0.01)
		assert.Equal(t, "00:10:00", dp.AvgTime)
	})
	
	// Test GetTrends
	t.Run("GetTrends", func(t *testing.T) {
		trends, err := service.GetTrends()
		assert.NoError(t, err)
		
		// Check that we have daily trends
		assert.Equal(t, 7, len(trends.Daily))
		
		// Today should have some solves
		today := time.Now().Format("2006-01-02")
		var todaySolves int
		for _, daily := range trends.Daily {
			if daily.Date == today {
				todaySolves = daily.Solved
				break
			}
		}
		assert.Equal(t, 2, todaySolves)
		
		// Check that we have weekly trends
		assert.GreaterOrEqual(t, len(trends.Weekly), 1)
	})
	
	// Test RecordSession
	t.Run("RecordSession", func(t *testing.T) {
		// Add a new session
		newSession := SessionStats{
			ProblemID:    "problem3",
			StartTime:    now,
			EndTime:      now.Add(20 * time.Minute),
			Duration:     20 * time.Minute,
			Solved:       true,
			Mode:         "practice",
			HintsUsed:    false,
			SolutionUsed: false,
			Patterns:     []string{"hash-map"},
			Difficulty:   "easy",
		}
		
		err := service.RecordSession(newSession)
		assert.NoError(t, err)
		
		// Verify the session was added
		sessions, err := service.GetAllSessions()
		assert.NoError(t, err)
		assert.Equal(t, 4, len(sessions))
		
		// Check summary again
		summary, err := service.GetSummary()
		assert.NoError(t, err)
		assert.Equal(t, 4, summary.TotalAttempted)
		assert.Equal(t, 3, summary.TotalSolved)
	})
	
	// Test Reset
	t.Run("Reset", func(t *testing.T) {
		err := service.Reset()
		assert.NoError(t, err)
		
		// Verify all sessions are gone
		sessions, err := service.GetAllSessions()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(sessions))
		
		// Summary should be empty
		summary, err := service.GetSummary()
		assert.NoError(t, err)
		assert.Equal(t, 0, summary.TotalAttempted)
	})
}