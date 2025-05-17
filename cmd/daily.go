// Daily scale practice command
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lancekrogers/algo-scales/internal/daily"
	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/spf13/cobra"
)

// dailyCmd represents the daily command for daily scale practice
var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Start daily scale practice",
	Long: `Start your daily scale practice session.
	
Just as musicians practice scales daily to build technique, developers can practice
algorithm patterns daily to build problem-solving intuition. This command will guide
you through one problem from each major algorithm pattern (musical scale).`,
	Run: func(cmd *cobra.Command, args []string) {
		startDailyScale()
	},
}

// startDailyScale starts the daily scale practice session
func startDailyScale() {
	// Display welcome message
	fmt.Println("╭───────────────────────────────────────────────────────────────╮")
	fmt.Println("│                🎵 AlgoScales Daily Practice 🎵                │")
	fmt.Println("╰───────────────────────────────────────────────────────────────╯")
	fmt.Println("")
	fmt.Println("Practice one problem from each algorithm pattern (scale) to build your skills.")
	fmt.Println("Just as a musician practices scales daily, this routine will help you master")
	fmt.Println("the fundamental patterns of algorithm problem-solving.")
	fmt.Println("")

	// Load progress or start fresh
	progress, err := daily.LoadProgress()
	if err != nil {
		fmt.Printf("Error loading progress: %v\n", err)
		fmt.Println("Starting with fresh progress")
		progress = daily.ScaleProgress{
			Current:       0,
			LastPracticed: time.Time{},
			Completed:     []string{},
			Streak:        0,
			LongestStreak: 0,
		}
	}

	// Update streak based on last practice date
	daily.UpdateStreak(&progress)

	// Display streak information
	displayStreakInfo(progress)

	// Check if we're continuing from a previous day
	today := time.Now().Format("2006-01-02")
	var lastPracticedDay string
	if !progress.LastPracticed.IsZero() {
		lastPracticedDay = progress.LastPracticed.Format("2006-01-02")
	}

	if lastPracticedDay != today {
		// Starting a new day, reset completed scales
		progress.Completed = []string{}
		progress.Current = 0
	}

	// Find the next scale to practice
	nextScale := daily.GetNextScale(progress.Completed)

	// If all scales are completed, show congratulations
	if nextScale == nil {
		fmt.Println("🎉 Congratulations! You've completed your daily scales practice for all patterns!")
		fmt.Println("Feel free to practice more specific patterns or try a different mode.")
		fmt.Println("")
		fmt.Println("Your current streak: " + fmt.Sprintf("%d days", progress.Streak))
		fmt.Println("Your longest streak: " + fmt.Sprintf("%d days", progress.LongestStreak))
		return
	}

	// Show information about remaining patterns
	remaining := daily.GetRemainingPatterns(progress.Completed)
	fmt.Printf("Patterns completed today: %d/11\n", len(progress.Completed))
	fmt.Printf("Patterns remaining: %d/11\n\n", remaining)

	// Show current scale information
	fmt.Printf("Now practicing: %s (%s)\n", nextScale.MusicalName, nextScale.Pattern)
	fmt.Printf("Description: %s\n\n", nextScale.Description)

	// Update progress with current pattern
	progress.Current = daily.GetPatternIndex(nextScale.Pattern)
	progress.LastPracticed = time.Now()

	// Save progress
	if err := daily.SaveProgress(progress); err != nil {
		fmt.Printf("Warning: Error saving progress: %v\n", err)
	}

	// Start practice session with this pattern
	opts := session.Options{
		Mode:       session.PracticeMode,
		Language:   language,
		Timer:      timer,
		Pattern:    nextScale.Pattern,
		Difficulty: difficulty,
	}

	// Create session for Vim mode or start directly
	if vimMode {
		s, err := session.CreateSession(opts)
		if err != nil {
			fmt.Printf("Error creating session: %v\n", err)
			os.Exit(1)
		}
		handleVimModeStart(s)
	} else {
		if err := session.Start(opts); err != nil {
			fmt.Printf("Error starting session: %v\n", err)
			os.Exit(1)
		}

		// When the session completes, mark this scale as completed
		progress.Completed = append(progress.Completed, nextScale.Pattern)

		// Save updated progress
		if err := daily.SaveProgress(progress); err != nil {
			fmt.Printf("Warning: Error saving progress: %v\n", err)
		}

		// Get the next scale (if any)
		nextScale = daily.GetNextScale(progress.Completed)

		// If there are more scales, ask if user wants to continue
		if nextScale != nil {
			fmt.Println()
			fmt.Printf("Pattern completed! Next scale: %s (%s)\n", 
				nextScale.MusicalName, nextScale.Pattern)
			fmt.Print("Continue to next scale? (y/n): ")

			var response string
			fmt.Scanln(&response)

			if response == "y" || response == "Y" {
				startDailyScale() // Recursively start the next scale
			} else {
				fmt.Println("Practice session paused. You can continue later with 'algo-scales daily'")
				fmt.Printf("Patterns completed today: %d/11\n", len(progress.Completed))
			}
		} else {
			// All scales completed!
			fmt.Println()
			fmt.Println("╭───────────────────────────────────────────────────────────────╮")
			fmt.Println("│         🎵 Congratulations! Daily Scales Complete! 🎵         │")
			fmt.Println("╰───────────────────────────────────────────────────────────────╯")
			fmt.Println()
			fmt.Println("You've completed all 11 algorithm pattern scales for today!")
			fmt.Println("Keep up the good work and maintain your practice streak.")
			fmt.Println()
			fmt.Printf("Current streak: %d days\n", progress.Streak)
			fmt.Printf("Longest streak: %d days\n", progress.LongestStreak)

			// Reset completion list for tomorrow but keep streak data
			progress.Completed = []string{}
			progress.Current = 0
			if err := daily.SaveProgress(progress); err != nil {
				fmt.Printf("Warning: Error saving progress: %v\n", err)
			}
		}
	}
}

// displayStreakInfo shows information about the user's practice streak
func displayStreakInfo(progress daily.ScaleProgress) {
	// Create a streak indicator
	var streakDisplay string
	if progress.Streak > 0 {
		flames := strings.Repeat("🔥", progress.Streak)
		if progress.Streak > 10 {
			flames = "🔥 x" + fmt.Sprintf("%d", progress.Streak)
		}
		streakDisplay = fmt.Sprintf("Current streak: %d days %s", progress.Streak, flames)
	} else {
		streakDisplay = "Start your streak today! 🎯"
	}

	// Display streak info
	fmt.Println(streakDisplay)
	if progress.LongestStreak > progress.Streak {
		fmt.Printf("Longest streak: %d days\n", progress.LongestStreak)
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(dailyCmd)

	// Use the same flags as start command for consistency
	dailyCmd.Flags().StringVarP(&language, "language", "l", "go", "Programming language (go, python, javascript)")
	dailyCmd.Flags().IntVarP(&timer, "timer", "t", 45, "Timer duration in minutes (15, 30, 45, 60)")
	dailyCmd.Flags().StringVarP(&difficulty, "difficulty", "d", "", "Problem difficulty (easy, medium, hard)")
}