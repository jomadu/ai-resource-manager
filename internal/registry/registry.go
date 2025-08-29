package registry

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/cache"
)

// GitRegistry implements Git-based registry access with caching.
type GitRegistry struct {
	cache        cache.Cache
	repo         Repository
	keyGen       cache.KeyGenerator
	registryURL  string
	registryType string
}

// NewGitRegistry creates a new Git-based registry.
func NewGitRegistry(cache cache.Cache, repo Repository, keyGen cache.KeyGenerator, registryURL, registryType string) *GitRegistry {
	return &GitRegistry{
		cache:        cache,
		repo:         repo,
		keyGen:       keyGen,
		registryURL:  registryURL,
		registryType: registryType,
	}
}

func (g *GitRegistry) ListVersions(ctx context.Context) ([]arm.VersionRef, error) {
	tags, err := g.repo.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	branches, err := g.repo.GetBranches(ctx)
	if err != nil {
		return nil, err
	}

	var versions []arm.VersionRef
	for _, tag := range tags {
		versions = append(versions, arm.VersionRef{ID: tag, Type: arm.Tag})
	}
	for _, branch := range branches {
		versions = append(versions, arm.VersionRef{ID: branch, Type: arm.Branch})
	}

	return versions, nil
}

func (g *GitRegistry) GetContent(ctx context.Context, version arm.VersionRef, selector arm.ContentSelector) ([]arm.File, error) {
	registryKey := g.keyGen.RegistryKey(g.registryURL, g.registryType)
	rulesetKey := g.keyGen.RulesetKey(selector)

	// Try cache first
	files, err := g.cache.Get(ctx, registryKey, rulesetKey, version.ID)
	if err == nil {
		return files, nil
	}

	// Checkout version and get files from repo
	if err := g.repo.Checkout(ctx, version.ID); err != nil {
		return nil, err
	}

	files, err = g.repo.GetFiles(ctx, selector)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = g.cache.Set(ctx, registryKey, rulesetKey, version.ID, files)

	return files, nil
}
