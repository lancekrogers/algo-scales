// Tests for session module

package session

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test fixtures
func getTestProblem() *problem.Problem {
	return &problem.Problem{
		ID:            "test-problem",
		Title:         "Test Problem",
		Difficulty:    "Easy",
		Patterns:      []string{"hash-map"},
		EstimatedTime: 15,
		Companies:     []string{"Test Company"},
		Description:   "This is a test problem description.",
		Examples: []problem.Example{
			{
				Input:       "nums = [1,2,3], target = 5",
				Output:      "[1,2]",
				Explanation: "Because nums[1] + nums[2] == 5",
			},
		},
		Constraints:         []string{"1 <= nums.length <= 100"},
		PatternExplanation:  "This is a test pattern explanation.",
		SolutionWalkthrough: []string{"Step 1", "Step 2"},
		StarterCode: map[string]string{
			"go":         "func solution() {}\n",
			"python":     "def solution():\n    pass\n",
			"javascript": "function solution() {}\n",
		},
		Solutions: map[string]string{
			"go":         "func solution() { return [1, 2] }\n",
			"python":     "def solution():\n    return [1, 2]\n",
			"javascript": "function solution() { return [1, 2]; }\n",
		},
		TestCases: []problem.TestCase{
			{
				Input:    "[1,2,3], 5",
				Expected: "[1,2]",
			},
		},
	}
}

// Mock selectProblem for testing
func mockSelectProblem(p *problem.Problem, err error) func() {
	original := selectProblem
	selectProblem = func(pattern, difficulty string) (*problem.Problem, error) {
		return p, err
	}
	return func() {
		selectProblem = original
	}
}

// Mock selectCramProblem for testing
func mockSelectCramProblem(p *problem.Problem, err error) func() {
	original := selectCramProblem
	selectCramProblem = func() (*problem.Problem, error) {
		return p, err
	}
	return func() {
		selectCramProblem = original
	}
}

// Mock problem.GetByID for testing
func mockGetByID(p *problem.Problem, err error) func() {
	original := problem.GetByID
	problem.GetByID = func(id string) (*problem.Problem, error) {
		return p, err
	}
	return func() {
		problem.GetByID = original
	}
}

func TestCreateWorkspace(t *testing.T) {
	testProblem := getTestProblem()

	// Create a session
	session := &Session{
		Problem: testProblem,
		Options: Options{
			Mode:     LearnMode,
			Language: "go",
			Timer:    30,
		},
		ShowHints:    true,
		ShowPattern:  true,
		ShowSolution: false,
	}

	// Try to create the workspace
	err := session.createWorkspace()
	require.NoError(t, err)
	defer os.RemoveAll(session.Workspace)

	// Verify workspace was created
	_, err = os.Stat(session.Workspace)
	assert.NoError(t, err, "Workspace directory should exist")

	// Verify problem description file was created
	descriptionFile := filepath.Join(session.Workspace, "problem.md")
	_, err = os.Stat(descriptionFile)
	assert.NoError(t, err, "Description file should exist")

	// Verify code file was created
	codeFile := filepath.Join(session.Workspace, "solution.go")
	_, err = os.Stat(codeFile)
	assert.NoError(t, err, "Code file should exist")

	// Verify code file content
	codeContent, err := ioutil.ReadFile(codeFile)
	require.NoError(t, err)
	assert.Equal(t, testProblem.StarterCode["go"], string(codeContent))
}

func TestFormatProblemDescription(t *testing.T) {
	testProblem := getTestProblem()

	t.Run("LearnMode", func(t *testing.T) {
		// Create a session in Learn mode
		session := &Session{
			Problem: testProblem,
			Options: Options{
				Mode: LearnMode,
			},
			ShowHints:    true,
			ShowPattern:  true,
			ShowSolution: false,
		}

		// Format the description
		description := session.FormatProblemDescription()

		// Verify content
		assert.Contains(t, description, testProblem.Title)
		assert.Contains(t, description, testProblem.Description)
		assert.Contains(t, description, testProblem.PatternExplanation) // Should include pattern in Learn mode
		assert.NotContains(t, description, "Solution Walkthrough")      // Should not include solution
	})

	t.Run("PracticeMode", func(t *testing.T) {
		// Create a session in Practice mode
		session := &Session{
			Problem: testProblem,
			Options: Options{
				Mode: PracticeMode,
			},
			ShowHints:    false,
			ShowPattern:  false,
			ShowSolution: false,
		}

		// Format the description
		description := session.FormatProblemDescription()

		// Verify content
		assert.Contains(t, description, testProblem.Title)
		assert.Contains(t, description, testProblem.Description)
		assert.NotContains(t, description, testProblem.PatternExplanation) // Should not include pattern in Practice mode
		assert.NotContains(t, description, "Solution Walkthrough")         // Should not include solution
	})

	t.Run("WithSolution", func(t *testing.T) {
		// Create a session with solution shown
		session := &Session{
			Problem: testProblem,
			Options: Options{
				Mode: PracticeMode,
			},
			ShowHints:    false,
			ShowPattern:  false,
			ShowSolution: true,
		}

		// Format the description
		description := session.FormatProblemDescription()

		// Verify content
		assert.Contains(t, description, testProblem.Title)
		assert.Contains(t, description, testProblem.Description)
		assert.Contains(t, description, "Solution Walkthrough") // Should include solution
		assert.Contains(t, description, "Step 1")
		assert.Contains(t, description, "Step 2")
	})
}

