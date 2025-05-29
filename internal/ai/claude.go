package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/claude-code-go/pkg/claude"
)

// ClaudeProvider implements the Agent interface using claude-code-go
type ClaudeProvider struct {
	config    ClaudeConfig
	client    *claude.ClaudeClient
	sessionID string // Track current session for multi-turn conversations
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(config ClaudeConfig) (*ClaudeProvider, error) {
	// Initialize the claude-code-go client
	// The client wraps the Claude Code CLI, which must be installed
	cliPath := config.CLIPath
	if cliPath == "" {
		cliPath = "claude"
	}

	client := claude.NewClient(cliPath)

	// Set default options based on config
	format := claude.JSONOutput
	switch config.DefaultFormat {
	case "stream-json":
		format = claude.StreamJSONOutput
	case "text":
		format = claude.TextOutput
	}

	client.DefaultOptions = &claude.RunOptions{
		Format: format,
	}

	// Create session directory if it doesn't exist
	if config.SaveSessions && config.SessionDir != "" {
		if err := os.MkdirAll(config.SessionDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create session directory: %w", err)
		}
	}

	return &ClaudeProvider{
		config: config,
		client: client,
	}, nil
}

// Chat implements the Agent interface
func (c *ClaudeProvider) Chat(ctx context.Context, messages []Message, opts ChatOptions) (<-chan ChatResponse, error) {
	respChan := make(chan ChatResponse)

	go func() {
		defer close(respChan)

		// Build the prompt from messages
		prompt := c.buildPromptFromMessages(messages)

		// Configure run options
		runOpts := &claude.RunOptions{
			Format:   claude.StreamJSONOutput, // Stream for real-time responses
			MaxTurns: c.config.MaxTurns,
		}

		// Add system prompt if present
		if len(messages) > 0 && messages[0].Role == "system" {
			runOpts.SystemPrompt = messages[0].Content
		}

		// Use existing session if continuing conversation
		if c.sessionID != "" {
			runOpts.ResumeID = c.sessionID
		}

		// Set up MCP if configured
		if c.config.MCP != nil && c.config.MCP.Enabled {
			mcpFile, err := c.writeMCPConfig(c.config.MCP)
			if err != nil {
				respChan <- ChatResponse{Error: fmt.Errorf("failed to write MCP config: %w", err)}
				return
			}
			defer os.Remove(mcpFile)
			runOpts.MCPConfigPath = mcpFile
		}

		// Set allowed tools
		if len(c.config.AllowedTools) > 0 {
			runOpts.AllowedTools = c.config.AllowedTools
		}

		// Stream the response
		messageCh, errCh := c.client.StreamPrompt(ctx, prompt, runOpts)

		// Handle errors in separate goroutine
		go func() {
			for err := range errCh {
				respChan <- ChatResponse{Error: err}
			}
		}()

		// Process streaming messages
		var fullResponse strings.Builder
		for msg := range messageCh {
			switch msg.Type {
			case "assistant":
				// Assistant's response content
				respChan <- ChatResponse{
					Content: msg.Result,
					Done:    false,
				}
				fullResponse.WriteString(msg.Result)
			case "result":
				// Final result with metadata
				c.sessionID = msg.SessionID // Save for continuation
				respChan <- ChatResponse{
					Content:   msg.Result,
					Done:      true,
					SessionID: msg.SessionID,
					Cost:      msg.CostUSD,
				}
			case "tool_use":
				// Show tool usage for transparency
				respChan <- ChatResponse{
					Content: fmt.Sprintf("[Using tool: %s]", msg.Type),
					Done:    false,
				}
			}
		}

		// Save session if configured
		if c.config.SaveSessions && c.sessionID != "" && fullResponse.Len() > 0 {
			c.saveSession(messages, fullResponse.String())
		}
	}()

	return respChan, nil
}

// GetHint implements progressive hint generation
func (c *ClaudeProvider) GetHint(ctx context.Context, prob problem.Problem, userCode string, level int) (<-chan string, error) {
	hintChan := make(chan string)

	go func() {
		defer close(hintChan)

		// Build progressive hint prompt based on level
		systemPrompt := c.buildHintSystemPrompt(prob, level)
		userPrompt := c.buildHintUserPrompt(prob, userCode, level)

		// For hints, we want focused, non-streaming responses
		result, err := c.client.RunPrompt(userPrompt, &claude.RunOptions{
			SystemPrompt: systemPrompt,
			Format:       claude.JSONOutput,
			MaxTurns:     1, // Single response for hints
		})

		if err != nil {
			hintChan <- fmt.Sprintf("Error generating hint: %v", err)
			return
		}

		// Send the complete hint
		hintChan <- result.Result
	}()

	return hintChan, nil
}

// ReviewCode provides AI-powered code review
func (c *ClaudeProvider) ReviewCode(ctx context.Context, prob problem.Problem, code string) (<-chan string, error) {
	reviewChan := make(chan string)

	go func() {
		defer close(reviewChan)

		// Create temporary file with the code
		tmpFile, err := os.CreateTemp("", fmt.Sprintf("review-%s-*%s", prob.ID, getFileExtension("go")))
		if err != nil {
			reviewChan <- fmt.Sprintf("Error creating temp file: %v", err)
			return
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.WriteString(code); err != nil {
			reviewChan <- fmt.Sprintf("Error writing code: %v", err)
			return
		}
		tmpFile.Close()

		// Use MCP filesystem tool to analyze the code
		mcpConfig := c.createCodeReviewMCPConfig()
		mcpFile, err := c.writeMCPConfig(mcpConfig)
		if err != nil {
			reviewChan <- fmt.Sprintf("Error creating MCP config: %v", err)
			return
		}
		defer os.Remove(mcpFile)

		// Review prompt
		prompt := fmt.Sprintf(`Review the code in %s for the problem "%s". 
Focus on:
1. Correctness for the given problem
2. Code quality and style
3. Performance considerations
4. Edge case handling
5. Suggestions for improvement

Problem details:
- Pattern: %s
- Difficulty: %s
- Description: %s`,
			tmpFile.Name(), prob.Title, getPrimaryPattern(prob), prob.Difficulty, prob.Description)

		// Stream the review
		messageCh, errCh := c.client.StreamPrompt(ctx, prompt, &claude.RunOptions{
			SystemPrompt:  "You are a senior software engineer conducting a thorough code review. Focus on educational feedback that helps the student improve.",
			MCPConfigPath: mcpFile,
			AllowedTools:  []string{"mcp__filesystem__read_file"},
			Format:        claude.StreamJSONOutput,
		})

		// Handle errors
		go func() {
			for err := range errCh {
				reviewChan <- fmt.Sprintf("Review error: %v", err)
			}
		}()

		// Stream review feedback
		for msg := range messageCh {
			if msg.Type == "assistant" {
				reviewChan <- msg.Result
			}
		}
	}()

	return reviewChan, nil
}

// ExplainPattern provides detailed pattern explanations
func (c *ClaudeProvider) ExplainPattern(ctx context.Context, pattern string, examples []problem.Problem) (<-chan string, error) {
	explainChan := make(chan string)

	go func() {
		defer close(explainChan)

		// Build examples context
		examplesText := ""
		for i, ex := range examples {
			if i >= 3 {
				break
			} // Limit to 3 examples
			examplesText += fmt.Sprintf("\nExample %d: %s (Difficulty: %s)\n", i+1, ex.Title, ex.Difficulty)
		}

		prompt := fmt.Sprintf(`Explain the "%s" algorithm pattern in detail.
Include:
1. When to use this pattern
2. Key characteristics
3. Common implementation approaches
4. Time and space complexity
5. Common pitfalls
6. Tips for recognition in interviews

Related problems:%s`, pattern, examplesText)

		result, err := c.client.RunPrompt(prompt, &claude.RunOptions{
			SystemPrompt: "You are an algorithm expert teaching computer science students. Make your explanations clear, practical, and interview-focused.",
			Format:       claude.JSONOutput,
		})

		if err != nil {
			explainChan <- fmt.Sprintf("Error: %v", err)
			return
		}

		explainChan <- result.Result
	}()

	return explainChan, nil
}

// Helper methods

func (c *ClaudeProvider) buildPromptFromMessages(messages []Message) string {
	var parts []string
	for _, msg := range messages {
		if msg.Role != "system" { // System prompt handled separately
			parts = append(parts, fmt.Sprintf("%s: %s", strings.Title(msg.Role), msg.Content))
		}
	}
	return strings.Join(parts, "\n\n")
}

func (c *ClaudeProvider) buildHintSystemPrompt(prob problem.Problem, level int) string {
	base := fmt.Sprintf(`You are a patient algorithm tutor helping a student with the "%s" problem.
Pattern: %s
Difficulty: %s

Your goal is to guide the student to discover the solution themselves.`,
		prob.Title, getPrimaryPattern(prob), prob.Difficulty)

	switch level {
	case 1:
		return base + "\nProvide a gentle hint about the general approach without revealing specifics. Focus on helping them recognize the pattern."
	case 2:
		return base + "\nProvide more specific guidance about the algorithm and data structures to use. You can mention specific techniques but don't give away the implementation."
	case 3:
		return base + "\nProvide detailed pseudocode or step-by-step implementation guidance. Help them understand exactly how to implement the solution."
	default:
		return base
	}
}

func (c *ClaudeProvider) buildHintUserPrompt(prob problem.Problem, userCode string, level int) string {
	if userCode != "" {
		return fmt.Sprintf("I'm working on this problem and here's my current code:\n```\n%s\n```\n\nI need a level %d hint.",
			userCode, level)
	}
	return fmt.Sprintf("I need a level %d hint for this problem.", level)
}

func (c *ClaudeProvider) createCodeReviewMCPConfig() *MCPConfig {
	return &MCPConfig{
		Enabled: true,
		Servers: map[string]MCPServerConfig{
			"filesystem": {
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "./"},
			},
		},
	}
}

func (c *ClaudeProvider) writeMCPConfig(config *MCPConfig) (string, error) {
	mcpFile, err := os.CreateTemp("", "mcp-config-*.json")
	if err != nil {
		return "", err
	}

	// Convert to format expected by claude-code-go
	mcpData := map[string]interface{}{
		"mcpServers": map[string]interface{}{},
	}

	for name, server := range config.Servers {
		serverConfig := map[string]interface{}{
			"command": server.Command,
			"args":    server.Args,
		}
		if len(server.Env) > 0 {
			serverConfig["env"] = server.Env
		}
		mcpData["mcpServers"].(map[string]interface{})[name] = serverConfig
	}

	encoder := json.NewEncoder(mcpFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(mcpData); err != nil {
		mcpFile.Close()
		os.Remove(mcpFile.Name())
		return "", err
	}

	mcpFile.Close()
	return mcpFile.Name(), nil
}

func (c *ClaudeProvider) saveSession(messages []Message, response string) {
	if c.config.SessionDir == "" {
		return
	}

	session := map[string]interface{}{
		"session_id": c.sessionID,
		"timestamp":  time.Now().Format(time.RFC3339),
		"messages":   messages,
		"response":   response,
	}

	filename := filepath.Join(c.config.SessionDir, fmt.Sprintf("session-%s.json", c.sessionID))
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(filename, data, 0600)
}

func getFileExtension(language string) string {
	switch strings.ToLower(language) {
	case "python":
		return ".py"
	case "javascript":
		return ".js"
	case "go":
		return ".go"
	case "java":
		return ".java"
	default:
		return ".txt"
	}
}

// Helper to get the primary pattern from a problem
func getPrimaryPattern(prob problem.Problem) string {
	if len(prob.Patterns) > 0 {
		return prob.Patterns[0]
	}
	return "unknown"
}