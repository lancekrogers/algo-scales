#!/bin/bash

# Setup and record with terminalizer (alternative approach)

echo "🎬 Setting up terminalizer for demo recording"
echo ""

# Install terminalizer
if ! command -v terminalizer >/dev/null 2>&1; then
    echo "Installing terminalizer..."
    npm install -g terminalizer
fi

# Create optimized config
cat > terminalizer-config.yml << 'EOF'
# Configuration for AlgoScales demo recording

# Specify a command to be executed
command: bash -c "DEMO_SPEED=fast ./full-workflow-demo.sh"

# Specify the current working directory path
cwd: null

# Export additional ENV variables
env:
  recording: true

# Explicitly set the number of columns
cols: 120

# Explicitly set the number of rows
rows: 30

# Amount of times to repeat GIF
repeat: 0

# Quality
quality: 100

# Delay between frames in ms
frameDelay: auto

# Maximum delay between frames in ms
maxIdleTime: 2000

# The surrounding frame
frameBox:
  type: floating
  title: 'AlgoScales Workflow Demo'
  style:
    border: '1px solid #ffffff22'
    borderRadius: 4px
    boxShadow: '0 5px 10px rgba(0,0,0,0.4)'

# Add a watermark image to the rendered gif
watermark:
  imagePath: null
  style:
    position: absolute
    right: 15px
    bottom: 15px
    width: 100px
    opacity: 0.9

# Cursor style can be one of
cursorStyle: block

# Font family
fontFamily: "Monaco, Lucida Console, Ubuntu Mono, monospace"

# The size of the font
fontSize: 12

# The height of lines
lineHeight: 1

# The spacing between letters
letterSpacing: 0

# Theme
theme:
  background: "transparent"
  foreground: "#ffffff"
  cursor: "#ffffff"
  black: "#000000"
  red: "#ff6b6b"
  green: "#51cf66"
  yellow: "#ffd93d"
  blue: "#74c0fc"
  magenta: "#da77f2"
  cyan: "#4dabf7"
  white: "#ffffff"
  brightBlack: "#495057"
  brightRed: "#ff7979"
  brightGreen: "#6bcf7f"
  brightYellow: "#ffeaa7"
  brightBlue: "#74b9ff"
  brightMagenta: "#fd79a8"
  brightCyan: "#0984e3"
  brightWhite: "#ffffff"
EOF

echo "✅ Terminalizer config created"
echo ""
echo "🎥 To record the demo:"
echo "   terminalizer record algoscales-demo --config terminalizer-config.yml"
echo ""
echo "🎨 To render as GIF:"
echo "   terminalizer render algoscales-demo"
echo ""
echo "⚙️  Manual recording steps:"
echo "1. Run: terminalizer record algoscales-demo --config terminalizer-config.yml"
echo "2. When recording starts, run: DEMO_SPEED=fast ./full-workflow-demo.sh"
echo "3. Press Ctrl+C when done"
echo "4. Run: terminalizer render algoscales-demo"