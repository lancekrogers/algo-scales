// Package utils provides utility functions for Algo Scales
package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Function variables that can be mocked in tests

// GetConfigDir returns the configuration directory
var GetConfigDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}

// CreateDirectory creates a directory if it doesn't exist
var CreateDirectory = func(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists checks if a file or directory exists
var FileExists = func(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ReadFile reads a file's contents
var ReadFile = func(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to a file
var WriteFile = func(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// TempDir returns a temporary directory path
var TempDir = func() string {
	return os.TempDir()
}

// OpenEditor opens the user's editor
var OpenEditor = func(filePath string) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default to vim
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}