// Refactored Stats command using service layer
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/lancekrogers/algo-scales/internal/services"
	"github.com/spf13/cobra"
)

// RefactoredStatsCmd demonstrates separation of command handling from business logic
var RefactoredStatsCmd = &cobra.Command{
	Use:   "stats-refactored",
	Short: "Show performance statistics (refactored version)",
	Long:  `Display your algorithm practice performance statistics using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		statsService := services.DefaultRegistry.GetStatsService()
		
		// Get overall stats
		ctx := context.Background()
		overall, err := statsService.GetOverallStats(ctx)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting overall stats: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "=== Overall Statistics ===")
		if overall.Summary != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Total Attempted: %d\n", overall.Summary.TotalAttempted)
			fmt.Fprintf(cmd.OutOrStdout(), "Total Solved: %d\n", overall.Summary.TotalSolved)
			fmt.Fprintf(cmd.OutOrStdout(), "Success Rate: %.1f%%\n", overall.Summary.SuccessRate*100)
			fmt.Fprintf(cmd.OutOrStdout(), "Average Solve Time: %s\n", overall.Summary.AvgSolveTime)
			
			if overall.Summary.FastestSolve.ProblemID != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Fastest Solve: %s (%s)\n", overall.Summary.FastestSolve.ProblemID, overall.Summary.FastestSolve.Time)
			}
			
			if overall.Summary.MostChallenging.ProblemID != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Most Challenging: %s (%d attempts)\n", overall.Summary.MostChallenging.ProblemID, overall.Summary.MostChallenging.Attempts)
			}
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "No summary data available")
		}
	},
}

// RefactoredPatternStatsCmd demonstrates pattern stats using service layer
var RefactoredPatternStatsCmd = &cobra.Command{
	Use:   "patterns-refactored",
	Short: "Show pattern-specific statistics (refactored version)",
	Long:  `Display performance statistics by algorithm pattern using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		statsService := services.DefaultRegistry.GetStatsService()
		
		ctx := context.Background()
		patternStats, err := statsService.GetPatternStats(ctx)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting pattern stats: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "=== Pattern Statistics ===")
		for pattern, stats := range patternStats {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n", pattern)
			fmt.Fprintf(cmd.OutOrStdout(), "  Attempted: %d\n", stats.Attempted)
			fmt.Fprintf(cmd.OutOrStdout(), "  Solved: %d\n", stats.Solved)
			fmt.Fprintf(cmd.OutOrStdout(), "  Success Rate: %.1f%%\n", stats.SuccessRate*100)
			fmt.Fprintf(cmd.OutOrStdout(), "  Average Time: %s\n", stats.AvgTime)
		}
	},
}

// RefactoredDifficultyStatsCmd demonstrates difficulty stats using service layer
var RefactoredDifficultyStatsCmd = &cobra.Command{
	Use:   "difficulty-refactored",
	Short: "Show difficulty-specific statistics (refactored version)",
	Long:  `Display performance statistics by difficulty level using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		statsService := services.DefaultRegistry.GetStatsService()
		
		ctx := context.Background()
		difficultyStats, err := statsService.GetDifficultyStats(ctx)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting difficulty stats: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "=== Difficulty Statistics ===")
		for difficulty, stats := range difficultyStats {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n", difficulty)
			fmt.Fprintf(cmd.OutOrStdout(), "  Attempted: %d\n", stats.Attempted)
			fmt.Fprintf(cmd.OutOrStdout(), "  Solved: %d\n", stats.Solved)
			fmt.Fprintf(cmd.OutOrStdout(), "  Success Rate: %.1f%%\n", stats.SuccessRate*100)
			fmt.Fprintf(cmd.OutOrStdout(), "  Average Time: %s\n", stats.AvgTime)
		}
	},
}

// RefactoredRecentActivityCmd demonstrates recent activity using service layer
var RefactoredRecentActivityCmd = &cobra.Command{
	Use:   "recent-refactored [days]",
	Short: "Show recent activity (refactored version)",
	Long:  `Display recent practice activity using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		statsService := services.DefaultRegistry.GetStatsService()
		
		days := 7 // default to 7 days
		if len(args) > 0 {
			// Parse days argument if provided
			if parsed, err := time.ParseDuration(args[0] + "h"); err == nil {
				days = int(parsed.Hours() / 24)
			}
		}
		
		ctx := context.Background()
		recentActivity, err := statsService.GetRecentActivity(ctx, days)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting recent activity: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "=== Recent Activity (Last %d days) ===\n", days)
		for _, daily := range recentActivity {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %d problems today (streak: %d days)\n",
				daily.Date, daily.ProblemsToday, daily.StreakDays)
			if len(daily.PatternsToday) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  Patterns: %v\n", daily.PatternsToday)
			}
		}
	},
}

func init() {
	// Add refactored subcommands to demonstrate service layer usage
	RefactoredStatsCmd.AddCommand(RefactoredPatternStatsCmd)
	RefactoredStatsCmd.AddCommand(RefactoredDifficultyStatsCmd)
	RefactoredStatsCmd.AddCommand(RefactoredRecentActivityCmd)
}