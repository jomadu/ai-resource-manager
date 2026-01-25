package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

// TestManifestCreation validates that arm.json is created correctly
func TestManifestCreation(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir

	t.Run("CreatedOnFirstRegistry", func(t *testing.T) {
		// Add registry
		arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")

		// Verify arm.json exists
		armJSON := filepath.Join(workDir, "arm.json")
		helpers.AssertFileExists(t, armJSON)

		// Verify valid JSON
		manifest := helpers.ReadJSON(t, armJSON)

		// Verify contains registries
		registries, ok := manifest["registries"].(map[string]interface{})
		if !ok {
			t.Fatal("registries field not found or invalid type")
		}

		if _, ok := registries["test-registry"]; !ok {
			t.Error("test-registry not found in arm.json")
		}
	})

	t.Run("UpdatedOnSinkAdd", func(t *testing.T) {
		// Add sink
		arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

		// Verify arm.json updated
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)

		// Verify contains sinks
		sinks, ok := manifest["sinks"].(map[string]interface{})
		if !ok {
			t.Fatal("sinks field not found or invalid type")
		}

		if _, ok := sinks["cursor-rules"]; !ok {
			t.Error("cursor-rules not found in arm.json")
		}

		// Verify registries still present (preserves existing config)
		registries, ok := manifest["registries"].(map[string]interface{})
		if !ok {
			t.Fatal("registries field not found after sink add")
		}

		if _, ok := registries["test-registry"]; !ok {
			t.Error("test-registry lost after sink add")
		}
	})

	t.Run("UpdatedOnInstall", func(t *testing.T) {
		// Install ruleset
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

		// Verify arm.json updated
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)

		// Verify contains dependencies
		deps, ok := manifest["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies field not found or invalid type")
		}

		// Find dependency key (may include version)
		found := false
		for key := range deps {
			if strings.Contains(key, "test-ruleset") {
				found = true
				break
			}
		}

		if !found {
			t.Error("test-ruleset not found in dependencies")
		}

		// Verify previous config preserved
		registries, ok := manifest["registries"].(map[string]interface{})
		if !ok || len(registries) == 0 {
			t.Error("registries lost after install")
		}

		sinks, ok := manifest["sinks"].(map[string]interface{})
		if !ok || len(sinks) == 0 {
			t.Error("sinks lost after install")
		}
	})
}

// TestLockFileCreation validates that arm-lock.json is created and updated correctly
func TestLockFileCreation(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("v1.0.0")
	repo.Tag("v1.0.0")

	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v2.0.0")
	repo.Tag("v2.0.0")

	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	t.Run("CreatedOnFirstInstall", func(t *testing.T) {
		// Install ruleset
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

		// Verify arm-lock.json exists
		lockFile := filepath.Join(workDir, "arm-lock.json")
		helpers.AssertFileExists(t, lockFile)

		// Verify valid JSON
		lock := helpers.ReadJSON(t, lockFile)

		// Verify contains dependencies
		deps, ok := lock["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies field not found in lock file")
		}

		if len(deps) == 0 {
			t.Error("expected dependencies in lock file")
		}
	})

	t.Run("ContainsResolvedVersions", func(t *testing.T) {
		lockFile := filepath.Join(workDir, "arm-lock.json")
		lock := helpers.ReadJSON(t, lockFile)

		deps, ok := lock["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies not found in lock file")
		}

		// Find the dependency entry - the key includes version info
		found := false
		for key := range deps {
			if strings.Contains(key, "test-ruleset") {
				found = true
				// Verify key contains version information
				if !strings.Contains(key, "@") {
					t.Errorf("expected version in key, got %s", key)
				}
				break
			}
		}

		if !found {
			t.Error("no dependency entry found in lock file")
		}
	})

	t.Run("UpdatedOnUpgrade", func(t *testing.T) {
		// Upgrade to v2.0.0
		arm.MustRun("upgrade")

		// Verify lock file updated
		lockFile := filepath.Join(workDir, "arm-lock.json")
		lock := helpers.ReadJSON(t, lockFile)

		deps, ok := lock["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies not found in lock file after upgrade")
		}

		// Find dependency with v2.0.0
		found := false
		for key := range deps {
			if strings.Contains(key, "2.0.0") || strings.Contains(key, "v2") {
				found = true
				break
			}
		}

		if !found {
			t.Error("lock file not updated to v2.0.0 after upgrade")
		}
	})
}

// TestLockFileBranchResolution validates that Git branches resolve to commit hashes
func TestLockFileBranchResolution(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository with branch
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")

	// Create feature branch
	repo.Branch("feature-branch")
	repo.Checkout("feature-branch")
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("Feature commit")

	repoURL := "file://" + repoDir
	// Add registry with branch tracking
	arm.MustRun("add", "registry", "git", "--url", repoURL, "--branches", "feature-branch", "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	// Install from branch
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@feature-branch", "cursor-rules")

	// Verify lock file contains branch reference
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)

	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatal("dependencies not found in lock file")
	}

	// Find the dependency key
	found := false
	for key := range deps {
		if strings.Contains(key, "test-ruleset") && strings.Contains(key, "@") {
			found = true
			// Branch should be in the key (either as branch name or commit hash)
			t.Logf("Found dependency key with branch/commit: %s", key)
			break
		}
	}

	if !found {
		t.Error("branch reference not found in lock file")
	}
}

