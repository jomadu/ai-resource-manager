package installer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestNewFileInstaller(t *testing.T) {
	installer := NewFileInstaller()
	if installer == nil {
		t.Error("Expected non-nil installer")
	}
}

func TestFileInstaller_Install(t *testing.T) {
	installer := NewFileInstaller()
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "installer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	files := []types.File{
		{Path: "rule1.json", Content: []byte(`{"rule": "test"}`), Size: 16},
		{Path: "subdir/rule2.json", Content: []byte(`{"rule": "nested"}`), Size: 18},
	}

	err = installer.Install(ctx, tempDir, "test-registry", "test-ruleset", "1.0.0", files)
	if err != nil {
		t.Errorf("Install failed: %v", err)
	}

	// Verify files were created in arm/registry/ruleset/version directory
	rule1Path := filepath.Join(tempDir, "arm", "test-registry", "test-ruleset", "1.0.0", "rule1.json")
	if _, err := os.Stat(rule1Path); os.IsNotExist(err) {
		t.Error("Expected rule1.json to be created")
	}

	rule2Path := filepath.Join(tempDir, "arm", "test-registry", "test-ruleset", "1.0.0", "subdir", "rule2.json")
	if _, err := os.Stat(rule2Path); os.IsNotExist(err) {
		t.Error("Expected subdir/rule2.json to be created")
	}

	// Verify content
	content, err := os.ReadFile(rule1Path)
	if err != nil {
		t.Errorf("Failed to read rule1.json: %v", err)
	}
	if string(content) != `{"rule": "test"}` {
		t.Errorf("Expected content %q, got %q", `{"rule": "test"}`, string(content))
	}
}

func TestFileInstaller_Uninstall(t *testing.T) {
	installer := NewFileInstaller()
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "installer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files in arm subdirectory
	rulesetDir := filepath.Join(tempDir, "arm", "test-registry", "test-ruleset")
	_ = os.MkdirAll(rulesetDir, 0o755)
	_ = os.WriteFile(filepath.Join(rulesetDir, "rule.json"), []byte("test"), 0o644)

	err = installer.Uninstall(ctx, tempDir, "test-registry", "test-ruleset")
	if err != nil {
		t.Errorf("Uninstall failed: %v", err)
	}

	// Verify ruleset directory was removed
	if _, err := os.Stat(rulesetDir); !os.IsNotExist(err) {
		t.Error("Expected ruleset directory to be removed")
	}
}

func TestFileInstaller_ListInstalled(t *testing.T) {
	installer := NewFileInstaller()
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "installer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test installations with version directories in arm/registry/ruleset subdirectory
	ruleset1Dir := filepath.Join(tempDir, "arm", "test-registry", "ruleset1", "1.0.0")
	ruleset2Dir := filepath.Join(tempDir, "arm", "test-registry", "ruleset2", "2.1.0")
	_ = os.MkdirAll(ruleset1Dir, 0o755)
	_ = os.MkdirAll(ruleset2Dir, 0o755)

	installations, err := installer.ListInstalled(ctx, tempDir)
	if err != nil {
		t.Errorf("ListInstalled failed: %v", err)
	}

	if len(installations) != 2 {
		t.Errorf("Expected 2 installations, got %d", len(installations))
	}

	// Verify installation details
	found := make(map[string]string)
	for _, inst := range installations {
		found[inst.Ruleset] = inst.Version
	}

	if found["ruleset1"] != "1.0.0" {
		t.Errorf("Expected ruleset1 version 1.0.0, got %s", found["ruleset1"])
	}
	if found["ruleset2"] != "2.1.0" {
		t.Errorf("Expected ruleset2 version 2.1.0, got %s", found["ruleset2"])
	}
}

func TestFileInstaller_InstallEmptyFiles(t *testing.T) {
	installer := NewFileInstaller()
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "installer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	err = installer.Install(ctx, tempDir, "test-registry", "empty-ruleset", "1.0.0", []types.File{})
	if err != nil {
		t.Errorf("Install with empty files failed: %v", err)
	}
}

func TestFileInstaller_InstallInvalidPath(t *testing.T) {
	installer := NewFileInstaller()
	ctx := context.Background()

	files := []types.File{{Path: "test.json", Content: []byte("test"), Size: 4}}

	err := installer.Install(ctx, "/nonexistent/path", "test-registry", "test-ruleset", "1.0.0", files)
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestFileInstaller_ImplementsInterface(t *testing.T) {
	var _ Installer = (*FileInstaller)(nil)
}

func TestFileInstaller_UninstallCleansUpEmptyDirectories(t *testing.T) {
	installer := NewFileInstaller()
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "installer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files for last ruleset in registry
	rulesetDir := filepath.Join(tempDir, "arm", "test-registry", "last-ruleset")
	_ = os.MkdirAll(rulesetDir, 0o755)
	_ = os.WriteFile(filepath.Join(rulesetDir, "rule.json"), []byte("test"), 0o644)

	err = installer.Uninstall(ctx, tempDir, "test-registry", "last-ruleset")
	if err != nil {
		t.Errorf("Uninstall failed: %v", err)
	}

	// Verify registry directory was removed
	registryDir := filepath.Join(tempDir, "arm", "test-registry")
	if _, err := os.Stat(registryDir); !os.IsNotExist(err) {
		t.Error("Expected empty registry directory to be removed")
	}

	// Verify arm directory was removed
	armDir := filepath.Join(tempDir, "arm")
	if _, err := os.Stat(armDir); !os.IsNotExist(err) {
		t.Error("Expected empty arm directory to be removed")
	}
}

func TestFileInstaller_UninstallKeepsNonEmptyDirectories(t *testing.T) {
	installer := NewFileInstaller()
	ctx := context.Background()

	tempDir, err := os.MkdirTemp("", "installer_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files for multiple rulesets
	ruleset1Dir := filepath.Join(tempDir, "arm", "test-registry", "ruleset1")
	ruleset2Dir := filepath.Join(tempDir, "arm", "test-registry", "ruleset2")
	_ = os.MkdirAll(ruleset1Dir, 0o755)
	_ = os.MkdirAll(ruleset2Dir, 0o755)
	_ = os.WriteFile(filepath.Join(ruleset1Dir, "rule.json"), []byte("test"), 0o644)
	_ = os.WriteFile(filepath.Join(ruleset2Dir, "rule.json"), []byte("test"), 0o644)

	err = installer.Uninstall(ctx, tempDir, "test-registry", "ruleset1")
	if err != nil {
		t.Errorf("Uninstall failed: %v", err)
	}

	// Verify registry directory still exists (has ruleset2)
	registryDir := filepath.Join(tempDir, "arm", "test-registry")
	if _, err := os.Stat(registryDir); os.IsNotExist(err) {
		t.Error("Expected non-empty registry directory to remain")
	}

	// Verify arm directory still exists
	armDir := filepath.Join(tempDir, "arm")
	if _, err := os.Stat(armDir); os.IsNotExist(err) {
		t.Error("Expected non-empty arm directory to remain")
	}
}
