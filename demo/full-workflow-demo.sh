#!/bin/bash

# AlgoScales Full Workflow Demo
# Simulates a complete interactive learning session with realistic timing

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# Demo settings
DEMO_SPEED=${DEMO_SPEED:-"normal"}  # fast, normal, slow
WORKSPACE_DIR="$HOME/AlgoScalesDemo"

# Timing functions
case $DEMO_SPEED in
    "fast")
        SHORT_PAUSE=0.5
        MEDIUM_PAUSE=1
        LONG_PAUSE=2
        TYPING_DELAY=0.01
        ;;
    "slow")
        SHORT_PAUSE=2
        MEDIUM_PAUSE=4
        LONG_PAUSE=6
        TYPING_DELAY=0.1
        ;;
    *)  # normal
        SHORT_PAUSE=1
        MEDIUM_PAUSE=2
        LONG_PAUSE=3
        TYPING_DELAY=0.03
        ;;
esac

# Utility functions
typewriter() {
    local text="$1"
    local delay=${2:-$TYPING_DELAY}
    
    for (( i=0; i<${#text}; i++ )); do
        printf "%s" "${text:$i:1}"
        sleep "$delay"
    done
    echo
}

print_header() {
    clear
    echo -e "${BLUE}╭────────────────────────────────────────────────────────────────╮${NC}"
    echo -e "${BLUE}│${NC}${BOLD}${CYAN}           🎵 AlgoScales Full Workflow Demo 🎵                ${NC}${BLUE}│${NC}"
    echo -e "${BLUE}╰────────────────────────────────────────────────────────────────╯${NC}"
    echo ""
}

print_section() {
    echo ""
    echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${YELLOW}  $1${NC}"
    echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
    sleep $MEDIUM_PAUSE
}

simulate_user_input() {
    echo -e "${CYAN}User:${NC} $1"
    sleep $SHORT_PAUSE
}

simulate_system_output() {
    echo -e "${GREEN}System:${NC} $1"
    sleep $SHORT_PAUSE
}

simulate_typing() {
    echo -e "${CYAN}Typing code...${NC}"
    typewriter "$1" 0.02
    sleep $MEDIUM_PAUSE
}

run_real_command() {
    echo -e "${BLUE}$${NC} ${BOLD}$1${NC}"
    sleep $SHORT_PAUSE
    eval "$1" 2>/dev/null || echo "(Command output simulated for demo)"
    sleep $MEDIUM_PAUSE
}

setup_demo_workspace() {
    # Create demo workspace
    mkdir -p "$WORKSPACE_DIR"
    cd "$WORKSPACE_DIR"
    
    # Create sample solutions
    mkdir -p demo_solutions
    
    # Two Sum solution
    cat > demo_solutions/two_sum.go << 'EOF'
package main

import "fmt"

func twoSum(nums []int, target int) []int {
    hashMap := make(map[int]int)
    
    for i, num := range nums {
        complement := target - num
        if j, exists := hashMap[complement]; exists {
            return []int{j, i}
        }
        hashMap[num] = i
    }
    
    return []int{}
}

func main() {
    nums := []int{2, 7, 11, 15}
    target := 9
    result := twoSum(nums, target)
    fmt.Printf("Output: %v\n", result)
}
EOF

    # Sliding window solution
    cat > demo_solutions/sliding_window.go << 'EOF'
package main

import "fmt"

func maxSumSubarray(nums []int, k int) int {
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += nums[i]
    }
    
    maxSum := windowSum
    for i := k; i < len(nums); i++ {
        windowSum = windowSum - nums[i-k] + nums[i]
        if windowSum > maxSum {
            maxSum = windowSum
        }
    }
    
    return maxSum
}

func main() {
    nums := []int{2, 1, 5, 1, 3, 2}
    k := 3
    result := maxSumSubarray(nums, k)
    fmt.Printf("Max sum: %d\n", result)
}
EOF
}

cleanup_demo() {
    echo ""
    echo -e "${YELLOW}Cleaning up demo workspace...${NC}"
    rm -rf "$WORKSPACE_DIR"
    echo -e "${GREEN}Demo complete!${NC}"
}

simulate_full_workflow() {
    print_header
    
    echo -e "${BOLD}🎵 Welcome to AlgoScales!${NC}"
    echo ""
    echo "This demo simulates a complete learning session with realistic timing."
    echo "You'll see exactly how a user would interact with AlgoScales."
    echo ""
    echo -e "${YELLOW}Demo Speed: $DEMO_SPEED (set DEMO_SPEED=fast|normal|slow to adjust)${NC}"
    echo ""
    
    simulate_user_input "Let me try this AlgoScales tool I heard about..."
    sleep $LONG_PAUSE
    
    # Phase 1: Discovery
    print_section "Phase 1: Discovering Algorithm Patterns"
    
    simulate_user_input "What patterns are available?"
    run_real_command "algo-scales list patterns"
    
    simulate_system_output "Wow! 11 fundamental patterns organized like musical scales"
    
    simulate_user_input "Let me see what hash-map problems look like"
    run_real_command "algo-scales list hash-map"
    
    simulate_system_output "I see problems organized by difficulty and pattern"
    sleep $LONG_PAUSE
    
    # Phase 2: Learning Mode
    print_section "Phase 2: Learning Mode - Understanding Patterns"
    
    simulate_user_input "I'm new to this, let me start with learn mode"
    simulate_user_input "algo-scales start learn two_sum"
    
    echo ""
    echo -e "${CYAN}AlgoScales starts learn mode...${NC}"
    sleep $MEDIUM_PAUSE
    
    # Simulate the learn mode interface with realistic problem display
    echo "╭─────────────────────────────────────────────────────────────────╮"
    echo "│                    🎵 Learn Mode: Two Sum 🎵                    │"
    echo "╰─────────────────────────────────────────────────────────────────╯"
    echo ""
    echo "📚 Problem: Two Sum (Hash Map Pattern - A Major Scale)"
    echo "🎯 Difficulty: Easy"
    echo "⏱️  Estimated Time: 15-30 minutes"
    echo "🏢 Companies: Amazon, Facebook, Apple"
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "PROBLEM STATEMENT"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "Given an array of integers and a target sum, find a pair of numbers"
    echo "that add up to the target and return their indices."
    echo ""
    echo "Example 1:"
    echo "  Input: nums = [2, 7, 11, 15], target = 9"
    echo "  Output: [0, 1]"
    echo "  Explanation: nums[0] + nums[1] = 2 + 7 = 9"
    echo ""
    echo "Example 2:"
    echo "  Input: nums = [3, 2, 4], target = 6"
    echo "  Output: [1, 2]"
    echo ""
    echo "Constraints:"
    echo "  • 2 <= nums.length <= 10^4"
    echo "  • -10^9 <= nums[i] <= 10^9"
    echo "  • Only one valid answer exists"
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "PATTERN EXPLANATION"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "💡 Hash Map Pattern (A Major Scale):"
    echo "The Hash Map pattern excels at lookup operations. Instead of checking"
    echo "every pair (O(n²)), we store complements in a hash map for O(1) lookup."
    echo ""
    echo "🔑 Key Insight: For each number, check if its complement exists in our map"
    echo ""
    echo "🎯 When to use:"
    echo "  • Need fast lookups (O(1) average case)"
    echo "  • Finding pairs, triplets with target sums"
    echo "  • Counting frequencies"
    echo "  • Detecting duplicates"
    echo ""
    
    sleep $LONG_PAUSE
    
    simulate_user_input "This explanation is really helpful! Let me see how users actually work on this"
    
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "WORKSPACE AND SOLUTION INTERFACE"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "AlgoScales creates a workspace file where you code your solution:"
    echo ""
    echo "📁 Workspace: ~/AlgoScalesPractice/two_sum.go"
    echo ""
    echo "```go"
    echo "package main"
    echo ""
    echo "import \"fmt\""
    echo ""
    echo "// Two Sum"
    echo "// Pattern: Hash Map (A Major Scale)"
    echo "// Difficulty: Easy"
    echo "//"
    echo "// Given an array of integers and a target sum, find a pair of numbers"
    echo "// that add up to the target and return their indices."
    echo "//"
    echo "// Example: nums = [2, 7, 11, 15], target = 9 -> [0, 1]"
    echo "//"
    echo "// Hint: Use a hash map to store complements for O(1) lookup"
    echo ""
    echo "func twoSum(nums []int, target int) []int {"
    echo "    // Your solution here"
    echo "    return nil"
    echo "}"
    echo ""
    echo "func main() {"
    echo "    // Test cases"
    echo "    nums := []int{2, 7, 11, 15}"
    echo "    target := 9"
    echo "    result := twoSum(nums, target)"
    echo "    fmt.Printf(\"Result: %v\\n\", result)"
    echo "}"
    echo "```"
    echo ""
    
    simulate_user_input "Perfect! I can see the problem details, hints, and starter code all in one place"
    sleep $MEDIUM_PAUSE
    
    echo ""
    echo "💡 Learn Mode Solution Walkthrough:"
    echo "1. Create a hash map to store number -> index mappings"
    echo "2. For each number, calculate its complement (target - number)"
    echo "3. Check if complement exists in hash map"
    echo "4. If yes: return indices, if no: add current number to map"
    echo ""
    echo "Complete Solution (shown in learn mode):"
    echo "```go"
    echo "func twoSum(nums []int, target int) []int {"
    echo "    hashMap := make(map[int]int)"
    echo "    "
    echo "    for i, num := range nums {"
    echo "        complement := target - num"
    echo "        if j, exists := hashMap[complement]; exists {"
    echo "            return []int{j, i}"
    echo "        }"
    echo "        hashMap[num] = i"
    echo "    }"
    echo "    "
    echo "    return []int{}"
    echo "}"
    echo "```"
    echo ""
    
    simulate_user_input "I can see the full solution here, this helps me understand the pattern"
    
    echo ""
    echo "🎮 Interactive Session Commands:"
    echo "  • 'e' - Open your editor to work on the solution"
    echo "  • 'h' - Get progressive hints (3 levels available)"
    echo "  • 's' - Show solution (learn mode only)"
    echo "  • Enter - Submit solution for testing"
    echo "  • 'n' - Skip to next problem"
    echo "  • 'q' - Quit session"
    echo "  • '?' - Show help"
    echo ""
    
    simulate_user_input "These commands make it easy to navigate while staying focused on coding"
    sleep $LONG_PAUSE
    
    # Phase 3: Practice Mode  
    print_section "Phase 3: Practice Mode - Applying the Pattern"
    
    simulate_user_input "Now I want to practice! Let me try a sliding window problem"
    simulate_user_input "algo-scales start practice max_sum_subarray"
    
    echo ""
    echo -e "${CYAN}AlgoScales starts practice mode...${NC}"
    sleep $MEDIUM_PAUSE
    
    echo "╭─────────────────────────────────────────────────────────────────╮"
    echo "│                🎵 Practice Mode: Maximum Sum Subarray 🎵        │"
    echo "╰─────────────────────────────────────────────────────────────────╯"
    echo ""
    echo "📚 Problem: Maximum Sum Subarray of Size K"
    echo "🎯 Pattern: Sliding Window (C Major Scale)"
    echo "⏱️  Timer: 30:00 (started)"
    echo ""
    echo "Find the maximum sum of any contiguous subarray of size k."
    echo ""
    echo "Example: nums = [2, 1, 5, 1, 3, 2], k = 3"
    echo "Output: 9 (subarray [5, 1, 3])"
    echo ""
    echo "🚫 Solution hidden in practice mode"
    echo "💡 Hint available (press 'h')"
    echo ""
    
    sleep $LONG_PAUSE
    
    simulate_user_input "Hmm, let me think about this sliding window pattern..."
    sleep $MEDIUM_PAUSE
    
    simulate_user_input "I'll get a hint to make sure I'm on the right track"
    simulate_user_input "Pressing 'h' for hint..."
    
    echo ""
    echo "💡 Hint Level 1:"
    echo "Think about maintaining a 'window' of exactly k elements."
    echo "Instead of recalculating the sum each time, can you slide"
    echo "the window by removing the leftmost element and adding"
    echo "the new rightmost element?"
    echo ""
    
    simulate_user_input "Ah! So I maintain a running sum and slide the window!"
    sleep $MEDIUM_PAUSE
    
    simulate_user_input "Let me open my editor to work on this"
    simulate_user_input "Pressing 'e' to edit..."
    
    echo ""
    echo -e "${CYAN}Opening editor with problem workspace...${NC}"
    sleep $SHORT_PAUSE
    
    echo ""
    echo "📁 File: ~/AlgoScalesPractice/max_sum_subarray.go"
    echo ""
    echo "```go"
    echo "package main"
    echo ""
    echo "import \"fmt\""
    echo ""
    echo "// Maximum Sum Subarray of Size K"
    echo "// Pattern: Sliding Window (C Major Scale)"
    echo "// Difficulty: Easy"
    echo "//"
    echo "// Find the maximum sum of any contiguous subarray of size k."
    echo "//"
    echo "// Example: nums = [2, 1, 5, 1, 3, 2], k = 3 -> 9"
    echo "//"
    echo "// Constraints: 1 <= k <= nums.length <= 10^5"
    echo ""
    echo "func maxSumSubarray(nums []int, k int) int {"
    echo "    // Your solution here"
    echo "    return 0"
    echo "}"
    echo ""
    echo "func main() {"
    echo "    nums := []int{2, 1, 5, 1, 3, 2}"
    echo "    k := 3"
    echo "    result := maxSumSubarray(nums, k)"
    echo "    fmt.Printf(\"Max sum: %d\\n\", result)"
    echo "}"
    echo "```"
    echo ""
    
    simulate_user_input "Great! I can see the problem setup with examples and test cases"
    sleep $MEDIUM_PAUSE
    
    simulate_user_input "Now I'll implement the sliding window solution..."
    echo ""
    echo -e "${CYAN}User starts coding...${NC}"
    sleep $SHORT_PAUSE
    
    # Simulate typing the solution
    echo "```go"
    simulate_typing "func maxSumSubarray(nums []int, k int) int {"
    simulate_typing "    // Calculate sum of first window"
    simulate_typing "    windowSum := 0"
    simulate_typing "    for i := 0; i < k; i++ {"
    simulate_typing "        windowSum += nums[i]"
    simulate_typing "    }"
    echo ""
    simulate_typing "    maxSum := windowSum"
    echo ""
    simulate_typing "    // Slide the window"
    simulate_typing "    for i := k; i < len(nums); i++ {"
    simulate_typing "        windowSum = windowSum - nums[i-k] + nums[i]"
    simulate_typing "        if windowSum > maxSum {"
    simulate_typing "            maxSum = windowSum"
    simulate_typing "        }"
    simulate_typing "    }"
    echo ""
    simulate_typing "    return maxSum"
    simulate_typing "}"
    echo "```"
    
    sleep $LONG_PAUSE
    
    simulate_user_input "Solution complete! Let me submit it"
    simulate_user_input "Pressing Enter to submit..."
    
    echo ""
    echo -e "${GREEN}🧪 Running tests...${NC}"
    sleep $MEDIUM_PAUSE
    
    echo "Test 1: [2, 1, 5, 1, 3, 2], k=3 -> Expected: 9, Got: 9 ✅"
    echo "Test 2: [2, 3, 4, 1, 5], k=2 -> Expected: 7, Got: 7 ✅"
    echo "Test 3: [1, 4, 2, 9, 3], k=3 -> Expected: 15, Got: 15 ✅"
    echo ""
    echo -e "${GREEN}🎉 All tests passed! Solution accepted!${NC}"
    echo "⏱️  Solve time: 08:42"
    echo "🎯 Pattern mastery: Sliding Window +1"
    
    sleep $LONG_PAUSE
    
    # Phase 4: AI Integration
    print_section "Phase 4: AI-Powered Learning"
    
    simulate_user_input "Let me try the AI features for deeper understanding"
    simulate_user_input "algo-scales ai review"
    
    echo ""
    echo -e "${CYAN}🤖 AI Code Review:${NC}"
    sleep $MEDIUM_PAUSE
    
    echo "Your sliding window solution is excellent! Here's my analysis:"
    echo ""
    echo "✅ Strengths:"
    echo "• Optimal O(n) time complexity"
    echo "• Space efficient O(1) solution"
    echo "• Clean, readable variable names"
    echo "• Proper sliding window technique"
    echo ""
    echo "💡 Learning Notes:"
    echo "• You've mastered the core sliding window pattern"
    echo "• This technique applies to many problems with fixed-size subarrays"
    echo "• Consider variable-size sliding windows as your next challenge"
    echo ""
    echo "🎯 Pattern Recommendation:"
    echo "Try 'Two Pointers' next - it builds on similar concepts!"
    
    sleep $LONG_PAUSE
    
    simulate_user_input "Wow, personalized feedback! Let me ask the AI about patterns"
    simulate_user_input "algo-scales ai chat"
    
    echo ""
    echo -e "${CYAN}🤖 AI Chat Mode:${NC}"
    echo ""
    
    simulate_user_input "What's the difference between sliding window and two pointers?"
    
    echo ""
    echo "🤖 AI: Great question! Both techniques use pointers, but:"
    echo ""
    echo "**Sliding Window (Fixed Size):**"
    echo "• Maintains a window of constant size k"
    echo "• Both pointers move together (window slides)"
    echo "• Best for: max/min of fixed-size subarrays"
    echo ""
    echo "**Two Pointers (Variable):**"
    echo "• Pointers move independently based on conditions"
    echo "• Window size changes dynamically"
    echo "• Best for: sorted arrays, finding pairs/triplets"
    echo ""
    echo "Think of sliding window as a 'rigid frame' and two pointers"
    echo "as 'flexible boundaries' that adapt to find optimal solutions!"
    
    sleep $LONG_PAUSE
    
    # Phase 5: Daily Practice
    print_section "Phase 5: Daily Practice System"
    
    simulate_user_input "I want to build consistent habits. Let me try daily practice"
    simulate_user_input "algo-scales daily"
    
    echo ""
    echo "╭─────────────────────────────────────────────────────────────────╮"
    echo "│                 🎵 AlgoScales Daily Practice 🎵                 │"
    echo "╰─────────────────────────────────────────────────────────────────╯"
    echo ""
    echo "Practice one problem from each algorithm pattern (scale) to build skills."
    echo "Just as musicians practice scales daily, this builds your fundamentals."
    echo ""
    echo "Current streak: 1 day 🔥"
    echo "Patterns completed today: 2/11"
    echo "Patterns remaining: 9"
    echo ""
    echo "Now practicing: G Major (Two Pointers)"
    echo "Description: Balanced and efficient, the workhorse of array manipulation"
    echo ""
    echo "Problem: Pair with Target Sum (easy)"
    echo ""
    
    simulate_user_input "Perfect! This ensures I practice all patterns systematically"
    sleep $LONG_PAUSE
    
    # Phase 6: Progress Tracking
    print_section "Phase 6: Progress and Statistics"
    
    simulate_user_input "Let me check my overall progress"
    run_real_command "algo-scales stats"
    
    simulate_system_output "I can see my progress across all patterns and my improvement over time!"
    
    echo ""
    echo "📊 Simulated Progress After One Week:"
    echo ""
    echo "Overall Statistics:"
    echo "Total Problems Attempted: 23"
    echo "Total Problems Solved: 18"
    echo "Success Rate: 78%"
    echo "Average Solve Time: 12:34"
    echo "Current Streak: 7 days 🔥🔥🔥🔥🔥🔥🔥"
    echo ""
    echo "Pattern Mastery:"
    echo "  🎵 C Major (Sliding Window): ████████░░ 80% (4/5)"
    echo "  🎵 G Major (Two Pointers): ██████████ 100% (3/3)"
    echo "  🎵 A Major (Hash Map): ██████████ 100% (2/2)"
    echo "  🎵 E Major (Binary Search): ██████░░░░ 60% (3/5)"
    echo "  🎵 D Major (Fast/Slow Pointers): ████░░░░░░ 40% (2/5)"
    echo ""
    echo "💡 AI Recommendation: Focus on Binary Search patterns next!"
    
    sleep $LONG_PAUSE
    
    # Phase 7: Conclusion
    print_section "Complete Learning Journey Summary"
    
    echo -e "${GREEN}🎉 You've experienced the complete AlgoScales workflow!${NC}"
    echo ""
    echo "What you accomplished in this demo:"
    echo ""
    echo "✅ **Pattern Discovery**: Explored 11 fundamental algorithm patterns"
    echo "✅ **Learn Mode**: Understood hash map pattern with full explanations"
    echo "✅ **Practice Mode**: Solved sliding window problem independently"
    echo "✅ **AI Integration**: Got personalized code review and pattern advice"
    echo "✅ **Daily Practice**: Established systematic learning routine"
    echo "✅ **Progress Tracking**: Monitored improvement and mastery"
    echo ""
    echo -e "${CYAN}🎵 Key AlgoScales Advantages:${NC}"
    echo ""
    echo "🎯 **Pattern-Focused**: Learn transferable techniques, not just solutions"
    echo "🤖 **AI-Powered**: Personalized hints and explanations"
    echo "⚡ **Efficient**: No grinding - smart, focused practice"
    echo "📈 **Measurable**: Clear progress tracking and mastery metrics"
    echo "🎼 **Musical**: Memorable pattern organization"
    echo ""
    echo -e "${YELLOW}Built to help developers master interview fundamentals efficiently!${NC}"
    echo ""
    echo -e "${BOLD}Ready to start your algorithm mastery journey?${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Install: make install-user"
    echo "2. Configure AI: algo-scales ai config"
    echo "3. Start learning: algo-scales start learn"
    echo "4. Build habits: algo-scales daily"
    echo ""
    
    sleep $LONG_PAUSE
}

# Handle interruption gracefully
trap 'echo -e "\n${YELLOW}Demo interrupted.${NC}"; cleanup_demo; exit 0' INT

# Main execution
main() {
    if ! command -v algo-scales >/dev/null 2>&1; then
        echo -e "${RED}❌ algo-scales not found. Please install first:${NC}"
        echo "   make install-user"
        exit 1
    fi
    
    setup_demo_workspace
    simulate_full_workflow
    cleanup_demo
}

echo -e "${BOLD}🎵 AlgoScales Full Workflow Demo${NC}"
echo ""
echo "This demo simulates a complete learning session showing:"
echo "• Real user interactions and timing"
echo "• Full problem-solving workflow"
echo "• AI-powered learning assistance"
echo "• Progress tracking and habit building"
echo ""
echo "Estimated time: 8-12 minutes"
echo ""
echo -e "${CYAN}Press Enter to start the demo (Ctrl+C to exit)${NC}"
read -r

main "$@"