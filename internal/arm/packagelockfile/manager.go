package packagelockfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type LockFile struct {
	Version      int                              `json:"version"`
	Dependencies map[string]DependencyLockConfig `json:"dependencies,omitempty"`
}

type DependencyLockConfig struct {
	Integrity string `json:"integrity"`
}

type FileManager struct {
	lockPath string
}

func NewFileManager() *FileManager {
	return &FileManager{lockPath: "arm-lock.json"}
}

func NewFileManagerWithPath(lockPath string) *FileManager {
	return &FileManager{lockPath: lockPath}
}

func (f *FileManager) GetDependencyLock(ctx context.Context, registry, packageName, version string) (*DependencyLockConfig, error) {
	key := lockKey(registry, packageName, version)
	lockfile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}

	lockInfo, exists := lockfile.Dependencies[key]
	if !exists {
		return nil, errors.New("dependency not found")
	}

	return &lockInfo, nil
}

func (f *FileManager) GetLockFile(ctx context.Context) (*LockFile, error) {
	return f.readLockFile()
}

func (f *FileManager) UpsertDependencyLock(ctx context.Context, registry, packageName, version string, config *DependencyLockConfig) error {
	key := lockKey(registry, packageName, version)
	var lockfile *LockFile
	var err error

	// Try to read existing lockfile
	lockfile, err = f.readLockFile()
	if err != nil {
		// If file doesn't exist, create new lockfile
		if os.IsNotExist(err) {
			lockfile = &LockFile{
				Version:      1,
				Dependencies: make(map[string]DependencyLockConfig),
			}
		} else {
			return err
		}
	}

	// Upsert the dependency
	lockfile.Dependencies[key] = *config

	// Write lockfile back to disk
	return f.writeLockFile(lockfile)
}

func (f *FileManager) RemoveDependencyLock(ctx context.Context, registry, packageName, version string) error {
	key := lockKey(registry, packageName, version)
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	_, exists := lockfile.Dependencies[key]
	if !exists {
		return errors.New("dependency not found")
	}

	// Remove the dependency
	delete(lockfile.Dependencies, key)

	// Delete lockfile if empty
	if f.isLockfileEmpty(lockfile) {
		return f.deleteLockFile()
	}

	return f.writeLockFile(lockfile)
}

func (f *FileManager) UpdateRegistryName(ctx context.Context, oldName, newName string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	// Update dependencies with new registry name
	for key, lockInfo := range lockfile.Dependencies {
		if strings.HasPrefix(key, oldName+"/") {
			newKey := newName + key[len(oldName):]
			lockfile.Dependencies[newKey] = lockInfo
			delete(lockfile.Dependencies, key)
		}
	}

	return f.writeLockFile(lockfile)
}

// readLockFile reads and parses the lockfile from disk
func (f *FileManager) readLockFile() (*LockFile, error) {
	data, err := os.ReadFile(f.lockPath)
	if err != nil {
		// Return error if file doesn't exist (no auto-create)
		return nil, err
	}

	var lockfile LockFile
	err = json.Unmarshal(data, &lockfile)
	if err != nil {
		return nil, err
	}

	// Initialize nil map
	if lockfile.Dependencies == nil {
		lockfile.Dependencies = make(map[string]DependencyLockConfig)
	}

	return &lockfile, nil
}

// writeLockFile writes the lockfile to disk
func (f *FileManager) writeLockFile(lockfile *LockFile) error {
	data, err := json.MarshalIndent(lockfile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.lockPath, data, 0644)
}

// isLockfileEmpty checks if the lockfile has no packages
func (f *FileManager) isLockfileEmpty(lockfile *LockFile) bool {
	return len(lockfile.Dependencies) == 0
}

// deleteLockFile removes the lockfile from disk
func (f *FileManager) deleteLockFile() error {
	return os.Remove(f.lockPath)
}

// lockKey creates a lock key in format "registry/package@version"
func lockKey(registry, packageName, version string) string {
	return registry + "/" + packageName + "@" + version
}