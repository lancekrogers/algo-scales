package services

import (
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// ProblemService provides business logic for problem management
type ProblemService interface {
	// ListAll returns all available problems
	ListAll() ([]problem.Problem, error)
	
	// ListByPattern returns problems organized by pattern
	ListByPattern() (map[string][]problem.Problem, error)
	
	// ListByDifficulty returns problems organized by difficulty
	ListByDifficulty() (map[string][]problem.Problem, error)
	
	// GetByID retrieves a specific problem by ID
	GetByID(id string) (*problem.Problem, error)
	
	// GetRandom returns a random problem with optional filters
	GetRandom(pattern, difficulty string) (*problem.Problem, error)
}

// ProblemServiceImpl implements ProblemService
type ProblemServiceImpl struct {
	repo interfaces.ProblemRepository
}

// NewProblemService creates a new problem service
func NewProblemService(repo interfaces.ProblemRepository) ProblemService {
	if repo == nil {
		// Return legacy implementation for backward compatibility
		return &LegacyProblemService{}
	}
	
	return &ProblemServiceImpl{
		repo: repo,
	}
}

// ListAll returns all available problems
func (s *ProblemServiceImpl) ListAll() ([]problem.Problem, error) {
	return s.repo.GetAll()
}

// ListByPattern returns problems organized by pattern
func (s *ProblemServiceImpl) ListByPattern() (map[string][]problem.Problem, error) {
	allProblems, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	
	patternMap := make(map[string][]problem.Problem)
	for _, prob := range allProblems {
		for _, pattern := range prob.Patterns {
			patternMap[pattern] = append(patternMap[pattern], prob)
		}
	}
	
	return patternMap, nil
}

// ListByDifficulty returns problems organized by difficulty
func (s *ProblemServiceImpl) ListByDifficulty() (map[string][]problem.Problem, error) {
	allProblems, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	
	difficultyMap := make(map[string][]problem.Problem)
	for _, prob := range allProblems {
		difficultyMap[prob.Difficulty] = append(difficultyMap[prob.Difficulty], prob)
	}
	
	return difficultyMap, nil
}

// GetByID retrieves a specific problem by ID
func (s *ProblemServiceImpl) GetByID(id string) (*problem.Problem, error) {
	return s.repo.GetByID(id)
}

// GetRandom returns a random problem with optional filters
func (s *ProblemServiceImpl) GetRandom(pattern, difficulty string) (*problem.Problem, error) {
	if pattern != "" && difficulty != "" {
		// Filter by both pattern and difficulty
		byPattern, err := s.repo.GetByPattern(pattern)
		if err != nil {
			return nil, err
		}
		
		var filtered []problem.Problem
		for _, prob := range byPattern {
			if prob.Difficulty == difficulty {
				filtered = append(filtered, prob)
			}
		}
		
		if len(filtered) == 0 {
			return s.repo.GetRandom()
		}
		
		// Return random from filtered set
		return s.repo.GetRandomByTags([]string{pattern})
	} else if pattern != "" {
		return s.repo.GetRandomByPattern(pattern)
	} else if difficulty != "" {
		return s.repo.GetRandomByDifficulty(difficulty)
	}
	
	return s.repo.GetRandom()
}

// LegacyProblemService provides backward compatibility with legacy problem functions
type LegacyProblemService struct{}

// ListAll returns all available problems using legacy functions
func (s *LegacyProblemService) ListAll() ([]problem.Problem, error) {
	return problem.ListAll()
}

// ListByPattern returns problems organized by pattern using legacy functions
func (s *LegacyProblemService) ListByPattern() (map[string][]problem.Problem, error) {
	return problem.ListPatterns()
}

// ListByDifficulty returns problems organized by difficulty using legacy functions
func (s *LegacyProblemService) ListByDifficulty() (map[string][]problem.Problem, error) {
	return problem.ListByDifficulty()
}

// GetByID retrieves a specific problem by ID using legacy functions
func (s *LegacyProblemService) GetByID(id string) (*problem.Problem, error) {
	return problem.GetProblemByID(id)
}

// GetRandom returns a random problem with optional filters using legacy functions
func (s *LegacyProblemService) GetRandom(pattern, difficulty string) (*problem.Problem, error) {
	if pattern != "" && difficulty != "" {
		return problem.GetRandomProblemByPatternAndDifficulty(pattern, difficulty)
	} else if pattern != "" {
		return problem.GetRandomProblemByPattern(pattern)
	} else if difficulty != "" {
		return problem.GetRandomProblemByDifficulty(difficulty)
	}
	
	return problem.GetRandomProblem()
}