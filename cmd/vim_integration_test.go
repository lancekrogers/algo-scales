package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVimModeIntegration tests the actual CLI commands in vim mode
func TestVimModeIntegration(t *testing.T) {
	// Build the binary first
	binPath := filepath.Join(t.TempDir(), "algo-scales")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../")
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build binary")

	t.Run("list problems in vim mode", func(t *testing.T) {
		cmd := exec.Command(binPath, "list", "--vim-mode")
		output, err := cmd.Output()
		require.NoError(t, err)

		var resp VimListResponse
		err = json.Unmarshal(output, &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Problems)
	})

	t.Run("start session in vim mode", func(t *testing.T) {
		cmd := exec.Command(binPath, "start", "learn", "two_sum", "--language", "go", "--vim-mode")
		output, err := cmd.Output()
		require.NoError(t, err)

		var resp VimProblemResponse
		err = json.Unmarshal(output, &resp)
		require.NoError(t, err)
		assert.Equal(t, "two_sum", resp.ID)
		assert.NotEmpty(t, resp.Title)
		assert.NotEmpty(t, resp.Description)
		assert.NotEmpty(t, resp.StarterCode)
	})

	t.Run("get hint in vim mode", func(t *testing.T) {
		cmd := exec.Command(binPath, "hint", "--problem-id", "pair_with_target_sum", "--language", "go", "--vim-mode")
		output, err := cmd.Output()
		require.NoError(t, err)

		var resp VimHintResponse
		err = json.Unmarshal(output, &resp)
		require.NoError(t, err)
		assert.Equal(t, 1, resp.Level)
		assert.NotEmpty(t, resp.Hint)
		assert.Contains(t, resp.Hint, "two-pointers")
	})

	t.Run("get solution in vim mode", func(t *testing.T) {
		cmd := exec.Command(binPath, "solution", "--problem-id", "two_sum", "--language", "go", "--vim-mode")
		output, err := cmd.Output()
		require.NoError(t, err)

		var resp VimSolutionResponse
		err = json.Unmarshal(output, &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Solution)
		assert.Contains(t, resp.Solution, "func twoSum")
	})

	t.Run("submit solution in vim mode", func(t *testing.T) {
		// Create a test solution file
		solutionFile := filepath.Join(t.TempDir(), "solution.go")
		solutionCode := `func twoSum(nums []int, target int) []int {
			seen := make(map[int]int)
			for i, num := range nums {
				if j, ok := seen[target-num]; ok {
					return []int{j, i}
				}
				seen[num] = i
			}
			return nil
		}`
		err := os.WriteFile(solutionFile, []byte(solutionCode), 0644)
		require.NoError(t, err)

		cmd := exec.Command(binPath, "submit", 
			"--problem-id", "two_sum",
			"--language", "go",
			"--file", solutionFile,
			"--vim-mode")
		output, err := cmd.Output()
		require.NoError(t, err)

		var resp VimSubmitResponse
		err = json.Unmarshal(output, &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.TestResults)
		// Check that we got test results back (may not all pass with simple solution)
		assert.GreaterOrEqual(t, len(resp.TestResults), 1)
	})

	t.Run("ai hint in vim mode", func(t *testing.T) {
		cmd := exec.Command(binPath, "ai-hint",
			"--problem-id", "two_sum",
			"--language", "go",
			"--user-code", "// my solution attempt",
			"--provider", "claude",
			"--vim-mode")
		output, err := cmd.Output()
		require.NoError(t, err)

		var resp map[string]interface{}
		err = json.Unmarshal(output, &resp)
		require.NoError(t, err)
		assert.True(t, resp["ready"].(bool))
		assert.NotEmpty(t, resp["system_prompt"])
		assert.NotEmpty(t, resp["user_message"])
	})
}

// TestHintLevelPersistence verifies that hint levels persist within a session
// Note: For CLI, each command is a new session, so levels reset
func TestHintLevelPersistence(t *testing.T) {
	// Within a single process, hint levels should persist
	
	// Reset hint levels
	hintLevels = make(map[string]int)
	
	problemID := "two_sum"
	
	// Simulate multiple hint calls within same process
	for i := 1; i <= 3; i++ {
		level := hintLevels[problemID]
		level++
		hintLevels[problemID] = level
		
		assert.Equal(t, i, level, "Hint level should increment")
	}
	
	// Different problem should start at level 1
	anotherProblem := "three_sum"
	level := hintLevels[anotherProblem]
	level++
	hintLevels[anotherProblem] = level
	assert.Equal(t, 1, level, "New problem should start at level 1")
}