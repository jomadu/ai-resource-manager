package rc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Manager interface {
	GetAllSections(ctx context.Context) (map[string]map[string]string, error)
	GetSection(ctx context.Context, section string) (map[string]string, error)
	GetValue(ctx context.Context, section, key string) (string, error)
}

// FileManager implements Manager interface for reading .armrc files
type FileManager struct {
	workingDir  string
	userHomeDir string
}

// NewFileManager creates a new file manager with default OS paths
func NewFileManager() *FileManager {
	workingDir, _ := os.Getwd()
	userHomeDir, _ := os.UserHomeDir()
	return &FileManager{
		workingDir:  workingDir,
		userHomeDir: userHomeDir,
	}
}

// NewFileManagerWithPaths creates a new file manager with custom paths for testing
func NewFileManagerWithPaths(workingDir, userHomeDir string) *FileManager {
	return &FileManager{
		workingDir:  workingDir,
		userHomeDir: userHomeDir,
	}
}

// GetAllSections retrieves all sections from .armrc files
// Looks in project .armrc first, then user home .armrc
// Project sections completely override user home sections with the same name
func (f *FileManager) GetAllSections(ctx context.Context) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// Load user home file first (base)
	if f.userHomeDir != "" {
		userRcPath := filepath.Join(f.userHomeDir, ".armrc")
		userSections, err := f.getAllSectionsFromFile(userRcPath)
		if err != nil {
			return nil, err
		}
		// Copy user sections to result
		for section, values := range userSections {
			result[section] = values
		}
	}

	// Load project file and override user home sections
	if f.workingDir != "" {
		projectRcPath := filepath.Join(f.workingDir, ".armrc")
		projectSections, err := f.getAllSectionsFromFile(projectRcPath)
		if err != nil {
			return nil, err
		}
		// Project sections override user home sections
		for section, values := range projectSections {
			result[section] = values
		}
	}

	return result, nil
}

// GetSection retrieves all key-value pairs from a section using hierarchical lookup
// Looks in project .armrc first, then user home .armrc
func (f *FileManager) GetSection(ctx context.Context, section string) (map[string]string, error) {
	// Try project .armrc first
	if f.workingDir != "" {
		projectRcPath := filepath.Join(f.workingDir, ".armrc")
		values, err := f.getSectionFromFile(projectRcPath, section)
		if err == nil {
			return values, nil
		}
		// If file doesn't exist, continue to user home
		// If section not found, continue to user home
	}

	// Fallback to user home .armrc
	if f.userHomeDir == "" {
		return nil, fmt.Errorf("section not found: %s", section)
	}
	userRcPath := filepath.Join(f.userHomeDir, ".armrc")
	return f.getSectionFromFile(userRcPath, section)
}

// GetValue retrieves a configuration value from .armrc file
func (f *FileManager) GetValue(ctx context.Context, section, key string) (string, error) {
	values, err := f.GetSection(ctx, section)
	if err != nil {
		return "", err
	}

	value, exists := values[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in section %s", key, section)
	}

	return value, nil
}

// Helper functions

// loadFileIfExists loads an INI file if it exists, returns nil if file doesn't exist
func (f *FileManager) loadFileIfExists(filePath string) (*ini.File, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to load .armrc: %w", err)
	}
	return cfg, nil
}

// getAllSectionsFromFile retrieves all sections from a specific INI file
func (f *FileManager) getAllSectionsFromFile(filePath string) (map[string]map[string]string, error) {
	cfg, err := f.loadFileIfExists(filePath)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return make(map[string]map[string]string), nil
	}

	sections := make(map[string]map[string]string)
	for _, section := range cfg.Sections() {
		// Skip default section if empty
		if section.Name() == ini.DEFAULT_SECTION && len(section.Keys()) == 0 {
			continue
		}

		values := make(map[string]string)
		for _, key := range section.Keys() {
			values[key.Name()] = f.expandEnvVars(key.Value())
		}
		if len(values) > 0 {
			sections[section.Name()] = values
		}
	}

	return sections, nil
}

// getSectionFromFile retrieves a section from a specific INI file
func (f *FileManager) getSectionFromFile(filePath, section string) (map[string]string, error) {
	cfg, err := f.loadFileIfExists(filePath)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("section not found: %s", section)
	}

	sec, err := cfg.GetSection(section)
	if err != nil {
		return nil, fmt.Errorf("section not found: %s", section)
	}

	values := make(map[string]string)
	for _, key := range sec.Keys() {
		values[key.Name()] = f.expandEnvVars(key.Value())
	}

	return values, nil
}

// expandEnvVars expands environment variables in the format ${VAR_NAME}
func (f *FileManager) expandEnvVars(value string) string {
	return os.ExpandEnv(value)
}
