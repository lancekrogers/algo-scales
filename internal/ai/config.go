package ai

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the AI assistant configuration
type Config struct {
	Version         string         `yaml:"version"`
	DefaultProvider string         `yaml:"default_provider"`
	Claude          *ClaudeConfig  `yaml:"claude,omitempty"`
	Ollama          *OllamaConfig  `yaml:"ollama,omitempty"`
	Prompts         *PromptConfig  `yaml:"prompts,omitempty"`
	Features        *FeatureConfig `yaml:"features,omitempty"`
	Logging         *LoggingConfig `yaml:"logging,omitempty"`
}

// ClaudeConfig configures the Claude Code integration
type ClaudeConfig struct {
	CLIPath         string            `yaml:"cli_path"`
	DefaultFormat   string            `yaml:"default_format"`
	SaveSessions    bool              `yaml:"save_sessions"`
	SessionDir      string            `yaml:"session_directory"`
	MaxTurns        int               `yaml:"max_turns"`
	Verbose         bool              `yaml:"verbose"`
	MCP             *MCPConfig        `yaml:"mcp,omitempty"`
	AllowedTools    []string          `yaml:"allowed_tools"`
	DisallowedTools []string          `yaml:"disallowed_tools"`
}

// MCPConfig represents Model Context Protocol configuration
type MCPConfig struct {
	Enabled    bool                       `yaml:"enabled"`
	ConfigFile string                     `yaml:"config_file"`
	Servers    map[string]MCPServerConfig `yaml:"servers"`
}

// MCPServerConfig represents an MCP server configuration
type MCPServerConfig struct {
	Command string            `yaml:"command"`
	Args    []string          `yaml:"args"`
	Env     map[string]string `yaml:"env,omitempty"`
}

// OllamaConfig configures Ollama integration
type OllamaConfig struct {
	Host        string  `yaml:"host"`
	Model       string  `yaml:"model"`
	Timeout     int     `yaml:"timeout"`
	NumCtx      int     `yaml:"num_ctx"`
	Temperature float64 `yaml:"temperature"`
}

// PromptConfig contains prompt templates
type PromptConfig struct {
	SystemPrefix string `yaml:"system_prefix"`
	HintTemplate string `yaml:"hint_template"`
}

// FeatureConfig toggles features
type FeatureConfig struct {
	CodeReview          bool `yaml:"code_review"`
	PatternExplanations bool `yaml:"pattern_explanations"`
	InteractiveREPL     bool `yaml:"interactive_repl"`
	AutoReview          bool `yaml:"auto_review"`
}

// LoggingConfig configures logging
type LoggingConfig struct {
	Level           string `yaml:"level"`
	LogInteractions bool   `yaml:"log_interactions"`
	LogFile         string `yaml:"log_file"`
}

// LoadConfig loads the AI configuration from file
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".algo-scales", "ai-config.yaml")

	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig(configPath)
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Expand paths
	config.expandPaths(homeDir)

	return &config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".algo-scales", "ai-config.yaml")
	configDir := filepath.Dir(configPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// createDefaultConfig creates and saves a default configuration
func createDefaultConfig(configPath string) (*Config, error) {
	config := &Config{
		Version:         "1.0",
		DefaultProvider: "claude",
		Claude: &ClaudeConfig{
			CLIPath:       "claude",
			DefaultFormat: "json",
			SaveSessions:  true,
			SessionDir:    "~/.algo-scales/claude-sessions",
			MaxTurns:      5,
			Verbose:       false,
			MCP: &MCPConfig{
				Enabled: true,
				Servers: map[string]MCPServerConfig{
					"filesystem": {
						Command: "npx",
						Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "./"},
					},
				},
			},
			AllowedTools: []string{
				"mcp__filesystem__read_file",
				"mcp__filesystem__list_directory",
				"Bash",
				"Read",
			},
			DisallowedTools: []string{
				"mcp__filesystem__write_file",
				"mcp__filesystem__delete_file",
			},
		},
		Ollama: &OllamaConfig{
			Host:        "http://localhost:11434",
			Model:       "llama3.2:latest",
			Timeout:     300,
			NumCtx:      4096,
			Temperature: 0.7,
		},
		Prompts: &PromptConfig{
			SystemPrefix: "You are an expert algorithm tutor helping students learn data structures and algorithms. Focus on teaching concepts and patterns rather than just providing solutions.",
		},
		Features: &FeatureConfig{
			CodeReview:          true,
			PatternExplanations: true,
			InteractiveREPL:     true,
			AutoReview:          false,
		},
		Logging: &LoggingConfig{
			Level:           "info",
			LogInteractions: false,
			LogFile:         "~/.algo-scales/ai-assistant.log",
		},
	}

	// Create config directory
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save the default config
	if err := SaveConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// expandPaths expands ~ in paths to the home directory
func (c *Config) expandPaths(homeDir string) {
	if c.Claude != nil && c.Claude.SessionDir != "" {
		c.Claude.SessionDir = expandPath(c.Claude.SessionDir, homeDir)
	}
	if c.Claude != nil && c.Claude.MCP != nil && c.Claude.MCP.ConfigFile != "" {
		c.Claude.MCP.ConfigFile = expandPath(c.Claude.MCP.ConfigFile, homeDir)
	}
	if c.Logging != nil && c.Logging.LogFile != "" {
		c.Logging.LogFile = expandPath(c.Logging.LogFile, homeDir)
	}
}

// expandPath expands ~ to home directory
func expandPath(path, homeDir string) string {
	if len(path) > 0 && path[0] == '~' {
		return filepath.Join(homeDir, path[1:])
	}
	return path
}