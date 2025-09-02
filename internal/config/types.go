package config

// RegistryConfig defines a registry configuration.
type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// SinkConfig defines a sink configuration for rule deployment.
type SinkConfig struct {
	Directories []string `json:"directories"`
	Include     []string `json:"include"`
	Exclude     []string `json:"exclude"`
}

// Config represents the .armrc.json configuration file structure.
type Config struct {
	Sinks map[string]SinkConfig `json:"sinks"`
}
