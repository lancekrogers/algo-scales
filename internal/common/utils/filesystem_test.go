package utils

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

func TestRealFileSystem(t *testing.T) {
	// Create a real file system instance
	fs := NewFileSystem()
	
	// Verify it implements the interface
	var _ interfaces.FileSystem = fs
	
	// Test basic operations in a controlled environment
	tempDir := os.TempDir()
	testDir := filepath.Join(tempDir, "algo-scales-test-"+GetRandomString(8))
	testFile := filepath.Join(testDir, "test-file.txt")
	testContent := []byte("test content")
	
	// Clean up after test
	defer func() {
		os.RemoveAll(testDir)
	}()
	
	// Test MkdirAll
	err := fs.MkdirAll(testDir, 0755)
	assert.NoError(t, err)
	
	// Test Exists
	exists := fs.Exists(testDir)
	assert.True(t, exists)
	
	// Test WriteFile
	err = fs.WriteFile(testFile, testContent, 0644)
	assert.NoError(t, err)
	
	// Test ReadFile
	content, err := fs.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, testContent, content)
	
	// Test Stat
	info, err := fs.Stat(testFile)
	assert.NoError(t, err)
	assert.Equal(t, "test-file.txt", info.Name())
	assert.Equal(t, int64(len(testContent)), info.Size())
	
	// Test ReadDir
	entries, err := fs.ReadDir(testDir)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "test-file.txt", entries[0].Name())
	
	// Test TempDir
	tmpDir := fs.TempDir()
	assert.NotEmpty(t, tmpDir)
	
	// Test GetConfigDir
	configDir := fs.GetConfigDir()
	assert.NotEmpty(t, configDir)
	
	// Test RemoveAll
	err = fs.RemoveAll(testDir)
	assert.NoError(t, err)
	assert.False(t, fs.Exists(testDir))
}

func TestMockFileSystem(t *testing.T) {
	// Create a mock file system
	mockFs := NewMockFileSystem()
	
	// Verify it implements the interface
	var _ interfaces.FileSystem = mockFs
	
	// Test directory creation
	testDir := "/mock/test-dir"
	err := mockFs.MkdirAll(testDir, 0755)
	assert.NoError(t, err)
	assert.True(t, mockFs.Dirs[testDir])
	
	// Test file creation
	testFile := "/mock/test-dir/test-file.txt"
	testContent := []byte("test content")
	err = mockFs.WriteFile(testFile, testContent, 0644)
	assert.NoError(t, err)
	assert.Equal(t, testContent, mockFs.Files[testFile])
	
	// Test file exists
	exists := mockFs.Exists(testFile)
	assert.True(t, exists)
	
	// Test directory exists
	exists = mockFs.Exists(testDir)
	assert.True(t, exists)
	
	// Test file read
	content, err := mockFs.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, testContent, content)
	
	// Test file stat
	info, err := mockFs.Stat(testFile)
	assert.NoError(t, err)
	assert.Equal(t, "test-file.txt", info.Name())
	assert.Equal(t, int64(len(testContent)), info.Size())
	assert.False(t, info.IsDir())
	
	// Test directory stat
	info, err = mockFs.Stat(testDir)
	assert.NoError(t, err)
	assert.Equal(t, "test-dir", info.Name())
	assert.True(t, info.IsDir())
	
	// Test read directory
	entries, err := mockFs.ReadDir(testDir)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "test-file.txt", entries[0].Name())
	assert.False(t, entries[0].IsDir())
	
	// Test open editor
	cmd := mockFs.OpenEditor(testFile)
	assert.NotNil(t, cmd)
	assert.Equal(t, mockFs.EditorCommand, cmd.Path)
	assert.Equal(t, []string{testFile}, mockFs.EditorCalls)
	
	// Test get config dir
	configDir := mockFs.GetConfigDir()
	assert.Equal(t, "/home/mockuser/.algo-scales", configDir)
	
	// Test user home dir
	homeDir, err := mockFs.UserHomeDir()
	assert.NoError(t, err)
	assert.Equal(t, "/home/mockuser", homeDir)
	
	// Test get working dir
	workingDir, err := mockFs.Getwd()
	assert.NoError(t, err)
	assert.Equal(t, "/home/mockuser/projects/algo-scales", workingDir)
	
	// Test get executable
	execPath, err := mockFs.Executable()
	assert.NoError(t, err)
	assert.Equal(t, "/usr/local/bin/algo-scales", execPath)
	
	// Test remove all
	err = mockFs.RemoveAll(testDir)
	assert.NoError(t, err)
	assert.False(t, mockFs.Exists(testDir))
	assert.False(t, mockFs.Exists(testFile))
}

// Helper function for creating random strings
func GetRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}