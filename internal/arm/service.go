package arm

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"sort"

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
	InstallRuleset(ctx context.Context, registry, ruleset, version string, include, exclude []string) error
	Install(ctx context.Context) error
	Uninstall(ctx context.Context, registry, ruleset string) error
	UpdateRuleset(ctx context.Context, registry, ruleset string) error
	Update(ctx context.Context) error
	Outdated(ctx context.Context) ([]OutdatedRuleset, error)
	List(ctx context.Context) ([]InstalledRuleset, error)
	Info(ctx context.Context, registry, ruleset string) (*RulesetInfo, error)
	InfoAll(ctx context.Context) ([]*RulesetInfo, error)
	Version() version.VersionInfo
}

// ArmService orchestrates all ARM operations.
type ArmService struct {
	configManager   config.Manager
	manifestManager manifest.Manager
	lockFileManager lockfile.Manager
	installer       installer.Installer
}

// NewArmService creates a new ARM service instance with all dependencies.
func NewArmService() *ArmService {
	return &ArmService{
		configManager:   config.NewFileManager(),
		manifestManager: manifest.NewFileManager(),
		lockFileManager: lockfile.NewFileManager(),
		installer:       installer.NewFileInstaller(),
	}
}

func (a *ArmService) InstallRuleset(ctx context.Context, registryName, ruleset, version string, include, exclude []string) error {
	// Normalize empty version to "latest"
	if version == "" {
		version = "latest"
	}

	// Expand shorthand constraints for storage
	version = expandVersionShorthand(version)

	// Validate registry exists in config
	registries, err := a.configManager.GetRegistries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get registries: %w", err)
	}
	registryConfig, exists := registries[registryName]
	if !exists {
		slog.ErrorContext(ctx, "Registry not found", "registry", registryName)
		return fmt.Errorf("registry %s not found", registryName)
	}

	// Create registry client
	registryClient, err := registry.NewRegistry(registryName, registryConfig)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	// Resolve version from registry (resolver expects "latest", not empty string)
	resolvedVersion, err := registryClient.ResolveVersion(ctx, version)
	if err != nil {
		return err
	}

	// Download ruleset files
	selector := types.ContentSelector{Include: include, Exclude: exclude}
	files, err := registryClient.GetContent(ctx, *resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to get content: %w", err)
	}

	// Update manifest
	manifestEntry := manifest.Entry{
		Version: version,
		Include: include,
		Exclude: exclude,
	}
	if err := a.manifestManager.CreateEntry(ctx, registryName, ruleset, manifestEntry); err != nil {
		if err := a.manifestManager.UpdateEntry(ctx, registryName, ruleset, manifestEntry); err != nil {
			return fmt.Errorf("failed to update manifest: %w", err)
		}
	}

	// Update lockfile
	lockEntry := &lockfile.Entry{
		URL:        registryConfig.URL,
		Type:       registryConfig.Type,
		Constraint: version,
		Resolved:   resolvedVersion.ID,
		Include:    include,
		Exclude:    exclude,
	}
	if err := a.lockFileManager.CreateEntry(ctx, registryName, ruleset, lockEntry); err != nil {
		if err := a.lockFileManager.UpdateEntry(ctx, registryName, ruleset, lockEntry); err != nil {
			return fmt.Errorf("failed to update lockfile: %w", err)
		}
	}

	// Install files to sink directories
	slog.InfoContext(ctx, "Installing ruleset", "registry", registryName, "ruleset", ruleset, "version", resolvedVersion.ID)
	sinks, err := a.configManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	rulesetKey := registryName + "/" + ruleset
	for _, sink := range sinks {
		if a.matchesSink(rulesetKey, sink) {
			for _, dir := range sink.Directories {
				if err := a.installer.Install(ctx, dir, ruleset, resolvedVersion.ID, files); err != nil {
					slog.ErrorContext(ctx, "Failed to install to directory", "dir", dir, "error", err)
					return err
				}
			}
		}
	}

	return nil
}

func (a *ArmService) Install(ctx context.Context) error {
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
			if err := a.InstallRuleset(ctx, registryName, rulesetName, entry.Version, entry.Include, entry.Exclude); err != nil {
				slog.ErrorContext(ctx, "Failed to install ruleset", "registry", registryName, "ruleset", rulesetName, "error", err)
				return err
			}
		}
	}

	return nil
}

