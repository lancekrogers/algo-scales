// Tests for stats module

package stats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Override config dir for testing
func withTestDir(t *testing.T) (string, func()) {
	// Create a temporary test directory
	tempDir, err := ioutil.TempDir("", "algo-scales-test")
	require.NoError(t, err)

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	getConfigDir = func() string {
		return tempDir
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
		getConfigDir = origGetConfigDir
	}
}

// Create sample session stats
func createSampleSessions(t *testing.T, dir string, count int) []SessionStats {
	statsDir := filepath.Join(dir, "stats")
	err := os.MkdirAll(statsDir, 0755)
	require.NoError(t, err)

	sessions := make([]SessionStats, count)
	for i := 0; i < count; i++ {
		// Alternate solved/unsolved
		solved := i%2 == 0

		// Create session stats
		startTime := time.Now().Add(-time.Duration(i) * time.Hour)
		endTime := startTime.Add(30 * time.Minute)

		sessions[i] = SessionStats{
			ProblemID:    "problem" + string(rune('1'+i%3)),
			StartTime:    startTime,
			EndTime:      endTime,
			Duration:     endTime.Sub(startTime),
			Solved:       solved,
			Mode:         "practice",
			HintsUsed:    !solved,
			SolutionUsed: !solved,
			Patterns:     []string{"hash-map", "two-pointers"},
			Difficulty:   "Medium",
		}

		// Save to file
		filename := filepath.Join(statsDir, "session_"+sessions[i].ProblemID+"_"+startTime.Format("20060102_150405")+".json")
		data, err := json.MarshalIndent(sessions[i], "", "  ")
		require.NoError(t, err)
		err = ioutil.WriteFile(filename, data, 0644)
		require.NoError(t, err)
	}

	return sessions
}

func TestRecordSession(t *testing.T) {
	tempDir, cleanup := withTestDir(t)
	defer cleanup()

	// Create test session stats
	stats := SessionStats{
		ProblemID:    "test-problem",
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(20 * time.Minute),
		Duration:     20 * time.Minute,
		Solved:       true,
		Mode:         "learn",
		HintsUsed:    true,
		SolutionUsed: false,
		Patterns:     []string{"hash-map"},
		Difficulty:   "Easy",
	}

	// Record the session
	err := RecordSession(stats)
	require.NoError(t, err)

	// Check that the stats were saved
	statsDir := filepath.Join(tempDir, "stats")
	files, err := ioutil.ReadDir(statsDir)
	require.NoError(t, err)
	assert.Equal(t, 1, len(files))

	// Verify file content
	data, err := ioutil.ReadFile(filepath.Join(statsDir, files[0].Name()))
	require.NoError(t, err)

	var savedStats SessionStats
	err = json.Unmarshal(data, &savedStats)
	require.NoError(t, err)

	assert.Equal(t, stats.ProblemID, savedStats.ProblemID)
	assert.Equal(t, stats.Solved, savedStats.Solved)
	assert.Equal(t, stats.Duration, savedStats.Duration)
}

func TestGetSummary(t *testing.T) {
	tempDir, cleanup := withTestDir(t)
	defer cleanup()

	// Create sample sessions (6 sessions, 3 solved)
	sessions := createSampleSessions(t, tempDir, 6)

	// Get summary
	summary, err := GetSummary()
	require.NoError(t, err)

	// Verify summary
	assert.Equal(t, 6, summary.TotalAttempted)
	assert.Equal(t, 3, summary.TotalSolved)
	assert.NotEmpty(t, summary.AvgSolveTime)
	assert.Equal(t, float64(50), summary.SuccessRate) // 3/6 = 50%

	assert.NotEmpty(t, summary.FastestSolve.ProblemID)
	assert.NotEmpty(t, summary.FastestSolve.Time)

	// Check for most challenging problem
	// In our sample data, we'd expect a problem that was attempted and never solved
	assert.NotEmpty(t, summary.MostChallenging.ProblemID)
	assert.Greater(t, summary.MostChallenging.Attempts, 0)
}

