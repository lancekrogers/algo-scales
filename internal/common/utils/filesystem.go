package utils

import (
	"io/fs"
	"os"
	"os/exec"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// RealFileSystem implements the FileSystem interface using the actual OS file system
type RealFileSystem struct{}

// NewFileSystem creates a new FileSystem implementation
func NewFileSystem() interfaces.FileSystem {
	return &RealFileSystem{}
}

// ReadFile reads the named file and returns its contents
func (fs *RealFileSystem) ReadFile(path string) ([]byte, error) {
	return ReadFile(path)
}

// WriteFile writes data to the named file
func (fs *RealFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	return WriteFile(path, data, perm)
}

// MkdirAll creates a directory and all necessary parents
func (fs *RealFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return CreateDirectory(path)
}

// RemoveAll removes path and any children it contains
func (fs *RealFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Stat returns a FileInfo describing the named file
func (fs *RealFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Exists checks if a file or directory exists
func (fs *RealFileSystem) Exists(path string) bool {
	return FileExists(path)
}

// ReadDir reads the named directory and returns a list of directory entries
func (fs *RealFileSystem) ReadDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(path)
}

// TempDir returns the default directory for temporary files
func (fs *RealFileSystem) TempDir() string {
	return TempDir()
}

// GetConfigDir returns the application config directory
func (fs *RealFileSystem) GetConfigDir() string {
	return GetConfigDir()
}

// OpenEditor opens the user's editor with the given file
func (fs *RealFileSystem) OpenEditor(filePath string) *exec.Cmd {
	return OpenEditor(filePath)
}

// UserHomeDir returns the current user's home directory
func (fs *RealFileSystem) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

// Getwd returns the current working directory
func (fs *RealFileSystem) Getwd() (string, error) {
	return os.Getwd()
}

// Executable returns the path to the current executable
func (fs *RealFileSystem) Executable() (string, error) {
	return os.Executable()
}