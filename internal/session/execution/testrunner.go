// Package execution provides code execution functionality
package execution

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/utils"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// BaseTestRunner contains common functionality for test runners
type BaseTestRunner struct {
	language string
	fs       interfaces.FileSystem
}

// NewBaseTestRunner creates a new base test runner
func NewBaseTestRunner(language string) BaseTestRunner {
	return BaseTestRunner{
		language: language,
		fs:       utils.NewFileSystem(),
	}
}

// WithFileSystem sets a custom file system for the test runner
func (b *BaseTestRunner) WithFileSystem(fs interfaces.FileSystem) *BaseTestRunner {
	b.fs = fs
	return b
}

// GetLanguage returns the language this runner supports
func (b *BaseTestRunner) GetLanguage() string {
	return b.language
}

// parseTestOutput parses the test output to extract results
func parseTestOutput(output string, testCases []problem.TestCase) []interfaces.TestResult {
	// Create results array
	results := make([]interfaces.TestResult, len(testCases))
	
	// Initialize with basic information
	for i, tc := range testCases {
		results[i] = interfaces.TestResult{
			Input:    tc.Input,
			Expected: tc.Expected,
			Actual:   "No output captured",
			Passed:   false,
		}
	}
	
	// Try to parse the output to get actual results
	// This is a simple parser that looks for test numbers and PASSED/FAILED markers
	lines := strings.Split(output, "\n")
	currentTest := -1
	
	for _, line := range lines {
		// Check if this is a test header line
		if strings.HasPrefix(line, "Test ") {
			testNumStr := strings.TrimPrefix(line, "Test ")
			var testNum int
			_, err := fmt.Sscanf(testNumStr, "%d", &testNum)
			if err == nil && testNum > 0 && testNum <= len(results) {
				currentTest = testNum - 1
			}
			continue
		}
		
		// If we have a current test, look for PASSED/FAILED
		if currentTest >= 0 && currentTest < len(results) {
			if strings.Contains(line, "✅ PASSED") {
				results[currentTest].Passed = true
				results[currentTest].Actual = results[currentTest].Expected // Assume correct if passed
			} else if strings.Contains(line, "❌ FAILED") {
				results[currentTest].Passed = false
				// Try to extract the actual output
				if idx := strings.Index(line, "Got: "); idx >= 0 {
					results[currentTest].Actual = strings.TrimSpace(line[idx+5:])
				}
			} else if strings.HasPrefix(line, "Got: ") {
				results[currentTest].Actual = strings.TrimPrefix(line, "Got: ")
			}
		}
	}
	
	return results
}

// addErrorToResults adds error messages to failed test results
func addErrorToResults(results []interfaces.TestResult, errorMsg string) []interfaces.TestResult {
	// Add error message to all failed tests
	for i := range results {
		if !results[i].Passed {
			results[i].Actual = fmt.Sprintf("Error: %s", errorMsg)
		}
	}
	return results
}

// allTestsPassed checks if all tests passed
func allTestsPassed(results []interfaces.TestResult) bool {
	for _, r := range results {
		if !r.Passed {
			return false
		}
	}
	return true
}

// runCommandWithTimeout runs a command with a timeout
func runCommandWithTimeout(cmd *exec.Cmd, timeout time.Duration) (stdout, stderr bytes.Buffer, err error) {
	// Set up stdout and stderr
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Start the command
	err = cmd.Start()
	if err != nil {
		return stdout, stderr, err
	}
	
	// Use a channel to signal completion
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	
	// Wait for the command to complete or timeout
	select {
	case <-time.After(timeout):
		// Try to kill the process
		cmd.Process.Kill()
		return stdout, stderr, fmt.Errorf("command timed out after %v", timeout)
	case err = <-done:
		// Command completed
		return stdout, stderr, err
	}
}