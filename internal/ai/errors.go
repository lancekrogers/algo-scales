package ai

import (
	"context"
	"errors"
	"fmt"
)

// Common errors
var (
	// ErrProviderNotConfigured indicates the AI provider is not properly configured
	ErrProviderNotConfigured = errors.New("AI provider not configured")

	// ErrNoAPIKey indicates the API key is missing
	ErrNoAPIKey = errors.New("API key not found")

	// ErrInvalidProvider indicates an unknown provider was specified
	ErrInvalidProvider = errors.New("invalid AI provider")

	// ErrConnectionFailed indicates a connection error to the AI service
	ErrConnectionFailed = errors.New("failed to connect to AI service")

	// ErrRateLimited indicates the API rate limit was exceeded
	ErrRateLimited = errors.New("rate limit exceeded")

	// ErrContextTooLong indicates the context exceeds the model's limit
	ErrContextTooLong = errors.New("context exceeds model limit")

	// ErrInvalidResponse indicates the AI returned an invalid response
	ErrInvalidResponse = errors.New("invalid response from AI service")

	// ErrSessionNotFound indicates a session ID was not found
	ErrSessionNotFound = errors.New("session not found")

	// ErrClaudeNotInstalled indicates Claude Code CLI is not installed
	ErrClaudeNotInstalled = errors.New("Claude Code CLI not installed or not in PATH")

	// ErrOllamaNotRunning indicates Ollama server is not running
	ErrOllamaNotRunning = errors.New("Ollama server is not running")
)

// ConfigError represents a configuration-related error
type ConfigError struct {
	Field   string
	Message string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("config error in %s: %s", e.Field, e.Message)
}

// APIError represents an error from an AI API
type APIError struct {
	Provider   string
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("%s API error (status %d): %s", e.Provider, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("%s API error: %s", e.Provider, e.Message)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s (value: %v): %s", e.Field, e.Value, e.Message)
}

// WrapError wraps an error with additional context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// IsRetryable determines if an error is retryable
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific retryable errors
	switch {
	case errors.Is(err, ErrConnectionFailed):
		return true
	case errors.Is(err, ErrRateLimited):
		return true
	case errors.Is(err, context.DeadlineExceeded):
		return true
	}

	// Check for API errors with retryable status codes
	var apiErr APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case 429, 502, 503, 504: // Rate limit, bad gateway, service unavailable, timeout
			return true
		}
	}

	return false
}

// HandleProviderError provides user-friendly error messages and solutions
func HandleProviderError(provider Provider, err error) string {
	if err == nil {
		return ""
	}

	switch provider {
	case ProviderClaude:
		switch {
		case errors.Is(err, ErrClaudeNotInstalled):
			return "Claude Code CLI is not installed. Please install it following the instructions at https://docs.anthropic.com/en/docs/claude-code/getting-started"
		case errors.Is(err, ErrProviderNotConfigured):
			return "Claude is not configured. Run 'algo-scales ai config' to set up Claude."
		case errors.Is(err, ErrRateLimited):
			return "Claude rate limit exceeded. Please wait a moment before trying again."
		}

	case ProviderOllama:
		switch {
		case errors.Is(err, ErrOllamaNotRunning):
			return "Ollama server is not running. Start it with 'ollama serve' in another terminal."
		case errors.Is(err, ErrConnectionFailed):
			return "Cannot connect to Ollama. Make sure Ollama is running on http://localhost:11434"
		case errors.Is(err, ErrProviderNotConfigured):
			return "Ollama is not configured. Run 'algo-scales ai config' to set up Ollama."
		}
	}

	// Generic error handling
	return fmt.Sprintf("AI Error (%s): %v", provider, err)
}