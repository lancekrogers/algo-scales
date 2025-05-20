package template

import (
	"fmt"
	"regexp"
	"strings"
	
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// JavaScriptGenerator generates JavaScript code templates
type JavaScriptGenerator struct{}

// NewJavaScriptGenerator creates a new JavaScript code generator
func NewJavaScriptGenerator() *JavaScriptGenerator {
	return &JavaScriptGenerator{}
}

// GetLanguage returns the language this generator supports
func (g *JavaScriptGenerator) GetLanguage() string {
	return "javascript"
}

// GetTemplate returns a code template for a problem
func (g *JavaScriptGenerator) GetTemplate(prob *problem.Problem) string {
	// First check if a starter code is provided
	if starterCode, ok := prob.StarterCode["javascript"]; ok && starterCode != "" {
		return starterCode
	}
	
	// Otherwise generate a default template
	return fmt.Sprintf(`// %s
// %s

// Implement your solution here
function solution() {
    // Your code goes here
}

// Test your solution here
function runTests() {
    console.log("Running tests for solution...");
    // Example:
    // const result = solution(...);
    // console.log("Result:", result);
}

// Run tests
runTests();
`, prob.Title, sanitizeCommentText(prob.Description))
}

// GetTestHarness generates a test harness for JavaScript
func (g *JavaScriptGenerator) GetTestHarness(prob *problem.Problem, solutionCode string) string {
	// Extract function name from solution code
	funcName := g.GetFunctionName(solutionCode)
	if funcName == "" {
		funcName = "solution" // Default function name
	}
	
	// Create a test harness template
	testHarness := `
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
    console.log("Some tests failed.");
    process.exit(1);
} else {
    console.log("All tests passed!");
}
`
	
	// Generate test code for each test case
	var testCases strings.Builder
	for i, tc := range prob.TestCases {
		testCases.WriteString(fmt.Sprintf("\n    // Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("    process.stdout.write(\"Test %d: \");\n", i+1))
		testCases.WriteString(fmt.Sprintf("    const inputStr = '%s';\n", tc.Input))
		testCases.WriteString(fmt.Sprintf("    const expectedStr = '%s';\n", tc.Expected))
		
		// Call function with parameters
		testCases.WriteString("    try {\n")
		testCases.WriteString(fmt.Sprintf("        const result = %s(); // Add parameters as needed\n", funcName))
		
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
	
	return fmt.Sprintf(testHarness, solutionCode, testCases.String())
}

// GetFunctionName extracts the function name from JavaScript code
func (g *JavaScriptGenerator) GetFunctionName(code string) string {
	// Match both function declarations and arrow functions
	re := regexp.MustCompile(`function\s+([a-zA-Z0-9_$]+)\s*\(|const\s+([a-zA-Z0-9_$]+)\s*=\s*\(?.*\)?\s*=>`)
	matches := re.FindStringSubmatch(code)
	if len(matches) >= 2 {
		if matches[1] != "" {
			return matches[1] // Regular function
		}
		if len(matches) >= 3 && matches[2] != "" {
			return matches[2] // Arrow function
		}
	}
	return ""
}