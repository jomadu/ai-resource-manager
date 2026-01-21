package service

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/arm/registry"
)

func TestInstallAll(t *testing.T) {
	ctx := context.Background()

	// Setup manifest with dependencies
	rulesetDep := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypeRuleset,
			Version: "1.0.0",
			Sinks:   []string{"test-sink"},
		},
		Priority: 100,
	}
	rulesetDepMap, _ := json.Marshal(rulesetDep)
	var rulesetDepInterface map[string]interface{}
	json.Unmarshal(rulesetDepMap, &rulesetDepInterface)

	promptsetDep := manifest.PromptsetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypePromptset,
			Version: "2.0.0",
			Sinks:   []string{"test-sink"},
		},
	}
	promptsetDepMap, _ := json.Marshal(promptsetDep)
	var promptsetDepInterface map[string]interface{}
	json.Unmarshal(promptsetDepMap, &promptsetDepInterface)

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Registries: map[string]map[string]interface{}{
				"test-registry": {
					"type": "git",
					"url":  "https://github.com/test/repo",
				},
			},
			Sinks: map[string]manifest.SinkConfig{
				"test-sink": {
					Directory: "/tmp/test",
					Tool:      compiler.Cursor,
				},
			},
			Dependencies: map[string]map[string]interface{}{
				"test-registry/ruleset1":   rulesetDepInterface,
				"test-registry/promptset1": promptsetDepInterface,
			},
		},
	}

	lockfileMgr := &mockLockfileManager{
		locks: make(map[string]*packagelockfile.DependencyLockConfig),
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"ruleset1": {
				{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			},
			"promptset1": {
				{Name: "promptset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"ruleset1@1.0.0": {
				Metadata:  core.PackageMetadata{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "hash1",
			},
			"promptset1@2.0.0": {
				Metadata:  core.PackageMetadata{Name: "promptset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "prompt.yml", Content: []byte("test")}},
				Integrity: "hash2",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}

	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	// Execute
	err := service.InstallAll(ctx)
	if err != nil {
		t.Fatalf("InstallAll() error = %v", err)
	}

	// Verify lockfile was updated
	if len(lockfileMgr.locks) != 2 {
		t.Errorf("Expected 2 lockfile entries, got %d", len(lockfileMgr.locks))
	}

	if _, exists := lockfileMgr.locks["test-registry/ruleset1@1.0.0"]; !exists {
		t.Error("Expected lockfile entry for ruleset1")
	}

	if _, exists := lockfileMgr.locks["test-registry/promptset1@2.0.0"]; !exists {
		t.Error("Expected lockfile entry for promptset1")
	}
}

// Mock implementations

type mockLockfileManager struct {
	locks map[string]*packagelockfile.DependencyLockConfig
}

func (m *mockLockfileManager) GetDependencyLock(ctx context.Context, registry, packageName, version string) (*packagelockfile.DependencyLockConfig, error) {
	key := registry + "/" + packageName + "@" + version
	if lock, ok := m.locks[key]; ok {
		return lock, nil
	}
	return nil, nil
}

func (m *mockLockfileManager) GetLockFile(ctx context.Context) (*packagelockfile.LockFile, error) {
	lockFile := &packagelockfile.LockFile{
		Version:      1,
		Dependencies: make(map[string]packagelockfile.DependencyLockConfig),
	}
	for key, config := range m.locks {
		lockFile.Dependencies[key] = *config
	}
	return lockFile, nil
}

func (m *mockLockfileManager) UpsertDependencyLock(ctx context.Context, registry, packageName, version string, config *packagelockfile.DependencyLockConfig) error {
	key := registry + "/" + packageName + "@" + version
	m.locks[key] = config
	return nil
}

func (m *mockLockfileManager) RemoveDependencyLock(ctx context.Context, registry, packageName, version string) error {
	key := registry + "/" + packageName + "@" + version
	delete(m.locks, key)
	return nil
}

func (m *mockLockfileManager) UpdateRegistryName(ctx context.Context, oldName, newName string) error {
	return nil
}

type mockRegistry struct {
	versions map[string][]core.PackageMetadata
	packages map[string]*core.Package
}

func (m *mockRegistry) ListPackages(ctx context.Context) ([]*core.PackageMetadata, error) {
	return nil, nil
}

func (m *mockRegistry) ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error) {
	metadata := m.versions[packageName]
	versions := make([]core.Version, len(metadata))
	for i, v := range metadata {
		versions[i] = v.Version
	}
	// Sort versions highest first (like real registries do)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) > 0
	})
	return versions, nil
}

func (m *mockRegistry) GetPackage(ctx context.Context, packageName string, version core.Version, include, exclude []string) (*core.Package, error) {
	key := packageName + "@" + version.Version
	return m.packages[key], nil
}

type mockRegistryFactory struct {
	registry registry.Registry
}

func (f *mockRegistryFactory) CreateRegistry(name string, config map[string]interface{}) (registry.Registry, error) {
	return f.registry, nil
}
