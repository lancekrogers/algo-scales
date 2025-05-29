package ai

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/claude-code-go/pkg/claude"
)

// REPL provides an interactive chat interface for AI assistance
type REPL struct {
	agent        Agent
	claudeClient *claude.ClaudeClient // Direct client for Claude provider
	sessionID    string               // Track conversation session
	context      []Message
	style        REPLStyle
	usingClaude  bool
	problem      *problem.Problem // Current problem context
}

// REPLStyle defines the visual styling for the REPL
type REPLStyle struct {
	User      lipgloss.Style
	Assistant lipgloss.Style
	System    lipgloss.Style
	Error     lipgloss.Style
	Tool      lipgloss.Style
	Cost      lipgloss.Style
}

// NewREPL creates a new REPL instance
func NewREPL(agent Agent) *REPL {
	repl := &REPL{
		agent:   agent,
		context: []Message{},
		style: REPLStyle{
			User:      lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true),
			Assistant: lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
			System:    lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Italic(true),
			Error:     lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true),
			Tool:      lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Italic(true),
			Cost:      lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		},
	}

	// Check if we're using Claude provider to leverage streaming
	if provider, ok := agent.(*ClaudeProvider); ok {
		repl.usingClaude = true
		repl.claudeClient = provider.client
		repl.sessionID = provider.sessionID
	}

	return repl
}

// Start begins an interactive chat session
func (r *REPL) Start(ctx context.Context, prob *problem.Problem) error {
	r.problem = prob

	// Build system context
	systemPrompt := r.buildSystemPrompt(prob)

	// Set up signal handling for graceful exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Handle signals in a goroutine
	go func() {
		<-sigChan
		fmt.Println(r.style.System.Render("\n\nInterrupted. Goodbye! Keep practicing! 👋"))
		os.Exit(0)
	}()

	fmt.Println(r.style.System.Render("🤖 AI Assistant Ready! Type 'help' for commands or 'exit' to quit."))
	fmt.Println(r.style.System.Render("   (Press Enter on empty line or use :q to exit)"))
	if prob != nil {
		pattern := ""
		if len(prob.Patterns) > 0 {
			pattern = prob.Patterns[0]
		}
		fmt.Println(r.style.System.Render(fmt.Sprintf("Problem: %s (%s pattern)", prob.Title, pattern)))
	}
	fmt.Println()

	// Exit commands based on claude-code-go demo
	exitCommands := []string{
		"exit", "quit", "bye", "goodbye", "q", ":q", ":quit", ":exit",
		"done", "finish", "end", "stop", "close", "leave",
		"/exit", "/quit", "\\q", "\\quit", // Common variations
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(r.style.User.Render("You> "))

		// Read input
		if !scanner.Scan() {
			// EOF or error - exit gracefully
			fmt.Println(r.style.System.Render("\nGoodbye! Keep practicing! 👋"))
			return nil
		}

		input := strings.TrimSpace(scanner.Text())
		
		// Empty input exits
		if input == "" {
			fmt.Println(r.style.System.Render("Goodbye! Keep practicing! 👋"))
			return nil
		}

		// Check for exit commands
		lowInput := strings.ToLower(input)
		isExit := false
		for _, cmd := range exitCommands {
			if lowInput == cmd {
				isExit = true
				break
			}
		}
		if isExit {
			fmt.Println(r.style.System.Render("Goodbye! Keep practicing! 👋"))
			return nil
		}

		// Handle other commands
		switch lowInput {
		case "help", "h", "?":
			r.showHelp()
			continue
		case "clear", "reset":
			r.context = []Message{}
			r.sessionID = ""
			fmt.Println(r.style.System.Render("Conversation cleared."))
			continue
		case "code":
			if prob != nil {
				r.showStarterCode(prob)
			} else {
				fmt.Println(r.style.Error.Render("No problem context available."))
			}
			continue
		case "hint", "hint 1":
			if prob != nil {
				r.getHint(ctx, prob, 1)
			} else {
				fmt.Println(r.style.Error.Render("No problem context available."))
			}
			continue
		case "hint 2":
			if prob != nil {
				r.getHint(ctx, prob, 2)
			} else {
				fmt.Println(r.style.Error.Render("No problem context available."))
			}
			continue
		case "hint 3":
			if prob != nil {
				r.getHint(ctx, prob, 3)
			} else {
				fmt.Println(r.style.Error.Render("No problem context available."))
			}
			continue
		case "pattern":
			if prob != nil {
				pattern := ""
				if len(prob.Patterns) > 0 {
					pattern = prob.Patterns[0]
				}
				r.explainPattern(ctx, pattern)
			} else {
				fmt.Println(r.style.Error.Render("No problem context available."))
			}
			continue
		case "":
			continue
		}

		// Regular chat message
		fmt.Print(r.style.Assistant.Render("Assistant> "))

		// Add to context
		userMsg := Message{Role: "user", Content: input}
		r.context = append(r.context, userMsg)

		// Prepare messages with system prompt
		messages := append([]Message{{Role: "system", Content: systemPrompt}}, r.context...)

		// Get response
		respChan, err := r.agent.Chat(ctx, messages, ChatOptions{
			Temperature: 0.7,
			MaxTokens:   2048,
			Stream:      true,
		})

		if err != nil {
			fmt.Println(r.style.Error.Render(fmt.Sprintf("\nError: %v", err)))
			continue
		}

		// Process streaming response
		var fullResponse strings.Builder
		for resp := range respChan {
			if resp.Error != nil {
				fmt.Println(r.style.Error.Render(fmt.Sprintf("\nError: %v", resp.Error)))
				break
			}

			// Handle special content (tool usage, etc.)
			if strings.HasPrefix(resp.Content, "[Using tool:") {
				fmt.Println()
				fmt.Println(r.style.Tool.Render(resp.Content))
				continue
			}

			// Stream the response
			fmt.Print(resp.Content)
			fullResponse.WriteString(resp.Content)

			// Handle completion
			if resp.Done {
				fmt.Println()
				if resp.SessionID != "" {
					r.sessionID = resp.SessionID
				}
				if resp.Cost > 0 {
					fmt.Println(r.style.Cost.Render(fmt.Sprintf("💰 Cost: $%.4f", resp.Cost)))
				}
			}
		}
		fmt.Println()

		// Save assistant response to context
		if fullResponse.Len() > 0 {
			r.context = append(r.context, Message{
				Role:    "assistant",
				Content: fullResponse.String(),
			})
		}
	}

	return nil
}

