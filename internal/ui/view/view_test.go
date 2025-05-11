package view

import (
	"strings"
	"testing"
)

func TestPatternVisualization(t *testing.T) {
	// Create a pattern visualization
	pv := NewPatternVisualization()
	
	// Test visualizations for different patterns
	patterns := []string{
		"sliding-window",
		"two-pointers",
		"fast-slow-pointers",
		"hash-map",
		"binary-search",
		"dfs",
		"bfs",
		"dynamic-programming",
		"greedy",
		"union-find",
		"heap",
	}
	
	for _, pattern := range patterns {
		// Generate visualization for this pattern
		viz := pv.VisualizePattern(pattern, "1,2,3,4,5", 40)
		
		// Simple test - just make sure it generates something non-empty
		if len(viz) == 0 {
			t.Errorf("Visualization for pattern %s is empty", pattern)
		}
		
		// Check that the visualization contains expected content for the pattern
		switch pattern {
		case "sliding-window":
			if !strings.Contains(viz, "──") {
				t.Errorf("Sliding window visualization missing window indicator")
			}
		case "two-pointers":
			if !strings.Contains(viz, "▼") {
				t.Errorf("Two pointers visualization missing pointer indicators")
			}
		case "binary-search":
			if !strings.Contains(viz, "mid") {
				t.Errorf("Binary search visualization missing mid indicator")
			}
		}
	}
}

func TestProgressBar(t *testing.T) {
	// Test progress bars with different values
	progressValues := []float64{0.0, 0.25, 0.5, 0.75, 1.0}
	
	for _, progress := range progressValues {
		bar := ProgressBar(10, progress, "sliding-window")
		
		// Check that bar has correct length
		if len(bar) < 10 {
			t.Errorf("Progress bar too short for value %.2f", progress)
		}
		
		// For visual inspection
		t.Logf("Progress %.2f: %s", progress, bar)
	}
}

func TestMusicScales(t *testing.T) {
	// Test that all patterns have music scale definitions
	patterns := []string{
		"sliding-window",
		"two-pointers",
		"fast-slow-pointers",
		"hash-map",
		"binary-search",
		"dfs",
		"bfs",
		"dynamic-programming",
		"greedy",
		"union-find",
		"heap",
	}
	
	for _, pattern := range patterns {
		scale, ok := MusicScales[pattern]
		if !ok {
			t.Errorf("Missing music scale definition for pattern: %s", pattern)
			continue
		}
		
		// Check that scale has all required fields
		if scale.Name == "" {
			t.Errorf("Missing name for pattern %s", pattern)
		}
		if scale.Description == "" {
			t.Errorf("Missing description for pattern %s", pattern)
		}
		// Colors are automatically initialized so no need to check them
	}
}

func TestGetPatternStyle(t *testing.T) {
	// Test that all patterns have styles
	patterns := []string{
		"sliding-window",
		"two-pointers",
		"fast-slow-pointers",
		"hash-map",
		"binary-search",
		"dfs",
		"bfs",
		"dynamic-programming",
		"greedy",
		"union-find",
		"heap",
	}
	
	for _, pattern := range patterns {
		primary, secondary, accent := GetPatternStyle(pattern)
		
		// Simple test - just make sure styles are created
		if primary.String() == "" || secondary.String() == "" || accent.String() == "" {
			t.Errorf("Missing style for pattern %s", pattern)
		}
	}
}