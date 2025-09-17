package converter

import (
	"os"

	"github.com/jomadu/ai-rules-manager/internal/urf"
	"gopkg.in/yaml.v3"
)

// URFRuleset wraps the official URF types with convenience methods
type URFRuleset struct {
	*urf.Ruleset
}

// ConversionConfig contains parameters for conversion
type ConversionConfig struct {
	RulesetID   string
	RulesetName string
	Version     string
	InputFile   string
}

// Converter interface for converting external formats to URF
type Converter interface {
	Convert(content string, config ConversionConfig) (*URFRuleset, error)
}

// SaveToFile saves the URF ruleset to a YAML file
func (r *URFRuleset) SaveToFile(filename string) error {
	data, err := yaml.Marshal(r.Ruleset)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644)
}
