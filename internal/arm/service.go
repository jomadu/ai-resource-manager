package arm

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/jomadu/ai-rules-manager/internal/urf"
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
	ListInstalledRulesets(ctx context.Context) ([]*RulesetInfo, error)
	GetRulesetInfo(ctx context.Context, registry, ruleset string) (*RulesetInfo, error)
	GetAllRulesetInfo(ctx context.Context) ([]*RulesetInfo, error)

	UpdateRulesetConfig(ctx context.Context, registry, ruleset, field, value string) error
	AddSink(ctx context.Context, name, directory, sinkType, layout, compileTarget string, force bool) error
	SyncRemovedSink(ctx context.Context, removedSink *manifest.SinkConfig) error
	CompileFiles(ctx context.Context, req *CompileRequest) (*CompileResult, error)
	Version() version.VersionInfo
}

// ArmService orchestrates all ARM operations.
type ArmService struct {
	manifestManager manifest.Manager
	lockFileManager lockfile.Manager
}

// NewArmService creates a new ARM service instance with all dependencies.
func NewArmService() *ArmService {
	return &ArmService{
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
	resolvedVersionResult, err := registryClient.ResolveVersion(ctx, req.Ruleset, versionStr)
	if err != nil {
		return fmt.Errorf("failed to resolve version: %w", err)
	}
	resolvedVersion := resolvedVersionResult.Version

	// Download content
	selector := types.ContentSelector{Include: req.Include, Exclude: req.Exclude}
	files, err := registryClient.GetContent(ctx, req.Ruleset, resolvedVersion, selector)
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
		Version:  manifestVersion,
		Priority: req.Priority,
		Include:  req.Include,
		Exclude:  req.Exclude,
		Sinks:    req.Sinks,
	}
	if err := a.manifestManager.CreateEntry(ctx, req.Registry, req.Ruleset, &manifestEntry); err != nil {
		if err := a.manifestManager.UpdateEntry(ctx, req.Registry, req.Ruleset, &manifestEntry); err != nil {
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

func (a *ArmService) cleanPreviousInstallation(ctx context.Context, registry, ruleset string) error {
	// Get current manifest entry to find previous sinks
	manifestEntry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
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
			if err := installer.Uninstall(ctx, registry, ruleset); err != nil {
				slog.WarnContext(ctx, "Failed to clean previous installation", "sink", sinkName, "error", err)
			}
		}
	}

	return nil
}

func (a *ArmService) installToSinks(ctx context.Context, req *InstallRequest, version types.Version, files []types.File) error {
	slog.InfoContext(ctx, "Installing ruleset", "registry", req.Registry, "ruleset", req.Ruleset, "version", version.Display)

	// First, remove from previous sink locations if this is a reinstall
	if err := a.cleanPreviousInstallation(ctx, req.Registry, req.Ruleset); err != nil {
		slog.WarnContext(ctx, "Failed to clean previous installation", "error", err)
	}

	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	// Validate that all requested sinks exist
	for _, sinkName := range req.Sinks {
		if _, exists := sinks[sinkName]; !exists {
			return fmt.Errorf("sink %s not configured", sinkName)
		}
	}

	// Install to explicitly specified sinks
	for _, sinkName := range req.Sinks {
		sink := sinks[sinkName]
		installer := installer.NewInstaller(&sink)
		if err := installer.Install(ctx, req.Registry, req.Ruleset, version.Display, req.Priority, files); err != nil {
			slog.ErrorContext(ctx, "Failed to install to sink", "sink", sinkName, "directory", sink.Directory, "error", err)
			return err
		}
	}

	slog.InfoContext(ctx, "Ruleset installed to sinks", "registry", req.Registry, "ruleset", req.Ruleset, "sinks", req.Sinks)
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
				Sinks:    entry.Sinks,
			}); err != nil {
				slog.ErrorContext(ctx, "Failed to install ruleset", "registry", registryName, "ruleset", rulesetName, "error", err)
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) UninstallRuleset(ctx context.Context, registry, ruleset string) error {
	// Get manifest entry to find target sinks
	manifestEntry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	// Remove installed files from sink directories
	slog.InfoContext(ctx, "Uninstalling ruleset", "registry", registry, "ruleset", ruleset)
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	for _, sinkName := range manifestEntry.Sinks {
		if sink, exists := sinks[sinkName]; exists {
			installer := installer.NewInstaller(&sink)
			if err := installer.Uninstall(ctx, registry, ruleset); err != nil {
				slog.ErrorContext(ctx, "Failed to uninstall from sink", "sink", sinkName, "directory", sink.Directory, "error", err)
				return err
			}
		}
	}

	// Remove from manifest
	if err := a.manifestManager.RemoveEntry(ctx, registry, ruleset); err != nil {
		return fmt.Errorf("failed to remove from manifest: %w", err)
	}

	// Remove from lockfile
	if err := a.lockFileManager.RemoveEntry(ctx, registry, ruleset); err != nil {
		return fmt.Errorf("failed to remove from lockfile: %w", err)
	}

	return nil
}

