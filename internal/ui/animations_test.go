package ui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAnimation(t *testing.T) {
	anim := NewAnimation(AnimationFadeIn, 500*time.Millisecond)
	
	assert.Equal(t, AnimationFadeIn, anim.Type)
	assert.Equal(t, 500*time.Millisecond, anim.Duration)
	assert.False(t, anim.Complete)
	assert.Equal(t, float64(0), anim.Progress)
	assert.NotZero(t, anim.StartTime)
}

func TestAnimation_Update(t *testing.T) {
	anim := Animation{
		Type:      AnimationFadeIn,
		Duration:  100 * time.Millisecond,
		StartTime: time.Now(),
	}
	
	// Initial state
	assert.Equal(t, float64(0), anim.Progress)
	assert.False(t, anim.Complete)
	
	// Wait a bit and update
	time.Sleep(50 * time.Millisecond)
	anim.Update()
	
	// Progress should be around 0.5
	assert.True(t, anim.Progress > 0.3)
	assert.True(t, anim.Progress < 0.7)
	assert.False(t, anim.Complete)
	
	// Wait for completion
	time.Sleep(60 * time.Millisecond)
	anim.Update()
	
	// Should be complete
	assert.Equal(t, float64(1.0), anim.Progress)
	assert.True(t, anim.Complete)
	
	// Update on complete animation should not change anything
	anim.Update()
	assert.Equal(t, float64(1.0), anim.Progress)
	assert.True(t, anim.Complete)
}

func TestAnimation_Apply(t *testing.T) {
	testCases := []struct {
		name     string
		animType AnimationType
		progress float64
		complete bool
		content  string
		expected func(result string) bool
	}{
		{
			name:     "no animation",
			animType: AnimationNone,
			progress: 0.5,
			content:  "Test Content",
			expected: func(result string) bool {
				return result == "Test Content"
			},
		},
		{
			name:     "completed animation",
			animType: AnimationFadeIn,
			progress: 1.0,
			complete: true,
			content:  "Test Content",
			expected: func(result string) bool {
				return result == "Test Content"
			},
		},
		{
			name:     "fade in start",
			animType: AnimationFadeIn,
			progress: 0.0,
			content:  "Line 1\nLine 2\nLine 3",
			expected: func(result string) bool {
				return result == ""
			},
		},
		{
			name:     "fade in middle",
			animType: AnimationFadeIn,
			progress: 0.25,
			content:  "Line 1\nLine 2\nLine 3",
			expected: func(result string) bool {
				return result == "Line 1"
			},
		},
		{
			name:     "fade in end",
			animType: AnimationFadeIn,
			progress: 0.75,
			content:  "Line 1\nLine 2\nLine 3",
			expected: func(result string) bool {
				return result == "Line 1\nLine 2\nLine 3"
			},
		},
		{
			name:     "slide left start",
			animType: AnimationSlideLeft,
			progress: 0.0,
			content:  "Test",
			expected: func(result string) bool {
				// Should have padding on the left
				return len(result) > len("Test")
			},
		},
		{
			name:     "slide left complete",
			animType: AnimationSlideLeft,
			progress: 1.0,
			content:  "Test",
			expected: func(result string) bool {
				// Should have no padding
				return result == "Test"
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			anim := Animation{
				Type:     tc.animType,
				Progress: tc.progress,
				Complete: tc.complete,
			}
			
			result := anim.Apply(tc.content, 80, 24)
			
			assert.True(t, tc.expected(result))
		})
	}
}

func TestAnimationTick(t *testing.T) {
	cmd := AnimationTick()
	assert.NotNil(t, cmd)
}

func TestPulseIndicator(t *testing.T) {
	indicator := NewPulseIndicator("●", primaryColor)
	
	assert.Equal(t, "●", indicator.symbol)
	assert.Equal(t, primaryColor, indicator.color)
	assert.Equal(t, 0, indicator.frame)
	
	// Update should increment frame
	indicator.Update()
	assert.Equal(t, 1, indicator.frame)
	
	// View should return styled content
	view := indicator.View()
	assert.NotEmpty(t, view)
}

func TestSelectionHighlight(t *testing.T) {
	highlight := NewSelectionHighlight()
	
	assert.False(t, highlight.active)
	assert.Equal(t, 0, highlight.frame)
	
	// Set active
	highlight.SetActive(true)
	assert.True(t, highlight.active)
	assert.Equal(t, 0, highlight.frame)
	
	// Update should increment frame when active
	highlight.Update()
	assert.Equal(t, 1, highlight.frame)
	
	// Apply when inactive should return content as-is
	highlight.SetActive(false)
	result := highlight.Apply("Test")
	assert.Equal(t, "Test", result)
	
	// Apply when active should style content
	highlight.SetActive(true)
	result = highlight.Apply("Test")
	// The result should at least contain the original content
	assert.Contains(t, result, "Test")
	// The result might be styled differently
	assert.NotNil(t, result)
}