package config

import (
	"encoding/json"
	"os"
)

// ConfigManager handles three-tier project configuration system
type ConfigManager interface {
	// LoadInfraConfig loads .armrc.json (registry definitions and sink mappings)
	LoadInfraConfig() (*InfraConfig, error)

	// SaveInfraConfig saves .armrc.json
	SaveInfraConfig(config *InfraConfig) error

	// LoadManifest loads arm.json (ruleset dependencies)
	LoadManifest() (*Manifest, error)

	// SaveManifest saves arm.json
	SaveManifest(manifest *Manifest) error

	// LoadLockFile loads arm.lock (resolved versions)
	LoadLockFile() (*LockFile, error)

	// SaveLockFile saves arm.lock
	SaveLockFile(lockFile *LockFile) error
}

// InfraConfig represents .armrc.json - "Where can I find rulesets and where should I install them?"
type InfraConfig struct {
	Registries map[string]*RegistryConfig `json:"registries"`
	Sinks      map[string]*SinkConfig     `json:"sinks"`
}

// RegistryConfig defines a registry
type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// SinkConfig defines where to install rulesets for AI tools
type SinkConfig struct {
	Directories []string `json:"directories"`
	Rulesets    []string `json:"rulesets"`
}

// Manifest represents arm.json - "What rulesets do I want and what versions?"
type Manifest struct {
	Rulesets map[string]map[string]*ManifestEntry `json:"rulesets"`
}

// ManifestEntry represents a ruleset dependency
type ManifestEntry struct {
	Version  string   `json:"version"`
	Patterns []string `json:"patterns"`
}

// LockFile represents arm.lock - "Exactly what was installed and from where?"
type LockFile struct {
	Rulesets map[string]map[string]*LockEntry `json:"rulesets"`
}

// LockEntry represents a locked ruleset version
type LockEntry struct {
	URL        string   `json:"url"`
	Type       string   `json:"type"`
	Constraint string   `json:"constraint"`
	Resolved   string   `json:"resolved"`
	Patterns   []string `json:"patterns"`
}

// FileConfigManager implements ConfigManager for file-based config
type FileConfigManager struct {
	infraPath    string
	manifestPath string
	lockPath     string
}

func NewFileConfigManager(infraPath, manifestPath, lockPath string) *FileConfigManager {
	return &FileConfigManager{
		infraPath:    infraPath,
		manifestPath: manifestPath,
		lockPath:     lockPath,
	}
}

func (f *FileConfigManager) LoadInfraConfig() (*InfraConfig, error) {
	data, err := os.ReadFile(f.infraPath)
	if os.IsNotExist(err) {
		return &InfraConfig{Registries: make(map[string]*RegistryConfig), Sinks: make(map[string]*SinkConfig)}, nil
	}
	if err != nil {
		return nil, err
	}
	var config InfraConfig
	err = json.Unmarshal(data, &config)
	return &config, err
}

func (f *FileConfigManager) SaveInfraConfig(config *InfraConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.infraPath, data, 0644)
}

func (f *FileConfigManager) LoadManifest() (*Manifest, error) {
	data, err := os.ReadFile(f.manifestPath)
	if os.IsNotExist(err) {
		return &Manifest{Rulesets: make(map[string]map[string]*ManifestEntry)}, nil
	}
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	return &manifest, err
}

func (f *FileConfigManager) SaveManifest(manifest *Manifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.manifestPath, data, 0644)
}

func (f *FileConfigManager) LoadLockFile() (*LockFile, error) {
	data, err := os.ReadFile(f.lockPath)
	if os.IsNotExist(err) {
		return &LockFile{Rulesets: make(map[string]map[string]*LockEntry)}, nil
	}
	if err != nil {
		return nil, err
	}
	var lockFile LockFile
	err = json.Unmarshal(data, &lockFile)
	return &lockFile, err
}

func (f *FileConfigManager) SaveLockFile(lockFile *LockFile) error {
	data, err := json.MarshalIndent(lockFile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.lockPath, data, 0644)
}
