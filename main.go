// Main entry point for the application

package main

import (
	"fmt"
	"os"

	"github.com/lancekrogers/algo-scales/cmd"
	"github.com/lancekrogers/algo-scales/internal/ui"
)

func main() {
	// Check if CLI mode is explicitly requested via flag
	useLegacyCLI := false
	for _, arg := range os.Args {
		if arg == "--cli" || arg == "--legacy" {
			useLegacyCLI = true
			break
		}
	}

	// Check if this is being run from Neovim via the vim plugin
	fromVim := false
	for _, arg := range os.Args {
		if arg == "--vim-mode" {
			fromVim = true
			break
		}
	}

	if useLegacyCLI || fromVim {
		// Use traditional CLI for legacy mode or Neovim integration
		cmd.Execute()
	} else {
		// Launch the enhanced interactive TUI by default
		app := ui.NewApp()
		if err := app.Start(); err != nil {
			fmt.Printf("Error launching UI: %v\n", err)
			os.Exit(1)
		}
	}
}