func (a *ArmService) UpdateRuleset(ctx context.Context, registryName, rulesetName string) error {
	// Get manifest entry for version constraint
	manifestEntry, err := a.manifestManager.GetEntry(ctx, registryName, rulesetName)
	if err != nil {
		return fmt.Errorf("failed to get manifest entry: %w", err)
	}

	// Resolve what version we should have
	registries, err := a.manifestManager.GetRawRegistries(ctx)
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
		installed, version, err := installer.IsInstalled(ctx, registryName, rulesetName)
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
		slog.InfoContext(ctx, "Installing ruleset (not currently installed)", "registry", registryName, "ruleset", rulesetName)
		return a.InstallRuleset(ctx, &InstallRequest{
			Registry: registryName,
			Ruleset:  rulesetName,
			Version:  manifestEntry.Version,
			Include:  manifestEntry.Include,
			Exclude:  manifestEntry.Exclude,
		})
	}

	// Check if installed version matches what we want
	if installedVersion == resolvedVersionResult.Version.Display {
		// Get lockfile entry to verify checksum
		currentLockEntry, err := a.lockFileManager.GetEntry(ctx, registryName, rulesetName)
		if err == nil {
			// Verify checksum to ensure integrity
			selector := types.ContentSelector{Include: manifestEntry.Include, Exclude: manifestEntry.Exclude}
			files, err := registryClient.GetContent(ctx, rulesetName, resolvedVersionResult.Version, selector)
			if err == nil && lockfile.VerifyChecksum(files, currentLockEntry.Checksum) {
				slog.InfoContext(ctx, "Ruleset already up to date", "registry", registryName, "ruleset", rulesetName, "version", installedVersion)
				return nil
			}
		}
		// If we can't verify checksum, fall through to reinstall
		slog.InfoContext(ctx, "Cannot verify integrity, reinstalling", "registry", registryName, "ruleset", rulesetName, "version", installedVersion)
	} else {
		slog.InfoContext(ctx, "Updating ruleset", "registry", registryName, "ruleset", rulesetName, "from", installedVersion, "to", resolvedVersionResult.Version.Display)
	}

	// Version changed or integrity check failed, proceed with update
	return a.InstallRuleset(ctx, &InstallRequest{
		Registry: registryName,
		Ruleset:  rulesetName,
		Version:  manifestEntry.Version,
		Include:  manifestEntry.Include,
		Exclude:  manifestEntry.Exclude,
		Sinks:    manifestEntry.Sinks,
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
			latestVersion, err := registryClient.ResolveVersion(ctx, rulesetName, "latest")
			if err != nil {
				continue
			}

			wantedVersion, err := registryClient.ResolveVersion(ctx, rulesetName, constraint)
			if err != nil {
				continue
			}

			if lockEntry.Version != latestVersion.Version.Version || lockEntry.Version != wantedVersion.Version.Version {
				rulesetInfo, err := a.GetRulesetInfo(ctx, registryName, rulesetName)
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

func (a *ArmService) ListInstalledRulesets(ctx context.Context) ([]*RulesetInfo, error) {
	lockEntries, err := a.lockFileManager.GetEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	var rulesets []*RulesetInfo
	for registryName, rulesetMap := range lockEntries {
		for rulesetName := range rulesetMap {
			rulesetInfo, err := a.GetRulesetInfo(ctx, registryName, rulesetName)
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

func (a *ArmService) GetRulesetInfo(ctx context.Context, registry, ruleset string) (*RulesetInfo, error) {
	// Get manifest entry
	manifestEntry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
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
			installations, err := installer.ListInstalled(ctx)
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
			Priority:   manifestEntry.Priority,
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
	return a.ListInstalledRulesets(ctx)
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
			if err := installer.Install(ctx, registryName, ruleset, lockEntry.Display, manifestEntry.Priority, files); err != nil {
				slog.ErrorContext(ctx, "Failed to install exact version to sink", "sink", sinkName, "directory", sink.Directory, "error", err)
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) SyncRemovedSink(ctx context.Context, removedSink *manifest.SinkConfig) error {
	// Scan removed sink directory to find installed rulesets
	installer := installer.NewInstaller(removedSink)
	installations, err := installer.ListInstalled(ctx)
	if err != nil {
		return nil // Skip directory that can't be scanned
	}

	// Uninstall all found rulesets from this directory
	for _, installation := range installations {
		if err := installer.Uninstall(ctx, installation.Registry, installation.Ruleset); err != nil {
			slog.ErrorContext(ctx, "Failed to uninstall from removed sink", "registry", installation.Registry, "ruleset", installation.Ruleset, "error", err)
		}
	}

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

func (a *ArmService) UpdateRulesetConfig(ctx context.Context, registry, ruleset, field, value string) error {
	// Get current manifest entry
	entry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
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
		entry.Priority = priority
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
	if err := a.manifestManager.UpdateEntry(ctx, registry, ruleset, entry); err != nil {
		return fmt.Errorf("failed to update manifest: %w", err)
	}

	// Trigger reinstall
	return a.InstallRuleset(ctx, &InstallRequest{
		Registry: registry,
		Ruleset:  ruleset,
		Version:  entry.Version,
		Priority: entry.Priority,
		Include:  entry.Include,
		Exclude:  entry.Exclude,
		Sinks:    entry.Sinks,
	})
}

func (a *ArmService) CompileFiles(ctx context.Context, req *CompileRequest) (*CompileResult, error) {
	// Input validation
	if req == nil {
		return nil, fmt.Errorf("compile request is required")
	}
	if len(req.Files) == 0 {
		return nil, fmt.Errorf("no files specified for compilation")
	}
	if len(req.Targets) == 0 {
		return nil, fmt.Errorf("no compilation targets specified")
	}
	if req.OutputDir == "" {
		return nil, fmt.Errorf("output directory is required")
	}

	// Initialize result
	result := &CompileResult{
		CompiledFiles: make([]CompiledFile, 0),
		Skipped:       make([]SkippedFile, 0),
		Errors:        make([]CompileError, 0),
		Stats: CompileStats{
			TargetStats: make(map[string]int),
		},
	}

	// Discover and load files
	files, err := a.discoverFiles(req.Files, req.Recursive, req.Include, req.Exclude)
	if err != nil {
		return nil, fmt.Errorf("failed to discover files: %w", err)
	}

	result.Stats.FilesProcessed = len(files)

	// Filter for valid URF files and handle validate-only mode
	parser := urf.NewParser()
	var urfFiles []types.File
	for _, file := range files {
		if parser.IsURF(&file) {
			if req.ValidateOnly {
				// Just validate, don't compile
				_, err := parser.Parse(&file)
				if err != nil {
					result.Errors = append(result.Errors, CompileError{
						FilePath: file.Path,
						Error:    fmt.Sprintf("URF validation failed: %v", err),
					})
					result.Stats.Errors++
					if req.FailFast {
						break
					}
				}
			} else {
				urfFiles = append(urfFiles, file)
			}
		} else {
			result.Skipped = append(result.Skipped, SkippedFile{
				Path:   file.Path,
				Reason: "not a valid URF file",
			})
			result.Stats.FilesSkipped++
		}
	}

	if req.ValidateOnly {
		slog.InfoContext(ctx, "URF validation completed",
			"files_processed", result.Stats.FilesProcessed,
			"errors", result.Stats.Errors)
		return result, nil
	}

	// Compile for each target
	for _, target := range req.Targets {
		targetErr := a.compileForTarget(ctx, urfFiles, target, req, result)
		if targetErr != nil && req.FailFast {
			return nil, targetErr
		}
	}

	slog.InfoContext(ctx, "Compile operation completed",
		"files_processed", result.Stats.FilesProcessed,
		"files_compiled", result.Stats.FilesCompiled,
		"errors", result.Stats.Errors)

	return result, nil
}

func (a *ArmService) Version() version.VersionInfo {
	return version.GetVersionInfo()
}

// discoverFiles finds and loads files from the given patterns
func (a *ArmService) discoverFiles(patterns []string, recursive bool, include, exclude []string) ([]types.File, error) {
	var allFiles []types.File
	seen := make(map[string]bool)

	// Create content selector for filtering
	selector := types.ContentSelector{
		Include: include,
		Exclude: exclude,
	}

	for _, pattern := range patterns {
		files, err := a.discoverFromPattern(pattern, recursive, selector)
		if err != nil {
			return nil, err
		}

		// Add files, avoiding duplicates
		for _, file := range files {
			if !seen[file.Path] {
				allFiles = append(allFiles, file)
				seen[file.Path] = true
			}
		}
	}

	return allFiles, nil
}

// discoverFromPattern discovers files from a single pattern
func (a *ArmService) discoverFromPattern(pattern string, recursive bool, selector types.ContentSelector) ([]types.File, error) {
	// Check if pattern is a directory
	if info, err := os.Stat(pattern); err == nil && info.IsDir() {
		return a.discoverFromDirectory(pattern, recursive, selector)
	}

	// Check if pattern is a single file
	if info, err := os.Stat(pattern); err == nil && !info.IsDir() {
		if selector.Matches(pattern) {
			content, err := os.ReadFile(pattern)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", pattern, err)
			}
			return []types.File{{
				Path:    pattern,
				Content: content,
				Size:    int64(len(content)),
			}}, nil
		}
		return nil, nil
	}

	// Handle glob pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
	}

	var files []types.File
	for _, match := range matches {
		if info, err := os.Stat(match); err == nil {
			if info.IsDir() {
				dirFiles, err := a.discoverFromDirectory(match, recursive, selector)
				if err != nil {
					return nil, err
				}
				files = append(files, dirFiles...)
			} else if selector.Matches(match) {
				content, err := os.ReadFile(match)
				if err != nil {
					return nil, fmt.Errorf("failed to read file %s: %w", match, err)
				}
				files = append(files, types.File{
					Path:    match,
					Content: content,
					Size:    int64(len(content)),
				})
			}
		}
	}

	return files, nil
}

// discoverFromDirectory discovers files from a directory
func (a *ArmService) discoverFromDirectory(dirPath string, recursive bool, selector types.ContentSelector) ([]types.File, error) {
	var files []types.File

	if !recursive {
		// Top-level only
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue // Skip subdirectories
			}

			filePath := filepath.Join(dirPath, entry.Name())
			if selector.Matches(filePath) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
				}
				files = append(files, types.File{
					Path:    filePath,
					Content: content,
					Size:    int64(len(content)),
				})
			}
		}
	} else {
		// Recursive traversal
		err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil // Continue traversal
			}

			if selector.Matches(path) {
				content, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", path, err)
				}
				files = append(files, types.File{
					Path:    path,
					Content: content,
					Size:    int64(len(content)),
				})
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("failed to walk directory %s: %w", dirPath, err)
		}
	}

	return files, nil
}

