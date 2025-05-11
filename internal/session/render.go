// Problem rendering with highlighting

package session

import (
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/common/highlight"
)

// FormatProblemDescriptionWithHighlighting creates a formatted markdown description with syntax highlighting
func (s *SessionImpl) FormatProblemDescriptionWithHighlighting() string {
	// Create a syntax highlighter
	highlighter := highlight.NewSyntaxHighlighter("monokai")

	var description string

	// Problem header
	description += fmt.Sprintf("# %s\n\n", s.Problem.Title)
	description += fmt.Sprintf("**Difficulty**: %s\n", s.Problem.Difficulty)
	description += fmt.Sprintf("**Estimated Time**: %d minutes\n", s.Problem.EstimatedTime)
	description += fmt.Sprintf("**Companies**: %s\n\n", JoinStrings(s.Problem.Companies))

	// Problem description
	description += fmt.Sprintf("## Problem Statement\n\n%s\n\n", s.Problem.Description)

	// Examples
	description += "## Examples\n\n"
	for i, example := range s.Problem.Examples {
		description += fmt.Sprintf("### Example %d\n\n", i+1)
		description += fmt.Sprintf("**Input**: %s\n\n", example.Input)
		description += fmt.Sprintf("**Output**: %s\n\n", example.Output)
		if example.Explanation != "" {
			description += fmt.Sprintf("**Explanation**: %s\n\n", example.Explanation)
		}
	}

	// Constraints
	description += "## Constraints\n\n"
	for _, constraint := range s.Problem.Constraints {
		description += fmt.Sprintf("- %s\n", constraint)
	}
	description += "\n"

	// Pattern explanation (if in Learn mode or hints shown)
	if s.ShowPattern {
		description += fmt.Sprintf("## Pattern: %s\n\n", JoinStrings(s.Problem.Patterns))
		description += fmt.Sprintf("%s\n\n", s.Problem.PatternExplanation)
	}

	// Add starter code with syntax highlighting
	description += "## Your Task\n\n"
	if starterCode, ok := s.Problem.StarterCode[s.Options.Language]; ok {
		description += highlighter.RenderCodeBlock(starterCode, s.Options.Language)
		description += "\n\n"
	}

	// Solution walkthrough (if requested)
	if s.solutionShown {
		description += "## Solution Walkthrough\n\n"
		for i, step := range s.Problem.SolutionWalkthrough {
			description += fmt.Sprintf("%d. %s\n", i+1, step)
		}
		description += "\n"

		// Show solution code with syntax highlighting
		if solution, ok := s.Problem.Solutions[s.Options.Language]; ok {
			description += "## Solution Code\n\n"
			description += highlighter.RenderCodeBlock(solution, s.Options.Language)
			description += "\n\n"
		}

		// Show test cases
		description += "## Test Cases\n\n"
		for i, testCase := range s.Problem.TestCases {
			description += fmt.Sprintf("**Test %d**:\n", i+1)
			description += fmt.Sprintf("- Input: `%s`\n", testCase.Input)
			description += fmt.Sprintf("- Expected: `%s`\n\n", testCase.Expected)
		}
	}

	return description
}
