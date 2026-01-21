package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/arm/registry"
)

func TestListAll(t *testing.T) {
	t.Run("list all dependencies with config and lock info", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, "lock.json")

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks: map[string]manifest.SinkConfig{
					"test-sink": {Directory: "/path", Tool: compiler.Cursor},
				},
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
					"test-reg/promptset1": {
						"type":    "promptset",
						"version": "2.0.0",
						"sinks":   []interface{}{"test-sink"},
					},
				},
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(lockPath)
		lockMgr.UpsertDependencyLock(context.Background(), "test-reg", "ruleset1", "1.0.0", &packagelockfile.DependencyLockConfig{
			Integrity: "sha256-ruleset1",
		})
		lockMgr.UpsertDependencyLock(context.Background(), "test-reg", "promptset1", "2.0.0", &packagelockfile.DependencyLockConfig{
			Integrity: "sha256-promptset1",
		})

		svc := NewArmService(mgr, lockMgr, nil)

		deps, err := svc.ListAll(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(deps) != 2 {
			t.Fatalf("expected 2 dependencies, got %d", len(deps))
		}
	})

	t.Run("list all when no dependencies", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, "lock.json")

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries:   make(map[string]map[string]interface{}),
				Sinks:        make(map[string]manifest.SinkConfig),
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(lockPath)
		svc := NewArmService(mgr, lockMgr, nil)

		deps, err := svc.ListAll(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(deps) != 0 {
			t.Fatalf("expected 0 dependencies, got %d", len(deps))
		}
	})
}

func TestGetDependencyInfo(t *testing.T) {
	t.Run("get existing dependency info", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, "lock.json")

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks: map[string]manifest.SinkConfig{
					"test-sink": {Directory: "/path", Tool: compiler.Cursor},
				},
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
				},
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(lockPath)
		lockMgr.UpsertDependencyLock(context.Background(), "test-reg", "ruleset1", "1.0.0", &packagelockfile.DependencyLockConfig{
			Integrity: "sha256-test",
		})

		svc := NewArmService(mgr, lockMgr, nil)

		info, err := svc.GetDependencyInfo(context.Background(), "test-reg", "ruleset1")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if info == nil {
			t.Fatal("expected dependency info, got nil")
		}
	})

	t.Run("get non-existent dependency", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, "lock.json")

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries:   make(map[string]map[string]interface{}),
				Sinks:        make(map[string]manifest.SinkConfig),
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(lockPath)
		svc := NewArmService(mgr, lockMgr, nil)

		_, err := svc.GetDependencyInfo(context.Background(), "test-reg", "nonexistent")

		if err == nil {
			t.Fatal("expected error when dependency does not exist")
		}
	})
}

func TestListOutdated(t *testing.T) {
	t.Run("list outdated dependencies", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, "lock.json")

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks: map[string]manifest.SinkConfig{
					"test-sink": {Directory: "/path", Tool: compiler.Cursor},
				},
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "^1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
				},
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(lockPath)
		lockMgr.UpsertDependencyLock(context.Background(), "test-reg", "ruleset1", "1.0.0", &packagelockfile.DependencyLockConfig{
			Integrity: "sha256-test",
		})

		mockReg := registry.NewMockRegistry()
		pkg1 := &core.Package{
			Metadata: core.PackageMetadata{
				Name:         "ruleset1",
				RegistryName: "test-reg",
				Version:      core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			},
		}
		pkg2 := &core.Package{
			Metadata: core.PackageMetadata{
				Name:         "ruleset1",
				RegistryName: "test-reg",
				Version:      core.Version{Major: 1, Minor: 5, Patch: 0, Version: "1.5.0", IsSemver: true},
			},
		}
		pkg3 := &core.Package{
			Metadata: core.PackageMetadata{
				Name:         "ruleset1",
				RegistryName: "test-reg",
				Version:      core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true},
			},
		}
		mockReg.AddPackage(pkg1)
		mockReg.AddPackage(pkg2)
		mockReg.AddPackage(pkg3)

		svc := NewArmService(mgr, lockMgr, registry.NewMockFactory(mockReg))

		outdated, err := svc.ListOutdated(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(outdated) != 1 {
			t.Fatalf("expected 1 outdated dependency, got %d", len(outdated))
		}
		if outdated[0].Current.Version.Version != "1.0.0" {
			t.Errorf("expected current version 1.0.0, got %s", outdated[0].Current.Version.Version)
		}
		if outdated[0].Wanted.Version.Version != "1.5.0" {
			t.Errorf("expected wanted version 1.5.0, got %s", outdated[0].Wanted.Version.Version)
		}
		if outdated[0].Latest.Version.Version != "2.0.0" {
			t.Errorf("expected latest version 2.0.0, got %s", outdated[0].Latest.Version.Version)
		}
	})

	t.Run("list outdated when all up to date", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, "lock.json")

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks: map[string]manifest.SinkConfig{
					"test-sink": {Directory: "/path", Tool: compiler.Cursor},
				},
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
				},
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(lockPath)
		lockMgr.UpsertDependencyLock(context.Background(), "test-reg", "ruleset1", "1.0.0", &packagelockfile.DependencyLockConfig{
			Integrity: "sha256-test",
		})

		mockReg := registry.NewMockRegistry()
		pkg := &core.Package{
			Metadata: core.PackageMetadata{
				Name:         "ruleset1",
				RegistryName: "test-reg",
				Version:      core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			},
		}
		mockReg.AddPackage(pkg)

		svc := NewArmService(mgr, lockMgr, registry.NewMockFactory(mockReg))

		outdated, err := svc.ListOutdated(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(outdated) != 0 {
			t.Fatalf("expected 0 outdated dependencies, got %d", len(outdated))
		}
	})
}
