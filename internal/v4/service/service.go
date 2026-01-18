package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
	"github.com/jomadu/ai-resource-manager/internal/v4/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/v4/registry"
	"github.com/jomadu/ai-resource-manager/internal/v4/sink"
)

type DependencyInfo struct {
	Installation sink.PackageInstallation
	LockInfo     packagelockfile.DependencyLockConfig
	Config       map[string]interface{}
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

	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, config)
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
func (s *ArmService) SetRegistryName(ctx context.Context, name string, newName string) error {
	return s.manifestMgr.UpdateRegistryConfigName(ctx, name, newName)
}

// SetRegistryURL sets registry URL
func (s *ArmService) SetRegistryURL(ctx context.Context, name string, url string) error {
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
func (s *ArmService) SetGitLabRegistryProjectID(ctx context.Context, name string, projectID string) error {
	config, err := s.manifestMgr.GetGitLabRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.ProjectID = projectID
	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, config)
}

// SetGitLabRegistryGroupID sets GitLab registry group ID
func (s *ArmService) SetGitLabRegistryGroupID(ctx context.Context, name string, groupID string) error {
	config, err := s.manifestMgr.GetGitLabRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.GroupID = groupID
	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, config)
}

// SetGitLabRegistryAPIVersion sets GitLab registry API version
func (s *ArmService) SetGitLabRegistryAPIVersion(ctx context.Context, name string, apiVersion string) error {
	config, err := s.manifestMgr.GetGitLabRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.APIVersion = apiVersion
	return s.manifestMgr.UpsertGitLabRegistryConfig(ctx, name, config)
}

// SetCloudsmithRegistryOwner sets Cloudsmith registry owner
func (s *ArmService) SetCloudsmithRegistryOwner(ctx context.Context, name string, owner string) error {
	config, err := s.manifestMgr.GetCloudsmithRegistryConfig(ctx, name)
	if err != nil {
		return err
	}

	config.Owner = owner
	return s.manifestMgr.UpsertCloudsmithRegistryConfig(ctx, name, config)
}

// SetCloudsmithRegistryRepository sets Cloudsmith registry repository
func (s *ArmService) SetCloudsmithRegistryRepository(ctx context.Context, name string, repository string) error {
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
func (s *ArmService) SetSinkName(ctx context.Context, name string, newName string) error {
	return s.manifestMgr.UpdateSinkConfigName(ctx, name, newName)
}

// SetSinkDirectory sets sink directory
func (s *ArmService) SetSinkDirectory(ctx context.Context, name string, directory string) error {
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

// InstallAll installs all dependencies
func (s *ArmService) InstallAll(ctx context.Context) error {
	deps, err := s.manifestMgr.GetAllDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	for key, depConfig := range deps {
		depType, ok := depConfig["type"].(string)
		if !ok {
			return fmt.Errorf("dependency %s missing type", key)
		}

		registryName, packageName := manifest.ParseDependencyKey(key)

		if depType == "ruleset" {
			rulesetCfg, err := s.manifestMgr.GetRulesetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				return err
			}
			if err := s.InstallRuleset(ctx, registryName, packageName, rulesetCfg.Version, rulesetCfg.Priority, rulesetCfg.Include, rulesetCfg.Exclude, rulesetCfg.Sinks); err != nil {
				return err
			}
		} else if depType == "promptset" {
			promptsetCfg, err := s.manifestMgr.GetPromptsetDependencyConfig(ctx, registryName, packageName)
			if err != nil {
				return err
			}
			if err := s.InstallPromptset(ctx, registryName, packageName, promptsetCfg.Version, promptsetCfg.Include, promptsetCfg.Exclude, promptsetCfg.Sinks); err != nil {
				return err
			}
		}
	}

	return nil
}

// resolveAndFetchPackage validates registry/sinks, resolves version, and fetches package
func (s *ArmService) resolveAndFetchPackage(ctx context.Context, registryName, packageName, version string, include, exclude, sinks []string) (*core.Package, string, map[string]manifest.SinkConfig, error) {
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

	resolvedVersion, err := core.ResolveVersion(version, availableVersions)
	if err != nil {
		return nil, "", nil, err
	}

	pkg, err := reg.GetPackage(ctx, packageName, resolvedVersion, include, exclude)
	if err != nil {
		return nil, "", nil, err
	}

	return pkg, resolvedVersion.Version, allSinks, nil
}

// InstallRuleset installs a ruleset
func (s *ArmService) InstallRuleset(ctx context.Context, registryName, ruleset, version string, priority int, include []string, exclude []string, sinks []string) error {
	pkg, resolvedVersion, allSinks, err := s.resolveAndFetchPackage(ctx, registryName, ruleset, version, include, exclude, sinks)
	if err != nil {
		return err
	}

	// 5. Update manifest with dependency
	depConfig := manifest.RulesetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Version: version,
			Sinks:   sinks,
			Include: include,
			Exclude: exclude,
		},
		Priority: priority,
	}

	if err := s.manifestMgr.UpsertRulesetDependencyConfig(ctx, registryName, ruleset, depConfig); err != nil {
		return err
	}

	// 6. Update lock file
	if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, ruleset, resolvedVersion, &packagelockfile.DependencyLockConfig{
		Integrity: pkg.Integrity,
	}); err != nil {
		return err
	}

	// 7. Install to each sink
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
func (s *ArmService) InstallPromptset(ctx context.Context, registryName, promptset, version string, include []string, exclude []string, sinks []string) error {
	pkg, resolvedVersion, allSinks, err := s.resolveAndFetchPackage(ctx, registryName, promptset, version, include, exclude, sinks)
	if err != nil {
		return err
	}

	// 5. Update manifest with dependency
	depConfig := manifest.PromptsetDependencyConfig{
		BaseDependencyConfig: manifest.BaseDependencyConfig{
			Version: version,
			Sinks:   sinks,
			Include: include,
			Exclude: exclude,
		},
	}

	if err := s.manifestMgr.UpsertPromptsetDependencyConfig(ctx, registryName, promptset, depConfig); err != nil {
		return err
	}

	// 6. Update lock file
	if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, promptset, resolvedVersion, &packagelockfile.DependencyLockConfig{
		Integrity: pkg.Integrity,
	}); err != nil {
		return err
	}

	// 7. Install to each sink
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
	deps, err := s.manifestMgr.GetAllDependenciesConfig(ctx)
	if err != nil {
		return err
	}

	allSinks, err := s.manifestMgr.GetAllSinksConfig(ctx)
	if err != nil {
		return err
	}

	for key := range deps {
		registryName, packageName := manifest.ParseDependencyKey(key)
		depConfig, err := s.manifestMgr.GetDependencyConfig(ctx, registryName, packageName)
		if err != nil {
			return err
		}

		sinks, ok := depConfig["sinks"].([]interface{})
		if ok {
			version, ok := depConfig["version"].(string)
			if ok {
				for _, sinkInterface := range sinks {
					sinkName, ok := sinkInterface.(string)
					if !ok {
						continue
					}

					sinkConfig, exists := allSinks[sinkName]
					if !exists {
						continue
					}

					sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
					if err := sinkMgr.Uninstall(core.PackageMetadata{
						RegistryName: registryName,
						Name:         packageName,
						Version:      core.Version{Version: version},
					}); err != nil {
						return err
					}
				}

				if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName, version); err != nil {
					return err
				}
			}
		}

		if err := s.manifestMgr.RemoveDependencyConfig(ctx, registryName, packageName); err != nil {
			return err
		}
	}

	return nil
}

