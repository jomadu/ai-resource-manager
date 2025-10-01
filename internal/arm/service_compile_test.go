package arm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/ui"
	"github.com/jomadu/ai-rules-manager/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUI implements ui.Interface for testing
type MockUI struct {
	mock.Mock
}

func (m *MockUI) InstallStep(step string) { m.Called(step) }
func (m *MockUI) InstallStepWithSpinner(step string) func(result string) {
	return m.Called(step).Get(0).(func(result string))
}
func (m *MockUI) InstallComplete(registry, ruleset, version string, sinks []string) {
	m.Called(registry, ruleset, version, sinks)
}
func (m *MockUI) Success(msg string) { m.Called(msg) }
func (m *MockUI) Error(err error)    { m.Called(err) }
func (m *MockUI) Warning(msg string) { m.Called(msg) }
func (m *MockUI) ConfigList(registries map[string]map[string]interface{}, sinks map[string]manifest.SinkConfig) {
	m.Called(registries, sinks)
}
func (m *MockUI) RulesetList(rulesets []*ui.RulesetInfo) { m.Called(rulesets) }
func (m *MockUI) RulesetInfoGrouped(rulesets []*ui.RulesetInfo, detailed bool) {
	m.Called(rulesets, detailed)
}
func (m *MockUI) OutdatedTable(outdated []ui.OutdatedRuleset, outputFormat string) {
	m.Called(outdated, outputFormat)
}
func (m *MockUI) VersionInfo(info version.VersionInfo) { m.Called(info) }
func (m *MockUI) CompileStep(step string)              { m.Called(step) }
func (m *MockUI) CompileComplete(stats ui.CompileStats, validateOnly bool) {
	m.Called(stats, validateOnly)
}

func TestArmService_isURFFile(t *testing.T) {
	service := &ArmService{}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"YAML file", "test.yaml", true},
		{"YML file", "test.yml", true},
		{"Uppercase YAML", "test.YAML", true},
		{"Uppercase YML", "test.YML", true},
		{"JSON file", "test.json", false},
		{"Text file", "test.txt", false},
		{"No extension", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isURFFile(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestArmService_discoverFiles(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Create test files
	yamlFile := filepath.Join(tempDir, "test.yaml")
	ymlFile := filepath.Join(tempDir, "test.yml")
	txtFile := filepath.Join(tempDir, "test.txt")
	subDir := filepath.Join(tempDir, "subdir")
	subYamlFile := filepath.Join(subDir, "sub.yaml")

	err := os.WriteFile(yamlFile, []byte("test"), 0o644)
	assert.NoError(t, err)
	err = os.WriteFile(ymlFile, []byte("test"), 0o644)
	assert.NoError(t, err)
	err = os.WriteFile(txtFile, []byte("test"), 0o644)
	assert.NoError(t, err)
	err = os.MkdirAll(subDir, 0o755)
	assert.NoError(t, err)
	err = os.WriteFile(subYamlFile, []byte("test"), 0o644)
	assert.NoError(t, err)

	service := &ArmService{}

	t.Run("discover single file", func(t *testing.T) {
		files, err := service.discoverFiles([]string{yamlFile}, false, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Contains(t, files, yamlFile)
	})

	t.Run("discover directory non-recursive", func(t *testing.T) {
		files, err := service.discoverFiles([]string{tempDir}, false, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, files, 2) // test.yaml and test.yml
		assert.Contains(t, files, yamlFile)
		assert.Contains(t, files, ymlFile)
		assert.NotContains(t, files, txtFile)
		assert.NotContains(t, files, subYamlFile)
	})

	t.Run("discover directory recursive", func(t *testing.T) {
		files, err := service.discoverFiles([]string{tempDir}, true, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, files, 3) // test.yaml, test.yml, and subdir/sub.yaml
		assert.Contains(t, files, yamlFile)
		assert.Contains(t, files, ymlFile)
		assert.Contains(t, files, subYamlFile)
		assert.NotContains(t, files, txtFile)
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := service.discoverFiles([]string{"nonexistent.yaml"}, false, nil, nil)
		assert.Error(t, err)
	})
}

func TestArmService_validateURFFile(t *testing.T) {
	tempDir := t.TempDir()
	service := &ArmService{}

	t.Run("valid URF file", func(t *testing.T) {
		validURF := `version: "1.0"
metadata:
  id: "test-rules"
  name: "Test Rules"
rules:
  rule1:
    name: "Test Rule"
    body: "This is a test rule"`

		filePath := filepath.Join(tempDir, "valid.yaml")
		err := os.WriteFile(filePath, []byte(validURF), 0o644)
		assert.NoError(t, err)

		err = service.validateURFFile(filePath)
		assert.NoError(t, err)
	})

	t.Run("invalid URF file", func(t *testing.T) {
		invalidURF := `invalid: yaml`

		filePath := filepath.Join(tempDir, "invalid.yaml")
		err := os.WriteFile(filePath, []byte(invalidURF), 0o644)
		assert.NoError(t, err)

		err = service.validateURFFile(filePath)
		assert.Error(t, err)
	})

	t.Run("non-existent file", func(t *testing.T) {
		err := service.validateURFFile("nonexistent.yaml")
		assert.Error(t, err)
	})
}

func TestArmService_CompileFiles_ValidateOnly(t *testing.T) {
	tempDir := t.TempDir()
	mockUI := &MockUI{}
	service := &ArmService{ui: mockUI}

	// Create valid URF file
	validURF := `version: "1.0"
metadata:
  id: "test-rules"
  name: "Test Rules"
rules:
  rule1:
    name: "Test Rule"
    body: "This is a test rule"`

	filePath := filepath.Join(tempDir, "test.yaml")
	err := os.WriteFile(filePath, []byte(validURF), 0o644)
	assert.NoError(t, err)

	// Set up mock expectations
	mockUI.On("CompileStep", mock.AnythingOfType("string")).Maybe()
	mockUI.On("Success", mock.AnythingOfType("string")).Maybe()
	mockUI.On("CompileComplete", mock.AnythingOfType("ui.CompileStats"), true).Once()

	req := &CompileRequest{
		Files:        []string{filePath},
		Target:       "cursor",
		ValidateOnly: true,
		Verbose:      true,
	}

	err = service.CompileFiles(context.Background(), req)
	assert.NoError(t, err)

	mockUI.AssertExpectations(t)
}

func TestArmService_CompileFiles_EmptyFiles(t *testing.T) {
	tempDir := t.TempDir()
	mockUI := &MockUI{}
	service := &ArmService{ui: mockUI}

	// Set up mock expectations
	mockUI.On("Warning", "No URF files found matching the criteria").Once()

	req := &CompileRequest{
		Files:  []string{tempDir}, // Empty directory
		Target: "cursor",
	}

	err := service.CompileFiles(context.Background(), req)
	assert.NoError(t, err)

	mockUI.AssertExpectations(t)
}
