package installer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// FlatInstaller implements flat file installation where all files are placed directly in the target directory
// with hashed names and a lookup index to preserve namespace information.
type FlatInstaller struct{}

// IndexEntry represents a file entry in the flat installer index
type IndexEntry struct {
	Registry string `json:"registry"`
	Ruleset  string `json:"ruleset"`
	Version  string `json:"version"`
	FilePath string `json:"filePath"`
}

// Index represents the flat installer lookup index
type Index map[string]IndexEntry

// NewFlatInstaller creates a new flat installer.
func NewFlatInstaller() *FlatInstaller {
	return &FlatInstaller{}
}

func (f *FlatInstaller) Install(ctx context.Context, dir, registry, ruleset, version string, files []types.File) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// Load existing index
	index, err := f.loadIndex(dir)
	if err != nil {
		return err
	}

	// Remove existing versions before installing new one
	if err := f.removeExistingVersions(ctx, dir, registry, ruleset, index); err != nil {
		return err
	}

	for _, file := range files {
		hash := f.hashFile(registry, ruleset, version, file.Path)
		fileName := "arm_" + hash + "_" + strings.ReplaceAll(file.Path, "/", "_")
		filePath := filepath.Join(dir, fileName)

		if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
			return err
		}

		// Update index
		index[fileName] = IndexEntry{
			Registry: registry,
			Ruleset:  ruleset,
			Version:  version,
			FilePath: file.Path,
		}
	}

	// Save updated index
	if err := f.saveIndex(dir, index); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Installed files (flat)", "count", len(files), "path", dir)
	return nil
}

func (f *FlatInstaller) Uninstall(ctx context.Context, dir, registry, ruleset string) error {
	// Load index
	index, err := f.loadIndex(dir)
	if err != nil {
		return err
	}

	// If index is empty, no previous installation exists
	if len(index) == 0 {
		return nil
	}

	// Find and remove files for this registry/ruleset
	var removedCount int
	for fileName, entry := range index {
		if entry.Registry == registry && entry.Ruleset == ruleset {
			filePath := filepath.Join(dir, fileName)
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				return err
			}
			delete(index, fileName)
			removedCount++
		}
	}

	// Save updated index
	if err := f.saveIndex(dir, index); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Uninstalled files (flat)", "count", removedCount, "path", dir)
	return nil
}

func (f *FlatInstaller) ListInstalled(ctx context.Context, dir string) ([]Installation, error) {
	// Load index
	index, err := f.loadIndex(dir)
	if err != nil {
		return nil, err
	}

	var installations []Installation
	for fileName, entry := range index {
		installations = append(installations, Installation{
			Ruleset: entry.Ruleset,
			Version: entry.Version,
			Path:    filepath.Join(dir, fileName),
		})
	}

	return installations, nil
}

func (f *FlatInstaller) IsInstalled(ctx context.Context, dir, registry, ruleset string) (installed bool, version string, err error) {
	// Load index
	index, err := f.loadIndex(dir)
	if err != nil {
		return false, "", err
	}

	// Find any file for this registry/ruleset
	for _, entry := range index {
		if entry.Registry == registry && entry.Ruleset == ruleset {
			return true, entry.Version, nil
		}
	}

	return false, "", nil
}

// hashFile creates a SHA256 hash for the file identifier
func (f *FlatInstaller) hashFile(registry, ruleset, version, filePath string) string {
	identifier := fmt.Sprintf("%s/%s@%s:%s", registry, ruleset, version, filePath)
	hash := sha256.Sum256([]byte(identifier))
	return hex.EncodeToString(hash[:])[:8]
}

// loadIndex loads the index file from the directory
func (f *FlatInstaller) loadIndex(dir string) (Index, error) {
	indexPath := filepath.Join(dir, "arm-index.json")
	index := make(Index)

	data, err := os.ReadFile(indexPath)
	if os.IsNotExist(err) {
		return index, nil // Return empty index if file doesn't exist
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	return index, nil
}

// saveIndex saves the index file to the directory
func (f *FlatInstaller) saveIndex(dir string, index Index) error {
	indexPath := filepath.Join(dir, "arm-index.json")

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0o644)
}

// removeExistingVersions removes ALL existing files for the same registry/ruleset
// before installing a new version, ensuring only one version exists at any time.
func (f *FlatInstaller) removeExistingVersions(ctx context.Context, dir, registry, ruleset string, index Index) error {
	// Find and remove all files for this registry/ruleset
	for fileName, entry := range index {
		if entry.Registry == registry && entry.Ruleset == ruleset {
			filePath := filepath.Join(dir, fileName)
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				return err
			}
			delete(index, fileName)
			slog.InfoContext(ctx, "Removed old version file", "path", filePath)
		}
	}

	return nil
}
