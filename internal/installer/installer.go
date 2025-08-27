package installer

import (
	"context"
	"errors"
)

// Installer manages physical file deployment to sink directories.
type Installer interface {
	Install(ctx context.Context, dir, ruleset, version string, files []File) error
	Uninstall(ctx context.Context, dir, ruleset string) error
	ListInstalled(ctx context.Context, dir string) ([]Installation, error)
}

// FileInstaller implements file-based installation to sink directories.
type FileInstaller struct{}

// NewFileInstaller creates a new file-based installer.
func NewFileInstaller() *FileInstaller {
	return &FileInstaller{}
}

func (f *FileInstaller) Install(ctx context.Context, dir, ruleset, version string, files []File) error {
	return errors.New("not implemented")
}

func (f *FileInstaller) Uninstall(ctx context.Context, dir, ruleset string) error {
	return errors.New("not implemented")
}

func (f *FileInstaller) ListInstalled(ctx context.Context, dir string) ([]Installation, error) {
	return nil, errors.New("not implemented")
}