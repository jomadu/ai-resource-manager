package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
)

// Normal: uninstall multiple dependencies
func TestUninstallAll_Normal(t *testing.T) {
	ctx := context.Background()

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

	service := NewArmService(manifestMgr, lockfileMgr, nil)

	err := service.UninstallAll(ctx)
	if err != nil {
		t.Fatalf("UninstallAll() error = %v", err)
	}

	if len(manifestMgr.manifest.Dependencies) != 0 {
		t.Errorf("Expected 0 dependencies, got %d", len(manifestMgr.manifest.Dependencies))
	}

	if len(lockfileMgr.locks) != 0 {
		t.Errorf("Expected 0 locks, got %d", len(lockfileMgr.locks))
	}
}

// Edge: empty dependencies
func TestUninstallAll_EmptyDependencies(t *testing.T) {
	ctx := context.Background()

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Sinks:        map[string]manifest.SinkConfig{},
			Dependencies: map[string]map[string]interface{}{},
		},
	}

	lockfileMgr := &mockLockfileManager{
		locks: make(map[string]*packagelockfile.DependencyLockConfig),
	}

	service := NewArmService(manifestMgr, lockfileMgr, nil)

	err := service.UninstallAll(ctx)
	if err != nil {
		t.Fatalf("UninstallAll() error = %v", err)
	}
}

// Edge: missing sink
func TestUninstallAll_MissingSink(t *testing.T) {
	ctx := context.Background()

	rulesetDep := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Type:    manifest.ResourceTypeRuleset,
			Version: "1.0.0",
			Sinks:   []string{"nonexistent-sink"},
		},
		Priority: 100,
	}
	rulesetDepMap, _ := json.Marshal(rulesetDep)
	var rulesetDepInterface map[string]interface{}
	_ = json.Unmarshal(rulesetDepMap, &rulesetDepInterface)

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Sinks: map[string]manifest.SinkConfig{},
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

	service := NewArmService(manifestMgr, lockfileMgr, nil)

	err := service.UninstallAll(ctx)
	if err != nil {
		t.Fatalf("UninstallAll() error = %v", err)
	}

	if len(manifestMgr.manifest.Dependencies) != 0 {
		t.Errorf("Expected 0 dependencies, got %d", len(manifestMgr.manifest.Dependencies))
	}
}

// Edge: malformed dependency (missing version)
func TestUninstallAll_MalformedDependency(t *testing.T) {
	ctx := context.Background()

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Sinks: map[string]manifest.SinkConfig{
				"test-sink": {
					Directory: "/tmp/test",
					Tool:      compiler.Cursor,
				},
			},
			Dependencies: map[string]map[string]interface{}{
				"test-registry/bad": {
					"type": "ruleset",
					// missing version and sinks
				},
			},
		},
	}

	lockfileMgr := &mockLockfileManager{
		locks: make(map[string]*packagelockfile.DependencyLockConfig),
	}

	service := NewArmService(manifestMgr, lockfileMgr, nil)

	err := service.UninstallAll(ctx)
	if err != nil {
		t.Fatalf("UninstallAll() error = %v", err)
	}

	// Should still remove the malformed dependency
	if len(manifestMgr.manifest.Dependencies) != 0 {
		t.Errorf("Expected 0 dependencies, got %d", len(manifestMgr.manifest.Dependencies))
	}
}

// Extreme: large number of dependencies
func TestUninstallAll_ManyDependencies(t *testing.T) {
	ctx := context.Background()

	deps := make(map[string]map[string]interface{})
	locks := make(map[string]*packagelockfile.DependencyLockConfig)

	for i := 0; i < 100; i++ {
		pkgName := "pkg" + string(rune('0'+i%10)) + string(rune('0'+i/10))
		key := manifest.DependencyKey("registry", pkgName)
		dep := manifest.RulesetDependencyConfig{
			BaseDependencyConfig: manifest.BaseDependencyConfig{
				Type:    manifest.ResourceTypeRuleset,
				Version: "1.0.0",
				Sinks:   []string{"test-sink"},
			},
			Priority: 100,
		}
		depMap, _ := json.Marshal(dep)
		var depInterface map[string]interface{}
		_ = json.Unmarshal(depMap, &depInterface)
		deps[key] = depInterface
		locks[key+"@1.0.0"] = &packagelockfile.DependencyLockConfig{Integrity: "hash"}
	}

	manifestMgr := &mockManifestManager{
		manifest: &manifest.Manifest{
			Sinks: map[string]manifest.SinkConfig{
				"test-sink": {
					Directory: "/tmp/test",
					Tool:      compiler.Cursor,
				},
			},
			Dependencies: deps,
		},
	}

	lockfileMgr := &mockLockfileManager{locks: locks}

	service := NewArmService(manifestMgr, lockfileMgr, nil)

	err := service.UninstallAll(ctx)
	if err != nil {
		t.Fatalf("UninstallAll() error = %v", err)
	}

	if len(manifestMgr.manifest.Dependencies) != 0 {
		t.Errorf("Expected 0 dependencies, got %d", len(manifestMgr.manifest.Dependencies))
	}
}