// UpdateAll updates all dependencies
func (s *ArmService) UpdateAll(ctx context.Context) error {
	deps, err := s.manifestMgr.GetAllDependenciesConfig(ctx)
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

	for key, depConfig := range deps {
		depType, ok := depConfig["type"].(string)
		if !ok {
			continue
		}

		registryName, packageName := manifest.ParseDependencyKey(key)

		version, ok := depConfig["version"].(string)
		if !ok {
			continue
		}

		regConfig, err := s.manifestMgr.GetRegistryConfig(ctx, registryName)
		if err != nil {
			return err
		}

		reg, err := s.registryFactory.CreateRegistry(registryName, regConfig)
		if err != nil {
			return err
		}

		availableVersions, err := reg.ListPackageVersions(ctx, packageName)
		if err != nil {
			return err
		}

		resolvedVersion, err := core.ResolveVersion(version, availableVersions)
		if err != nil {
			return err
		}

		var oldVersion string
		if lockFile != nil && lockFile.Dependencies != nil {
			for lockKey := range lockFile.Dependencies {
				if strings.HasPrefix(lockKey, key+"@") {
					oldVersion = lockKey[len(key)+1:]
					break
				}
			}
		}

		if oldVersion == resolvedVersion.Version {
			continue
		}

		include, _ := depConfig["include"].([]string)
		exclude, _ := depConfig["exclude"].([]string)

		pkg, err := reg.GetPackage(ctx, packageName, resolvedVersion, include, exclude)
		if err != nil {
			return err
		}

		sinksInterface, ok := depConfig["sinks"].([]interface{})
		if !ok {
			continue
		}

		if oldVersion != "" {
			for _, sinkInterface := range sinksInterface {
				sinkName, ok := sinkInterface.(string)
				if !ok {
					continue
				}

				sinkConfig, exists := allSinks[sinkName]
				if !exists {
					continue
				}

				sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)
				if err := sinkMgr.Uninstall(core.PackageMetadata{
					RegistryName: registryName,
					Name:         packageName,
					Version:      core.Version{Version: oldVersion},
				}); err != nil {
					return err
				}
			}
			if err := s.lockfileMgr.RemoveDependencyLock(ctx, registryName, packageName, oldVersion); err != nil {
				return err
			}
		}

		for _, sinkInterface := range sinksInterface {
			sinkName, ok := sinkInterface.(string)
			if !ok {
				continue
			}

			sinkConfig, exists := allSinks[sinkName]
			if !exists {
				continue
			}

			sinkMgr := sink.NewManager(sinkConfig.Directory, sinkConfig.Tool)

			if depType == "ruleset" {
				priority, _ := depConfig["priority"].(int)
				if err := sinkMgr.InstallRuleset(pkg, priority); err != nil {
					return err
				}
			} else if depType == "promptset" {
				if err := sinkMgr.InstallPromptset(pkg); err != nil {
					return err
				}
			}
		}

		if err := s.lockfileMgr.UpsertDependencyLock(ctx, registryName, packageName, resolvedVersion.Version, &packagelockfile.DependencyLockConfig{
			Integrity: pkg.Integrity,
		}); err != nil {
			return err
		}
	}

	return nil
}

