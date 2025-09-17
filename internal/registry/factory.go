package registry

import (
	"encoding/json"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/cache"
)

// NewRegistry creates a registry instance based on the registry configuration type.
// Accepts raw config data and handles type-specific parsing internally.
func NewRegistry(name string, rawConfig map[string]interface{}) (Registry, error) {
	registryType, ok := rawConfig["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid registry type")
	}

	switch registryType {
	case "git":
		return newGitRegistry(name, rawConfig)
	case "gitlab":
		return newGitLabRegistry(name, rawConfig)
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", registryType)
	}
}

func newGitRegistry(name string, rawConfig map[string]interface{}) (*GitRegistry, error) {
	// Parse raw config into GitRegistryConfig
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var gitConfig GitRegistryConfig
	if err := json.Unmarshal(configBytes, &gitConfig); err != nil {
		return nil, fmt.Errorf("failed to parse git registry config: %w", err)
	}

	registryKeyObj := map[string]string{
		"url":  gitConfig.URL,
		"type": gitConfig.Type,
	}

	rulesetCache, err := cache.NewRegistryRulesetCache(registryKeyObj)
	if err != nil {
		return nil, err
	}

	repoCache, err := cache.NewGitRepoCache(registryKeyObj, name, gitConfig.URL)
	if err != nil {
		return nil, err
	}

	return NewGitRegistry(gitConfig, rulesetCache, repoCache), nil
}

func newGitLabRegistry(name string, rawConfig map[string]interface{}) (*GitLabRegistry, error) {
	// Parse raw config into GitLabRegistryConfig
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var gitlabConfig GitLabRegistryConfig
	if err := json.Unmarshal(configBytes, &gitlabConfig); err != nil {
		return nil, fmt.Errorf("failed to parse gitlab registry config: %w", err)
	}

	// Build registry key object for cache uniqueness
	registryKeyObj := map[string]string{
		"url":  gitlabConfig.URL,
		"type": gitlabConfig.Type,
	}
	if gitlabConfig.ProjectID != "" {
		registryKeyObj["project_id"] = gitlabConfig.ProjectID
	}
	if gitlabConfig.GroupID != "" {
		registryKeyObj["group_id"] = gitlabConfig.GroupID
	}

	rulesetCache, err := cache.NewRegistryRulesetCache(registryKeyObj)
	if err != nil {
		return nil, err
	}

	return NewGitLabRegistry(name, &gitlabConfig, rulesetCache), nil
}
