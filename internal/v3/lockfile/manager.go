package lockfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

// Manager handles arm-lock.json file operations.
type Manager interface {
	// Ruleset operations
	GetRuleset(ctx context.Context, registry, ruleset string) (*Entry, error)
	GetRulesets(ctx context.Context) (map[string]map[string]Entry, error)
	CreateOrUpdateRuleset(ctx context.Context, registry, ruleset string, entry *Entry) error
	RemoveRuleset(ctx context.Context, registry, ruleset string) error

	// Promptset operations
	GetPromptset(ctx context.Context, registry, promptset string) (*Entry, error)
	GetPromptsets(ctx context.Context) (map[string]map[string]Entry, error)
	CreateOrUpdatePromptset(ctx context.Context, registry, promptset string, entry *Entry) error
	RemovePromptset(ctx context.Context, registry, promptset string) error

	// General operations
	GetAllEntries(ctx context.Context) (*LockFile, error)
}

// FileManager implements file-based lock file management.
type FileManager struct {
	lockPath string
}

// NewFileManager creates a new file-based lock file manager.
func NewFileManager() *FileManager {
	return &FileManager{lockPath: "arm-lock.json"}
}

// NewFileManagerWithPath creates a new file-based lock file manager with custom path.
func NewFileManagerWithPath(lockPath string) *FileManager {
	return &FileManager{lockPath: lockPath}
}

func (f *FileManager) GetRuleset(ctx context.Context, registry, ruleset string) (*Entry, error) {
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

func (f *FileManager) GetPromptset(ctx context.Context, registry, promptset string) (*Entry, error) {
	lockFile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}

	registryEntries, exists := lockFile.Promptsets[registry]
	if !exists {
		return nil, errors.New("registry not found")
	}

	entry, exists := registryEntries[promptset]
	if !exists {
		return nil, errors.New("promptset not found")
	}

	return &entry, nil
}

func (f *FileManager) GetRulesets(ctx context.Context) (map[string]map[string]Entry, error) {
	lockFile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}
	return lockFile.Rulesets, nil
}

func (f *FileManager) GetPromptsets(ctx context.Context) (map[string]map[string]Entry, error) {
	lockFile, err := f.readLockFile()
	if err != nil {
		return nil, err
	}
	return lockFile.Promptsets, nil
}

func (f *FileManager) CreateOrUpdateRuleset(ctx context.Context, registry, ruleset string, entry *Entry) error {
	lockFile, err := f.readLockFile()
	if err != nil {
		lockFile = &LockFile{
			Rulesets:   make(map[string]map[string]Entry),
			Promptsets: make(map[string]map[string]Entry),
		}
	}

	if lockFile.Rulesets[registry] == nil {
		lockFile.Rulesets[registry] = make(map[string]Entry)
	}

	lockFile.Rulesets[registry][ruleset] = *entry
	return f.writeLockFile(lockFile)
}

func (f *FileManager) CreateOrUpdatePromptset(ctx context.Context, registry, promptset string, entry *Entry) error {
	lockFile, err := f.readLockFile()
	if err != nil {
		lockFile = &LockFile{
			Rulesets:   make(map[string]map[string]Entry),
			Promptsets: make(map[string]map[string]Entry),
		}
	}

	if lockFile.Promptsets[registry] == nil {
		lockFile.Promptsets[registry] = make(map[string]Entry)
	}

	lockFile.Promptsets[registry][promptset] = *entry
	return f.writeLockFile(lockFile)
}

func (f *FileManager) RemoveRuleset(ctx context.Context, registry, ruleset string) error {
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

func (f *FileManager) RemovePromptset(ctx context.Context, registry, promptset string) error {
	lockFile, err := f.readLockFile()
	if err != nil {
		return err
	}

	if lockFile.Promptsets[registry] == nil {
		return errors.New("registry not found")
	}

	delete(lockFile.Promptsets[registry], promptset)

	// Remove registry if empty
	if len(lockFile.Promptsets[registry]) == 0 {
		delete(lockFile.Promptsets, registry)
	}

	return f.writeLockFile(lockFile)
}

func (f *FileManager) GetAllEntries(ctx context.Context) (*LockFile, error) {
	return f.readLockFile()
}

func (f *FileManager) readLockFile() (*LockFile, error) {
	data, err := os.ReadFile(f.lockPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty lock file if it doesn't exist
			return &LockFile{
				Rulesets:   make(map[string]map[string]Entry),
				Promptsets: make(map[string]map[string]Entry),
			}, nil
		}
		return nil, err
	}

	var lockFile LockFile
	err = json.Unmarshal(data, &lockFile)
	if err != nil {
		return nil, err
	}

	// Ensure both sections exist even if they're empty
	if lockFile.Rulesets == nil {
		lockFile.Rulesets = make(map[string]map[string]Entry)
	}
	if lockFile.Promptsets == nil {
		lockFile.Promptsets = make(map[string]map[string]Entry)
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
