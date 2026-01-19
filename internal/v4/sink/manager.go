package sink

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/filetype"
	"github.com/jomadu/ai-resource-manager/internal/v4/parser"
)

type Layout string

const (
	LayoutFlat         Layout = "flat"
	LayoutHierarchical Layout = "hierarchical"
)

type PackageInstallation struct {
	Metadata core.PackageMetadata
	Sinks    []string
}

type Manager struct {
	directory               string
	layout                  Layout
	armDir                  string
	indexPath               string
	rulesetIndexRulePath    string
	ruleGenerator           compiler.RuleGenerator
	promptGenerator         compiler.PromptGenerator
	ruleFilenameGenerator   compiler.RuleFilenameGenerator
	promptFilenameGenerator compiler.PromptFilenameGenerator
}

type Index struct {
	Version    int                            `json:"version"`
	Rulesets   map[string]RulesetIndexEntry   `json:"rulesets,omitempty"`
	Promptsets map[string]PromptsetIndexEntry `json:"promptsets,omitempty"`
}

type RulesetIndexEntry struct {
	Priority int      `json:"priority"`
	Files    []string `json:"files"`
}

type PromptsetIndexEntry struct {
	Files []string `json:"files"`
}

type InstalledRuleset struct {
	Metadata core.PackageMetadata
	Priority int
	Files    []string
}

type InstalledPromptset struct {
	Metadata core.PackageMetadata
	Files    []string
}

// NewManager creates a new sink manager
func NewManager(directory string, tool compiler.Tool) *Manager {
	// Create directory if it doesn't exist
	os.MkdirAll(directory, 0755)

	// Determine layout: copilot is flat, others are hierarchical
	layout := LayoutHierarchical
	if tool == compiler.Copilot {
		layout = LayoutFlat
	}

	// Set up paths based on layout
	var armDir, indexPath string
	if layout == LayoutFlat {
		armDir = directory
		indexPath = filepath.Join(directory, "arm-index.json")
	} else {
		armDir = filepath.Join(directory, "arm")
		indexPath = filepath.Join(armDir, "arm-index.json")
	}

	// Initialize generators
	ruleGenFactory := compiler.NewRuleGeneratorFactory()
	promptGenFactory := compiler.NewPromptGeneratorFactory()
	ruleFilenameGenFactory := compiler.NewRuleFilenameGeneratorFactory()
	promptFilenameGenFactory := compiler.NewPromptFilenameGeneratorFactory()

	ruleGen, _ := ruleGenFactory.NewRuleGenerator(tool)
	promptGen, _ := promptGenFactory.NewPromptGenerator(tool)
	ruleFilenameGen, _ := ruleFilenameGenFactory.NewRuleFilenameGenerator(tool)
	promptFilenameGen, _ := promptFilenameGenFactory.NewPromptFilenameGenerator(tool)

	// Generate ruleset index rule filename
	rulesetIndexFilename, _ := ruleFilenameGen.GenerateRuleFilename("arm", "index")
	rulesetIndexRulePath := filepath.Join(armDir, rulesetIndexFilename)

	return &Manager{
		directory:               directory,
		layout:                  layout,
		armDir:                  armDir,
		indexPath:               indexPath,
		rulesetIndexRulePath:    rulesetIndexRulePath,
		ruleGenerator:           ruleGen,
		promptGenerator:         promptGen,
		ruleFilenameGenerator:   ruleFilenameGen,
		promptFilenameGenerator: promptFilenameGen,
	}
}

