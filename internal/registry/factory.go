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
		return newGitRegistry(name, config)
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", config.Type)
	}
}

func newGitRegistry(name string, config config.RegistryConfig) (*GitRegistry, error) {
	registryKeyObj := map[string]string{
		"url":  config.URL,
		"type": config.Type,
	}

	rulesetCache, err := cache.NewRegistryRulesetCache(registryKeyObj)
	if err != nil {
		return nil, err
	}

	repoCache, err := cache.NewGitRepoCache(registryKeyObj, name, config.URL)
	if err != nil {
		return nil, err
	}

	return NewGitRegistry(rulesetCache, repoCache, config.URL, config.Type), nil
}
