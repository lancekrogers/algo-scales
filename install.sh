#!/bin/bash

# AlgoScales Installation Script
# This script installs AlgoScales without requiring sudo

set -e

echo "🎵 AlgoScales Installation Script"
echo "================================="
echo ""

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    echo "❌ Go is not installed. Please install Go 1.16+ first:"
    echo "   https://golang.org/doc/install"
    exit 1
fi

echo "✅ Go found: $(go version)"
echo ""

# Option 1: Use go install (recommended)
echo "🚀 Installing AlgoScales using 'go install'..."
if go install github.com/lancekrogers/algo-scales@latest; then
    echo "✅ AlgoScales installed successfully!"
    echo ""
    echo "📍 Binary location: $GOPATH/bin/algo-scales (or $HOME/go/bin/algo-scales)"
    echo ""
    
    # Check if GOPATH/bin or $HOME/go/bin is in PATH
    GOBIN="${GOPATH:-$HOME/go}/bin"
    if echo "$PATH" | grep -q "$GOBIN"; then
        echo "✅ $GOBIN is already in your PATH"
        echo "🎉 You can now run: algo-scales"
    else
        echo "⚠️  $GOBIN is not in your PATH"
        echo ""
        echo "🔧 Add this to your ~/.bashrc or ~/.zshrc:"
        echo "   export PATH=\$PATH:$GOBIN"
        echo ""
        echo "Then reload your shell or run:"
        echo "   source ~/.bashrc  # or source ~/.zshrc"
    fi
    
    echo ""
    echo "🚀 Quick start:"
    echo "   algo-scales ai config   # Configure AI assistant (optional)"
    echo "   algo-scales             # Start practicing!"
    
else
    echo "❌ 'go install' failed. Falling back to building from source..."
    echo ""
    
    # Option 2: Build from source
    echo "📥 Cloning repository..."
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    if git clone https://github.com/lancekrogers/algo-scales.git; then
        cd algo-scales
        echo "🔨 Building AlgoScales..."
        
        if make install-user; then
            echo ""
            echo "🎉 Installation complete!"
            echo "🗑️  Cleaning up temporary files..."
            cd /
            rm -rf "$TEMP_DIR"
        else
            echo "❌ Build failed"
            exit 1
        fi
    else
        echo "❌ Failed to clone repository"
        exit 1
    fi
fi

echo ""
echo "🎵 Happy coding with AlgoScales!"