// InstallRuleset installs a ruleset package with priority
func (m *Manager) InstallRuleset(pkg *core.Package, priority int) error {
	// Uninstall old version if exists
	if m.IsInstalled(pkg.Metadata) {
		if err := m.Uninstall(pkg.Metadata); err != nil {
			return err
		}
	}

	var installedFiles []string
	namespace := fmt.Sprintf("%s/%s@%s", pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version)

	for _, file := range pkg.Files {
		if filetype.IsResourceFile(file.Path) {
			if filetype.IsRulesetFile(file.Path) {
				rulesetResource, err := parser.ParseRuleset(file)
				if err != nil {
					return err
				}

				for ruleID := range rulesetResource.Spec.Rules {
					content, err := m.ruleGenerator.GenerateRule(namespace, rulesetResource, ruleID)
					if err != nil {
						return err
					}

					filename, err := m.ruleFilenameGenerator.GenerateRuleFilename(rulesetResource.Metadata.ID, ruleID)
					if err != nil {
						return err
					}

					combinedPath := filepath.Join(filepath.Dir(file.Path), filename)
					fullPath := m.getFilePath(pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version, combinedPath)

					if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
						return err
					}
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}

					relPath, _ := filepath.Rel(m.directory, fullPath)
					installedFiles = append(installedFiles, relPath)
				}
			}
		} else {
			fullPath := m.getFilePath(pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version, file.Path)

			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(fullPath, file.Content, 0644); err != nil {
				return err
			}

			relPath, _ := filepath.Rel(m.directory, fullPath)
			installedFiles = append(installedFiles, relPath)
		}
	}

	index, err := m.loadIndex()
	if err != nil {
		return err
	}

	key := pkgKey(pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version)
	index.Rulesets[key] = RulesetIndexEntry{
		Priority: priority,
		Files:    installedFiles,
	}

	if err := m.saveIndex(index); err != nil {
		return err
	}

	return m.generateRulesetIndexRuleFile()
}

// InstallPromptset installs a promptset package
func (m *Manager) InstallPromptset(pkg *core.Package) error {
	// Uninstall old version if exists
	if m.IsInstalled(pkg.Metadata) {
		if err := m.Uninstall(pkg.Metadata); err != nil {
			return err
		}
	}

	var installedFiles []string
	namespace := fmt.Sprintf("%s/%s@%s", pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version)

	for _, file := range pkg.Files {
		if filetype.IsResourceFile(file.Path) {
			if filetype.IsPromptsetFile(file.Path) {
				promptsetResource, err := parser.ParsePromptset(file)
				if err != nil {
					return err
				}

				for promptID := range promptsetResource.Spec.Prompts {
					content, err := m.promptGenerator.GeneratePrompt(namespace, promptsetResource, promptID)
					if err != nil {
						return err
					}

					filename, err := m.promptFilenameGenerator.GeneratePromptFilename(promptsetResource.Metadata.ID, promptID)
					if err != nil {
						return err
					}

					combinedPath := filepath.Join(filepath.Dir(file.Path), filename)
					fullPath := m.getFilePath(pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version, combinedPath)

					if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
						return err
					}
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}

					relPath, _ := filepath.Rel(m.directory, fullPath)
					installedFiles = append(installedFiles, relPath)
				}
			}
		} else {
			fullPath := m.getFilePath(pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version, file.Path)

			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(fullPath, file.Content, 0644); err != nil {
				return err
			}

			relPath, _ := filepath.Rel(m.directory, fullPath)
			installedFiles = append(installedFiles, relPath)
		}
	}

	index, err := m.loadIndex()
	if err != nil {
		return err
	}

	key := pkgKey(pkg.Metadata.RegistryName, pkg.Metadata.Name, pkg.Metadata.Version.Version)
	index.Promptsets[key] = PromptsetIndexEntry{
		Files: installedFiles,
	}

	return m.saveIndex(index)
}

// Uninstall removes a package from the sink
func (m *Manager) Uninstall(metadata core.PackageMetadata) error {
	index, err := m.loadIndex()
	if err != nil {
		return err
	}

	key := pkgKey(metadata.RegistryName, metadata.Name, metadata.Version.Version)

	// Remove files for rulesets
	if entry, exists := index.Rulesets[key]; exists {
		for _, filePath := range entry.Files {
			fullPath := filepath.Join(m.directory, filePath)
			os.Remove(fullPath) // Ignore errors
		}
		delete(index.Rulesets, key)
	}

	// Remove files for promptsets
	if entry, exists := index.Promptsets[key]; exists {
		for _, filePath := range entry.Files {
			fullPath := filepath.Join(m.directory, filePath)
			os.Remove(fullPath) // Ignore errors
		}
		delete(index.Promptsets, key)
	}

	return m.saveIndex(index)
}

