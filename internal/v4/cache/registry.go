package cache

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// RegistryMetadata represents the metadata stored in registry metadata.json
type RegistryMetadata struct {
	Metadata       interface{} `json:"metadata"`
	CreatedOn      time.Time   `json:"created_on"`
	LastUpdatedOn  time.Time   `json:"last_updated_on"`
	LastAccessedOn time.Time   `json:"last_accessed_on"`
}

// RegistryCache handles registry directory and metadata
type RegistryCache struct {
	registryKey interface{}
	registryDir string
}

// NewRegistryCache creates registry directory and metadata.json
func NewRegistryCache(registryKey interface{}) (*RegistryCache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Join(homeDir, ".arm")
	return NewRegistryCacheWithPath(baseDir, registryKey)
}

func NewRegistryCacheWithPath(baseDir string, registryKey interface{}) (*RegistryCache, error) {
	// Generate key from registryKey
	key, err := GenerateKey(registryKey)
	if err != nil {
		return nil, err
	}
	
	// Create registry directory
	registryDir := filepath.Join(baseDir, "registries", key)
	if err := os.MkdirAll(registryDir, 0755); err != nil {
		return nil, err
	}
	
	// Create metadata.json
	now := time.Now().UTC()
	metadata := RegistryMetadata{
		Metadata:       registryKey,
		CreatedOn:      now,
		LastUpdatedOn:  now,
		LastAccessedOn: now,
	}
	
	metadataPath := filepath.Join(registryDir, "metadata.json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return nil, err
	}
	
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return nil, err
	}
	
	return &RegistryCache{
		registryKey: registryKey,
		registryDir: registryDir,
	}, nil
}

// GetRegistryDir returns the registry directory path
func (r *RegistryCache) GetRegistryDir() string {
	return r.registryDir
}

// GetRepoDir returns the git repository directory path
func (r *RegistryCache) GetRepoDir() string {
	return filepath.Join(r.registryDir, "repo")
}

// GetPackagesDir returns the packages directory path
func (r *RegistryCache) GetPackagesDir() string {
	return filepath.Join(r.registryDir, "packages")
}

// UpdateAccessTime updates registry metadata access time
func (r *RegistryCache) UpdateAccessTime(ctx context.Context) error {
	return r.updateTimestamp("last_accessed_on")
}

// UpdateUpdatedTime updates registry metadata updated time
func (r *RegistryCache) UpdateUpdatedTime(ctx context.Context) error {
	return r.updateTimestamp("last_updated_on")
}

// updateTimestamp updates a specific timestamp field in metadata.json
func (r *RegistryCache) updateTimestamp(field string) error {
	metadataPath := filepath.Join(r.registryDir, "metadata.json")
	
	// Read existing metadata
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return err
	}
	
	var metadata RegistryMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return err
	}
	
	// Update timestamp
	now := time.Now().UTC()
	switch field {
	case "last_accessed_on":
		metadata.LastAccessedOn = now
	case "last_updated_on":
		metadata.LastUpdatedOn = now
	}
	
	// Write back to file
	updatedData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(metadataPath, updatedData, 0644)
}