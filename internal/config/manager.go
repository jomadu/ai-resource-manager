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
	GetSinks(ctx context.Context) (map[string]SinkConfig, error)
	AddSink(ctx context.Context, name string, dirs, include, exclude []string) error
	RemoveSink(ctx context.Context, name string) error
}

// FileManager implements file-based configuration management.
type FileManager struct{}

// NewFileManager creates a new file-based configuration manager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f *FileManager) GetSinks(ctx context.Context) (map[string]SinkConfig, error) {
	config, err := f.loadConfig()
	if err != nil {
		return nil, err
	}
	return config.Sinks, nil
}

func (f *FileManager) AddSink(ctx context.Context, name string, dirs, include, exclude []string) error {
	config, err := f.loadConfig()
	if err != nil {
		config = &Config{
			Sinks: make(map[string]SinkConfig),
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
