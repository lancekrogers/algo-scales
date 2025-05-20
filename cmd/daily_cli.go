// Daily CLI commands
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/daily"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/lancekrogers/algo-scales/internal/session/execution"
	"github.com/spf13/cobra"
)

// dailyTestCmd represents the test command for daily practice
var dailyTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test your solution for the current daily problem",
	Long: `Test your solution for the current problem in daily practice.
This command will verify if your solution passes all test cases.
The problem will only be marked as completed when all tests pass.`,
	Run: func(cmd *cobra.Command, args []string) {
		testDailySolution()
	},
}

// dailySkipCmd represents the skip command for daily practice
var dailySkipCmd = &cobra.Command{
	Use:   "skip",
	Short: "Skip the current daily problem",
	Long: `Skip the current problem in daily practice.
The problem will be marked as skipped and you'll move to the next one.
You can come back to skipped problems later using the resume-skipped command.`,
	Run: func(cmd *cobra.Command, args []string) {
		skipDailyProblem()
	},
}

// dailyResumeSkippedCmd represents the resume-skipped command
var dailyResumeSkippedCmd = &cobra.Command{
	Use:   "resume-skipped",
	Short: "Resume skipped problems",
	Long: `Resume working on problems you previously skipped.
This command will show you a list of skipped problems and let you
choose which one to resume.`,
	Run: func(cmd *cobra.Command, args []string) {
		resumeSkippedProblem()
	},
}

// dailyStatusCmd represents the status command
var dailyStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of your daily practice",
	Long: `Shows the status of your daily practice session.
Displays which problems you've completed, which ones are skipped,
and which one you're currently working on.`,
	Run: func(cmd *cobra.Command, args []string) {
		showDailyStatus()
	},
}

func init() {
	// Add subcommands to daily command
	dailyCmd.AddCommand(dailyTestCmd)
	dailyCmd.AddCommand(dailySkipCmd)
	dailyCmd.AddCommand(dailyResumeSkippedCmd)
	dailyCmd.AddCommand(dailyStatusCmd)
}

// startDailyCliMode starts the CLI-based daily practice session
func startDailyCliMode() {
	// Display welcome message
	fmt.Println("╭───────────────────────────────────────────────────────────────╮")
	fmt.Println("│                🎵 AlgoScales Daily Practice 🎵                │")
	fmt.Println("╰───────────────────────────────────────────────────────────────╯")
	fmt.Println("")
	fmt.Println("Practice one problem from each algorithm pattern (scale) to build your skills.")
	fmt.Println("Each problem will be saved to ~/Dev/AlgoScalesPractice/Daily/{today's date}/")
	fmt.Println("")

	// Load or create daily session
	dailySession, err := daily.GetOrCreateSession()
	if err != nil {
		fmt.Printf("Error initializing daily session: %v\n", err)
		os.Exit(1)
	}

	// Update streak information
	progress, err := daily.LoadProgress()
	if err != nil {
		progress = daily.ScaleProgress{
			Current:       0,
			LastPracticed: time.Time{},
			Completed:     []string{},
			Streak:        0,
			LongestStreak: 0,
		}
	}

	// Update streak based on last practice date
	daily.UpdateStreak(&progress)
	progress.LastPracticed = time.Now()
	
	// Save progress
	if err := daily.SaveProgress(progress); err != nil {
		fmt.Printf("Warning: Error saving progress: %v\n", err)
	}

	// Display streak information
	displayStreakInfo(progress)

	// Display session status
	fmt.Printf("Problems completed today: %d/%d\n", 
		dailySession.GetCompletedCount(), dailySession.GetTotalProblems())
	fmt.Printf("Problems skipped: %d\n", dailySession.GetSkippedCount())
	fmt.Printf("Problems remaining: %d\n\n", 
		dailySession.GetTotalProblems() - 
		dailySession.GetCompletedCount() - 
		dailySession.GetSkippedCount())

	// Check if all problems are completed
	if dailySession.Completed {
		fmt.Println("🎉 Congratulations! You've completed your daily scales practice for all patterns!")
		fmt.Println("Feel free to practice more specific patterns or try a different mode.")
		fmt.Println("")
		fmt.Println("Your current streak: " + fmt.Sprintf("%d days", progress.Streak))
		fmt.Println("Your longest streak: " + fmt.Sprintf("%d days", progress.LongestStreak))
		return
	}

	// Get the next pattern to practice
	nextPattern := dailySession.GetNextPendingPattern()
	if nextPattern == "" {
		fmt.Println("No pending patterns found. Try using 'algo-scales daily resume-skipped'")
		fmt.Println("to work on problems you've skipped.")
		return
	}

	// Get the scale information
	scale := daily.GetScaleByPattern(nextPattern)
	if scale == nil {
		fmt.Printf("Error: Pattern '%s' not found\n", nextPattern)
		return
	}

	// Show current scale information
	fmt.Printf("Now practicing: %s (%s)\n", scale.MusicalName, scale.Pattern)
	fmt.Printf("Description: %s\n\n", scale.Description)

	// Select a problem for this pattern
	prob, err := problem.GetRandomProblemByPattern(scale.Pattern)
	if err != nil {
		fmt.Printf("Error selecting problem: %v\n", err)
		return
	}

	// Update session with problem
	if err := dailySession.StartProblem(scale.Pattern, prob.ID); err != nil {
		fmt.Printf("Error updating session: %v\n", err)
		return
	}

	// Create a problem file with embedded problem text
	filePath, err := daily.CreateProblemFile(prob, language)
	if err != nil {
		fmt.Printf("Error creating problem file: %v\n", err)
		return
	}

	// Show instructions
	fmt.Printf("Problem: %s (%s)\n", prob.Title, prob.Difficulty)
	fmt.Printf("A file has been created at: %s\n\n", filePath)
	
	fmt.Println("Instructions:")
	fmt.Println("1. Open the file to see the problem description in the comments")
	fmt.Println("2. Implement your solution in the file")
	fmt.Println("3. Run 'algo-scales daily test' to test your solution")
	fmt.Println("4. If you want to skip this problem, run 'algo-scales daily skip'")
	
	// Offer to open the editor
	fmt.Print("\nWould you like to open the file in your editor now? (y/n): ")
	var response string
	fmt.Scanln(&response)
	
	if response == "y" || response == "Y" {
		openEditorForDaily(filePath)
	}
}

