// Vim mode extensions for start command

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/services"
	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/spf13/cobra"
)

// handleVimModeSession handles starting a session in vim mode
func handleVimModeSession(opts session.Options) {
	// Get problem service
	problemService := services.DefaultRegistry.GetProblemService()

	var prob *problem.Problem
	var err error

	ctx := context.Background()
	if opts.ProblemID != "" {
		// Get specific problem
		prob, err = problemService.GetByID(ctx, opts.ProblemID)
		if err != nil {
			outputVimError(fmt.Errorf("failed to get problem: %v", err))
			return
		}
	} else {
		// Get a random problem based on filters
		problems, err := problemService.ListAll(ctx)
		if err != nil {
			outputVimError(fmt.Errorf("failed to get problems: %v", err))
			return
		}

		// Filter by pattern if specified
		if opts.Pattern != "" {
			var filtered []problem.Problem
			for _, p := range problems {
				for _, pat := range p.Patterns {
					if pat == opts.Pattern {
						filtered = append(filtered, p)
						break
					}
				}
			}
			problems = filtered
		}

		// Filter by difficulty if specified
		if opts.Difficulty != "" {
			var filtered []problem.Problem
			for _, p := range problems {
				if p.Difficulty == opts.Difficulty {
					filtered = append(filtered, p)
				}
			}
			problems = filtered
		}

		if len(problems) == 0 {
			outputVimError(fmt.Errorf("no problems found matching criteria"))
			return
		}

		// Select first problem (could be randomized)
		prob = &problems[0]
	}

	// Create a session and workspace for this problem
	sess, err := session.CreateSession(opts)
	if err != nil {
		outputVimError(fmt.Errorf("failed to create session: %v", err))
		return
	}

	// Create the response with workspace information
	resp := VimProblemResponse{
		ID:          prob.ID,
		Title:       prob.Title,
		Difficulty:  prob.Difficulty,
		Patterns:    prob.Patterns,
		Language:    opts.Language,
		Description: formatProblemDescriptionForVim(prob),
		WorkspacePath: sess.Workspace,
		SessionID:   sess.Problem.ID, // Use problem ID as session identifier
	}

	// Get the starter code
	if starterCode, ok := prob.StarterCode[opts.Language]; ok {
		resp.StarterCode = starterCode
	} else {
		// Fallback to a default language if the requested one isn't available
		for lang, code := range prob.StarterCode {
			resp.StarterCode = code
			resp.Language = lang
			break
		}
	}

	// Add musical scale information if available
	if len(prob.Patterns) > 0 {
		pattern := prob.Patterns[0]
		if scale, ok := musicalScales[pattern]; ok {
			resp.Scale = scale.Name
			resp.ScaleDesc = scale.Description
		}
	}

	// Return JSON response
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		outputVimError(fmt.Errorf("error marshaling JSON: %v", err))
		return
	}

	fmt.Println(string(jsonResp))
}

// formatProblemDescriptionForVim formats the problem description for vim mode
func formatProblemDescriptionForVim(p *problem.Problem) string {
	desc := p.Title + "\n\n"
	desc += p.Description + "\n\n"

	if len(p.Examples) > 0 {
		desc += "Examples:\n"
		for i, ex := range p.Examples {
			desc += fmt.Sprintf("\nExample %d:\n", i+1)
			desc += fmt.Sprintf("Input: %s\n", ex.Input)
			desc += fmt.Sprintf("Output: %s\n", ex.Output)
			if ex.Explanation != "" {
				desc += fmt.Sprintf("Explanation: %s\n", ex.Explanation)
			}
		}
	}

	if len(p.Constraints) > 0 {
		desc += "\nConstraints:\n"
		for _, c := range p.Constraints {
			desc += "- " + c + "\n"
		}
	}

	return desc
}

// Override the start command handlers in init
func init() {
	// We'll modify the existing commands by updating their Run functions
	// This is done in a post-init hook
	oldLearnRun := learnCmd.Run
	learnCmd.Run = func(cmd *cobra.Command, args []string) {
		isVimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
		if isVimMode {
			var problemID string
			if len(args) > 0 {
				problemID = args[0]
			}

			opts := session.Options{
				Mode:       session.LearnMode,
				Language:   language,
				Timer:      timer,
				Pattern:    pattern,
				Difficulty: difficulty,
				ProblemID:  problemID,
			}
			handleVimModeSession(opts)
			return
		}
		// Call original handler
		oldLearnRun(cmd, args)
	}

	oldPracticeRun := practiceCmd.Run
	practiceCmd.Run = func(cmd *cobra.Command, args []string) {
		isVimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
		if isVimMode {
			var problemID string
			if len(args) > 0 {
				problemID = args[0]
			}

			opts := session.Options{
				Mode:       session.PracticeMode,
				Language:   language,
				Timer:      timer,
				Pattern:    pattern,
				Difficulty: difficulty,
				ProblemID:  problemID,
			}
			handleVimModeSession(opts)
			return
		}
		// Call original handler
		oldPracticeRun(cmd, args)
	}

	oldCramRun := cramCmd.Run
	cramCmd.Run = func(cmd *cobra.Command, args []string) {
		isVimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
		if isVimMode {
			opts := session.Options{
				Mode:       session.CramMode,
				Language:   language,
				Timer:      timer,
				Pattern:    pattern,
				Difficulty: difficulty,
			}
			handleVimModeSession(opts)
			return
		}
		// Call original handler
		oldCramRun(cmd, args)
	}
}