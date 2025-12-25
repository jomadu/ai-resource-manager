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
		name          string
		directory     string
		compileTarget compiler.CompileTarget
	}{
		{
			name:          "valid cursor target",
			directory:     "/tmp/test-sink",
			compileTarget: compiler.TargetCursor,
		},
		{
			name:          "valid amazonq target",
			directory:     "/tmp/test-sink",
			compileTarget: compiler.TargetAmazonQ,
		},
		{
			name:          "valid copilot target",
			directory:     "/tmp/test-sink",
			compileTarget: compiler.TargetCopilot,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.RemoveAll(tt.directory)
			
			manager := NewManager(tt.directory, tt.compileTarget)
			
			if manager == nil {
				t.Errorf("NewManager() returned nil for valid input")
				return
			}

			// Check manager fields
			if manager.directory != tt.directory {
				t.Errorf("directory = %v, want %v", manager.directory, tt.directory)
			}
			if manager.compileTarget != string(tt.compileTarget) {
				t.Errorf("compileTarget = %v, want %v", manager.compileTarget, string(tt.compileTarget))
			}

			// Check paths are set correctly
			expectedIndexPath := filepath.Join(tt.directory, "arm", "arm-index.json")
			if manager.indexPath != expectedIndexPath {
				t.Errorf("indexPath = %v, want %v", manager.indexPath, expectedIndexPath)
			}

			expectedArmDir := filepath.Join(tt.directory, "arm")
			if manager.armDir != expectedArmDir {
				t.Errorf("armDir = %v, want %v", manager.armDir, expectedArmDir)
			}

			// Check directory was created
			if _, err := os.Stat(tt.directory); os.IsNotExist(err) {
				t.Errorf("directory %v was not created", tt.directory)
			}

			// Clean up after test
			os.RemoveAll(tt.directory)
		})
	}
}

func TestManager_IsInstalled(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, compiler.TargetCursor)
	
	metadata := core.PackageMetadata{
		PackageId:  core.PackageId{ID: "test-pkg", Name: "test-pkg"},
		RegistryId: core.RegistryId{ID: "test-reg", Name: "test-reg"},
		Version:    core.Version{Version: "1.0.0"},
	}

	// Package not installed initially
	if manager.IsInstalled(metadata) {
		t.Errorf("IsInstalled() = true, want false for non-existent package")
	}
}

func TestManager_InstallRuleset(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, compiler.TargetCursor)

	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			PackageId:  core.PackageId{ID: "clean-code", Name: "clean-code"},
			RegistryId: core.RegistryId{ID: "test-reg", Name: "test-reg"},
			Version:    core.Version{Version: "1.0.0"},
		},
		Files: []core.File{
			{
				Path:    "rule1.yml",
				Content: []byte("test rule content"),
				Size:    17,
			},
		},
	}

	err := manager.InstallRuleset(pkg, 100)
	if err != nil {
		t.Errorf("InstallRuleset() error = %v, want nil", err)
	}

	// Check package is now installed
	if !manager.IsInstalled(pkg.Metadata) {
		t.Errorf("IsInstalled() = false, want true after install")
	}
}

func TestManager_InstallPromptset(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, compiler.TargetAmazonQ)

	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			PackageId:  core.PackageId{ID: "code-review", Name: "code-review"},
			RegistryId: core.RegistryId{ID: "test-reg", Name: "test-reg"},
			Version:    core.Version{Version: "1.0.0"},
		},
		Files: []core.File{
			{
				Path:    "prompt1.yml",
				Content: []byte("test prompt content"),
				Size:    19,
			},
		},
	}

	err := manager.InstallPromptset(pkg)
	if err != nil {
		t.Errorf("InstallPromptset() error = %v, want nil", err)
	}

	// Check package is now installed
	if !manager.IsInstalled(pkg.Metadata) {
		t.Errorf("IsInstalled() = false, want true after install")
	}
}

