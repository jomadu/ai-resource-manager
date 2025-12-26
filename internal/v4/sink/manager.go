package sink

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

type PackageInstallation struct {
	Metadata core.PackageMetadata
	Sinks    []string
}

type Manager struct {
	directory            string
	compileTarget        string
	armDir               string
	indexPath            string
	rulesetIndexRulePath string
}

type Index struct {
	Packages map[string]map[string]interface{} `json:"packages"`
}

type PackageEntry struct {
	Version      string   `json:"version"`
	ResourceType string   `json:"resourceType"`
	Files        []string `json:"files"`
}

type RulesetEntry struct {
	PackageEntry
	Priority int `json:"priority"`
}

type PromptsetEntry struct {
	PackageEntry
}

// NewManager creates a new sink manager
func NewManager(directory string, compileTarget compiler.CompileTarget) *Manager {
	// Create directory if it doesn't exist
	os.MkdirAll(directory, 0755)

	// Set up paths
	armDir := filepath.Join(directory, "arm")
	indexPath := filepath.Join(armDir, "arm-index.json")

	return &Manager{
		directory:     directory,
		compileTarget: string(compileTarget),
		armDir:        armDir,
		indexPath:     indexPath,
	}
}

// InstallRuleset installs a ruleset package with priority
func (m *Manager) InstallRuleset(pkg *core.Package, priority int) error {
	return nil // TODO
}

// InstallPromptset installs a promptset package
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
func (m *Manager) ListRulesets() (map[string]map[string]*RulesetEntry, error) {
	return nil, nil // TODO
}

// ListPromptsets returns all installed promptsets
func (m *Manager) ListPromptsets() (map[string]map[string]*PromptsetEntry, error) {
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

// generateIndexFile creates arm_index.* file for rulesets (priority explanation for AI agents)
func (m *Manager) generateIndexFile() error {
	return nil // TODO
}

// Index helpers
func (idx *Index) GetRuleset(registry, pkg string) (*RulesetEntry, error) {
	if idx.Packages[registry] == nil {
		return nil, fmt.Errorf("registry %s not found", registry)
	}
	entry, exists := idx.Packages[registry][pkg]
	if !exists {
		return nil, fmt.Errorf("package %s not found in registry %s", pkg, registry)
	}
	ruleset, ok := entry.(*RulesetEntry)
	if !ok {
		return nil, fmt.Errorf("package %s is not a ruleset", pkg)
	}
	return ruleset, nil
}

func (idx *Index) GetPromptset(registry, pkg string) (*PromptsetEntry, error) {
	if idx.Packages[registry] == nil {
		return nil, fmt.Errorf("registry %s not found", registry)
	}
	entry, exists := idx.Packages[registry][pkg]
	if !exists {
		return nil, fmt.Errorf("package %s not found in registry %s", pkg, registry)
	}
	promptset, ok := entry.(*PromptsetEntry)
	if !ok {
		return nil, fmt.Errorf("package %s is not a promptset", pkg)
	}
	return promptset, nil
}
