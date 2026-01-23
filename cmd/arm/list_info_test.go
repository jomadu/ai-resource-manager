package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestListRegistry(t *testing.T) {
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
		expectedOutput []string
	}{
		{
			name:           "no registries",
			setupManifest:  `{"version":1}`,
			expectedOutput: []string{"No registries configured"},
		},
		{
			name: "single registry",
			setupManifest: `{
				"version": 1,
				"registries": {
					"my-registry": {
						"type": "git",
						"url": "https://github.com/test/repo"
					}
				}
			}`,
			expectedOutput: []string{"my-registry"},
		},
		{
			name: "multiple registries",
			setupManifest: `{
				"version": 1,
				"registries": {
					"git-reg": {
						"type": "git",
						"url": "https://github.com/test/repo"
					},
					"gitlab-reg": {
						"type": "gitlab",
						"url": "https://gitlab.com/test/repo"
					}
				}
			}`,
			expectedOutput: []string{"git-reg", "gitlab-reg"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")
			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			// Run command
			cmd := exec.Command(binaryPath, "list", "registry")
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

func TestInfoRegistry(t *testing.T) {
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
	}{
		{
			name:           "no registries",
			setupManifest:  `{"version":1}`,
			args:           []string{},
			expectedOutput: []string{"No registries configured"},
		},
		{
			name: "git registry info",
			setupManifest: `{
				"version": 1,
				"registries": {
					"my-git": {
						"type": "git",
						"url": "https://github.com/test/repo",
						"branches": ["main", "develop"]
					}
				}
			}`,
			args: []string{"my-git"},
			expectedOutput: []string{
				"Registry: my-git",
				"Type: git",
				"URL: https://github.com/test/repo",
				"Branches: main, develop",
			},
		},
		{
			name: "gitlab registry info",
			setupManifest: `{
				"version": 1,
				"registries": {
					"my-gitlab": {
						"type": "gitlab",
						"url": "https://gitlab.com/test/repo",
						"projectId": "123",
						"groupId": "456",
						"apiVersion": "v4"
					}
				}
			}`,
			args: []string{"my-gitlab"},
			expectedOutput: []string{
				"Registry: my-gitlab",
				"Type: gitlab",
				"URL: https://gitlab.com/test/repo",
				"Project ID: 123",
				"Group ID: 456",
				"API Version: v4",
			},
		},
		{
			name: "cloudsmith registry info",
			setupManifest: `{
				"version": 1,
				"registries": {
					"my-cloudsmith": {
						"type": "cloudsmith",
						"url": "https://cloudsmith.io",
						"owner": "myorg",
						"repository": "myrepo"
					}
				}
			}`,
			args: []string{"my-cloudsmith"},
			expectedOutput: []string{
				"Registry: my-cloudsmith",
				"Type: cloudsmith",
				"URL: https://cloudsmith.io",
				"Owner: myorg",
				"Repository: myrepo",
			},
		},
		{
			name: "all registries info",
			setupManifest: `{
				"version": 1,
				"registries": {
					"git-reg": {
						"type": "git",
						"url": "https://github.com/test/repo"
					},
					"gitlab-reg": {
						"type": "gitlab",
						"url": "https://gitlab.com/test/repo"
					}
				}
			}`,
			args: []string{},
			expectedOutput: []string{
				"Registry: git-reg",
				"Type: git",
				"Registry: gitlab-reg",
				"Type: gitlab",
			},
		},
		{
			name: "multiple specific registries",
			setupManifest: `{
				"version": 1,
				"registries": {
					"reg1": {
						"type": "git",
						"url": "https://github.com/test/repo1"
					},
					"reg2": {
						"type": "git",
						"url": "https://github.com/test/repo2"
					},
					"reg3": {
						"type": "git",
						"url": "https://github.com/test/repo3"
					}
				}
			}`,
			args: []string{"reg1", "reg3"},
			expectedOutput: []string{
				"Registry: reg1",
				"https://github.com/test/repo1",
				"Registry: reg3",
				"https://github.com/test/repo3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")
			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			// Run command
			cmdArgs := []string{"info", "registry"}
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
