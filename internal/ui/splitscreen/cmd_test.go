package splitscreen

import (
	"fmt"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Mock tea.Program to avoid actually running the UI
type mockProgram struct {
	model tea.Model
	err   error
}

func (m *mockProgram) Run() (tea.Model, error) {
	return m.model, m.err
}

// TestStartWithSampleProblem tests the sample problem initialization
func TestStartWithSampleProblem(t *testing.T) {
	// Override os.Exit for testing
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()
	
	// Using a custom variable name to avoid redeclaration
	var testExitCode int
	osExit = func(code int) {
		testExitCode = code
		// Use the variable to avoid unused variable warning
		_ = testExitCode
	}
	
	// Mock out RunCLI to avoid actually running the program
	originalRunProg := runProgram
	defer func() { runProgram = originalRunProg }()
	
	// Use a temporary local function to avoid redeclaration
	runProgram = func(m tea.Model, opts ...tea.ProgramOption) (tea.Model, error) {
		// Access the model to verify it was correctly set up
		splitModel, ok := m.(Model)
		if !ok {
			t.Errorf("expected model to be a Model, got %T", m)
		}
		
		// Verify the problem got loaded into the model
		problem := splitModel.currentProblem
		if problem == nil {
			t.Errorf("problem not loaded into model")
		}
		
		// Return immediately to avoid actually running the UI
		return m, nil
	}
	
	// Run the function and check results
	err := StartWithSampleProblem()
	if err != nil {
		t.Errorf("StartWithSampleProblem returned error: %v", err)
	}
}

// TestRunCLI tests the RunCLI function by replacing the UI functions with test doubles
func TestRunCLI(t *testing.T) {
	// We need to create a fully isolated testing environment without TTY requirements
	
	// Save original function pointers
	origStartWithSample := runStartWithSampleProblem
	origGetProblem := runGetProblemAndStartUI
	
	// Create replacement functions for testing
	var (
		sampleCalled  bool
		problemCalled bool
		problemID     string
		sampleErr     error
		problemErr    error
	)
	
	// Create temporary mock functions
	runStartWithSampleProblem = func() error { 
		sampleCalled = true
		return sampleErr
	}
	
	runGetProblemAndStartUI = func(id string) error {
		problemCalled = true
		problemID = id
		return problemErr
	}
	
	// Restore original functions when test is done
	defer func() {
		runStartWithSampleProblem = origStartWithSample
		runGetProblemAndStartUI = origGetProblem
	}()
	
	// Capture stderr output
	oldStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()
	
	// Test Case 1: No args - should call StartWithSampleProblem and return 0
	sampleCalled = false
	problemCalled = false
	sampleErr = nil
	
	exitCode := RunCLI([]string{})
	
	if !sampleCalled {
		t.Errorf("RunCLI with no args should call StartWithSampleProblem")
	}
	
	if problemCalled {
		t.Errorf("RunCLI with no args should not call GetProblemAndStartUI")
	}
	
	if exitCode != 0 {
		t.Errorf("RunCLI should return 0 when successful, got %d", exitCode)
	}
	
	// Test Case 2: With problem ID - should call GetProblemAndStartUI and return 0
	sampleCalled = false
	problemCalled = false
	problemID = ""
	problemErr = nil
	
	exitCode = RunCLI([]string{"test-problem"})
	
	if sampleCalled {
		t.Errorf("RunCLI with problem ID should not call StartWithSampleProblem")
	}
	
	if !problemCalled {
		t.Errorf("RunCLI with problem ID should call GetProblemAndStartUI")
	}
	
	if problemID != "test-problem" {
		t.Errorf("RunCLI should pass correct problem ID, expected 'test-problem', got '%s'", problemID)
	}
	
	if exitCode != 0 {
		t.Errorf("RunCLI should return 0 when successful, got %d", exitCode)
	}
	
	// Test Case 3: Error from StartWithSampleProblem - should return 1
	sampleCalled = false
	problemCalled = false
	sampleErr = fmt.Errorf("sample error")
	
	exitCode = RunCLI([]string{})
	
	if !sampleCalled {
		t.Errorf("RunCLI with no args should call StartWithSampleProblem even if it will error")
	}
	
	if exitCode != 1 {
		t.Errorf("RunCLI should return 1 when StartWithSampleProblem fails, got %d", exitCode)
	}
	
	// Test Case 4: Error from GetProblemAndStartUI - should return 1
	sampleCalled = false
	problemCalled = false
	problemErr = fmt.Errorf("problem error")
	
	exitCode = RunCLI([]string{"test-problem"})
	
	if !problemCalled {
		t.Errorf("RunCLI with problem ID should call GetProblemAndStartUI even if it will error")
	}
	
	if exitCode != 1 {
		t.Errorf("RunCLI should return 1 when GetProblemAndStartUI fails, got %d", exitCode)
	}
	
	// Close the pipe to flush stderr
	w.Close()
}

// TestGetProblemAndStartUI tests the GetProblemAndStartUI function
func TestGetProblemAndStartUI(t *testing.T) {
	// Save and restore original functions
	originalStartUI := startUIProg
	originalGetProblem := getProblemByID
	defer func() {
		startUIProg = originalStartUI
		getProblemByID = originalGetProblem
	}()
	
	// Mock the dependencies
	var uiStarted bool
	var problemID string
	
	getProblemByID = func(id string) (*problem.Problem, error) {
		problemID = id
		return &problem.Problem{
			ID:    id,
			Title: "Test Problem",
		}, nil
	}
	
	startUIProg = func(p *problem.Problem) error {
		uiStarted = true
		return nil
	}
	
	// Call the function
	err := GetProblemAndStartUI("test-problem")
	
	// Verify the results
	if err != nil {
		t.Errorf("GetProblemAndStartUI returned error: %v", err)
	}
	
	if problemID != "test-problem" {
		t.Errorf("expected getProblemByID to be called with 'test-problem', got '%s'", problemID)
	}
	
	if !uiStarted {
		t.Errorf("expected startUIProg to be called")
	}
}

// TestGetProblemAndStartUIError tests the error handling in GetProblemAndStartUI
func TestGetProblemAndStartUIError(t *testing.T) {
	// Save and restore original functions
	originalGetProblem := getProblemByID
	defer func() {
		getProblemByID = originalGetProblem
	}()
	
	// Mock the dependencies to return an error
	getProblemByID = func(id string) (*problem.Problem, error) {
		return nil, fmt.Errorf("problem not found")
	}
	
	// Call the function
	err := GetProblemAndStartUI("non-existent")
	
	// Verify the error is returned
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// TestRunCLI needs access to these original function references for testing
var (
	// These variables are used within tests to store original references
	// that can be restored in defer statements
	originalOsExit = os.Exit
	originalStartUIProg = StartUI
	originalGetProblemByID = problem.GetByID
)