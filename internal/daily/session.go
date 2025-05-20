package daily

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

const (
	// SessionBucketName is the BoltDB bucket for storing daily session info
	SessionBucketName = "daily_sessions"
	
	// ActiveSessionKey is the key for storing the current session
	ActiveSessionKey = "active_session"
	
	// SessionDBFileName is the name of the session database file
	SessionDBFileName = "daily_sessions.db"
)

// DailySession represents a daily practice session
type DailySession struct {
	Date      string                  `json:"date"`
	Problems  map[string]DailyProblem `json:"problems"`
	StartTime time.Time               `json:"start_time"`
	EndTime   time.Time               `json:"end_time,omitempty"`
	Completed bool                    `json:"completed"`
}

// CreateNewSession creates a new daily session
func CreateNewSession() (*DailySession, error) {
	today := time.Now().Format("2006-01-02")
	
	// Initialize with all patterns as pending
	problems := make(map[string]DailyProblem)
	
	for _, scale := range Scales {
		problems[scale.Pattern] = DailyProblem{
			Pattern:    scale.Pattern,
			ProblemID:  "", // Will be populated when we select a problem
			State:      StatePending,
			StartedAt:  time.Time{},
			Attempts:   0,
		}
	}
	
	session := &DailySession{
		Date:      today,
		Problems:  problems,
		StartTime: time.Now(),
		Completed: false,
	}
	
	// Save the session
	if err := SaveSession(session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}
	
	return session, nil
}

// LoadSession loads the active daily session
func LoadSession() (*DailySession, error) {
	dbPath := GetSessionDBPath()
	
	// Create dirs if needed
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("error creating directories: %w", err)
	}
	
	// Open database file (will be created if it doesn't exist)
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()
	
	// Initialize the bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(SessionBucketName))
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}
	
	// Load session data
	var session DailySession
	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(SessionBucketName))
		data := bucket.Get([]byte(ActiveSessionKey))
		
		if data == nil {
			// No active session
			return fmt.Errorf("no active session found")
		}
		
		// Unmarshal the JSON data
		if err := json.Unmarshal(data, &session); err != nil {
			return fmt.Errorf("error unmarshaling session data: %w", err)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return &session, nil
}

// SaveSession saves the daily session to the database
func SaveSession(session *DailySession) error {
	dbPath := GetSessionDBPath()
	
	// Open database file
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()
	
	// Marshal the session struct to JSON
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("error marshaling session data: %w", err)
	}
	
	// Save to database
	err = db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(SessionBucketName))
		err := bucket.Put([]byte(ActiveSessionKey), data)
		if err != nil {
			return fmt.Errorf("error saving session data: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}
	
	return nil
}

// GetOrCreateSession loads the active session or creates a new one if needed
func GetOrCreateSession() (*DailySession, error) {
	// Try to load existing session
	session, err := LoadSession()
	if err == nil {
		// Check if this session is for today
		today := time.Now().Format("2006-01-02")
		if session.Date == today {
			return session, nil
		}
		
		// Session exists but it's from a previous day
		// Save it as completed if it wasn't already
		if !session.Completed {
			session.Completed = true
			session.EndTime = time.Now()
			if err := SaveSession(session); err != nil {
				return nil, fmt.Errorf("error completing previous session: %w", err)
			}
		}
	}
	
	// Create a new session for today
	return CreateNewSession()
}

// StartProblem marks a problem as in progress and assigns a specific problem
func (s *DailySession) StartProblem(pattern string, problemID string) error {
	// Check if pattern exists
	prob, ok := s.Problems[pattern]
	if !ok {
		return fmt.Errorf("pattern not found: %s", pattern)
	}
	
	// Update problem information
	prob.ProblemID = problemID
	prob.State = StateInProgress
	prob.StartedAt = time.Now()
	prob.Attempts++
	
	// Save back to map
	s.Problems[pattern] = prob
	
	// Save the session
	return SaveSession(s)
}

// CompleteProblem marks a problem as completed
func (s *DailySession) CompleteProblem(pattern string) error {
	// Check if pattern exists
	prob, ok := s.Problems[pattern]
	if !ok {
		return fmt.Errorf("pattern not found: %s", pattern)
	}
	
	// Update problem information
	prob.State = StateCompleted
	prob.CompletedAt = time.Now()
	
	// Save back to map
	s.Problems[pattern] = prob
	
	// Check if all problems are completed
	allCompleted := true
	for _, p := range s.Problems {
		if p.State != StateCompleted && p.State != StateSkipped {
			allCompleted = false
			break
		}
	}
	
	// Mark session as completed if all problems are done
	if allCompleted {
		s.Completed = true
		s.EndTime = time.Now()
	}
	
	// Save the session
	return SaveSession(s)
}

// SkipProblem marks a problem as skipped
func (s *DailySession) SkipProblem(pattern string) error {
	// Check if pattern exists
	prob, ok := s.Problems[pattern]
	if !ok {
		return fmt.Errorf("pattern not found: %s", pattern)
	}
	
	// Update problem information
	prob.State = StateSkipped
	
	// Save back to map
	s.Problems[pattern] = prob
	
	// Save the session
	return SaveSession(s)
}

// GetNextPendingPattern returns the next pattern that is pending
func (s *DailySession) GetNextPendingPattern() string {
	// Check if any pattern is in progress
	for pattern, prob := range s.Problems {
		if prob.State == StateInProgress {
			return pattern
		}
	}
	
	// If not, find the first pending pattern
	for _, scale := range Scales {
		pattern := scale.Pattern
		if prob, ok := s.Problems[pattern]; ok && prob.State == StatePending {
			return pattern
		}
	}
	
	return ""
}

// GetCompletedCount returns the number of completed problems
func (s *DailySession) GetCompletedCount() int {
	count := 0
	for _, prob := range s.Problems {
		if prob.State == StateCompleted {
			count++
		}
	}
	return count
}

// GetSkippedCount returns the number of skipped problems
func (s *DailySession) GetSkippedCount() int {
	count := 0
	for _, prob := range s.Problems {
		if prob.State == StateSkipped {
			count++
		}
	}
	return count
}

// GetPendingCount returns the number of pending problems
func (s *DailySession) GetPendingCount() int {
	count := 0
	for _, prob := range s.Problems {
		if prob.State == StatePending {
			count++
		}
	}
	return count
}

// GetInProgressCount returns the number of problems in progress
func (s *DailySession) GetInProgressCount() int {
	count := 0
	for _, prob := range s.Problems {
		if prob.State == StateInProgress {
			count++
		}
	}
	return count
}

// GetTotalProblems returns the total number of problems
func (s *DailySession) GetTotalProblems() int {
	return len(s.Problems)
}

// GetSessionDBPath returns the path to the session database
func GetSessionDBPath() string {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if can't get home
		return SessionDBFileName
	}
	
	// Create path for database file in .algo-scales directory
	return filepath.Join(homeDir, ".algo-scales", "stats", SessionDBFileName)
}