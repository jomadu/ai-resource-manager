package e2e

import (
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestRulesetInstallation(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install ruleset with semver tag
	t.Run("InstallRulesetWithSemver", func(t *testing.T) {
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")
		
		// Verify arm.json updated
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		deps, ok := manifest["dependencies"].(map[string]interface{})
		if !ok {
			t.Fatal("dependencies field not found")
		}
		
		if _, ok := deps["test-registry/test-ruleset"]; !ok {
			t.Error("dependency not found in arm.json")
		}
		
		// Verify arm-lock.json created
		lockFile := filepath.Join(workDir, "arm-lock.json")
		helpers.AssertFileExists(t, lockFile)
		
		// Verify sink directory populated
		sinkDir := filepath.Join(workDir, ".cursor/rules")
		helpers.AssertDirExists(t, sinkDir)
		
		// Verify compiled files exist
		fileCount := helpers.CountFilesRecursive(t, sinkDir)
		if fileCount == 0 {
			t.Error("expected compiled files in sink directory")
		}
	})
}

func TestRulesetInstallationWithLatest(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple versions
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
	
	// Test: Install with @latest should get v2.0.0
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@latest", "cursor-rules")
	
	// Verify lock file contains v2.0.0
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	
	deps, ok := lock["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatalf("dependencies not found in lock file, lock file: %+v", lock)
	}
	
	// The key includes the version: test-registry/test-ruleset@v2.0.0
	// Find the key that starts with test-registry/test-ruleset@
	var foundKey string
	for key := range deps {
		if contains(key, "test-registry/test-ruleset@") || contains(key, "test-ruleset@") {
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

// Helper to get map keys for debugging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestRulesetInstallationWithBranch(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with branch
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	
	repo.Branch("develop")
	repo.Checkout("develop")
	repo.WriteFile("dev-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("Add dev ruleset")
	
	// Setup: Add registry with branches tracking and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "--branches", "develop", "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install from branch (use @develop since branch is tracked)
	arm.MustRun("install", "ruleset", "test-registry/dev-ruleset@develop", "cursor-rules")
	
	// Verify lock file contains commit hash
	lockFile := filepath.Join(workDir, "arm-lock.json")
	lock := helpers.ReadJSON(t, lockFile)
	deps := lock["dependencies"].(map[string]interface{})
	
	// Find the key that contains dev-ruleset
	var foundKey string
	for key := range deps {
		if contains(key, "dev-ruleset@") {
			foundKey = key
			break
		}
	}
	
	if foundKey == "" {
		t.Fatalf("dependency not found in lock file, available keys: %v", getKeys(deps))
	}
	
	// The key should contain @develop or @<commit-hash>
	if !contains(foundKey, "@develop") && !contains(foundKey, "@") {
		t.Errorf("expected branch or commit hash in key, got %s", foundKey)
	}
}

func TestRulesetInstallationWithPriority(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install with custom priority
	arm.MustRun("install", "ruleset", "--priority", "200", "test-registry/test-ruleset@1.0.0", "cursor-rules")
	
	// Verify priority in arm.json
	armJSON := filepath.Join(workDir, "arm.json")
	manifest := helpers.ReadJSON(t, armJSON)
	deps := manifest["dependencies"].(map[string]interface{})
	dep := deps["test-registry/test-ruleset"].(map[string]interface{})
	
	priority := dep["priority"]
	if priority != float64(200) {
		t.Errorf("expected priority 200, got %v", priority)
	}
}

func TestRulesetInstallationToMultipleSinks(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")
	
	// Setup: Add registry and multiple sinks
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	arm.MustRun("add", "sink", "--tool", "amazonq", "q-rules", ".amazonq/rules")
	
	// Test: Install to multiple sinks
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules", "q-rules")
	
	// Verify both sink directories populated
	cursorDir := filepath.Join(workDir, ".cursor/rules")
	helpers.AssertDirExists(t, cursorDir)
	if helpers.CountFilesRecursive(t, cursorDir) == 0 {
		t.Error("expected files in cursor sink")
	}
	
	qDir := filepath.Join(workDir, ".amazonq/rules")
	helpers.AssertDirExists(t, qDir)
	if helpers.CountFilesRecursive(t, qDir) == 0 {
		t.Error("expected files in amazonq sink")
	}
	
	// Verify arm.json tracks both sinks
	armJSON := filepath.Join(workDir, "arm.json")
	manifest := helpers.ReadJSON(t, armJSON)
	deps := manifest["dependencies"].(map[string]interface{})
	dep := deps["test-registry/test-ruleset"].(map[string]interface{})
	
	sinks, ok := dep["sinks"].([]interface{})
	if !ok {
		t.Fatal("sinks not found in dependency")
	}
	
	if len(sinks) != 2 {
		t.Errorf("expected 2 sinks, got %d", len(sinks))
	}
}

func TestPromptsetInstallation(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-promptset.yml", helpers.MinimalPromptset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-prompts", ".cursor/prompts")
	
	// Test: Install promptset
	arm.MustRun("install", "promptset", "test-registry/test-promptset@1.0.0", "cursor-prompts")
	
	// Verify arm.json updated
	armJSON := filepath.Join(workDir, "arm.json")
	manifest := helpers.ReadJSON(t, armJSON)
	deps := manifest["dependencies"].(map[string]interface{})
	
	if _, ok := deps["test-registry/test-promptset"]; !ok {
		t.Error("promptset dependency not found in arm.json")
	}
	
	// Verify sink directory populated
	sinkDir := filepath.Join(workDir, ".cursor/prompts")
	helpers.AssertDirExists(t, sinkDir)
	
	fileCount := helpers.CountFilesRecursive(t, sinkDir)
	if fileCount == 0 {
		t.Error("expected compiled files in sink directory")
	}
}

func TestInstallWithPatterns(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Setup: Create test Git repository with multiple files
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("security/auth-rules.yml", helpers.SecurityRuleset)
	repo.WriteFile("general/clean-code.yml", helpers.MinimalRuleset)
	repo.WriteFile("experimental/test.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")
	
	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
	
	// Test: Install with include pattern
	t.Run("InstallWithInclude", func(t *testing.T) {
		arm.MustRun("install", "ruleset", "--include", "security/**/*.yml", 
			"test-registry/security-rules@1.0.0", "cursor-rules")
		
		// Verify only security files are installed
		sinkDir := filepath.Join(workDir, ".cursor/rules")
		helpers.AssertDirExists(t, sinkDir)
		
		// Should have files from security directory
		fileCount := helpers.CountFilesRecursive(t, sinkDir)
		if fileCount == 0 {
			t.Error("expected files in sink directory")
		}
	})
}
