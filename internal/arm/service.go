package arm

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"sort"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/jomadu/ai-rules-manager/internal/version"
)

// Service provides the main ARM functionality for managing AI rule rulesets.
type Service interface {
	InstallRuleset(ctx context.Context, req *InstallRequest) error
	InstallManifest(ctx context.Context) error
	UninstallRuleset(ctx context.Context, registry, ruleset string) error
	UpdateRuleset(ctx context.Context, registry, ruleset string) error
	UpdateAllRulesets(ctx context.Context) error
	GetOutdatedRulesets(ctx context.Context) ([]OutdatedRuleset, error)
	ListInstalledRulesets(ctx context.Context) ([]InstalledRuleset, error)
	GetRulesetInfo(ctx context.Context, registry, ruleset string) (*RulesetInfo, error)
	GetAllRulesetInfo(ctx context.Context) ([]*RulesetInfo, error)
	SyncSink(ctx context.Context, sinkName string, sink *config.SinkConfig) error
	SyncRemovedSink(ctx context.Context, removedSink *config.SinkConfig) error
	Version() version.VersionInfo
}

// ArmService orchestrates all ARM operations.
type ArmService struct {
	configManager   config.Manager
	manifestManager manifest.Manager
	lockFileManager lockfile.Manager
}

// NewArmService creates a new ARM service instance with all dependencies.
func NewArmService() *ArmService {
	return &ArmService{
		configManager:   config.NewFileManager(),
		manifestManager: manifest.NewFileManager(),
		lockFileManager: lockfile.NewFileManager(),
	}
}

func (a *ArmService) InstallRuleset(ctx context.Context, req *InstallRequest) error {

	// Load registries once and validate
	registries, err := a.manifestManager.GetRawRegistries(ctx)
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
	resolvedVersionResult, err := registryClient.ResolveVersion(ctx, versionStr)
	if err != nil {
		return fmt.Errorf("failed to resolve version: %w", err)
	}
	resolvedVersion := resolvedVersionResult.Version

	// Download content
	selector := types.ContentSelector{Include: req.Include, Exclude: req.Exclude}
	files, err := registryClient.GetContent(ctx, resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to download content: %w", err)
	}

	if err := a.updateTrackingFiles(ctx, req, resolvedVersion, files); err != nil {
		return fmt.Errorf("failed to update tracking files: %w", err)
	}

	return a.installToSinks(ctx, req, resolvedVersion, files)
}

func (a *ArmService) updateTrackingFiles(ctx context.Context, req *InstallRequest, version types.Version, files []types.File) error {
	// Store normalized version in manifest
	manifestVersion := req.Version
	if manifestVersion == "" {
		manifestVersion = "latest"
	}
	manifestVersion = expandVersionShorthand(manifestVersion)

	manifestEntry := manifest.Entry{
		Version: manifestVersion,
		Include: req.Include,
		Exclude: req.Exclude,
	}
	if err := a.manifestManager.CreateEntry(ctx, req.Registry, req.Ruleset, manifestEntry); err != nil {
		if err := a.manifestManager.UpdateEntry(ctx, req.Registry, req.Ruleset, manifestEntry); err != nil {
			return fmt.Errorf("failed to update manifest: %w", err)
		}
	}

	checksum := lockfile.GenerateChecksum(files)
	lockEntry := &lockfile.Entry{
		Version:  version.Version,
		Display:  version.Display,
		Checksum: checksum,
	}
	if err := a.lockFileManager.CreateEntry(ctx, req.Registry, req.Ruleset, lockEntry); err != nil {
		if err := a.lockFileManager.UpdateEntry(ctx, req.Registry, req.Ruleset, lockEntry); err != nil {
			return fmt.Errorf("failed to update lockfile: %w", err)
		}
	}

	return nil
}

