package ui

// State represents the current UI state
type State int

const (
	StateHome State = iota
	StatePatternSelection
	StateProblemList
	StateProblemDetail
	StateSession
	StateStats
	StateDaily
	StateSettings
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateHome:
		return "home"
	case StatePatternSelection:
		return "pattern_selection"
	case StateProblemList:
		return "problem_list"
	case StateProblemDetail:
		return "problem_detail"
	case StateSession:
		return "session"
	case StateStats:
		return "stats"
	case StateDaily:
		return "daily"
	case StateSettings:
		return "settings"
	default:
		return "unknown"
	}
}