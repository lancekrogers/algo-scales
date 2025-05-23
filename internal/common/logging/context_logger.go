package logging

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ContextKey represents a key for context values
type ContextKey string

const (
	// TraceIDKey is the key for trace IDs in context
	TraceIDKey ContextKey = "trace_id"
	
	// OperationKey is the key for operation names in context
	OperationKey ContextKey = "operation"
	
	// ComponentKey is the key for component names in context
	ComponentKey ContextKey = "component"
)

// Logger provides structured logging with context support
type Logger struct {
	prefix string
}

// NewLogger creates a new context-aware logger
func NewLogger(component string) *Logger {
	return &Logger{
		prefix: fmt.Sprintf("[%s]", component),
	}
}

// WithContext adds context information to the logger
func (l *Logger) WithContext(ctx context.Context) *ContextLogger {
	return &ContextLogger{
		logger: l,
		ctx:    ctx,
	}
}

// ContextLogger provides logging with context information
type ContextLogger struct {
	logger *Logger
	ctx    context.Context
}

// Info logs an info message with context
func (cl *ContextLogger) Info(msg string, args ...interface{}) {
	cl.log("INFO", msg, args...)
}

// Warn logs a warning message with context
func (cl *ContextLogger) Warn(msg string, args ...interface{}) {
	cl.log("WARN", msg, args...)
}

// Error logs an error message with context
func (cl *ContextLogger) Error(msg string, args ...interface{}) {
	cl.log("ERROR", msg, args...)
}

// Debug logs a debug message with context
func (cl *ContextLogger) Debug(msg string, args ...interface{}) {
	cl.log("DEBUG", msg, args...)
}

// StartOperation logs the start of an operation and returns a function to log completion
func (cl *ContextLogger) StartOperation(operation string) func(error) {
	start := time.Now()
	cl.Info("Starting operation: %s", operation)
	
	return func(err error) {
		duration := time.Since(start)
		if err != nil {
			cl.Error("Operation failed: %s (took %v) - %v", operation, duration, err)
		} else {
			cl.Info("Operation completed: %s (took %v)", operation, duration)
		}
	}
}

// log formats and outputs the log message with context information
func (cl *ContextLogger) log(level, msg string, args ...interface{}) {
	// Format the message
	formatted := fmt.Sprintf(msg, args...)
	
	// Extract context values
	var contextInfo string
	if traceID := cl.ctx.Value(TraceIDKey); traceID != nil {
		contextInfo += fmt.Sprintf(" trace=%v", traceID)
	}
	if operation := cl.ctx.Value(OperationKey); operation != nil {
		contextInfo += fmt.Sprintf(" op=%v", operation)
	}
	if component := cl.ctx.Value(ComponentKey); component != nil {
		contextInfo += fmt.Sprintf(" comp=%v", component)
	}
	
	// Log the message
	log.Printf("%s [%s]%s %s", cl.logger.prefix, level, contextInfo, formatted)
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithOperation adds an operation name to the context
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// WithComponent adds a component name to the context
func WithComponent(ctx context.Context, component string) context.Context {
	return context.WithValue(ctx, ComponentKey, component)
}

// Global logger instances
var (
	TestRunnerLogger = NewLogger("TestRunner")
	ProblemLogger    = NewLogger("Problem")
	SessionLogger    = NewLogger("Session")
	FileOpsLogger    = NewLogger("FileOps")
)