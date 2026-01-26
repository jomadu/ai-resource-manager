package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/filetype"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/arm/parser"
	"github.com/jomadu/ai-resource-manager/internal/arm/registry"
	"github.com/jomadu/ai-resource-manager/internal/arm/sink"
	"github.com/jomadu/ai-resource-manager/internal/arm/storage"
)

type DependencyInfo struct {
	Installation sink.PackageInstallation
	LockInfo     packagelockfile.DependencyLockConfig
	Config       map[string]interface{}
	Version      string // Installed version from lockfile
}

type OutdatedDependency struct {
	Current    core.PackageMetadata
	Constraint string
	Wanted     core.PackageMetadata
	Latest     core.PackageMetadata
}

// ArmService handles all ARM operations
type ArmService struct {
	manifestMgr     manifest.Manager
	lockfileMgr     packagelockfile.Manager
	registryFactory registry.Factory
}

// NewArmService creates a new ARM service
func NewArmService(manifestMgr manifest.Manager, lockfileMgr packagelockfile.Manager, registryFactory registry.Factory) *ArmService {
	if registryFactory == nil {
		registryFactory = &registry.DefaultFactory{}
	}
	return &ArmService{
		manifestMgr:     manifestMgr,
		lockfileMgr:     lockfileMgr,
		registryFactory: registryFactory,
	}
}

// ---------------------------------------------
// Registry Management (Git, GitLab, Cloudsmith)
// ---------------------------------------------

// AddGitRegistry adds a Git registry
func (s *ArmService) AddGitRegistry(ctx context.Context, name, url string, branches []string, force bool) error {
	registries, err := s.manifestMgr.GetAllRegistriesConfig(ctx)
	if err != nil {
		return err
	}

	if _, exists := registries[name]; !force && exists {
		return errors.New("registry already exists")
	}

	config := manifest.GitRegistryConfig{
		URL:      url,
		Branches: branches,
	}

	return s.manifestMgr.UpsertGitRegistryConfig(ctx, name, config)
}

// AddGitLabRegistry adds a GitLab registry
func (s *ArmService) AddGitLabRegistry(ctx context.Context, name, url, projectID, groupID, apiVersion string, force bool) error {
	registries, err := s.manifestMgr.GetAllRegistriesConfig(ctx)
	if err != nil {
		return err
	}

	if _, exists := registries[name]; !force && exists {
		return errors.New("registry already exists")
	}

	config := manifest.GitLabRegistryConfig{
		URL:        url,
		ProjectID:  projectID,
		GroupID:    groupID,
		APIVersion: apiVersion,
	}

	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, &config)
}

// AddCloudsmithRegistry adds a Cloudsmith registry
func (s *ArmService) AddCloudsmithRegistry(ctx context.Context, name, url, owner, repository string, force bool) error {
	registries, err := s.manifestMgr.GetAllRegistriesConfig(ctx)
	if err != nil {
		return err
	}

	if _, exists := registries[name]; !force && exists {
		return errors.New("registry already exists")
	}

	config := manifest.CloudsmithRegistryConfig{
		URL:        url,
		Owner:      owner,
		Repository: repository,
	}

	return s.manifestMgr.UpsertCloudsmithRegistryConfig(ctx, name, config)
}

// RemoveRegistry removes a registry
func (s *ArmService) RemoveRegistry(ctx context.Context, name string) error {
	return s.manifestMgr.RemoveRegistryConfig(ctx, name)
}

// SetRegistryName sets registry name
func (s *ArmService) SetRegistryName(ctx context.Context, name, newName string) error {
	return s.manifestMgr.UpdateRegistryConfigName(ctx, name, newName)
}

// SetRegistryURL sets registry URL
func (s *ArmService) SetRegistryURL(ctx context.Context, name, url string) error {
	reg, err := s.manifestMgr.GetRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	reg["url"] = url
	return s.manifestMgr.UpsertRegistryConfig(ctx, name, reg)
}

// SetGitRegistryBranches sets Git registry branches
func (s *ArmService) SetGitRegistryBranches(ctx context.Context, name string, branches []string) error {
	config, err := s.manifestMgr.GetGitRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.Branches = branches
	return s.manifestMgr.UpsertGitRegistryConfig(ctx, name, config)
}

// SetGitLabRegistryProjectID sets GitLab registry project ID
func (s *ArmService) SetGitLabRegistryProjectID(ctx context.Context, name, projectID string) error {
	config, err := s.manifestMgr.GetGitLabRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.ProjectID = projectID
	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, &config)
}

// SetGitLabRegistryGroupID sets GitLab registry group ID
func (s *ArmService) SetGitLabRegistryGroupID(ctx context.Context, name, groupID string) error {
	config, err := s.manifestMgr.GetGitLabRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.GroupID = groupID
	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, &config)
}

// SetGitLabRegistryAPIVersion sets GitLab registry API version
func (s *ArmService) SetGitLabRegistryAPIVersion(ctx context.Context, name, apiVersion string) error {
	config, err := s.manifestMgr.GetGitLabRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.APIVersion = apiVersion
	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, &config)
}

// SetCloudsmithRegistryOwner sets Cloudsmith registry owner
func (s *ArmService) SetCloudsmithRegistryOwner(ctx context.Context, name, owner string) error {
	config, err := s.manifestMgr.GetCloudsmithRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.Owner = owner
	return s.manifestMgr.UpsertCloudsmithRegistryConfig(ctx, name, config)
}

// SetCloudsmithRegistryRepository sets Cloudsmith registry repository
func (s *ArmService) SetCloudsmithRegistryRepository(ctx context.Context, name, repository string) error {
	config, err := s.manifestMgr.GetCloudsmithRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.Repository = repository
	return s.manifestMgr.UpsertCloudsmithRegistryConfig(ctx, name, config)
}

// GetRegistryConfig gets registry configuration
func (s *ArmService) GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error) {
	return s.manifestMgr.GetRegistryConfig(ctx, name)
}

// GetAllRegistriesConfig gets all registry configurations
func (s *ArmService) GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	return s.manifestMgr.GetAllRegistriesConfig(ctx)
}

// ---------------
// Sink Management
// ---------------

// AddSink adds a sink
func (s *ArmService) AddSink(ctx context.Context, name, directory string, tool compiler.Tool, force bool) error {
	sinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	if _, exists := sinks[name]; !force && exists {
		return errors.New("sink already exists")
	}

	return s.manifestMgr.UpsertSinkConfig(ctx, name, manifest.SinkConfig{
		Directory: directory,
		Tool:      tool,
	})
}