func (a *ArmService) installToSinks(ctx context.Context, req *InstallRequest, version types.Version, files []types.File) error {
	slog.InfoContext(ctx, "Installing ruleset", "registry", req.Registry, "ruleset", req.Ruleset, "version", version.Display)

	sinks, err := a.configManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	rulesetKey := req.Registry + "/" + req.Ruleset
	var matchingSinkNames []string
	for sinkName, sink := range sinks {
		if a.matchesSink(rulesetKey, &sink) {
			matchingSinkNames = append(matchingSinkNames, sinkName)
			installer := installer.NewInstaller(&sink)
			for _, dir := range sink.Directories {
				if err := installer.Install(ctx, dir, req.Registry, req.Ruleset, version.Display, files); err != nil {
					slog.ErrorContext(ctx, "Failed to install to directory", "dir", dir, "error", err)
					return err
				}
			}
		}
	}

	if len(matchingSinkNames) == 0 {
		slog.WarnContext(ctx, "Ruleset not targeted by any sinks", "registry", req.Registry, "ruleset", req.Ruleset)
	} else {
		slog.InfoContext(ctx, "Ruleset targeted by sinks", "registry", req.Registry, "ruleset", req.Ruleset, "sinks", matchingSinkNames)
	}

	return nil
}

func (a *ArmService) InstallManifest(ctx context.Context) error {
	manifestEntries, manifestErr := a.manifestManager.GetEntries(ctx)
	lockEntries, lockErr := a.lockFileManager.GetEntries(ctx)

	// Case: No manifest, no lockfile
	if manifestErr != nil && lockErr != nil {
		slog.ErrorContext(ctx, "No manifest or lockfile found")
		return fmt.Errorf("neither arm.json nor arm-lock.json found")
	}

	// Case: No manifest, lockfile exists
	if manifestErr != nil && lockErr == nil {
		slog.ErrorContext(ctx, "Manifest file not found")
		return fmt.Errorf("arm.json not found")
	}

	// Case: Manifest exists, lockfile exists - use exact lockfile versions
	if manifestErr == nil && lockErr == nil {
		return a.installFromLockfile(ctx, lockEntries)
	}

	// Case: Manifest exists, no lockfile - resolve from manifest and create lockfile
	for registryName, rulesets := range manifestEntries {
		for rulesetName, entry := range rulesets {
			if err := a.InstallRuleset(ctx, &InstallRequest{
				Registry: registryName,
				Ruleset:  rulesetName,
				Version:  entry.Version,
				Include:  entry.Include,
				Exclude:  entry.Exclude,
			}); err != nil {
				slog.ErrorContext(ctx, "Failed to install ruleset", "registry", registryName, "ruleset", rulesetName, "error", err)
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) UninstallRuleset(ctx context.Context, registry, ruleset string) error {
	// Remove from manifest
	if err := a.manifestManager.RemoveEntry(ctx, registry, ruleset); err != nil {
		return fmt.Errorf("failed to remove from manifest: %w", err)
	}

	// Remove from lockfile
	if err := a.lockFileManager.RemoveEntry(ctx, registry, ruleset); err != nil {
		return fmt.Errorf("failed to remove from lockfile: %w", err)
	}

	// Remove installed files from sink directories
	slog.InfoContext(ctx, "Uninstalling ruleset", "registry", registry, "ruleset", ruleset)
	sinks, err := a.configManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	rulesetKey := registry + "/" + ruleset
	for _, sink := range sinks {
		if a.matchesSink(rulesetKey, &sink) {
			installer := installer.NewInstaller(&sink)
			for _, dir := range sink.Directories {
				if err := installer.Uninstall(ctx, dir, registry, ruleset); err != nil {
					slog.ErrorContext(ctx, "Failed to uninstall from directory", "dir", dir, "error", err)
					return err
				}
			}
		}
	}

	return nil
}

func (a *ArmService) UpdateRuleset(ctx context.Context, registry, ruleset string) error {
	manifestEntry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	slog.InfoContext(ctx, "Updating ruleset", "registry", registry, "ruleset", ruleset)
	return a.InstallRuleset(ctx, &InstallRequest{
		Registry: registry,
		Ruleset:  ruleset,
		Version:  manifestEntry.Version,
		Include:  manifestEntry.Include,
		Exclude:  manifestEntry.Exclude,
	})
}

func (a *ArmService) GetOutdatedRulesets(ctx context.Context) ([]OutdatedRuleset, error) {
	lockEntries, err := a.lockFileManager.GetEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	manifestEntries, err := a.manifestManager.GetEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest entries: %w", err)
	}

	registryConfigs, err := a.manifestManager.GetRawRegistries(ctx)
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
			latestVersion, err := registryClient.ResolveVersion(ctx, "latest")
			if err != nil {
				continue
			}

			wantedVersion, err := registryClient.ResolveVersion(ctx, constraint)
			if err != nil {
				continue
			}

			if lockEntry.Version != latestVersion.Version.Version || lockEntry.Version != wantedVersion.Version.Version {
				outdated = append(outdated, OutdatedRuleset{
					Registry:   registryName,
					Name:       rulesetName,
					Constraint: constraint,
					Current:    lockEntry.Display,
					Wanted:     wantedVersion.Version.Display,
					Latest:     latestVersion.Version.Display,
				})
			}
		}
	}

	return outdated, nil
}

