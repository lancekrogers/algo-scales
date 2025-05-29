#!/bin/bash

# AlgoScales Interactive Demo
# This script demonstrates the full AlgoScales workflow without requiring manual coding

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Demo solutions for automatic submission
DEMO_SOLUTIONS_DIR="demo_solutions"

# Utility functions
print_header() {
    echo ""
    echo -e "${BLUE}╭────────────────────────────────────────────────────────────────╮${NC}"
    echo -e "${BLUE}│${NC}${BOLD}${CYAN}                    🎵 AlgoScales Demo 🎵                     ${NC}${BLUE}│${NC}"
    echo -e "${BLUE}╰────────────────────────────────────────────────────────────────╯${NC}"
    echo ""
}

print_section() {
    echo ""
    echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${YELLOW}  $1${NC}"
    echo -e "${YELLOW}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
}

print_step() {
    echo -e "${GREEN}➤${NC} ${BOLD}$1${NC}"
    echo ""
}

wait_for_user() {
    echo -e "${CYAN}Press Enter to continue...${NC}"
    read -r
}

run_command() {
    echo -e "${BLUE}$${NC} ${BOLD}$1${NC}"
    echo ""
    eval "$1"
    echo ""
}

create_demo_solutions() {
    print_step "Setting up demo solutions..."
    
    mkdir -p "$DEMO_SOLUTIONS_DIR"
    
    # Two Sum solution (Go)
    cat > "$DEMO_SOLUTIONS_DIR/two_sum.go" << 'EOF'
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
    // Test the solution
    nums := []int{2, 7, 11, 15}
    target := 9
    result := twoSum(nums, target)
    fmt.Printf("Input: nums = %v, target = %d\n", nums, target)
    fmt.Printf("Output: %v\n", result)
}
EOF

    # Sliding Window solution (Go)
    cat > "$DEMO_SOLUTIONS_DIR/max_sum_subarray.go" << 'EOF'
package main

import "fmt"

func maxSumSubarray(nums []int, k int) int {
    if len(nums) < k {
        return 0
    }
    
    // Calculate sum of first window
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += nums[i]
    }
    
    maxSum := windowSum
    
    // Slide the window
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
    fmt.Printf("Input: nums = %v, k = %d\n", nums, k)
    fmt.Printf("Max sum subarray of size %d: %d\n", k, result)
}
EOF

    # Binary Search solution (Go)
    cat > "$DEMO_SOLUTIONS_DIR/search_in_rotated_array.go" << 'EOF'
package main

import "fmt"

func search(nums []int, target int) int {
    left, right := 0, len(nums)-1
    
    for left <= right {
        mid := left + (right-left)/2
        
        if nums[mid] == target {
            return mid
        }
        
        // Left half is sorted
        if nums[left] <= nums[mid] {
            if target >= nums[left] && target < nums[mid] {
                right = mid - 1
            } else {
                left = mid + 1
            }
        } else { // Right half is sorted
            if target > nums[mid] && target <= nums[right] {
                left = mid + 1
            } else {
                right = mid - 1
            }
        }
    }
    
    return -1
}

func main() {
    nums := []int{4, 5, 6, 7, 0, 1, 2}
    target := 0
    result := search(nums, target)
    fmt.Printf("Input: nums = %v, target = %d\n", nums, target)
    fmt.Printf("Index of target: %d\n", result)
}
EOF
}

cleanup_demo() {
    print_step "Cleaning up demo files..."
    rm -rf "$DEMO_SOLUTIONS_DIR"
    echo "Demo cleanup complete!"
}

# Trap to cleanup on exit
trap cleanup_demo EXIT

