package sink

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name     string
		tool     compiler.Tool
		wantFlat bool
	}{
		{"cursor creates hierarchical", compiler.Cursor, false},
		{"amazonq creates hierarchical", compiler.AmazonQ, false},
		{"copilot creates flat", compiler.Copilot, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			m := NewManager(tmpDir, tt.tool)

			if tt.wantFlat && m.layout != LayoutFlat {
				t.Errorf("expected flat layout, got %v", m.layout)
			}
			if !tt.wantFlat && m.layout != LayoutHierarchical {
				t.Errorf("expected hierarchical layout, got %v", m.layout)
			}

			// Check directory was created
			if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
				t.Errorf("directory not created: %v", err)
			}
		})
	}
}

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		name         string
		layout       Layout
		registry     string
		packageName  string
		version      string
		relativePath string
		wantPattern  string // Pattern to match instead of exact string
	}{
		{
			name:         "flat layout short path",
			layout:       LayoutFlat,
			registry:     "test-reg",
			packageName:  "test-pkg",
			version:      "1.0.0",
			relativePath: "rules/test.yml",
			wantPattern:  "arm_[a-f0-9]{4}_[a-f0-9]{4}_rules_test.yml",
		},
		{
			name:         "flat layout long path gets truncated",
			layout:       LayoutFlat,
			registry:     "test-reg",
			packageName:  "test-pkg",
			version:      "1.0.0",
			relativePath: "very/long/path/that/exceeds/the/maximum/length/limit/for/filename/test.yml",
			wantPattern:  "arm_[a-f0-9]{4}_[a-f0-9]{4}_test.yml", // Should use just filename
		},
		{
			name:         "flat layout very long filename gets truncated",
			layout:       LayoutFlat,
			registry:     "test-reg",
			packageName:  "test-pkg",
			version:      "1.0.0",
			relativePath: "rules/very_long_filename_that_definitely_exceeds_the_maximum_allowed_length_for_a_single_filename_component.instructions.md",
			wantPattern:  "arm_[a-f0-9]{4}_[a-f0-9]{4}_very_long_filename_that_definitely_exceeds_the_maximum_allowed_length_for_a_single_filename_component.instructions.md", // Truncated but keeps extension
		},
		{
			name:         "hierarchical layout",
			layout:       LayoutHierarchical,
			registry:     "test-reg",
			packageName:  "test-pkg",
			version:      "1.0.0",
			relativePath: "rules/test.yml",
			wantPattern:  "arm/test-reg/test-pkg/1.0.0/rules/test.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			m := &Manager{
				directory: tmpDir,
				layout:    tt.layout,
				armDir:    filepath.Join(tmpDir, "arm"),
			}

			got := m.getFilePath(tt.registry, tt.packageName, tt.version, tt.relativePath)

			if tt.layout == LayoutFlat {
				// Check it starts with arm_ and has dual hash pattern
				basename := filepath.Base(got)
				if !filepath.HasPrefix(basename, "arm_") {
					t.Errorf("flat layout should start with arm_, got %v", basename)
				}
				// Check dual hash pattern: arm_xxxx_xxxx_...
				if len(basename) < 14 { // arm_ + 4 + _ + 4 + _ = 14 minimum
					t.Errorf("filename too short for dual hash pattern: %v", basename)
				}
				// Verify total length doesn't exceed 100 chars
				if len(basename) > 100 {
					t.Errorf("filename exceeds 100 chars: %d", len(basename))
				}
			} else {
				expected := filepath.Join(tmpDir, tt.wantPattern)
				if got != expected {
					t.Errorf("getFilePath() = %v, want %v", got, expected)
				}
			}
		})
	}
}

