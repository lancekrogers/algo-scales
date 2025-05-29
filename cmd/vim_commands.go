// Vim mode commands for Neovim plugin integration

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/lancekrogers/algo-scales/internal/ai"
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

		// Get current hint level for this problem
		currentLevel := hintLevels[problemID]
		currentLevel++ // Increment for this request
		hintLevels[problemID] = currentLevel

		// Create response with appropriate level of detail
		resp := VimHintResponse{
			Level: currentLevel,
		}

		// Level 1: Pattern explanation
		if currentLevel >= 1 {
			if prob.PatternExplanation != "" {
				resp.Hint = prob.PatternExplanation
			} else {
				// Fallback to generic pattern hint
				resp.Hint = "Think about the pattern: " + getPatternHint(prob.Patterns)
			}
		}

		// Level 2: Add solution walkthrough
		if currentLevel >= 2 && len(prob.SolutionWalkthrough) > 0 {
			resp.Walkthrough = prob.SolutionWalkthrough
		}

		// Level 3: Add actual solution code
		if currentLevel >= 3 {
			// Get solution in the requested language
			if prob.Solutions != nil {
				if solution, ok := prob.Solutions[language]; ok {
					resp.Solution = solution
					resp.Language = language
				} else {
					// Try to get any solution
					for lang, sol := range prob.Solutions {
						resp.Solution = sol
						resp.Language = lang
						break
					}
				}
			}
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

// aiHintCmd represents the AI hint command for vim mode
var aiHintCmd = &cobra.Command{
	Use:   "ai-hint",
	Short: "Get AI-powered hints for problem (vim mode)",
	Long:  `Get AI-powered hints using claude-code-go or Ollama. Used by the Neovim plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		problemID, _ := cmd.Flags().GetString("problem-id")
		language, _ := cmd.Flags().GetString("language")
		userCode, _ := cmd.Flags().GetString("user-code")
		filePath, _ := cmd.Flags().GetString("file")
		provider, _ := cmd.Flags().GetString("provider") // "claude" or "ollama"
		model, _ := cmd.Flags().GetString("model")
		isVimMode, _ := cmd.Flags().GetBool("vim-mode")
		chatMode, _ := cmd.Flags().GetBool("chat")

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

		// Get user code from either flag or file
		if filePath != "" && userCode == "" {
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				outputVimError(fmt.Errorf("failed to read file: %v", err))
				return
			}
			userCode = string(content)
		}

		// Load AI configuration
		aiConfig, err := ai.LoadConfig()
		if err != nil {
			outputVimError(fmt.Errorf("AI not configured. Run 'algo-scales ai config' to set up: %v", err))
			return
		}

		// Determine provider (use flag if provided, otherwise default)
		var aiProvider ai.Provider
		if provider != "" {
			switch provider {
			case "claude":
				aiProvider = ai.ProviderClaude
			case "ollama":
				aiProvider = ai.ProviderOllama
			default:
				outputVimError(fmt.Errorf("unsupported provider: %s", provider))
				return
			}
		} else {
			switch aiConfig.DefaultProvider {
			case "claude":
				aiProvider = ai.ProviderClaude
			case "ollama":
				aiProvider = ai.ProviderOllama
			default:
				outputVimError(fmt.Errorf("no valid default provider configured"))
				return
			}
		}

		// Override model if specified via flags
		if model != "" {
			if aiProvider == ai.ProviderOllama && aiConfig.Ollama != nil {
				aiConfig.Ollama.Model = model
			}
			// Claude uses default model from CLI, no need to override
		}

		// Create AI agent
		agent, err := ai.NewAgent(aiProvider, aiConfig)
		if err != nil {
			outputVimError(fmt.Errorf("failed to create AI agent: %v", err))
			return
		}

		if chatMode {
			// Launch interactive chat mode
			resp := map[string]interface{}{
				"mode": "chat",
				"command": fmt.Sprintf("%s ai repl --problem-id %s --language %s --provider %s", 
					os.Args[0], problemID, language, string(aiProvider)),
				"provider": string(aiProvider),
				"model": model,
			}
			
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				outputVimError(fmt.Errorf("failed to marshal response: %v", err))
				return
			}
			
			fmt.Println(string(jsonResp))
			return
		}

		// Single hint mode
		hintStream, err := agent.GetHint(ctx, *prob, userCode, 1)
		if err != nil {
			outputVimError(fmt.Errorf("AI hint failed: %v", err))
			return
		}

		// Collect the streaming response
		var hintContent strings.Builder
		hasContent := false
		for chunk := range hintStream {
			if chunk != "" {
				hintContent.WriteString(chunk)
				hasContent = true
			}
		}
		
		// If no content received, provide a fallback message
		if !hasContent {
			hintContent.WriteString("AI hint service is available but no content was generated. Try using 'algo-scales ai repl' for interactive chat.")
		}

		// Create response with the AI hint
		resp := map[string]interface{}{
			"mode": "hint",
			"hint": hintContent.String(),
			"provider": string(aiProvider),
			"model": model,
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
	rootCmd.AddCommand(aiHintCmd)

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
	hintCmd.Flags().String("language", "go", "Programming language")
	hintCmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
	hintCmd.MarkFlagRequired("problem-id")

	// Add flags for solution command
	solutionCmd.Flags().String("problem-id", "", "Problem ID")
	solutionCmd.Flags().String("language", "go", "Programming language")
	solutionCmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
	solutionCmd.MarkFlagRequired("problem-id")

	// Add flags for AI hint command
	aiHintCmd.Flags().String("problem-id", "", "Problem ID")
	aiHintCmd.Flags().String("language", "go", "Programming language")
	aiHintCmd.Flags().String("user-code", "", "User's current solution code")
	aiHintCmd.Flags().String("file", "", "Path to solution file (alternative to --user-code)")
	aiHintCmd.Flags().String("provider", "claude", "AI provider (claude or ollama)")
	aiHintCmd.Flags().String("model", "", "AI model to use")
	aiHintCmd.Flags().Bool("vim-mode", false, "Enable vim mode output")
	aiHintCmd.Flags().Bool("chat", false, "Launch interactive chat mode")
	aiHintCmd.MarkFlagRequired("problem-id")
}