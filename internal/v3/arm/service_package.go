package arm

import (
	"context"
	"fmt"
)

// InstallAll installs all configured packages (rulesets and promptsets)
func (a *ArmService) InstallAll(ctx context.Context) error {
	manifestRulesets, manifestRulesetsErr := a.manifestManager.GetRulesets(ctx)
	lockRulesets, lockRulesetsErr := a.lockFileManager.GetRulesets(ctx)

	manifestPromptsets, manifestPromptsetsErr := a.manifestManager.GetPromptsets(ctx)
	lockPromptsets, lockPromptsetsErr := a.lockFileManager.GetPromptsets(ctx)

	// Determine installation strategy based on available files
	type installCase int
	const (
		noManifestNoLock installCase = iota
		noManifestWithLock
		manifestWithLock
		manifestNoLock
	)

	var currentCase installCase
	switch {
	case manifestRulesetsErr != nil && lockRulesetsErr != nil && manifestPromptsetsErr != nil && lockPromptsetsErr != nil:
		currentCase = noManifestNoLock
	case manifestRulesetsErr != nil && lockRulesetsErr == nil:
		currentCase = noManifestWithLock
	case manifestRulesetsErr == nil && lockRulesetsErr == nil:
		currentCase = manifestWithLock
	default:
		currentCase = manifestNoLock
	}

	switch currentCase {
	case noManifestNoLock:
		return fmt.Errorf("neither arm.json nor arm-lock.json found")
	case noManifestWithLock:
		return fmt.Errorf("arm.json not found")
	case manifestWithLock:
		// Install rulesets from lockfile
		for registryName, rulesets := range lockRulesets {
			for rulesetName, lockEntry := range rulesets {
				if err := a.installExactRulesetVersion(ctx, registryName, rulesetName, &lockEntry); err != nil {
					return err
				}
			}
		}
	case manifestNoLock:
		// Install rulesets from manifest and create lockfile
		for registryName, rulesets := range manifestRulesets {
			for rulesetName, entry := range rulesets {
				priority := 100
				if entry.Priority != nil {
					priority = *entry.Priority
				}

				req := NewInstallRulesetRequest(
					registryName,
					rulesetName,
					entry.Version,
					entry.Sinks,
				).WithPriority(priority).
					WithInclude(entry.GetIncludePatterns()).
					WithExclude(entry.Exclude)

				if err := a.InstallRuleset(ctx, req); err != nil {
					return err
				}
			}
		}
	}

	// Handle promptsets - install from lockfile if available, otherwise from manifest
	switch {
	case manifestPromptsetsErr == nil && lockPromptsetsErr == nil:
		// Install promptsets from lockfile
		for registryName, promptsets := range lockPromptsets {
			for promptsetName, lockEntry := range promptsets {
				if err := a.installPromptsetExactVersion(ctx, registryName, promptsetName, &lockEntry); err != nil {
					return err
				}
			}
		}
	default:
		// Install promptsets from manifest
		for registryName, promptsets := range manifestPromptsets {
			for promptsetName, entry := range promptsets {
				if err := a.InstallPromptset(ctx, &InstallPromptsetRequest{
					Registry:  registryName,
					Promptset: promptsetName,
					Version:   entry.Version,
					Include:   entry.Include,
					Exclude:   entry.Exclude,
					Sinks:     entry.Sinks,
				}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UpdateAll updates all installed packages to their latest available versions within constraints
func (a *ArmService) UpdateAll(ctx context.Context) error {
	// Update all rulesets
	if err := a.UpdateAllRulesets(ctx); err != nil {
		return fmt.Errorf("failed to update rulesets: %w", err)
	}

	// Update all promptsets
	if err := a.UpdateAllPromptsets(ctx); err != nil {
		return fmt.Errorf("failed to update promptsets: %w", err)
	}

	return nil
}

// UpgradeAll upgrades all installed packages to their latest available versions, ignoring constraints
func (a *ArmService) UpgradeAll(ctx context.Context) error {
	// Get all installed rulesets and promptsets
	rulesets, err := a.manifestManager.GetRulesets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rulesets: %w", err)
	}

	promptsets, err := a.manifestManager.GetPromptsets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get promptsets: %w", err)
	}

	// Upgrade all rulesets
	for registry, registryRulesets := range rulesets {
		for ruleset := range registryRulesets {
			err = a.UpgradeRuleset(ctx, registry, ruleset)
			if err != nil {
				return fmt.Errorf("failed to upgrade ruleset %s/%s: %w", registry, ruleset, err)
			}
		}
	}

	// Upgrade all promptsets
	for registry, registryPromptsets := range promptsets {
		for promptset := range registryPromptsets {
			err = a.UpgradePromptset(ctx, registry, promptset)
			if err != nil {
				return fmt.Errorf("failed to upgrade promptset %s/%s: %w", registry, promptset, err)
			}
		}
	}
	return nil
}

// UninstallAll uninstalls all configured packages from their assigned sinks
func (a *ArmService) UninstallAll(ctx context.Context) error {
	// Get all rulesets
	allRulesets, err := a.manifestManager.GetRulesets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rulesets: %w", err)
	}

	// Get all promptsets
	allPromptsets, err := a.manifestManager.GetPromptsets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get promptsets: %w", err)
	}

	// Uninstall all rulesets
	for registry, registryRulesets := range allRulesets {
		for ruleset := range registryRulesets {
			if err := a.UninstallRuleset(ctx, registry, ruleset); err != nil {
				return fmt.Errorf("failed to uninstall ruleset %s/%s: %w", registry, ruleset, err)
			}
		}
	}

	// Uninstall all promptsets
	for registry, registryPromptsets := range allPromptsets {
		for promptset := range registryPromptsets {
			if err := a.UninstallPromptset(ctx, registry, promptset); err != nil {
				return fmt.Errorf("failed to uninstall promptset %s/%s: %w", registry, promptset, err)
			}
		}
	}

	return nil
}

// GetOutdatedPackages returns all outdated packages (rulesets and promptsets)
func (a *ArmService) GetOutdatedPackages(ctx context.Context) ([]*OutdatedPackage, error) {
	rulesetOutdated, err := a.getOutdatedRulesets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get outdated rulesets: %w", err)
	}

	promptsetOutdated, err := a.getOutdatedPromptsets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get outdated promptsets: %w", err)
	}

	var allPackages []*OutdatedPackage

	for _, ruleset := range rulesetOutdated {
		allPackages = append(allPackages, &OutdatedPackage{
			Package:    fmt.Sprintf("%s/%s", ruleset.RulesetInfo.Registry, ruleset.RulesetInfo.Name),
			Type:       "ruleset",
			Constraint: ruleset.RulesetInfo.Manifest.Constraint,
			Current:    ruleset.RulesetInfo.Installation.Version,
			Wanted:     ruleset.Wanted,
			Latest:     ruleset.Latest,
		})
	}

	for _, promptset := range promptsetOutdated {
		allPackages = append(allPackages, &OutdatedPackage{
			Package:    fmt.Sprintf("%s/%s", promptset.PromptsetInfo.Registry, promptset.PromptsetInfo.Name),
			Type:       "promptset",
			Constraint: promptset.PromptsetInfo.Manifest.Constraint,
			Current:    promptset.PromptsetInfo.Installation.Version,
			Wanted:     promptset.Wanted,
			Latest:     promptset.Latest,
		})
	}

	return allPackages, nil
}