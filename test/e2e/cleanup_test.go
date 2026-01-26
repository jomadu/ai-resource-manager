package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestUninstallCleanup(t *testing.T) {
	t.Run("removes empty directories after uninstall", func(t *testing.T) {
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

		// Install package
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

		sinkDir := filepath.Join(workDir, ".cursor", "rules")

		// Verify files and directories exist
		indexPath := filepath.Join(sinkDir, "arm", "arm-index.json")
		helpers.AssertFileExists(t, indexPath)

		// Find at least one directory was created
		var foundDir bool
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() && path != sinkDir {
				foundDir = true
			}
			return nil
		})
		if !foundDir {
			t.Fatalf("expected directories to be created after install")
		}

		// Uninstall package
		arm.MustRun("uninstall", "test-registry/test-ruleset")

		// Verify empty directories are removed
		var emptyDirs []string
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() && path != sinkDir {
				entries, _ := os.ReadDir(path)
				if len(entries) == 0 {
					emptyDirs = append(emptyDirs, path)
				}
			}
			return nil
		})
		if len(emptyDirs) > 0 {
			t.Errorf("found empty directories after uninstall: %v", emptyDirs)
		}

		// Verify arm-index.json is removed
		if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
			t.Errorf("arm-index.json should be removed after uninstalling all packages")
		}

		// Verify sink root directory still exists
		if _, err := os.Stat(sinkDir); os.IsNotExist(err) {
			t.Errorf("sink root directory should never be removed")
		}
	})

	t.Run("removes arm-index.json when all packages uninstalled", func(t *testing.T) {
		workDir := t.TempDir()
		arm := helpers.NewARMRunner(t, workDir)

		// Setup: Create test Git repository with two packages
		repoDir := t.TempDir()
		repo := helpers.NewGitRepo(t, repoDir)
		repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
		repo.WriteFile("test-promptset.yml", helpers.MinimalPromptset)
		repo.Commit("Initial commit")
		repo.Tag("v1.0.0")

		// Setup: Add registry and sink
		repoURL := "file://" + repoDir
		arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
		arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

		// Install two packages
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")
		arm.MustRun("install", "promptset", "test-registry/test-promptset@1.0.0", "cursor-rules")

		sinkDir := filepath.Join(workDir, ".cursor", "rules")
		indexPath := filepath.Join(sinkDir, "arm", "arm-index.json")

		// Verify index exists
		helpers.AssertFileExists(t, indexPath)

		// Uninstall first package
		arm.MustRun("uninstall", "test-registry/test-ruleset")

		// Index should still exist (one package remaining)
		helpers.AssertFileExists(t, indexPath)

		// Uninstall second package
		arm.MustRun("uninstall", "test-registry/test-promptset")

		// Index should be removed (no packages remaining)
		if _, err := os.Stat(indexPath); !os.IsNotExist(err) {
			t.Errorf("arm-index.json should be removed when all packages uninstalled")
		}
	})

	t.Run("removes priority index when all rulesets uninstalled", func(t *testing.T) {
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

		// Install ruleset
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

		sinkDir := filepath.Join(workDir, ".cursor", "rules")

		// Find priority index file
		var priorityIndexPath string
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Base(path) == "arm_index.mdc" {
				priorityIndexPath = path
			}
			return nil
		})
		if priorityIndexPath == "" {
			t.Fatalf("priority index file should exist after installing ruleset")
		}

		// Uninstall ruleset
		arm.MustRun("uninstall", "test-registry/test-ruleset")

		// Priority index should be removed
		if _, err := os.Stat(priorityIndexPath); !os.IsNotExist(err) {
			t.Errorf("priority index file should be removed when all rulesets uninstalled")
		}
	})

	t.Run("handles multiple packages in same sink", func(t *testing.T) {
		workDir := t.TempDir()
		arm := helpers.NewARMRunner(t, workDir)

		// Setup: Create test Git repository with two packages
		repoDir := t.TempDir()
		repo := helpers.NewGitRepo(t, repoDir)
		repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
		repo.WriteFile("test-promptset.yml", helpers.MinimalPromptset)
		repo.Commit("Initial commit")
		repo.Tag("v1.0.0")

		// Setup: Add registry and sink
		repoURL := "file://" + repoDir
		arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
		arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

		// Install two packages
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")
		arm.MustRun("install", "promptset", "test-registry/test-promptset@1.0.0", "cursor-rules")

		sinkDir := filepath.Join(workDir, ".cursor", "rules")
		indexPath := filepath.Join(sinkDir, "arm", "arm-index.json")

		// Uninstall first package
		arm.MustRun("uninstall", "test-registry/test-ruleset")

		// Index should still exist
		helpers.AssertFileExists(t, indexPath)

		// Verify second package still installed
		output := arm.MustRun("info", "dependency")
		if !strings.Contains(output, "test-promptset") {
			t.Errorf("second package should still be installed")
		}
	})
}