func TestHashPackageAndPath(t *testing.T) {
	m := &Manager{}

	// Test package hash
	pkgHash1 := m.hashPackage("reg", "pkg", "1.0.0")
	pkgHash2 := m.hashPackage("reg", "pkg", "1.0.0")
	pkgHash3 := m.hashPackage("reg", "pkg", "2.0.0")

	// Same inputs should produce same hash
	if pkgHash1 != pkgHash2 {
		t.Errorf("same package inputs should produce same hash: %v != %v", pkgHash1, pkgHash2)
	}

	// Different versions should produce different hash
	if pkgHash1 == pkgHash3 {
		t.Errorf("different package versions should produce different hash: %v == %v", pkgHash1, pkgHash3)
	}

	// Package hash should be 4 characters
	if len(pkgHash1) != 4 {
		t.Errorf("package hash should be 4 characters, got %d", len(pkgHash1))
	}

	// Test path hash
	pathHash1 := m.hashPath("rules/test.yml")
	pathHash2 := m.hashPath("rules/test.yml")
	pathHash3 := m.hashPath("rules/other.yml")

	// Same inputs should produce same hash
	if pathHash1 != pathHash2 {
		t.Errorf("same path inputs should produce same hash: %v != %v", pathHash1, pathHash2)
	}

	// Different paths should produce different hash
	if pathHash1 == pathHash3 {
		t.Errorf("different paths should produce different hash: %v == %v", pathHash1, pathHash3)
	}

	// Path hash should be 4 characters
	if len(pathHash1) != 4 {
		t.Errorf("path hash should be 4 characters, got %d", len(pathHash1))
	}

	// Test that same package produces same hash for different files
	pkgHashA := m.hashPackage("reg", "pkg", "1.0.0")
	pkgHashB := m.hashPackage("reg", "pkg", "1.0.0")
	if pkgHashA != pkgHashB {
		t.Errorf("same package should produce same hash regardless of file: %v != %v", pkgHashA, pkgHashB)
	}
}

func TestIsInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	metadata := core.PackageMetadata{
		RegistryName: "test-reg",
		Name:         "test-pkg",
		Version:      core.Version{Version: "1.0.0"},
	}

	// Should return false for non-installed package
	if m.IsInstalled(metadata) {
		t.Errorf("IsInstalled should return false for non-installed package")
	}
}

func TestLoadSaveIndex(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Load non-existent index should return empty index
	index, err := m.loadIndex()
	if err != nil {
		t.Errorf("loadIndex failed: %v", err)
	}
	if index.Version != 1 {
		t.Errorf("expected version 1, got %d", index.Version)
	}
	if len(index.Rulesets) != 0 {
		t.Errorf("expected empty rulesets, got %d", len(index.Rulesets))
	}

	// Add some data and save
	index.Rulesets["test-reg/test-pkg@1.0.0"] = RulesetIndexEntry{
		Priority: 100,
		Files:    []string{"test.mdc"},
	}

	if err := m.saveIndex(index); err != nil {
		t.Errorf("saveIndex failed: %v", err)
	}

	// Load again and verify
	index2, err := m.loadIndex()
	if err != nil {
		t.Errorf("loadIndex failed: %v", err)
	}
	if len(index2.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(index2.Rulesets))
	}
	entry := index2.Rulesets["test-reg/test-pkg@1.0.0"]
	if entry.Priority != 100 {
		t.Errorf("expected priority 100, got %d", entry.Priority)
	}
}

func TestListRulesets(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Empty list initially
	rulesets, err := m.ListRulesets()
	if err != nil {
		t.Errorf("ListRulesets failed: %v", err)
	}
	if len(rulesets) != 0 {
		t.Errorf("expected empty list, got %d", len(rulesets))
	}

	// Add ruleset to index
	index, _ := m.loadIndex()
	index.Rulesets["test-reg/test-pkg@1.0.0"] = RulesetIndexEntry{
		Priority: 200,
		Files:    []string{"rule1.mdc", "rule2.mdc"},
	}
	m.saveIndex(index)

	// List should return the ruleset
	rulesets, err = m.ListRulesets()
	if err != nil {
		t.Errorf("ListRulesets failed: %v", err)
	}
	if len(rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(rulesets))
	}
	if rulesets[0].Metadata.RegistryName != "test-reg" {
		t.Errorf("expected registry test-reg, got %s", rulesets[0].Metadata.RegistryName)
	}
	if rulesets[0].Priority != 200 {
		t.Errorf("expected priority 200, got %d", rulesets[0].Priority)
	}
}

