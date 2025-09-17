package registry

// RegistryConfig defines base registry configuration.
type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// GitRegistryConfig defines git-specific registry configuration.
type GitRegistryConfig struct {
	RegistryConfig
	Branches []string `json:"branches"`
}

// GitLabRegistryConfig defines gitlab-specific registry configuration.
type GitLabRegistryConfig struct {
	RegistryConfig
	ProjectID  string `json:"project_id,omitempty"`
	GroupID    string `json:"group_id,omitempty"`
	APIVersion string `json:"api_version"`
}
