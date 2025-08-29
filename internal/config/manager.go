package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

const CONFIG_FILE = ".armrc.json"

// Manager handles .armrc.json configuration file operations.
type Manager interface {
	GetRegistries(ctx context.Context) (map[string]RegistryConfig, error)
	GetSinks(ctx context.Context) (map[string]SinkConfig, error)
	AddRegistry(ctx context.Context, name, url, registryType string) error
	AddSink(ctx context.Context, name string, dirs, include, exclude []string) error
	RemoveRegistry(ctx context.Context, name string) error
	RemoveSink(ctx context.Context, name string) error
}

// FileManager implements file-based configuration management.
type FileManager struct{}

// NewFileManager creates a new file-based configuration manager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f *FileManager) GetRegistries(ctx context.Context) (map[string]RegistryConfig, error) {
	config, err := f.loadConfig()
	if err != nil {
		return nil, err
	}
	return config.Registries, nil
}

func (f *FileManager) GetSinks(ctx context.Context) (map[string]SinkConfig, error) {
	config, err := f.loadConfig()
	if err != nil {
		return nil, err
	}
	return config.Sinks, nil
}

func (f *FileManager) AddRegistry(ctx context.Context, name, url, registryType string) error {
	config, err := f.loadConfig()
	if err != nil {
		config = &Config{
			Registries: make(map[string]RegistryConfig),
			Sinks:      make(map[string]SinkConfig),
		}
	}

	if _, exists := config.Registries[name]; exists {
		return fmt.Errorf("registry %s already exists", name)
	}

	config.Registries[name] = RegistryConfig{
		URL:  url,
		Type: registryType,
	}

	return f.saveConfig(config)
}

func (f *FileManager) AddSink(ctx context.Context, name string, dirs, include, exclude []string) error {
	config, err := f.loadConfig()
	if err != nil {
		config = &Config{
			Registries: make(map[string]RegistryConfig),
			Sinks:      make(map[string]SinkConfig),
		}
	}

	if _, exists := config.Sinks[name]; exists {
		return fmt.Errorf("sink %s already exists", name)
	}

	config.Sinks[name] = SinkConfig{
		Directories: dirs,
		Include:     include,
		Exclude:     exclude,
	}

	return f.saveConfig(config)
}

func (f *FileManager) RemoveRegistry(ctx context.Context, name string) error {
	config, err := f.loadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Registries[name]; !exists {
		return fmt.Errorf("registry %s not found", name)
	}

	delete(config.Registries, name)
	return f.saveConfig(config)
}

func (f *FileManager) RemoveSink(ctx context.Context, name string) error {
	config, err := f.loadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Sinks[name]; !exists {
		return fmt.Errorf("sink %s not found", name)
	}

	delete(config.Sinks, name)
	return f.saveConfig(config)
}

func (f *FileManager) loadConfig() (*Config, error) {
	data, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (f *FileManager) saveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(CONFIG_FILE, data, 0o644)
}
