package view

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// CustomSpinners defines specialized spinners for Algo Scales
type CustomSpinners struct {
	// Musical note spinner frames
	MusicNotes spinner.Spinner

	// Pattern-specific spinners
	SlidingWindow spinner.Spinner
	TwoPointers   spinner.Spinner
	FastSlow      spinner.Spinner
	DFS           spinner.Spinner
	BFS           spinner.Spinner
	
	// Timer spinners
	TimerSpinner spinner.Spinner
}

// NewCustomSpinners creates a set of custom spinners
func NewCustomSpinners() CustomSpinners {
	return CustomSpinners{
		// Musical notes spinner (default)
		MusicNotes: spinner.Spinner{
			Frames: []string{"♩", "♪", "♫", "♬", "♫", "♪"},
			FPS:    time.Second / 6,
		},
		
		// Sliding window spinner (moves a window across)
		SlidingWindow: spinner.Spinner{
			Frames: []string{
				"[⬜⬜⬜⬜⬜]",
				"[⬜⬜⬜⬜⬜]",
				"[🟦⬜⬜⬜⬜]",
				"[🟦🟦⬜⬜⬜]",
				"[🟦🟦🟦⬜⬜]",
				"[⬜🟦🟦🟦⬜]",
				"[⬜⬜🟦🟦🟦]",
				"[⬜⬜⬜🟦🟦]",
				"[⬜⬜⬜⬜🟦]",
				"[⬜⬜⬜⬜⬜]",
			},
			FPS: time.Second / 5,
		},
		
		// Two pointers spinner (converging arrows)
		TwoPointers: spinner.Spinner{
			Frames: []string{
				"[◀️     ▶️]",
				"[ ◀️   ▶️ ]",
				"[  ◀️ ▶️  ]",
				"[   ⚡   ]",
				"[  ◀️ ▶️  ]",
				"[ ◀️   ▶️ ]",
				"[◀️     ▶️]",
			},
			FPS: time.Second / 4,
		},
		
		// Fast-slow spinner (two dots moving at different speeds)
		FastSlow: spinner.Spinner{
			Frames: []string{
				"[🔵     🔴]",
				"[🔵    🔴 ]",
				"[🔵   🔴  ]",
				"[🔵  🔴   ]",
				"[🔵 🔴    ]",
				"[🔵🔴     ]",
				"[🔴🔵     ]",
				"[🔴 🔵    ]",
				"[🔴  🔵   ]",
				"[🔴   🔵  ]",
				"[🔴    🔵 ]",
				"[🔴     🔵]",
			},
			FPS: time.Second / 8,
		},
		
		// DFS spinner (descending levels)
		DFS: spinner.Spinner{
			Frames: []string{
				"[↓     ]",
				"[↓→    ]",
				"[↓→↓   ]",
				"[↓→↓→  ]",
				"[↓→↓→↓ ]",
				"[↓→↓→↓→]",
				"[↑→↓→↓→]",
				"[↑←↑→↓→]",
				"[↑←↑←↑→]",
				"[↑←↑←↑←]",
				"[↑←↑←  ]",
				"[↑←    ]",
			},
			FPS: time.Second / 6,
		},
		
		// BFS spinner (expanding levels)
		BFS: spinner.Spinner{
			Frames: []string{
				"[●      ]",
				"[●→●    ]",
				"[●→●→●  ]",
				"[●→●→●→●]",
				"[●→●→●  ]",
				"[●→●    ]",
				"[●      ]",
			},
			FPS: time.Second / 4,
		},
		
		// Timer spinner (for countdowns)
		TimerSpinner: spinner.Spinner{
			Frames: []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"},
			FPS:    time.Second,
		},
	}
}

// GetPatternSpinner returns a spinner for a specific algorithm pattern
func (s CustomSpinners) GetPatternSpinner(pattern string) spinner.Spinner {
	switch pattern {
	case "sliding-window":
		return s.SlidingWindow
	case "two-pointers":
		return s.TwoPointers
	case "fast-slow-pointers":
		return s.FastSlow
	case "dfs":
		return s.DFS
	case "bfs":
		return s.BFS
	default:
		return s.MusicNotes
	}
}

// StyleSpinnerForPattern styles a spinner for a specific pattern
func StyleSpinnerForPattern(s spinner.Model, pattern string) spinner.Model {
	// Get colors for this pattern
	scale, ok := MusicScales[pattern]
	if !ok {
		scale = MusicScales["sliding-window"] // Default to C Major
	}
	
	// Apply the pattern color to the spinner
	s.Style = lipgloss.NewStyle().Foreground(scale.PrimaryColor)
	return s
}

// ProgressIndicator creates a progress bar with the pattern's theme
func ProgressIndicator(width int, percent float64, pattern string) string {
	// Get the musical scale for this pattern
	scale, ok := MusicScales[pattern]
	if !ok {
		// Default to C Major if pattern not found
		scale = MusicScales["sliding-window"]
	}
	
	// Calculate filled and empty portions
	filled := int(float64(width) * percent)
	if filled > width {
		filled = width
	}
	
	// Create styled bar segments
	filledStyle := lipgloss.NewStyle().Foreground(scale.PrimaryColor)
	emptyStyle := lipgloss.NewStyle().Foreground(subtleGray)
	
	// Generate progress bar
	bar := filledStyle.Render(strings.Repeat("█", filled)) + 
		   emptyStyle.Render(strings.Repeat("░", width-filled))
	
	// Add percentage
	percentText := lipgloss.NewStyle().
		Foreground(scale.SecondaryColor).
		Render(fmt.Sprintf(" %d%%", int(percent*100)))
	
	return bar + percentText
}

// LoadingBar creates an animated loading bar
func LoadingBar(width int, progress float64, pattern string) string {
	// Get colors for this pattern
	scale, ok := MusicScales[pattern]
	if !ok {
		scale = MusicScales["sliding-window"] // Default to C Major
	}
	
	// Calculate number of blocks to fill
	numBlocks := width
	numFilled := int(float64(numBlocks) * progress)
	
	// Generate the bar
	var bar strings.Builder
	
	// Filled blocks
	fillStyle := lipgloss.NewStyle().Foreground(scale.PrimaryColor)
	for i := 0; i < numFilled; i++ {
		bar.WriteString(fillStyle.Render("█"))
	}
	
	// Empty blocks
	emptyStyle := lipgloss.NewStyle().Foreground(subtleGray)
	for i := numFilled; i < numBlocks; i++ {
		bar.WriteString(emptyStyle.Render("░"))
	}
	
	// Add musical note at the end based on progress
	if progress < 1.0 {
		// Calculate which note to show based on progress
		noteIndex := int(progress * 4) % 4
		notes := []string{"♩", "♪", "♫", "♬"}
		note := notes[noteIndex]
		
		// Replace the last character with the note
		barString := bar.String()
		if len(barString) > 0 {
			barRunes := []rune(barString)
			barRunes[numFilled-1] = []rune(note)[0]
			return string(barRunes)
		}
	}
	
	return bar.String()
}