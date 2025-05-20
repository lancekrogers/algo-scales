// Package daily provides functionality for handling daily scale practice
package daily

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

// ProblemState represents the current state of a problem in daily practice
type ProblemState string

const (
	// StatePending means the problem hasn't been started yet
	StatePending ProblemState = "pending"
	
	// StateInProgress means the problem has been started but not completed
	StateInProgress ProblemState = "in_progress"
	
	// StateCompleted means the problem has been completed successfully
	StateCompleted ProblemState = "completed"
	
	// StateSkipped means the problem was skipped by the user
	StateSkipped ProblemState = "skipped"
)

// DailyProblem represents a problem in the daily practice session
type DailyProblem struct {
	Pattern    string       `json:"pattern"`
	ProblemID  string       `json:"problem_id"`
	State      ProblemState `json:"state"`
	StartedAt  time.Time    `json:"started_at"`
	CompletedAt time.Time   `json:"completed_at,omitempty"`
	Attempts   int          `json:"attempts"`
}

// GetDailyWorkspacePath returns the path to the daily workspace directory
func GetDailyWorkspacePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temporary directory
		return filepath.Join(os.TempDir(), "AlgoScalesPractice", "Daily")
	}
	
	// Use the requested path in home directory
	return filepath.Join(homeDir, "Dev", "AlgoScalesPractice", "Daily")
}

// GetTodayWorkspacePath returns the path for today's practice directory
func GetTodayWorkspacePath() string {
	todayStr := time.Now().Format("2006-01-02")
	return filepath.Join(GetDailyWorkspacePath(), todayStr)
}

// CreateDailyWorkspace creates the daily practice workspace directory
func CreateDailyWorkspace() error {
	path := GetTodayWorkspacePath()
	return os.MkdirAll(path, 0755)
}

// FormatProblemAsComment formats a problem description as source code comments
// for the given programming language
func FormatProblemAsComment(prob *problem.Problem, language string) string {
	// Determine comment style based on language
	var lineComment string
	var blockStart string
	var blockEnd string
	
	switch language {
	case "python":
		lineComment = "# "
		blockStart = "'''\n"
		blockEnd = "'''\n"
	case "javascript":
		lineComment = "// "
		blockStart = "/**\n"
		blockEnd = " */\n"
	case "go":
		lineComment = "// "
		blockStart = "/*\n"
		blockEnd = " */\n"
	default:
		// Default to C-style comments
		lineComment = "// "
		blockStart = "/*\n"
		blockEnd = " */\n"
	}
	
	var builder strings.Builder
	
	// Use block comment for header and description
	builder.WriteString(blockStart)
	builder.WriteString(fmt.Sprintf(" # %s\n", prob.Title))
	builder.WriteString(fmt.Sprintf(" **Difficulty**: %s\n", prob.Difficulty))
	builder.WriteString(fmt.Sprintf(" **Pattern**: %s\n", strings.Join(prob.Patterns, ", ")))
	builder.WriteString(fmt.Sprintf(" **Estimated Time**: %d minutes\n", prob.EstimatedTime))
	builder.WriteString("\n")
	
	// Problem statement
	builder.WriteString(" ## Problem Statement\n\n")
	for _, line := range strings.Split(prob.Description, "\n") {
		builder.WriteString(fmt.Sprintf(" %s\n", line))
	}
	builder.WriteString("\n")
	
	// Examples
	builder.WriteString(" ## Examples\n\n")
	for i, example := range prob.Examples {
		builder.WriteString(fmt.Sprintf(" ### Example %d\n\n", i+1))
		builder.WriteString(fmt.Sprintf(" **Input**: %s\n", example.Input))
		builder.WriteString(fmt.Sprintf(" **Output**: %s\n", example.Output))
		if example.Explanation != "" {
			builder.WriteString(fmt.Sprintf(" **Explanation**: %s\n", example.Explanation))
		}
		builder.WriteString("\n")
	}
	
	// Constraints
	builder.WriteString(" ## Constraints\n\n")
	for _, constraint := range prob.Constraints {
		builder.WriteString(fmt.Sprintf(" - %s\n", constraint))
	}
	builder.WriteString("\n")
	
	// Test cases
	builder.WriteString(" ## Test Cases\n\n")
	for i, test := range prob.TestCases {
		builder.WriteString(fmt.Sprintf(" Test %d:\n", i+1))
		builder.WriteString(fmt.Sprintf(" - Input: %s\n", test.Input))
		builder.WriteString(fmt.Sprintf(" - Expected: %s\n", test.Expected))
		builder.WriteString("\n")
	}
	
	// Add pattern information if available
	if prob.PatternExplanation != "" {
		builder.WriteString(" ## Pattern Description\n\n")
		for _, line := range strings.Split(prob.PatternExplanation, "\n") {
			builder.WriteString(fmt.Sprintf(" %s\n", line))
		}
		builder.WriteString("\n")
	}
	
	builder.WriteString(blockEnd)
	builder.WriteString("\n")
	
	// Add starter code
	starterCode, ok := prob.StarterCode[language]
	if !ok {
		// Fallback to any available language
		for _, code := range prob.StarterCode {
			starterCode = code
			break
		}
	}
	
	builder.WriteString(starterCode)
	builder.WriteString("\n\n")
	
	// Add test marker section
	builder.WriteString(lineComment + "Do not modify below this line\n")
	builder.WriteString(lineComment + "AlgoScales: Test Section\n")
	
	return builder.String()
}

// CreateProblemFile creates a file for the problem in the daily workspace
func CreateProblemFile(prob *problem.Problem, language string) (string, error) {
	// Ensure workspace exists
	if err := CreateDailyWorkspace(); err != nil {
		return "", fmt.Errorf("failed to create workspace: %w", err)
	}
	
	// Get file extension for the language
	ext := GetFileExtension(language)
	
	// Create the file path
	filePath := filepath.Join(GetTodayWorkspacePath(), fmt.Sprintf("%s.%s", prob.ID, ext))
	
	// Format the problem as comments
	content := FormatProblemAsComment(prob, language)
	
	// Write to file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write problem file: %w", err)
	}
	
	return filePath, nil
}

// GetFileExtension returns the file extension for a programming language
func GetFileExtension(language string) string {
	switch language {
	case "go":
		return "go"
	case "python":
		return "py"
	case "javascript":
		return "js"
	default:
		return "txt"
	}
}

// GetProblemFilePath returns the path to the problem file for a specific problem
func GetProblemFilePath(problemID, language string) string {
	ext := GetFileExtension(language)
	return filepath.Join(GetTodayWorkspacePath(), fmt.Sprintf("%s.%s", problemID, ext))
}

// ProblemFileExists checks if a problem file exists
func ProblemFileExists(problemID, language string) bool {
	path := GetProblemFilePath(problemID, language)
	_, err := os.Stat(path)
	return err == nil
}