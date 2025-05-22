// Session creation functionality
package session

import (
	"fmt"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// CreateSession creates a session without starting the UI
func CreateSession(opts Options) (*Session, error) {
	// Convert old Options to new SessionOptions
	sessionOpts := interfaces.SessionOptions{
		Mode:       interfaces.SessionMode(opts.Mode),
		Language:   opts.Language,
		Timer:      opts.Timer,
		Pattern:    opts.Pattern,
		Difficulty: opts.Difficulty,
		ProblemID:  opts.ProblemID,
	}
	
	// Create a manager to handle session creation
	manager := NewManager()
	
	// Use the manager to start the session
	interfaceSession, err := manager.StartSession(sessionOpts)
	if err != nil {
		return nil, err
	}
	
	// Convert the interface session back to legacy Session type
	// Note: This maintains backward compatibility with existing code
	sessionImpl, ok := interfaceSession.(*SessionImpl)
	if !ok {
		return nil, fmt.Errorf("unexpected session implementation type")
	}
	
	// Create legacy Session struct for backward compatibility
	legacySession := &Session{
		Options:      opts,
		Problem:      sessionImpl.Problem,
		StartTime:    sessionImpl.StartTime,
		ShowHints:    sessionImpl.hintsShown,
		ShowPattern:  sessionImpl.ShowPattern,
		ShowSolution: sessionImpl.solutionShown,
		Workspace:    sessionImpl.Workspace,
		CodeFile:     sessionImpl.CodeFile,
	}

	// Workspace is already created by the manager, so we can return directly
	return legacySession, nil
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
