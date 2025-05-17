// Main entry point for the application

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lancekrogers/algo-scales/cmd"
	"github.com/lancekrogers/algo-scales/internal/ui"
	"github.com/lancekrogers/algo-scales/internal/ui/splitscreen"
)

func main() {
	// Set up global signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Create a channel that can be closed to stop the goroutine
	stopSignalHandler := make(chan struct{})
	
	// Handle Ctrl+C in a separate goroutine for all modes
	go func() {
		select {
		case <-sigChan:
			fmt.Println("\nExiting AlgoScales. Thanks for practicing!")
			os.Exit(0)
		case <-stopSignalHandler:
			// Clean exit if we need to stop the handler
			return
		}
	}()

	// Check command line flags
	useLegacyCLI := false
	fromVim := false
	splitScreenUI := false // Split-screen UI mode

	for _, arg := range os.Args {
		switch arg {
		case "--cli", "--legacy":
			useLegacyCLI = true
		case "--debug":
			// Set DEBUG environment variable for the UI code
			os.Setenv("DEBUG", "1")
		case "--vim-mode":
			fromVim = true
		case "--split", "--tui", "--splitscreen":
			splitScreenUI = true
		}
	}

	// First attempt to detect if we're in a non-interactive terminal
	// or one that doesn't support the features we need
	if !isInteractiveTerminal() {
		fmt.Println("Non-interactive terminal detected. Using CLI mode.")
		useLegacyCLI = true
	}

	// Filter out our custom flags before passing to cobra
	filteredArgs := make([]string, 0)
	for _, arg := range os.Args {
		if arg != "--cli" && arg != "--legacy" && arg != "--debug" &&
		   arg != "--vim-mode" && arg != "--split" && arg != "--tui" && arg != "--splitscreen" {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	os.Args = filteredArgs

	// Decide which mode to use
	if useLegacyCLI || fromVim {
		// Use traditional CLI for legacy mode or Neovim integration
		cmd.Execute()
	} else if splitScreenUI {
		// Use the split-screen UI
		fmt.Println("Starting AlgoScales split-screen UI...")
		
		// Try to start the UI with a sample problem
		err := splitscreen.StartWithSampleProblem()
		if err != nil {
			fmt.Printf("Error launching split-screen UI: %v\n", err)
			fmt.Println("Falling back to CLI mode. Use --cli for command line.")
			cmd.Execute()
		}
	} else {
		// Launch the original interactive TUI by default
		// With improved error handling and terminal recovery
		fmt.Println("Starting AlgoScales interactive UI...")
		fmt.Println("If your terminal freezes, press Ctrl+C and restart with --split")

		// Add an extra debug message if in debug mode
		if os.Getenv("DEBUG") == "1" {
			fmt.Println("Debug mode: Using interactive TUI with enhanced error handling")
		}

		// Use channels for coordinating UI startup and potential fallback
		done := make(chan bool)
		errChan := make(chan error)

		// Use a separate goroutine to prevent blocking the main thread
		go func() {
			// Create and start the UI app with enhanced error handling
			app := ui.NewApp()
			err := app.Start()
			if err != nil {
				// Return error to main thread for proper handling
				errChan <- err
			} else {
				done <- true
			}
		}()

		// Set a timeout to prevent indefinite freezing with better fallback strategy
		select {
		case <-done:
			// UI exited normally
			fmt.Println("UI session complete")
		
		case err := <-errChan:
			// UI reported an error, handle more gracefully
			if os.Getenv("DEBUG") == "1" {
				fmt.Printf("Debug mode: UI Error details: %v\n", err)
			} else {
				fmt.Printf("Error launching interactive UI: %v\n", err)
			}
			
			// More detailed user guidance
			fmt.Println("\nIt looks like your terminal may not support the interactive UI.")
			fmt.Println("Trying split-screen UI as a fallback...")
			
			// Force terminal reset before trying alternative UI
			fmt.Print("\033c") // Terminal reset escape code
			time.Sleep(200 * time.Millisecond)
			
			// Try the split-screen UI as fallback - it has better terminal compatibility
			err = splitscreen.StartWithSampleProblem()
			if err != nil {
				if os.Getenv("DEBUG") == "1" {
					fmt.Printf("Debug mode: Split-screen UI error: %v\n", err)
				} else {
					fmt.Printf("Split-screen UI also failed: %v\n", err)
				}
				fmt.Println("Falling back to CLI mode.")
				cmd.Execute()
			}
		
		case <-time.After(3 * time.Second): // Reduced timeout for faster feedback
			// Terminal reset to prevent frozen state and clear the screen
			fmt.Print("\033c") 
			fmt.Println("\nUI initialization timed out. Your terminal might not fully support the interactive UI.")
			fmt.Println("Trying split-screen UI as a fallback...")
			
			// Small delay to ensure terminal is reset
			time.Sleep(200 * time.Millisecond)
			
			// Try the split-screen UI as fallback - it has better terminal compatibility
			err := splitscreen.StartWithSampleProblem()
			if err != nil {
				fmt.Printf("Split-screen UI also failed: %v\n", err)
				fmt.Println("Falling back to CLI mode.")
				cmd.Execute()
			}
		}
	}
}

// isInteractiveTerminal checks if we're running in an interactive terminal
// that supports the features we need
func isInteractiveTerminal() bool {
	// Use a more comprehensive check for terminal capabilities
	
	// Try opening /dev/tty first, but don't fail if it doesn't exist
	tty, err := os.Open("/dev/tty")
	if err == nil {
		tty.Close()
	} else {
		// If we can't open /dev/tty, we're likely not in an interactive terminal
		// but we'll continue with additional checks before giving up completely
		if os.Getenv("DEBUG") == "1" {
			fmt.Println("Debug mode: TTY check failed:", err)
		}
	}
	
	// Check for environment variables that indicate we're running in a CI system
	if os.Getenv("CI") != "" || os.Getenv("CONTINUOUS_INTEGRATION") != "" {
		return false
	}

	// Check if TERM is set to something useful
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}
	
	// Additional checks for terminal type
	termProgram := os.Getenv("TERM_PROGRAM")
	if termProgram == "iTerm.app" || 
	   termProgram == "Apple_Terminal" || 
	   strings.Contains(term, "xterm") || 
	   strings.Contains(term, "screen") {
		// These are generally well-supported terminals
		return true
	}
	
	// Check if we can get terminal size as a final test
	// This usually works even in minimal terminals
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil || len(output) == 0 {
		// If stty fails, we probably don't have a proper terminal
		return false
	}
	
	return true
}