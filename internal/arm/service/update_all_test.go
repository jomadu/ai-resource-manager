package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
)

// Normal: update multiple dependencies with newer versions available
func TestUpdateAll_Normal(t *testing.T) {
	ctx := context.Background()

	rulesetDep := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypeRuleset,
			Version: "^1.0.0",
			Sinks:   []string{"test-sink"},
		},
		Priority: 100,
	}
	rulesetDepMap, _ := json.Marshal(rulesetDep)
	var rulesetDepInterface map[string]interface{}
	_ = json.Unmarshal(rulesetDepMap, &rulesetDepInterface)

	promptsetDep := manifest.PromptsetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypePromptset,
			Version: "^2.0.0",
			Sinks:   []string{"test-sink"},
		},
	}
	promptsetDepMap, _ := json.Marshal(promptsetDep)
	var promptsetDepInterface map[string]interface{}
	_ = json.Unmarshal(promptsetDepMap, &promptsetDepInterface)

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
		locks: map[string]*packagelockfile.DependencyLockConfig{
			"test-registry/ruleset1@1.0.0":   {Integrity: "hash1"},
			"test-registry/promptset1@2.0.0": {Integrity: "hash2"},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"ruleset1": {
				{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 2, Patch: 0, Version: "1.2.0", IsSemver: true}},
			},
			"promptset1": {
				{Name: "promptset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
				{Name: "promptset1", Version: core.Version{Major: 2, Minor: 1, Patch: 0, Version: "2.1.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"ruleset1@1.2.0": {
				Metadata:  core.PackageMetadata{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 2, Patch: 0, Version: "1.2.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "hash1-new",
			},
			"promptset1@2.1.0": {
				Metadata:  core.PackageMetadata{Name: "promptset1", Version: core.Version{Major: 2, Minor: 1, Patch: 0, Version: "2.1.0", IsSemver: true}},
				Files:     []*core.File{{Path: "prompt.yml", Content: []byte("test")}},
				Integrity: "hash2-new",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpdateAll(ctx)
	if err != nil {
		t.Fatalf("UpdateAll() error = %v", err)
	}

	// Verify lockfile was updated with new versions
	if lock, exists := lockfileMgr.locks["test-registry/ruleset1@1.2.0"]; !exists || lock.Integrity != "hash1-new" {
		t.Error("Expected lockfile entry for ruleset1@1.2.0")
	}

	if lock, exists := lockfileMgr.locks["test-registry/promptset1@2.1.0"]; !exists || lock.Integrity != "hash2-new" {
		t.Error("Expected lockfile entry for promptset1@2.1.0")
	}

	// Old locks should be removed
	if _, exists := lockfileMgr.locks["test-registry/ruleset1@1.0.0"]; exists {
		t.Error("Old lockfile entry for ruleset1@1.0.0 should be removed")
	}
}

// Normal: no updates available (already at latest)
func TestUpdateAll_AlreadyUpToDate(t *testing.T) {
	ctx := context.Background()

	rulesetDep := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypeRuleset,
			Version: "^1.0.0",
			Sinks:   []string{"test-sink"},
		},
		Priority: 100,
	}
	rulesetDepMap, _ := json.Marshal(rulesetDep)
	var rulesetDepInterface map[string]interface{}
	_ = json.Unmarshal(rulesetDepMap, &rulesetDepInterface)

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Registries: map[string]map[string]interface{}{
				"test-registry": {"type": "git", "url": "https://github.com/test/repo"},
			},
			Sinks: map[string]manifest.SinkConfig{
				"test-sink": {Directory: "/tmp/test", Tool: compiler.Cursor},
			},
			Dependencies: map[string]map[string]interface{}{
				"test-registry/ruleset1": rulesetDepInterface,
			},
		},
	}

	lockfileMgr := &mockLockfileManager{
		locks: map[string]*packagelockfile.DependencyLockConfig{
			"test-registry/ruleset1@1.0.0": {Integrity: "hash1"},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"ruleset1": {
				{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpdateAll(ctx)
	if err != nil {
		t.Fatalf("UpdateAll() error = %v", err)
	}

	// Verify lockfile unchanged
	if lock, exists := lockfileMgr.locks["test-registry/ruleset1@1.0.0"]; !exists || lock.Integrity != "hash1" {
		t.Error("Lockfile should remain unchanged when already up to date")
	}
}

// Edge: empty dependencies
func TestUpdateAll_EmptyDependencies(t *testing.T) {
	ctx := context.Background()

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Registries:   map[string]map[string]interface{}{},
			Sinks:        map[string]manifest.SinkConfig{},
			Dependencies: map[string]map[string]interface{}{},
		},
	}

	lockfileMgr := &mockLockfileManager{
		locks: make(map[string]*packagelockfile.DependencyLockConfig),
	}

	service := NewArmService(manifestMgr, lockfileMgr, nil)

	err := service.UpdateAll(ctx)
	if err != nil {
		t.Fatalf("UpdateAll() error = %v", err)
	}
}

// Edge: missing lockfile entry
func TestUpdateAll_MissingLockfileEntry(t *testing.T) {
	ctx := context.Background()

	rulesetDep := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypeRuleset,
			Version: "^1.0.0",
			Sinks:   []string{"test-sink"},
		},
		Priority: 100,
	}
	rulesetDepMap, _ := json.Marshal(rulesetDep)
	var rulesetDepInterface map[string]interface{}
	_ = json.Unmarshal(rulesetDepMap, &rulesetDepInterface)

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Registries: map[string]map[string]interface{}{
				"test-registry": {"type": "git", "url": "https://github.com/test/repo"},
			},
			Sinks: map[string]manifest.SinkConfig{
				"test-sink": {Directory: "/tmp/test", Tool: compiler.Cursor},
			},
			Dependencies: map[string]map[string]interface{}{
				"test-registry/ruleset1": rulesetDepInterface,
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
		},
		packages: map[string]*core.Package{
			"ruleset1@1.0.0": {
				Metadata:  core.PackageMetadata{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "hash1",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpdateAll(ctx)
	if err != nil {
		t.Fatalf("UpdateAll() error = %v", err)
	}

	// Should install the package
	if _, exists := lockfileMgr.locks["test-registry/ruleset1@1.0.0"]; !exists {
		t.Error("Expected lockfile entry to be created")
	}
}

// Edge: mixed rulesets and promptsets
func TestUpdateAll_MixedTypes(t *testing.T) {
	ctx := context.Background()

	rulesetDep := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypeRuleset,
			Version: "^1.0.0",
			Sinks:   []string{"test-sink"},
		},
		Priority: 150,
	}
	rulesetDepMap, _ := json.Marshal(rulesetDep)
	var rulesetDepInterface map[string]interface{}
	_ = json.Unmarshal(rulesetDepMap, &rulesetDepInterface)

	promptsetDep := manifest.PromptsetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypePromptset,
			Version: "2.0.0",
			Sinks:   []string{"test-sink"},
		},
	}
	promptsetDepMap, _ := json.Marshal(promptsetDep)
	var promptsetDepInterface map[string]interface{}
	_ = json.Unmarshal(promptsetDepMap, &promptsetDepInterface)

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Registries: map[string]map[string]interface{}{
				"test-registry": {"type": "git", "url": "https://github.com/test/repo"},
			},
			Sinks: map[string]manifest.SinkConfig{
				"test-sink": {Directory: "/tmp/test", Tool: compiler.Cursor},
			},
			Dependencies: map[string]map[string]interface{}{
				"test-registry/ruleset1":   rulesetDepInterface,
				"test-registry/promptset1": promptsetDepInterface,
			},
		},
	}

	lockfileMgr := &mockLockfileManager{
		locks: map[string]*packagelockfile.DependencyLockConfig{
			"test-registry/ruleset1@1.0.0":   {Integrity: "hash1"},
			"test-registry/promptset1@2.0.0": {Integrity: "hash2"},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"ruleset1": {
				{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 1, Patch: 0, Version: "1.1.0", IsSemver: true}},
			},
			"promptset1": {
				{Name: "promptset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"ruleset1@1.1.0": {
				Metadata:  core.PackageMetadata{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 1, Patch: 0, Version: "1.1.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "hash1-new",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpdateAll(ctx)
	if err != nil {
		t.Fatalf("UpdateAll() error = %v", err)
	}

	// Ruleset should be updated
	if _, exists := lockfileMgr.locks["test-registry/ruleset1@1.1.0"]; !exists {
		t.Error("Expected updated ruleset lockfile entry")
	}

	// Promptset should remain unchanged
	if _, exists := lockfileMgr.locks["test-registry/promptset1@2.0.0"]; !exists {
		t.Error("Expected promptset lockfile entry to remain")
	}
}

// Extreme: many dependencies
func TestUpdateAll_ManyDependencies(t *testing.T) {
	ctx := context.Background()

	deps := make(map[string]map[string]interface{})
	locks := make(map[string]*packagelockfile.DependencyLockConfig)
	versions := make(map[string][]core.PackageMetadata)
	packages := make(map[string]*core.Package)

	for i := 0; i < 50; i++ {
		pkgName := "pkg" + string(rune('0'+i%10)) + string(rune('0'+i/10))
		key := manifest.DependencyKey("registry", pkgName)

		dep := manifest.RulesetDependencyConfig{
			BaseDependencyConfig: manifest.BaseDependencyConfig{
				Type:    manifest.ResourceTypeRuleset,
				Version: "^1.0.0",
				Sinks:   []string{"test-sink"},
			},
			Priority: 100,
		}
		depMap, _ := json.Marshal(dep)
		var depInterface map[string]interface{}
		_ = json.Unmarshal(depMap, &depInterface)
		deps[key] = depInterface

		locks[key+"@1.0.0"] = &packagelockfile.DependencyLockConfig{Integrity: "hash-old"}

		versions[pkgName] = []core.PackageMetadata{
			{Name: pkgName, Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			{Name: pkgName, Version: core.Version{Major: 1, Minor: 1, Patch: 0, Version: "1.1.0", IsSemver: true}},
		}

		packages[pkgName+"@1.1.0"] = &core.Package{
			Metadata:  core.PackageMetadata{Name: pkgName, Version: core.Version{Major: 1, Minor: 1, Patch: 0, Version: "1.1.0", IsSemver: true}},
			Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
			Integrity: "hash-new",
		}
	}

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Registries: map[string]map[string]interface{}{
				"registry": {"type": "git", "url": "https://github.com/test/repo"},
			},
			Sinks: map[string]manifest.SinkConfig{
				"test-sink": {Directory: "/tmp/test", Tool: compiler.Cursor},
			},
			Dependencies: deps,
		},
	}

	lockfileMgr := &mockLockfileManager{locks: locks}
	mockReg := &mockRegistry{versions: versions, packages: packages}
	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpdateAll(ctx)
	if err != nil {
		t.Fatalf("UpdateAll() error = %v", err)
	}

	// All should be updated to 1.1.0
	for i := 0; i < 50; i++ {
		pkgName := "pkg" + string(rune('0'+i%10)) + string(rune('0'+i/10))
		key := "registry/" + pkgName + "@1.1.0"
		if _, exists := lockfileMgr.locks[key]; !exists {
			t.Errorf("Expected lockfile entry for %s", key)
		}
	}
}
