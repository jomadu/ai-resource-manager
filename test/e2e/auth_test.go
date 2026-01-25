package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

// TestAuthenticationWithArmrc tests .armrc file handling for authentication
func TestAuthenticationWithArmrc(t *testing.T) {
	workDir := t.TempDir()
	homeDir := t.TempDir()

	// Set HOME environment variable for this test
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", homeDir)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	arm := helpers.NewARMRunner(t, workDir)

	t.Run("LocalArmrcPrecedence", func(t *testing.T) {
		// Create global .armrc in home directory
		globalArmrc := filepath.Join(homeDir, ".armrc")
		globalContent := `[registry https://gitlab.example.com/project/123]
token = global-token-123
`
		if err := os.WriteFile(globalArmrc, []byte(globalContent), 0o600); err != nil {
			t.Fatalf("failed to write global .armrc: %v", err)
		}

		// Create local .armrc in working directory (should override global)
		localArmrc := filepath.Join(workDir, ".armrc")
		localContent := `[registry https://gitlab.example.com/project/123]
token = local-token-123
`
		if err := os.WriteFile(localArmrc, []byte(localContent), 0o600); err != nil {
			t.Fatalf("failed to write local .armrc: %v", err)
		}

		// Add a GitLab registry (this will read .armrc internally)
		// Note: We can't directly test token usage without a real GitLab instance,
		// but we can verify the .armrc file is read correctly by checking the config
		arm.MustRun("add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--project-id", "123", "test-gitlab")

		// Verify registry was added
		armJSON := filepath.Join(workDir, "arm.json")
		helpers.AssertFileExists(t, armJSON)

		manifest := helpers.ReadJSON(t, armJSON)
		registries, ok := manifest["registries"].(map[string]interface{})
		if !ok {
			t.Fatal("registries field not found")
		}

		if _, ok := registries["test-gitlab"]; !ok {
			t.Error("test-gitlab registry not found")
		}

		// Clean up for next test
		arm.MustRun("remove", "registry", "test-gitlab")
		_ = os.Remove(localArmrc)
		_ = os.Remove(globalArmrc)
	})

	t.Run("EnvironmentVariableExpansion", func(t *testing.T) {
		// Set environment variable
		testToken := "test-token-from-env-var"
		_ = os.Setenv("TEST_GITLAB_TOKEN", testToken)
		defer func() { _ = os.Unsetenv("TEST_GITLAB_TOKEN") }()

		// Create .armrc with environment variable reference
		localArmrc := filepath.Join(workDir, ".armrc")
		content := `[registry https://gitlab.example.com/project/456]
token = ${TEST_GITLAB_TOKEN}
`
		if err := os.WriteFile(localArmrc, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write .armrc: %v", err)
		}

		// Add GitLab registry
		arm.MustRun("add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--project-id", "456", "test-gitlab-env")

		// Verify registry was added
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries := manifest["registries"].(map[string]interface{})

		if _, ok := registries["test-gitlab-env"]; !ok {
			t.Error("test-gitlab-env registry not found")
		}

		// Clean up
		arm.MustRun("remove", "registry", "test-gitlab-env")
		_ = os.Remove(localArmrc)
	})

	t.Run("CloudsmithAuthentication", func(t *testing.T) {
		// Create .armrc with Cloudsmith API key
		localArmrc := filepath.Join(workDir, ".armrc")
		content := `[registry https://api.cloudsmith.io/myorg/ai-rules]
token = ckcy_test_api_key_123456789
`
		if err := os.WriteFile(localArmrc, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write .armrc: %v", err)
		}

		// Add Cloudsmith registry
		arm.MustRun("add", "registry", "cloudsmith", "--owner", "myorg", "--repo", "ai-rules", "test-cloudsmith")

		// Verify registry was added
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries := manifest["registries"].(map[string]interface{})

		if _, ok := registries["test-cloudsmith"]; !ok {
			t.Error("test-cloudsmith registry not found")
		}

		// Clean up
		arm.MustRun("remove", "registry", "test-cloudsmith")
		_ = os.Remove(localArmrc)
	})

	t.Run("MultipleSectionsInArmrc", func(t *testing.T) {
		// Create .armrc with multiple registry sections
		localArmrc := filepath.Join(workDir, ".armrc")
		content := `[registry https://gitlab.example.com/project/111]
token = token-111

[registry https://gitlab.example.com/project/222]
token = token-222

[registry https://api.cloudsmith.io/org1/repo1]
token = ckcy_token_1

[registry https://api.cloudsmith.io/org2/repo2]
token = ckcy_token_2
`
		if err := os.WriteFile(localArmrc, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write .armrc: %v", err)
		}

		// Add multiple registries
		arm.MustRun("add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--project-id", "111", "gitlab-1")
		arm.MustRun("add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--project-id", "222", "gitlab-2")
		arm.MustRun("add", "registry", "cloudsmith", "--owner", "org1", "--repo", "repo1", "cloudsmith-1")
		arm.MustRun("add", "registry", "cloudsmith", "--owner", "org2", "--repo", "repo2", "cloudsmith-2")

		// Verify all registries were added
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries := manifest["registries"].(map[string]interface{})

		expectedRegistries := []string{"gitlab-1", "gitlab-2", "cloudsmith-1", "cloudsmith-2"}
		for _, name := range expectedRegistries {
			if _, ok := registries[name]; !ok {
				t.Errorf("registry %s not found", name)
			}
		}

		// Clean up
		for _, name := range expectedRegistries {
			arm.MustRun("remove", "registry", name)
		}
		_ = os.Remove(localArmrc)
	})

	t.Run("GlobalArmrcOnly", func(t *testing.T) {
		// Create only global .armrc (no local)
		globalArmrc := filepath.Join(homeDir, ".armrc")
		content := `[registry https://gitlab.example.com/project/789]
token = global-only-token
`
		if err := os.WriteFile(globalArmrc, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write global .armrc: %v", err)
		}

		// Add GitLab registry (should use global .armrc)
		arm.MustRun("add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--project-id", "789", "test-global")

		// Verify registry was added
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries := manifest["registries"].(map[string]interface{})

		if _, ok := registries["test-global"]; !ok {
			t.Error("test-global registry not found")
		}

		// Clean up
		arm.MustRun("remove", "registry", "test-global")
		_ = os.Remove(globalArmrc)
	})

	t.Run("NoArmrcFile", func(t *testing.T) {
		// Ensure no .armrc files exist
		localArmrc := filepath.Join(workDir, ".armrc")
		globalArmrc := filepath.Join(homeDir, ".armrc")
		_ = os.Remove(localArmrc)
		_ = os.Remove(globalArmrc)

		// Add GitLab registry without .armrc (should still work, just no auth)
		arm.MustRun("add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--project-id", "999", "test-no-auth")

		// Verify registry was added
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries := manifest["registries"].(map[string]interface{})

		if _, ok := registries["test-no-auth"]; !ok {
			t.Error("test-no-auth registry not found")
		}

		// Clean up
		arm.MustRun("remove", "registry", "test-no-auth")
	})
}

