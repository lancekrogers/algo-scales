package template

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

func TestTemplateService(t *testing.T) {
	// Create a service
	service := NewService()
	
	// Verify the service implements the TemplateService interface
	var _ interfaces.TemplateService = service
	
	// Create a test problem
	testProblem := &interfaces.Problem{
		ID:          "test-problem",
		Title:       "Test Problem",
		Description: "This is a test problem",
		TestCases: []interfaces.TestCase{
			{Input: "5", Expected: "25"},
			{Input: "10", Expected: "100"},
		},
		Languages: []string{"go", "python", "javascript"},
	}
	
	// Test GetSupportedLanguages
	t.Run("GetSupportedLanguages", func(t *testing.T) {
		languages := service.GetSupportedLanguages()
		assert.Contains(t, languages, "go")
		assert.Contains(t, languages, "python")
		assert.Contains(t, languages, "javascript")
	})
	
	// Test GetTemplate for Go
	t.Run("GetTemplate_Go", func(t *testing.T) {
		template, err := service.GetTemplate(testProblem, "go")
		assert.NoError(t, err)
		assert.Contains(t, template, "package main")
		assert.Contains(t, template, "Test Problem")
		assert.Contains(t, template, "func solution()")
	})
	
	// Test GetTemplate for Python
	t.Run("GetTemplate_Python", func(t *testing.T) {
		template, err := service.GetTemplate(testProblem, "python")
		assert.NoError(t, err)
		assert.Contains(t, template, "def solution():")
		assert.Contains(t, template, "Test Problem")
		assert.Contains(t, template, "if __name__ == \"__main__\":")
	})
	
	// Test GetTemplate for JavaScript
	t.Run("GetTemplate_JavaScript", func(t *testing.T) {
		template, err := service.GetTemplate(testProblem, "javascript")
		assert.NoError(t, err)
		assert.Contains(t, template, "function solution()")
		assert.Contains(t, template, "Test Problem")
		assert.Contains(t, template, "function runTests()")
	})
	
	// Test GetTemplate for an unsupported language
	t.Run("GetTemplate_Unsupported", func(t *testing.T) {
		template, err := service.GetTemplate(testProblem, "unsupported")
		assert.NoError(t, err) // Should not error, but provide a generic template
		assert.Contains(t, template, "Test Problem")
		assert.Contains(t, template, "TODO")
	})
	
	// Test GetTestHarness for Go
	t.Run("GetTestHarness_Go", func(t *testing.T) {
		solutionCode := `package main

func square(n int) int {
	return n * n
}

func main() {
	// Test
}
`
		harness, err := service.GetTestHarness(testProblem, solutionCode, "go")
		assert.NoError(t, err)
		assert.Contains(t, harness, "square") // Should extract the function name
		assert.Contains(t, harness, "Test 1:")
		assert.Contains(t, harness, "Test 2:")
		assert.Contains(t, harness, "allPassed")
	})
	
	// Test GetTestHarness for Python
	t.Run("GetTestHarness_Python", func(t *testing.T) {
		solutionCode := `def square(n):
	return n * n

if __name__ == "__main__":
	print(square(5))
`
		harness, err := service.GetTestHarness(testProblem, solutionCode, "python")
		assert.NoError(t, err)
		assert.Contains(t, harness, "square") // Should extract the function name
		assert.Contains(t, harness, "Test 1:")
		assert.Contains(t, harness, "Test 2:")
		assert.Contains(t, harness, "all_passed")
	})
	
	// Test GetTestHarness for JavaScript
	t.Run("GetTestHarness_JavaScript", func(t *testing.T) {
		solutionCode := `function square(n) {
	return n * n;
}

function runTests() {
	console.log(square(5));
}
`
		harness, err := service.GetTestHarness(testProblem, solutionCode, "javascript")
		assert.NoError(t, err)
		assert.Contains(t, harness, "square") // Should extract the function name
		assert.Contains(t, harness, "Test 1:")
		assert.Contains(t, harness, "Test 2:")
		assert.Contains(t, harness, "allPassed")
	})
	
	// Test GetFunctionName extraction
	t.Run("GetFunctionName", func(t *testing.T) {
		// Go
		goCode := `func calculateSum(a, b int) int {
			return a + b
		}`
		goGenerator := NewGoGenerator()
		assert.Equal(t, "calculateSum", goGenerator.GetFunctionName(goCode))
		
		// Python
		pythonCode := `def calculate_sum(a, b):
			return a + b`
		pythonGenerator := NewPythonGenerator()
		assert.Equal(t, "calculate_sum", pythonGenerator.GetFunctionName(pythonCode))
		
		// JavaScript function
		jsCode := `function calculateSum(a, b) {
			return a + b;
		}`
		jsGenerator := NewJavaScriptGenerator()
		assert.Equal(t, "calculateSum", jsGenerator.GetFunctionName(jsCode))
		
		// JavaScript arrow function
		jsArrowCode := `const calculateSum = (a, b) => {
			return a + b;
		}`
		assert.Equal(t, "calculateSum", jsGenerator.GetFunctionName(jsArrowCode))
	})
}

func TestMockTemplateService(t *testing.T) {
	// Create a mock service
	mockService := NewMockService()
	
	// Verify the mock service implements the TemplateService interface
	var _ interfaces.TemplateService = mockService
	
	// Create a test problem
	testProblem := &interfaces.Problem{
		ID:          "test-problem",
		Title:       "Test Problem",
		Description: "This is a test problem",
	}
	
	// Set a custom template
	mockService.SetTemplate("test-problem", "go", "Custom Go template")
	template, err := mockService.GetTemplate(testProblem, "go")
	assert.NoError(t, err)
	assert.Equal(t, "Custom Go template", template)
	
	// Test default template for unset language
	template, err = mockService.GetTemplate(testProblem, "python")
	assert.NoError(t, err)
	assert.Contains(t, template, "Mock template for test-problem in python")
	
	// Set a custom test harness
	mockService.SetTestHarness("test-problem", "go", "Custom Go test harness")
	harness, err := mockService.GetTestHarness(testProblem, "code", "go")
	assert.NoError(t, err)
	assert.Equal(t, "Custom Go test harness", harness)
	
	// Set a custom function name
	mockService.SetFunctionName("func test()", "go", "testFunction")
	funcName, err := mockService.GetFunctionName("func test()", "go")
	assert.NoError(t, err)
	assert.Equal(t, "testFunction", funcName)
}