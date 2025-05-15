package splitscreen

import (
	"github.com/charmbracelet/lipgloss"
)

// ScaleTheme represents a color theme based on a musical scale
type ScaleTheme struct {
	Name          string
	BaseColor     string
	AccentColor   string
	ContrastColor string
	MutedColor    string
	BrightColor   string
}

// Define themes based on musical scales
var (
	// Major scale (bright, confident)
	MajorTheme = ScaleTheme{
		Name:          "Major",
		BaseColor:     "#4A90E2", // Bright blue
		AccentColor:   "#50E3C2", // Cyan
		ContrastColor: "#F5A623", // Gold
		MutedColor:    "#9B9B9B", // Grey
		BrightColor:   "#FFFFFF", // White
	}

	// Minor scale (mysterious, introspective)
	MinorTheme = ScaleTheme{
		Name:          "Minor",
		BaseColor:     "#8B572A", // Brown
		AccentColor:   "#7ED321", // Green
		ContrastColor: "#BD10E0", // Purple
		MutedColor:    "#4A4A4A", // Dark Grey
		BrightColor:   "#F8E71C", // Yellow
	}
	
	// Blues scale (expressive, emotional)
	BluesTheme = ScaleTheme{
		Name:          "Blues",
		BaseColor:     "#2C3E50", // Deep blue
		AccentColor:   "#3498DB", // Bright blue
		ContrastColor: "#E74C3C", // Red
		MutedColor:    "#95A5A6", // Grey
		BrightColor:   "#ECF0F1", // Light grey
	}
	
	// Pentatonic scale (balanced, harmonious)
	PentatonicTheme = ScaleTheme{
		Name:          "Pentatonic",
		BaseColor:     "#27AE60", // Green
		AccentColor:   "#F1C40F", // Yellow
		ContrastColor: "#8E44AD", // Purple
		MutedColor:    "#7F8C8D", // Grey
		BrightColor:   "#ECF0F1", // Light grey
	}
)

// ThemeStyles generates Lipgloss styles from a theme
func ThemeStyles(theme ScaleTheme) map[string]lipgloss.Style {
	return map[string]lipgloss.Style{
		"title": lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.BrightColor)).
			Background(lipgloss.Color(theme.BaseColor)).
			Bold(true).
			Padding(0, 1),
			
		"panel": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.AccentColor)).
			Padding(1, 2),
			
		"activePanel": lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.BrightColor)).
			Padding(1, 2),
			
		"heading": lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.AccentColor)).
			Bold(true),
			
		"subheading": lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ContrastColor)).
			Bold(true),
			
		"text": lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.BrightColor)),
			
		"mutedText": lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.MutedColor)),
			
		"statusBar": lipgloss.NewStyle().
			Background(lipgloss.Color(theme.BaseColor)).
			Foreground(lipgloss.Color(theme.BrightColor)).
			Padding(0, 1),
			
		"timer": lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ContrastColor)).
			Bold(true),
			
		"success": lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2ECC71")),
			
		"error": lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")),
			
		"infoBlock": lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.AccentColor)).
			Padding(1, 2),
	}
}