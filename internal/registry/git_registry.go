package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/archive"
	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/registry/common"
	"github.com/jomadu/ai-rules-manager/internal/resolver"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// GitRegistry implements Git-based registry access with caching.
// It abstracts all git repository operations - users simply create the registry
// and call methods like GetTags() without worrying about cloning, fetching, etc.
type GitRegistry struct {
	cache     cache.RegistryRulesetCache
	repo      cache.GitRepoCache
	config    GitRegistryConfig
	resolver  resolver.ConstraintResolver
	semver    *common.SemverHelper
	extractor *archive.Extractor
}

// NewGitRegistry creates a new Git-based registry that handles all git operations internally.
// The registry will automatically clone the repository on first use and fetch updates as needed.
func NewGitRegistry(config GitRegistryConfig, rulesetCache cache.RegistryRulesetCache, repoCache cache.GitRepoCache) *GitRegistry {
	return &GitRegistry{
		cache:     rulesetCache,
		repo:      repoCache,
		config:    config,
		resolver:  resolver.NewGitConstraintResolver(),
		semver:    common.NewSemverHelper(),
		extractor: archive.NewExtractor(),
	}
}

// isBranchConstraint checks if a constraint refers to a branch name.
// Uses permissive detection: anything that's not a semantic version or "latest" is treated as a potential branch.
func (g *GitRegistry) isBranchConstraint(constraint string) bool {
	// Exclude "latest" special keyword
	if constraint == "latest" {
		return false
	}

	// Exclude semantic version patterns (contains dots and follows semver format)
	if strings.Contains(constraint, ".") {
		// Check if it's a valid semantic version
		if g.semver.IsSemverVersion(constraint) || g.semver.IsSemverVersion("v"+constraint) {
			return false
		}
		// Also exclude patterns with dots that aren't valid semver (these are invalid)
		return false
	}

	// Exclude version constraint prefixes
	if strings.HasPrefix(constraint, "^") || strings.HasPrefix(constraint, "~") {
		return false
	}

	// Everything else is treated as a potential branch name
	return true
}

func (g *GitRegistry) ListVersions(ctx context.Context, ruleset string) ([]types.Version, error) {
	// ruleset parameter ignored - Git registries return all repository versions
	tags, err := g.repo.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	// Sort tags by semantic version (descending)
	sortedTags := g.semver.SortVersionsBySemver(tags)

	var versions []types.Version
	// Add tags first (priority ordering)
	for _, tag := range sortedTags {
		versions = append(versions, types.Version{Version: tag, Display: tag})
	}

	// Add branch commits in configuration order
	for _, branch := range g.config.Branches {
		hash, err := g.repo.GetBranchHeadCommitHash(ctx, branch)
		if err != nil {
			continue // Skip branches that don't exist
		}
		versions = append(versions, types.Version{Version: hash, Display: hash[:7]})
	}

	return versions, nil
}

func (g *GitRegistry) GetContent(ctx context.Context, ruleset string, version types.Version, selector types.ContentSelector) ([]types.File, error) {
	// ruleset parameter ignored - Git registries use selector for caching
	// Try cache first
	files, err := g.cache.GetRulesetVersion(ctx, selector, version.Version)
	if err == nil {
		return files, nil
	}

	// Get all files from git repo (no selector filtering yet)
	allSelector := types.ContentSelector{} // Empty selector matches all files
	rawFiles, err := g.repo.GetFilesFromCommit(ctx, version.Version, allSelector)
	if err != nil {
		return nil, err
	}

	// Extract and merge archives with loose files
	mergedFiles, err := g.extractor.ExtractAndMerge(rawFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to extract and merge content: %w", err)
	}

	// Apply selector patterns to merged content
	var filteredFiles []types.File
	for _, file := range mergedFiles {
		if selector.Matches(file.Path) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Cache the result
	_ = g.cache.SetRulesetVersion(ctx, selector, version.Version, filteredFiles)

	return filteredFiles, nil
}

// GetTags returns all available tags from the registry.
func (g *GitRegistry) GetTags(ctx context.Context) ([]string, error) {
	return g.repo.GetTags(ctx)
}

// GetBranches returns all available branches from the registry.
func (g *GitRegistry) GetBranches(ctx context.Context) ([]string, error) {
	return g.repo.GetBranches(ctx)
}

func (g *GitRegistry) ResolveVersion(ctx context.Context, ruleset, constraint string) (*resolver.ResolvedVersion, error) {
	// ruleset parameter ignored - Git registries resolve versions for entire repository
	// Parse constraint first
	parsedConstraint, err := g.resolver.ParseConstraint(constraint)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint %s: %w", constraint, err)
	}

	// Handle branch constraints by resolving to commit hash
	if g.isBranchConstraint(constraint) {
		hash, err := g.repo.GetBranchHeadCommitHash(ctx, constraint)
		if err != nil {
			// Provide helpful error with available branches
			branches, branchErr := g.repo.GetBranches(ctx)
			if branchErr == nil && len(branches) > 0 {
				return nil, fmt.Errorf("branch %s not found. Available branches: %v", constraint, branches)
			}
			return nil, fmt.Errorf("failed to resolve branch %s: %w", constraint, err)
		}
		return &resolver.ResolvedVersion{
			Constraint: parsedConstraint,
			Version: types.Version{
				Version: hash,
				Display: hash[:7], // Use first 7 chars for display
			},
		}, nil
	}

	// Handle semantic version constraints
	versions, err := g.ListVersions(ctx, ruleset)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	resolvedVersion, err := g.resolver.FindBestMatch(parsedConstraint, versions)
	if err != nil {
		return nil, fmt.Errorf("no matching version found for %s: %w", constraint, err)
	}

	return &resolver.ResolvedVersion{
		Constraint: parsedConstraint,
		Version:    *resolvedVersion,
	}, nil
}
