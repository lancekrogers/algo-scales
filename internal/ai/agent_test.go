package ai

import (
	"testing"
)

func TestNewAgent(t *testing.T) {
	// Create test config
	config := &Config{
		Claude: &ClaudeConfig{
			CLIPath: "claude",
		},
		Ollama: &OllamaConfig{
			Host:  "http://localhost:11434",
			Model: "llama3",
		},
	}

	tests := []struct {
		name     string
		provider Provider
		wantErr  bool
	}{
		{
			name:     "Claude provider",
			provider: ProviderClaude,
			wantErr:  false,
		},
		{
			name:     "Ollama provider",
			provider: ProviderOllama,
			wantErr:  false,
		},
		{
			name:     "Unknown provider",
			provider: Provider("unknown"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewAgent(tt.provider, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && agent == nil {
				t.Error("Expected agent to be created")
			}
		})
	}
}

func TestNewAgentMissingConfig(t *testing.T) {
	// Test with missing provider configs
	tests := []struct {
		name     string
		provider Provider
		config   *Config
	}{
		{
			name:     "Missing Claude config",
			provider: ProviderClaude,
			config:   &Config{}, // No Claude config
		},
		{
			name:     "Missing Ollama config",
			provider: ProviderOllama,
			config:   &Config{}, // No Ollama config
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAgent(tt.provider, tt.config)
			if err == nil {
				t.Error("Expected error for missing config")
			}
		})
	}
}