func TestListPromptsets(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Empty list initially
	promptsets, err := m.ListPromptsets()
	if err != nil {
		t.Errorf("ListPromptsets failed: %v", err)
	}
	if len(promptsets) != 0 {
		t.Errorf("expected empty list, got %d", len(promptsets))
	}

	// Add promptset to index
	index, _ := m.loadIndex()
	index.Promptsets["test-reg/test-prompts@1.0.0"] = PromptsetIndexEntry{
		Files: []string{"prompt1.md", "prompt2.md"},
	}
	m.saveIndex(index)

	// List should return the promptset
	promptsets, err = m.ListPromptsets()
	if err != nil {
		t.Errorf("ListPromptsets failed: %v", err)
	}
	if len(promptsets) != 1 {
		t.Errorf("expected 1 promptset, got %d", len(promptsets))
	}
	if promptsets[0].Metadata.RegistryName != "test-reg" {
		t.Errorf("expected registry test-reg, got %s", promptsets[0].Metadata.RegistryName)
	}
}

func TestUninstall(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Create some test files
	testFile1 := filepath.Join(tmpDir, "arm", "test-reg", "test-pkg", "1.0.0", "rule1.mdc")
	os.MkdirAll(filepath.Dir(testFile1), 0755)
	os.WriteFile(testFile1, []byte("test"), 0644)

	// Add to index
	index, _ := m.loadIndex()
	index.Rulesets["test-reg/test-pkg@1.0.0"] = RulesetIndexEntry{
		Priority: 100,
		Files:    []string{"arm/test-reg/test-pkg/1.0.0/rule1.mdc"},
	}
	m.saveIndex(index)

	// Uninstall
	metadata := core.PackageMetadata{
		RegistryName: "test-reg",
		Name:         "test-pkg",
		Version:      core.Version{Version: "1.0.0"},
	}
	err := m.Uninstall(metadata)
	if err != nil {
		t.Errorf("Uninstall failed: %v", err)
	}

	// File should be removed
	if _, err := os.Stat(testFile1); !os.IsNotExist(err) {
		t.Errorf("file should be removed")
	}

	// Should not be installed anymore
	if m.IsInstalled(metadata) {
		t.Errorf("package should not be installed after uninstall")
	}
}

