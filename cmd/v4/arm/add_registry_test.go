package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddGitRegistry(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid git registry",
			args:    []string{"add", "registry", "git", "--url", "https://github.com/test/repo", "test-reg"},
			wantErr: false,
		},
		{
			name:    "with branches",
			args:    []string{"add", "registry", "git", "--url", "https://github.com/test/repo", "--branches", "main,dev", "test-reg2"},
			wantErr: false,
		},
		{
			name:        "missing url",
			args:        []string{"add", "registry", "git", "test-reg"},
			wantErr:     true,
			errContains: "--url is required",
		},
		{
			name:        "missing name",
			args:        []string{"add", "registry", "git", "--url", "https://github.com/test/repo"},
			wantErr:     true,
			errContains: "NAME is required",
		},
		{
			name:        "duplicate without force",
			args:        []string{"add", "registry", "git", "--url", "https://github.com/test/repo", "test-reg"},
			wantErr:     true,
			errContains: "registry already exists",
		},
		{
			name:    "duplicate with force",
			args:    []string{"add", "registry", "git", "--url", "https://github.com/test/repo2", "--force", "test-reg"},
			wantErr: false,
		},
	}

	// Build the binary once
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join(".") // Current directory is cmd/v4/arm
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Each test gets its own manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm-manifest.json")

			// For duplicate tests, create the registry first
			if strings.Contains(tt.name, "duplicate") {
				setupCmd := exec.Command(binPath, "add", "registry", "git", "--url", "https://github.com/test/repo", "test-reg")
				setupCmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
				if err := setupCmd.Run(); err != nil {
					t.Fatalf("Failed to setup duplicate test: %v", err)
				}
			}

			cmd := exec.Command(binPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errContains != "" && !strings.Contains(string(output), tt.errContains) {
					t.Errorf("Expected error containing %q, got: %s", tt.errContains, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
			}
		})
	}
}
