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
func (r *GoTestRunner) ExecuteTests(prob *problem.Problem, code string, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
	// Create a temporary directory for test execution
	testDir, err := os.MkdirTemp("", "algo-scales-go-test")
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
	mainFile := filepath.Join(testDir, "main.go")
	err = os.WriteFile(mainFile, []byte(testCode), 0644)
	if err != nil {
		return nil, false, fmt.Errorf("failed to write test file: %v", err)
	}
	
	// Build and run the test
	cmd := exec.Command("go", "run", mainFile)
	
	// Run the command with timeout
	stdout, stderr, err := runCommandWithTimeout(cmd, timeout)
	
	// Parse the results from stdout
	output := stdout.String()
	results := parseTestOutput(output, prob.TestCases)
	
	// If there were compile errors, include them in the results
	if err != nil && len(stderr.String()) > 0 {
		results = addErrorToResults(results, stderr.String())
	}
	
	return results, allTestsPassed(results), nil
}

// GenerateTestCode creates Go test code for a given problem
func (r *GoTestRunner) GenerateTestCode(prob *problem.Problem, solutionCode string) (string, error) {
	// Create the test file content template
	testTemplate := `package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// User's solution
%s

func main() {
	// Run tests
	allPassed := true
	
	%s
	
	if !allPassed {
		os.Exit(1)
	}
}
`
	
	// Generate test code for each test case
	var testCases strings.Builder
	for i, tc := range prob.TestCases {
		testCases.WriteString(fmt.Sprintf("\n\t// Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("\tfmt.Printf(\"Test %d\\n\")\n", i+1))
		testCases.WriteString(fmt.Sprintf("\t{\n\t\tinputStr := `%s`\n", tc.Input))
		testCases.WriteString(fmt.Sprintf("\t\texpectedStr := `%s`\n", tc.Expected))
		
		// Parse input based on the problem
		testCases.WriteString("\t\t// Parse input\n")
		testCases.WriteString("\t\t// Note: This parsing logic needs to be customized for each problem type\n")
		testCases.WriteString("\t\t// For example, parsing \"[1,2,3], 5\" into a slice and an int\n")
		testCases.WriteString("\t\tparts := strings.Split(inputStr, \", \")\n")
		testCases.WriteString("\t\tif len(parts) < 1 {\n\t\t\tfmt.Printf(\"Error parsing input: %s\\n\", inputStr)\n\t\t\tallPassed = false\n\t\t\tcontinue\n\t\t}\n")
		
		// Actual test execution - simplified for demonstration
		// This would need to parse the actual parameters and call the function
		testCases.WriteString("\t\t// Execute solution with the input\n")
		testCases.WriteString("\t\t// Note: This would need to call the actual solution function with parsed parameters\n")
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