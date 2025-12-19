package parser

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
	"gopkg.in/yaml.v3"
)

var validate = validator.New()

// ParseRuleset parses a file into a ruleset resource
func ParseRuleset(file *core.File) (*resource.RulesetResource, error) {
	var ruleset resource.RulesetResource
	if err := yaml.Unmarshal(file.Content, &ruleset); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %s: %w", file.Path, err)
	}
	if err := validate.Struct(&ruleset); err != nil {
		return nil, fmt.Errorf("invalid ruleset in %s: %w", file.Path, err)
	}
	return &ruleset, nil
}

// ParsePromptset parses a file into a promptset resource
func ParsePromptset(file *core.File) (*resource.PromptsetResource, error) {
	var promptset resource.PromptsetResource
	if err := yaml.Unmarshal(file.Content, &promptset); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %s: %w", file.Path, err)
	}
	if err := validate.Struct(&promptset); err != nil {
		return nil, fmt.Errorf("invalid promptset in %s: %w", file.Path, err)
	}
	return &promptset, nil
}


