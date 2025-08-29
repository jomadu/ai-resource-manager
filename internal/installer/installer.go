package installer

import (
	"context"
	"os"
	"path/filepath"
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
	rulesetDir := filepath.Join(dir, ruleset, version)
	if err := os.MkdirAll(rulesetDir, 0o755); err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(rulesetDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func (f *FileInstaller) Uninstall(ctx context.Context, dir, ruleset string) error {
	rulesetDir := filepath.Join(dir, ruleset)
	return os.RemoveAll(rulesetDir)
}

func (f *FileInstaller) ListInstalled(ctx context.Context, dir string) ([]Installation, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var installations []Installation
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		rulesetPath := filepath.Join(dir, entry.Name())
		versionEntries, err := os.ReadDir(rulesetPath)
		if err != nil {
			continue
		}

		for _, versionEntry := range versionEntries {
			if !versionEntry.IsDir() {
				continue
			}

			installations = append(installations, Installation{
				Ruleset: entry.Name(),
				Version: versionEntry.Name(),
				Path:    filepath.Join(rulesetPath, versionEntry.Name()),
			})
		}
	}

	return installations, nil
}
