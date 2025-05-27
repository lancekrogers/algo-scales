package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVimCommands(t *testing.T) {
	// Create a temporary solution file for testing
	tmpDir := t.TempDir()
	solutionFile := filepath.Join(tmpDir, "solution.go")
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
	err := ioutil.WriteFile(solutionFile, []byte(solutionCode), 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		cmd      *cobra.Command
		args     []string
		flags    map[string]string
		wantJSON bool
		checkFn  func(t *testing.T, output string)
	}{
		{
			name: "submit command with vim mode",
			cmd:  submitCmd,
			args: []string{},
			flags: map[string]string{
				"problem-id": "two_sum",
				"language":   "go",
				"file":       solutionFile,
				"vim-mode":   "true",
			},
			wantJSON: true,
			checkFn: func(t *testing.T, output string) {
				var resp VimSubmitResponse
				err := json.Unmarshal([]byte(output), &resp)
				assert.NoError(t, err, "output should be valid JSON")
				assert.NotNil(t, resp.TestResults, "should have test results")
			},
		},
		{
			name: "hint command with vim mode",
			cmd:  hintCmd,
			args: []string{},
			flags: map[string]string{
				"problem-id": "two_sum",
				"vim-mode":   "true",
			},
			wantJSON: true,
			checkFn: func(t *testing.T, output string) {
				var resp VimHintResponse
				err := json.Unmarshal([]byte(output), &resp)
				assert.NoError(t, err, "output should be valid JSON")
				assert.NotEmpty(t, resp.Hint, "should have hint text")
			},
		},
		{
			name: "solution command with vim mode",
			cmd:  solutionCmd,
			args: []string{},
			flags: map[string]string{
				"problem-id": "two_sum",
				"language":   "go",
				"vim-mode":   "true",
			},
			wantJSON: true,
			checkFn: func(t *testing.T, output string) {
				var resp VimSolutionResponse
				err := json.Unmarshal([]byte(output), &resp)
				assert.NoError(t, err, "output should be valid JSON")
				assert.NotEmpty(t, resp.Solution, "should have solution code")
			},
		},
		{
			name: "submit command without vim mode",
			cmd:  submitCmd,
			args: []string{},
			flags: map[string]string{
				"problem-id": "two_sum",
				"language":   "go",
				"file":       solutionFile,
				"vim-mode":   "false",
			},
			wantJSON: false,
			checkFn: func(t *testing.T, output string) {
				assert.Contains(t, output, "This command is for vim mode only")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset command flags
			tt.cmd.ResetFlags()
			
			// Re-initialize flags (since ResetFlags removes them)
			switch tt.cmd {
			case submitCmd, testCmd:
				tt.cmd.Flags().String("problem-id", "", "Problem ID")
				tt.cmd.Flags().String("language", "go", "Programming language")
				tt.cmd.Flags().String("file", "", "Solution file path")
				tt.cmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
			case hintCmd:
				tt.cmd.Flags().String("problem-id", "", "Problem ID")
				tt.cmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
			case solutionCmd:
				tt.cmd.Flags().String("problem-id", "", "Problem ID")
				tt.cmd.Flags().String("language", "go", "Programming language")
				tt.cmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
			}

			// Set flags
			for flag, value := range tt.flags {
				err := tt.cmd.Flags().Set(flag, value)
				require.NoError(t, err)
			}

			// Capture output
			buf := new(bytes.Buffer)
			tt.cmd.SetOut(buf)
			tt.cmd.SetErr(buf)

			// Redirect stdout to capture fmt.Println output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute command
			tt.cmd.Run(tt.cmd, tt.args)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout
			output, _ := ioutil.ReadAll(r)

			// Check output
			outputStr := string(output)
			if tt.wantJSON {
				// For JSON output, verify it's valid JSON
				var js json.RawMessage
				err := json.Unmarshal(output, &js)
				if err != nil {
					t.Logf("Output: %s", outputStr)
				}
			}
			
			if tt.checkFn != nil {
				tt.checkFn(t, outputStr)
			}
		})
	}
}

func TestVimModeList(t *testing.T) {
	// Test the vim mode list functionality
	// Mock the problem.ListAll function to return test data
	oldListAll := problem.ListAll
	problem.ListAll = func() ([]problem.Problem, error) {
		return []problem.Problem{
			{
				ID:         "test_problem",
				Title:      "Test Problem",
				Difficulty: "easy",
			},
		}, nil
	}
	
	// Restore original function after test
	defer func() {
		problem.ListAll = oldListAll
	}()

	// Create a new root command instance to avoid state pollution
	testRootCmd := &cobra.Command{
		Use:   "algo-scales",
		Short: "Test root command",
	}
	testRootCmd.PersistentFlags().Bool("vim-mode", false, "Enable vim mode output")
	
	// Create a test list command
	testListCmd := &cobra.Command{
		Use:   "list",
		Short: "List problems",
		Run:   listCmd.Run, // Use the actual wrapped run function
	}
	testRootCmd.AddCommand(testListCmd)
	
	// Set vim mode
	err := testRootCmd.PersistentFlags().Set("vim-mode", "true")
	require.NoError(t, err)
	
	// Capture output using a buffer and custom writer
	output := &bytes.Buffer{}
	testListCmd.SetOut(output)
	
	// Save and redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Execute command
	testListCmd.Run(testListCmd, []string{})
	
	// Restore stdout and read output
	w.Close()
	os.Stdout = oldStdout
	capturedOutput, _ := ioutil.ReadAll(r)
	
	// Use captured output if buffer is empty
	if output.Len() == 0 && len(capturedOutput) > 0 {
		output.Write(capturedOutput)
	}

	// Verify JSON output
	outputStr := output.String()
	t.Logf("Output: %s", outputStr)
	
	var resp VimListResponse
	err = json.Unmarshal([]byte(outputStr), &resp)
	require.NoError(t, err, "Failed to unmarshal JSON output")
	require.Len(t, resp.Problems, 1, "Expected 1 problem in response")
	assert.Equal(t, "test_problem", resp.Problems[0].ID)
}

func TestVimModeStart(t *testing.T) {
	// Test the vim mode start functionality
	tests := []struct {
		name    string
		mode    string
		wantErr bool
	}{
		{
			name:    "learn mode",
			mode:    "learn",
			wantErr: false,
		},
		{
			name:    "practice mode", 
			mode:    "practice",
			wantErr: false,
		},
		{
			name:    "cram mode",
			mode:    "cram",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would test handleVimModeSession but we need to mock the services
			// For now, we just verify the commands are properly initialized
			var cmd *cobra.Command
			switch tt.mode {
			case "learn":
				cmd = learnCmd
			case "practice":
				cmd = practiceCmd
			case "cram":
				cmd = cramCmd
			}
			
			assert.NotNil(t, cmd)
			assert.NotNil(t, cmd.Run)
		})
	}
}

func TestGetPatternHint(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		want     string
	}{
		{
			name:     "empty patterns",
			patterns: []string{},
			want:     "general problem-solving techniques",
		},
		{
			name:     "known pattern",
			patterns: []string{"two-pointers"},
			want:     "Balanced and efficient, the workhorse of array manipulation",
		},
		{
			name:     "unknown pattern",
			patterns: []string{"unknown-pattern"},
			want:     "unknown-pattern pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPatternHint(tt.patterns)
			assert.Equal(t, tt.want, got)
		})
	}
}