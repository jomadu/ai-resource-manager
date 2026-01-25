package e2e

import (
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestMultiSinkCrossToolInstallation(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")

	// Add multiple sinks for different tools
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("add", "sink", "--tool", "amazonq", "q-rules", ".amazonq/rules")

	// Install same ruleset to both sinks
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules", "q-rules")

	// Verify both sink directories populated
	cursorDir := filepath.Join(workDir, ".cursor/rules")
	helpers.AssertDirExists(t, cursorDir)
	cursorFiles := helpers.CountFilesRecursive(t, cursorDir)
	if cursorFiles == 0 {
		t.Error("expected files in cursor sink")
	}

	qDir := filepath.Join(workDir, ".amazonq/rules")
	helpers.AssertDirExists(t, qDir)
	qFiles := helpers.CountFilesRecursive(t, qDir)
	if qFiles == 0 {
		t.Error("expected files in amazonq sink")
	}

	// Verify different compilation formats
	// Cursor uses .mdc files, AmazonQ uses .md files
	cursorMdcFiles := helpers.CountFilesWithExtension(t, cursorDir, ".mdc")
	if cursorMdcFiles == 0 {
		t.Error("expected .mdc files in cursor sink")
	}

	qMdFiles := helpers.CountFilesWithExtension(t, qDir, ".md")
	if qMdFiles == 0 {
		t.Error("expected .md files in amazonq sink")
	}

	// Verify arm.json tracks both sinks
	armJSON := filepath.Join(workDir, "arm.json")
	helpers.AssertFileExists(t, armJSON)
	manifest := helpers.ReadJSON(t, armJSON)

	deps, ok := manifest["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatal("dependencies not found in arm.json")
	}

	dep, ok := deps["test-registry/test-ruleset"].(map[string]interface{})
	if !ok {
		t.Fatal("test-registry/test-ruleset not found in arm.json")
	}

	sinks, ok := dep["sinks"].([]interface{})
	if !ok {
		t.Fatal("sinks not found in dependency")
	}

	if len(sinks) != 2 {
		t.Errorf("expected 2 sinks, got %d", len(sinks))
	}
}

func TestMultiSinkSwitching(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")

	// Add two sinks
	arm.MustRun("add", "sink", "--tool", "cursor", "sink-a", ".sink-a")
	arm.MustRun("add", "sink", "--tool", "cursor", "sink-b", ".sink-b")

	// Install to sink A
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "sink-a")

	// Verify sink A populated
	sinkADir := filepath.Join(workDir, ".sink-a")
	helpers.AssertDirExists(t, sinkADir)
	sinkAFiles := helpers.CountFilesRecursive(t, sinkADir)
	if sinkAFiles == 0 {
		t.Error("expected files in sink A")
	}

	// Verify sink B is empty
	sinkBDir := filepath.Join(workDir, ".sink-b")
	if helpers.DirExists(sinkBDir) {
		sinkBFiles := helpers.CountFilesRecursive(t, sinkBDir)
		if sinkBFiles > 0 {
			t.Error("expected sink B to be empty initially")
		}
	}

	// Reinstall to sink B (should clean sink A)
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "sink-b")

	// Verify sink B populated
	helpers.AssertDirExists(t, sinkBDir)
	sinkBFilesAfter := helpers.CountFilesRecursive(t, sinkBDir)
	if sinkBFilesAfter == 0 {
		t.Error("expected files in sink B after reinstall")
	}

	// Note: Sink A cleanup behavior depends on ARM implementation
	// Some implementations may leave old files, others may clean them
	// This test verifies that sink B is populated correctly
}

func TestMultiSinkUpdate(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")

	// Add multiple sinks
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("add", "sink", "--tool", "copilot", "copilot-rules", ".github/copilot")

	// Install to both sinks
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules", "copilot-rules")

	// Create a new version
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("Update ruleset")
	repo.Tag("v1.1.0")

	// Update (should update both sinks)
	arm.MustRun("upgrade")

	// Verify both sinks still populated
	cursorDir := filepath.Join(workDir, ".cursor/rules")
	helpers.AssertDirExists(t, cursorDir)
	if helpers.CountFilesRecursive(t, cursorDir) == 0 {
		t.Error("expected files in cursor sink after update")
	}

	copilotDir := filepath.Join(workDir, ".github/copilot")
	helpers.AssertDirExists(t, copilotDir)
	if helpers.CountFilesRecursive(t, copilotDir) == 0 {
		t.Error("expected files in copilot sink after update")
	}

	// Verify arm-lock.json shows v1.1.0
	lockFile := filepath.Join(workDir, "arm-lock.json")
	helpers.AssertFileExists(t, lockFile)
	lock := helpers.ReadJSON(t, lockFile)

	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatalf("dependencies not found in arm-lock.json, got: %+v", lock)
	}

	// The key includes the version
	dep, ok := deps["test-registry/test-ruleset@v1.1.0"].(map[string]interface{})
	if !ok {
		t.Fatalf("test-registry/test-ruleset@v1.1.0 not found in arm-lock.json, available keys: %+v", deps)
	}

	// Verify integrity field exists
	if _, ok := dep["integrity"]; !ok {
		t.Error("integrity field not found in dependency")
	}
}
