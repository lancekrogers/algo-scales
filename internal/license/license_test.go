// Tests for license module
package license

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateLicense(t *testing.T) {
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
	t.Run("NoLicenseFile", func(t *testing.T) {
		valid, err := ValidateLicense()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "license file not found")
		assert.False(t, valid)
	})

	t.Run("ValidLicense", func(t *testing.T) {
		// Create a valid license
		license := License{
			LicenseKey:   "valid-key",
			Email:        "test@example.com",
			PurchaseDate: time.Now(),
			ExpiryDate:   time.Now().AddDate(1, 0, 0), // Valid for 1 year
			Signature:    "valid-signature",
		}

		// Save license to file
		licenseFile := filepath.Join(tempDir, "license.json")
		licenseData, err := json.MarshalIndent(license, "", "  ")
		require.NoError(t, err)
		err = os.WriteFile(licenseFile, licenseData, 0644)
		require.NoError(t, err)

		// Override verify signature for testing
		origVerifySignature := verifySignature
		defer func() { verifySignature = origVerifySignature }()
		verifySignature = func(lic License) bool {
			return true
		}

		// Validate the license
		valid, err := ValidateLicense()
		require.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("ExpiredLicense", func(t *testing.T) {
		// Create an expired license
		license := License{
			LicenseKey:   "expired-key",
			Email:        "test@example.com",
			PurchaseDate: time.Now().AddDate(-2, 0, 0), // 2 years ago
			ExpiryDate:   time.Now().AddDate(-1, 0, 0), // Expired 1 year ago
			Signature:    "valid-signature",
		}

		// Save license to file
		licenseFile := filepath.Join(tempDir, "license.json")
		licenseData, err := json.MarshalIndent(license, "", "  ")
		require.NoError(t, err)
		err = os.WriteFile(licenseFile, licenseData, 0644)
		require.NoError(t, err)

		// Validate the license
		valid, err := ValidateLicense()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "license expired")
		assert.False(t, valid)
	})

	t.Run("InvalidSignature", func(t *testing.T) {
		// Create a license with invalid signature
		license := License{
			LicenseKey:   "invalid-sig-key",
			Email:        "test@example.com",
			PurchaseDate: time.Now(),
			ExpiryDate:   time.Now().AddDate(1, 0, 0), // Valid for 1 year
			Signature:    "invalid-signature",
		}

		// Save license to file
		licenseFile := filepath.Join(tempDir, "license.json")
		licenseData, err := json.MarshalIndent(license, "", "  ")
		require.NoError(t, err)
		err = os.WriteFile(licenseFile, licenseData, 0644)
		require.NoError(t, err)

		// Override verify signature for testing
		origVerifySignature := verifySignature
		defer func() { verifySignature = origVerifySignature }()
		verifySignature = func(lic License) bool {
			return false
		}

		// Validate the license
		valid, err := ValidateLicense()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid license signature")
		assert.False(t, valid)
	})

	t.Run("CorruptLicenseFile", func(t *testing.T) {
		// Create a corrupt license file
		licenseFile := filepath.Join(tempDir, "license.json")
		err = os.WriteFile(licenseFile, []byte("corrupt json"), 0644)
		require.NoError(t, err)

		// Validate the license
		valid, err := ValidateLicense()
		require.Error(t, err)
		assert.False(t, valid)
	})
}

func TestRequestLicense(t *testing.T) {
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

	// We can't easily test the interactive parts of RequestLicense
	// that require user input, but we can test the file writing functionality
	// by mocking the console input.

	// This would be a more complete implementation in a real test:
	// - Mock os.Stdin to provide predetermined input
	// - Capture os.Stdout to verify prompts
	// - Test the creation and validation of the license file

	// For now, we'll just test the helper functions:

	t.Run("GenerateSignature", func(t *testing.T) {
		sig := generateSignature("test-key", "test@example.com")
		assert.NotEmpty(t, sig)
		assert.Contains(t, sig, "test")
	})
}

func TestVerifySignature(t *testing.T) {
	// For the MVP version, verifySignature always returns true if the signature exists
	license := License{
		LicenseKey:   "test-key",
		Email:        "test@example.com",
		PurchaseDate: time.Now(),
		ExpiryDate:   time.Now().AddDate(1, 0, 0),
		Signature:    "valid-signature",
	}

	assert.True(t, verifySignature(license))

	// Test with empty signature
	license.Signature = ""
	assert.False(t, verifySignature(license))
}