// RemoveSink removes a sink
func (s *ArmService) RemoveSink(ctx context.Context, name string) error {
	return s.manifestMgr.RemoveSinkConfig(ctx, name)
}

// GetSinkConfig gets sink configuration
func (s *ArmService) GetSinkConfig(ctx context.Context, name string) (manifest.SinkConfig, error) {
	return s.manifestMgr.GetSinkConfig(ctx, name)
}

// GetAllSinkConfigs gets all sink configurations
func (s *ArmService) GetAllSinkConfigs(ctx context.Context) (map[string]manifest.SinkConfig, error) {
	return s.manifestMgr.GetAllSinksConfig(ctx)
}

// SetSinkName sets sink name
func (s *ArmService) SetSinkName(ctx context.Context, name, newName string) error {
	return s.manifestMgr.UpdateSinkConfigName(ctx, name, newName)
}

// SetSinkDirectory sets sink directory
func (s *ArmService) SetSinkDirectory(ctx context.Context, name, directory string) error {
	sink, err := s.manifestMgr.GetSinkConfig(ctx, name)
	if err != nil {
		return err
	}

	sink.Directory = directory
	return s.manifestMgr.UpsertSinkConfig(ctx, name, sink)
}

// SetSinkTool sets sink tool
func (s *ArmService) SetSinkTool(ctx context.Context, name string, tool compiler.Tool) error {
	sink, err := s.manifestMgr.GetSinkConfig(ctx, name)
	if err != nil {
		return err
	}

	sink.Tool = tool
	return s.manifestMgr.UpsertSinkConfig(ctx, name, sink)
}

// ---------------------
// Dependency Management
// ---------------------

// GetAllDependenciesConfig gets all dependency configurations
func (s *ArmService) GetAllDependenciesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	return s.manifestMgr.GetAllDependenciesConfig(ctx)
}

// InstallAll installs all dependencies
func (s *ArmService) InstallAll(ctx context.Context) error {
	rulesets, err := s.manifestMgr.GetAllRulesetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	promptsets, err := s.manifestMgr.GetAllPromptsetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	for key, rulesetCfg := range rulesets {
		registryName, packageName := manifest.ParseDependencyKey(key)
		if err := s.InstallRuleset(ctx, registryName, packageName, rulesetCfg.Version, rulesetCfg.Priority, rulesetCfg.Include, rulesetCfg.Exclude, rulesetCfg.Sinks); err != nil {
			return err
		}
	}

	for key, promptsetCfg := range promptsets {
		registryName, packageName := manifest.ParseDependencyKey(key)
		if err := s.InstallPromptset(ctx, registryName, packageName, promptsetCfg.Version, promptsetCfg.Include, promptsetCfg.Exclude, promptsetCfg.Sinks); err != nil {
			return err
		}
	}

	return nil
}

// resolveAndFetchPackage validates registry/sinks, resolves version, and fetches package
func (s *ArmService) resolveAndFetchPackage(ctx context.Context, registryName, packageName, version string, include, exclude, sinks []string) (pkg *core.Package, resolvedVersion string, sinkConfigs map[string]manifest.SinkConfig, err error) {
	regConfig, err := s.manifestMgr.GetRegistryConfig(ctx, registryName)
	if err != nil {
		return nil, "", nil, err
	}

	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, sinkName := range sinks {
		if _, exists := allSinks[sinkName]; !exists {
			return nil, "", nil, fmt.Errorf("sink does not exist: %s", sinkName)
		}
	}

	reg, err := s.registryFactory.CreateRegistry(registryName, regConfig)
	if err != nil {
		return nil, "", nil, err
	}

	availableVersions, err := reg.ListPackageVersions(ctx, packageName)
	if err != nil {
		return nil, "", nil, err
	}

	resolvedVer, err := core.ResolveVersion(version, availableVersions)
	if err != nil {
		return nil, "", nil, err
	}

	pkg, err = reg.GetPackage(ctx, packageName, &resolvedVer, include, exclude)
	if err != nil {
		return nil, "", nil, err
	}

	// Verify integrity if package is already locked
	lockedInfo, err := s.lockfileMgr.GetDependencyLock(ctx, registryName, packageName, resolvedVer.Version)
	if err == nil && lockedInfo != nil && lockedInfo.Integrity != "" {
		// Lock exists with integrity - verify it matches
		if pkg.Integrity != lockedInfo.Integrity {
			return nil, "", nil, fmt.Errorf("integrity verification failed for %s/%s@%s\n  Expected: %s\n  Got:      %s\n\nThis indicates the package has been modified since it was locked.\nTo resolve:\n  1. If you trust the new package: delete arm-lock.json and reinstall\n  2. If you suspect tampering: investigate the package source",
				registryName, packageName, resolvedVer.Version,
				lockedInfo.Integrity, pkg.Integrity)
		}
	}
	// If lock doesn't exist or has no integrity field, skip verification (backwards compatibility)

	return pkg, resolvedVer.Version, allSinks, nil
}

// InstallRuleset installs a ruleset
func (s *ArmService) InstallRuleset(ctx context.Context, registryName, ruleset, version string, priority int, include, exclude, sinks []string) error {
	pkg, resolvedVersion, allSinks, err := s.resolveAndFetchPackage(ctx, registryName, ruleset, version, include, exclude, sinks)
	if err != nil {
		return err
	}

	depConfig := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Version: version,
			Sinks:   sinks,
			Include: include,
			Exclude: exclude,
		},
		Priority: priority,
	}

	if err := s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registryName, ruleset, &depConfig); err != nil {
		return err
	}

	if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, ruleset, resolvedVersion, &packagelockfile.DependencyLockConfig{
		Integrity: pkg.Integrity,
	}); err != nil {
		return err
	}

	for _, sinkName := range sinks {
		sinkConfig := allSinks[sinkName]
		sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
		if err := sinkMgr.InstallRuleset(pkg, priority); err != nil {
			return err
		}
	}

	return nil
}

