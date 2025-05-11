// API client for problem downloads

// internal/api/client.go
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/lancekrogers/algo-scales/internal/license"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

const (
	// In a real implementation, this would be a remote API
	// For MVP, we'll use local sample data
	BaseURL = "https://api.algo-scales.com/v1"
)

// DownloadProblems fetches problems from the API and stores them locally
// For MVP, we'll use sample data instead of making actual API calls
func DownloadProblems(force bool) error {
	// Check license validity
	valid, err := license.ValidateLicense()
	if err != nil || !valid {
		return fmt.Errorf("invalid license: %v", err)
	}

	// Check if we need to update
	if !force && !shouldUpdate() {
		return nil // No update needed
	}

	// Create necessary directories
	configDir := getConfigDir()
	problemsDir := filepath.Join(configDir, "problems")

	if err := os.MkdirAll(problemsDir, 0755); err != nil {
		return err
	}

	// For MVP, we'll use embedded sample problems instead of API call
	problemSet := getSampleProblems()

	// Save version info
	versionFile := filepath.Join(configDir, "version.json")
	versionData, err := json.MarshalIndent(map[string]interface{}{
		"version":      "1.0.0",
		"last_updated": time.Now(),
	}, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(versionFile, versionData, 0644); err != nil {
		return err
	}

	// Save individual problems
	for _, p := range problemSet.Problems {
		// Create pattern directories
		for _, pattern := range p.Patterns {
			patternDir := filepath.Join(problemsDir, pattern)
			if err := os.MkdirAll(patternDir, 0755); err != nil {
				return err
			}

			// Save problem file
			problemData, err := json.MarshalIndent(p, "", "  ")
			if err != nil {
				return err
			}

			problemFile := filepath.Join(patternDir, fmt.Sprintf("%s.json", p.ID))
			if err := ioutil.WriteFile(problemFile, problemData, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

// shouldUpdate checks if we need to update problem sets
func shouldUpdate() bool {
	configDir := getConfigDir()
	versionFile := filepath.Join(configDir, "version.json")

	// If version file doesn't exist, we need to update
	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		return true
	}

	// Read version file
	data, err := ioutil.ReadFile(versionFile)
	if err != nil {
		return true
	}

	var version struct {
		Version     string    `json:"version"`
		LastUpdated time.Time `json:"last_updated"`
	}

	if err := json.Unmarshal(data, &version); err != nil {
		return true
	}

	// For MVP, we'll always return false after initial download
	// In a real implementation, you'd check against the server version
	return false
}

// getConfigDir returns the configuration directory
// Exported as variable for testing
var getConfigDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".algo-scales")
}

// ProblemSet represents a set of problems
type ProblemSet struct {
	Version     string            `json:"version"`
	LastUpdated time.Time         `json:"last_updated"`
	Problems    []problem.Problem `json:"problems"`
}

// getSampleProblems returns a set of sample problems for MVP
func getSampleProblems() ProblemSet {
	// Generate sample problems for different patterns
	return ProblemSet{
		Version:     "1.0.0",
		LastUpdated: time.Now(),
		Problems: []problem.Problem{
			{
				ID:            "two-sum",
				Title:         "Two Sum",
				Difficulty:    "Easy",
				Patterns:      []string{"hash-map"},
				EstimatedTime: 15,
				Companies:     []string{"Amazon", "Google", "Microsoft"},
				Description:   "Given an array of integers `nums` and an integer `target`, return indices of the two numbers such that they add up to `target`.\n\nYou may assume that each input would have exactly one solution, and you may not use the same element twice.",
				Examples: []problem.Example{
					{
						Input:       "nums = [2,7,11,15], target = 9",
						Output:      "[0,1]",
						Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
					},
				},
				Constraints: []string{
					"2 <= nums.length <= 10^4",
					"-10^9 <= nums[i] <= 10^9",
					"-10^9 <= target <= 10^9",
					"Only one valid answer exists.",
				},
				PatternExplanation: "This problem demonstrates the Hash Map pattern, which provides O(1) lookups. For each element, we check if its complement (target - current) exists in our hash map. If not, we add the current element to the map and continue.",
				SolutionWalkthrough: []string{
					"Initialize an empty hash map to store values and their indices",
					"Iterate through the array:",
					"  - For each element, calculate the complement (target - current)",
					"  - Check if the complement exists in the hash map",
					"  - If found, return the current index and the complement's index",
					"  - Otherwise, add the current element and its index to the hash map",
					"If no solution is found, return an empty array (though problem states a solution always exists)",
				},
				StarterCode: map[string]string{
					"go":         "func twoSum(nums []int, target int) []int {\n    // Your code here\n}",
					"python":     "def two_sum(nums, target):\n    # Your code here\n    pass",
					"javascript": "function twoSum(nums, target) {\n    // Your code here\n}",
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
					{
						Input:    "[3,3], 6",
						Expected: "[0,1]",
					},
				},
			},
			// Add more sample problems for different patterns
			{
				ID:            "max-subarray",
				Title:         "Maximum Subarray",
				Difficulty:    "Easy",
				Patterns:      []string{"dynamic-programming", "sliding-window"},
				EstimatedTime: 20,
				Companies:     []string{"Apple", "Microsoft", "Amazon"},
				Description:   "Given an integer array nums, find the contiguous subarray (containing at least one number) which has the largest sum and return its sum.",
				Examples: []problem.Example{
					{
						Input:       "nums = [-2,1,-3,4,-1,2,1,-5,4]",
						Output:      "6",
						Explanation: "The subarray [4,-1,2,1] has the largest sum 6.",
					},
				},
				Constraints: []string{
					"1 <= nums.length <= 3 * 10^4",
					"-10^5 <= nums[i] <= 10^5",
				},
				PatternExplanation: "This problem demonstrates the dynamic programming pattern. We can solve it by keeping track of the maximum sum ending at the current position and the global maximum sum seen so far.",
				SolutionWalkthrough: []string{
					"Initialize two variables: currentSum and maxSum, both set to the first element.",
					"Iterate through the array starting from the second element:",
					"  - Update currentSum as the maximum of the current element and the sum of currentSum and the current element.",
					"  - Update maxSum as the maximum of maxSum and currentSum.",
					"Return maxSum as the result.",
				},
				StarterCode: map[string]string{
					"go":         "func maxSubArray(nums []int) int {\n    // Your code here\n}",
					"python":     "def max_subarray(nums):\n    # Your code here\n    pass",
					"javascript": "function maxSubArray(nums) {\n    // Your code here\n}",
				},
				Solutions: map[string]string{
					"go":         "func maxSubArray(nums []int) int {\n    if len(nums) == 0 {\n        return 0\n    }\n    \n    currentSum := nums[0]\n    maxSum := nums[0]\n    \n    for i := 1; i < len(nums); i++ {\n        currentSum = max(nums[i], currentSum + nums[i])\n        maxSum = max(maxSum, currentSum)\n    }\n    \n    return maxSum\n}\n\nfunc max(a, b int) int {\n    if a > b {\n        return a\n    }\n    return b\n}",
					"python":     "def max_subarray(nums):\n    if not nums:\n        return 0\n        \n    current_sum = max_sum = nums[0]\n    \n    for num in nums[1:]:\n        current_sum = max(num, current_sum + num)\n        max_sum = max(max_sum, current_sum)\n        \n    return max_sum",
					"javascript": "function maxSubArray(nums) {\n    if (nums.length === 0) {\n        return 0;\n    }\n    \n    let currentSum = nums[0];\n    let maxSum = nums[0];\n    \n    for (let i = 1; i < nums.length; i++) {\n        currentSum = Math.max(nums[i], currentSum + nums[i]);\n        maxSum = Math.max(maxSum, currentSum);\n    }\n    \n    return maxSum;\n}",
				},
				TestCases: []problem.TestCase{
					{
						Input:    "[-2,1,-3,4,-1,2,1,-5,4]",
						Expected: "6",
					},
					{
						Input:    "[1]",
						Expected: "1",
					},
					{
						Input:    "[5,4,-1,7,8]",
						Expected: "23",
					},
				},
			},
		},
	}
}
