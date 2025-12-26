package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type ResourceType string

const (
	ResourceTypeRuleset  ResourceType = "ruleset"
	ResourceTypePromptset ResourceType = "promptset"
)

type Manifest struct {
	Version string `json:"version"`
	Registries map[string]map[string]interface{} `json:"registries,omitempty"`
	Sinks      map[string]SinkConfig             `json:"sinks,omitempty"`
	Packages   map[string]map[string]interface{} `json:"packages"`
}

type SinkConfig struct {
	Directory string `json:"directory"`
	Layout string `json:"layout,omitempty"`
	CompileTarget string `json:"compileTarget"`
}

type PackageConfig struct {
	Version string `json:"version"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
	Sinks []string `json:"sinks"`
	ResourceType ResourceType `json:"resourceType,omitempty"`
}

type RulesetConfig struct {
	PackageConfig
	Priority int `json:"priority,omitempty"`
}

type PromptsetConfig struct {
	PackageConfig
}

type RegistryConfig struct {
	URL string `json:"url"`
	Type string `json:"type"`
}

type GitRegistryConfig struct {
	RegistryConfig
	Branches []string `json:"branches,omitempty"`
}

type GitLabRegistryConfig struct {
	RegistryConfig
	ProjectID string `json:"projectId,omitempty"`
	GroupID string `json:"groupId,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type CloudsmithRegistryConfig struct {
	RegistryConfig
	Owner string `json:"owner"`
	Repository string `json:"repository"`
}
type Manager interface {
	// Registry operations
	GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error)
	GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error)
	UpdateRegistryConfigName(ctx context.Context, name string, newName string) error
	RemoveRegistryConfig(ctx context.Context, name string) error
	GetGitRegistryConfig(ctx context.Context, name string) (*GitRegistryConfig, error)
	GetGitLabRegistryConfig(ctx context.Context, name string) (*GitLabRegistryConfig, error)
	GetCloudsmithRegistryConfig(ctx context.Context, name string) (*CloudsmithRegistryConfig, error)
	AddGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error
	AddGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error
	AddCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error
	UpdateGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error
	UpdateGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error
	UpdateCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error

	// Sink operations
	GetAllSinksConfig(ctx context.Context) (map[string]*SinkConfig, error)
	GetSinkConfig(ctx context.Context, name string) (*SinkConfig, error)
	AddSinkConfig(ctx context.Context, name string, config *SinkConfig) error
	UpdateSinkConfigName(ctx context.Context, name string, newName string) error
	UpdateSinkConfig(ctx context.Context, name string, config *SinkConfig) error
	RemoveSinkConfig(ctx context.Context, name string) error

	// Package operations
	GetAllPackagesConfig(ctx context.Context) (map[string]map[string]interface{}, error)
	GetPackageConfig(ctx context.Context, registryName, packageName string) (map[string]interface{}, error)
	UpdatePackageConfigName(ctx context.Context, registryName, packageName string, newPackageName string) error
	RemovePackageConfig(ctx context.Context, registryName, packageName string) error
	GetRulesetConfig(ctx context.Context, registryName, packageName string) (*RulesetConfig, error)
	GetPromptsetConfig(ctx context.Context, registryName, packageName string) (*PromptsetConfig, error)
	AddRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error
	AddPromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error
	UpdateRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error
	UpdatePromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error
}

// FileManager implements file-based manifest management.
// It reads from and writes to arm.json in the current directory.
type FileManager struct {
	manifestPath string
}

// NewFileManager creates a new file-based manifest manager.
// Uses "arm.json" in the current directory.
func NewFileManager() *FileManager {
	return &FileManager{manifestPath: "arm.json"}
}

// NewFileManagerWithPath creates a new file-based manifest manager with a custom path.
// Useful for testing.
func NewFileManagerWithPath(manifestPath string) *FileManager {
	return &FileManager{manifestPath: manifestPath}
}

// Registry operations (generic)

func (f *FileManager) GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Registries, nil
}

func (f *FileManager) GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	config, exists := manifest.Registries[name]
	if !exists {
		return nil, fmt.Errorf("registry %s not found", name)
	}

	return config, nil
}



