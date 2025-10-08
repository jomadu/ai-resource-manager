package arm

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/pterm/pterm"
)

func (a *ArmService) InstallRuleset(ctx context.Context, req *InstallRulesetRequest) error {
	// Load registries once and validate
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registries: %w", err)
	}
	if req.Registry == "" {
		return fmt.Errorf("registry is required")
	}
	if req.Ruleset == "" {
		return fmt.Errorf("ruleset is required")
	}
	if _, exists := registries[req.Registry]; !exists {
		return fmt.Errorf("registry %s not configured", req.Registry)
	}

	// Resolve version
	finishResolving := a.ui.InstallStepWithSpinner("Resolving version...")
	registryConfig := registries[req.Registry]
	registryClient, err := registry.NewRegistry(req.Registry, registryConfig)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}
	versionStr := req.Version
	if versionStr == "" {
		versionStr = "latest"
	}
	versionStr = expandVersionShorthand(versionStr)
	resolvedVersionResult, err := registryClient.ResolveVersion(ctx, req.Ruleset, versionStr)
	if err != nil {
		return fmt.Errorf("failed to resolve version: %w", err)
	}
	resolvedVersion := resolvedVersionResult.Version
	finishResolving(fmt.Sprintf("Version resolved... %s (from %s)", resolvedVersion.Display, versionStr))

	// Download content
	finishDownloading := a.ui.InstallStepWithSpinner("Downloading content...")
	selector := types.ContentSelector{Include: req.Include, Exclude: req.Exclude}
	files, err := registryClient.GetContent(ctx, req.Ruleset, resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to download content: %w", err)
	}
	finishDownloading(fmt.Sprintf("Downloaded content... %d files", len(files)))

	if err := a.updateTrackingFiles(ctx, req, resolvedVersion, files); err != nil {
		return fmt.Errorf("failed to update tracking files: %w", err)
	}

	_, err = a.installToSinks(ctx, req, resolvedVersion, files)
	if err != nil {
		return err
	}

	a.ui.InstallComplete(req.Registry, req.Ruleset, resolvedVersion.Display, req.Sinks)
	return nil
}

