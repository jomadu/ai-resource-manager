package resource

import (
	"fmt"
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

// IsResource checks if a file is a resource file by attempting to parse it
func (p *YAMLParser) IsResource(file *types.File) bool {
	ext := strings.ToLower(filepath.Ext(file.Path))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	// Try to parse as either ruleset or promptset
	_, errRuleset := p.ParseRuleset(file)
	_, errPromptset := p.ParsePromptset(file)
	return errRuleset == nil || errPromptset == nil
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
