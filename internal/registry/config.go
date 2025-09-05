package registry

// RegistryConfig defines base registry configuration.
type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// GitRegistryConfig defines git-specific registry configuration.
type GitRegistryConfig struct {
	RegistryConfig
	Branches []string `json:"branches,omitempty"`
}
