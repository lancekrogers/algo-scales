package main

import "testing"

func TestGenerateLicenseKey_ShortEmail(t *testing.T) {
	key := generateLicenseKey("a@b")
	if key == "" {
		t.Fatal("expected non-empty license key")
	}
}
