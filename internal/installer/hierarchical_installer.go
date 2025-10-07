package installer

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/index"
	"github.com/jomadu/ai-rules-manager/internal/resource"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// HierarchicalInstaller implements hierarchical file installation preserving directory structure.
type HierarchicalInstaller struct {
	baseDir      string
	indexManager *index.IndexManager
	compiler     resource.Compiler
}

// NewHierarchicalInstaller creates a new hierarchical installer.
func NewHierarchicalInstaller(baseDir string, target resource.CompileTarget) *HierarchicalInstaller {
	compiler, _ := resource.NewCompiler(target)
	return &HierarchicalInstaller{
		baseDir:      baseDir,
		indexManager: index.NewIndexManager(baseDir, "hierarchical", target),
		compiler:     compiler,
	}
}

func (h *HierarchicalInstaller) Install(ctx context.Context, registry, ruleset, version string, priority int, files []types.File) error {
	// Remove existing versions
	if err := h.Uninstall(ctx, registry, ruleset); err != nil {
		return err
	}

	rulesetDir := filepath.Join(h.baseDir, "arm", registry, ruleset, version)
	if err := os.MkdirAll(rulesetDir, 0o755); err != nil {
		return err
	}

	// Process resource files
	processedFiles, err := processResourceFiles(files, registry, ruleset, version, h.compiler)
	if err != nil {
		return err
	}

	var filePaths []string
	for _, file := range processedFiles {
		filePath := filepath.Join(rulesetDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}
		relativePath := filepath.Join("arm", registry, ruleset, version, file.Path)
		filePaths = append(filePaths, relativePath)
	}

	// Update index manager
	if err := h.indexManager.CreateRuleset(ctx, registry, ruleset, version, priority, filePaths); err != nil {
		return err
	}

	return nil
}

func (h *HierarchicalInstaller) Uninstall(ctx context.Context, registry, ruleset string) error {
	// Remove from index first
	if err := h.indexManager.DeleteRuleset(ctx, registry, ruleset); err != nil {
		// Continue even if not in index
		_ = err // Explicitly ignore error
	}

	rulesetDir := filepath.Join(h.baseDir, "arm", registry, ruleset)
	if err := os.RemoveAll(rulesetDir); err != nil {
		return err
	}

	// Clean up empty parent directories
	registryDir := filepath.Join(h.baseDir, "arm", registry)
	if isEmpty(registryDir) {
		if err := os.Remove(registryDir); err != nil {
			return err
		}

		armDir := filepath.Join(h.baseDir, "arm")
		if isEmpty(armDir) {
			if err := os.Remove(armDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *HierarchicalInstaller) ListInstalled(ctx context.Context) ([]Ruleset, error) {
	rulesets, err := h.indexManager.ListRulesets(ctx)
	if err != nil {
		return nil, err
	}

	var installations []Ruleset
	for registry, rulesetMap := range rulesets {
		for name, info := range rulesetMap {
			var filePaths []string
			for _, filePath := range info.FilePaths {
				fullPath := filepath.Join(h.baseDir, strings.TrimPrefix(filePath, "./"))
				filePaths = append(filePaths, fullPath)
			}
			versionPath := filepath.Join(h.baseDir, "arm", registry, name, info.Version)
			installations = append(installations, Ruleset{
				Registry:  registry,
				Ruleset:   name,
				Version:   info.Version,
				Priority:  info.Priority,
				Path:      versionPath,
				FilePaths: filePaths,
			})
		}
	}

	return installations, nil
}

func (h *HierarchicalInstaller) IsInstalled(ctx context.Context, registry, ruleset string) (installed bool, version string, err error) {
	info, err := h.indexManager.ReadRuleset(ctx, registry, ruleset)
	if err != nil {
		return false, "", nil
	}
	return true, info.Version, nil
}
