package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingScreen displays a loading animation with a message
type LoadingScreen struct {
	message      string
	spinnerFrame int
	width        int
	height       int
}

// NewLoadingScreen creates a new loading screen
func NewLoadingScreen(message string) LoadingScreen {
	return LoadingScreen{
		message: message,
	}
}

// Update handles loading screen updates
func (l LoadingScreen) Update(msg tea.Msg) (LoadingScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width
		l.height = msg.Height
		return l, nil

	case spinnerTickMsg:
		l.spinnerFrame++
		return l, tickSpinner()
	}

	return l, nil
}

// View renders the loading screen
func (l LoadingScreen) View() string {
	if l.width == 0 || l.height == 0 {
		return ""
	}

	var b strings.Builder

	// Center content vertically
	topPadding := (l.height - 5) / 2
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Spinner
	spinner := getSpinnerFrame(l.spinnerFrame)
	spinnerLine := fmt.Sprintf("%s %s", spinner, l.message)
	
	// Center horizontally
	leftPadding := (l.width - lipgloss.Width(spinnerLine)) / 2
	if leftPadding > 0 {
		spinnerLine = strings.Repeat(" ", leftPadding) + spinnerLine
	}

	b.WriteString(loadingStyle.Render(spinnerLine))

	return b.String()
}

// Transitions

// FadeTransition creates a fade effect between screens
type FadeTransition struct {
	from      string
	to        string
	progress  float64
	duration  time.Duration
	startTime time.Time
}

// NewFadeTransition creates a new fade transition
func NewFadeTransition(from, to string, duration time.Duration) FadeTransition {
	return FadeTransition{
		from:      from,
		to:        to,
		duration:  duration,
		startTime: time.Now(),
	}
}

// Update handles transition updates
func (f FadeTransition) Update(msg tea.Msg) (FadeTransition, tea.Cmd) {
	switch msg.(type) {
	case transitionTickMsg:
		elapsed := time.Since(f.startTime)
		f.progress = float64(elapsed) / float64(f.duration)
		
		if f.progress >= 1.0 {
			f.progress = 1.0
			return f, nil
		}
		
		return f, tickTransition()
	}
	
	return f, nil
}

// View renders the transition
func (f FadeTransition) View() string {
	if f.progress >= 1.0 {
		return f.to
	}
	
	// Simple fade by showing the target with increasing opacity
	// In a real implementation, we'd blend the two screens
	if f.progress > 0.5 {
		return f.to
	}
	return f.from
}

// isDone checks if the transition is complete
func (f FadeTransition) isDone() bool {
	return f.progress >= 1.0
}

// Message types
type spinnerTickMsg time.Time
type transitionTickMsg time.Time

// Commands
func tickSpinner() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}

func tickTransition() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { // ~60fps
		return transitionTickMsg(t)
	})
}