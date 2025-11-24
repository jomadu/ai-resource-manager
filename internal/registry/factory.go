package registry

import (
	"encoding/json"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/cache"
)

// NewRegistry creates a registry instance based on the registry configuration type.
// Accepts raw config data and handles type-specific parsing internally.
func NewRegistry(rawConfig map[string]interface{}) (Registry, error) {
	registryType, ok := rawConfig["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid registry type")
	}

	switch registryType {
	case "git":
		return newGitRegistry(rawConfig)
	case "gitlab":
		return newGitLabRegistry(rawConfig)
	case "cloudsmith":
		return newCloudsmithRegistry(rawConfig)
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", registryType)
	}
}

func newGitRegistry(rawConfig map[string]interface{}) (*GitRegistry, error) {
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

	packageCache, err := cache.NewRegistryPackageCache(registryKeyObj)
	if err != nil {
		return nil, err
	}

	repoCache, err := cache.NewGitRepoCache(registryKeyObj, gitConfig.URL)
	if err != nil {
		return nil, err
	}

	return NewGitRegistry(gitConfig, packageCache, repoCache), nil
}

func newGitLabRegistry(rawConfig map[string]interface{}) (*GitLabRegistry, error) {
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

	packageCache, err := cache.NewRegistryPackageCache(registryKeyObj)
	if err != nil {
		return nil, err
	}

	return NewGitLabRegistry(&gitlabConfig, packageCache), nil
}

func newCloudsmithRegistry(rawConfig map[string]interface{}) (*CloudsmithRegistry, error) {
	// Parse raw config into CloudsmithRegistryConfig
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	var cloudsmithConfig CloudsmithRegistryConfig
	if err := json.Unmarshal(configBytes, &cloudsmithConfig); err != nil {
		return nil, fmt.Errorf("failed to parse cloudsmith registry config: %w", err)
	}

	// Build registry key object for cache uniqueness
	registryKeyObj := map[string]string{
		"url":        cloudsmithConfig.GetBaseURL(),
		"type":       cloudsmithConfig.Type,
		"owner":      cloudsmithConfig.Owner,
		"repository": cloudsmithConfig.Repository,
	}

	packageCache, err := cache.NewRegistryPackageCache(registryKeyObj)
	if err != nil {
		return nil, err
	}

	return NewCloudsmithRegistry(&cloudsmithConfig, packageCache), nil
}
