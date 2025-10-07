package index

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jomadu/ai-rules-manager/internal/resource"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"gopkg.in/yaml.v3"
)

type IndexManager struct {
	sinkDir   string
	compiler  resource.Compiler
	generator IndexGenerator
	layout    string
}

type IndexData struct {
	Rulesets map[string]map[string]RulesetInfo `json:"rulesets"`
	Files    map[string]FileInfo               `json:"files"`
}

type RulesetInfo struct {
	Version   string   `json:"version"`
	Priority  int      `json:"priority"`
	FilePaths []string `json:"filePaths"`
}

type FileInfo struct {
	Registry string `json:"registry"`
	Ruleset  string `json:"ruleset"`
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

func (m *IndexManager) Create(ctx context.Context, registry, ruleset, version string, priority int, files []string) error {
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

	for _, file := range files {
		data.Files[file] = FileInfo{Registry: registry, Ruleset: ruleset}
	}

	return m.sync(data)
}

func (m *IndexManager) Read(ctx context.Context, registry, ruleset string) (*RulesetInfo, error) {
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

func (m *IndexManager) Delete(ctx context.Context, registry, ruleset string) error {
	data, err := m.loadJSON()
	if err != nil {
		return err
	}
	if rulesets, ok := data.Rulesets[registry]; ok {
		if info, ok := rulesets[ruleset]; ok {
			for _, file := range info.FilePaths {
				delete(data.Files, file)
			}
			delete(rulesets, ruleset)
			if len(rulesets) == 0 {
				delete(data.Rulesets, registry)
			}
			return m.sync(data)
		}
	}
	return fmt.Errorf("ruleset %s/%s not found", registry, ruleset)
}

func (m *IndexManager) List(ctx context.Context) (map[string]map[string]RulesetInfo, error) {
	data, err := m.loadJSON()
	if err != nil {
		return nil, err
	}
	return data.Rulesets, nil
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
		return &IndexData{Rulesets: make(map[string]map[string]RulesetInfo), Files: make(map[string]FileInfo)}, nil
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
	ruleset := m.generator.CreateRuleset(data)

	// Convert ruleset to resource file format
	resourceContent, err := yaml.Marshal(ruleset)
	if err != nil {
		return fmt.Errorf("failed to marshal ruleset to YAML: %w", err)
	}

	resourceFile := &types.File{
		Path:    "arm-rulesets.yml",
		Content: resourceContent,
		Size:    int64(len(resourceContent)),
	}

	files, err := m.compiler.CompileRuleset("arm", resourceFile)
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
