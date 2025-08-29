package cache

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

func TestFileCache_Set_Get(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	ctx := context.Background()

	registryKey := "sha256_registry_key"
	rulesetKey := "sha256_ruleset_key"
	version := "1111111"

	files := []arm.File{
		{
			Path:    "rules/amazonq/grug-brained-dev.md",
			Content: []byte("# Grug Brain Dev Rules"),
			Size:    22,
		},
		{
			Path:    "rules/cursor/grug-brained-dev.mdc",
			Content: []byte("// Grug Brain Dev Rules"),
			Size:    23,
		},
	}

	// Test Set
	err := cache.Set(ctx, registryKey, rulesetKey, version, files)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Verify directory structure
	expectedDir := filepath.Join(tempDir, "registries", registryKey, "rulesets", rulesetKey, version)
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %s was not created", expectedDir)
	}

	// Test Get
	gotFiles, err := cache.Get(ctx, registryKey, rulesetKey, version)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if len(gotFiles) != len(files) {
		t.Errorf("Get() returned %d files, want %d", len(gotFiles), len(files))
	}

	// Sort for comparison
	sort.Slice(gotFiles, func(i, j int) bool { return gotFiles[i].Path < gotFiles[j].Path })
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })

	for i, got := range gotFiles {
		want := files[i]
		if got.Path != want.Path {
			t.Errorf("File[%d].Path = %s, want %s", i, got.Path, want.Path)
		}
		if !reflect.DeepEqual(got.Content, want.Content) {
			t.Errorf("File[%d].Content = %s, want %s", i, got.Content, want.Content)
		}
		if got.Size != want.Size {
			t.Errorf("File[%d].Size = %d, want %d", i, got.Size, want.Size)
		}
	}
}

func TestFileCache_ListVersions(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	ctx := context.Background()

	registryKey := "sha256_registry_key"
	rulesetKey := "sha256_ruleset_key"

	// Initially empty
	versions, err := cache.ListVersions(ctx, registryKey, rulesetKey)
	if err != nil {
		t.Fatalf("ListVersions() error = %v", err)
	}
	if len(versions) != 0 {
		t.Errorf("ListVersions() = %v, want empty slice", versions)
	}

	// Add versions from PRD timeline
	testVersions := []string{"1111111", "2222222", "3333333", "6666666"}
	files := []arm.File{{Path: "test.md", Content: []byte("test"), Size: 4}}

	for _, version := range testVersions {
		err := cache.Set(ctx, registryKey, rulesetKey, version, files)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	// List versions
	gotVersions, err := cache.ListVersions(ctx, registryKey, rulesetKey)
	if err != nil {
		t.Fatalf("ListVersions() error = %v", err)
	}

	sort.Strings(gotVersions)
	sort.Strings(testVersions)

	if !reflect.DeepEqual(gotVersions, testVersions) {
		t.Errorf("ListVersions() = %v, want %v", gotVersions, testVersions)
	}
}

func TestFileCache_Get_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	ctx := context.Background()

	_, err := cache.Get(ctx, "nonexistent", "ruleset", "version")
	if err == nil {
		t.Error("Get() expected error for nonexistent version")
	}
}

func TestFileCache_InvalidateVersion(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	ctx := context.Background()

	registryKey := "sha256_registry_key"
	rulesetKey := "sha256_ruleset_key"
	version := "1111111"

	files := []arm.File{{Path: "test.md", Content: []byte("test"), Size: 4}}

	// Set up cache
	err := cache.Set(ctx, registryKey, rulesetKey, version, files)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Verify it exists
	_, err = cache.Get(ctx, registryKey, rulesetKey, version)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	// Invalidate version
	err = cache.InvalidateVersion(ctx, registryKey, rulesetKey, version)
	if err != nil {
		t.Fatalf("InvalidateVersion() error = %v", err)
	}

	// Verify it's gone
	_, err = cache.Get(ctx, registryKey, rulesetKey, version)
	if err == nil {
		t.Error("Get() should fail after InvalidateVersion()")
	}
}

func TestFileCache_InvalidateRuleset(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	ctx := context.Background()

	registryKey := "sha256_registry_key"
	rulesetKey := "sha256_ruleset_key"
	files := []arm.File{{Path: "test.md", Content: []byte("test"), Size: 4}}

	// Set up multiple versions
	versions := []string{"1111111", "2222222"}
	for _, version := range versions {
		err := cache.Set(ctx, registryKey, rulesetKey, version, files)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	// Verify they exist
	gotVersions, err := cache.ListVersions(ctx, registryKey, rulesetKey)
	if err != nil {
		t.Fatalf("ListVersions() error = %v", err)
	}
	if len(gotVersions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(gotVersions))
	}

	// Invalidate ruleset
	err = cache.InvalidateRuleset(ctx, registryKey, rulesetKey)
	if err != nil {
		t.Fatalf("InvalidateRuleset() error = %v", err)
	}

	// Verify all versions are gone
	gotVersions, err = cache.ListVersions(ctx, registryKey, rulesetKey)
	if err != nil {
		t.Fatalf("ListVersions() error = %v", err)
	}
	if len(gotVersions) != 0 {
		t.Errorf("Expected 0 versions after invalidation, got %d", len(gotVersions))
	}
}

