package e2e

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

// TestUpdateWithinConstraints verifies that update respects version constraints
func TestUpdateWithinConstraints(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple versions
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create v1.0.0
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("v1.0.0")
	repo.Tag("v1.0.0")
	
	// Create v1.1.0
	updatedRuleset := helpers.MinimalRuleset + `
    ruleThree:
      body: "This is rule three added in v1.1.0."
`
	repo.WriteFile("test-ruleset.yml", updatedRuleset)
	repo.Commit("v1.1.0")
	repo.Tag("v1.1.0")
	
	// Create v2.0.0
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v2.0.0")
	repo.Tag("v2.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Install v1.0.0 with major constraint (@1)
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1", "cursor-rules")
	
	// Verify initial installation is v1.1.0 (highest 1.x.x)
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	deps := lock["dependencies"].(map[string]interface{})
	
	var installedVersion string
	for key := range deps {
		if strings.Contains(key, "test-ruleset@") {
			parts := strings.Split(key, "@")
			if len(parts) >= 2 {
				installedVersion = parts[len(parts)-1]
			}
			break
		}
	}
	
	if !strings.HasPrefix(installedVersion, "v1.1") && !strings.HasPrefix(installedVersion, "1.1") {
		t.Errorf("expected v1.1.x to be installed, got %s", installedVersion)
	}
	
	// Run update - should stay at v1.1.0 (not upgrade to v2.0.0)
	arm.MustRun("update")
	
	// Verify still on v1.x.x
	lock = helpers.ReadJSON(t, lockFile)
	deps = lock["dependencies"].(map[string]interface{})
	
	for key := range deps {
		if strings.Contains(key, "test-ruleset@") {
			parts := strings.Split(key, "@")
			if len(parts) >= 2 {
				installedVersion = parts[len(parts)-1]
			}
			break
		}
	}
	
	if strings.HasPrefix(installedVersion, "v2") || strings.HasPrefix(installedVersion, "2") {
		t.Errorf("update should not upgrade to v2.x.x, got %s", installedVersion)
	}
}

// TestUpgradeIgnoresConstraints verifies that upgrade ignores version constraints
func TestUpgradeIgnoresConstraints(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple versions
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create v1.0.0
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("v1.0.0")
	repo.Tag("v1.0.0")
	
	// Create v2.0.0
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v2.0.0")
	repo.Tag("v2.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Install v1.0.0 with exact constraint
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")
	
	// Verify initial installation is v1.0.0
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	deps := lock["dependencies"].(map[string]interface{})
	
	var installedVersion string
	for key := range deps {
		if strings.Contains(key, "test-ruleset@") {
			parts := strings.Split(key, "@")
			if len(parts) >= 2 {
				installedVersion = parts[len(parts)-1]
			}
			break
		}
	}
	
	if !strings.Contains(installedVersion, "1.0.0") {
		t.Errorf("expected v1.0.0 to be installed, got %s", installedVersion)
	}
	
	// Run upgrade - should upgrade to v2.0.0 ignoring constraint
	arm.MustRun("upgrade")
	
	// Verify upgraded to v2.0.0
	lock = helpers.ReadJSON(t, lockFile)
	deps = lock["dependencies"].(map[string]interface{})
	
	for key := range deps {
		if strings.Contains(key, "test-ruleset@") {
			parts := strings.Split(key, "@")
			if len(parts) >= 2 {
				installedVersion = parts[len(parts)-1]
			}
			break
		}
	}
	
	if !strings.Contains(installedVersion, "2.0.0") {
		t.Errorf("upgrade should upgrade to v2.0.0, got %s", installedVersion)
	}
	
	// Verify arm.json constraint updated
	armJSON := filepath.Join(workDir, "arm.json")
	manifest := helpers.ReadJSON(t, armJSON)
	deps = manifest["dependencies"].(map[string]interface{})
	
	var constraint string
	for key, val := range deps {
		if strings.Contains(key, "test-ruleset") {
			depInfo := val.(map[string]interface{})
			constraint = depInfo["version"].(string)
			break
		}
	}
	
	// Constraint should be updated to ^2.0.0 (major constraint for v2.x.x)
	if !strings.HasPrefix(constraint, "^2") {
		t.Errorf("expected constraint to be ^2.0.0, got %s", constraint)
	}
}

// TestUpdatePromptset verifies update works with promptsets
func TestUpdatePromptset(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple versions
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create v1.0.0
	repo.WriteFile("test-promptset.yml", helpers.MinimalPromptset)
	repo.Commit("v1.0.0")
	repo.Tag("v1.0.0")
	
	// Create v1.1.0
	updatedPromptset := helpers.MinimalPromptset + `
    promptThree:
      body: "This is prompt three added in v1.1.0."
`
	repo.WriteFile("test-promptset.yml", updatedPromptset)
	repo.Commit("v1.1.0")
	repo.Tag("v1.1.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-prompts", ".cursor/prompts")
	
	// Install v1.0.0
	arm.MustRun("install", "promptset", "test-registry/test-promptset@1.0.0", "cursor-prompts")
	
	// Run update - should stay at v1.0.0 (exact constraint)
	arm.MustRun("update")
	
	// Verify still on v1.0.0
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	deps := lock["dependencies"].(map[string]interface{})
	
	var installedVersion string
	for key := range deps {
		if strings.Contains(key, "test-promptset@") {
			parts := strings.Split(key, "@")
			if len(parts) >= 2 {
				installedVersion = parts[len(parts)-1]
			}
			break
		}
	}
	
	if !strings.Contains(installedVersion, "1.0.0") {
		t.Errorf("update with exact constraint should stay at v1.0.0, got %s", installedVersion)
	}
}

// TestUpgradePromptset verifies upgrade works with promptsets
func TestUpgradePromptset(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple versions
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create v1.0.0
	repo.WriteFile("test-promptset.yml", helpers.MinimalPromptset)
	repo.Commit("v1.0.0")
	repo.Tag("v1.0.0")
	
	// Create v2.0.0
	repo.WriteFile("test-promptset.yml", helpers.CodeReviewPromptset)
	repo.Commit("v2.0.0")
	repo.Tag("v2.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-prompts", ".cursor/prompts")
	
	// Install v1.0.0
	arm.MustRun("install", "promptset", "test-registry/test-promptset@1.0.0", "cursor-prompts")
	
	// Run upgrade - should upgrade to v2.0.0
	arm.MustRun("upgrade")
	
	// Verify upgraded to v2.0.0
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	deps := lock["dependencies"].(map[string]interface{})
	
	var installedVersion string
	for key := range deps {
		if strings.Contains(key, "test-promptset@") {
			parts := strings.Split(key, "@")
			if len(parts) >= 2 {
				installedVersion = parts[len(parts)-1]
			}
			break
		}
	}
	
	if !strings.Contains(installedVersion, "2.0.0") {
		t.Errorf("upgrade should upgrade to v2.0.0, got %s", installedVersion)
	}
}

// TestManifestFilesUpdated verifies arm.json and arm-lock.json are updated correctly
func TestManifestFilesUpdated(t *testing.T) {
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
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Install v1.0.0
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")
	
	// Verify arm.json has v1.0.0 constraint
	armJSON := filepath.Join(workDir, "arm.json")
	manifest := helpers.ReadJSON(t, armJSON)
	deps := manifest["dependencies"].(map[string]interface{})
	
	var depKey string
	for key := range deps {
		if strings.Contains(key, "test-ruleset") {
			depKey = key
			break
		}
	}
	
	if depKey == "" {
		t.Fatal("dependency not found in arm.json")
	}
	
	depInfo := deps[depKey].(map[string]interface{})
	version := depInfo["version"].(string)
	
	if !strings.Contains(version, "1.0.0") {
		t.Errorf("expected version constraint to contain 1.0.0, got %s", version)
	}
	
	// Verify arm-lock.json has v1.0.0
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	lockDeps := lock["dependencies"].(map[string]interface{})
	
	var lockKey string
	for key := range lockDeps {
		if strings.Contains(key, "test-ruleset@") {
			lockKey = key
			break
		}
	}
	
	if lockKey == "" {
		t.Fatal("dependency not found in arm-lock.json")
	}
	
	if !strings.Contains(lockKey, "1.0.0") {
		t.Errorf("expected lock file to contain v1.0.0, got %s", lockKey)
	}
	
	// Run upgrade
	arm.MustRun("upgrade")
	
	// Verify arm-lock.json updated to v2.0.0
	lock = helpers.ReadJSON(t, lockFile)
	lockDeps = lock["dependencies"].(map[string]interface{})
	
	var newLockKey string
	for key := range lockDeps {
		if strings.Contains(key, "test-ruleset@") {
			newLockKey = key
			break
		}
	}
	
	if !strings.Contains(newLockKey, "2.0.0") {
		t.Errorf("expected lock file to be updated to v2.0.0, got %s", newLockKey)
	}
}
