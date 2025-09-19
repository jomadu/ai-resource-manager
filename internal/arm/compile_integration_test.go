package arm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/urf"
)

// TestCompileIntegration_EndToEnd tests complete compilation workflows
func TestCompileIntegration_EndToEnd(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "arm-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test URF files
	testURF := `version: "1.0"
metadata:
  id: "integration-test-rules"
  name: "Integration Test Rules"
  version: "1.0.0"
  description: "Test rules for integration testing"
rules:
  - id: "test-rule-1"
    name: "Test Rule One"
    description: "First test rule"
    priority: 100
    enforcement: "must"
    scope:
      - files: ["**/*.go", "**/*.js"]
    body: |
      This is a test rule for integration testing.

      Example: Write clear, simple code that follows best practices.

  - id: "test-rule-2"
    name: "Test Rule Two"
    description: "Second test rule"
    priority: 80
    enforcement: "should"
    body: |
      This is the second test rule with lower priority.
`

	// Write test URF file
	urfFile := filepath.Join(tempDir, "test-rules.yaml")
	if err := os.WriteFile(urfFile, []byte(testURF), 0o644); err != nil {
		t.Fatalf("Failed to write URF file: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	service := NewArmService()
	ctx := context.Background()

	tests := []struct {
		name            string
		request         *CompileRequest
		expectedFiles   int
		expectedTargets []string
	}{
		{
			name: "single target compilation",
			request: &CompileRequest{
				Files:     []string{urfFile},
				Targets:   []urf.CompileTarget{urf.TargetCursor},
				OutputDir: outputDir,
				Include:   []string{"**/*.yaml"},
				Namespace: "integration-test@1.0.0",
			},
			expectedFiles:   2, // 2 rules
			expectedTargets: []string{"cursor"},
		},
		{
			name: "multi-target compilation",
			request: &CompileRequest{
				Files:     []string{urfFile},
				Targets:   []urf.CompileTarget{urf.TargetCursor, urf.TargetAmazonQ, urf.TargetMarkdown},
				OutputDir: filepath.Join(outputDir, "multi"),
				Include:   []string{"**/*.yaml"},
				Namespace: "integration-test@1.0.0",
				Force:     true, // Allow overwriting
			},
			expectedFiles:   6, // 2 rules × 3 targets
			expectedTargets: []string{"cursor", "amazonq", "markdown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CompileFiles(ctx, tt.request)
			if err != nil {
				t.Fatalf("Compilation failed: %v", err)
			}

			// Verify result statistics
			if result.Stats.FilesProcessed != 1 {
				t.Errorf("Expected 1 file processed, got %d", result.Stats.FilesProcessed)
			}

			if result.Stats.FilesCompiled != tt.expectedFiles {
				t.Errorf("Expected %d files compiled, got %d", tt.expectedFiles, result.Stats.FilesCompiled)
			}

			if result.Stats.Errors != 0 {
				t.Errorf("Expected 0 errors, got %d", result.Stats.Errors)
				for _, err := range result.Errors {
					t.Logf("Error: %s", err.Error)
				}
			}

			// Verify target statistics
			for _, target := range tt.expectedTargets {
				if count, exists := result.Stats.TargetStats[target]; !exists {
					t.Errorf("Expected target %s in statistics", target)
				} else if count != 2 { // 2 rules per target
					t.Errorf("Expected 2 files for target %s, got %d", target, count)
				}
			}

			// Verify actual output files exist
			for _, compiledFile := range result.CompiledFiles {
				if _, err := os.Stat(compiledFile.TargetPath); os.IsNotExist(err) {
					t.Errorf("Compiled file does not exist: %s", compiledFile.TargetPath)
				}

				// Verify file has content
				content, err := os.ReadFile(compiledFile.TargetPath)
				if err != nil {
					t.Errorf("Failed to read compiled file %s: %v", compiledFile.TargetPath, err)
					continue
				}

				if len(content) == 0 {
					t.Errorf("Compiled file is empty: %s", compiledFile.TargetPath)
				}

				// Verify content contains expected elements
				contentStr := string(content)
				if !containsAll(contentStr, []string{"integration-test@1.0.0", "Test Rule"}) {
					t.Errorf("Compiled file missing expected content: %s", compiledFile.TargetPath)
					t.Logf("Content: %s", contentStr)
				}
			}
		})
	}
}

// TestCompileIntegration_DirectoryProcessing tests directory-based compilation
func TestCompileIntegration_DirectoryProcessing(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "arm-dir-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files in different directories
	testFiles := map[string]string{
		"rules1.yaml": `version: "1.0"
metadata:
  id: "rules1"
  name: "Rules One"
  version: "1.0.0"
rules:
  - id: "rule1"
    name: "Rule One"
    priority: 100
    enforcement: "must"
    body: "Rule one content"`,
		"subdir/rules2.yml": `version: "1.0"
metadata:
  id: "rules2"
  name: "Rules Two"
  version: "1.0.0"
rules:
  - id: "rule2"
    name: "Rule Two"
    priority: 90
    enforcement: "should"
    body: "Rule two content"`,
		"ignored.txt": "This should be ignored",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("Failed to create dir for %s: %v", path, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	service := NewArmService()
	ctx := context.Background()
	outputDir := filepath.Join(tempDir, "output")

	tests := []struct {
		name              string
		recursive         bool
		expectedProcessed int
		expectedCompiled  int
	}{
		{
			name:              "non-recursive",
			recursive:         false,
			expectedProcessed: 1, // Only top-level .yaml file
			expectedCompiled:  1,
		},
		{
			name:              "recursive",
			recursive:         true,
			expectedProcessed: 2, // Both .yaml and .yml files
			expectedCompiled:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &CompileRequest{
				Files:     []string{tempDir},
				Targets:   []urf.CompileTarget{urf.TargetCursor},
				OutputDir: filepath.Join(outputDir, tt.name),
				Include:   []string{"**/*.yaml", "**/*.yml"},
				Recursive: tt.recursive,
				Namespace: "dir-test@1.0.0",
			}

			result, err := service.CompileFiles(ctx, request)
			if err != nil {
				t.Fatalf("Compilation failed: %v", err)
			}

			if result.Stats.FilesProcessed != tt.expectedProcessed {
				t.Errorf("Expected %d files processed, got %d", tt.expectedProcessed, result.Stats.FilesProcessed)
			}

			if result.Stats.FilesCompiled != tt.expectedCompiled {
				t.Errorf("Expected %d files compiled, got %d", tt.expectedCompiled, result.Stats.FilesCompiled)
			}

			// Verify .txt file was skipped
			found := false
			for _, skipped := range result.Skipped {
				if filepath.Base(skipped.Path) == "ignored.txt" {
					found = true
					break
				}
			}
			if !found {
				t.Error("Expected ignored.txt to be skipped")
			}
		})
	}
}

// TestCompileIntegration_ErrorHandling tests error scenarios
func TestCompileIntegration_ErrorHandling(t *testing.T) {
	service := NewArmService()
	ctx := context.Background()

	tests := []struct {
		name           string
		request        *CompileRequest
		expectError    bool
		expectedErrors int
	}{
		{
			name: "nonexistent file",
			request: &CompileRequest{
				Files:     []string{"nonexistent.yaml"},
				Targets:   []urf.CompileTarget{urf.TargetCursor},
				OutputDir: "/tmp/test",
				Include:   []string{"**/*.yaml"},
			},
			expectError:    false, // Should not error on file discovery, just process 0 files
			expectedErrors: 0,
		},
		{
			name: "invalid target",
			request: &CompileRequest{
				Files:     []string{"test.yaml"},
				Targets:   []urf.CompileTarget{urf.CompileTarget("invalid")},
				OutputDir: "/tmp/test",
				Include:   []string{"**/*.yaml"},
			},
			expectError:    false, // Error should be captured in result
			expectedErrors: 0,     // No files to process, so no errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CompileFiles(ctx, tt.request)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != nil && result.Stats.Errors != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tt.expectedErrors, result.Stats.Errors)
			}
		})
	}
}

// Helper function to check if content contains all required strings
func containsAll(content string, required []string) bool {
	for _, req := range required {
		if !contains(content, req) {
			return false
		}
	}
	return true
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || substr == "" || indexOfSubstring(s, substr) >= 0)
}

// Simple substring search
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
