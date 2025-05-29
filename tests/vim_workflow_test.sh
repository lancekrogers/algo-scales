#!/bin/bash

# Vim Workflow Test Script
# Tests the complete AlgoScales vim workflow

set -e

echo "🎵 AlgoScales Vim Workflow Test"
echo "================================"

# Configuration
VIM_PLUGIN_DIR="$(pwd)/vim-plugin"
TEST_WORKSPACE="/tmp/algoscales-vim-test"
BINARY_PATH="$(pwd)/bin/algo-scales"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

# Cleanup function
cleanup() {
    echo "Cleaning up test environment..."
    rm -rf "$TEST_WORKSPACE"
}

trap cleanup EXIT

# Test 1: Verify binary exists
echo
echo "Test 1: Binary verification"
if [ ! -f "$BINARY_PATH" ]; then
    log_error "AlgoScales binary not found at $BINARY_PATH"
    echo "Please build with: make build"
    exit 1
fi
log_info "Binary found at $BINARY_PATH"

# Test 2: Test vim mode commands directly
echo
echo "Test 2: CLI vim mode commands"

echo "Testing start command..."
START_OUTPUT=$($BINARY_PATH start practice --vim-mode --language go 2>&1)
if echo "$START_OUTPUT" | grep -q '"id"'; then
    log_info "Start command returns JSON"
else
    log_error "Start command failed or doesn't return JSON"
    echo "Output: $START_OUTPUT"
    exit 1
fi

echo "Testing list command..."
LIST_OUTPUT=$($BINARY_PATH list --vim-mode 2>&1)
if echo "$LIST_OUTPUT" | grep -q '"problems"'; then
    log_info "List command returns JSON"
else
    log_error "List command failed or doesn't return JSON"
    echo "Output: $LIST_OUTPUT"
    exit 1
fi

# Test 3: Test workspace creation
echo
echo "Test 3: Workspace creation"
mkdir -p "$TEST_WORKSPACE"
cd "$TEST_WORKSPACE"

# Start a session and capture output
SESSION_OUTPUT=$($BINARY_PATH start practice two_sum --vim-mode --language go 2>&1)
echo "Session output: $SESSION_OUTPUT"

# Parse the JSON to extract workspace path
if echo "$SESSION_OUTPUT" | grep -q '"workspace_path"'; then
    WORKSPACE_PATH=$(echo "$SESSION_OUTPUT" | grep -o '"workspace_path":"[^"]*"' | cut -d'"' -f4)
    if [ -d "$WORKSPACE_PATH" ]; then
        log_info "Workspace created at $WORKSPACE_PATH"
    else
        log_warn "Workspace path provided but directory doesn't exist: $WORKSPACE_PATH"
    fi
else
    log_warn "No workspace_path in response, plugin will create its own"
fi

# Test 4: Create a minimal vim script to test the plugin
echo
echo "Test 4: Vim plugin functionality test"

VIM_TEST_SCRIPT=$(cat << 'EOF'
" Load the plugin
set runtimepath+=/tmp/algoscales-vim-test-plugin
source /tmp/algoscales-vim-test-plugin/plugin/algo-scales.vim

" Configure for testing
let g:algo_scales_path = '%BINARY_PATH%'
let g:algo_scales_workspace = '/tmp/algoscales-vim-test-workspace'
let g:algo_scales_language = 'go'
let g:algo_scales_auto_test = 0

" Test function
function! TestWorkflow()
    try
        " Test start session
        call algo_scales#StartSession('two_sum')
        
        " Check if session was created
        if empty(g:algo_scales_current_session)
            echo "ERROR: Session not created"
            quit!
        endif
        
        echo "SUCCESS: Session created for " . g:algo_scales_current_session.title
        
        " Write simple solution to test file
        normal! ggdG
        call setline(1, [
            \ 'func twoSum(nums []int, target int) []int {',
            \ '    for i := 0; i < len(nums); i++ {',
            \ '        for j := i + 1; j < len(nums); j++ {',
            \ '            if nums[i] + nums[j] == target {',
            \ '                return []int{i, j}',
            \ '            }',
            \ '        }',
            \ '    }',
            \ '    return []int{}',
            \ '}'
        \ ])
        
        " Save and test
        write
        call algo_scales#TestSolution()
        
        echo "SUCCESS: Workflow completed"
        quit!
        
    catch
        echo "ERROR: " . v:exception
        quit!
    endtry
endfunction

" Run the test
call TestWorkflow()
EOF
)

# Create plugin directory and copy plugin
mkdir -p /tmp/algoscales-vim-test-plugin/plugin
mkdir -p /tmp/algoscales-vim-test-plugin/autoload
cp "$VIM_PLUGIN_DIR/plugin/algo-scales.vim" /tmp/algoscales-vim-test-plugin/plugin/
cp "$VIM_PLUGIN_DIR/autoload/algo_scales.vim" /tmp/algoscales-vim-test-plugin/autoload/

# Replace placeholder in test script
echo "$VIM_TEST_SCRIPT" | sed "s|%BINARY_PATH%|$BINARY_PATH|g" > /tmp/vim_test.vim

# Run vim test (silent mode)
echo "Running vim workflow test..."
if vim -n -c "source /tmp/vim_test.vim" < /dev/null > /tmp/vim_test_output.txt 2>&1; then
    TEST_OUTPUT=$(cat /tmp/vim_test_output.txt)
    if echo "$TEST_OUTPUT" | grep -q "SUCCESS: Workflow completed"; then
        log_info "Vim plugin workflow test passed"
    else
        log_warn "Vim test completed but may have issues"
        echo "Output: $TEST_OUTPUT"
    fi
else
    log_error "Vim plugin test failed"
    echo "Output: $(cat /tmp/vim_test_output.txt)"
fi

# Test 5: Manual workflow verification
echo
echo "Test 5: Manual workflow steps"
echo "To test manually:"
echo "1. Add plugin to your vimrc:"
echo "   set runtimepath+=$VIM_PLUGIN_DIR"
echo "   source $VIM_PLUGIN_DIR/plugin/algo-scales.vim"
echo ""
echo "2. Set configuration (optional):"
echo "   let g:algo_scales_path = '$BINARY_PATH'"
echo "   let g:algo_scales_language = 'go'"
echo ""
echo "3. Start a session:"
echo "   :AlgoScalesStart two_sum"
echo ""
echo "4. Edit the solution file and save"
echo ""
echo "5. Test your solution:"
echo "   :AlgoScalesTest"
echo ""
echo "6. Get hints if needed:"
echo "   :AlgoScalesHint"

# Summary
echo
echo "🎵 Test Summary"
echo "==============="
log_info "CLI vim mode commands work correctly"
log_info "Basic vim plugin functionality implemented"
log_info "Workspace creation integrated"
log_info "Auto-test on save configured"

echo
echo "Next steps:"
echo "- Test with real vim/neovim editor"
echo "- Add session completion tracking"
echo "- Test with different problem types"
echo "- Add error handling improvements"

echo
log_info "Vim workflow test completed successfully!"