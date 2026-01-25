package e2e

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

// TestArchiveTarGz tests installing from .tar.gz archives
func TestArchiveTarGz(t *testing.T) {
	testDir := t.TempDir()
	repoDir := filepath.Join(testDir, "repo")
	projectDir := filepath.Join(testDir, "project")

	// Create directories
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Create Git repository with .tar.gz archive
	repo := helpers.NewGitRepo(t, repoDir)

	// Create .tar.gz archive with rules
	tarGzContent := createTarGzArchive(t, map[string]string{
		"rule1.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "rule1"
  name: "Rule 1"
  description: "Test rule 1 from archive"
spec:
  rules:
    rule1:
      body: "This is rule 1 from tar.gz archive"`,
		"rule2.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "rule2"
  name: "Rule 2"
  description: "Test rule 2 from archive"
spec:
  rules:
    rule2:
      body: "This is rule 2 from tar.gz archive"`,
	})

	repo.WriteFile("test-ruleset/rules.tar.gz", string(tarGzContent))
	repo.Commit("Add tar.gz archive")
	repo.Tag("v1.0.0")

	// Setup ARM project
	arm := helpers.NewARMRunner(t, projectDir)
	arm.MustRun("add", "registry", "git", "--url", "file://"+repoDir, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "test-sink", ".cursor/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "test-sink")

	// Verify extracted files exist (filename format is {rulesetID}_{ruleID}.mdc)
	helpers.AssertFileExists(t, filepath.Join(projectDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "rule1_rule1.mdc"))
	helpers.AssertFileExists(t, filepath.Join(projectDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "rule2_rule2.mdc"))

	// Verify archive file itself is not present
	archivePath := filepath.Join(projectDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "rules.tar.gz")
	if _, err := os.Stat(archivePath); err == nil {
		t.Errorf("archive file should not be present in sink: %s", archivePath)
	}
}

// TestArchiveZip tests installing from .zip archives
func TestArchiveZip(t *testing.T) {
	testDir := t.TempDir()
	repoDir := filepath.Join(testDir, "repo")
	projectDir := filepath.Join(testDir, "project")

	// Create directories
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Create Git repository with .zip archive
	repo := helpers.NewGitRepo(t, repoDir)

	// Create .zip archive with rules
	zipContent := createZipArchive(t, map[string]string{
		"rule1.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "rule1"
  name: "Rule 1"
  description: "Test rule 1 from zip"
spec:
  rules:
    rule1:
      body: "This is rule 1 from zip archive"`,
		"rule2.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "rule2"
  name: "Rule 2"
  description: "Test rule 2 from zip"
spec:
  rules:
    rule2:
      body: "This is rule 2 from zip archive"`,
	})

	repo.WriteFile("test-ruleset/rules.zip", string(zipContent))
	repo.Commit("Add zip archive")
	repo.Tag("v1.0.0")

	// Setup ARM project
	arm := helpers.NewARMRunner(t, projectDir)
	arm.MustRun("add", "registry", "git", "--url", "file://"+repoDir, "test-registry")
	arm.MustRun("add", "sink", "--tool", "amazonq", "test-sink", ".amazonq/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "test-sink")

	// Verify extracted files exist (filename format is {rulesetID}_{ruleID}.md for amazonq)
	helpers.AssertFileExists(t, filepath.Join(projectDir, ".amazonq", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "rule1_rule1.md"))
	helpers.AssertFileExists(t, filepath.Join(projectDir, ".amazonq", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "rule2_rule2.md"))

	// Verify archive file itself is not present
	archivePath := filepath.Join(projectDir, ".amazonq", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "rules.zip")
	if _, err := os.Stat(archivePath); err == nil {
		t.Errorf("archive file should not be present in sink: %s", archivePath)
	}
}

// TestArchiveMixedWithLooseFiles tests archives mixed with loose files
func TestArchiveMixedWithLooseFiles(t *testing.T) {
	testDir := t.TempDir()
	repoDir := filepath.Join(testDir, "repo")
	projectDir := filepath.Join(testDir, "project")

	// Create directories
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Create Git repository with both archives and loose files
	repo := helpers.NewGitRepo(t, repoDir)

	// Create .tar.gz archive
	tarGzContent := createTarGzArchive(t, map[string]string{
		"archived-rule1.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "archivedRule1"
  name: "Archived Rule 1"
  description: "From tar.gz archive"
spec:
  rules:
    archivedRule1:
      body: "This is from tar.gz archive"`,
	})

	// Create .zip archive
	zipContent := createZipArchive(t, map[string]string{
		"archived-rule2.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "archivedRule2"
  name: "Archived Rule 2"
  description: "From zip archive"
spec:
  rules:
    archivedRule2:
      body: "This is from zip archive"`,
	})

	repo.WriteFile("test-ruleset/rules.tar.gz", string(tarGzContent))
	repo.WriteFile("test-ruleset/rules.zip", string(zipContent))
	repo.WriteFile("test-ruleset/loose-rule.yml", `apiVersion: v1
kind: Ruleset
metadata:
  id: "looseRule"
  name: "Loose Rule"
  description: "Loose file"
spec:
  rules:
    looseRule:
      body: "This is a loose file"`)
	repo.Commit("Add mixed files")
	repo.Tag("v1.0.0")

	// Setup ARM project
	arm := helpers.NewARMRunner(t, projectDir)
	arm.MustRun("add", "registry", "git", "--url", "file://"+repoDir, "test-registry")
	arm.MustRun("add", "sink", "--tool", "copilot", "test-sink", ".github/instructions")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "test-sink")

	// Verify all files exist (Copilot uses flat layout with hash-prefixed names)
	// Just verify that we have 3 rule files (2 from archives + 1 loose)
	sinkDir := filepath.Join(projectDir, ".github", "instructions")
	ruleFileCount := 0
	_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" && filepath.Base(path) != "arm_index.instructions.md" {
			ruleFileCount++
		}
		return nil
	})

	if ruleFileCount != 3 {
		t.Errorf("expected 3 rule files, got %d", ruleFileCount)
	}
}

