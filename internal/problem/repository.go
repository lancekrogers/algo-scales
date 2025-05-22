package problem

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/common/utils"
)

// Repository implements the ProblemRepository interface
type Repository struct {
	fs interfaces.FileSystem
}

// NewRepository creates a new problem repository with the default file system
func NewRepository() interfaces.ProblemRepository {
	return &Repository{
		fs: utils.NewFileSystem(),
	}
}

// WithFileSystem returns a repository with a custom file system
func (r *Repository) WithFileSystem(fs interfaces.FileSystem) *Repository {
	return &Repository{fs: fs}
}

// GetAll returns all available problems
func (r *Repository) GetAll() ([]interfaces.Problem, error) {
	problems, err := r.getAllLocal()
	if err != nil {
		return nil, err
	}
	
	// Convert to interface types
	result := make([]interfaces.Problem, len(problems))
	for i, p := range problems {
		result[i] = r.convertToInterface(p)
	}
	return result, nil
}

// getAllLocal returns all problems as local Problem types
func (r *Repository) getAllLocal() ([]Problem, error) {
	// First try the standard config dir location
	configDir := r.fs.GetConfigDir()
	problemsDir := filepath.Join(configDir, "problems")
	
	// If problems directory doesn't exist in config dir,
	// try the local problems directory relative to binary
	if !r.fs.Exists(problemsDir) {
		// Get the executable directory
		exePath, err := r.fs.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %v", err)
		}
		
		exeDir := filepath.Dir(exePath)
		problemsDir = filepath.Join(exeDir, "problems")
		
		// If still no problems directory, try current directory
		if !r.fs.Exists(problemsDir) {
			curDir, err := r.fs.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get current directory: %v", err)
			}
			
			problemsDir = filepath.Join(curDir, "problems")
			
			// If still no problems directory, try project root
			if !r.fs.Exists(problemsDir) {
				// Try going up one directory (assuming we're in algo-scales/...)
				rootDir := filepath.Dir(curDir)
				problemsDir = filepath.Join(rootDir, "algo-scales", "problems")
				
				// If still no problems directory, return empty result
				if !r.fs.Exists(problemsDir) {
					return []Problem{}, nil
				}
			}
		}
	}
	
	// Track processed problem IDs to avoid duplicates
	var problems []Problem
	processedIDs := make(map[string]bool)
	
	// Get pattern directories
	patternDirs, err := r.fs.ReadDir(problemsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read problems directory: %v", err)
	}
	
	// Iterate through pattern directories
	for _, patternDir := range patternDirs {
		if !patternDir.IsDir() {
			continue
		}
		
		patternName := patternDir.Name()
		patternPath := filepath.Join(problemsDir, patternName)
		
		// Read problem files in the pattern directory
		problemFiles, err := r.fs.ReadDir(patternPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read pattern directory %s: %v", patternName, err)
		}
		
		for _, file := range problemFiles {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}
			
			// Read problem file
			problemPath := filepath.Join(patternPath, file.Name())
			data, err := r.fs.ReadFile(problemPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read problem file %s: %v", problemPath, err)
			}
			
			// Parse problem data
			var problem Problem
			err = json.Unmarshal(data, &problem)
			if err != nil {
				return nil, fmt.Errorf("failed to parse problem file %s: %v", problemPath, err)
			}
			
			// Skip if already processed
			if processedIDs[problem.ID] {
				continue
			}
			
			// Add problem to the list
			problems = append(problems, problem)
			processedIDs[problem.ID] = true
		}
	}
	
	// Sort problems by difficulty (easy, medium, hard)
	sort.Slice(problems, func(i, j int) bool {
		// Define difficulty order
		difficultyOrder := map[string]int{
			"easy":   0,
			"medium": 1,
			"hard":   2,
		}
		
		// Get difficulty values
		diffI := difficultyOrder[problems[i].Difficulty]
		diffJ := difficultyOrder[problems[j].Difficulty]
		
		// Sort by difficulty first
		if diffI != diffJ {
			return diffI < diffJ
		}
		
		// Then by ID for consistent ordering
		return problems[i].ID < problems[j].ID
	})
	
	return problems, nil
}

// GetByID retrieves a specific problem by its ID
func (r *Repository) GetByID(id string) (*interfaces.Problem, error) {
	problem, err := r.getByIDLocal(id)
	if err != nil {
		return nil, err
	}
	
	converted := r.convertToInterface(*problem)
	return &converted, nil
}

// getByIDLocal retrieves a specific problem by its ID as local type
func (r *Repository) getByIDLocal(id string) (*Problem, error) {
	configDir := r.fs.GetConfigDir()
	
	// Search in all pattern directories
	patternDirs, err := r.fs.ReadDir(filepath.Join(configDir, "problems"))
	if err != nil {
		return nil, err
	}
	
	for _, patternDir := range patternDirs {
		if !patternDir.IsDir() {
			continue
		}
		
		problemPath := filepath.Join(configDir, "problems", patternDir.Name(), fmt.Sprintf("%s.json", id))
		if !r.fs.Exists(problemPath) {
			continue
		}
		
		// Found the problem file
		data, err := r.fs.ReadFile(problemPath)
		if err != nil {
			return nil, err
		}
		
		var problem Problem
		if err := json.Unmarshal(data, &problem); err != nil {
			return nil, err
		}
		
		return &problem, nil
	}
	
	return nil, ErrProblemNotFound
}

// GetByPattern returns problems matching a specific pattern
func (r *Repository) GetByPattern(pattern string) ([]interfaces.Problem, error) {
	allProblems, err := r.getAllLocal()
	if err != nil {
		return nil, err
	}
	
	var filtered []Problem
	if pattern == "" {
		filtered = allProblems
	} else {
		for _, p := range allProblems {
			for _, patternName := range p.Patterns {
				if patternName == pattern {
					filtered = append(filtered, p)
					break
				}
			}
		}
	}
	
	// Convert to interface types
	result := make([]interfaces.Problem, len(filtered))
	for i, p := range filtered {
		result[i] = r.convertToInterface(p)
	}
	return result, nil
}