// InstallPromptset installs a promptset
func (s *ArmService) InstallPromptset(ctx context.Context, registryName, promptset, version string, include, exclude, sinks []string) error {
	pkg, resolvedVersion, allSinks, err := s.resolveAndFetchPackage(ctx, registryName, promptset, version, include, exclude, sinks)
	if err != nil {
		return err
	}

	depConfig := manifest.PromptsetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Version: version,
			Sinks:   sinks,
			Include: include,
			Exclude: exclude,
		},
	}

	if err := s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registryName, promptset, &depConfig); err != nil {
		return err
	}

	if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, promptset, resolvedVersion, &packagelockfile.DependencyLockConfig{
		Integrity: pkg.Integrity,
	}); err != nil {
		return err
	}

	for _, sinkName := range sinks {
		sinkConfig := allSinks[sinkName]
		sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
		if err := sinkMgr.InstallPromptset(pkg); err != nil {
			return err
		}
	}

	return nil
}

// UninstallAll uninstalls all dependencies
func (s *ArmService) UninstallAll(ctx context.Context) error {
	rulesets, err := s.manifestMgr.GetAllRulesetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	promptsets, err := s.manifestMgr.GetAllPromptsetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	for key, rulesetConfig := range rulesets {
		registryName, packageName := manifest.ParseDependencyKey(key)

		for _, sinkName := range rulesetConfig.Sinks {
			sinkConfig, exists := allSinks[sinkName]
			if !exists {
				continue
			}

			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if err := sinkMgr.Uninstall(registryName, packageName); err != nil {
				return err
			}
		}

		if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
			return err
		}

		if err := s.manifestMgr.RemoveDependencyConfig(ctx, registryName, packageName); err != nil {
			return err
		}
	}

	for key, promptsetConfig := range promptsets {
		registryName, packageName := manifest.ParseDependencyKey(key)

		for _, sinkName := range promptsetConfig.Sinks {
			sinkConfig, exists := allSinks[sinkName]
			if !exists {
				continue
			}

			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if err := sinkMgr.Uninstall(registryName, packageName); err != nil {
				return err
			}
		}

		if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
			return err
		}

		if err := s.manifestMgr.RemoveDependencyConfig(ctx, registryName, packageName); err != nil {
			return err
		}
	}

	return nil
}

// UninstallPackages uninstalls specific packages
func (s *ArmService) UninstallPackages(ctx context.Context, packages []string) error {
	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	successCount := 0
	var lastErr error

	for _, pkg := range packages {
		registryName, packageName := manifest.ParseDependencyKey(pkg)
		if registryName == "" || packageName == "" {
			fmt.Fprintf(os.Stderr, "Warning: invalid package format '%s', expected registry/package\n", pkg)
			lastErr = fmt.Errorf("invalid package format: %s", pkg)
			continue
		}

		depConfig, err := s.manifestMgr.GetDependencyConfig(ctx, registryName, packageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: package not found '%s'\n", pkg)
			lastErr = fmt.Errorf("package not found: %s", pkg)
			continue
		}

		depType, ok := depConfig["type"].(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: invalid package type for '%s'\n", pkg)
			lastErr = fmt.Errorf("invalid package type: %s", pkg)
			continue
		}

		var sinks []string

		switch depType {
		case "ruleset":
			rulesetConfig, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				lastErr = err
				continue
			}
			sinks = rulesetConfig.Sinks
		case "promptset":
			promptsetConfig, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				lastErr = err
				continue
			}
			sinks = promptsetConfig.Sinks
		default:
			fmt.Fprintf(os.Stderr, "Warning: unknown package type '%s' for '%s'\n", depType, pkg)
			lastErr = fmt.Errorf("unknown package type: %s", depType)
			continue
		}

		for _, sinkName := range sinks {
			sinkConfig, exists := allSinks[sinkName]
			if !exists {
				continue
			}

			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if err := sinkMgr.Uninstall(registryName, packageName); err != nil {
				lastErr = err
				continue
			}
		}

		if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
			lastErr = err
			continue
		}

		if err := s.manifestMgr.RemoveDependencyConfig(ctx, registryName, packageName); err != nil {
			lastErr = err
			continue
		}

		successCount++
	}

	if successCount == 0 && lastErr != nil {
		return lastErr
	}

	return nil
}

// UpdatePackages updates specific packages
func (s *ArmService) UpdatePackages(ctx context.Context, packages []string) error {
	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	lockFile, err := s.lockfileMgr.GetLockFile(ctx)
	if err != nil {
		return err
	}

	successCount := 0
	var lastErr error

	for _, pkg := range packages {
		registryName, packageName := manifest.ParseDependencyKey(pkg)
		if registryName == "" || packageName == "" {
			fmt.Fprintf(os.Stderr, "Warning: invalid package format '%s', expected registry/package\n", pkg)
			lastErr = fmt.Errorf("invalid package format: %s", pkg)
			continue
		}

		depConfig, err := s.manifestMgr.GetDependencyConfig(ctx, registryName, packageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: package not found '%s'\n", pkg)
			lastErr = fmt.Errorf("package not found: %s", pkg)
			continue
		}

		depType, ok := depConfig["type"].(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: invalid package type for '%s'\n", pkg)
			lastErr = fmt.Errorf("invalid package type: %s", pkg)
			continue
		}

		var sinks []string
		var version string
		var include []string
		var exclude []string
		var priority int

		switch depType {
		case "ruleset":
			rulesetConfig, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				lastErr = err
				continue
			}
			sinks = rulesetConfig.Sinks
			version = rulesetConfig.Version
			include = rulesetConfig.Include
			exclude = rulesetConfig.Exclude
			priority = rulesetConfig.Priority
		case "promptset":
			promptsetConfig, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				lastErr = err
				continue
			}
			sinks = promptsetConfig.Sinks
			version = promptsetConfig.Version
			include = promptsetConfig.Include
			exclude = promptsetConfig.Exclude
		default:
			fmt.Fprintf(os.Stderr, "Warning: unknown package type '%s' for '%s'\n", depType, pkg)
			lastErr = fmt.Errorf("unknown package type: %s", depType)
			continue
		}

		oldVersion := s.getOldVersionFromLock(lockFile, pkg)
		newVersion, fetchedPkg, err := s.resolveAndFetchUpdate(ctx, registryName, packageName, version, include, exclude)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to update '%s': %v\n", pkg, err)
			lastErr = err
			continue
		}

		if oldVersion == newVersion {
			successCount++
			continue
		}

		if oldVersion != "" {
			if err := s.uninstallFromSinks(sinks, allSinks, registryName, packageName); err != nil {
				lastErr = err
				continue
			}
			if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
				lastErr = err
				continue
			}
		}

		for _, sinkName := range sinks {
			sinkConfig := allSinks[sinkName]
			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if depType == "ruleset" {
				if err := sinkMgr.InstallRuleset(fetchedPkg, priority); err != nil {
					lastErr = err
					continue
				}
			} else {
				if err := sinkMgr.InstallPromptset(fetchedPkg); err != nil {
					lastErr = err
					continue
				}
			}
		}

		if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, packageName, newVersion, &packagelockfile.DependencyLockConfig{
			Integrity: fetchedPkg.Integrity,
		}); err != nil {
			lastErr = err
			continue
		}

		successCount++
	}

	if successCount == 0 && lastErr != nil {
		return lastErr
	}

	return nil
}

