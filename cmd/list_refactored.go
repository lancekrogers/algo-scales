// Refactored List command using service layer
package cmd

import (
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/services"
	"github.com/spf13/cobra"
)

// RefactoredListCmd demonstrates separation of command handling from business logic
var RefactoredListCmd = &cobra.Command{
	Use:   "list-refactored",
	Short: "List available problems (refactored version)",
	Long:  `List the available algorithm problems by various criteria using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		problemService := services.DefaultRegistry.GetProblemService()
		
		problems, err := problemService.ListAll()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error listing problems: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Available Problems:")
		for _, p := range problems {
			fmt.Fprintf(cmd.OutOrStdout(), "- %s (%s): %s\n", p.ID, p.Difficulty, p.Title)
		}
	},
}

// RefactoredPatternsCmd demonstrates pattern listing using service layer
var RefactoredPatternsCmd = &cobra.Command{
	Use:   "patterns-refactored",
	Short: "List problems by pattern (refactored version)",
	Long:  `List the available algorithm problems organized by pattern using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		problemService := services.DefaultRegistry.GetProblemService()
		
		patterns, err := problemService.ListByPattern()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error listing patterns: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Algorithm Patterns:")
		for pattern, problems := range patterns {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n", pattern)
			for _, p := range problems {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s): %s\n", p.ID, p.Difficulty, p.Title)
			}
		}
	},
}

// RefactoredDifficultyCmd demonstrates difficulty listing using service layer
var RefactoredDifficultyCmd = &cobra.Command{
	Use:   "difficulty-refactored",
	Short: "List problems by difficulty (refactored version)",
	Long:  `List the available algorithm problems organized by difficulty using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		problemService := services.DefaultRegistry.GetProblemService()
		
		difficulties, err := problemService.ListByDifficulty()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error listing by difficulty: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Problems by Difficulty:")
		for difficulty, problems := range difficulties {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n", difficulty)
			for _, p := range problems {
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s: %s\n", p.ID, p.Title)
			}
		}
	},
}

func init() {
	// Add refactored subcommands to demonstrate service layer usage
	RefactoredListCmd.AddCommand(RefactoredPatternsCmd)
	RefactoredListCmd.AddCommand(RefactoredDifficultyCmd)
}