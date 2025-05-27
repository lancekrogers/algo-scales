// Vim mode extensions for list command

package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/services"
	"github.com/spf13/cobra"
)

// Override the list command handlers in init
func init() {
	// Override main list command
	oldListRun := listCmd.Run
	listCmd.Run = func(cmd *cobra.Command, args []string) {
		isVimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
		if isVimMode {
			handleVimModeList(cmd)
			return
		}
		// Call original handler
		oldListRun(cmd, args)
	}

	// Override patterns subcommand
	oldPatternsRun := patternsCmd.Run
	patternsCmd.Run = func(cmd *cobra.Command, args []string) {
		isVimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
		if isVimMode {
			handleVimModeList(cmd)
			return
		}
		// Call original handler
		oldPatternsRun(cmd, args)
	}

	// Override difficulties subcommand
	oldDifficultiesRun := difficultiesCmd.Run
	difficultiesCmd.Run = func(cmd *cobra.Command, args []string) {
		isVimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
		if isVimMode {
			handleVimModeList(cmd)
			return
		}
		// Call original handler
		oldDifficultiesRun(cmd, args)
	}

	// Override companies subcommand
	oldCompaniesRun := companiesCmd.Run
	companiesCmd.Run = func(cmd *cobra.Command, args []string) {
		isVimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
		if isVimMode {
			handleVimModeList(cmd)
			return
		}
		// Call original handler
		oldCompaniesRun(cmd, args)
	}
}

// handleVimModeList handles the list command in vim mode
func handleVimModeList(cmd *cobra.Command) {
	// Get problem service
	problemService := services.DefaultRegistry.GetProblemService()

	// Get all problems
	ctx := context.Background()
	problems, err := problemService.ListAll(ctx)
	if err != nil {
		// Try fallback to direct problem listing
		problems, err = problem.ListAll()
		if err != nil {
			outputVimError(fmt.Errorf("failed to get problems: %v", err))
			return
		}
	}

	// Create and output response
	resp := VimListResponse{
		Problems: problems,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		outputVimError(fmt.Errorf("failed to marshal response: %v", err))
		return
	}

	fmt.Println(string(jsonResp))
}