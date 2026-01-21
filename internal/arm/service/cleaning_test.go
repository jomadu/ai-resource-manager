package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/sink"
	"github.com/jomadu/ai-resource-manager/internal/arm/storage"
)

func TestCleanCacheByAge(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, ".arm", "storage")

	// Create test registry with packages
	registryKey := map[string]interface{}{
		"url":  "https://example.com",
		"type": "git",
	}
	reg, err := storage.NewRegistryWithPath(filepath.Join(tempDir, ".arm"), registryKey)
	if err != nil {
		t.Fatal(err)
	}

	packageCache := storage.NewPackageCache(reg.GetPackagesDir())
	packageKey := map[string]interface{}{"name": "test-package"}

	// Create old version (8 days old)
	oldVersion := core.Version{Major: 1, Minor: 0, Patch: 0}
	if err := packageCache.SetPackageVersion(ctx, packageKey, oldVersion, []*core.File{
		{Path: "test.txt", Content: []byte("old")},
	}); err != nil {
		t.Fatal(err)
	}

	// Manually set old timestamp
	oldMetadataPath := filepath.Join(reg.GetPackagesDir(), mustGenerateKey(packageKey), "v1.0.0", "metadata.json")
	oldTime := time.Now().Add(-8 * 24 * time.Hour)
	setMetadataTime(t, oldMetadataPath, oldTime)

	// Create new version (1 day old)
	newVersion := core.Version{Major: 2, Minor: 0, Patch: 0}
	if err := packageCache.SetPackageVersion(ctx, packageKey, newVersion, []*core.File{
		{Path: "test.txt", Content: []byte("new")},
	}); err != nil {
		t.Fatal(err)
	}

	// Clean cache older than 7 days
	service := &ArmService{}
	if err := service.cleanCacheByAgeWithPath(ctx, 7*24*time.Hour, storageDir); err != nil {
		t.Fatal(err)
	}

	// Check old version removed
	versions, err := packageCache.ListPackageVersions(ctx, packageKey)
	if err != nil {
		t.Fatal(err)
	}

	if len(versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(versions))
	}

	if versions[0].Major != 2 {
		t.Errorf("expected version 2.0.0, got %d.%d.%d", versions[0].Major, versions[0].Minor, versions[0].Patch)
	}
}

func TestNukeCache(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, ".arm", "storage")

	// Create test registry with packages
	registryKey := map[string]interface{}{
		"url":  "https://example.com",
		"type": "git",
	}
	reg, err := storage.NewRegistryWithPath(filepath.Join(tempDir, ".arm"), registryKey)
	if err != nil {
		t.Fatal(err)
	}

	packageCache := storage.NewPackageCache(reg.GetPackagesDir())
	packageKey := map[string]interface{}{"name": "test-package"}
	version := core.Version{Major: 1, Minor: 0, Patch: 0}

	if err := packageCache.SetPackageVersion(ctx, packageKey, version, []*core.File{
		{Path: "test.txt", Content: []byte("test")},
	}); err != nil {
		t.Fatal(err)
	}

	// Nuke cache
	service := &ArmService{}
	if err := service.nukeCacheWithPath(ctx, storageDir); err != nil {
		t.Fatal(err)
	}

	// Check storage directory removed
	if _, err := os.Stat(storageDir); !os.IsNotExist(err) {
		t.Error("storage directory should be removed")
	}
}

// Helper functions
func mustGenerateKey(obj interface{}) string {
	key, err := storage.GenerateKey(obj)
	if err != nil {
		panic(err)
	}
	return key
}

func setMetadataTime(t *testing.T, metadataPath string, timestamp time.Time) {
	t.Helper()
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		t.Fatal(err)
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		t.Fatal(err)
	}

	metadata["updatedAt"] = timestamp.Format(time.RFC3339)
	newData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(metadataPath, newData, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestCleanSinks(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// Create sink directories
	sinkDir := filepath.Join(tempDir, ".cursor", "rules")

	mgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Sinks: map[string]manifest.SinkConfig{
				"cursor-rules": {Directory: sinkDir, Tool: compiler.Cursor},
			},
		},
	}

	// Create sink manager and install package
	sinkMgr := sink.NewManager(sinkDir, compiler.Cursor)
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-registry",
			Name:         "test-package",
			Version:      core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
		},
		Files: []*core.File{
			{Path: "test.md", Content: []byte("test content")},
		},
	}
	if err := sinkMgr.InstallRuleset(pkg, 100); err != nil {
		t.Fatal(err)
	}

	// Create orphaned file
	armDir := filepath.Join(sinkDir, "arm")
	orphanedFile := filepath.Join(armDir, "orphaned.md")
	if err := os.MkdirAll(filepath.Dir(orphanedFile), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(orphanedFile, []byte("orphaned"), 0644); err != nil {
		t.Fatal(err)
	}

	// Clean sinks
	service := NewArmService(mgr, nil, nil)
	if err := service.CleanSinks(ctx); err != nil {
		t.Fatal(err)
	}

	// Check orphaned file removed
	if _, err := os.Stat(orphanedFile); !os.IsNotExist(err) {
		t.Error("orphaned file should be removed")
	}

	// Check tracked files still exist
	index, err := sinkMgr.ListRulesets()
	if err != nil {
		t.Fatal(err)
	}
	if len(index) != 1 {
		t.Errorf("expected 1 tracked package, got %d", len(index))
	}
	if len(index) > 0 {
		for _, file := range index[0].Files {
			fullPath := filepath.Join(sinkDir, file)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("tracked file should exist: %s", file)
			}
		}
	}
}

