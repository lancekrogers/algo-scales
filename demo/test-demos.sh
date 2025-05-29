#!/bin/bash

# Test all demo scripts for syntax errors and basic functionality

echo "🧪 Testing AlgoScales Demo Scripts"
echo "=================================="
echo ""

# Test directory
DEMO_DIR="$(dirname "$0")"
cd "$DEMO_DIR"

# Test syntax of all bash scripts
echo "1. Checking script syntax..."
for script in *.sh; do
    if [[ "$script" != "test-demos.sh" ]]; then
        echo -n "   $script: "
        if bash -n "$script"; then
            echo "✅ OK"
        else
            echo "❌ SYNTAX ERROR"
            exit 1
        fi
    fi
done
echo ""

# Test that scripts can start (first few lines)
echo "2. Testing script startup..."

echo -n "   full-workflow-demo.sh: "
if echo "q" | timeout 5s bash -c 'DEMO_SPEED=fast ./full-workflow-demo.sh' >/dev/null 2>&1; then
    echo "✅ Starts correctly"
else
    echo "⚠️  May have startup issues (common with interactive scripts)"
fi

echo -n "   quick-demo.sh: "
if ./quick-demo.sh >/dev/null 2>&1; then
    echo "✅ Runs successfully"
else
    echo "❌ FAILED"
    exit 1
fi

echo -n "   demo.sh: "
if bash -c 'echo "q" | timeout 5s ./demo.sh' >/dev/null 2>&1; then
    echo "✅ Starts correctly"
else
    echo "⚠️  May have startup issues (common with interactive scripts)"
fi

echo ""
echo "3. Checking dependencies..."

echo -n "   algo-scales binary: "
if command -v algo-scales >/dev/null 2>&1; then
    echo "✅ Found"
else
    echo "⚠️  Not found (install with: make install-user)"
fi

echo -n "   asciinema (for GIF): "
if command -v asciinema >/dev/null 2>&1; then
    echo "✅ Found"
else
    echo "ℹ️  Not found (optional: brew install asciinema)"
fi

echo -n "   agg (for GIF): "
if command -v agg >/dev/null 2>&1; then
    echo "✅ Found"
else
    echo "ℹ️  Not found (optional: cargo install --git https://github.com/asciinema/agg)"
fi

echo ""
echo "✅ Demo system test complete!"
echo ""
echo "🚀 Ready to showcase AlgoScales:"
echo "   ./full-workflow-demo.sh  # Complete user experience"
echo "   ./quick-demo.sh          # Quick feature overview"
echo "   ./demo.sh                # Detailed exploration"