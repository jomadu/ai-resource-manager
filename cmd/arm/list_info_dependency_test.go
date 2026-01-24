package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInfoDependency(t *testing.T) {
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
		setupLockfile  string
		args           []string
		expectedOutput []string
	}{
		{
			name:           "no dependencies",
			setupManifest:  `{"version":1}`,
			setupLockfile:  `{"version":1}`,
			args:           []string{},
			expectedOutput: []string{"No dependencies configured"},
		},
		{
			name: "single ruleset dependency",
			setupManifest: `{
				"version": 1,
				"dependencies": {
					"sample-registry/clean-code-ruleset": {
						"type": "ruleset",
						"version": "^1.0.0",
						"priority": 100,
						"sinks": ["cursor-rules", "amazonq-rules"],
						"include": ["**/*.yml"],
						"exclude": ["**/experimental/**"]
					}
				}
			}`,
			setupLockfile: `{
				"version": 1,
				"dependencies": {
					"sample-registry/clean-code-ruleset@1.0.0": {
						"version": "1.0.0"
					}
				}
			}`,
			args: []string{"sample-registry/clean-code-ruleset"},
			expectedOutput: []string{
				"sample-registry/clean-code-ruleset:",
				"type: ruleset",
				"version: 1.0.0",
				"constraint: ^1.0.0",
				"priority: 100",
				"sinks:",
				"- cursor-rules",
				"- amazonq-rules",
				"include:",
				`"**/*.yml"`,
				"exclude:",
				`"**/experimental/**"`,
			},
		},
		{
			name: "single promptset dependency",
			setupManifest: `{
				"version": 1,
				"dependencies": {
					"sample-registry/code-review-promptset": {
						"type": "promptset",
						"version": "^1.0.0",
						"sinks": ["cursor-commands", "amazonq-prompts"],
						"include": ["review/**/*.yml"]
					}
				}
			}`,
			setupLockfile: `{
				"version": 1,
				"dependencies": {
					"sample-registry/code-review-promptset@1.0.0": {
						"version": "1.0.0"
					}
				}
			}`,
			args: []string{"sample-registry/code-review-promptset"},
			expectedOutput: []string{
				"sample-registry/code-review-promptset:",
				"type: promptset",
				"version: 1.0.0",
				"constraint: ^1.0.0",
				"sinks:",
				"- cursor-commands",
				"- amazonq-prompts",
				"include:",
				`"review/**/*.yml"`,
			},
		},
		{
			name: "all dependencies",
			setupManifest: `{
				"version": 1,
				"dependencies": {
					"sample-registry/clean-code-ruleset": {
						"type": "ruleset",
						"version": "^1.0.0",
						"priority": 100,
						"sinks": ["cursor-rules"]
					},
					"sample-registry/code-review-promptset": {
						"type": "promptset",
						"version": "^1.0.0",
						"sinks": ["cursor-commands"]
					}
				}
			}`,
			setupLockfile: `{
				"version": 1,
				"dependencies": {
					"sample-registry/clean-code-ruleset@1.0.0": {
						"version": "1.0.0"
					},
					"sample-registry/code-review-promptset@1.0.0": {
						"version": "1.0.0"
					}
				}
			}`,
			args: []string{},
			expectedOutput: []string{
				"sample-registry/clean-code-ruleset:",
				"type: ruleset",
				"sample-registry/code-review-promptset:",
				"type: promptset",
			},
		},
		{
			name: "multiple specific dependencies",
			setupManifest: `{
				"version": 1,
				"dependencies": {
					"reg1/pkg1": {
						"type": "ruleset",
						"version": "^1.0.0",
						"sinks": ["sink1"]
					},
					"reg2/pkg2": {
						"type": "promptset",
						"version": "^2.0.0",
						"sinks": ["sink2"]
					},
					"reg3/pkg3": {
						"type": "ruleset",
						"version": "^3.0.0",
						"sinks": ["sink3"]
					}
				}
			}`,
			setupLockfile: `{
				"version": 1,
				"dependencies": {
					"reg1/pkg1@1.0.0": {"version": "1.0.0"},
					"reg2/pkg2@2.0.0": {"version": "2.0.0"},
					"reg3/pkg3@3.0.0": {"version": "3.0.0"}
				}
			}`,
			args: []string{"reg1/pkg1", "reg3/pkg3"},
			expectedOutput: []string{
				"reg1/pkg1:",
				"type: ruleset",
				"version: 1.0.0",
				"reg3/pkg3:",
				"type: ruleset",
				"version: 3.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test manifest and lockfile
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")
			lockfilePath := filepath.Join(testDir, "arm-lock.json")

			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			if err := os.WriteFile(lockfilePath, []byte(tt.setupLockfile), 0o644); err != nil {
				t.Fatalf("Failed to write lockfile: %v", err)
			}

			// Run command
			cmdArgs := []string{"info", "dependency"}
			cmdArgs = append(cmdArgs, tt.args...)
			cmd := exec.Command(binaryPath, cmdArgs...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}

			outputStr := strings.TrimSpace(string(output))

			// Check all expected strings are present
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, outputStr)
				}
			}
		})
	}
}

func TestInfoDependencyInvalidFormat(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Create test manifest
	testDir := t.TempDir()
	manifestPath := filepath.Join(testDir, "arm.json")
	lockfilePath := filepath.Join(testDir, "arm-lock.json")

	manifest := `{
		"version": 1,
		"dependencies": {
			"sample-registry/clean-code-ruleset": {
				"type": "ruleset",
				"version": "^1.0.0",
				"sinks": ["cursor-rules"]
			}
		}
	}`

	lockfile := `{
		"version": 1,
		"dependencies": {
			"sample-registry/clean-code-ruleset@1.0.0": {
				"version": "1.0.0"
			}
		}
	}`

	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	if err := os.WriteFile(lockfilePath, []byte(lockfile), 0o644); err != nil {
		t.Fatalf("Failed to write lockfile: %v", err)
	}

	// Run command with invalid format
	cmd = exec.Command(binaryPath, "info", "dependency", "invalid-format")
	cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
	output, _ := cmd.CombinedOutput()

	// Should print error message about invalid format
	outputStr := string(output)
	if !strings.Contains(outputStr, "Invalid dependency format") {
		t.Errorf("Expected error message about invalid format, got: %s", outputStr)
	}
}
