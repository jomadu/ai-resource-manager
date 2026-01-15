package service

import (
	"context"
	"errors"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
	"github.com/jomadu/ai-resource-manager/internal/v4/packagelockfile"
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
	manifestMgr manifest.Manager
}

// NewArmService creates a new ARM service
func NewArmService(manifestMgr manifest.Manager) *ArmService {
	return &ArmService{
		manifestMgr: manifestMgr,
	}
}

// ---------------------------------------------
// Registry Management (Git, GitLab, Cloudsmith)
// ---------------------------------------------

// AddGitRegistry adds a Git registry
func (s *ArmService) AddGitRegistry(ctx context.Context, name, url string, branches []string, force bool) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Registries[name]; !force && exists {
		return errors.New("registry already exists")
	}

	config := map[string]interface{}{
		"url":  url,
		"type": "git",
	}
	if len(branches) > 0 {
		config["branches"] = branches
	}

	m.Registries[name] = config
	return s.manifestMgr.Save(m)
}

// AddGitLabRegistry adds a GitLab registry
func (s *ArmService) AddGitLabRegistry(ctx context.Context, name, url, projectID, groupID, apiVersion string, force bool) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Registries[name]; !force && exists {
		return errors.New("registry already exists")
	}

	config := map[string]interface{}{
		"url":  url,
		"type": "gitlab",
	}
	if projectID != "" {
		config["projectId"] = projectID
	}
	if groupID != "" {
		config["groupId"] = groupID
	}
	if apiVersion != "" {
		config["apiVersion"] = apiVersion
	}

	m.Registries[name] = config
	return s.manifestMgr.Save(m)
}

// AddCloudsmithRegistry adds a Cloudsmith registry
func (s *ArmService) AddCloudsmithRegistry(ctx context.Context, name, url, owner, repository string, force bool) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Registries[name]; !force && exists {
		return errors.New("registry already exists")
	}

	m.Registries[name] = map[string]interface{}{
		"url":        url,
		"type":       "cloudsmith",
		"owner":      owner,
		"repository": repository,
	}

	return s.manifestMgr.Save(m)
}

// RemoveRegistry removes a registry
func (s *ArmService) RemoveRegistry(ctx context.Context, name string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Registries[name]; !exists {
		return errors.New("registry does not exist")
	}

	delete(m.Registries, name)
	return s.manifestMgr.Save(m)
}

// SetRegistryName sets registry name
func (s *ArmService) SetRegistryName(ctx context.Context, name string, newName string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Registries[name]; !exists {
		return errors.New("registry does not exist")
	}

	if _, exists := m.Registries[newName]; exists {
		return errors.New("registry with new name already exists")
	}

	m.Registries[newName] = m.Registries[name]
	delete(m.Registries, name)

	return s.manifestMgr.Save(m)
}

// SetRegistryURL sets registry URL
func (s *ArmService) SetRegistryURL(ctx context.Context, name string, url string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	reg, exists := m.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}

	reg["url"] = url
	m.Registries[name] = reg

	return s.manifestMgr.Save(m)
}

// SetGitRegistryBranches sets Git registry branches
func (s *ArmService) SetGitRegistryBranches(ctx context.Context, name string, branches []string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	reg, exists := m.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}

	if regType, ok := reg["type"].(string); !ok || regType != "git" {
		return errors.New("registry is not a git registry")
	}

	reg["branches"] = branches
	m.Registries[name] = reg

	return s.manifestMgr.Save(m)
}

// SetGitLabRegistryProjectID sets GitLab registry project ID
func (s *ArmService) SetGitLabRegistryProjectID(ctx context.Context, name string, projectID string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	reg, exists := m.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}

	if regType, ok := reg["type"].(string); !ok || regType != "gitlab" {
		return errors.New("registry is not a gitlab registry")
	}

	reg["projectId"] = projectID
	m.Registries[name] = reg

	return s.manifestMgr.Save(m)
}

// SetGitLabRegistryGroupID sets GitLab registry group ID
func (s *ArmService) SetGitLabRegistryGroupID(ctx context.Context, name string, groupID string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	reg, exists := m.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}

	if regType, ok := reg["type"].(string); !ok || regType != "gitlab" {
		return errors.New("registry is not a gitlab registry")
	}

	reg["groupId"] = groupID
	m.Registries[name] = reg

	return s.manifestMgr.Save(m)
}

// SetGitLabRegistryAPIVersion sets GitLab registry API version
func (s *ArmService) SetGitLabRegistryAPIVersion(ctx context.Context, name string, apiVersion string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	reg, exists := m.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}

	if regType, ok := reg["type"].(string); !ok || regType != "gitlab" {
		return errors.New("registry is not a gitlab registry")
	}

	reg["apiVersion"] = apiVersion
	m.Registries[name] = reg

	return s.manifestMgr.Save(m)
}

