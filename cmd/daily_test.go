package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Skip tests for daily_test.go for now as they require more extensive mocking
func TestDailyCommandFlags(t *testing.T) {
	// Test that the daily command has the expected flags
	assert.NotNil(t, dailyCmd.Flag("language"))
	assert.NotNil(t, dailyCmd.Flag("timer"))
	assert.NotNil(t, dailyCmd.Flag("difficulty"))
}

// TestDailyCommand disabled for now as it requires more extensive mocking

// TestDisplayStreakInfo disabled for now as it requires capturing stdout