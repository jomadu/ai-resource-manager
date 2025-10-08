package installer

import (
	"context"
	"os"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestHierarchicalInstaller_ListInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	installer := NewHierarchicalInstaller(tmpDir, "markdown")
	ctx := context.Background()

	// Install test files
	files := []types.File{
		{Path: "rule1.md", Content: []byte("rule1")},
		{Path: "subdir/rule2.md", Content: []byte("rule2")},
	}

	err := installer.InstallRuleset(ctx, "registry1", "ruleset1", "1.0.0", 100, files)
	if err != nil {
		t.Fatalf("InstallRuleset failed: %v", err)
	}

	// Test ListInstalledRulesets
	installations, err := installer.ListInstalledRulesets(ctx)
	if err != nil {
		t.Fatalf("ListInstalledRulesets failed: %v", err)
	}

	if len(installations) != 1 {
		t.Fatalf("Expected 1 installation, got %d", len(installations))
	}

	installation := installations[0]
	if installation.Ruleset != "ruleset1" {
		t.Errorf("Expected ruleset 'ruleset1', got '%s'", installation.Ruleset)
	}
	if installation.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", installation.Version)
	}
	if len(installation.FilePaths) != 2 {
		t.Errorf("Expected 2 file paths, got %d", len(installation.FilePaths))
	}

	// Verify file paths exist
	for _, filePath := range installation.FilePaths {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("File path does not exist: %s", filePath)
		}
	}
}
