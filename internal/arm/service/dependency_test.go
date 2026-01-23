package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/arm/registry"
)

func TestInstallPromptset(t *testing.T) {
	t.Run("install promptset successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		sinkDir := filepath.Join(tmpDir, "sink")
		_ = os.MkdirAll(sinkDir, 0o755)

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks: map[string]manifest.SinkConfig{
					"test-sink": {Directory: sinkDir, Tool: compiler.Cursor},
				},
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(filepath.Join(tmpDir, "lock.json"))

		mockReg := registry.NewMockRegistry()
		pkg := &core.Package{
			Metadata: core.PackageMetadata{
				Name:         "test-promptset",
				RegistryName: "test-reg",
				Version:      core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			},
			Integrity: "sha256-test",
			Files:     []*core.File{},
		}
		mockReg.AddPackage(pkg)

		svc := NewArmService(mgr, lockMgr, registry.NewMockFactory(mockReg))

		err := svc.InstallPromptset(context.Background(), "test-reg", "test-promptset", "1.0.0", nil, nil, []string{"test-sink"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if mgr.manifest.Dependencies["test-reg/test-promptset"] == nil {
			t.Fatal("dependency not added to manifest")
		}
	})

	t.Run("fail when registry does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries:   make(map[string]map[string]interface{}),
				Sinks:        make(map[string]manifest.SinkConfig),
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.InstallPromptset(context.Background(), "nonexistent", "test-promptset", "1.0.0", nil, nil, []string{"test-sink"})

		if err == nil {
			t.Fatal("expected error when registry does not exist")
		}
	})

	t.Run("fail when sink does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks:        make(map[string]manifest.SinkConfig),
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.InstallPromptset(context.Background(), "test-reg", "test-promptset", "1.0.0", nil, nil, []string{"nonexistent"})

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})
}

func TestInstallRuleset(t *testing.T) {
	t.Run("install ruleset successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		sinkDir := filepath.Join(tmpDir, "sink")
		_ = os.MkdirAll(sinkDir, 0o755)

		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks: map[string]manifest.SinkConfig{
					"test-sink": {Directory: sinkDir, Tool: compiler.Cursor},
				},
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		lockMgr := packagelockfile.NewFileManagerWithPath(filepath.Join(tmpDir, "lock.json"))

		mockReg := registry.NewMockRegistry()
		pkg := &core.Package{
			Metadata: core.PackageMetadata{
				Name:         "test-ruleset",
				RegistryName: "test-reg",
				Version:      core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			},
			Integrity: "sha256-test",
			Files:     []*core.File{},
		}
		mockReg.AddPackage(pkg)

		svc := NewArmService(mgr, lockMgr, registry.NewMockFactory(mockReg))

		err := svc.InstallRuleset(context.Background(), "test-reg", "test-ruleset", "1.0.0", 100, nil, nil, []string{"test-sink"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if mgr.manifest.Dependencies["test-reg/test-ruleset"] == nil {
			t.Fatal("dependency not added to manifest")
		}
	})

	t.Run("fail when registry does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries:   make(map[string]map[string]interface{}),
				Sinks:        make(map[string]manifest.SinkConfig),
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.InstallRuleset(context.Background(), "nonexistent", "test-ruleset", "1.0.0", 100, nil, nil, []string{"test-sink"})

		if err == nil {
			t.Fatal("expected error when registry does not exist")
		}
	})

	t.Run("fail when sink does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test-reg": {"type": "mock"},
				},
				Sinks:        make(map[string]manifest.SinkConfig),
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.InstallRuleset(context.Background(), "test-reg", "test-ruleset", "1.0.0", 100, nil, nil, []string{"nonexistent"})

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})
}