// UpdateAll updates all dependencies
func (s *ArmService) UpdateAll(ctx context.Context) error {
	rulesets, err := s.manifestMgr.GetAllRulesetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	promptsets, err := s.manifestMgr.GetAllPromptsetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	lockFile, err := s.lockfileMgr.GetLockFile(ctx)
	if err != nil {
		return err
	}

	for key, rulesetConfig := range rulesets {
		registryName, packageName := manifest.ParseDependencyKey(key)

		oldVersion := s.getOldVersionFromLock(lockFile, key)
		newVersion, pkg, err := s.resolveAndFetchUpdate(ctx, registryName, packageName, rulesetConfig.Version, rulesetConfig.Include, rulesetConfig.Exclude)
		if err != nil {
			return err
		}

		if oldVersion == newVersion {
			continue
		}

		if oldVersion != "" {
			if err := s.uninstallFromSinks(rulesetConfig.Sinks, allSinks, registryName, packageName); err != nil {
				return err
			}
			if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
				return err
			}
		}

		for _, sinkName := range rulesetConfig.Sinks {
			sinkConfig := allSinks[sinkName]
			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if err := sinkMgr.InstallRuleset(pkg, rulesetConfig.Priority); err != nil {
				return err
			}
		}

		if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, packageName, newVersion, &packagelockfile.DependencyLockConfig{
			Integrity: pkg.Integrity,
		}); err != nil {
			return err
		}
	}

	for key, promptsetConfig := range promptsets {
		registryName, packageName := manifest.ParseDependencyKey(key)

		oldVersion := s.getOldVersionFromLock(lockFile, key)
		newVersion, pkg, err := s.resolveAndFetchUpdate(ctx, registryName, packageName, promptsetConfig.Version, promptsetConfig.Include, promptsetConfig.Exclude)
		if err != nil {
			return err
		}

		if oldVersion == newVersion {
			continue
		}

		if oldVersion != "" {
			if err := s.uninstallFromSinks(promptsetConfig.Sinks, allSinks, registryName, packageName); err != nil {
				return err
			}
			if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
				return err
			}
		}

		for _, sinkName := range promptsetConfig.Sinks {
			sinkConfig := allSinks[sinkName]
			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if err := sinkMgr.InstallPromptset(pkg); err != nil {
				return err
			}
		}

		if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, packageName, newVersion, &packagelockfile.DependencyLockConfig{
			Integrity: pkg.Integrity,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *ArmService) getOldVersionFromLock(lockFile *packagelockfile.LockFile, key string) string {
	if lockFile == nil || lockFile.Dependencies == nil {
		return ""
	}

	for lockKey := range lockFile.Dependencies {
		if strings.HasPrefix(lockKey, key+"@") {
			return lockKey[len(key)+1:]
		}
	}
	return ""
}

func (s *ArmService) resolveAndFetchUpdate(ctx context.Context, registryName, packageName, version string, include, exclude []string) (string, *core.Package, error) {
	regConfig, err := s.manifestMgr.GetRegistryConfig(ctx, registryName)
	if err != nil {
		return "", nil, err
	}

	reg, err := s.registryFactory.CreateRegistry(registryName, regConfig)
	if err != nil {
		return "", nil, err
	}

	availableVersions, err := reg.ListPackageVersions(ctx, packageName)
	if err != nil {
		return "", nil, err
	}

	resolvedVersion, err := core.ResolveVersion(version, availableVersions)
	if err != nil {
		return "", nil, err
	}

	pkg, err := reg.GetPackage(ctx, packageName, &resolvedVersion, include, exclude)
	if err != nil {
		return "", nil, err
	}

	return resolvedVersion.Version, pkg, nil
}

func (s *ArmService) uninstallFromSinks(sinkNames []string, allSinks map[string]manifest.SinkConfig, registryName, packageName string) error {
	for _, sinkName := range sinkNames {
		sinkConfig, exists := allSinks[sinkName]
		if !exists {
			continue
		}

		sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
		if err := sinkMgr.Uninstall(registryName, packageName); err != nil {
			return err
		}
	}
	return nil
}

// UpgradeAll upgrades all dependencies to latest versions
func (s *ArmService) UpgradeAll(ctx context.Context) error {
	rulesets, err := s.manifestMgr.GetAllRulesetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	promptsets, err := s.manifestMgr.GetAllPromptsetDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	lockFile, err := s.lockfileMgr.GetLockFile(ctx)
	if err != nil {
		return err
	}

	for key, rulesetConfig := range rulesets {
		registryName, packageName := manifest.ParseDependencyKey(key)

		oldVersion := s.getOldVersionFromLock(lockFile, key)
		latestVersion, pkg, err := s.fetchLatest(ctx, registryName, packageName, rulesetConfig.Include, rulesetConfig.Exclude)
		if err != nil {
			return err
		}

		if oldVersion == latestVersion {
			continue
		}

		if oldVersion != "" {
			if err := s.uninstallFromSinks(rulesetConfig.Sinks, allSinks, registryName, packageName); err != nil {
				return err
			}
			if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
				return err
			}
		}

		for _, sinkName := range rulesetConfig.Sinks {
			sinkConfig := allSinks[sinkName]
			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if err := sinkMgr.InstallRuleset(pkg, rulesetConfig.Priority); err != nil {
				return err
			}
		}

		if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, packageName, latestVersion, &packagelockfile.DependencyLockConfig{
			Integrity: pkg.Integrity,
		}); err != nil {
			return err
		}

		newConstraint := fmt.Sprintf("^%d.0.0", pkg.Metadata.Version.Major)
		rulesetConfig.Version = newConstraint
		if err := s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registryName, packageName, rulesetConfig); err != nil {
			return err
		}
	}

	for key, promptsetConfig := range promptsets {
		registryName, packageName := manifest.ParseDependencyKey(key)

		oldVersion := s.getOldVersionFromLock(lockFile, key)
		latestVersion, pkg, err := s.fetchLatest(ctx, registryName, packageName, promptsetConfig.Include, promptsetConfig.Exclude)
		if err != nil {
			return err
		}

		if oldVersion == latestVersion {
			continue
		}

		if oldVersion != "" {
			if err := s.uninstallFromSinks(promptsetConfig.Sinks, allSinks, registryName, packageName); err != nil {
				return err
			}
			if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
				return err
			}
		}

		for _, sinkName := range promptsetConfig.Sinks {
			sinkConfig := allSinks[sinkName]
			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if err := sinkMgr.InstallPromptset(pkg); err != nil {
				return err
			}
		}

		if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, packageName, latestVersion, &packagelockfile.DependencyLockConfig{
			Integrity: pkg.Integrity,
		}); err != nil {
			return err
		}

		newConstraint := fmt.Sprintf("^%d.0.0", pkg.Metadata.Version.Major)
		promptsetConfig.Version = newConstraint
		if err := s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registryName, packageName, promptsetConfig); err != nil {
			return err
		}
	}

	return nil
}