// GetPatterns returns all available algorithm patterns
func (r *Repository) GetPatterns() ([]string, error) {
	allProblems, err := r.getAllLocal()
	if err != nil {
		return nil, err
	}
	
	// Use map to track unique patterns
	patterns := make(map[string]bool)
	
	for _, problem := range allProblems {
		for _, pattern := range problem.Patterns {
			// Convert kebab-case to Title Case for display
			displayPattern := convertPatternToDisplay(pattern)
			patterns[displayPattern] = true
		}
	}
	
	// Convert map to sorted slice
	result := make([]string, 0, len(patterns))
	for pattern := range patterns {
		result = append(result, pattern)
	}
	
	// Sort patterns for consistent ordering
	sort.Strings(result)
	
	return result, nil
}

// GetLanguages returns all available programming languages
func (r *Repository) GetLanguages() ([]string, error) {
	allProblems, err := r.getAllLocal()
	if err != nil {
		return nil, err
	}
	
	// Use map to track unique languages
	languages := make(map[string]bool)
	
	for _, problem := range allProblems {
		for lang := range problem.StarterCode {
			languages[lang] = true
		}
	}
	
	// Convert map to sorted slice
	result := make([]string, 0, len(languages))
	for lang := range languages {
		result = append(result, lang)
	}
	
	// Sort languages for consistent ordering
	sort.Strings(result)
	
	return result, nil
}

// GetByDifficulty returns problems with a specific difficulty level
func (r *Repository) GetByDifficulty(difficulty string) ([]interfaces.Problem, error) {
	allProblems, err := r.getAllLocal()
	if err != nil {
		return nil, err
	}
	
	var filtered []Problem
	for _, p := range allProblems {
		if p.Difficulty == difficulty {
			filtered = append(filtered, p)
		}
	}
	
	// Convert to interface types
	result := make([]interfaces.Problem, len(filtered))
	for i, p := range filtered {
		result[i] = r.convertToInterface(p)
	}
	return result, nil
}

// GetByCompany returns problems from a specific company
func (r *Repository) GetByCompany(company string) ([]Problem, error) {
	allProblems, err := r.getAllLocal()
	if err != nil {
		return nil, err
	}
	
	var filtered []Problem
	for _, p := range allProblems {
		for _, c := range p.Companies {
			if c == company {
				filtered = append(filtered, p)
				break
			}
		}
	}
	
	return filtered, nil
}

// convertToInterface converts a local Problem to interfaces.Problem
func (r *Repository) convertToInterface(p Problem) interfaces.Problem {
	var pattern string
	if len(p.Patterns) > 0 {
		pattern = p.Patterns[0] // Use first pattern for simplicity
	}
	
	var languages []string
	for lang := range p.StarterCode {
		languages = append(languages, lang)
	}
	
	// Convert test cases
	testCases := make([]interfaces.TestCase, len(p.TestCases))
	for i, tc := range p.TestCases {
		testCases[i] = interfaces.TestCase{
			Input:    tc.Input,
			Expected: tc.Expected,
		}
	}
	
	return interfaces.Problem{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Pattern:     pattern,
		Difficulty:  p.Difficulty,
		Companies:   p.Companies,
		Tags:        p.Patterns, // Use patterns as tags
		TestCases:   testCases,
		Languages:   languages,
	}
}

// GetByTags returns problems matching any of the specified tags
func (r *Repository) GetByTags(tags []string) ([]interfaces.Problem, error) {
	allProblems, err := r.getAllLocal()
	if err != nil {
		return nil, err
	}
	
	var filtered []Problem
	for _, p := range allProblems {
		// Check if any pattern matches any tag
		for _, pattern := range p.Patterns {
			for _, tag := range tags {
				if pattern == tag {
					filtered = append(filtered, p)
					goto next_problem
				}
			}
		}
		next_problem:
	}
	
	// Convert to interface types
	result := make([]interfaces.Problem, len(filtered))
	for i, p := range filtered {
		result[i] = r.convertToInterface(p)
	}
	return result, nil
}

// GetRandom returns a random problem
func (r *Repository) GetRandom() (*interfaces.Problem, error) {
	problems, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple random selection (could be improved with proper randomization)
	randomIndex := len(problems) / 2 // Simple selection for now
	return &problems[randomIndex], nil
}

// GetRandomByPattern returns a random problem matching a pattern
func (r *Repository) GetRandomByPattern(pattern string) (*interfaces.Problem, error) {
	problems, err := r.GetByPattern(pattern)
	if err != nil {
		return nil, err
	}
	
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple random selection
	randomIndex := len(problems) / 2
	return &problems[randomIndex], nil
}

// GetRandomByDifficulty returns a random problem with a difficulty
func (r *Repository) GetRandomByDifficulty(difficulty string) (*interfaces.Problem, error) {
	problems, err := r.GetByDifficulty(difficulty)
	if err != nil {
		return nil, err
	}
	
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple random selection
	randomIndex := len(problems) / 2
	return &problems[randomIndex], nil
}

// GetRandomByTags returns a random problem matching tags
func (r *Repository) GetRandomByTags(tags []string) (*interfaces.Problem, error) {
	problems, err := r.GetByTags(tags)
	if err != nil {
		return nil, err
	}
	
	if len(problems) == 0 {
		return nil, ErrProblemNotFound
	}
	
	// Use simple random selection
	randomIndex := len(problems) / 2
	return &problems[randomIndex], nil
}