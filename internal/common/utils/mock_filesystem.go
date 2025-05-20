package utils

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// MockFileInfo implements os.FileInfo for testing
type MockFileInfo struct {
	FileName    string
	FileSize    int64
	FileMode    os.FileMode
	FileModTime time.Time
	FileIsDir   bool
}

func (m MockFileInfo) Name() string       { return m.FileName }
func (m MockFileInfo) Size() int64        { return m.FileSize }
func (m MockFileInfo) Mode() os.FileMode  { return m.FileMode }
func (m MockFileInfo) ModTime() time.Time { return m.FileModTime }
func (m MockFileInfo) IsDir() bool        { return m.FileIsDir }
func (m MockFileInfo) Sys() interface{}   { return nil }

// MockDirEntry implements fs.DirEntry for testing
type MockDirEntry struct {
	EntryName  string
	EntryIsDir bool
	EntryType  fs.FileMode
	EntryInfo  fs.FileInfo
}

func (m MockDirEntry) Name() string               { return m.EntryName }
func (m MockDirEntry) IsDir() bool                { return m.EntryIsDir }
func (m MockDirEntry) Type() fs.FileMode          { return m.EntryType }
func (m MockDirEntry) Info() (fs.FileInfo, error) { return m.EntryInfo, nil }

// MockFileSystem implements the FileSystem interface for testing
type MockFileSystem struct {
	Files      map[string][]byte
	Dirs       map[string]bool
	FileInfos  map[string]os.FileInfo
	DirEntries map[string][]fs.DirEntry
	
	HomeDir      string
	ConfigDir    string
	WorkingDir   string
	ExecutablePath string
	TempDirPath  string
	
	EditorCommand string
	EditorCalls   []string
}

// NewMockFileSystem creates a new mock file system
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files:      make(map[string][]byte),
		Dirs:       make(map[string]bool),
		FileInfos:  make(map[string]os.FileInfo),
		DirEntries: make(map[string][]fs.DirEntry),
		
		HomeDir:      "/home/mockuser",
		ConfigDir:    "/home/mockuser/.algo-scales",
		WorkingDir:   "/home/mockuser/projects/algo-scales",
		ExecutablePath: "/usr/local/bin/algo-scales",
		TempDirPath:  "/tmp",
		
		EditorCommand: "mock-editor",
		EditorCalls:   []string{},
	}
}

// Ensure MockFileSystem implements FileSystem
var _ interfaces.FileSystem = (*MockFileSystem)(nil)

// ReadFile reads the named file and returns its contents
func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	if data, ok := m.Files[path]; ok {
		return data, nil
	}
	return nil, os.ErrNotExist
}

// WriteFile writes data to the named file
func (m *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	if !m.Exists(dir) {
		m.Dirs[dir] = true
	}
	
	m.Files[path] = data
	
	// Create a file info entry if it doesn't exist
	if _, ok := m.FileInfos[path]; !ok {
		m.FileInfos[path] = &MockFileInfo{
			FileName:    filepath.Base(path),
			FileSize:    int64(len(data)),
			FileMode:    perm,
			FileModTime: time.Now(),
			FileIsDir:   false,
		}
	}
	
	return nil
}

// MkdirAll creates a directory and all necessary parents
func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	m.Dirs[path] = true
	
	// Create file info for the directory
	m.FileInfos[path] = &MockFileInfo{
		FileName:    filepath.Base(path),
		FileSize:    0,
		FileMode:    perm | os.ModeDir,
		FileModTime: time.Now(),
		FileIsDir:   true,
	}
	
	// Create parent directories
	parent := filepath.Dir(path)
	if parent != path {
		m.Dirs[parent] = true
	}
	
	return nil
}

// RemoveAll removes path and any children it contains
func (m *MockFileSystem) RemoveAll(path string) error {
	// Remove directories
	delete(m.Dirs, path)
	
	// Remove file infos
	delete(m.FileInfos, path)
	
	// Remove all files with this path as prefix
	for filePath := range m.Files {
		if filePath == path || strings.HasPrefix(filePath, path+"/") {
			delete(m.Files, filePath)
		}
	}
	
	// Remove all dirs with this path as prefix
	for dirPath := range m.Dirs {
		if dirPath == path || strings.HasPrefix(dirPath, path+"/") {
			delete(m.Dirs, dirPath)
		}
	}
	
	return nil
}

