package installer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/index"
	"github.com/jomadu/ai-rules-manager/internal/resource"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// FlatInstaller implements flat file installation where all files are placed directly in the target directory
// with hashed names and a lookup index to preserve namespace information.
type FlatInstaller struct {
	baseDir      string
	indexManager *index.IndexManager
	compiler     resource.Compiler
}

// NewFlatInstaller creates a new flat installer.
func NewFlatInstaller(baseDir string, target resource.CompileTarget) *FlatInstaller {
	compiler, _ := resource.NewCompiler(target)
	return &FlatInstaller{
		baseDir:      baseDir,
		indexManager: index.NewIndexManager(baseDir, "flat", target),
		compiler:     compiler,
	}
}

func (f *FlatInstaller) InstallRuleset(ctx context.Context, registry, ruleset, version string, priority int, files []types.File) error {
	if err := os.MkdirAll(f.baseDir, 0o755); err != nil {
		return err
	}

	// Remove existing versions
	if err := f.UninstallRuleset(ctx, registry, ruleset); err != nil {
		return err
	}

	// Process resource files
	processedFiles, err := compileResourceFiles(files, registry, ruleset, version, f.compiler)
	if err != nil {
		return err
	}

	var filePaths []string
	for _, file := range processedFiles {
		hash := f.hashFile(registry, ruleset, version, file.Path)
		fileName := "arm_" + hash + "_" + strings.ReplaceAll(file.Path, "/", "_")
		filePath := filepath.Join(f.baseDir, fileName)

		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}

		filePaths = append(filePaths, fileName)
	}

	// Update index manager
	if err := f.indexManager.CreateRuleset(ctx, registry, ruleset, version, priority, filePaths); err != nil {
		return err
	}

	return nil
}

func (f *FlatInstaller) UninstallRuleset(ctx context.Context, registry, ruleset string) error {
	info, err := f.indexManager.ReadRuleset(ctx, registry, ruleset)
	if err != nil {
		return nil // Not installed
	}

	// Remove files
	var removedCount int
	for _, filePath := range info.FilePaths {
		fullPath := filepath.Join(f.baseDir, filePath)
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return err
		}
		removedCount++
	}

	// Remove from index
	if err := f.indexManager.DeleteRuleset(ctx, registry, ruleset); err != nil {
		return err
	}

	return nil
}

func (f *FlatInstaller) ListInstalledRulesets(ctx context.Context) ([]Ruleset, error) {
	rulesets, err := f.indexManager.ListRulesets(ctx)
	if err != nil {
		return nil, err
	}

	var installations []Ruleset
	for registry, rulesetMap := range rulesets {
		for name, info := range rulesetMap {
			var filePaths []string
			for _, filePath := range info.FilePaths {
				fullPath := filepath.Join(f.baseDir, filePath)
				filePaths = append(filePaths, fullPath)
			}
			installations = append(installations, Ruleset{
				Registry:  registry,
				Ruleset:   name,
				Version:   info.Version,
				Priority:  info.Priority,
				Path:      f.baseDir,
				FilePaths: filePaths,
			})
		}
	}

	return installations, nil
}

func (f *FlatInstaller) IsRulesetInstalled(ctx context.Context, registry, ruleset string) (installed bool, version string, err error) {
	info, err := f.indexManager.ReadRuleset(ctx, registry, ruleset)
	if err != nil {
		return false, "", nil
	}
	return true, info.Version, nil
}

func (f *FlatInstaller) InstallPromptset(ctx context.Context, registry, promptset, version string, files []types.File) error {
	if err := os.MkdirAll(f.baseDir, 0o755); err != nil {
		return err
	}

	// Remove existing versions
	if err := f.UninstallPromptset(ctx, registry, promptset); err != nil {
		return err
	}

	// Process resource files
	processedFiles, err := compileResourceFiles(files, registry, promptset, version, f.compiler)
	if err != nil {
		return err
	}

	var filePaths []string
	for _, file := range processedFiles {
		hash := f.hashFile(registry, promptset, version, file.Path)
		fileName := "arm_" + hash + "_" + strings.ReplaceAll(file.Path, "/", "_")
		filePath := filepath.Join(f.baseDir, fileName)

		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}

		filePaths = append(filePaths, fileName)
	}

	// Update index manager (promptsets don't impact arm_index.* generation)
	if err := f.indexManager.CreatePromptset(ctx, registry, promptset, version, filePaths); err != nil {
		return err
	}

	return nil
}

func (f *FlatInstaller) UninstallPromptset(ctx context.Context, registry, promptset string) error {
	info, err := f.indexManager.ReadPromptset(ctx, registry, promptset)
	if err != nil {
		return nil // Not installed
	}

	// Remove files
	var removedCount int
	for _, filePath := range info.FilePaths {
		fullPath := filepath.Join(f.baseDir, filePath)
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return err
		}
		removedCount++
	}

	// Remove from index
	if err := f.indexManager.DeletePromptset(ctx, registry, promptset); err != nil {
		return err
	}

	return nil
}

func (f *FlatInstaller) ListInstalledPromptsets(ctx context.Context) ([]Promptset, error) {
	promptsets, err := f.indexManager.ListPromptsets(ctx)
	if err != nil {
		return nil, err
	}

	var installations []Promptset
	for registry, promptsetMap := range promptsets {
		for name, info := range promptsetMap {
			var filePaths []string
			for _, filePath := range info.FilePaths {
				fullPath := filepath.Join(f.baseDir, filePath)
				filePaths = append(filePaths, fullPath)
			}
			installations = append(installations, Promptset{
				Registry:  registry,
				Promptset: name,
				Version:   info.Version,
				Path:      f.baseDir,
				FilePaths: filePaths,
			})
		}
	}

	return installations, nil
}

func (f *FlatInstaller) IsPromptsetInstalled(ctx context.Context, registry, promptset string) (installed bool, version string, err error) {
	info, err := f.indexManager.ReadPromptset(ctx, registry, promptset)
	if err != nil {
		return false, "", nil
	}
	return true, info.Version, nil
}

// hashFile creates a SHA256 hash for the file identifier
func (f *FlatInstaller) hashFile(registry, resourceName, version, filePath string) string {
	identifier := fmt.Sprintf("%s/%s@%s:%s", registry, resourceName, version, filePath)
	hash := sha256.Sum256([]byte(identifier))
	return hex.EncodeToString(hash[:])[:8]
}
