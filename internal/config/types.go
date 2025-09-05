package config

// SinkConfig defines a sink configuration for rule deployment.
type SinkConfig struct {
	Directories []string `json:"directories"`
	Include     []string `json:"include"`
	Exclude     []string `json:"exclude"`
	Layout      string   `json:"layout,omitempty"`
}

// Config represents the .armrc.json configuration file structure.
type Config struct {
	Sinks map[string]SinkConfig `json:"sinks"`
}