main() {
    print_header
    
    echo -e "${BOLD}Welcome to the AlgoScales Interactive Demo!${NC}"
    echo ""
    echo "This demo will walk you through the complete AlgoScales workflow:"
    echo "• Pattern-based algorithm learning"
    echo "• Different practice modes (learn, practice, daily)"
    echo "• AI-powered assistance"
    echo "• Progress tracking and statistics"
    echo ""
    echo -e "${YELLOW}Note: This demo uses pre-written solutions to show the full workflow${NC}"
    echo -e "${YELLOW}without requiring you to write code manually.${NC}"
    echo ""
    
    wait_for_user
    
    # Check if algo-scales is installed
    if ! command -v algo-scales >/dev/null 2>&1; then
        echo -e "${RED}❌ AlgoScales not found in PATH${NC}"
        echo ""
        echo "Please install AlgoScales first:"
        echo "  make install-user"
        echo "  # or"
        echo "  go install github.com/lancekrogers/algo-scales@latest"
        exit 1
    fi
    
    echo -e "${GREEN}✅ AlgoScales found: $(algo-scales --version 2>/dev/null || echo 'latest')${NC}"
    echo ""
    
    # Setup demo solutions
    create_demo_solutions
    
    # Demo Section 1: Listing Problems and Patterns
    print_section "1. Exploring Algorithm Patterns (Musical Scales)"
    
    print_step "AlgoScales organizes algorithms by patterns, like musical scales"
    echo "Each pattern represents a fundamental problem-solving technique:"
    echo ""
    
    run_command "algo-scales list patterns"
    
    wait_for_user
    
    print_step "Let's see what problems are available for the 'hash-map' pattern:"
    run_command "algo-scales list hash-map"
    
    wait_for_user
    
    # Demo Section 2: Learn Mode
    print_section "2. Learn Mode - Understanding the Pattern"
    
    print_step "Starting with 'learn' mode to understand the Two Sum problem"
    echo "Learn mode provides hints and lets you see the solution immediately."
    echo ""
    
    # Start learn mode in background and simulate interaction
    echo -e "${BLUE}$${NC} ${BOLD}algo-scales start learn two_sum${NC}"
    echo ""
    echo "📚 Problem: Two Sum"
    echo "Given an array of integers and a target sum, find two numbers that add up to the target."
    echo ""
    echo "💡 This demonstrates the Hash Map pattern (A Major scale)"
    echo "   - Use a hash map to store complements"
    echo "   - Check if current number's complement exists"
    echo "   - O(n) time complexity vs O(n²) brute force"
    echo ""
    echo "🎯 Pattern Learning: Hash maps excel at lookup operations"
    echo ""
    
    wait_for_user
    
    # Demo Section 3: Practice Mode with Solution Submission
    print_section "3. Practice Mode - Applying the Pattern"
    
    print_step "Now let's practice with a sliding window problem"
    echo "Practice mode includes a timer and tests your implementation."
    echo ""
    
    echo -e "${BLUE}$${NC} ${BOLD}algo-scales start practice max_sum_subarray${NC}"
    echo ""
    echo "📚 Problem: Maximum Sum Subarray of Size K"
    echo "Find the maximum sum of any contiguous subarray of size k."
    echo ""
    echo "💡 This demonstrates the Sliding Window pattern (C Major scale)"
    echo "   - Maintain a window of fixed size"
    echo "   - Slide the window by removing left element and adding right element"
    echo "   - Track maximum sum seen so far"
    echo ""
    
    # Simulate working on the problem
    echo "⏰ Timer started... (simulated)"
    echo ""
    echo "🔧 Writing solution to practice the sliding window technique..."
    echo ""
    
    # Copy our demo solution to the expected location
    PRACTICE_DIR="$HOME/AlgoScalesPractice"
    mkdir -p "$PRACTICE_DIR"
    cp "$DEMO_SOLUTIONS_DIR/max_sum_subarray.go" "$PRACTICE_DIR/"
    
    echo "✅ Solution written! Let's test it:"
    echo ""
    
    # Show the solution briefly
    echo -e "${CYAN}Preview of sliding window solution:${NC}"
    head -20 "$DEMO_SOLUTIONS_DIR/max_sum_subarray.go"
    echo "..."
    echo ""
    
    wait_for_user
    
    # Demo Section 4: AI Integration
    print_section "4. AI-Powered Learning Assistant"
    
    print_step "AlgoScales includes an AI assistant for personalized help"
    echo "The AI can provide hints, review code, and explain patterns."
    echo ""
    
    # Check if AI is configured
    echo -e "${BLUE}$${NC} ${BOLD}algo-scales ai config${NC}"
    echo ""
    echo "🤖 AI Configuration:"
    echo "   Provider: Claude (via claude-code)"
    echo "   Status: Available for hints and code review"
    echo "   Features: Progressive hints, pattern explanations, code analysis"
    echo ""
    echo "💡 AI capabilities:"
    echo "   • Progressive hints (3 levels of increasing detail)"
    echo "   • Code review with suggestions"
    echo "   • Pattern-specific explanations"
    echo "   • Interactive chat for questions"
    echo ""
    
    wait_for_user
    
    print_step "Example: Getting an AI hint for binary search pattern"
    echo -e "${BLUE}$${NC} ${BOLD}algo-scales ai hint search_in_rotated_array${NC}"
    echo ""
    echo "🤖 AI Hint (Level 1):"
    echo "\"Think about which half of the array is sorted. In a rotated sorted array,"
    echo "at least one half is always sorted normally. Use this property to decide"
    echo "which direction to search.\""
    echo ""
    echo "💡 The AI provides increasingly detailed hints as you need them."
    echo ""
    
    wait_for_user
    
    # Demo Section 5: Daily Scale Practice
    print_section "5. Daily Scale Practice - Building Consistency"
    
    print_step "Daily practice helps you master all patterns systematically"
    echo "Like a musician practicing scales, this builds fundamental skills."
    echo ""
    
    echo -e "${BLUE}$${NC} ${BOLD}algo-scales daily${NC}"
    echo ""
    echo "🎵 Daily Scale Practice"
    echo "Today's practice: C Major (Sliding Window)"
    echo ""
    echo "📊 Your Progress:"
    echo "   Current streak: 3 days 🔥🔥🔥"
    echo "   Patterns completed today: 2/11"
    echo "   Next pattern: G Major (Two Pointers)"
    echo ""
    echo "💡 Daily practice ensures you:"
    echo "   • Stay sharp on all patterns"
    echo "   • Build consistent coding habits"
    echo "   • Progress through all fundamental techniques"
    echo ""
    
    wait_for_user
    
    # Demo Section 6: Statistics and Progress Tracking
    print_section "6. Progress Tracking and Statistics"
    
    print_step "AlgoScales tracks your learning progress across all patterns"
    echo ""
    
    run_command "algo-scales stats"
    
    wait_for_user
    
    # Demo Section 7: Different Languages
    print_section "7. Multi-Language Support"
    
    print_step "Practice the same patterns in different programming languages"
    echo ""
    
    echo -e "${BLUE}$${NC} ${BOLD}algo-scales start practice two_sum --language python${NC}"
    echo ""
    echo "🐍 Python version of Two Sum:"
    echo ""
    cat << 'EOF'
def two_sum(nums, target):
    hash_map = {}
    
    for i, num in enumerate(nums):
        complement = target - num
        if complement in hash_map:
            return [hash_map[complement], i]
        hash_map[num] = i
    
    return []
EOF
    echo ""
    echo "💡 Same pattern, different syntax - reinforces understanding!"
    echo ""
    
    wait_for_user
    
    # Demo Section 8: Advanced Features
    print_section "8. Advanced Features Preview"
    
    print_step "AlgoScales includes powerful features for serious learners"
    echo ""
    
    echo "🚀 Advanced Features:"
    echo ""
    echo "📈 Progress Analytics:"
    echo "   • Pattern mastery tracking"
    echo "   • Difficulty progression"
    echo "   • Time-to-solve metrics"
    echo "   • Streak maintenance"
    echo ""
    echo "🎯 Adaptive Learning:"
    echo "   • AI suggests next problems based on progress"
    echo "   • Difficulty adjustment based on performance"
    echo "   • Weak pattern identification"
    echo ""
    echo "🔗 Integration Options:"
    echo "   • Neovim plugin (in development)"
    echo "   • VS Code extension (planned)"
    echo "   • CLI automation scripts"
    echo ""
    echo "💼 Interview Preparation:"
    echo "   • Company-specific problem sets"
    echo "   • Mock interview mode"
    echo "   • Pattern-based study plans"
    echo ""
    
    wait_for_user
    
    # Demo Conclusion
    print_section "Demo Complete! 🎉"
    
    echo -e "${GREEN}✅ You've seen the complete AlgoScales workflow:${NC}"
    echo ""
    echo "1. ✅ Pattern-based problem organization"
    echo "2. ✅ Multiple learning modes (learn, practice, daily)"
    echo "3. ✅ AI-powered assistance and hints"
    echo "4. ✅ Progress tracking and statistics"
    echo "5. ✅ Multi-language support"
    echo "6. ✅ Consistent daily practice system"
    echo ""
    echo -e "${CYAN}🎵 Ready to start your algorithm learning journey?${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Configure AI assistant: ${BOLD}algo-scales ai config${NC}"
    echo "  2. Start with learn mode: ${BOLD}algo-scales start learn${NC}"
    echo "  3. Try daily practice: ${BOLD}algo-scales daily${NC}"
    echo "  4. Check your progress: ${BOLD}algo-scales stats${NC}"
    echo ""
    echo -e "${YELLOW}💡 Remember: Focus on learning patterns, not memorizing solutions!${NC}"
    echo ""
    
    print_header
    echo -e "${BOLD}${CYAN}Thank you for trying AlgoScales! 🎵${NC}"
    echo ""
}

# Handle Ctrl+C gracefully
trap 'echo -e "\n${YELLOW}Demo interrupted. Cleaning up...${NC}"; cleanup_demo; exit 0' INT

# Run the demo
main "$@"