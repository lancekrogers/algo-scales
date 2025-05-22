package stats

import (
	"sort"
	"time"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// Service implements the StatsService interface
type Service struct {
	storage interfaces.StatsStorage
}

// NewService creates a new stats service with the default file storage
func NewService() *Service {
	return &Service{
		storage: NewFileStorage(),
	}
}

// WithStorage sets a custom storage for the service
func (s *Service) WithStorage(storage interfaces.StatsStorage) *Service {
	s.storage = storage
	return s
}

// RecordSession records a session's statistics
func (s *Service) RecordSession(sessionStats interfaces.SessionStats) error {
	return s.storage.SaveSession(sessionStats)
}

// GetSummary returns summary statistics
func (s *Service) GetSummary() (*interfaces.Summary, error) {
	// Load all session stats
	sessions, err := s.storage.LoadAllSessions()
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return &interfaces.Summary{
			TotalAttempted: 0,
			TotalSolved:    0,
			SuccessRate:    0,
			AvgSolveTime:   "00:00",
			FastestSolve: struct {
				ProblemID string `json:"problem_id"`
				Time      string `json:"time"`
			}{
				ProblemID: "",
				Time:      "",
			},
			MostChallenging: struct {
				ProblemID string `json:"problem_id"`
				Attempts  int    `json:"attempts"`
			}{
				ProblemID: "",
				Attempts:  0,
			},
		}, nil
	}

	// Calculate summary stats
	summary := &interfaces.Summary{}
	summary.TotalAttempted = len(sessions)

	var totalSolveTime time.Duration
	var solvedCount int

	// Track problem attempts
	problemAttempts := make(map[string]int)
	problemSolved := make(map[string]bool)

	// Track fastest solve
	var fastestTime time.Duration
	var fastestProblem string

	for _, session := range sessions {
		problemAttempts[session.ProblemID]++

		if session.Solved {
			solvedCount++
			totalSolveTime += session.Duration

			// Update fastest solve if this is faster
			if fastestProblem == "" || session.Duration < fastestTime {
				fastestTime = session.Duration
				fastestProblem = session.ProblemID
			}

			problemSolved[session.ProblemID] = true
		}
	}

	summary.TotalSolved = solvedCount
	if solvedCount > 0 {
		avgTime := totalSolveTime / time.Duration(solvedCount)
		summary.AvgSolveTime = formatDuration(avgTime)
	}

	if summary.TotalAttempted > 0 {
		summary.SuccessRate = float64(solvedCount) / float64(summary.TotalAttempted) * 100
	}

	// Find most challenging problem
	var maxAttempts int
	var mostChallengingProblem string

	for problem, attempts := range problemAttempts {
		if attempts > maxAttempts && !problemSolved[problem] {
			maxAttempts = attempts
			mostChallengingProblem = problem
		}
	}

	if fastestProblem != "" {
		summary.FastestSolve.ProblemID = fastestProblem
		summary.FastestSolve.Time = formatDuration(fastestTime)
	}

	if mostChallengingProblem != "" {
		summary.MostChallenging.ProblemID = mostChallengingProblem
		summary.MostChallenging.Attempts = maxAttempts
	}

	return summary, nil
}

// GetByPattern returns statistics by pattern
func (s *Service) GetByPattern() (map[string]interfaces.PatternStats, error) {
	// Load all session stats
	sessions, err := s.storage.LoadAllSessions()
	if err != nil {
		return nil, err
	}

	// Group by pattern
	patternStats := make(map[string]interfaces.PatternStats)

	// Track pattern solve times
	patternTimes := make(map[string][]time.Duration)

	for _, session := range sessions {
		for _, pattern := range session.Patterns {
			// Initialize if not exists
			if _, ok := patternStats[pattern]; !ok {
				patternStats[pattern] = interfaces.PatternStats{
					Pattern: pattern,
				}
			}

			// Update stats
			stats := patternStats[pattern]
			stats.Attempted++

			if session.Solved {
				stats.Solved++
				patternTimes[pattern] = append(patternTimes[pattern], session.Duration)
			}

			// Calculate success rate
			stats.SuccessRate = float64(stats.Solved) / float64(stats.Attempted) * 100

			// Update in map
			patternStats[pattern] = stats
		}
	}

	// Calculate average solve times
	for pattern, times := range patternTimes {
		if len(times) > 0 {
			var total time.Duration
			for _, t := range times {
				total += t
			}
			avgTime := total / time.Duration(len(times))

			stats := patternStats[pattern]
			stats.AvgTime = formatDuration(avgTime)
			patternStats[pattern] = stats
		}
	}

	return patternStats, nil
}

