package installer

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// HierarchicalInstaller implements hierarchical file installation preserving directory structure.
type HierarchicalInstaller struct{}

// NewHierarchicalInstaller creates a new hierarchical installer.
func NewHierarchicalInstaller() *HierarchicalInstaller {
	return &HierarchicalInstaller{}
}

func (h *HierarchicalInstaller) Install(ctx context.Context, dir, registry, ruleset, version string, files []types.File) error {
	// Remove existing versions before installing new one
	if err := h.removeExistingVersions(ctx, dir, registry, ruleset); err != nil {
		return err
	}

	rulesetDir := filepath.Join(dir, "arm", registry, ruleset, version)
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

	slog.InfoContext(ctx, "Installed files (hierarchical)", "count", len(files), "path", rulesetDir)
	return nil
}

func (h *HierarchicalInstaller) Uninstall(ctx context.Context, dir, registry, ruleset string) error {
	rulesetDir := filepath.Join(dir, "arm", registry, ruleset)
	if err := os.RemoveAll(rulesetDir); err != nil {
		return err
	}
	slog.InfoContext(ctx, "Removed directory", "path", rulesetDir)

	// Clean up empty parent directories
	registryDir := filepath.Join(dir, "arm", registry)
	if isEmpty(registryDir) {
		if err := os.Remove(registryDir); err != nil {
			return err
		}
		slog.InfoContext(ctx, "Removed empty registry directory", "path", registryDir)

		armDir := filepath.Join(dir, "arm")
		if isEmpty(armDir) {
			if err := os.Remove(armDir); err != nil {
				return err
			}
			slog.InfoContext(ctx, "Removed empty arm directory", "path", armDir)
		}
	}

	return nil
}

func (h *HierarchicalInstaller) ListInstalled(ctx context.Context, dir string) ([]Installation, error) {
	armDir := filepath.Join(dir, "arm")
	registryEntries, err := os.ReadDir(armDir)
	if err != nil {
		return nil, err
	}

	var installations []Installation
	for _, registryEntry := range registryEntries {
		if !registryEntry.IsDir() {
			continue
		}

		registryPath := filepath.Join(armDir, registryEntry.Name())
		rulesetEntries, err := os.ReadDir(registryPath)
		if err != nil {
			continue
		}

		for _, rulesetEntry := range rulesetEntries {
			if !rulesetEntry.IsDir() {
				continue
			}

			rulesetPath := filepath.Join(registryPath, rulesetEntry.Name())
			versionEntries, err := os.ReadDir(rulesetPath)
			if err != nil {
				continue
			}

			for _, versionEntry := range versionEntries {
				if !versionEntry.IsDir() {
					continue
				}

				installations = append(installations, Installation{
					Ruleset: rulesetEntry.Name(),
					Version: versionEntry.Name(),
					Path:    filepath.Join(rulesetPath, versionEntry.Name()),
				})
			}
		}
	}

	return installations, nil
}

func (h *HierarchicalInstaller) IsInstalled(ctx context.Context, dir, registry, ruleset string) (installed bool, version string, err error) {
	rulesetPath := filepath.Join(dir, "arm", registry, ruleset)
	versionEntries, err := os.ReadDir(rulesetPath)
	if err != nil {
		return false, "", nil // Directory doesn't exist, not installed
	}

	for _, versionEntry := range versionEntries {
		if versionEntry.IsDir() {
			return true, versionEntry.Name(), nil
		}
	}

	return false, "", nil
}

// removeExistingVersions removes ALL existing version directories for the same registry/ruleset
// before installing a new version, ensuring only one version exists at any time.
func (h *HierarchicalInstaller) removeExistingVersions(ctx context.Context, dir, registry, ruleset string) error {
	rulesetBaseDir := filepath.Join(dir, "arm", registry, ruleset)
	versionEntries, err := os.ReadDir(rulesetBaseDir)
	if err != nil {
		// Directory doesn't exist yet, nothing to clean up
		return nil
	}

	for _, versionEntry := range versionEntries {
		if !versionEntry.IsDir() {
			continue
		}

		// Remove ALL existing version directories
		oldVersionDir := filepath.Join(rulesetBaseDir, versionEntry.Name())
		if err := os.RemoveAll(oldVersionDir); err != nil {
			return err
		}
		slog.InfoContext(ctx, "Removed old version directory", "path", oldVersionDir)
	}

	return nil
}
