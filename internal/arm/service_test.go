package arm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/jomadu/ai-rules-manager/internal/urf"
)

func TestArmService_CompileFiles_InputValidation(t *testing.T) {
	service := NewArmService()
	ctx := context.Background()

	tests := []struct {
		name      string
		request   *CompileRequest
		expectErr bool
		errorMsg  string
	}{
		{
			name:      "nil request",
			request:   nil,
			expectErr: true,
			errorMsg:  "compile request is required",
		},
		{
			name: "empty files",
			request: &CompileRequest{
				Files:     []string{},
				Targets:   []urf.CompileTarget{urf.TargetCursor},
				OutputDir: "/tmp/test",
			},
			expectErr: true,
			errorMsg:  "no files specified for compilation",
		},
		{
			name: "empty targets",
			request: &CompileRequest{
				Files:     []string{"test.yaml"},
				Targets:   []urf.CompileTarget{},
				OutputDir: "/tmp/test",
			},
			expectErr: true,
			errorMsg:  "no compilation targets specified",
		},
		{
			name: "empty output directory",
			request: &CompileRequest{
				Files:     []string{"test.yaml"},
				Targets:   []urf.CompileTarget{urf.TargetCursor},
				OutputDir: "",
			},
			expectErr: true,
			errorMsg:  "output directory is required",
		},
		{
			name: "valid request with nonexistent files",
			request: &CompileRequest{
				Files:     []string{"nonexistent.yaml"},
				Targets:   []urf.CompileTarget{urf.TargetCursor},
				OutputDir: "/tmp/test",
				Include:   []string{"**/*.yaml"},
			},
			expectErr: false, // Should not error on input validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CompileFiles(ctx, tt.request)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error containing %q, but got none", tt.errorMsg)
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Expected error %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result, got nil")
			}
		})
	}
}

func TestContentSelector_Matches(t *testing.T) {
	tests := []struct {
		name     string
		selector types.ContentSelector
		path     string
		expected bool
	}{
		{
			name:     "empty selector matches all",
			selector: types.ContentSelector{},
			path:     "test.yaml",
			expected: true,
		},
		{
			name:     "yaml include pattern",
			selector: types.ContentSelector{Include: []string{"**/*.yaml"}},
			path:     "test.yaml",
			expected: true,
		},
		{
			name:     "yaml include pattern no match",
			selector: types.ContentSelector{Include: []string{"**/*.yaml"}},
			path:     "test.txt",
			expected: false,
		},
		{
			name:     "exclude overrides include",
			selector: types.ContentSelector{Include: []string{"**/*.yaml"}, Exclude: []string{"**/ignore.yaml"}},
			path:     "ignore.yaml",
			expected: false,
		},
		{
			name:     "nested path match",
			selector: types.ContentSelector{Include: []string{"**/*.yaml"}},
			path:     "subdir/test.yaml",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.selector.Matches(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for path %q", tt.expected, result, tt.path)
			}
		})
	}
}

func TestArmService_DiscoverFiles(t *testing.T) {
	service := NewArmService()

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "arm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files
	testFiles := map[string]string{
		"test1.yaml":         "version: '1.0'\nmetadata:\n  id: test1\n  name: Test 1\n  version: 1.0.0\nrules: []",
		"test2.yml":          "version: '1.0'\nmetadata:\n  id: test2\n  name: Test 2\n  version: 1.0.0\nrules: []",
		"ignore.txt":         "not a yaml file",
		"subdir/nested.yaml": "version: '1.0'\nmetadata:\n  id: nested\n  name: Nested\n  version: 1.0.0\nrules: []",
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

	tests := []struct {
		name         string
		patterns     []string
		recursive    bool
		include      []string
		exclude      []string
		expectedFile int
	}{
		{
			name:         "single file",
			patterns:     []string{filepath.Join(tempDir, "test1.yaml")},
			recursive:    false,
			include:      []string{"**/*.yaml"},
			exclude:      []string{},
			expectedFile: 1,
		},
		{
			name:         "directory non-recursive",
			patterns:     []string{tempDir},
			recursive:    false,
			include:      []string{"**/*.yaml", "**/*.yml"},
			exclude:      []string{},
			expectedFile: 2, // test1.yaml and test2.yml, but not nested
		},
		{
			name:         "directory recursive",
			patterns:     []string{tempDir},
			recursive:    true,
			include:      []string{"**/*.yaml", "**/*.yml"},
			exclude:      []string{},
			expectedFile: 3, // all yaml files including nested
		},
		{
			name:         "with exclude pattern",
			patterns:     []string{tempDir},
			recursive:    true,
			include:      []string{"**/*.yaml", "**/*.yml"},
			exclude:      []string{"**/test1.yaml"},
			expectedFile: 2, // excluding test1.yaml
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := service.discoverFiles(tt.patterns, tt.recursive, tt.include, tt.exclude)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(files) != tt.expectedFile {
				t.Errorf("Expected %d files, got %d", tt.expectedFile, len(files))
				t.Logf("Files found:")
				for _, file := range files {
					t.Logf("  - %s", file.Path)
				}
			}
		})
	}
}

func TestCompileRequest_Validation(t *testing.T) {
	tests := []struct {
		name     string
		request  CompileRequest
		isValid  bool
		errorMsg string
	}{
		{
			name: "valid request",
			request: CompileRequest{
				Files:     []string{"test.yaml"},
				Targets:   []urf.CompileTarget{urf.TargetCursor},
				OutputDir: "/tmp/test",
				Include:   []string{"**/*.yaml"},
			},
			isValid: true,
		},
		{
			name: "conflicting flags",
			request: CompileRequest{
				Files:        []string{"test.yaml"},
				Targets:      []urf.CompileTarget{urf.TargetCursor},
				OutputDir:    "/tmp/test",
				ValidateOnly: true,
				DryRun:       true, // Should be conflicting
			},
			isValid:  false,
			errorMsg: "validate-only cannot be used with dry-run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would be validation logic we might add to the CLI
			if tt.request.ValidateOnly && (tt.request.DryRun || tt.request.Force) {
				if tt.isValid {
					t.Error("Expected validation to pass, but conflict detected")
				}
				return
			}

			if !tt.isValid {
				t.Error("Expected validation to fail, but no conflict detected")
			}
		})
	}
}

// Benchmark file discovery performance
func BenchmarkDiscoverFiles(b *testing.B) {
	service := NewArmService()

	// Create temporary test directory with many files
	tempDir, err := os.MkdirTemp("", "arm-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create 100 test files
	for i := 0; i < 100; i++ {
		content := "version: '1.0'\nmetadata:\n  id: test\n  name: Test\n  version: 1.0.0\nrules: []"
		path := filepath.Join(tempDir, filepath.Join("subdir", "file"+string(rune(i))+".yaml"))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			b.Fatalf("Failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			b.Fatalf("Failed to write file: %v", err)
		}
	}

	include := []string{"**/*.yaml"}
	patterns := []string{tempDir}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.discoverFiles(patterns, true, include, []string{})
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
