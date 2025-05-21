package session

import (
	"strings"
	"testing"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

func TestProblemFormatter_FormatDescription(t *testing.T) {
	formatter := NewProblemFormatter()
	
	prob := &problem.Problem{
		ID:                   "test_problem",
		Title:                "Test Problem",
		Description:          "This is a test problem",
		Difficulty:           "Easy",
		EstimatedTime:        30,
		Companies:            []string{"Test Corp", "Example Inc"},
		Patterns:             []string{"Array", "Two Pointers"},
		PatternExplanation:   "This problem uses array manipulation and two pointers technique.",
		Examples: []problem.Example{
			{
				Input:       "[1,2,3]",
				Output:      "[3,2,1]",
				Explanation: "Reverse the array",
			},
		},
		Constraints: []string{
			"1 <= n <= 1000",
			"Array elements are integers",
		},
		SolutionWalkthrough: []string{
			"Initialize two pointers",
			"Swap elements and move pointers",
			"Continue until pointers meet",
		},
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
	if !strings.Contains(description, "30 minutes") {
		t.Error("Expected estimated time to be present")
	}
	if !strings.Contains(description, "Test Corp") {
		t.Error("Expected companies to be present")
	}
	if !strings.Contains(description, "This is a test problem") {
		t.Error("Expected description to be present")
	}
	if !strings.Contains(description, "[1,2,3]") {
		t.Error("Expected example input to be present")
	}
	if !strings.Contains(description, "1 <= n <= 1000") {
		t.Error("Expected constraints to be present")
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
	
	prob := &problem.Problem{
		Title:              "Test Problem",
		Patterns:           []string{"Array", "Two Pointers"},
		PatternExplanation: "This problem uses array manipulation and two pointers technique.",
	}
	
	description := formatter.FormatDescription(prob, true, false)
	
	if !strings.Contains(description, "Pattern: Array, Two Pointers") {
		t.Error("Expected pattern to be present when showPattern is true")
	}
	if !strings.Contains(description, "array manipulation and two pointers") {
		t.Error("Expected pattern explanation to be present")
	}
}

func TestProblemFormatter_WithSolution(t *testing.T) {
	formatter := NewProblemFormatter()
	
	prob := &problem.Problem{
		Title: "Test Problem",
		SolutionWalkthrough: []string{
			"Initialize two pointers",
			"Swap elements and move pointers",
			"Continue until pointers meet",
		},
	}
	
	description := formatter.FormatDescription(prob, false, true)
	
	if !strings.Contains(description, "Solution Walkthrough") {
		t.Error("Expected solution walkthrough header to be present when showSolution is true")
	}
	if !strings.Contains(description, "1. Initialize two pointers") {
		t.Error("Expected solution steps to be present and numbered")
	}
	if !strings.Contains(description, "2. Swap elements and move pointers") {
		t.Error("Expected second solution step to be present")
	}
}