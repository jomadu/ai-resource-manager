package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestIntegrityVerification_E2E(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Create a local git repository with a ruleset
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Add test ruleset")
	repo.Tag("v1.0.0")

	// Setup registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "test-sink", ".cursor/rules")

	// First install - should succeed and create lock file
	t.Run("FirstInstall", func(t *testing.T) {
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@v1.0.0", "test-sink")

		// Verify lock file was created with integrity
		lockFile := filepath.Join(workDir, "arm-lock.json")
		helpers.AssertFileExists(t, lockFile)

		lockData := helpers.ReadJSON(t, lockFile)
		deps, ok := lockData["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies field not found in lock file")
		}

		// Find the lock entry
		var lockEntry map[string]interface{}
		for key, val := range deps {
			if strings.Contains(key, "test-registry/test-ruleset@v1.0.0") {
				lockEntry = val.(map[string]interface{})
				break
			}
		}

		if lockEntry == nil {
			t.Fatal("Lock entry not found for test-registry/test-ruleset@v1.0.0")
		}

		integrity, ok := lockEntry["integrity"].(string)
		if !ok || integrity == "" {
			t.Fatal("Integrity should be stored in lock file")
		}

		if !strings.HasPrefix(integrity, "sha256-") {
			t.Errorf("Integrity should start with 'sha256-', got: %s", integrity)
		}

		t.Logf("First install created lock with integrity: %s", integrity)
	})

	// Second install - should succeed with same integrity
	t.Run("ReinstallSameContent", func(t *testing.T) {
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@v1.0.0", "test-sink")
	})

	// Modify the repository content (simulate tampering)
	t.Run("ModifiedContent", func(t *testing.T) {
		// Verify lock file still exists from previous install
		lockFile := filepath.Join(workDir, "arm-lock.json")
		if _, err := os.Stat(lockFile); os.IsNotExist(err) {
			t.Fatal("Lock file should exist from previous install")
		}

		// Read the current lock file to get the expected integrity
		lockDataBefore := helpers.ReadJSON(t, lockFile)
		depsBefore, ok := lockDataBefore["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies field not found in lock file")
		}

		var expectedIntegrity string
		for key, val := range depsBefore {
			if strings.Contains(key, "test-registry/test-ruleset@v1.0.0") {
				lockEntry := val.(map[string]interface{})
				expectedIntegrity = lockEntry["integrity"].(string)
				break
			}
		}

		if expectedIntegrity == "" {
			t.Fatal("Could not find expected integrity from lock file")
		}

		t.Logf("Expected integrity from lock: %s", expectedIntegrity)

		// Modify the ruleset file
		modifiedRuleset := `
version: 1
metadata:
  id: test-ruleset-modified
  name: Test Ruleset Modified
  description: MODIFIED - This has been tampered with
rules:
  - id: rule1-modified
    description: Modified test rule
    content: |
      This is a MODIFIED test rule with different content
`
		repo.WriteFile("test-ruleset.yml", modifiedRuleset)
		repo.Commit("Modify ruleset (simulate tampering)")

		// Force update the tag to point to the new commit
		cmd := exec.Command("git", "tag", "-f", "v1.0.0")
		cmd.Dir = repoDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to force update tag: %v", err)
		}

		// Clear the entire ARM storage to force re-fetch
		storagePath := filepath.Join(os.Getenv("HOME"), ".arm", "storage")
		if err := os.RemoveAll(storagePath); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to clear storage: %v", err)
		}

		// Try to install again - should fail with integrity verification error
		_, stderr, err := arm.Run("install", "ruleset", "test-registry/test-ruleset@v1.0.0", "test-sink")
		if err == nil {
			t.Fatal("Install should fail when content has been modified")
		}

		// Verify error message
		if !strings.Contains(stderr, "integrity verification failed") {
			t.Errorf("Error should mention integrity verification, got: %s", stderr)
		}
		if !strings.Contains(stderr, "Expected:") {
			t.Errorf("Error should show expected hash, got: %s", stderr)
		}
		if !strings.Contains(stderr, "Got:") {
			t.Errorf("Error should show actual hash, got: %s", stderr)
		}
		if !strings.Contains(stderr, expectedIntegrity) {
			t.Errorf("Error should show expected integrity %s, got: %s", expectedIntegrity, stderr)
		}

		t.Logf("Integrity verification correctly detected tampering")
	})

	// Delete lock file and reinstall - should succeed
	t.Run("ReinstallAfterLockDelete", func(t *testing.T) {
		// Delete lock file
		lockFile := filepath.Join(workDir, "arm-lock.json")
		if err := os.Remove(lockFile); err != nil {
			t.Fatal(err)
		}

		// Install should succeed now (no lock to verify against)
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@v1.0.0", "test-sink")

		// Verify new lock file was created with new integrity
		helpers.AssertFileExists(t, lockFile)
		lockData := helpers.ReadJSON(t, lockFile)
		deps, ok := lockData["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies field not found in lock file")
		}

		// Find the lock entry
		var lockEntry map[string]interface{}
		for key, val := range deps {
			if strings.Contains(key, "test-registry/test-ruleset@v1.0.0") {
				lockEntry = val.(map[string]interface{})
				break
			}
		}

		if lockEntry == nil {
			t.Fatal("Lock entry not found")
		}

		integrity, ok := lockEntry["integrity"].(string)
		if !ok || integrity == "" {
			t.Fatal("Integrity should be stored in new lock file")
		}

		t.Logf("New lock file created with integrity: %s", integrity)
	})
}

func TestIntegrityVerification_BackwardsCompatibility(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Create a local git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Add ruleset")
	repo.Tag("v1.0.0")

	// Setup registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "test-sink", ".cursor/rules")

	// Create a legacy lock file without integrity field
	legacyLock := map[string]interface{}{
		"version": 1,
		"dependencies": map[string]interface{}{
			"test-registry/test-ruleset@v1.0.0": map[string]interface{}{
				"integrity": "", // Empty integrity (legacy)
			},
		},
	}

	lockData, err := json.MarshalIndent(legacyLock, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	lockFile := filepath.Join(workDir, "arm-lock.json")
	if err := os.WriteFile(lockFile, lockData, 0o644); err != nil {
		t.Fatal(err)
	}

	// Install should succeed even with legacy lock file (backwards compatibility)
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@v1.0.0", "test-sink")

	// Verify integrity was added to lock file
	lockData2 := helpers.ReadJSON(t, lockFile)
	deps, ok := lockData2["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatal("dependencies field not found")
	}

	var lockEntry map[string]interface{}
	for key, val := range deps {
		if strings.Contains(key, "test-registry/test-ruleset@v1.0.0") {
			lockEntry = val.(map[string]interface{})
			break
		}
	}

	if lockEntry == nil {
		t.Fatal("Lock entry not found")
	}

	integrity, ok := lockEntry["integrity"].(string)
	if !ok || integrity == "" {
		t.Error("Integrity should be updated in lock file after install")
	}

	t.Logf("Legacy lock file upgraded with integrity: %s", integrity)
}
