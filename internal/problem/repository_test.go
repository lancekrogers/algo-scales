package problem

import (
	"context"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
)

func TestMockRepository(t *testing.T) {
	// Create a mock repository
	mockRepo := NewMockRepository()
	
	// Add test problems
	mockRepo.AddProblem(Problem{
		ID:         "test-problem-1",
		Title:      "Test Problem 1",
		Difficulty: "easy",
		Patterns:   []string{"two-pointers", "sliding-window"},
		Companies:  []string{"google", "amazon"},
		StarterCode: map[string]string{
			"go":       "func testProblem() {}",
			"python":   "def test_problem(): pass",
		},
	})
	
	mockRepo.AddProblem(Problem{
		ID:         "test-problem-2",
		Title:      "Test Problem 2",
		Difficulty: "medium",
		Patterns:   []string{"dynamic-programming"},
		Companies:  []string{"facebook", "microsoft"},
		StarterCode: map[string]string{
			"go":       "func testProblem2() {}",
			"python":   "def test_problem2(): pass",
			"javascript": "function testProblem2() {}",
		},
	})
	
	// Set patterns and languages
	mockRepo.SetPatterns([]string{"two-pointers", "sliding-window", "dynamic-programming"})
	mockRepo.SetLanguages([]string{"go", "python", "javascript"})
	
	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		problems, err := mockRepo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, problems, 2)
	})
	
	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		problem, err := mockRepo.GetByID(context.Background(), "test-problem-1")
		assert.NoError(t, err)
		assert.Equal(t, "Test Problem 1", problem.Title)
		
		// Non-existent problem
		_, err = mockRepo.GetByID(context.Background(), "non-existent")
		assert.Error(t, err)
		assert.Equal(t, ErrProblemNotFound, err)
	})
	
	// Test GetByPattern
	t.Run("GetByPattern", func(t *testing.T) {
		problems, err := mockRepo.GetByPattern(context.Background(), "sliding-window")
		assert.NoError(t, err)
		assert.Len(t, problems, 1)
		assert.Equal(t, "test-problem-1", problems[0].ID)
		
		problems, err = mockRepo.GetByPattern(context.Background(), "dynamic-programming")
		assert.NoError(t, err)
		assert.Len(t, problems, 1)
		assert.Equal(t, "test-problem-2", problems[0].ID)
		
		// Empty pattern returns all
		problems, err = mockRepo.GetByPattern(context.Background(), "")
		assert.NoError(t, err)
		assert.Len(t, problems, 2)
	})
	
	// Test GetPatterns
	t.Run("GetPatterns", func(t *testing.T) {
		patterns, err := mockRepo.GetPatterns(context.Background())
		assert.NoError(t, err)
		assert.Len(t, patterns, 3)
		assert.Contains(t, patterns, "two-pointers")
		assert.Contains(t, patterns, "sliding-window")
		assert.Contains(t, patterns, "dynamic-programming")
	})
	
	// Test GetLanguages
	t.Run("GetLanguages", func(t *testing.T) {
		languages, err := mockRepo.GetLanguages(context.Background())
		assert.NoError(t, err)
		assert.Len(t, languages, 3)
		assert.Contains(t, languages, "go")
		assert.Contains(t, languages, "python")
		assert.Contains(t, languages, "javascript")
	})
	
	// Test GetByDifficulty
	t.Run("GetByDifficulty", func(t *testing.T) {
		problems, err := mockRepo.GetByDifficulty(context.Background(), "easy")
		assert.NoError(t, err)
		assert.Len(t, problems, 1)
		assert.Equal(t, "test-problem-1", problems[0].ID)
		
		problems, err = mockRepo.GetByDifficulty(context.Background(), "medium")
		assert.NoError(t, err)
		assert.Len(t, problems, 1)
		assert.Equal(t, "test-problem-2", problems[0].ID)
		
		// No problems with this difficulty
		problems, err = mockRepo.GetByDifficulty(context.Background(), "hard")
		assert.NoError(t, err)
		assert.Len(t, problems, 0)
	})
	
	// Test GetByCompany
	t.Run("GetByCompany", func(t *testing.T) {
		problems, err := mockRepo.GetByCompany(context.Background(), "google")
		assert.NoError(t, err)
		assert.Len(t, problems, 1)
		assert.Equal(t, "test-problem-1", problems[0].ID)
		
		problems, err = mockRepo.GetByCompany(context.Background(), "microsoft")
		assert.NoError(t, err)
		assert.Len(t, problems, 1)
		assert.Equal(t, "test-problem-2", problems[0].ID)
		
		// No problems for this company
		problems, err = mockRepo.GetByCompany(context.Background(), "apple")
		assert.NoError(t, err)
		assert.Len(t, problems, 0)
	})
}

// Test the Service with the mock repository
func TestServiceWithMockRepository(t *testing.T) {
	// Create a mock repository
	mockRepo := NewMockRepository()
	
	// Add test problems
	mockRepo.AddProblem(Problem{
		ID:         "test-problem-1",
		Title:      "Test Problem 1",
		Difficulty: "easy",
		Patterns:   []string{"two-pointers"},
		Companies:  []string{"google"},
		TestCases:  []TestCase{
			{Input: "input1", Expected: "output1"},
			{Input: "input2", Expected: "output2"},
		},
	})
	
	// Create service with mock repository
	service := NewService().WithRepository(mockRepo)
	
	// Test ListAll
	t.Run("ListAll", func(t *testing.T) {
		problems, err := service.ListAll()
		assert.NoError(t, err)
		assert.Len(t, problems, 1)
		assert.Equal(t, "test-problem-1", problems[0].ID)
	})
	
	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		problem, err := service.GetByID("test-problem-1")
		assert.NoError(t, err)
		assert.Equal(t, "Test Problem 1", problem.Title)
	})
	
	// Test TestSolution
	t.Run("TestSolution", func(t *testing.T) {
		results, err := service.TestSolution("test-problem-1", "test code", "go")
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "input1", results[0].Input)
		assert.Equal(t, "output1", results[0].Expected)
		assert.True(t, results[0].Passed)
	})
}

// Verify the Repository implements ProblemRepository
func TestRepositoryInterface(t *testing.T) {
	var repo interfaces.ProblemRepository = NewRepository()
	assert.NotNil(t, repo)
}