// GetTrends returns usage trends over time
func (s *Service) GetTrends() (*interfaces.Trends, error) {
	// Load all session stats
	sessions, err := s.storage.LoadAllSessions()
	if err != nil {
		return nil, err
	}

	// Sort sessions by start time
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.Before(sessions[j].StartTime)
	})

	// Group by day and week
	trends := &interfaces.Trends{}
	dailyStats := make(map[string]struct {
		Solved    int
		TotalTime time.Duration
		Count     int
	})

	weeklyStats := make(map[string]struct {
		Solved    int
		Attempted int
		StartDate time.Time
		EndDate   time.Time
	})

	for _, session := range sessions {
		// Daily stats
		dateStr := session.StartTime.Format("2006-01-02")
		daily := dailyStats[dateStr]
		if session.Solved {
			daily.Solved++
			daily.TotalTime += session.Duration
			daily.Count++
		}
		dailyStats[dateStr] = daily

		// Weekly stats
		yearWeek := getYearWeek(session.StartTime)
		weekly := weeklyStats[yearWeek]
		if weekly.StartDate.IsZero() || session.StartTime.Before(weekly.StartDate) {
			weekly.StartDate = startOfWeek(session.StartTime)
		}
		if weekly.EndDate.IsZero() || session.StartTime.After(weekly.EndDate) {
			weekly.EndDate = endOfWeek(session.StartTime)
		}
		weekly.Attempted++
		if session.Solved {
			weekly.Solved++
		}
		weeklyStats[yearWeek] = weekly
	}

	// Convert maps to slices for the last 7 days
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		dailyTrend := interfaces.DailyTrend{
			Date:   dateStr,
			Solved: 0,
		}

		// Add data if available
		if stats, ok := dailyStats[dateStr]; ok {
			dailyTrend.Solved = stats.Solved
			if stats.Count > 0 {
				avgTime := stats.TotalTime / time.Duration(stats.Count)
				dailyTrend.AvgTime = formatDuration(avgTime)
			}
		}

		trends.Daily = append(trends.Daily, dailyTrend)
	}

	// Sort daily trends by date
	sort.Slice(trends.Daily, func(i, j int) bool {
		return trends.Daily[i].Date < trends.Daily[j].Date
	})

	// Convert weekly stats map to slice
	for _, stats := range weeklyStats {
		weeklyTrend := interfaces.WeeklyTrend{
			StartDate:   stats.StartDate.Format("2006-01-02"),
			EndDate:     stats.EndDate.Format("2006-01-02"),
			Solved:      stats.Solved,
			SuccessRate: 0,
		}

		if stats.Attempted > 0 {
			weeklyTrend.SuccessRate = float64(stats.Solved) / float64(stats.Attempted) * 100
		}

		trends.Weekly = append(trends.Weekly, weeklyTrend)
	}

	// Sort weekly trends by start date
	sort.Slice(trends.Weekly, func(i, j int) bool {
		return trends.Weekly[i].StartDate < trends.Weekly[j].StartDate
	})

	// Limit to last 4 weeks
	if len(trends.Weekly) > 4 {
		trends.Weekly = trends.Weekly[len(trends.Weekly)-4:]
	}

	return trends, nil
}

// Reset resets all statistics
func (s *Service) Reset() error {
	return s.storage.ClearAllSessions()
}

// GetAllSessions returns all recorded sessions
func (s *Service) GetAllSessions() ([]interfaces.SessionStats, error) {
	sessions, err := s.storage.LoadAllSessions()
	if err != nil {
		return nil, err
	}
	
	// Convert to interfaces types
	result := make([]interfaces.SessionStats, len(sessions))
	for i, session := range sessions {
		result[i] = interfaces.SessionStats{
			ProblemID:    session.ProblemID,
			StartTime:    session.StartTime,
			EndTime:      session.EndTime,
			Duration:     session.Duration,
			Solved:       session.Solved,
			Mode:         session.Mode,
			HintsUsed:    session.HintsUsed,
			SolutionUsed: session.SolutionUsed,
			Patterns:     session.Patterns,
			Difficulty:   session.Difficulty,
		}
	}
	return result, nil
}