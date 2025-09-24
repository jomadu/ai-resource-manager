//go:build integration

package git_integration

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitRegistry_Integration(t *testing.T) {
	_ = godotenv.Load(".env")
	config := loadGitRegistryConfig(t)
	if config == nil {
		t.Skip("Git registry integration test skipped - set GIT_REGISTRY_URL environment variable")
	}

	reg := createGitRegistry(t, *config)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ruleset := getEnvOrDefault("GIT_TEST_RULESET", "test-ruleset")

	t.Run("ListVersions", func(t *testing.T) {
		versions, err := reg.ListVersions(ctx, ruleset)
		require.NoError(t, err)
		assert.NotEmpty(t, versions, "should return at least one version")

		for _, version := range versions {
			assert.NotEmpty(t, version.Version, "version should not be empty")
			assert.NotEmpty(t, version.Display, "display should not be empty")
		}
	})

	t.Run("ResolveVersion", func(t *testing.T) {
		version := getEnvOrDefault("GIT_TEST_VERSION", "latest")
		resolved, err := reg.ResolveVersion(ctx, ruleset, version)
		require.NoError(t, err)
		assert.NotNil(t, resolved)
		assert.NotEmpty(t, resolved.Version.Version)
	})

	t.Run("GetContent", func(t *testing.T) {
		versions, err := reg.ListVersions(ctx, ruleset)
		require.NoError(t, err)
		require.NotEmpty(t, versions)

		selector := loadContentSelector()

		files, err := reg.GetContent(ctx, ruleset, versions[0], selector)
		require.NoError(t, err)
		assert.NotEmpty(t, files, "should return at least one file")

		for _, file := range files {
			assert.NotEmpty(t, file.Path, "file path should not be empty")
			assert.NotEmpty(t, file.Content, "file content should not be empty")
		}
	})
}

func loadGitRegistryConfig(t *testing.T) *registry.GitRegistryConfig {
	url := os.Getenv("GIT_REGISTRY_URL")
	if url == "" {
		return nil
	}

	config := &registry.GitRegistryConfig{
		RegistryConfig: registry.RegistryConfig{
			URL:  url,
			Type: "git",
		},
	}

	if branches := os.Getenv("GIT_REGISTRY_BRANCHES"); branches != "" {
		config.Branches = []string{branches}
	}

	return config
}

func createGitRegistry(t *testing.T, config registry.GitRegistryConfig) registry.Registry {
	repoCache, err := cache.NewGitRepoCache(config, "test-repo", config.URL)
	require.NoError(t, err)

	return registry.NewGitRegistryNoCache(config, repoCache)
}

func loadContentSelector() types.ContentSelector {
	var includes []string
	var excludes []string

	if includeStr := os.Getenv("GIT_TEST_INCLUDES"); includeStr != "" {
		includes = strings.Split(includeStr, ",")
	} else {
		includes = []string{"**/*.yml", "**/*.yaml", "**/*.md"}
	}

	if excludeStr := os.Getenv("GIT_TEST_EXCLUDES"); excludeStr != "" {
		excludes = strings.Split(excludeStr, ",")
	}

	return types.ContentSelector{
		Include: includes,
		Exclude: excludes,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
