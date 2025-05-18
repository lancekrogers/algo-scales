package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateString(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateHome, "home"},
		{StatePatternSelection, "pattern_selection"},
		{StateProblemList, "problem_list"},
		{StateProblemDetail, "problem_detail"},
		{StateSession, "session"},
		{StateStats, "stats"},
		{StateDaily, "daily"},
		{StateSettings, "settings"},
		{State(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.state.String())
		})
	}
}