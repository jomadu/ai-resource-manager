package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/resource"
)

// Manager handles arm.json manifest file operations.
type Manager interface {
	// Ruleset operations
	GetRuleset(ctx context.Context, registry, ruleset string) (*RulesetConfig, error)
	GetRulesets(ctx context.Context) (map[string]map[string]RulesetConfig, error)
	AddRuleset(ctx context.Context, registry, ruleset string, entry *RulesetConfig) error
	UpdateRuleset(ctx context.Context, registry, ruleset string, entry *RulesetConfig) error
	RemoveRuleset(ctx context.Context, registry, ruleset string) error

	// Promptset operations
	GetPromptset(ctx context.Context, registry, promptset string) (*PromptsetConfig, error)
	GetPromptsets(ctx context.Context) (map[string]map[string]PromptsetConfig, error)
	AddPromptset(ctx context.Context, registry, promptset string, entry *PromptsetConfig) error
	UpdatePromptset(ctx context.Context, registry, promptset string, entry *PromptsetConfig) error
	RemovePromptset(ctx context.Context, registry, promptset string) error

	// Registry operations
	GetRegistries(ctx context.Context) (map[string]map[string]interface{}, error)
	GetGitRegistry(ctx context.Context, name string) (*registry.GitRegistryConfig, error)
	AddGitRegistry(ctx context.Context, name string, config registry.GitRegistryConfig, force bool) error
	UpdateGitRegistry(ctx context.Context, name string, config registry.GitRegistryConfig, force bool) error
	GetGitLabRegistry(ctx context.Context, name string) (*registry.GitLabRegistryConfig, error)
	AddGitLabRegistry(ctx context.Context, name string, config *registry.GitLabRegistryConfig, force bool) error
	UpdateGitLabRegistry(ctx context.Context, name string, config *registry.GitLabRegistryConfig, force bool) error
	GetCloudsmithRegistry(ctx context.Context, name string) (*registry.CloudsmithRegistryConfig, error)
	AddCloudsmithRegistry(ctx context.Context, name string, config *registry.CloudsmithRegistryConfig, force bool) error
	UpdateCloudsmithRegistry(ctx context.Context, name string, config *registry.CloudsmithRegistryConfig, force bool) error
	RemoveRegistry(ctx context.Context, name string) error

	// Sink operations
	GetSinks(ctx context.Context) (map[string]SinkConfig, error)
	GetSink(ctx context.Context, name string) (*SinkConfig, error)
	AddSink(ctx context.Context, name string, sink SinkConfig, force bool) error
	UpdateSink(ctx context.Context, name, field, value string) error
	RemoveSink(ctx context.Context, name string) error
}

// FileManager implements file-based manifest management.
type FileManager struct{}

// NewFileManager creates a new file-based manifest manager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f *FileManager) GetRuleset(ctx context.Context, registry, ruleset string) (*RulesetConfig, error) {
	entries, err := f.GetRulesets(ctx)
	if err != nil {
		return nil, err
	}
	registryMap, exists := entries[registry]
	if !exists {
		return nil, fmt.Errorf("ruleset %s/%s not found in manifest (registry %s has no rulesets)", registry, ruleset, registry)
	}
	entry, exists := registryMap[ruleset]
	if !exists {
		return nil, fmt.Errorf("ruleset %s/%s not found in manifest (not installed or already uninstalled)", registry, ruleset)
	}
	return &entry, nil
}

func (f *FileManager) GetRulesets(ctx context.Context) (map[string]map[string]RulesetConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Packages.Rulesets, nil
}

func (f *FileManager) GetRegistries(ctx context.Context) (map[string]map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Registries, nil
}

