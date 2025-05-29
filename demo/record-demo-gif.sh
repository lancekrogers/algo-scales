#!/bin/bash

# Script to automatically record the demo as a GIF
# Requires: asciinema, agg (asciinema gif generator)

set -e

echo "🎬 AlgoScales Demo GIF Generator"
echo "================================"
echo ""
echo "Usage: $0 [demo-script] [title]"
echo "Example: $0 ./quick-demo.sh 'AlgoScales Quick Demo'"
echo ""

# Check dependencies
check_dependency() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "❌ $1 not found. Installing..."
        case "$1" in
            "asciinema")
                if command -v brew >/dev/null 2>&1; then
                    brew install asciinema
                elif command -v pip >/dev/null 2>&1; then
                    pip install asciinema
                else
                    echo "Please install asciinema: https://asciinema.org/docs/installation"
                    exit 1
                fi
                ;;
            "agg")
                if command -v cargo >/dev/null 2>&1; then
                    cargo install --git https://github.com/asciinema/agg
                else
                    echo "Please install Rust/Cargo, then: cargo install --git https://github.com/asciinema/agg"
                    exit 1
                fi
                ;;
        esac
    else
        echo "✅ $1 found"
    fi
}

# Install dependencies
echo "Checking dependencies..."
check_dependency "asciinema"
check_dependency "agg"
echo ""

# Configuration
DEMO_NAME="algoscales-workflow-demo"
OUTPUT_DIR="demo-assets"
ASCIICAST_FILE="$OUTPUT_DIR/$DEMO_NAME.cast"
GIF_FILE="$OUTPUT_DIR/$DEMO_NAME.gif"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "🎥 Recording demo..."
echo "This will run the full workflow demo automatically."
echo "Press Ctrl+C if you need to stop."
echo ""

# Select which demo to record
DEMO_SCRIPT="${1:-./full-workflow-demo.sh}"
DEMO_TITLE="${2:-AlgoScales Demo}"

echo "Recording: $DEMO_SCRIPT"

# Record the demo
echo "" | DEMO_SPEED=fast asciinema rec "$ASCIICAST_FILE" \
    --command="$DEMO_SCRIPT" \
    --title="$DEMO_TITLE" \
    --idle-time-limit=2 \
    --overwrite

echo ""
echo "🎨 Converting to GIF..."

# Convert to GIF with optimized settings
agg \
    --theme=monokai \
    --font-size=14 \
    --cols=100 \
    --rows=30 \
    --speed=1.5 \
    "$ASCIICAST_FILE" \
    "$GIF_FILE"

echo ""
echo "✅ Demo GIF created: $GIF_FILE"
echo ""
echo "📊 File size: $(du -h "$GIF_FILE" | cut -f1)"
echo ""
echo "🚀 Ready for README and social media!"

# Also create a shorter version for social media
echo ""
echo "🎬 Creating short version for social media..."

# Create a condensed version of the demo
cat > "condensed-demo.sh" << 'EOF'
#!/bin/bash
set -e

echo "🎵 AlgoScales - Algorithm Learning with Musical Patterns"
sleep 1

echo ""
echo "$ algo-scales list patterns"
echo "Available patterns: Hash Map, Sliding Window, Two Pointers..."
sleep 2

echo ""
echo "$ algo-scales start learn two_sum"
echo "🎵 Learn Mode: Two Sum (Hash Map Pattern)"
echo "💡 Use hash map for O(1) lookups instead of O(n²) nested loops"
sleep 3

echo ""
echo "$ algo-scales start practice max_sum_subarray"
echo "🎵 Practice Mode: Sliding Window Pattern"
echo "⏱️  Timer: 30:00"
echo "💻 Coding solution..."
sleep 2

echo ""
echo "✅ All tests passed! Solution accepted!"
echo "🎯 Pattern mastery: Sliding Window +1"
sleep 2

echo ""
echo "$ algo-scales ai review"
echo "🤖 AI: Excellent sliding window technique!"
echo "💡 Try Two Pointers pattern next"
sleep 2

echo ""
echo "$ algo-scales daily"
echo "🎵 Daily Practice: 3/11 patterns completed"
echo "🔥 Current streak: 7 days"
sleep 2

echo ""
echo "🎵 Master algorithms through musical patterns + AI guidance!"
echo "⭐ github.com/your-username/algo-scales"
EOF

chmod +x condensed-demo.sh

# Record short version
SHORT_ASCIICAST="$OUTPUT_DIR/$DEMO_NAME-short.cast"
SHORT_GIF="$OUTPUT_DIR/$DEMO_NAME-short.gif"

echo "" | asciinema rec "$SHORT_ASCIICAST" \
    --command="./condensed-demo.sh" \
    --title="AlgoScales Quick Demo" \
    --idle-time-limit=1 \
    --overwrite

agg \
    --theme=monokai \
    --font-size=16 \
    --cols=80 \
    --rows=20 \
    --speed=1.2 \
    "$SHORT_ASCIICAST" \
    "$SHORT_GIF"

rm condensed-demo.sh

echo "✅ Short demo GIF created: $SHORT_GIF"
echo "📊 File size: $(du -h "$SHORT_GIF" | cut -f1)"
echo ""

echo "📁 Generated files:"
echo "  - $GIF_FILE (full demo)"
echo "  - $SHORT_GIF (social media version)"
echo "  - $ASCIICAST_FILE (source recording)"
echo ""
echo "🎯 Usage suggestions:"
echo "  - Add full GIF to README"
echo "  - Use short GIF for social media"
echo "  - Upload asciicast to asciinema.org for web players"