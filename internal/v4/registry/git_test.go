package registry

import (
	"context"
	"os"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/storage"
)

// Test cases based on git-registry.md documentation

// Basic Operations
func TestGitRegistry_ListPackageVersions(t *testing.T) {
	// Test finding tags and branches
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Tag("v2.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(versions))
	}
}

func TestGitRegistry_GetPackage(t *testing.T) {
	// Test retrieving files from specific version
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	pkg, err := registry.GetPackage(ctx, "test-package", versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if pkg.Metadata.RegistryName != "test-registry" {
		t.Errorf("expected registry name 'test-registry', got %s", pkg.Metadata.RegistryName)
	}

	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(pkg.Files))
	}

	if pkg.Files[0].Path != "test.yml" {
		t.Errorf("expected file path 'test.yml', got %s", pkg.Files[0].Path)
	}
}

// Version Resolution
func TestGitRegistry_SemanticVersionTags(t *testing.T) {
	// Test v1.0.0, v2.1.0, test sorting
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v2.1.0").
		Tag("v1.0.0").
		Tag("v1.2.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	if len(versions) != 3 {
		t.Errorf("expected 3 versions, got %d", len(versions))
	}

	// Check versions are parsed correctly
	for _, v := range versions {
		if v.Version == "" {
			t.Error("version should not be empty")
		}
	}
}

func TestGitRegistry_VersionWithoutVPrefix(t *testing.T) {
	// Test 1.0.0, 2.1.0 (no v prefix)
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("1.0.0").
		Tag("2.1.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(versions))
	}

	// Check versions without v prefix are parsed
	for _, v := range versions {
		if v.Version != "1.0.0" && v.Version != "2.1.0" {
			t.Errorf("unexpected version: %s", v.Version)
		}
	}
}

func TestGitRegistry_PartialVersions(t *testing.T) {
	// Test 1.0, 2 (missing patch/minor)
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("1.0").
		Tag("2").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(versions))
	}
}

func TestGitRegistry_MixedVersionFormats(t *testing.T) {
	// Test v1.0.0, 2.1.0, 3.0, v4 mixed together
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Tag("2.1.0").
		Tag("3.0").
		Tag("v4").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	if len(versions) != 4 {
		t.Errorf("expected 4 versions, got %d", len(versions))
	}
}

func TestGitRegistry_NonSemanticTags(t *testing.T) {
	// Test that non-semantic tags are ignored
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("test-tag").
		Tag("release-candidate").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Non-semantic tags should be ignored
	if len(versions) != 0 {
		t.Errorf("expected 0 versions (non-semantic tags ignored), got %d", len(versions))
	}
}

func TestGitRegistry_BranchSupport(t *testing.T) {
	// Test main, develop branches (if configured)
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Branch("develop").
		AddFile("dev.yml", "dev content").
		Commit("Dev commit").
		Checkout("main").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
		Branches: []string{"main", "develop"},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Should have 2 branches as versions (main + develop)
	if len(versions) < 2 {
		t.Errorf("expected at least 2 versions, got %d", len(versions))
	}

	// Verify branch names are found
	foundMain := false
	foundDevelop := false
	for _, v := range versions {
		if v.Version == "main" {
			foundMain = true
		}
		if v.Version == "develop" {
			foundDevelop = true
		}
	}

	if !foundMain {
		t.Error("expected to find 'main' branch in versions")
	}
	if !foundDevelop {
		t.Error("expected to find 'develop' branch in versions")
	}
}

func TestGitRegistry_BranchNotFound(t *testing.T) {
	// Test configured branch that doesn't exist
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
		Branches: []string{"main", "nonexistent"},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Should have at least 1 version (main), not nonexistent
	if len(versions) < 1 {
		t.Errorf("expected at least 1 version, got %d", len(versions))
		return
	}

	// Verify main branch found but not nonexistent
	foundMain := false
	foundNonexistent := false
	for _, v := range versions {
		if v.Version == "main" {
			foundMain = true
		}
		if v.Version == "nonexistent" {
			foundNonexistent = true
		}
	}

	if !foundMain {
		t.Error("expected to find 'main' branch in versions")
	}
	if foundNonexistent {
		t.Error("should not find 'nonexistent' branch in versions")
	}
}

