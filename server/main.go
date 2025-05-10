// API server implementation

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// License represents a user license
type License struct {
	LicenseKey   string    `json:"license_key"`
	Email        string    `json:"email"`
	PurchaseDate time.Time `json:"purchase_date"`
	ExpiryDate   time.Time `json:"expiry_date"` // For potential subscription model
	Signature    string    `json:"signature"`
}

// Problem represents an algorithm problem
type Problem struct {
	ID                  string            `json:"id"`
	Title               string            `json:"title"`
	Difficulty          string            `json:"difficulty"`
	Patterns            []string          `json:"patterns"`
	EstimatedTime       int               `json:"estimated_time"` // in minutes
	Companies           []string          `json:"companies"`
	Description         string            `json:"description"`
	Examples            []Example         `json:"examples"`
	Constraints         []string          `json:"constraints"`
	PatternExplanation  string            `json:"pattern_explanation"`
	SolutionWalkthrough []string          `json:"solution_walkthrough"`
	StarterCode         map[string]string `json:"starter_code"`
	Solutions           map[string]string `json:"solutions"`
	TestCases           []TestCase        `json:"test_cases"`
}

// Example represents an example for a problem
type Example struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation,omitempty"`
}

// TestCase represents a test case for a problem
type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

// ProblemSet represents a set of problems
type ProblemSet struct {
	Version     string    `json:"version"`
	LastUpdated time.Time `json:"last_updated"`
	Problems    []Problem `json:"problems"`
}

// Database would normally be a real database, but for demo we'll use in-memory
var (
	problemsDB = getSampleProblems()
	licensesDB = make(map[string]License)
)

func main() {
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Routes
	r.GET("/v1/problems", getProblems)
	r.POST("/v1/validate-license", validateLicense)
	r.POST("/v1/register-license", registerLicense)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	r.Run(":" + port)
}

// getProblems returns all problems
func getProblems(c *gin.Context) {
	// Verify license in request
	licenseKey := c.Query("license")
	if !isValidLicense(licenseKey) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid license",
		})
		return
	}

	c.JSON(http.StatusOK, problemsDB)
}

// validateLicense validates a license
func validateLicense(c *gin.Context) {
	// Parse request
	var req struct {
		LicenseKey string `json:"license_key"`
		Email      string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Validate license
	valid := isValidLicense(req.LicenseKey)

	c.JSON(http.StatusOK, gin.H{
		"valid": valid,
	})
}

// registerLicense registers a new license
func registerLicense(c *gin.Context) {
	// Parse request
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Generate license key
	licenseKey := generateLicenseKey(req.Email)

	// Create license
	license := License{
		LicenseKey:   licenseKey,
		Email:        req.Email,
		PurchaseDate: time.Now(),
		ExpiryDate:   time.Now().AddDate(1, 0, 0), // Valid for 1 year
		Signature:    generateSignature(licenseKey, req.Email),
	}

	// Save license
	licensesDB[licenseKey] = license

	c.JSON(http.StatusOK, gin.H{
		"license_key": licenseKey,
		"email":       req.Email,
		"expiry_date": license.ExpiryDate,
	})
}

// Helper functions

// isValidLicense checks if a license is valid
func isValidLicense(licenseKey string) bool {
	// In a real implementation, this would check a database
	// For demo, we'll validate any non-empty license key
	return licenseKey != ""
}

// generateLicenseKey generates a license key from an email
func generateLicenseKey(email string) string {
	// In a real implementation, this would generate a cryptographically secure key
	// For demo, we'll just use a simple hash
	return fmt.Sprintf("LICENSE-%s-%d", email[:4], time.Now().Unix())
}

// generateSignature creates a signature for a license
func generateSignature(licenseKey, email string) string {
	// In a real implementation, this would use private key cryptography
	// For demo, we'll just return a dummy signature
	return "valid-signature-" + licenseKey[:4]
}

// getSampleProblems returns a set of sample problems
func getSampleProblems() ProblemSet {
	// For demo, we'll use the same sample problems as in the client
	return ProblemSet{
		Version:     "1.0.0",
		LastUpdated: time.Now(),
		Problems: []Problem{
			{
				ID:            "two-sum",
				Title:         "Two Sum",
				Difficulty:    "Easy",
				Patterns:      []string{"hash-map"},
				EstimatedTime: 15,
				Companies:     []string{"Amazon", "Google", "Microsoft"},
				Description:   "Given an array of integers `nums` and an integer `target`, return indices of the two numbers such that they add up to `target`.\n\nYou may assume that each input would have exactly one solution, and you may not use the same element twice.",
				Examples: []Example{
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
				TestCases: []TestCase{
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
			{
				ID:            "max-subarray",
				Title:         "Maximum Subarray",
				Difficulty:    "Easy",
				Patterns:      []string{"dynamic-programming", "sliding-window"},
				EstimatedTime: 20,
				Companies:     []string{"Apple", "Microsoft", "Amazon"},
				Description:   "Given an integer array nums, find the contiguous subarray (containing at least one number) which has the largest sum and return its sum.",
				Examples: []Example{
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
				TestCases: []TestCase{
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
			// Add more problems here for a complete library
		},
	}
}