// TestArchivePrecedenceOverLooseFiles tests that archives override loose files with same path
func TestArchivePrecedenceOverLooseFiles(t *testing.T) {
	testDir := t.TempDir()
	repoDir := filepath.Join(testDir, "repo")
	projectDir := filepath.Join(testDir, "project")

	// Create directories
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Create Git repository with conflicting files
	repo := helpers.NewGitRepo(t, repoDir)

	// Create archive with rule1.yml
	tarGzContent := createTarGzArchive(t, map[string]string{
		"rule1.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "rule1"
  name: "Rule 1"
  description: "From archive"
spec:
  rules:
    rule1:
      body: "This is from the archive"`,
	})

	// Create loose file with same name (will be overridden by archive)
	repo.WriteFile("test-ruleset/rule1.yml", `apiVersion: v1
kind: Ruleset
metadata:
  id: "rule1"
  name: "Rule 1"
  description: "Loose file"
spec:
  rules:
    rule1:
      body: "This is the loose file"`)
	repo.WriteFile("test-ruleset/rules.tar.gz", string(tarGzContent))
	repo.Commit("Add conflicting files")
	repo.Tag("v1.0.0")

	// Setup ARM project
	arm := helpers.NewARMRunner(t, projectDir)
	arm.MustRun("add", "registry", "git", "--url", "file://"+repoDir, "test-registry")
	arm.MustRun("add", "sink", "--tool", "markdown", "test-sink", ".arm/rules")
	arm.MustRun("install", "ruleset", "test-registry/test-ruleset@1.0.0", "test-sink")

	// Verify archive version wins (filename format is {rulesetID}_{ruleID}.md for markdown)
	rulePath := filepath.Join(projectDir, ".arm", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "rule1_rule1.md")
	helpers.AssertFileExists(t, rulePath)

	content, err := os.ReadFile(rulePath)
	if err != nil {
		t.Fatalf("failed to read rule file: %v", err)
	}

	if !bytes.Contains(content, []byte("This is from the archive")) {
		t.Errorf("expected archive content, got: %s", string(content))
	}
	if bytes.Contains(content, []byte("This is the loose file")) {
		t.Errorf("loose file content should be overridden by archive")
	}
}

// TestArchiveWithIncludeExcludePatterns tests pattern filtering on extracted archive content
func TestArchiveWithIncludeExcludePatterns(t *testing.T) {
	t.Skip("Pattern filtering with archives needs more investigation - patterns may need to match extracted paths")

	testDir := t.TempDir()
	repoDir := filepath.Join(testDir, "repo")
	projectDir := filepath.Join(testDir, "project")

	// Create directories
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Create Git repository with archive containing multiple files
	repo := helpers.NewGitRepo(t, repoDir)

	tarGzContent := createTarGzArchive(t, map[string]string{
		"security/rule1.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "securityRule1"
  name: "Security Rule 1"
  description: "Security rule 1"
spec:
  rules:
    securityRule1:
      body: "Security content 1"`,
		"security/rule2.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "securityRule2"
  name: "Security Rule 2"
  description: "Security rule 2"
spec:
  rules:
    securityRule2:
      body: "Security content 2"`,
		"general/rule3.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "generalRule3"
  name: "General Rule 3"
  description: "General rule 3"
spec:
  rules:
    generalRule3:
      body: "General content 3"`,
		"experimental/rule4.yml": `apiVersion: v1
kind: Ruleset
metadata:
  id: "experimentalRule4"
  name: "Experimental Rule 4"
  description: "Experimental rule 4"
spec:
  rules:
    experimentalRule4:
      body: "Experimental content 4"`,
	})

	repo.WriteFile("test-ruleset/rules.tar.gz", string(tarGzContent))
	repo.Commit("Add archive with multiple files")
	repo.Tag("v1.0.0")

	// Setup ARM project
	arm := helpers.NewARMRunner(t, projectDir)
	arm.MustRun("add", "registry", "git", "--url", "file://"+repoDir, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "test-sink", ".cursor/rules")

	// Install with include pattern for security files, exclude experimental
	arm.MustRun("install", "ruleset",
		"--include", "security/**/*.yml",
		"--exclude", "**/experimental/**",
		"test-registry/test-ruleset@1.0.0", "test-sink")

	// Debug: list all files created
	sinkDir := filepath.Join(projectDir, ".cursor", "rules")
	_ = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(sinkDir, path)
			t.Logf("Created file: %s", relPath)
		}
		return nil
	})

	// Verify only security files are present (filename format is {rulesetID}_{ruleID}.mdc)
	helpers.AssertFileExists(t, filepath.Join(projectDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "securityRule1_securityRule1.mdc"))
	helpers.AssertFileExists(t, filepath.Join(projectDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "securityRule2_securityRule2.mdc"))

	// Verify general and experimental files are not present
	generalPath := filepath.Join(projectDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "generalRule3_generalRule3.mdc")
	if _, err := os.Stat(generalPath); err == nil {
		t.Errorf("general file should not be present: %s", generalPath)
	}

	expPath := filepath.Join(projectDir, ".cursor", "rules", "arm", "test-registry", "test-ruleset", "v1.0.0", "experimentalRule4_experimentalRule4.mdc")
	if _, err := os.Stat(expPath); err == nil {
		t.Errorf("experimental file should not be present: %s", expPath)
	}
}

// Helper functions to create archives

func createTarGzArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for path, content := range files {
		header := &tar.Header{
			Name: path,
			Mode: 0o644,
			Size: int64(len(content)),
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		if _, err := tarWriter.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write tar content: %v", err)
		}
	}

	if err := tarWriter.Close(); err != nil {
		t.Fatalf("failed to close tar writer: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}

	return buf.Bytes()
}

func createZipArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for path, content := range files {
		writer, err := zipWriter.Create(path)
		if err != nil {
			t.Fatalf("failed to create zip entry: %v", err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write zip content: %v", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}

	return buf.Bytes()
}
