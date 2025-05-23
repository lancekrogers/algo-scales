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
func (f *ProblemFormatterImpl) FormatDescription(prob *interfaces.Problem, showPattern bool, showSolution bool) string {
	// Convert to local problem type to access detailed fields
	// Note: This is a simplified conversion - in a real implementation,
	// we might want to fetch the full problem details from the repository
	localProb := f.convertToLocalProblem(prob)
	
	var description string

	// Problem header
	description += fmt.Sprintf("# %s\n\n", localProb.Title)
	description += fmt.Sprintf("**Difficulty**: %s\n", localProb.Difficulty)
	description += fmt.Sprintf("**Estimated Time**: %d minutes\n", localProb.EstimatedTime)
	description += fmt.Sprintf("**Companies**: %s\n\n", JoinStrings(localProb.Companies))

	// Problem description
	description += fmt.Sprintf("## Problem Statement\n\n%s\n\n", localProb.Description)

	// Examples
	description += "## Examples\n\n"
	for i, example := range localProb.Examples {
		description += fmt.Sprintf("### Example %d\n\n", i+1)
		description += fmt.Sprintf("**Input**: %s\n\n", example.Input)
		description += fmt.Sprintf("**Output**: %s\n\n", example.Output)
		if example.Explanation != "" {
			description += fmt.Sprintf("**Explanation**: %s\n\n", example.Explanation)
		}
	}

	// Constraints
	description += "## Constraints\n\n"
	for _, constraint := range localProb.Constraints {
		description += fmt.Sprintf("- %s\n", constraint)
	}
	description += "\n"

	// Pattern explanation (if in Learn mode)
	if showPattern {
		description += fmt.Sprintf("## Pattern: %s\n\n", JoinStrings(localProb.Patterns))
		description += fmt.Sprintf("%s\n\n", localProb.PatternExplanation)
	}

	// Solution walkthrough (if requested)
	if showSolution {
		description += "## Solution Walkthrough\n\n"
		for i, step := range localProb.SolutionWalkthrough {
			description += fmt.Sprintf("%d. %s\n", i+1, step)
		}
		description += "\n"
	}

	return description
}
// convertToLocalProblem converts an interfaces.Problem to a local problem.Problem
func (f *ProblemFormatterImpl) convertToLocalProblem(p *interfaces.Problem) problem.Problem {
	// Convert test cases
	testCases := make([]problem.TestCase, len(p.TestCases))
	for i, tc := range p.TestCases {
		testCases[i] = problem.TestCase{
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}
	
	// Create starter code map
	starterCode := make(map[string]string)
	for _, lang := range p.Languages {
		starterCode[lang] = ""
	}
	
	return problem.Problem{
		ID:                  p.ID,
		Title:               p.Title,
		Description:         p.Description,
		Difficulty:          p.Difficulty,
		Patterns:            p.Tags,
		Companies:           p.Companies,
		TestCases:           testCases,
		StarterCode:         starterCode,
		Solutions:           make(map[string]string),
		EstimatedTime:       30, // Default value
		Examples:            []problem.Example{}, // Empty for now
		Constraints:         []string{}, // Empty for now
		PatternExplanation:  "", // Empty for now
		SolutionWalkthrough: []string{}, // Empty for now
	}
}
