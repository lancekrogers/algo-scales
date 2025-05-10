// Tests for problem module

// internal/problem/problem_test.go
package problem

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByID(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := ioutil.TempDir("", "algo-scales-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up test problem
	problem := Problem{
		ID:         "test-problem",
		Title:      "Test Problem",
		Difficulty: "Easy",
		Patterns:   []string{"hash-map"},
	}

	// Create pattern directory
	patternDir := filepath.Join(tempDir, "problems", "hash-map")
	err = os.MkdirAll(patternDir, 0755)
	require.NoError(t, err)

	// Save problem file
	problemData, err := json.MarshalIndent(problem, "", "  ")
	require.NoError(t, err)
	problemFile := filepath.Join(patternDir, "test-problem.json")
	err = ioutil.WriteFile(problemFile, problemData, 0644)
	require.NoError(t, err)

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	defer func() { getConfigDir = origGetConfigDir }()
	getConfigDir = func() string {
		return tempDir
	}

	// Test GetByID
	t.Run("ExistingProblem", func(t *testing.T) {
		result, err := GetByID("test-problem")
		require.NoError(t, err)
		assert.Equal(t, problem.ID, result.ID)
		assert.Equal(t, problem.Title, result.Title)
		assert.Equal(t, problem.Difficulty, result.Difficulty)
		assert.Equal(t, problem.Patterns, result.Patterns)
	})

	t.Run("NonExistentProblem", func(t *testing.T) {
		_, err := GetByID("non-existent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "problem not found")
	})
}

func TestListAll(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := ioutil.TempDir("", "algo-scales-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up test problems in multiple pattern directories
	problems := []Problem{
		{
			ID:         "problem1",
			Title:      "Problem 1",
			Difficulty: "Easy",
			Patterns:   []string{"hash-map", "two-pointers"},
		},
		{
			ID:         "problem2",
			Title:      "Problem 2",
			Difficulty: "Medium",
			Patterns:   []string{"sliding-window"},
		},
	}

	// Create pattern directories and save problems
	for _, problem := range problems {
		for _, pattern := range problem.Patterns {
			patternDir := filepath.Join(tempDir, "problems", pattern)
			err = os.MkdirAll(patternDir, 0755)
			require.NoError(t, err)

			problemData, err := json.MarshalIndent(problem, "", "  ")
			require.NoError(t, err)
			problemFile := filepath.Join(patternDir, problem.ID+".json")
			err = ioutil.WriteFile(problemFile, problemData, 0644)
			require.NoError(t, err)
		}
	}

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	defer func() { getConfigDir = origGetConfigDir }()
	getConfigDir = func() string {
		return tempDir
	}

	// Test ListAll
	t.Run("ListAllProblems", func(t *testing.T) {
		results, err := ListAll()
		require.NoError(t, err)
		assert.Len(t, results, 2) // Should deduplicate across pattern directories

		// Check problem IDs
		ids := make(map[string]bool)
		for _, p := range results {
			ids[p.ID] = true
		}
		assert.True(t, ids["problem1"])
		assert.True(t, ids["problem2"])
	})
}

func TestListPatterns(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := ioutil.TempDir("", "algo-scales-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up test problems in multiple pattern directories
	problems := []Problem{
		{
			ID:         "problem1",
			Title:      "Problem 1",
			Difficulty: "Easy",
			Patterns:   []string{"hash-map", "two-pointers"},
		},
		{
			ID:         "problem2",
			Title:      "Problem 2",
			Difficulty: "Medium",
			Patterns:   []string{"sliding-window"},
		},
	}

	// Create pattern directories and save problems
	for _, problem := range problems {
		for _, pattern := range problem.Patterns {
			patternDir := filepath.Join(tempDir, "problems", pattern)
			err = os.MkdirAll(patternDir, 0755)
			require.NoError(t, err)

			problemData, err := json.MarshalIndent(problem, "", "  ")
			require.NoError(t, err)
			problemFile := filepath.Join(patternDir, problem.ID+".json")
			err = ioutil.WriteFile(problemFile, problemData, 0644)
			require.NoError(t, err)
		}
	}

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	defer func() { getConfigDir = origGetConfigDir }()
	getConfigDir = func() string {
		return tempDir
	}

	// Test ListPatterns
	t.Run("ListByPattern", func(t *testing.T) {
		patterns, err := ListPatterns()
		require.NoError(t, err)

		// Check pattern counts
		assert.Len(t, patterns["hash-map"], 1)
		assert.Len(t, patterns["two-pointers"], 1)
		assert.Len(t, patterns["sliding-window"], 1)
	})
}

func TestListByDifficulty(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := ioutil.TempDir("", "algo-scales-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up test problems with different difficulties
	problems := []Problem{
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
			Patterns:   []string{"sliding-window"},
		},
		{
			ID:         "problem3",
			Title:      "Problem 3",
			Difficulty: "Hard",
			Patterns:   []string{"dynamic-programming"},
		},
	}

	// Create pattern directories and save problems
	for _, problem := range problems {
		for _, pattern := range problem.Patterns {
			patternDir := filepath.Join(tempDir, "problems", pattern)
			err = os.MkdirAll(patternDir, 0755)
			require.NoError(t, err)

			problemData, err := json.MarshalIndent(problem, "", "  ")
			require.NoError(t, err)
			problemFile := filepath.Join(patternDir, problem.ID+".json")
			err = ioutil.WriteFile(problemFile, problemData, 0644)
			require.NoError(t, err)
		}
	}

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	defer func() { getConfigDir = origGetConfigDir }()
	getConfigDir = func() string {
		return tempDir
	}

	// Test ListByDifficulty
	t.Run("ListByDifficulty", func(t *testing.T) {
		difficulties, err := ListByDifficulty()
		require.NoError(t, err)

		// Check difficulty counts
		assert.Len(t, difficulties["Easy"], 1)
		assert.Len(t, difficulties["Medium"], 1)
		assert.Len(t, difficulties["Hard"], 1)

		// Check problem assignments
		assert.Equal(t, "problem1", difficulties["Easy"][0].ID)
		assert.Equal(t, "problem2", difficulties["Medium"][0].ID)
		assert.Equal(t, "problem3", difficulties["Hard"][0].ID)
	})
}

func TestListByCompany(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := ioutil.TempDir("", "algo-scales-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up test problems with different companies
	problems := []Problem{
		{
			ID:         "problem1",
			Title:      "Problem 1",
			Difficulty: "Easy",
			Patterns:   []string{"hash-map"},
			Companies:  []string{"Google", "Amazon"},
		},
		{
			ID:         "problem2",
			Title:      "Problem 2",
			Difficulty: "Medium",
			Patterns:   []string{"sliding-window"},
			Companies:  []string{"Microsoft", "Amazon"},
		},
	}

	// Create pattern directories and save problems
	for _, problem := range problems {
		for _, pattern := range problem.Patterns {
			patternDir := filepath.Join(tempDir, "problems", pattern)
			err = os.MkdirAll(patternDir, 0755)
			require.NoError(t, err)

			problemData, err := json.MarshalIndent(problem, "", "  ")
			require.NoError(t, err)
			problemFile := filepath.Join(patternDir, problem.ID+".json")
			err = ioutil.WriteFile(problemFile, problemData, 0644)
			require.NoError(t, err)
		}
	}

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	defer func() { getConfigDir = origGetConfigDir }()
	getConfigDir = func() string {
		return tempDir
	}

	// Test ListByCompany
	t.Run("ListByCompany", func(t *testing.T) {
		companies, err := ListByCompany()
		require.NoError(t, err)

		// Check company counts
		assert.Len(t, companies["Google"], 1)
		assert.Len(t, companies["Microsoft"], 1)
		assert.Len(t, companies["Amazon"], 2) // Both problems are from Amazon

		// Check problem assignment
		assert.Equal(t, "problem1", companies["Google"][0].ID)
		assert.Equal(t, "problem2", companies["Microsoft"][0].ID)
	})
}
