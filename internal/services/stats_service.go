package services

import (
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/stats"
)

// StatsCommandService provides business logic for stats command operations
type StatsCommandService interface {
	// GetOverallStats returns overall performance statistics
	GetOverallStats() (*interfaces.OverallStats, error)
	
	// GetPatternStats returns performance by algorithm pattern
	GetPatternStats() (map[string]*interfaces.PatternStats, error)
	
	// GetDifficultyStats returns performance by difficulty level
	GetDifficultyStats() (map[string]*interfaces.DifficultyStats, error)
	
	// GetLanguageStats returns performance by programming language
	GetLanguageStats() (map[string]*interfaces.LanguageStats, error)
	
	// GetRecentActivity returns recent session activity
	GetRecentActivity(days int) ([]*interfaces.DailyStats, error)
}

// StatsCommandServiceImpl implements StatsCommandService
type StatsCommandServiceImpl struct {
	statsService interfaces.StatsService
}

// NewStatsCommandService creates a new stats command service
func NewStatsCommandService(statsService interfaces.StatsService) StatsCommandService {
	if statsService == nil {
		// Fallback to legacy stats for compatibility
		return &LegacyStatsCommandService{}
	}
	
	return &StatsCommandServiceImpl{
		statsService: statsService,
	}
}

// GetOverallStats returns overall performance statistics
func (s *StatsCommandServiceImpl) GetOverallStats() (*interfaces.OverallStats, error) {
	return s.statsService.GetOverallStats()
}

// GetPatternStats returns performance by algorithm pattern
func (s *StatsCommandServiceImpl) GetPatternStats() (map[string]*interfaces.PatternStats, error) {
	return s.statsService.GetPatternStats()
}

// GetDifficultyStats returns performance by difficulty level
func (s *StatsCommandServiceImpl) GetDifficultyStats() (map[string]*interfaces.DifficultyStats, error) {
	return s.statsService.GetDifficultyStats()
}

// GetLanguageStats returns performance by programming language
func (s *StatsCommandServiceImpl) GetLanguageStats() (map[string]*interfaces.LanguageStats, error) {
	return s.statsService.GetLanguageStats()
}

// GetRecentActivity returns recent session activity
func (s *StatsCommandServiceImpl) GetRecentActivity(days int) ([]*interfaces.DailyStats, error) {
	return s.statsService.GetRecentActivity(days)
}

// LegacyStatsCommandService provides backward compatibility with legacy stats functions
type LegacyStatsCommandService struct{}

// GetOverallStats returns overall performance statistics using legacy functions
func (s *LegacyStatsCommandService) GetOverallStats() (*interfaces.OverallStats, error) {
	legacy, err := stats.GetOverallStats()
	if err != nil {
		return nil, err
	}
	
	// Convert legacy format to interface format
	return &interfaces.OverallStats{
		TotalSessions:     legacy.TotalSessions,
		SolvedProblems:    legacy.SolvedProblems,
		UnsolvedProblems:  legacy.UnsolvedProblems,
		AverageTime:       legacy.AverageTime,
		TotalTime:         legacy.TotalTime,
		SuccessRate:       legacy.SuccessRate,
		CurrentStreak:     legacy.CurrentStreak,
		LongestStreak:     legacy.LongestStreak,
		FavoritePattern:   legacy.FavoritePattern,
		FavoriteLanguage:  legacy.FavoriteLanguage,
	}, nil
}

// GetPatternStats returns performance by algorithm pattern using legacy functions
func (s *LegacyStatsCommandService) GetPatternStats() (map[string]*interfaces.PatternStats, error) {
	legacy, err := stats.GetPatternStats()
	if err != nil {
		return nil, err
	}
	
	// Convert legacy format to interface format
	result := make(map[string]*interfaces.PatternStats)
	for pattern, legacyStats := range legacy {
		result[pattern] = &interfaces.PatternStats{
			Pattern:       legacyStats.Pattern,
			TotalSessions: legacyStats.TotalSessions,
			Solved:        legacyStats.Solved,
			Unsolved:      legacyStats.Unsolved,
			SuccessRate:   legacyStats.SuccessRate,
			AverageTime:   legacyStats.AverageTime,
		}
	}
	
	return result, nil
}

// GetDifficultyStats returns performance by difficulty level using legacy functions
func (s *LegacyStatsCommandService) GetDifficultyStats() (map[string]*interfaces.DifficultyStats, error) {
	legacy, err := stats.GetDifficultyStats()
	if err != nil {
		return nil, err
	}
	
	// Convert legacy format to interface format
	result := make(map[string]*interfaces.DifficultyStats)
	for difficulty, legacyStats := range legacy {
		result[difficulty] = &interfaces.DifficultyStats{
			Difficulty:    legacyStats.Difficulty,
			TotalSessions: legacyStats.TotalSessions,
			Solved:        legacyStats.Solved,
			Unsolved:      legacyStats.Unsolved,
			SuccessRate:   legacyStats.SuccessRate,
			AverageTime:   legacyStats.AverageTime,
		}
	}
	
	return result, nil
}

// GetLanguageStats returns performance by programming language using legacy functions
func (s *LegacyStatsCommandService) GetLanguageStats() (map[string]*interfaces.LanguageStats, error) {
	legacy, err := stats.GetLanguageStats()
	if err != nil {
		return nil, err
	}
	
	// Convert legacy format to interface format
	result := make(map[string]*interfaces.LanguageStats)
	for language, legacyStats := range legacy {
		result[language] = &interfaces.LanguageStats{
			Language:      legacyStats.Language,
			TotalSessions: legacyStats.TotalSessions,
			Solved:        legacyStats.Solved,
			Unsolved:      legacyStats.Unsolved,
			SuccessRate:   legacyStats.SuccessRate,
			AverageTime:   legacyStats.AverageTime,
		}
	}
	
	return result, nil
}

// GetRecentActivity returns recent session activity using legacy functions
func (s *LegacyStatsCommandService) GetRecentActivity(days int) ([]*interfaces.DailyStats, error) {
	legacy, err := stats.GetRecentActivity(days)
	if err != nil {
		return nil, err
	}
	
	// Convert legacy format to interface format
	result := make([]*interfaces.DailyStats, len(legacy))
	for i, legacyStats := range legacy {
		result[i] = &interfaces.DailyStats{
			Date:     legacyStats.Date,
			Sessions: legacyStats.Sessions,
			Solved:   legacyStats.Solved,
			Unsolved: legacyStats.Unsolved,
		}
	}
	
	return result, nil
}