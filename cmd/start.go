// Start command for practice sessions

package cmd

import (
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/spf13/cobra"
)

var (
	language   string
	timer      int
	pattern    string
	difficulty string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new practice session",
	Long:  `Start a new algorithm practice session in the specified mode.`,
}

// learnCmd represents the learn subcommand
var learnCmd = &cobra.Command{
	Use:   "learn [problem]",
	Short: "Start in Learn mode",
	Long: `Start a session in Learn mode, which includes pattern explanations 
and step-by-step solutions to help you understand the algorithm patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if err := session.Start(opts); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error starting session: %v\n", err)
		}
	},
}

// practiceCmd represents the practice subcommand
var practiceCmd = &cobra.Command{
	Use:   "practice [problem]",
	Short: "Start in Practice mode",
	Long: `Start a session in Practice mode, which hides solutions and hints 
but allows you to request them when needed.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if err := session.Start(opts); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error starting session: %v\n", err)
		}
	},
}

// cramCmd represents the cram subcommand
var cramCmd = &cobra.Command{
	Use:   "cram",
	Short: "Start in Cram mode",
	Long: `Start a session in Cram mode, which quickly cycles through problems
from the most common algorithm patterns, with a timer for each problem.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := session.Options{
			Mode:       session.CramMode,
			Language:   language,
			Timer:      timer,
			Pattern:    pattern,
			Difficulty: difficulty,
		}

		if err := session.Start(opts); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error starting session: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(learnCmd)
	startCmd.AddCommand(practiceCmd)
	startCmd.AddCommand(cramCmd)

	// Add flags to the start command and all subcommands
	startCmd.PersistentFlags().StringVarP(&language, "language", "l", "go", "Programming language (go, python, javascript)")
	startCmd.PersistentFlags().IntVarP(&timer, "timer", "t", 45, "Timer duration in minutes (15, 30, 45, 60)")
	startCmd.PersistentFlags().StringVarP(&pattern, "pattern", "p", "", "Algorithm pattern to focus on")
	startCmd.PersistentFlags().StringVarP(&difficulty, "difficulty", "d", "", "Problem difficulty (easy, medium, hard)")
}
