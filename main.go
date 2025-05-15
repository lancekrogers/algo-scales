// Main entry point for the application

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lancekrogers/algo-scales/cmd"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/session"
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
	simpleUI := false     // A simple UI mode for limited terminals
	splitScreenUI := false // New split-screen UI mode

	for _, arg := range os.Args {
		switch arg {
		case "--cli", "--legacy":
			useLegacyCLI = true
		case "--debug":
			// Set DEBUG environment variable for the UI code
			os.Setenv("DEBUG", "1")
		case "--vim-mode":
			fromVim = true
		case "--simple":
			simpleUI = true
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
		   arg != "--vim-mode" && arg != "--simple" && 
		   arg != "--split" && arg != "--tui" && arg != "--splitscreen" {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	os.Args = filteredArgs

	// Decide which mode to use
	if useLegacyCLI || fromVim {
		// Use traditional CLI for legacy mode or Neovim integration
		cmd.Execute()
	} else if simpleUI {
		// Use a simple text-based UI that works in most terminals
		runSimpleTextUI()
	} else if splitScreenUI {
		// Use the new split-screen UI
		fmt.Println("Starting AlgoScales split-screen UI...")
		
		// Try to start the UI with a sample problem
		err := splitscreen.StartWithSampleProblem()
		if err != nil {
			fmt.Printf("Error launching split-screen UI: %v\n", err)
			fmt.Println("Falling back to CLI mode. Use --cli for command line or --simple for basic UI.")
			cmd.Execute()
		}
	} else {
		// Launch the original interactive TUI by default
		// With improved error handling and terminal recovery
		fmt.Println("Starting AlgoScales interactive UI...")
		fmt.Println("If your terminal freezes, press Ctrl+C and restart with --simple or --split")

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
			
			// Try the split-screen UI as first fallback - it has better terminal compatibility
			err = splitscreen.StartWithSampleProblem()
			if err != nil {
				if os.Getenv("DEBUG") == "1" {
					fmt.Printf("Debug mode: Split-screen UI error: %v\n", err)
				} else {
					fmt.Printf("Split-screen UI also failed: %v\n", err)
				}
				fmt.Println("Falling back to simple text mode...")
				runSimpleTextUI()
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
				fmt.Println("Falling back to simple mode which works in almost all terminals...")
				runSimpleTextUI()
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

// runSimpleTextUI runs a simple text-based UI that works in most terminals
func runSimpleTextUI() {
	// Reset signal handling specifically for the simple UI
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Create a reader with a timeout to allow checking for signals
	reader := bufio.NewReader(os.Stdin)

	// Print welcome banner
	fmt.Println("╭───────────────────────────────────────────────────────────────╮")
	fmt.Println("│         🎵  Welcome to AlgoScales - Simple Mode  🎵           │")
	fmt.Println("╰───────────────────────────────────────────────────────────────╯")
	fmt.Println()
	fmt.Println("Press Ctrl+C at any time to exit.")

	// Function to safely read user input with interrupt detection
	readInput := func() string {
		inputCh := make(chan string, 1)
		interruptCh := make(chan bool, 1)

		// Start a goroutine to read input
		go func() {
			text, err := reader.ReadString('\n')
			if err != nil {
				inputCh <- ""
				return
			}
			inputCh <- strings.TrimSpace(text)
		}()

		// Another goroutine to monitor for ctrl+c
		go func() {
			<-sigChan
			interruptCh <- true
		}()

		// Wait for either input or interrupt
		select {
		case input := <-inputCh:
			return input
		case <-interruptCh:
			fmt.Println("\nExiting AlgoScales. Thanks for practicing!")
			os.Exit(0)
			return ""
		}
	}

	for {
		// Show main menu
		fmt.Println("What would you like to do?")
		fmt.Println("1. Start Practice")
		fmt.Println("2. List Problems")
		fmt.Println("3. View Statistics")
		fmt.Println("4. Exit")
		fmt.Print("> ")

		// Use our custom input reader that handles interrupts
		choice := readInput()

		switch choice {
		case "1":
			startPractice(reader)
		case "2":
			listProblems()
		case "3":
			viewStats()
		case "4", "q", "quit", "exit":
			fmt.Println("Thanks for practicing with AlgoScales!")
			return
		default:
			if choice != "" {
				fmt.Println("Invalid choice. Please try again.")
			}
		}

		fmt.Println() // Add a blank line for spacing
	}
}

// startPractice handles the practice workflow in simple UI mode
func startPractice(reader *bufio.Reader) {
	// Set up local signal handling for the practice session
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Function to safely read input with interrupt handling
	readInput := func() string {
		inputCh := make(chan string, 1)
		interruptCh := make(chan bool, 1)

		// Start a goroutine to read input
		go func() {
			text, err := reader.ReadString('\n')
			if err != nil {
				inputCh <- ""
				return
			}
			inputCh <- strings.TrimSpace(text)
		}()

		// Another goroutine to monitor for ctrl+c
		go func() {
			<-sigChan
			interruptCh <- true
		}()

		// Wait for either input or interrupt
		select {
		case input := <-inputCh:
			return input
		case <-interruptCh:
			fmt.Println("\nExiting AlgoScales. Thanks for practicing!")
			os.Exit(0)
			return ""
		}
	}
	
	// Choose practice mode
	fmt.Println("\nSelect practice mode:")
	fmt.Println("1. Learn Mode (with guidance)")
	fmt.Println("2. Practice Mode (test yourself)")
	fmt.Println("3. Cram Mode (rapid fire practice)")
	fmt.Println("4. Back to Main Menu")
	fmt.Print("> ")

	choice := readInput()

	var mode session.Mode
	switch choice {
	case "1":
		mode = session.LearnMode
	case "2":
		mode = session.PracticeMode
	case "3":
		mode = session.CramMode
	case "4", "back", "b":
		return
	default:
		fmt.Println("Invalid choice. Going back to main menu.")
		return
	}

	// Choose programming language
	fmt.Println("\nSelect programming language:")
	fmt.Println("1. Go")
	fmt.Println("2. Python")
	fmt.Println("3. JavaScript")
	fmt.Println("4. Back to Main Menu")
	fmt.Print("> ")

	choice = readInput()

	var language string
	switch choice {
	case "1":
		language = "go"
	case "2":
		language = "python"
	case "3":
		language = "javascript"
	case "4", "back", "b":
		return
	default:
		fmt.Println("Invalid choice. Using Go as default.")
		language = "go"
	}

	// Find a problem
	problems, err := problem.ListAll()
	if err != nil {
		fmt.Printf("Error loading problems: %v\n", err)
		return
	}

	if len(problems) == 0 {
		fmt.Println("No problems available. Please check your installation.")
		return
	}

	// For simplicity, just pick the first problem
	// A complete implementation would allow selection from a list
	selectedProblem := problems[0]

	fmt.Printf("\nSelected problem: %s (%s)\n", selectedProblem.Title, selectedProblem.Difficulty)
	fmt.Printf("Patterns: %s\n", strings.Join(selectedProblem.Patterns, ", "))
	fmt.Println("Press Enter to start the practice session...")
	readInput()

	// Create session options
	opts := session.Options{
		Mode:       mode,
		Language:   language,
		Timer:      45, // Default timer
		ProblemID:  selectedProblem.ID,
	}

	// Set up a goroutine to handle Ctrl+C during the session
	done := make(chan bool, 1)
	go func() {
		// Wait for an interrupt signal
		<-sigChan
		fmt.Println("\nSession interrupted. Returning to main menu...")
		done <- true
	}()

	// Start the session in a goroutine
	go func() {
		if err := session.Start(opts); err != nil {
			fmt.Printf("Error starting session: %v\n", err)
		}
		done <- true
	}()

	// Wait for either completion or interrupt
	<-done
}

// listProblems shows available problems
func listProblems() {
	problems, err := problem.ListAll()
	if err != nil {
		fmt.Printf("Error listing problems: %v\n", err)
		return
	}

	fmt.Println("\nAvailable Problems:")
	for i, p := range problems {
		fmt.Printf("%d. %s (%s) - Patterns: %s\n",
			i+1, p.Title, p.Difficulty, strings.Join(p.Patterns, ", "))
	}
}

// viewStats shows user statistics
func viewStats() {
	// This would use the stats package in a full implementation
	fmt.Println("\nStatistics: (Sample data)")
	fmt.Println("Total Problems Attempted: 10")
	fmt.Println("Total Problems Solved: 7")
	fmt.Println("Current Streak: 3 days")
	fmt.Println("Longest Streak: 5 days")

	fmt.Println("\nPerformance by Pattern:")
	fmt.Println("- Sliding Window: 80% success rate")
	fmt.Println("- Two Pointers: 75% success rate")
	fmt.Println("- Hash Map: 100% success rate")
}