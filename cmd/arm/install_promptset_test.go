package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallPromptset(t *testing.T) {
	tmpBin := filepath.Join(t.TempDir(), "arm")
	cmd := exec.Command("go", "build", "-o", tmpBin, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing package spec",
			args:        []string{"install", "promptset"},
			wantErr:     true,
			errContains: "Package spec required",
		},
		{
			name:        "missing sinks",
			args:        []string{"install", "promptset", "registry/promptset"},
			wantErr:     true,
			errContains: "At least one sink required",
		},
		{
			name:        "invalid package spec format",
			args:        []string{"install", "promptset", "invalid", "sink1"},
			wantErr:     true,
			errContains: "invalid format",
		},
		{
			name:        "empty registry name",
			args:        []string{"install", "promptset", "/promptset", "sink1"},
			wantErr:     true,
			errContains: "registry name cannot be empty",
		},
		{
			name:        "empty package name",
			args:        []string{"install", "promptset", "registry/", "sink1"},
			wantErr:     true,
			errContains: "package name cannot be empty",
		},
		{
			name:        "empty version after @",
			args:        []string{"install", "promptset", "registry/promptset@", "sink1"},
			wantErr:     true,
			errContains: "version cannot be empty after @",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			manifestPath := filepath.Join(tmpDir, "arm-manifest.json")

			cmd := exec.Command(tmpBin, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errContains != "" && !strings.Contains(string(output), tt.errContains) {
					t.Errorf("Expected error message to contain %q, got: %s", tt.errContains, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
			}
		})
	}
}
