// Tests for API client

package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lancekrogers/algo-scales/internal/license"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock license validation
func mockLicenseValidation(valid bool, err error) func() {
	original := license.ValidateLicense
	license.ValidateLicense = func() (bool, error) {
		return valid, err
	}
	return func() {
		license.ValidateLicense = original
	}
}

func TestDownloadProblems(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "algo-scales-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	defer func() { getConfigDir = origGetConfigDir }()
	getConfigDir = func() string {
		return tempDir
	}

	// Test cases
	t.Run("InvalidLicense", func(t *testing.T) {
		// Mock invalid license
		restore := mockLicenseValidation(false, nil)
		defer restore()

		// Try to download problems
		err := DownloadProblems(true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid license")
	})

	t.Run("ValidLicense", func(t *testing.T) {
		// Mock valid license
		restore := mockLicenseValidation(true, nil)
		defer restore()

		// Try to download problems
		err := DownloadProblems(true)
		require.NoError(t, err)

		// Verify problem files were created
		problemsDir := filepath.Join(tempDir, "problems")
		for _, problem := range getSampleProblems().Problems {
			for _, pattern := range problem.Patterns {
				problemFile := filepath.Join(problemsDir, pattern, problem.ID+".json")
				_, err := os.Stat(problemFile)
				assert.NoError(t, err, "Problem file should exist: %s", problemFile)
			}
		}

		// Verify version file was created
		versionFile := filepath.Join(tempDir, "version.json")
		_, err = os.Stat(versionFile)
		assert.NoError(t, err, "Version file should exist")

		// Check version file content
		data, err := os.ReadFile(versionFile)
		require.NoError(t, err)
		var version struct {
			Version     string    `json:"version"`
			LastUpdated time.Time `json:"last_updated"`
		}
		err = json.Unmarshal(data, &version)
		require.NoError(t, err)
		assert.Equal(t, "1.0.0", version.Version)
	})

	t.Run("NoUpdateNeeded", func(t *testing.T) {
		// Mock valid license
		restore := mockLicenseValidation(true, nil)
		defer restore()

		// Create a version file with recent timestamp
		versionFile := filepath.Join(tempDir, "version.json")
		versionData, err := json.MarshalIndent(map[string]interface{}{
			"version":      "1.0.0",
			"last_updated": time.Now(),
		}, "", "  ")
		require.NoError(t, err)
		err = os.WriteFile(versionFile, versionData, 0644)
		require.NoError(t, err)

		// Try to download problems without force flag
		err = DownloadProblems(false)
		require.NoError(t, err)

		// No files should have been modified
		fileInfo, err := os.Stat(versionFile)
		require.NoError(t, err)
		assert.Equal(t, int64(len(versionData)), fileInfo.Size(), "Version file should not have been modified")
	})
}

func TestShouldUpdate(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "algo-scales-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override config dir for testing
	origGetConfigDir := getConfigDir
	defer func() { getConfigDir = origGetConfigDir }()
	getConfigDir = func() string {
		return tempDir
	}

	t.Run("NoVersionFile", func(t *testing.T) {
		// If no version file exists, we should update
		assert.True(t, shouldUpdate())
	})

	t.Run("CorruptVersionFile", func(t *testing.T) {
		// Create a corrupt version file
		versionFile := filepath.Join(tempDir, "version.json")
		err = os.WriteFile(versionFile, []byte("corrupt json"), 0644)
		require.NoError(t, err)

		// With a corrupt version file, we should update
		assert.True(t, shouldUpdate())
	})

	t.Run("ValidVersionFile", func(t *testing.T) {
		// Create a valid version file
		versionFile := filepath.Join(tempDir, "version.json")
		versionData, err := json.MarshalIndent(map[string]interface{}{
			"version":      "1.0.0",
			"last_updated": time.Now(),
		}, "", "  ")
		require.NoError(t, err)
		err = os.WriteFile(versionFile, versionData, 0644)
		require.NoError(t, err)

		// In MVP, we always return false after initial download
		assert.False(t, shouldUpdate())
	})
}

func TestGetSampleProblems(t *testing.T) {
	// Test the sample problems generation
	problems := getSampleProblems()

	// Verify structure
	assert.Equal(t, "1.0.0", problems.Version)
	assert.NotNil(t, problems.LastUpdated)
	assert.NotEmpty(t, problems.Problems)

	// Check first problem
	firstProblem := problems.Problems[0]
	assert.Equal(t, "two-sum", firstProblem.ID)
	assert.Equal(t, "Two Sum", firstProblem.Title)
	assert.NotEmpty(t, firstProblem.Description)
	assert.NotEmpty(t, firstProblem.StarterCode)
	assert.NotEmpty(t, firstProblem.Solutions)
}
