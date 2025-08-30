package registry

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/cache"
)

// GitRegistry implements Git-based registry access with caching.
// It abstracts all git repository operations - users simply create the registry
// and call methods like GetTags() without worrying about cloning, fetching, etc.
type GitRegistry struct {
	cache        cache.Cache
	repo         Repository
	keyGen       cache.KeyGenerator
	registryURL  string
	registryType string
	initialized  bool
}

// NewGitRegistry creates a new Git-based registry that handles all git operations internally.
// The registry will automatically clone the repository on first use and fetch updates as needed.
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
	if err := g.ensureInitialized(ctx); err != nil {
		return nil, err
	}

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
	if err := g.ensureInitialized(ctx); err != nil {
		return nil, err
	}

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

// GetTags returns all available tags from the registry.
func (g *GitRegistry) GetTags(ctx context.Context) ([]string, error) {
	if err := g.ensureInitialized(ctx); err != nil {
		return nil, err
	}
	return g.repo.GetTags(ctx)
}

// GetBranches returns all available branches from the registry.
func (g *GitRegistry) GetBranches(ctx context.Context) ([]string, error) {
	if err := g.ensureInitialized(ctx); err != nil {
		return nil, err
	}
	return g.repo.GetBranches(ctx)
}

func (g *GitRegistry) ensureInitialized(ctx context.Context) error {
	if g.initialized {
		return g.repo.Fetch(ctx)
	}

	workDir := g.repo.(*GitRepo).workDir
	if _, err := os.Stat(filepath.Join(workDir, ".git")); os.IsNotExist(err) {
		if err := g.repo.Clone(ctx, g.registryURL); err != nil {
			return err
		}
	} else {
		if err := g.repo.Fetch(ctx); err != nil {
			return err
		}
	}

	g.initialized = true
	return nil
}
