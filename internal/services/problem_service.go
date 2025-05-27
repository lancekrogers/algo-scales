package services

import (
	"context"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// ProblemService provides business logic for problem management
type ProblemService interface {
	// ListAll returns all available problems
	ListAll(ctx context.Context) ([]problem.Problem, error)
	
	// ListByPattern returns problems organized by pattern
	ListByPattern(ctx context.Context) (map[string][]problem.Problem, error)
	
	// ListByDifficulty returns problems organized by difficulty
	ListByDifficulty(ctx context.Context) (map[string][]problem.Problem, error)
	
	// GetByID retrieves a specific problem by ID
	GetByID(ctx context.Context, id string) (*problem.Problem, error)
	
	// GetRandom returns a random problem with optional filters
	GetRandom(ctx context.Context, pattern, difficulty string) (*problem.Problem, error)
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
func (s *ProblemServiceImpl) ListAll(ctx context.Context) ([]problem.Problem, error) {
	interfaceProblems, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert from interfaces.Problem to problem.Problem
	problems := make([]problem.Problem, len(interfaceProblems))
	for i, p := range interfaceProblems {
		problems[i] = s.convertFromInterface(p)
	}
	
	return problems, nil
}

// ListByPattern returns problems organized by pattern
func (s *ProblemServiceImpl) ListByPattern(ctx context.Context) (map[string][]problem.Problem, error) {
	interfaceProblems, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	patternMap := make(map[string][]problem.Problem)
	for _, interfaceProb := range interfaceProblems {
		prob := s.convertFromInterface(interfaceProb)
		for _, pattern := range prob.Patterns {
			patternMap[pattern] = append(patternMap[pattern], prob)
		}
	}
	
	return patternMap, nil
}

// ListByDifficulty returns problems organized by difficulty
func (s *ProblemServiceImpl) ListByDifficulty(ctx context.Context) (map[string][]problem.Problem, error) {
	interfaceProblems, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	difficultyMap := make(map[string][]problem.Problem)
	for _, interfaceProb := range interfaceProblems {
		prob := s.convertFromInterface(interfaceProb)
		difficultyMap[prob.Difficulty] = append(difficultyMap[prob.Difficulty], prob)
	}
	
	return difficultyMap, nil
}

// GetByID retrieves a specific problem by ID
func (s *ProblemServiceImpl) GetByID(ctx context.Context, id string) (*problem.Problem, error) {
	interfaceProb, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if interfaceProb == nil {
		return nil, nil
	}
	
	// Convert to local problem type
	localProb := s.convertFromInterface(*interfaceProb)
	return &localProb, nil
}

// GetRandom returns a random problem with optional filters
func (s *ProblemServiceImpl) GetRandom(ctx context.Context, pattern, difficulty string) (*problem.Problem, error) {
	var interfaceProb *interfaces.Problem
	var err error
	
	if pattern != "" && difficulty != "" {
		// Filter by both pattern and difficulty
		byPattern, err := s.repo.GetByPattern(ctx, pattern)
		if err != nil {
			return nil, err
		}
		
		var filtered []interfaces.Problem
		for _, prob := range byPattern {
			if prob.Difficulty == difficulty {
				filtered = append(filtered, prob)
			}
		}
		
		if len(filtered) == 0 {
			interfaceProb, err = s.repo.GetRandom(ctx)
		} else {
			// Return random from filtered set
			interfaceProb, err = s.repo.GetRandomByTags(ctx, []string{pattern})
		}
	} else if pattern != "" {
		interfaceProb, err = s.repo.GetRandomByPattern(ctx, pattern)
	} else if difficulty != "" {
		interfaceProb, err = s.repo.GetRandomByDifficulty(ctx, difficulty)
	} else {
		interfaceProb, err = s.repo.GetRandom(ctx)
	}
	
	if err != nil {
		return nil, err
	}
	
	if interfaceProb == nil {
		return nil, nil
	}
	
	// Convert to local problem type
	localProb := s.convertFromInterface(*interfaceProb)
	return &localProb, nil
}

// LegacyProblemService provides backward compatibility with legacy problem functions
type LegacyProblemService struct{}

// ListAll returns all available problems using legacy functions
func (s *LegacyProblemService) ListAll(ctx context.Context) ([]problem.Problem, error) {
	return problem.ListAll()
}

// ListByPattern returns problems organized by pattern using legacy functions
func (s *LegacyProblemService) ListByPattern(ctx context.Context) (map[string][]problem.Problem, error) {
	return problem.ListPatterns()
}

// ListByDifficulty returns problems organized by difficulty using legacy functions
func (s *LegacyProblemService) ListByDifficulty(ctx context.Context) (map[string][]problem.Problem, error) {
	return problem.ListByDifficulty()
}

// GetByID retrieves a specific problem by ID using legacy functions
func (s *LegacyProblemService) GetByID(ctx context.Context, id string) (*problem.Problem, error) {
	return problem.GetByID(id)
}

// convertFromInterface converts an interfaces.Problem to a problem.Problem
func (s *ProblemServiceImpl) convertFromInterface(p interfaces.Problem) problem.Problem {
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
	if p.StarterCode != nil {
		starterCode = p.StarterCode
	} else {
		for _, lang := range p.Languages {
			starterCode[lang] = ""
		}
	}
	
	return problem.Problem{
		ID:                  p.ID,
		Title:               p.Title,
		Description:         p.Description,
		Difficulty:          p.Difficulty,
		Patterns:            p.Tags, // Map Tags to Patterns
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

// GetRandom returns a random problem with optional filters using legacy functions
func (s *LegacyProblemService) GetRandom(ctx context.Context, pattern, difficulty string) (*problem.Problem, error) {
	// Use simple random selection for legacy service
	// This is a simplified implementation
	return problem.GetRandomProblem()
}