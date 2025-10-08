package resource

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"gopkg.in/yaml.v3"
)

// YAMLParser implements resource file parsing
type YAMLParser struct {
	validator *validator.Validate
}

// NewParser creates a new resource parser
func NewParser() Parser {
	return &YAMLParser{
		validator: validator.New(),
	}
}

// IsRuleset checks if a file is a ruleset file by attempting to parse it
func (p *YAMLParser) IsRuleset(file *types.File) bool {
	ext := strings.ToLower(filepath.Ext(file.Path))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	// Try to parse as ruleset
	_, err := p.ParseRuleset(file)
	return err == nil
}

// IsPromptset checks if a file is a promptset file by attempting to parse it
func (p *YAMLParser) IsPromptset(file *types.File) bool {
	ext := strings.ToLower(filepath.Ext(file.Path))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	// Try to parse as promptset
	_, err := p.ParsePromptset(file)
	return err == nil
}

// IsRulesetFile checks if a file at the given path is a ruleset file
func (p *YAMLParser) IsRulesetFile(path string) bool {
	// First check file extension for performance
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	file := &types.File{
		Path:    path,
		Content: content,
		Size:    int64(len(content)),
	}

	return p.IsRuleset(file)
}

// IsPromptsetFile checks if a file at the given path is a promptset file
func (p *YAMLParser) IsPromptsetFile(path string) bool {
	// First check file extension for performance
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	file := &types.File{
		Path:    path,
		Content: content,
		Size:    int64(len(content)),
	}

	return p.IsPromptset(file)
}

// ParseRuleset parses and validates a ruleset file
func (p *YAMLParser) ParseRuleset(file *types.File) (*Ruleset, error) {
	var ruleset Ruleset
	if err := yaml.Unmarshal(file.Content, &ruleset); err != nil {
		return nil, fmt.Errorf("failed to parse ruleset file %s: %w", file.Path, err)
	}

	// Validate structure using validator
	if err := p.validator.Struct(&ruleset); err != nil {
		return nil, fmt.Errorf("%s: validation failed: %w", file.Path, err)
	}

	return &ruleset, nil
}

// ParsePromptset parses and validates a promptset file
func (p *YAMLParser) ParsePromptset(file *types.File) (*Promptset, error) {
	var promptset Promptset
	if err := yaml.Unmarshal(file.Content, &promptset); err != nil {
		return nil, fmt.Errorf("failed to parse promptset file %s: %w", file.Path, err)
	}

	// Validate structure using validator
	if err := p.validator.Struct(&promptset); err != nil {
		return nil, fmt.Errorf("%s: validation failed: %w", file.Path, err)
	}

	return &promptset, nil
}

// ParseRulesets parses all ruleset files from the given directories
func (p *YAMLParser) ParseRulesets(dirs []string, recursive bool, include, exclude []string) ([]*Ruleset, error) {
	var rulesets []*Ruleset

	for _, dir := range dirs {
		files, err := p.discoverFiles(dir, recursive, include, exclude, p.IsRulesetFile)
		if err != nil {
			return nil, fmt.Errorf("failed to discover files in %s: %w", dir, err)
		}

		for _, filePath := range files {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
			}

			file := &types.File{
				Path:    filePath,
				Content: content,
				Size:    int64(len(content)),
			}

			ruleset, err := p.ParseRuleset(file)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ruleset %s: %w", filePath, err)
			}

			rulesets = append(rulesets, ruleset)
		}
	}

	return rulesets, nil
}

// ParsePromptsets parses all promptset files from the given directories
func (p *YAMLParser) ParsePromptsets(dirs []string, recursive bool, include, exclude []string) ([]*Promptset, error) {
	var promptsets []*Promptset

	for _, dir := range dirs {
		files, err := p.discoverFiles(dir, recursive, include, exclude, p.IsPromptsetFile)
		if err != nil {
			return nil, fmt.Errorf("failed to discover files in %s: %w", dir, err)
		}

		for _, filePath := range files {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
			}

			file := &types.File{
				Path:    filePath,
				Content: content,
				Size:    int64(len(content)),
			}

			promptset, err := p.ParsePromptset(file)
			if err != nil {
				return nil, fmt.Errorf("failed to parse promptset %s: %w", filePath, err)
			}

			promptsets = append(promptsets, promptset)
		}
	}

	return promptsets, nil
}

// discoverFiles discovers files in a directory using the given file type checker
func (p *YAMLParser) discoverFiles(dir string, recursive bool, include, exclude []string, fileChecker func(string) bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && fileChecker(path) && p.matchesPatterns(path, include, exclude) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				path := filepath.Join(dir, entry.Name())
				if fileChecker(path) && p.matchesPatterns(path, include, exclude) {
					files = append(files, path)
				}
			}
		}
	}

	return files, nil
}

// matchesPatterns checks if a file path matches the include/exclude patterns
func (p *YAMLParser) matchesPatterns(path string, include, exclude []string) bool {
	// First check if explicitly excluded
	for _, pattern := range exclude {
		if matched, _ := filepath.Match(pattern, path); matched {
			return false
		}
	}

	// If no include patterns, exclude everything
	if len(include) == 0 {
		return false
	}

	// Check if explicitly included
	for _, pattern := range include {
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
	}

	// Not explicitly included, so exclude
	return false
}
