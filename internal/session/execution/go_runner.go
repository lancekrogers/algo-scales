package execution

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/logging"
)

// GoTestRunner implements the TestRunner interface for Go code
type GoTestRunner struct {
	BaseTestRunner
}

// NewGoTestRunner creates a new Go test runner
func NewGoTestRunner() *GoTestRunner {
	return &GoTestRunner{
		BaseTestRunner: NewBaseTestRunner("go"),
	}
}

// ExecuteTests runs tests for a Go solution
func (r *GoTestRunner) ExecuteTests(ctx context.Context, prob *interfaces.Problem, code string, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
	// Create a context with timeout for the entire operation
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	// Add logging context
	ctx = logging.WithOperation(ctx, "ExecuteGoTests")
	ctx = logging.WithComponent(ctx, "GoTestRunner")
	logger := logging.TestRunnerLogger.WithContext(ctx)
	
	// Create session snapshot for error logging
	sessionState := &logging.SessionSnapshot{
		ProblemID:    prob.ID,
		Language:     "go",
		Mode:         "test_execution",
		UserCode:     code,
		StartTime:    time.Now(),
		Patterns:     prob.Tags,
		Difficulty:   prob.Difficulty,
		CustomFields: map[string]string{
			"timeout":      timeout.String(),
			"test_count":   fmt.Sprintf("%d", len(prob.TestCases)),
		},
	}
	
	// Log operation start
	finishLog := logger.StartOperation(fmt.Sprintf("Execute Go tests for problem %s", prob.ID))
	defer func() {
		if r := recover(); r != nil {
			if logging.GlobalErrorLogger != nil {
				logging.GlobalErrorLogger.LogPanic(ctx, r, "execute_go_tests", sessionState)
			}
			finishLog(fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()
	
	logger.Info("Creating temporary directory for Go test execution")
	// Create a temporary directory for test execution
	testDir, err := os.MkdirTemp("", "algo-scales-go-test")
	if err != nil {
		if logging.GlobalErrorLogger != nil {
			logging.GlobalErrorLogger.LogFileOperationError(ctx, err, "create_temp_dir", testDir, sessionState)
		}
		finishLog(err)
		return nil, false, fmt.Errorf("failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up when done
	
	logger.Info("Generating test code")
	// Generate test code
	testCode, err := r.GenerateTestCode(prob, code)
	if err != nil {
		if logging.GlobalErrorLogger != nil {
			logging.GlobalErrorLogger.LogTestExecutionError(ctx, err, "go", code, "", sessionState)
		}
		finishLog(err)
		return nil, false, fmt.Errorf("failed to generate test code: %v", err)
	}
	
	logger.Info("Writing test file to temporary directory")
	// Write the test file
	mainFile := filepath.Join(testDir, "main.go")
	err = os.WriteFile(mainFile, []byte(testCode), 0644)
	if err != nil {
		if logging.GlobalErrorLogger != nil {
			logging.GlobalErrorLogger.LogFileOperationError(ctx, err, "write_test_file", mainFile, sessionState)
		}
		finishLog(err)
		return nil, false, fmt.Errorf("failed to write test file: %v", err)
	}
	
	logger.Info("Executing Go test with timeout of %v", timeout)
	// Build and run the test
	cmd := exec.CommandContext(ctx, "go", "run", mainFile)
	
	// Update session state with test file info
	sessionState.CodeFile = mainFile
	sessionState.Workspace = testDir
	
	// Run the command with timeout
	stdout, stderr, err := runCommandWithTimeout(cmd, timeout)
	
	// Parse the results from stdout
	output := stdout.String()
	results := parseTestOutput(output, prob.TestCases)
	
	// If there were compile errors, include them in the results
	if err != nil && len(stderr.String()) > 0 {
		logger.Warn("Test execution failed with errors: %v", stderr.String())
		
		// Log detailed test execution error
		if logging.GlobalErrorLogger != nil {
			testError := fmt.Errorf("test execution failed: %v\nSTDOUT:\n%s\nSTDERR:\n%s", err, stdout.String(), stderr.String())
			logging.GlobalErrorLogger.LogTestExecutionError(ctx, testError, "go", code, "", sessionState)
		}
		
		results = addErrorToResults(results, stderr.String())
	}
	
	allPassed := allTestsPassed(results)
	logger.Info("Test execution completed: %d tests, %t all passed", len(results), allPassed)
	
	// Log success or failure
	if !allPassed && logging.GlobalErrorLogger != nil {
		failedTests := 0
		for _, result := range results {
			if !result.Passed {
				failedTests++
			}
		}
		testError := fmt.Errorf("test execution failed: %d of %d tests failed", failedTests, len(results))
		logging.GlobalErrorLogger.LogTestExecutionError(ctx, testError, "go", code, "", sessionState)
	}
	
	finishLog(nil)
	return results, allPassed, nil
}
// GenerateTestCode creates test code for a given problem
func (r *GoTestRunner) GenerateTestCode(prob *interfaces.Problem, solutionCode string) (string, error) {
	return r.generateTestTemplate(prob, solutionCode)
}

// generateTestTemplate generates the Go test template (extracted from ExecuteTests)
func (r *GoTestRunner) generateTestTemplate(prob *interfaces.Problem, solutionCode string) (string, error) {
	testTemplate := `package main

import (
	"fmt"
	"os"
	"strings"
)

// User's solution
%s

func main() {
	// Run tests
	allPassed := true
	
	%s
	
	if \!allPassed {
		os.Exit(1)
	}
}
`
	
	// Generate test code for each test case
	var testCases strings.Builder
	for i, tc := range prob.TestCases {
		// Use string fields directly
		inputStr := tc.Input
		expectedStr := tc.Expected
		
		testCases.WriteString(fmt.Sprintf("\n\t// Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("\tfmt.Printf(\"Test %d\\n\")\n", i+1))
		testCases.WriteString(fmt.Sprintf("\t{\n\t\tinputStr := `%s`\n", inputStr))
		testCases.WriteString(fmt.Sprintf("\t\texpectedStr := `%s`\n", expectedStr))
		
		// Parse input based on the problem
		testCases.WriteString("\t\t// Parse input\n")
		testCases.WriteString("\t\t// Note: This parsing logic needs to be customized for each problem type\n")
		testCases.WriteString("\t\tparts := strings.Split(inputStr, \", \")\n")
		testCases.WriteString("\t\tif len(parts) < 1 {\n\t\t\tfmt.Printf(\"Error parsing input: %s\\n\", inputStr)\n\t\t\tallPassed = false\n\t\t\tcontinue\n\t\t}\n")
		
		// Actual test execution - simplified for demonstration
		testCases.WriteString("\t\t// Execute solution with the input\n")
		testCases.WriteString("\t\tresult := \"PLACEHOLDER\" // Replace with actual function call\n")
		
		// Check result
		testCases.WriteString("\t\t// Check result\n")
		testCases.WriteString("\t\tif fmt.Sprintf(\"%v\", result) == expectedStr {\n")
		testCases.WriteString("\t\t\tfmt.Println(\"✅ PASSED\")\n")
		testCases.WriteString("\t\t} else {\n")
		testCases.WriteString("\t\t\tfmt.Printf(\"❌ FAILED\\nExpected: %s\\nGot: %v\\n\", expectedStr, result)\n")
		testCases.WriteString("\t\t\tallPassed = false\n")
		testCases.WriteString("\t\t}\n")
		testCases.WriteString("\t}\n")
	}
	
	// Complete the test code
	return fmt.Sprintf(testTemplate, solutionCode, testCases.String()), nil
}