func (f *FileManager) UpdateRegistryConfigName(ctx context.Context, name string, newName string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureRegistryExists(manifest, name); err != nil {
		return err
	}

	if _, exists := manifest.Registries[newName]; exists {
		return fmt.Errorf("registry %s already exists", newName)
	}

	// Move registry config
	manifest.Registries[newName] = manifest.Registries[name]
	delete(manifest.Registries, name)

	// Move all packages from old registry to new registry
	if packages, exists := manifest.Packages[name]; exists {
		manifest.Packages[newName] = packages
		delete(manifest.Packages, name)
	}

	return f.saveManifest(manifest)
}



func (f *FileManager) RemoveRegistryConfig(ctx context.Context, name string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureRegistryExists(manifest, name); err != nil {
		return err
	}

	delete(manifest.Registries, name)
	return f.saveManifest(manifest)
}

// Registry operations (type-safe helpers)

func (f *FileManager) GetGitRegistryConfig(ctx context.Context, name string) (*GitRegistryConfig, error) {
	rawConfig, err := f.GetRegistryConfig(ctx, name)
	if err != nil {
		return nil, err
	}

	// Check registry type
	regType, ok := rawConfig["type"].(string)
	if !ok || regType != "git" {
		return nil, fmt.Errorf("registry %s is not a git registry", name)
	}

	return convertMapToGitRegistry(rawConfig)
}

func (f *FileManager) GetGitLabRegistryConfig(ctx context.Context, name string) (*GitLabRegistryConfig, error) {
	rawConfig, err := f.GetRegistryConfig(ctx, name)
	if err != nil {
		return nil, err
	}

	// Check registry type
	regType, ok := rawConfig["type"].(string)
	if !ok || regType != "gitlab" {
		return nil, fmt.Errorf("registry %s is not a gitlab registry", name)
	}

	return convertMapToGitLabRegistry(rawConfig)
}

func (f *FileManager) GetCloudsmithRegistryConfig(ctx context.Context, name string) (*CloudsmithRegistryConfig, error) {
	rawConfig, err := f.GetRegistryConfig(ctx, name)
	if err != nil {
		return nil, err
	}

	// Check registry type
	regType, ok := rawConfig["type"].(string)
	if !ok || regType != "cloudsmith" {
		return nil, fmt.Errorf("registry %s is not a cloudsmith registry", name)
	}

	return convertMapToCloudsmithRegistry(rawConfig)
}