// testDailySolution tests the solution for the current daily problem
func testDailySolution() {
	// Load session
	dailySession, err := daily.LoadSession()
	if err != nil {
		fmt.Printf("Error loading session: %v\n", err)
		fmt.Println("Please start a daily session first with 'algo-scales daily'")
		return
	}
	
	// Find the in-progress problem
	var currentPattern string
	var currentProblem daily.DailyProblem
	
	for pattern, prob := range dailySession.Problems {
		if prob.State == daily.StateInProgress {
			currentPattern = pattern
			currentProblem = prob
			break
		}
	}
	
	if currentPattern == "" {
		fmt.Println("No problem is currently in progress.")
		fmt.Println("Start a new problem with 'algo-scales daily'")
		return
	}
	
	// Load the problem details
	prob, err := problem.GetByID(currentProblem.ProblemID)
	if err != nil {
		fmt.Printf("Error loading problem: %v\n", err)
		return
	}
	
	// Get the file path
	filePath := daily.GetProblemFilePath(currentProblem.ProblemID, language)
	
	// Check if file exists
	if !daily.ProblemFileExists(currentProblem.ProblemID, language) {
		fmt.Printf("Problem file not found at: %s\n", filePath)
		fmt.Println("Please run 'algo-scales daily' to create the problem file")
		return
	}
	
	fmt.Printf("Testing solution for %s (%s)...\n\n", prob.Title, currentPattern)
	
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading solution file: %v\n", err)
		return
	}
	
	// Create a temporary session to run tests
	tempSession := &session.SessionImpl{
		Problem: prob,
		Options: interfaces.SessionOptions{
			Language: language,
			Mode:     interfaces.SessionMode(session.PracticeMode),
		},
		CodeFile: filePath,
		Code:     string(content),
	}
	
	// First try to run the file directly since it has test code
	var cmd *exec.Cmd
	var allPassed bool
	var results []interfaces.TestResult
	
	// Create results array for output
	results = make([]interfaces.TestResult, len(prob.TestCases))
	for i, tc := range prob.TestCases {
		results[i] = interfaces.TestResult{
			Input:    tc.Input,
			Expected: tc.Expected,
			Actual:   "",
			Passed:   false,
		}
	}
	
	// Execute based on language
	switch language {
	case "go":
		cmd = exec.Command("go", "run", filePath)
	case "python":
		cmd = exec.Command("python", filePath)
	case "javascript":
		cmd = exec.Command("node", filePath)
	default:
		fmt.Printf("Unsupported language: %s\n", language)
		return
	}
	
	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Run the command
	err = cmd.Run()
	
	// Parse test results from output
	output := stdout.String()
	fmt.Println("\nTest Results:")
	fmt.Println(output)
	
	// Check if all tests passed
	allPassed = err == nil && strings.Contains(output, "All tests passed")
	
	// If direct execution fails, fall back to the execution engine
	if err != nil && !strings.Contains(output, "FAILED") {
		fmt.Println("Direct execution failed, falling back to test runner...")
		results, allPassed, err = execution.ExecuteTests(tempSession, 30*time.Second)
		if err != nil {
			fmt.Printf("Error executing tests: %v\n", err)
			return
		}
	}
	
	// Display test results
	fmt.Println("--- Test Results ---")
	
	for i, result := range results {
		passed := "❌ FAILED"
		if result.Passed {
			passed = "✅ PASSED"
		}
		
		fmt.Printf("\nTest %d: %s\n", i+1, passed)
		fmt.Printf("Input: %s\n", result.Input)
		fmt.Printf("Expected: %s\n", result.Expected)
		fmt.Printf("Actual: %s\n", result.Actual)
	}
	
	// If all tests pass, mark the problem as completed
	if allPassed {
		fmt.Println("\n🎉 All tests passed! Problem solved! 🎉")
		
		// Mark problem as completed
		if err := dailySession.CompleteProblem(currentPattern); err != nil {
			fmt.Printf("Error updating session: %v\n", err)
			return
		}
		
		// Check if all problems are completed
		completedCount := dailySession.GetCompletedCount()
		totalProblems := dailySession.GetTotalProblems()
		skippedCount := dailySession.GetSkippedCount()
		
		fmt.Printf("\nProgress: %d/%d problems completed\n", 
			completedCount, totalProblems)
		
		// If there are more problems to solve
		if completedCount + skippedCount < totalProblems {
			fmt.Println("\nWould you like to continue to the next problem? (y/n): ")
			var response string
			fmt.Scanln(&response)
			
			if response == "y" || response == "Y" {
				// Start the next problem
				startDailyCliMode()
			} else {
				fmt.Println("You can continue later with 'algo-scales daily'")
			}
		} else if skippedCount > 0 {
			// All problems either completed or skipped
			fmt.Printf("\nAll problems are either completed (%d) or skipped (%d).\n", 
				completedCount, skippedCount)
			fmt.Println("You can resume skipped problems with 'algo-scales daily resume-skipped'")
		} else {
			// All problems completed
			fmt.Println("\n╭───────────────────────────────────────────────────────────────╮")
			fmt.Println("│         🎵 Congratulations! Daily Scales Complete! 🎵         │")
			fmt.Println("╰───────────────────────────────────────────────────────────────╯")
			fmt.Println("\nYou've completed all algorithm pattern scales for today!")
			
			// Load progress for streak info
			progress, err := daily.LoadProgress()
			if err == nil {
				fmt.Printf("\nCurrent streak: %d days\n", progress.Streak)
				fmt.Printf("Longest streak: %d days\n", progress.LongestStreak)
			}
		}
	} else {
		fmt.Println("\n❌ Some tests failed. Keep working on your solution!")
		fmt.Println("Edit your solution and run 'algo-scales daily test' again when ready.")
	}
}

