package execution

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// PythonTestRunner implements the TestRunner interface for Python code
type PythonTestRunner struct {
	BaseTestRunner
}

// NewPythonTestRunner creates a new Python test runner
func NewPythonTestRunner() *PythonTestRunner {
	return &PythonTestRunner{
		BaseTestRunner: NewBaseTestRunner("python"),
	}
}

// ExecuteTests runs tests for a Python solution
func (r *PythonTestRunner) ExecuteTests(prob *interfaces.Problem, code string, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
	// Create a temporary directory for test execution
	testDir, err := os.MkdirTemp("", "algo-scales-python-test")
	if err != nil {
		return nil, false, fmt.Errorf("failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up when done
	
	// Generate test code
	testCode, err := r.GenerateTestCode(prob, code)
	if err != nil {
		return nil, false, fmt.Errorf("failed to generate test code: %v", err)
	}
	
	// Write the test file
	testFile := filepath.Join(testDir, "test_solution.py")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		return nil, false, fmt.Errorf("failed to write test file: %v", err)
	}
	
	// Run the test
	cmd := exec.Command("python", testFile)
	
	// Run the command with timeout
	stdout, stderr, err := runCommandWithTimeout(cmd, timeout)
	
	// Parse the results from stdout
	output := stdout.String()
	results := parseTestOutput(output, prob.TestCases)
	
	// If there were errors, include them in the results
	if err != nil && len(stderr.String()) > 0 {
		results = addErrorToResults(results, stderr.String())
	}
	
	return results, allTestsPassed(results), nil
}

// GenerateTestCode creates Python test code for a given problem
func (r *PythonTestRunner) GenerateTestCode(prob *interfaces.Problem, solutionCode string) (string, error) {
	// Create the test file content template
	testTemplate := `
# User's solution
%s

# Test cases
def main():
    all_passed = True
    
    %s
    
    return all_passed

if __name__ == "__main__":
    success = main()
    if not success:
        exit(1)
`
	
	// Generate test code for each test case
	var testCases strings.Builder
	for i, tc := range prob.TestCases {
		// Use string fields directly
		inputStr := tc.Input
		expectedStr := tc.Expected
		
		testCases.WriteString(fmt.Sprintf("\n    # Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("    print(\"Test %d\")\n", i+1))
		testCases.WriteString(fmt.Sprintf("    input_str = '%s'\n", inputStr))
		testCases.WriteString(fmt.Sprintf("    expected_str = '%s'\n", expectedStr))
		
		// Parse input (very simplified - would need to be customized)
		testCases.WriteString("    # Parse input (simplified)\n")
		testCases.WriteString("    # This would need to be customized based on the problem\n")
		testCases.WriteString("    try:\n")
		testCases.WriteString("        # Simplified parsing logic - would need to be customized\n")
		testCases.WriteString("        # For example, parsing \"[1,2,3], 5\" for a two_sum problem\n")
		testCases.WriteString("        # result = two_sum(parsed_array, parsed_target)\n")
		testCases.WriteString("        result = \"PLACEHOLDER\"\n")
		
		// Check result
		testCases.WriteString("        # Check result\n")
		testCases.WriteString("        if str(result) == expected_str:\n")
		testCases.WriteString("            print(\"✅ PASSED\")\n")
		testCases.WriteString("        else:\n")
		testCases.WriteString("            print(f\"❌ FAILED\\nExpected: {expected_str}\\nGot: {result}\")\n")
		testCases.WriteString("            all_passed = False\n")
		testCases.WriteString("    except Exception as e:\n")
		testCases.WriteString("        print(f\"❌ ERROR: {e}\")\n")
		testCases.WriteString("        all_passed = False\n")
	}
	
	// Complete the test code
	return fmt.Sprintf(testTemplate, solutionCode, testCases.String()), nil
}