func (s *ArmService) UpgradePackages(ctx context.Context, packages []string) error {
	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	lockFile, err := s.lockfileMgr.GetLockFile(ctx)
	if err != nil {
		return err
	}

	successCount := 0
	var lastErr error

	for _, pkg := range packages {
		registryName, packageName := manifest.ParseDependencyKey(pkg)
		if registryName == "" || packageName == "" {
			fmt.Fprintf(os.Stderr, "Warning: invalid package format '%s', expected registry/package\n", pkg)
			lastErr = fmt.Errorf("invalid package format: %s", pkg)
			continue
		}

		depConfig, err := s.manifestMgr.GetDependencyConfig(ctx, registryName, packageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: package not found '%s'\n", pkg)
			lastErr = fmt.Errorf("package not found: %s", pkg)
			continue
		}

		depType, ok := depConfig["type"].(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: invalid package type for '%s'\n", pkg)
			lastErr = fmt.Errorf("invalid package type: %s", pkg)
			continue
		}

		var sinks []string
		var include []string
		var exclude []string
		var priority int

		switch depType {
		case "ruleset":
			rulesetConfig, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				lastErr = err
				continue
			}
			sinks = rulesetConfig.Sinks
			include = rulesetConfig.Include
			exclude = rulesetConfig.Exclude
			priority = rulesetConfig.Priority
		case "promptset":
			promptsetConfig, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				lastErr = err
				continue
			}
			sinks = promptsetConfig.Sinks
			include = promptsetConfig.Include
			exclude = promptsetConfig.Exclude
		default:
			fmt.Fprintf(os.Stderr, "Warning: unknown package type '%s' for '%s'\n", depType, pkg)
			lastErr = fmt.Errorf("unknown package type: %s", depType)
			continue
		}

		oldVersion := s.getOldVersionFromLock(lockFile, pkg)
		newVersion, fetchedPkg, err := s.fetchLatest(ctx, registryName, packageName, include, exclude)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to upgrade '%s': %v\n", pkg, err)
			lastErr = err
			continue
		}

		if oldVersion == newVersion {
			successCount++
			continue
		}

		if oldVersion != "" {
			if err := s.uninstallFromSinks(sinks, allSinks, registryName, packageName); err != nil {
				lastErr = err
				continue
			}
			if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName); err != nil {
				lastErr = err
				continue
			}
		}

		for _, sinkName := range sinks {
			sinkConfig := allSinks[sinkName]
			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
			if depType == "ruleset" {
				if err := sinkMgr.InstallRuleset(fetchedPkg, priority); err != nil {
					lastErr = err
					continue
				}
			} else {
				if err := sinkMgr.InstallPromptset(fetchedPkg); err != nil {
					lastErr = err
					continue
				}
			}
		}

		if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, packageName, newVersion, &packagelockfile.DependencyLockConfig{
			Integrity: fetchedPkg.Integrity,
		}); err != nil {
			lastErr = err
			continue
		}

		newConstraint := fmt.Sprintf("^%d.0.0", fetchedPkg.Metadata.Version.Major)
		if depType == "ruleset" {
			rulesetConfig, _ := s.manifestMgr.GetRulesetDependencyConfig(ctx, registryName, packageName)
			rulesetConfig.Version = newConstraint
			if err := s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registryName, packageName, rulesetConfig); err != nil {
				lastErr = err
				continue
			}
		} else {
			promptsetConfig, _ := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registryName, packageName)
			promptsetConfig.Version = newConstraint
			if err := s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registryName, packageName, promptsetConfig); err != nil {
				lastErr = err
				continue
			}
		}

		successCount++
	}

	if successCount == 0 && lastErr != nil {
		return lastErr
	}

	return nil
}

func (s *ArmService) fetchLatest(ctx context.Context, registryName, packageName string, include, exclude []string) (string, *core.Package, error) {
	regConfig, err := s.manifestMgr.GetRegistryConfig(ctx, registryName)
	if err != nil {
		return "", nil, err
	}

	reg, err := s.registryFactory.CreateRegistry(registryName, regConfig)
	if err != nil {
		return "", nil, err
	}

	availableVersions, err := reg.ListPackageVersions(ctx, packageName)
	if err != nil {
		return "", nil, err
	}

	if len(availableVersions) == 0 {
		return "", nil, fmt.Errorf("no versions available for %s", packageName)
	}

	latestVersion := availableVersions[0]

	pkg, err := reg.GetPackage(ctx, packageName, &latestVersion, include, exclude)
	if err != nil {
		return "", nil, err
	}

	return latestVersion.Version, pkg, nil
}

