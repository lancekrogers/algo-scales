package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// ErrorCategory represents different types of errors in the system
type ErrorCategory string

const (
	TestExecution     ErrorCategory = "test_execution"
	FileOperations    ErrorCategory = "file_operations"
	UIInteraction     ErrorCategory = "ui_interaction"
	EditorIntegration ErrorCategory = "editor_integration"
	NetworkAPI        ErrorCategory = "network_api"
	UserValidation    ErrorCategory = "user_validation"
	SystemSetup       ErrorCategory = "system_setup"
	Unknown           ErrorCategory = "unknown"
)

// ErrorSeverity represents the severity level of errors
type ErrorSeverity string

const (
	SeverityCritical ErrorSeverity = "critical"
	SeverityHigh     ErrorSeverity = "high" 
	SeverityMedium   ErrorSeverity = "medium"
	SeverityLow      ErrorSeverity = "low"
	SeverityInfo     ErrorSeverity = "info"
)

// SessionSnapshot captures the current session state
type SessionSnapshot struct {
	ProblemID     string            `json:"problem_id"`
	Language      string            `json:"language"`
	Mode          string            `json:"mode"`
	UserCode      string            `json:"user_code,omitempty"` // omitempty for large code
	TestCase      string            `json:"test_case,omitempty"`
	StartTime     time.Time         `json:"start_time"`
	HintsUsed     bool              `json:"hints_used"`
	SolutionUsed  bool              `json:"solution_used"`
	Patterns      []string          `json:"patterns"`
	Difficulty    string            `json:"difficulty"`
	Workspace     string            `json:"workspace"`
	CodeFile      string            `json:"code_file"`
	CustomFields  map[string]string `json:"custom_fields,omitempty"`
}

// SystemSnapshot captures current system state
type SystemSnapshot struct {
	OS               string            `json:"os"`
	Arch             string            `json:"arch"`
	GoVersion        string            `json:"go_version"`
	TerminalType     string            `json:"terminal_type,omitempty"`
	WorkingDirectory string            `json:"working_directory"`
	DiskSpace        int64             `json:"disk_space_mb"`
	MemoryUsage      int64             `json:"memory_usage_mb"`
	ProcessID        int               `json:"process_id"`
	Environment      map[string]string `json:"environment,omitempty"`
	LanguageVersions map[string]string `json:"language_versions,omitempty"`
	Timestamp        time.Time         `json:"timestamp"`
}

