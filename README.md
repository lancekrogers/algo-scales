# Algo Scales

An algorithm study tool designed for efficient interview preparation, focusing on common patterns used in technical coding interviews.

## Overview

Algo Scales is a terminal-based application that helps developers practice and master algorithm patterns for technical interviews. It emphasizes focused learning through curated problems organized by algorithm patterns, offering different learning modes to accommodate various study needs.

## Features

- **Multiple Learning Modes**:

  - **Learn Mode**: Shows pattern explanations and detailed walkthroughs
  - **Practice Mode**: Hides solutions but allows hints on demand
  - **Cram Mode**: Rapid-fire practice with timers for interview preparation

- **Pattern-Based Learning**: Problems organized by common algorithm patterns (sliding window, two pointers, DFS, etc.)

- **Distraction-Free Environment**: Terminal-based UI keeps you focused on the problem at hand

- **Statistics Tracking**: Records your progress and performance over time

- **Configurable Timer**: Set time limits to simulate interview conditions

- **Multiple Language Support**: Practice in Go, Python, or JavaScript

## Installation

### Prerequisites

- Go 1.16 or higher
- A terminal editor (vim, nano, etc.) configured through the `EDITOR` environment variable

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Blockhead-Consulting/algo-scales.git
cd algo-scales

# Build the CLI
go build -o algo-scales

# Build the API server (optional)
cd server
go build -o algo-scales-server
```

### Download Binary

Pre-built binaries are available for the following platforms:

- Linux (x86_64)
- macOS (x86_64, arm64)
- Windows (x86_64)

Download the appropriate binary for your platform from the [releases page](https://github.com/lancekrogers/algo-scales/releases).

## Usage

### First Run

The first time you run Algo Scales, you'll be asked to enter your license key and email. After validation, the tool will download the problem sets.

```bash
./algo-scales
```

### Available Commands

```bash
# Start a session in Learn mode
./algo-scales start learn [problem-id]

# Start a session in Practice mode
./algo-scales start practice [problem-id]

# Start a session in Cram mode (rapid-fire problems)
./algo-scales start cram

# List all available problems
./algo-scales list

# List problems by pattern
./algo-scales list patterns

# List problems by difficulty
./algo-scales list difficulties

# List problems by company
./algo-scales list companies

# View your statistics
./algo-scales stats

# View statistics by pattern
./algo-scales stats patterns

# View your progress trends
./algo-scales stats trends

# Reset your statistics
./algo-scales stats reset
```

### Options

```bash
# Set the programming language (default: go)
./algo-scales start learn --language python

# Set the timer duration in minutes (default: 45)
./algo-scales start practice --timer 30

# Focus on a specific algorithm pattern
./algo-scales start learn --pattern sliding-window

# Select by difficulty
./algo-scales start practice --difficulty medium
```

## In-Session Commands

When in a practice session, you can use the following keyboard shortcuts:

- `e`: Open your code in the configured editor
- `h`: Show hints (if available)
- `s`: Show solution (if available)
- `Enter`: Submit your solution
- `n`: Skip to the next problem
- `q` or `Ctrl+C`: Quit the session
- `?`: Show help

## API Server (Optional)

For license validation and problem downloads, you can run the API server:

```bash
cd server
./algo-scales-server
```

The server runs on port 8080 by default, but you can change this with the `PORT` environment variable.

## Licensing

Algo Scales is a commercial product licensed on a per-user basis. Each license is valid for one year from the purchase date.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Roadmap

- Neovim/LazyVim plugin integration
- AI-powered hints and code reviews
- More algorithm patterns and problems
- Interactive visualizations for algorithms
- Simulated interview mode with AI interviewers

## Credits

Created by Blockhead Consulting

## Support

For support, please email <lance@blockhead.consulting> or open an issue on GitHub.
