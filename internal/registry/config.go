package registry

import (
	"fmt"
	neturl "net/url"
	"strings"
)

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

// GetBaseURL returns base URL with default if none specified
func (c *CloudsmithRegistryConfig) GetBaseURL() string {
	if c.URL == "" {
		return "https://api.cloudsmith.io"
	}
	return c.URL
}

// ParseCloudsmithURL extracts owner and repository from a Cloudsmith URL
// Expected format: https://app.cloudsmith.com/[owner]/[repository]
// Returns: owner, repository, error
func ParseCloudsmithURL(url string) (owner, repository string, err error) {
	// Parse the URL
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if it's a Cloudsmith URL
	if parsedURL.Host != "app.cloudsmith.com" {
		return "", "", fmt.Errorf("expected Cloudsmith URL format: https://app.cloudsmith.com/[owner]/[repository], got: %s", url)
	}

	// Split the path to extract owner and repository
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) != 2 {
		return "", "", fmt.Errorf("expected URL format: https://app.cloudsmith.com/[owner]/[repository], got: %s", url)
	}

	owner = pathParts[0]
	repository = pathParts[1]

	// Validate non-empty
	if owner == "" {
		return "", "", fmt.Errorf("owner cannot be empty in URL: %s", url)
	}
	if repository == "" {
		return "", "", fmt.Errorf("repository cannot be empty in URL: %s", url)
	}

	// Return owner and repository
	return owner, repository, nil
}
