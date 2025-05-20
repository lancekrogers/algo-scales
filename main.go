// Main entry point for the application

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lancekrogers/algo-scales/cmd"
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

	// Let Cobra handle the command-line parsing and execution
	cmd.Execute()
	
	// Signal handler cleanup
	close(stopSignalHandler)
}