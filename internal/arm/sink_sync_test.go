package arm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
)

func TestSinkLifecycleSync(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(oldWd) }()

	// Create ARM service
	armService := NewArmService()
	ctx := context.Background()

	// Setup test registry and manifest
	manifestManager := manifest.NewFileManager()
	gitConfig := registry.GitRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			URL:  "https://github.com/test/repo",
			Type: "git",
		},
		Branches: []string{"main"},
	}
	err := manifestManager.AddGitRegistry(ctx, "test-registry", gitConfig, false)
	if err != nil {
		t.Fatalf("Failed to add registry: %v", err)
	}

	// Add test ruleset to manifest
	err = manifestManager.CreateEntry(ctx, "test-registry", "test-ruleset", manifest.Entry{
		Version: "latest",
		Include: []string{"*.md"},
		Exclude: []string{},
	})
	if err != nil {
		t.Fatalf("Failed to create manifest entry: %v", err)
	}

	// Create lockfile entry to simulate installed ruleset
	lockFileManager := lockfile.NewFileManager()
	err = lockFileManager.CreateEntry(ctx, "test-registry", "test-ruleset", &lockfile.Entry{
		Version:  "1.0.0",
		Display:  "1.0.0",
		Checksum: "test-checksum",
	})
	if err != nil {
		t.Fatalf("Failed to create lockfile entry: %v", err)
	}

	// Create test directories
	testDir1 := ".test1/rules"
	testDir2 := ".test2/rules"
	err = os.MkdirAll(testDir1, 0o755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	err = os.MkdirAll(testDir2, 0o755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Simulate installed files in test directory 1
	installedDir := filepath.Join(testDir1, "arm", "test-registry", "test-ruleset", "1.0.0")
	err = os.MkdirAll(installedDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create installed directory: %v", err)
	}
	err = os.WriteFile(filepath.Join(installedDir, "test.md"), []byte("test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test SyncSink - should install ruleset to new sink
	newSink := &config.SinkConfig{
		Directories: []string{testDir2},
		Include:     []string{"test-registry/*"},
		Exclude:     []string{},
		Layout:      "hierarchical",
	}

	// Mock the InstallRuleset method by creating the expected directory structure
	// Since we can't actually download from the registry in this test
	expectedDir := filepath.Join(testDir2, "arm", "test-registry", "test-ruleset", "1.0.0")
	err = os.MkdirAll(expectedDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create expected directory: %v", err)
	}

	// Note: In a real scenario, SyncSink would call InstallRuleset which would fail
	// without a real registry, so we're mainly testing SyncRemovedSink here
	_ = newSink // Use the variable to avoid compiler error

	// Test SyncRemovedSink - should remove files from removed sink
	removedSink := &config.SinkConfig{
		Directories: []string{testDir1},
		Include:     []string{"test-registry/*"},
		Exclude:     []string{},
		Layout:      "hierarchical",
	}

	err = armService.SyncRemovedSink(ctx, removedSink)
	if err != nil {
		t.Fatalf("SyncRemovedSink failed: %v", err)
	}

	// Verify that files were removed from the removed sink directory
	installer := installer.NewInstaller(removedSink)
	installations, err := installer.ListInstalled(ctx, testDir1)
	if err == nil && len(installations) > 0 {
		t.Errorf("Expected no installations in removed sink directory, but found %d", len(installations))
	}
}

func TestSinkMatching(t *testing.T) {
	armService := NewArmService()

	tests := []struct {
		name       string
		rulesetKey string
		sink       *config.SinkConfig
		expected   bool
	}{
		{
			name:       "matches include pattern",
			rulesetKey: "ai-rules/amazonq-rules",
			sink: &config.SinkConfig{
				Include: []string{"ai-rules/amazonq-*"},
				Exclude: []string{},
			},
			expected: true,
		},
		{
			name:       "excluded by exclude pattern",
			rulesetKey: "ai-rules/cursor-rules",
			sink: &config.SinkConfig{
				Include: []string{"ai-rules/*"},
				Exclude: []string{"ai-rules/cursor-*"},
			},
			expected: false,
		},
		{
			name:       "no include patterns matches all",
			rulesetKey: "any-registry/any-ruleset",
			sink: &config.SinkConfig{
				Include: []string{},
				Exclude: []string{},
			},
			expected: true,
		},
		{
			name:       "doesn't match include pattern",
			rulesetKey: "other-registry/rules",
			sink: &config.SinkConfig{
				Include: []string{"ai-rules/*"},
				Exclude: []string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := armService.matchesSink(tt.rulesetKey, tt.sink)
			if result != tt.expected {
				t.Errorf("matchesSink(%s, %+v) = %v, expected %v", tt.rulesetKey, tt.sink, result, tt.expected)
			}
		})
	}
}
