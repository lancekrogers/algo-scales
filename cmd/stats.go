// Stats command for viewing progress

package cmd

import (
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/stats"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your problem-solving statistics",
	Long:  `View statistics about your algorithm problem-solving performance.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior shows summary stats
		statistics, err := stats.GetSummary()
		if err != nil {
			fmt.Printf("Error retrieving stats: %v\n", err)
			return
		}

		fmt.Println("Overall Statistics:")
		fmt.Printf("Total Problems Attempted: %d\n", statistics.TotalAttempted)
		fmt.Printf("Total Problems Solved: %d\n", statistics.TotalSolved)
		fmt.Printf("Average Solve Time: %s\n", statistics.AvgSolveTime)
		fmt.Printf("Fastest Solve: %s (%s)\n", statistics.FastestSolve.Time, statistics.FastestSolve.ProblemID)
		fmt.Printf("Most Challenging: %s (attempts: %d)\n", statistics.MostChallenging.ProblemID, statistics.MostChallenging.Attempts)
	},
}

// patternStatsCmd represents the patterns subcommand for stats
var patternStatsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "View stats by pattern",
	Long:  `View your algorithm problem-solving statistics organized by pattern.`,
	Run: func(cmd *cobra.Command, args []string) {
		patternStats, err := stats.GetByPattern()
		if err != nil {
			fmt.Printf("Error retrieving pattern stats: %v\n", err)
			return
		}

		fmt.Println("Stats by Pattern:")
		for pattern, pstat := range patternStats {
			fmt.Printf("\n%s:\n", pattern)
			fmt.Printf("  Attempted: %d, Solved: %d\n", pstat.Attempted, pstat.Solved)
			fmt.Printf("  Success Rate: %.1f%%\n", pstat.SuccessRate)
			fmt.Printf("  Average Time: %s\n", pstat.AvgTime)
		}
	},
}

// trendsCmd represents the trends subcommand for stats
var trendsCmd = &cobra.Command{
	Use:   "trends",
	Short: "View progress trends",
	Long:  `View your progress trends over time.`,
	Run: func(cmd *cobra.Command, args []string) {
		trends, err := stats.GetTrends()
		if err != nil {
			fmt.Printf("Error retrieving trend stats: %v\n", err)
			return
		}

		fmt.Println("Progress Trends:")
		fmt.Println("Last 7 Days:")
		for _, day := range trends.Daily {
			fmt.Printf("  %s: %d solved (avg time: %s)\n", day.Date, day.Solved, day.AvgTime)
		}

		fmt.Println("\nWeekly Progress:")
		for _, week := range trends.Weekly {
			fmt.Printf("  Week of %s: %d solved (success rate: %.1f%%)\n", week.StartDate, week.Solved, week.SuccessRate)
		}
	},
}

// resetStatsCmd represents the reset subcommand for stats
var resetStatsCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset statistics",
	Long:  `Reset your problem-solving statistics.`,
	Run: func(cmd *cobra.Command, args []string) {
		confirm := false
		fmt.Print("Are you sure you want to reset all statistics? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			confirm = true
		}

		if !confirm {
			fmt.Println("Operation cancelled.")
			return
		}

		if err := stats.Reset(); err != nil {
			fmt.Printf("Error resetting stats: %v\n", err)
			return
		}

		fmt.Println("Statistics have been reset.")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.AddCommand(patternStatsCmd)
	statsCmd.AddCommand(trendsCmd)
	statsCmd.AddCommand(resetStatsCmd)
}
