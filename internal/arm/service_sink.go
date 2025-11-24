package arm

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/resource"
)

// AddSink adds a new sink to the ARM configuration
func (a *ArmService) AddSink(ctx context.Context, name, directory, layout, compileTarget string, force bool) error {
	// Use manifest manager to add sink
	sink := manifest.SinkConfig{
		Directory:     directory,
		Layout:        layout,
		CompileTarget: resource.CompileTarget(compileTarget),
	}
	if err := a.manifestManager.AddSink(ctx, name, sink, force); err != nil {
		return err
	}

	a.ui.Success(fmt.Sprintf("Sink %s added", name))
	return nil
}

// RemoveSink removes a sink from the ARM configuration
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

// SetSinkConfig sets configuration values for a specific sink
func (a *ArmService) SetSinkConfig(ctx context.Context, name, field, value string) error {
	err := a.manifestManager.UpdateSink(ctx, name, field, value)
	if err != nil {
		return err
	}
	a.ui.Success(fmt.Sprintf("Sink %s updated", name))
	return nil
}

// ShowSinkList lists all configured sinks
func (a *ArmService) ShowSinkList(ctx context.Context) error {
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return err
	}

	a.ui.SinkList(sinks)
	return nil
}

// ShowSinkInfo displays detailed information about one or more sinks
func (a *ArmService) ShowSinkInfo(ctx context.Context, sinks []string) error {
	allSinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return err
	}

	if len(sinks) == 0 {
		// Show info for all sinks
		for name, config := range allSinks {
			a.ui.SinkInfo(name, config)
		}
		return nil
	}

	// Show info for specific sinks
	for _, name := range sinks {
		config, exists := allSinks[name]
		if !exists {
			return fmt.Errorf("sink '%s' not found", name)
		}
		a.ui.SinkInfo(name, config)
	}
	return nil
}

// syncRemovedSink cleans up files from a removed sink directory
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
