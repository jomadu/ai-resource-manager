package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestListAll(t *testing.T) {
	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		setupManifest  string
		expectedOutput []string
	}{
		{
			name: "empty manifest",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"dependencies": {}
			}`,
			expectedOutput: []string{
				"Registries:",
				"(none)",
				"Sinks:",
				"(none)",
				"Dependencies:",
				"(none)",
			},
		},
		{
			name: "with registries and sinks",
			setupManifest: `{
				"version": 1,
				"registries": {
					"test-git": {
						"type": "git",
						"url": "https://github.com/test/repo"
					}
				},
				"sinks": {
					"cursor-rules": {
						"tool": "cursor",
						"directory": ".cursor/rules"
					}
				},
				"dependencies": {}
			}`,
			expectedOutput: []string{
				"Registries:",
				"test-git",
				"Sinks:",
				"cursor-rules",
				"Dependencies:",
				"(none)",
			},
		},
		{
			name: "with dependencies",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"dependencies": {
					"test-git/clean-code": {
						"type": "ruleset",
						"version": "^1.0.0",
						"priority": 100,
						"sinks": ["cursor-rules"]
					},
					"test-git/prompts": {
						"type": "promptset",
						"version": "^1.0.0",
						"sinks": ["cursor-commands"]
					}
				}
			}`,
			expectedOutput: []string{
				"Dependencies:",
				"test-git/clean-code (ruleset)",
				"test-git/prompts (promptset)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp manifest
			manifestPath := filepath.Join(tmpDir, "arm-"+tt.name+".json")
			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			// Run command
			cmd := exec.Command(binaryPath, "list")
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain %q, got:\n%s", expected, outputStr)
				}
			}
		})
	}
}

func TestInfoAll(t *testing.T) {
	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		setupManifest  string
		expectedOutput []string
	}{
		{
			name: "empty manifest",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"dependencies": {}
			}`,
			expectedOutput: []string{
				"Registries:",
				"(none)",
				"Sinks:",
				"(none)",
				"Dependencies:",
				"(none)",
			},
		},
		{
			name: "with all entities",
			setupManifest: `{
				"version": 1,
				"registries": {
					"test-git": {
						"type": "git",
						"url": "https://github.com/test/repo"
					}
				},
				"sinks": {
					"cursor-rules": {
						"tool": "cursor",
						"directory": ".cursor/rules"
					}
				},
				"dependencies": {
					"test-git/clean-code": {
						"type": "ruleset",
						"version": "^1.0.0",
						"priority": 100,
						"sinks": ["cursor-rules"]
					}
				}
			}`,
			expectedOutput: []string{
				"Registries:",
				"test-git:",
				"type: git",
				"url: https://github.com/test/repo",
				"Sinks:",
				"cursor-rules:",
				"tool: cursor",
				"directory: .cursor/rules",
				"Dependencies:",
				"test-git/clean-code:",
				"type: ruleset",
				"version: ^1.0.0",
				"priority: 100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp manifest
			manifestPath := filepath.Join(tmpDir, "arm-"+tt.name+".json")
			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			// Run command
			cmd := exec.Command(binaryPath, "info")
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain %q, got:\n%s", expected, outputStr)
				}
			}
		})
	}
}