func TestFileCache_InvalidateRegistry(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	ctx := context.Background()

	registryKey := "sha256_registry_key"
	files := []arm.File{{Path: "test.md", Content: []byte("test"), Size: 4}}

	// Set up multiple rulesets
	rulesets := []string{"ruleset1", "ruleset2"}
	for _, rulesetKey := range rulesets {
		err := cache.Set(ctx, registryKey, rulesetKey, "1111111", files)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	// Verify they exist
	for _, rulesetKey := range rulesets {
		versions, err := cache.ListVersions(ctx, registryKey, rulesetKey)
		if err != nil {
			t.Fatalf("ListVersions() error = %v", err)
		}
		if len(versions) != 1 {
			t.Errorf("Expected 1 version for %s, got %d", rulesetKey, len(versions))
		}
	}

	// Invalidate registry
	err := cache.InvalidateRegistry(ctx, registryKey)
	if err != nil {
		t.Fatalf("InvalidateRegistry() error = %v", err)
	}

	// Verify all rulesets are gone
	for _, rulesetKey := range rulesets {
		versions, err := cache.ListVersions(ctx, registryKey, rulesetKey)
		if err != nil {
			t.Fatalf("ListVersions() error = %v", err)
		}
		if len(versions) != 0 {
			t.Errorf("Expected 0 versions for %s after registry invalidation, got %d", rulesetKey, len(versions))
		}
	}
}

func TestFileCache_PRDScenario(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	ctx := context.Background()

	// PRD registry key: sha256("https://github.com/my-user/ai-rules" + "git")
	registryKey := "sha256_github_my_user_ai_rules_git"

	// PRD ruleset keys
	amazonqRulesetKey := "sha256_rules_amazonq_md"

	// PRD timeline versions
	v100Files := []arm.File{
		{Path: "rules/amazonq/grug-brained-dev.md", Content: []byte("# Grug Brain Dev"), Size: 15},
	}
	v110Files := []arm.File{
		{Path: "rules/amazonq/grug-brained-dev.md", Content: []byte("# Grug Brain Dev"), Size: 15},
		{Path: "rules/amazonq/generate-tasks.md", Content: []byte("# Generate Tasks"), Size: 16},
		{Path: "rules/amazonq/process-tasks.md", Content: []byte("# Process Tasks"), Size: 15},
	}

	// Test v1.0.0 (commit 1111111)
	err := cache.Set(ctx, registryKey, amazonqRulesetKey, "1111111", v100Files)
	if err != nil {
		t.Fatalf("Set v1.0.0 error = %v", err)
	}

	// Test v1.1.0 (commit 3333333)
	err = cache.Set(ctx, registryKey, amazonqRulesetKey, "3333333", v110Files)
	if err != nil {
		t.Fatalf("Set v1.1.0 error = %v", err)
	}

	// Verify versions
	versions, err := cache.ListVersions(ctx, registryKey, amazonqRulesetKey)
	if err != nil {
		t.Fatalf("ListVersions() error = %v", err)
	}

	expectedVersions := []string{"1111111", "3333333"}
	sort.Strings(versions)
	sort.Strings(expectedVersions)

	if !reflect.DeepEqual(versions, expectedVersions) {
		t.Errorf("ListVersions() = %v, want %v", versions, expectedVersions)
	}

	// Verify v1.1.0 content
	files, err := cache.Get(ctx, registryKey, amazonqRulesetKey, "3333333")
	if err != nil {
		t.Fatalf("Get v1.1.0 error = %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files in v1.1.0, got %d", len(files))
	}

	// Verify cache directory structure matches PRD
	expectedPath := filepath.Join(tempDir, "registries", registryKey, "rulesets", amazonqRulesetKey, "3333333", "rules", "amazonq")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected PRD cache structure not found: %s", expectedPath)
	}
}

func TestFileCache_NewFileCache(t *testing.T) {
	cache := NewFileCache()
	if cache == nil {
		t.Error("NewFileCache() returned nil")
	}
}

func TestFileCache_NewFileCacheWithDir(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewFileCacheWithDir(tempDir)
	if cache == nil {
		t.Error("NewFileCacheWithDir() returned nil")
	}
}
