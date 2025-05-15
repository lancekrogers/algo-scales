// Package view contains UI components and rendering
package view

import (
	"strings"
	"github.com/charmbracelet/lipgloss"
)

// MusicScale represents a musical scale and its associated colors
type MusicScale struct {
	Name          string
	Pattern       string
	Description   string
	PrimaryColor  lipgloss.Color
	SecondaryColor lipgloss.Color
	AccentColor   lipgloss.Color
}

// Defining colors for each musical scale/algorithm pattern
var (
	// Color definitions
	cMajorBlue         = lipgloss.Color("#3498db")
	cMajorLightBlue    = lipgloss.Color("#5dade2")
	cMajorDarkBlue     = lipgloss.Color("#2874a6")

	gMajorGreen        = lipgloss.Color("#2ecc71")
	gMajorLightGreen   = lipgloss.Color("#58d68d")
	gMajorDarkGreen    = lipgloss.Color("#239b56")

	dMajorOrange       = lipgloss.Color("#e67e22")
	dMajorLightOrange  = lipgloss.Color("#eb984e")
	dMajorDarkOrange   = lipgloss.Color("#af601a")

	aMajorRed          = lipgloss.Color("#e74c3c")
	aMajorLightRed     = lipgloss.Color("#ec7063")
	aMajorDarkRed      = lipgloss.Color("#b03a2e")

	eMajorPurple       = lipgloss.Color("#9b59b6")
	eMajorLightPurple  = lipgloss.Color("#af7ac5")
	eMajorDarkPurple   = lipgloss.Color("#7d3c98")

	bMajorDeepBlue     = lipgloss.Color("#1b4f72")
	bMajorMediumBlue   = lipgloss.Color("#2874a6")
	bMajorLightBlue    = lipgloss.Color("#3498db")

	fSharpMajorTeal    = lipgloss.Color("#16a085")
	fSharpMajorLightTeal = lipgloss.Color("#45b39d")
	fSharpMajorDarkTeal = lipgloss.Color("#117a65")

	dbMajorYellow      = lipgloss.Color("#f1c40f")
	dbMajorLightYellow = lipgloss.Color("#f4d03f")
	dbMajorDarkYellow  = lipgloss.Color("#b7950b")

	abMajorMagenta     = lipgloss.Color("#8e44ad")
	abMajorLightMagenta = lipgloss.Color("#a569bd")
	abMajorDarkMagenta = lipgloss.Color("#6c3483")

	ebMajorCyan        = lipgloss.Color("#00bcd4")
	ebMajorLightCyan   = lipgloss.Color("#4dd0e1")
	ebMajorDarkCyan    = lipgloss.Color("#0097a7")

	bbMajorBrown       = lipgloss.Color("#795548")
	bbMajorLightBrown  = lipgloss.Color("#a1887f")
	bbMajorDarkBrown   = lipgloss.Color("#5d4037")

	// Default colors
	subtleGray       = lipgloss.Color("#6c757d")
	defaultFg        = lipgloss.Color("#eeeeee")
	defaultBg        = lipgloss.Color("#333333")
	successGreen     = lipgloss.Color("#28a745")
	errorRed         = lipgloss.Color("#dc3545")
	warningYellow    = lipgloss.Color("#ffc107")
	infoBlue         = lipgloss.Color("#17a2b8")
)

// Musical scale definitions with their pattern associations
var MusicScales = map[string]MusicScale{
	"sliding-window": {
		Name:          "C Major (Sliding Window)",
		Pattern:       "sliding-window",
		Description:   "The fundamental scale, elegant and versatile",
		PrimaryColor:  cMajorBlue,
		SecondaryColor: cMajorLightBlue,
		AccentColor:   cMajorDarkBlue,
	},
	"two-pointers": {
		Name:          "G Major (Two Pointers)",
		Pattern:       "two-pointers",
		Description:   "Balanced and efficient, the workhorse of array manipulation",
		PrimaryColor:  gMajorGreen,
		SecondaryColor: gMajorLightGreen,
		AccentColor:   gMajorDarkGreen,
	},
	"fast-slow-pointers": {
		Name:          "D Major (Fast & Slow Pointers)",
		Pattern:       "fast-slow-pointers",
		Description:   "The cycle detector, bright and revealing",
		PrimaryColor:  dMajorOrange,
		SecondaryColor: dMajorLightOrange,
		AccentColor:   dMajorDarkOrange,
	},
	"hash-map": {
		Name:          "A Major (Hash Maps)",
		Pattern:       "hash-map",
		Description:   "The lookup accelerator, crisp and direct",
		PrimaryColor:  aMajorRed,
		SecondaryColor: aMajorLightRed,
		AccentColor:   aMajorDarkRed,
	},
	"binary-search": {
		Name:          "E Major (Binary Search)",
		Pattern:       "binary-search",
		Description:   "The divider and conqueror, precise and logarithmic",
		PrimaryColor:  eMajorPurple,
		SecondaryColor: eMajorLightPurple,
		AccentColor:   eMajorDarkPurple,
	},
	"dfs": {
		Name:          "B Major (DFS)",
		Pattern:       "dfs",
		Description:   "The deep explorer, rich and thorough",
		PrimaryColor:  bMajorDeepBlue,
		SecondaryColor: bMajorMediumBlue,
		AccentColor:   bMajorLightBlue,
	},
	"bfs": {
		Name:          "F# Major (BFS)",
		Pattern:       "bfs",
		Description:   "The level-by-level discoverer, methodical and complete",
		PrimaryColor:  fSharpMajorTeal,
		SecondaryColor: fSharpMajorLightTeal,
		AccentColor:   fSharpMajorDarkTeal,
	},
	"dynamic-programming": {
		Name:          "Db Major (Dynamic Programming)",
		Pattern:       "dynamic-programming",
		Description:   "The optimizer, complex and powerful",
		PrimaryColor:  dbMajorYellow,
		SecondaryColor: dbMajorLightYellow,
		AccentColor:   dbMajorDarkYellow,
	},
	"greedy": {
		Name:          "Ab Major (Greedy)",
		Pattern:       "greedy",
		Description:   "The local maximizer, bold and decisive",
		PrimaryColor:  abMajorMagenta,
		SecondaryColor: abMajorLightMagenta,
		AccentColor:   abMajorDarkMagenta,
	},
	"union-find": {
		Name:          "Eb Major (Union-Find)",
		Pattern:       "union-find",
		Description:   "The connector, structured and organized",
		PrimaryColor:  ebMajorCyan,
		SecondaryColor: ebMajorLightCyan,
		AccentColor:   ebMajorDarkCyan,
	},
	"heap": {
		Name:          "Bb Major (Heap / Priority Queue)",
		Pattern:       "heap",
		Description:   "The sorter, flexible and maintaining order",
		PrimaryColor:  bbMajorBrown,
		SecondaryColor: bbMajorLightBrown,
		AccentColor:   bbMajorDarkBrown,
	},
}

