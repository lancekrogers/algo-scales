#!/bin/bash

# AlgoScales Quick Demo - Shows core commands without interaction
# Perfect for screenshots, documentation, or quick evaluation

set -e

echo "🎵 AlgoScales Quick Demo"
echo "========================"
echo ""

# Check if algo-scales is available
if ! command -v algo-scales >/dev/null 2>&1; then
    echo "❌ algo-scales not found. Please install first:"
    echo "   make install-user"
    exit 1
fi

echo "1. List algorithm patterns (musical scales):"
echo "$ algo-scales list patterns"
algo-scales list patterns
echo ""

echo "2. List problems for a specific pattern:"
echo "$ algo-scales list hash-map"
algo-scales list hash-map
echo ""

echo "3. Check AI assistant status:"
echo "$ algo-scales ai config"
algo-scales ai config 2>/dev/null || echo "AI not configured (optional)"
echo ""

echo "4. View practice statistics:"
echo "$ algo-scales stats"
algo-scales stats 2>/dev/null || echo "No practice sessions yet"
echo ""

echo "5. Start daily practice (shows what would happen):"
echo "$ algo-scales daily --dry-run"
echo "🎵 Daily Scale Practice"
echo "Today: C Major (Sliding Window) - Maximum Sum Subarray"
echo "Progress: 0/11 patterns completed today"
echo "Current streak: Start your streak today! 🎯"
echo ""

echo "6. Available learning modes:"
echo "$ algo-scales start --help"
algo-scales start --help 2>/dev/null | head -15 || echo "Start command help not available"
echo ""

echo "✅ Quick demo complete!"
echo ""
echo "Key features demonstrated:"
echo "• Pattern-based problem organization"
echo "• AI assistant integration"
echo "• Progress tracking"
echo "• Multiple learning modes"
echo ""
echo "Run './demo.sh' for the full interactive experience!"