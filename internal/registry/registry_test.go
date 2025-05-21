package registry

import (
	"testing"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/session/execution"
)

// MockProblemRepository for testing
type MockProblemRepository struct{}

func (m *MockProblemRepository) GetAll() ([]interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetByID(id string) (*interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetByPattern(pattern string) ([]interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetByDifficulty(difficulty string) ([]interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetByTags(tags []string) ([]interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetRandom() (*interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetRandomByPattern(pattern string) (*interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetRandomByDifficulty(difficulty string) (*interfaces.Problem, error) { return nil, nil }
func (m *MockProblemRepository) GetRandomByTags(tags []string) (*interfaces.Problem, error) { return nil, nil }

// MockFileSystem for testing
type MockFileSystem struct{}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) { return nil, nil }
func (m *MockFileSystem) WriteFile(path string, data []byte, perm interface{}) error { return nil }
func (m *MockFileSystem) MkdirAll(path string, perm interface{}) error { return nil }
func (m *MockFileSystem) Remove(path string) error { return nil }
func (m *MockFileSystem) RemoveAll(path string) error { return nil }
func (m *MockFileSystem) Exists(path string) bool { return false }
func (m *MockFileSystem) IsDir(path string) bool { return false }
func (m *MockFileSystem) IsFile(path string) bool { return false }
func (m *MockFileSystem) List(path string) ([]string, error) { return nil, nil }
func (m *MockFileSystem) Glob(pattern string) ([]string, error) { return nil, nil }
func (m *MockFileSystem) Getwd() (string, error) { return "", nil }
func (m *MockFileSystem) Chdir(path string) error { return nil }
func (m *MockFileSystem) TempDir() string { return "" }
func (m *MockFileSystem) Join(paths ...string) string { return "" }
func (m *MockFileSystem) Base(path string) string { return "" }
func (m *MockFileSystem) Dir(path string) string { return "" }

// MockStatsService for testing
type MockStatsService struct{}

func (m *MockStatsService) RecordSolve(problemID, pattern, difficulty, language string, success bool, timeSpent int, codeLength int) error { return nil }
func (m *MockStatsService) GetOverallStats() (*interfaces.OverallStats, error) { return nil, nil }
func (m *MockStatsService) GetPatternStats() (map[string]*interfaces.PatternStats, error) { return nil, nil }
func (m *MockStatsService) GetDifficultyStats() (map[string]*interfaces.DifficultyStats, error) { return nil, nil }
func (m *MockStatsService) GetLanguageStats() (map[string]*interfaces.LanguageStats, error) { return nil, nil }
func (m *MockStatsService) GetRecentActivity(days int) ([]*interfaces.DailyStats, error) { return nil, nil }

// MockTemplateService for testing
type MockTemplateService struct{}

func (m *MockTemplateService) GenerateTemplate(problem *interfaces.Problem, language string) (string, error) { return "", nil }
func (m *MockTemplateService) GetSupportedLanguages() []string { return nil }

func TestNewServiceRegistry(t *testing.T) {
	registry := NewServiceRegistry()
	
	if registry == nil {
		t.Fatal("Expected non-nil registry")
	}
	
	if registry.GetTestRunnerRegistry() == nil {
		t.Error("Expected test runner registry to be initialized")
	}
}

func TestServiceRegistryWithMethods(t *testing.T) {
	registry := NewServiceRegistry()
	
	// Test WithFileSystem
	mockFS := &MockFileSystem{}
	registry.WithFileSystem(mockFS)
	
	fs, err := registry.GetFileSystem()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if fs != mockFS {
		t.Error("Expected same filesystem instance")
	}
	
	// Test WithProblemRepository
	mockRepo := &MockProblemRepository{}
	registry.WithProblemRepository(mockRepo)
	
	repo, err := registry.GetProblemRepository()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if repo != mockRepo {
		t.Error("Expected same repository instance")
	}
	
	// Test WithStatsService
	mockStats := &MockStatsService{}
	registry.WithStatsService(mockStats)
	
	service, err := registry.GetStatsService()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if service != mockStats {
		t.Error("Expected same stats service instance")
	}
	
	// Test WithTemplateService
	mockTemplate := &MockTemplateService{}
	registry.WithTemplateService(mockTemplate)
	
	tmplService, err := registry.GetTemplateService()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tmplService != mockTemplate {
		t.Error("Expected same template service instance")
	}
	
	// Test WithTestRunnerRegistry
	testRunnerReg := execution.NewRunnerRegistry()
	registry.WithTestRunnerRegistry(testRunnerReg)
	
	if registry.GetTestRunnerRegistry() != testRunnerReg {
		t.Error("Expected same test runner registry instance")
	}
}

func TestServiceRegistryGetters_Uninitialized(t *testing.T) {
	registry := NewServiceRegistry()
	
	// Test uninitialized services return errors
	_, err := registry.GetFileSystem()
	if err == nil {
		t.Error("Expected error for uninitialized file system")
	}
	
	_, err = registry.GetProblemRepository()
	if err == nil {
		t.Error("Expected error for uninitialized problem repository")
	}
	
	_, err = registry.GetStatsService()
	if err == nil {
		t.Error("Expected error for uninitialized stats service")
	}
	
	_, err = registry.GetTemplateService()
	if err == nil {
		t.Error("Expected error for uninitialized template service")
	}
	
	// Test runner registry is always available
	if registry.GetTestRunnerRegistry() == nil {
		t.Error("Expected test runner registry to be available")
	}
}

func TestInitializeDefaults(t *testing.T) {
	// Save original state
	originalRegistry := DefaultRegistry
	defer func() {
		DefaultRegistry = originalRegistry
	}()
	
	// Reset registry for testing
	DefaultRegistry = NewServiceRegistry()
	
	// Initialize with mock services
	InitializeDefaults(
		&MockFileSystem{},
		&MockProblemRepository{},
		&MockStatsService{},
		&MockTemplateService{},
	)
	
	// Verify all services are initialized
	_, err := DefaultRegistry.GetFileSystem()
	if err != nil {
		t.Errorf("Expected file system to be initialized, got %v", err)
	}
	
	_, err = DefaultRegistry.GetProblemRepository()
	if err != nil {
		t.Errorf("Expected problem repository to be initialized, got %v", err)
	}
	
	_, err = DefaultRegistry.GetStatsService()
	if err != nil {
		t.Errorf("Expected stats service to be initialized, got %v", err)
	}
	
	_, err = DefaultRegistry.GetTemplateService()
	if err != nil {
		t.Errorf("Expected template service to be initialized, got %v", err)
	}
	
	if DefaultRegistry.GetTestRunnerRegistry() == nil {
		t.Error("Expected test runner registry to be available")
	}
}