package arm

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/jomadu/ai-rules-manager/internal/ui"
	"github.com/pterm/pterm"
)

func (a *ArmService) InstallPromptset(ctx context.Context, req *InstallPromptsetRequest) error {
	// Load registries once and validate
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registries: %w", err)
	}
	if req.Registry == "" {
		return fmt.Errorf("registry is required")
	}
	if req.Promptset == "" {
		return fmt.Errorf("promptset is required")
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
	resolvedVersionResult, err := registryClient.ResolveVersion(ctx, req.Promptset, versionStr)
	if err != nil {
		return fmt.Errorf("failed to resolve version: %w", err)
	}
	resolvedVersion := resolvedVersionResult.Version
	finishResolving(fmt.Sprintf("Version resolved... %s (from %s)", resolvedVersion.Display, versionStr))

	// Download content
	finishDownloading := a.ui.InstallStepWithSpinner("Downloading content...")
	selector := types.ContentSelector{Include: req.Include, Exclude: req.Exclude}
	files, err := registryClient.GetContent(ctx, req.Promptset, resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to download content: %w", err)
	}
	finishDownloading(fmt.Sprintf("Downloaded content... %d files", len(files)))

	if err := a.updatePromptsetTrackingFiles(ctx, req, resolvedVersion, files); err != nil {
		return fmt.Errorf("failed to update tracking files: %w", err)
	}

	_, err = a.installPromptsetToSinks(ctx, req, resolvedVersion, files)
	if err != nil {
		return fmt.Errorf("failed to install to sinks: %w", err)
	}

	a.ui.InstallComplete(req.Registry, req.Promptset, resolvedVersion.Display, "promptset", req.Sinks)
	return nil
}

func (a *ArmService) UninstallPromptset(ctx context.Context, registry, promptset string) error {
	// Get promptset info from lockfile to know what files to remove
	_, err := a.lockFileManager.GetPromptset(ctx, registry, promptset)
	if err != nil {
		return fmt.Errorf("promptset %s/%s not found in lockfile", registry, promptset)
	}

	// Get promptset config from manifest to know which sinks to clean
	promptsetConfig, err := a.manifestManager.GetPromptset(ctx, registry, promptset)
	if err != nil {
		return fmt.Errorf("promptset %s/%s not found in manifest", registry, promptset)
	}

	// Remove from all sinks
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	for _, sinkName := range promptsetConfig.Sinks {
		sinkConfig, exists := sinks[sinkName]
		if !exists {
			continue // Sink no longer exists, skip
		}

		// Create installer for this sink
		installer := installer.NewInstaller(&sinkConfig)

		// Uninstall promptset from sink
		err = installer.UninstallPromptset(ctx, registry, promptset)
		if err != nil {
			return fmt.Errorf("failed to uninstall promptset from sink %s: %w", sinkName, err)
		}
	}

	// Remove from manifest
	err = a.manifestManager.RemovePromptset(ctx, registry, promptset)
	if err != nil {
		return fmt.Errorf("failed to remove promptset from manifest: %w", err)
	}

	// Remove from lockfile
	err = a.lockFileManager.RemovePromptset(ctx, registry, promptset)
	if err != nil {
		return fmt.Errorf("failed to remove promptset from lockfile: %w", err)
	}

	a.ui.Success(fmt.Sprintf("Uninstalled promptset %s/%s", registry, promptset))
	return nil
}

// updatePromptsetTrackingFiles updates manifest and lockfile for promptset installation

func (a *ArmService) updatePromptsetTrackingFiles(ctx context.Context, req *InstallPromptsetRequest, version types.Version, files []types.File) error {
	// Store normalized version in manifest
	manifestVersion := req.Version
	if manifestVersion == "" {
		manifestVersion = "latest"
	}
	manifestVersion = expandVersionShorthand(manifestVersion)

	// Create or update manifest entry
	promptsetConfig := &manifest.PromptsetConfig{
		Version: manifestVersion,
		Include: req.Include,
		Exclude: req.Exclude,
		Sinks:   req.Sinks,
	}

	if err := a.manifestManager.CreateOrUpdatePromptset(ctx, req.Registry, req.Promptset, promptsetConfig); err != nil {
		return fmt.Errorf("failed to update manifest: %w", err)
	}

	// Generate checksum and store both version fields
	checksum := lockfile.GenerateChecksum(files)
	lockEntry := &lockfile.Entry{
		Version:  version.Version,
		Display:  version.Display,
		Checksum: checksum,
	}
	err := a.lockFileManager.CreateOrUpdatePromptset(ctx, req.Registry, req.Promptset, lockEntry)
	if err != nil {
		return fmt.Errorf("failed to update lockfile: %w", err)
	}

	return nil
}

