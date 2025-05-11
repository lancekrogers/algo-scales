package session

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/utils"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Manager manages practice sessions
type Manager struct {
	// Map of active sessions by ID
	sessions     map[string]interfaces.Session
	sessionMutex sync.RWMutex
}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]interfaces.Session),
	}
}

// StartSession begins a new practice session
func (m *Manager) StartSession(opts interfaces.SessionOptions) (interfaces.Session, error) {
	// Choose problem based on options
	var p *problem.Problem
	var err error
	
	if opts.ProblemID != "" {
		// Specific problem requested
		p, err = problem.GetByID(opts.ProblemID)
		if err != nil {
			return nil, fmt.Errorf("failed to load problem: %v", err)
		}
	} else if opts.Mode == interfaces.CramMode {
		// Cram mode - choose problems from common patterns
		p, err = selectCramProblem()
		if err != nil {
			return nil, fmt.Errorf("failed to select problem for cram mode: %v", err)
		}
	} else {
		// Filter by pattern/difficulty if specified
		p, err = selectProblem(opts.Pattern, opts.Difficulty)
		if err != nil {
			return nil, fmt.Errorf("failed to select problem: %v", err)
		}
	}
	
	// Initialize session
	session := &SessionImpl{
		Options:      opts,
		Problem:      p,
		StartTime:    time.Now(),
		hintsShown:   opts.Mode == interfaces.LearnMode,
		ShowPattern:  opts.Mode == interfaces.LearnMode,
		solutionShown: false,
	}
	
	// Create workspace
	if err := createWorkspace(session); err != nil {
		return nil, fmt.Errorf("failed to create workspace: %v", err)
	}
	
	// Generate a session ID
	sessionID := fmt.Sprintf("%s-%d", p.ID, time.Now().Unix())
	
	// Store session
	m.sessionMutex.Lock()
	m.sessions[sessionID] = session
	m.sessionMutex.Unlock()
	
	return session, nil
}

// GetSessionByID retrieves an active session
func (m *Manager) GetSessionByID(id string) (interfaces.Session, bool) {
	m.sessionMutex.RLock()
	defer m.sessionMutex.RUnlock()
	
	session, exists := m.sessions[id]
	return session, exists
}

// FinishSession completes a session
func (m *Manager) FinishSession(sessionID string, solved bool) error {
	m.sessionMutex.Lock()
	defer m.sessionMutex.Unlock()
	
	session, exists := m.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	
	// Finish the session
	err := session.Finish(solved)
	if err != nil {
		return err
	}
	
	// Remove from active sessions
	delete(m.sessions, sessionID)
	
	return nil
}

// createWorkspace sets up a workspace for the problem
func createWorkspace(s *SessionImpl) error {
	// Create workspace directory
	workspaceDir := filepath.Join(utils.TempDir(), "algo-scales", s.Problem.ID)
	if err := utils.CreateDirectory(workspaceDir); err != nil {
		return err
	}

	s.Workspace = workspaceDir

	// Create problem description file
	descriptionFile := filepath.Join(workspaceDir, "problem.md")
	description := s.FormatDescription()
	if err := utils.WriteFile(descriptionFile, []byte(description), 0644); err != nil {
		return err
	}

	// Create code file with starter code
	ext := languageExtension(s.Options.Language)
	codeFile := filepath.Join(workspaceDir, fmt.Sprintf("solution.%s", ext))

	starterCode, ok := s.Problem.StarterCode[s.Options.Language]
	if !ok {
		// Fallback to a default language if the requested one isn't available
		for lang, code := range s.Problem.StarterCode {
			starterCode = code
			s.Options.Language = lang
			break
		}
	}

	if err := utils.WriteFile(codeFile, []byte(starterCode), 0644); err != nil {
		return err
	}

	s.CodeFile = codeFile
	s.Code = starterCode

	return nil
}

// selectProblem chooses a problem based on filters
var selectProblem = func(pattern, difficulty string) (*problem.Problem, error) {
	// Get all problems
	problems, err := problem.ListAll()
	if err != nil {
		return nil, err
	}

	// Filter problems
	var filtered []problem.Problem
	for _, p := range problems {
		matchesPattern := pattern == "" || containsPattern(p.Patterns, pattern)
		matchesDifficulty := difficulty == "" || p.Difficulty == difficulty

		if matchesPattern && matchesDifficulty {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no problems match the specified filters")
	}

	// Choose a random problem from filtered list
	rand.Seed(time.Now().UnixNano())
	selected := filtered[rand.Intn(len(filtered))]

	return &selected, nil
}

// selectCramProblem chooses a problem for cram mode
var selectCramProblem = func() (*problem.Problem, error) {
	// For cram mode, we focus on the most common patterns
	commonPatterns := []string{
		"sliding-window",
		"two-pointers",
		"fast-slow-pointers",
		"hash-map",
		"binary-search",
		"dfs",
		"bfs",
		"dynamic-programming",
		"greedy",
		"union-find",
		"heap",
	}

	// Choose a random pattern
	rand.Seed(time.Now().UnixNano())
	pattern := commonPatterns[rand.Intn(len(commonPatterns))]

	// Get a problem with this pattern
	return selectProblem(pattern, "")
}

// Helper functions

// containsPattern checks if a pattern is in a list
func containsPattern(patterns []string, pattern string) bool {
	for _, p := range patterns {
		if p == pattern {
			return true
		}
	}
	return false
}

// JoinStrings joins a string slice with commas
func JoinStrings(strings []string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += ", " + strings[i]
	}

	return result
}

// languageExtension returns the file extension for a language
func languageExtension(language string) string {
	switch language {
	case "go":
		return "go"
	case "python":
		return "py"
	case "javascript":
		return "js"
	default:
		return "txt"
	}
}