package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/registry"
)

// Manager handles arm.json manifest file operations.
type Manager interface {
	GetEntry(ctx context.Context, registry, ruleset string) (*Entry, error)
	GetEntries(ctx context.Context) (map[string]map[string]Entry, error)
	GetRawRegistries(ctx context.Context) (map[string]map[string]interface{}, error)
	AddGitRegistry(ctx context.Context, name string, config registry.GitRegistryConfig, force bool) error
	UpdateGitRegistry(ctx context.Context, name, field, value string) error
	RemoveRegistry(ctx context.Context, name string) error
	CreateEntry(ctx context.Context, registry, ruleset string, entry Entry) error
	UpdateEntry(ctx context.Context, registry, ruleset string, entry Entry) error
	RemoveEntry(ctx context.Context, registry, ruleset string) error
}

// FileManager implements file-based manifest management.
type FileManager struct{}

// NewFileManager creates a new file-based manifest manager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f *FileManager) GetEntry(ctx context.Context, registry, ruleset string) (*Entry, error) {
	entries, err := f.GetEntries(ctx)
	if err != nil {
		return nil, err
	}
	registryMap, exists := entries[registry]
	if !exists {
		return nil, errors.New("registry not found")
	}
	entry, exists := registryMap[ruleset]
	if !exists {
		return nil, errors.New("ruleset not found")
	}
	return &entry, nil
}

func (f *FileManager) GetEntries(ctx context.Context) (map[string]map[string]Entry, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Rulesets, nil
}

func (f *FileManager) GetRawRegistries(ctx context.Context) (map[string]map[string]interface{}, error) {
	manifest, err := f.loadManifest()
	if err != nil {
		return nil, err
	}
	return manifest.Registries, nil
}

func (f *FileManager) CreateEntry(ctx context.Context, registry, ruleset string, entry Entry) error {
	manifest, err := f.loadManifest()
	if err != nil {
		manifest = &Manifest{
			Registries: make(map[string]map[string]interface{}),
			Rulesets:   make(map[string]map[string]Entry),
		}
	}
	if manifest.Rulesets[registry] == nil {
		manifest.Rulesets[registry] = make(map[string]Entry)
	}
	if _, exists := manifest.Rulesets[registry][ruleset]; exists {
		return errors.New("entry already exists")
	}
	manifest.Rulesets[registry][ruleset] = entry
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateEntry(ctx context.Context, registry, ruleset string, entry Entry) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Rulesets[registry] == nil {
		return errors.New("registry not found")
	}
	if _, exists := manifest.Rulesets[registry][ruleset]; !exists {
		return errors.New("entry not found")
	}
	manifest.Rulesets[registry][ruleset] = entry
	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveEntry(ctx context.Context, registry, ruleset string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}
	if manifest.Rulesets[registry] == nil {
		return errors.New("registry not found")
	}
	if _, exists := manifest.Rulesets[registry][ruleset]; !exists {
		return errors.New("entry not found")
	}
	delete(manifest.Rulesets[registry], ruleset)
	if len(manifest.Rulesets[registry]) == 0 {
		delete(manifest.Rulesets, registry)
	}
	return f.saveManifest(manifest)
}

func (f *FileManager) loadManifest() (*Manifest, error) {
	data, err := os.ReadFile("arm.json")
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (f *FileManager) saveManifest(manifest *Manifest) error {
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
		manifest = &Manifest{
			Registries: make(map[string]map[string]interface{}),
			Rulesets:   make(map[string]map[string]Entry),
		}
	}

	if _, exists := manifest.Registries[name]; exists && !force {
		return errors.New("registry already exists (use --force to overwrite)")
	}

	// Apply default branches if not specified
	if len(config.Branches) == 0 {
		config.Branches = []string{"main", "master"} // Default branches for "latest" resolution
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

	slog.InfoContext(ctx, "Adding git registry configuration",
		"action", "git_registry_add",
		"name", name,
		"url", config.URL,
		"type", config.Type,
		"branches", config.Branches)

	return f.saveManifest(manifest)
}

func (f *FileManager) RemoveRegistry(ctx context.Context, name string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	removedRegistry, exists := manifest.Registries[name]
	if !exists {
		return errors.New("registry not found")
	}

	slog.InfoContext(ctx, "Removing registry configuration",
		"action", "registry_remove",
		"name", name,
		"removed_url", removedRegistry["url"],
		"removed_type", removedRegistry["type"])

	delete(manifest.Registries, name)
	return f.saveManifest(manifest)
}

func (f *FileManager) UpdateGitRegistry(ctx context.Context, name, field, value string) error {
	manifest, err := f.loadManifest()
	if err != nil {
		return err
	}

	regConfig, exists := manifest.Registries[name]
	if !exists {
		return errors.New("registry not found")
	}

	switch field {
	case "url":
		regConfig["url"] = value
	case "type":
		if value != "git" {
			return errors.New("type must be 'git'")
		}
		regConfig["type"] = value
	case "branches":
		branches := strings.Split(value, ",")
		for i, branch := range branches {
			branches[i] = strings.TrimSpace(branch)
		}
		regConfig["branches"] = branches
	default:
		return fmt.Errorf("unknown field '%s' (valid: url, type, branches)", field)
	}

	slog.InfoContext(ctx, "Updating git registry field",
		"action", "git_registry_update",
		"name", name,
		"field", field,
		"value", value)

	return f.saveManifest(manifest)
}
