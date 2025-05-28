# AlgoScales: Daily Scale Practice

## The Musician's Approach to Algorithm Mastery

Just as musicians develop technical proficiency by practicing scales daily, developers can build algorithm intuition through regular practice of fundamental patterns. The Daily Scale Practice feature guides you through this process in a structured, motivating way.

## Getting Started with Daily Scales

To begin your daily scale practice from the command line:

```bash
# Use CLI mode (default):
algo-scales daily

# TUI mode (work in progress - not recommended):
# algo-scales daily --tui
```

This will start a sequence of problem-solving sessions, one for each of the 11 core algorithm patterns ("scales").

## CLI Mode Workflow

In CLI mode (default), AlgoScales creates problem files in a dedicated workspace:

```
~/Dev/AlgoScalesPractice/Daily/{YYYY-MM-DD}/
```

The CLI workflow is:

1. **Start Session**: Run `algo-scales daily` to start or resume a daily session
2. **Solve Problems**: Edit the problem file in your editor of choice
3. **Test Solutions**: Run `algo-scales daily test` to verify your solution
4. **Skip Problems**: Run `algo-scales daily skip` if you want to skip the current problem
5. **Resume Skipped**: Run `algo-scales daily resume-skipped` to return to skipped problems
6. **Check Status**: Run `algo-scales daily status` to view your current progress

Problems are only marked as complete when your solution passes all tests.

### Command-Line Options

```bash
# Set the programming language
algo-scales daily --language python

# Change the timer duration (in minutes)
algo-scales daily --timer 30

# Focus on a specific difficulty
algo-scales daily --difficulty medium
```

## The 11 Core Scales (Algorithm Patterns)

Each "scale" represents a fundamental algorithm pattern that appears frequently in interviews:

1. **C Major (Sliding Window)** - The fundamental scale, elegant and versatile
2. **G Major (Two Pointers)** - Balanced and efficient, the workhorse of array manipulation
3. **D Major (Fast & Slow Pointers)** - The cycle detector, bright and revealing
4. **A Major (Hash Maps)** - The lookup accelerator, crisp and direct
5. **E Major (Binary Search)** - The divider and conqueror, precise and logarithmic
6. **B Major (DFS)** - The deep explorer, rich and thorough
7. **F# Major (BFS)** - The level-by-level discoverer, methodical and complete
8. **Db Major (Dynamic Programming)** - The optimizer, complex and powerful
9. **Ab Major (Greedy)** - The local maximizer, bold and decisive
10. **Eb Major (Union-Find)** - The connector, structured and organized
11. **Bb Major (Heap / Priority Queue)** - The sorter, flexible and maintaining order

## Practice Workflow

1. **Start a Scale**: Begin with the first pattern you haven't practiced today
2. **Solve the Problem**: Work on a problem exemplifying that pattern
3. **Complete the Scale**: When you solve the problem, that scale is marked complete
4. **Progress to Next Scale**: You'll be prompted to continue to the next pattern
5. **Daily Completion**: When you finish all scales, you've completed your daily practice

## Progression System

As you consistently practice your scales, you'll advance through skill levels:

- **Apprentice**: Just beginning your journey
- **Journeyman**: Recognizing and applying patterns with growing confidence
- **Virtuoso**: Quickly identifying and solving pattern-based problems
- **Maestro**: Mastery of all patterns with the ability to combine them effortlessly

## Tracking Your Progress

Your practice sessions are tracked automatically:

- **Daily Streaks**: Consecutive days of practice
- **Pattern Mastery**: Progress in each algorithm pattern
- **Achievement Badges**: Recognition for significant milestones

## Tips for Effective Practice

1. **Consistency Over Quantity**: Better to practice a little each day than cram occasionally
2. **Focus on Understanding**: Don't just memorize solutions, understand the pattern
3. **Verbalize Your Approach**: Explain your solution process out loud
4. **Review Past Solutions**: Periodically review patterns you've already practiced
5. **Apply Variations**: Try solving the same problem with different approaches

## Using with Neovim

If you use Neovim, try our dedicated plugin for an even smoother experience:

```lua
-- Start daily scale practice from within Neovim
:AlgoScalesDailyScale
```

Remember: Just as a pianist must practice scales daily to perform Chopin brilliantly, a developer must practice algorithm patterns to shine in technical interviews.
