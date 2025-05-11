// Problem model and management
package problem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Problem represents an algorithm problem
type Problem struct {
	ID                  string            `json:"id"`
	Title               string            `json:"title"`
	Difficulty          string            `json:"difficulty"`
	Patterns            []string          `json:"patterns"`
	EstimatedTime       int               `json:"estimated_time"` // in minutes
	Companies           []string          `json:"companies"`
	Description         string            `json:"description"`
	Examples            []Example         `json:"examples"`
	Constraints         []string          `json:"constraints"`
	PatternExplanation  string            `json:"pattern_explanation"`
	SolutionWalkthrough []string          `json:"solution_walkthrough"`
	StarterCode         map[string]string `json:"starter_code"`
	Solutions           map[string]string `json:"solutions"`
	TestCases           []TestCase        `json:"test_cases"`
}

// Example represents an example for a problem
type Example struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation,omitempty"`
}

// TestCase represents a test case for a problem
type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

// GetByID retrieves a problem by its ID
// Exported as variable for testing
var GetByID = func(id string) (*Problem, error) {
	configDir := getConfigDir()

	// Search in all pattern directories
	patternDirs, err := ioutil.ReadDir(filepath.Join(configDir, "problems"))
	if err != nil {
		return nil, err
	}

	for _, patternDir := range patternDirs {
		if !patternDir.IsDir() {
			continue
		}

		problemPath := filepath.Join(configDir, "problems", patternDir.Name(), fmt.Sprintf("%s.json", id))
		if _, err := os.Stat(problemPath); os.IsNotExist(err) {
			continue
		}

		// Found the problem file
		data, err := ioutil.ReadFile(problemPath)
		if err != nil {
			return nil, err
		}

		var problem Problem
		if err := json.Unmarshal(data, &problem); err != nil {
			return nil, err
		}

		return &problem, nil
	}

	return nil, fmt.Errorf("problem not found: %s", id)
}

// ListAll lists all available problems
// Exported as variable for testing
var ListAll = func() ([]Problem, error) {
	var problems []Problem
	configDir := getConfigDir()

	// Get all pattern directories
	patternDirs, err := ioutil.ReadDir(filepath.Join(configDir, "problems"))
	if err != nil {
		return nil, err
	}

	// Track processed problem IDs to avoid duplicates
	processedIDs := make(map[string]bool)

	for _, patternDir := range patternDirs {
		if !patternDir.IsDir() {
			continue
		}

		// Read all problem files in this pattern directory
		problemFiles, err := ioutil.ReadDir(filepath.Join(configDir, "problems", patternDir.Name()))
		if err != nil {
			return nil, err
		}

		for _, problemFile := range problemFiles {
			if problemFile.IsDir() || !strings.HasSuffix(problemFile.Name(), ".json") {
				continue
			}

			// Read problem file
			data, err := ioutil.ReadFile(filepath.Join(configDir, "problems", patternDir.Name(), problemFile.Name()))
			if err != nil {
				return nil, err
			}

			var problem Problem
			if err := json.Unmarshal(data, &problem); err != nil {
				return nil, err
			}

			// Skip if already processed
			if processedIDs[problem.ID] {
				continue
			}

			problems = append(problems, problem)
			processedIDs[problem.ID] = true
		}
	}

	return problems, nil
}

// ListPatterns lists problems organized by pattern
// Exported as variable for testing
var ListPatterns = func() (map[string][]Problem, error) {
	patterns := make(map[string][]Problem)

	// Get all problems
	problems, err := ListAll()
	if err != nil {
		return nil, err
	}

	// Organize by pattern
	for _, problem := range problems {
		for _, pattern := range problem.Patterns {
			patterns[pattern] = append(patterns[pattern], problem)
		}
	}

	return patterns, nil
}

// ListByDifficulty lists problems organized by difficulty
// Exported as variable for testing
var ListByDifficulty = func() (map[string][]Problem, error) {
	difficulties := make(map[string][]Problem)

	// Get all problems
	problems, err := ListAll()
	if err != nil {
		return nil, err
	}

	// Organize by difficulty
	for _, problem := range problems {
		difficulties[problem.Difficulty] = append(difficulties[problem.Difficulty], problem)
	}

	return difficulties, nil
}

// ListByCompany lists problems organized by company
// Exported as variable for testing
var ListByCompany = func() (map[string][]Problem, error) {
	companies := make(map[string][]Problem)

	// Get all problems
	problems, err := ListAll()
	if err != nil {
		return nil, err
	}

	// Organize by company
	for _, problem := range problems {
		for _, company := range problem.Companies {
			companies[company] = append(companies[company], problem)
		}
	}

	return companies, nil
}

// GetConfigDir returns the configuration directory
// Exported as variable for testing
var getConfigDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}