func (a *ArmService) ListInstalledRulesets(ctx context.Context) ([]InstalledRuleset, error) {
	lockEntries, err := a.lockFileManager.GetEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	manifestEntries, err := a.manifestManager.GetEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest entries: %w", err)
	}

	var rulesets []InstalledRuleset
	for registryName, rulesetMap := range lockEntries {
		for rulesetName, lockEntry := range rulesetMap {
			// Get include/exclude and constraint from manifest
			var include, exclude []string
			var constraint string
			if manifestRegistry, exists := manifestEntries[registryName]; exists {
				if manifestEntry, exists := manifestRegistry[rulesetName]; exists {
					include = manifestEntry.Include
					exclude = manifestEntry.Exclude
					constraint = manifestEntry.Version
				}
			}

			rulesets = append(rulesets, InstalledRuleset{
				Registry:   registryName,
				Name:       rulesetName,
				Version:    lockEntry.Display,
				Constraint: constraint,
				Include:    include,
				Exclude:    exclude,
			})
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

func (a *ArmService) GetRulesetInfo(ctx context.Context, registry, ruleset string) (*RulesetInfo, error) {
	// Get lockfile entry
	lockEntry, err := a.lockFileManager.GetEntry(ctx, registry, ruleset)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entry: %w", err)
	}

	// Get manifest entry
	manifestEntry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest entry: %w", err)
	}

	// Get sinks and find installation paths
	sinks, err := a.configManager.GetSinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sinks: %w", err)
	}

	var installedPaths []string
	var sinkNames []string
	for sinkName, sink := range sinks {
		installer := installer.NewInstaller(&sink)
		for _, dir := range sink.Directories {
			installations, err := installer.ListInstalled(ctx, dir)
			if err != nil {
				continue
			}
			for _, installation := range installations {
				if installation.Ruleset == ruleset {
					installedPaths = append(installedPaths, installation.Path)
					sinkNames = append(sinkNames, sinkName)
					break
				}
			}
		}
	}

	return &RulesetInfo{
		Registry:       registry,
		Name:           ruleset,
		Include:        manifestEntry.Include,
		Exclude:        manifestEntry.Exclude,
		InstalledPaths: installedPaths,
		Sinks:          sinkNames,
		Constraint:     manifestEntry.Version,
		Resolved:       lockEntry.Display,
	}, nil
}