// installPromptsetToSinks installs promptset files to specified sinks

func (a *ArmService) installPromptsetToSinks(ctx context.Context, req *InstallPromptsetRequest, version types.Version, files []types.File) (int, error) {
	// First, remove from previous sink locations if this is a reinstall
	if err := a.cleanPreviousPromptsetInstallation(ctx, req.Registry, req.Promptset); err != nil {
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
		if err := installer.InstallPromptset(ctx, req.Registry, req.Promptset, version.Display, files); err != nil {
			return 0, err
		}
		finishInstalling(fmt.Sprintf("Installed to %s... %d files", sinkName, len(files)))
		totalInstalledFiles += len(files) // Each sink gets all files
	}

	return totalInstalledFiles, nil
}

func (a *ArmService) UpdatePromptset(ctx context.Context, registry, promptset string) error {
	// Get current promptset config
	promptsetConfig, err := a.manifestManager.GetPromptset(ctx, registry, promptset)
	if err != nil {
		return fmt.Errorf("promptset %s/%s not found in manifest", registry, promptset)
	}

	// Create install request with current config
	req := &InstallPromptsetRequest{
		Registry:  registry,
		Promptset: promptset,
		Version:   promptsetConfig.Version,
		Include:   promptsetConfig.Include,
		Exclude:   promptsetConfig.Exclude,
		Sinks:     promptsetConfig.Sinks,
	}

	// Install with updated version (this will resolve to latest version that satisfies constraint)
	return a.InstallPromptset(ctx, req)
}

