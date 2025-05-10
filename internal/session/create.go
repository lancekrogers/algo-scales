// Session creation functionality
package session

import (
	"fmt"
	"time"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

// CreateSession creates a session without starting the UI
func CreateSession(opts Options) (*Session, error) {
	// Initialize session
	session := &Session{
		Options:      opts,
		StartTime:    time.Now(),
		ShowHints:    opts.Mode == LearnMode,
		ShowPattern:  opts.Mode == LearnMode,
		ShowSolution: false,
	}

	// Choose problem based on options
	var err error
	if opts.ProblemID != "" {
		// Specific problem requested
		session.Problem, err = problem.GetByID(opts.ProblemID)
		if err != nil {
			return nil, fmt.Errorf("failed to load problem: %v", err)
		}
	} else if opts.Mode == CramMode {
		// Cram mode - choose problems from common patterns
		session.Problem, err = selectCramProblem()
		if err != nil {
			return nil, fmt.Errorf("failed to select problem for cram mode: %v", err)
		}
	} else {
		// Filter by pattern/difficulty if specified
		session.Problem, err = selectProblem(opts.Pattern, opts.Difficulty)
		if err != nil {
			return nil, fmt.Errorf("failed to select problem: %v", err)
		}
	}

	// Create workspace
	if err := session.createWorkspace(); err != nil {
		return nil, fmt.Errorf("failed to create workspace: %v", err)
	}

	return session, nil
}

// getDefaultProblem returns a fallback problem if no problems match filters
func getDefaultProblem() *problem.Problem {
	// Create a simple default problem as a fallback
	return &problem.Problem{
		ID:            "two-sum",
		Title:         "Two Sum",
		Difficulty:    "Easy",
		Patterns:      []string{"hash-map"},
		EstimatedTime: 15,
		Companies:     []string{"Many Companies"},
		Description:   "Given an array of integers and a target sum, return the indices of the two numbers that add up to the target.",
		Examples: []problem.Example{
			{
				Input:       "nums = [2,7,11,15], target = 9",
				Output:      "[0,1]",
				Explanation: "Because nums[0] + nums[1] = 2 + 7 = 9",
			},
		},
		Constraints: []string{
			"2 <= nums.length <= 10^4",
			"-10^9 <= nums[i] <= 10^9",
			"Only one valid answer exists.",
		},
		PatternExplanation: "This problem demonstrates the hash map pattern, which allows us to find elements in constant time.",
		SolutionWalkthrough: []string{
			"Create a hash map to store values we've seen",
			"Iterate through the array",
			"For each element, check if its complement exists in the map",
			"If it does, return the current index and the complement's index",
			"Otherwise, add the current element to the map",
		},
		StarterCode: map[string]string{
			"go":         "func twoSum(nums []int, target int) []int {\n    // Your solution here\n}",
			"python":     "def two_sum(nums, target):\n    # Your solution here\n    pass",
			"javascript": "function twoSum(nums, target) {\n    // Your solution here\n}",
		},
		Solutions: map[string]string{
			"go":         "func twoSum(nums []int, target int) []int {\n    seen := make(map[int]int)\n    for i, num := range nums {\n        complement := target - num\n        if j, ok := seen[complement]; ok {\n            return []int{j, i}\n        }\n        seen[num] = i\n    }\n    return []int{}\n}",
			"python":     "def two_sum(nums, target):\n    seen = {}\n    for i, num in enumerate(nums):\n        complement = target - num\n        if complement in seen:\n            return [seen[complement], i]\n        seen[num] = i\n    return []",
			"javascript": "function twoSum(nums, target) {\n    const seen = {};\n    for (let i = 0; i < nums.length; i++) {\n        const complement = target - nums[i];\n        if (complement in seen) {\n            return [seen[complement], i];\n        }\n        seen[nums[i]] = i;\n    }\n    return [];\n}",
		},
		TestCases: []problem.TestCase{
			{
				Input:    "[2,7,11,15], 9",
				Expected: "[0,1]",
			},
			{
				Input:    "[3,2,4], 6",
				Expected: "[1,2]",
			},
		},
	}
}
