package session

import (
	"path/filepath"

	"github.com/lancekrogers/algo-scales/internal/common/interfaces"
	"github.com/lancekrogers/algo-scales/internal/problem"
)

// CodeManagerImpl implements the CodeManager interface
type CodeManagerImpl struct {
	fs            interfaces.FileSystem
	templateSvc   interfaces.TemplateService
	workspace     string
	codeFile      string
	currentCode   string
}

// NewCodeManager creates a new code manager
func NewCodeManager(fs interfaces.FileSystem, templateSvc interfaces.TemplateService) interfaces.CodeManager {
	return &CodeManagerImpl{
		fs:          fs,
		templateSvc: templateSvc,
	}
}

// GetCode returns the current user code
func (cm *CodeManagerImpl) GetCode() string {
	if cm.currentCode != "" {
		return cm.currentCode
	}
	
	if cm.codeFile != "" && cm.fs.Exists(cm.codeFile) {
		data, err := cm.fs.ReadFile(cm.codeFile)
		if err == nil {
			cm.currentCode = string(data)
			return cm.currentCode
		}
	}
	
	return ""
}

// SetCode updates the user code
func (cm *CodeManagerImpl) SetCode(code string) error {
	cm.currentCode = code
	
	if cm.codeFile != "" {
		err := cm.fs.WriteFile(cm.codeFile, []byte(code), 0644)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// GetWorkspace returns the workspace directory
func (cm *CodeManagerImpl) GetWorkspace() string {
	return cm.workspace
}

// GetCodeFile returns the path to the code file
func (cm *CodeManagerImpl) GetCodeFile() string {
	return cm.codeFile
}

// SetWorkspace sets the workspace directory
func (cm *CodeManagerImpl) SetWorkspace(workspace string) error {
	cm.workspace = workspace
	return nil
}

// InitializeWorkspace creates workspace and initial code file
func (cm *CodeManagerImpl) InitializeWorkspace(prob *problem.Problem, language string) error {
	// Create workspace directory
	if cm.workspace == "" {
		tempDir := cm.fs.TempDir()
		cm.workspace = filepath.Join(tempDir, "algo-scales", prob.ID)
	}
	
	err := cm.fs.MkdirAll(cm.workspace, 0755)
	if err != nil {
		return err
	}
	
	// Determine file extension based on language
	var extension string
	switch language {
	case "go":
		extension = ".go"
	case "python":
		extension = ".py"
	case "javascript":
		extension = ".js"
	default:
		extension = ".txt"
	}
	
	cm.codeFile = filepath.Join(cm.workspace, "solution"+extension)
	
	// Generate initial template code
	if cm.templateSvc != nil {
		templateCode, err := cm.templateSvc.GenerateTemplate(prob, language)
		if err == nil {
			cm.currentCode = templateCode
		}
	}
	
	// Write initial code to file
	if cm.currentCode != "" {
		err = cm.fs.WriteFile(cm.codeFile, []byte(cm.currentCode), 0644)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// CleanupWorkspace removes temporary files
func (cm *CodeManagerImpl) CleanupWorkspace() error {
	if cm.workspace != "" {
		return cm.fs.RemoveAll(cm.workspace)
	}
	return nil
}

// WithFileSystem sets a custom file system
func (cm *CodeManagerImpl) WithFileSystem(fs interfaces.FileSystem) *CodeManagerImpl {
	cm.fs = fs
	return cm
}

// WithTemplateService sets a custom template service
func (cm *CodeManagerImpl) WithTemplateService(svc interfaces.TemplateService) *CodeManagerImpl {
	cm.templateSvc = svc
	return cm
}