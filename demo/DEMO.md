# AlgoScales Demo Guide

This document explains how to demonstrate AlgoScales functionality for evaluation, testing, or showcasing the project.

## Quick Demo (5 minutes)

For a rapid overview of core features:

```bash
./quick-demo.sh
```

This shows:
- Pattern-based problem organization
- Available commands and modes
- AI integration status
- Statistics tracking

## Interactive Demo (15-20 minutes)

For a comprehensive walkthrough including simulated problem solving:

```bash
./demo.sh
```

This demonstrates:
- Complete learning workflow
- Multiple practice modes
- AI-powered assistance
- Progress tracking
- Multi-language support
- Daily practice system

## Manual Testing

### Basic Commands

```bash
# List all algorithm patterns
algo-scales list patterns

# List problems for specific pattern
algo-scales list sliding-window

# View current statistics
algo-scales stats

# Check AI configuration
algo-scales ai config
```

### Learning Modes

```bash
# Learn mode - shows solutions and hints
algo-scales start learn two_sum

# Practice mode - timed practice with hints available
algo-scales start practice max_sum_subarray

# Daily practice - systematic pattern coverage
algo-scales daily
```

### AI Features

```bash
# Configure AI assistant
algo-scales ai config

# Get AI hint for a problem
algo-scales ai hint sliding-window

# Request AI code review
algo-scales ai review

# Start AI chat session
algo-scales ai chat
```

## Demo Solutions

The interactive demo includes pre-written solutions for several problems:

- **Two Sum** (Hash Map pattern)
- **Maximum Sum Subarray** (Sliding Window pattern)  
- **Search in Rotated Array** (Binary Search pattern)

These solutions demonstrate the expected code quality and pattern implementation.

## Expected Behavior

### Pattern Organization
- 11 fundamental algorithm patterns (musical scales)
- Problems categorized by their primary solving technique
- Clear pattern explanations and use cases

### Learning Progression
- **Learn Mode**: Full solutions visible, pattern explanations provided
- **Practice Mode**: Solutions hidden, hints available on demand
- **Daily Mode**: One problem per pattern for consistent practice

### AI Integration
- Progressive hints (3 levels of increasing detail)
- Code review with suggestions and improvements
- Pattern-specific explanations and guidance
- Interactive chat for questions

### Progress Tracking
- Statistics per pattern and overall
- Solve times and attempt counts
- Daily practice streaks
- Difficulty progression

## Troubleshooting Demo Issues

### "Command not found" errors
Ensure AlgoScales is installed and in PATH:
```bash
make install-user
export PATH=$PATH:~/.local/bin  # or ~/go/bin
```

### AI features not working
Configure an AI provider first:
```bash
algo-scales ai config
```

### No problems available
The demo uses built-in problems. If issues persist:
```bash
algo-scales list patterns  # Should show 11 patterns
```

## Demo Script Details

### Quick Demo (`quick-demo.sh`)
- **Purpose**: Rapid feature overview
- **Duration**: 2-3 minutes
- **Audience**: Evaluators, potential users
- **Output**: Command examples and explanations

### Interactive Demo (`demo.sh`)
- **Purpose**: Complete workflow demonstration
- **Duration**: 15-20 minutes
- **Audience**: Detailed evaluation, feature exploration
- **Output**: Full simulated learning session

### Demo Solutions Directory
- **Location**: `demo_solutions/`
- **Content**: Working code examples for major patterns
- **Languages**: Go implementations with proper documentation
- **Cleanup**: Automatically removed after demo

## Customizing Demos

### Adding New Demo Problems
1. Create solution in `demo_solutions/`
2. Add section to demo script
3. Include pattern explanation

### Modifying Demo Flow
1. Edit section order in `demo.sh`
2. Adjust wait times between sections
3. Customize explanatory text

### Creating Language-Specific Demos
1. Add solutions in target language to `demo_solutions/`
2. Update demo script to use `--language` flag
3. Highlight language-specific features

## Production Readiness

### Current Status
- ✅ Core CLI functionality complete
- ✅ AI integration working
- ✅ Pattern-based learning implemented
- ✅ Progress tracking functional
- ✅ Multi-language support
- 🚧 Neovim plugin in development
- 📋 VS Code extension planned

### Demo Confidence Level
- **CLI Features**: Ready for production use
- **AI Integration**: Fully functional with proper setup
- **Learning System**: Proven pedagogical approach
- **Problem Quality**: Curated algorithm fundamentals

This demo system provides multiple ways to evaluate AlgoScales without requiring manual coding, making it ideal for showcasing the project's capabilities and educational value.