func TestSelectProblem(t *testing.T) {
	// This would be a more comprehensive test with actual problem data in testdata
	// For now, we'll just mock the behavior

	t.Run("FilterByPattern", func(t *testing.T) {
		// Mock ListAll to return specific test problems
		origListAll := problem.ListAll
		defer func() { problem.ListAll = origListAll }()

		problem.ListAll = func() ([]problem.Problem, error) {
			return []problem.Problem{
				{
					ID:         "problem1",
					Title:      "Problem 1",
					Difficulty: "Easy",
					Patterns:   []string{"hash-map"},
				},
				{
					ID:         "problem2",
					Title:      "Problem 2",
					Difficulty: "Medium",
					Patterns:   []string{"two-pointers"},
				},
			}, nil
		}

		// Select a problem by pattern
		result, err := selectProblem("hash-map", "")
		require.NoError(t, err)
		assert.Equal(t, "problem1", result.ID)
	})

	t.Run("FilterByDifficulty", func(t *testing.T) {
		// Mock ListAll to return specific test problems
		origListAll := problem.ListAll
		defer func() { problem.ListAll = origListAll }()

		problem.ListAll = func() ([]problem.Problem, error) {
			return []problem.Problem{
				{
					ID:         "problem1",
					Title:      "Problem 1",
					Difficulty: "Easy",
					Patterns:   []string{"hash-map"},
				},
				{
					ID:         "problem2",
					Title:      "Problem 2",
					Difficulty: "Medium",
					Patterns:   []string{"two-pointers"},
				},
			}, nil
		}

		// Select a problem by difficulty
		result, err := selectProblem("", "Medium")
		require.NoError(t, err)
		assert.Equal(t, "problem2", result.ID)
	})

	t.Run("NoMatchingProblems", func(t *testing.T) {
		// Mock ListAll to return specific test problems
		origListAll := problem.ListAll
		defer func() { problem.ListAll = origListAll }()

		problem.ListAll = func() ([]problem.Problem, error) {
			return []problem.Problem{
				{
					ID:         "problem1",
					Title:      "Problem 1",
					Difficulty: "Easy",
					Patterns:   []string{"hash-map"},
				},
			}, nil
		}

		// Try to select a problem with non-matching criteria
		_, err := selectProblem("non-existent", "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no problems match")
	})
}

func TestLanguageExtension(t *testing.T) {
	testCases := []struct {
		language string
		expected string
	}{
		{"go", "go"},
		{"python", "py"},
		{"javascript", "js"},
		{"unknown", "txt"},
	}

	for _, tc := range testCases {
		t.Run(tc.language, func(t *testing.T) {
			result := languageExtension(tc.language)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestContainsPattern(t *testing.T) {
	patterns := []string{"hash-map", "two-pointers", "sliding-window"}

	assert.True(t, containsPattern(patterns, "hash-map"))
	assert.True(t, containsPattern(patterns, "two-pointers"))
	assert.True(t, containsPattern(patterns, "sliding-window"))
	assert.False(t, containsPattern(patterns, "dfs"))
	assert.False(t, containsPattern(patterns, ""))
}

func TestJoinStrings(t *testing.T) {
	testCases := []struct {
		strings  []string
		expected string
	}{
		{[]string{"a", "b", "c"}, "a, b, c"},
		{[]string{"test"}, "test"},
		{[]string{}, ""},
	}

	for _, tc := range testCases {
		t.Run("join", func(t *testing.T) {
			result := joinStrings(tc.strings)
			assert.Equal(t, tc.expected, result)
		})
	}
}
