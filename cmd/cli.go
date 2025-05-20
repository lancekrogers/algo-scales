// CLI mode implementation for AlgoScales
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lancekrogers/algo-scales/internal/session"
	"github.com/spf13/cobra"
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "solve [problem]",
	Short: "Solve a problem in CLI mode",
	Long: `Solve a problem using the command-line interface:
1. Choose a problem (or specify one)
2. View the problem statement
3. Edit the solution code in your editor
4. Test your solution
5. Submit when ready`,
	Run: func(cmd *cobra.Command, args []string) {
		var problemID string
		if len(args) > 0 {
			problemID = args[0]
		}

		opts := session.Options{
			Mode:       session.PracticeMode, // Default to practice mode
			Language:   language,
			Timer:      timer,
			Pattern:    pattern,
			Difficulty: difficulty,
			ProblemID:  problemID,
		}

		// Create session without starting UI
		sess, err := session.CreateSession(opts)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error creating session: %v\n", err)
			return
		}

		// Create a session adapter
		adapter := &SessionAdapter{Session: sess}
		
		// Run CLI problem solving workflow
		if err := runCliWorkflow(adapter); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error in CLI workflow: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)

	// Add flags to the cli command
	cliCmd.Flags().StringVarP(&language, "language", "l", "go", "Programming language (go, python, javascript)")
	cliCmd.Flags().IntVarP(&timer, "timer", "t", 45, "Timer duration in minutes (15, 30, 45, 60)")
	cliCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Algorithm pattern to focus on")
	cliCmd.Flags().StringVarP(&difficulty, "difficulty", "d", "", "Problem difficulty (easy, medium, hard)")
}

// runCliWorkflow handles the CLI problem-solving workflow
func runCliWorkflow(s *SessionAdapter) error {
	// Display welcome message
	fmt.Println("🎵 AlgoScales CLI Mode 🎵")
	fmt.Println("—————————————————————————")

	// Show problem details
	fmt.Printf("Problem: %s (%s)\n", s.Problem.Title, s.Problem.Difficulty)
	fmt.Printf("Pattern: %s\n", JoinStrings(s.Problem.Patterns))
	fmt.Printf("Estimated Time: %d minutes\n\n", s.Problem.EstimatedTime)

	// Path to files
	descFile := filepath.Join(s.Workspace, "problem.md")
	codeFile := s.CodeFile

	// Main interaction loop
	for {
		// Display menu
		fmt.Println("\nOptions:")
		fmt.Println("1. View problem description")
		fmt.Println("2. Edit solution")
		fmt.Println("3. Test solution")
		if s.Options.Mode == session.LearnMode {
			fmt.Println("4. View hints")
			fmt.Println("5. View solution")
			fmt.Println("6. Exit")
		} else {
			fmt.Println("4. Exit")
		}

		// Get user choice
		fmt.Print("\nEnter your choice: ")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1": // View problem
			viewFile(descFile)

		case "2": // Edit solution
			// Open in user's preferred editor
			openEditor(codeFile)

			// Update the session code after editing
			code, err := os.ReadFile(codeFile)
			if err != nil {
				fmt.Printf("Error reading code file: %v\n", err)
				continue
			}
			s.SetCode(string(code))

		case "3": // Test solution
			// Run tests
			results, allPassed, err := s.RunTests()
			if err != nil {
				fmt.Printf("Error running tests: %v\n", err)
				continue
			}

			// Display test results
			fmt.Println("\n--- Test Results ---")
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

			if allPassed {
				fmt.Println("\n🎉 All tests passed! Problem solved! 🎉")

				// Record completion
				s.FinishSession(true)
				return nil
			}

		case "4":
			if s.Options.Mode == session.LearnMode {
				// View hints - we'll just show the pattern explanation
				fmt.Println("\n--- Pattern Information ---")
				fmt.Println(s.Problem.PatternExplanation)
				s.ShowHints(true)
			} else {
				// Exit
				fmt.Println("Exiting session...")
				s.FinishSession(false)
				return nil
			}

		case "5":
			if s.Options.Mode == session.LearnMode {
				// View solution
				fmt.Println("\n--- Solution ---")

				if solution, ok := s.Problem.Solutions[s.Options.Language]; ok {
					fmt.Println(solution)
				} else {
					// Try to find a solution in any language
					for _, solution := range s.Problem.Solutions {
						fmt.Println(solution)
						break
					}
				}

				s.ShowSolution(true)
			} else {
				fmt.Println("Invalid choice. Please try again.")
			}

		case "6":
			if s.Options.Mode == session.LearnMode {
				// Exit
				fmt.Println("Exiting session...")
				s.FinishSession(false)
				return nil
			} else {
				fmt.Println("Invalid choice. Please try again.")
			}

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// viewFile displays the contents of a file
func viewFile(path string) {
	// Check for common pager programs
	pagers := []string{"less", "more", "cat"}
	var pager string

	for _, p := range pagers {
		if _, err := exec.LookPath(p); err == nil {
			pager = p
			break
		}
	}

	if pager == "" {
		// Fall back to reading file directly
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Println(string(data))
		return
	}

	// Use the pager
	cmd := exec.Command(pager, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running pager: %v\n", err)

		// Fall back to reading file directly
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Println(string(data))
	}
}

// openEditor opens the file in the user's preferred editor
func openEditor(path string) {
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

// JoinStrings joins a string slice with commas (redefined to avoid import circular references)
func JoinStrings(strs []string) string {
	return strings.Join(strs, ", ")
}