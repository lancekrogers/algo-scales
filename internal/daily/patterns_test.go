package daily

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNextScale(t *testing.T) {
	tests := []struct {
		name      string
		completed []string
		want      *Scale
		wantNil   bool
	}{
		{
			name:      "no completed patterns",
			completed: []string{},
			want:      &Scales[0], // Should return the first scale
		},
		{
			name:      "some completed patterns",
			completed: []string{"sliding-window", "two-pointers"},
			want:      &Scales[2], // Should return the third scale (fast-slow-pointers)
		},
		{
			name:      "all completed patterns",
			completed: []string{
				"sliding-window", "two-pointers", "fast-slow-pointers",
				"hash-map", "binary-search", "dfs", "bfs",
				"dynamic-programming", "greedy", "union-find", "heap",
			},
			wantNil: true, // Should return nil when all scales are completed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNextScale(tt.completed)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.Pattern, result.Pattern)
				assert.Equal(t, tt.want.MusicalName, result.MusicalName)
			}
		})
	}
}

func TestGetScaleByPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    *Scale
		wantNil bool
	}{
		{
			name:    "existing pattern",
			pattern: "sliding-window",
			want:    &Scales[0],
		},
		{
			name:    "another existing pattern",
			pattern: "dynamic-programming",
			want:    &Scales[7],
		},
		{
			name:    "non-existent pattern",
			pattern: "not-a-real-pattern",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetScaleByPattern(tt.pattern)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.Pattern, result.Pattern)
				assert.Equal(t, tt.want.MusicalName, result.MusicalName)
				assert.Equal(t, tt.want.Description, result.Description)
			}
		})
	}
}

func TestGetPatternIndex(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected int
	}{
		{
			name:     "first pattern",
			pattern:  "sliding-window",
			expected: 0,
		},
		{
			name:     "middle pattern",
			pattern:  "bfs",
			expected: 6,
		},
		{
			name:     "last pattern",
			pattern:  "heap",
			expected: 10,
		},
		{
			name:     "non-existent pattern",
			pattern:  "not-a-real-pattern",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPatternIndex(tt.pattern)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRemainingPatterns(t *testing.T) {
	tests := []struct {
		name      string
		completed []string
		expected  int
	}{
		{
			name:      "no completed patterns",
			completed: []string{},
			expected:  11, // All patterns remain
		},
		{
			name:      "some completed patterns",
			completed: []string{"sliding-window", "two-pointers", "hash-map"},
			expected:  8, // 11 - 3 = 8 patterns remain
		},
		{
			name: "all completed patterns",
			completed: []string{
				"sliding-window", "two-pointers", "fast-slow-pointers",
				"hash-map", "binary-search", "dfs", "bfs",
				"dynamic-programming", "greedy", "union-find", "heap",
			},
			expected: 0, // No patterns remain
		},
		{
			name:      "duplicate completed patterns",
			completed: []string{"sliding-window", "sliding-window", "two-pointers"},
			expected:  9, // Only unique patterns are considered
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRemainingPatterns(tt.completed)
			assert.Equal(t, tt.expected, result)
		})
	}
}