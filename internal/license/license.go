// License validation functionality

package license

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// License represents a user license
type License struct {
	LicenseKey   string    `json:"license_key"`
	Email        string    `json:"email"`
	PurchaseDate time.Time `json:"purchase_date"`
	ExpiryDate   time.Time `json:"expiry_date"` // For potential subscription model
	Signature    string    `json:"signature"`
}

// ValidateLicense checks if the license is valid
// Exported as variable for testing
var ValidateLicense = func() (bool, error) {
	licenseFile := filepath.Join(getConfigDir(), "license.json")

	// Check if license file exists
	if _, err := os.Stat(licenseFile); os.IsNotExist(err) {
		return false, fmt.Errorf("license file not found")
	}

	// Read license file
	data, err := os.ReadFile(licenseFile)
	if err != nil {
		return false, err
	}

	// Parse license
	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return false, err
	}

	// Check expiry (for subscription model)
	if !license.ExpiryDate.IsZero() && time.Now().After(license.ExpiryDate) {
		return false, fmt.Errorf("license expired")
	}

	// Verify signature (simplified for MVP)
	isValid := verifySignature(license)
	if !isValid {
		return false, fmt.Errorf("invalid license signature")
	}

	return true, nil
}

// RequestLicense prompts the user for their license key
// For MVP, we'll just create a dummy license
func RequestLicense() error {
	var email, licenseKey string

	// In a real implementation, you'd validate this with an API call
	// For MVP, we'll just create a dummy license
	fmt.Print("Enter your email: ")
	fmt.Scanln(&email)

	fmt.Print("Enter your license key: ")
	fmt.Scanln(&licenseKey)

	// Create a license (for demo purposes this is always valid)
	license := License{
		LicenseKey:   licenseKey,
		Email:        email,
		PurchaseDate: time.Now(),
		ExpiryDate:   time.Now().AddDate(1, 0, 0), // Valid for 1 year
		Signature:    generateSignature(licenseKey, email),
	}

	// Save license to file
	licenseFile := filepath.Join(getConfigDir(), "license.json")
	licenseData, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(licenseFile, licenseData, 0644)
}

// Helper functions - exported as variables for testing
var getConfigDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}

// verifySignature checks if a license signature is valid
// This is a simplified version for MVP - exported as variable for testing
var verifySignature = func(license License) bool {
	// In a real implementation, this would use public key cryptography
	// For MVP, we'll just check if the signature exists
	return license.Signature != ""
}

// generateSignature creates a signature for a license
// This is a simplified version for MVP
func generateSignature(licenseKey, email string) string {
	// In a real implementation, this would use private key cryptography
	// For MVP, we'll just return a dummy signature
	return "valid-signature-" + licenseKey[:4]
}