// IsInstalled checks if a package is installed
func (m *Manager) IsInstalled(metadata core.PackageMetadata) bool {
	index, err := m.loadIndex()
	if err != nil {
		return false
	}

	key := pkgKey(metadata.RegistryName, metadata.Name, metadata.Version.Version)
	_, rulesetExists := index.Rulesets[key]
	_, promptsetExists := index.Promptsets[key]
	return rulesetExists || promptsetExists
}

// ListRulesets returns all installed rulesets
func (m *Manager) ListRulesets() ([]*InstalledRuleset, error) {
	index, err := m.loadIndex()
	if err != nil {
		return nil, err
	}

	var rulesets []*InstalledRuleset
	for key, entry := range index.Rulesets {
		registry, name, version, err := parsePkgKey(key)
		if err != nil {
			continue // Skip invalid entries
		}

		rulesets = append(rulesets, &InstalledRuleset{
			Metadata: core.PackageMetadata{
				RegistryName: registry,
				Name:         name,
				Version:      core.Version{Version: version},
			},
			Priority: entry.Priority,
			Files:    entry.Files,
		})
	}

	return rulesets, nil
}

// ListPromptsets returns all installed promptsets
func (m *Manager) ListPromptsets() ([]*InstalledPromptset, error) {
	index, err := m.loadIndex()
	if err != nil {
		return nil, err
	}

	var promptsets []*InstalledPromptset
	for key, entry := range index.Promptsets {
		registry, name, version, err := parsePkgKey(key)
		if err != nil {
			continue // Skip invalid entries
		}

		promptsets = append(promptsets, &InstalledPromptset{
			Metadata: core.PackageMetadata{
				RegistryName: registry,
				Name:         name,
				Version:      core.Version{Version: version},
			},
			Files: entry.Files,
		})
	}

	return promptsets, nil
}

// Clean removes orphaned files
func (m *Manager) Clean() error {
	index, err := m.loadIndex()
	if err != nil {
		return err
	}

	// Collect all tracked files
	trackedFiles := make(map[string]bool)
	for _, entry := range index.Rulesets {
		for _, filePath := range entry.Files {
			trackedFiles[filePath] = true
		}
	}
	for _, entry := range index.Promptsets {
		for _, filePath := range entry.Files {
			trackedFiles[filePath] = true
		}
	}

	// Walk arm directory and remove untracked files
	if _, err := os.Stat(m.armDir); !os.IsNotExist(err) {
		return filepath.Walk(m.armDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue on errors
			}
			if info.IsDir() {
				return nil // Skip directories
			}

			// Skip ARM system files
			filename := filepath.Base(path)
			if filename == "arm-index.json" || strings.HasPrefix(filename, "arm_index.") {
				return nil
			}

			// Get relative path from directory
			relPath, err := filepath.Rel(m.directory, path)
			if err != nil {
				return nil
			}

			// Remove if not tracked
			if !trackedFiles[relPath] {
				os.Remove(path)
			}

			return nil
		})
	}

	return nil
}

// Index management
func (m *Manager) loadIndex() (*Index, error) {
	if _, err := os.Stat(m.indexPath); os.IsNotExist(err) {
		// Return empty index if file doesn't exist
		return &Index{
			Version:    1,
			Rulesets:   make(map[string]RulesetIndexEntry),
			Promptsets: make(map[string]PromptsetIndexEntry),
		}, nil
	}

	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		return nil, err
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	// Initialize maps if nil
	if index.Rulesets == nil {
		index.Rulesets = make(map[string]RulesetIndexEntry)
	}
	if index.Promptsets == nil {
		index.Promptsets = make(map[string]PromptsetIndexEntry)
	}

	return &index, nil
}

