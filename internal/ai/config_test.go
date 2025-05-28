package ai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".algo-scales", "ai-config.yaml")

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Test 1: Create default config when none exists
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	if config.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", config.Version)
	}

	if config.DefaultProvider != "claude" {
		t.Errorf("Expected default provider claude, got %s", config.DefaultProvider)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Test 2: Load existing config
	config2, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load existing config: %v", err)
	}

	if config2.Version != config.Version {
		t.Error("Loaded config doesn't match original")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	config := &Config{
		Version:         "1.0",
		DefaultProvider: "ollama",
		Ollama: &OllamaConfig{
			Host:  "http://localhost:11434",
			Model: "codellama",
		},
	}

	err := SaveConfig(config)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load and verify
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loaded.DefaultProvider != "ollama" {
		t.Errorf("Expected provider ollama, got %s", loaded.DefaultProvider)
	}

	if loaded.Ollama.Model != "codellama" {
		t.Errorf("Expected model codellama, got %s", loaded.Ollama.Model)
	}
}

func TestExpandPaths(t *testing.T) {
	config := &Config{
		Claude: &ClaudeConfig{
			SessionDir: "~/sessions",
		},
		Logging: &LoggingConfig{
			LogFile: "~/logs/ai.log",
		},
	}

	config.expandPaths("/home/user")

	expectedSession := "/home/user/sessions"
	if config.Claude.SessionDir != expectedSession {
		t.Errorf("Expected %s, got %s", expectedSession, config.Claude.SessionDir)
	}

	expectedLog := "/home/user/logs/ai.log"
	if config.Logging.LogFile != expectedLog {
		t.Errorf("Expected %s, got %s", expectedLog, config.Logging.LogFile)
	}
}