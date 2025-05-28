package ai

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

func TestPromptBuilder(t *testing.T) {
	pb := NewPromptBuilder()

	// Test problem for prompts
	testProblem := problem.Problem{
		ID:          "two_sum",
		Title:       "Two Sum",
		Patterns:    []string{"hash-map"},
		Difficulty:  "easy",
		Description: "Find two numbers that add up to target",
		Constraints: []string{"Array length >= 2", "One solution exists"},
	}

	t.Run("BuildHintPrompt", func(t *testing.T) {
		tests := []struct {
			name     string
			level    int
			userCode string
			contains []string
		}{
			{
				name:     "Level 1 without code",
				level:    1,
				userCode: "",
				contains: []string{"Two Sum", "hash-map", "Level 1:"},
			},
			{
				name:     "Level 2 with code",
				level:    2,
				userCode: "func twoSum() {}",
				contains: []string{"Two Sum", "User's current approach", "func twoSum()"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				prompt, err := pb.BuildHintPrompt(testProblem, tt.userCode, tt.level)
				if err != nil {
					t.Fatalf("Failed to build hint prompt: %v", err)
				}

				for _, expected := range tt.contains {
					if !strings.Contains(prompt, expected) {
						t.Errorf("Prompt missing expected content: %s", expected)
					}
				}
			})
		}
	})

	t.Run("BuildReviewPrompt", func(t *testing.T) {
		code := `func twoSum(nums []int, target int) []int {
    return []int{}
}`
		prompt, err := pb.BuildReviewPrompt(testProblem, code, "go")
		if err != nil {
			t.Fatalf("Failed to build review prompt: %v", err)
		}

		expectedContents := []string{
			"Two Sum",
			"hash-map",
			"func twoSum",
			"Correctness",
			"Code quality",
		}

		for _, expected := range expectedContents {
			if !strings.Contains(prompt, expected) {
				t.Errorf("Review prompt missing: %s", expected)
			}
		}
	})

	t.Run("BuildPatternPrompt", func(t *testing.T) {
		examples := []problem.Problem{testProblem}
		prompt, err := pb.BuildPatternPrompt("hash-map", examples)
		if err != nil {
			t.Fatalf("Failed to build pattern prompt: %v", err)
		}

		if !strings.Contains(prompt, "hash-map") {
			t.Error("Pattern prompt missing pattern name")
		}
		if !strings.Contains(prompt, "Two Sum") {
			t.Error("Pattern prompt missing example problem")
		}
	})
}

func TestSystemPrompts(t *testing.T) {
	sp := NewSystemPrompts()

	prompts := map[string]string{
		"Tutor":       sp.GetTutorPrompt(),
		"Interviewer": sp.GetInterviewerPrompt(),
		"Reviewer":    sp.GetReviewerPrompt(),
		"Debugger":    sp.GetDebuggerPrompt(),
	}

	for name, prompt := range prompts {
		if len(prompt) < 50 {
			t.Errorf("%s prompt seems too short: %d chars", name, len(prompt))
		}
		if !strings.Contains(strings.ToLower(prompt), strings.ToLower(name)) && name != "Tutor" {
			t.Errorf("%s prompt doesn't mention its role", name)
		}
	}
}

func TestResponseFormatter(t *testing.T) {
	rf := NewResponseFormatter()

	t.Run("FormatHint", func(t *testing.T) {
		tests := []struct {
			level    int
			expected string
		}{
			{1, "💡 General Approach"},
			{2, "🔍 Specific Guidance"},
			{3, "📝 Implementation Details"},
		}

		for _, tt := range tests {
			result := rf.FormatHint(tt.level, "Test hint")
			if !strings.HasPrefix(result, tt.expected) {
				t.Errorf("Level %d: expected prefix %s, got %s", tt.level, tt.expected, result)
			}
		}
	})

	t.Run("FormatError", func(t *testing.T) {
		err := fmt.Errorf("test error")
		result := rf.FormatError(err)
		if !strings.Contains(result, "❌") || !strings.Contains(result, "test error") {
			t.Errorf("Error format incorrect: %s", result)
		}
	})

	t.Run("FormatSuccess", func(t *testing.T) {
		result := rf.FormatSuccess("Operation completed")
		if !strings.Contains(result, "✅") || !strings.Contains(result, "Operation completed") {
			t.Errorf("Success format incorrect: %s", result)
		}
	})
}