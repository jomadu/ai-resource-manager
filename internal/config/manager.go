package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const CONFIG_FILE = ".armrc.json"

// Manager handles .armrc.json configuration file operations.
type Manager interface {
	GetSinks(ctx context.Context) (map[string]SinkConfig, error)
	GetSink(ctx context.Context, name string) (*SinkConfig, error)
	AddSink(ctx context.Context, name string, dirs, include, exclude []string, layout string, force bool) error
	UpdateSink(ctx context.Context, name, field, value string) error
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

func (f *FileManager) GetSink(ctx context.Context, name string) (*SinkConfig, error) {
	config, err := f.loadConfig()
	if err != nil {
		return nil, err
	}
	sink, exists := config.Sinks[name]
	if !exists {
		return nil, fmt.Errorf("sink %s not found", name)
	}
	return &sink, nil
}

func (f *FileManager) AddSink(ctx context.Context, name string, dirs, include, exclude []string, layout string, force bool) error {
	if layout == "" {
		layout = "hierarchical"
	}
	config, err := f.loadConfig()
	if err != nil {
		config = &Config{
			Sinks: make(map[string]SinkConfig),
		}
	}

	if _, exists := config.Sinks[name]; exists && !force {
		return fmt.Errorf("sink %s already exists (use --force to overwrite)", name)
	}

	newSink := SinkConfig{
		Directories: dirs,
		Include:     include,
		Exclude:     exclude,
		Layout:      layout,
	}
	config.Sinks[name] = newSink

	slog.InfoContext(ctx, "Adding sink configuration",
		"action", "sink_add",
		"name", name,
		"directories", dirs,
		"include", include,
		"exclude", exclude,
		"layout", layout)

	return f.saveConfig(config)
}

func (f *FileManager) RemoveSink(ctx context.Context, name string) error {
	config, err := f.loadConfig()
	if err != nil {
		return err
	}

	removedSink, exists := config.Sinks[name]
	if !exists {
		return fmt.Errorf("sink %s not found", name)
	}

	slog.InfoContext(ctx, "Removing sink configuration",
		"action", "sink_remove",
		"name", name,
		"removed_directories", removedSink.Directories,
		"removed_include", removedSink.Include,
		"removed_exclude", removedSink.Exclude)

	delete(config.Sinks, name)
	return f.saveConfig(config)
}

func (f *FileManager) UpdateSink(ctx context.Context, name, field, value string) error {
	config, err := f.loadConfig()
	if err != nil {
		return err
	}

	sink, exists := config.Sinks[name]
	if !exists {
		return fmt.Errorf("sink %s not found", name)
	}

	switch field {
	case "directories":
		dirs := strings.Split(value, ",")
		for i, dir := range dirs {
			dirs[i] = strings.TrimSpace(dir)
		}
		sink.Directories = dirs
	case "include":
		include := strings.Split(value, ",")
		for i, pattern := range include {
			include[i] = strings.TrimSpace(pattern)
		}
		sink.Include = include
	case "exclude":
		exclude := strings.Split(value, ",")
		for i, pattern := range exclude {
			exclude[i] = strings.TrimSpace(pattern)
		}
		sink.Exclude = exclude
	case "layout":
		if value != "hierarchical" && value != "flat" {
			return fmt.Errorf("layout must be 'hierarchical' or 'flat'")
		}
		sink.Layout = value
	default:
		return fmt.Errorf("unknown field '%s' (valid: directories, include, exclude, layout)", field)
	}

	config.Sinks[name] = sink

	slog.InfoContext(ctx, "Updating sink field",
		"action", "sink_update",
		"name", name,
		"field", field,
		"value", value)

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
