# AlgoScales CLI Mode

AlgoScales operates primarily in command-line interface (CLI) mode, which allows you to practice algorithm problems efficiently. The CLI mode is the default way to interact with AlgoScales, providing a streamlined experience for practicing algorithms in your preferred code editor.

## Getting Started with CLI Mode

CLI mode is the default operating mode for AlgoScales. Simply run any command without flags to use CLI mode:

```bash
# Start a practice session
algo-scales start practice

# Start a learn session (shows solutions)
algo-scales start learn

# Use the CLI solve command for single problems
algo-scales cli solve
```

> **Note**: The Terminal UI mode (`--tui` flag) is currently a work in progress and may not function as expected.

## Available Commands in CLI Mode

### Starting Sessions

```bash
# Start a practice session (no solutions shown)
algo-scales start practice

# Start a learn session (solutions available)
algo-scales start learn

# Start a cram session (rapid-fire practice)
algo-scales start cram

# Start with a specific pattern
algo-scales start practice --pattern sliding-window

# Start with a specific difficulty
algo-scales start practice --difficulty medium

# Start in a specific language
algo-scales start practice --language python  # Options: go, python, javascript
```

### CLI Solve Command

```bash
# Solve a single problem
algo-scales cli solve

# Solve a specific problem by ID
algo-scales cli solve two-sum

# Solve with specific options
algo-scales cli solve --pattern hash-map --language python
```

### Daily Practice

```bash
# Start daily practice (one problem from each pattern)
algo-scales daily

# Test your solution for the current problem
algo-scales daily test

# Skip the current problem
algo-scales daily skip

# Resume a skipped problem
algo-scales daily resume-skipped

# Check your daily practice status
algo-scales daily status
```

### Listing Problems

```bash
# List all available problems
algo-scales list

# List problems by pattern
algo-scales list patterns

# List problems by difficulty
algo-scales list difficulties

# List problems by company
algo-scales list companies
```

### AI Assistant

```bash
# Configure AI provider
algo-scales ai config

# Test AI configuration
algo-scales ai test

# Get AI hints (in vim)
:AlgoScalesAIHint

# Start AI chat (in vim)
:AlgoScalesAIChat

# Review code with AI
algo-scales ai review [file]
```

### Viewing Statistics

```bash
# View your problem-solving statistics
algo-scales stats

# View statistics by pattern
algo-scales stats patterns

# View your progress trends
algo-scales stats trends

# Reset your statistics
algo-scales stats reset
```

## CLI Problem-Solving Workflow

When you start a problem in CLI mode, AlgoScales creates a temporary workspace with the problem description and starter code. The workflow is:

1. **Choose a problem**: Either specify one or let AlgoScales pick a random one
2. **View the problem**: Read the problem description
3. **Edit your solution**: Write your code in your preferred editor
4. **Test your solution**: Run tests to see if your solution passes
5. **Submit**: When your code passes all tests, the problem is marked as solved

### Interactive Menu

The CLI mode provides an interactive menu to guide you through the problem-solving process:

```
🎵 AlgoScales CLI Mode 🎵
—————————————————————————
Problem: Two Sum (Easy)
Pattern: hash-map
Estimated Time: 15 minutes

Options:
1. View problem description
2. Edit solution
3. Test solution
4. Exit
```

## Working with Files

AlgoScales CLI mode creates and manages these files for you:

- `problem.md`: Contains the formatted problem description
- `solution.{ext}`: Your solution file with the appropriate extension for your language

## Testing Your Solution

When you test your solution through the CLI menu:

1. AlgoScales runs your code against the problem's test cases
2. You'll see a detailed report of test results
3. If all tests pass, the problem is marked as solved and statistics are recorded

```
--- Test Results ---

Test 1: ✅ PASSED
Input: [2,7,11,15], 9
Expected: [0,1]
Actual: [0,1]

Test 2: ✅ PASSED
Input: [3,2,4], 6
Expected: [1,2]
Actual: [1,2]

🎉 All tests passed! Problem solved! 🎉
```

## Statistics Tracking

AlgoScales tracks your progress even in CLI mode. You can view your statistics with:

```bash
algo-scales stats
```

This shows:
- Overall progress and success rate
- Breakdown by algorithm pattern
- Recent activity
- Current practice streak

For more detailed views:
```bash
algo-scales stats patterns  # Pattern-specific stats
algo-scales stats trends    # Progress over time
```

## Environment Variables

- `EDITOR`: Set this to your preferred text editor for editing solutions

## Vim Integration

If you're using the vim plugin, additional commands are available:

```vim
" Start a new session
:AlgoScalesStart

" List problems
:AlgoScalesList

" Get progressive hints
:AlgoScalesHint

" Get AI hints
:AlgoScalesAIHint

" Start AI chat
:AlgoScalesAIChat

" Run tests
:AlgoScalesTest
```

## Tips for CLI Mode

1. **Language Selection**: Use the `--language` flag to choose your preferred programming language
2. **Session Management**: Sessions are automatically saved and statistics tracked
3. **Pattern Practice**: Use the `--pattern` flag to focus on specific algorithm patterns
4. **View Statistics**: Check your progress regularly with the stats command

## Troubleshooting

If you encounter issues:

1. **Command Not Found**: Make sure AlgoScales is in your PATH
2. **Editor Not Opening**: Set the EDITOR environment variable
3. **Test Failures**: Check the error messages for syntax or logic issues
4. **Language Issues**: Ensure you have the appropriate language runtime installed (Go, Python, Node.js)