func TestGetByPattern(t *testing.T) {
	tempDir, cleanup := withTestDir(t)
	defer cleanup()

	// Create sample sessions
	sessions := createSampleSessions(t, tempDir, 6)

	// Get pattern stats
	patternStats, err := GetByPattern()
	require.NoError(t, err)

	// Verify pattern stats
	assert.Len(t, patternStats, 2) // We used 2 patterns in our sample data

	// Check "hash-map" pattern
	hashMapStats, ok := patternStats["hash-map"]
	assert.True(t, ok)
	assert.Equal(t, "hash-map", hashMapStats.Pattern)
	assert.Equal(t, 6, hashMapStats.Attempted) // All sessions used hash-map
	assert.Equal(t, 3, hashMapStats.Solved)    // Half are solved
	assert.Equal(t, float64(50), hashMapStats.SuccessRate)
	assert.NotEmpty(t, hashMapStats.AvgTime)

	// Check "two-pointers" pattern
	twoPointersStats, ok := patternStats["two-pointers"]
	assert.True(t, ok)
	assert.Equal(t, "two-pointers", twoPointersStats.Pattern)
	assert.Equal(t, 6, twoPointersStats.Attempted)
	assert.Equal(t, 3, twoPointersStats.Solved)
	assert.Equal(t, float64(50), twoPointersStats.SuccessRate)
	assert.NotEmpty(t, twoPointersStats.AvgTime)
}

func TestGetTrends(t *testing.T) {
	tempDir, cleanup := withTestDir(t)
	defer cleanup()

	// Create sample sessions spanning multiple days
	createSampleSessions(t, tempDir, 10)

	// Get trends
	trends, err := GetTrends()
	require.NoError(t, err)

	// Verify daily trends
	assert.NotEmpty(t, trends.Daily)
	assert.LessOrEqual(t, len(trends.Daily), 7) // Should have up to 7 days

	// Verify weekly trends
	assert.NotEmpty(t, trends.Weekly)
	assert.LessOrEqual(t, len(trends.Weekly), 4) // Should have up to 4 weeks

	// Check trend data
	for _, day := range trends.Daily {
		assert.NotEmpty(t, day.Date)
		// We can't reliably check exact counts here due to the time-based nature
	}

	for _, week := range trends.Weekly {
		assert.NotEmpty(t, week.StartDate)
		assert.NotEmpty(t, week.EndDate)
		// Similarly, can't check exact counts
	}
}

func TestReset(t *testing.T) {
	tempDir, cleanup := withTestDir(t)
	defer cleanup()

	// Create sample sessions
	createSampleSessions(t, tempDir, 5)

	// Verify sessions were created
	statsDir := filepath.Join(tempDir, "stats")
	files, err := ioutil.ReadDir(statsDir)
	require.NoError(t, err)
	assert.Equal(t, 5, len(files))

	// Reset stats
	err = Reset()
	require.NoError(t, err)

	// Verify all files were removed
	files, err = ioutil.ReadDir(statsDir)
	require.NoError(t, err)
	assert.Equal(t, 0, len(files))
}

func TestLoadAllSessions(t *testing.T) {
	tempDir, cleanup := withTestDir(t)
	defer cleanup()

	// Create sample sessions
	expectedSessions := createSampleSessions(t, tempDir, 3)

	// Load sessions
	sessions, err := loadAllSessions()
	require.NoError(t, err)

	// Verify sessions were loaded
	assert.Equal(t, 3, len(sessions))

	// Check if all problem IDs are present
	problemIDs := make(map[string]bool)
	for _, session := range sessions {
		problemIDs[session.ProblemID] = true
	}

	for _, expected := range expectedSessions {
		assert.True(t, problemIDs[expected.ProblemID])
	}
}

func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "00:30"},
		{2 * time.Minute, "02:00"},
		{2*time.Minute + 30*time.Second, "02:30"},
		{1*time.Hour + 15*time.Minute + 45*time.Second, "75:45"},
	}

	for _, tc := range testCases {
		t.Run("format", func(t *testing.T) {
			result := formatDuration(tc.duration)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTimeHelpers(t *testing.T) {
	now := time.Now()

	t.Run("GetYearWeek", func(t *testing.T) {
		yearWeek := getYearWeek(now)
		year, week := now.ISOWeek()
		expected := fmt.Sprintf("%d-W%02d", year, week)
		assert.Equal(t, expected, yearWeek)
	})

	t.Run("StartOfWeek", func(t *testing.T) {
		start := startOfWeek(now)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekday-- // Adjust to make Monday the first day

		// Start of week should be the correct number of days before now
		expected := time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())
		assert.Equal(t, expected, start)
	})

	t.Run("EndOfWeek", func(t *testing.T) {
		end := endOfWeek(now)
		start := startOfWeek(now)
		expected := start.AddDate(0, 0, 6)
		assert.Equal(t, expected, end)
	})
}
