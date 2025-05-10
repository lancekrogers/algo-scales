// Session management core

package session

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/stats"
	"github.com/lancekrogers/algo-scales/internal/ui"
)

// Mode represents a session mode
type Mode string

const (
	LearnMode    Mode = "learn"
	PracticeMode Mode = "practice"
	CramMode     Mode = "cram"
)

// Options represents options for a session
type Options struct {
	Mode       Mode
	Language   string
	Timer      int
	Pattern    string
	Difficulty string
	ProblemID  string
}

// Session represents a practice session
type Session struct {
	Options      Options
	Problem      *problem.Problem
	StartTime    time.Time
	EndTime      time.Time
	Workspace    string
	CodeFile     string
	ShowHints    bool
	ShowPattern  bool
	ShowSolution bool
}

// Start begins a new practice session
func Start(opts Options) error {
	// Initialize session
	session := &Session{
		Options:      opts,
		StartTime:    time.Now(),
		ShowHints:    opts.Mode == LearnMode,
		ShowPattern:  opts.Mode == LearnMode,
		ShowSolution: false,
	}

	// Choose problem based on options
	var err error
	if opts.ProblemID != "" {
		// Specific problem requested
		session.Problem, err = problem.GetByID(opts.ProblemID)
		if err != nil {
			return fmt.Errorf("failed to load problem: %v", err)
		}
	} else if opts.Mode == CramMode {
		// Cram mode - choose problems from common patterns
		session.Problem, err = selectCramProblem()
		if err != nil {
			return fmt.Errorf("failed to select problem for cram mode: %v", err)
		}
	} else {
		// Filter by pattern/difficulty if specified
		session.Problem, err = selectProblem(opts.Pattern, opts.Difficulty)
		if err != nil {
			return fmt.Errorf("failed to select problem: %v", err)
		}
	}

	// Create workspace
	if err := session.createWorkspace(); err != nil {
		return fmt.Errorf("failed to create workspace: %v", err)
	}

	// Start UI
	return ui.StartSession(session)
}

// selectProblem chooses a problem based on filters
func selectProblem(pattern, difficulty string) (*problem.Problem, error) {
	// Get all problems
	problems, err := problem.ListAll()
	if err != nil {
		return nil, err
	}

	// Filter problems
	var filtered []problem.Problem
	for _, p := range problems {
		matchesPattern := pattern == "" || containsPattern(p.Patterns, pattern)
		matchesDifficulty := difficulty == "" || p.Difficulty == difficulty

		if matchesPattern && matchesDifficulty {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no problems match the specified filters")
	}

	// Choose a random problem from filtered list
	rand.Seed(time.Now().UnixNano())
	selected := filtered[rand.Intn(len(filtered))]

	return &selected, nil
}

// selectCramProblem chooses a problem for cram mode
func selectCramProblem() (*problem.Problem, error) {
	// For cram mode, we focus on the most common patterns
	commonPatterns := []string{
		"sliding-window",
		"two-pointers",
		"fast-slow-pointers",
		"hash-map",
		"binary-search",
		"dfs",
		"bfs",
		"dynamic-programming",
		"greedy",
		"union-find",
		"heap",
	}

	// Choose a random pattern
	rand.Seed(time.Now().UnixNano())
	pattern := commonPatterns[rand.Intn(len(commonPatterns))]

	// Get a problem with this pattern
	return selectProblem(pattern, "")
}

// createWorkspace sets up a workspace for the problem
func (s *Session) createWorkspace() error {
	// Create workspace directory
	workspaceDir := filepath.Join(os.TempDir(), "algo-scales", s.Problem.ID)
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return err
	}

	s.Workspace = workspaceDir

	// Create problem description file
	descriptionFile := filepath.Join(workspaceDir, "problem.md")
	description := s.FormatProblemDescription()
	if err := ioutil.WriteFile(descriptionFile, []byte(description), 0644); err != nil {
		return err
	}

	// Create code file with starter code
	ext := languageExtension(s.Options.Language)
	codeFile := filepath.Join(workspaceDir, fmt.Sprintf("solution.%s", ext))

	starterCode, ok := s.Problem.StarterCode[s.Options.Language]
	if !ok {
		// Fallback to a default language if the requested one isn't available
		for lang, code := range s.Problem.StarterCode {
			starterCode = code
			s.Options.Language = lang
			break
		}
	}

	if err := ioutil.WriteFile(codeFile, []byte(starterCode), 0644); err != nil {
		return err
	}

	s.CodeFile = codeFile

	return nil
}

// FormatProblemDescription creates a formatted markdown description
func (s *Session) FormatProblemDescription() string {
	var description string

	// Problem header
	description += fmt.Sprintf("# %s\n\n", s.Problem.Title)
	description += fmt.Sprintf("**Difficulty**: %s\n", s.Problem.Difficulty)
	description += fmt.Sprintf("**Estimated Time**: %d minutes\n", s.Problem.EstimatedTime)
	description += fmt.Sprintf("**Companies**: %s\n\n", joinStrings(s.Problem.Companies))

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

	// Pattern explanation (if in Learn mode)
	if s.ShowPattern {
		description += fmt.Sprintf("## Pattern: %s\n\n", joinStrings(s.Problem.Patterns))
		description += fmt.Sprintf("%s\n\n", s.Problem.PatternExplanation)
	}

	// Solution walkthrough (if requested)
	if s.ShowSolution {
		description += "## Solution Walkthrough\n\n"
		for i, step := range s.Problem.SolutionWalkthrough {
			description += fmt.Sprintf("%d. %s\n", i+1, step)
		}
		description += "\n"
	}

	return description
}

// FinishSession completes a session and records stats
func (s *Session) FinishSession(solved bool) error {
	s.EndTime = time.Now()

	// Record stats
	sessionStats := stats.SessionStats{
		ProblemID:    s.Problem.ID,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		Duration:     s.EndTime.Sub(s.StartTime),
		Solved:       solved,
		Mode:         string(s.Options.Mode),
		HintsUsed:    s.ShowHints,
		SolutionUsed: s.ShowSolution,
		Patterns:     s.Problem.Patterns,
		Difficulty:   s.Problem.Difficulty,
	}

	return stats.RecordSession(sessionStats)
}

// Helper functions

// containsPattern checks if a pattern is in a list
func containsPattern(patterns []string, pattern string) bool {
	for _, p := range patterns {
		if p == pattern {
			return true
		}
	}
	return false
}

// joinStrings joins a string slice with commas
func joinStrings(strings []string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += ", " + strings[i]
	}

	return result
}

// languageExtension returns the file extension for a language
func languageExtension(language string) string {
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