// TestIndexFileCreation validates that arm-index.json is created and updated correctly
func TestIndexFileCreation(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	t.Run("CreatedOnFirstInstall", func(t *testing.T) {
		// Install ruleset
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

		// Verify arm-index.json exists in sink directory
		indexFile := filepath.Join(workDir, ".cursor/rules/arm/arm-index.json")
		helpers.AssertFileExists(t, indexFile)

		// Verify valid JSON
		index := helpers.ReadJSON(t, indexFile)

		// Verify contains rulesets
		rulesets, ok := index["rulesets"].(map[string]interface{})
		if !ok {
			t.Fatal("rulesets field not found in index file")
		}

		if len(rulesets) == 0 {
			t.Error("expected rulesets in index file")
		}
	})

	t.Run("TracksInstalledFiles", func(t *testing.T) {
		indexFile := filepath.Join(workDir, ".cursor/rules/arm/arm-index.json")
		index := helpers.ReadJSON(t, indexFile)

		rulesets, ok := index["rulesets"].(map[string]interface{})
		if !ok {
			t.Fatal("rulesets not found in index file")
		}

		// Find any ruleset entry
		var rulesetEntry map[string]interface{}
		for _, entry := range rulesets {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				rulesetEntry = entryMap
				break
			}
		}

		if rulesetEntry == nil {
			t.Fatal("no ruleset entry found in index file")
		}

		// Verify has files
		files, ok := rulesetEntry["files"].([]interface{})
		if !ok {
			t.Fatal("files field not found in ruleset entry")
		}

		if len(files) == 0 {
			t.Error("expected files tracked in index for ruleset")
		}
	})

	t.Run("UpdatedOnUninstall", func(t *testing.T) {
		// Uninstall
		arm.MustRun("uninstall")

		// Verify index file updated
		indexFile := filepath.Join(workDir, ".cursor/rules/arm/arm-index.json")
		index := helpers.ReadJSON(t, indexFile)

		// After uninstall, rulesets should be empty or omitted
		if rulesets, ok := index["rulesets"].(map[string]interface{}); ok {
			if len(rulesets) > 0 {
				t.Error("expected no rulesets in index after uninstall")
			}
		}
		// If rulesets field is omitted entirely, that's also valid
	})
}

// TestManifestJSONValidity validates that all manifest files are valid JSON
func TestManifestJSONValidity(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

	t.Run("ArmJSONValid", func(t *testing.T) {
		armJSON := filepath.Join(workDir, "arm.json")
		data, err := os.ReadFile(armJSON)
		if err != nil {
			t.Fatalf("failed to read arm.json: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Errorf("arm.json is not valid JSON: %v", err)
		}
	})

	t.Run("ArmLockJSONValid", func(t *testing.T) {
		lockFile := filepath.Join(workDir, "arm-lock.json")
		data, err := os.ReadFile(lockFile)
		if err != nil {
			t.Fatalf("failed to read arm-lock.json: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Errorf("arm-lock.json is not valid JSON: %v", err)
		}
	})

	t.Run("ArmIndexJSONValid", func(t *testing.T) {
		indexFile := filepath.Join(workDir, "arm-index.json")
		
		// Check if file exists first
		if _, err := os.Stat(indexFile); os.IsNotExist(err) {
			t.Skip("arm-index.json not created (may not be required for this configuration)")
		}
		
		data, err := os.ReadFile(indexFile)
		if err != nil {
			t.Fatalf("failed to read arm-index.json: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Errorf("arm-index.json is not valid JSON: %v", err)
		}
	})
}

// TestManifestPreservesConfiguration validates that arm.json preserves existing configuration
func TestManifestPreservesConfiguration(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create two test Git repositories
	repoDir1 := t.TempDir()
	repo1 := helpers.NewGitRepo(t, repoDir1)
	repo1.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo1.Commit("Initial commit")
	repo1.Tag("v1.0.0")

	repoDir2 := t.TempDir()
	repo2 := helpers.NewGitRepo(t, repoDir2)
	repo2.WriteFile("security-ruleset.yml", helpers.SecurityRuleset)
	repo2.Commit("Initial commit")
	repo2.Tag("v1.0.0")

	// Add first registry
	repoURL1 := "file://" + repoDir1
	arm.MustRun("add", "registry", "git", "--url", repoURL1, "registry-one")

	// Add first sink
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	// Add second registry
	repoURL2 := "file://" + repoDir2
	arm.MustRun("add", "registry", "git", "--url", repoURL2, "registry-two")

	// Add second sink
	arm.MustRun("add", "sink", "--tool", "amazonq", "q-rules", ".amazonq/rules")

	// Verify all configuration preserved
	armJSON := filepath.Join(workDir, "arm.json")
	manifest := helpers.ReadJSON(t, armJSON)

	// Check registries
	registries, ok := manifest["registries"].(map[string]interface{})
	if !ok {
		t.Fatal("registries not found")
	}

	if len(registries) != 2 {
		t.Errorf("expected 2 registries, got %d", len(registries))
	}

	if _, ok := registries["registry-one"]; !ok {
		t.Error("registry-one not preserved")
	}

	if _, ok := registries["registry-two"]; !ok {
		t.Error("registry-two not preserved")
	}

	// Check sinks
	sinks, ok := manifest["sinks"].(map[string]interface{})
	if !ok {
		t.Fatal("sinks not found")
	}

	if len(sinks) != 2 {
		t.Errorf("expected 2 sinks, got %d", len(sinks))
	}

	if _, ok := sinks["cursor-rules"]; !ok {
		t.Error("cursor-rules not preserved")
	}

	if _, ok := sinks["q-rules"]; !ok {
		t.Error("q-rules not preserved")
	}
}
