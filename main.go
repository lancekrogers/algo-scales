// Main entry point for the application

package main

import (
	"fmt"
	"os"

	"github.com/lancekrogers/algo-scales/cmd"
	"github.com/lancekrogers/algo-scales/internal/ui"
)

func main() {
	// Check if UI flag is provided
	useUI := false
	for _, arg := range os.Args {
		if arg == "--ui" || arg == "-u" {
			useUI = true
			break
		}
	}

	if useUI {
		// Launch the interactive TUI
		app := ui.NewApp()
		if err := app.Start(); err != nil {
			fmt.Printf("Error launching UI: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Use traditional CLI
		cmd.Execute()
	}
}
