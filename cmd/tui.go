package cmd

import (
	"github.com/lancekrogers/algo-scales/internal/ui"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start Algo Scales with terminal UI",
	Long: `Start Algo Scales with a full-featured terminal UI that provides
language selection, timer configuration, and split-screen problem solving.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.StartTUI()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	
	// Add flags if needed
	tuiCmd.Flags().BoolP("config", "c", false, "Start in configuration mode")
	tuiCmd.Flags().StringP("language", "l", "", "Specify programming language")
	tuiCmd.Flags().IntP("timer", "t", 0, "Set timer duration in minutes")
	tuiCmd.Flags().StringP("mode", "m", "", "Set learning mode (learn, practice, cram)")
	tuiCmd.Flags().StringP("pattern", "p", "", "Focus on specific algorithm pattern")
}