// ListAll lists all dependencies
func (s *ArmService) ListAll(ctx context.Context) ([]*DependencyInfo, error) {
	allDeps, err := s.manifestMgr.GetAllDependenciesConfig(ctx)
	if err != nil {
		return nil, err
	}

	lockFile, err := s.lockfileMgr.GetLockFile(ctx)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var result []*DependencyInfo
	for key, config := range allDeps {
		registryName, packageName := manifest.ParseDependencyKey(key)

		var lockInfo packagelockfile.DependencyLockConfig
		if lockFile != nil && lockFile.Dependencies != nil {
			for lockKey, lockCfg := range lockFile.Dependencies {
				if strings.HasPrefix(lockKey, key+"@") {
					lockInfo = lockCfg
					break
				}
			}
		}

		result = append(result, &DependencyInfo{
			Installation: sink.PackageInstallation{
				Metadata: core.PackageMetadata{
					RegistryName: registryName,
					Name:         packageName,
				},
			},
			LockInfo: lockInfo,
			Config:   config,
		})
	}

	return result, nil
}

// GetDependencyInfo gets dependency information
func (s *ArmService) GetDependencyInfo(ctx context.Context, registry, dependencyName string) (*DependencyInfo, error) {
	config, err := s.manifestMgr.GetDependencyConfig(ctx, registry, dependencyName)
	if err != nil {
		return nil, err
	}

	lockFile, err := s.lockfileMgr.GetLockFile(ctx)
	if err != nil {
		return nil, err
	}

	key := registry + "/" + dependencyName
	var lockInfo packagelockfile.DependencyLockConfig
	var version string
	if lockFile != nil && lockFile.Dependencies != nil {
		for lockKey, lockCfg := range lockFile.Dependencies {
			if strings.HasPrefix(lockKey, key+"@") {
				lockInfo = lockCfg
				// Extract version from key (format: registry/package@version)
				parts := strings.Split(lockKey, "@")
				if len(parts) == 2 {
					version = parts[1]
				}
				break
			}
		}
	}

	return &DependencyInfo{
		Installation: sink.PackageInstallation{
			Metadata: core.PackageMetadata{
				RegistryName: registry,
				Name:         dependencyName,
			},
		},
		LockInfo: lockInfo,
		Config:   config,
		Version:  version,
	}, nil
}

// ListOutdated lists outdated dependencies
func (s *ArmService) ListOutdated(ctx context.Context) ([]*OutdatedDependency, error) {
	allDeps, err := s.manifestMgr.GetAllDependenciesConfig(ctx)
	if err != nil {
		return nil, err
	}

	lockFile, err := s.lockfileMgr.GetLockFile(ctx)
	if err != nil {
		return nil, err
	}

	var result []*OutdatedDependency
	for key, config := range allDeps {
		registryName, packageName := manifest.ParseDependencyKey(key)

		versionConstraint, ok := config["version"].(string)
		if !ok {
			continue
		}

		var currentVersion string
		if lockFile != nil && lockFile.Dependencies != nil {
			for lockKey := range lockFile.Dependencies {
				if strings.HasPrefix(lockKey, key+"@") {
					currentVersion = lockKey[len(key)+1:]
					break
				}
			}
		}

		if currentVersion == "" {
			continue
		}

		regConfig, err := s.manifestMgr.GetRegistryConfig(ctx, registryName)
		if err != nil {
			continue
		}

		reg, err := s.registryFactory.CreateRegistry(registryName, regConfig)
		if err != nil {
			continue
		}

		availableVersions, err := reg.ListPackageVersions(ctx, packageName)
		if err != nil {
			continue
		}

		if len(availableVersions) == 0 {
			continue
		}

		wantedVersion, err := core.ResolveVersion(versionConstraint, availableVersions)
		if err != nil {
			continue
		}

		latestVersion := availableVersions[0]

		if wantedVersion.Version != currentVersion || latestVersion.Version != currentVersion {
			result = append(result, &OutdatedDependency{
				Current: core.PackageMetadata{
					Name:         packageName,
					RegistryName: registryName,
					Version:      core.Version{Version: currentVersion},
				},
				Constraint: versionConstraint,
				Wanted: core.PackageMetadata{
					Name:         packageName,
					RegistryName: registryName,
					Version:      wantedVersion,
				},
				Latest: core.PackageMetadata{
					Name:         packageName,
					RegistryName: registryName,
					Version:      latestVersion,
				},
			})
		}
	}

	return result, nil
}

// SetRulesetName sets ruleset name
func (s *ArmService) SetRulesetName(ctx context.Context, registry, ruleset, newName string) error {
	return s.manifestMgr.UpdateDependencyConfigName(ctx, registry, ruleset, registry, newName)
}

// SetRulesetVersion sets ruleset version
func (s *ArmService) SetRulesetVersion(ctx context.Context, registry, ruleset, version string) error {
	config, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registry, ruleset)
	if err != nil {
		return err
	}
	config.Version = version
	return s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registry, ruleset, config)
}

// SetRulesetPriority sets ruleset priority
func (s *ArmService) SetRulesetPriority(ctx context.Context, registry, ruleset string, priority int) error {
	config, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registry, ruleset)
	if err != nil {
		return err
	}
	config.Priority = priority
	return s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registry, ruleset, config)
}

// SetRulesetInclude sets ruleset include patterns
func (s *ArmService) SetRulesetInclude(ctx context.Context, registry, ruleset string, include []string) error {
	config, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registry, ruleset)
	if err != nil {
		return err
	}
	config.Include = include
	return s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registry, ruleset, config)
}

// SetRulesetExclude sets ruleset exclude patterns
func (s *ArmService) SetRulesetExclude(ctx context.Context, registry, ruleset string, exclude []string) error {
	config, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registry, ruleset)
	if err != nil {
		return err
	}
	config.Exclude = exclude
	return s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registry, ruleset, config)
}

// SetRulesetSinks sets ruleset sinks
func (s *ArmService) SetRulesetSinks(ctx context.Context, registry, ruleset string, sinks []string) error {
	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}
	for _, sinkName := range sinks {
		if _, exists := allSinks[sinkName]; !exists {
			return fmt.Errorf("sink %s does not exist", sinkName)
		}
	}
	config, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registry, ruleset)
	if err != nil {
		return err
	}
	config.Sinks = sinks
	return s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registry, ruleset, config)
}

// SetPromptsetName sets promptset name
func (s *ArmService) SetPromptsetName(ctx context.Context, registry, ruleset, newName string) error {
	return s.manifestMgr.UpdateDependencyConfigName(ctx, registry, ruleset, registry, newName)
}

// SetPromptsetVersion sets promptset version
func (s *ArmService) SetPromptsetVersion(ctx context.Context, registry, promptset, version string) error {
	config, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registry, promptset)
	if err != nil {
		return err
	}
	config.Version = version
	return s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registry, promptset, config)
}