// skipDailyProblem skips the current daily problem
func skipDailyProblem() {
	// Load session
	dailySession, err := daily.LoadSession()
	if err != nil {
		fmt.Printf("Error loading session: %v\n", err)
		fmt.Println("Please start a daily session first with 'algo-scales daily'")
		return
	}
	
	// Find the in-progress problem
	var currentPattern string
	var currentProblem daily.DailyProblem
	
	for pattern, prob := range dailySession.Problems {
		if prob.State == daily.StateInProgress {
			currentPattern = pattern
			currentProblem = prob
			break
		}
	}
	
	if currentPattern == "" {
		fmt.Println("No problem is currently in progress.")
		fmt.Println("Start a new problem with 'algo-scales daily'")
		return
	}
	
	// Get the scale information
	scale := daily.GetScaleByPattern(currentPattern)
	if scale == nil {
		fmt.Printf("Error: Pattern '%s' not found\n", currentPattern)
		return
	}
	
	// Confirm skip
	fmt.Printf("Are you sure you want to skip the current problem (%s - %s)? (y/n): ", 
		scale.MusicalName, currentProblem.ProblemID)
	var response string
	fmt.Scanln(&response)
	
	if response != "y" && response != "Y" {
		fmt.Println("Skip cancelled.")
		return
	}
	
	// Mark problem as skipped
	if err := dailySession.SkipProblem(currentPattern); err != nil {
		fmt.Printf("Error updating session: %v\n", err)
		return
	}
	
	fmt.Printf("Problem %s (%s) has been skipped.\n", 
		currentProblem.ProblemID, scale.MusicalName)
	
	// Check if all problems are either completed or skipped
	completedCount := dailySession.GetCompletedCount()
	totalProblems := dailySession.GetTotalProblems()
	skippedCount := dailySession.GetSkippedCount()
	
	if completedCount + skippedCount >= totalProblems {
		fmt.Println("\nAll problems are either completed or skipped.")
		fmt.Println("You can resume skipped problems with 'algo-scales daily resume-skipped'")
	} else {
		// Ask if user wants to continue to next problem
		fmt.Print("\nWould you like to continue to the next problem? (y/n): ")
		fmt.Scanln(&response)
		
		if response == "y" || response == "Y" {
			// Start the next problem
			startDailyCliMode()
		} else {
			fmt.Println("You can continue later with 'algo-scales daily'")
		}
	}
}

