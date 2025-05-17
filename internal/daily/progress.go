// Package daily provides functionality for handling daily scale practice
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
	// BucketName is the BoltDB bucket for storing daily progress
	BucketName = "daily_progress"
	
	// ProgressKey is the key for storing the ScaleProgress struct
	ProgressKey = "progress"
	
	// DBFileName is the name of the BoltDB database file
	DBFileName = "daily.db"
)

// ScaleProgress tracks progress through scales
type ScaleProgress struct {
	Current       int       `json:"current"`
	LastPracticed time.Time `json:"last_practiced"`
	Completed     []string  `json:"completed"`
	Streak        int       `json:"streak"`
	LongestStreak int       `json:"longest_streak"`
}

// LoadProgress loads the scale progress from BoltDB
func LoadProgress() (ScaleProgress, error) {
	dbPath := GetDBPath()
	
	// Create dirs if needed
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return ScaleProgress{}, fmt.Errorf("error creating directories: %w", err)
	}
	
	// Default progress (starting fresh)
	defaultProgress := ScaleProgress{
		Current:       0,
		LastPracticed: time.Time{}, // Zero time (never practiced)
		Completed:     []string{},
		Streak:        0,
		LongestStreak: 0,
	}
	
	// Open database file (will be created if it doesn't exist)
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return defaultProgress, fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()
	
	// Initialize the bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return defaultProgress, fmt.Errorf("error initializing database: %w", err)
	}
	
	// Load progress data
	var progress ScaleProgress
	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		data := bucket.Get([]byte(ProgressKey))
		
		if data == nil {
			// No data yet, return default
			progress = defaultProgress
			return nil
		}
		
		// Unmarshal the JSON data
		if err := json.Unmarshal(data, &progress); err != nil {
			return fmt.Errorf("error unmarshaling progress data: %w", err)
		}
		
		return nil
	})
	
	if err != nil {
		return defaultProgress, fmt.Errorf("error loading progress: %w", err)
	}
	
	return progress, nil
}

// SaveProgress saves the scale progress to BoltDB
func SaveProgress(progress ScaleProgress) error {
	dbPath := GetDBPath()
	
	// Open database file
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()
	
	// Marshal the progress struct to JSON
	data, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("error marshaling progress data: %w", err)
	}
	
	// Save to database
	err = db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		err := bucket.Put([]byte(ProgressKey), data)
		if err != nil {
			return fmt.Errorf("error saving progress data: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("error saving progress: %w", err)
	}
	
	return nil
}

// UpdateStreak updates the practice streak based on the last practice time
func UpdateStreak(progress *ScaleProgress) {
	// If this is the first practice session ever
	if progress.LastPracticed.IsZero() {
		progress.Streak = 1
		progress.LongestStreak = 1
		return
	}
	
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	lastPracticed := progress.LastPracticed.Truncate(24 * time.Hour)
	
	// If practiced today, don't update streak
	if lastPracticed.Equal(today) {
		return
	}
	
	// If practiced yesterday, increment streak
	if lastPracticed.Equal(yesterday) {
		progress.Streak++
		if progress.Streak > progress.LongestStreak {
			progress.LongestStreak = progress.Streak
		}
	} else {
		// Break in streak, reset to 1
		progress.Streak = 1
	}
}

// Make getDBPath a variable for testing
var GetDBPath = func() string {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if can't get home
		return DBFileName
	}
	
	// Create path for database file in .algo-scales directory
	return filepath.Join(homeDir, ".algo-scales", "stats", DBFileName)
}

// Contains checks if a string is in a slice
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}