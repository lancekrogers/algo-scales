// CLI display for statistics
package cmd

import (
	"fmt"
	"time"

	"github.com/lancekrogers/algo-scales/internal/stats"
	"github.com/spf13/cobra"
)

// cliStatsCmd represents the stats display for CLI mode
var cliStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View problem-solving statistics in CLI mode",
	Long:  `Display your problem-solving statistics and progress in CLI mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayCLIStats()
	},
}

func init() {
	cliCmd.AddCommand(cliStatsCmd)
}

// displayCLIStats shows statistics in CLI mode
func displayCLIStats() {
	fmt.Println("\n📊 AlgoScales Statistics 📊")
	fmt.Println("—————————————————————————")
	
	// Get overall stats
	sessions, err := stats.GetAllSessions()
	if err != nil {
		fmt.Printf("Error loading statistics: %v\n", err)
		return
	}
	
	if len(sessions) == 0 {
		fmt.Println("No sessions recorded yet. Start solving problems to build statistics!")
		return
	}
	
	// Overall summary
	totalSessions := len(sessions)
	solvedCount := 0
	var totalTime time.Duration
	
	for _, s := range sessions {
		if s.Solved {
			solvedCount++
		}
		totalTime += s.Duration
	}
	
	successRate := float64(solvedCount) / float64(totalSessions) * 100
	avgTime := totalTime / time.Duration(totalSessions)
	
	fmt.Printf("\n## Overall Progress\n")
	fmt.Printf("Total Problems Attempted: %d\n", totalSessions)
	fmt.Printf("Problems Solved: %d (%.1f%%)\n", solvedCount, successRate)
	fmt.Printf("Average Session Time: %s\n", formatDuration(avgTime))
	
	// Pattern breakdown
	fmt.Printf("\n## Progress by Pattern\n")
	patternStats := make(map[string]struct {
		Attempted int
		Solved    int
	})
	
	for _, s := range sessions {
		for _, pattern := range s.Patterns {
			stats := patternStats[pattern]
			stats.Attempted++
			if s.Solved {
				stats.Solved++
			}
			patternStats[pattern] = stats
		}
	}
	
	for pattern, stats := range patternStats {
		successRate := float64(stats.Solved) / float64(stats.Attempted) * 100
		fmt.Printf("- %s: %d/%d solved (%.1f%%)\n", pattern, stats.Solved, stats.Attempted, successRate)
	}
	
	// Recent activity
	fmt.Printf("\n## Recent Activity\n")
	recent := getRecentSessions(sessions, 5)
	for i, s := range recent {
		solved := "❌"
		if s.Solved {
			solved = "✅"
		}
		fmt.Printf("%d. %s %s [%s] - %s\n", i+1, solved, s.ProblemID, JoinStrings(s.Patterns), formatTime(s.EndTime))
	}
	
	// Streak information
	streak, lastDate := calculateStreak(sessions)
	fmt.Printf("\n## Practice Streak\n")
	fmt.Printf("Current Streak: %d days\n", streak)
	if streak > 0 {
		fmt.Printf("Last Practice: %s\n", formatDate(lastDate))
	}
}

// getRecentSessions returns the n most recent sessions
func getRecentSessions(sessions []stats.SessionStats, n int) []stats.SessionStats {
	// Sort sessions by end time (most recent first)
	// Note: In a real implementation, we'd sort properly
	// For simplicity, we'll assume they're already sorted
	
	if len(sessions) <= n {
		return sessions
	}
	return sessions[len(sessions)-n:]
}

// calculateStreak determines the current streak of consecutive days
func calculateStreak(sessions []stats.SessionStats) (int, time.Time) {
	if len(sessions) == 0 {
		return 0, time.Time{}
	}
	
	// Find the most recent session
	var mostRecent time.Time
	for _, s := range sessions {
		if s.EndTime.After(mostRecent) {
			mostRecent = s.EndTime
		}
	}
	
	// Check if it's today or yesterday
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	
	mostRecentDay := mostRecent.Format("2006-01-02")
	todayStr := today.Format("2006-01-02")
	yesterdayStr := yesterday.Format("2006-01-02")
	
	if mostRecentDay != todayStr && mostRecentDay != yesterdayStr {
		// Streak broken - more than a day since last session
		return 0, mostRecent
	}
	
	// Count consecutive days
	streak := 1
	currentDay := mostRecent
	
	// Map to track dates with sessions
	datesWithSessions := make(map[string]bool)
	for _, s := range sessions {
		dateStr := s.EndTime.Format("2006-01-02")
		datesWithSessions[dateStr] = true
	}
	
	// Count back from the most recent day
	for {
		prevDay := currentDay.AddDate(0, 0, -1)
		prevDayStr := prevDay.Format("2006-01-02")
		
		if datesWithSessions[prevDayStr] {
			streak++
			currentDay = prevDay
		} else {
			break
		}
	}
	
	return streak, mostRecent
}

// formatDuration formats a duration nicely
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// formatTime formats time nicely
func formatTime(t time.Time) string {
	return t.Format("Jan 2, 2006 3:04 PM")
}

// formatDate formats date nicely
func formatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}