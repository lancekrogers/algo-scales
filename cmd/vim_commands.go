// Vim mode commands for Neovim plugin integration

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/services"
	"github.com/lancekrogers/algo-scales/internal/session/execution"
	"github.com/spf13/cobra"
)

// submitCmd represents the submit command for vim mode
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit solution for testing (vim mode)",
	Long:  `Submit a solution file for testing. Used by the Neovim plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		problemID, _ := cmd.Flags().GetString("problem-id")
		language, _ := cmd.Flags().GetString("language")
		filePath, _ := cmd.Flags().GetString("file")
		isVimMode, _ := cmd.Flags().GetBool("vim-mode")

		if !isVimMode {
			fmt.Println("This command is for vim mode only")
			return
		}

		// Create context - in production, this should have timeout
		ctx := context.Background()

		// Read the solution file
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			outputVimError(fmt.Errorf("failed to read file: %v", err))
			return
		}

		// Get problem from repository
		problemService := services.DefaultRegistry.GetProblemService()
		prob, err := problemService.GetByID(ctx, problemID)
		if err != nil {
			outputVimError(fmt.Errorf("failed to get problem: %v", err))
			return
		}

		// Get test runner registry
		registry := execution.NewRunnerRegistry()
		runner, err := registry.GetRunner(language)
		if err != nil {
			outputVimError(fmt.Errorf("unsupported language: %v", err))
			return
		}

		// Run tests directly
		// Convert test cases to interface type
		var interfaceTestCases []interfaces.TestCase
		for _, tc := range prob.TestCases {
			interfaceTestCases = append(interfaceTestCases, interfaces.TestCase{
				Input:    tc.Input,
				Expected: tc.Expected,
			})
		}
		
		interfaceProb := &interfaces.Problem{
			ID:          prob.ID,
			Title:       prob.Title,
			Description: prob.Description,
			TestCases:   interfaceTestCases,
		}
		
		results, _, err := runner.ExecuteTests(ctx, interfaceProb, string(content), 30*time.Second)
		if err != nil {
			outputVimError(fmt.Errorf("failed to run tests: %v", err))
			return
		}

		// Convert to vim response format
		var testResults []TestResult
		allPassed := true
		for _, result := range results {
			tr := TestResult{
				Input:    fmt.Sprintf("%v", result.Input),
				Expected: fmt.Sprintf("%v", result.Expected),
				Actual:   fmt.Sprintf("%v", result.Actual),
				Passed:   result.Passed,
			}
			testResults = append(testResults, tr)
			if !result.Passed {
				allPassed = false
			}
		}

		// Create and output response
		resp := VimSubmitResponse{
			Passed:      allPassed,
			TestResults: testResults,
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			outputVimError(fmt.Errorf("failed to marshal response: %v", err))
			return
		}

		fmt.Println(string(jsonResp))
	},
}

// testCmd represents the test command for vim mode
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests on solution (vim mode)",
	Long:  `Run tests on a solution file. Used by the Neovim plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Same implementation as submit for now
		submitCmd.Run(cmd, args)
	},
}

// hintCmd represents the hint command for vim mode
var hintCmd = &cobra.Command{
	Use:   "hint",
	Short: "Get hint for problem (vim mode)",
	Long:  `Get a hint for the specified problem. Used by the Neovim plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		problemID, _ := cmd.Flags().GetString("problem-id")
		isVimMode, _ := cmd.Flags().GetBool("vim-mode")

		if !isVimMode {
			fmt.Println("This command is for vim mode only")
			return
		}

		// Create context
		ctx := context.Background()

		// Get problem from repository
		problemService := services.DefaultRegistry.GetProblemService()
		prob, err := problemService.GetByID(ctx, problemID)
		if err != nil {
			outputVimError(fmt.Errorf("failed to get problem: %v", err))
			return
		}

		// Get hint text
		hintText := prob.PatternExplanation
		if hintText == "" {
			hintText = "Think about the pattern: " + getPatternHint(prob.Patterns)
		}

		// Create and output response
		resp := VimHintResponse{
			Hint: hintText,
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			outputVimError(fmt.Errorf("failed to marshal response: %v", err))
			return
		}

		fmt.Println(string(jsonResp))
	},
}

// solutionCmd represents the solution command for vim mode
var solutionCmd = &cobra.Command{
	Use:   "solution",
	Short: "Get solution for problem (vim mode)",
	Long:  `Get the solution for the specified problem. Used by the Neovim plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		problemID, _ := cmd.Flags().GetString("problem-id")
		language, _ := cmd.Flags().GetString("language")
		isVimMode, _ := cmd.Flags().GetBool("vim-mode")

		if !isVimMode {
			fmt.Println("This command is for vim mode only")
			return
		}

		// Create context
		ctx := context.Background()

		// Get problem from repository
		problemService := services.DefaultRegistry.GetProblemService()
		prob, err := problemService.GetByID(ctx, problemID)
		if err != nil {
			outputVimError(fmt.Errorf("failed to get problem: %v", err))
			return
		}

		// Get solution code
		solutionCode := ""
		if prob.Solutions != nil {
			if code, ok := prob.Solutions[language]; ok {
				solutionCode = code
			} else {
				// Try to get any solution
				for _, code := range prob.Solutions {
					solutionCode = code
					break
				}
			}
		}

		if solutionCode == "" {
			solutionCode = "// Solution not available for this problem"
		}

		// Create and output response
		resp := VimSolutionResponse{
			Solution: solutionCode,
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			outputVimError(fmt.Errorf("failed to marshal response: %v", err))
			return
		}

		fmt.Println(string(jsonResp))
	},
}

// Helper function to output vim mode errors
func outputVimError(err error) {
	errResp := map[string]string{
		"error": err.Error(),
	}
	jsonResp, _ := json.Marshal(errResp)
	fmt.Println(string(jsonResp))
	os.Exit(1)
}

// Helper function to get pattern hint
func getPatternHint(patterns []string) string {
	if len(patterns) == 0 {
		return "general problem-solving techniques"
	}
	pattern := patterns[0]
	if scale, ok := musicalScales[pattern]; ok {
		return scale.Description
	}
	return pattern + " pattern"
}

func init() {
	// Add commands to root
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(hintCmd)
	rootCmd.AddCommand(solutionCmd)

	// Add flags for submit/test commands
	submitCmd.Flags().String("problem-id", "", "Problem ID")
	submitCmd.Flags().String("language", "go", "Programming language")
	submitCmd.Flags().String("file", "", "Solution file path")
	submitCmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
	submitCmd.MarkFlagRequired("problem-id")
	submitCmd.MarkFlagRequired("file")

	testCmd.Flags().String("problem-id", "", "Problem ID")
	testCmd.Flags().String("language", "go", "Programming language")
	testCmd.Flags().String("file", "", "Solution file path")
	testCmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
	testCmd.MarkFlagRequired("problem-id")
	testCmd.MarkFlagRequired("file")

	// Add flags for hint command
	hintCmd.Flags().String("problem-id", "", "Problem ID")
	hintCmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
	hintCmd.MarkFlagRequired("problem-id")

	// Add flags for solution command
	solutionCmd.Flags().String("problem-id", "", "Problem ID")
	solutionCmd.Flags().String("language", "go", "Programming language")
	solutionCmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
	solutionCmd.MarkFlagRequired("problem-id")
}