// resumeSkippedProblem resumes a skipped problem
func resumeSkippedProblem() {
	// Load session
	dailySession, err := daily.LoadSession()
	if err != nil {
		fmt.Printf("Error loading session: %v\n", err)
		fmt.Println("Please start a daily session first with 'algo-scales daily'")
		return
	}
	
	// Find skipped problems
	skippedProblems := make(map[string]daily.DailyProblem)
	for pattern, prob := range dailySession.Problems {
		if prob.State == daily.StateSkipped {
			skippedProblems[pattern] = prob
		}
	}
	
	if len(skippedProblems) == 0 {
		fmt.Println("No skipped problems found.")
		return
	}
	
	// Display skipped problems
	fmt.Println("Skipped problems:")
	
	patternList := make([]string, 0, len(skippedProblems))
	i := 1
	for pattern, prob := range skippedProblems {
		scale := daily.GetScaleByPattern(pattern)
		if scale == nil {
			continue
		}
		
		patternList = append(patternList, pattern)
		fmt.Printf("%d. %s (%s) - %s\n", i, scale.MusicalName, pattern, prob.ProblemID)
		i++
	}
	
	// Ask which problem to resume
	fmt.Print("\nWhich problem would you like to resume? (Enter number): ")
	var choice int
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(patternList) {
		fmt.Println("Invalid choice.")
		return
	}
	
	// Get the selected pattern
	selectedPattern := patternList[choice-1]
	problemInfo := skippedProblems[selectedPattern]
	
	// Load the problem
	prob, err := problem.GetByID(problemInfo.ProblemID)
	if err != nil {
		fmt.Printf("Error loading problem: %v\n", err)
		return
	}
	
	// Check if problem file exists
	if !daily.ProblemFileExists(problemInfo.ProblemID, language) {
		// Create new file
		filePath, err := daily.CreateProblemFile(prob, language)
		if err != nil {
			fmt.Printf("Error creating problem file: %v\n", err)
			return
		}
		fmt.Printf("Problem file created at: %s\n", filePath)
	} else {
		filePath := daily.GetProblemFilePath(problemInfo.ProblemID, language)
		fmt.Printf("Problem file already exists at: %s\n", filePath)
	}
	
	// Update problem state to in-progress
	problemInfo.State = daily.StateInProgress
	dailySession.Problems[selectedPattern] = problemInfo
	
	// Save session
	if err := daily.SaveSession(dailySession); err != nil {
		fmt.Printf("Error saving session: %v\n", err)
		return
	}
	
	scale := daily.GetScaleByPattern(selectedPattern)
	fmt.Printf("\nYou are now working on: %s (%s)\n", 
		scale.MusicalName, selectedPattern)
	fmt.Printf("Problem: %s\n\n", prob.Title)
	
	// Offer to open the editor
	filePath := daily.GetProblemFilePath(problemInfo.ProblemID, language)
	fmt.Print("Would you like to open the file in your editor now? (y/n): ")
	var response string
	fmt.Scanln(&response)
	
	if response == "y" || response == "Y" {
		openEditorForDaily(filePath)
	}
}

