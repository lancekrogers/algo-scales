package session

import (
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// ProblemFormatterImpl implements the ProblemFormatter interface
type ProblemFormatterImpl struct{}

// NewProblemFormatter creates a new problem formatter
func NewProblemFormatter() interfaces.ProblemFormatter {
	return &ProblemFormatterImpl{}
}

// FormatDescription returns formatted problem description
func (f *ProblemFormatterImpl) FormatDescription(prob *problem.Problem, showPattern bool, showSolution bool) string {
	var description string

	// Problem header
	description += fmt.Sprintf("# %s\n\n", prob.Title)
	description += fmt.Sprintf("**Difficulty**: %s\n", prob.Difficulty)
	description += fmt.Sprintf("**Estimated Time**: %d minutes\n", prob.EstimatedTime)
	description += fmt.Sprintf("**Companies**: %s\n\n", JoinStrings(prob.Companies))

	// Problem description
	description += fmt.Sprintf("## Problem Statement\n\n%s\n\n", prob.Description)

	// Examples
	description += "## Examples\n\n"
	for i, example := range prob.Examples {
		description += fmt.Sprintf("### Example %d\n\n", i+1)
		description += fmt.Sprintf("**Input**: %s\n\n", example.Input)
		description += fmt.Sprintf("**Output**: %s\n\n", example.Output)
		if example.Explanation != "" {
			description += fmt.Sprintf("**Explanation**: %s\n\n", example.Explanation)
		}
	}

	// Constraints
	description += "## Constraints\n\n"
	for _, constraint := range prob.Constraints {
		description += fmt.Sprintf("- %s\n", constraint)
	}
	description += "\n"

	// Pattern explanation (if in Learn mode)
	if showPattern {
		description += fmt.Sprintf("## Pattern: %s\n\n", JoinStrings(prob.Patterns))
		description += fmt.Sprintf("%s\n\n", prob.PatternExplanation)
	}

	// Solution walkthrough (if requested)
	if showSolution {
		description += "## Solution Walkthrough\n\n"
		for i, step := range prob.SolutionWalkthrough {
			description += fmt.Sprintf("%d. %s\n", i+1, step)
		}
		description += "\n"
	}

	return description
}