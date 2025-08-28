package config

import (
	"context"
	"errors"
)

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
	return nil, errors.New("not implemented")
}

func (f *FileManager) GetSinks(ctx context.Context) (map[string]SinkConfig, error) {
	return nil, errors.New("not implemented")
}

func (f *FileManager) AddRegistry(ctx context.Context, name, url, registryType string) error {
	return errors.New("not implemented")
}

func (f *FileManager) AddSink(ctx context.Context, name string, dirs, include, exclude []string) error {
	return errors.New("not implemented")
}

func (f *FileManager) RemoveRegistry(ctx context.Context, name string) error {
	return errors.New("not implemented")
}

func (f *FileManager) RemoveSink(ctx context.Context, name string) error {
	return errors.New("not implemented")
}
