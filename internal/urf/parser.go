package urf

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"gopkg.in/yaml.v3"
)

// YAMLParser implements URF file parsing
type YAMLParser struct {
	validator *validator.Validate
}

// NewParser creates a new URF parser
func NewParser() Parser {
	return &YAMLParser{
		validator: validator.New(),
	}
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
func (p *YAMLParser) Parse(file *types.File) (*Ruleset, error) {
	var ruleset Ruleset
	if err := yaml.Unmarshal(file.Content, &ruleset); err != nil {
		return nil, fmt.Errorf("failed to parse URF file %s: %w", file.Path, err)
	}

	// Validate structure using validator
	if err := p.validator.Struct(&ruleset); err != nil {
		return nil, fmt.Errorf("%s: validation failed: %w", file.Path, err)
	}

	// Additional custom validation
	if err := p.validateCustomRules(&ruleset, file.Path); err != nil {
		return nil, err
	}

	return &ruleset, nil
}

// validateCustomRules performs additional business logic validation
func (p *YAMLParser) validateCustomRules(ruleset *Ruleset, filePath string) error {
	// No need to validate rule ID uniqueness since map keys are naturally unique
	// Additional custom validation can be added here if needed
	return nil
}
