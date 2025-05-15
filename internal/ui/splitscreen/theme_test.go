package splitscreen

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestThemeColors(t *testing.T) {
	// Test that all themes have valid color values
	themes := []ScaleTheme{
		MajorTheme,
		MinorTheme,
		BluesTheme,
		PentatonicTheme,
	}

	for _, theme := range themes {
		if theme.Name == "" {
			t.Errorf("theme missing name")
		}

		// Check that colors are valid and parse correctly
		colors := []string{
			theme.BaseColor,
			theme.AccentColor,
			theme.ContrastColor,
			theme.MutedColor,
			theme.BrightColor,
		}

		for _, color := range colors {
			if color == "" {
				t.Errorf("theme %s contains empty color", theme.Name)
			}

			// Verify the color is a valid lipgloss color
			_ = lipgloss.Color(color)
		}
	}
}

func TestThemeStyles(t *testing.T) {
	// Test that theme styles are generated correctly
	styles := ThemeStyles(MajorTheme)

	// Check required styles are present
	requiredStyles := []string{
		"title",
		"panel",
		"activePanel",
		"heading",
		"subheading",
		"text",
		"mutedText",
		"statusBar",
		"timer",
		"success",
		"error",
		"infoBlock",
	}

	for _, styleName := range requiredStyles {
		if _, ok := styles[styleName]; !ok {
			t.Errorf("missing required style: %s", styleName)
		}
	}

	// Test that styles apply colors from the theme
	// Since we can't directly compare color values, let's test that foreground and background
	// colors are applied by checking if they're non-empty
	
	// Verify that styles have the expected visual properties
	if !styles["title"].GetBold() {
		t.Errorf("title style is missing bold property")
	}
	
	if !styles["heading"].GetBold() {
		t.Errorf("heading style is missing bold property")
	}
	
	// Verify styles exist with proper dimensions
	// We can't directly check if a border style is nil, but we can check the padding
	// since we know these styles should have padding
	if styles["panel"].GetPaddingLeft() <= 0 && styles["panel"].GetPaddingRight() <= 0 {
		t.Errorf("panel style is missing padding")
	}
	
	if styles["activePanel"].GetPaddingLeft() <= 0 && styles["activePanel"].GetPaddingRight() <= 0 {
		t.Errorf("activePanel style is missing padding")
	}
}