package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestListSink(t *testing.T) {
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
		expectedOutput string
		expectError    bool
	}{
		{
			name: "list sinks - multiple sinks",
			setupManifest: `{
				"sinks": {
					"cursor-rules": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					},
					"q-rules": {
						"directory": ".amazonq/rules",
						"tool": "amazonq"
					}
				}
			}`,
			expectedOutput: "cursor-rules\nq-rules",
			expectError:    false,
		},
		{
			name: "list sinks - single sink",
			setupManifest: `{
				"sinks": {
					"cursor-rules": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			expectedOutput: "cursor-rules",
			expectError:    false,
		},
		{
			name:           "list sinks - no sinks",
			setupManifest:  `{"sinks": {}}`,
			expectedOutput: "No sinks configured",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			// Write manifest
			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			// Run command
			cmd := exec.Command(binaryPath, "list", "sink")
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v\nOutput: %s", err, output)
			}

			outputStr := strings.TrimSpace(string(output))
			if !strings.Contains(outputStr, tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, outputStr)
			}
		})
	}
}

func TestInfoSink(t *testing.T) {
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
		expectedOutput []string
		expectError    bool
	}{
		{
			name: "info sink - specific sink",
			setupManifest: `{
				"sinks": {
					"cursor-rules": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					},
					"q-rules": {
						"directory": ".amazonq/rules",
						"tool": "amazonq"
					}
				}
			}`,
			args: []string{"cursor-rules"},
			expectedOutput: []string{
				"Sink: cursor-rules",
				"Tool: cursor",
				"Directory: .cursor/rules",
			},
			expectError: false,
		},
		{
			name: "info sink - multiple sinks",
			setupManifest: `{
				"sinks": {
					"cursor-rules": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					},
					"q-rules": {
						"directory": ".amazonq/rules",
						"tool": "amazonq"
					}
				}
			}`,
			args: []string{"cursor-rules", "q-rules"},
			expectedOutput: []string{
				"Sink: cursor-rules",
				"Tool: cursor",
				"Directory: .cursor/rules",
				"Sink: q-rules",
				"Tool: amazonq",
				"Directory: .amazonq/rules",
			},
			expectError: false,
		},
		{
			name: "info sink - all sinks (no args)",
			setupManifest: `{
				"sinks": {
					"cursor-rules": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			args: []string{},
			expectedOutput: []string{
				"Sink: cursor-rules",
				"Tool: cursor",
				"Directory: .cursor/rules",
			},
			expectError: false,
		},
		{
			name:           "info sink - no sinks",
			setupManifest:  `{"sinks": {}}`,
			args:           []string{},
			expectedOutput: []string{"No sinks configured"},
			expectError:    false,
		},
		{
			name: "info sink - nonexistent sink",
			setupManifest: `{
				"sinks": {
					"cursor-rules": {
						"directory": ".cursor/rules",
						"tool": "cursor"
					}
				}
			}`,
			args:           []string{"nonexistent"},
			expectedOutput: []string{"Error getting sink 'nonexistent'"},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			// Write manifest
			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			// Run command
			cmdArgs := append([]string{"info", "sink"}, tt.args...)
			cmd := exec.Command(binaryPath, cmdArgs...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v\nOutput: %s", err, output)
			}

			outputStr := strings.TrimSpace(string(output))
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain %q, got %q", expected, outputStr)
				}
			}
		})
	}
}
