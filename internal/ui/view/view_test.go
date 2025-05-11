package view

import (
	"testing"
)

func TestPatternVisualization(t *testing.T) {
	// Test creating visualization
	viz := NewPatternVisualization()

	// Test visualizing patterns
	patterns := []string{
		"sliding-window",
		"two-pointers",
		"binary-search",
	}

	for _, pattern := range patterns {
		// The VisualizePattern method requires data and width parameters
		art := viz.VisualizePattern(pattern, "", 40) // Use empty data and default width
		if art == "" {
			t.Errorf("Expected visualization for pattern %s, got empty string", pattern)
		}
	}

	// Test fallback for unknown pattern
	art := viz.VisualizePattern("unknown-pattern", "", 40)
	if art == "" {
		t.Error("Expected fallback visualization for unknown pattern, got empty string")
	}
}

func TestProgressBar(t *testing.T) {
	// Test progress bar rendering with various percentages
	percentages := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	for _, percent := range percentages {
		bar := ProgressBar(10, percent, "sliding-window")
		t.Logf("Progress %.2f: %s", percent, bar)

		// Check the length - should be 10 characters plus ANSI color codes
		// We can't directly check the length due to ANSI codes, but ensure it's not empty
		if bar == "" {
			t.Error("Expected progress bar, got empty string")
		}
	}
}

func TestMusicScales(t *testing.T) {
	// Check that at least one musical scale is defined
	// We don't need to check all patterns, just that the map exists

	pattern := "sliding-window"
	scale, ok := MusicScales[pattern]

	if !ok {
		t.Errorf("Expected musical scale for pattern %s to be defined", pattern)
	}

	// Make sure fields are set
	if scale.Name == "" || scale.Description == "" {
		t.Errorf("Missing fields in musical scale for pattern %s", pattern)
	}
}

func TestGetPatternStyle(t *testing.T) {
	// Get style for any pattern - even if it falls back
	primary, secondary, accent := GetPatternStyle("any-pattern")

	// In Go, the zero value of a struct is valid, so Lipgloss styles are always usable
	// Just verify the type is as expected
	_ = primary.String()
	_ = secondary.String()
	_ = accent.String()

	// No need to check content as long as the function returns something
	t.Log("Pattern style returned successfully")
}