// SetCloudsmithRegistryOwner sets Cloudsmith registry owner
func (s *ArmService) SetCloudsmithRegistryOwner(ctx context.Context, name string, owner string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	reg, exists := m.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}

	if regType, ok := reg["type"].(string); !ok || regType != "cloudsmith" {
		return errors.New("registry is not a cloudsmith registry")
	}

	reg["owner"] = owner
	m.Registries[name] = reg

	return s.manifestMgr.Save(m)
}

// SetCloudsmithRegistryRepository sets Cloudsmith registry repository
func (s *ArmService) SetCloudsmithRegistryRepository(ctx context.Context, name string, repository string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	reg, exists := m.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}

	if regType, ok := reg["type"].(string); !ok || regType != "cloudsmith" {
		return errors.New("registry is not a cloudsmith registry")
	}

	reg["repository"] = repository
	m.Registries[name] = reg

	return s.manifestMgr.Save(m)
}

// GetRegistryConfig gets registry configuration
func (s *ArmService) GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error) {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return nil, err
	}

	cfg, exists := m.Registries[name]
	if !exists {
		return nil, errors.New("registry does not exist")
	}

	return cfg, nil
}

// GetAllRegistriesConfig gets all registry configurations
func (s *ArmService) GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return nil, err
	}

	return m.Registries, nil
}

// ---------------
// Sink Management
// ---------------

// AddSink adds a sink
func (s *ArmService) AddSink(ctx context.Context, name, directory string, tool compiler.Tool, force bool) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Sinks[name]; !force && exists {
		return errors.New("sink already exists")
	}

	m.Sinks[name] = manifest.SinkConfig{
		Directory: directory,
		Tool:      tool,
	}

	return s.manifestMgr.Save(m)
}

// RemoveSink removes a sink
func (s *ArmService) RemoveSink(ctx context.Context, name string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Sinks[name]; !exists {
		return errors.New("sink does not exist")
	}

	delete(m.Sinks, name)

	return s.manifestMgr.Save(m)
}

// GetSinkConfig gets sink configuration
func (s *ArmService) GetSinkConfig(ctx context.Context, name string) (manifest.SinkConfig, error) {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return manifest.SinkConfig{}, err
	}

	cfg, exists := m.Sinks[name]
	if !exists {
		return manifest.SinkConfig{}, errors.New("sink does not exist")
	}

	return cfg, nil
}

// GetAllSinkConfigs gets all sink configurations
func (s *ArmService) GetAllSinkConfigs(ctx context.Context) (map[string]manifest.SinkConfig, error) {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return nil, err
	}

	return m.Sinks, nil
}

// SetSinkName sets sink name
func (s *ArmService) SetSinkName(ctx context.Context, name string, newName string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	if _, exists := m.Sinks[name]; !exists {
		return errors.New("sink does not exist")
	}

	if _, exists := m.Sinks[newName]; exists {
		return errors.New("sink with new name already exists")
	}

	m.Sinks[newName] = m.Sinks[name]
	delete(m.Sinks, name)

	return s.manifestMgr.Save(m)
}

// SetSinkDirectory sets sink directory
func (s *ArmService) SetSinkDirectory(ctx context.Context, name string, directory string) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	sink, exists := m.Sinks[name]
	if !exists {
		return errors.New("sink does not exist")
	}

	sink.Directory = directory
	m.Sinks[name] = sink

	return s.manifestMgr.Save(m)
}

// SetSinkTool sets sink tool
func (s *ArmService) SetSinkTool(ctx context.Context, name string, tool compiler.Tool) error {
	m, err := s.manifestMgr.Load()
	if err != nil {
		return err
	}

	sink, exists := m.Sinks[name]
	if !exists {
		return errors.New("sink does not exist")
	}

	sink.Tool = tool
	m.Sinks[name] = sink

	return s.manifestMgr.Save(m)
}

// ---------------------
// Dependency Management
// ---------------------

// InstallAll installs all dependencies
func (s *ArmService) InstallAll(ctx context.Context) error {
	// TODO: implement
	return nil
}

// InstallRuleset installs a ruleset
func (s *ArmService) InstallRuleset(ctx context.Context, registry, ruleset, version string, priority int, include []string, exclude []string, sinks []string) error {
	// TODO: implement
	return nil
}

// InstallPromptset installs a promptset
func (s *ArmService) InstallPromptset(ctx context.Context, registry, promptset, version string, include []string, exclude []string, sinks []string) error {
	// TODO: implement
	return nil
}

// UninstallAll uninstalls all dependencies
func (s *ArmService) UninstallAll(ctx context.Context) error {
	// TODO: implement
	return nil
}

// UpdateAll updates all dependencies
func (s *ArmService) UpdateAll(ctx context.Context) error {
	// TODO: implement
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
