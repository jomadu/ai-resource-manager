package index

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jomadu/ai-rules-manager/internal/resource"
)

type IndexManager struct {
	sinkDir   string
	compiler  resource.Compiler
	generator IndexGenerator
	layout    string
}

type IndexData struct {
	Rulesets   map[string]map[string]RulesetInfo   `json:"rulesets"`
	Promptsets map[string]map[string]PromptsetInfo `json:"promptsets"`
}

type RulesetInfo struct {
	Version   string   `json:"version"`
	Priority  int      `json:"priority"`
	FilePaths []string `json:"file_paths"`
}

type PromptsetInfo struct {
	Version   string   `json:"version"`
	FilePaths []string `json:"file_paths"`
}

func NewIndexManager(sinkDir, layout string, target resource.CompileTarget) *IndexManager {
	compiler, _ := resource.NewCompiler(target)
	return &IndexManager{
		sinkDir:   sinkDir,
		compiler:  compiler,
		generator: &DefaultIndexGenerator{},
		layout:    layout,
	}
}

func (m *IndexManager) CreateRuleset(ctx context.Context, registry, ruleset, version string, priority int, files []string) error {
	data, err := m.loadJSON()
	if err != nil {
		return err
	}

	if data.Rulesets[registry] == nil {
		data.Rulesets[registry] = make(map[string]RulesetInfo)
	}

	data.Rulesets[registry][ruleset] = RulesetInfo{
		Version:   version,
		Priority:  priority,
		FilePaths: files,
	}

	return m.sync(data)
}

func (m *IndexManager) CreatePromptset(ctx context.Context, registry, promptset, version string, files []string) error {
	data, err := m.loadJSON()
	if err != nil {
		return err
	}

	if data.Promptsets[registry] == nil {
		data.Promptsets[registry] = make(map[string]PromptsetInfo)
	}

	data.Promptsets[registry][promptset] = PromptsetInfo{
		Version:   version,
		FilePaths: files,
	}

	return m.sync(data)
}

func (m *IndexManager) ReadRuleset(ctx context.Context, registry, ruleset string) (*RulesetInfo, error) {
	data, err := m.loadJSON()
	if err != nil {
		return nil, err
	}
	if rulesets, ok := data.Rulesets[registry]; ok {
		if info, ok := rulesets[ruleset]; ok {
			return &info, nil
		}
	}
	return nil, fmt.Errorf("ruleset %s/%s not found", registry, ruleset)
}

func (m *IndexManager) ReadPromptset(ctx context.Context, registry, promptset string) (*PromptsetInfo, error) {
	data, err := m.loadJSON()
	if err != nil {
		return nil, err
	}
	if promptsets, ok := data.Promptsets[registry]; ok {
		if info, ok := promptsets[promptset]; ok {
			return &info, nil
		}
	}
	return nil, fmt.Errorf("promptset %s/%s not found", registry, promptset)
}

func (m *IndexManager) DeleteRuleset(ctx context.Context, registry, ruleset string) error {
	data, err := m.loadJSON()
	if err != nil {
		return err
	}
	if rulesets, ok := data.Rulesets[registry]; ok {
		if _, ok := rulesets[ruleset]; ok {
			delete(rulesets, ruleset)
			if len(rulesets) == 0 {
				delete(data.Rulesets, registry)
			}
			return m.sync(data)
		}
	}
	return fmt.Errorf("ruleset %s/%s not found", registry, ruleset)
}

func (m *IndexManager) DeletePromptset(ctx context.Context, registry, promptset string) error {
	data, err := m.loadJSON()
	if err != nil {
		return err
	}
	if promptsets, ok := data.Promptsets[registry]; ok {
		if _, ok := promptsets[promptset]; ok {
			delete(promptsets, promptset)
			if len(promptsets) == 0 {
				delete(data.Promptsets, registry)
			}
			return m.sync(data)
		}
	}
	return fmt.Errorf("promptset %s/%s not found", registry, promptset)
}

func (m *IndexManager) ListRulesets(ctx context.Context) (map[string]map[string]RulesetInfo, error) {
	data, err := m.loadJSON()
	if err != nil {
		return nil, err
	}
	return data.Rulesets, nil
}

func (m *IndexManager) ListPromptsets(ctx context.Context) (map[string]map[string]PromptsetInfo, error) {
	data, err := m.loadJSON()
	if err != nil {
		return nil, err
	}
	return data.Promptsets, nil
}

