package ai

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/lancekrogers/algo-scales/internal/problem"
)

// PromptBuilder helps construct prompts for various AI tasks
type PromptBuilder struct {
	templates map[string]*template.Template
}

// NewPromptBuilder creates a new prompt builder with default templates
func NewPromptBuilder() *PromptBuilder {
	pb := &PromptBuilder{
		templates: make(map[string]*template.Template),
	}
	pb.loadDefaultTemplates()
	return pb
}

// loadDefaultTemplates loads the default prompt templates
func (pb *PromptBuilder) loadDefaultTemplates() {
	// Hint template
	hintTemplate := `Problem: {{.Problem.Title}}
Pattern: {{if .Problem.Patterns}}{{index .Problem.Patterns 0}}{{else}}unknown{{end}}
Difficulty: {{.Problem.Difficulty}}

{{if .UserCode}}
User's current approach:
` + "```\n{{.UserCode}}\n```" + `
{{end}}

Provide a helpful hint at level {{.Level}}:
- Level 1: General approach and pattern guidance
- Level 2: More specific algorithmic hints
- Level 3: Detailed pseudocode or implementation tips

IMPORTANT: Do not provide the complete code implementation. Guide the student to implement it themselves.`

	// Code review template
	reviewTemplate := `Review this {{.Language}} code for the problem "{{.Problem.Title}}":

Problem details:
- Pattern: {{if .Problem.Patterns}}{{index .Problem.Patterns 0}}{{else}}unknown{{end}}
- Difficulty: {{.Problem.Difficulty}}
- Key constraints: {{.Problem.Constraints}}

Code to review:
` + "```{{.Language}}\n{{.Code}}\n```" + `

Please provide feedback on:
1. Correctness for the given problem
2. Code quality and {{.Language}} best practices
3. Time and space complexity analysis
4. Edge case handling
5. Suggestions for improvement

Focus on educational feedback that helps the student learn.`

	// Pattern explanation template
	patternTemplate := `Explain the "{{.Pattern}}" algorithm pattern for technical interviews.

{{if .Examples}}
Related problems:
{{range .Examples}}- {{.Title}} ({{.Difficulty}})
{{end}}{{end}}

Cover these aspects:
1. Core concept and intuition
2. When to recognize and use this pattern
3. Common variations and implementations
4. Time and space complexity considerations
5. Common mistakes and how to avoid them
6. Interview tips for this pattern`

	// Solution walkthrough template
	walkthroughTemplate := `Walk through solving "{{.Problem.Title}}" step by step.

Problem summary: {{.Problem.Description}}

Provide:
1. Understanding the problem (inputs, outputs, constraints)
2. Recognizing the pattern (why {{if .Problem.Patterns}}{{index .Problem.Patterns 0}}{{else}}unknown{{end}} applies)
3. Developing the approach
4. Implementation considerations
5. Complexity analysis
6. Testing strategy`

	// Load templates
	pb.templates["hint"] = template.Must(template.New("hint").Parse(hintTemplate))
	pb.templates["review"] = template.Must(template.New("review").Parse(reviewTemplate))
	pb.templates["pattern"] = template.Must(template.New("pattern").Parse(patternTemplate))
	pb.templates["walkthrough"] = template.Must(template.New("walkthrough").Parse(walkthroughTemplate))
}

// BuildHintPrompt creates a hint prompt
func (pb *PromptBuilder) BuildHintPrompt(prob problem.Problem, userCode string, level int) (string, error) {
	data := map[string]interface{}{
		"Problem":  prob,
		"UserCode": userCode,
		"Level":    level,
	}
	return pb.executeTemplate("hint", data)
}

// BuildReviewPrompt creates a code review prompt
func (pb *PromptBuilder) BuildReviewPrompt(prob problem.Problem, code string, language string) (string, error) {
	data := map[string]interface{}{
		"Problem":  prob,
		"Code":     code,
		"Language": language,
	}
	return pb.executeTemplate("review", data)
}

// BuildPatternPrompt creates a pattern explanation prompt
func (pb *PromptBuilder) BuildPatternPrompt(pattern string, examples []problem.Problem) (string, error) {
	data := map[string]interface{}{
		"Pattern":  pattern,
		"Examples": examples,
	}
	return pb.executeTemplate("pattern", data)
}

