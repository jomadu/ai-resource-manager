package rcfile

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Service handles reading configuration from .armrc files
type Service struct {
	workingDir  string
	userHomeDir string
}

// NewService creates a new RC file service with default OS paths
func NewService() *Service {
	workingDir, _ := os.Getwd()
	userHomeDir, _ := os.UserHomeDir()
	return &Service{
		workingDir:  workingDir,
		userHomeDir: userHomeDir,
	}
}

// NewServiceWithPaths creates a new RC file service with custom paths for testing
func NewServiceWithPaths(workingDir, userHomeDir string) *Service {
	return &Service{
		workingDir:  workingDir,
		userHomeDir: userHomeDir,
	}
}

// GetSection retrieves all key-value pairs from a section using hierarchical lookup
func (s *Service) GetSection(section string) (map[string]string, error) {
	// Try project .armrc first
	if s.workingDir != "" {
		projectRcPath := filepath.Join(s.workingDir, ".armrc")
		if values, err := s.getSectionFromFile(projectRcPath, section); err == nil {
			return values, nil
		}
	}

	// Fallback to user .armrc
	if s.userHomeDir == "" {
		return nil, fmt.Errorf(".armrc file not found")
	}
	userRcPath := filepath.Join(s.userHomeDir, ".armrc")
	return s.getSectionFromFile(userRcPath, section)
}

// getSectionFromFile retrieves section from specific file
func (s *Service) getSectionFromFile(filePath, section string) (map[string]string, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(".armrc file not found")
		}
		return nil, fmt.Errorf("failed to load .armrc: %w", err)
	}

	sec, err := cfg.GetSection(section)
	if err != nil {
		return nil, fmt.Errorf("section %s not found", section)
	}

	values := make(map[string]string)
	for _, key := range sec.Keys() {
		values[key.Name()] = s.expandEnvVars(key.Value())
	}

	return values, nil
}

// GetValue retrieves a configuration value from .armrc file
func (s *Service) GetValue(section, key string) (string, error) {
	values, err := s.GetSection(section)
	if err != nil {
		return "", err
	}

	value, exists := values[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in section %s", key, section)
	}

	return value, nil
}

// expandEnvVars expands environment variables in the format ${VAR_NAME}
func (s *Service) expandEnvVars(value string) string {
	return os.ExpandEnv(value)
}