// Stat returns a FileInfo describing the named file
func (m *MockFileSystem) Stat(path string) (os.FileInfo, error) {
	if info, ok := m.FileInfos[path]; ok {
		return info, nil
	}
	
	// Check if it's a directory
	if m.Dirs[path] {
		return &MockFileInfo{
			FileName:    filepath.Base(path),
			FileSize:    0,
			FileMode:    0755 | os.ModeDir,
			FileModTime: time.Now(),
			FileIsDir:   true,
		}, nil
	}
	
	// Check if it's a file
	if _, ok := m.Files[path]; ok {
		return &MockFileInfo{
			FileName:    filepath.Base(path),
			FileSize:    int64(len(m.Files[path])),
			FileMode:    0644,
			FileModTime: time.Now(),
			FileIsDir:   false,
		}, nil
	}
	
	return nil, os.ErrNotExist
}

// Exists checks if a file or directory exists
func (m *MockFileSystem) Exists(path string) bool {
	_, fileOk := m.Files[path]
	_, dirOk := m.Dirs[path]
	return fileOk || dirOk
}

// ReadDir reads the named directory and returns a list of directory entries
func (m *MockFileSystem) ReadDir(path string) ([]fs.DirEntry, error) {
	if entries, ok := m.DirEntries[path]; ok {
		return entries, nil
	}
	
	// If no explicit entries, generate them based on files and dirs
	entries := []fs.DirEntry{}
	
	// Check if the directory exists
	if !m.Dirs[path] && !strings.HasSuffix(path, "/") {
		return nil, os.ErrNotExist
	}
	
	// Add all immediate children
	pathPrefix := path
	if !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix += "/"
	}
	
	// Process files
	for filePath, fileData := range m.Files {
		if strings.HasPrefix(filePath, pathPrefix) {
			// Only include immediate children
			relPath := strings.TrimPrefix(filePath, pathPrefix)
			if !strings.Contains(relPath, "/") {
				info := &MockFileInfo{
					FileName:    filepath.Base(filePath),
					FileSize:    int64(len(fileData)),
					FileMode:    0644,
					FileModTime: time.Now(),
					FileIsDir:   false,
				}
				
				entry := MockDirEntry{
					EntryName:  filepath.Base(filePath),
					EntryIsDir: false,
					EntryType:  0644,
					EntryInfo:  info,
				}
				
				entries = append(entries, entry)
			}
		}
	}
	
	// Process directories
	for dirPath := range m.Dirs {
		if strings.HasPrefix(dirPath, pathPrefix) && dirPath != path {
			// Only include immediate children
			relPath := strings.TrimPrefix(dirPath, pathPrefix)
			if !strings.Contains(relPath, "/") {
				info := &MockFileInfo{
					FileName:    filepath.Base(dirPath),
					FileSize:    0,
					FileMode:    0755 | os.ModeDir,
					FileModTime: time.Now(),
					FileIsDir:   true,
				}
				
				entry := MockDirEntry{
					EntryName:  filepath.Base(dirPath),
					EntryIsDir: true,
					EntryType:  0755 | os.ModeDir,
					EntryInfo:  info,
				}
				
				entries = append(entries, entry)
			}
		}
	}
	
	return entries, nil
}

// TempDir returns the default directory for temporary files
func (m *MockFileSystem) TempDir() string {
	return m.TempDirPath
}

// GetConfigDir returns the application config directory
func (m *MockFileSystem) GetConfigDir() string {
	return m.ConfigDir
}

// OpenEditor opens the user's editor with the given file
func (m *MockFileSystem) OpenEditor(filePath string) *exec.Cmd {
	m.EditorCalls = append(m.EditorCalls, filePath)
	
	// Create a mock command that does nothing
	cmd := exec.Command(m.EditorCommand, filePath)
	return cmd
}

// UserHomeDir returns the current user's home directory
func (m *MockFileSystem) UserHomeDir() (string, error) {
	if m.HomeDir == "" {
		return "", errors.New("home directory not configured in mock")
	}
	return m.HomeDir, nil
}

// Getwd returns the current working directory
func (m *MockFileSystem) Getwd() (string, error) {
	if m.WorkingDir == "" {
		return "", errors.New("working directory not configured in mock")
	}
	return m.WorkingDir, nil
}

// Executable returns the path to the current executable
func (m *MockFileSystem) Executable() (string, error) {
	if m.ExecutablePath == "" {
		return "", errors.New("executable path not configured in mock")
	}
	return m.ExecutablePath, nil
}