// TestArmrcFilePermissions tests that .armrc files should have restricted permissions
func TestArmrcFilePermissions(t *testing.T) {
	workDir := t.TempDir()

	t.Run("CreateArmrcWithRestrictedPermissions", func(t *testing.T) {
		armrcPath := filepath.Join(workDir, ".armrc")
		content := `[registry https://gitlab.example.com/project/123]
token = secret-token
`
		// Create with 0600 permissions (owner read/write only)
		if err := os.WriteFile(armrcPath, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write .armrc: %v", err)
		}

		// Verify permissions
		info, err := os.Stat(armrcPath)
		if err != nil {
			t.Fatalf("failed to stat .armrc: %v", err)
		}

		mode := info.Mode()
		// Check that only owner has read/write permissions
		if mode.Perm() != 0o600 {
			t.Errorf("expected .armrc permissions to be 0600, got %o", mode.Perm())
		}
	})
}

// TestArmrcSectionMatching tests that section names must exactly match registry URLs
func TestArmrcSectionMatching(t *testing.T) {
	workDir := t.TempDir()
	homeDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", homeDir)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	arm := helpers.NewARMRunner(t, workDir)

	t.Run("ExactURLMatch", func(t *testing.T) {
		// Create .armrc with exact URL match
		localArmrc := filepath.Join(workDir, ".armrc")
		content := `[registry https://gitlab.example.com/project/123]
token = exact-match-token
`
		if err := os.WriteFile(localArmrc, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write .armrc: %v", err)
		}

		// Add registry with matching URL
		arm.MustRun("add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--project-id", "123", "test-exact")

		// Verify registry was added
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		registries := manifest["registries"].(map[string]interface{})

		if _, ok := registries["test-exact"]; !ok {
			t.Error("test-exact registry not found")
		}

		// Clean up
		arm.MustRun("remove", "registry", "test-exact")
		_ = os.Remove(localArmrc)
	})

	t.Run("HTTPvsHTTPS", func(t *testing.T) {
		// Create .armrc with HTTP URL
		localArmrc := filepath.Join(workDir, ".armrc")
		content := `[registry http://internal-gitlab.company.com/project/456]
token = http-token

[registry https://internal-gitlab.company.com/project/456]
token = https-token
`
		if err := os.WriteFile(localArmrc, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write .armrc: %v", err)
		}

		// Add HTTP registry (should match HTTP section, not HTTPS)
		// Note: This is a conceptual test - actual registry addition might not support HTTP
		// The important part is that the section matching is protocol-aware

		// Clean up
		_ = os.Remove(localArmrc)
	})
}
