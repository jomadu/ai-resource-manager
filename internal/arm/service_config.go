package arm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/resource"
)

func (a *ArmService) ShowConfig(ctx context.Context) error {
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		registries = make(map[string]map[string]interface{})
	}

	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		sinks = make(map[string]manifest.SinkConfig)
	}

	a.ui.ConfigList(registries, sinks)
	return nil
}

func (a *ArmService) AddRegistry(ctx context.Context, name, url, regType string, options map[string]interface{}) error {
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
		err = a.manifestManager.AddGitRegistry(ctx, name, config, false)
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
		err = a.manifestManager.AddGitLabRegistry(ctx, name, config, false)
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
		err = a.manifestManager.AddCloudsmithRegistry(ctx, name, config, false)
	default:
		return fmt.Errorf("unsupported registry type: %s", regType)
	}

	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Registry %s added", name))
	return nil
}

func (a *ArmService) RemoveRegistry(ctx context.Context, name string) error {
	err := a.manifestManager.RemoveRegistry(ctx, name)
	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Registry %s removed", name))
	return nil
}

func (a *ArmService) AddSink(ctx context.Context, name, directory, sinkType, layout, compileTarget string, force bool) error {
	// Apply type-based defaults if sinkType is specified
	if sinkType != "" {
		switch sinkType {
		case "cursor":
			if layout == "" {
				layout = "hierarchical"
			}
			if compileTarget == "" {
				compileTarget = "cursor"
			}
		case "copilot":
			if layout == "" {
				layout = "flat"
			}
			if compileTarget == "" {
				compileTarget = "copilot"
			}
		case "amazonq":
			if layout == "" {
				layout = "hierarchical"
			}
			if compileTarget == "" {
				compileTarget = "amazonq"
			}
		default:
			return fmt.Errorf("type must be one of: cursor, copilot, amazonq")
		}
	}

	// Require either sinkType or compileTarget
	if sinkType == "" && compileTarget == "" {
		return fmt.Errorf("either --type or --compile-to is required")
	}

	// Validate compileTarget
	if compileTarget != "" && compileTarget != "cursor" && compileTarget != "amazonq" && compileTarget != "markdown" && compileTarget != "copilot" {
		return fmt.Errorf("compile-to must be one of: cursor, amazonq, markdown, copilot")
	}

	// Use manifest manager to add sink
	sink := manifest.SinkConfig{
		Directory:     directory,
		Layout:        layout,
		CompileTarget: resource.CompileTarget(compileTarget),
	}
	return a.manifestManager.AddSink(ctx, name, sink, force)
}

func (a *ArmService) RemoveSink(ctx context.Context, name string) error {
	// Get sink before removal for cleanup
	sink, err := a.manifestManager.GetSink(ctx, name)
	if err != nil {
		return err
	}

	// Remove from manifest
	err = a.manifestManager.RemoveSink(ctx, name)
	if err != nil {
		return err
	}

	// Clean files from sink directory
	a.syncRemovedSink(ctx, sink)

	a.ui.Success(fmt.Sprintf("Sink %s removed", name))
	return nil
}

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

func (a *ArmService) SetSinkConfig(ctx context.Context, name, field, value string) error {
	err := a.manifestManager.UpdateSink(ctx, name, field, value)
	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Sink %s updated", name))
	return nil
}

func (a *ArmService) syncRemovedSink(ctx context.Context, removedSink *manifest.SinkConfig) {
	installer := installer.NewInstaller(removedSink)

	// Scan removed sink directory to find installed rulesets
	rulesetInstallations, err := installer.ListInstalledRulesets(ctx)
	if err != nil {
		// Continue even if ruleset scan fails
		_ = err
	} else {
		// Uninstall all found rulesets from this directory
		for _, installation := range rulesetInstallations {
			if err := installer.UninstallRuleset(ctx, installation.Registry, installation.Ruleset); err != nil {
				// Continue on uninstall failure
				_ = err
			}
		}
	}

	// Scan removed sink directory to find installed promptsets
	promptsetInstallations, err := installer.ListInstalledPromptsets(ctx)
	if err != nil {
		// Continue even if promptset scan fails
		_ = err
	} else {
		// Uninstall all found promptsets from this directory
		for _, installation := range promptsetInstallations {
			if err := installer.UninstallPromptset(ctx, installation.Registry, installation.Promptset); err != nil {
				// Continue on uninstall failure
				_ = err
			}
		}
	}
}