// ErrorContext contains comprehensive error information
type ErrorContext struct {
	ID            string           `json:"id"`
	Category      ErrorCategory    `json:"category"`
	Severity      ErrorSeverity    `json:"severity"`
	TraceID       string           `json:"trace_id,omitempty"`
	UserAction    string           `json:"user_action"`
	ErrorMessage  string           `json:"error_message"`
	StackTrace    string           `json:"stack_trace,omitempty"`
	SessionState  *SessionSnapshot `json:"session_state,omitempty"`
	SystemState   *SystemSnapshot  `json:"system_state,omitempty"`
	Timestamp     time.Time        `json:"timestamp"`
	RelatedErrors []string         `json:"related_errors,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// CentralErrorLogger handles all structured error logging
type CentralErrorLogger struct {
	logFile    *os.File
	errorCount map[ErrorCategory]int
}

// NewCentralErrorLogger creates a new centralized error logger
func NewCentralErrorLogger(logPath string) (*CentralErrorLogger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	
	// Open error log file
	logFile := fmt.Sprintf("%s/errors_%s.log", logPath, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open error log file: %w", err)
	}
	
	return &CentralErrorLogger{
		logFile:    file,
		errorCount: make(map[ErrorCategory]int),
	}, nil
}

// LogError logs a structured error with full context
func (cel *CentralErrorLogger) LogError(ctx context.Context, category ErrorCategory, severity ErrorSeverity, err error, userAction string, sessionState *SessionSnapshot) string {
	errorID := generateErrorID()
	
	errorCtx := ErrorContext{
		ID:           errorID,
		Category:     category,
		Severity:     severity,
		TraceID:      getTraceID(ctx),
		UserAction:   userAction,
		ErrorMessage: err.Error(),
		StackTrace:   captureStackTrace(),
		SessionState: sessionState,
		SystemState:  captureSystemSnapshot(),
		Timestamp:    time.Now(),
		Tags:         extractTagsFromContext(ctx),
		Metadata:     make(map[string]interface{}),
	}
	
	// Log to file and console
	cel.writeErrorLog(errorCtx)
	cel.errorCount[category]++
	
	// Also log to context logger for immediate visibility
	logger := NewLogger("CentralErrorLogger").WithContext(ctx)
	logger.Error("Error logged - ID: %s, Category: %s, Action: %s - %v", errorID, category, userAction, err)
	
	return errorID
}

// LogPanic logs panic recovery with full context
func (cel *CentralErrorLogger) LogPanic(ctx context.Context, recovered interface{}, userAction string, sessionState *SessionSnapshot) string {
	errorID := generateErrorID()
	
	errorCtx := ErrorContext{
		ID:           errorID,
		Category:     Unknown, // Panics are often unexpected
		Severity:     SeverityCritical,
		TraceID:      getTraceID(ctx),
		UserAction:   userAction,
		ErrorMessage: fmt.Sprintf("PANIC: %v", recovered),
		StackTrace:   string(debug.Stack()),
		SessionState: sessionState,
		SystemState:  captureSystemSnapshot(),
		Timestamp:    time.Now(),
		Tags:         extractTagsFromContext(ctx),
		Metadata: map[string]interface{}{
			"panic_value": recovered,
			"panic_type":  fmt.Sprintf("%T", recovered),
		},
	}
	
	cel.writeErrorLog(errorCtx)
	cel.errorCount[Unknown]++
	
	// Critical logging to console
	logger := NewLogger("PanicRecovery").WithContext(ctx)
	logger.Error("PANIC RECOVERED - ID: %s, Action: %s - %v", errorID, userAction, recovered)
	
	return errorID
}

// LogTestExecutionError logs test execution failures with enhanced context
func (cel *CentralErrorLogger) LogTestExecutionError(ctx context.Context, err error, language string, userCode string, testCase string, sessionState *SessionSnapshot) string {
	// Enhance session state with test-specific information
	if sessionState != nil {
		sessionState.Language = language
		sessionState.UserCode = userCode
		sessionState.TestCase = testCase
	}
	
	errorID := cel.LogError(ctx, TestExecution, SeverityHigh, err, "execute_tests", sessionState)
	
	// Add test-specific metadata
	cel.addMetadata(errorID, map[string]interface{}{
		"language":        language,
		"code_length":     len(userCode),
		"test_case":       testCase,
		"has_syntax_error": strings.Contains(err.Error(), "syntax"),
		"has_timeout":      strings.Contains(err.Error(), "timeout"),
	})
	
	return errorID
}

// LogFileOperationError logs file operation errors with system context
func (cel *CentralErrorLogger) LogFileOperationError(ctx context.Context, err error, operation string, filePath string, sessionState *SessionSnapshot) string {
	errorID := cel.LogError(ctx, FileOperations, SeverityMedium, err, fmt.Sprintf("file_%s", operation), sessionState)
	
	// Add file-specific metadata
	fileInfo := captureFileSystemContext(filePath)
	cel.addMetadata(errorID, map[string]interface{}{
		"operation":      operation,
		"file_path":      filePath,
		"file_exists":    fileInfo["exists"],
		"file_size":      fileInfo["size"],
		"permissions":    fileInfo["permissions"],
		"disk_space_mb":  fileInfo["disk_space_mb"],
	})
	
	return errorID
}

// LogUIError logs UI interaction errors
func (cel *CentralErrorLogger) LogUIError(ctx context.Context, err error, screen string, action string, sessionState *SessionSnapshot) string {
	errorID := cel.LogError(ctx, UIInteraction, SeverityMedium, err, fmt.Sprintf("ui_%s_%s", screen, action), sessionState)
	
	cel.addMetadata(errorID, map[string]interface{}{
		"screen":         screen,
		"ui_action":      action,
		"terminal_type":  os.Getenv("TERM"),
		"terminal_size":  getTerminalSize(),
	})
	
	return errorID
}

// LogEditorError logs editor integration errors
func (cel *CentralErrorLogger) LogEditorError(ctx context.Context, err error, editor string, filePath string, sessionState *SessionSnapshot) string {
	errorID := cel.LogError(ctx, EditorIntegration, SeverityHigh, err, "open_editor", sessionState)
	
	cel.addMetadata(errorID, map[string]interface{}{
		"editor":       editor,
		"file_path":    filePath,
		"editor_env":   os.Getenv("EDITOR"),
		"visual_env":   os.Getenv("VISUAL"),
		"has_display":  os.Getenv("DISPLAY") != "",
	})
	
	return errorID
}

// GetErrorStats returns error statistics
func (cel *CentralErrorLogger) GetErrorStats() map[ErrorCategory]int {
	return cel.errorCount
}

// Close closes the error logger
func (cel *CentralErrorLogger) Close() error {
	if cel.logFile != nil {
		return cel.logFile.Close()
	}
	return nil
}

// writeErrorLog writes the error context to the log file
func (cel *CentralErrorLogger) writeErrorLog(errorCtx ErrorContext) {
	jsonData, err := json.Marshal(errorCtx)
	if err != nil {
		log.Printf("Failed to marshal error context: %v", err)
		return
	}
	
	// Write to file
	if cel.logFile != nil {
		cel.logFile.WriteString(string(jsonData) + "\n")
		cel.logFile.Sync()
	}
	
	// Also write to console in development
	if os.Getenv("ALGO_SCALES_DEBUG") == "true" {
		log.Printf("ERROR[%s]: %s", errorCtx.Category, errorCtx.ErrorMessage)
	}
}

// addMetadata adds metadata to an existing error
func (cel *CentralErrorLogger) addMetadata(errorID string, metadata map[string]interface{}) {
	// In a more sophisticated implementation, this would update the stored error
	// For now, we'll log the additional metadata
	jsonData, _ := json.Marshal(metadata)
	if cel.logFile != nil {
		cel.logFile.WriteString(fmt.Sprintf(`{"error_id":"%s","additional_metadata":%s}`, errorID, string(jsonData)) + "\n")
		cel.logFile.Sync()
	}
}

// Helper functions
func generateErrorID() string {
	return fmt.Sprintf("err_%d_%d", time.Now().UnixNano(), os.Getpid())
}

func getTraceID(ctx context.Context) string {
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		return fmt.Sprintf("%v", traceID)
	}
	return ""
}

func captureStackTrace() string {
	buf := make([]byte, 1024*4)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func captureSystemSnapshot() *SystemSnapshot {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	
	wd, _ := os.Getwd()
	
	return &SystemSnapshot{
		OS:               runtime.GOOS,
		Arch:             runtime.GOARCH,
		GoVersion:        runtime.Version(),
		TerminalType:     os.Getenv("TERM"),
		WorkingDirectory: wd,
		DiskSpace:        getDiskSpace(),
		MemoryUsage:      int64(m.Alloc / 1024 / 1024), // MB
		ProcessID:        os.Getpid(),
		Environment:      captureRelevantEnv(),
		LanguageVersions: captureLanguageVersions(),
		Timestamp:        time.Now(),
	}
}

func extractTagsFromContext(ctx context.Context) map[string]string {
	tags := make(map[string]string)
	
	if op := ctx.Value(OperationKey); op != nil {
		tags["operation"] = fmt.Sprintf("%v", op)
	}
	if comp := ctx.Value(ComponentKey); comp != nil {
		tags["component"] = fmt.Sprintf("%v", comp)
	}
	
	return tags
}

func captureFileSystemContext(filePath string) map[string]interface{} {
	context := make(map[string]interface{})
	
	if stat, err := os.Stat(filePath); err == nil {
		context["exists"] = true
		context["size"] = stat.Size()
		context["permissions"] = stat.Mode().String()
	} else {
		context["exists"] = false
		context["error"] = err.Error()
	}
	
	context["disk_space_mb"] = getDiskSpace()
	return context
}

func captureRelevantEnv() map[string]string {
	relevantVars := []string{
		"TERM", "EDITOR", "VISUAL", "DISPLAY", "PATH", 
		"GOPATH", "GOROOT", "NODE_PATH", "PYTHON_PATH",
		"ALGO_SCALES_DEBUG", "ALGO_SCALES_LOG_LEVEL",
	}
	
	env := make(map[string]string)
	for _, v := range relevantVars {
		if val := os.Getenv(v); val != "" {
			env[v] = val
		}
	}
	return env
}

func captureLanguageVersions() map[string]string {
	versions := make(map[string]string)
	
	// This would be enhanced to actually check language versions
	// For now, we'll add placeholders
	versions["go"] = runtime.Version()
	// TODO: Add actual version detection for Python, Node.js, etc.
	
	return versions
}

func getDiskSpace() int64 {
	// Simplified disk space check - would need platform-specific implementation
	return 0 // Placeholder
}

func getTerminalSize() string {
	// Simplified terminal size - would use terminal packages for real implementation
	return fmt.Sprintf("%sx%s", os.Getenv("COLUMNS"), os.Getenv("LINES"))
}

// Global error logger instance
var GlobalErrorLogger *CentralErrorLogger

// InitializeGlobalErrorLogger initializes the global error logger
func InitializeGlobalErrorLogger(logPath string) error {
	var err error
	GlobalErrorLogger, err = NewCentralErrorLogger(logPath)
	return err
}