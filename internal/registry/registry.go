package registry

import (
	"context"
	"errors"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/cache"
)

// GitRegistry implements Git-based registry access with caching.
type GitRegistry struct {
	cache  cache.Cache
	repo   Repository
	keyGen cache.KeyGenerator
}

// NewGitRegistry creates a new Git-based registry.
func NewGitRegistry(cache cache.Cache, repo Repository, keyGen cache.KeyGenerator) *GitRegistry {
	return &GitRegistry{
		cache:  cache,
		repo:   repo,
		keyGen: keyGen,
	}
}

func (g *GitRegistry) ListVersions(ctx context.Context) ([]arm.VersionRef, error) {
	return nil, errors.New("not implemented")
}

func (g *GitRegistry) GetContent(ctx context.Context, version arm.VersionRef, selector arm.ContentSelector) ([]arm.File, error) {
	return nil, errors.New("not implemented")
}