func TestGitRegistry_VersionPriority(t *testing.T) {
	// Test semantic tags > branches
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Tag("stable").  // Non-semantic, should be ignored
		Tag("v2.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
		Branches: []string{"main"},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Should have 2 semantic tags + 1 branch = 3 versions
	// Note: "stable" non-semantic tag should be filtered out
	if len(versions) < 3 {
		t.Errorf("expected at least 3 versions (2 semantic tags + 1 branch), got %d", len(versions))
		for i, v := range versions {
			t.Logf("  version[%d]: %s (Major:%d Minor:%d Patch:%d)", i, v.Version, v.Major, v.Minor, v.Patch)
		}
		return
	}

	// Verify we have the semantic versions
	foundV1 := false
	foundV2 := false
	for _, v := range versions {
		if v.Version == "v1.0.0" {
			foundV1 = true
		}
		if v.Version == "v2.0.0" {
			foundV2 = true
		}
	}

	if !foundV1 || !foundV2 {
		t.Error("expected to find v1.0.0 and v2.0.0")
	}

	// Verify non-semantic tag excluded
	for _, v := range versions {
		if v.Version == "stable" {
			t.Error("non-semantic tag 'stable' should be excluded")
		}
	}
}

// Cache Key Normalization
func TestNormalizePatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name:     "single pattern",
			input:    []string{"*.yml"},
			expected: []string{"*.yml"},
		},
		{
			name:     "sort patterns",
			input:    []string{"*.yaml", "*.yml", "*.json"},
			expected: []string{"*.json", "*.yaml", "*.yml"},
		},
		{
			name:     "normalize backslashes",
			input:    []string{"dir\\*.yml", "test\\**"},
			expected: []string{"dir/*.yml", "test/**"},
		},
		{
			name:     "trim whitespace",
			input:    []string{" *.yml ", "\t*.yaml\n"},
			expected: []string{"*.yaml", "*.yml"},
		},
		{
			name:     "mixed separators and whitespace",
			input:    []string{" build\\** ", "test/**", " *.yml"},
			expected: []string{"*.yml", "build/**", "test/**"},
		},
		{
			name:     "duplicates preserved",
			input:    []string{"*.yml", "*.yml"},
			expected: []string{"*.yml", "*.yml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePatterns(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("at index %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestGitRegistry_CacheKeyNormalization(t *testing.T) {
	// Test that different pattern orders produce same cache key
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Get package with patterns in different orders
	pkg1, err := registry.GetPackage(ctx, "test-package", versions[0], []string{"*.yml", "*.yaml"}, []string{"test/**", "build/**"})
	if err != nil {
		t.Fatalf("failed to get package 1: %v", err)
	}

	pkg2, err := registry.GetPackage(ctx, "test-package", versions[0], []string{"*.yaml", "*.yml"}, []string{"build/**", "test/**"})
	if err != nil {
		t.Fatalf("failed to get package 2: %v", err)
	}

	// Both should return same result (cache hit on second call)
	if len(pkg1.Files) != len(pkg2.Files) {
		t.Errorf("expected same number of files, got %d vs %d", len(pkg1.Files), len(pkg2.Files))
	}
}

// File Filtering
func TestGitRegistry_IncludePatterns(t *testing.T) {
	// Test --include "*.yml"
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("doc.md", "md content").
		AddFile("config.json", "json content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Test include pattern
	pkg, err := registry.GetPackage(ctx, "test-package", versions[0], []string{"*.yml"}, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(pkg.Files))
	}

	if pkg.Files[0].Path != "rule.yml" {
		t.Errorf("expected file 'rule.yml', got %s", pkg.Files[0].Path)
	}
}

func TestGitRegistry_ExcludePatterns(t *testing.T) {
	// Test --exclude "build/**"
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("build/output.js", "build content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Test exclude pattern
	pkg, err := registry.GetPackage(ctx, "test-package", versions[0], nil, []string{"build/*"})
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(pkg.Files))
	}

	if pkg.Files[0].Path != "rule.yml" {
		t.Errorf("expected file 'rule.yml', got %s", pkg.Files[0].Path)
	}
}

func TestGitRegistry_CombinedPatterns(t *testing.T) {
	// Test include + exclude together
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("test.yml", "test yml content").
		AddFile("doc.md", "md content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Test include *.yml but exclude test*
	pkg, err := registry.GetPackage(ctx, "test-package", versions[0], []string{"*.yml"}, []string{"test*"})
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(pkg.Files))
	}

	if pkg.Files[0].Path != "rule.yml" {
		t.Errorf("expected file 'rule.yml', got %s", pkg.Files[0].Path)
	}
}

