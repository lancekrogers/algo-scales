// Refactored Stats command using service layer
package cmd

import (
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
		overall, err := statsService.GetOverallStats()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting overall stats: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "=== Overall Statistics ===")
		fmt.Fprintf(cmd.OutOrStdout(), "Total Sessions: %d\n", overall.TotalSessions)
		fmt.Fprintf(cmd.OutOrStdout(), "Solved Problems: %d\n", overall.SolvedProblems)
		fmt.Fprintf(cmd.OutOrStdout(), "Unsolved Problems: %d\n", overall.UnsolvedProblems)
		fmt.Fprintf(cmd.OutOrStdout(), "Success Rate: %.1f%%\n", overall.SuccessRate*100)
		fmt.Fprintf(cmd.OutOrStdout(), "Average Time: %s\n", overall.AverageTime.String())
		fmt.Fprintf(cmd.OutOrStdout(), "Total Time: %s\n", overall.TotalTime.String())
		fmt.Fprintf(cmd.OutOrStdout(), "Current Streak: %d\n", overall.CurrentStreak)
		fmt.Fprintf(cmd.OutOrStdout(), "Longest Streak: %d\n", overall.LongestStreak)
		fmt.Fprintf(cmd.OutOrStdout(), "Favorite Pattern: %s\n", overall.FavoritePattern)
		fmt.Fprintf(cmd.OutOrStdout(), "Favorite Language: %s\n", overall.FavoriteLanguage)
	},
}

// RefactoredPatternStatsCmd demonstrates pattern stats using service layer
var RefactoredPatternStatsCmd = &cobra.Command{
	Use:   "patterns-refactored",
	Short: "Show pattern-specific statistics (refactored version)",
	Long:  `Display performance statistics by algorithm pattern using the service layer.`,
	Run: func(cmd *cobra.Command, args []string) {
		statsService := services.DefaultRegistry.GetStatsService()
		
		patternStats, err := statsService.GetPatternStats()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting pattern stats: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "=== Pattern Statistics ===")
		for pattern, stats := range patternStats {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n", pattern)
			fmt.Fprintf(cmd.OutOrStdout(), "  Sessions: %d\n", stats.TotalSessions)
			fmt.Fprintf(cmd.OutOrStdout(), "  Solved: %d\n", stats.Solved)
			fmt.Fprintf(cmd.OutOrStdout(), "  Unsolved: %d\n", stats.Unsolved)
			fmt.Fprintf(cmd.OutOrStdout(), "  Success Rate: %.1f%%\n", stats.SuccessRate*100)
			fmt.Fprintf(cmd.OutOrStdout(), "  Average Time: %s\n", stats.AverageTime.String())
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
		
		difficultyStats, err := statsService.GetDifficultyStats()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting difficulty stats: %v\n", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "=== Difficulty Statistics ===")
		for difficulty, stats := range difficultyStats {
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n", difficulty)
			fmt.Fprintf(cmd.OutOrStdout(), "  Sessions: %d\n", stats.TotalSessions)
			fmt.Fprintf(cmd.OutOrStdout(), "  Solved: %d\n", stats.Solved)
			fmt.Fprintf(cmd.OutOrStdout(), "  Unsolved: %d\n", stats.Unsolved)
			fmt.Fprintf(cmd.OutOrStdout(), "  Success Rate: %.1f%%\n", stats.SuccessRate*100)
			fmt.Fprintf(cmd.OutOrStdout(), "  Average Time: %s\n", stats.AverageTime.String())
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
		
		recentActivity, err := statsService.GetRecentActivity(days)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error getting recent activity: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "=== Recent Activity (Last %d days) ===\n", days)
		for _, daily := range recentActivity {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %d sessions (%d solved, %d unsolved)\n",
				daily.Date.Format("2006-01-02"), daily.Sessions, daily.Solved, daily.Unsolved)
		}
	},
}

func init() {
	// Add refactored subcommands to demonstrate service layer usage
	RefactoredStatsCmd.AddCommand(RefactoredPatternStatsCmd)
	RefactoredStatsCmd.AddCommand(RefactoredDifficultyStatsCmd)
	RefactoredStatsCmd.AddCommand(RefactoredRecentActivityCmd)
}