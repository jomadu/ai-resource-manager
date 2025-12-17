//go:build integration

package cloudsmith_integration

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

func TestCloudsmithRegistry_Integration(t *testing.T) {
	_ = godotenv.Load(".env")
	config := loadCloudsmithRegistryConfig(t)
	if config == nil {
		t.Skip("Cloudsmith registry integration test skipped - set CLOUDSMITH_URL, CLOUDSMITH_OWNER, and CLOUDSMITH_REPOSITORY environment variables")
	}

	reg := createCloudsmithRegistry(t, *config)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ruleset := getEnvOrDefault("CLOUDSMITH_TEST_RULESET", "test-ruleset")

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
		require.NotEmpty(t, files, "should return at least one file")

		for _, file := range files {
			assert.NotEmpty(t, file.Path, "file path should not be empty")
			assert.NotEmpty(t, file.Content, "file content should not be empty")
		}
	})
}

func loadCloudsmithRegistryConfig(t *testing.T) *registry.CloudsmithRegistryConfig {
	url := os.Getenv("CLOUDSMITH_URL")
	owner := os.Getenv("CLOUDSMITH_OWNER")
	repository := os.Getenv("CLOUDSMITH_REPOSITORY")

	if url == "" || owner == "" || repository == "" {
		return nil
	}

	config := &registry.CloudsmithRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			URL:  url,
			Type: "cloudsmith",
		},
		Owner:      owner,
		Repository: repository,
	}

	return config
}

func createCloudsmithRegistry(t *testing.T, config registry.CloudsmithRegistryConfig) registry.Registry {
	return registry.NewCloudsmithRegistryNoCache(&config)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
