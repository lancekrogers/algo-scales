// Main entry point for the application

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lancekrogers/algo-scales/cmd"
	"github.com/lancekrogers/algo-scales/internal/common/logging"
)

func main() {
	// Initialize global error handling first
	ctx := context.Background()
	globalHandler, err := logging.InitializeGlobalErrorHandling(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize error handling: %v", err)
	}
	defer globalHandler.Close()
	
	// Store global handler reference
	logging.GlobalHandler = globalHandler
	
	// Wrap main execution with global error handling
	err = globalHandler.WrapMainFunction(func() error {
		// Set up global signal handling for Ctrl+C
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		
		// Create a channel that can be closed to stop the goroutine
		stopSignalHandler := make(chan struct{})
		
		// Handle Ctrl+C in a separate goroutine for all modes
		go func() {
			defer func() {
				if r := recover(); r != nil {
					globalHandler.HandlePanic(r, "signal_handler", nil)
				}
			}()
			
			select {
			case <-sigChan:
				fmt.Println("\nExiting AlgoScales. Thanks for practicing!")
				
				// Log graceful shutdown
				logger := logging.NewLogger("Main").WithContext(ctx)
				logger.Info("Application shutdown initiated by user signal")
				
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
		
		return nil
	})
	
	if err != nil {
		globalHandler.LogCriticalError(err, "main_execution", nil)
		log.Fatalf("Application execution failed: %v", err)
	}
}