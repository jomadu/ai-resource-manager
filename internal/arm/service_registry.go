package arm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/registry"
)

// AddRegistry adds a new registry to the ARM configuration
func (a *ArmService) AddRegistry(ctx context.Context, name, url, regType string, options map[string]interface{}, force bool) error {
	var err error
	switch regType {
	case "git":
		config := registry.GitRegistryConfig{
			RegistryConfig: registry.RegistryConfig{
				Type: "git",
				URL:  url,
			},
		}
		if branches, ok := options["branches"].([]string); ok {
			config.Branches = branches
		}
		err = a.manifestManager.AddGitRegistry(ctx, name, config, force)
	case "gitlab":
		config := &registry.GitLabRegistryConfig{
			RegistryConfig: registry.RegistryConfig{
				Type: "gitlab",
				URL:  url,
			},
		}
		if projectID, ok := options["project_id"].(string); ok {
			config.ProjectID = projectID
		}
		if groupID, ok := options["group_id"].(string); ok {
			config.GroupID = groupID
		}
		if apiVersion, ok := options["api_version"].(string); ok {
			config.APIVersion = apiVersion
		}
		err = a.manifestManager.AddGitLabRegistry(ctx, name, config, force)
	case "cloudsmith":
		config := &registry.CloudsmithRegistryConfig{
			RegistryConfig: registry.RegistryConfig{
				Type: "cloudsmith",
				URL:  url,
			},
		}
		if owner, ok := options["owner"].(string); ok {
			config.Owner = owner
		}
		if repository, ok := options["repository"].(string); ok {
			config.Repository = repository
		}
		err = a.manifestManager.AddCloudsmithRegistry(ctx, name, config, force)
	default:
		return fmt.Errorf("unsupported registry type: %s", regType)
	}

	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Registry %s added", name))
	return nil
}

// RemoveRegistry removes a registry from the ARM configuration
func (a *ArmService) RemoveRegistry(ctx context.Context, name string) error {
	err := a.manifestManager.RemoveRegistry(ctx, name)
	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Registry %s removed", name))
	return nil
}

// ListRegistries lists all configured registries
func (a *ArmService) ListRegistries(ctx context.Context) error {
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return err
	}

	a.ui.RegistryList(registries)
	return nil
}

// ShowRegistryInfo displays detailed information about one or more registries
func (a *ArmService) ShowRegistryInfo(ctx context.Context, registries []string) error {
	allRegistries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return err
	}

	if len(registries) == 0 {
		// Show info for all registries
		for name, config := range allRegistries {
			a.ui.RegistryInfo(name, config)
		}
		return nil
	}

	// Show info for specific registries
	for _, name := range registries {
		config, exists := allRegistries[name]
		if !exists {
			return fmt.Errorf("registry '%s' not found", name)
		}
		a.ui.RegistryInfo(name, config)
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
	a.ui.Success(fmt.Sprintf("Git registry %s updated", name))
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
	a.ui.Success(fmt.Sprintf("GitLab registry %s updated", name))
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
	a.ui.Success(fmt.Sprintf("Cloudsmith registry %s updated", name))
	return nil
}
