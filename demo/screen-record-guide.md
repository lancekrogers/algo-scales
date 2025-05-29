# Screen Recording Guide for AlgoScales Demo

If automated tools don't work, you can create a GIF using screen recording:

## macOS (Built-in)

### Step 1: Prepare Terminal
```bash
# Set optimal terminal size
# Resize terminal to ~120x30 characters
# Use a clean theme (dark background recommended)

# Set fast demo speed
export DEMO_SPEED=fast
```

### Step 2: Record Screen
```bash
# Use QuickTime Player or built-in screen recording
# Cmd+Shift+5 -> Select area -> Record

# Or use ffmpeg if installed:
ffmpeg -f avfoundation -i "1" -r 30 -t 600 algoscales-demo.mov
```

### Step 3: Run Demo
```bash
./full-workflow-demo.sh
```

### Step 4: Convert to GIF
```bash
# Using ffmpeg (install with: brew install ffmpeg)
ffmpeg -i algoscales-demo.mov \
  -vf "fps=10,scale=1200:-1:flags=lanczos,palettegen" \
  palette.png

ffmpeg -i algoscales-demo.mov -i palette.png \
  -vf "fps=10,scale=1200:-1:flags=lanczos" \
  -loop 0 algoscales-demo.gif

# Or use online converters:
# - ezgif.com
# - convertio.co
# - cloudconvert.com
```

## Linux

### Using `byzanz` (Ubuntu/Debian)
```bash
# Install
sudo apt install byzanz

# Record (adjust coordinates for your terminal window)
byzanz-record --duration=600 --x=100 --y=100 --width=1200 --height=800 algoscales-demo.gif

# Then run demo in terminal
./full-workflow-demo.sh
```

### Using `peek`
```bash
# Install
sudo apt install peek

# Run peek, select terminal window, start recording
# Run demo, stop recording when done
```

## Windows

### Using Windows Game Bar
1. Win+G to open Game Bar
2. Click record button
3. Run demo in terminal
4. Stop recording
5. Convert .mp4 to .gif using online tools

### Using OBS Studio
1. Install OBS Studio
2. Create scene with window capture of terminal
3. Start recording
4. Run demo
5. Convert output to GIF

## Optimization Tips

### Terminal Setup
- Use 120x30 character window size
- Dark theme with high contrast
- Increase font size to 14-16pt for readability
- Clear terminal before recording

### Demo Settings
```bash
# Fast demo for shorter GIF
export DEMO_SPEED=fast

# Or create custom speed
export DEMO_SPEED=recording
# Then edit full-workflow-demo.sh to add even faster timing
```

### GIF Optimization
- Target 10-15 FPS for good balance
- Scale to max 1200px width
- Use palette optimization
- Compress with tools like gifsicle

### File Size Management
```bash
# Reduce colors and optimize
gifsicle -O3 --colors 64 algoscales-demo.gif -o algoscales-demo-optimized.gif

# Create multiple versions
ffmpeg -i demo.mov -vf "scale=800:-1" small-demo.gif    # Small for README
ffmpeg -i demo.mov -vf "scale=1200:-1" large-demo.gif   # Large for website
```

## Recommended Final Specs

- **Duration**: 8-12 minutes (full demo) or 2-3 minutes (condensed)
- **Size**: 1200x800px (large) or 800x600px (compact)
- **FPS**: 10-15 fps
- **Colors**: 64-128 colors
- **File size**: Target <10MB for GitHub, <5MB for social media
- **Format**: GIF or WebM for modern browsers

## Quick Commands

### macOS One-liner
```bash
# Record, run demo, convert (requires ffmpeg)
./record-demo-gif.sh
```

### Manual Recording
```bash
# 1. Start recording
# 2. Run this:
DEMO_SPEED=fast ./full-workflow-demo.sh

# 3. Stop recording and convert
```

The automated `record-demo-gif.sh` script is the best option if the dependencies install properly!