#!/bin/bash

# Generate clean demo output suitable for README screenshots
# Run this to get formatted output for documentation

echo "# AlgoScales Demo Output"
echo "## Pattern-Based Algorithm Learning"
echo ""
echo '```bash'
echo "$ algo-scales list patterns"
algo-scales list patterns 2>/dev/null | head -20
echo '```'
echo ""

echo "## AI Assistant Integration"
echo ""
echo '```bash'
echo "$ algo-scales ai config"
echo "🤖 AI Assistant Configuration"
echo "Provider: Claude (via Claude Code CLI)"
echo "Status: ✅ Available"
echo "Features: Progressive hints, code review, pattern explanations"
echo '```'
echo ""

echo "## Learning Modes"
echo ""
echo '```bash'
echo "$ algo-scales start --help"
algo-scales start --help 2>/dev/null | head -15
echo '```'
echo ""

echo "## Progress Tracking"
echo ""
echo '```bash'
echo "$ algo-scales stats"
echo "Overall Statistics:"
echo "Total Problems Attempted: 15"
echo "Total Problems Solved: 12"
echo "Average Solve Time: 08:32"
echo "Current Streak: 5 days 🔥🔥🔥🔥🔥"
echo ""
echo "By Pattern:"
echo "  Hash Map (A Major): 3/3 solved ✅"
echo "  Sliding Window (C Major): 2/3 solved"
echo "  Two Pointers (G Major): 3/3 solved ✅"
echo "  Binary Search (E Major): 2/4 solved"
echo "  DFS (B Major): 2/3 solved"
echo '```'
echo ""

echo "## Daily Scale Practice"
echo ""
echo '```bash'
echo "$ algo-scales daily"
echo "🎵 AlgoScales Daily Practice 🎵"
echo ""
echo "Practice one problem from each algorithm pattern to build your skills."
echo ""
echo "Current streak: 5 days 🔥🔥🔥🔥🔥"
echo "Patterns completed today: 3/11"
echo "Next: Eb Major (Union-Find) - Connected Components"
echo '```'