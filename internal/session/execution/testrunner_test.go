package execution

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

func TestRunnerRegistry(t *testing.T) {
	// Create registry
	registry := NewRunnerRegistry()
	
	// Verify default runners were registered
	langs := registry.GetSupportedLanguages()
	assert.Contains(t, langs, "go")
	assert.Contains(t, langs, "python")
	assert.Contains(t, langs, "javascript")
	
	// Get runner for each language
	goRunner, err := registry.GetRunner("go")
	assert.NoError(t, err)
	assert.Equal(t, "go", goRunner.GetLanguage())
	
	pyRunner, err := registry.GetRunner("python")
	assert.NoError(t, err)
	assert.Equal(t, "python", pyRunner.GetLanguage())
	
	jsRunner, err := registry.GetRunner("javascript")
	assert.NoError(t, err)
	assert.Equal(t, "javascript", jsRunner.GetLanguage())
	
	// Try getting a non-existent runner
	_, err = registry.GetRunner("nonexistent")
	assert.Error(t, err)
}

// MockTestRunner for testing
type MockTestRunner struct {
	BaseTestRunner
	executeFn func(ctx context.Context, prob *interfaces.Problem, code string, timeout time.Duration) ([]interfaces.TestResult, bool, error)
}

// Implement TestRunner interface
func (m *MockTestRunner) ExecuteTests(ctx context.Context, prob *interfaces.Problem, code string, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
	return m.executeFn(ctx, prob, code, timeout)
}

func (m *MockTestRunner) GenerateTestCode(prob *interfaces.Problem, solutionCode string) (string, error) {
	return "mock test code", nil
}

func TestMockTestRunner(t *testing.T) {
	
	// Create a mock test runner
	mockRunner := &MockTestRunner{
		BaseTestRunner: NewBaseTestRunner("mock"),
		executeFn: func(ctx context.Context, prob *interfaces.Problem, code string, timeout time.Duration) ([]interfaces.TestResult, bool, error) {
			results := []interfaces.TestResult{
				{
					Input:    "input1",
					Expected: "expected1",
					Actual:   "expected1",
					Passed:   true,
				},
				{
					Input:    "input2",
					Expected: "expected2",
					Actual:   "wrong",
					Passed:   false,
				},
			}
			return results, false, nil
		},
	}
	
	// Register the mock runner
	registry := NewRunnerRegistry()
	err := registry.RegisterRunner(mockRunner)
	assert.NoError(t, err)
	
	// Get the mock runner
	runner, err := registry.GetRunner("mock")
	assert.NoError(t, err)
	assert.Equal(t, "mock", runner.GetLanguage())
	
	// Execute tests with the mock runner
	testProblem := &interfaces.Problem{
		ID:    "test-problem",
		Title: "Test Problem",
		TestCases: []interfaces.TestCase{
			{Input: "input1", Expected: "expected1"},
			{Input: "input2", Expected: "expected2"},
		},
	}
	
	results, allPassed, err := runner.ExecuteTests(context.Background(), testProblem, "mock code", 1*time.Second)
	assert.NoError(t, err)
	assert.False(t, allPassed)
	assert.Len(t, results, 2)
	assert.True(t, results[0].Passed)
	assert.False(t, results[1].Passed)
}

func TestHelperFunctions(t *testing.T) {
	// Test parsing test output
	testOutput := `Test 1
✅ PASSED
Test 2
❌ FAILED
Expected: result2
Got: wrong
Test 3
✅ PASSED`

	testCases := []interfaces.TestCase{
		{Input: "input1", Expected: "result1"},
		{Input: "input2", Expected: "result2"},
		{Input: "input3", Expected: "result3"},
	}
	
	results := parseTestOutput(testOutput, testCases)
	assert.Len(t, results, 3)
	assert.True(t, results[0].Passed)
	assert.False(t, results[1].Passed)
	assert.Equal(t, "wrong", results[1].Actual)
	assert.True(t, results[2].Passed)
	
	// Test adding error to results
	errorMsg := "compilation error"
	results = addErrorToResults(results, errorMsg)
	assert.Equal(t, "Error: compilation error", results[1].Actual)
	
	// Test all tests passed
	assert.False(t, allTestsPassed(results))
	
	// Make all tests pass
	for i := range results {
		results[i].Passed = true
	}
	assert.True(t, allTestsPassed(results))
}