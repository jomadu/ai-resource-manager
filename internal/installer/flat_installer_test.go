package installer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestFlatInstaller_HashFile(t *testing.T) {
	installer := NewFlatInstaller()

	tests := []struct {
		registry, ruleset, version, filePath string
	}{
		{
			registry: "ai-rules",
			ruleset:  "cursor-rules",
			version:  "1.0.0",
			filePath: "rules/cursor/grug-brained-dev.mdc",
		},
		{
			registry: "my_registry",
			ruleset:  "test_ruleset_with_underscores",
			version:  "v1.0.0-beta.2",
			filePath: "src/main.js",
		},
	}

	for _, test := range tests {
		hash := installer.hashFile(test.registry, test.ruleset, test.version, test.filePath)
		if len(hash) != 8 { // Truncated SHA256 hex string length
			t.Errorf("hashFile should return 8-char hex string, got %d chars", len(hash))
		}

		// Test consistency - same input should produce same hash
		hash2 := installer.hashFile(test.registry, test.ruleset, test.version, test.filePath)
		if hash != hash2 {
			t.Errorf("hashFile should be deterministic, got different hashes")
		}
	}
}

func TestFlatInstaller_IndexOperations(t *testing.T) {
	installer := NewFlatInstaller()

	tmpDir, err := os.MkdirTemp("", "index_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Test empty index
	index, err := installer.loadIndex(tmpDir)
	if err != nil {
		t.Fatalf("loadIndex failed: %v", err)
	}
	if len(index) != 0 {
		t.Errorf("Expected empty index, got %d entries", len(index))
	}

	// Add entry and save
	index["test.txt"] = IndexEntry{
		Registry: "test-registry",
		Ruleset:  "test-ruleset",
		Version:  "1.0.0",
		FilePath: "test.txt",
	}

	if err := installer.saveIndex(tmpDir, index); err != nil {
		t.Fatalf("saveIndex failed: %v", err)
	}

	// Load and verify
	loadedIndex, err := installer.loadIndex(tmpDir)
	if err != nil {
		t.Fatalf("loadIndex after save failed: %v", err)
	}
	if len(loadedIndex) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(loadedIndex))
	}
	if entry, exists := loadedIndex["test.txt"]; !exists || entry.Registry != "test-registry" {
		t.Errorf("Index entry not preserved correctly")
	}
}

func TestFlatInstaller_InstallUninstall(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "flat_installer_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	installer := NewFlatInstaller()
	ctx := context.Background()

	files := []types.File{
		{
			Path:    "rules/cursor/test.mdc",
			Content: []byte("test content"),
		},
	}

	// Test Install
	err = installer.Install(ctx, tmpDir, "ai-rules", "cursor-rules", "1.0.0", files)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Check file exists with hashed name
	hash := installer.hashFile("ai-rules", "cursor-rules", "1.0.0", "rules/cursor/test.mdc")
	expectedFile := filepath.Join(tmpDir, hash+"_rules_cursor_test.mdc")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s does not exist", expectedFile)
	}

	// Test ListInstalled
	installations, err := installer.ListInstalled(ctx, tmpDir)
	if err != nil {
		t.Fatalf("ListInstalled failed: %v", err)
	}
	if len(installations) != 1 {
		t.Errorf("Expected 1 installation, got %d", len(installations))
	}
	if installations[0].Ruleset != "cursor-rules" || installations[0].Version != "1.0.0" {
		t.Errorf("Unexpected installation: %+v", installations[0])
	}

	// Test Uninstall
	err = installer.Uninstall(ctx, tmpDir, "ai-rules", "cursor-rules")
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Check file is removed
	if _, err := os.Stat(expectedFile); !os.IsNotExist(err) {
		t.Errorf("File %s should have been removed", expectedFile)
	}
}
