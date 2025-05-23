package session

import (
	"strings"
	"testing"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

func TestProblemFormatter_FormatDescription(t *testing.T) {
	formatter := NewProblemFormatter()
	
	prob := &interfaces.Problem{
		ID:          "test_problem",
		Title:       "Test Problem",
		Description: "This is a test problem",
		Difficulty:  "Easy",
		Companies:   []string{"Test Corp", "Example Inc"},
		Tags:        []string{"Array", "Two Pointers"},
		TestCases: []interfaces.TestCase{
			{
				Input:    "[1,2,3]",
				Expected: "[3,2,1]",
			},
		},
		Languages: []string{"go", "python", "javascript"},
	}
	
	// Test basic formatting without pattern or solution
	description := formatter.FormatDescription(prob, false, false)
	
	// Verify basic content is present
	if !strings.Contains(description, "Test Problem") {
		t.Error("Expected title to be present")
	}
	if !strings.Contains(description, "Easy") {
		t.Error("Expected difficulty to be present")
	}
	if !strings.Contains(description, "Test Corp") {
		t.Error("Expected companies to be present")
	}
	if !strings.Contains(description, "This is a test problem") {
		t.Error("Expected description to be present")
	}
	
	// Verify pattern and solution are NOT present when disabled
	if strings.Contains(description, "Pattern:") {
		t.Error("Pattern should not be present when showPattern is false")
	}
	if strings.Contains(description, "Solution Walkthrough") {
		t.Error("Solution should not be present when showSolution is false")
	}
}

func TestProblemFormatter_WithPattern(t *testing.T) {
	formatter := NewProblemFormatter()
	
	prob := &interfaces.Problem{
		Title:     "Test Problem",
		Tags:      []string{"Array", "Two Pointers"},
	}
	
	description := formatter.FormatDescription(prob, true, false)
	
	if !strings.Contains(description, "Pattern: Array, Two Pointers") {
		t.Error("Expected pattern to be present when showPattern is true")
	}
}

func TestProblemFormatter_WithSolution(t *testing.T) {
	formatter := NewProblemFormatter()
	
	prob := &interfaces.Problem{
		Title: "Test Problem",
	}
	
	description := formatter.FormatDescription(prob, false, true)
	
	// Since interfaces.Problem doesn't have SolutionWalkthrough,
	// the formatter will convert it to an empty array
	if !strings.Contains(description, "Solution Walkthrough") {
		t.Error("Expected solution walkthrough header to be present when showSolution is true")
	}
}