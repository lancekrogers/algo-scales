// Package ui implements the terminal user interface for Algo Scales
package ui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/ui/controller"
	"github.com/lancekrogers/algo-scales/internal/ui/model"
	"github.com/lancekrogers/algo-scales/internal/ui/view"
)

// App represents the UI application
type App struct {
	model      *model.UIModel
	view       *view.View
	controller *controller.Controller
}

// NewApp creates a new UI application
func NewApp() *App {
	// Create model
	m := model.NewModel()

	// Create view and controller
	v := view.NewView(&m)
	c := controller.NewController(&m)

	return &App{
		model:      &m,
		view:       v,
		controller: c,
	}
}

// Start launches the UI application
func (a *App) Start() error {
	// Check if debugging is enabled
	debug := false
	if os.Getenv("DEBUG") == "1" {
		debug = true
	}

	// Setup for terminal compatibility
	// Use a safer set of options that work in most terminals
	opts := []tea.ProgramOption{}

	// Check if TTY is available before trying to use terminal features
	// This helps prevent freezing when running in non-interactive environments
	if isTTYAvailable() {
		// Check terminal capability explicitly to help prevent race conditions
		term := os.Getenv("TERM")
		termProgram := os.Getenv("TERM_PROGRAM")

		// Terminal capability detection with specific handling for different terminal types
		if termProgram == "iTerm.app" {
			// iTerm2 specific settings - should be more reliable
			opts = append(opts, tea.WithAltScreen())
			opts = append(opts, tea.WithMouseCellMotion())
		} else if term == "xterm-256color" || term == "screen-256color" {
			// Standard terminal settings
			opts = append(opts, tea.WithAltScreen())
		}
	} else {
		// For non-TTY environments, return an error immediately
		// This is better than freezing or hanging
		err := fmt.Errorf("no TTY available for interactive UI")
		if debug {
			fmt.Println("Debug mode: UI error:", err)
		}
		return fmt.Errorf("terminal UI error: %v\n\nTry running with --cli flag for command-line mode", err)
	}

	// Always use standard output for debug mode
	if debug {
		// Use standard renderer that doesn't require all TTY capabilities
		opts = append(opts, tea.WithoutRenderer())
	}

	// Setup a channel to detect if tea initialization is hanging
	initDone := make(chan struct{})
	go func() {
		// This will unblock shortly after initialization starts
		// If the UI is successfully initialized
		time.Sleep(500 * time.Millisecond)
		select {
		case <-initDone:
			// Already completed, do nothing
		default:
			// Force terminal reset to prevent freezing
			// This helps recover from potential lock-ups
			fmt.Print("\033c") // Terminal reset escape code
		}
	}()

	// Create the Bubble Tea program with options and a quit function
	p := tea.NewProgram(a, opts...)

	// Setup a channel for safe communication between goroutines
	done := make(chan struct{})
	var err error

	// Run the program in a goroutine to prevent blocking
	go func() {
		close(initDone) // Signal initialization has started
		_, err = p.Run()
		close(done)
	}()

	// Wait for the program to complete with a timeout
	select {
	case <-done:
		// Normal completion
		if err != nil {
			if debug {
				fmt.Println("Debug mode: UI error:", err)
			}
			return fmt.Errorf("terminal UI error: %v\n\nTry running with --cli flag for command-line mode", err)
		}
		return nil
	case <-time.After(2 * time.Second):
		// If initialization is taking too long, try to gracefully terminate
		fmt.Print("\033c") // Reset terminal
		p.Quit()           // Tell bubbletea to quit
		return fmt.Errorf("UI initialization timed out. Try running with --simple or --split flag")
	}
}

// isTTYAvailable checks if a TTY device is available for interactive UI
func isTTYAvailable() bool {
	// Check if we can open /dev/tty
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return false
	}
	tty.Close()

	// Additional checks for terminal capabilities
	// Check for environment variables that indicate we're running in a CI system
	if os.Getenv("CI") != "" || os.Getenv("CONTINUOUS_INTEGRATION") != "" {
		return false
	}

	// Check if TERM is set to something useful
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}

	return true
}

// Implement tea.Model interface for the App

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return a.controller.Initialize()
}

// Update handles updates to the application
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := a.controller.Update(msg)
	// Cast the returned model back to our App type or return self if nil
	if model == nil {
		return a, cmd
	}
	// We only need the updated UIModel state from the controller
	// but we need to return the App as our tea.Model
	return a, cmd
}

// View renders the application
func (a *App) View() string {
	return a.view.Render()
}

// Messages are defined in the model/messages.go file
// and are used for communication between components

