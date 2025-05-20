package problem

import (
	"fmt"
	"math/rand"
	"time"
)

// GetRandomProblemByPattern finds a random problem with the specified pattern
var GetRandomProblemByPattern = func(pattern string) (*Problem, error) {
	// Load all problems
	problems, err := ListAll()
	if err != nil {
		return nil, err
	}

	// Filter problems by pattern
	var filteredProblems []Problem
	for _, p := range problems {
		for _, pat := range p.Patterns {
			if pat == pattern {
				filteredProblems = append(filteredProblems, p)
				break
			}
		}
	}

	if len(filteredProblems) == 0 {
		return nil, fmt.Errorf("no problems found with pattern: %s", pattern)
	}

	// Pick a random problem
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(filteredProblems))
	return &filteredProblems[randomIndex], nil
}

// GetRandomProblem selects a random problem from all available problems
var GetRandomProblem = func() (*Problem, error) {
	// Load all problems
	problems, err := ListAll()
	if err != nil {
		return nil, err
	}

	if len(problems) == 0 {
		return nil, fmt.Errorf("no problems found")
	}

	// Pick a random problem
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(problems))
	return &problems[randomIndex], nil
}

// GetRandomProblemByDifficulty finds a random problem with the specified difficulty
var GetRandomProblemByDifficulty = func(difficulty string) (*Problem, error) {
	// Load all problems
	problems, err := ListAll()
	if err != nil {
		return nil, err
	}

	// Filter problems by difficulty
	var filteredProblems []Problem
	for _, p := range problems {
		if p.Difficulty == difficulty {
			filteredProblems = append(filteredProblems, p)
		}
	}

	if len(filteredProblems) == 0 {
		return nil, fmt.Errorf("no problems found with difficulty: %s", difficulty)
	}

	// Pick a random problem
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(filteredProblems))
	return &filteredProblems[randomIndex], nil
}

// GetRandomProblemExcluding finds a random problem that is not in the excluded list
var GetRandomProblemExcluding = func(excludedIDs []string) (*Problem, error) {
	// Load all problems
	problems, err := ListAll()
	if err != nil {
		return nil, err
	}

	// Filter out excluded problems
	var filteredProblems []Problem
	for _, p := range problems {
		excluded := false
		for _, id := range excludedIDs {
			if p.ID == id {
				excluded = true
				break
			}
		}
		if !excluded {
			filteredProblems = append(filteredProblems, p)
		}
	}

	if len(filteredProblems) == 0 {
		return nil, fmt.Errorf("no problems available after exclusions")
	}

	// Pick a random problem
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(filteredProblems))
	return &filteredProblems[randomIndex], nil
}