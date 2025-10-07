package manifest

import "github.com/jomadu/ai-rules-manager/internal/resource"

// Entry represents a single resource entry in the manifest.
type Entry struct {
	Version  string   `json:"version"`
	Priority *int     `json:"priority,omitempty"` // Only for rulesets
	Include  []string `json:"include,omitempty"`
	Exclude  []string `json:"exclude,omitempty"`
	Sinks    []string `json:"sinks"`
}

// GetIncludePatterns returns include patterns with defaults if none specified
func (e *Entry) GetIncludePatterns() []string {
	if len(e.Include) == 0 {
		return []string{"*.yml", "*.yaml"}
	}
	return e.Include
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
	Rulesets   map[string]map[string]Entry `json:"rulesets,omitempty"`
	Promptsets map[string]map[string]Entry `json:"promptsets,omitempty"`
}