func TestGitRegistry_NoPatterns(t *testing.T) {
	// Test return all files
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("doc.md", "md content").
		AddFile("config.json", "json content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Test no patterns - should return all files
	pkg, err := registry.GetPackage(ctx, "test-package", versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(pkg.Files))
	}
}

// Repository Structure
func TestGitRegistry_MultipleFileTypes(t *testing.T) {
	// Test .yml, .md, .json files
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("readme.md", "md content").
		AddFile("config.json", "json content").
		AddFile("script.js", "js content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	pkg, err := registry.GetPackage(ctx, "test-package", versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 4 {
		t.Errorf("expected 4 files, got %d", len(pkg.Files))
	}

	// Verify all file types present
	fileTypes := make(map[string]bool)
	for _, file := range pkg.Files {
		fileTypes[file.Path] = true
	}

	expected := []string{"rule.yml", "readme.md", "config.json", "script.js"}
	for _, exp := range expected {
		if !fileTypes[exp] {
			t.Errorf("expected file %s not found", exp)
		}
	}
}

func TestGitRegistry_NestedDirectories(t *testing.T) {
	// Test build/cursor/, rules/security/
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("build/cursor/rule.mdc", "cursor content").
		AddFile("rules/security/auth.yml", "security content").
		AddFile("root.yml", "root content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	pkg, err := registry.GetPackage(ctx, "test-package", versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(pkg.Files))
	}

	// Verify nested paths preserved
	paths := make(map[string]bool)
	for _, file := range pkg.Files {
		paths[file.Path] = true
	}

	expected := []string{"build/cursor/rule.mdc", "rules/security/auth.yml", "root.yml"}
	for _, exp := range expected {
		if !paths[exp] {
			t.Errorf("expected path %s not found", exp)
		}
	}
}

func TestGitRegistry_ArchiveSupport(t *testing.T) {
	// Test .zip and .tar.gz files (if implemented)
	t.Skip("TODO: implement")
}

// Edge Cases
func TestGitRegistry_EmptyRepository(t *testing.T) {
	// Test no files, should not crash
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// Empty repo should have no versions
	if len(versions) != 0 {
		t.Errorf("expected 0 versions, got %d", len(versions))
	}
}

func TestGitRegistry_NoTags(t *testing.T) {
	// Test only commits, handle gracefully
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()
	versions, err := registry.ListPackageVersions(ctx, "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	// No tags should result in empty versions
	if len(versions) != 0 {
		t.Errorf("expected 0 versions, got %d", len(versions))
	}
}

func TestGitRegistry_NonExistentVersion(t *testing.T) {
	// Test request v99.0.0 that doesn't exist
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testRepo := storage.NewTestRepo(t, tempDir)
	testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Build()

	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
	}

	registry, err := NewGitRegistry("test-registry", config)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	ctx := context.Background()

	// Try to get non-existent version
	nonExistentVersion, _ := core.ParseVersion("v99.0.0")
	_, err = registry.GetPackage(ctx, "test-package", nonExistentVersion, nil, nil)
	if err == nil {
		t.Error("expected error for non-existent version, got nil")
	}
}
