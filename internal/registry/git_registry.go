package registry

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/resolver"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// GitRegistry implements Git-based registry access with caching.
// It abstracts all git repository operations - users simply create the registry
// and call methods like GetTags() without worrying about cloning, fetching, etc.
type GitRegistry struct {
	cache        cache.RegistryRulesetCache
	repo         cache.GitRepoCache
	registryURL  string
	registryType string
	resolver     resolver.ConstraintResolver
}

// NewGitRegistry creates a new Git-based registry that handles all git operations internally.
// The registry will automatically clone the repository on first use and fetch updates as needed.
func NewGitRegistry(rulesetCache cache.RegistryRulesetCache, repoCache cache.GitRepoCache, registryURL, registryType string) *GitRegistry {
	return &GitRegistry{
		cache:        rulesetCache,
		repo:         repoCache,
		registryURL:  registryURL,
		registryType: registryType,
		resolver:     resolver.NewGitConstraintResolver(),
	}
}

func (g *GitRegistry) ListVersions(ctx context.Context) ([]types.VersionRef, error) {
	tags, err := g.repo.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	branches, err := g.repo.GetBranches(ctx)
	if err != nil {
		return nil, err
	}

	var versions []types.VersionRef
	for _, tag := range tags {
		versions = append(versions, types.VersionRef{ID: tag, Type: types.Tag})
	}
	for _, branch := range branches {
		versions = append(versions, types.VersionRef{ID: branch, Type: types.Branch})
	}

	return versions, nil
}

func (g *GitRegistry) GetContent(ctx context.Context, version types.VersionRef, selector types.ContentSelector) ([]types.File, error) {
	// Try cache first
	files, err := g.cache.GetRulesetVersion(ctx, selector, version.ID)
	if err == nil {
		return files, nil
	}

	// Get files from git repo
	files, err = g.repo.GetFiles(ctx, version.ID, selector)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = g.cache.SetRulesetVersion(ctx, selector, version.ID, files)

	return files, nil
}

// GetTags returns all available tags from the registry.
func (g *GitRegistry) GetTags(ctx context.Context) ([]string, error) {
	return g.repo.GetTags(ctx)
}

// GetBranches returns all available branches from the registry.
func (g *GitRegistry) GetBranches(ctx context.Context) ([]string, error) {
	return g.repo.GetBranches(ctx)
}

func (g *GitRegistry) ResolveVersion(ctx context.Context, constraint string) (*types.VersionRef, error) {
	parsedConstraint, err := g.resolver.ParseConstraint(constraint)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint %s: %w", constraint, err)
	}

	versions, err := g.ListVersions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	resolvedVersion, err := g.resolver.FindBestMatch(parsedConstraint, versions)
	if err != nil {
		return nil, fmt.Errorf("no matching version found for %s: %w", constraint, err)
	}

	return resolvedVersion, nil
}
