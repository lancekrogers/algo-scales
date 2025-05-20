package stats

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/utils"
)

// FileStorage implements the StatsStorage interface using the file system
type FileStorage struct {
	fs interfaces.FileSystem
}

// NewFileStorage creates a new file storage
func NewFileStorage() *FileStorage {
	return &FileStorage{
		fs: utils.NewFileSystem(),
	}
}

// WithFileSystem sets a custom file system
func (s *FileStorage) WithFileSystem(fs interfaces.FileSystem) *FileStorage {
	s.fs = fs
	return s
}

// SaveSession saves a session's statistics
func (s *FileStorage) SaveSession(session SessionStats) error {
	// Get the stats directory
	statsDir := filepath.Join(s.fs.GetConfigDir(), "stats")
	if err := s.fs.MkdirAll(statsDir, 0755); err != nil {
		return err
	}

	// Generate a unique filename
	filename := fmt.Sprintf("session_%s_%s.json", session.ProblemID, session.StartTime.Format("20060102_150405"))
	statsFile := filepath.Join(statsDir, filename)

	// Save stats to file
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return s.fs.WriteFile(statsFile, data, 0644)
}

// LoadAllSessions loads all session statistics
func (s *FileStorage) LoadAllSessions() ([]SessionStats, error) {
	var sessions []SessionStats

	statsDir := filepath.Join(s.fs.GetConfigDir(), "stats")

	// Check if directory exists
	if !s.fs.Exists(statsDir) {
		return sessions, nil
	}

	// Read all stats files
	files, err := s.fs.ReadDir(statsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !isStatsFile(file.Name()) {
			continue
		}

		// Read file
		data, err := s.fs.ReadFile(filepath.Join(statsDir, file.Name()))
		if err != nil {
			return nil, err
		}

		var session SessionStats
		if err := json.Unmarshal(data, &session); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// ClearAllSessions removes all session statistics
func (s *FileStorage) ClearAllSessions() error {
	statsDir := filepath.Join(s.fs.GetConfigDir(), "stats")

	// Check if directory exists
	if !s.fs.Exists(statsDir) {
		return nil // Nothing to reset
	}

	// Get all files in the directory
	files, err := s.fs.ReadDir(statsDir)
	if err != nil {
		return err
	}

	// Remove each stats file
	for _, file := range files {
		if file.IsDir() || !isStatsFile(file.Name()) {
			continue
		}
		
		if err := s.fs.RemoveAll(filepath.Join(statsDir, file.Name())); err != nil {
			return err
		}
	}

	return nil
}