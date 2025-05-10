// Daily scale practice command
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/spf13/cobra"
)

// Patterns organized as musical scales
var scales = []struct {
	Pattern     string
	MusicalName string
	Description string
}{
	{
		Pattern:     "sliding-window",
		MusicalName: "C Major",
		Description: "The fundamental scale, elegant and versatile",
	},
	{
		Pattern:     "two-pointers",
		MusicalName: "G Major",
		Description: "Balanced and efficient, the workhorse of array manipulation",
	},
	{
		Pattern:     "fast-slow-pointers",
		MusicalName: "D Major",
		Description: "The cycle detector, bright and revealing",
	},
	{
		Pattern:     "hash-map",
		MusicalName: "A Major",
		Description: "The lookup accelerator, crisp and direct",
	},
	{
		Pattern:     "binary-search",
		MusicalName: "E Major",
		Description: "The divider and conqueror, precise and logarithmic",
	},
	{
		Pattern:     "dfs",
		MusicalName: "B Major",
		Description: "The deep explorer, rich and thorough",
	},
	{
		Pattern:     "bfs",
		MusicalName: "F# Major",
		Description: "The level-by-level discoverer, methodical and complete",
	},
	{
		Pattern:     "dynamic-programming",
		MusicalName: "Db Major",
		Description: "The optimizer, complex and powerful",
	},
	{
		Pattern:     "greedy",
		MusicalName: "Ab Major",
		Description: "The local maximizer, bold and decisive",
	},
	{
		Pattern:     "union-find",
		MusicalName: "Eb Major",
		Description: "The connector, structured and organized",
	},
	{
		Pattern:     "heap",
		MusicalName: "Bb Major",
		Description: "The sorter, flexible and maintaining order",
	},
}

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

// ScaleProgress tracks progress through scales
type ScaleProgress struct {
	Current       int       `json:"current"`
	LastPracticed time.Time `json:"last_practiced"`
	Completed     []string  `json:"completed"`
}

// startDailyScale starts the daily scale practice
func startDailyScale() {
	// First, display welcome message
	fmt.Println("♪ AlgoScales Daily Practice ♪")
	fmt.Println("----------------------------")
	fmt.Println("Practice one problem from each algorithm pattern (scale) to build your skills.")
	fmt.Println("Just as a musician practices scales daily, this routine will help you master")
	fmt.Println("the fundamental patterns of algorithm problem-solving.")
	fmt.Println()

	// Load progress if it exists
	progress := loadProgress()

	// Get starting index
	startIdx := progress.Current
	if startIdx >= len(scales) {
		// If we've completed all scales, start over
		startIdx = 0
		progress.Completed = []string{}
	}

	// Check if we're continuing from a previous day
	today := time.Now().Format("2006-01-02")
	lastPracticedDay := progress.LastPracticed.Format("2006-01-02")

	if lastPracticedDay != today {
		// Starting a new day, reset completed scales
		progress.Completed = []string{}
		progress.Current = 0
		startIdx = 0
	}

	// Find the first scale that hasn't been completed
	for startIdx < len(scales) {
		if !contains(progress.Completed, scales[startIdx].Pattern) {
			break
		}
		startIdx++
	}

	// If all scales are completed, show congratulations
	if startIdx >= len(scales) {
		fmt.Println("🎉 Congratulations! You've completed your daily scales practice for all patterns!")
		fmt.Println("Feel free to practice more specific patterns or try a different mode.")
		return
	}

	// Update and save progress
	progress.Current = startIdx
	progress.LastPracticed = time.Now()
	saveProgress(progress)

	// Show current scale information
	scale := scales[startIdx]
	fmt.Printf("Now practicing: %s (%s)\n", scale.MusicalName, scale.Pattern)
	fmt.Printf("Description: %s\n\n", scale.Description)

	// Start practice session with this pattern
	opts := session.Options{
		Mode:       session.PracticeMode,
		Language:   language,
		Timer:      timer,
		Pattern:    scale.Pattern,
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
		progress.Completed = append(progress.Completed, scale.Pattern)

		// Move to next scale
		progress.Current = startIdx + 1
		saveProgress(progress)

		// If there are more scales, ask if user wants to continue
		if progress.Current < len(scales) {
			nextScale := scales[progress.Current]
			fmt.Println()
			fmt.Printf("You've completed %s! Next scale: %s (%s)\n",
				scale.MusicalName, nextScale.MusicalName, nextScale.Pattern)
			fmt.Print("Continue to next scale? (y/n): ")

			var response string
			fmt.Scanln(&response)

			if response == "y" || response == "Y" {
				startDailyScale() // Recursively start the next scale
			} else {
				fmt.Println("Practice session paused. You can continue later with 'algo-scales daily'")
			}
		} else {
			// All scales completed!
			fmt.Println()
			fmt.Println("🎵 Congratulations! You've completed your daily scales practice for all patterns! 🎵")
			fmt.Println("Keep up the good work and maintain your practice streak.")

			// Reset for tomorrow
			progress.Current = 0
			progress.Completed = []string{}
			saveProgress(progress)
		}
	}
}

// loadProgress loads the scale progress from a file
func loadProgress() ScaleProgress {
	// For MVP, this could be a simple file in the config directory
	// A more complex implementation would use a database

	// Default progress (starting fresh)
	return ScaleProgress{
		Current:       0,
		LastPracticed: time.Time{}, // Zero time (never practiced)
		Completed:     []string{},
	}
}

// saveProgress saves the scale progress to a file
func saveProgress(progress ScaleProgress) {
	// For MVP, just pretend we saved it
	// In a real implementation, this would write to a file or database
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(dailyCmd)

	// Use the same flags as start command for consistency
	dailyCmd.Flags().StringVarP(&language, "language", "l", "go", "Programming language (go, python, javascript)")
	dailyCmd.Flags().IntVarP(&timer, "timer", "t", 45, "Timer duration in minutes (15, 30, 45, 60)")
	dailyCmd.Flags().StringVarP(&difficulty, "difficulty", "d", "", "Problem difficulty (easy, medium, hard)")
}
