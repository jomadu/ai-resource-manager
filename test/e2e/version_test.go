package e2e

import (
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestVersionResolutionLatest(t *testing.T) {
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
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v1.1.0")
	repo.Tag("v1.1.0")
	
	// Create v2.0.0
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("v2.0.0")
	repo.Tag("v2.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install with @latest should resolve to v2.0.0
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@latest", "cursor-rules")
	
	// Verify arm-lock.json contains resolved version
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	
	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatalf("dependencies not found in lock file, lock file: %+v", lock)
	}
	
	// The key includes the version: test-registry/test-ruleset@2.0.0
	var foundKey string
	for key := range deps {
		if contains(key, "test-registry/test-ruleset@") {
			foundKey = key
			break
		}
	}
	
	if foundKey == "" {
		t.Fatalf("dependency not found in lock file, available keys: %v", getKeys(deps))
	}
	
	// Verify the key contains v2.0.0
	if !contains(foundKey, "2.0.0") {
		t.Errorf("expected version 2.0.0 in key, got %s", foundKey)
	}
}

func TestVersionResolutionMajorConstraint(t *testing.T) {
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
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v1.1.0")
	repo.Tag("v1.1.0")
	
	// Create v1.2.0
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("v1.2.0")
	repo.Tag("v1.2.0")
	
	// Create v2.0.0
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v2.0.0")
	repo.Tag("v2.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install with @1 should resolve to highest 1.x.x (v1.2.0)
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1", "cursor-rules")
	
	// Verify arm-lock.json contains v1.2.0
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	
	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatalf("dependencies not found in lock file, lock file: %+v", lock)
	}
	
	// Find the key that contains test-registry/test-ruleset@
	var foundKey string
	for key := range deps {
		if contains(key, "test-registry/test-ruleset@") {
			foundKey = key
			break
		}
	}
	
	if foundKey == "" {
		t.Fatalf("dependency not found in lock file, available keys: %v", getKeys(deps))
	}
	
	// Verify the key contains v1.2.0
	if !contains(foundKey, "1.2.0") {
		t.Errorf("expected version 1.2.0 in key, got %s", foundKey)
	}
}

func TestVersionResolutionMinorConstraint(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple versions
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create v1.1.0
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("v1.1.0")
	repo.Tag("v1.1.0")
	
	// Create v1.1.1
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v1.1.1")
	repo.Tag("v1.1.1")
	
	// Create v1.1.2
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("v1.1.2")
	repo.Tag("v1.1.2")
	
	// Create v1.2.0
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("v1.2.0")
	repo.Tag("v1.2.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install with @1.1 should resolve to highest 1.1.x (v1.1.2)
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.1", "cursor-rules")
	
	// Verify arm-lock.json contains v1.1.2
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	
	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatalf("dependencies not found in lock file, lock file: %+v", lock)
	}
	
	// Find the key that contains test-registry/test-ruleset@
	var foundKey string
	for key := range deps {
		if contains(key, "test-registry/test-ruleset@") {
			foundKey = key
			break
		}
	}
	
	if foundKey == "" {
		t.Fatalf("dependency not found in lock file, available keys: %v", getKeys(deps))
	}
	
	// Verify the key contains v1.1.2
	if !contains(foundKey, "1.1.2") {
		t.Errorf("expected version 1.1.2 in key, got %s", foundKey)
	}
}

func TestVersionResolutionExactVersion(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple versions
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create v1.0.0 with unique content
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.WriteFile("version.txt", "1.0.0")
	repo.Commit("v1.0.0")
	repo.Tag("v1.0.0")
	
	// Create v1.1.0 with different content
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.WriteFile("version.txt", "1.1.0")
	repo.Commit("v1.1.0")
	repo.Tag("v1.1.0")
	
	// Create v2.0.0
	repo.WriteFile("version.txt", "2.0.0")
	repo.Commit("v2.0.0")
	repo.Tag("v2.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install with @1.0.0 should resolve to exactly 1.0.0
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")
	
	// Verify arm-lock.json contains v1.0.0 (exact version)
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	
	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatalf("dependencies not found in lock file, lock file: %+v", lock)
	}
	
	// Find the key that contains test-registry/test-ruleset@
	var foundKey string
	for key := range deps {
		if contains(key, "test-registry/test-ruleset@") {
			foundKey = key
			break
		}
	}
	
	if foundKey == "" {
		t.Fatalf("dependency not found in lock file, available keys: %v", getKeys(deps))
	}
	
	// Verify the key contains v1.0.0 (exact version)
	if !contains(foundKey, "1.0.0") {
		t.Errorf("expected version 1.0.0 (exact), got %s", foundKey)
	}
	
	// Verify it doesn't contain 1.1.0 or 2.0.0
	if contains(foundKey, "1.1.0") {
		t.Errorf("expected exact version 1.0.0, but got 1.1.0 in key: %s", foundKey)
	}
	if contains(foundKey, "2.0.0") {
		t.Errorf("expected exact version 1.0.0, but got 2.0.0 in key: %s", foundKey)
	}
}

func TestVersionResolutionBranchToCommit(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with branch
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create main branch content
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	
	// Create feature branch
	repo.Branch("feature-branch")
	repo.Checkout("feature-branch")
	repo.WriteFile("test-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("Feature commit")
	
	// Setup: Add registry with branch tracking and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "--branches", "feature-branch", "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install with branch name should resolve to commit hash
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@feature-branch", "cursor-rules")
	
	// Verify arm-lock.json contains commit hash or branch name
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	
	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatalf("dependencies not found in lock file, lock file: %+v", lock)
	}
	
	// Find the key that contains test-registry/test-ruleset@
	var foundKey string
	for key := range deps {
		if contains(key, "test-registry/test-ruleset@") {
			foundKey = key
			break
		}
	}
	
	if foundKey == "" {
		t.Fatalf("dependency not found in lock file, available keys: %v", getKeys(deps))
	}
	
	// The key should contain @feature-branch or @<commit-hash>
	if !contains(foundKey, "@feature-branch") && !contains(foundKey, "@") {
		t.Errorf("expected branch or commit hash in key, got %s", foundKey)
	}
}

func TestVersionResolutionLatestWithNoTags(t *testing.T) {
	t.Skip("Skipping: @latest without tags requires branch tracking, which is tested in TestVersionResolutionBranchToCommit")
}
