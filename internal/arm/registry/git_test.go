package registry

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"os"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/storage"
)

// Test cases based on git-registry.md documentation

// Basic Operations
func TestGitRegistry_ListPackageVersions(t *testing.T) {
	// Test finding tags and branches
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Tag("v2.0.0")
	_ = builder.Build()

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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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

	pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], nil, nil)
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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v2.1.0").
		Tag("v1.0.0").
		Tag("v1.2.0")
	_ = builder.Build()

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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("1.0.0").
		Tag("2.1.0")
	_ = builder.Build()

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
	// Test 1.0, 2 (missing patch/minor) - should be ignored
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("1.0").
		Tag("2")
	_ = builder.Build()

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

	// Partial versions should be ignored (not valid semver)
	if len(versions) != 0 {
		t.Errorf("expected 0 versions (partial versions ignored), got %d", len(versions))
	}
}

func TestGitRegistry_MixedVersionFormats(t *testing.T) {
	// Test v1.0.0, 2.1.0, 3.0, v4 mixed together - only full semver accepted
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Tag("2.1.0").
		Tag("3.0").
		Tag("v4")
	_ = builder.Build()

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

	// Only full semver tags accepted (v1.0.0 and 2.1.0), partial versions ignored
	if len(versions) != 2 {
		t.Errorf("expected 2 versions (only full semver), got %d", len(versions))
	}
}

func TestGitRegistry_NonSemanticTags(t *testing.T) {
	// Test that non-semantic tags are ignored
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("test-tag").
		Tag("release-candidate")
	_ = builder.Build()

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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Branch("develop").
		AddFile("dev.yml", "dev content").
		Commit("Dev commit").
		Checkout("main")
	_ = builder.Build()

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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit")
	_ = builder.Build()

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

func TestGitRegistry_BranchConfigOrder(t *testing.T) {
	// Test that branches are ordered by config, not alphabetically
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Branch("alpha").
		AddFile("alpha.yml", "alpha content").
		Commit("Alpha commit").
		Checkout("main").
		Branch("zebra").
		AddFile("zebra.yml", "zebra content").
		Commit("Zebra commit").
		Checkout("main")
	_ = builder.Build()

	// Configure branches in specific order: zebra, main, alpha
	// (not alphabetical)
	config := GitRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "file://" + tempDir,
			Type: "git",
		},
		Branches: []string{"zebra", "main", "alpha"},
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

	// Should have 3 branches
	if len(versions) != 3 {
		t.Errorf("expected 3 versions, got %d", len(versions))
		return
	}

	// Verify branches are in config order: zebra, main, alpha
	if versions[0].Version != "zebra" {
		t.Errorf("expected first branch to be zebra, got %s", versions[0].Version)
	}
	if versions[1].Version != "main" {
		t.Errorf("expected second branch to be main, got %s", versions[1].Version)
	}
	if versions[2].Version != "alpha" {
		t.Errorf("expected third branch to be alpha, got %s", versions[2].Version)
	}
}

