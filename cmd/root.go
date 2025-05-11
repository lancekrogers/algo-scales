// Root command implementation

package cmd


import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lancekrogers/algo-scales/internal/license"
	"github.com/lancekrogers/algo-scales/internal/api"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "algo-scales",
	Short: "Algorithm study tool for interview preparation",
	Long: `Algo Scales is a terminal-based algorithm study tool designed to help developers
prepare for coding interviews efficiently. It focuses on teaching common algorithm
patterns through curated problems and features different learning modes.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

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
