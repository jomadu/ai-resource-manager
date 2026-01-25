package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestCompilationToolFormats(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository with ruleset
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")

	t.Run("CursorFormat", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "cursor-rules")

		// Verify .mdc file with frontmatter in hierarchical structure
		sinkDir := filepath.Join(workDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset")
		helpers.AssertDirExists(t, sinkDir)

		// Find .mdc files recursively
		foundMDC := false
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".mdc") {
				foundMDC = true
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				// Verify frontmatter exists
				if !strings.Contains(string(content), "---") {
					t.Error("expected frontmatter in .mdc file")
				}
			}
			return nil
		})
		if !foundMDC {
			t.Error("expected .mdc file in cursor sink")
		}
	})

	t.Run("AmazonQFormat", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "amazonq", "q-rules", ".amazonq/rules")
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "q-rules")

		// Verify .md file (pure markdown) in hierarchical structure
		sinkDir := filepath.Join(workDir, ".amazonq", "rules", "arm", "test-registry", "test-ruleset")
		helpers.AssertDirExists(t, sinkDir)

		// Find .md files recursively
		foundMD := false
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && !strings.Contains(info.Name(), "arm_index") {
				foundMD = true
			}
			return nil
		})
		if !foundMD {
			t.Error("expected .md file in amazonq sink")
		}
	})

	t.Run("CopilotFormat", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "copilot", "copilot-rules", ".github/instructions")
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "copilot-rules")

		// Copilot uses flat layout - files are in root with hash prefixes
		sinkDir := filepath.Join(workDir, ".github", "instructions")
		helpers.AssertDirExists(t, sinkDir)

		// Find .instructions.md files in root directory
		files, err := os.ReadDir(sinkDir)
		if err != nil {
			t.Fatalf("failed to read sink directory: %v", err)
		}

		foundInstructions := false
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".instructions.md") && !strings.Contains(file.Name(), "arm_index") {
				foundInstructions = true
				break
			}
		}
		if !foundInstructions {
			t.Error("expected .instructions.md file in copilot sink")
		}
	})

	t.Run("MarkdownFormat", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "markdown", "md-rules", ".rules")
		arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "md-rules")

		// Verify .md file in hierarchical structure
		sinkDir := filepath.Join(workDir, ".rules", "arm", "test-registry", "test-ruleset")
		helpers.AssertDirExists(t, sinkDir)

		// Find .md files recursively
		foundMD := false
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && !strings.Contains(info.Name(), "arm_index") {
				foundMD = true
			}
			return nil
		})
		if !foundMD {
			t.Error("expected .md file in markdown sink")
		}
	})
}

func TestCompilationPromptsets(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository with promptset
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-promptset.yml", helpers.MinimalPromptset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")

	t.Run("CursorPrompts", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "cursor", "cursor-prompts", ".cursor/prompts")
		arm.MustRun("install", "promptset", "test-registry/test-promptset@1.0.0", "cursor-prompts")

		// Verify .md file (no frontmatter for prompts) in hierarchical structure
		sinkDir := filepath.Join(workDir, ".cursor", "prompts", "arm", "test-registry", "test-promptset")
		helpers.AssertDirExists(t, sinkDir)

		// Find .md files recursively
		foundMD := false
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				foundMD = true
			}
			return nil
		})
		if !foundMD {
			t.Error("expected .md file in cursor prompts sink")
		}
	})

	t.Run("AmazonQPrompts", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "amazonq", "q-prompts", ".amazonq/prompts")
		arm.MustRun("install", "promptset", "test-registry/test-promptset@1.0.0", "q-prompts")

		// Verify .md file in hierarchical structure
		sinkDir := filepath.Join(workDir, ".amazonq", "prompts", "arm", "test-registry", "test-promptset")
		helpers.AssertDirExists(t, sinkDir)

		// Find .md files recursively
		foundMD := false
		_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				foundMD = true
			}
			return nil
		})
		if !foundMD {
			t.Error("expected .md file in amazonq prompts sink")
		}
	})
}

