package sink

import (
	"os"
	"path/filepath"
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
		want         string
	}{
		{
			name:         "flat layout",
			layout:       LayoutFlat,
			registry:     "test-reg",
			packageName:  "test-pkg",
			version:      "1.0.0",
			relativePath: "rules/test.yml",
			want:         "arm_12345678_rules_test.yml", // hash will be computed
		},
		{
			name:         "hierarchical layout",
			layout:       LayoutHierarchical,
			registry:     "test-reg",
			packageName:  "test-pkg",
			version:      "1.0.0",
			relativePath: "rules/test.yml",
			want:         "arm/test-reg/test-pkg/1.0.0/rules/test.yml",
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
				// Check it starts with arm_ and has hash
				if !filepath.HasPrefix(filepath.Base(got), "arm_") {
					t.Errorf("flat layout should start with arm_, got %v", got)
				}
			} else {
				expected := filepath.Join(tmpDir, tt.want)
				if got != expected {
					t.Errorf("getFilePath() = %v, want %v", got, expected)
				}
			}
		})
	}
}

func TestHashFile(t *testing.T) {
	m := &Manager{}

	hash1 := m.hashFile("reg", "pkg", "1.0.0", "test.yml")
	hash2 := m.hashFile("reg", "pkg", "1.0.0", "test.yml")
	hash3 := m.hashFile("reg", "pkg", "2.0.0", "test.yml")

	// Same inputs should produce same hash
	if hash1 != hash2 {
		t.Errorf("same inputs should produce same hash: %v != %v", hash1, hash2)
	}

	// Different inputs should produce different hash
	if hash1 == hash3 {
		t.Errorf("different inputs should produce different hash: %v == %v", hash1, hash3)
	}

	// Hash should be 8 characters
	if len(hash1) != 8 {
		t.Errorf("hash should be 8 characters, got %d", len(hash1))
	}
}

func TestIsInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	m := NewManager(tmpDir, compiler.Cursor)

	metadata := core.PackageMetadata{
		Registry: "test-reg",
		Name:     "test-pkg",
		Version:  "1.0.0",
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
	if rulesets[0].Metadata.Registry != "test-reg" {
		t.Errorf("expected registry test-reg, got %s", rulesets[0].Metadata.Registry)
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
	if promptsets[0].Metadata.Registry != "test-reg" {
		t.Errorf("expected registry test-reg, got %s", promptsets[0].Metadata.Registry)
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
		Registry: "test-reg",
		Name:     "test-pkg",
		Version:  "1.0.0",
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