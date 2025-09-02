package arm

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
)

const testRegistry = "https://github.com/jomadu/ai-rules-manager-sample-git-registry"

func setupTest(t *testing.T) (*ArmService, context.Context) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(oldDir) })
	_ = os.Chdir(tmpDir)

	service := NewArmService()
	ctx := context.Background()

	// Setup manifest with registry
	manifestFile := manifest.Manifest{
		Registries: map[string]manifest.RegistryConfig{
			"ai-rules": {URL: testRegistry, Type: "git"},
		},
		Rulesets: map[string]map[string]manifest.Entry{},
	}
	data, _ := json.MarshalIndent(manifestFile, "", "  ")
	err := os.WriteFile("arm.json", data, 0o644)
	if err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}

	err = service.configManager.AddSink(ctx, "q", []string{".amazonq/rules"}, []string{"ai-rules/amazonq-*"}, []string{"ai-rules/cursor-*"})
	if err != nil {
		t.Fatalf("Failed to add q sink: %v", err)
	}

	err = service.configManager.AddSink(ctx, "cursor", []string{".cursor/rules"}, []string{"ai-rules/cursor-*"}, []string{"ai-rules/amazonq-*"})
	if err != nil {
		t.Fatalf("Failed to add cursor sink: %v", err)
	}

	return service, ctx
}

func TestIntegrationInstallLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service, ctx := setupTest(t)

	// Install latest versions with empty constraint (should find latest)
	err := service.InstallRuleset(ctx, "ai-rules", "amazonq-rules", "", []string{"rules/amazonq/*.md"}, nil)
	if err != nil {
		t.Fatalf("Failed to install amazonq-rules: %v", err)
	}

	err = service.InstallRuleset(ctx, "ai-rules", "cursor-rules", "", []string{"rules/cursor/*.mdc"}, nil)
	if err != nil {
		t.Fatalf("Failed to install cursor-rules: %v", err)
	}

	// Verify files created
	assertFileExists(t, "arm.json")
	assertFileExists(t, "arm-lock.json")
	assertFileExists(t, ".amazonq/rules/arm")
	assertFileExists(t, ".cursor/rules/arm")

	// Verify manifest content
	manifestData, err := os.ReadFile("arm.json")
	if err != nil {
		t.Fatalf("Failed to read manifest: %v", err)
	}
	var manifestFile manifest.Manifest
	if err := json.Unmarshal(manifestData, &manifestFile); err != nil {
		t.Fatalf("Failed to parse manifest: %v", err)
	}
	if len(manifestFile.Rulesets["ai-rules"]) != 2 {
		t.Errorf("Expected 2 rulesets, got %d", len(manifestFile.Rulesets["ai-rules"]))
	}
}

func TestIntegrationList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service, ctx := setupTest(t)

	// Install rulesets first
	err := service.InstallRuleset(ctx, "ai-rules", "amazonq-rules", "^2.1.0", []string{"rules/amazonq/*.md"}, nil)
	if err != nil {
		t.Fatalf("Failed to install amazonq-rules: %v", err)
	}

	err = service.InstallRuleset(ctx, "ai-rules", "cursor-rules", "^2.1.0", []string{"rules/cursor/*.mdc"}, nil)
	if err != nil {
		t.Fatalf("Failed to install cursor-rules: %v", err)
	}

	// Test list
	rulesets, err := service.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list rulesets: %v", err)
	}
	if len(rulesets) != 2 {
		t.Errorf("Expected 2 rulesets, got %d", len(rulesets))
	}
}

func TestIntegrationInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service, ctx := setupTest(t)

	// Install ruleset first
	err := service.InstallRuleset(ctx, "ai-rules", "amazonq-rules", "^2.1.0", []string{"rules/amazonq/*.md"}, nil)
	if err != nil {
		t.Fatalf("Failed to install amazonq-rules: %v", err)
	}

	// Test info
	info, err := service.Info(ctx, "ai-rules", "amazonq-rules")
	if err != nil {
		t.Fatalf("Failed to get info: %v", err)
	}
	if info.Registry != "ai-rules" {
		t.Errorf("Expected registry ai-rules, got %s", info.Registry)
	}
	if info.Name != "amazonq-rules" {
		t.Errorf("Expected name amazonq-rules, got %s", info.Name)
	}
}

