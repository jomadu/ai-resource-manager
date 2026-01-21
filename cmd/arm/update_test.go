package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateCommand(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "arm-manifest.json")

	// Build the binary inline
	binaryPath := filepath.Join(tmpDir, "arm")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		setupManifest  string
		setupLockfile  string
		expectedOutput string
		expectError    bool
	}{
		{
			name: "update with no manifest",
			setupManifest: "",
			setupLockfile: "",
			expectedOutput: "",
			expectError: true,
		},
		{
			name: "update with empty manifest",
			setupManifest: `{
				"registries": {},
				"sinks": {},
				"rulesets": {},
				"promptsets": {}
			}`,
			setupLockfile: `{
				"dependencies": {}
			}`,
			expectedOutput: "All packages updated successfully",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup manifest if provided
			if tt.setupManifest != "" {
				if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0644); err != nil {
					t.Fatalf("Failed to write manifest: %v", err)
				}
			} else {
				// Remove manifest if it exists
				os.Remove(manifestPath)
			}

			// Setup lockfile if provided
			lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"
			if tt.setupLockfile != "" {
				if err := os.WriteFile(lockfilePath, []byte(tt.setupLockfile), 0644); err != nil {
					t.Fatalf("Failed to write lockfile: %v", err)
				}
			} else {
				// Remove lockfile if it exists
				os.Remove(lockfilePath)
			}

			// Run update command
			cmd := exec.Command(binaryPath, "update")
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v\nOutput: %s", err, output)
			}

			if tt.expectedOutput != "" && !strings.Contains(string(output), tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got: %s", tt.expectedOutput, output)
			}
		})
	}
}
