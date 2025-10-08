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

func TestArmService_CompileFiles_ValidateOnly(t *testing.T) {
	tempDir := t.TempDir()
	mockUI := &MockUI{}
	service := &ArmService{ui: mockUI}

	// Create valid resource file
	validResource := `apiVersion: v1
kind: Ruleset
metadata:
  id: "test-rules"
  name: "Test Rules"
spec:
  rules:
    rule1:
      name: "Test Rule"
      body: "This is a test rule"`

	filePath := filepath.Join(tempDir, "test.yaml")
	err := os.WriteFile(filePath, []byte(validResource), 0o644)
	assert.NoError(t, err)

	// Set up mock expectations
	mockUI.On("CompileStep", mock.AnythingOfType("string")).Maybe()
	mockUI.On("Success", mock.AnythingOfType("string")).Maybe()
	mockUI.On("CompileComplete", mock.AnythingOfType("ui.CompileStats"), true).Once()

	req := &CompileRequest{
		Paths:        []string{filePath},
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
	mockUI.On("Warning", "No resource files found matching the criteria").Once()

	req := &CompileRequest{
		Paths:  []string{tempDir}, // Empty directory
		Target: "cursor",
	}

	err := service.CompileFiles(context.Background(), req)
	assert.NoError(t, err)

	mockUI.AssertExpectations(t)
}
