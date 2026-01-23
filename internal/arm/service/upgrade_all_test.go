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

// Normal: upgrade to latest ignoring constraints
func TestUpgradeAll_Normal(t *testing.T) {
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
			Version: "^1.0.0",
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
			"test-registry/promptset1@1.0.0": {Integrity: "hash2"},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"ruleset1": {
				{Name: "ruleset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				{Name: "ruleset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
			},
			"promptset1": {
				{Name: "promptset1", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				{Name: "promptset1", Version: core.Version{Major: 3, Minor: 0, Patch: 0, Version: "3.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"ruleset1@2.0.0": {
				Metadata:  core.PackageMetadata{Name: "ruleset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "hash1-new",
			},
			"promptset1@3.0.0": {
				Metadata:  core.PackageMetadata{Name: "promptset1", Version: core.Version{Major: 3, Minor: 0, Patch: 0, Version: "3.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "prompt.yml", Content: []byte("test")}},
				Integrity: "hash2-new",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpgradeAll(ctx)
	if err != nil {
		t.Fatalf("UpgradeAll() error = %v", err)
	}

	// Verify lockfile updated to latest versions
	if lock, exists := lockfileMgr.locks["test-registry/ruleset1@2.0.0"]; !exists || lock.Integrity != "hash1-new" {
		t.Error("Expected lockfile entry for ruleset1@2.0.0")
	}

	if lock, exists := lockfileMgr.locks["test-registry/promptset1@3.0.0"]; !exists || lock.Integrity != "hash2-new" {
		t.Error("Expected lockfile entry for promptset1@3.0.0")
	}

	// Old locks removed
	if _, exists := lockfileMgr.locks["test-registry/ruleset1@1.0.0"]; exists {
		t.Error("Old lockfile entry should be removed")
	}

	// Verify manifest constraints updated to ^2.0.0 and ^3.0.0
	rulesetConfig, _ := manifestMgr.GetRulesetDependencyConfig(ctx, "test-registry", "ruleset1")
	if rulesetConfig.Version != "^2.0.0" {
		t.Errorf("Expected ruleset constraint ^2.0.0, got %s", rulesetConfig.Version)
	}

	promptsetConfig, _ := manifestMgr.GetPromptsetDependencyConfig(ctx, "test-registry", "promptset1")
	if promptsetConfig.Version != "^3.0.0" {
		t.Errorf("Expected promptset constraint ^3.0.0, got %s", promptsetConfig.Version)
	}
}

// Normal: already at latest
func TestUpgradeAll_AlreadyLatest(t *testing.T) {
	ctx := context.Background()

	rulesetDep := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypeRuleset,
			Version: "^2.0.0",
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
			"test-registry/ruleset1@2.0.0": {Integrity: "hash1"},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"ruleset1": {
				{Name: "ruleset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpgradeAll(ctx)
	if err != nil {
		t.Fatalf("UpgradeAll() error = %v", err)
	}

	// Lockfile unchanged
	if lock, exists := lockfileMgr.locks["test-registry/ruleset1@2.0.0"]; !exists || lock.Integrity != "hash1" {
		t.Error("Lockfile should remain unchanged")
	}

	// Constraint unchanged
	rulesetConfig, _ := manifestMgr.GetRulesetDependencyConfig(ctx, "test-registry", "ruleset1")
	if rulesetConfig.Version != "^2.0.0" {
		t.Errorf("Constraint should remain ^2.0.0, got %s", rulesetConfig.Version)
	}
}

// Edge: empty dependencies
func TestUpgradeAll_EmptyDependencies(t *testing.T) {
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

	err := service.UpgradeAll(ctx)
	if err != nil {
		t.Fatalf("UpgradeAll() error = %v", err)
	}
}

// Edge: missing lockfile entry
func TestUpgradeAll_MissingLockfileEntry(t *testing.T) {
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
				{Name: "ruleset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"ruleset1@2.0.0": {
				Metadata:  core.PackageMetadata{Name: "ruleset1", Version: core.Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "hash1",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	err := service.UpgradeAll(ctx)
	if err != nil {
		t.Fatalf("UpgradeAll() error = %v", err)
	}

	// Should install latest
	if _, exists := lockfileMgr.locks["test-registry/ruleset1@2.0.0"]; !exists {
		t.Error("Expected lockfile entry to be created")
	}

	// Constraint updated
	rulesetConfig, _ := manifestMgr.GetRulesetDependencyConfig(ctx, "test-registry", "ruleset1")
	if rulesetConfig.Version != "^2.0.0" {
		t.Errorf("Expected constraint ^2.0.0, got %s", rulesetConfig.Version)
	}
}
