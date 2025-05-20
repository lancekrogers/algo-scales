// Root command implementation

package cmd


import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lancekrogers/algo-scales/internal/api"
	"github.com/lancekrogers/algo-scales/internal/license"
	"github.com/lancekrogers/algo-scales/internal/ui"
	"github.com/lancekrogers/algo-scales/internal/ui/splitscreen"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "algo-scales",
	Short: "Algorithm study tool for interview preparation",
	Long: `Algo Scales is a terminal-based algorithm study tool designed to help developers
prepare for coding interviews efficiently. It focuses on teaching common algorithm
patterns through curated problems and features different learning modes.`,
	
	// Run the TUI by default
	Run: func(cmd *cobra.Command, args []string) {
		// Don't attempt to run TUI in test mode
		if os.Getenv("TESTING") == "1" {
			fmt.Fprintln(cmd.OutOrStdout(), "algo-scales - Algorithm study tool for interview preparation")
			return
		}
		
		// Check the flags to determine the UI mode
		useCLI, _ := cmd.Flags().GetBool("cli")
		legacyCLI, _ := cmd.Flags().GetBool("legacy") 
		vimMode, _ := cmd.Flags().GetBool("vim-mode")
		splitUI, _ := cmd.Flags().GetBool("split")
		tuiFlag, _ := cmd.Flags().GetBool("tui")
		splitscreenFlag, _ := cmd.Flags().GetBool("splitscreen")
		debugFlag, _ := cmd.Flags().GetBool("debug")
		
		// Set debug mode if flag is used
		if debugFlag {
			os.Setenv("DEBUG", "1")
		}
		
		// Determine final split screen mode from any of the aliased flags
		useSplitScreen := splitUI || tuiFlag || splitscreenFlag
		
		// Set VIM_MODE environment variable if needed
		if vimMode {
			os.Setenv("VIM_MODE", "1")
		}
		
		// Determine if we should use CLI mode
		useCliMode := useCLI || legacyCLI || vimMode || !isTerminal()
		
		if useCliMode {
			// Use CLI mode - just show help since no specific command was given
			if err := cmd.Help(); err != nil {
				fmt.Println("Error displaying help:", err)
			}
			return
		}
		
		// Use split-screen UI if requested
		if useSplitScreen {
			if err := splitscreen.StartUI(nil); err != nil {
				fmt.Printf("Error running split-screen UI: %v\n", err)
				fmt.Println("Falling back to standard TUI...")
				// Fall through to standard TUI
			} else {
				return // Split-screen successful, exit
			}
		}
		
		// Default to the terminal UI (TUI)
		if isTerminal() {
			// Start the TUI
			err := ui.StartTUI()
			if err != nil {
				fmt.Printf("Error starting TUI: %v\n", err)
				fmt.Println("Falling back to CLI mode...")
				// Show help as fallback
				if err := cmd.Help(); err != nil {
					fmt.Println("Error displaying help:", err)
				}
			}
		} else {
			// Show help for non-terminal environments
			if err := cmd.Help(); err != nil {
				fmt.Println("Error displaying help:", err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// hideFlags hides flags that shouldn't show in help output
func hideFlags(cmd *cobra.Command) {
	// Hide these flags from the command and all its subcommands
	flagsToHide := []string{"legacy", "tui", "splitscreen"}
	
	for _, flag := range flagsToHide {
		// Skip if the flag doesn't exist
		if cmd.PersistentFlags().Lookup(flag) != nil {
			cmd.PersistentFlags().MarkHidden(flag)
		} else if cmd.Flags().Lookup(flag) != nil {
			cmd.Flags().MarkHidden(flag)
		}
	}
	
	// Do the same for all subcommands
	for _, subCmd := range cmd.Commands() {
		hideFlags(subCmd)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	
	// Add global flags
	rootCmd.PersistentFlags().Bool("cli", false, "Use command-line interface mode (no TUI)")
	rootCmd.PersistentFlags().Bool("legacy", false, "Alias for --cli")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().Bool("vim-mode", false, "Use VIM-optimized mode")
	rootCmd.PersistentFlags().Bool("split", false, "Use split-screen TUI mode")
	rootCmd.PersistentFlags().Bool("tui", false, "Alias for --split")
	rootCmd.PersistentFlags().Bool("splitscreen", false, "Alias for --split")
	
	// Hide aliases from help output to avoid confusion
	hideFlags(rootCmd)
	
	// Check if first run and perform setup if needed
	if isFirstRun() {
		fmt.Println("Welcome to Algo Scales!")
		fmt.Println("Setting up your environment...")
		
		if err := setupConfigDir(); err != nil {
			fmt.Printf("Setup failed: %v\n", err)
			os.Exit(1)
		}
		
		if err := license.RequestLicense(); err != nil {
			fmt.Printf("License setup failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Downloading problem sets...")
		if err := api.DownloadProblems(true); err != nil {
			fmt.Printf("Problem download failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Setup complete! You're ready to start practicing.")
	}
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	// Set up config if needed
}

// isFirstRun checks if this is the first time the app is run
func isFirstRun() bool {
	// Skip setup during tests
	if os.Getenv("TESTING") == "1" {
		return false
	}
	configDir := getConfigDir()
	return !fileExists(configDir)
}

// setupConfigDir creates the necessary directories
func setupConfigDir() error {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(configDir, "problems"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(configDir, "stats"), 0755); err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file or directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// getConfigDir returns the configuration directory
// Exported as variable for testing
var getConfigDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}

// isTerminal checks if we're running in an actual terminal
func isTerminal() bool {
	// Check if we're running from vim
	if os.Getenv("VIM") != "" || os.Getenv("VIM_MODE") == "1" {
		return false
	}
	
	// Check if stdin is a terminal
	if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return false
	}
	
	// Additional terminal detection for different environments
	if term := os.Getenv("TERM"); term == "" || term == "dumb" {
		return false
	}
	
	// Check for specific CI/non-interactive environments
	if os.Getenv("CI") != "" || os.Getenv("JENKINS_URL") != "" {
		return false
	}
	
	// Platform-specific terminal checks
	if err := checkPlatformTerminal(); err != nil {
		return false
	}
	
	return true
}

// Platform-specific terminal checking
func checkPlatformTerminal() error {
	// Simple platform detection using 'which' command
	cmd := exec.Command("which", "tput")
	if err := cmd.Run(); err != nil {
		// tput not available, might not be a full terminal
		return err
	}
	
	// Check if terminal supports colors/interaction
	cmd = exec.Command("tput", "colors")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	
	// Parse color count
	colorCount := strings.TrimSpace(string(output))
	if colorCount == "" || colorCount == "0" {
		return fmt.Errorf("terminal does not support colors")
	}
	
	return nil
}