func (a *ArmService) UpdateAllRulesets(ctx context.Context) error {
	manifestEntries, manifestErr := a.manifestManager.GetEntries(ctx)
	_, lockErr := a.lockFileManager.GetEntries(ctx)

	// Case: No manifest, no lockfile
	if manifestErr != nil && lockErr != nil {
		slog.ErrorContext(ctx, "No manifest or lockfile found for update")
		return fmt.Errorf("neither arm.json nor arm-lock.json found")
	}

	// Case: No manifest, lockfile exists
	if manifestErr != nil && lockErr == nil {
		slog.ErrorContext(ctx, "Manifest file not found for update")
		return fmt.Errorf("arm.json not found")
	}

	// Case: Manifest exists - update within version constraints
	for registryName, rulesets := range manifestEntries {
		for rulesetName := range rulesets {
			if err := a.UpdateRuleset(ctx, registryName, rulesetName); err != nil {
				slog.ErrorContext(ctx, "Failed to update ruleset", "registry", registryName, "ruleset", rulesetName, "error", err)
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) GetAllRulesetInfo(ctx context.Context) ([]*RulesetInfo, error) {
	installed, err := a.ListInstalledRulesets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list installed rulesets: %w", err)
	}

	var infos []*RulesetInfo
	for i := range installed {
		ruleset := &installed[i]
		info, err := a.GetRulesetInfo(ctx, ruleset.Registry, ruleset.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get info for %s/%s: %w", ruleset.Registry, ruleset.Name, err)
		}
		infos = append(infos, info)
	}

	return infos, nil
}

func (a *ArmService) installFromLockfile(ctx context.Context, lockEntries map[string]map[string]lockfile.Entry) error {
	for registryName, rulesets := range lockEntries {
		for rulesetName, lockEntry := range rulesets {
			if err := a.installExactVersion(ctx, registryName, rulesetName, &lockEntry); err != nil {
				slog.ErrorContext(ctx, "Failed to install exact version", "registry", registryName, "ruleset", rulesetName, "error", err)
				return err
			}
		}
	}
	return nil
}

func (a *ArmService) installExactVersion(ctx context.Context, registryName, ruleset string, lockEntry *lockfile.Entry) error {
	// Get registry config from manifest
	registries, err := a.manifestManager.GetRawRegistries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registries: %w", err)
	}
	registryConfig, exists := registries[registryName]
	if !exists {
		return fmt.Errorf("registry %s not found in manifest", registryName)
	}

	// Get manifest entry for include/exclude patterns
	manifestEntry, err := a.manifestManager.GetEntry(ctx, registryName, ruleset)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	registryClient, err := registry.NewRegistry(registryName, registryConfig)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	resolvedVersion := types.Version{Version: lockEntry.Version, Display: lockEntry.Display}
	selector := types.ContentSelector{Include: manifestEntry.Include, Exclude: manifestEntry.Exclude}
	files, err := registryClient.GetContent(ctx, resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to get content: %w", err)
	}

	// Verify checksum for integrity
	if !lockfile.VerifyChecksum(files, lockEntry.Checksum) {
		return fmt.Errorf("checksum verification failed for %s/%s@%s", registryName, ruleset, lockEntry.Version)
	}

	sinks, err := a.configManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	rulesetKey := registryName + "/" + ruleset
	for _, sink := range sinks {
		if a.matchesSink(rulesetKey, &sink) {
			installer := installer.NewInstaller(&sink)
			for _, dir := range sink.Directories {
				// Use display version for directory names
				if err := installer.Install(ctx, dir, registryName, ruleset, lockEntry.Display, files); err != nil {
					slog.ErrorContext(ctx, "Failed to install exact version to directory", "dir", dir, "error", err)
					return err
				}
			}
		}
	}

	return nil
}

func (a *ArmService) matchesSink(rulesetKey string, sink *config.SinkConfig) bool {
	// Check exclude patterns first
	for _, pattern := range sink.Exclude {
		if matched, _ := doublestar.Match(pattern, rulesetKey); matched {
			return false
		}
	}

	// If no include patterns, allow all (that aren't excluded)
	if len(sink.Include) == 0 {
		return true
	}

	// Check include patterns
	for _, pattern := range sink.Include {
		if matched, _ := doublestar.Match(pattern, rulesetKey); matched {
			return true
		}
	}

	return false
}

func (a *ArmService) SyncSink(ctx context.Context, sinkName string, sink *config.SinkConfig) error {
	// Get manifest entries to determine what should be installed
	manifestEntries, err := a.manifestManager.GetEntries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get manifest entries: %w", err)
	}

	// Scan sink directories to discover what's actually installed
	installedRulesets := make(map[string]bool)
	for _, dir := range sink.Directories {
		installer := installer.NewInstaller(sink)
		installations, err := installer.ListInstalled(ctx, dir)
		if err != nil {
			continue // Skip directories that can't be scanned
		}
		for _, installation := range installations {
			installedRulesets[installation.Ruleset] = true
		}
	}

	// Find rulesets that should be installed for this sink
	for registryName, rulesets := range manifestEntries {
		for rulesetName, entry := range rulesets {
			rulesetKey := registryName + "/" + rulesetName
			if a.matchesSink(rulesetKey, sink) {
				if !installedRulesets[rulesetName] {
					// Install missing ruleset
					if err := a.InstallRuleset(ctx, &InstallRequest{
						Registry: registryName,
						Ruleset:  rulesetName,
						Version:  entry.Version,
						Include:  entry.Include,
						Exclude:  entry.Exclude,
					}); err != nil {
						slog.ErrorContext(ctx, "Failed to sync install ruleset", "registry", registryName, "ruleset", rulesetName, "error", err)
					}
				}
				delete(installedRulesets, rulesetName) // Mark as handled
			}
		}
	}

	// Remove rulesets that no longer match sink patterns
	for rulesetName := range installedRulesets {
		// Find registry for this ruleset from lockfile
		lockEntries, err := a.lockFileManager.GetEntries(ctx)
		if err != nil {
			continue
		}
		for registryName, rulesets := range lockEntries {
			if _, exists := rulesets[rulesetName]; exists {
				rulesetKey := registryName + "/" + rulesetName
				if !a.matchesSink(rulesetKey, sink) {
					// Uninstall from this sink only
					for _, dir := range sink.Directories {
						installer := installer.NewInstaller(sink)
						if err := installer.Uninstall(ctx, dir, registryName, rulesetName); err != nil {
							slog.ErrorContext(ctx, "Failed to sync uninstall ruleset", "registry", registryName, "ruleset", rulesetName, "error", err)
						}
					}
				}
				break
			}
		}
	}

	return nil
}

