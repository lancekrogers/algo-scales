// Tests for start command

package cmd

import (
	"errors"
	"testing"

	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/stretchr/testify/assert"
)

// Mock session.Start for testing
func mockSessionStart(err error) func() {
	original := session.Start
	session.Start = func(opts session.Options) error {
		return err
	}
	return func() {
		session.Start = original
	}
}

func TestStartCommand(t *testing.T) {
	t.Run("Learn", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute learn command
		output, err := executeCommand(rootCmd, "start", "learn")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("Practice", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute practice command
		output, err := executeCommand(rootCmd, "start", "practice")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("Cram", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute cram command
		output, err := executeCommand(rootCmd, "start", "cram")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("WithLanguageFlag", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute command with language flag
		output, err := executeCommand(rootCmd, "start", "learn", "--language", "python")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("WithTimerFlag", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute command with timer flag
		output, err := executeCommand(rootCmd, "start", "practice", "--timer", "30")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("WithPatternFlag", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute command with pattern flag
		output, err := executeCommand(rootCmd, "start", "learn", "--pattern", "sliding-window")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("WithDifficultyFlag", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute command with difficulty flag
		output, err := executeCommand(rootCmd, "start", "practice", "--difficulty", "hard")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("WithProblemArg", func(t *testing.T) {
		// Mock session.Start to return no error
		restore := mockSessionStart(nil)
		defer restore()

		// Execute command with problem argument
		output, err := executeCommand(rootCmd, "start", "learn", "two-sum")
		assert.NoError(t, err)
		assert.NotContains(t, output, "Error starting session")
	})

	t.Run("SessionError", func(t *testing.T) {
		// Mock session.Start to return an error
		restore := mockSessionStart(errors.New("session error"))
		defer restore()

		// Execute command
		output, err := executeCommand(rootCmd, "start", "learn")
		assert.NoError(t, err)                               // Command itself should not error
		assert.Contains(t, output, "Error starting session") // But output should contain error message
	})
}
