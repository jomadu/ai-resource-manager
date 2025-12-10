package manifest

import "github.com/jomadu/ai-rules-manager/internal/resource"

// RulesetConfig represents a ruleset entry in the manifest.
type RulesetConfig struct {
	Version  string   `json:"version"`
	Priority *int     `json:"priority,omitempty"`
	Include  []string `json:"include,omitempty"`
	Exclude  []string `json:"exclude,omitempty"`
	Sinks    []string `json:"sinks"`
}

// PromptsetConfig represents a promptset entry in the manifest.
type PromptsetConfig struct {
	Version string   `json:"version"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
	Sinks   []string `json:"sinks"`
}

// GetIncludePatterns returns include patterns with defaults if none specified
func (r *RulesetConfig) GetIncludePatterns() []string {
	if len(r.Include) == 0 {
		return []string{"*.yml", "*.yaml"}
	}
	return r.Include
}

// GetIncludePatterns returns include patterns with defaults if none specified
func (p *PromptsetConfig) GetIncludePatterns() []string {
	if len(p.Include) == 0 {
		return []string{"*.yml", "*.yaml"}
	}
	return p.Include
}

// SinkConfig defines a sink configuration for resource deployment.
type SinkConfig struct {
	Directory     string                 `json:"directory"`
	Layout        string                 `json:"layout,omitempty"`
	CompileTarget resource.CompileTarget `json:"compileTarget"`
}

// GetLayout returns layout with default if none specified
func (s *SinkConfig) GetLayout() string {
	if s.Layout == "" {
		return "hierarchical"
	}
	return s.Layout
}

// Manifest represents the arm.json manifest file structure.
type Manifest struct {
	Registries map[string]map[string]interface{} `json:"registries,omitempty"`
	Packages   PackageConfig                     `json:"packages"`
	Sinks      map[string]SinkConfig             `json:"sinks,omitempty"`
}

// PackageConfig contains both rulesets and promptsets
type PackageConfig struct {
	Rulesets   map[string]map[string]RulesetConfig   `json:"rulesets,omitempty"`
	Promptsets map[string]map[string]PromptsetConfig `json:"promptsets,omitempty"`
}
