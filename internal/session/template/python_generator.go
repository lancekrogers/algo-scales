package template

import (
	"fmt"
	"regexp"
	"strings"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// PythonGenerator generates Python code templates
type PythonGenerator struct{}

// NewPythonGenerator creates a new Python code generator
func NewPythonGenerator() *PythonGenerator {
	return &PythonGenerator{}
}

// GetLanguage returns the language this generator supports
func (g *PythonGenerator) GetLanguage() string {
	return "python"
}

// GetTemplate returns a code template for a problem
func (g *PythonGenerator) GetTemplate(prob interfaces.Problem) string {
	// First check if a starter code is provided
	if starterCode, ok := prob.StarterCode["python"]; ok && starterCode != "" {
		return starterCode
	}
	
	// Otherwise generate a default template
	return fmt.Sprintf(`# %s
# %s

def solution():
    """
    Implement your solution here.
    
    Step 1: Understand the problem
    - Read the problem description carefully
    - Identify input/output requirements
    - Consider edge cases
    
    Step 2: Plan your approach
    - What algorithm pattern applies here?
    - What data structures do you need?
    - What's the time/space complexity?
    
    Step 3: Implement your solution
    - Replace this with your actual implementation
    """
    # Your implementation here
    return None  # Update return value as needed

# Test your solution here
if __name__ == "__main__":
    # Example:
    # result = solution(...)
    # print(f"Result: {result}")
    print("Running tests for solution...")
`, prob.Title, sanitizeCommentText(prob.Description))
}

// GetTestHarness generates a test harness for Python
func (g *PythonGenerator) GetTestHarness(prob interfaces.Problem, solutionCode string) string {
	// Extract function name from solution code
	funcName := g.GetFunctionName(solutionCode)
	if funcName == "" {
		funcName = "solution" // Default function name
	}
	
	// Create a test harness template
	testHarness := `
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
    else:
        print("All tests passed!")
`
	
	// Generate test code for each test case
	var testCases strings.Builder
	for i, tc := range prob.TestCases {
		testCases.WriteString(fmt.Sprintf("\n    # Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("    print(\"Test %d: \", end=\"\")\n", i+1))
		testCases.WriteString(fmt.Sprintf("    input_str = '%s'\n", tc.Input))
		testCases.WriteString(fmt.Sprintf("    expected_str = '%s'\n", tc.Expected))
		
		// Call function with parameters
		testCases.WriteString("    try:\n")
		testCases.WriteString(fmt.Sprintf("        result = %s() # Add parameters as needed\n", funcName))
		
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
	
	return fmt.Sprintf(testHarness, solutionCode, testCases.String())
}

// GetFunctionName extracts the function name from Python code
func (g *PythonGenerator) GetFunctionName(code string) string {
	re := regexp.MustCompile(`def\s+([a-zA-Z0-9_]+)\s*\(`)
	matches := re.FindStringSubmatch(code)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}