// showDailyStatus displays the status of the daily practice session
func showDailyStatus() {
	// Load session
	dailySession, err := daily.LoadSession()
	if err != nil {
		fmt.Printf("Error loading session: %v\n", err)
		fmt.Println("Please start a daily session first with 'algo-scales daily'")
		return
	}
	
	fmt.Println("╭───────────────────────────────────────────────────────────────╮")
	fmt.Println("│                 🎵 Daily Practice Status 🎵                   │")
	fmt.Println("╰───────────────────────────────────────────────────────────────╯")
	
	// Display progress information
	fmt.Printf("\nSession date: %s\n", dailySession.Date)
	fmt.Printf("Problems completed: %d/%d\n", 
		dailySession.GetCompletedCount(), dailySession.GetTotalProblems())
	fmt.Printf("Problems skipped: %d/%d\n", 
		dailySession.GetSkippedCount(), dailySession.GetTotalProblems())
	fmt.Printf("Problems pending: %d/%d\n\n", 
		dailySession.GetPendingCount(), dailySession.GetTotalProblems())
	
	// Display problem status
	fmt.Println("Problem Status:")
	fmt.Println("------------------------------------------")
	
	// Sort by patterns as defined in Scales
	for _, scale := range daily.Scales {
		prob, ok := dailySession.Problems[scale.Pattern]
		if !ok {
			continue
		}
		
		// Get status indicator
		var status string
		switch prob.State {
		case daily.StateCompleted:
			status = "✅ COMPLETED"
		case daily.StateSkipped:
			status = "⏭️ SKIPPED"
		case daily.StateInProgress:
			status = "🔄 IN PROGRESS"
		case daily.StatePending:
			status = "⏳ PENDING"
		}
		
		// Get problem ID or placeholder
		problemID := prob.ProblemID
		if problemID == "" {
			problemID = "<not assigned>"
		}
		
		fmt.Printf("%-20s %-15s %s\n", scale.MusicalName, status, problemID)
	}
	
	// Load progress for streak info
	progress, err := daily.LoadProgress()
	if err == nil {
		fmt.Printf("\nCurrent streak: %d days\n", progress.Streak)
		fmt.Printf("Longest streak: %d days\n", progress.LongestStreak)
	}
	
	// Show what to do next
	fmt.Println("\nNext steps:")
	
	// Check for in-progress problems
	inProgressCount := dailySession.GetInProgressCount()
	if inProgressCount > 0 {
		fmt.Println("- Continue working on your in-progress problem")
		fmt.Println("- Run 'algo-scales daily test' when you're ready to test your solution")
	} else if dailySession.GetPendingCount() > 0 {
		fmt.Println("- Start a new problem with 'algo-scales daily'")
	} else if dailySession.GetSkippedCount() > 0 {
		fmt.Println("- Resume a skipped problem with 'algo-scales daily resume-skipped'")
	} else {
		fmt.Println("- All problems completed! Come back tomorrow for a new set of problems.")
	}
}

// openEditorForDaily opens the file in the user's preferred editor
// This is a renamed version of openEditor to avoid conflict with cli.go
func openEditorForDaily(path string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Try to find a common editor
		editors := []string{"vim", "nano", "emacs", "code", "subl", "pico"}
		for _, e := range editors {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
		
		if editor == "" {
			fmt.Println("No editor found. Please set the EDITOR environment variable.")
			return
		}
	}
	
	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running editor: %v\n", err)
	}
}