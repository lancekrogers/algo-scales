# AI Assistant for AlgoScales

AlgoScales includes an AI-powered assistant that helps you learn algorithm patterns more effectively. The assistant can provide intelligent hints, review your code, explain patterns, and guide you through problem-solving without simply giving away the answers.

## Quick Start

### 1. Install Prerequisites

#### For Claude Support
1. **Claude Max Subscription**: Required for Claude Code access at [claude.ai](https://claude.ai)
2. **Claude Code CLI**: Install following [official docs](https://docs.anthropic.com/en/docs/claude-code/getting-started)
3. **Verify installation**: `claude --help`

#### For Ollama Support (Local AI)
1. **Install Ollama**: `curl -fsSL https://ollama.com/install.sh | sh`
2. **Pull a model**: `ollama pull llama3`

### 2. Configure Your AI Provider

```bash
# Interactive configuration
algo-scales ai config

# Or set directly for Claude Code
algo-scales ai config set claude.cli_path "claude"  # or full path to claude binary
algo-scales ai config set default_provider claude

# For Ollama
algo-scales ai config set default_provider ollama
algo-scales ai config set ollama.model "llama3"
```

### 3. Get AI-Powered Hints

When working on a problem, request AI assistance:

```bash
# Get a smart hint for your current problem
algo-scales hint --ai

# Start an interactive chat session
algo-scales hint --ai --interactive
```

## Supported AI Providers

### Claude (via Claude Code)

Claude Code is a powerful CLI tool that provides access to Claude with enhanced capabilities:
- **Real-time tool usage**: See exactly what Claude is doing (reading files, running commands)
- **MCP (Model Context Protocol)**: Extended capabilities for code analysis
- **Session persistence**: Continue conversations across commands
- **Streaming responses**: See Claude's thinking in real-time

**Setup:**
1. Subscribe to Claude Max at [claude.ai](https://claude.ai)
2. Install Claude Code CLI (see Prerequisites above)
3. Configure: `algo-scales ai config set claude.cli_path "claude"`
4. Optional: Install MCP servers for enhanced capabilities:
   ```bash
   npm install -g @modelcontextprotocol/server-filesystem
   npm install -g @modelcontextprotocol/server-github
   ```

### Ollama (Local Models)

Ollama allows you to run AI models locally, providing:
- Complete privacy - your code never leaves your machine
- No API costs
- Works offline
- Support for various open-source models

**Setup:**
1. Install Ollama: `curl -fsSL https://ollama.com/install.sh | sh`
2. Pull a model: `ollama pull llama3`
3. Configure: `algo-scales ai config set ollama.model "llama3"`
4. Set as default: `algo-scales ai config set default_provider ollama`

## Features

### 1. Progressive Hints

The AI assistant provides hints at different levels:

- **Level 1**: General approach and pattern recognition
- **Level 2**: Specific algorithmic guidance
- **Level 3**: Detailed implementation tips

Example:
```bash
# First hint - gentle nudge
$ algo-scales hint --ai
"Consider using a hash map to track elements you've seen..."

# Need more help?
$ algo-scales hint --ai
"You'll want to iterate through the array once, storing each element..."

# Still stuck?
$ algo-scales hint --ai
"Here's the approach: As you iterate, check if target - current exists..."
```

### 2. Interactive Chat Mode

Start a conversation with the AI about your current problem:

```bash
$ algo-scales hint --ai --interactive

🤖 AI Assistant Ready! Type 'help' for commands or 'exit' to quit.

You> I'm not sure how to handle the edge case where the array is empty

Assistant> Good thinking about edge cases! For an empty array, consider what your 
function should return. In the Two Sum problem, if the array is empty, there 
can't be any two numbers that sum to the target. What would be an appropriate 
return value in this case?

You> Should I return null or an empty array?

Assistant> Great question! Check the problem requirements. Most implementations 
expect a consistent return type. If the normal case returns an array of indices 
[i, j], then returning an empty array [] for the "no solution" case maintains 
type consistency. This is generally preferred over null in most languages...
```

### 3. Code Review with Tool Usage

Get feedback on your solution with full transparency:

```bash
$ algo-scales review --ai solution.py

🔍 AI Code Review for: Two Sum Problem

[Using tool: mcp__filesystem__read_file]
[Using tool: Bash - running tests]

✅ Strengths:
- Correct use of hash map for O(n) time complexity
- Good variable naming
- Handles the basic case correctly

💡 Suggestions:
1. Edge Cases: Consider adding a check for null/empty array
2. Optimization: You're checking `if complement in seen` after adding 
   the current number. This could miss cases where the same element 
   is used twice.
3. Style: Consider using more descriptive variable names than 'i' and 'j'

📝 Pattern Recognition:
You've correctly identified this as a hash map problem! This pattern 
appears in many other problems like "Three Sum" and "Subarray Sum".

💰 Session Cost: $0.0234 | Session ID: abc123
```

### 4. Pattern Explanations

Learn about algorithm patterns with tailored examples:

```bash
$ algo-scales learn sliding-window --ai

The AI will explain:
- What the sliding window pattern is
- When to use it
- Common variations
- Practice problems that use this pattern
- Implementation tips
```

## Configuration Options

The AI assistant is configured via `~/.algo-scales/ai-config.yaml`:

```yaml
# Basic settings
default_provider: claude  # or 'ollama'

# Claude Code settings
claude:
  cli_path: "claude"  # Path to claude binary
  default_format: "json"  # Output format (text, json, stream-json)
  save_sessions: true  # Persist conversations
  session_directory: "~/.algo-scales/claude-sessions"
  max_turns: 5  # Max conversation turns per interaction
  
  # MCP (Model Context Protocol) settings
  mcp:
    enabled: true
    servers:
      filesystem:
        command: "npx"
        args: ["-y", "@modelcontextprotocol/server-filesystem", "./"]
  
  # Safety settings
  allowed_tools:
    - "mcp__filesystem__read_file"
    - "Bash"  # Allow bash commands
    - "Read"  # Allow file reading

# Ollama settings  
ollama:
  host: "http://localhost:11434"  # Default Ollama server
  model: "llama3"  # or codellama, mixtral, etc.
  temperature: 0.7

# Behavior settings
features:
  code_review: true  # Enable AI code review
  interactive_repl: true  # Enable chat mode
  auto_review: false  # Auto-review before submission
```

## Best Practices

### 1. Learning, Not Cheating

The AI assistant is designed to help you learn, not to solve problems for you:

- ✅ Ask for explanations of concepts
- ✅ Request hints when stuck
- ✅ Get feedback on your approach
- ❌ Don't ask for complete solutions
- ❌ Don't copy-paste AI code without understanding

### 2. Effective Prompting

Get better help by being specific:

```bash
# Less effective
"Help me solve this"

# More effective  
"I'm trying to use a two-pointer approach but I'm not sure how to handle duplicates"

# Even better
"My solution works for most cases but fails when the array has negative numbers. 
Here's my approach: [explain your logic]. What am I missing?"
```

### 3. Use the Right Provider

- **Use Claude for**: Complex algorithm explanations, nuanced hints, code review
- **Use Ollama for**: Privacy-sensitive code, offline work, quick hints

## Privacy and Security

### Data Handling

- **Claude Code**: Uses the official Claude Code CLI tool. Your code is processed according to Anthropic's privacy policy. Sessions can be saved locally for continuity.
- **Ollama**: Everything runs locally. Your code never leaves your machine.

### Security Notes

```bash
# Claude Code handles authentication automatically
# No API keys needed in config files!

# For config file permissions
chmod 600 ~/.algo-scales/ai-config.yaml

# Session data is stored locally
chmod -R 700 ~/.algo-scales/claude-sessions/
```

## Troubleshooting

### Common Issues

**"Claude command not found"**
```bash
# Make sure Claude Code is installed
claude --help

# Check your configuration
algo-scales ai config show

# Set the correct path if needed
algo-scales ai config set claude.cli_path "/path/to/claude"
```

**"Not authenticated to Claude"**
```bash
# Claude Code handles authentication automatically
# Make sure you have a Claude Max subscription
# The CLI will guide you through login if needed
```

**"Ollama connection refused"**
```bash
# Make sure Ollama is running
ollama serve

# Check if model is downloaded
ollama list
```

**"Rate limit exceeded"**
- Claude Max provides generous limits for Claude Code usage
- For unlimited requests, consider using Ollama locally

### Debug Mode

Enable detailed logging for troubleshooting:

```bash
algo-scales ai config set logging.level debug
algo-scales ai config set logging.log_interactions true
```

Logs are saved to `~/.algo-scales/ai-assistant.log`.

## Examples

### Example 1: Learning Two Pointers Pattern

```bash
$ algo-scales start practice --pattern two-pointers
$ algo-scales hint --ai --interactive

You> I don't understand when to move which pointer

[Using tool: mcp__filesystem__read_file - examining your current code]

Assistant> Great question! Looking at your code, I can see you're working with 
the two-pointer pattern. The key is to understand what each pointer represents 
and what condition should trigger their movement.

Think of it like this:
- Left pointer: Represents the "smaller" or "start" element
- Right pointer: Represents the "larger" or "end" element

The decision to move a pointer usually depends on comparing your current sum/product/
difference with your target...

💰 Cost: $0.0023 | Continuing session: xyz789
```

### Example 2: Debugging with AI

```bash
$ algo-scales test solution.go
❌ Test case 3 failed

$ algo-scales hint --ai
You> My solution fails on test case 3 but I can't see the input

Assistant> When a specific test case fails, let's debug systematically:
1. First, what does your current solution assume about the input?
2. Test case 3 often involves edge cases. Could it be:
   - An empty array?
   - Single element?
   - All duplicates?
   - Negative numbers?

Try adding logging to see what your function receives...
```

### Example 3: Pattern Recognition

```bash
$ algo-scales review --ai solution.py

Assistant> I notice you're using a nested loop here (O(n²) complexity). 
This problem actually follows the "sliding window" pattern. 

Here's why: You're looking for a subarray with a specific property, and 
that subarray is contiguous. Instead of checking every possible subarray, 
you can maintain a "window" and slide it across the array...
```

## Integration with Editors

### Neovim Integration

If you're using the AlgoScales Neovim plugin:

```vim
" Get AI hint for current problem
:AlgoScalesAIHint

" Start interactive AI chat in split
:AlgoScalesAIChat

" Review current buffer with AI
:AlgoScalesAIReview
```

### VS Code Integration (Coming Soon)

The AI assistant will be available in the VS Code extension with:
- Inline hints
- Code lens suggestions  
- Problem-specific chat panel

## Feedback

The AI assistant is continuously improving. Help us make it better:

- Report issues: `algo-scales feedback "AI gave incorrect hint for..."`
- Suggest improvements: Open an issue on GitHub
- Share success stories: Let us know how AI helped you learn!

Remember: The goal is to help you become a better problem solver, not to solve problems for you. Use the AI assistant as a learning tool, and you'll find yourself needing it less over time as you master algorithm patterns!