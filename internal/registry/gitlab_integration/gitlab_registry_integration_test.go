//go:build integration

package gitlab_integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitLabRegistry_Integration(t *testing.T) {
	_ = godotenv.Load(".env")
	config := loadGitLabRegistryConfig(t)
	if config == nil {
		t.Skip("GitLab registry integration test skipped - set GITLAB_REGISTRY_URL and GITLAB_PROJECT_ID or GITLAB_GROUP_ID environment variables")
	}

	reg := createGitLabRegistry(t, *config)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ruleset := getEnvOrDefault("GITLAB_TEST_RULESET", "test-ruleset")

	t.Run("ListVersions", func(t *testing.T) {
		t.Logf("Testing ruleset: %s", ruleset)

		versions, err := reg.ListVersions(ctx, ruleset)
		if err != nil {
			t.Logf("ListVersions error: %v", err)
		}
		t.Logf("Found %d versions: %+v", len(versions), versions)
		require.NoError(t, err)
		require.NotEmpty(t, versions, "should return at least one version")

		for _, version := range versions {
			assert.NotEmpty(t, version.Version, "version should not be empty")
			assert.NotEmpty(t, version.Display, "display should not be empty")
		}
	})

	t.Run("ResolveVersion", func(t *testing.T) {
		resolved, err := reg.ResolveVersion(ctx, ruleset, "latest")
		require.NoError(t, err)
		assert.NotNil(t, resolved)
		assert.NotEmpty(t, resolved.Version.Version)
	})

	t.Run("GetContent", func(t *testing.T) {
		versions, err := reg.ListVersions(ctx, ruleset)
		require.NoError(t, err)
		require.NotEmpty(t, versions)

		t.Logf("Using version: %+v", versions[0])

		selector := types.ContentSelector{
			Include: []string{"**/*.yml", "**/*.yaml", "**/*.md"},
		}

		files, err := reg.GetContent(ctx, ruleset, versions[0], selector)
		if err != nil {
			t.Logf("GetContent error: %v", err)
		}
		require.NoError(t, err)

		for _, file := range files {
			assert.NotEmpty(t, file.Path, "file path should not be empty")
			assert.NotEmpty(t, file.Content, "file content should not be empty")
		}
	})
}

func loadGitLabRegistryConfig(t *testing.T) *registry.GitLabRegistryConfig {
	url := os.Getenv("GITLAB_REGISTRY_URL")
	if url == "" {
		return nil
	}

	config := &registry.GitLabRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			URL:  url,
			Type: "gitlab",
		},
	}

	if projectID := os.Getenv("GITLAB_PROJECT_ID"); projectID != "" {
		config.ProjectID = projectID
	}

	if groupID := os.Getenv("GITLAB_GROUP_ID"); groupID != "" {
		config.GroupID = groupID
	}

	if apiVersion := os.Getenv("GITLAB_API_VERSION"); apiVersion != "" {
		config.APIVersion = apiVersion
	}

	return config
}

func createGitLabRegistry(t *testing.T, config registry.GitLabRegistryConfig) registry.Registry {
	return registry.NewGitLabRegistryNoCache("test-registry", &config)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
