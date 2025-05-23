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
	"github.com/lancekrogers/algo-scales/internal/session/execution"
)

// Manager manages practice sessions
type Manager struct {
	// Map of active sessions by ID
	sessions     map[string]interfaces.Session
	sessionMutex sync.RWMutex
	fs           interfaces.FileSystem
	problemRepo  interfaces.ProblemRepository
	testRegistry interfaces.TestRunnerRegistry
}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{
		sessions:     make(map[string]interfaces.Session),
		fs:           utils.NewFileSystem(),
		problemRepo:  problem.NewRepository(),
		testRegistry: execution.DefaultRegistry,
	}
}

// WithFileSystem sets a custom file system for the manager
func (m *Manager) WithFileSystem(fs interfaces.FileSystem) *Manager {
	m.fs = fs
	return m
}

// WithProblemRepository sets a custom problem repository
func (m *Manager) WithProblemRepository(repo interfaces.ProblemRepository) *Manager {
	m.problemRepo = repo
	return m
}

// WithTestRunnerRegistry sets a custom test runner registry
func (m *Manager) WithTestRunnerRegistry(registry interfaces.TestRunnerRegistry) *Manager {
	m.testRegistry = registry
	return m
}

// StartSession begins a new practice session
func (m *Manager) StartSession(opts interfaces.SessionOptions) (interfaces.Session, error) {
	// Choose problem based on options
	var p *problem.Problem
	var err error
	
	if opts.ProblemID != "" {
		// Specific problem requested
		interfaceProb, err := m.problemRepo.GetByID(opts.ProblemID)
		if err != nil {
			return nil, fmt.Errorf("failed to load problem: %v", err)
		}
		
		// Convert to local problem type
		localProb := m.convertInterfaceToLocalProblem(*interfaceProb)
		p = &localProb
	} else if opts.Mode == interfaces.CramMode {
		// Cram mode - choose problems from common patterns
		p, err = m.selectCramProblem()
		if err != nil {
			return nil, fmt.Errorf("failed to select problem for cram mode: %v", err)
		}
	} else {
		// Filter by pattern/difficulty if specified
		p, err = m.selectProblem(opts.Pattern, opts.Difficulty)
		if err != nil {
			return nil, fmt.Errorf("failed to select problem: %v", err)
		}
	}
	
	// Initialize session
	session := NewSessionImpl(opts, p)
	session.hintsShown = opts.Mode == interfaces.LearnMode
	session.ShowPattern = opts.Mode == interfaces.LearnMode
	session.WithFileSystem(m.fs)
	session.WithTestRegistry(m.testRegistry)
	
	// Create workspace
	if err := m.createWorkspace(session); err != nil {
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
func (m *Manager) createWorkspace(s *SessionImpl) error {
	// Create workspace directory
	workspaceDir := filepath.Join(m.fs.TempDir(), "algo-scales", s.Problem.ID)
	if err := m.fs.MkdirAll(workspaceDir, 0755); err != nil {
		return err
	}

	s.Workspace = workspaceDir

	// Create problem description file
	descriptionFile := filepath.Join(workspaceDir, "problem.md")
	description := s.FormatDescription()
	if err := m.fs.WriteFile(descriptionFile, []byte(description), 0644); err != nil {
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

	if err := m.fs.WriteFile(codeFile, []byte(starterCode), 0644); err != nil {
		return err
	}

	s.CodeFile = codeFile
	s.Code = starterCode

	return nil
}

// selectProblem chooses a problem based on pattern and difficulty
func (m *Manager) selectProblem(pattern, difficulty string) (*problem.Problem, error) {
	// Get all problems
	var problems []problem.Problem
	
	if pattern != "" && difficulty != "" {
		// Filter by both pattern and difficulty
		interfaceProbs, err := m.problemRepo.GetByPattern(pattern)
		if err != nil {
			return nil, err
		}
		
		// Convert and filter by difficulty
		for _, p := range interfaceProbs {
			if p.Difficulty == difficulty {
				localProb := m.convertInterfaceToLocalProblem(p)
				problems = append(problems, localProb)
			}
		}
	} else if pattern != "" {
		// Filter by pattern only
		interfaceProbs, err := m.problemRepo.GetByPattern(pattern)
		if err != nil {
			return nil, err
		}
		problems = m.convertInterfaceProblemsToLocal(interfaceProbs)
	} else if difficulty != "" {
		// Filter by difficulty only
		interfaceProbs, err := m.problemRepo.GetByDifficulty(difficulty)
		if err != nil {
			return nil, err
		}
		problems = m.convertInterfaceProblemsToLocal(interfaceProbs)
	} else {
		// No filters, get all problems
		interfaceProbs, err := m.problemRepo.GetAll()
		if err != nil {
			return nil, err
		}
		problems = m.convertInterfaceProblemsToLocal(interfaceProbs)
	}
	
	if len(problems) == 0 {
		return nil, fmt.Errorf("no problems found matching criteria")
	}
	
	// Select random problem
	rand.Seed(time.Now().UnixNano())
	selectedIndex := rand.Intn(len(problems))
	return &problems[selectedIndex], nil
}

// selectCramProblem chooses a problem for cram mode
func (m *Manager) selectCramProblem() (*problem.Problem, error) {
	// For cram mode, we typically want to focus on common patterns
	// This is a simplified implementation - may be improved in the future
	commonPatterns := []string{
		"two-pointers",
		"sliding-window",
		"hash-map",
		"binary-search",
		"dfs",
		"bfs",
		"dynamic-programming",
	}
	
	// Choose a random pattern
	rand.Seed(time.Now().UnixNano())
	patternIndex := rand.Intn(len(commonPatterns))
	selectedPattern := commonPatterns[patternIndex]
	
	// Get problems for this pattern
	interfaceProbs, err := m.problemRepo.GetByPattern(selectedPattern)
	if err != nil {
		return nil, err
	}
	patternProblems := m.convertInterfaceProblemsToLocal(interfaceProbs)
	
	if len(patternProblems) == 0 {
		return nil, fmt.Errorf("no problems found for pattern: %s", selectedPattern)
	}
	
	// Select random problem from this pattern
	selectedIndex := rand.Intn(len(patternProblems))
	return &patternProblems[selectedIndex], nil
}

// JoinStrings joins a slice of strings with commas
func JoinStrings(items []string) string {
	if len(items) == 0 {
		return ""
	}
	
	result := items[0]
	for i := 1; i < len(items); i++ {
		result += ", " + items[i]
	}
	
	return result
}

// languageExtension returns the file extension for a language
func languageExtension(language string) string {
	extensions := map[string]string{
		"go":         "go",
		"python":     "py",
		"javascript": "js",
		"java":       "java",
		"c++":        "cpp",
		"typescript": "ts",
	}
	
	if ext, ok := extensions[language]; ok {
		return ext
	}
	
	// Default to .txt if language not recognized
	return "txt"
}
// convertInterfaceToLocalProblem converts an interfaces.Problem to a local problem.Problem
func (m *Manager) convertInterfaceToLocalProblem(p interfaces.Problem) problem.Problem {
	// Convert test cases
	testCases := make([]problem.TestCase, len(p.TestCases))
	for i, tc := range p.TestCases {
		testCases[i] = problem.TestCase{
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}
	
	// Create starter code map
	starterCode := make(map[string]string)
	for _, lang := range p.Languages {
		starterCode[lang] = ""
	}
	
	return problem.Problem{
		ID:                  p.ID,
		Title:               p.Title,
		Description:         p.Description,
		Difficulty:          p.Difficulty,
		Patterns:            p.Tags,
		Companies:           p.Companies,
		TestCases:           testCases,
		StarterCode:         starterCode,
		Solutions:           make(map[string]string),
		EstimatedTime:       30, // Default value
		Examples:            []problem.Example{}, // Empty for now
		Constraints:         []string{}, // Empty for now
		PatternExplanation:  "", // Empty for now
		SolutionWalkthrough: []string{}, // Empty for now
	}
}

// convertInterfaceProblemsToLocal converts a slice of interfaces.Problem to local problem.Problem
func (m *Manager) convertInterfaceProblemsToLocal(probs []interfaces.Problem) []problem.Problem {
	result := make([]problem.Problem, len(probs))
	for i, p := range probs {
		result[i] = m.convertInterfaceToLocalProblem(p)
	}
	return result
}
