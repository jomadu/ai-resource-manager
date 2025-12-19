package filetype

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
	"gopkg.in/yaml.v3"
)

var validate = validator.New()

// IsRulesetFile checks extension and validates file content as RulesetResource
func IsRulesetFile(path string) bool {
	if !hasYAMLExtension(path) {
		return false
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var ruleset resource.RulesetResource
	if err := yaml.Unmarshal(content, &ruleset); err != nil {
		return false
	}

	return validate.Struct(&ruleset) == nil
}

// IsPromptsetFile checks extension and validates file content as PromptsetResource
func IsPromptsetFile(path string) bool {
	if !hasYAMLExtension(path) {
		return false
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var promptset resource.PromptsetResource
	if err := yaml.Unmarshal(content, &promptset); err != nil {
		return false
	}

	return validate.Struct(&promptset) == nil
}

func hasYAMLExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yml" || ext == ".yaml"
}
