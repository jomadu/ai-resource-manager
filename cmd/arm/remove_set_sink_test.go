package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveSink(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		setupManifest  string
		args           []string
		wantErr        bool
		wantOutput     string
		wantErrContain string
	}{
		{
			name: "remove existing sink",
			setupManifest: `{
				"sinks": {
					"test-sink": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			args:       []string{"remove", "sink", "test-sink"},
			wantOutput: "Removed sink 'test-sink'",
		},
		{
			name: "remove non-existent sink",
			setupManifest: `{
				"sinks": {}
			}`,
			args:           []string{"remove", "sink", "nonexistent"},
			wantErr:        true,
			wantErrContain: "sink nonexistent not found",
		},
		{
			name:           "missing sink name",
			setupManifest:  `{"sinks": {}}`,
			args:           []string{"remove", "sink"},
			wantErr:        true,
			wantErrContain: "Usage: arm remove sink NAME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.wantErrContain != "" && !strings.Contains(string(output), tt.wantErrContain) {
					t.Errorf("Expected error containing %q, got: %s", tt.wantErrContain, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
				if tt.wantOutput != "" && !strings.Contains(string(output), tt.wantOutput) {
					t.Errorf("Expected output containing %q, got: %s", tt.wantOutput, output)
				}
			}
		})
	}
}

func TestSetSink(t *testing.T) {
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		setupManifest  string
		args           []string
		wantErr        bool
		wantOutput     string
		wantErrContain string
	}{
		{
			name: "set sink tool",
			setupManifest: `{
				"sinks": {
					"test-sink": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			args:       []string{"set", "sink", "test-sink", "tool", "amazonq"},
			wantOutput: "Updated sink 'test-sink' tool",
		},
		{
			name: "set sink directory",
			setupManifest: `{
				"sinks": {
					"test-sink": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			args:       []string{"set", "sink", "test-sink", "directory", ".amazonq/rules"},
			wantOutput: "Updated sink 'test-sink' directory",
		},
		{
			name: "set invalid tool",
			setupManifest: `{
				"sinks": {
					"test-sink": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			args:           []string{"set", "sink", "test-sink", "tool", "invalid"},
			wantErr:        true,
			wantErrContain: "Invalid tool",
		},
		{
			name: "set unknown key",
			setupManifest: `{
				"sinks": {
					"test-sink": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			args:           []string{"set", "sink", "test-sink", "unknown", "value"},
			wantErr:        true,
			wantErrContain: "Unknown key",
		},
		{
			name: "set non-existent sink",
			setupManifest: `{
				"sinks": {}
			}`,
			args:           []string{"set", "sink", "nonexistent", "tool", "cursor"},
			wantErr:        true,
			wantErrContain: "sink nonexistent not found",
		},
		{
			name:           "missing arguments",
			setupManifest:  `{"sinks": {}}`,
			args:           []string{"set", "sink", "test-sink", "tool"},
			wantErr:        true,
			wantErrContain: "Usage: arm set sink NAME KEY VALUE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.wantErrContain != "" && !strings.Contains(string(output), tt.wantErrContain) {
					t.Errorf("Expected error containing %q, got: %s", tt.wantErrContain, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
				if tt.wantOutput != "" && !strings.Contains(string(output), tt.wantOutput) {
					t.Errorf("Expected output containing %q, got: %s", tt.wantOutput, output)
				}
			}
		})
	}
}
