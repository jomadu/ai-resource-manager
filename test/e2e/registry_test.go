package e2e

import (
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestGitRegistryManagement(t *testing.T) {
	// Create isolated test environment
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Create a test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Add test resources
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")
	
	// Test: Add Git registry with file:// URL
	t.Run("AddGitRegistry", func(t *testing.T) {
		repoURL := "file://" + repoDir
		arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
		
		// Verify arm.json was created
		armJSON := filepath.Join(workDir, "arm.json")
		helpers.AssertFileExists(t, armJSON)
		
		// Verify registry is in the manifest
		manifest := helpers.ReadJSON(t, armJSON)
		registries, ok := manifest["registries"].(map[string]interface{})
		if !ok {
			t.Fatal("registries field not found or invalid type")
		}
		
		if _, ok := registries["test-registry"]; !ok {
			t.Error("test-registry not found in registries")
		}
	})
	
	// Test: List registries
	t.Run("ListRegistries", func(t *testing.T) {
		output := arm.MustRun("list", "registry")
		
		// Should contain the registry name
		if !contains(output, "test-registry") {
			t.Errorf("list registry output should contain 'test-registry', got: %s", output)
		}
	})
	
	// Test: Info for specific registry
	t.Run("InfoRegistry", func(t *testing.T) {
		output := arm.MustRun("info", "registry", "test-registry")
		
		// Should contain registry details
		if !contains(output, "test-registry") {
			t.Errorf("info registry output should contain 'test-registry', got: %s", output)
		}
		if !contains(output, repoDir) {
			t.Errorf("info registry output should contain repo path, got: %s", output)
		}
	})
	
	// Test: Set registry configuration
	t.Run("SetRegistry", func(t *testing.T) {
		// Create another repo for testing set
		newRepoDir := t.TempDir()
		newRepo := helpers.NewGitRepo(t, newRepoDir)
		newRepo.WriteFile("test.yml", helpers.MinimalRuleset)
		newRepo.Commit("Initial commit")
		
		newRepoURL := "file://" + newRepoDir
		arm.MustRun("set", "registry", "test-registry", "url", newRepoURL)
		
		// Verify the URL was updated
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries := manifest["registries"].(map[string]interface{})
		registry := registries["test-registry"].(map[string]interface{})
		
		if registry["url"] != newRepoURL {
			t.Errorf("registry URL not updated: expected %s, got %s", newRepoURL, registry["url"])
		}
	})
	
	// Test: Remove registry
	t.Run("RemoveRegistry", func(t *testing.T) {
		arm.MustRun("remove", "registry", "test-registry")
		
		// Verify registry is removed from manifest
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries, ok := manifest["registries"].(map[string]interface{})
		
		// Registries might be nil or empty after removing the last one
		if ok && registries != nil {
			if _, ok := registries["test-registry"]; ok {
				t.Error("test-registry should be removed from registries")
			}
		}
	})
	
	// Test: Add duplicate registry should fail without --force
	t.Run("AddDuplicateRegistryFails", func(t *testing.T) {
		repoURL := "file://" + repoDir
		arm.MustRun("add", "registry", "git", "--url", repoURL, "dup-registry")
		
		// Try to add again without --force
		stderr := arm.MustFail("add", "registry", "git", "--url", repoURL, "dup-registry")
		
		if !contains(stderr, "already exists") && !contains(stderr, "duplicate") {
			t.Errorf("expected error about duplicate registry, got: %s", stderr)
		}
	})
}

func TestGitRegistryWithBranches(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)
	
	// Create a test Git repository with branches
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	
	// Create main branch content
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")
	
	// Create develop branch
	repo.Branch("develop")
	repo.Checkout("develop")
	repo.WriteFile("dev-ruleset.yml", helpers.SecurityRuleset)
	repo.Commit("Add dev ruleset")
	
	// Switch back to main
	repo.Checkout("main")
	
	// Add registry with branches specification
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "--branches", "develop", "test-registry")
	
	// Verify branches is stored in manifest
	armJSON := filepath.Join(workDir, "arm.json")
	manifest := helpers.ReadJSON(t, armJSON)
	registries := manifest["registries"].(map[string]interface{})
	registry := registries["test-registry"].(map[string]interface{})
	
	// Branches should be stored as an array
	branches, ok := registry["branches"]
	if !ok {
		t.Error("branches field not found in registry")
		return
	}
	
	// Check if branches contains develop
	branchesArray, ok := branches.([]interface{})
	if !ok {
		t.Errorf("branches should be an array, got %T", branches)
		return
	}
	
	found := false
	for _, b := range branchesArray {
		if b == "develop" {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("expected branch 'develop' in branches, got %v", branches)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
