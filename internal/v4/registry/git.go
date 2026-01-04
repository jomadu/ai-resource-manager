package registry

import (
	"context"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/storage"
)

type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type GitRegistryConfig struct {
	RegistryConfig
	Branches []string `json:"branches,omitempty"`
}

type GitLabRegistryConfig struct {
	RegistryConfig
	ProjectID  string `json:"projectId,omitempty"`
	GroupID    string `json:"groupId,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type CloudsmithRegistryConfig struct {
	RegistryConfig
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
}

type GitRegistry struct {
	config       GitRegistryConfig
	repo         storage.RepoInterface
	packageCache *storage.PackageCache
}

func NewGitRegistry(config GitRegistryConfig) (*GitRegistry, error) {
	registry, err := storage.NewRegistry(config)
	if err != nil {
		return nil, err
	}
	
	return &GitRegistry{
		config:       config,
		repo:         storage.NewRepo(registry.GetRepoDir()),
		packageCache: storage.NewPackageCache(registry.GetPackagesDir()),
	}, nil
}

// ListPackages returns empty list for git registries.
// Git registries don't have predefined packages - users define package boundaries
// via include/exclude patterns when installing.
func (g *GitRegistry) ListPackages(ctx context.Context) ([]*core.PackageMetadata, error) {
	return []*core.PackageMetadata{}, nil
}

// ListPackageVersions returns all available versions for the git repository.
// packageName is ignored since git versions are repository-wide.
// Returns both git tags and branch names as versions.
func (g *GitRegistry) ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error) {
	// Get git tags
	tags, err := g.repo.GetTags(ctx, g.config.URL)
	if err != nil {
		return nil, err
	}
	
	var versions []core.Version
	for _, tag := range tags {
		version, _ := core.ParseVersion(tag)
		versions = append(versions, version)
	}
	
	// Get branches if configured
	if len(g.config.Branches) > 0 {
		for _, branch := range g.config.Branches {
			version, _ := core.ParseVersion(branch)
			versions = append(versions, version)
		}
	}
	
	return versions, nil
}

// GetPackage returns files from git repository filtered by include/exclude patterns.
// packageName is used only for response metadata, not for caching or filtering.
// Cache key is based on version + include + exclude patterns, not package name.
// This allows multiple "packages" with same patterns to share cached results.
func (g *GitRegistry) GetPackage(ctx context.Context, packageName string, version core.Version, include []string, exclude []string) (*core.Package, error) {
	// Create cache key from version and patterns (not package name)
	cacheKey := struct {
		Version core.Version `json:"version"`
		Include []string     `json:"include"`
		Exclude []string     `json:"exclude"`
	}{version, include, exclude}
	
	// Try cache first
	if files, err := g.packageCache.GetPackageVersion(ctx, cacheKey, version); err == nil {
		return &core.Package{
			Metadata: core.PackageMetadata{Name: packageName, Version: version},
			Files:    files,
		}, nil
	}
	
	// Get all files from git
	files, err := g.repo.GetFilesFromCommit(ctx, g.config.URL, version.Version)
	if err != nil {
		return nil, err
	}
	
	// Apply include/exclude filtering
	var filteredFiles []*core.File
	for _, file := range files {
		if g.matchesPatterns(file.Path, include, exclude) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	
	// Cache the filtered result
	g.packageCache.SetPackageVersion(ctx, cacheKey, version, filteredFiles)
	
	return &core.Package{
		Metadata: core.PackageMetadata{Name: packageName, Version: version},
		Files:    filteredFiles,
	}, nil
}

// matchesPatterns checks if file path matches include/exclude patterns.
// Simple implementation - can be enhanced with glob patterns later.
func (g *GitRegistry) matchesPatterns(filePath string, include []string, exclude []string) bool {
	// If no patterns, include all files
	if len(include) == 0 && len(exclude) == 0 {
		return true
	}
	
	// Check exclude patterns first
	for _, pattern := range exclude {
		if strings.Contains(filePath, pattern) {
			return false
		}
	}
	
	// If no include patterns, file is included (not excluded)
	if len(include) == 0 {
		return true
	}
	
	// Check include patterns
	for _, pattern := range include {
		if strings.Contains(filePath, pattern) {
			return true
		}
	}
	
	return false
}