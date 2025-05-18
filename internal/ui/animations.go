package ui

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AnimationType defines the type of animation
type AnimationType int

const (
	AnimationNone AnimationType = iota
	AnimationFadeIn
	AnimationFadeOut
	AnimationSlideLeft
	AnimationSlideRight
	AnimationSlideUp
	AnimationSlideDown
	AnimationPulse
	AnimationBounce
)

// Animation manages screen animations
type Animation struct {
	Type      AnimationType
	Duration  time.Duration
	StartTime time.Time
	Progress  float64
	Complete  bool
}

// NewAnimation creates a new animation
func NewAnimation(animType AnimationType, duration time.Duration) Animation {
	return Animation{
		Type:      animType,
		Duration:  duration,
		StartTime: time.Now(),
	}
}

// Update updates the animation progress
func (a *Animation) Update() {
	if a.Complete {
		return
	}

	elapsed := time.Since(a.StartTime)
	a.Progress = float64(elapsed) / float64(a.Duration)

	if a.Progress >= 1.0 {
		a.Progress = 1.0
		a.Complete = true
	}
}

// Apply applies the animation to content
func (a Animation) Apply(content string, width, height int) string {
	if a.Type == AnimationNone || a.Complete {
		return content
	}

	switch a.Type {
	case AnimationFadeIn:
		return a.applyFadeIn(content)
	case AnimationFadeOut:
		return a.applyFadeOut(content)
	case AnimationSlideLeft:
		return a.applySlideLeft(content, width)
	case AnimationSlideRight:
		return a.applySlideRight(content, width)
	case AnimationSlideUp:
		return a.applySlideUp(content, height)
	case AnimationSlideDown:
		return a.applySlideDown(content, height)
	case AnimationPulse:
		return a.applyPulse(content)
	case AnimationBounce:
		return a.applyBounce(content)
	default:
		return content
	}
}

// Fade in animation
func (a Animation) applyFadeIn(content string) string {
	if a.Progress < 0.5 {
		// Show partial content based on progress
		lines := strings.Split(content, "\n")
		visibleLines := int(float64(len(lines)) * (a.Progress * 2))
		if visibleLines > 0 {
			return strings.Join(lines[:visibleLines], "\n")
		}
		return ""
	}
	return content
}

// Fade out animation
func (a Animation) applyFadeOut(content string) string {
	if a.Progress > 0.5 {
		// Hide partial content based on progress
		lines := strings.Split(content, "\n")
		visibleLines := int(float64(len(lines)) * ((1.0 - a.Progress) * 2))
		if visibleLines > 0 {
			return strings.Join(lines[:visibleLines], "\n")
		}
		return ""
	}
	return content
}

// Slide left animation
func (a Animation) applySlideLeft(content string, width int) string {
	offset := int(float64(width) * (1.0 - a.Progress))
	padding := strings.Repeat(" ", offset)
	
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = padding + line
	}
	
	return strings.Join(lines, "\n")
}

// Slide right animation
func (a Animation) applySlideRight(content string, width int) string {
	offset := int(float64(width) * a.Progress)
	padding := strings.Repeat(" ", offset)
	
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = padding + line
	}
	
	return strings.Join(lines, "\n")
}

// Slide up animation
func (a Animation) applySlideUp(content string, height int) string {
	offset := int(float64(height) * (1.0 - a.Progress))
	padding := strings.Repeat("\n", offset)
	return padding + content
}

// Slide down animation
func (a Animation) applySlideDown(content string, height int) string {
	offset := int(float64(height) * a.Progress)
	padding := strings.Repeat("\n", offset)
	return padding + content
}

// Pulse animation
func (a Animation) applyPulse(content string) string {
	// Create a pulsing effect by varying the style intensity
	intensity := 0.5 + 0.5*math.Sin(a.Progress*math.Pi*4)
	
	// Apply a style with varying intensity
	pulseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primaryColor)).
		Bold(intensity > 0.7)
	
	return pulseStyle.Render(content)
}

// Bounce animation
func (a Animation) applyBounce(content string) string {
	// Simple bounce effect
	bounce := math.Abs(math.Sin(a.Progress * math.Pi * 3))
	offset := int(bounce * 3)
	
	lines := strings.Split(content, "\n")
	paddedLines := make([]string, offset)
	paddedLines = append(paddedLines, lines...)
	
	return strings.Join(paddedLines, "\n")
}

// AnimationTick generates animation tick messages
func AnimationTick() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { // ~60fps
		return animationTickMsg(t)
	})
}

// animationTickMsg is sent on animation ticks
type animationTickMsg time.Time

// Visual feedback components

// PulseIndicator creates a pulsing indicator
type PulseIndicator struct {
	symbol   string
	frame    int
	color    lipgloss.Color
}

// NewPulseIndicator creates a new pulse indicator
func NewPulseIndicator(symbol string, color lipgloss.Color) PulseIndicator {
	return PulseIndicator{
		symbol: symbol,
		color:  color,
	}
}

// Update updates the pulse indicator
func (p *PulseIndicator) Update() {
	p.frame++
}

// View renders the pulse indicator
func (p PulseIndicator) View() string {
	intensity := 0.5 + 0.5*math.Sin(float64(p.frame)*0.1)
	
	style := lipgloss.NewStyle().
		Foreground(p.color).
		Bold(intensity > 0.7)
	
	return style.Render(p.symbol)
}

// SelectionHighlight creates a selection highlight effect
type SelectionHighlight struct {
	active bool
	frame  int
}

// NewSelectionHighlight creates a new selection highlight
func NewSelectionHighlight() SelectionHighlight {
	return SelectionHighlight{}
}

// SetActive sets the highlight active state
func (s *SelectionHighlight) SetActive(active bool) {
	s.active = active
	if active {
		s.frame = 0
	}
}

// Update updates the selection highlight
func (s *SelectionHighlight) Update() {
	if s.active {
		s.frame++
	}
}

// Apply applies the highlight to content
func (s SelectionHighlight) Apply(content string) string {
	if !s.active {
		return content
	}
	
	// Create a glowing effect
	intensity := 0.7 + 0.3*math.Sin(float64(s.frame)*0.1)
	
	highlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(primaryColor)).
		Foreground(lipgloss.Color("255")).
		Bold(true)
	
	if intensity < 0.85 {
		highlightStyle = highlightStyle.Faint(true)
	}
	
	return highlightStyle.Render(content)
}