func TestGitRegistry_VersionPriority(t *testing.T) {
	// Test semantic tags > branches
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0").
		Tag("stable"). // Non-semantic, should be ignored
		Tag("v2.0.0")
	_ = builder.Build()

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

	// Verify ordering: semver tags first (descending), then branches
	// Expected order: v2.0.0, v1.0.0, main
	if versions[0].Version != "v2.0.0" && versions[0].Version != "2.0.0" {
		t.Errorf("expected first version to be v2.0.0, got %s", versions[0].Version)
	}
	if versions[1].Version != "v1.0.0" && versions[1].Version != "1.0.0" {
		t.Errorf("expected second version to be v1.0.0, got %s", versions[1].Version)
	}
	if versions[2].Version != "main" {
		t.Errorf("expected third version to be main (branch), got %s", versions[2].Version)
	}

	// Verify all semver versions come before all branches
	lastSemverIndex := -1
	firstBranchIndex := len(versions)
	for i, v := range versions {
		if v.IsSemver {
			lastSemverIndex = i
		} else if i < firstBranchIndex {
			firstBranchIndex = i
		}
	}
	if lastSemverIndex >= firstBranchIndex {
		t.Errorf("semver versions should come before branches, but found semver at index %d and branch at index %d", lastSemverIndex, firstBranchIndex)
	}

	// Verify we have the semantic versions
	foundV1 := false
	foundV2 := false
	for _, v := range versions {
		if v.Version == "v1.0.0" || v.Version == "1.0.0" {
			foundV1 = true
		}
		if v.Version == "v2.0.0" || v.Version == "2.0.0" {
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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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
	pkg1, err := registry.GetPackage(ctx, "test-package", &versions[0], []string{"*.yml", "*.yaml"}, []string{"test/**", "build/**"})
	if err != nil {
		t.Fatalf("failed to get package 1: %v", err)
	}

	pkg2, err := registry.GetPackage(ctx, "test-package", &versions[0], []string{"*.yaml", "*.yml"}, []string{"build/**", "test/**"})
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
	t.Run("patterns with loose files", func(t *testing.T) {
		// Test --include "*.yml"
		tempDir, err := os.MkdirTemp("", "git-registry-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		testRepo := storage.NewTestRepo(t, tempDir)
		builder := testRepo.Builder().
			Init().
			AddFile("rule.yml", "yml content").
			AddFile("doc.md", "md content").
			AddFile("config.json", "json content").
			Commit("Initial commit").
			Tag("v1.0.0")
		_ = builder.Build()

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
		pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], []string{"*.yml"}, nil)
		if err != nil {
			t.Fatalf("failed to get package: %v", err)
		}

		if len(pkg.Files) != 1 {
			t.Errorf("expected 1 file, got %d", len(pkg.Files))
		}

		if pkg.Files[0].Path != "rule.yml" {
			t.Errorf("expected file 'rule.yml', got %s", pkg.Files[0].Path)
		}
	})
}

func TestGitRegistry_ExcludePatterns(t *testing.T) {
	// Test --exclude "build/**"
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("build/output.js", "build content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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
	pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], nil, []string{"build/*"})
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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("test.yml", "test yml content").
		AddFile("doc.md", "md content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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
	pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], []string{"*.yml"}, []string{"test*"})
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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("doc.md", "md content").
		AddFile("config.json", "json content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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

	// Test no patterns - should return only YAML files by default
	pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 YAML file, got %d", len(pkg.Files))
	}
}

// Repository Structure
func TestGitRegistry_MultipleFileTypes(t *testing.T) {
	// Test .yml, .md, .json files
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("rule.yml", "yml content").
		AddFile("readme.md", "md content").
		AddFile("config.json", "json content").
		AddFile("script.js", "js content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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

	pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	// With no patterns, should only return YAML files by default
	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 YAML file, got %d", len(pkg.Files))
	}

	// Verify only YAML file present
	if pkg.Files[0].Path != "rule.yml" {
		t.Errorf("expected file rule.yml, got %s", pkg.Files[0].Path)
	}
}

