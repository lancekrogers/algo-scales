package daily

// Scale represents a pattern practice "scale" from musical scales
type Scale struct {
	Pattern     string
	MusicalName string
	Description string
}

// Scales is a slice of all algorithm pattern scales
var Scales = []Scale{
	{
		Pattern:     "sliding-window",
		MusicalName: "C Major",
		Description: "The fundamental scale, elegant and versatile",
	},
	{
		Pattern:     "two-pointers",
		MusicalName: "G Major",
		Description: "Balanced and efficient, the workhorse of array manipulation",
	},
	{
		Pattern:     "fast-slow-pointers",
		MusicalName: "D Major",
		Description: "The cycle detector, bright and revealing",
	},
	{
		Pattern:     "hash-map",
		MusicalName: "A Major",
		Description: "The lookup accelerator, crisp and direct",
	},
	{
		Pattern:     "binary-search",
		MusicalName: "E Major",
		Description: "The divider and conqueror, precise and logarithmic",
	},
	{
		Pattern:     "dfs",
		MusicalName: "B Major",
		Description: "The deep explorer, rich and thorough",
	},
	{
		Pattern:     "bfs",
		MusicalName: "F# Major",
		Description: "The level-by-level discoverer, methodical and complete",
	},
	{
		Pattern:     "dynamic-programming",
		MusicalName: "Db Major",
		Description: "The optimizer, complex and powerful",
	},
	{
		Pattern:     "greedy",
		MusicalName: "Ab Major",
		Description: "The local maximizer, bold and decisive",
	},
	{
		Pattern:     "union-find",
		MusicalName: "Eb Major",
		Description: "The connector, structured and organized",
	},
	{
		Pattern:     "heap",
		MusicalName: "Bb Major",
		Description: "The sorter, flexible and maintaining order",
	},
}

// GetNextScale finds the next scale to practice based on completed patterns
func GetNextScale(completed []string) *Scale {
	for _, scale := range Scales {
		if !Contains(completed, scale.Pattern) {
			return &scale
		}
	}
	return nil
}

// GetScaleByPattern returns a scale by its pattern name
func GetScaleByPattern(pattern string) *Scale {
	for _, scale := range Scales {
		if scale.Pattern == pattern {
			return &scale
		}
	}
	return nil
}

// GetPatternIndex returns the index of a pattern in the scales slice
func GetPatternIndex(pattern string) int {
	for i, scale := range Scales {
		if scale.Pattern == pattern {
			return i
		}
	}
	return -1
}

// GetRemainingPatterns returns the number of patterns not yet completed
func GetRemainingPatterns(completed []string) int {
	count := 0
	for _, scale := range Scales {
		if !Contains(completed, scale.Pattern) {
			count++
		}
	}
	return count
}