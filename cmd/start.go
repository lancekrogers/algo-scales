// Start command for practice sessions

package cmd

import (
	"fmt"
	"os"

	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/lancekrogers/algo-scales/internal/ui"
	"github.com/lancekrogers/algo-scales/internal/ui/splitscreen"
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
			return
		}
		
		// Launch the appropriate UI
		if err := launchUI(cmd); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error launching UI: %v\n", err)
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
			return
		}
		
		// Launch the appropriate UI
		if err := launchUI(cmd); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error launching UI: %v\n", err)
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
			return
		}
		
		// Launch the appropriate UI
		if err := launchUI(cmd); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error launching UI: %v\n", err)
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

// launchUI determines which UI to launch based on flags
func launchUI(cmd *cobra.Command) error {
	// Skip UI launch during tests
	if os.Getenv("TESTING") == "1" {
		return nil
	}
	
	// Check flags to determine UI mode
	useTUI, _ := cmd.Root().PersistentFlags().GetBool("tui")
	useSplit, _ := cmd.Root().PersistentFlags().GetBool("split")
	splitscreenFlag, _ := cmd.Root().PersistentFlags().GetBool("splitscreen")
	vimMode, _ := cmd.Root().PersistentFlags().GetBool("vim-mode")
	
	// Set VIM_MODE environment variable if needed
	if vimMode {
		os.Setenv("VIM_MODE", "1")
	}
	
	// Determine if any TUI mode is requested
	useSplitScreen := useSplit || splitscreenFlag
	
	// Use split-screen UI if requested
	if useSplitScreen && isTerminal() {
		return splitscreen.StartUI(nil)
	} else if useTUI && isTerminal() {
		// Use standard TUI if requested
		return ui.StartTUI()
	}
	
	// Default to TUI mode for start commands (interactive problem solving)
	if isTerminal() {
		return ui.StartTUI()
	}
	
	// If not in terminal, print a message
	fmt.Println("Session created successfully!")
	fmt.Println("Run with --tui flag for interactive mode.")
	return nil
}

