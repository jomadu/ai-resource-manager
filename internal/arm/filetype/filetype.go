package filetype

import (
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/resource"
	"gopkg.in/yaml.v3"
)

var validate = validator.New()

// IsResourceFile checks if file is any ARM resource type
func IsResourceFile(file *core.File) bool {
	return IsRulesetFile(file) || IsPromptsetFile(file)
}

// IsRulesetFile checks extension and validates file content as RulesetResource
func IsRulesetFile(file *core.File) bool {
	if !hasYAMLExtension(file.Path) {
		return false
	}

	var ruleset resource.RulesetResource
	if err := yaml.Unmarshal(file.Content, &ruleset); err != nil {
		return false
	}

	return validate.Struct(&ruleset) == nil
}

// IsPromptsetFile checks extension and validates file content as PromptsetResource
func IsPromptsetFile(file *core.File) bool {
	if !hasYAMLExtension(file.Path) {
		return false
	}

	var promptset resource.PromptsetResource
	if err := yaml.Unmarshal(file.Content, &promptset); err != nil {
		return false
	}

	return validate.Struct(&promptset) == nil
}

func hasYAMLExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yml" || ext == ".yaml"
}
