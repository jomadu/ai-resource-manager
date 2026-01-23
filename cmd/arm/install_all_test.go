package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallAll(t *testing.T) {
	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Create test manifest with dependencies
	manifestPath := filepath.Join(tmpDir, "arm-manifest.json")
	manifestContent := `{
		"registries": {
			"test-registry": {
				"type": "git",
				"url": "https://github.com/test/repo"
			}
		},
		"sinks": {
			"test-sink": {
				"tool": "markdown",
				"directory": "` + filepath.Join(tmpDir, "output") + `"
			}
		},
		"dependencies": {
			"rulesets": {},
			"promptsets": {}
		}
	}`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0o644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	// Run install command
	cmd = exec.Command(binaryPath, "install")
	cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "All dependencies installed successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

func TestInstallAllNoManifest(t *testing.T) {
	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	manifestPath := filepath.Join(tmpDir, "nonexistent.json")

	// Run install command - should succeed with empty manifest
	cmd = exec.Command(binaryPath, "install")
	cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "All dependencies installed successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

func TestUninstall(t *testing.T) {
	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Create test manifest with dependencies
	manifestPath := filepath.Join(tmpDir, "arm-manifest.json")
	manifestContent := `{
		"registries": {
			"test-registry": {
				"type": "git",
				"url": "https://github.com/test/repo"
			}
		},
		"sinks": {
			"test-sink": {
				"tool": "markdown",
				"directory": "` + filepath.Join(tmpDir, "output") + `"
			}
		},
		"dependencies": {
			"rulesets": {},
			"promptsets": {}
		}
	}`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0o644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	// Run uninstall command
	cmd = exec.Command(binaryPath, "uninstall")
	cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "All packages uninstalled successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

func TestUninstallNoManifest(t *testing.T) {
	// Build the binary
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	manifestPath := filepath.Join(tmpDir, "nonexistent.json")

	// Run uninstall command - should succeed with empty manifest
	cmd = exec.Command(binaryPath, "uninstall")
	cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "All packages uninstalled successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}
}
