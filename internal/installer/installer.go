package installer

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// Installer manages physical file deployment to sink directories.
type Installer interface {
	Install(ctx context.Context, dir, ruleset, version string, files []types.File) error
	Uninstall(ctx context.Context, dir, ruleset string) error
	ListInstalled(ctx context.Context, dir string) ([]Installation, error)
}

// FileInstaller implements file-based installation to sink directories.
type FileInstaller struct{}

// NewFileInstaller creates a new file-based installer.
func NewFileInstaller() *FileInstaller {
	return &FileInstaller{}
}

func (f *FileInstaller) Install(ctx context.Context, dir, ruleset, version string, files []types.File) error {
	rulesetDir := filepath.Join(dir, "arm", ruleset, version)
	if err := os.MkdirAll(rulesetDir, 0o755); err != nil {
		return err
	}
	slog.InfoContext(ctx, "Created directory", "path", rulesetDir)

	for _, file := range files {
		filePath := filepath.Join(rulesetDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}
	}

	slog.InfoContext(ctx, "Installed files", "count", len(files), "path", rulesetDir)
	return nil
}

func (f *FileInstaller) Uninstall(ctx context.Context, dir, ruleset string) error {
	rulesetDir := filepath.Join(dir, "arm", ruleset)
	if err := os.RemoveAll(rulesetDir); err != nil {
		return err
	}
	slog.InfoContext(ctx, "Removed directory", "path", rulesetDir)
	return nil
}

func (f *FileInstaller) ListInstalled(ctx context.Context, dir string) ([]Installation, error) {
	armDir := filepath.Join(dir, "arm")
	entries, err := os.ReadDir(armDir)
	if err != nil {
		return nil, err
	}

	var installations []Installation
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		rulesetPath := filepath.Join(armDir, entry.Name())
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
