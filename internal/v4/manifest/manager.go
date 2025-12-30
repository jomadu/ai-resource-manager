package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
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
	Dependencies Dependencies                      `json:"dependencies"`
}

type Dependencies struct {
	Rulesets   map[string]RulesetConfig   `json:"rulesets,omitempty"`
	Promptsets map[string]PromptsetConfig `json:"promptsets,omitempty"`
}

type SinkConfig struct {
	Directory string `json:"directory"`
	Tool      string `json:"tool"`
}

type RulesetConfig struct {
	Version  string   `json:"version"`
	Priority int      `json:"priority,omitempty"`
	Sinks    []string `json:"sinks"`
	Include  []string `json:"include,omitempty"`
	Exclude  []string `json:"exclude,omitempty"`
}

type PromptsetConfig struct {
	Version string   `json:"version"`
	Sinks   []string `json:"sinks"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
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
	for key, rulesetConfig := range manifest.Dependencies.Rulesets {
		regName, pkgName := core.ParsePackageKey(key)
		if regName == name {
			newKey := core.PackageKey(newName, pkgName)
			manifest.Dependencies.Rulesets[newKey] = rulesetConfig
			delete(manifest.Dependencies.Rulesets, key)
		}
	}

	for key, promptsetConfig := range manifest.Dependencies.Promptsets {
		regName, pkgName := core.ParsePackageKey(key)
		if regName == name {
			newKey := core.PackageKey(newName, pkgName)
			manifest.Dependencies.Promptsets[newKey] = promptsetConfig
			delete(manifest.Dependencies.Promptsets, key)
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
	for key := range manifest.Dependencies.Rulesets {
		regName, _ := core.ParsePackageKey(key)
		if regName == name {
			delete(manifest.Dependencies.Rulesets, key)
		}
	}
	
	for key := range manifest.Dependencies.Promptsets {
		regName, _ := core.ParsePackageKey(key)
		if regName == name {
			delete(manifest.Dependencies.Promptsets, key)
		}
	}
	
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

// Package operations

func (f *FileManager) GetRulesetConfig(ctx context.Context, packageKey string) (*RulesetConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	config, exists := manifest.Dependencies.Rulesets[packageKey]
	if !exists {
		return nil, fmt.Errorf("ruleset %s not found", packageKey)
	}

	return &config, nil
}

func (f *FileManager) GetPromptsetConfig(ctx context.Context, packageKey string) (*PromptsetConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	config, exists := manifest.Dependencies.Promptsets[packageKey]
	if !exists {
		return nil, fmt.Errorf("promptset %s not found", packageKey)
	}

	return &config, nil
}

func (f *FileManager) AddRulesetConfig(ctx context.Context, packageKey string, config *RulesetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Dependencies.Rulesets[packageKey]; exists {
		return fmt.Errorf("ruleset %s already exists", packageKey)
	}

	manifest.Dependencies.Rulesets[packageKey] = *config
	return f.saveManifest(manifest)
}

func (f *FileManager) AddPromptsetConfig(ctx context.Context, packageKey string, config *PromptsetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Dependencies.Promptsets[packageKey]; exists {
		return fmt.Errorf("promptset %s already exists", packageKey)
	}

	manifest.Dependencies.Promptsets[packageKey] = *config
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateRulesetConfig(ctx context.Context, packageKey string, config *RulesetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Dependencies.Rulesets[packageKey]; !exists {
		return fmt.Errorf("ruleset %s not found", packageKey)
	}

	manifest.Dependencies.Rulesets[packageKey] = *config
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdatePromptsetConfig(ctx context.Context, packageKey string, config *PromptsetConfig) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Dependencies.Promptsets[packageKey]; !exists {
		return fmt.Errorf("promptset %s not found", packageKey)
	}

	manifest.Dependencies.Promptsets[packageKey] = *config
	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveRulesetConfig(ctx context.Context, packageKey string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Dependencies.Rulesets[packageKey]; !exists {
		return fmt.Errorf("ruleset %s not found", packageKey)
	}

	delete(manifest.Dependencies.Rulesets, packageKey)
	return f.saveManifest(manifest)
}

func (f *FileManager) RemovePromptsetConfig(ctx context.Context, packageKey string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Dependencies.Promptsets[packageKey]; !exists {
		return fmt.Errorf("promptset %s not found", packageKey)
	}

	delete(manifest.Dependencies.Promptsets, packageKey)
	return f.saveManifest(manifest)
}

func (f *FileManager) GetAllRulesets(ctx context.Context) (map[string]*RulesetConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*RulesetConfig)
	for key, config := range manifest.Dependencies.Rulesets {
		configCopy := config
		result[key] = &configCopy
	}

	return result, nil
}

func (f *FileManager) GetAllPromptsets(ctx context.Context) (map[string]*PromptsetConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*PromptsetConfig)
	for key, config := range manifest.Dependencies.Promptsets {
		configCopy := config
		result[key] = &configCopy
	}

	return result, nil
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
	if manifest.Dependencies.Rulesets == nil {
		manifest.Dependencies.Rulesets = make(map[string]RulesetConfig)
	}
	if manifest.Dependencies.Promptsets == nil {
		manifest.Dependencies.Promptsets = make(map[string]PromptsetConfig)
	}

	return &manifest, nil
}

// newEmptyManifest creates a new empty manifest with all maps initialized.
// Used when manifest file doesn't exist yet.
func (f *FileManager) newEmptyManifest() *Manifest {
	return &Manifest{
		Version:    1,
		Registries: make(map[string]map[string]interface{}),
		Sinks:      make(map[string]SinkConfig),
		Dependencies: Dependencies{
			Rulesets:   make(map[string]RulesetConfig),
			Promptsets: make(map[string]PromptsetConfig),
		},
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
	if len(manifest.Dependencies.Rulesets) == 0 && len(manifest.Dependencies.Promptsets) == 0 {
		manifest.Dependencies = Dependencies{}
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