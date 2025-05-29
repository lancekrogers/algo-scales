#!/bin/bash

# Manual Vim Plugin Test Script
# Run this to test the vim plugin without modifying your vim config

echo "🎵 AlgoScales Vim Plugin Manual Test"
echo "====================================="

# Build binary if needed
if [ ! -f "bin/algo-scales" ]; then
    echo "Building AlgoScales binary..."
    make build
fi

# Create a temporary vim config
TEMP_DIR="/tmp/algoscales-vim-manual-test"
mkdir -p "$TEMP_DIR"

cat > "$TEMP_DIR/vimrc" << EOF
" Temporary vim config for AlgoScales testing
set runtimepath+=$(pwd)/vim-plugin

" Plugin configuration
let g:algo_scales_path = '$(pwd)/bin/algo-scales'
let g:algo_scales_language = 'go'
let g:algo_scales_auto_test = 1

" Enable better colors
syntax on
set number
set background=dark

" Show available commands
echo "AlgoScales Vim Plugin Loaded!"
echo "Available commands:"
echo "  :AlgoScalesStart [problem] - Start a session"
echo "  :AlgoScalesList - List all problems"
echo "  :AlgoScalesDaily - Start daily practice"
echo "  :AlgoScalesTest - Test current solution"
echo "  :AlgoScalesHint - Get a hint"
echo ""
echo "Try: :AlgoScalesStart two_sum"
EOF

echo ""
echo "Starting vim with AlgoScales plugin..."
echo "Commands to try:"
echo "  :AlgoScalesStart two_sum"
echo "  :AlgoScalesList"  
echo "  :AlgoScalesDaily"
echo ""
echo "Press any key to continue..."
read -n 1

# Start vim with our temporary config
vim -u "$TEMP_DIR/vimrc"

# Cleanup
rm -rf "$TEMP_DIR"
echo "Test completed!"