// SetPromptsetInclude sets promptset include patterns
func (s *ArmService) SetPromptsetInclude(ctx context.Context, registry, promptset string, include []string) error {
	config, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registry, promptset)
	if err != nil {
		return err
	}
	config.Include = include
	return s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registry, promptset, config)
}

// SetPromptsetExclude sets promptset exclude patterns
func (s *ArmService) SetPromptsetExclude(ctx context.Context, registry, promptset string, exclude []string) error {
	config, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registry, promptset)
	if err != nil {
		return err
	}
	config.Exclude = exclude
	return s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registry, promptset, config)
}

// SetPromptsetSinks sets promptset sinks
func (s *ArmService) SetPromptsetSinks(ctx context.Context, registry, promptset string, sinks []string) error {
	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}
	for _, sinkName := range sinks {
		if _, exists := allSinks[sinkName]; !exists {
			return fmt.Errorf("sink %s does not exist", sinkName)
		}
	}
	config, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registry, promptset)
	if err != nil {
		return err
	}
	config.Sinks = sinks
	return s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registry, promptset, config)
}

// --------
// Cleaning
// --------

// CleanCacheByAge cleans cache by age
func (s *ArmService) CleanCacheByAge(ctx context.Context, maxAge time.Duration) error {
	return s.CleanCacheByAgeWithHomeDir(ctx, maxAge, "")
}

// CleanCacheByAgeWithHomeDir cleans cache by age with custom home directory
func (s *ArmService) CleanCacheByAgeWithHomeDir(ctx context.Context, maxAge time.Duration, homeDir string) error {
	if homeDir == "" {
		homeDir = os.Getenv("ARM_HOME")
		if homeDir == "" {
			var err error
			homeDir, err = os.UserHomeDir()
			if err != nil {
				return err
			}
		}
	}
	storageDir := filepath.Join(homeDir, ".arm", "storage")
	return s.cleanCacheByAgeWithPath(ctx, maxAge, storageDir)
}

func (s *ArmService) cleanCacheByAgeWithPath(ctx context.Context, maxAge time.Duration, storageDir string) error {
	registriesDir := filepath.Join(storageDir, "registries")
	if _, err := os.Stat(registriesDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			packagesDir := filepath.Join(registriesDir, entry.Name(), "packages")
			if _, err := os.Stat(packagesDir); os.IsNotExist(err) {
				continue
			}
			packageCache := storage.NewPackageCache(packagesDir)
			if err := packageCache.RemoveOldVersions(ctx, maxAge); err != nil {
				return err
			}
		}
	}
	return nil
}

// CleanCacheByTimeSinceLastAccess cleans cache by time since last access
func (s *ArmService) CleanCacheByTimeSinceLastAccess(ctx context.Context, maxTimeSinceLastAccess time.Duration) error {
	return s.CleanCacheByTimeSinceLastAccessWithHomeDir(ctx, maxTimeSinceLastAccess, "")
}

// CleanCacheByTimeSinceLastAccessWithHomeDir cleans cache by time since last access with custom home directory
func (s *ArmService) CleanCacheByTimeSinceLastAccessWithHomeDir(ctx context.Context, maxTimeSinceLastAccess time.Duration, homeDir string) error {
	if homeDir == "" {
		homeDir = os.Getenv("ARM_HOME")
		if homeDir == "" {
			var err error
			homeDir, err = os.UserHomeDir()
			if err != nil {
				return err
			}
		}
	}
	storageDir := filepath.Join(homeDir, ".arm", "storage")
	return s.cleanCacheByTimeSinceLastAccessWithPath(ctx, maxTimeSinceLastAccess, storageDir)
}

func (s *ArmService) cleanCacheByTimeSinceLastAccessWithPath(ctx context.Context, maxTimeSinceLastAccess time.Duration, storageDir string) error {
	registriesDir := filepath.Join(storageDir, "registries")
	if _, err := os.Stat(registriesDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(registriesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			packagesDir := filepath.Join(registriesDir, entry.Name(), "packages")
			if _, err := os.Stat(packagesDir); os.IsNotExist(err) {
				continue
			}
			packageCache := storage.NewPackageCache(packagesDir)
			if err := packageCache.RemoveUnusedVersions(ctx, maxTimeSinceLastAccess); err != nil {
				return err
			}
		}
	}
	return nil
}

// NukeCache nukes the cache
func (s *ArmService) NukeCache(ctx context.Context) error {
	return s.NukeCacheWithHomeDir(ctx, "")
}

// NukeCacheWithHomeDir nukes the cache with custom home directory
func (s *ArmService) NukeCacheWithHomeDir(ctx context.Context, homeDir string) error {
	if homeDir == "" {
		homeDir = os.Getenv("ARM_HOME")
		if homeDir == "" {
			var err error
			homeDir, err = os.UserHomeDir()
			if err != nil {
				return err
			}
		}
	}
	storageDir := filepath.Join(homeDir, ".arm", "storage")
	return s.nukeCacheWithPath(ctx, storageDir)
}

func (s *ArmService) nukeCacheWithPath(_ context.Context, storageDir string) error {
	return os.RemoveAll(storageDir)
}

// CleanSinks cleans sinks
func (s *ArmService) CleanSinks(ctx context.Context) error {
	sinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	for _, sinkConfig := range sinks {
		sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
		if err := sinkMgr.Clean(); err != nil {
			return err
		}
	}
	return nil
}

// NukeSinks nukes sinks
func (s *ArmService) NukeSinks(ctx context.Context) error {
	sinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	for _, sinkConfig := range sinks {
		if sinkConfig.Tool == compiler.Copilot {
			// Flat layout: remove arm_* files and arm-index.json
			entries, err := os.ReadDir(sinkConfig.Directory)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return err
			}
			for _, entry := range entries {
				name := entry.Name()
				if len(name) >= 4 && name[:4] == "arm_" || name == "arm-index.json" {
					_ = os.Remove(filepath.Join(sinkConfig.Directory, name))
				}
			}
		} else {
			// Hierarchical layout: remove arm/ directory
			armDir := filepath.Join(sinkConfig.Directory, "arm")
			_ = os.RemoveAll(armDir)
		}
	}
	return nil
}

// CompileRequest groups compile parameters following ARM patterns
type CompileRequest struct {
	Paths        []string
	Tool         string
	OutputDir    string
	Namespace    string
	Force        bool
	Recursive    bool
	Verbose      bool
	ValidateOnly bool
	Include      []string
	Exclude      []string
	FailFast     bool
}