func (f *FileManager) AddRuleset(ctx context.Context, registry, ruleset string, entry *RulesetConfig) error {
	if err := f.validateRuleset(entry); err != nil {
		return err
	}
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Packages.Rulesets == nil {
		manifest.Packages.Rulesets = make(map[string]map[string]RulesetConfig)
	}
	if manifest.Packages.Rulesets[registry] == nil {
		manifest.Packages.Rulesets[registry] = make(map[string]RulesetConfig)
	}
	if _, exists := manifest.Packages.Rulesets[registry][ruleset]; exists {
		return errors.New("ruleset already exists")
	}
	manifest.Packages.Rulesets[registry][ruleset] = *entry
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateRuleset(ctx context.Context, registry, ruleset string, entry *RulesetConfig) error {
	if err := f.validateRuleset(entry); err != nil {
		return err
	}
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Packages.Rulesets == nil || manifest.Packages.Rulesets[registry] == nil {
		return errors.New("registry not found")
	}
	if _, exists := manifest.Packages.Rulesets[registry][ruleset]; !exists {
		return errors.New("ruleset not found")
	}
	manifest.Packages.Rulesets[registry][ruleset] = *entry
	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveRuleset(ctx context.Context, registry, ruleset string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Packages.Rulesets == nil || manifest.Packages.Rulesets[registry] == nil {
		return errors.New("registry not found")
	}
	if _, exists := manifest.Packages.Rulesets[registry][ruleset]; !exists {
		return errors.New("ruleset not found")
	}
	delete(manifest.Packages.Rulesets[registry], ruleset)
	if len(manifest.Packages.Rulesets[registry]) == 0 {
		delete(manifest.Packages.Rulesets, registry)
	}
	return f.saveManifest(manifest)
}

// Promptset methods
func (f *FileManager) GetPromptset(ctx context.Context, registry, promptset string) (*PromptsetConfig, error) {
	entries, err := f.GetPromptsets(ctx)
	if err != nil {
		return nil, err
	}
	registryMap, exists := entries[registry]
	if !exists {
		return nil, fmt.Errorf("promptset %s/%s not found in manifest (registry %s has no promptsets)", registry, promptset, registry)
	}
	entry, exists := registryMap[promptset]
	if !exists {
		return nil, fmt.Errorf("promptset %s/%s not found in manifest (not installed or already uninstalled)", registry, promptset)
	}
	return &entry, nil
}

func (f *FileManager) GetPromptsets(ctx context.Context) (map[string]map[string]PromptsetConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Packages.Promptsets, nil
}

func (f *FileManager) AddPromptset(ctx context.Context, registry, promptset string, entry *PromptsetConfig) error {
	if err := f.validatePromptset(entry); err != nil {
		return err
	}
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Packages.Promptsets == nil {
		manifest.Packages.Promptsets = make(map[string]map[string]PromptsetConfig)
	}
	if manifest.Packages.Promptsets[registry] == nil {
		manifest.Packages.Promptsets[registry] = make(map[string]PromptsetConfig)
	}
	if _, exists := manifest.Packages.Promptsets[registry][promptset]; exists {
		return errors.New("promptset already exists")
	}
	manifest.Packages.Promptsets[registry][promptset] = *entry
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdatePromptset(ctx context.Context, registry, promptset string, entry *PromptsetConfig) error {
	if err := f.validatePromptset(entry); err != nil {
		return err
	}
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Packages.Promptsets == nil || manifest.Packages.Promptsets[registry] == nil {
		return errors.New("registry not found")
	}
	if _, exists := manifest.Packages.Promptsets[registry][promptset]; !exists {
		return errors.New("promptset not found")
	}
	manifest.Packages.Promptsets[registry][promptset] = *entry
	return f.saveManifest(manifest)
}

func (f *FileManager) RemovePromptset(ctx context.Context, registry, promptset string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Packages.Promptsets == nil || manifest.Packages.Promptsets[registry] == nil {
		return errors.New("registry not found")
	}
	if _, exists := manifest.Packages.Promptsets[registry][promptset]; !exists {
		return errors.New("promptset not found")
	}
	delete(manifest.Packages.Promptsets[registry], promptset)
	if len(manifest.Packages.Promptsets[registry]) == 0 {
		delete(manifest.Packages.Promptsets, registry)
	}
	return f.saveManifest(manifest)
}

// Validation methods
func (f *FileManager) validateRuleset(entry *RulesetConfig) error {
	// Rulesets can have priority, so no validation needed for priority field
	return nil
}

func (f *FileManager) validatePromptset(entry *PromptsetConfig) error {
	// Promptsets should not have priority (this is enforced by the type system)
	return nil
}

func (f *FileManager) loadManifest() (*Manifest, error) {
	data, err := os.ReadFile("arm.json")
	if err != nil {
		// File doesn't exist, return initialized manifest
		return &Manifest{
			Registries: make(map[string]map[string]interface{}),
			Packages: PackageConfig{
				Rulesets:   make(map[string]map[string]RulesetConfig),
				Promptsets: make(map[string]map[string]PromptsetConfig),
			},
			Sinks: make(map[string]SinkConfig),
		}, nil
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	// Initialize maps only if they're nil to support minimal configurations
	if manifest.Registries == nil {
		manifest.Registries = make(map[string]map[string]interface{})
	}
	if manifest.Packages.Rulesets == nil {
		manifest.Packages.Rulesets = make(map[string]map[string]RulesetConfig)
	}
	if manifest.Packages.Promptsets == nil {
		manifest.Packages.Promptsets = make(map[string]map[string]PromptsetConfig)
	}
	if manifest.Sinks == nil {
		manifest.Sinks = make(map[string]SinkConfig)
	}
	return &manifest, nil
}

func (f *FileManager) saveManifest(manifest *Manifest) error {
	// Clean up empty maps to keep JSON minimal
	if len(manifest.Registries) == 0 {
		manifest.Registries = nil
	}
	if len(manifest.Packages.Rulesets) == 0 {
		manifest.Packages.Rulesets = nil
	}
	if len(manifest.Packages.Promptsets) == 0 {
		manifest.Packages.Promptsets = nil
	}
	if len(manifest.Sinks) == 0 {
		manifest.Sinks = nil
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("arm.json", data, 0o644)
}

// LoadManifest loads the manifest file (public method)
func (f *FileManager) LoadManifest() (*Manifest, error) {
	return f.loadManifest()
}

// SaveManifest saves the manifest file (public method)
func (f *FileManager) SaveManifest(manifest *Manifest) error {
	return f.saveManifest(manifest)
}

func (f *FileManager) AddGitRegistry(ctx context.Context, name string, config registry.GitRegistryConfig, force bool) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists && !force {
		return errors.New("registry already exists (use --force to overwrite)")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(configBytes, &rawConfig); err != nil {
		return err
	}

	manifest.Registries[name] = rawConfig

	return f.saveManifest(manifest)
}

func (f *FileManager) AddGitLabRegistry(ctx context.Context, name string, config *registry.GitLabRegistryConfig, force bool) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists && !force {
		return errors.New("registry already exists (use --force to overwrite)")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(configBytes, &rawConfig); err != nil {
		return err
	}

	manifest.Registries[name] = rawConfig

	return f.saveManifest(manifest)
}

func (f *FileManager) AddCloudsmithRegistry(ctx context.Context, name string, config *registry.CloudsmithRegistryConfig, force bool) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists && !force {
		return errors.New("registry already exists (use --force to overwrite)")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(configBytes, &rawConfig); err != nil {
		return err
	}

	manifest.Registries[name] = rawConfig

	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveRegistry(ctx context.Context, name string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	_, exists := manifest.Registries[name]
	if !exists {
		return errors.New("registry not found")
	}

	delete(manifest.Registries, name)
	return f.saveManifest(manifest)
}

// New registry methods
func (f *FileManager) GetGitRegistry(ctx context.Context, name string) (*registry.GitRegistryConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	rawConfig, exists := manifest.Registries[name]
	if !exists {
		return nil, fmt.Errorf("registry %s not found", name)
	}

	// Convert raw config to GitRegistryConfig
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, err
	}

	var config registry.GitRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (f *FileManager) UpdateGitRegistry(ctx context.Context, name string, config registry.GitRegistryConfig, force bool) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists && !force {
		return errors.New("registry already exists (use --force to overwrite)")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(configBytes, &rawConfig); err != nil {
		return err
	}

	manifest.Registries[name] = rawConfig
	return f.saveManifest(manifest)
}

func (f *FileManager) GetGitLabRegistry(ctx context.Context, name string) (*registry.GitLabRegistryConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	rawConfig, exists := manifest.Registries[name]
	if !exists {
		return nil, fmt.Errorf("registry %s not found", name)
	}

	// Convert raw config to GitLabRegistryConfig
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, err
	}

	var config registry.GitLabRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (f *FileManager) UpdateGitLabRegistry(ctx context.Context, name string, config *registry.GitLabRegistryConfig, force bool) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists && !force {
		return errors.New("registry already exists (use --force to overwrite)")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(configBytes, &rawConfig); err != nil {
		return err
	}

	manifest.Registries[name] = rawConfig
	return f.saveManifest(manifest)
}

func (f *FileManager) GetCloudsmithRegistry(ctx context.Context, name string) (*registry.CloudsmithRegistryConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}

	rawConfig, exists := manifest.Registries[name]
	if !exists {
		return nil, fmt.Errorf("registry %s not found", name)
	}

	// Convert raw config to CloudsmithRegistryConfig
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, err
	}

	var config registry.CloudsmithRegistryConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (f *FileManager) UpdateCloudsmithRegistry(ctx context.Context, name string, config *registry.CloudsmithRegistryConfig, force bool) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Registries[name]; exists && !force {
		return errors.New("registry already exists (use --force to overwrite)")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(configBytes, &rawConfig); err != nil {
		return err
	}

	manifest.Registries[name] = rawConfig
	return f.saveManifest(manifest)
}

func (f *FileManager) GetSinks(ctx context.Context) (map[string]SinkConfig, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Sinks, nil
}

func (f *FileManager) GetSink(ctx context.Context, name string) (*SinkConfig, error) {
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

func (f *FileManager) AddSink(ctx context.Context, name string, sink SinkConfig, force bool) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	if _, exists := manifest.Sinks[name]; exists && !force {
		return fmt.Errorf("sink %s already exists (use --force to overwrite)", name)
	}

	manifest.Sinks[name] = sink
	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveSink(ctx context.Context, name string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	_, exists := manifest.Sinks[name]
	if !exists {
		return fmt.Errorf("sink %s not found", name)
	}

	// Check if any rulesets or promptsets are using this sink
	var usingResources []string
	for registryName, rulesets := range manifest.Packages.Rulesets {
		for rulesetName, ruleset := range rulesets {
			for _, sink := range ruleset.Sinks {
				if sink == name {
					usingResources = append(usingResources, fmt.Sprintf("ruleset %s/%s", registryName, rulesetName))
					break
				}
			}
		}
	}
	for registryName, promptsets := range manifest.Packages.Promptsets {
		for promptsetName, promptset := range promptsets {
			for _, sink := range promptset.Sinks {
				if sink == name {
					usingResources = append(usingResources, fmt.Sprintf("promptset %s/%s", registryName, promptsetName))
					break
				}
			}
		}
	}

	if len(usingResources) > 0 {
		return fmt.Errorf("cannot remove sink %s: it is being used by %v. Uninstall these resources first", name, usingResources)
	}

	// Remove sink from configuration
	delete(manifest.Sinks, name)

	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateSink(ctx context.Context, name, field, value string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	sink, exists := manifest.Sinks[name]
	if !exists {
		return fmt.Errorf("sink %s not found", name)
	}

	switch field {
	case "directory":
		sink.Directory = strings.TrimSpace(value)
	case "layout":
		if value != "hierarchical" && value != "flat" {
			return fmt.Errorf("layout must be 'hierarchical' or 'flat'")
		}
		sink.Layout = value
	case "compileTarget":
		sink.CompileTarget = resource.CompileTarget(strings.TrimSpace(value))
	default:
		return fmt.Errorf("unknown field '%s' (valid: directory, layout, compileTarget)", field)
	}

	manifest.Sinks[name] = sink

	return f.saveManifest(manifest)
}
