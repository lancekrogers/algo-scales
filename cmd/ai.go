package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/ai"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI assistant configuration and management",
	Long:  `Configure and interact with AI assistants (Claude or Ollama) for algorithm learning support.`,
}

var aiConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure AI assistant settings",
	Long:  `Interactive configuration wizard for setting up AI providers.`,
	Run: func(cmd *cobra.Command, args []string) {
		configureAIAssistant()
	},
}

var aiConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current AI configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := ai.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println("Run 'algo-scales ai config' to create a configuration.")
			return
		}

		// Display configuration (hide sensitive data)
		displayConfig(config)
	},
}

var aiConfigSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long: `Set a specific configuration value.
Examples:
  algo-scales ai config set default_provider claude
  algo-scales ai config set claude.cli_path /usr/local/bin/claude
  algo-scales ai config set ollama.model llama3`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]

		if err := updateConfig(key, value); err != nil {
			fmt.Printf("Error updating config: %v\n", err)
			return
		}

		fmt.Printf("✅ Updated %s = %s\n", key, value)
	},
}

var aiTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test AI configuration",
	Long:  `Test your AI provider configuration to ensure it's working correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		testAIProvider(provider)
	},
}

func init() {
	// Add config subcommands
	aiConfigCmd.AddCommand(aiConfigShowCmd)
	aiConfigCmd.AddCommand(aiConfigSetCmd)
	
	// Add subcommands to ai command
	aiCmd.AddCommand(aiConfigCmd)
	aiCmd.AddCommand(aiTestCmd)

	// Add flags
	aiTestCmd.Flags().StringP("provider", "p", "", "AI provider to test (claude or ollama)")

	// Add ai command to root
	rootCmd.AddCommand(aiCmd)

	// Modify hint command to support AI
	enhanceHintCommand()

	// Add review command
	rootCmd.AddCommand(reviewCmd)
}

// reviewCmd provides AI-powered code review
var reviewCmd = &cobra.Command{
	Use:   "review [file]",
	Short: "Get AI-powered code review",
	Long:  `Submit your solution for AI-powered code review and feedback.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useAI, _ := cmd.Flags().GetBool("ai")
		if !useAI {
			fmt.Println("Code review requires AI. Use --ai flag to enable.")
			return
		}

		// For now, we'll require the user to specify a problem
		// In the future, we'll track the current session
		problemID, _ := cmd.Flags().GetString("problem")
		if problemID == "" {
			fmt.Println("Please specify a problem with --problem flag")
			fmt.Println("Example: algo-scales review --ai --problem two_sum")
			return
		}

		// Get problem
		prob, err := problem.GetByID(problemID)
		if err != nil {
			fmt.Printf("Error loading problem: %v\n", err)
			return
		}

		// Get code to review
		var code string
		if len(args) > 0 {
			// Read from file
			content, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}
			code = string(content)
		} else {
			fmt.Println("Please specify a file to review")
			fmt.Println("Example: algo-scales review --ai --problem two_sum solution.go")
			return
		}

		// Perform AI review (default to Go for now)
		reviewCode(prob, code, "go")
	},
}

func enhanceHintCommand() {
	// Find the hint command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "hint" {
			// Add AI flags
			cmd.Flags().Bool("ai", false, "Use AI assistant for hints")
			cmd.Flags().Bool("interactive", false, "Start interactive AI chat")
			cmd.Flags().StringP("problem", "p", "", "Problem ID for AI hints")

			// Override the run function
			originalRun := cmd.Run
			cmd.Run = func(cmd *cobra.Command, args []string) {
				useAI, _ := cmd.Flags().GetBool("ai")
				interactive, _ := cmd.Flags().GetBool("interactive")

				if useAI {
					// For now, require problem ID
					problemID, _ := cmd.Flags().GetString("problem")
					if problemID == "" {
						fmt.Println("Please specify a problem with --problem flag")
						fmt.Println("Example: algo-scales hint --ai --problem two_sum")
						return
					}

					// Get problem
					prob, err := problem.GetByID(problemID)
					if err != nil {
						fmt.Printf("Error loading problem: %v\n", err)
						return
					}

					if interactive {
						// Start interactive REPL
						startAIChat(prob)
					} else {
						// Get single AI hint (start at level 1)
						getAIHint(prob, "", 1)
					}
				} else {
					// Use traditional hint system
					originalRun(cmd, args)
				}
			}
			break
		}
	}
}