// Helper methods

func (r *REPL) buildSystemPrompt(prob *problem.Problem) string {
	if prob == nil {
		return "You are an expert algorithm tutor helping students learn data structures and algorithms. Focus on teaching concepts and patterns rather than just providing solutions."
	}

	pattern := "unknown"
	if len(prob.Patterns) > 0 {
		pattern = prob.Patterns[0]
	}
	
	return fmt.Sprintf(
		`You are helping with the algorithm problem "%s" which uses the %s pattern.
        
Problem Description: %s

Guide the student through understanding and solving this problem. 
Focus on teaching the pattern and approach rather than just giving the answer.
If they share code, help them debug and improve it.
Be encouraging and patient.`,
		prob.Title, pattern, prob.Description,
	)
}

func (r *REPL) showHelp() {
	help := `
Available Commands:
  help      - Show this help message
  clear     - Clear conversation history
  code      - Show starter code for current problem
  hint      - Get a level 1 hint (general approach)
  hint 2    - Get a level 2 hint (specific guidance)
  hint 3    - Get a level 3 hint (detailed pseudocode)
  pattern   - Explain the algorithm pattern
  
Exit Commands:
  exit, quit, q, :q, bye, done, stop
  Or press Enter on empty line
  Or press Ctrl+C

Tips:
  - Share your code for specific feedback
  - Ask about edge cases and optimizations
  - Request step-by-step walkthroughs
  - Discuss time/space complexity
`
	fmt.Println(r.style.System.Render(help))
}

func (r *REPL) showStarterCode(prob *problem.Problem) {
	fmt.Println(r.style.System.Render("\nStarter Code:"))
	for lang, code := range prob.StarterCode {
		fmt.Printf("\n=== %s ===\n%s\n", lang, code)
	}
}

func (r *REPL) getHint(ctx context.Context, prob *problem.Problem, level int) {
	fmt.Println(r.style.System.Render(fmt.Sprintf("\nGenerating level %d hint...", level)))

	hintChan, err := r.agent.GetHint(ctx, *prob, "", level)
	if err != nil {
		fmt.Println(r.style.Error.Render(fmt.Sprintf("Error: %v", err)))
		return
	}

	fmt.Print(r.style.Assistant.Render("Hint: "))
	for hint := range hintChan {
		fmt.Println(hint)
	}
}

func (r *REPL) explainPattern(ctx context.Context, pattern string) {
	fmt.Println(r.style.System.Render(fmt.Sprintf("\nExplaining %s pattern...", pattern)))

	examples := []problem.Problem{}
	if r.problem != nil {
		examples = append(examples, *r.problem)
	}

	explainChan, err := r.agent.ExplainPattern(ctx, pattern, examples)
	if err != nil {
		fmt.Println(r.style.Error.Render(fmt.Sprintf("Error: %v", err)))
		return
	}

	for explanation := range explainChan {
		fmt.Println(r.style.Assistant.Render(explanation))
	}
}