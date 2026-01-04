package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/registry"
)

type ResourceType string

const (
	ResourceTypeRuleset  ResourceType = "ruleset"
	ResourceTypePromptset ResourceType = "promptset"
)

type Manifest struct {
	Version      int                               `json:"version"`
	Registries   map[string]map[string]interface{} `json:"registries,omitempty"`
	Sinks        map[string]SinkConfig             `json:"sinks,omitempty"`
	Dependencies map[string]map[string]interface{} `json:"dependencies,omitempty"`
}

type SinkConfig struct {
	Directory string        `json:"directory"`
	Tool      compiler.Tool `json:"tool"`
}

// Type-safe dependency configs for API
type BaseDependencyConfig struct {
	Type    ResourceType `json:"type"`
	Version string       `json:"version"`
	Sinks   []string     `json:"sinks"`
	Include []string     `json:"include,omitempty"`
	Exclude []string     `json:"exclude,omitempty"`
}

type RulesetDependencyConfig struct {
	BaseDependencyConfig
	Priority int `json:"priority,omitempty"`
}

type PromptsetDependencyConfig struct {
	BaseDependencyConfig
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

	// Update package keys from "oldName/package" to "newName/package"
	for key, depConfig := range manifest.Dependencies {
		regName, pkgName := parseDependencyKey(key)
		if regName == name {
			newKey := dependencyKey(newName, pkgName)
			manifest.Dependencies[newKey] = depConfig
			delete(manifest.Dependencies, key)
		}
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
	
	// Remove all packages from this registry
	for key := range manifest.Dependencies {
		regName, _ := parseDependencyKey(key)
		if regName == name {
			delete(manifest.Dependencies, key)
		}
	}
	
	return f.saveManifest(manifest)
}

// Registry operations (type-safe helpers)

func (f *FileManager) GetGitRegistryConfig(ctx context.Context, name string) (*registry.GitRegistryConfig, error) {
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

func (f *FileManager) GetGitLabRegistryConfig(ctx context.Context, name string) (*registry.GitLabRegistryConfig, error) {
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

func (f *FileManager) GetCloudsmithRegistryConfig(ctx context.Context, name string) (*registry.CloudsmithRegistryConfig, error) {
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

func (f *FileManager) UpsertGitRegistryConfig(ctx context.Context, name string, config *registry.GitRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpsertGitLabRegistryConfig(ctx context.Context, name string, config *registry.GitLabRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	configMap, err := convertRegistryToMap(config)
	if err != nil {
		return err
	}

	manifest.Registries[name] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpsertCloudsmithRegistryConfig(ctx context.Context, name string, config *registry.CloudsmithRegistryConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
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

func (f *FileManager) UpsertSinkConfig(ctx context.Context, name string, config *SinkConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
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

// Dependencies operations (generic)

func (f *FileManager) GetAllDependenciesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Dependencies, nil
}

func (f *FileManager) GetDependencyConfig(ctx context.Context, key string) (map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	config, exists := manifest.Dependencies[key]
	if !exists {
		return nil, fmt.Errorf("dependency %s not found", key)
	}

	return config, nil
}

func (f *FileManager) UpdateDependencyConfigName(ctx context.Context, key string, newKey string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	config, exists := manifest.Dependencies[key]
	if !exists {
		return fmt.Errorf("dependency %s not found", key)
	}

	if _, exists := manifest.Dependencies[newKey]; exists {
		return fmt.Errorf("dependency %s already exists", newKey)
	}

	manifest.Dependencies[newKey] = config
	delete(manifest.Dependencies, key)

	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveDependencyConfig(ctx context.Context, key string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Dependencies[key]; !exists {
		return fmt.Errorf("dependency %s not found", key)
	}

	delete(manifest.Dependencies, key)
	return f.saveManifest(manifest)
}

// Dependencies operations (type-safe helpers)

func (f *FileManager) GetRulesetDependencyConfig(ctx context.Context, key string) (*RulesetDependencyConfig, error) {
	rawConfig, err := f.GetDependencyConfig(ctx, key)
	if err != nil {
		return nil, err
	}

	// Check dependency type
	depType, ok := rawConfig["type"].(string)
	if !ok || depType != "ruleset" {
		return nil, fmt.Errorf("dependency %s is not a ruleset", key)
	}

	return convertMapToRulesetDependency(rawConfig)
}

func (f *FileManager) GetPromptsetDependencyConfig(ctx context.Context, key string) (*PromptsetDependencyConfig, error) {
	rawConfig, err := f.GetDependencyConfig(ctx, key)
	if err != nil {
		return nil, err
	}

	// Check dependency type
	depType, ok := rawConfig["type"].(string)
	if !ok || depType != "promptset" {
		return nil, fmt.Errorf("dependency %s is not a promptset", key)
	}

	return convertMapToPromptsetDependency(rawConfig)
}

func (f *FileManager) UpsertRulesetDependencyConfig(ctx context.Context, key string, config *RulesetDependencyConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	// Set type to ruleset
	config.Type = ResourceTypeRuleset

	configMap, err := convertDependencyToMap(config)
	if err != nil {
		return err
	}

	manifest.Dependencies[key] = configMap
	return f.saveManifest(manifest)
}

func (f *FileManager) UpsertPromptsetDependencyConfig(ctx context.Context, key string, config *PromptsetDependencyConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	// Set type to promptset
	config.Type = ResourceTypePromptset

	configMap, err := convertDependencyToMap(config)
	if err != nil {
		return err
	}

	manifest.Dependencies[key] = configMap
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
	if manifest.Sinks == nil {
		manifest.Sinks = make(map[string]SinkConfig)
	}
	if manifest.Dependencies == nil {
		manifest.Dependencies = make(map[string]map[string]interface{})
	}

	return &manifest, nil
}

// newEmptyManifest creates a new empty manifest with all maps initialized.
// Used when manifest file doesn't exist yet.
func (f *FileManager) newEmptyManifest() *Manifest {
	return &Manifest{
		Version:      1,
		Registries:   make(map[string]map[string]interface{}),
		Sinks:        make(map[string]SinkConfig),
		Dependencies: make(map[string]map[string]interface{}),
	}
}

// saveManifest saves the manifest to arm.json file.
// Cleans up empty maps before saving to keep JSON minimal.
func (f *FileManager) saveManifest(manifest *Manifest) error {
	// Clean up empty maps to keep JSON minimal
	if len(manifest.Registries) == 0 {
		manifest.Registries = nil
	}
	if len(manifest.Sinks) == 0 {
		manifest.Sinks = nil
	}
	if len(manifest.Dependencies) == 0 {
		manifest.Dependencies = nil
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

// convertDependencyToMap converts a typed dependency config to map[string]interface{}.
// Used when storing typed configs in the generic manifest structure.
func convertDependencyToMap(config interface{}) (map[string]interface{}, error) {
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

// convertMapToGitRegistry converts map[string]interface{} to registry.GitRegistryConfig.
func convertMapToGitRegistry(m map[string]interface{}) (*registry.GitRegistryConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config registry.GitRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToGitLabRegistry converts map[string]interface{} to registry.GitLabRegistryConfig.
func convertMapToGitLabRegistry(m map[string]interface{}) (*registry.GitLabRegistryConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config registry.GitLabRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToCloudsmithRegistry converts map[string]interface{} to registry.CloudsmithRegistryConfig.
func convertMapToCloudsmithRegistry(m map[string]interface{}) (*registry.CloudsmithRegistryConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config registry.CloudsmithRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToRulesetDependency converts map[string]interface{} to RulesetDependencyConfig.
func convertMapToRulesetDependency(m map[string]interface{}) (*RulesetDependencyConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config RulesetDependencyConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// convertMapToPromptsetDependency converts map[string]interface{} to PromptsetDependencyConfig.
func convertMapToPromptsetDependency(m map[string]interface{}) (*PromptsetDependencyConfig, error) {
	configBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var config PromptsetDependencyConfig
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

// ensureSinkExists ensures a sink exists in the manifest, returns error if not.
func (f *FileManager) ensureSinkExists(manifest *Manifest, name string) error {
	if _, exists := manifest.Sinks[name]; !exists {
		return fmt.Errorf("sink %s not found", name)
	}
	return nil
}

// Local dependency key helpers (manifest uses registry/package format without version)

// dependencyKey creates a dependency key in format "registry/package"
func dependencyKey(registry, packageName string) string {
	return fmt.Sprintf("%s/%s", registry, packageName)
}

// parseDependencyKey parses a dependency key and returns registry, package name
func parseDependencyKey(key string) (registry, packageName string) {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return "", "" // Invalid format
	}
	return parts[0], parts[1]
}