func TestNukeSinks(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// Create sink directories
	cursorDir := filepath.Join(tempDir, ".cursor", "rules")
	copilotDir := filepath.Join(tempDir, ".github", "copilot")

	mgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Sinks: map[string]manifest.SinkConfig{
				"cursor-rules":  {Directory: cursorDir, Tool: compiler.Cursor},
				"copilot-rules": {Directory: copilotDir, Tool: compiler.Copilot},
			},
		},
	}

	// Install packages
	cursorMgr := sink.NewManager(cursorDir, compiler.Cursor)
	pkg := &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: "test-registry",
			Name:         "test-package",
			Version:      core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
		},
		Files: []*core.File{
			{Path: "test.md", Content: []byte("test content")},
		},
	}
	if err := cursorMgr.InstallRuleset(pkg, 100); err != nil {
		t.Fatal(err)
	}

	copilotMgr := sink.NewManager(copilotDir, compiler.Copilot)
	if err := copilotMgr.InstallRuleset(pkg, 100); err != nil {
		t.Fatal(err)
	}

	// Nuke sinks
	service := NewArmService(mgr, nil, nil)
	if err := service.NukeSinks(ctx); err != nil {
		t.Fatal(err)
	}

	// Check hierarchical sink arm directory removed
	armDir := filepath.Join(cursorDir, "arm")
	if _, err := os.Stat(armDir); !os.IsNotExist(err) {
		t.Error("arm directory should be removed")
	}

	// Check flat sink arm files removed
	entries, err := os.ReadDir(copilotDir)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if len(name) >= 4 && name[:4] == "arm_" || name == "arm-index.json" {
			t.Errorf("arm file should be removed: %s", name)
		}
	}
}

func TestCleanCacheByTimeSinceLastAccess(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()
	storageDir := filepath.Join(tempDir, ".arm", "storage")

	registryKey := map[string]interface{}{
		"url":  "https://example.com",
		"type": "git",
	}
	reg, err := storage.NewRegistryWithPath(filepath.Join(tempDir, ".arm"), registryKey)
	if err != nil {
		t.Fatal(err)
	}

	packageCache := storage.NewPackageCache(reg.GetPackagesDir())
	packageKey := map[string]interface{}{"name": "test-package"}

	// Create old accessed version (8 days since last access)
	oldVersion := core.Version{Major: 1, Minor: 0, Patch: 0}
	if err := packageCache.SetPackageVersion(ctx, packageKey, oldVersion, []*core.File{
		{Path: "test.txt", Content: []byte("old")},
	}); err != nil {
		t.Fatal(err)
	}

	oldMetadataPath := filepath.Join(reg.GetPackagesDir(), mustGenerateKey(packageKey), "v1.0.0", "metadata.json")
	oldAccessTime := time.Now().Add(-8 * 24 * time.Hour)
	setMetadataAccessTime(t, oldMetadataPath, oldAccessTime)

	// Create recently accessed version (1 day since last access)
	newVersion := core.Version{Major: 2, Minor: 0, Patch: 0}
	if err := packageCache.SetPackageVersion(ctx, packageKey, newVersion, []*core.File{
		{Path: "test.txt", Content: []byte("new")},
	}); err != nil {
		t.Fatal(err)
	}

	// Clean cache not accessed in 7 days
	service := &ArmService{}
	if err := service.cleanCacheByTimeSinceLastAccessWithPath(ctx, 7*24*time.Hour, storageDir); err != nil {
		t.Fatal(err)
	}

	// Check old version removed
	versions, err := packageCache.ListPackageVersions(ctx, packageKey)
	if err != nil {
		t.Fatal(err)
	}

	if len(versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(versions))
	}

	if versions[0].Major != 2 {
		t.Errorf("expected version 2.0.0, got %d.%d.%d", versions[0].Major, versions[0].Minor, versions[0].Patch)
	}
}

func setMetadataAccessTime(t *testing.T, metadataPath string, timestamp time.Time) {
	t.Helper()
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		t.Fatal(err)
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		t.Fatal(err)
	}

	metadata["accessedAt"] = timestamp.Format(time.RFC3339)
	newData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(metadataPath, newData, 0644); err != nil {
		t.Fatal(err)
	}
}
