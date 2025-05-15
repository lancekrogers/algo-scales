package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/lancekrogers/algo-scales/internal/common/config"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/ui/screens"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start Algo Scales with terminal UI",
	Long: `Start Algo Scales with a full-featured terminal UI that provides
language selection, timer configuration, and split-screen problem solving.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if configuration mode is enabled
		configMode, _ := cmd.Flags().GetBool("config")
		
		if configMode {
			// Start with setup screen
			startSetupUI()
		} else {
			// Load config and start practice session
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading config: %v. Starting setup...\n", err)
				startSetupUI()
				return
			}
			
			// Start session UI
			startSessionUI(cfg)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	
	// Add flags
	tuiCmd.Flags().BoolP("config", "c", false, "Start in configuration mode")
	tuiCmd.Flags().StringP("language", "l", "", "Specify programming language")
	tuiCmd.Flags().IntP("timer", "t", 0, "Set timer duration in minutes")
	tuiCmd.Flags().StringP("mode", "m", "", "Set learning mode (learn, practice, cram)")
	tuiCmd.Flags().StringP("pattern", "p", "", "Focus on specific algorithm pattern")
}

// startSetupUI starts the setup user interface
func startSetupUI() {
	setupModel := screens.NewSetupModel()
	p := tea.NewProgram(setupModel, tea.WithAltScreen())
	
	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running setup UI: %v\n", err)
		os.Exit(1)
	}
	
	// Load saved config and start session
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	
	startSessionUI(cfg)
}

// startSessionUI starts the session user interface
func startSessionUI(cfg config.UserConfig) {
	// Load problems
	problems, err := problem.LoadLocalProblems()
	if err != nil {
		fmt.Printf("Error loading problems: %v\n", err)
		os.Exit(1)
	}
	
	// Check if we have problems
	if len(problems) == 0 {
		fmt.Println("No problems found. Please check your installation.")
		os.Exit(1)
	}
	
	// Start problem selection screen
	problemSelectionModel := screens.NewProblemSelectionModel(problems, cfg.Language, cfg.Mode)
	p := tea.NewProgram(problemSelectionModel, tea.WithAltScreen())
	
	// Run the problem selection program
	selection, err := p.Run()
	if err != nil {
		fmt.Printf("Error running problem selection UI: %v\n", err)
		os.Exit(1)
	}
	
	// Check if we got a problem selection
	if selectionModel, ok := selection.(screens.ProblemSelectionModel); ok {
		if selectionModel.SelectedProblem == nil {
			fmt.Println("No problem selected. Exiting.")
			os.Exit(0)
		}
		
		// Create session model with selected problem
		sessionModel := screens.NewSessionModel(
			selectionModel.SelectedProblem,
			cfg.Mode,
			cfg.Language,
			selectionModel.SelectedPattern,
		)
		
		// Create program for the session
		sessionProgram := tea.NewProgram(sessionModel, tea.WithAltScreen())
		
		// Run the session program
		if _, err := sessionProgram.Run(); err != nil {
			fmt.Printf("Error running session UI: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Problem selection aborted. Exiting.")
		os.Exit(0)
	}
}