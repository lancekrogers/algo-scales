// Tests for list command

package cmd

import (
	"testing"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/stretchr/testify/assert"
)

// Mock problem.ListAll for testing
func mockListAll(problems []problem.Problem, err error) func() {
	original := problem.ListAll
	problem.ListAll = func() ([]problem.Problem, error) {
		return problems, err
	}
	return func() {
		problem.ListAll = original
	}
}

// Mock problem.ListPatterns for testing
func mockListPatterns(patterns map[string][]problem.Problem, err error) func() {
	original := problem.ListPatterns
	problem.ListPatterns = func() (map[string][]problem.Problem, error) {
		return patterns, err
	}
	return func() {
		problem.ListPatterns = original
	}
}

// Mock problem.ListByDifficulty for testing
func mockListByDifficulty(difficulties map[string][]problem.Problem, err error) func() {
	original := problem.ListByDifficulty
	problem.ListByDifficulty = func() (map[string][]problem.Problem, error) {
		return difficulties, err
	}
	return func() {
		problem.ListByDifficulty = original
	}
}

// Mock problem.ListByCompany for testing
func mockListByCompany(companies map[string][]problem.Problem, err error) func() {
	original := problem.ListByCompany
	problem.ListByCompany = func() (map[string][]problem.Problem, error) {
		return companies, err
	}
	return func() {
		problem.ListByCompany = original
	}
}

func TestListCommand(t *testing.T) {
	// Create sample problems for testing
	sampleProblems := []problem.Problem{
		{
			ID:         "test1",
			Title:      "Test Problem 1",
			Difficulty: "Easy",
			Patterns:   []string{"hash-map"},
		},
		{
			ID:         "test2",
			Title:      "Test Problem 2",
			Difficulty: "Medium",
			Patterns:   []string{"two-pointers"},
		},
	}

	t.Run("ListAll", func(t *testing.T) {
		// Mock ListAll to return sample problems
		restore := mockListAll(sampleProblems, nil)
		defer restore()

		// Execute list command
		output, err := executeCommand(rootCmd, "list")
		assert.NoError(t, err)

		// Check output contains problem info
		assert.Contains(t, output, "test1")
		assert.Contains(t, output, "Test Problem 1")
		assert.Contains(t, output, "Easy")
		assert.Contains(t, output, "test2")
		assert.Contains(t, output, "Test Problem 2")
		assert.Contains(t, output, "Medium")
	})

	t.Run("ListPatterns", func(t *testing.T) {
		// Create sample patterns
		patterns := map[string][]problem.Problem{
			"hash-map":     {sampleProblems[0]},
			"two-pointers": {sampleProblems[1]},
		}

		// Mock ListPatterns
		restore := mockListPatterns(patterns, nil)
		defer restore()

		// Execute patterns command
		output, err := executeCommand(rootCmd, "list", "patterns")
		assert.NoError(t, err)

		// Check output contains pattern info
		assert.Contains(t, output, "hash-map")
		assert.Contains(t, output, "two-pointers")
		assert.Contains(t, output, "test1")
		assert.Contains(t, output, "test2")
	})

	t.Run("ListDifficulties", func(t *testing.T) {
		// Create sample difficulties
		difficulties := map[string][]problem.Problem{
			"Easy":   {sampleProblems[0]},
			"Medium": {sampleProblems[1]},
		}

		// Mock ListByDifficulty
		restore := mockListByDifficulty(difficulties, nil)
		defer restore()

		// Execute difficulties command
		output, err := executeCommand(rootCmd, "list", "difficulties")
		assert.NoError(t, err)

		// Check output contains difficulty info
		assert.Contains(t, output, "Easy")
		assert.Contains(t, output, "Medium")
		assert.Contains(t, output, "test1")
		assert.Contains(t, output, "test2")
	})

	t.Run("ListCompanies", func(t *testing.T) {
		// Create sample companies
		companies := map[string][]problem.Problem{
			"Google": {sampleProblems[0]},
			"Amazon": {sampleProblems[1]},
		}

		// Mock ListByCompany
		restore := mockListByCompany(companies, nil)
		defer restore()

		// Execute companies command
		output, err := executeCommand(rootCmd, "list", "companies")
		assert.NoError(t, err)

		// Check output contains company info
		assert.Contains(t, output, "Google")
		assert.Contains(t, output, "Amazon")
		assert.Contains(t, output, "test1")
		assert.Contains(t, output, "test2")
	})
}
