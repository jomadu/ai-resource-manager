package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RegistryMetadata represents the metadata stored in registry metadata.json
type RegistryMetadata struct {
	URL        string `json:"url"`
	Type       string `json:"type"`
	GroupID    string `json:"group_id,omitempty"`
	ProjectID  string `json:"project_id,omitempty"`
	Owner      string `json:"owner,omitempty"`
	Repository string `json:"repository,omitempty"`
}

// Registry handles registry directory and metadata with cross-process locking
type Registry struct {
	registryKey interface{}
	registryDir string
	lock        *FileLock // Protects metadata.json operations
}

// NewRegistry creates registry directory and metadata.json
func NewRegistry(registryKey interface{}) (*Registry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Join(homeDir, ".arm")
	return NewRegistryWithPath(baseDir, registryKey)
}

func NewRegistryWithPath(baseDir string, registryKey interface{}) (*Registry, error) {
	// Generate key from registryKey
	key, err := GenerateKey(registryKey)
	if err != nil {
		return nil, err
	}

	// Create registry directory
	registryDir := filepath.Join(baseDir, "storage", "registries", key)
	if err := os.MkdirAll(registryDir, 0o755); err != nil {
		return nil, err
	}

	// Create metadata.json from registryKey fields
	var metadata RegistryMetadata
	if keyMap, ok := registryKey.(map[string]interface{}); ok {
		if url, ok := keyMap["url"].(string); ok {
			metadata.URL = url
		}
		if typ, ok := keyMap["type"].(string); ok {
			metadata.Type = typ
		}
		if groupID, ok := keyMap["group_id"].(string); ok {
			metadata.GroupID = groupID
		}
		if projectID, ok := keyMap["project_id"].(string); ok {
			metadata.ProjectID = projectID
		}
		if owner, ok := keyMap["owner"].(string); ok {
			metadata.Owner = owner
		}
		if repository, ok := keyMap["repository"].(string); ok {
			metadata.Repository = repository
		}
	}

	metadataPath := filepath.Join(registryDir, "metadata.json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(metadataPath, data, 0o644); err != nil {
		return nil, err
	}

	return &Registry{
		registryKey: registryKey,
		registryDir: registryDir,
		lock:        NewFileLock(registryDir),
	}, nil
}

// GetRegistryDir returns the registry directory path
func (r *Registry) GetRegistryDir() string {
	return r.registryDir
}

// GetRepoDir returns the git repository directory path
func (r *Registry) GetRepoDir() string {
	return filepath.Join(r.registryDir, "repo")
}

// GetPackagesDir returns the packages directory path
func (r *Registry) GetPackagesDir() string {
	return filepath.Join(r.registryDir, "packages")
}