func TestGitRegistry_NestedDirectories(t *testing.T) {
	// Test build/cursor/, rules/security/
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("build/cursor/rule.mdc", "cursor content").
		AddFile("rules/security/auth.yml", "security content").
		AddFile("root.yml", "root content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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

	pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	// With no patterns, should only return YAML files by default
	if len(pkg.Files) != 2 {
		t.Errorf("expected 2 YAML files, got %d", len(pkg.Files))
	}

	// Verify nested paths preserved for YAML files only
	paths := make(map[string]bool)
	for _, file := range pkg.Files {
		paths[file.Path] = true
	}

	expected := []string{"rules/security/auth.yml", "root.yml"}
	for _, exp := range expected {
		if !paths[exp] {
			t.Errorf("expected path %s not found", exp)
		}
	}
}

func TestGitRegistry_ArchiveSupport(t *testing.T) {
	// Test .zip and .tar.gz files are extracted and merged with loose files
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test archives
	zipContent := createZipArchive(map[string][]byte{
		"from-zip/rule1.yml": []byte("zip rule 1"),
		"from-zip/rule2.yml": []byte("zip rule 2"),
	})

	tarGzContent := createTarGzArchive(map[string][]byte{
		"from-tar/rule3.yml": []byte("tar rule 3"),
		"from-tar/rule4.yml": []byte("tar rule 4"),
	})

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test-package/loose-file.yml", "loose file content").
		AddFile("test-package/archive.zip", string(zipContent)).
		AddFile("test-package/archive.tar.gz", string(tarGzContent)).
		Commit("Add files and archives").
		Tag("v1.0.0")
	_ = builder.Build()

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

	pkg, err := registry.GetPackage(ctx, "test-package", &versions[0], nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	// Should have: 1 loose file + 2 from zip + 2 from tar.gz = 5 files
	// (archive files themselves are not included after extraction)
	// Archives extract to subdirectories named after the archive
	if len(pkg.Files) != 5 {
		t.Errorf("expected 5 files (1 loose + 2 from zip + 2 from tar.gz), got %d", len(pkg.Files))
		for _, f := range pkg.Files {
			t.Logf("  - %s", f.Path)
		}
	}

	// Verify all expected files are present
	paths := make(map[string]bool)
	for _, file := range pkg.Files {
		paths[file.Path] = true
	}

	expected := []string{
		"test-package/loose-file.yml",
		"archive/from-zip/rule1.yml",      // archive.zip extracts to archive/ subdirectory
		"archive/from-zip/rule2.yml",
		"archive/from-tar/rule3.yml",      // archive.tar.gz extracts to archive/ subdirectory
		"archive/from-tar/rule4.yml",
	}

	for _, exp := range expected {
		if !paths[exp] {
			t.Errorf("expected path %s not found in extracted files", exp)
		}
	}

	// Verify archive files themselves are not in the list
	if paths["test-package/archive.zip"] {
		t.Errorf("archive.zip should not be in extracted files")
	}
	if paths["test-package/archive.tar.gz"] {
		t.Errorf("archive.tar.gz should not be in extracted files")
	}

	// Verify content of extracted files
	for _, file := range pkg.Files {
		switch file.Path {
		case "test-package/loose-file.yml":
			if string(file.Content) != "loose file content" {
				t.Errorf("loose-file.yml content mismatch: got %s", string(file.Content))
			}
		case "from-zip/rule1.yml":
			if string(file.Content) != "zip rule 1" {
				t.Errorf("from-zip/rule1.yml content mismatch: got %s", string(file.Content))
			}
		case "from-zip/rule2.yml":
			if string(file.Content) != "zip rule 2" {
				t.Errorf("from-zip/rule2.yml content mismatch: got %s", string(file.Content))
			}
		case "from-tar/rule3.yml":
			if string(file.Content) != "tar rule 3" {
				t.Errorf("from-tar/rule3.yml content mismatch: got %s", string(file.Content))
			}
		case "from-tar/rule4.yml":
			if string(file.Content) != "tar rule 4" {
				t.Errorf("from-tar/rule4.yml content mismatch: got %s", string(file.Content))
			}
		}
	}
}

// Helper to create zip archive for testing
func createZipArchive(files map[string][]byte) []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for path, content := range files {
		writer, _ := zipWriter.Create(path)
		_, _ = writer.Write(content)
	}

	_ = zipWriter.Close()
	return buf.Bytes()
}

// Helper to create tar.gz archive for testing
func createTarGzArchive(files map[string][]byte) []byte {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for path, content := range files {
		header := &tar.Header{
			Name: path,
			Mode: 0o644,
			Size: int64(len(content)),
		}
		_ = tarWriter.WriteHeader(header)
		_, _ = tarWriter.Write(content)
	}

	_ = tarWriter.Close()
	_ = gzWriter.Close()
	return buf.Bytes()
}

// Edge Cases
func TestGitRegistry_EmptyRepository(t *testing.T) {
	// Test no files, should not crash
	tempDir, err := os.MkdirTemp("", "git-registry-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init()
	_ = builder.Build()

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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit")
	_ = builder.Build()

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
	defer func() { _ = os.RemoveAll(tempDir) }()

	testRepo := storage.NewTestRepo(t, tempDir)
	builder := testRepo.Builder().
		Init().
		AddFile("test.yml", "test content").
		Commit("Initial commit").
		Tag("v1.0.0")
	_ = builder.Build()

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
	_, err = registry.GetPackage(ctx, "test-package", &nonExistentVersion, nil, nil)
	if err == nil {
		t.Error("expected error for non-existent version, got nil")
	}
}
