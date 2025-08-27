package manifest

import (
	"context"
	"errors"
)

// Manager handles arm.json manifest file operations.
type Manager interface {
	GetEntry(ctx context.Context, registry, ruleset string) (*Entry, error)
	GetEntries(ctx context.Context) (map[string]map[string]Entry, error)
	CreateEntry(ctx context.Context, registry, ruleset string, entry Entry) error
	UpdateEntry(ctx context.Context, registry, ruleset string, entry Entry) error
	RemoveEntry(ctx context.Context, registry, ruleset string) error
}

// FileManager implements file-based manifest management.
type FileManager struct{}

// NewFileManager creates a new file-based manifest manager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f *FileManager) GetEntry(ctx context.Context, registry, ruleset string) (*Entry, error) {
	return nil, errors.New("not implemented")
}

func (f *FileManager) GetEntries(ctx context.Context) (map[string]map[string]Entry, error) {
	return nil, errors.New("not implemented")
}

func (f *FileManager) CreateEntry(ctx context.Context, registry, ruleset string, entry Entry) error {
	return errors.New("not implemented")
}

func (f *FileManager) UpdateEntry(ctx context.Context, registry, ruleset string, entry Entry) error {
	return errors.New("not implemented")
}

func (f *FileManager) RemoveEntry(ctx context.Context, registry, ruleset string) error {
	return errors.New("not implemented")
}