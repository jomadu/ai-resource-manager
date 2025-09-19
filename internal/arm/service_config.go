package arm

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/urf"
)

func (a *ArmService) ShowConfig(ctx context.Context) error {
	registries, err := a.manifestManager.GetRawRegistries(ctx)
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
	return a.manifestManager.AddSink(ctx, name, directory, layout, urf.CompileTarget(compileTarget), force)
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

func (a *ArmService) UpdateRegistry(ctx context.Context, name, field, value string) error {
	err := a.manifestManager.UpdateGitRegistry(ctx, name, field, value)
	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Registry %s updated", name))
	return nil
}

func (a *ArmService) UpdateSink(ctx context.Context, name, field, value string) error {
	err := a.manifestManager.UpdateSink(ctx, name, field, value)
	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Sink %s updated", name))
	return nil
}

func (a *ArmService) syncRemovedSink(ctx context.Context, removedSink *manifest.SinkConfig) {
	// Scan removed sink directory to find installed rulesets
	installer := installer.NewInstaller(removedSink)
	installations, err := installer.ListInstalled(ctx)
	if err != nil {
		return // Skip directory that can't be scanned
	}

	// Uninstall all found rulesets from this directory
	for _, installation := range installations {
		if err := installer.Uninstall(ctx, installation.Registry, installation.Ruleset); err != nil {
			// Continue on uninstall failure
			_ = err
		}
	}
}
