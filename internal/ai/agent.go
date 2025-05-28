// Package ai provides AI assistant functionality for AlgoScales
package ai

import (
	"context"
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Agent defines the interface for AI providers
type Agent interface {
	// Chat sends a message and streams the response
	Chat(ctx context.Context, messages []Message, opts ChatOptions) (<-chan ChatResponse, error)

	// GetHint generates a hint for a problem
	GetHint(ctx context.Context, problem problem.Problem, userCode string, level int) (<-chan string, error)

	// ReviewCode provides feedback on submitted code
	ReviewCode(ctx context.Context, problem problem.Problem, code string) (<-chan string, error)

	// ExplainPattern explains an algorithm pattern
	ExplainPattern(ctx context.Context, pattern string, examples []problem.Problem) (<-chan string, error)
}

// Message represents a chat message
type Message struct {
	Role    string // "user", "assistant", "system"
	Content string
}

// ChatOptions configures chat behavior
type ChatOptions struct {
	Temperature float64
	MaxTokens   int
	Stream      bool
}

// ChatResponse represents a streaming response from the AI
type ChatResponse struct {
	Content   string
	Done      bool
	Error     error
	SessionID string  // For session continuation
	Cost      float64 // Cost in USD (if applicable)
}

// Provider represents the available AI providers
type Provider string

const (
	ProviderClaude Provider = "claude"
	ProviderOllama Provider = "ollama"
)

// NewAgent creates a new AI agent based on the configuration
func NewAgent(provider Provider, config *Config) (Agent, error) {
	switch provider {
	case ProviderClaude:
		if config.Claude == nil {
			return nil, fmt.Errorf("claude configuration not found")
		}
		return NewClaudeProvider(*config.Claude)
	case ProviderOllama:
		if config.Ollama == nil {
			return nil, fmt.Errorf("ollama configuration not found")
		}
		return NewOllamaProvider(*config.Ollama)
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

// GetDefaultAgent returns an agent using the default provider from config
func GetDefaultAgent() (Agent, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	provider := Provider(config.DefaultProvider)
	if provider == "" {
		provider = ProviderClaude
	}

	return NewAgent(provider, config)
}