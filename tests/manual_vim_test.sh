#!/bin/bash

# Manual test script for vim mode commands

echo "Testing AlgoScales vim mode commands..."
echo "======================================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test function
test_command() {
    local description=$1
    local command=$2
    local expected_field=$3
    
    echo -n "Testing: $description... "
    
    output=$(eval "$command" 2>/dev/null)
    
    if echo "$output" | jq -e ".$expected_field" > /dev/null 2>&1; then
        echo -e "${GREEN}PASS${NC}"
        return 0
    else
        echo -e "${RED}FAIL${NC}"
        echo "Output: $output"
        return 1
    fi
}

# Build the binary
echo "Building algo-scales..."
go build -o bin/algo-scales ./

# Run tests
echo ""
echo "Running vim mode tests:"
echo "----------------------"

test_command "List problems" \
    "./bin/algo-scales list --vim-mode" \
    "problems"

test_command "Start session" \
    "./bin/algo-scales start learn two_sum --language go --vim-mode" \
    "id"

test_command "Get hint (level 1)" \
    "./bin/algo-scales hint --problem-id pair_with_target_sum --language go --vim-mode" \
    "hint"

test_command "Get solution" \
    "./bin/algo-scales solution --problem-id two_sum --language go --vim-mode" \
    "solution"

test_command "AI hint" \
    "./bin/algo-scales ai-hint --problem-id two_sum --language go --vim-mode" \
    "ready"

# Test multi-level hints with a simple script
echo ""
echo "Testing multi-level hints:"
echo "-------------------------"

cat > /tmp/test_hints.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
)

type HintResponse struct {
    Hint        string   `json:"hint"`
    Level       int      `json:"level"`
    Walkthrough []string `json:"walkthrough,omitempty"`
    Solution    string   `json:"solution,omitempty"`
}

func main() {
    // Note: In real usage, hint levels would be tracked within a session
    // For CLI, each call is independent
    cmd := exec.Command("./bin/algo-scales", "hint", 
        "--problem-id", "pair_with_target_sum",
        "--language", "go",
        "--vim-mode")
    
    output, err := cmd.Output()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    var resp HintResponse
    if err := json.Unmarshal(output, &resp); err != nil {
        fmt.Printf("Parse error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Hint Level: %d\n", resp.Level)
    fmt.Printf("Has hint: %v\n", resp.Hint != "")
    fmt.Printf("Has walkthrough: %v\n", len(resp.Walkthrough) > 0)
    fmt.Printf("Has solution: %v\n", resp.Solution != "")
}
EOF

go run /tmp/test_hints.go

# Test submit with a valid solution
echo ""
echo "Testing solution submission:"
echo "---------------------------"

cat > /tmp/solution.go << 'EOF'
func twoSum(nums []int, target int) []int {
    seen := make(map[int]int)
    for i, num := range nums {
        if j, ok := seen[target-num]; ok {
            return []int{j, i}
        }
        seen[num] = i
    }
    return nil
}
EOF

test_command "Submit solution" \
    "./bin/algo-scales submit --problem-id two_sum --language go --file /tmp/solution.go --vim-mode" \
    "test_results"

# Clean up
rm -f /tmp/test_hints.go /tmp/solution.go

echo ""
echo "Testing complete!"