func TestIntegrationUninstall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service, ctx := setupTest(t)

	// Install rulesets first
	err := service.InstallRuleset(ctx, "ai-rules", "amazonq-rules", "^2.1.0", []string{"rules/amazonq/*.md"}, nil)
	if err != nil {
		t.Fatalf("Failed to install amazonq-rules: %v", err)
	}

	err = service.InstallRuleset(ctx, "ai-rules", "cursor-rules", "^2.1.0", []string{"rules/cursor/*.mdc"}, nil)
	if err != nil {
		t.Fatalf("Failed to install cursor-rules: %v", err)
	}

	// Uninstall one
	err = service.Uninstall(ctx, "ai-rules", "cursor-rules")
	if err != nil {
		t.Fatalf("Failed to uninstall: %v", err)
	}

	// Verify only one ruleset remains
	rulesets, err := service.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list after uninstall: %v", err)
	}
	if len(rulesets) != 1 {
		t.Errorf("Expected 1 ruleset after uninstall, got %d", len(rulesets))
	}
	if rulesets[0].Name != "amazonq-rules" {
		t.Errorf("Expected amazonq-rules to remain, got %s", rulesets[0].Name)
	}
}

func TestIntegrationSpecificVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service, ctx := setupTest(t)

	// Install specific version
	err := service.InstallRuleset(ctx, "ai-rules", "cursor-rules", "1.0.0", []string{"rules/cursor/*.mdc"}, nil)
	if err != nil {
		t.Fatalf("Failed to install specific version: %v", err)
	}

	// Verify lockfile has specific version
	lockData, err := os.ReadFile("arm-lock.json")
	if err != nil {
		t.Fatalf("Failed to read lockfile: %v", err)
	}
	var lockFile lockfile.LockFile
	if err := json.Unmarshal(lockData, &lockFile); err != nil {
		t.Fatalf("Failed to parse lockfile: %v", err)
	}
	cursorEntry := lockFile.Rulesets["ai-rules"]["cursor-rules"]
	if cursorEntry.Constraint != "1.0.0" {
		t.Errorf("Expected constraint 1.0.0, got %s", cursorEntry.Constraint)
	}
}

func TestIntegrationLatestConstraint(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service, ctx := setupTest(t)

	// Install without version constraint (should normalize to "latest")
	err := service.InstallRuleset(ctx, "ai-rules", "amazonq-rules", "", []string{"rules/amazonq/*.md"}, nil)
	if err != nil {
		t.Fatalf("Failed to install with empty constraint: %v", err)
	}

	// Verify manifest shows "latest"
	manifestData, err := os.ReadFile("arm.json")
	if err != nil {
		t.Fatalf("Failed to read manifest: %v", err)
	}
	var manifestFile manifest.Manifest
	if err := json.Unmarshal(manifestData, &manifestFile); err != nil {
		t.Fatalf("Failed to parse manifest: %v", err)
	}
	entry := manifestFile.Rulesets["ai-rules"]["amazonq-rules"]
	if entry.Version != "latest" {
		t.Errorf("Expected manifest version 'latest', got '%s'", entry.Version)
	}

	// Verify lockfile shows "latest"
	lockData, err := os.ReadFile("arm-lock.json")
	if err != nil {
		t.Fatalf("Failed to read lockfile: %v", err)
	}
	var lockFile lockfile.LockFile
	if err := json.Unmarshal(lockData, &lockFile); err != nil {
		t.Fatalf("Failed to parse lockfile: %v", err)
	}
	lockEntry := lockFile.Rulesets["ai-rules"]["amazonq-rules"]
	if lockEntry.Constraint != "latest" {
		t.Errorf("Expected lockfile constraint 'latest', got '%s'", lockEntry.Constraint)
	}
}

func TestIntegrationNpmLikeBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service, ctx := setupTest(t)

	// Test no files
	err := service.Install(ctx)
	if err == nil {
		t.Fatal("Expected error with no files")
	}
	if err.Error() != "neither arm.json nor arm-lock.json found" {
		t.Errorf("Expected 'neither arm.json nor arm-lock.json found', got: %v", err)
	}

	// Test manifest only
	manifestFile := manifest.Manifest{
		Registries: map[string]manifest.RegistryConfig{
			"ai-rules": {URL: testRegistry, Type: "git"},
		},
		Rulesets: map[string]map[string]manifest.Entry{
			"ai-rules": {
				"amazonq-rules": {
					Version: "^1.0.0",
					Include: []string{"rules/amazonq/*.md"},
				},
			},
		},
	}
	data, _ := json.MarshalIndent(manifestFile, "", "  ")
	_ = os.WriteFile("arm.json", data, 0o644)

	err = service.Install(ctx)
	if err != nil {
		t.Fatalf("Failed to install from manifest only: %v", err)
	}

	// Should create lockfile
	assertFileExists(t, "arm-lock.json")

	// Test both files exist - should use lockfile
	err = service.Install(ctx)
	if err != nil {
		t.Fatalf("Failed to install with both files: %v", err)
	}
	assertFileExists(t, ".amazonq/rules/arm")

	// Test lockfile only
	_ = os.Remove("arm.json")
	_ = os.RemoveAll(".amazonq")

	err = service.Install(ctx)
	if err == nil {
		t.Fatal("Expected error with lockfile only")
	}
	if err.Error() != "arm.json not found" {
		t.Errorf("Expected 'arm.json not found', got: %v", err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file/directory %s to exist", path)
	}
}