func (a *ArmService) SetPromptsetConfig(ctx context.Context, registry, promptset, field, value string) error {
	// Get current promptset config
	promptsetConfig, err := a.manifestManager.GetPromptset(ctx, registry, promptset)
	if err != nil {
		return fmt.Errorf("promptset %s/%s not found in manifest", registry, promptset)
	}

	// Update the specified field
	switch field {
	case "version":
		promptsetConfig.Version = value
	case "sinks":
		// Parse comma-separated sinks
		sinks := strings.Split(value, ",")
		for i, sink := range sinks {
			sinks[i] = strings.TrimSpace(sink)
		}
		promptsetConfig.Sinks = sinks
	case "include":
		// Parse comma-separated include
		include := strings.Split(value, ",")
		for i, entry := range include {
			include[i] = strings.TrimSpace(entry)
		}
		promptsetConfig.Include = include
	case "exclude":
		// Parse comma-separated exclude
		exclude := strings.Split(value, ",")
		for i, entry := range exclude {
			exclude[i] = strings.TrimSpace(entry)
		}
		promptsetConfig.Exclude = exclude
	default:
		return fmt.Errorf("unknown field: %s (supported: version, sinks, include, exclude)", field)
	}

	// Update the manifest
	err = a.manifestManager.CreateOrUpdatePromptset(ctx, registry, promptset, promptsetConfig)
	if err != nil {
		return fmt.Errorf("failed to update promptset config: %w", err)
	}

	// Show what changed and that reinstall is needed
	switch field {
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
	return a.InstallPromptset(ctx, &InstallPromptsetRequest{
		Registry:  registry,
		Promptset: promptset,
		Version:   promptsetConfig.Version,
		Include:   promptsetConfig.Include,
		Exclude:   promptsetConfig.Exclude,
		Sinks:     promptsetConfig.Sinks,
	})
}

// Unified operations

func (a *ArmService) UpgradePromptset(ctx context.Context, registry, promptset string) error {
	// Get current promptset config
	promptsetConfig, err := a.manifestManager.GetPromptset(ctx, registry, promptset)
	if err != nil {
		return fmt.Errorf("promptset %s/%s not found in manifest", registry, promptset)
	}

	// Create install request with "latest" version to ignore constraints
	req := &InstallPromptsetRequest{
		Registry:  registry,
		Promptset: promptset,
		Version:   "latest", // This will ignore version constraints
		Include:   promptsetConfig.Include,
		Exclude:   promptsetConfig.Exclude,
		Sinks:     promptsetConfig.Sinks,
	}

	// Install with latest version
	return a.InstallPromptset(ctx, req)
}

// Info operations

func (a *ArmService) ShowPromptsetInfo(ctx context.Context, promptsets []string) error {
	if len(promptsets) == 0 {
		infos, err := a.listInstalledPromptsets(ctx)
		if err != nil {
			return err
		}
		a.ui.PromptsetInfoGrouped(infos, false)
		return nil
	}

	// Show info for specific promptsets
	var infos []*PromptsetInfo
	for _, promptsetArg := range promptsets {
		parts := strings.Split(promptsetArg, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid promptset format: %s (expected registry/promptset)", promptsetArg)
		}
		info, err := a.getPromptsetInfo(ctx, parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("failed to get promptset info for %s: %w", promptsetArg, err)
		}
		infos = append(infos, info)
	}

	detailed := len(promptsets) == 1
	a.ui.PromptsetInfoGrouped(infos, detailed)
	return nil
}

// getPromptsetInfo gets detailed information about a specific promptset
func (a *ArmService) getPromptsetInfo(ctx context.Context, registry, promptset string) (*PromptsetInfo, error) {
	// Get manifest entry
	manifestEntry, err := a.manifestManager.GetPromptset(ctx, registry, promptset)
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
		sink, exists := sinks[sinkName]
		if !exists {
			continue
		}
		installer := installer.NewInstaller(&sink)
		if err != nil {
			continue // Skip sinks with invalid configurations
		}
		installations, err := installer.ListInstalledPromptsets(ctx)
		if err != nil {
			continue
		}
		for _, installation := range installations {
			if installation.Registry == registry && installation.Promptset == promptset {
				installedPaths = append(installedPaths, installation.Path)
				if resolvedVersion == "" {
					resolvedVersion = installation.Version
				}
				break
			}
		}
	}

	return &PromptsetInfo{
		Registry: registry,
		Name:     promptset,
		Manifest: ui.PromptsetManifestInfo{
			Constraint: manifestEntry.Version,
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

// ShowPromptsetList shows a list of all installed promptsets
func (a *ArmService) ShowPromptsetList(ctx context.Context) error {
	promptsets, err := a.listInstalledPromptsets(ctx)
	if err != nil {
		return err
	}

	a.ui.PromptsetList(promptsets)
	return nil
}

// ShowPromptsetOutdated shows outdated promptsets
func (a *ArmService) ShowPromptsetOutdated(ctx context.Context, outputFormat string, noSpinner bool) error {
	if noSpinner || outputFormat == "json" {
		outdated, err := a.getOutdatedPromptsets(ctx)
		if err != nil {
			return err
		}
		// Convert to unified format
		packages := a.convertOutdatedPromptsetsToPackages(outdated)
		a.ui.OutdatedTable(packages, outputFormat)
	} else {
		finishChecking := a.ui.InstallStepWithSpinner("Checking for updates...")
		outdated, err := a.getOutdatedPromptsets(ctx)
		if err != nil {
			return err
		}
		finishChecking(fmt.Sprintf("Found %d outdated promptsets", len(outdated)))
		fmt.Println() // Add spacing between spinner and table
		// Convert to unified format
		packages := a.convertOutdatedPromptsetsToPackages(outdated)
		a.ui.OutdatedTable(packages, outputFormat)
	}
	return nil
}

// getOutdatedPromptsets returns a list of outdated promptsets
func (a *ArmService) getOutdatedPromptsets(ctx context.Context) ([]OutdatedPromptset, error) {
	lockEntries, err := a.lockFileManager.GetPromptsets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	var outdated []OutdatedPromptset
	for registryName, promptsets := range lockEntries {
		for promptsetName, lockEntry := range promptsets {
			// Get latest and wanted versions using the unified method
			latestVersion, wantedVersion, err := a.getLatestVersions(ctx, registryName, promptsetName)
			if err != nil {
				continue
			}

			// Only add if there's a newer version available
			if latestVersion != lockEntry.Version {
				promptsetInfo, err := a.getPromptsetInfo(ctx, registryName, promptsetName)
				if err != nil {
					continue
				}
				outdated = append(outdated, OutdatedPromptset{
					PromptsetInfo: promptsetInfo,
					Wanted:        wantedVersion,
					Latest:        latestVersion,
				})
			}
		}
	}

	return outdated, nil
}

func (a *ArmService) listInstalledPromptsets(ctx context.Context) ([]*PromptsetInfo, error) {
	lockEntries, err := a.lockFileManager.GetPromptsets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	var promptsets []*PromptsetInfo
	for registryName, promptsetMap := range lockEntries {
		for promptsetName := range promptsetMap {
			promptsetInfo, err := a.getPromptsetInfo(ctx, registryName, promptsetName)
			if err != nil {
				continue
			}
			promptsets = append(promptsets, promptsetInfo)
		}
	}

	// Sort by registry then promptset name
	sort.Slice(promptsets, func(i, j int) bool {
		if promptsets[i].Registry != promptsets[j].Registry {
			return promptsets[i].Registry < promptsets[j].Registry
		}
		return promptsets[i].Name < promptsets[j].Name
	})

	return promptsets, nil
}

// UpdateAllPromptsets updates all installed promptsets to their latest available versions within constraints
func (a *ArmService) UpdateAllPromptsets(ctx context.Context) error {
	manifestEntries, manifestErr := a.manifestManager.GetPromptsets(ctx)
	_, lockErr := a.lockFileManager.GetPromptsets(ctx)

	// Case: No manifest, no lockfile
	if manifestErr != nil && lockErr != nil {
		return fmt.Errorf("neither arm.json nor arm-lock.json found")
	}

	// Case: No manifest, lockfile exists
	if manifestErr != nil && lockErr == nil {
		return fmt.Errorf("arm.json not found")
	}

	// Case: Manifest exists - update within version constraints
	for registryName, promptsets := range manifestEntries {
		for promptsetName := range promptsets {
			if err := a.UpdatePromptset(ctx, registryName, promptsetName); err != nil {
				return err
			}
		}
	}

	return nil
}

// convertOutdatedPromptsetsToPackages converts OutdatedPromptset to OutdatedPackage format
func (a *ArmService) convertOutdatedPromptsetsToPackages(outdated []OutdatedPromptset) []ui.OutdatedPackage {
	packages := make([]ui.OutdatedPackage, len(outdated))
	for i, promptset := range outdated {
		packages[i] = ui.OutdatedPackage{
			Package:    fmt.Sprintf("%s/%s", promptset.PromptsetInfo.Registry, promptset.PromptsetInfo.Name),
			Type:       "promptset",
			Constraint: promptset.PromptsetInfo.Manifest.Constraint,
			Current:    promptset.PromptsetInfo.Installation.Version,
			Wanted:     promptset.Wanted,
			Latest:     promptset.Latest,
		}
	}
	return packages
}

// installPromptsetExactVersion installs a promptset from lockfile with exact version and checksum verification
func (a *ArmService) installPromptsetExactVersion(ctx context.Context, registryName, promptset string, lockEntry *lockfile.Entry) error {
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
	manifestEntry, err := a.manifestManager.GetPromptset(ctx, registryName, promptset)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	registryClient, err := registry.NewRegistry(registryName, registryConfig)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	resolvedVersion := types.Version{Version: lockEntry.Version, Display: lockEntry.Display}
	selector := types.ContentSelector{Include: manifestEntry.Include, Exclude: manifestEntry.Exclude}
	files, err := registryClient.GetContent(ctx, promptset, resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to get content: %w", err)
	}

	// Verify checksum for integrity
	if !lockfile.VerifyChecksum(files, lockEntry.Checksum) {
		return fmt.Errorf("checksum verification failed for %s/%s@%s", registryName, promptset, lockEntry.Version)
	}

	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	// Install to sinks specified in manifest entry
	for _, sinkName := range manifestEntry.Sinks {
		if sink, exists := sinks[sinkName]; exists {
			installer := installer.NewInstaller(&sink)
			if err != nil {
				return fmt.Errorf("failed to create installer for sink %s: %w", sinkName, err)
			}
			// Use display version for directory names
			if err := installer.InstallPromptset(ctx, registryName, promptset, lockEntry.Display, files); err != nil {
				return err
			}
		}
	}

	return nil
}

// cleanPreviousPromptsetInstallation removes previous promptset installation from sink directories
func (a *ArmService) cleanPreviousPromptsetInstallation(ctx context.Context, registry, promptset string) error {
	// Get current manifest entry to find previous sinks
	manifestEntry, err := a.manifestManager.GetPromptset(ctx, registry, promptset)
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
			if err != nil {
				continue // Skip sinks with invalid configurations
			}
			if err := installer.UninstallPromptset(ctx, registry, promptset); err != nil {
				// Continue on cleanup failure
				_ = err
			}
		}
	}

	return nil
}
