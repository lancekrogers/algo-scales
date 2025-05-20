// Code execution logic for testing solutions
package session

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// ExecuteTests runs tests for the current solution
// This replaces the simulation in RunTests
func ExecuteTests(s interfaces.Session, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
	// Get the problem and language
	prob := s.GetProblem()
	language := s.GetLanguage()
	
	// Get the current user code
	code := s.GetCode()
	
	// Create a temporary directory for test execution
	testDir, err := ioutil.TempDir("", "algo-scales-test")
	if err != nil {
		return nil, false, fmt.Errorf("failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up when done
	
	// Run tests based on the language
	var results []interfaces.TestResult
	
	switch language {
	case "go":
		results, err = executeGoTests(testDir, prob, code)
	case "python":
		results, err = executePythonTests(testDir, prob, code)
	case "javascript":
		results, err = executeJavaScriptTests(testDir, prob, code)
	default:
		return nil, false, fmt.Errorf("unsupported language: %s", language)
	}
	
	if err != nil {
		return nil, false, err
	}
	
	// Check if all tests passed
	allPassed := true
	for _, result := range results {
		if !result.Passed {
			allPassed = false
			break
		}
	}
	
	return results, allPassed, nil
}

// executeGoTests runs tests for Go solutions
func executeGoTests(testDir string, prob *problem.Problem, code string) ([]interfaces.TestResult, error) {
	// Create main.go with the solution and test code
	mainFile := filepath.Join(testDir, "main.go")
	
	// Create the test file content
	testContent := `package main

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
		// Note: This would need to be customized based on the actual problem's function signature
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

	// Write the test file
	testFileContent := fmt.Sprintf(testContent, code, testCases.String())
	err := ioutil.WriteFile(mainFile, []byte(testFileContent), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write test file: %v", err)
	}
	
	// Build and run the test
	cmd := exec.Command("go", "run", mainFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	
	// Parse the results from stdout
	output := stdout.String()
	
	// Create results array
	results := make([]interfaces.TestResult, len(prob.TestCases))
	
	// Initialize with basic information
	for i, tc := range prob.TestCases {
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
	
	// If there were compile errors, include them in the results
	if err != nil && len(stderr.String()) > 0 {
		errorMsg := stderr.String()
		// Add error message to all failed tests
		for i := range results {
			if !results[i].Passed {
				results[i].Actual = fmt.Sprintf("Error: %s", errorMsg)
			}
		}
	}
	
	return results, nil
}

// executePythonTests runs tests for Python solutions
func executePythonTests(testDir string, prob *problem.Problem, code string) ([]interfaces.TestResult, error) {
	// Create a Python file with the solution and test code
	testFile := filepath.Join(testDir, "test_solution.py")
	
	// Create the test file content
	testContent := `
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
		testCases.WriteString(fmt.Sprintf("\n    # Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("    print(\"Test %d\")\n", i+1))
		testCases.WriteString(fmt.Sprintf("    input_str = '%s'\n", tc.Input))
		testCases.WriteString(fmt.Sprintf("    expected_str = '%s'\n", tc.Expected))
		
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

	// Write the test file
	testFileContent := fmt.Sprintf(testContent, code, testCases.String())
	err := ioutil.WriteFile(testFile, []byte(testFileContent), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write test file: %v", err)
	}
	
	// Run the test
	cmd := exec.Command("python", testFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	
	// Parse the results from stdout
	output := stdout.String()
	
	// Create results array
	results := make([]interfaces.TestResult, len(prob.TestCases))
	
	// Initialize with basic information
	for i, tc := range prob.TestCases {
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
	
	// If there were errors, include them in the results
	if err != nil && len(stderr.String()) > 0 {
		errorMsg := stderr.String()
		// Add error message to all failed tests
		for i := range results {
			if !results[i].Passed {
				results[i].Actual = fmt.Sprintf("Error: %s", errorMsg)
			}
		}
	}
	
	return results, nil
}

// executeJavaScriptTests runs tests for JavaScript solutions
func executeJavaScriptTests(testDir string, prob *problem.Problem, code string) ([]interfaces.TestResult, error) {
	// Create a JavaScript file with the solution and test code
	testFile := filepath.Join(testDir, "test_solution.js")
	
	// Create the test file content
	testContent := `
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

	// Write the test file
	testFileContent := fmt.Sprintf(testContent, code, testCases.String())
	err := ioutil.WriteFile(testFile, []byte(testFileContent), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write test file: %v", err)
	}
	
	// Run the test
	cmd := exec.Command("node", testFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	
	// Parse the results from stdout
	output := stdout.String()
	
	// Create results array
	results := make([]interfaces.TestResult, len(prob.TestCases))
	
	// Initialize with basic information
	for i, tc := range prob.TestCases {
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
	
	// If there were errors, include them in the results
	if err != nil && len(stderr.String()) > 0 {
		errorMsg := stderr.String()
		// Add error message to all failed tests
		for i := range results {
			if !results[i].Passed {
				results[i].Actual = fmt.Sprintf("Error: %s", errorMsg)
			}
		}
	}
	
	return results, nil
}