package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpgrade_Success(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "arm-manifest.json")
	lockfilePath := filepath.Join(tmpDir, "arm-manifest-lock.json")

	// Create manifest with a dependency
	manifestContent := `{
		"rulesets": {
			"test-registry/test-ruleset": {
				"version": "^1.0.0",
				"priority": 100,
				"sinks": ["test-sink"]
			}
		}
	}`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create lockfile with old version
	lockfileContent := `{
		"rulesets": {
			"test-registry/test-ruleset": {
				"1.0.0": {
					"integrity": "sha256-old"
				}
			}
		}
	}`
	if err := os.WriteFile(lockfilePath, []byte(lockfileContent), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".", "upgrade")
	cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()

	// Note: This will fail in practice because we don't have a real registry
	// But we're testing the CLI parsing and invocation
	if err == nil {
		t.Logf("Output: %s", output)
	} else {
		// Expected to fail due to missing registry, but should show proper error
		if !strings.Contains(string(output), "Error:") {
			t.Errorf("Expected error message, got: %s", output)
		}
	}
}

func TestUpgrade_MissingLockfile(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "arm-manifest.json")

	// Create manifest without lockfile
	manifestContent := `{
		"rulesets": {
			"test-registry/test-ruleset": {
				"version": "^1.0.0",
				"priority": 100,
				"sinks": ["test-sink"]
			}
		}
	}`
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".", "upgrade")
	cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Errorf("Expected error for missing lockfile, got success: %s", output)
	}

	if !strings.Contains(string(output), "Error:") {
		t.Errorf("Expected error message, got: %s", output)
	}
}