// GetAllInstalledFiles returns all files that should exist in the sink directory
// This includes both resource files (from rulesets/promptsets) and index files
func (m *IndexManager) GetAllInstalledFiles(ctx context.Context) ([]string, error) {
	var allFiles []string

	// Get all ruleset files
	rulesets, err := m.ListRulesets(ctx)
	if err != nil {
		return nil, err
	}

	for _, registryRulesets := range rulesets {
		for _, info := range registryRulesets {
			for _, filePath := range info.FilePaths {
				fullPath := filepath.Join(m.sinkDir, filePath)
				allFiles = append(allFiles, fullPath)
			}
		}
	}

	// Get all promptset files
	promptsets, err := m.ListPromptsets(ctx)
	if err != nil {
		return nil, err
	}

	for _, registryPromptsets := range promptsets {
		for _, info := range registryPromptsets {
			for _, filePath := range info.FilePaths {
				fullPath := filepath.Join(m.sinkDir, filePath)
				allFiles = append(allFiles, fullPath)
			}
		}
	}

	// Add index files
	indexFiles := m.getIndexFilePaths()
	allFiles = append(allFiles, indexFiles...)

	return allFiles, nil
}

// getIndexFilePaths returns the paths of all index files managed by this IndexManager
func (m *IndexManager) getIndexFilePaths() []string {
	var paths []string

	// JSON index file
	if m.layout == "flat" {
		paths = append(paths, filepath.Join(m.sinkDir, "arm-index.json"))
	} else {
		paths = append(paths, filepath.Join(m.sinkDir, "arm", "arm-index.json"))
	}

	// Get compiled index file paths by using the generator and compiler
	// Create a dummy IndexData to generate a ruleset and get the actual compiled file paths
	dummyData := &IndexData{
		Rulesets:   make(map[string]map[string]RulesetInfo),
		Promptsets: make(map[string]map[string]PromptsetInfo),
	}

	// Generate a dummy ruleset using the generator
	dummyRuleset := m.generator.CreateRuleset(dummyData)

	// Compile the dummy ruleset to get the actual file paths
	compiledFiles, err := m.compiler.CompileRuleset("arm", dummyRuleset)
	if err == nil {
		for _, compiledFile := range compiledFiles {
			if m.layout == "flat" {
				paths = append(paths, filepath.Join(m.sinkDir, compiledFile.Path))
			} else {
				paths = append(paths, filepath.Join(m.sinkDir, "arm", compiledFile.Path))
			}
		}
	}

	return paths
}

func (m *IndexManager) sync(data *IndexData) error {
	if err := m.writeJSON(data); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	if err := m.writeCompiled(data); err != nil {
		return fmt.Errorf("failed to write compiled format: %w", err)
	}
	return nil
}

func (m *IndexManager) loadJSON() (*IndexData, error) {
	var jsonPath string
	if m.layout == "flat" {
		jsonPath = filepath.Join(m.sinkDir, "arm-index.json")
	} else {
		jsonPath = filepath.Join(m.sinkDir, "arm", "arm-index.json")
	}

	fileData, err := os.ReadFile(jsonPath)
	if os.IsNotExist(err) {
		return &IndexData{
			Rulesets:   make(map[string]map[string]RulesetInfo),
			Promptsets: make(map[string]map[string]PromptsetInfo),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	data := &IndexData{}
	err = json.Unmarshal(fileData, data)
	return data, err
}

func (m *IndexManager) writeJSON(data *IndexData) error {
	var jsonPath string
	if m.layout == "flat" {
		jsonPath = filepath.Join(m.sinkDir, "arm-index.json")
	} else {
		jsonPath = filepath.Join(m.sinkDir, "arm", "arm-index.json")
	}

	if err := os.MkdirAll(filepath.Dir(jsonPath), 0o755); err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(jsonPath, jsonData, 0o644)
}

func (m *IndexManager) writeCompiled(data *IndexData) error {
	// Only generate arm_index.* files for rulesets (not promptsets)
	// This is because promptsets don't have priority conflicts that need resolution
	if len(data.Rulesets) == 0 {
		return nil // No rulesets to generate index for
	}

	ruleset := m.generator.CreateRuleset(data)

	files, err := m.compiler.CompileRuleset("arm", ruleset)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		var compiledPath string
		if m.layout == "flat" {
			compiledPath = filepath.Join(m.sinkDir, files[0].Path)
		} else {
			compiledPath = filepath.Join(m.sinkDir, "arm", files[0].Path)
		}
		return os.WriteFile(compiledPath, files[0].Content, 0o644)
	}
	return nil
}