// Helper functions

func configureAIAssistant() {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	fmt.Println(style.Render("🤖 AI Assistant Configuration"))
	fmt.Println(strings.Repeat("─", 50))

	// Load existing config or create new
	config, _ := ai.LoadConfig()
	if config == nil {
		config = &ai.Config{
			Version: "1.0",
		}
	}

	// Select provider
	fmt.Println("\nSelect your AI provider:")
	fmt.Println("1. Claude (via Claude Code CLI)")
	fmt.Println("2. Ollama (local AI)")
	fmt.Print("\nChoice (1-2): ")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		config.DefaultProvider = "claude"
		configureClaude(config)
	case "2":
		config.DefaultProvider = "ollama"
		configureOllama(config)
	default:
		fmt.Println(errorStyle.Render("Invalid choice"))
		return
	}

	// Save configuration
	if err := ai.SaveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	fmt.Println(style.Render("\n✅ Configuration saved successfully!"))
	fmt.Println("\nNext steps:")
	fmt.Println("- Test your configuration: algo-scales ai test")
	fmt.Println("- Get AI hints: algo-scales hint --ai")
	fmt.Println("- Start AI chat: algo-scales hint --ai --interactive")
}

func configureClaude(config *ai.Config) {
	if config.Claude == nil {
		config.Claude = &ai.ClaudeConfig{
			CLIPath:       "claude",
			DefaultFormat: "json",
			SaveSessions:  true,
			SessionDir:    "~/.algo-scales/claude-sessions",
			MaxTurns:      5,
		}
	}

	fmt.Println("\nClaude Configuration")
	fmt.Println("Note: Claude Code CLI must be installed. See: https://docs.anthropic.com/en/docs/claude-code/getting-started")

	fmt.Printf("\nClaude CLI path [%s]: ", config.Claude.CLIPath)
	var cliPath string
	fmt.Scanln(&cliPath)
	if cliPath != "" {
		config.Claude.CLIPath = cliPath
	}

	fmt.Println("\nEnable MCP (Model Context Protocol) for enhanced capabilities? (y/n) [y]: ")
	var enableMCP string
	fmt.Scanln(&enableMCP)
	if enableMCP != "n" {
		config.Claude.MCP = &ai.MCPConfig{
			Enabled: true,
			Servers: map[string]ai.MCPServerConfig{
				"filesystem": {
					Command: "npx",
					Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "./"},
				},
			},
		}
		config.Claude.AllowedTools = []string{
			"mcp__filesystem__read_file",
			"mcp__filesystem__list_directory",
			"Read",
		}
	}
}

func configureOllama(config *ai.Config) {
	if config.Ollama == nil {
		config.Ollama = &ai.OllamaConfig{
			Host:        "http://localhost:11434",
			Model:       "llama3.2:latest",
			Temperature: 0.7,
			NumCtx:      4096,
			Timeout:     300,
		}
	}

	fmt.Println("\nOllama Configuration")
	fmt.Println("Note: Ollama must be installed and running. See: https://ollama.com")

	fmt.Printf("\nOllama host [%s]: ", config.Ollama.Host)
	var host string
	fmt.Scanln(&host)
	if host != "" {
		config.Ollama.Host = host
	}

	fmt.Printf("Model name [%s]: ", config.Ollama.Model)
	var model string
	fmt.Scanln(&model)
	if model != "" {
		config.Ollama.Model = model
	}
}

func displayConfig(config *ai.Config) {
	// Convert to YAML for pretty display
	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Printf("Error displaying config: %v\n", err)
		return
	}

	fmt.Println("Current AI Configuration:")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(string(data))
}