func (a *ArmService) InstallAll(ctx context.Context) error {
	manifestEntries, manifestErr := a.manifestManager.GetRulesets(ctx)
	lockEntries, lockErr := a.lockFileManager.GetRulesets(ctx)

	// Case: No manifest, no lockfile
	if manifestErr != nil && lockErr != nil {
		return fmt.Errorf("neither arm.json nor arm-lock.json found")
	}

	// Case: No manifest, lockfile exists
	if manifestErr != nil && lockErr == nil {
		return fmt.Errorf("arm.json not found")
	}

	// Case: Manifest exists, lockfile exists - use exact lockfile versions
	if manifestErr == nil && lockErr == nil {
		return a.installFromLockfile(ctx, lockEntries)
	}

	// Case: Manifest exists, no lockfile - resolve from manifest and create lockfile
	for registryName, rulesets := range manifestEntries {
		for rulesetName, entry := range rulesets {
			if err := a.InstallRuleset(ctx, &InstallRulesetRequest{
				Registry: registryName,
				Ruleset:  rulesetName,
				Version:  entry.Version,
				Include:  entry.GetIncludePatterns(),
				Exclude:  entry.Exclude,
				Sinks:    entry.Sinks,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) UninstallRuleset(ctx context.Context, registry, ruleset string) error {
	// Get manifest entry to find target sinks
	manifestEntry, err := a.manifestManager.GetRuleset(ctx, registry, ruleset)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	// Remove installed files from sink directories
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	for _, sinkName := range manifestEntry.Sinks {
		if sink, exists := sinks[sinkName]; exists {
			finishUninstalling := a.ui.InstallStepWithSpinner(fmt.Sprintf("Uninstalling from %s...", sinkName))
			installer := installer.NewInstaller(&sink)
			if err := installer.UninstallRuleset(ctx, registry, ruleset); err != nil {
				return err
			}
			finishUninstalling(fmt.Sprintf("Uninstalled from %s", sinkName))
		}
	}

	// Remove from manifest
	if err := a.manifestManager.RemoveRuleset(ctx, registry, ruleset); err != nil {
		return fmt.Errorf("failed to remove from manifest: %w", err)
	}

	// Remove from lockfile
	if err := a.lockFileManager.RemoveRuleset(ctx, registry, ruleset); err != nil {
		return fmt.Errorf("failed to remove from lockfile: %w", err)
	}

	a.ui.Success(fmt.Sprintf("Uninstalled %s/%s", registry, ruleset))
	return nil
}

func (a *ArmService) UpdateRuleset(ctx context.Context, registryName, rulesetName string) error {
	// Get manifest entry for version constraint
	manifestEntry, err := a.manifestManager.GetRuleset(ctx, registryName, rulesetName)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	// Resolve what version we should have
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registries: %w", err)
	}
	registryConfig, exists := registries[registryName]
	if !exists {
		return fmt.Errorf("registry %s not configured", registryName)
	}

	registryClient, err := registry.NewRegistry(registryName, registryConfig)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	versionStr := manifestEntry.Version
	if versionStr == "" {
		versionStr = "latest"
	}
	versionStr = expandVersionShorthand(versionStr)

	resolvedVersionResult, err := registryClient.ResolveVersion(ctx, rulesetName, versionStr)
	if err != nil {
		return fmt.Errorf("failed to resolve version: %w", err)
	}

	// Check what's actually installed in the filesystem
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	var isCurrentlyInstalled bool
	var installedVersion string

	// Check filesystem to see what's actually installed
	for _, sink := range sinks {
		installer := installer.NewInstaller(&sink)
		installed, version, err := installer.IsRulesetInstalled(ctx, registryName, rulesetName)
		if err != nil {
			continue
		}
		if installed {
			isCurrentlyInstalled = true
			installedVersion = version
			break
		}
	}

	if !isCurrentlyInstalled {
		// Nothing installed, proceed with install
		return a.InstallRuleset(ctx, &InstallRulesetRequest{
			Registry: registryName,
			Ruleset:  rulesetName,
			Version:  manifestEntry.Version,
			Include:  manifestEntry.GetIncludePatterns(),
			Exclude:  manifestEntry.Exclude,
		})
	}

	// Check if installed version matches what we want
	if installedVersion == resolvedVersionResult.Version.Display {
		// Get lockfile entry to verify checksum
		currentLockEntry, err := a.lockFileManager.GetRuleset(ctx, registryName, rulesetName)
		if err == nil {
			// Verify checksum to ensure integrity
			selector := types.ContentSelector{Include: manifestEntry.GetIncludePatterns(), Exclude: manifestEntry.Exclude}
			files, err := registryClient.GetContent(ctx, rulesetName, resolvedVersionResult.Version, selector)
			if err == nil && lockfile.VerifyChecksum(files, currentLockEntry.Checksum) {
				return nil
			}
		}
		// If we can't verify checksum, fall through to reinstall
	}

	// Version changed or integrity check failed, proceed with update
	return a.InstallRuleset(ctx, &InstallRulesetRequest{
		Registry: registryName,
		Ruleset:  rulesetName,
		Version:  manifestEntry.Version,
		Include:  manifestEntry.GetIncludePatterns(),
		Exclude:  manifestEntry.Exclude,
		Sinks:    manifestEntry.Sinks,
	})
}

func (a *ArmService) UpdateAll(ctx context.Context) error {
	return a.UpdateAllRulesets(ctx)
}

func (a *ArmService) UpdateAllRulesets(ctx context.Context) error {
	manifestEntries, manifestErr := a.manifestManager.GetRulesets(ctx)
	_, lockErr := a.lockFileManager.GetRulesets(ctx)

	// Case: No manifest, no lockfile
	if manifestErr != nil && lockErr != nil {
		return fmt.Errorf("neither arm.json nor arm-lock.json found")
	}

	// Case: No manifest, lockfile exists
	if manifestErr != nil && lockErr == nil {
		return fmt.Errorf("arm.json not found")
	}

	// Case: Manifest exists - update within version constraints
	for registryName, rulesets := range manifestEntries {
		for rulesetName := range rulesets {
			if err := a.UpdateRuleset(ctx, registryName, rulesetName); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) SetRulesetConfig(ctx context.Context, registry, ruleset, field, value string) error {

	// Get current manifest entry
	entry, err := a.manifestManager.GetRuleset(ctx, registry, ruleset)
	if err != nil {
		return fmt.Errorf("failed to get ruleset entry: %w", err)
	}

	// Update the specified field
	switch field {
	case "priority":
		priority := 0
		if _, err := fmt.Sscanf(value, "%d", &priority); err != nil {
			return fmt.Errorf("priority must be a number: %w", err)
		}
		entry.Priority = &priority
	case "version":
		entry.Version = value
	case "sinks":
		entry.Sinks = strings.Split(value, ",")
		for i, sink := range entry.Sinks {
			entry.Sinks[i] = strings.TrimSpace(sink)
		}
	case "include":
		entry.Include = strings.Split(value, ",")
		for i, pattern := range entry.Include {
			entry.Include[i] = strings.TrimSpace(pattern)
		}
	case "exclude":
		entry.Exclude = strings.Split(value, ",")
		for i, pattern := range entry.Exclude {
			entry.Exclude[i] = strings.TrimSpace(pattern)
		}
	default:
		return fmt.Errorf("unknown field '%s' (valid: priority, version, sinks, include, exclude)", field)
	}

	// Update manifest
	if err := a.manifestManager.UpdateRuleset(ctx, registry, ruleset, entry); err != nil {
		return fmt.Errorf("failed to update manifest: %w", err)
	}

	// Show what changed and that reinstall is needed
	switch field {
	case "priority":
		pterm.Info.Printf("Priority change to %s requires reinstall...\n", value)
	case "version":
		pterm.Info.Printf("Version change to %s requires reinstall...\n", value)
	case "sinks":
		pterm.Info.Printf("Sink change to %s requires reinstall...\n", value)
	case "include":
		pterm.Info.Printf("Include pattern change to %s requires reinstall...\n", value)
	case "exclude":
		pterm.Info.Printf("Exclude pattern change to %s requires reinstall...\n", value)
	}

	// Trigger reinstall
	return a.InstallRuleset(ctx, &InstallRulesetRequest{
		Registry: registry,
		Ruleset:  ruleset,
		Version:  entry.Version,
		Priority: *entry.Priority,
		Include:  entry.GetIncludePatterns(),
		Exclude:  entry.Exclude,
		Sinks:    entry.Sinks,
	})
}

// Private helper methods

func (a *ArmService) updateTrackingFiles(ctx context.Context, req *InstallRulesetRequest, version types.Version, files []types.File) error {
	// Store normalized version in manifest
	manifestVersion := req.Version
	if manifestVersion == "" {
		manifestVersion = "latest"
	}
	manifestVersion = expandVersionShorthand(manifestVersion)

	manifestEntry := manifest.RulesetConfig{
		Version:  manifestVersion,
		Priority: &req.Priority,
		Include:  req.Include,
		Exclude:  req.Exclude,
		Sinks:    req.Sinks,
	}
	if err := a.manifestManager.AddRuleset(ctx, req.Registry, req.Ruleset, &manifestEntry); err != nil {
		if err := a.manifestManager.UpdateRuleset(ctx, req.Registry, req.Ruleset, &manifestEntry); err != nil {
			return fmt.Errorf("failed to update manifest: %w", err)
		}
	}

	checksum := lockfile.GenerateChecksum(files)
	lockEntry := &lockfile.Entry{
		Version:  version.Version,
		Display:  version.Display,
		Checksum: checksum,
	}
	if err := a.lockFileManager.CreateOrUpdateRuleset(ctx, req.Registry, req.Ruleset, lockEntry); err != nil {
		return fmt.Errorf("failed to update lockfile: %w", err)
	}

	return nil
}

func (a *ArmService) cleanPreviousInstallation(ctx context.Context, registry, ruleset string) error {
	// Get current manifest entry to find previous sinks
	manifestEntry, err := a.manifestManager.GetRuleset(ctx, registry, ruleset)
	if err != nil {
		// No previous installation
		return nil
	}

	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return err
	}

	// Remove from previous sink locations
	for _, sinkName := range manifestEntry.Sinks {
		if sink, exists := sinks[sinkName]; exists {
			installer := installer.NewInstaller(&sink)
			if err := installer.UninstallRuleset(ctx, registry, ruleset); err != nil {
				// Continue on cleanup failure
				_ = err
			}
		}
	}

	return nil
}

func (a *ArmService) installToSinks(ctx context.Context, req *InstallRulesetRequest, version types.Version, files []types.File) (int, error) {
	// First, remove from previous sink locations if this is a reinstall
	if err := a.cleanPreviousInstallation(ctx, req.Registry, req.Ruleset); err != nil {
		// Continue on cleanup failure
		_ = err
	}

	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get sinks: %w", err)
	}

	// Validate that all requested sinks exist
	for _, sinkName := range req.Sinks {
		if _, exists := sinks[sinkName]; !exists {
			return 0, fmt.Errorf("sink %s not configured", sinkName)
		}
	}

	// Install to explicitly specified sinks
	totalInstalledFiles := 0
	for _, sinkName := range req.Sinks {
		finishInstalling := a.ui.InstallStepWithSpinner(fmt.Sprintf("Installing to %s...", sinkName))
		sink := sinks[sinkName]
		installer := installer.NewInstaller(&sink)
		if err := installer.InstallRuleset(ctx, req.Registry, req.Ruleset, version.Display, req.Priority, files); err != nil {
			return 0, err
		}
		finishInstalling(fmt.Sprintf("Installed to %s... %d files", sinkName, len(files)))
		totalInstalledFiles += len(files) // Each sink gets all files
	}

	return totalInstalledFiles, nil
}

func (a *ArmService) installFromLockfile(ctx context.Context, lockEntries map[string]map[string]lockfile.Entry) error {
	for registryName, rulesets := range lockEntries {
		for rulesetName, lockEntry := range rulesets {
			if err := a.installExactVersion(ctx, registryName, rulesetName, &lockEntry); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *ArmService) installExactVersion(ctx context.Context, registryName, ruleset string, lockEntry *lockfile.Entry) error {
	// Get registry config from manifest
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registries: %w", err)
	}
	registryConfig, exists := registries[registryName]
	if !exists {
		return fmt.Errorf("registry %s not found in manifest", registryName)
	}

	// Get manifest entry for include/exclude patterns
	manifestEntry, err := a.manifestManager.GetRuleset(ctx, registryName, ruleset)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	registryClient, err := registry.NewRegistry(registryName, registryConfig)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	resolvedVersion := types.Version{Version: lockEntry.Version, Display: lockEntry.Display}
	selector := types.ContentSelector{Include: manifestEntry.GetIncludePatterns(), Exclude: manifestEntry.Exclude}
	files, err := registryClient.GetContent(ctx, ruleset, resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to get content: %w", err)
	}

	// Verify checksum for integrity
	if !lockfile.VerifyChecksum(files, lockEntry.Checksum) {
		return fmt.Errorf("checksum verification failed for %s/%s@%s", registryName, ruleset, lockEntry.Version)
	}

	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	// Install to sinks specified in manifest entry
	for _, sinkName := range manifestEntry.Sinks {
		if sink, exists := sinks[sinkName]; exists {
			installer := installer.NewInstaller(&sink)
			// Use display version for directory names
			if err := installer.InstallRuleset(ctx, registryName, ruleset, lockEntry.Display, *manifestEntry.Priority, files); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) getOutdatedRulesets(ctx context.Context) ([]OutdatedRuleset, error) {
	lockEntries, err := a.lockFileManager.GetRulesets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	manifestEntries, err := a.manifestManager.GetRulesets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest entries: %w", err)
	}

	registryConfigs, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registries: %w", err)
	}

	// Pre-create registry clients
	registryClients := make(map[string]registry.Registry)
	for registryName, registryConfig := range registryConfigs {
		if client, err := registry.NewRegistry(registryName, registryConfig); err == nil {
			registryClients[registryName] = client
		}
	}

	var outdated []OutdatedRuleset
	for registryName, rulesets := range lockEntries {
		registryClient, exists := registryClients[registryName]
		if !exists {
			continue
		}

		for rulesetName, lockEntry := range rulesets {
			// Get constraint from manifest
			var constraint string
			if manifestRegistry, exists := manifestEntries[registryName]; exists {
				if manifestEntry, exists := manifestRegistry[rulesetName]; exists {
					constraint = manifestEntry.Version
				} else {
					continue
				}
			} else {
				continue
			}

			// Get latest version using proper resolution (prefers latest tag, falls back to default branch)
			latestVersion, err := registryClient.ResolveVersion(ctx, rulesetName, "latest")
			if err != nil {
				continue
			}

			wantedVersion, err := registryClient.ResolveVersion(ctx, rulesetName, constraint)
			if err != nil {
				continue
			}

			if lockEntry.Version != latestVersion.Version.Version || lockEntry.Version != wantedVersion.Version.Version {
				rulesetInfo, err := a.getRulesetInfo(ctx, registryName, rulesetName)
				if err != nil {
					continue
				}
				outdated = append(outdated, OutdatedRuleset{
					RulesetInfo: rulesetInfo,
					Wanted:      wantedVersion.Version.Display,
					Latest:      latestVersion.Version.Display,
				})
			}
		}
	}

	return outdated, nil
}

func (a *ArmService) listInstalledRulesets(ctx context.Context) ([]*RulesetInfo, error) {
	lockEntries, err := a.lockFileManager.GetRulesets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	var rulesets []*RulesetInfo
	for registryName, rulesetMap := range lockEntries {
		for rulesetName := range rulesetMap {
			rulesetInfo, err := a.getRulesetInfo(ctx, registryName, rulesetName)
			if err != nil {
				continue
			}
			rulesets = append(rulesets, rulesetInfo)
		}
	}

	// Sort by registry then ruleset name
	sort.Slice(rulesets, func(i, j int) bool {
		if rulesets[i].Registry != rulesets[j].Registry {
			return rulesets[i].Registry < rulesets[j].Registry
		}
		return rulesets[i].Name < rulesets[j].Name
	})

	return rulesets, nil
}

func (a *ArmService) getRulesetInfo(ctx context.Context, registry, ruleset string) (*RulesetInfo, error) {
	// Get manifest entry
	manifestEntry, err := a.manifestManager.GetRuleset(ctx, registry, ruleset)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest entry: %w", err)
	}

	// Get sinks and find installation paths and version
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sinks: %w", err)
	}

	var installedPaths []string
	var resolvedVersion string
	// Use sinks from manifest entry
	for _, sinkName := range manifestEntry.Sinks {
		if sink, exists := sinks[sinkName]; exists {
			installer := installer.NewInstaller(&sink)
			installations, err := installer.ListInstalledRulesets(ctx)
			if err != nil {
				continue
			}
			for _, installation := range installations {
				if installation.Registry == registry && installation.Ruleset == ruleset {
					installedPaths = append(installedPaths, installation.Path)
					if resolvedVersion == "" {
						resolvedVersion = installation.Version
					}
					break
				}
			}
		}
	}

	return &RulesetInfo{
		Registry: registry,
		Name:     ruleset,
		Manifest: ManifestInfo{
			Constraint: manifestEntry.Version,
			Priority:   *manifestEntry.Priority,
			Include:    manifestEntry.Include,
			Exclude:    manifestEntry.Exclude,
			Sinks:      manifestEntry.Sinks,
		},
		Installation: InstallationInfo{
			Version:        resolvedVersion,
			InstalledPaths: installedPaths,
		},
	}, nil
}

// expandVersionShorthand expands npm-style version shorthands to proper semantic version constraints.
// "1" -> "^1.0.0", "1.0" -> "~1.0.0"
func expandVersionShorthand(constraint string) string {
	// Match pure major version (e.g., "1")
	if matched, _ := regexp.MatchString(`^\d+$`, constraint); matched {
		return "^" + constraint + ".0.0"
	}
	// Match major.minor version (e.g., "1.0")
	if matched, _ := regexp.MatchString(`^\d+\.\d+$`, constraint); matched {
		return "~" + constraint + ".0"
	}
	return constraint
}

// Promptset operations - TODO: Implement these methods
func (a *ArmService) InstallPromptset(ctx context.Context, req *InstallPromptsetRequest) error {
	return fmt.Errorf("InstallPromptset not yet implemented")
}

func (a *ArmService) UninstallPromptset(ctx context.Context, registry, promptset string) error {
	return fmt.Errorf("UninstallPromptset not yet implemented")
}

func (a *ArmService) UpdatePromptset(ctx context.Context, registry, promptset string) error {
	return fmt.Errorf("UpdatePromptset not yet implemented")
}

func (a *ArmService) SetPromptsetConfig(ctx context.Context, registry, promptset, field, value string) error {
	return fmt.Errorf("SetPromptsetConfig not yet implemented")
}

// Unified operations - TODO: Implement these methods
func (a *ArmService) UpgradeAll(ctx context.Context) error {
	return fmt.Errorf("UpgradeAll not yet implemented")
}

func (a *ArmService) UninstallAll(ctx context.Context) error {
	return fmt.Errorf("UninstallAll not yet implemented")
}

// Info operations - TODO: Implement this method
func (a *ArmService) ShowPromptsetInfo(ctx context.Context, promptsets []string) error {
	return fmt.Errorf("ShowPromptsetInfo not yet implemented")
}