// Base styles for UI components
var (
	// Text styles
	BaseStyle = lipgloss.NewStyle().
		Foreground(defaultFg).
		Background(defaultBg)

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1).
		Width(60).
		Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BBDEFB")).
		Italic(true).
		MarginBottom(1)

	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A5D6A7")).
		Background(lipgloss.Color("#1B5E20")).
		Padding(0, 1).
		MarginTop(1)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2E7D32")).
		Bold(true).
		Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#C62828")).
		Bold(true).
		Padding(0, 1)

	WarningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#F9A825")).
		Bold(true).
		Padding(0, 1)

	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#1565C0")).
		Padding(0, 1)

	// Box styles
	BorderedBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#64B5F6")).
		Padding(1).
		MarginTop(1).
		MarginBottom(1)

	CodeBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF9800")).
		Background(lipgloss.Color("#263238")).
		Foreground(lipgloss.Color("#ECEFF1")).
		Padding(1).
		MarginTop(1).
		MarginBottom(1)

	ProblemBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#64B5F6")).
		Background(lipgloss.Color("#1A237E")).
		Foreground(lipgloss.Color("#E8EAF6")).
		Padding(1).
		MarginTop(1).
		MarginBottom(1)

	HorizontalLine = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64B5F6")).
		Render("─────────────────────────────────────")

	// Menu styles
	FocusedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#6200EA")).
		Bold(true).
		Padding(0, 1)

	UnfocusedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E1F5FE")).
		Padding(0, 1)

	MenuBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#B39DDB")).
		Padding(1).
		Width(60)

	// Status bar style
	StatusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#303F9F")).
		Bold(true).
		Padding(0, 1)
		
	TimerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#F57C00")).
		Bold(true).
		Padding(0, 1)
		
	TimerWarningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#D32F2F")).
		Bold(true).
		Padding(0, 1)
	
	// Button styles	
	ButtonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#1976D2")).
		Padding(0, 2).
		MarginRight(1).
		Bold(true)
		
	ActiveButtonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#673AB7")).
		Padding(0, 2).
		MarginRight(1).
		Bold(true)
		
	HeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B3E5FC")).
		Bold(true)
)

// GetPatternStyle returns styles for a specific algorithm pattern
func GetPatternStyle(pattern string) (lipgloss.Style, lipgloss.Style, lipgloss.Style) {
	// Get the musical scale for this pattern
	scale, ok := MusicScales[pattern]
	if !ok {
		// Default to C Major if pattern not found
		scale = MusicScales["sliding-window"]
	}

	// Create styles based on the scale's colors
	primaryStyle := lipgloss.NewStyle().
		Foreground(scale.PrimaryColor)

	secondaryStyle := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor)

	accentStyle := lipgloss.NewStyle().
		Foreground(scale.AccentColor)

	return primaryStyle, secondaryStyle, accentStyle
}

// Musical note spinners
var (
	MusicNoteSpinner = []string{"♩", "♪", "♫", "♬"}
)

// Helper function to create a progress bar
func ProgressBar(width int, percent float64, pattern string) string {
	scale, ok := MusicScales[pattern]
	if !ok {
		scale = MusicScales["sliding-window"]
	}

	// Calculate the number of filled blocks
	filledWidth := int(float64(width) * percent)
	if filledWidth > width {
		filledWidth = width
	}

	// Create the filled part of the progress bar
	filledStyle := lipgloss.NewStyle().
		Foreground(scale.PrimaryColor).
		Background(scale.PrimaryColor)
	
	// Create the empty part of the progress bar
	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#333333")).
		Background(lipgloss.Color("#333333"))

	filled := filledStyle.Render(strings.Repeat("█", filledWidth))
	empty := emptyStyle.Render(strings.Repeat("░", width-filledWidth))

	return filled + empty
}