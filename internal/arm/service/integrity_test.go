package service

import (
	"context"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
)

func TestIntegrityVerification_Success(t *testing.T) {
	ctx := context.Background()

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
					Directory: t.TempDir(),
					Tool:      compiler.Cursor,
				},
			},
		},
	}

	// Lock file with matching integrity
	lockfileMgr := &mockLockfileManager{
		locks: map[string]*packagelockfile.DependencyLockConfig{
			"test-registry/test-ruleset@1.0.0": {
				Integrity: "sha256-correct-hash",
			},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"test-ruleset": {
				{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"test-ruleset@1.0.0": {
				Metadata:  core.PackageMetadata{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "sha256-correct-hash", // Matches lock file
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	// Should succeed because integrity matches
	err := service.InstallRuleset(ctx, "test-registry", "test-ruleset", "1.0.0", 100, nil, nil, []string{"test-sink"})
	if err != nil {
		t.Fatalf("InstallRuleset() should succeed with matching integrity, got error: %v", err)
	}
}

func TestIntegrityVerification_Failure(t *testing.T) {
	ctx := context.Background()

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
					Directory: t.TempDir(),
					Tool:      compiler.Cursor,
				},
			},
		},
	}

	// Lock file with different integrity
	lockfileMgr := &mockLockfileManager{
		locks: map[string]*packagelockfile.DependencyLockConfig{
			"test-registry/test-ruleset@1.0.0": {
				Integrity: "sha256-expected-hash",
			},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"test-ruleset": {
				{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"test-ruleset@1.0.0": {
				Metadata:  core.PackageMetadata{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("modified content")}},
				Integrity: "sha256-different-hash", // Does NOT match lock file
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	// Should fail because integrity doesn't match
	err := service.InstallRuleset(ctx, "test-registry", "test-ruleset", "1.0.0", 100, nil, nil, []string{"test-sink"})
	if err == nil {
		t.Fatal("InstallRuleset() should fail with mismatched integrity")
	}

	// Verify error message contains expected information
	errMsg := err.Error()
	if !strings.Contains(errMsg, "integrity verification failed") {
		t.Errorf("Error should mention integrity verification, got: %v", err)
	}
	if !strings.Contains(errMsg, "sha256-expected-hash") {
		t.Errorf("Error should show expected hash, got: %v", err)
	}
	if !strings.Contains(errMsg, "sha256-different-hash") {
		t.Errorf("Error should show actual hash, got: %v", err)
	}
	if !strings.Contains(errMsg, "test-registry/test-ruleset@1.0.0") {
		t.Errorf("Error should show package identifier, got: %v", err)
	}
}

func TestIntegrityVerification_NoLockFile(t *testing.T) {
	ctx := context.Background()

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
					Directory: t.TempDir(),
					Tool:      compiler.Cursor,
				},
			},
		},
	}

	// Empty lock file (no existing locks)
	lockfileMgr := &mockLockfileManager{
		locks: map[string]*packagelockfile.DependencyLockConfig{},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"test-ruleset": {
				{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"test-ruleset@1.0.0": {
				Metadata:  core.PackageMetadata{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "sha256-some-hash",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	// Should succeed because no lock exists (first install)
	err := service.InstallRuleset(ctx, "test-registry", "test-ruleset", "1.0.0", 100, nil, nil, []string{"test-sink"})
	if err != nil {
		t.Fatalf("InstallRuleset() should succeed with no lock file, got error: %v", err)
	}

	// Verify integrity was stored in lock file
	lock, err := lockfileMgr.GetDependencyLock(ctx, "test-registry", "test-ruleset", "1.0.0")
	if err != nil {
		t.Fatalf("Failed to get lock: %v", err)
	}
	if lock.Integrity != "sha256-some-hash" {
		t.Errorf("Expected integrity to be stored, got: %s", lock.Integrity)
	}
}

func TestIntegrityVerification_EmptyIntegrity(t *testing.T) {
	ctx := context.Background()

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
					Directory: t.TempDir(),
					Tool:      compiler.Cursor,
				},
			},
		},
	}

	// Lock file exists but has empty integrity (backwards compatibility)
	lockfileMgr := &mockLockfileManager{
		locks: map[string]*packagelockfile.DependencyLockConfig{
			"test-registry/test-ruleset@1.0.0": {
				Integrity: "", // Empty integrity
			},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"test-ruleset": {
				{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"test-ruleset@1.0.0": {
				Metadata:  core.PackageMetadata{Name: "test-ruleset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "rule.yml", Content: []byte("test")}},
				Integrity: "sha256-some-hash",
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	// Should succeed because lock has empty integrity (backwards compatibility)
	err := service.InstallRuleset(ctx, "test-registry", "test-ruleset", "1.0.0", 100, nil, nil, []string{"test-sink"})
	if err != nil {
		t.Fatalf("InstallRuleset() should succeed with empty integrity (backwards compatibility), got error: %v", err)
	}
}

func TestIntegrityVerification_Promptset(t *testing.T) {
	ctx := context.Background()

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
					Directory: t.TempDir(),
					Tool:      compiler.Cursor,
				},
			},
		},
	}

	// Lock file with different integrity for promptset
	lockfileMgr := &mockLockfileManager{
		locks: map[string]*packagelockfile.DependencyLockConfig{
			"test-registry/test-promptset@1.0.0": {
				Integrity: "sha256-expected-hash",
			},
		},
	}

	mockReg := &mockRegistry{
		versions: map[string][]core.PackageMetadata{
			"test-promptset": {
				{Name: "test-promptset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
			},
		},
		packages: map[string]*core.Package{
			"test-promptset@1.0.0": {
				Metadata:  core.PackageMetadata{Name: "test-promptset", Version: core.Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true}},
				Files:     []*core.File{{Path: "prompt.yml", Content: []byte("modified")}},
				Integrity: "sha256-different-hash", // Does NOT match
			},
		},
	}

	regFactory := &mockRegistryFactory{registry: mockReg}
	service := NewArmService(manifestMgr, lockfileMgr, regFactory)

	// Should fail for promptsets too
	err := service.InstallPromptset(ctx, "test-registry", "test-promptset", "1.0.0", nil, nil, []string{"test-sink"})
	if err == nil {
		t.Fatal("InstallPromptset() should fail with mismatched integrity")
	}

	if !strings.Contains(err.Error(), "integrity verification failed") {
		t.Errorf("Error should mention integrity verification, got: %v", err)
	}
}
