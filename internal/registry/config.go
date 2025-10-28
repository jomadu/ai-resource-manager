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

// GetBranches returns branches with defaults if none specified
func (g *GitRegistryConfig) GetBranches() []string {
	if len(g.Branches) == 0 {
		return []string{"main", "master"}
	}
	return g.Branches
}

// GitLabRegistryConfig defines gitlab-specific registry configuration.
type GitLabRegistryConfig struct {
	RegistryConfig
	ProjectID  string `json:"projectId,omitempty"`
	GroupID    string `json:"groupId,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

// GetAPIVersion returns API version with default if none specified
func (g *GitLabRegistryConfig) GetAPIVersion() string {
	if g.APIVersion == "" {
		return "v4"
	}
	return g.APIVersion
}

// CloudsmithRegistryConfig defines cloudsmith-specific registry configuration.
type CloudsmithRegistryConfig struct {
	RegistryConfig
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
}

// GetBaseURL returns the API base URL with default if none specified
func (c *CloudsmithRegistryConfig) GetBaseURL() string {
	if c.URL == "" {
		return "https://api.cloudsmith.io"
	}
	return c.URL
}
