package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveRegistry(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name          string
		setupManifest string
		args          []string
		wantErr       bool
		wantOutput    string
		checkManifest func(t *testing.T, manifestPath string)
	}{
		{
			name: "remove existing registry",
			setupManifest: `{
				"registries": {
					"test-registry": {
						"type": "git",
						"url": "https://github.com/test/repo"
					}
				}
			}`,
			args:       []string{"remove", "registry", "test-registry"},
			wantErr:    false,
			wantOutput: "Removed registry 'test-registry'",
			checkManifest: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				if strings.Contains(string(data), "test-registry") {
					t.Error("Registry still exists in manifest")
				}
			},
		},
		{
			name: "remove non-existent registry",
			setupManifest: `{
				"registries": {}
			}`,
			args:    []string{"remove", "registry", "nonexistent"},
			wantErr: true,
		},
		{
			name:    "missing registry name",
			args:    []string{"remove", "registry"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create isolated manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			if tt.setupManifest != "" {
				if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
					t.Fatalf("Failed to write manifest: %v", err)
				}
			}

			// Run command
			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v. Output: %s", err, output)
				return
			}

			if tt.wantOutput != "" && !strings.Contains(string(output), tt.wantOutput) {
				t.Errorf("Expected output to contain %q, got: %s", tt.wantOutput, output)
			}

			if tt.checkManifest != nil {
				tt.checkManifest(t, manifestPath)
			}
		})
	}
}

func TestSetRegistry(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name          string
		setupManifest string
		args          []string
		wantErr       bool
		wantOutput    string
		checkManifest func(t *testing.T, manifestPath string)
	}{
		{
			name: "set registry url",
			setupManifest: `{
				"registries": {
					"test-registry": {
						"type": "git",
						"url": "https://github.com/test/old"
					}
				}
			}`,
			args:       []string{"set", "registry", "test-registry", "url", "https://github.com/test/new"},
			wantErr:    false,
			wantOutput: "Updated registry 'test-registry' url",
			checkManifest: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				if !strings.Contains(string(data), "https://github.com/test/new") {
					t.Error("URL not updated in manifest")
				}
			},
		},
		{
			name: "set registry name",
			setupManifest: `{
				"registries": {
					"old-name": {
						"type": "git",
						"url": "https://github.com/test/repo"
					}
				}
			}`,
			args:       []string{"set", "registry", "old-name", "name", "new-name"},
			wantErr:    false,
			wantOutput: "Updated registry 'old-name' name",
			checkManifest: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				if !strings.Contains(string(data), "new-name") {
					t.Error("Name not updated in manifest")
				}
				if strings.Contains(string(data), "old-name") {
					t.Error("Old name still exists in manifest")
				}
			},
		},
		{
			name: "set non-existent registry",
			setupManifest: `{
				"registries": {}
			}`,
			args:    []string{"set", "registry", "nonexistent", "url", "https://github.com/test/repo"},
			wantErr: true,
		},
		{
			name: "set invalid key",
			setupManifest: `{
				"registries": {
					"test-registry": {
						"type": "git",
						"url": "https://github.com/test/repo"
					}
				}
			}`,
			args:    []string{"set", "registry", "test-registry", "invalid", "value"},
			wantErr: true,
		},
		{
			name:    "missing arguments",
			args:    []string{"set", "registry", "test-registry"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create isolated manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			if tt.setupManifest != "" {
				if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
					t.Fatalf("Failed to write manifest: %v", err)
				}
			}

			// Run command
			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v. Output: %s", err, output)
				return
			}

			if tt.wantOutput != "" && !strings.Contains(string(output), tt.wantOutput) {
				t.Errorf("Expected output to contain %q, got: %s", tt.wantOutput, output)
			}

			if tt.checkManifest != nil {
				tt.checkManifest(t, manifestPath)
			}
		})
	}
}