func TestManager_Uninstall(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, compiler.TargetCursor)

	// Install a package first
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			PackageId:  core.PackageId{ID: "test-pkg", Name: "test-pkg"},
			RegistryId: core.RegistryId{ID: "test-reg", Name: "test-reg"},
			Version:    core.Version{Version: "1.0.0"},
		},
		Files: []core.File{
			{
				Path:    "rule1.yml",
				Content: []byte("test content"),
				Size:    12,
			},
		},
	}

	manager.InstallRuleset(pkg, 100)

	// Uninstall the package
	err := manager.Uninstall(pkg.Metadata)
	if err != nil {
		t.Errorf("Uninstall() error = %v, want nil", err)
	}

	// Check package is no longer installed
	if manager.IsInstalled(pkg.Metadata) {
		t.Errorf("IsInstalled() = true, want false after uninstall")
	}
}

func TestManager_ListRulesets(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, compiler.TargetCursor)

	// Install a ruleset
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			PackageId:  core.PackageId{ID: "clean-code", Name: "clean-code"},
			RegistryId: core.RegistryId{ID: "test-reg", Name: "test-reg"},
			Version:    core.Version{Version: "1.0.0"},
		},
		Files: []core.File{
			{
				Path:    "rule1.yml",
				Content: []byte("test rule"),
				Size:    9,
			},
		},
	}

	manager.InstallRuleset(pkg, 100)

	rulesets, err := manager.ListRulesets()
	if err != nil {
		t.Errorf("ListRulesets() error = %v, want nil", err)
	}

	if len(rulesets) == 0 {
		t.Errorf("ListRulesets() returned empty map, want 1 ruleset")
	}

	if rulesets["test-reg"] == nil || rulesets["test-reg"]["clean-code"] == nil {
		t.Errorf("ListRulesets() missing expected ruleset")
	}
}

func TestManager_ListPromptsets(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir, compiler.TargetAmazonQ)

	// Install a promptset
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			PackageId:  core.PackageId{ID: "code-review", Name: "code-review"},
			RegistryId: core.RegistryId{ID: "test-reg", Name: "test-reg"},
			Version:    core.Version{Version: "1.0.0"},
		},
		Files: []core.File{
			{
				Path:    "prompt1.yml",
				Content: []byte("test prompt"),
				Size:    11,
			},
		},
	}

	manager.InstallPromptset(pkg)

	promptsets, err := manager.ListPromptsets()
	if err != nil {
		t.Errorf("ListPromptsets() error = %v, want nil", err)
	}

	if len(promptsets) == 0 {
		t.Errorf("ListPromptsets() returned empty map, want 1 promptset")
	}

	if promptsets["test-reg"] == nil || promptsets["test-reg"]["code-review"] == nil {
		t.Errorf("ListPromptsets() missing expected promptset")
	}
}

func TestIndex_GetRuleset(t *testing.T) {
	index := &Index{
		Packages: map[string]map[string]interface{}{
			"test-reg": {
				"clean-code": &RulesetEntry{
					Version:      "1.0.0",
					ResourceType: "ruleset",
					Priority:     100,
					Files:        []string{"rule1.mdc"},
				},
			},
		},
	}

	entry, err := index.GetRuleset("test-reg", "clean-code")
	if err != nil {
		t.Errorf("GetRuleset() error = %v, want nil", err)
	}

	if entry.Priority != 100 {
		t.Errorf("GetRuleset() priority = %v, want 100", entry.Priority)
	}
}

func TestIndex_GetPromptset(t *testing.T) {
	index := &Index{
		Packages: map[string]map[string]interface{}{
			"test-reg": {
				"code-review": &PromptsetEntry{
					Version:      "1.0.0",
					ResourceType: "promptset",
					Files:        []string{"prompt1.md"},
				},
			},
		},
	}

	entry, err := index.GetPromptset("test-reg", "code-review")
	if err != nil {
		t.Errorf("GetPromptset() error = %v, want nil", err)
	}

	if entry.Version != "1.0.0" {
		t.Errorf("GetPromptset() version = %v, want 1.0.0", entry.Version)
	}
}