// compileForTarget compiles URF files for a specific target
func (a *ArmService) compileForTarget(ctx context.Context, files []types.File, target urf.CompileTarget, req *CompileRequest, result *CompileResult) error {
	// Create compiler for this target
	compiler, err := urf.NewCompiler(target)
	if err != nil {
		return fmt.Errorf("failed to create compiler for target %s: %w", target, err)
	}

	// Create target directory if multi-target
	targetDir := req.OutputDir
	if len(req.Targets) > 1 {
		targetDir = filepath.Join(req.OutputDir, string(target))
		if !req.DryRun {
			if err := os.MkdirAll(targetDir, 0o755); err != nil {
				return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
			}
		}
	}

	// Generate namespace
	namespace := req.Namespace
	if namespace == "" {
		namespace = "compile@local"
	}

	// Process each URF file
	for _, file := range files {
		// Compile URF file
		compiledFiles, err := compiler.Compile(namespace, &file)
		if err != nil {
			result.Errors = append(result.Errors, CompileError{
				FilePath: file.Path,
				Target:   string(target),
				Error:    err.Error(),
			})
			result.Stats.Errors++
			continue
		}

		// Write compiled files
		for _, compiledFile := range compiledFiles {
			outputPath := filepath.Join(targetDir, compiledFile.Path)

			if req.DryRun {
				slog.InfoContext(ctx, "Would compile file",
					"source", file.Path,
					"output", outputPath,
					"target", target)
			} else {
				// Check if file exists and handle force flag
				if _, err := os.Stat(outputPath); err == nil && !req.Force {
					result.Skipped = append(result.Skipped, SkippedFile{
						Path:   outputPath,
						Reason: "file exists (use --force to overwrite)",
					})
					result.Stats.FilesSkipped++
					continue
				}

				// Ensure directory exists
				if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
					result.Errors = append(result.Errors, CompileError{
						FilePath: file.Path,
						Target:   string(target),
						Error:    fmt.Sprintf("failed to create output directory: %v", err),
					})
					result.Stats.Errors++
					continue
				}

				// Write file
				if err := os.WriteFile(outputPath, compiledFile.Content, 0o644); err != nil {
					result.Errors = append(result.Errors, CompileError{
						FilePath: file.Path,
						Target:   string(target),
						Error:    fmt.Sprintf("failed to write output file: %v", err),
					})
					result.Stats.Errors++
					continue
				}

				if req.Verbose {
					slog.InfoContext(ctx, "Compiled file",
						"source", file.Path,
						"output", outputPath,
						"target", target)
				}
			}

			// Count successful compilation
			result.CompiledFiles = append(result.CompiledFiles, CompiledFile{
				SourcePath: file.Path,
				TargetPath: outputPath,
				Target:     target,
				RuleCount:  1, // TODO: Count actual rules
			})
			result.Stats.FilesCompiled++
			result.Stats.TargetStats[string(target)]++
		}
	}

	return nil
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
