package config

// RegistryConfig defines a registry configuration.
type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// SinkConfig defines a sink configuration for rule deployment.
type SinkConfig struct {
	Directories []string `json:"directories"`
	Rulesets    []string `json:"rulesets"`
}

// RCConfig represents the .armrc.json configuration file structure.
type RCConfig struct {
	Registries map[string]RegistryConfig `json:"registries"`
	Sinks      map[string]SinkConfig     `json:"sinks"`
}