package installer

import (
	"context"
	"os"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestFlatInstallerHashFile(t *testing.T) {
	installer := NewFlatInstaller("/tmp", "markdown")

	tests := []struct {
		name     string
		registry string
		ruleset  string
		version  string
		filePath string
	}{
		{
			name:     "basic hash",
			registry: "test-registry",
			ruleset:  "test-ruleset",
			version:  "1.0.0",
			filePath: "rules/test.md",
		},
		{
			name:     "different registry",
			registry: "other-registry",
			ruleset:  "test-ruleset",
			version:  "1.0.0",
			filePath: "rules/test.md",
		},
		{
			name:     "different version",
			registry: "test-registry",
			ruleset:  "test-ruleset",
			version:  "2.0.0",
			filePath: "rules/test.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := installer.hashFile(tt.registry, tt.ruleset, tt.version, tt.filePath)
			if len(got) != 8 {
				t.Errorf("hashFile() = %v, want 8 character hash", got)
			}
		})
	}
}

func TestFlatInstallerHashFileDeterministic(t *testing.T) {
	installer := NewFlatInstaller("/tmp", "markdown")

	hash1 := installer.hashFile("registry", "ruleset", "1.0.0", "file.md")
	hash2 := installer.hashFile("registry", "ruleset", "1.0.0", "file.md")

	if hash1 != hash2 {
		t.Errorf("hashFile() not deterministic: %v != %v", hash1, hash2)
	}
}

func TestFlatInstallerHashFileUnique(t *testing.T) {
	installer := NewFlatInstaller("/tmp", "markdown")

	hash1 := installer.hashFile("registry1", "ruleset", "1.0.0", "file.md")
	hash2 := installer.hashFile("registry2", "ruleset", "1.0.0", "file.md")

	if hash1 == hash2 {
		t.Errorf("hashFile() not unique for different registries: %v == %v", hash1, hash2)
	}
}

func TestFlatInstaller_ListInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	installer := NewFlatInstaller(tmpDir, "markdown")
	ctx := context.Background()

	// Install test files
	files := []types.File{
		{Path: "rule1.md", Content: []byte("rule1")},
		{Path: "rule2.md", Content: []byte("rule2")},
	}

	err := installer.Install(ctx, "registry1", "ruleset1", "1.0.0", 100, files)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Test ListInstalled
	installations, err := installer.ListInstalled(ctx)
	if err != nil {
		t.Fatalf("ListInstalled failed: %v", err)
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
