package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	// Build the binary once for all tests
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, output)
	}

	// Create test working directory with sample input files
	workDir := t.TempDir()

	// Create sample ARM resource files for testing
	sampleRuleset := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset
  name: Test Ruleset
spec:
  rules:
    testRule:
      name: Test Rule
      body: |
        This is a test rule
`

	sampleRuleset1 := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset-1
  name: Test Ruleset 1
spec:
  rules:
    testRule1:
      name: Test Rule 1
      body: |
        This is test rule 1
`

	sampleRuleset2 := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset-2
  name: Test Ruleset 2
spec:
  rules:
    testRule2:
      name: Test Rule 2
      body: |
        This is test rule 2
`

	if err := os.WriteFile(filepath.Join(workDir, "input.yml"), []byte(sampleRuleset), 0o644); err != nil {
		t.Fatalf("Failed to create input.yml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workDir, "input1.yml"), []byte(sampleRuleset1), 0o644); err != nil {
		t.Fatalf("Failed to create input1.yml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workDir, "input2.yml"), []byte(sampleRuleset2), 0o644); err != nil {
		t.Fatalf("Failed to create input2.yml: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:        "no arguments",
			args:        []string{"compile"},
			wantErr:     true,
			errContains: "at least one INPUT_PATH is required",
		},
		{
			name:        "missing output path without validate-only",
			args:        []string{"compile", "input.yml"},
			wantErr:     true,
			errContains: "OUTPUT_PATH is required",
		},
		{
			name:    "validate-only without output path",
			args:    []string{"compile", "--validate-only", "input.yml"},
			wantErr: false,
		},
		{
			name:    "with output path",
			args:    []string{"compile", "--tool", "cursor", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "with tool flag",
			args:    []string{"compile", "--tool", "cursor", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "with namespace flag",
			args:    []string{"compile", "--tool", "cursor", "--namespace", "test", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "with force flag",
			args:    []string{"compile", "--tool", "cursor", "--force", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "with recursive flag",
			args:    []string{"compile", "--tool", "cursor", "--recursive", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "with include flag",
			args:    []string{"compile", "--tool", "cursor", "--include", "*.yml", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "with exclude flag",
			args:    []string{"compile", "--tool", "cursor", "--exclude", "test*", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "with fail-fast flag",
			args:    []string{"compile", "--tool", "cursor", "--fail-fast", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "multiple input paths",
			args:    []string{"compile", "--tool", "cursor", "input1.yml", "input2.yml", "output"},
			wantErr: false,
		},
		{
			name:    "multiple include patterns",
			args:    []string{"compile", "--tool", "cursor", "--include", "*.yml", "--include", "*.yaml", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "multiple exclude patterns",
			args:    []string{"compile", "--tool", "cursor", "--exclude", "test*", "--exclude", "tmp*", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:    "all flags combined",
			args:    []string{"compile", "--tool", "amazonq", "--namespace", "test", "--force", "--recursive", "--include", "*.yml", "--exclude", "test*", "--fail-fast", "input.yml", "output"},
			wantErr: false,
		},
		{
			name:        "unknown flag",
			args:        []string{"compile", "--unknown", "input.yml", "output"},
			wantErr:     true,
			errContains: "Unknown flag",
		},
		{
			name:        "tool flag without value",
			args:        []string{"compile", "--tool"},
			wantErr:     true,
			errContains: "--tool requires a value",
		},
		{
			name:        "namespace flag without value",
			args:        []string{"compile", "--namespace"},
			wantErr:     true,
			errContains: "--namespace requires a value",
		},
		{
			name:        "include flag without value",
			args:        []string{"compile", "--include"},
			wantErr:     true,
			errContains: "--include requires a value",
		},
		{
			name:        "exclude flag without value",
			args:        []string{"compile", "--exclude"},
			wantErr:     true,
			errContains: "--exclude requires a value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create unique output directory for this subtest
			subtestDir := t.TempDir()

			// Replace "output" in args with unique output path
			args := make([]string, len(tt.args))
			for i, arg := range tt.args {
				if arg == "output" {
					args[i] = filepath.Join(subtestDir, "output")
				} else {
					args[i] = arg
				}
			}

			cmd := exec.Command(binaryPath, args...)
			cmd.Dir = workDir // Set working directory to where test files are
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errContains != "" && !strings.Contains(string(output), tt.errContains) {
					t.Errorf("Expected error containing %q, got: %s", tt.errContains, output)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v\nOutput: %s", err, output)
			}
		})
	}
}

func TestCompileHelp(t *testing.T) {
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, output)
	}

	cmd = exec.Command(binaryPath, "help", "compile")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help compile failed: %v\n%s", err, output)
	}

	expectedStrings := []string{
		"Compile rulesets and promptsets",
		"arm compile INPUT_PATH",
		"--tool",
		"--namespace",
		"--force",
		"--recursive",
		"--validate-only",
		"--include",
		"--exclude",
		"--fail-fast",
	}

	outputStr := string(output)
	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Expected help output to contain %q, got:\n%s", expected, outputStr)
		}
	}
}
