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
// 
// Algorithm:
// 1. For each file in package:
//    - If filetype.IsResourceFile(file.Path):
//      - If filetype.IsRulesetFile(file.Path):
//        - Parse with parser.ParseRuleset(&file)
//        - namespace = fmt.Sprintf("%s/%s/%s", registry, name, version)
//        - For each rule in rulesetResource.Spec.Rules:
//          - Generate content with m.ruleGenerator.GenerateRule(namespace, resource, ruleID)
//          - Generate filename with m.ruleFilenameGenerator.GenerateRuleFilename(rulesetID, ruleID)
//          - Combine: filepath.Join(filepath.Dir(file.Path), filename)
//          - Write to disk at m.getFilePath(registry, name, version, combinedPath)
//          - Add combinedPath to installedFiles
//      - Else (promptset file): skip
//    - Else (regular file):
//      - Write file.Content to m.getFilePath(registry, name, version, file.Path)
//      - Add file.Path to installedFiles
// 2. Update index:
//    - packageID = core.PackageID(registry, name, version)
//    - index.Rulesets[packageID] = {Priority: priority, Files: installedFiles}
// 3. Generate ruleset index rule file for AI agents
func (m *Manager) InstallRuleset(pkg *core.Package, priority int) error {
	return nil // TODO
}

// InstallPromptset installs a promptset package
//
// Algorithm:
// 1. For each file in package:
//    - If filetype.IsResourceFile(file.Path):
//      - If filetype.IsPromptsetFile(file.Path):
//        - Parse with parser.ParsePromptset(&file)
//        - namespace = fmt.Sprintf("%s/%s/%s", registry, name, version)
//        - For each prompt in promptsetResource.Spec.Prompts:
//          - Generate content with m.promptGenerator.GeneratePrompt(namespace, resource, promptID)
//          - Generate filename with m.promptFilenameGenerator.GeneratePromptFilename(promptsetID, promptID)
//          - Combine: filepath.Join(filepath.Dir(file.Path), filename)
//          - Write to disk at m.getFilePath(registry, name, version, combinedPath)
//          - Add combinedPath to installedFiles
//      - Else (ruleset file): skip
//    - Else (regular file):
//      - Write file.Content to m.getFilePath(registry, name, version, file.Path)
//      - Add file.Path to installedFiles
// 2. Update index:
//    - packageID = core.PackageID(registry, name, version)
//    - index.Promptsets[packageID] = {Files: installedFiles}
func (m *Manager) InstallPromptset(pkg *core.Package) error {
	return nil // TODO
}

// Uninstall removes a package from the sink
func (m *Manager) Uninstall(metadata core.PackageMetadata) error {
	index, err := m.loadIndex()
	if err != nil {
		return err
	}

	packageID := core.PackageID(metadata.RegistryName, metadata.Name, metadata.Version.Version)

	// Remove files for rulesets
	if entry, exists := index.Rulesets[packageID]; exists {
		for _, filePath := range entry.Files {
			fullPath := filepath.Join(m.directory, filePath)
			os.Remove(fullPath) // Ignore errors
		}
		delete(index.Rulesets, packageID)
	}

	// Remove files for promptsets
	if entry, exists := index.Promptsets[packageID]; exists {
		for _, filePath := range entry.Files {
			fullPath := filepath.Join(m.directory, filePath)
			os.Remove(fullPath) // Ignore errors
		}
		delete(index.Promptsets, packageID)
	}

	return m.saveIndex(index)
}

// IsInstalled checks if a package is installed
func (m *Manager) IsInstalled(metadata core.PackageMetadata) bool {
	index, err := m.loadIndex()
	if err != nil {
		return false
	}

	packageID := core.PackageID(metadata.RegistryName, metadata.Name, metadata.Version.Version)
	_, rulesetExists := index.Rulesets[packageID]
	_, promptsetExists := index.Promptsets[packageID]
	return rulesetExists || promptsetExists
}

// ListRulesets returns all installed rulesets
func (m *Manager) ListRulesets() ([]*InstalledRuleset, error) {
	index, err := m.loadIndex()
	if err != nil {
		return nil, err
	}

	var rulesets []*InstalledRuleset
	for packageID, entry := range index.Rulesets {
		registry, name, version, err := core.ParsePackageID(packageID)
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
	for packageID, entry := range index.Promptsets {
		registry, name, version, err := core.ParsePackageID(packageID)
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
	return nil // TODO
}

// Helper functions for path computation

// getFilePath returns the appropriate path for a file based on layout
func (m *Manager) getFilePath(registry, name, version, relativePath string) string {
	if m.layout == LayoutFlat {
		hash := m.hashFile(registry, name, version, relativePath)
		fileName := "arm_" + hash + "_" + strings.ReplaceAll(strings.ReplaceAll(relativePath, "/", "_"), "\\", "_")
		return filepath.Join(m.directory, fileName)
	}
	return filepath.Join(m.armDir, registry, name, version, relativePath)
}

// hashFile creates SHA256 hash for flat layout file naming
func (m *Manager) hashFile(registry, name, version, filePath string) string {
	identifier := fmt.Sprintf("%s/%s@%s:%s", registry, name, version, filePath)
	hash := sha256.Sum256([]byte(identifier))
	return hex.EncodeToString(hash[:])[:8]
}


