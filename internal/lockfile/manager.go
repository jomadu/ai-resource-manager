package lockfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

// Manager handles arm.lock file operations.
type Manager interface {
	GetEntry(ctx context.Context, registry, ruleset string) (*Entry, error)
	GetEntries(ctx context.Context) (map[string]map[string]Entry, error)
	CreateEntry(ctx context.Context, registry, ruleset string, entry *Entry) error
	UpdateEntry(ctx context.Context, registry, ruleset string, entry *Entry) error
	RemoveEntry(ctx context.Context, registry, ruleset string) error
}

// FileManager implements file-based lock file management.
type FileManager struct {
	lockPath string
}

// NewFileManager creates a new file-based lock file manager.
func NewFileManager() *FileManager {
	return &FileManager{lockPath: "arm.lock"}
}

// NewFileManagerWithPath creates a new file-based lock file manager with custom path.
func NewFileManagerWithPath(lockPath string) *FileManager {
	return &FileManager{lockPath: lockPath}
}

func (f *FileManager) GetEntry(ctx context.Context, registry, ruleset string) (*Entry, error) {
	lockFile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}

	registryEntries, exists := lockFile.Rulesets[registry]
	if !exists {
		return nil, errors.New("registry not found")
	}

	entry, exists := registryEntries[ruleset]
	if !exists {
		return nil, errors.New("ruleset not found")
	}

	return &entry, nil
}

func (f *FileManager) GetEntries(ctx context.Context) (map[string]map[string]Entry, error) {
	lockFile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}
	return lockFile.Rulesets, nil
}

func (f *FileManager) CreateEntry(ctx context.Context, registry, ruleset string, entry *Entry) error {
	lockFile, err := f.readLockFile()
	if err != nil {
		lockFile = &LockFile{Rulesets: make(map[string]map[string]Entry)}
	}

	if lockFile.Rulesets[registry] == nil {
		lockFile.Rulesets[registry] = make(map[string]Entry)
	}

	lockFile.Rulesets[registry][ruleset] = *entry
	return f.writeLockFile(lockFile)
}

func (f *FileManager) UpdateEntry(ctx context.Context, registry, ruleset string, entry *Entry) error {
	lockFile, err := f.readLockFile()
	if err != nil {
		return err
	}

	if lockFile.Rulesets[registry] == nil {
		return errors.New("registry not found")
	}

	lockFile.Rulesets[registry][ruleset] = *entry
	return f.writeLockFile(lockFile)
}

func (f *FileManager) RemoveEntry(ctx context.Context, registry, ruleset string) error {
	lockFile, err := f.readLockFile()
	if err != nil {
		return err
	}

	if lockFile.Rulesets[registry] == nil {
		return errors.New("registry not found")
	}

	delete(lockFile.Rulesets[registry], ruleset)

	// Remove registry if empty
	if len(lockFile.Rulesets[registry]) == 0 {
		delete(lockFile.Rulesets, registry)
	}

	return f.writeLockFile(lockFile)
}

func (f *FileManager) readLockFile() (*LockFile, error) {
	data, err := os.ReadFile(f.lockPath)
	if err != nil {
		return nil, err
	}

	var lockFile LockFile
	err = json.Unmarshal(data, &lockFile)
	if err != nil {
		return nil, err
	}

	return &lockFile, nil
}

func (f *FileManager) writeLockFile(lockFile *LockFile) error {
	data, err := json.MarshalIndent(lockFile, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.lockPath, data, 0o644)
}
