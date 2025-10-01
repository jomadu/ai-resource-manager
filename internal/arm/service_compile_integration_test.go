package arm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArmService_CompileFiles_Integration(t *testing.T) {
	tempDir := t.TempDir()
	mockUI := &MockUI{}
	service := &ArmService{ui: mockUI}

	// Create valid URF file
	validURF := `version: "1.0"
metadata:
  id: "integration-test"
  name: "Integration Test Rules"
rules:
  rule1:
    name: "Test Rule 1"
    body: "This is test rule 1"
  rule2:
    name: "Test Rule 2"
    body: "This is test rule 2"`

	filePath := filepath.Join(tempDir, "integration.yaml")
	err := os.WriteFile(filePath, []byte(validURF), 0o644)
	assert.NoError(t, err)

	outputDir := filepath.Join(tempDir, "output")

	// Set up mock expectations
	mockUI.On("CompileStep", mock.AnythingOfType("string")).Maybe()
	mockUI.On("CompileComplete", mock.MatchedBy(func(stats ui.CompileStats) bool {
		return stats.FilesProcessed == 1 && stats.FilesCompiled == 1 && stats.RulesGenerated == 2
	}), false).Once()

	req := &CompileRequest{
		Files:     []string{filePath},
		Target:    "cursor",
		OutputDir: outputDir,
		Verbose:   true,
	}

	err = service.CompileFiles(context.Background(), req)
	assert.NoError(t, err)

	// Verify output files were created
	expectedFiles := []string{
		filepath.Join(outputDir, "integration-test_rule1.mdc"),
		filepath.Join(outputDir, "integration-test_rule2.mdc"),
	}

	for _, expectedFile := range expectedFiles {
		_, err := os.Stat(expectedFile)
		assert.NoError(t, err, "Expected file %s should exist", expectedFile)

		// Verify file has content
		content, err := os.ReadFile(expectedFile)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
	}

	mockUI.AssertExpectations(t)
}

func TestArmService_CompileFiles_MultiTarget_Integration(t *testing.T) {
	tempDir := t.TempDir()
	mockUI := &MockUI{}
	service := &ArmService{ui: mockUI}

	// Create valid URF file
	validURF := `version: "1.0"
metadata:
  id: "multi-target-test"
  name: "Multi Target Test Rules"
rules:
  rule1:
    name: "Test Rule"
    body: "This is a test rule"`

	filePath := filepath.Join(tempDir, "multi.yaml")
	err := os.WriteFile(filePath, []byte(validURF), 0o644)
	assert.NoError(t, err)

	outputDir := filepath.Join(tempDir, "multi-output")

	// Set up mock expectations
	mockUI.On("CompileStep", mock.AnythingOfType("string")).Maybe()
	mockUI.On("CompileComplete", mock.MatchedBy(func(stats ui.CompileStats) bool {
		return stats.FilesProcessed == 1 && stats.FilesCompiled == 1 && stats.RulesGenerated == 2
	}), false).Once()

	req := &CompileRequest{
		Files:     []string{filePath},
		Target:    "cursor,amazonq",
		OutputDir: outputDir,
		Verbose:   true,
	}

	err = service.CompileFiles(context.Background(), req)
	assert.NoError(t, err)

	// Verify output files were created in target subdirectories
	expectedFiles := []string{
		filepath.Join(outputDir, "cursor", "multi-target-test_rule1.mdc"),
		filepath.Join(outputDir, "amazonq", "multi-target-test_rule1.md"),
	}

	for _, expectedFile := range expectedFiles {
		_, err := os.Stat(expectedFile)
		assert.NoError(t, err, "Expected file %s should exist", expectedFile)

		// Verify file has content
		content, err := os.ReadFile(expectedFile)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
	}

	mockUI.AssertExpectations(t)
}
