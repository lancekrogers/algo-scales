package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// StartTUI starts the terminal user interface
func StartTUI() error {
	// Check if debugging is enabled
	debug := false
	if os.Getenv("DEBUG") == "1" {
		debug = true
	}

	// Create the model
	model := NewModel()

	// Setup program options
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}

	// Add debug options if needed
	if debug {
		opts = append(opts, tea.WithoutSignals())
	}

	// Create and run the program
	p := tea.NewProgram(model, opts...)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}