func updateConfig(key, value string) error {
	config, err := ai.LoadConfig()
	if err != nil {
		return err
	}

	// Parse key path (e.g., "claude.model" -> ["claude", "model"])
	parts := strings.Split(key, ".")

	switch parts[0] {
	case "default_provider":
		config.DefaultProvider = value
	case "claude":
		if config.Claude == nil {
			config.Claude = &ai.ClaudeConfig{}
		}
		if len(parts) > 1 {
			switch parts[1] {
			case "cli_path":
				config.Claude.CLIPath = value
			case "default_format":
				config.Claude.DefaultFormat = value
			case "save_sessions":
				config.Claude.SaveSessions = value == "true"
			case "session_directory":
				config.Claude.SessionDir = value
			case "max_turns":
				fmt.Sscanf(value, "%d", &config.Claude.MaxTurns)
			default:
				return fmt.Errorf("unknown claude setting: %s", parts[1])
			}
		}
	case "ollama":
		if config.Ollama == nil {
			config.Ollama = &ai.OllamaConfig{}
		}
		if len(parts) > 1 {
			switch parts[1] {
			case "host":
				config.Ollama.Host = value
			case "model":
				config.Ollama.Model = value
			case "temperature":
				fmt.Sscanf(value, "%f", &config.Ollama.Temperature)
			case "num_ctx":
				fmt.Sscanf(value, "%d", &config.Ollama.NumCtx)
			case "timeout":
				fmt.Sscanf(value, "%d", &config.Ollama.Timeout)
			default:
				return fmt.Errorf("unknown ollama setting: %s", parts[1])
			}
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return ai.SaveConfig(config)
}

func testAIProvider(provider string) {
	// Load config
	config, err := ai.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Use default provider if not specified
	if provider == "" {
		provider = config.DefaultProvider
	}

	fmt.Printf("Testing %s provider...\n", provider)

	// Create agent
	agent, err := ai.NewAgent(ai.Provider(provider), config)
	if err != nil {
		fmt.Printf("❌ Failed to create agent: %v\n", err)
		fmt.Println(ai.HandleProviderError(ai.Provider(provider), err))
		return
	}

	// Test with a simple prompt
	ctx := context.Background()
	messages := []ai.Message{
		{Role: "user", Content: "Hello! Please respond with 'AI assistant is working!' if you can see this message."},
	}

	respChan, err := agent.Chat(ctx, messages, ai.ChatOptions{
		Temperature: 0.1,
		Stream:      true,
	})
	if err != nil {
		fmt.Printf("❌ Failed to send message: %v\n", err)
		return
	}

	// Collect response
	fmt.Print("Response: ")
	gotResponse := false
	for resp := range respChan {
		if resp.Error != nil {
			fmt.Printf("\n❌ Error: %v\n", resp.Error)
			return
		}
		fmt.Print(resp.Content)
		gotResponse = true
	}
	fmt.Println()

	if gotResponse {
		fmt.Printf("\n✅ %s provider is working correctly!\n", provider)
	} else {
		fmt.Println("\n❌ No response received")
	}
}

func getAIHint(prob *problem.Problem, userCode string, level int) {
	agent, err := ai.GetDefaultAgent()
	if err != nil {
		fmt.Printf("Error initializing AI: %v\n", err)
		fmt.Println("Run 'algo-scales ai config' to set up AI assistant.")
		return
	}

	fmt.Printf("🤔 Generating level %d hint...\n", level)

	ctx := context.Background()
	hintChan, err := agent.GetHint(ctx, *prob, userCode, level)
	if err != nil {
		fmt.Printf("Error getting hint: %v\n", err)
		return
	}

	formatter := ai.NewResponseFormatter()
	for hint := range hintChan {
		fmt.Println(formatter.FormatHint(level, hint))
	}
}

func startAIChat(prob *problem.Problem) {
	agent, err := ai.GetDefaultAgent()
	if err != nil {
		fmt.Printf("Error initializing AI: %v\n", err)
		fmt.Println("Run 'algo-scales ai config' to set up AI assistant.")
		return
	}

	repl := ai.NewREPL(agent)
	ctx := context.Background()
	if err := repl.Start(ctx, prob); err != nil {
		fmt.Printf("Error in AI chat: %v\n", err)
	}
}

func reviewCode(prob *problem.Problem, code string, language string) {
	agent, err := ai.GetDefaultAgent()
	if err != nil {
		fmt.Printf("Error initializing AI: %v\n", err)
		fmt.Println("Run 'algo-scales ai config' to set up AI assistant.")
		return
	}

	fmt.Println("🔍 Reviewing your code...")

	ctx := context.Background()
	reviewChan, err := agent.ReviewCode(ctx, *prob, code)
	if err != nil {
		fmt.Printf("Error reviewing code: %v\n", err)
		return
	}

	formatter := ai.NewResponseFormatter()
	var fullReview strings.Builder
	for review := range reviewChan {
		fullReview.WriteString(review)
	}

	fmt.Println(formatter.FormatCodeReview(fullReview.String()))
}

func init() {
	// Add flags to review command
	reviewCmd.Flags().Bool("ai", true, "Use AI for code review")
	reviewCmd.Flags().StringP("problem", "p", "", "Problem ID to review against")
}