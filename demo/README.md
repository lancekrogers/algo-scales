# AlgoScales Demo Collection

This directory contains various demos to showcase AlgoScales functionality without requiring manual coding.

## 🌟 Recommended: Full Workflow Demo

**`./full-workflow-demo.sh`** - Experience the complete AlgoScales learning journey (8-12 minutes)

Shows the realistic user experience including:

- **Problem Display**: How problems are formatted with examples, constraints, and pattern explanations
- **Workspace Setup**: How AlgoScales creates files with starter code and problem context
- **Interactive Commands**: The session controls users have (`e` for editor, `h` for hints, etc.)
- **AI Integration**: Real AI assistance, code review, and personalized guidance
- **Progress Tracking**: Statistics, streaks, and pattern mastery progression
- **Complete Workflow**: From problem discovery to solution submission

## Available Demos

### Core Demos

1. **Quick Demo** (`./quick-demo.sh`) - 2-3 minutes

   - Core commands and features
   - Pattern organization
   - AI integration overview
   - Perfect for quick evaluation

2. **Full Workflow Demo** (`./full-workflow-demo.sh`) - 8-12 minutes

   - Complete realistic user experience
   - All AI features integrated
   - Problem solving from start to finish
   - Shows progress tracking and daily practice

3. **Interactive Demo** (`./demo.sh`) - 15-20 minutes
   - Hands-on exploration mode
   - Pre-written solutions for demonstration
   - Detailed feature walkthrough
   - Advanced AI capabilities

### Tools & Utilities

**GIF Recording** (`./record-demo-gif.sh`)

- Automatically records any demo as a GIF
- Creates both full and social media versions
- Requires dependencies:

  ```bash
  # macOS
  brew install asciinema
  cargo install --git https://github.com/asciinema/agg

  # Linux
  pip install asciinema
  cargo install --git https://github.com/asciinema/agg
  ```

**Demo Testing** (`./test-demos.sh`)

- Validates all demo scripts work correctly
- Useful before recording or presenting

**Manual Recording Guide** (`screen-record-guide.md`)

- Fallback instructions for screen recording
- Platform-specific guides (macOS, Linux, Windows)
- GIF optimization tips

## Demo Configuration

All demos support speed adjustment:

```bash
# Fast playback
DEMO_SPEED=fast ./full-workflow-demo.sh

# Slow, detailed viewing
DEMO_SPEED=slow ./full-workflow-demo.sh

# Default timing
./full-workflow-demo.sh
```

## Key Demo Features

### Realistic User Experience

- Authentic problem-solving workflow
- Realistic timing between interactions
- Shows actual user thought process
- Demonstrates real editor integration

### No Coding Required

- Pre-written solutions are "typed" automatically
- Simulated test results show success scenarios
- AI interactions use realistic responses
- Progress tracking shows believable advancement

### Professional Presentation

- Color-coded output for different interaction types
- Clear section headers and progress indicators
- Proper terminal formatting and sizing
- Graceful interruption handling (Ctrl+C)

