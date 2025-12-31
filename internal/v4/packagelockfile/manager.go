package packagelockfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

type PackageLockfile struct {
	Version    int                            `json:"version"`
	Rulesets   map[string]*PackageLockInfo   `json:"rulesets,omitempty"`
	Promptsets map[string]*PackageLockInfo   `json:"promptsets,omitempty"`
}

type PackageLockInfo struct {
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

func (f *FileManager) GetRulesetLockInfo(ctx context.Context, packageKey string) (*PackageLockInfo, error) {
	lockfile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}

	lockInfo, exists := lockfile.Rulesets[packageKey]
	if !exists {
		return nil, errors.New("ruleset not found")
	}

	return lockInfo, nil
}

func (f *FileManager) GetPromptsetLockInfo(ctx context.Context, packageKey string) (*PackageLockInfo, error) {
	lockfile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}

	lockInfo, exists := lockfile.Promptsets[packageKey]
	if !exists {
		return nil, errors.New("promptset not found")
	}

	return lockInfo, nil
}

func (f *FileManager) GetPackageLockfile(ctx context.Context) (*PackageLockfile, error) {
	return f.readLockFile()
}

func (f *FileManager) UpsertRulesetLockInfo(ctx context.Context, packageKey string, lockInfo *PackageLockInfo) error {
	var lockfile *PackageLockfile
	var err error

	// Try to read existing lockfile
	lockfile, err = f.readLockFile()
	if err != nil {
		// If file doesn't exist, create new lockfile
		if os.IsNotExist(err) {
			lockfile = &PackageLockfile{
				Version:    1,
				Rulesets:   make(map[string]*PackageLockInfo),
				Promptsets: make(map[string]*PackageLockInfo),
			}
		} else {
			return err
		}
	}

	// Upsert the ruleset
	lockfile.Rulesets[packageKey] = lockInfo

	// Write lockfile back to disk
	return f.writeLockFile(lockfile)
}

func (f *FileManager) UpsertPromptsetLockInfo(ctx context.Context, packageKey string, lockInfo *PackageLockInfo) error {
	var lockfile *PackageLockfile
	var err error

	// Try to read existing lockfile
	lockfile, err = f.readLockFile()
	if err != nil {
		// If file doesn't exist, create new lockfile
		if os.IsNotExist(err) {
			lockfile = &PackageLockfile{
				Version:    1,
				Rulesets:   make(map[string]*PackageLockInfo),
				Promptsets: make(map[string]*PackageLockInfo),
			}
		} else {
			return err
		}
	}

	// Upsert the promptset
	lockfile.Promptsets[packageKey] = lockInfo

	// Write lockfile back to disk
	return f.writeLockFile(lockfile)
}

func (f *FileManager) RemoveRulesetLockInfo(ctx context.Context, packageKey string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	_, exists := lockfile.Rulesets[packageKey]
	if !exists {
		return errors.New("ruleset not found")
	}

	// Remove the ruleset
	delete(lockfile.Rulesets, packageKey)

	// Delete lockfile if empty
	if f.isLockfileEmpty(lockfile) {
		return f.deleteLockFile()
	}

	return f.writeLockFile(lockfile)
}

func (f *FileManager) RemovePromptsetLockInfo(ctx context.Context, packageKey string) error {
	lockfile, err := f.readLockFile()
	if err != nil {
		return err
	}

	_, exists := lockfile.Promptsets[packageKey]
	if !exists {
		return errors.New("promptset not found")
	}

	// Remove the promptset
	delete(lockfile.Promptsets, packageKey)

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

	// Update rulesets
	for key, lockInfo := range lockfile.Rulesets {
		regName, pkgName := core.ParsePackageKey(key)
		if regName == oldName {
			newKey := core.PackageKey(newName, pkgName)
			lockfile.Rulesets[newKey] = lockInfo
			delete(lockfile.Rulesets, key)
		}
	}

	// Update promptsets
	for key, lockInfo := range lockfile.Promptsets {
		regName, pkgName := core.ParsePackageKey(key)
		if regName == oldName {
			newKey := core.PackageKey(newName, pkgName)
			lockfile.Promptsets[newKey] = lockInfo
			delete(lockfile.Promptsets, key)
		}
	}

	return f.writeLockFile(lockfile)
}

// readLockFile reads and parses the lockfile from disk
func (f *FileManager) readLockFile() (*PackageLockfile, error) {
	data, err := os.ReadFile(f.lockPath)
	if err != nil {
		// Return error if file doesn't exist (no auto-create)
		return nil, err
	}

	var lockfile PackageLockfile
	err = json.Unmarshal(data, &lockfile)
	if err != nil {
		return nil, err
	}

	// Initialize nil maps
	if lockfile.Rulesets == nil {
		lockfile.Rulesets = make(map[string]*PackageLockInfo)
	}
	if lockfile.Promptsets == nil {
		lockfile.Promptsets = make(map[string]*PackageLockInfo)
	}

	return &lockfile, nil
}

// writeLockFile writes the lockfile to disk
func (f *FileManager) writeLockFile(lockfile *PackageLockfile) error {
	data, err := json.MarshalIndent(lockfile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.lockPath, data, 0644)
}

// isLockfileEmpty checks if the lockfile has no packages
func (f *FileManager) isLockfileEmpty(lockfile *PackageLockfile) bool {
	return len(lockfile.Rulesets) == 0 && len(lockfile.Promptsets) == 0
}

// deleteLockFile removes the lockfile from disk
func (f *FileManager) deleteLockFile() error {
	return os.Remove(f.lockPath)
}