func TestClean(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Create orphaned file (not in index)
	orphanFile := filepath.Join(tmpDir, "arm", "orphan.mdc")
	os.MkdirAll(filepath.Dir(orphanFile), 0755)
	os.WriteFile(orphanFile, []byte("orphan"), 0644)

	// Create tracked file
	trackedFile := filepath.Join(tmpDir, "arm", "test-reg", "test-pkg", "1.0.0", "rule1.mdc")
	os.MkdirAll(filepath.Dir(trackedFile), 0755)
	os.WriteFile(trackedFile, []byte("tracked"), 0644)

	// Add tracked file to index
	index, _ := m.loadIndex()
	index.Rulesets["test-reg/test-pkg@1.0.0"] = RulesetIndexEntry{
		Priority: 100,
		Files:    []string{"arm/test-reg/test-pkg/1.0.0/rule1.mdc"},
	}
	m.saveIndex(index)

	// Clean should remove orphaned file but keep tracked file
	err := m.Clean()
	if err != nil {
		t.Errorf("Clean failed: %v", err)
	}

	// Orphaned file should be removed
	if _, err := os.Stat(orphanFile); !os.IsNotExist(err) {
		t.Errorf("orphaned file should be removed")
	}

	// Tracked file should remain
	if _, err := os.Stat(trackedFile); os.IsNotExist(err) {
		t.Errorf("tracked file should remain")
	}
}
func TestFilenameTruncation(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Copilot) // Use Copilot for flat layout

	// Get the actual hashes for our test inputs
	pkgHash := m.hashPackage("test-reg", "test-pkg", "1.0.0")

	tests := []struct {
		name         string
		relativePath string
		wantFilename string
	}{
		{
			name:         "short path no truncation",
			relativePath: "rules/test.yml",
			wantFilename: "arm_" + pkgHash + "_" + m.hashPath("rules/test.yml") + "_rules_test.yml",
		},
		{
			name:         "long path uses filename only",
			relativePath: "very/long/directory/structure/that/exceeds/the/maximum/allowed/length/for/filename/generation/test.yml",
			wantFilename: "arm_" + pkgHash + "_" + m.hashPath("very/long/directory/structure/that/exceeds/the/maximum/allowed/length/for/filename/generation/test.yml") + "_test.yml",
		},
		{
			name:         "very long filename gets truncated",
			relativePath: "rules/this_is_a_very_long_filename_that_definitely_exceeds_the_maximum_allowed_length_for_filename_generation_and_should_be_truncated.instructions.md",
			// This will be truncated to fit in 100 chars total
			wantFilename: "arm_" + pkgHash + "_" + m.hashPath("rules/this_is_a_very_long_filename_that_definitely_exceeds_the_maximum_allowed_length_for_filename_generation_and_should_be_truncated.instructions.md") + "_this_is_a_very_long_filename_that_definitely_exceeds_the_maximum_allowed_length_for.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.getFilePath("test-reg", "test-pkg", "1.0.0", tt.relativePath)
			basename := filepath.Base(got)

			// Should not exceed 100 characters
			if len(basename) > 100 {
				t.Errorf("filename exceeds 100 chars: %d", len(basename))
			}

			// Check exact match
			if basename != tt.wantFilename {
				t.Errorf("getFilePath() = %v, want %v", basename, tt.wantFilename)
			}
		})
	}
}

func TestGenerateRulesetIndexRuleFile(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Test with no rulesets - should not create file
	err := m.generateRulesetIndexRuleFile()
	if err != nil {
		t.Errorf("generateRulesetIndexRuleFile failed: %v", err)
	}
	if _, err := os.Stat(m.rulesetIndexRulePath); !os.IsNotExist(err) {
		t.Errorf("index file should not exist when no rulesets")
	}

	// Add rulesets to index
	index, _ := m.loadIndex()
	index.Rulesets["reg1/pkg1@1.0.0"] = RulesetIndexEntry{
		Priority: 200,
		Files:    []string{"arm/reg1/pkg1/1.0.0/rule1.mdc", "arm/reg1/pkg1/1.0.0/rule2.mdc"},
	}
	index.Rulesets["reg2/pkg2@2.0.0"] = RulesetIndexEntry{
		Priority: 100,
		Files:    []string{"arm/reg2/pkg2/2.0.0/rule3.mdc"},
	}
	m.saveIndex(index)

	// Generate index file
	err = m.generateRulesetIndexRuleFile()
	if err != nil {
		t.Errorf("generateRulesetIndexRuleFile failed: %v", err)
	}

	// File should exist
	if _, err := os.Stat(m.rulesetIndexRulePath); os.IsNotExist(err) {
		t.Errorf("index file should exist")
	}

	// Read and verify content
	content, err := os.ReadFile(m.rulesetIndexRulePath)
	if err != nil {
		t.Errorf("failed to read index file: %v", err)
	}

	contentStr := string(content)

	// Check header
	if !strings.Contains(contentStr, "# ARM Rulesets") {
		t.Errorf("missing header")
	}
	if !strings.Contains(contentStr, "Priority Rules") {
		t.Errorf("missing priority rules section")
	}

	// Check both packages are listed
	if !strings.Contains(contentStr, "reg1/pkg1@1.0.0") {
		t.Errorf("missing reg1/pkg1@1.0.0")
	}
	if !strings.Contains(contentStr, "reg2/pkg2@2.0.0") {
		t.Errorf("missing reg2/pkg2@2.0.0")
	}

	// Check priorities
	if !strings.Contains(contentStr, "**Priority:** 200") {
		t.Errorf("missing priority 200")
	}
	if !strings.Contains(contentStr, "**Priority:** 100") {
		t.Errorf("missing priority 100")
	}

	// Check files are listed
	if !strings.Contains(contentStr, "arm/reg1/pkg1/1.0.0/rule1.mdc") {
		t.Errorf("missing rule1.mdc")
	}
	if !strings.Contains(contentStr, "arm/reg2/pkg2/2.0.0/rule3.mdc") {
		t.Errorf("missing rule3.mdc")
	}

	// Verify priority order (higher priority first)
	pkg1Idx := strings.Index(contentStr, "reg1/pkg1@1.0.0")
	pkg2Idx := strings.Index(contentStr, "reg2/pkg2@2.0.0")
	if pkg1Idx > pkg2Idx {
		t.Errorf("higher priority package should appear first")
	}
}

