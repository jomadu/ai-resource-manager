package registry

import (
	"context"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/storage"
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
// Returns semantic version tags and branch names as versions.
// Non-semantic tags are ignored.
func (g *GitRegistry) ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error) {
	// Get git tags
	tags, err := g.repo.GetTags(ctx, g.config.URL)
	if err != nil {
		return nil, err
	}

	var versions []core.Version
	for _, tag := range tags {
		version, _ := core.ParseVersion(tag)
		// Only include semantic versions
		if version.IsSemver {
			versions = append(versions, version)
		}
	}

	// Get branches if configured
	var branchVersions []core.Version
	if len(g.config.Branches) > 0 {
		// Get actual branches from repo
		actualBranches, err := g.repo.GetBranches(ctx, g.config.URL)
		if err != nil {
			return nil, err
		}

		// Only add configured branches that actually exist (supports glob patterns)
		// Preserve config order for branch priority
		for _, configBranch := range g.config.Branches {
			for _, actualBranch := range actualBranches {
				matched, err := filepath.Match(configBranch, actualBranch)
				if err != nil {
					// If pattern invalid, try exact match
					matched = (configBranch == actualBranch)
				}
				if matched {
					version, _ := core.ParseVersion(actualBranch)
					branchVersions = append(branchVersions, version)
				}
			}
		}
	}

	// Sort semver versions descending (highest first)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Compare(&versions[j]) > 0
	})

	// Append branches in config order (already in order from loop above)
	versions = append(versions, branchVersions...)

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
func (g *GitRegistry) GetPackage(ctx context.Context, packageName string, version *core.Version, include, exclude []string) (*core.Package, error) {
	// Create cache key from version and normalized patterns (not package name)
	cacheKey := struct {
		Version core.Version `json:"version"`
		Include []string     `json:"include"`
		Exclude []string     `json:"exclude"`
	}{*version, normalizePatterns(include), normalizePatterns(exclude)}

	// Try cache first
	if files, err := g.packageCache.GetPackageVersion(ctx, cacheKey, version); err == nil {
		integrity := calculateIntegrity(files)
		return &core.Package{
			Metadata: core.PackageMetadata{
				RegistryName: g.name,
				Name:         packageName,
				Version:      *version,
			},
			Files:     files,
			Integrity: integrity,
		}, nil
	}

	// Get all files from git
	files, err := g.repo.GetFilesFromCommit(ctx, g.config.URL, version.Version)
	if err != nil {
		return nil, err
	}

	// Extract archives and merge with loose files
	extractor := core.NewExtractor()
	files, err = extractor.ExtractAndMerge(files)
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
	_ = g.packageCache.SetPackageVersion(ctx, cacheKey, version, filteredFiles)

	integrity := calculateIntegrity(filteredFiles)

	return &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: g.name,
			Name:         packageName,
			Version:      *version,
		},
		Files:     filteredFiles,
		Integrity: integrity,
	}, nil
}

// matchesPatterns checks if file path matches include/exclude patterns.
// Uses core.MatchPattern for glob pattern support with ** for recursive matching.
func (g *GitRegistry) matchesPatterns(filePath string, include, exclude []string) bool {
	// Default to YAML files if no patterns specified
	if len(include) == 0 && len(exclude) == 0 {
		include = []string{"**/*.yml", "**/*.yaml"}
	}

	// Check exclude patterns first
	for _, pattern := range exclude {
		if core.MatchPattern(pattern, filePath) {
			return false
		}
	}

	// If no include patterns, file is included (not excluded)
	if len(include) == 0 {
		return true
	}

	// Check include patterns
	for _, pattern := range include {
		if core.MatchPattern(pattern, filePath) {
			return true
		}
	}

	return false
}
