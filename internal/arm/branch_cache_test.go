package arm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
)

func TestBranchConstraintCaching(t *testing.T) {
	// Skip if GitHub CLI not available or not authenticated
	if err := exec.Command("gh", "auth", "status").Run(); err != nil {
		t.Skip("GitHub CLI not authenticated - skipping integration test")
	}

	// Nuke cache at start of test to ensure clean state
	nukeCache(t)

	ctx := context.Background()
	tempDir := t.TempDir()
	repoName := fmt.Sprintf("arm-test-%d", time.Now().Unix())

	// Ensure cleanup happens even if test panics
	var repoCreated bool
	defer func() {
		if repoCreated {
			deleteTestRepo(t, repoName)
		}
	}()

	// Create temporary GitHub repository
	repoURL := createTestRepo(t, repoName, tempDir)
	repoCreated = true

	// Setup ARM in temp directory
	workDir := filepath.Join(tempDir, "work")
	setupARM(t, workDir, repoName, repoURL)

	service := NewArmService()
	// Ensure we're in the work directory for all operations
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(workDir)

	// First install from main branch
	err := service.InstallRuleset(ctx, "test-registry", "rules", "main", []string{"*.md"}, nil)
	if err != nil {
		t.Fatalf("First install failed: %v", err)
	}

	// Verify first content
	firstContent := readInstalledContent(t, workDir)
	if !strings.Contains(firstContent, "Initial content") {
		t.Fatalf("First install didn't get initial content")
	}

	// Push change to main branch
	pushChangeToRepo(t, repoName, "Updated content for cache test")

	// Ensure we're still in work directory
	os.Chdir(workDir)

	// Second install from main branch (should get new content but currently gets cached)
	err = service.InstallRuleset(ctx, "test-registry", "rules", "main", []string{"*.md"}, nil)
	if err != nil {
		t.Fatalf("Second install failed: %v", err)
	}

	// Verify second content - this will fail with current implementation
	secondContent := readInstalledContent(t, workDir)
	if !strings.Contains(secondContent, "Updated content") {
		t.Errorf("BUG REPRODUCED: Second install still has old cached content")
		preview := secondContent
		if len(preview) > 100 {
			preview = preview[:100]
		}
		t.Logf("Expected: 'Updated content', Got content containing: %s", preview)
	}
}

func createTestRepo(t *testing.T, repoName, tempDir string) string {
	repoDir := filepath.Join(tempDir, "repo")
	os.MkdirAll(repoDir, 0755)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(repoDir)

	// Initialize git repo
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	exec.Command("git", "config", "user.name", "ARM Test").Run()
	exec.Command("git", "config", "user.email", "test@arm.com").Run()

	// Create initial content
	os.WriteFile("test-rule.md", []byte("# Test Rule\n\nInitial content for testing.\n"), 0644)
	exec.Command("git", "add", ".").Run()
	if err := exec.Command("git", "commit", "-m", "Initial commit").Run(); err != nil {
		t.Fatalf("Failed to commit initial content: %v", err)
	}

	// Create GitHub repo and push
	cmd := exec.Command("gh", "repo", "create", repoName, "--public", "--source=.", "--remote=origin", "--push")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create GitHub repo: %v", err)
	}

	// Verify repo was created and get URL
	out, err := exec.Command("gh", "repo", "view", repoName, "--json", "url", "-q", ".url").Output()
	if err != nil {
		t.Fatalf("Failed to get repo URL: %v", err)
	}
	repoURL := strings.TrimSpace(string(out))
	if repoURL == "" {
		t.Fatalf("Got empty repo URL")
	}
	t.Logf("Created test repository: %s", repoURL)
	return repoURL
}

func deleteTestRepo(t *testing.T, repoName string) {
	// Try multiple times to ensure cleanup
	for i := 0; i < 3; i++ {
		cmd := exec.Command("gh", "repo", "delete", repoName, "--confirm")
		if err := cmd.Run(); err == nil {
			t.Logf("Successfully deleted test repository: %s", repoName)
			return
		}
		// Wait before retry
		time.Sleep(2 * time.Second)
	}
	t.Logf("Warning: Failed to delete test repository: %s (may need manual cleanup)", repoName)
}

func setupARM(t *testing.T, workDir, repoName, repoURL string) {
	os.MkdirAll(workDir, 0755)
	os.Chdir(workDir)

	// Create ARM config
	configManager := config.NewFileManager()
	err := configManager.AddSink(context.Background(), "test", []string{".rules"}, []string{"test-registry/*"}, nil, "hierarchical", true)
	if err != nil {
		t.Fatalf("Failed to add sink: %v", err)
	}

	// Create manifest with test registry
	manifestManager := manifest.NewFileManager()
	registryConfig := registry.GitRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			URL:  repoURL,
			Type: "git",
		},
		Branches: []string{"main"},
	}
	err = manifestManager.AddGitRegistry(context.Background(), "test-registry", registryConfig, true)
	if err != nil {
		t.Fatalf("Failed to add registry: %v", err)
	}
}

func readInstalledContent(t *testing.T, workDir string) string {
	// Find installed rule file
	ruleFile := filepath.Join(workDir, ".rules", "arm", "test-registry", "rules", "*", "test-rule.md")
	matches, _ := filepath.Glob(ruleFile)
	if len(matches) == 0 {
		t.Fatalf("No installed rule file found at %s", ruleFile)
	}

	content, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("Failed to read installed content: %v", err)
	}
	return string(content)
}

func pushChangeToRepo(t *testing.T, repoName, newContent string) {
	tempDir := t.TempDir()
	cloneDir := filepath.Join(tempDir, "clone")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Clone repo
	cmd := exec.Command("gh", "repo", "clone", repoName, cloneDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to clone repo: %v", err)
	}

	os.Chdir(cloneDir)

	// Update content
	os.WriteFile("test-rule.md", []byte(fmt.Sprintf("# Test Rule\n\n%s\n", newContent)), 0644)
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		t.Fatalf("Failed to add changes: %v", err)
	}
	if err := exec.Command("git", "commit", "-m", "Update content").Run(); err != nil {
		t.Fatalf("Failed to commit changes: %v", err)
	}
	if err := exec.Command("git", "push").Run(); err != nil {
		t.Fatalf("Failed to push changes: %v", err)
	}
	t.Logf("Successfully pushed changes to repository")
}

func nukeCache(t *testing.T) {
	cacheDir := cache.GetCacheDir()
	if err := os.RemoveAll(cacheDir); err != nil {
		t.Logf("Warning: Failed to remove cache directory %s: %v", cacheDir, err)
	} else {
		t.Logf("Nuked cache directory: %s", cacheDir)
	}
}
