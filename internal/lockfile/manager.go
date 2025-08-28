package lockfile

import (
	"context"
	"errors"
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
type FileManager struct{}

// NewFileManager creates a new file-based lock file manager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f *FileManager) GetEntry(ctx context.Context, registry, ruleset string) (*Entry, error) {
	return nil, errors.New("not implemented")
}

func (f *FileManager) GetEntries(ctx context.Context) (map[string]map[string]Entry, error) {
	return nil, errors.New("not implemented")
}

func (f *FileManager) CreateEntry(ctx context.Context, registry, ruleset string, entry *Entry) error {
	return errors.New("not implemented")
}

func (f *FileManager) UpdateEntry(ctx context.Context, registry, ruleset string, entry *Entry) error {
	return errors.New("not implemented")
}

func (f *FileManager) RemoveEntry(ctx context.Context, registry, ruleset string) error {
	return errors.New("not implemented")
}
