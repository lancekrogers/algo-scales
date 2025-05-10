// Vim mode integration for Neovim

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/spf13/cobra"
)

// Flag to enable vim mode output (JSON format)
var vimMode bool

// VimProblemResponse represents the JSON response for a problem in vim mode
type VimProblemResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Difficulty  string            `json:"difficulty"`
	Description string            `json:"description"`
	StarterCode string            `json:"starter_code"`
	Language    string            `json:"language"`
	Patterns    []string          `json:"patterns"`
	Scale       string            `json:"scale,omitempty"` // Musical scale name
	ScaleDesc   string            `json:"scale_desc,omitempty"` // Musical scale description
}

// VimTestResponse represents the JSON response for test results in vim mode
type VimTestResponse struct {
	TestResults []TestResult `json:"test_results"`
}

// TestResult represents a single test result
type TestResult struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Actual   string `json:"actual,omitempty"`
	Passed   bool   `json:"passed"`
}

// VimSubmitResponse represents the JSON response for a submission in vim mode
type VimSubmitResponse struct {
	Passed      bool         `json:"passed"`
	TestResults []TestResult `json:"test_results"`
}

// VimHintResponse represents the JSON response for a hint in vim mode
type VimHintResponse struct {
	Hint string `json:"hint"`
}

// VimSolutionResponse represents the JSON response for a solution in vim mode
type VimSolutionResponse struct {
	Solution string `json:"solution"`
}

// VimListResponse represents the JSON response for listing problems in vim mode
type VimListResponse struct {
	Problems []problem.Problem `json:"problems"`
}

// Musical scale definitions
var musicalScales = map[string]struct {
	Name        string
	Description string
}{
	"sliding-window": {
		Name:        "C Major (Sliding Window)",
		Description: "The fundamental scale, elegant and versatile",
	},
	"two-pointers": {
		Name:        "G Major (Two Pointers)",
		Description: "Balanced and efficient, the workhorse of array manipulation",
	},
	"fast-slow-pointers": {
		Name:        "D Major (Fast & Slow Pointers)",
		Description: "The cycle detector, bright and revealing",
	},
	"hash-map": {
		Name:        "A Major (Hash Maps)",
		Description: "The lookup accelerator, crisp and direct",
	},
	"binary-search": {
		Name:        "E Major (Binary Search)",
		Description: "The divider and conqueror, precise and logarithmic",
	},
	"dfs": {
		Name:        "B Major (DFS)",
		Description: "The deep explorer, rich and thorough",
	},
	"bfs": {
		Name:        "F# Major (BFS)",
		Description: "The level-by-level discoverer, methodical and complete",
	},
	"dynamic-programming": {
		Name:        "Db Major (Dynamic Programming)",
		Description: "The optimizer, complex and powerful",
	},
	"greedy": {
		Name:        "Ab Major (Greedy)",
		Description: "The local maximizer, bold and decisive",
	},
	"union-find": {
		Name:        "Eb Major (Union-Find)",
		Description: "The connector, structured and organized",
	},
	"heap": {
		Name:        "Bb Major (Heap / Priority Queue)",
		Description: "The sorter, flexible and maintaining order",
	},
}

// Initialize Vim Mode
func initVimMode() {
	// Add vim-mode flag to root command and all subcommands
	rootCmd.PersistentFlags().BoolVar(&vimMode, "vim-mode", false, "Output in Vim-friendly JSON format")
}

// handleVimModeStart handles the start command in Vim mode
func handleVimModeStart(s *session.Session) error {
	// Create the response
	resp := VimProblemResponse{
		ID:          s.Problem.ID,
		Title:       s.Problem.Title,
		Difficulty:  s.Problem.Difficulty,
		Patterns:    s.Problem.Patterns,
		Language:    s.Options.Language,
	}

	// Get the description
	resp.Description = s.FormatProblemDescription()

	// Get the starter code
	if starterCode, ok := s.Problem.StarterCode[s.Options.Language]; ok {
		resp.StarterCode = starterCode
	} else {
		// Fallback to a default language if the requested one isn't available
		for lang, code := range s.Problem.StarterCode {
			resp.StarterCode = code
			resp.Language = lang
			break
		}
	}

	// Add musical scale information if available
	if len(s.Problem.Patterns) > 0 {
		pattern := s.Problem.Patterns[0]
		if scale, ok := musicalScales[pattern]; ok {
			resp.Scale = scale.Name
			resp.ScaleDesc = scale.Description
		}
	}