func TestCompilationValidation(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)

	t.Run("ValidRuleset", func(t *testing.T) {
		repo.WriteFile("valid-ruleset.yml", helpers.MinimalRuleset)
		repo.Commit("Add valid ruleset")
		repo.Tag("v1.0.0")

		repoURL := "file://" + repoDir
		arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
		arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

		// Should succeed
		arm.MustRun("install", "ruleset", "test-registry/valid-ruleset@1.0.0", "cursor-rules")

		// Verify files were created
		sinkDir := filepath.Join(workDir, ".cursor", "rules", "arm", "test-registry", "valid-ruleset")
		helpers.AssertDirExists(t, sinkDir)
	})
}

func TestCompilationIndexGeneration(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository with ruleset
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

	t.Run("IndexFileGenerated", func(t *testing.T) {
		// Verify arm_index.mdc exists in the arm/ subdirectory
		indexFile := filepath.Join(workDir, ".cursor", "rules", "arm", "arm_index.mdc")
		helpers.AssertFileExists(t, indexFile)

		// Verify index contains metadata
		content, err := os.ReadFile(indexFile)
		if err != nil {
			t.Fatalf("failed to read index file: %v", err)
		}

		// Should contain priority information or metadata
		contentStr := string(content)
		if !strings.Contains(contentStr, "priority") && !strings.Contains(contentStr, "Priority") && !strings.Contains(contentStr, "test-ruleset") {
			t.Error("expected metadata in index file")
		}
	})

	t.Run("IndexJSONGenerated", func(t *testing.T) {
		// Verify arm-index.json exists in the arm/ subdirectory
		indexJSON := filepath.Join(workDir, ".cursor", "rules", "arm", "arm-index.json")
		helpers.AssertFileExists(t, indexJSON)

		// Verify it's valid JSON
		data := helpers.ReadJSON(t, indexJSON)
		if data == nil {
			t.Error("expected valid JSON in arm-index.json")
		}
	})
}

func TestCompilationHierarchicalLayout(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository with ruleset
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

	t.Run("HierarchicalStructure", func(t *testing.T) {
		// Verify hierarchical directory structure exists
		// arm/<registry>/<package>/<version>/
		// Note: version includes 'v' prefix
		expectedPath := filepath.Join(workDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0")
		helpers.AssertDirExists(t, expectedPath)

		// Verify files are in the hierarchical path
		files, err := os.ReadDir(expectedPath)
		if err != nil {
			t.Fatalf("failed to read hierarchical directory: %v", err)
		}

		if len(files) == 0 {
			t.Error("expected files in hierarchical directory")
		}
	})
}

func TestCompilationMultiplePriorities(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository with two rulesets
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("low-priority.yml", helpers.MinimalRuleset)
	repo.WriteFile("high-priority.yml", `
kind: Ruleset
metadata:
  name: high-priority
  version: 1.0.0
spec:
  rules:
    - id: high-rule
      title: High Priority Rule
      description: This rule has high priority
      content: High priority content
`)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	// Install with different priorities
	arm.MustRun("install", "ruleset", "--priority", "50", "test-registry/low-priority@1.0.0", "cursor-rules")
	arm.MustRun("install", "ruleset", "--priority", "200", "test-registry/high-priority@1.0.0", "cursor-rules")

	t.Run("IndexOrderedByPriority", func(t *testing.T) {
		// Verify arm_index.mdc exists in arm/ subdirectory
		indexFile := filepath.Join(workDir, ".cursor", "rules", "arm", "arm_index.mdc")
		content, err := os.ReadFile(indexFile)
		if err != nil {
			t.Fatalf("failed to read index file: %v", err)
		}

		// Higher priority (200) should appear before lower priority (50)
		contentStr := string(content)
		highPos := strings.Index(contentStr, "high-priority")
		lowPos := strings.Index(contentStr, "low-priority")

		if highPos == -1 || lowPos == -1 {
			t.Error("expected both rulesets in index")
		} else if highPos > lowPos {
			t.Error("expected high priority ruleset to appear before low priority")
		}
	})
}
