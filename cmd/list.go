// List command for displaying problems

package cmd

import (
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available problems",
	Long:  `List the available algorithm problems by various criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no subcommand is specified
		problems, err := problem.ListAll()
		if err != nil {
			fmt.Printf("Error listing problems: %v\n", err)
			return
		}

		fmt.Println("Available Problems:")
		for _, p := range problems {
			fmt.Printf("- %s (%s): %s\n", p.ID, p.Difficulty, p.Title)
		}
	},
}

// patternsCmd represents the patterns subcommand
var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "List problems by pattern",
	Long:  `List the available algorithm problems organized by pattern.`,
	Run: func(cmd *cobra.Command, args []string) {
		patterns, err := problem.ListPatterns()
		if err != nil {
			fmt.Printf("Error listing patterns: %v\n", err)
			return
		}

		fmt.Println("Algorithm Patterns:")
		for pattern, problems := range patterns {
			fmt.Printf("\n%s:\n", pattern)
			for _, p := range problems {
				fmt.Printf("  - %s (%s): %s\n", p.ID, p.Difficulty, p.Title)
			}
		}
	},
}

// difficultiesCmd represents the difficulties subcommand
var difficultiesCmd = &cobra.Command{
	Use:   "difficulties",
	Short: "List problems by difficulty",
	Long:  `List the available algorithm problems organized by difficulty.`,
	Run: func(cmd *cobra.Command, args []string) {
		difficulties, err := problem.ListByDifficulty()
		if err != nil {
			fmt.Printf("Error listing by difficulty: %v\n", err)
			return
		}

		fmt.Println("Problems by Difficulty:")
		for difficulty, problems := range difficulties {
			fmt.Printf("\n%s:\n", difficulty)
			for _, p := range problems {
				fmt.Printf("  - %s: %s\n", p.ID, p.Title)
			}
		}
	},
}

// companiesCmd represents the companies subcommand
var companiesCmd = &cobra.Command{
	Use:   "companies",
	Short: "List problems by company",
	Long:  `List the available algorithm problems organized by the companies that commonly ask them.`,
	Run: func(cmd *cobra.Command, args []string) {
		companies, err := problem.ListByCompany()
		if err != nil {
			fmt.Printf("Error listing by company: %v\n", err)
			return
		}

		fmt.Println("Problems by Company:")
		for company, problems := range companies {
			fmt.Printf("\n%s:\n", company)
			for _, p := range problems {
				fmt.Printf("  - %s (%s): %s\n", p.ID, p.Difficulty, p.Title)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(patternsCmd)
	listCmd.AddCommand(difficultiesCmd)
	listCmd.AddCommand(companiesCmd)
}
