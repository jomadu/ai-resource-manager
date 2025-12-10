package arm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/registry"
)

// AddGitRegistry adds a new Git registry to the ARM configuration
func (a *ArmService) AddGitRegistry(ctx context.Context, name, url string, branches []string, force bool) error {
	config := registry.GitRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			Type: "git",
			URL:  url,
		},
		Branches: branches,
	}

	err := a.manifestManager.AddGitRegistry(ctx, name, config, force)
	if err != nil {
		return err
	}
	return nil
}

// AddGitLabRegistry adds a new GitLab registry to the ARM configuration
func (a *ArmService) AddGitLabRegistry(ctx context.Context, name, url, projectID, groupID, apiVersion string, force bool) error {
	config := &registry.GitLabRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			Type: "gitlab",
			URL:  url,
		},
		ProjectID:  projectID,
		GroupID:    groupID,
		APIVersion: apiVersion,
	}

	err := a.manifestManager.AddGitLabRegistry(ctx, name, config, force)
	if err != nil {
		return err
	}
	return nil
}

// AddCloudsmithRegistry adds a new Cloudsmith registry to the ARM configuration
func (a *ArmService) AddCloudsmithRegistry(ctx context.Context, name, url, owner, repository string, force bool) error {
	config := &registry.CloudsmithRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			Type: "cloudsmith",
			URL:  url,
		},
		Owner:      owner,
		Repository: repository,
	}

	err := a.manifestManager.AddCloudsmithRegistry(ctx, name, config, force)
	if err != nil {
		return err
	}
	return nil
}

// RemoveRegistry removes a registry from the ARM configuration
func (a *ArmService) RemoveRegistry(ctx context.Context, name string) error {
	err := a.manifestManager.RemoveRegistry(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

// SetRegistryConfig sets configuration values for a specific registry
func (a *ArmService) SetRegistryConfig(ctx context.Context, name, field, value string) error {
	// Get raw registry config to determine type
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registries: %w", err)
	}

	rawConfig, exists := registries[name]
	if !exists {
		return fmt.Errorf("registry %s not found", name)
	}

	// Determine registry type
	regType, ok := rawConfig["type"].(string)
	if !ok {
		return fmt.Errorf("registry %s has no type field", name)
	}

	// Handle different registry types
	switch regType {
	case "git":
		return a.setGitRegistryConfig(ctx, name, field, value)
	case "gitlab":
		return a.setGitLabRegistryConfig(ctx, name, field, value)
	case "cloudsmith":
		return a.setCloudsmithRegistryConfig(ctx, name, field, value)
	default:
		return fmt.Errorf("unsupported registry type: %s", regType)
	}
}

// setGitRegistryConfig sets configuration for a Git registry
func (a *ArmService) setGitRegistryConfig(ctx context.Context, name, field, value string) error {
	config, err := a.manifestManager.GetGitRegistry(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get git registry config: %w", err)
	}

	switch field {
	case "url":
		config.URL = value
	case "type":
		if value != "git" {
			return fmt.Errorf("type must be 'git'")
		}
		config.Type = value
	case "branches":
		branches := strings.Split(value, ",")
		for i, branch := range branches {
			branches[i] = strings.TrimSpace(branch)
		}
		config.Branches = branches
	default:
		return fmt.Errorf("unknown field '%s' for git registry (valid: url, type, branches)", field)
	}

	err = a.manifestManager.UpdateGitRegistry(ctx, name, *config, true)
	if err != nil {
		return err
	}
	return nil
}

// setGitLabRegistryConfig sets configuration for a GitLab registry
func (a *ArmService) setGitLabRegistryConfig(ctx context.Context, name, field, value string) error {
	config, err := a.manifestManager.GetGitLabRegistry(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get gitlab registry config: %w", err)
	}

	switch field {
	case "url":
		config.URL = value
	case "type":
		if value != "gitlab" {
			return fmt.Errorf("type must be 'gitlab'")
		}
		config.Type = value
	case "projectId":
		config.ProjectID = value
	case "groupId":
		config.GroupID = value
	case "apiVersion":
		config.APIVersion = value
	default:
		return fmt.Errorf("unknown field '%s' for gitlab registry (valid: url, type, projectId, groupId, apiVersion)", field)
	}

	err = a.manifestManager.UpdateGitLabRegistry(ctx, name, config, true)
	if err != nil {
		return err
	}
	return nil
}

// setCloudsmithRegistryConfig sets configuration for a Cloudsmith registry
func (a *ArmService) setCloudsmithRegistryConfig(ctx context.Context, name, field, value string) error {
	config, err := a.manifestManager.GetCloudsmithRegistry(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get cloudsmith registry config: %w", err)
	}

	switch field {
	case "url":
		config.URL = value
	case "type":
		if value != "cloudsmith" {
			return fmt.Errorf("type must be 'cloudsmith'")
		}
		config.Type = value
	case "owner":
		config.Owner = value
	case "repository":
		config.Repository = value
	default:
		return fmt.Errorf("unknown field '%s' for cloudsmith registry (valid: url, type, owner, repository)", field)
	}

	err = a.manifestManager.UpdateCloudsmithRegistry(ctx, name, config, true)
	if err != nil {
		return err
	}
	return nil
}

// GetRegistries returns all configured registries
func (a *ArmService) GetRegistries(ctx context.Context) (map[string]map[string]interface{}, error) {
	return a.manifestManager.GetRegistries(ctx)
}