func TestGenerateRulesetIndexRuleFileRemovesWhenEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Create index file first
	index, _ := m.loadIndex()
	index.Rulesets["reg1/pkg1@1.0.0"] = RulesetIndexEntry{
		Priority: 100,
		Files:    []string{"rule1.mdc"},
	}
	m.saveIndex(index)
	m.generateRulesetIndexRuleFile()

	// Verify file exists
	if _, err := os.Stat(m.rulesetIndexRulePath); os.IsNotExist(err) {
		t.Errorf("index file should exist")
	}

	// Remove all rulesets
	index.Rulesets = make(map[string]RulesetIndexEntry)
	m.saveIndex(index)
	m.generateRulesetIndexRuleFile()

	// File should be removed
	if _, err := os.Stat(m.rulesetIndexRulePath); !os.IsNotExist(err) {
		t.Errorf("index file should be removed when no rulesets")
	}
}

func TestInstallRuleset(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Create test package with regular file
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "test-pkg",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{
			{
				Path:    "README.md",
				Content: []byte("# Test Package"),
			},
		},
	}

	err := m.InstallRuleset(pkg, 100)
	if err != nil {
		t.Errorf("InstallRuleset failed: %v", err)
	}

	// Check package is installed
	if !m.IsInstalled(pkg.Metadata) {
		t.Errorf("package should be installed")
	}

	// Check index
	index, _ := m.loadIndex()
	key := pkgKey("test-reg", "test-pkg", "1.0.0")
	entry, exists := index.Rulesets[key]
	if !exists {
		t.Errorf("package should be in index")
	}
	if entry.Priority != 100 {
		t.Errorf("expected priority 100, got %d", entry.Priority)
	}
	if len(entry.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(entry.Files))
	}

	// Check file exists
	expectedPath := filepath.Join(tmpDir, "arm", "test-reg", "test-pkg", "1.0.0", "README.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("file should exist at %s", expectedPath)
	}

	// Check index rule file exists
	if _, err := os.Stat(m.rulesetIndexRulePath); os.IsNotExist(err) {
		t.Errorf("index rule file should exist")
	}
}

