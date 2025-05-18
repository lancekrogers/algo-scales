package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Common colors
var (
	primaryColor   = lipgloss.Color("62")  // Cyan
	secondaryColor = lipgloss.Color("212") // Light blue
	successColor   = lipgloss.Color("46")  // Green
	warningColor   = lipgloss.Color("214") // Orange
	errorColor     = lipgloss.Color("196") // Red
	mutedColor     = lipgloss.Color("241") // Gray
	darkGray       = lipgloss.Color("238")
	lightGray      = lipgloss.Color("245")
	backgroundColor = lipgloss.Color("235")
)

// Common styles
var (
	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(2)

	// Subtitle styles
	subtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor).
			MarginBottom(1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(2)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(errorColor)

	// Success style
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor)

	// Warning style
	warningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(warningColor)

	// Selected item style
	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(secondaryColor)

	// Cursor style
	cursorStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Box styles
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	successBoxStyle = boxStyle.Copy().
			BorderForeground(successColor)

	warningBoxStyle = boxStyle.Copy().
			BorderForeground(warningColor)

	// Code block style
	codeBlockStyle = lipgloss.NewStyle().
			Foreground(lightGray).
			Background(backgroundColor).
			Padding(0, 1)

	// Loading style
	loadingStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	// Progress bar styles
	progressBarStyle = lipgloss.NewStyle().
				Foreground(successColor)

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(darkGray)

	// Timer styles
	timerNormalStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(successColor)

	timerWarningStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(warningColor)

	timerDangerStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(errorColor)

	// Difficulty styles
	easyStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	mediumStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	hardStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(backgroundColor).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 2)

	// Button styles
	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(primaryColor).
			Padding(0, 2).
			MarginRight(1)

	disabledButtonStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Background(darkGray).
				Padding(0, 2).
				MarginRight(1)
)

// Helper functions
func getDifficultyStyle(difficulty string) lipgloss.Style {
	switch difficulty {
	case "easy":
		return easyStyle
	case "medium":
		return mediumStyle
	case "hard":
		return hardStyle
	default:
		return lipgloss.NewStyle()
	}
}

func getTimerStyle(minutes int) lipgloss.Style {
	if minutes > 30 {
		return timerDangerStyle
	} else if minutes > 20 {
		return timerWarningStyle
	}
	return timerNormalStyle
}

// Loading spinner frames
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Get spinner frame based on time
func getSpinnerFrame(counter int) string {
	return spinnerFrames[counter%len(spinnerFrames)]
}

// Progress bar helper
func renderProgressBar(progress float64, width int) string {
	filledWidth := int(float64(width) * progress)
	emptyWidth := width - filledWidth

	filled := progressBarStyle.Render(stringRepeat("█", filledWidth))
	empty := progressEmptyStyle.Render(stringRepeat("░", emptyWidth))

	return filled + empty
}

// String repeat helper
func stringRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}