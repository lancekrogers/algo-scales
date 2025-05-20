package execution

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// JavaScriptTestRunner implements the TestRunner interface for JavaScript code
type JavaScriptTestRunner struct {
	BaseTestRunner
}

// NewJavaScriptTestRunner creates a new JavaScript test runner
func NewJavaScriptTestRunner() *JavaScriptTestRunner {
	return &JavaScriptTestRunner{
		BaseTestRunner: NewBaseTestRunner("javascript"),
	}
}

// ExecuteTests runs tests for a JavaScript solution
func (r *JavaScriptTestRunner) ExecuteTests(prob *problem.Problem, code string, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
	// Create a temporary directory for test execution
	testDir, err := os.MkdirTemp("", "algo-scales-js-test")
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
	testFile := filepath.Join(testDir, "test_solution.js")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		return nil, false, fmt.Errorf("failed to write test file: %v", err)
	}
	
	// Run the test
	cmd := exec.Command("node", testFile)
	
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

// GenerateTestCode creates JavaScript test code for a given problem
func (r *JavaScriptTestRunner) GenerateTestCode(prob *problem.Problem, solutionCode string) (string, error) {
	// Create the test file content template
	testTemplate := `
// User's solution
%s

// Test cases
function runTests() {
    let allPassed = true;
    
    %s
    
    return allPassed;
}

// Run tests
const success = runTests();
if (!success) {
    process.exit(1);
}
`
	
	// Generate test code for each test case
	var testCases strings.Builder
	for i, tc := range prob.TestCases {
		testCases.WriteString(fmt.Sprintf("\n    // Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("    console.log(\"Test %d\");\n", i+1))
		testCases.WriteString(fmt.Sprintf("    const inputStr = '%s';\n", tc.Input))
		testCases.WriteString(fmt.Sprintf("    const expectedStr = '%s';\n", tc.Expected))
		
		// Parse input (very simplified - would need to be customized)
		testCases.WriteString("    // Parse input (simplified)\n")
		testCases.WriteString("    // This would need to be customized based on the problem\n")
		testCases.WriteString("    try {\n")
		testCases.WriteString("        // Simplified parsing logic - would need to be customized\n")
		testCases.WriteString("        // For example, parsing \"[1,2,3], 5\" for a twoSum problem\n")
		testCases.WriteString("        // const result = twoSum(parsedArray, parsedTarget);\n")
		testCases.WriteString("        const result = \"PLACEHOLDER\";\n")
		
		// Check result
		testCases.WriteString("        // Check result\n")
		testCases.WriteString("        if (String(result) === expectedStr) {\n")
		testCases.WriteString("            console.log(\"✅ PASSED\");\n")
		testCases.WriteString("        } else {\n")
		testCases.WriteString("            console.log(`❌ FAILED\\nExpected: ${expectedStr}\\nGot: ${result}`);\n")
		testCases.WriteString("            allPassed = false;\n")
		testCases.WriteString("        }\n")
		testCases.WriteString("    } catch (e) {\n")
		testCases.WriteString("        console.log(`❌ ERROR: ${e.message}`);\n")
		testCases.WriteString("        allPassed = false;\n")
		testCases.WriteString("    }\n")
	}
	
	// Complete the test code
	return fmt.Sprintf(testTemplate, solutionCode, testCases.String()), nil
}