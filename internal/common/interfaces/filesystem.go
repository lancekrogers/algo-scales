// Package interfaces defines the core interfaces for Algo Scales
package interfaces

import (
	"io/fs"
	"os"
	"os/exec"
)

// FileSystem defines an interface for file system operations
// This allows for easier testing and mocking of file operations
type FileSystem interface {
	// ReadFile reads the named file and returns its contents
	ReadFile(path string) ([]byte, error)
	
	// WriteFile writes data to the named file
	WriteFile(path string, data []byte, perm os.FileMode) error
	
	// MkdirAll creates a directory and all necessary parents
	MkdirAll(path string, perm os.FileMode) error
	
	// RemoveAll removes path and any children it contains
	RemoveAll(path string) error
	
	// Stat returns a FileInfo describing the named file
	Stat(path string) (os.FileInfo, error)
	
	// Exists checks if a file or directory exists
	Exists(path string) bool
	
	// ReadDir reads the named directory and returns a list of directory entries
	ReadDir(path string) ([]fs.DirEntry, error)
	
	// TempDir returns the default directory for temporary files
	TempDir() string
	
	// GetConfigDir returns the application config directory
	GetConfigDir() string
	
	// OpenEditor opens the user's editor with the given file
	OpenEditor(filePath string) *exec.Cmd
	
	// UserHomeDir returns the current user's home directory
	UserHomeDir() (string, error)
	
	// Getwd returns the current working directory
	Getwd() (string, error)
	
	// Executable returns the path to the current executable
	Executable() (string, error)
}