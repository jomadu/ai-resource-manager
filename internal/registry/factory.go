package registry

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/config"
)

// NewRegistry creates a registry instance based on the registry configuration type.
func NewRegistry(name string, config config.RegistryConfig) (Registry, error) {
	switch config.Type {
	case "git":
		return newGitRegistry(name, config), nil
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", config.Type)
	}
}

func newGitRegistry(name string, config config.RegistryConfig) *GitRegistry {
	keyGen := cache.NewGitKeyGen()
	registryKey := keyGen.RegistryKey(config.URL, config.Type)
	rulesetCache := cache.NewRegistryRulesetCache(registryKey, config.URL, config.Type)
	repoCache := cache.NewGitRepoCache(registryKey, name, config.URL)

	return NewGitRegistry(rulesetCache, repoCache, config.URL, config.Type)
}
