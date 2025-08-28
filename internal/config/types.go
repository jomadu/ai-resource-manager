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

// Config represents the .armrc.json configuration file structure.
type Config struct {
	Registries map[string]RegistryConfig `json:"registries"`
	Sinks      map[string]SinkConfig     `json:"sinks"`
}
