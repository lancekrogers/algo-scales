package daily

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) (string, func()) {
	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "algoscales-test")
	require.NoError(t, err)

	// Set custom DBFileName for testing to avoid interfering with actual DB
	origDBPath := GetDBPath
	GetDBPath = func() string {
		return filepath.Join(tempDir, "test-daily.db")
	}

	// Return a cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
		GetDBPath = origDBPath
	}

	return tempDir, cleanup
}

func TestLoadDefaultProgress(t *testing.T) {
	// Setup test database with a new path that doesn't exist yet
	_, cleanup := setupTestDB(t)
	defer cleanup()

	// Load progress (should return default values)
	progress, err := LoadProgress()
	require.NoError(t, err)
	
	// Verify default values
	assert.Equal(t, 0, progress.Current)
	assert.Equal(t, 0, len(progress.Completed))
	assert.Equal(t, 0, progress.Streak)
	assert.Equal(t, 0, progress.LongestStreak)
	assert.True(t, progress.LastPracticed.IsZero())
}

func TestUpdateStreak(t *testing.T) {
	tests := []struct {
		name          string
		progress      ScaleProgress
		lastPracticed time.Time
		wantStreak    int
	}{
		{
			name: "first practice ever",
			progress: ScaleProgress{
				Streak:        0,
				LongestStreak: 0,
				LastPracticed: time.Time{}, // Zero time
			},
			wantStreak: 1,
		},
		{
			name: "practiced yesterday",
			progress: ScaleProgress{
				Streak:        2,
				LongestStreak: 5,
				LastPracticed: time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour),
			},
			wantStreak: 3,
		},
		{
			name: "practiced today already",
			progress: ScaleProgress{
				Streak:        3,
				LongestStreak: 5,
				LastPracticed: time.Now().Truncate(24 * time.Hour),
			},
			wantStreak: 3, // No change
		},
		{
			name: "break in streak",
			progress: ScaleProgress{
				Streak:        5,
				LongestStreak: 7,
				LastPracticed: time.Now().Add(-48 * time.Hour).Truncate(24 * time.Hour), // 2 days ago
			},
			wantStreak: 1, // Reset to 1
		},
		{
			name: "new longest streak",
			progress: ScaleProgress{
				Streak:        7,
				LongestStreak: 7,
				LastPracticed: time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour),
			},
			wantStreak: 8, // New longest
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the update streak function
			UpdateStreak(&tt.progress)

			// Check streak was updated correctly
			assert.Equal(t, tt.wantStreak, tt.progress.Streak)
			
			// Check longest streak logic
			if tt.name == "new longest streak" {
				assert.Equal(t, 8, tt.progress.LongestStreak)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "test",
			expected: false,
		},
		{
			name:     "item present",
			slice:    []string{"one", "two", "three"},
			item:     "two",
			expected: true,
		},
		{
			name:     "item not present",
			slice:    []string{"one", "two", "three"},
			item:     "four",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}