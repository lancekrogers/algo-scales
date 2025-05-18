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
			os.Setenv("VIM_MODE", "1")
		case "--split", "--tui", "--splitscreen":
			// Force split-screen UI mode
			splitScreenUI = true
		}
	}

	// CLI mode detection
	if fromVim || (!isTerminal() && !splitScreenUI) {
		// Automatically set simple mode if not in terminal
		os.Setenv("SIMPLE_MODE", "1")
		// Stop the signal handler for non-UI modes
		close(stopSignalHandler)
		// Execute CLI commands directly
		cmd.Execute()
		return
	}

	// If using legacy CLI or no terminal detected, use the cobra CLI
	if useLegacyCLI || !isTerminal() {
		// Stop the signal handler for non-UI modes
		close(stopSignalHandler)
		cmd.Execute()
		return
	}

	// If explicitly using split-screen UI
	if splitScreenUI {
		if err := splitscreen.StartUI(nil); err != nil {
			fmt.Printf("Error running split-screen UI: %v\n", err)
			fmt.Println("Falling back to standard TUI...")
			// Fall through to standard TUI
		} else {
			return // Split-screen successful, exit
		}
	}

	// Default to the terminal UI (TUI)
	// Stop the signal handler for UI mode - let UI handle its own signals
	close(stopSignalHandler)
	
	// Add timeout mechanism and error handling
	if isTerminal() {
		// Use the timeout approach only for terminal UIs
		startWithTimeout()
	} else {
		// For non-terminal environments, use CLI directly
		cmd.Execute()
	}
}

// Start TUI with timeout and error recovery
func startWithTimeout() {
	// Use channels for coordinating UI startup and potential fallback
	done := make(chan bool)
	errChan := make(chan error)

	// Use a separate goroutine to prevent blocking the main thread
	go func() {
		// Create and start the UI app with enhanced error handling
		err := ui.StartTUI()
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
		// UI started successfully, normal exit
		return
	case err := <-errChan:
		// UI failed - provide error information and fallback to CLI
		fmt.Printf("Error starting TUI: %v\n", err)
		fmt.Println("Falling back to CLI mode...")
		
		// Give brief pause for user to see error
		time.Sleep(1 * time.Second)
		cmd.Execute()
	case <-time.After(5 * time.Second):
		// Timeout - inform user and suggest CLI mode
		fmt.Println("\nUI initialization taking too long. This might be due to terminal compatibility issues.")
		fmt.Println("You can:")
		fmt.Println("1. Press Ctrl+C to exit and run with --cli flag for command-line mode")
		fmt.Println("2. Wait a bit longer for UI to load")
		fmt.Println("3. Check if your terminal supports interactive UIs")
		
		// Continue waiting for UI with extended timeout
		select {
		case <-done:
			return
		case err := <-errChan:
			fmt.Printf("\nError: %v\n", err)
			fmt.Println("Falling back to CLI mode...")
			time.Sleep(1 * time.Second)
			cmd.Execute()
		case <-time.After(10 * time.Second):
			// Final timeout - exit gracefully
			fmt.Println("\nUI failed to start. Please run with --cli flag for command-line mode.")
			os.Exit(1)
		}
	}
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
	
	// Platform-specific terminal checks (optional)
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