func (f *FileManager) AddGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists {
		return fmt.Errorf("registry %s already exists", name)
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) AddGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists {
		return fmt.Errorf("registry %s already exists", name)
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) AddCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists {
		return fmt.Errorf("registry %s already exists", name)
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureRegistryExists(manifest, name); err != nil {
		return err
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureRegistryExists(manifest, name); err != nil {
		return err
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureRegistryExists(manifest, name); err != nil {
		return err
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

// Sink operations

func (f *FileManager) GetAllSinksConfig(ctx context.Context) (map[string]*SinkConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*SinkConfig)
	for name, sink := range manifest.Sinks {
		sinkCopy := sink
		result[name] = &sinkCopy
	}

	return result, nil
}

func (f *FileManager) GetSinkConfig(ctx context.Context, name string) (*SinkConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	sink, exists := manifest.Sinks[name]
	if !exists {
		return nil, fmt.Errorf("sink %s not found", name)
	}

	return &sink, nil
}

func (f *FileManager) AddSinkConfig(ctx context.Context, name string, config *SinkConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Sinks[name]; exists {
		return fmt.Errorf("sink %s already exists", name)
	}

	manifest.Sinks[name] = *config
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateSinkConfigName(ctx context.Context, name string, newName string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureSinkExists(manifest, name); err != nil {
		return err
	}

	if _, exists := manifest.Sinks[newName]; exists {
		return fmt.Errorf("sink %s already exists", newName)
	}

	manifest.Sinks[newName] = manifest.Sinks[name]
	delete(manifest.Sinks, name)

	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateSinkConfig(ctx context.Context, name string, config *SinkConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureSinkExists(manifest, name); err != nil {
		return err
	}

	manifest.Sinks[name] = *config
	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveSinkConfig(ctx context.Context, name string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensureSinkExists(manifest, name); err != nil {
		return err
	}

	delete(manifest.Sinks, name)
	return f.saveManifest(manifest)
}

// Package operations (generic)

func (f *FileManager) GetAllPackagesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Packages, nil
}

func (f *FileManager) GetPackageConfig(ctx context.Context, registryName, packageName string) (map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	registry, exists := manifest.Packages[registryName]
	if !exists {
		return nil, fmt.Errorf("registry %s not found", registryName)
	}

	configInterface, exists := registry[packageName]
	if !exists {
		return nil, fmt.Errorf("package %s not found in registry %s", packageName, registryName)
	}

	config, ok := configInterface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("package %s/%s has invalid config format", registryName, packageName)
	}

	return config, nil
}



func (f *FileManager) UpdatePackageConfigName(ctx context.Context, registryName, packageName string, newPackageName string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensurePackageExists(manifest, registryName, packageName); err != nil {
		return err
	}

	if _, exists := manifest.Packages[registryName][newPackageName]; exists {
		return fmt.Errorf("package %s already exists in registry %s", newPackageName, registryName)
	}

	manifest.Packages[registryName][newPackageName] = manifest.Packages[registryName][packageName]
	delete(manifest.Packages[registryName], packageName)

	return f.saveManifest(manifest)
}



func (f *FileManager) RemovePackageConfig(ctx context.Context, registryName, packageName string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensurePackageExists(manifest, registryName, packageName); err != nil {
		return err
	}

	delete(manifest.Packages[registryName], packageName)

	// Remove empty registry
	if len(manifest.Packages[registryName]) == 0 {
		delete(manifest.Packages, registryName)
	}

	return f.saveManifest(manifest)
}

// Package operations (type-safe helpers)

func (f *FileManager) GetRulesetConfig(ctx context.Context, registryName, packageName string) (*RulesetConfig, error) {
	rawConfig, err := f.GetPackageConfig(ctx, registryName, packageName)
	if err != nil {
		return nil, err
	}

	// Check resource type
	resourceType, ok := rawConfig["resourceType"].(string)
	if !ok || resourceType != string(ResourceTypeRuleset) {
		return nil, fmt.Errorf("package %s/%s is not a ruleset", registryName, packageName)
	}

	return convertMapToRulesetConfig(rawConfig)
}

func (f *FileManager) GetPromptsetConfig(ctx context.Context, registryName, packageName string) (*PromptsetConfig, error) {
	rawConfig, err := f.GetPackageConfig(ctx, registryName, packageName)
	if err != nil {
		return nil, err
	}

	// Check resource type
	resourceType, ok := rawConfig["resourceType"].(string)
	if !ok || resourceType != string(ResourceTypePromptset) {
		return nil, fmt.Errorf("package %s/%s is not a promptset", registryName, packageName)
	}

	return convertMapToPromptsetConfig(rawConfig)
}

func (f *FileManager) AddRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if manifest.Packages[registryName] == nil {
		manifest.Packages[registryName] = make(map[string]interface{})
	}

	if _, exists := manifest.Packages[registryName][packageName]; exists {
		return fmt.Errorf("package %s already exists in registry %s", packageName, registryName)
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Packages[registryName][packageName] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) AddPromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if manifest.Packages[registryName] == nil {
		manifest.Packages[registryName] = make(map[string]interface{})
	}

	if _, exists := manifest.Packages[registryName][packageName]; exists {
		return fmt.Errorf("package %s already exists in registry %s", packageName, registryName)
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Packages[registryName][packageName] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensurePackageExists(manifest, registryName, packageName); err != nil {
		return err
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Packages[registryName][packageName] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdatePromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if err := f.ensurePackageExists(manifest, registryName, packageName); err != nil {
		return err
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Packages[registryName][packageName] = configMap
	return f.saveManifest(manifest)
}

// Helper functions for FileManager implementation

// loadManifest loads the manifest from arm.json file.
// If file doesn't exist, returns an initialized empty manifest (no error).
// This makes most operations work even if manifest hasn't been created yet.
// Initializes all maps to prevent nil pointer issues.
func (f *FileManager) loadManifest() (*Manifest, error) {
	data, err := os.ReadFile(f.manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return f.newEmptyManifest(), nil
		}
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	// Initialize nil maps to prevent nil pointer issues
	if manifest.Registries == nil {
		manifest.Registries = make(map[string]map[string]interface{})
	}
	if manifest.Packages == nil {
		manifest.Packages = make(map[string]map[string]interface{})
	}
	if manifest.Sinks == nil {
		manifest.Sinks = make(map[string]SinkConfig)
	}

	return &manifest, nil
}

// newEmptyManifest creates a new empty manifest with all maps initialized.
// Used when manifest file doesn't exist yet.
func (f *FileManager) newEmptyManifest() *Manifest {
	return &Manifest{
		Version:    "1.0.0",
		Registries: make(map[string]map[string]interface{}),
		Packages:   make(map[string]map[string]interface{}),
		Sinks:      make(map[string]SinkConfig),
	}
}

// saveManifest saves the manifest to arm.json file.
// Cleans up empty maps before saving to keep JSON minimal.
func (f *FileManager) saveManifest(manifest *Manifest) error {
	// Clean up empty maps to keep JSON minimal
	if len(manifest.Registries) == 0 {
		manifest.Registries = nil
	}
	if len(manifest.Packages) == 0 {
		manifest.Packages = nil
	}
	if len(manifest.Sinks) == 0 {
		manifest.Sinks = nil
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.manifestPath, data, 0644)
}

// convertRegistryToMap converts a typed registry config to map[string]interface{}.
// Used when storing typed configs in the generic manifest structure.
func convertRegistryToMap(config interface{}) (map[string]interface{}, error) {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(configBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// convertMapToGitRegistry converts map[string]interface{} to GitRegistryConfig.
func convertMapToGitRegistry(m map[string]interface{}) (*GitRegistryConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config GitRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToGitLabRegistry converts map[string]interface{} to GitLabRegistryConfig.
func convertMapToGitLabRegistry(m map[string]interface{}) (*GitLabRegistryConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config GitLabRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToCloudsmithRegistry converts map[string]interface{} to CloudsmithRegistryConfig.
func convertMapToCloudsmithRegistry(m map[string]interface{}) (*CloudsmithRegistryConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config CloudsmithRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToRulesetConfig converts map[string]interface{} to RulesetConfig.
func convertMapToRulesetConfig(m map[string]interface{}) (*RulesetConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config RulesetConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToPromptsetConfig converts map[string]interface{} to PromptsetConfig.
func convertMapToPromptsetConfig(m map[string]interface{}) (*PromptsetConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config PromptsetConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// ensureRegistryExists ensures a registry exists in the manifest, returns error if not.
func (f *FileManager) ensureRegistryExists(manifest *Manifest, name string) error {
	if _, exists := manifest.Registries[name]; !exists {
		return fmt.Errorf("registry %s not found", name)
	}
	return nil
}

// ensurePackageExists ensures a package exists in the manifest, returns error if not.
func (f *FileManager) ensurePackageExists(manifest *Manifest, registryName, packageName string) error {
	registry, exists := manifest.Packages[registryName]
	if !exists {
		return fmt.Errorf("registry %s not found", registryName)
	}
	if _, exists := registry[packageName]; !exists {
		return fmt.Errorf("package %s not found in registry %s", packageName, registryName)
	}
	return nil
}

// ensureSinkExists ensures a sink exists in the manifest, returns error if not.
func (f *FileManager) ensureSinkExists(manifest *Manifest, name string) error {
	if _, exists := manifest.Sinks[name]; !exists {
		return fmt.Errorf("sink %s not found", name)
	}
	return nil
}