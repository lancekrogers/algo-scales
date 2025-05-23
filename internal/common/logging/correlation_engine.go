package logging

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ErrorPattern represents a pattern of related errors
type ErrorPattern struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ErrorIDs    []string          `json:"error_ids"`
	Frequency   int               `json:"frequency"`
	FirstSeen   time.Time         `json:"first_seen"`
	LastSeen    time.Time         `json:"last_seen"`
	Categories  []ErrorCategory   `json:"categories"`
	UserActions []string          `json:"user_actions"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ErrorCorrelation represents a relationship between errors
type ErrorCorrelation struct {
	ErrorID1     string    `json:"error_id_1"`
	ErrorID2     string    `json:"error_id_2"`
	Correlation  float64   `json:"correlation"` // 0.0 to 1.0
	TimeWindow   time.Duration `json:"time_window"`
	CommonFields []string  `json:"common_fields"`
	Confidence   float64   `json:"confidence"`
}

// ErrorInsight provides actionable insights from error patterns
type ErrorInsight struct {
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    ErrorSeverity `json:"severity"`
	PatternID   string    `json:"pattern_id"`
	Suggestions []string  `json:"suggestions"`
	AffectedComponents []string `json:"affected_components"`
	Timestamp   time.Time `json:"timestamp"`
}

// CorrelationEngine analyzes error patterns and relationships
type CorrelationEngine struct {
	errors      map[string]*ErrorContext
	patterns    map[string]*ErrorPattern
	correlations []ErrorCorrelation
	insights    []ErrorInsight
}

// NewCorrelationEngine creates a new error correlation engine
func NewCorrelationEngine() *CorrelationEngine {
	return &CorrelationEngine{
		errors:      make(map[string]*ErrorContext),
		patterns:    make(map[string]*ErrorPattern),
		correlations: make([]ErrorCorrelation, 0),
		insights:    make([]ErrorInsight, 0),
	}
}

// AddError adds an error to the correlation engine for analysis
func (ce *CorrelationEngine) AddError(errorCtx *ErrorContext) {
	ce.errors[errorCtx.ID] = errorCtx
	
	// Analyze patterns and correlations
	ce.analyzePatterns(errorCtx)
	ce.findCorrelations(errorCtx)
	ce.generateInsights(errorCtx)
}

// analyzePatterns identifies patterns in the new error
func (ce *CorrelationEngine) analyzePatterns(newError *ErrorContext) {
	// Check for similar errors
	for _, existingError := range ce.errors {
		if ce.areErrorsSimilar(newError, existingError) {
			ce.updatePattern(newError, existingError)
			return
		}
	}
	
	// Create new pattern if no similar errors found
	ce.createNewPattern(newError)
}

// areErrorsSimilar checks if two errors are similar enough to be part of the same pattern
func (ce *CorrelationEngine) areErrorsSimilar(err1, err2 *ErrorContext) bool {
	// Same category
	if err1.Category != err2.Category {
		return false
	}
	
	// Similar user actions
	if ce.calculateStringSimilarity(err1.UserAction, err2.UserAction) < 0.7 {
		return false
	}
	
	// Similar error messages (fuzzy matching)
	if ce.calculateStringSimilarity(err1.ErrorMessage, err2.ErrorMessage) < 0.6 {
		return false
	}
	
	// Similar session context (if available)
	if err1.SessionState != nil && err2.SessionState != nil {
		if err1.SessionState.Language != err2.SessionState.Language {
			return false
		}
		if err1.SessionState.ProblemID != err2.SessionState.ProblemID {
			return false
		}
	}
	
	return true
}

// calculateStringSimilarity calculates similarity between two strings (simplified)
func (ce *CorrelationEngine) calculateStringSimilarity(s1, s2 string) float64 {
	// Simple similarity based on common words
	words1 := strings.Fields(strings.ToLower(s1))
	words2 := strings.Fields(strings.ToLower(s2))
	
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}
	
	commonWords := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 {
				commonWords++
				break
			}
		}
	}
	
	return float64(commonWords) / float64(max(len(words1), len(words2)))
}

// updatePattern updates an existing pattern with a new error
func (ce *CorrelationEngine) updatePattern(newError *ErrorContext, similarError *ErrorContext) {
	// Find the pattern that contains the similar error
	for _, pattern := range ce.patterns {
		for _, errorID := range pattern.ErrorIDs {
			if errorID == similarError.ID {
				// Add new error to this pattern
				pattern.ErrorIDs = append(pattern.ErrorIDs, newError.ID)
				pattern.Frequency++
				pattern.LastSeen = newError.Timestamp
				
				// Update pattern metadata
				ce.updatePatternMetadata(pattern, newError)
				return
			}
		}
	}
}

// createNewPattern creates a new error pattern
func (ce *CorrelationEngine) createNewPattern(errorCtx *ErrorContext) {
	patternID := fmt.Sprintf("pattern_%s_%d", errorCtx.Category, time.Now().Unix())
	
	pattern := &ErrorPattern{
		ID:          patternID,
		Name:        ce.generatePatternName(errorCtx),
		Description: ce.generatePatternDescription(errorCtx),
		ErrorIDs:    []string{errorCtx.ID},
		Frequency:   1,
		FirstSeen:   errorCtx.Timestamp,
		LastSeen:    errorCtx.Timestamp,
		Categories:  []ErrorCategory{errorCtx.Category},
		UserActions: []string{errorCtx.UserAction},
		Metadata:    make(map[string]interface{}),
	}
	
	ce.patterns[patternID] = pattern
}

// generatePatternName generates a human-readable name for a pattern
func (ce *CorrelationEngine) generatePatternName(errorCtx *ErrorContext) string {
	switch errorCtx.Category {
	case TestExecution:
		return fmt.Sprintf("Test Execution Failures (%s)", extractLanguageFromContext(errorCtx))
	case EditorIntegration:
		return "Editor Integration Issues"
	case FileOperations:
		return "File Operation Failures"
	case UIInteraction:
		return "UI Interaction Problems"
	default:
		return fmt.Sprintf("%s Errors", errorCtx.Category)
	}
}

// generatePatternDescription generates a description for a pattern
func (ce *CorrelationEngine) generatePatternDescription(errorCtx *ErrorContext) string {
	switch errorCtx.Category {
	case TestExecution:
		return "Errors occurring during test execution, compilation, or runtime"
	case EditorIntegration:
		return "Issues with opening external editors or editor configuration"
	case FileOperations:
		return "Problems with file creation, reading, writing, or permissions"
	case UIInteraction:
		return "Errors in terminal UI interactions or display issues"
	default:
		return fmt.Sprintf("Errors in the %s category", errorCtx.Category)
	}
}

// findCorrelations identifies correlations between the new error and existing errors
func (ce *CorrelationEngine) findCorrelations(newError *ErrorContext) {
	for _, existingError := range ce.errors {
		if existingError.ID == newError.ID {
			continue
		}
		
		correlation := ce.calculateCorrelation(newError, existingError)
		if correlation.Correlation > 0.5 { // Threshold for significant correlation
			ce.correlations = append(ce.correlations, correlation)
		}
	}
}

// calculateCorrelation calculates correlation between two errors
func (ce *CorrelationEngine) calculateCorrelation(err1, err2 *ErrorContext) ErrorCorrelation {
	correlation := ErrorCorrelation{
		ErrorID1:     err1.ID,
		ErrorID2:     err2.ID,
		CommonFields: make([]string, 0),
	}
	
	score := 0.0
	factors := 0
	
	// Time proximity (errors close in time are more likely related)
	timeDiff := err1.Timestamp.Sub(err2.Timestamp)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	correlation.TimeWindow = timeDiff
	
	if timeDiff < 5*time.Minute {
		score += 0.3
		correlation.CommonFields = append(correlation.CommonFields, "time_proximity")
	}
	factors++
	
	// Same trace ID
	if err1.TraceID != "" && err1.TraceID == err2.TraceID {
		score += 0.4
		correlation.CommonFields = append(correlation.CommonFields, "trace_id")
	}
	factors++
	
	// Same user action sequence
	if err1.UserAction == err2.UserAction {
		score += 0.2
		correlation.CommonFields = append(correlation.CommonFields, "user_action")
	}
	factors++
	
	// Same session context
	if err1.SessionState != nil && err2.SessionState != nil {
		if err1.SessionState.ProblemID == err2.SessionState.ProblemID {
			score += 0.2
			correlation.CommonFields = append(correlation.CommonFields, "problem_id")
		}
		if err1.SessionState.Language == err2.SessionState.Language {
			score += 0.1
			correlation.CommonFields = append(correlation.CommonFields, "language")
		}
	}
	factors += 2
	
	correlation.Correlation = score / float64(factors)
	correlation.Confidence = ce.calculateConfidence(correlation)
	
	return correlation
}

// calculateConfidence calculates confidence in a correlation
func (ce *CorrelationEngine) calculateConfidence(correlation ErrorCorrelation) float64 {
	// Confidence based on number of common fields and correlation strength
	fieldWeight := float64(len(correlation.CommonFields)) / 5.0 // Max 5 fields
	return (correlation.Correlation + fieldWeight) / 2.0
}

// generateInsights generates actionable insights from error patterns
func (ce *CorrelationEngine) generateInsights(newError *ErrorContext) {
	// Check for high-frequency patterns
	for _, pattern := range ce.patterns {
		if pattern.Frequency >= 3 && time.Since(pattern.LastSeen) < time.Hour {
			insight := ce.createFrequencyInsight(pattern)
			ce.insights = append(ce.insights, insight)
		}
	}
	
	// Check for error cascades (highly correlated errors in short time)
	recentCorrelations := ce.getRecentHighCorrelations(5 * time.Minute)
	if len(recentCorrelations) >= 2 {
		insight := ce.createCascadeInsight(recentCorrelations)
		ce.insights = append(ce.insights, insight)
	}
	
	// Generate category-specific insights
	ce.generateCategoryInsights(newError)
}

// createFrequencyInsight creates insight for high-frequency error patterns
func (ce *CorrelationEngine) createFrequencyInsight(pattern *ErrorPattern) ErrorInsight {
	return ErrorInsight{
		Type:        "high_frequency",
		Title:       fmt.Sprintf("Frequent %s", pattern.Name),
		Description: fmt.Sprintf("Pattern '%s' has occurred %d times in the last hour", pattern.Name, pattern.Frequency),
		Severity:    SeverityHigh,
		PatternID:   pattern.ID,
		Suggestions: ce.generateSuggestions(pattern),
		AffectedComponents: ce.extractAffectedComponents(pattern),
		Timestamp:   time.Now(),
	}
}

// createCascadeInsight creates insight for error cascades
func (ce *CorrelationEngine) createCascadeInsight(correlations []ErrorCorrelation) ErrorInsight {
	return ErrorInsight{
		Type:        "error_cascade",
		Title:       "Error Cascade Detected",
		Description: fmt.Sprintf("Multiple related errors occurred within a short time window"),
		Severity:    SeverityCritical,
		PatternID:   fmt.Sprintf("cascade_%d", time.Now().Unix()),
		Suggestions: []string{
			"Check for root cause in the first error of the sequence",
			"Review system resources and dependencies",
			"Consider implementing retry logic with backoff",
		},
		AffectedComponents: []string{"multiple"},
		Timestamp:          time.Now(),
	}
}

// generateCategoryInsights generates category-specific insights
func (ce *CorrelationEngine) generateCategoryInsights(errorCtx *ErrorContext) {
	switch errorCtx.Category {
	case TestExecution:
		ce.generateTestExecutionInsights(errorCtx)
	case EditorIntegration:
		ce.generateEditorInsights(errorCtx)
	case FileOperations:
		ce.generateFileOpInsights(errorCtx)
	}
}

// generateTestExecutionInsights generates insights for test execution errors
func (ce *CorrelationEngine) generateTestExecutionInsights(errorCtx *ErrorContext) {
	if strings.Contains(errorCtx.ErrorMessage, "timeout") {
		insight := ErrorInsight{
			Type:        "test_timeout",
			Title:       "Test Execution Timeout",
			Description: "Tests are timing out, possibly due to infinite loops or performance issues",
			Severity:    SeverityHigh,
			Suggestions: []string{
				"Review code for infinite loops",
				"Check algorithm complexity",
				"Increase timeout if appropriate",
				"Add debugging output to identify where code hangs",
			},
			AffectedComponents: []string{"test_execution"},
			Timestamp:          time.Now(),
		}
		ce.insights = append(ce.insights, insight)
	}
}

// generateEditorInsights generates insights for editor errors
func (ce *CorrelationEngine) generateEditorInsights(errorCtx *ErrorContext) {
	insight := ErrorInsight{
		Type:        "editor_failure",
		Title:       "Editor Integration Issue",
		Description: "Problems with external editor integration",
		Severity:    SeverityMedium,
		Suggestions: []string{
			"Check EDITOR environment variable",
			"Verify editor is installed and in PATH",
			"Try setting a different editor in settings",
			"Check file permissions in workspace directory",
		},
		AffectedComponents: []string{"editor_integration"},
		Timestamp:          time.Now(),
	}
	ce.insights = append(ce.insights, insight)
}

// generateFileOpInsights generates insights for file operation errors
func (ce *CorrelationEngine) generateFileOpInsights(errorCtx *ErrorContext) {
	if strings.Contains(errorCtx.ErrorMessage, "permission") {
		insight := ErrorInsight{
			Type:        "permission_error",
			Title:       "File Permission Issue",
			Description: "File operation failed due to insufficient permissions",
			Severity:    SeverityMedium,
			Suggestions: []string{
				"Check file and directory permissions",
				"Ensure write access to workspace directory",
				"Run with appropriate privileges if necessary",
				"Check disk space availability",
			},
			AffectedComponents: []string{"file_operations"},
			Timestamp:          time.Now(),
		}
		ce.insights = append(ce.insights, insight)
	}
}

// GetPatterns returns all detected error patterns
func (ce *CorrelationEngine) GetPatterns() []*ErrorPattern {
	patterns := make([]*ErrorPattern, 0, len(ce.patterns))
	for _, pattern := range ce.patterns {
		patterns = append(patterns, pattern)
	}
	
	// Sort by frequency (most frequent first)
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Frequency > patterns[j].Frequency
	})
	
	return patterns
}

// GetInsights returns actionable insights
func (ce *CorrelationEngine) GetInsights() []ErrorInsight {
	// Return most recent insights first
	sort.Slice(ce.insights, func(i, j int) bool {
		return ce.insights[i].Timestamp.After(ce.insights[j].Timestamp)
	})
	
	return ce.insights
}

// GetCorrelations returns error correlations
func (ce *CorrelationEngine) GetCorrelations() []ErrorCorrelation {
	// Sort by correlation strength
	sort.Slice(ce.correlations, func(i, j int) bool {
		return ce.correlations[i].Correlation > ce.correlations[j].Correlation
	})
	
	return ce.correlations
}

// Helper functions
func extractLanguageFromContext(errorCtx *ErrorContext) string {
	if errorCtx.SessionState != nil {
		return errorCtx.SessionState.Language
	}
	return "unknown"
}

func (ce *CorrelationEngine) generateSuggestions(pattern *ErrorPattern) []string {
	suggestions := make([]string, 0)
	
	for _, category := range pattern.Categories {
		switch category {
		case TestExecution:
			suggestions = append(suggestions, "Review test code and inputs")
			suggestions = append(suggestions, "Check language environment setup")
		case EditorIntegration:
			suggestions = append(suggestions, "Verify editor configuration")
			suggestions = append(suggestions, "Check file permissions")
		case FileOperations:
			suggestions = append(suggestions, "Check disk space and permissions")
			suggestions = append(suggestions, "Verify file paths exist")
		}
	}
	
	return suggestions
}

func (ce *CorrelationEngine) extractAffectedComponents(pattern *ErrorPattern) []string {
	components := make([]string, 0)
	for _, category := range pattern.Categories {
		components = append(components, string(category))
	}
	return components
}

func (ce *CorrelationEngine) getRecentHighCorrelations(timeWindow time.Duration) []ErrorCorrelation {
	recent := make([]ErrorCorrelation, 0)
	cutoff := time.Now().Add(-timeWindow)
	
	for _, correlation := range ce.correlations {
		if correlation.Correlation > 0.7 && correlation.TimeWindow < timeWindow {
			// Check if either error is recent
			if err1, exists := ce.errors[correlation.ErrorID1]; exists && err1.Timestamp.After(cutoff) {
				recent = append(recent, correlation)
			}
		}
	}
	
	return recent
}

func (ce *CorrelationEngine) updatePatternMetadata(pattern *ErrorPattern, newError *ErrorContext) {
	// Update user actions
	actionExists := false
	for _, action := range pattern.UserActions {
		if action == newError.UserAction {
			actionExists = true
			break
		}
	}
	if !actionExists {
		pattern.UserActions = append(pattern.UserActions, newError.UserAction)
	}
	
	// Update categories
	categoryExists := false
	for _, category := range pattern.Categories {
		if category == newError.Category {
			categoryExists = true
			break
		}
	}
	if !categoryExists {
		pattern.Categories = append(pattern.Categories, newError.Category)
	}
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}