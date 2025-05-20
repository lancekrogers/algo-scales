# AlgoScales CLI Mode

AlgoScales provides a powerful command-line interface (CLI) mode that allows you to practice algorithm problems without using the Terminal UI interface. This is useful for environments where a full TUI isn't practical, or for users who prefer to work in their own editor.

## Getting Started with CLI Mode

To use CLI mode, run AlgoScales with the `--cli` flag or use the `solve` command:

```bash
# Use the CLI flag with any command
algo-scales start practice --cli

# Or use the dedicated solve command
algo-scales solve
```

## Available Commands in CLI Mode

### Solving Problems

```bash
# Solve a random problem
algo-scales solve

# Solve a specific problem by ID
algo-scales solve two-sum

# Solve a problem with a specific pattern
algo-scales solve --pattern sliding-window

# Solve a problem with a specific difficulty
algo-scales solve --difficulty medium

# Solve in a specific language
algo-scales solve --language python  # Options: go, python, javascript
```

### Daily Practice

```bash
# Start daily practice (one problem from each pattern)
algo-scales daily --cli

# Test your solution for the current problem
algo-scales daily test

# Skip the current problem
algo-scales daily skip

# Resume a skipped problem
algo-scales daily resume-skipped

# Check your daily practice status
algo-scales daily status
```

### Viewing Statistics

```bash
# View your problem-solving statistics
algo-scales solve stats
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
algo-scales solve stats
```

This shows:
- Overall progress and success rate
- Breakdown by algorithm pattern
- Recent activity
- Current practice streak

## Environment Variables

- `EDITOR`: Set this to your preferred text editor for editing solutions

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