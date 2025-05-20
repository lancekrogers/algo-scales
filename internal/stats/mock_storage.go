package stats

import (
	"fmt"
	"sync"
)

// MockStorage implements the StatsStorage interface for testing
type MockStorage struct {
	sessions map[string]SessionStats
	mutex    sync.RWMutex
}

// NewMockStorage creates a new mock storage for statistics
func NewMockStorage() *MockStorage {
	return &MockStorage{
		sessions: make(map[string]SessionStats),
	}
}

// SaveSession saves a session's statistics
func (s *MockStorage) SaveSession(session SessionStats) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Generate a key based on problem ID and start time
	key := fmt.Sprintf("%s_%s", session.ProblemID, session.StartTime.Format("20060102_150405"))
	s.sessions[key] = session
	return nil
}

// LoadAllSessions loads all session statistics
func (s *MockStorage) LoadAllSessions() ([]SessionStats, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	sessions := make([]SessionStats, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessions = append(sessions, session)
	}
	
	return sessions, nil
}

// ClearAllSessions removes all session statistics
func (s *MockStorage) ClearAllSessions() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.sessions = make(map[string]SessionStats)
	return nil
}

// AddSession adds a session directly to the mock storage
// This is a helper method for testing, not part of the StatsStorage interface
func (s *MockStorage) AddSession(session SessionStats) *MockStorage {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	key := fmt.Sprintf("%s_%s", session.ProblemID, session.StartTime.Format("20060102_150405"))
	s.sessions[key] = session
	
	return s
}