package services

import (
	"context"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

// StatsCommandService provides business logic for stats command operations
type StatsCommandService interface {
	// GetOverallStats returns overall performance statistics
	GetOverallStats(ctx context.Context) (*interfaces.OverallStats, error)
	
	// GetPatternStats returns performance by algorithm pattern
	GetPatternStats(ctx context.Context) (map[string]*interfaces.PatternStats, error)
	
	// GetDifficultyStats returns performance by difficulty level
	GetDifficultyStats(ctx context.Context) (map[string]*interfaces.DifficultyStats, error)
	
	// GetLanguageStats returns performance by programming language
	GetLanguageStats(ctx context.Context) (map[string]*interfaces.LanguageStats, error)
	
	// GetRecentActivity returns recent session activity
	GetRecentActivity(ctx context.Context, days int) ([]*interfaces.DailyStats, error)
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
func (s *StatsCommandServiceImpl) GetOverallStats(ctx context.Context) (*interfaces.OverallStats, error) {
	summary, err := s.statsService.GetSummary(ctx)
	if err != nil {
		return nil, err
	}
	
	trends, err := s.statsService.GetTrends(ctx)
	if err != nil {
		return nil, err
	}
	
	return &interfaces.OverallStats{
		Summary: summary,
		Trends:  trends,
	}, nil
}

// GetPatternStats returns performance by algorithm pattern
func (s *StatsCommandServiceImpl) GetPatternStats(ctx context.Context) (map[string]*interfaces.PatternStats, error) {
	patternStats, err := s.statsService.GetByPattern(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert to pointer map
	result := make(map[string]*interfaces.PatternStats)
	for k, v := range patternStats {
		pattern := v // Create a copy
		result[k] = &pattern
	}
	
	return result, nil
}

// GetDifficultyStats returns performance by difficulty level
func (s *StatsCommandServiceImpl) GetDifficultyStats(ctx context.Context) (map[string]*interfaces.DifficultyStats, error) {
	// This would need to be implemented by analyzing sessions
	// For now, return empty map
	return make(map[string]*interfaces.DifficultyStats), nil
}

// GetLanguageStats returns performance by programming language
func (s *StatsCommandServiceImpl) GetLanguageStats(ctx context.Context) (map[string]*interfaces.LanguageStats, error) {
	// This would need to be implemented by analyzing sessions
	// For now, return empty map
	return make(map[string]*interfaces.LanguageStats), nil
}

// GetRecentActivity returns recent session activity
func (s *StatsCommandServiceImpl) GetRecentActivity(ctx context.Context, days int) ([]*interfaces.DailyStats, error) {
	// This would need to be implemented by analyzing recent sessions
	// For now, return empty slice
	return make([]*interfaces.DailyStats, 0), nil
}

// LegacyStatsCommandService provides backward compatibility with legacy stats functions
type LegacyStatsCommandService struct{}

// GetOverallStats returns overall performance statistics using legacy functions
func (s *LegacyStatsCommandService) GetOverallStats(ctx context.Context) (*interfaces.OverallStats, error) {
	// Legacy implementation with basic fallback
	summary := &interfaces.Summary{
		TotalAttempted: 0,
		TotalSolved:    0,
		SuccessRate:    0.0,
		AvgSolveTime:   "0s",
	}
	
	// Create empty trends for now
	trends := &interfaces.Trends{
		Daily:  []interfaces.DailyTrend{},
		Weekly: []interfaces.WeeklyTrend{},
	}
	
	return &interfaces.OverallStats{
		Summary: summary,
		Trends:  trends,
	}, nil
}

// GetPatternStats returns performance by algorithm pattern using legacy functions
func (s *LegacyStatsCommandService) GetPatternStats(ctx context.Context) (map[string]*interfaces.PatternStats, error) {
	// Legacy implementation with basic fallback
	return make(map[string]*interfaces.PatternStats), nil
}

// GetDifficultyStats returns performance by difficulty level using legacy functions
func (s *LegacyStatsCommandService) GetDifficultyStats(ctx context.Context) (map[string]*interfaces.DifficultyStats, error) {
	// Legacy implementation with basic fallback
	return make(map[string]*interfaces.DifficultyStats), nil
}

// GetLanguageStats returns performance by programming language using legacy functions
func (s *LegacyStatsCommandService) GetLanguageStats(ctx context.Context) (map[string]*interfaces.LanguageStats, error) {
	// Legacy implementation with basic fallback
	return make(map[string]*interfaces.LanguageStats), nil
}

// GetRecentActivity returns recent session activity using legacy functions
func (s *LegacyStatsCommandService) GetRecentActivity(ctx context.Context, days int) ([]*interfaces.DailyStats, error) {
	// Legacy implementation with basic fallback
	return make([]*interfaces.DailyStats, 0), nil
}