// CompileFiles compiles files
func (s *ArmService) CompileFiles(ctx context.Context, req *CompileRequest) error {
	// Discover files from input paths
	files, err := s.discoverFiles(req.Paths, req.Recursive, req.Include, req.Exclude)
	if err != nil {
		return fmt.Errorf("failed to discover files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found matching criteria")
	}

	// Process each file
	var errors []error
	for _, file := range files {
		if err := s.compileFile(ctx, file, req); err != nil {
			errors = append(errors, err)
			if req.FailFast {
				return err
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("compilation failed with %d error(s): %v", len(errors), errors[0])
	}

	return nil
}

func (s *ArmService) discoverFiles(paths []string, recursive bool, include, exclude []string) ([]*core.File, error) {
	// Default include patterns if none specified
	if len(include) == 0 {
		include = []string{"*.yml", "*.yaml"}
	}

	var files []*core.File
	seen := make(map[string]bool)

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s: %w", path, err)
		}

		if info.IsDir() {
			// Discover files in directory
			dirFiles, err := s.discoverInDirectory(path, recursive, include, exclude)
			if err != nil {
				return nil, err
			}
			for _, f := range dirFiles {
				if !seen[f.Path] {
					files = append(files, f)
					seen[f.Path] = true
				}
			}
		} else if !seen[path] {
			// Single file
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", path, err)
			}
			files = append(files, &core.File{
				Path:    path,
				Content: content,
				Size:    info.Size(),
			})
			seen[path] = true
		}
	}

	return files, nil
}

func (s *ArmService) discoverInDirectory(dir string, recursive bool, include, exclude []string) ([]*core.File, error) {
	var files []*core.File

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		if info.IsDir() {
			if !recursive && path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path for pattern matching
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return nil
		}

		// Check if file matches patterns
		if !s.matchesPatterns(relPath, include, exclude) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Continue on errors
		}

		files = append(files, &core.File{
			Path:    path,
			Content: content,
			Size:    info.Size(),
		})

		return nil
	}

	if err := filepath.Walk(dir, walkFn); err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", dir, err)
	}

	return files, nil
}

func (s *ArmService) matchesPatterns(filePath string, include, exclude []string) bool {
	// Check exclude patterns first
	for _, pattern := range exclude {
		if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
			return false
		}
	}

	// If no include patterns, accept all
	if len(include) == 0 {
		return true
	}

	// Check include patterns
	for _, pattern := range include {
		if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
			return true
		}
	}

	return false
}

func (s *ArmService) parseTool(toolStr string) (compiler.Tool, error) {
	switch toolStr {
	case "cursor":
		return compiler.Cursor, nil
	case "copilot":
		return compiler.Copilot, nil
	case "amazonq":
		return compiler.AmazonQ, nil
	case "markdown":
		return compiler.Markdown, nil
	default:
		return "", fmt.Errorf("invalid tool: %s (must be cursor, copilot, amazonq, or markdown)", toolStr)
	}
}

func (s *ArmService) compileFile(_ context.Context, file *core.File, req *CompileRequest) error {
	// Detect file type
	isRuleset := filetype.IsRulesetFile(file)
	isPromptset := filetype.IsPromptsetFile(file)

	if !isRuleset && !isPromptset {
		return fmt.Errorf("file %s is not a valid ruleset or promptset", file.Path)
	}

	// Parse and validate
	if isRuleset {
		ruleset, err := parser.ParseRuleset(file)
		if err != nil {
			return fmt.Errorf("failed to parse ruleset %s: %w", file.Path, err)
		}

		// If validate-only, we're done
		if req.ValidateOnly {
			if req.Verbose {
				fmt.Printf("✓ Validated ruleset: %s\n", file.Path)
			}
			return nil
		}

		// Determine namespace
		namespace := req.Namespace
		if namespace == "" {
			namespace = ruleset.Metadata.ID
		}

		// Compile
		tool, err := s.parseTool(req.Tool)
		if err != nil {
			return err
		}

		compiledFiles, err := compiler.CompileRuleset(tool, namespace, ruleset)
		if err != nil {
			return fmt.Errorf("failed to compile ruleset %s: %w", file.Path, err)
		}

		// Write output files
		if err := s.writeCompiledFiles(compiledFiles, req.OutputDir, req.Force); err != nil {
			return fmt.Errorf("failed to write compiled files for %s: %w", file.Path, err)
		}

		if req.Verbose {
			fmt.Printf("✓ Compiled ruleset: %s -> %d file(s)\n", file.Path, len(compiledFiles))
		}

	} else if isPromptset {
		promptset, err := parser.ParsePromptset(file)
		if err != nil {
			return fmt.Errorf("failed to parse promptset %s: %w", file.Path, err)
		}

		// If validate-only, we're done
		if req.ValidateOnly {
			if req.Verbose {
				fmt.Printf("✓ Validated promptset: %s\n", file.Path)
			}
			return nil
		}

		// Determine namespace
		namespace := req.Namespace
		if namespace == "" {
			namespace = promptset.Metadata.ID
		}

		// Compile
		tool, err := s.parseTool(req.Tool)
		if err != nil {
			return err
		}

		compiledFiles, err := compiler.CompilePromptset(tool, namespace, promptset)
		if err != nil {
			return fmt.Errorf("failed to compile promptset %s: %w", file.Path, err)
		}

		// Write output files
		if err := s.writeCompiledFiles(compiledFiles, req.OutputDir, req.Force); err != nil {
			return fmt.Errorf("failed to write compiled files for %s: %w", file.Path, err)
		}

		if req.Verbose {
			fmt.Printf("✓ Compiled promptset: %s -> %d file(s)\n", file.Path, len(compiledFiles))
		}
	}

	return nil
}

func (s *ArmService) writeCompiledFiles(files []*core.File, outputDir string, force bool) error {
	for _, file := range files {
		outputPath := filepath.Join(outputDir, file.Path)

		// Check if file exists
		if !force {
			if _, err := os.Stat(outputPath); err == nil {
				return fmt.Errorf("file %s already exists (use --force to overwrite)", outputPath)
			}
		}

		// Create directory if needed
		dir := filepath.Dir(outputPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(outputPath, file.Content, 0o644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputPath, err)
		}
	}

	return nil
}
