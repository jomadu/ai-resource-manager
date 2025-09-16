package manifest

import "github.com/jomadu/ai-rules-manager/internal/urf"

// Entry represents a single ruleset entry in the manifest.
type Entry struct {
	Version  string   `json:"version"`
	Priority int      `json:"priority"`
	Include  []string `json:"include"`
	Exclude  []string `json:"exclude"`
	Sinks    []string `json:"sinks"`
}

// SinkConfig defines a sink configuration for rule deployment.
type SinkConfig struct {
	Directory     string            `json:"directory"`
	Layout        string            `json:"layout"`
	CompileTarget urf.CompileTarget `json:"compileTarget"`
}

// Manifest represents the arm.json manifest file structure.
type Manifest struct {
	Registries map[string]map[string]interface{} `json:"registries"`
	Rulesets   map[string]map[string]Entry       `json:"rulesets"`
	Sinks      map[string]SinkConfig             `json:"sinks"`
}
