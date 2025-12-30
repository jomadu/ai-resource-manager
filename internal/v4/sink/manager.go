package sink

import (
	"fmt"
	"os"
	"path/filepath"

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
//        - For each rule in rulesetResource.Spec.Rules:
//          - Generate content with m.ruleGenerator.GenerateRule(namespace, resource, ruleID)
//          - Generate filename with m.ruleFilenameGenerator.GenerateRuleFilename(rulesetID, ruleID)
//          - Write to disk at filepath.Join(m.directory, filename)
//          - Add filename to installedFiles
//      - Else (promptset file): skip
//    - Else (regular file):
//      - Write file.Content to filepath.Join(m.directory, file.Path)
//      - Add file.Path to installedFiles
// 2. Update index:
//    - packageKey = registry/name@version
//    - index.Rulesets[packageKey] = {Priority: priority, Files: installedFiles}
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
//        - For each prompt in promptsetResource.Spec.Prompts:
//          - Generate content with m.promptGenerator.GeneratePrompt(namespace, resource, promptID)
//          - Generate filename with m.promptFilenameGenerator.GeneratePromptFilename(promptsetID, promptID)
//          - Write to disk at filepath.Join(m.directory, filename)
//          - Add filename to installedFiles
//      - Else (ruleset file): skip
//    - Else (regular file):
//      - Write file.Content to filepath.Join(m.directory, file.Path)
//      - Add file.Path to installedFiles
// 2. Update index:
//    - packageKey = registry/name@version
//    - index.Promptsets[packageKey] = {Files: installedFiles}
func (m *Manager) InstallPromptset(pkg *core.Package) error {
	return nil // TODO
}

// Uninstall removes a package from the sink
func (m *Manager) Uninstall(metadata core.PackageMetadata) error {
	return nil // TODO
}

// IsInstalled checks if a package is installed
func (m *Manager) IsInstalled(metadata core.PackageMetadata) bool {
	return false // TODO
}

// ListRulesets returns all installed rulesets
func (m *Manager) ListRulesets() (*InstalledRuleset[], error) {
	return nil, nil // TODO
}

// ListPromptsets returns all installed promptsets
func (m *Manager) ListPromptsets() (*InstalledPromptset[], error) {
	return nil, nil // TODO
}

// Clean removes orphaned files
func (m *Manager) Clean() error {
	return nil // TODO
}

// Index management
func (m *Manager) loadIndex() (*Index, error) {
	return nil, nil // TODO
}

func (m *Manager) saveIndex(index *Index) error {
	return nil // TODO
}

// generateRulesetIndexRuleFile creates arm_index.* file for rulesets (priority explanation for AI agents)
func (m *Manager) generateRulesetIndexRuleFile() error {
	return nil // TODO
}