func (a *ArmService) SyncRemovedSink(ctx context.Context, removedSink *config.SinkConfig) error {
	// Scan removed sink directories to find installed rulesets
	for _, dir := range removedSink.Directories {
		installer := installer.NewInstaller(removedSink)
		installations, err := installer.ListInstalled(ctx, dir)
		if err != nil {
			continue // Skip directories that can't be scanned
		}

		// Uninstall all found rulesets from this directory
		for _, installation := range installations {
			// Extract registry from installation path or use lockfile lookup
			lockEntries, err := a.lockFileManager.GetEntries(ctx)
			if err != nil {
				continue
			}
			for registryName, rulesets := range lockEntries {
				if _, exists := rulesets[installation.Ruleset]; exists {
					if err := installer.Uninstall(ctx, dir, registryName, installation.Ruleset); err != nil {
						slog.ErrorContext(ctx, "Failed to uninstall from removed sink", "registry", registryName, "ruleset", installation.Ruleset, "error", err)
					}
					break
				}
			}
		}
	}

	return nil
}

func (a *ArmService) Version() version.VersionInfo {
	return version.GetVersionInfo()
}

// expandVersionShorthand expands npm-style version shorthands to proper semantic version constraints.
// "1" -> "^1.0.0", "1.2" -> "^1.2.0"
func expandVersionShorthand(constraint string) string {
	// Match pure major version (e.g., "1")
	if matched, _ := regexp.MatchString(`^\d+$`, constraint); matched {
		return "^" + constraint + ".0.0"
	}
	// Match major.minor version (e.g., "1.2")
	if matched, _ := regexp.MatchString(`^\d+\.\d+$`, constraint); matched {
		return "^" + constraint + ".0"
	}
	return constraint
}
