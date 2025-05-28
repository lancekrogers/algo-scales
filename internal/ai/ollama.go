package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

// OllamaProvider implements the Agent interface using Ollama's local API
type OllamaProvider struct {
	config     OllamaConfig
	client     *http.Client
	apiBaseURL string
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(config OllamaConfig) (*OllamaProvider, error) {
	// Set default values
	if config.Host == "" {
		config.Host = "http://localhost:11434"
	}
	if config.Model == "" {
		config.Model = "llama3.2:latest"
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.NumCtx == 0 {
		config.NumCtx = 4096
	}
	if config.Timeout == 0 {
		config.Timeout = 300
	}

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	// Ensure host doesn't have trailing slash
	apiBaseURL := strings.TrimRight(config.Host, "/")

	return &OllamaProvider{
		config:     config,
		client:     client,
		apiBaseURL: apiBaseURL,
	}, nil
}

// Chat implements the Agent interface
func (o *OllamaProvider) Chat(ctx context.Context, messages []Message, opts ChatOptions) (<-chan ChatResponse, error) {
	respChan := make(chan ChatResponse)

	go func() {
		defer close(respChan)

		// Convert messages to Ollama format
		ollamaMessages := make([]ollamaMessage, len(messages))
		for i, msg := range messages {
			ollamaMessages[i] = ollamaMessage{
				Role:    msg.Role,
				Content: msg.Content,
			}
		}

		// Prepare request
		reqBody := ollamaChatRequest{
			Model:    o.config.Model,
			Messages: ollamaMessages,
			Options: map[string]interface{}{
				"temperature": opts.Temperature,
				"num_ctx":     o.config.NumCtx,
			},
			Stream: true, // Always stream for real-time responses
		}

		reqData, err := json.Marshal(reqBody)
		if err != nil {
			respChan <- ChatResponse{Error: fmt.Errorf("failed to marshal request: %w", err)}
			return
		}

		// Create HTTP request
		req, err := http.NewRequestWithContext(ctx, "POST", o.apiBaseURL+"/api/chat", bytes.NewReader(reqData))
		if err != nil {
			respChan <- ChatResponse{Error: fmt.Errorf("failed to create request: %w", err)}
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := o.client.Do(req)
		if err != nil {
			respChan <- ChatResponse{Error: fmt.Errorf("request failed: %w", err)}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			respChan <- ChatResponse{Error: fmt.Errorf("ollama API error: %s (status %d)", string(body), resp.StatusCode)}
			return
		}

		// Read streaming response
		decoder := json.NewDecoder(resp.Body)
		for {
			var streamResp ollamaChatResponse
			if err := decoder.Decode(&streamResp); err != nil {
				if err == io.EOF {
					break
				}
				respChan <- ChatResponse{Error: fmt.Errorf("failed to decode response: %w", err)}
				return
			}

			respChan <- ChatResponse{
				Content: streamResp.Message.Content,
				Done:    streamResp.Done,
			}
		}
	}()

	return respChan, nil
}

// GetHint implements progressive hint generation
func (o *OllamaProvider) GetHint(ctx context.Context, prob problem.Problem, userCode string, level int) (<-chan string, error) {
	hintChan := make(chan string)

	go func() {
		defer close(hintChan)

		// Build hint prompt
		systemPrompt := o.buildHintSystemPrompt(prob, level)
		userPrompt := o.buildHintUserPrompt(prob, userCode, level)

		messages := []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		}

		// Get response from Ollama
		respChan, err := o.Chat(ctx, messages, ChatOptions{
			Temperature: 0.7,
			Stream:      true,
		})
		if err != nil {
			hintChan <- fmt.Sprintf("Error generating hint: %v", err)
			return
		}

		// Collect full response
		var fullHint strings.Builder
		for resp := range respChan {
			if resp.Error != nil {
				hintChan <- fmt.Sprintf("Error: %v", resp.Error)
				return
			}
			fullHint.WriteString(resp.Content)
		}

		hintChan <- fullHint.String()
	}()

	return hintChan, nil
}

// ReviewCode provides AI-powered code review
func (o *OllamaProvider) ReviewCode(ctx context.Context, prob problem.Problem, code string) (<-chan string, error) {
	reviewChan := make(chan string)

	go func() {
		defer close(reviewChan)

		// Build review prompt
		systemPrompt := "You are a senior software engineer conducting a thorough code review. Focus on educational feedback that helps the student improve."
		userPrompt := fmt.Sprintf("Review this code for the problem \"%s\":\n\n" +
			"Problem details:\n" +
			"- Pattern: %s\n" +
			"- Difficulty: %s\n" +
			"- Description: %s\n\n" +
			"Code to review:\n" +
			"```%s\n" +
			"%s\n" +
			"```\n\n" +
			"Please provide feedback on:\n" +
			"1. Correctness for the given problem\n" +
			"2. Code quality and style\n" +
			"3. Performance considerations\n" +
			"4. Edge case handling\n" +
			"5. Suggestions for improvement",
			prob.Title, getPrimaryPattern(prob), prob.Difficulty, prob.Description, "go", code)

		messages := []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		}

		// Get response from Ollama
		respChan, err := o.Chat(ctx, messages, ChatOptions{
			Temperature: 0.7,
			Stream:      true,
		})
		if err != nil {
			reviewChan <- fmt.Sprintf("Error generating review: %v", err)
			return
		}

		// Stream the review
		for resp := range respChan {
			if resp.Error != nil {
				reviewChan <- fmt.Sprintf("Error: %v", resp.Error)
				return
			}
			reviewChan <- resp.Content
		}
	}()

	return reviewChan, nil
}

// ExplainPattern provides detailed pattern explanations
func (o *OllamaProvider) ExplainPattern(ctx context.Context, pattern string, examples []problem.Problem) (<-chan string, error) {
	explainChan := make(chan string)

	go func() {
		defer close(explainChan)

		// Build examples context
		examplesText := ""
		for i, ex := range examples {
			if i >= 3 {
				break
			} // Limit to 3 examples
			examplesText += fmt.Sprintf("\n- %s (Difficulty: %s)", ex.Title, ex.Difficulty)
		}

		systemPrompt := "You are an algorithm expert teaching computer science students. Make your explanations clear, practical, and interview-focused."
		userPrompt := fmt.Sprintf(`Explain the "%s" algorithm pattern in detail.

Include:
1. When to use this pattern
2. Key characteristics
3. Common implementation approaches
4. Time and space complexity
5. Common pitfalls
6. Tips for recognition in interviews

Related example problems:%s`, pattern, examplesText)

		messages := []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		}

		// Get response from Ollama
		respChan, err := o.Chat(ctx, messages, ChatOptions{
			Temperature: 0.7,
			Stream:      true,
		})
		if err != nil {
			explainChan <- fmt.Sprintf("Error generating explanation: %v", err)
			return
		}

		// Stream the explanation
		for resp := range respChan {
			if resp.Error != nil {
				explainChan <- fmt.Sprintf("Error: %v", resp.Error)
				return
			}
			explainChan <- resp.Content
		}
	}()

	return explainChan, nil
}

// Helper methods

func (o *OllamaProvider) buildHintSystemPrompt(prob problem.Problem, level int) string {
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

func (o *OllamaProvider) buildHintUserPrompt(prob problem.Problem, userCode string, level int) string {
	if userCode != "" {
		return fmt.Sprintf("I'm working on this problem and here's my current code:\n```\n%s\n```\n\nI need a level %d hint.",
			userCode, level)
	}
	return fmt.Sprintf("I need a level %d hint for this problem.", level)
}

// Ollama API types

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string                 `json:"model"`
	Messages []ollamaMessage        `json:"messages"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Stream   bool                   `json:"stream"`
}

type ollamaChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}