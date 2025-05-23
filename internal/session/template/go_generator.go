package template

import (
	"fmt"
	"regexp"
	"strings"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// GoGenerator generates Go code templates
type GoGenerator struct{}

// NewGoGenerator creates a new Go code generator
func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

// GetLanguage returns the language this generator supports
func (g *GoGenerator) GetLanguage() string {
	return "go"
}

// GetTemplate returns a code template for a problem
func (g *GoGenerator) GetTemplate(prob interfaces.Problem) string {
	// First check if a starter code is provided
	if starterCode, ok := prob.StarterCode["go"]; ok && starterCode != "" {
		return starterCode
	}
	
	// Otherwise generate a default template
	return fmt.Sprintf(`// %s
// %s

package main

import (
	"fmt"
)

// Implement your solution here
func solution() {
	// Your code goes here
}

func main() {
	// Test your solution here
	fmt.Println("Running tests for solution...")
	// Example:
	// result := solution(...)
	// fmt.Printf("Result: %%v\n", result)
}
`, prob.Title, sanitizeCommentText(prob.Description))
}

// GetTestHarness generates a test harness for Go
func (g *GoGenerator) GetTestHarness(prob interfaces.Problem, solutionCode string) string {
	// Extract function name from solution code
	funcName := g.GetFunctionName(solutionCode)
	if funcName == "" {
		funcName = "solution" // Default function name
	}
	
	// Create a test harness template
	testHarness := `package main

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
	} else {
		fmt.Println("All tests passed!")
	}
}
`
	
	// Generate test code for each test case
	var testCases strings.Builder
	for i, tc := range prob.TestCases {
		testCases.WriteString(fmt.Sprintf("\n\t// Test case %d\n", i+1))
		testCases.WriteString(fmt.Sprintf("\tfmt.Printf(\"Test %d: \")\n", i+1))
		testCases.WriteString(fmt.Sprintf("\t{\n\t\tinputStr := `%s`\n", tc.Input))
		testCases.WriteString(fmt.Sprintf("\t\texpectedStr := `%s`\n", tc.Expected))
		
		// Parse input and call function with parameters
		testCases.WriteString("\t\t// Call the solution function\n")
		testCases.WriteString(fmt.Sprintf("\t\tresult := %s() // Add parameters as needed\n", funcName))
		
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
	
	return fmt.Sprintf(testHarness, solutionCode, testCases.String())
}

// GetFunctionName extracts the function name from Go code
func (g *GoGenerator) GetFunctionName(code string) string {
	re := regexp.MustCompile(`func\s+([a-zA-Z0-9_]+)\s*\(`)
	matches := re.FindStringSubmatch(code)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}