func (a *ArmService) Uninstall(ctx context.Context, registry, ruleset string) error {
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
		if a.matchesSink(rulesetKey, sink) {
			for _, dir := range sink.Directories {
				if err := a.installer.Uninstall(ctx, dir, ruleset); err != nil {
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
	return a.InstallRuleset(ctx, registry, ruleset, manifestEntry.Version, manifestEntry.Include, manifestEntry.Exclude)
}

func (a *ArmService) Outdated(ctx context.Context) ([]OutdatedRuleset, error) {
	lockEntries, err := a.lockFileManager.GetEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	registryConfigs, err := a.configManager.GetRegistries(ctx)
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
			versions, err := registryClient.ListVersions(ctx)
			if err != nil || len(versions) == 0 {
				continue
			}
			latestVersion := versions[len(versions)-1] // Assume last version is latest

			wantedVersion, err := registryClient.ResolveVersion(ctx, lockEntry.Constraint)
			if err != nil {
				continue
			}

			if lockEntry.Resolved != latestVersion.ID || lockEntry.Resolved != wantedVersion.ID {
				outdated = append(outdated, OutdatedRuleset{
					Registry: registryName,
					Name:     rulesetName,
					Current:  lockEntry.Resolved,
					Wanted:   wantedVersion.ID,
					Latest:   latestVersion.ID,
				})
			}
		}
	}

	return outdated, nil
}

func (a *ArmService) List(ctx context.Context) ([]InstalledRuleset, error) {
	lockEntries, err := a.lockFileManager.GetEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockfile entries: %w", err)
	}

	var rulesets []InstalledRuleset
	for registryName, rulesetMap := range lockEntries {
		for rulesetName, entry := range rulesetMap {
			rulesets = append(rulesets, InstalledRuleset{
				Registry: registryName,
				Name:     rulesetName,
				Version:  entry.Resolved,
				Include:  entry.Include,
				Exclude:  entry.Exclude,
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

func (a *ArmService) Info(ctx context.Context, registry, ruleset string) (*RulesetInfo, error) {
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

	// Get registry config
	registries, err := a.configManager.GetRegistries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registries: %w", err)
	}
	registryConfig := registries[registry]

	// Get sinks and find installation paths
	sinks, err := a.configManager.GetSinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sinks: %w", err)
	}

	var installedPaths []string
	var sinkNames []string
	for sinkName, sink := range sinks {
		for _, dir := range sink.Directories {
			installations, err := a.installer.ListInstalled(ctx, dir)
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
		RegistryURL:    registryConfig.URL,
		RegistryType:   registryConfig.Type,
		Include:        manifestEntry.Include,
		Exclude:        manifestEntry.Exclude,
		InstalledPaths: installedPaths,
		Sinks:          sinkNames,
		Constraint:     lockEntry.Constraint,
		Resolved:       lockEntry.Resolved,
	}, nil
}

func (a *ArmService) Update(ctx context.Context) error {
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

func (a *ArmService) InfoAll(ctx context.Context) ([]*RulesetInfo, error) {
	installed, err := a.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list installed rulesets: %w", err)
	}

	var infos []*RulesetInfo
	for _, ruleset := range installed {
		info, err := a.Info(ctx, ruleset.Registry, ruleset.Name)
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
	registryConfig := config.RegistryConfig{
		URL:  lockEntry.URL,
		Type: lockEntry.Type,
	}

	registryClient, err := registry.NewRegistry(registryName, registryConfig)
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	resolvedVersion := &types.VersionRef{ID: lockEntry.Resolved}
	selector := types.ContentSelector{Include: lockEntry.Include, Exclude: lockEntry.Exclude}
	files, err := registryClient.GetContent(ctx, *resolvedVersion, selector)
	if err != nil {
		return fmt.Errorf("failed to get content: %w", err)
	}

	sinks, err := a.configManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	rulesetKey := registryName + "/" + ruleset
	for _, sink := range sinks {
		if a.matchesSink(rulesetKey, sink) {
			for _, dir := range sink.Directories {
				if err := a.installer.Install(ctx, dir, ruleset, lockEntry.Resolved, files); err != nil {
					slog.ErrorContext(ctx, "Failed to install exact version to directory", "dir", dir, "error", err)
					return err
				}
			}
		}
	}

	return nil
}

func (a *ArmService) matchesSink(rulesetKey string, sink config.SinkConfig) bool {
	// Check exclude patterns first
	for _, pattern := range sink.Exclude {
		if matched, _ := filepath.Match(pattern, rulesetKey); matched {
			return false
		}
	}

	// If no include patterns, allow all (that aren't excluded)
	if len(sink.Include) == 0 {
		return true
	}

	// Check include patterns
	for _, pattern := range sink.Include {
		if matched, _ := filepath.Match(pattern, rulesetKey); matched {
			return true
		}
	}

	return false
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
