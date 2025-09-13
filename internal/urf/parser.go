package urf

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
	"gopkg.in/yaml.v3"
)

// YAMLParser implements URF file parsing
type YAMLParser struct{}

// NewParser creates a new URF parser
func NewParser() Parser {
	return &YAMLParser{}
}

// IsURF checks if a file is a URF file by attempting to parse it
func (p *YAMLParser) IsURF(file *types.File) bool {
	ext := strings.ToLower(filepath.Ext(file.Path))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	_, err := p.Parse(file)
	return err == nil
}

// Parse parses and validates a URF file
func (p *YAMLParser) Parse(file *types.File) (*URFFile, error) {
	var urf URFFile
	if err := yaml.Unmarshal(file.Content, &urf); err != nil {
		return nil, fmt.Errorf("failed to parse URF file %s: %w", file.Path, err)
	}

	// Validate structure
	if err := p.validate(&urf, file.Path); err != nil {
		return nil, err
	}

	return &urf, nil
}

// validate validates a URF file structure
func (p *YAMLParser) validate(urf *URFFile, filePath string) error {
	if urf.Version == "" {
		return fmt.Errorf("invalid URF format in %s: missing required field 'version'", filePath)
	}

	if urf.Metadata.ID == "" {
		return fmt.Errorf("invalid URF format in %s: missing required field 'metadata.id'", filePath)
	}

	if urf.Metadata.Name == "" {
		return fmt.Errorf("invalid URF format in %s: missing required field 'metadata.name'", filePath)
	}

	if urf.Metadata.Version == "" {
		return fmt.Errorf("invalid URF format in %s: missing required field 'metadata.version'", filePath)
	}

	if len(urf.Rules) == 0 {
		return fmt.Errorf("invalid URF format in %s: missing required field 'rules'", filePath)
	}

	// Validate rule ID uniqueness
	ruleIDs := make(map[string]bool)
	for _, rule := range urf.Rules {
		if rule.ID == "" {
			return fmt.Errorf("invalid URF format in %s: rule missing required field 'id'", filePath)
		}
		if ruleIDs[rule.ID] {
			return fmt.Errorf("invalid URF format in %s: duplicate rule ID: %s", filePath, rule.ID)
		}
		ruleIDs[rule.ID] = true
	}

	return nil
}