// BuildWalkthroughPrompt creates a solution walkthrough prompt
func (pb *PromptBuilder) BuildWalkthroughPrompt(prob problem.Problem) (string, error) {
	data := map[string]interface{}{
		"Problem": prob,
	}
	return pb.executeTemplate("walkthrough", data)
}

// executeTemplate executes a template with the given data
func (pb *PromptBuilder) executeTemplate(name string, data interface{}) (string, error) {
	tmpl, ok := pb.templates[name]
	if !ok {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// SystemPrompts provides pre-defined system prompts for different contexts
type SystemPrompts struct{}

// NewSystemPrompts creates a new system prompts provider
func NewSystemPrompts() *SystemPrompts {
	return &SystemPrompts{}
}

// GetTutorPrompt returns the base tutor system prompt
func (sp *SystemPrompts) GetTutorPrompt() string {
	return `You are an expert algorithm tutor with years of experience teaching data structures and algorithms.
Your approach is patient, encouraging, and focused on helping students truly understand concepts rather than memorize solutions.
You excel at breaking down complex problems into manageable steps and relating abstract concepts to real-world examples.
Always prioritize learning over just getting the answer.`
}

// GetInterviewerPrompt returns the mock interviewer system prompt
func (sp *SystemPrompts) GetInterviewerPrompt() string {
	return `You are an experienced technical interviewer at a top tech company.
Your role is to assess the candidate's problem-solving skills, coding ability, and communication.
Ask clarifying questions, provide hints when the candidate is stuck, and evaluate their approach.
Be professional but friendly, and provide constructive feedback.`
}

// GetReviewerPrompt returns the code reviewer system prompt
func (sp *SystemPrompts) GetReviewerPrompt() string {
	return `You are a senior software engineer and code reviewer conducting educational code reviews.
Focus on helping students improve their coding skills through constructive feedback.
Point out both strengths and areas for improvement.
Explain the "why" behind your suggestions and provide examples when helpful.
Consider code correctness, efficiency, readability, and best practices.`
}

// GetDebuggerPrompt returns the debugging assistant system prompt
func (sp *SystemPrompts) GetDebuggerPrompt() string {
	return `You are a debugging expert and debugger helping students fix issues in their code.
Guide them through the debugging process step by step.
Help them understand error messages, identify logic errors, and test edge cases.
Teach debugging techniques they can use in the future.
Never just fix the code - help them understand and fix it themselves.`
}

// FormatResponse provides utilities for formatting AI responses
type ResponseFormatter struct{}

// NewResponseFormatter creates a new response formatter
func NewResponseFormatter() *ResponseFormatter {
	return &ResponseFormatter{}
}

// FormatHint formats a hint response with appropriate styling
func (rf *ResponseFormatter) FormatHint(level int, hint string) string {
	var prefix string
	switch level {
	case 1:
		prefix = "💡 General Approach"
	case 2:
		prefix = "🔍 Specific Guidance"
	case 3:
		prefix = "📝 Implementation Details"
	default:
		prefix = "💡 Hint"
	}

	return fmt.Sprintf("%s:\n%s", prefix, hint)
}

// FormatCodeReview formats a code review response
func (rf *ResponseFormatter) FormatCodeReview(review string) string {
	sections := []string{
		"✅ Strengths",
		"⚠️ Issues",
		"💡 Suggestions",
		"📊 Complexity Analysis",
		"🎯 Next Steps",
	}

	formatted := "🔍 Code Review\n" + strings.Repeat("─", 50) + "\n\n"
	
	// Try to identify and format sections
	for _, section := range sections {
		if strings.Contains(review, section) {
			formatted = strings.ReplaceAll(review, section, "\n"+section+"\n")
		}
	}

	if formatted == "🔍 Code Review\n" + strings.Repeat("─", 50) + "\n\n" {
		// If no sections found, just return the review as-is
		return formatted + review
	}

	return formatted
}

// FormatError formats an error message for display
func (rf *ResponseFormatter) FormatError(err error) string {
	return fmt.Sprintf("❌ Error: %v", err)
}

// FormatSuccess formats a success message
func (rf *ResponseFormatter) FormatSuccess(message string) string {
	return fmt.Sprintf("✅ %s", message)
}