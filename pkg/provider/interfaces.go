package provider

import (
	"github.com/jomadu/ai-rules-manager/pkg/cache"
	"github.com/jomadu/ai-rules-manager/pkg/config"
	"github.com/jomadu/ai-rules-manager/pkg/registry"
	"github.com/jomadu/ai-rules-manager/pkg/version"
)

// RegistryProvider creates registry-specific components
type RegistryProvider interface {
	CreateRegistry(config *config.RegistryConfig) (registry.Registry, error)
	CreateVersionResolver() (version.VersionResolver, error)
	CreateContentResolver() (version.ContentResolver, error)
	CreateCacheKeyGenerator() (cache.CacheKeyGenerator, error)
	CreateRepositoryCache(basePath string) (cache.RepositoryCache, error)
	CreateRulesetCache(basePath string) (cache.RulesetCache, error)
}

// GitRegistryProvider implements RegistryProvider for Git repositories
type GitRegistryProvider struct{}

func NewGitRegistryProvider() *GitRegistryProvider {
	return &GitRegistryProvider{}
}

func (g *GitRegistryProvider) CreateRegistry(config *config.RegistryConfig) (registry.Registry, error) {
	return registry.NewGitRegistry(config.URL), nil
}

func (g *GitRegistryProvider) CreateVersionResolver() (version.VersionResolver, error) {
	return version.NewSemVerResolver(), nil
}

func (g *GitRegistryProvider) CreateContentResolver() (version.ContentResolver, error) {
	return version.NewGitContentResolver(), nil
}

func (g *GitRegistryProvider) CreateCacheKeyGenerator() (cache.CacheKeyGenerator, error) {
	return cache.NewGitCacheKeyGenerator(), nil
}

func (g *GitRegistryProvider) CreateRepositoryCache(basePath string) (cache.RepositoryCache, error) {
	return cache.NewGitRepositoryCache(basePath), nil
}

func (g *GitRegistryProvider) CreateRulesetCache(basePath string) (cache.RulesetCache, error) {
	return cache.NewFileRulesetCache(basePath), nil
}
