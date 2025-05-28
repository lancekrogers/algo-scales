package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
)

// GlobalErrorHandler provides application-wide error handling and recovery
type GlobalErrorHandler struct {
	logger *CentralErrorLogger
	appCtx context.Context
}

// NewGlobalErrorHandler creates a new global error handler
func NewGlobalErrorHandler(ctx context.Context, logPath string) (*GlobalErrorHandler, error) {
	logger, err := NewCentralErrorLogger(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create central error logger: %w", err)
	}
	
	// Set global logger instance
	GlobalErrorLogger = logger
	
	return &GlobalErrorHandler{
		logger: logger,
		appCtx: ctx,
	}, nil
}

// SetupGlobalRecovery sets up global panic recovery for the entire application
func (geh *GlobalErrorHandler) SetupGlobalRecovery() {
	// Capture any panics that occur during startup or in goroutines
	defer func() {
		if r := recover(); r != nil {
			geh.HandlePanic(r, "application_startup", nil)
			
			// Log to stderr as well since we might not have full logging set up
			fmt.Fprintf(os.Stderr, "CRITICAL: Application panic during startup: %v\n", r)
			fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", debug.Stack())
			
			// Exit gracefully
			os.Exit(1)
		}
	}()
}

// HandlePanic handles panic recovery with comprehensive logging
func (geh *GlobalErrorHandler) HandlePanic(recovered interface{}, userAction string, sessionState *SessionSnapshot) string {
	// Create context for panic logging
	ctx := WithOperation(geh.appCtx, "panic_recovery")
	ctx = WithComponent(ctx, "GlobalErrorHandler")
	
	errorID := geh.logger.LogPanic(ctx, recovered, userAction, sessionState)
	
	// Log to console for immediate visibility
	fmt.Fprintf(os.Stderr, "PANIC RECOVERED [ID: %s]: %v\n", errorID, recovered)
	fmt.Fprintf(os.Stderr, "Action: %s\n", userAction)
	if sessionState != nil {
		fmt.Fprintf(os.Stderr, "Problem: %s, Language: %s\n", sessionState.ProblemID, sessionState.Language)
	}
	
	return errorID
}

// WrapMainFunction wraps the main application function with global error handling
func (geh *GlobalErrorHandler) WrapMainFunction(mainFunc func() error) error {
	defer func() {
		if r := recover(); r != nil {
			geh.HandlePanic(r, "main_application", nil)
			
			// Exit after logging
			os.Exit(1)
		}
	}()
	
	return mainFunc()
}

// WrapUIFunction wraps UI functions with error handling and session context
func (geh *GlobalErrorHandler) WrapUIFunction(uiFunc func() error, sessionState *SessionSnapshot) error {
	defer func() {
		if r := recover(); r != nil {
			geh.HandlePanic(r, "ui_interaction", sessionState)
		}
	}()
	
	return uiFunc()
}

// WrapTestFunction wraps test execution functions with enhanced error context
func (geh *GlobalErrorHandler) WrapTestFunction(testFunc func() error, language string, code string, sessionState *SessionSnapshot) error {
	defer func() {
		if r := recover(); r != nil {
			// Enhance session state with test-specific info
			if sessionState != nil {
				sessionState.Language = language
				sessionState.UserCode = code
			}
			geh.HandlePanic(r, fmt.Sprintf("test_execution_%s", language), sessionState)
		}
	}()
	
	return testFunc()
}

// LogCriticalError logs critical errors that should be investigated immediately
func (geh *GlobalErrorHandler) LogCriticalError(err error, context string, sessionState *SessionSnapshot) string {
	ctx := WithOperation(geh.appCtx, "critical_error")
	ctx = WithComponent(ctx, "GlobalErrorHandler")
	
	return geh.logger.LogError(ctx, Unknown, SeverityCritical, err, context, sessionState)
}

// LogUserFacingError logs errors that directly impact user experience
func (geh *GlobalErrorHandler) LogUserFacingError(err error, userAction string, sessionState *SessionSnapshot) string {
	ctx := WithOperation(geh.appCtx, "user_facing_error")
	ctx = WithComponent(ctx, "GlobalErrorHandler")
	
	// Determine category based on user action
	category := Unknown
	switch {
	case userAction == "execute_tests" || userAction == "run_tests":
		category = TestExecution
	case userAction == "open_editor" || userAction == "save_code":
		category = EditorIntegration
	case userAction == "load_problem" || userAction == "save_session":
		category = FileOperations
	case userAction == "ui_interaction":
		category = UIInteraction
	}
	
	return geh.logger.LogError(ctx, category, SeverityHigh, err, userAction, sessionState)
}

// GetErrorStats returns current error statistics
func (geh *GlobalErrorHandler) GetErrorStats() map[ErrorCategory]int {
	return geh.logger.GetErrorStats()
}

// Close gracefully closes the global error handler
func (geh *GlobalErrorHandler) Close() error {
	if geh.logger != nil {
		return geh.logger.Close()
	}
	return nil
}

// InitializeGlobalErrorHandling initializes global error handling for the application
func InitializeGlobalErrorHandling(ctx context.Context) (*GlobalErrorHandler, error) {
	// Create logs directory in user's config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to current directory
		configDir = "."
	}
	
	logPath := filepath.Join(configDir, "algo-scales", "logs")
	
	handler, err := NewGlobalErrorHandler(ctx, logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize global error handling: %w", err)
	}
	
	// Set up global recovery
	handler.SetupGlobalRecovery()
	
	// Only log initialization message if not in vim mode or JSON output mode
	// Check if we're running with --vim-mode flag or in a context that requires pure JSON output
	isVimMode := false
	isJSONOutput := false
	for _, arg := range os.Args {
		if arg == "--vim-mode" || arg == "--vim" {
			isVimMode = true
			break
		}
		if arg == "--json" {
			isJSONOutput = true
			break
		}
	}
	
	// Don't pollute stdout with log messages in vim/JSON modes
	if !isVimMode && !isJSONOutput {
		log.Printf("Global error handling initialized - logs at: %s", logPath)
	}
	
	return handler, nil
}

// Global instance
var GlobalHandler *GlobalErrorHandler