// UpgradeAll upgrades all dependencies
func (s *ArmService) UpgradeAll(ctx context.Context) error {
	// TODO: implement
	return nil
}

// ListAll lists all dependencies
func (s *ArmService) ListAll(ctx context.Context) ([]*DependencyInfo, error) {
	// TODO: implement
	return nil, nil
}

// GetDependencyInfo gets dependency information
func (s *ArmService) GetDependencyInfo(ctx context.Context, registry, dependencyName string) (*DependencyInfo, error) {
	// TODO: implement
	return nil, nil
}

// ListOutdated lists outdated dependencies
func (s *ArmService) ListOutdated(ctx context.Context) ([]*OutdatedDependency, error) {
	// TODO: implement
	return nil, nil
}

// SetRulesetName sets ruleset name
func (s *ArmService) SetRulesetName(ctx context.Context, registry, ruleset, newName string) error {
	// TODO: implement
	return nil
}

// SetRulesetVersion sets ruleset version
func (s *ArmService) SetRulesetVersion(ctx context.Context, registry, ruleset, version string) error {
	// TODO: implement
	return nil
}

// SetRulesetPriority sets ruleset priority
func (s *ArmService) SetRulesetPriority(ctx context.Context, registry, ruleset string, priority int) error {
	// TODO: implement
	return nil
}

// SetRulesetInclude sets ruleset include patterns
func (s *ArmService) SetRulesetInclude(ctx context.Context, registry, ruleset string, include []string) error {
	// TODO: implement
	return nil
}

// SetRulesetExclude sets ruleset exclude patterns
func (s *ArmService) SetRulesetExclude(ctx context.Context, registry, ruleset string, exclude []string) error {
	// TODO: implement
	return nil
}

// SetRulesetSinks sets ruleset sinks
func (s *ArmService) SetRulesetSinks(ctx context.Context, registry, ruleset string, sinks []string) error {
	// TODO: implement
	return nil
}

// SetPromptsetName sets promptset name
func (s *ArmService) SetPromptsetName(ctx context.Context, registry, ruleset, newName string) error {
	// TODO: implement
	return nil
}

// SetPromptsetVersion sets promptset version
func (s *ArmService) SetPromptsetVersion(ctx context.Context, registry, ruleset, version string) error {
	// TODO: implement
	return nil
}

// SetPromptsetInclude sets promptset include patterns
func (s *ArmService) SetPromptsetInclude(ctx context.Context, registry, ruleset string, include []string) error {
	// TODO: implement
	return nil
}

// SetPromptsetExclude sets promptset exclude patterns
func (s *ArmService) SetPromptsetExclude(ctx context.Context, registry, ruleset string, exclude []string) error {
	// TODO: implement
	return nil
}

// SetPromptsetSinks sets promptset sinks
func (s *ArmService) SetPromptsetSinks(ctx context.Context, registry, ruleset string, sinks []string) error {
	// TODO: implement
	return nil
}

// --------
// Cleaning
// --------

// CleanCacheByAge cleans cache by age
func (s *ArmService) CleanCacheByAge(ctx context.Context, maxAge time.Duration) error {
	// TODO: implement
	return nil
}

// CleanCacheByTimeSinceLastAccess cleans cache by time since last access
func (s *ArmService) CleanCacheByTimeSinceLastAccess(ctx context.Context, maxTimeSinceLastAccess time.Duration) error {
	// TODO: implement
	return nil
}

// NukeCache nukes the cache
func (s *ArmService) NukeCache(ctx context.Context) error {
	// TODO: implement
	return nil
}

// CleanSinks cleans sinks
func (s *ArmService) CleanSinks(ctx context.Context) error {
	// TODO: implement
	return nil
}

// NukeSinks nukes sinks
func (s *ArmService) NukeSinks(ctx context.Context) error {
	// TODO: implement
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
	// TODO: implement
	return nil
}
