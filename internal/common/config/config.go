// Package config handles configuration and settings
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// UserConfig represents the user's configuration
type UserConfig struct {
	// User preferences
	Language      string `json:"language"`      // Preferred programming language
	TimerDuration int    `json:"timerDuration"` // Timer duration in minutes
	Mode          string `json:"mode"`          // Default mode: "learn", "practice", "cram"
	
	// UI preferences
	Theme         string `json:"theme"`         // UI theme
	EditorCommand string `json:"editorCommand"` // External editor command
	
	// Focus settings
	FocusPatterns []string `json:"focusPatterns"` // Patterns to focus on
}

// DefaultConfig returns the default configuration
func DefaultConfig() UserConfig {
	return UserConfig{
		Language:      "go",
		TimerDuration: 30,
		Mode:          "practice",
		Theme:         "default",
		EditorCommand: getDefaultEditor(),
		FocusPatterns: []string{},
	}
}

// LoadConfig loads the user's configuration from file
func LoadConfig() (UserConfig, error) {
	configDir := getConfigDir()
	configFile := filepath.Join(configDir, "config.json")
	
	// If config file doesn't exist, create default
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := DefaultConfig()
		err := SaveConfig(config)
		return config, err
	}
	
	// Read config file
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return DefaultConfig(), err
	}
	
	// Parse config data
	var config UserConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return DefaultConfig(), err
	}
	
	return config, nil
}

// SaveConfig saves the user's configuration to file
func SaveConfig(config UserConfig) error {
	configDir := getConfigDir()
	
	// Create config directory if it doesn't exist
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	
	// Write config file
	configFile := filepath.Join(configDir, "config.json")
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	
	return nil
}

// getConfigDir returns the configuration directory path
func getConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}

// getDefaultEditor returns the default editor command based on the OS and environment
func getDefaultEditor() string {
	// Check EDITOR environment variable first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	
	// Check VISUAL environment variable
	if visual := os.Getenv("VISUAL"); visual != "" {
		return visual
	}
	
	// Default to vi/vim on Unix-like systems
	if _, err := exec.LookPath("vim"); err == nil {
		return "vim"
	}
	
	if _, err := exec.LookPath("vi"); err == nil {
		return "vi"
	}
	
	// On Windows, default to notepad
	if _, err := exec.LookPath("notepad"); err == nil {
		return "notepad"
	}
	
	// If all else fails, return a placeholder
	return "editor"
}

// ListLanguages returns the available programming languages
func ListLanguages() []string {
	return []string{
		"go",
		"python",
		"java",
		"javascript",
		"typescript",
		"c++",
		"c#",
	}
}

// ListModes returns the available learning modes
func ListModes() []string {
	return []string{
		"learn",
		"practice",
		"cram",
	}
}

// ListTimerOptions returns the available timer durations in minutes
func ListTimerOptions() []int {
	return []int{
		5,
		10,
		15, 
		20,
		30,
		45,
		60,
	}
}