func TestInstallRulesetReplacesOldVersion(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Install v1.0.0
	pkg1 := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "test-pkg",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{
			{
				Path:    "old.md",
				Content: []byte("old"),
			},
		},
	}
	m.InstallRuleset(pkg1, 100)

	// Install v1.0.0 again with different file
	pkg2 := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "test-pkg",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{
			{
				Path:    "new.md",
				Content: []byte("new"),
			},
		},
	}
	err := m.InstallRuleset(pkg2, 200)
	if err != nil {
		t.Errorf("InstallRuleset failed: %v", err)
	}

	// Check only new file exists
	oldPath := filepath.Join(tmpDir, "arm", "test-reg", "test-pkg", "1.0.0", "old.md")
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Errorf("old file should be removed")
	}

	newPath := filepath.Join(tmpDir, "arm", "test-reg", "test-pkg", "1.0.0", "new.md")
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("new file should exist")
	}

	// Check priority updated
	index, _ := m.loadIndex()
	key := pkgKey("test-reg", "test-pkg", "1.0.0")
	if index.Rulesets[key].Priority != 200 {
		t.Errorf("expected priority 200, got %d", index.Rulesets[packageID].Priority)
	}
}

func TestInstallRulesetEmptyPackage(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Install package with no files
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "empty-pkg",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{},
	}

	err := m.InstallRuleset(pkg, 100)
	if err != nil {
		t.Errorf("InstallRuleset should handle empty package: %v", err)
	}

	// Should still be in index
	if !m.IsInstalled(pkg.Metadata) {
		t.Errorf("empty package should be installed")
	}

	index, _ := m.loadIndex()
	key := pkgKey("test-reg", "empty-pkg", "1.0.0")
	if len(index.Rulesets[key].Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(index.Rulesets[packageID].Files))
	}
}

func TestInstallPromptset(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Create test package with regular file
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "test-prompts",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{
			{
				Path:    "README.md",
				Content: []byte("# Test Promptset"),
			},
		},
	}

	err := m.InstallPromptset(pkg)
	if err != nil {
		t.Errorf("InstallPromptset failed: %v", err)
	}

	// Check package is installed
	if !m.IsInstalled(pkg.Metadata) {
		t.Errorf("package should be installed")
	}

	// Check index
	index, _ := m.loadIndex()
	key := pkgKey("test-reg", "test-prompts", "1.0.0")
	entry, exists := index.Promptsets[key]
	if !exists {
		t.Errorf("package should be in index")
	}
	if len(entry.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(entry.Files))
	}

	// Check file exists
	expectedPath := filepath.Join(tmpDir, "arm", "test-reg", "test-prompts", "1.0.0", "README.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("file should exist at %s", expectedPath)
	}
}

func TestInstallPromptsetReplacesOldVersion(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Install v1.0.0
	pkg1 := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "test-prompts",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{
			{
				Path:    "old.md",
				Content: []byte("old"),
			},
		},
	}
	m.InstallPromptset(pkg1)

	// Install v1.0.0 again with different file
	pkg2 := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "test-prompts",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{
			{
				Path:    "new.md",
				Content: []byte("new"),
			},
		},
	}
	err := m.InstallPromptset(pkg2)
	if err != nil {
		t.Errorf("InstallPromptset failed: %v", err)
	}

	// Check only new file exists
	oldPath := filepath.Join(tmpDir, "arm", "test-reg", "test-prompts", "1.0.0", "old.md")
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Errorf("old file should be removed")
	}

	newPath := filepath.Join(tmpDir, "arm", "test-reg", "test-prompts", "1.0.0", "new.md")
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("new file should exist")
	}
}

func TestInstallPromptsetEmptyPackage(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	// Install package with no files
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-reg",
			Name:         "empty-prompts",
			Version:      core.Version{Version: "1.0.0"},
		},
		Files: []*core.File{},
	}

	err := m.InstallPromptset(pkg)
	if err != nil {
		t.Errorf("InstallPromptset should handle empty package: %v", err)
	}

	// Should still be in index
	if !m.IsInstalled(pkg.Metadata) {
		t.Errorf("empty package should be installed")
	}

	index, _ := m.loadIndex()
	key := pkgKey("test-reg", "empty-prompts", "1.0.0")
	if len(index.Promptsets[key].Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(index.Promptsets[packageID].Files))
	}
}