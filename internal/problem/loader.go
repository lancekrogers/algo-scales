// Package problem handles algorithm problems
package problem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LoadLocalProblems loads problems from the local problems directory
func LoadLocalProblems() ([]Problem, error) {
	return LoadLocalProblemsWithContext(context.Background())
}

// LoadLocalProblemsWithContext loads problems from the local problems directory with context
func LoadLocalProblemsWithContext(ctx context.Context) ([]Problem, error) {
	// Create a context with timeout for file operations
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// First try the standard config dir location
	configDir := getConfigDir()
	problemsDir := filepath.Join(configDir, "problems")
	
	// If problems directory doesn't exist in config dir,
	// try the local problems directory relative to binary
	if _, err := os.Stat(problemsDir); os.IsNotExist(err) {
		// Get the executable directory
		exePath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %v", err)
		}
		
		exeDir := filepath.Dir(exePath)
		problemsDir = filepath.Join(exeDir, "problems")
		
		// If still no problems directory, try current directory
		if _, err := os.Stat(problemsDir); os.IsNotExist(err) {
			curDir, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get current directory: %v", err)
			}
			
			problemsDir = filepath.Join(curDir, "problems")
			
			// If still no problems directory, try project root
			if _, err := os.Stat(problemsDir); os.IsNotExist(err) {
				// Try going up one directory (assuming we're in algo-scales/...)
				rootDir := filepath.Dir(curDir)
				problemsDir = filepath.Join(rootDir, "algo-scales", "problems")
				
				// If still no problems directory, return empty result
				if _, err := os.Stat(problemsDir); os.IsNotExist(err) {
					return []Problem{}, nil
				}
			}
		}
	}
	
	// Track processed problem IDs to avoid duplicates
	var problems []Problem
	processedIDs := make(map[string]bool)
	
	// Get pattern directories
	patternDirs, err := os.ReadDir(problemsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read problems directory: %v", err)
	}
	
	// Iterate through pattern directories
	for _, patternDir := range patternDirs {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("problem loading cancelled: %w", ctx.Err())
		default:
		}
		
		if !patternDir.IsDir() {
			continue
		}
		
		patternName := patternDir.Name()
		patternPath := filepath.Join(problemsDir, patternName)
		
		// Read problem files in the pattern directory
		problemFiles, err := os.ReadDir(patternPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read pattern directory %s: %v", patternName, err)
		}
		
		for _, file := range problemFiles {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}
			
			// Read problem file
			problemPath := filepath.Join(patternPath, file.Name())
			data, err := os.ReadFile(problemPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read problem file %s: %v", problemPath, err)
			}
			
			// Parse problem data
			var problem Problem
			err = json.Unmarshal(data, &problem)
			if err != nil {
				return nil, fmt.Errorf("failed to parse problem file %s: %v", problemPath, err)
			}
			
			// Skip if already processed
			if processedIDs[problem.ID] {
				continue
			}
			
			// Add problem to the list
			problems = append(problems, problem)
			processedIDs[problem.ID] = true
		}
	}
	
	// Sort problems by difficulty (easy, medium, hard)
	sort.Slice(problems, func(i, j int) bool {
		// Define difficulty order
		difficultyOrder := map[string]int{
			"easy":   0,
			"medium": 1,
			"hard":   2,
		}
		
		// Get difficulty values
		diffI := difficultyOrder[problems[i].Difficulty]
		diffJ := difficultyOrder[problems[j].Difficulty]
		
		// Sort by difficulty first
		if diffI != diffJ {
			return diffI < diffJ
		}
		
		// Then by ID for consistent ordering
		return problems[i].ID < problems[j].ID
	})
	
	return problems, nil
}

// GetProblemsByPattern returns all problems matching the given pattern
func GetProblemsByPattern(allProblems []Problem, pattern string) []Problem {
	if pattern == "" {
		return allProblems
	}
	
	var filtered []Problem
	for _, p := range allProblems {
		for _, patternName := range p.Patterns {
			if patternName == pattern {
				filtered = append(filtered, p)
				break
			}
		}
	}
	
	return filtered
}

// GetPatterns returns all available algorithm patterns
func GetPatterns(allProblems []Problem) []string {
	// Use map to track unique patterns
	patterns := make(map[string]bool)
	
	for _, problem := range allProblems {
		for _, pattern := range problem.Patterns {
			// Convert kebab-case to Title Case for display
			displayPattern := convertPatternToDisplay(pattern)
			patterns[displayPattern] = true
		}
	}
	
	// Convert map to sorted slice
	result := make([]string, 0, len(patterns))
	for pattern := range patterns {
		result = append(result, pattern)
	}
	
	// Sort patterns for consistent ordering
	sort.Strings(result)
	
	return result
}

// convertPatternToDisplay converts kebab-case pattern names to Title Case
func convertPatternToDisplay(pattern string) string {
	// Handle special cases
	patternMap := map[string]string{
		"two-pointers":       "Two Pointers",
		"sliding-window":     "Sliding Window",
		"fast-slow-pointers": "Fast & Slow Pointers",
		"hash-map":          "Hash Map",
		"binary-search":     "Binary Search",
		"bfs":               "BFS",
		"dfs":               "DFS",
		"dynamic-programming": "Dynamic Programming",
		"greedy":            "Greedy",
		"heap":              "Heap",
		"union-find":        "Union Find",
	}
	
	if display, ok := patternMap[pattern]; ok {
		return display
	}
	
	// Fallback: convert kebab-case to Title Case
	words := strings.Split(pattern, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// GetLanguages returns all available programming languages from problems
func GetLanguages(allProblems []Problem) []string {
	// Use map to track unique languages
	languages := make(map[string]bool)
	
	for _, problem := range allProblems {
		for lang := range problem.StarterCode {
			languages[lang] = true
		}
	}
	
	// Convert map to sorted slice
	result := make([]string, 0, len(languages))
	for lang := range languages {
		result = append(result, lang)
	}
	
	// Sort languages for consistent ordering
	sort.Strings(result)
	
	return result
}