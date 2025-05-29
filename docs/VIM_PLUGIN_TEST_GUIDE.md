# AlgoScales Vim Plugin Testing Guide

## Quick Test (Recommended)

```bash
cd /path/to/algo-scales
./test_vim_manual.sh
```

This will start vim with the plugin loaded and show you available commands.

## Manual Setup Test

### 1. Add to your vim config

**For Vim** (`~/.vimrc`):
```vim
set runtimepath+=/Users/lancerogers/Dev/BlockheadConsulting/AlgoScales/algo-scales/vim-plugin
let g:algo_scales_path = '/Users/lancerogers/Dev/BlockheadConsulting/AlgoScales/algo-scales/bin/algo-scales'
let g:algo_scales_language = 'go'
```

**For Neovim** (`~/.config/nvim/init.vim` or `init.lua`):
```vim
set runtimepath+=/Users/lancerogers/Dev/BlockheadConsulting/AlgoScales/algo-scales/vim-plugin
let g:algo_scales_path = '/Users/lancerogers/Dev/BlockheadConsulting/AlgoScales/algo-scales/bin/algo-scales'
let g:algo_scales_language = 'go'
```

### 2. Test the workflow

Open vim and run these commands step by step:

```vim
" 1. Start a session
:AlgoScalesStart two_sum

" You should see:
" - Left pane: Problem description
" - Right pane: Solution file with starter code
" - Bottom pane: Test output area
" - Message: "Started session: Two Sum (easy)"
```

```vim
" 2. Test the solution (write a simple solution first)
" Edit the solution in the right pane, then:
:AlgoScalesTest

" You should see test results in the bottom pane
```

```vim
" 3. Get a hint
:AlgoScalesHint

" Should open a new split with hint information
```

```vim
" 4. List problems
:AlgoScalesList

" Should open a new tab with all available problems
" Press Enter on a problem line to start it
```

```vim
" 5. Daily practice
:AlgoScalesDaily

" Should start a random problem for daily practice
```

## What to Look For

### ✅ Success Indicators

1. **Session Creation**:
   - Three-pane layout appears
   - Problem description in left pane
   - Editable solution file in right pane
   - Test output area at bottom
   - Success message displayed

2. **File Generation**:
   - Check that files are created in workspace directory
   - Solution file contains starter code
   - File is editable and saves correctly

3. **Testing**:
   - `:AlgoScalesTest` runs without errors
   - Test output appears in bottom pane
   - Success/failure message displayed

4. **Auto-test on save**:
   - Edit solution file and save (`:w`)
   - Tests should run automatically
   - Test output should update

### ❌ Troubleshooting

**"Unknown function" errors**:
```bash
# Check that both directories exist:
ls /path/to/algo-scales/vim-plugin/plugin/
ls /path/to/algo-scales/vim-plugin/autoload/
```

**"Binary not found" errors**:
```bash
# Build the binary:
cd /path/to/algo-scales
make build

# Or set absolute path:
let g:algo_scales_path = '/full/path/to/algo-scales/bin/algo-scales'
```

**"No problems found" errors**:
```bash
# Test CLI directly:
/path/to/algo-scales/bin/algo-scales list --vim-mode
```

## Complete Test Scenario

Here's a full workflow to test:

1. **Setup**: Load plugin in vim
2. **Start**: `:AlgoScalesStart sliding-window` 
3. **Code**: Write a solution in the right pane
4. **Test**: Save file (auto-test) or `:AlgoScalesTest`
5. **Hint**: `:AlgoScalesHint` if needed
6. **Complete**: `:AlgoScalesComplete` when done
7. **Stats**: Check if completion was tracked

## Expected File Structure

After starting a session, you should see:

```
~/AlgoScales/  (or your configured workspace)
└── sliding-window/
    └── solution.go  (or .py, .js based on language)
```

Or if workspace_path is provided by CLI:
```
/tmp/algo-scales/sliding-window/
└── solution.go
```

## Performance Test

For a more thorough test:

1. **Test multiple problems**: Try different algorithm patterns
2. **Test different languages**: Set `g:algo_scales_language` to 'python' or 'javascript'
3. **Test auto-test**: Edit and save multiple times
4. **Test hints**: Get progressive hints (run `:AlgoScalesHint` multiple times)
5. **Test completion**: Complete sessions and check stats

## Debugging

If something doesn't work:

1. **Check vim messages**: `:messages`
2. **Test CLI directly**: 
   ```bash
   ./bin/algo-scales start practice two_sum --vim-mode --language go
   ```
3. **Check plugin state**:
   ```vim
   :echo g:algo_scales_current_session
   ```

## Expected Output Examples

**Successful session start**:
```
Started session: Two Sum (easy)
```

**Successful test**:
```
✓ All tests passed!
```

**Failed test**:
```
✗ Some tests failed. Check test output below.
```

This testing approach will help you verify that the complete workflow is functional!