// Integration tests for CLI commands

package cmd

import (
	"os"
	"testing"

	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLICommandsIntegration tests that CLI commands can be executed without panics
func TestCLICommandsIntegration(t *testing.T) {
	// Set testing environment to skip UI launches
	os.Setenv("TESTING", "1")
	defer os.Unsetenv("TESTING")

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "start cram",
			args: []string{"start", "cram"},
		},
		{
			name: "start learn",
			args: []string{"start", "learn"},
		},
		{
			name: "start practice",
			args: []string{"start", "practice"},
		},
		{
			name: "solve command",
			args: []string{"solve"},
		},
		{
			name: "list problems",
			args: []string{"list"},
		},
		{
			name: "stats summary",
			args: []string{"stats"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh root command for each test
			cmd := &cobra.Command{
				Use: "algo-scales",
			}
			
			// Add all subcommands
			cmd.AddCommand(startCmd)
			cmd.AddCommand(cliCmd)
			cmd.AddCommand(listCmd)
			cmd.AddCommand(statsCmd)
			
			// Mock session.Start to prevent actual session creation during tests
			require.NotPanics(t, func() {
				output, err := executeCommand(cmd, tt.args...)
				// We expect some errors since we're not setting up full mock environment
				// but we should not panic
				t.Logf("Command: %v, Output: %s, Error: %v", tt.args, output, err)
			}, "Command should not panic: %v", tt.args)
		})
	}
}

// TestStartCommandsWithMockSession tests start commands with mocked session
func TestStartCommandsWithMockSession(t *testing.T) {
	// This test ensures that start commands call session.Start and then launchUI
	// without actually creating sessions or launching UI
	
	// Mock session.Start to track calls
	originalStart := session.Start
	var startCalled bool
	var startOpts session.Options
	
	session.Start = func(opts session.Options) error {
		startCalled = true
		startOpts = opts
		return nil // Return success
	}
	defer func() {
		session.Start = originalStart
	}()
	
	tests := []struct {
		name         string
		args         []string
		expectedMode session.Mode
	}{
		{
			name:         "cram mode",
			args:         []string{"start", "cram"},
			expectedMode: session.CramMode,
		},
		{
			name:         "learn mode",
			args:         []string{"start", "learn"},
			expectedMode: session.LearnMode,
		},
		{
			name:         "practice mode",
			args:         []string{"start", "practice"},
			expectedMode: session.PracticeMode,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startCalled = false
			startOpts = session.Options{}
			
			output, err := executeCommand(rootCmd, tt.args...)
			
			assert.NoError(t, err)
			assert.True(t, startCalled, "session.Start should have been called")
			assert.Equal(t, tt.expectedMode, startOpts.Mode, "Expected mode should match")
			assert.NotContains(t, output, "Error starting session")
		})
	}
}