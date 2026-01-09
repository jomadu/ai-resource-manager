package registry

import (
	"context"
	"path/filepath"
	"sort"
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
	name         string
	config       GitRegistryConfig
	repo         storage.RepoInterface
	packageCache *storage.PackageCache
}

func NewGitRegistry(name string, config GitRegistryConfig) (*GitRegistry, error) {
	registry, err := storage.NewRegistry(config)
	if err != nil {
		return nil, err
	}
	
	return &GitRegistry{
		name:         name,
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
		// Get actual branches from repo
		actualBranches, err := g.repo.GetBranches(ctx, g.config.URL)
		if err != nil {
			return nil, err
		}
		
		// Only add configured branches that actually exist (supports glob patterns)
		for _, configBranch := range g.config.Branches {
			for _, actualBranch := range actualBranches {
				matched, err := filepath.Match(configBranch, actualBranch)
				if err != nil {
					// If pattern invalid, try exact match
					matched = (configBranch == actualBranch)
				}
				if matched {
					version, _ := core.ParseVersion(actualBranch)
					versions = append(versions, version)
				}
			}
		}
	}
	
	return versions, nil
}

// normalizePatterns sorts patterns and normalizes path separators for consistent cache keys
func normalizePatterns(patterns []string) []string {
	if len(patterns) == 0 {
		return patterns
	}
	normalized := make([]string, len(patterns))
	for i, pattern := range patterns {
		normalized[i] = strings.ReplaceAll(strings.TrimSpace(pattern), "\\", "/")
	}
	sort.Strings(normalized)
	return normalized
}

// GetPackage returns files from git repository filtered by include/exclude patterns.
// packageName is used only for response metadata, not for caching or filtering.
// Cache key is based on version + include + exclude patterns, not package name.
// This allows multiple "packages" with same patterns to share cached results.
func (g *GitRegistry) GetPackage(ctx context.Context, packageName string, version core.Version, include []string, exclude []string) (*core.Package, error) {
	// Create cache key from version and normalized patterns (not package name)
	cacheKey := struct {
		Version core.Version `json:"version"`
		Include []string     `json:"include"`
		Exclude []string     `json:"exclude"`
	}{version, normalizePatterns(include), normalizePatterns(exclude)}
	
	// Try cache first
	if files, err := g.packageCache.GetPackageVersion(ctx, cacheKey, version); err == nil {
		return &core.Package{
			Metadata: core.PackageMetadata{
				RegistryName: g.name,
				Name:         packageName,
				Version:      version,
			},
			Files: files,
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
		Metadata: core.PackageMetadata{
			RegistryName: g.name,
			Name:         packageName,
			Version:      version,
		},
		Files: filteredFiles,
	}, nil
}

// matchesPatterns checks if file path matches include/exclude patterns.
// Uses filepath.Match for glob pattern support. Invalid patterns are skipped.
func (g *GitRegistry) matchesPatterns(filePath string, include []string, exclude []string) bool {
	// If no patterns, include all files
	if len(include) == 0 && len(exclude) == 0 {
		return true
	}
	
	// Check exclude patterns first
	for _, pattern := range exclude {
		matched, err := filepath.Match(pattern, filePath)
		if err != nil {
			// Invalid pattern - skip it (could log warning)
			continue
		}
		if matched {
			return false
		}
	}
	
	// If no include patterns, file is included (not excluded)
	if len(include) == 0 {
		return true
	}
	
	// Check include patterns
	for _, pattern := range include {
		matched, err := filepath.Match(pattern, filePath)
		if err != nil {
			// Invalid pattern - skip it (could log warning)
			continue
		}
		if matched {
			return true
		}
	}
	
	return false
}