func (m *Manager) saveIndex(index *Index) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(m.indexPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.indexPath, data, 0644)
}

// generateRulesetIndexRuleFile creates arm_index.* file for rulesets (priority explanation for AI agents)
func (m *Manager) generateRulesetIndexRuleFile() error {
	index, err := m.loadIndex()
	if err != nil {
		return err
	}

	if len(index.Rulesets) == 0 {
		if _, err := os.Stat(m.rulesetIndexRulePath); err == nil {
			os.Remove(m.rulesetIndexRulePath)
		}
		return nil
	}

	body := "# ARM Rulesets\n\n"
	body += "This file defines the installation priorities for rulesets managed by ARM.\n\n"
	body += "## Priority Rules\n\n"
	body += "**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:\n\n"
	body += "1. **Higher priority numbers take precedence** over lower priority numbers\n"
	body += "2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets\n"
	body += "3. **Always consult this index** to resolve any ambiguity about which rules to follow\n\n"
	body += "## Installed Rulesets\n\n"

	type entry struct {
		key      string
		priority int
		files    []string
	}

	var entries []entry
	for key, info := range index.Rulesets {
		entries = append(entries, entry{
			key:      key,
			priority: info.Priority,
			files:    info.Files,
		})
	}

	// Sort by priority (high to low)
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].priority > entries[i].priority {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	for _, e := range entries {
		body += fmt.Sprintf("### %s\n", e.key)
		body += fmt.Sprintf("- **Priority:** %d\n", e.priority)
		body += "- **Rules:**\n"
		for _, file := range e.files {
			body += fmt.Sprintf("  - %s\n", file)
		}
		body += "\n"
	}

	if err := os.MkdirAll(filepath.Dir(m.rulesetIndexRulePath), 0755); err != nil {
		return err
	}
	return os.WriteFile(m.rulesetIndexRulePath, []byte(body), 0644)
}

// Helper functions for path computation

// getFilePath returns the appropriate path for a file based on layout
func (m *Manager) getFilePath(registry, name, version, relativePath string) string {
	if m.layout == LayoutFlat {
		packageHash := m.hashPackage(registry, name, version)
		pathHash := m.hashPath(relativePath)
		pathPart := strings.ReplaceAll(strings.ReplaceAll(relativePath, "/", "_"), "\\", "_")

		// Calculate filename overhead: arm_xxxx_xxxx_
		filenameOverhead := len("arm_xxxx_xxxx_")
		maxPathLen := 100 - filenameOverhead

		if len(pathPart) > maxPathLen {
			filename := filepath.Base(relativePath)
			ext := filepath.Ext(filename)
			nameWithoutExt := strings.TrimSuffix(filename, ext)

			// Try progressively shorter options
			if len(filename) <= maxPathLen {
				// 2. Use just filename with extension
				pathPart = filename
			} else {
				// 3. Use truncated filename with extension
				availableForName := maxPathLen - len(ext)
				truncatedName := nameWithoutExt[:min(availableForName, len(nameWithoutExt))]
				pathPart = truncatedName + ext
			}
		}

		fileName := "arm_" + packageHash + "_" + pathHash + "_" + pathPart
		return filepath.Join(m.directory, fileName)
	}
	return filepath.Join(m.armDir, registry, name, version, relativePath)
}

// hashPackage creates 4-char hash for package identification
func (m *Manager) hashPackage(registry, name, version string) string {
	identifier := fmt.Sprintf("%s/%s@%s", registry, name, version)
	hash := sha256.Sum256([]byte(identifier))
	return hex.EncodeToString(hash[:])[:4]
}

// hashPath creates 4-char hash for path identification
func (m *Manager) hashPath(filePath string) string {
	hash := sha256.Sum256([]byte(filePath))
	return hex.EncodeToString(hash[:])[:4]
}
