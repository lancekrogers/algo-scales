package splitscreen

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// Variables to allow mocking for testing
var (
	osExit = os.Exit
	
	runProgram = func(m tea.Model, opts ...tea.ProgramOption) (tea.Model, error) {
		p := tea.NewProgram(m, opts...)
		return p.Run()
	}
)

// StartUI launches the split-screen UI with the given problem
func StartUI(p *problem.Problem) error {
	// Create the model
	m := NewModel()
	
	// Set the current problem if provided
	if p != nil {
		m.SetProblem(p)
	}
	
	// Program options
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	}
	
	// Create and run the program
	_, err := runProgram(m, opts...)
	
	// Check for errors
	if err != nil {
		return fmt.Errorf("error running split-screen UI: %v", err)
	}
	
	return nil
}

// StartWithSampleProblem starts the UI with a sample problem for demonstration
func StartWithSampleProblem() error {
	// Create a sample problem
	sampleProblem := &problem.Problem{
		ID:         "sample-two-sum",
		Title:      "Two Sum",
		Difficulty: "easy",
		Description: `Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.

You may assume that each input would have exactly one solution, and you may not use the same element twice.

You can return the answer in any order.`,
		Examples: []problem.Example{
			{
				Input:  "nums = [2,7,11,15], target = 9",
				Output: "[0,1]",
				Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
			},
			{
				Input:  "nums = [3,2,4], target = 6",
				Output: "[1,2]",
			},
		},
		Constraints: []string{
			"2 <= nums.length <= 10^4",
			"-10^9 <= nums[i] <= 10^9",
			"-10^9 <= target <= 10^9",
			"Only one valid answer exists.",
		},
		Solutions: map[string]string{
			"go": `func twoSum(nums []int, target int) []int {
    numMap := make(map[int]int)
    
    for i, num := range nums {
        complement := target - num
        if idx, found := numMap[complement]; found {
            return []int{idx, i}
        }
        numMap[num] = i
    }
    
    return nil // No solution found
}`,
			"python": `def twoSum(nums, target):
    num_map = {}
    
    for i, num in enumerate(nums):
        complement = target - num
        if complement in num_map:
            return [num_map[complement], i]
        num_map[num] = i
    
    return None  # No solution found`,
			"javascript": `function twoSum(nums, target) {
    const numMap = new Map();
    
    for (let i = 0; i < nums.length; i++) {
        const complement = target - nums[i];
        if (numMap.has(complement)) {
            return [numMap.get(complement), i];
        }
        numMap.set(nums[i], i);
    }
    
    return null; // No solution found
}`,
		},
		Patterns:          []string{"Hash Map", "Array"},
		EstimatedTime:     15,
		PatternExplanation: "This problem can be solved efficiently using a hash map. As we iterate through the array, we check if the complement (target - current number) exists in our hash map. If it does, we've found our pair. If not, we add the current number to the hash map and continue.",
	}
	
	// Launch UI with the sample problem
	return StartUI(sampleProblem)
}

// Variables for testing mocks
var (
	startUIProg = StartUI
	getProblemByID = problem.GetByID
	
	// For testing RunCLI
	runStartWithSampleProblem = func() error {
		return StartWithSampleProblem()
	}
	
	runGetProblemAndStartUI = func(id string) error {
		return GetProblemAndStartUI(id) 
	}
)

// GetProblemAndStartUI loads a problem by ID and starts the UI
func GetProblemAndStartUI(problemID string) error {
	// Try to load the problem
	p, err := getProblemByID(problemID)
	if err != nil {
		return fmt.Errorf("error loading problem: %v", err)
	}
	
	// Start the UI with the loaded problem
	return startUIProg(p)
}

// RunCLI provides a command-line interface for running the split-screen UI
func RunCLI(args []string) int {
	// Check for a problem ID argument
	if len(args) > 0 {
		problemID := args[0]
		err := runGetProblemAndStartUI(problemID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return 1
		}
		return 0
	}
	
	// No problem ID provided, start with a sample problem
	err := runStartWithSampleProblem()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	
	return 0
}