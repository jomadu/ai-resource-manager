package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCleanCache(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	t.Logf("DEBUG: tmpDir=%s, binaryPath=%s", tmpDir, binaryPath)
	
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v, output: %s", err, string(output))
	}
	
	// Verify binary exists and is executable
	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Fatalf("Binary not found after build: %v", err)
	}
	t.Logf("DEBUG: binary size=%d, mode=%v", info.Size(), info.Mode())

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no subcommand",
			args:        []string{"clean"},
			expectError: true,
			errorMsg:    "clean requires a subcommand",
		},
		{
			name:        "unknown subcommand",
			args:        []string{"clean", "unknown"},
			expectError: true,
			errorMsg:    "Unknown clean subcommand",
		},
		{
			name:        "cache with default max-age",
			args:        []string{"clean", "cache"},
			expectError: false,
		},
		{
			name:        "cache with custom max-age",
			args:        []string{"clean", "cache", "--max-age", "30d"},
			expectError: false,
		},
		{
			name:        "cache with nuke",
			args:        []string{"clean", "cache", "--nuke"},
			expectError: false,
		},
		{
			name:        "cache with both max-age and nuke",
			args:        []string{"clean", "cache", "--max-age", "7d", "--nuke"},
			expectError: true,
			errorMsg:    "mutually exclusive",
		},
		{
			name:        "cache with max-age missing value",
			args:        []string{"clean", "cache", "--max-age"},
			expectError: true,
			errorMsg:    "--max-age requires a value",
		},
		{
			name:        "cache with invalid duration",
			args:        []string{"clean", "cache", "--max-age", "invalid"},
			expectError: true,
			errorMsg:    "invalid duration",
		},
		{
			name:        "cache with unknown flag",
			args:        []string{"clean", "cache", "--unknown"},
			expectError: true,
			errorMsg:    "Unknown flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			t.Logf("DEBUG: subtest=%s, exitCode=%v, output=%s", tt.name, err, string(output))

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errorMsg != "" && !contains(string(output), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got: %s", tt.errorMsg, output)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v. Output: %s", err, output)
			}
		})
	}
}
