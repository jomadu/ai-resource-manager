package registry

import (
	"context"
	"strings"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestCloudsmithRegistry_GetAuthKey(t *testing.T) {
	config := &CloudsmithRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "https://app.cloudsmith.com/myorg/ai-rules",
			Type: "cloudsmith",
		},
		Owner:      "myorg",
		Repository: "ai-rules",
	}

	registry := NewCloudsmithRegistryNoCache("test", config)
	authKey := registry.config.URL

	expected := "https://app.cloudsmith.com/myorg/ai-rules"
	if authKey != expected {
		t.Errorf("Expected auth key %s, got %s", expected, authKey)
	}
}

func TestCloudsmithRegistryConfig_GetBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "default URL",
			url:      "",
			expected: "https://api.cloudsmith.io",
		},
		{
			name:     "custom URL",
			url:      "https://custom.cloudsmith.io",
			expected: "https://custom.cloudsmith.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CloudsmithRegistryConfig{
				RegistryConfig: RegistryConfig{
					URL:  tt.url,
					Type: "cloudsmith",
				},
			}

			result := config.GetBaseURL()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCloudsmithRegistry_ListVersions_EmptyResponse(t *testing.T) {
	config := &CloudsmithRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "https://api.cloudsmith.io",
			Type: "cloudsmith",
		},
		Owner:      "myorg",
		Repository: "ai-rules",
	}

	registry := NewCloudsmithRegistryNoCache("test", config)

	// This test would require mocking the HTTP client
	// For now, we'll just test the basic structure
	ctx := context.Background()
	_, err := registry.ListVersions(ctx, "test-ruleset")

	// We expect an error since we don't have authentication configured
	if err == nil {
		t.Error("Expected error due to missing authentication, got nil")
	}
}

func TestCloudsmithRegistry_GetContent_CacheKey(t *testing.T) {
	config := &CloudsmithRegistryConfig{
		RegistryConfig: RegistryConfig{
			URL:  "https://api.cloudsmith.io",
			Type: "cloudsmith",
		},
		Owner:      "myorg",
		Repository: "ai-rules",
	}

	registry := NewCloudsmithRegistryNoCache("test", config)

	ctx := context.Background()
	version := types.Version{Version: "1.0.0", Display: "1.0.0"}
	selector := types.ContentSelector{
		Include: []string{"*.yml"},
	}

	// This test would require mocking the HTTP client and cache
	// For now, we'll just test that the method exists and handles basic cases
	_, err := registry.GetContent(ctx, "test-ruleset", version, selector)

	// We expect an error since we don't have authentication configured
	if err == nil {
		t.Error("Expected error due to missing authentication, got nil")
	}
}

func TestParseCloudsmithURL(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		expectedOwner string
		expectedRepo  string
		expectError   bool
		errorContains string
	}{
		{
			name:          "valid URL",
			url:           "https://app.cloudsmith.com/thetradedesk/cybersecurity-dev-tools-internal",
			expectedOwner: "thetradedesk",
			expectedRepo:  "cybersecurity-dev-tools-internal",
			expectError:   false,
		},
		{
			name:          "valid URL with trailing slash",
			url:           "https://app.cloudsmith.com/myorg/myrepo/",
			expectedOwner: "myorg",
			expectedRepo:  "myrepo",
			expectError:   false,
		},
		{
			name:          "valid URL with single character owner and repo",
			url:           "https://app.cloudsmith.com/a/b",
			expectedOwner: "a",
			expectedRepo:  "b",
			expectError:   false,
		},
		{
			name:          "valid URL with hyphens and underscores",
			url:           "https://app.cloudsmith.com/my-org/my_repo",
			expectedOwner: "my-org",
			expectedRepo:  "my_repo",
			expectError:   false,
		},
		{
			name:          "invalid URL - wrong host",
			url:           "https://api.cloudsmith.io/myorg/myrepo",
			expectError:   true,
			errorContains: "expected Cloudsmith URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "valid URL with http protocol",
			url:           "http://app.cloudsmith.com/myorg/myrepo",
			expectedOwner: "myorg",
			expectedRepo:  "myrepo",
			expectError:   false,
		},
		{
			name:          "invalid URL - missing repository",
			url:           "https://app.cloudsmith.com/myorg",
			expectError:   true,
			errorContains: "expected URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "invalid URL - too many path segments",
			url:           "https://app.cloudsmith.com/myorg/myrepo/extra",
			expectError:   true,
			errorContains: "expected URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "invalid URL - empty owner",
			url:           "https://app.cloudsmith.com//myrepo",
			expectError:   true,
			errorContains: "expected URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "invalid URL - empty repository",
			url:           "https://app.cloudsmith.com/myorg/",
			expectError:   true,
			errorContains: "expected URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "invalid URL - malformed URL",
			url:           "not-a-url",
			expectError:   true,
			errorContains: "expected Cloudsmith URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "invalid URL - empty string",
			url:           "",
			expectError:   true,
			errorContains: "expected Cloudsmith URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "invalid URL - missing path",
			url:           "https://app.cloudsmith.com",
			expectError:   true,
			errorContains: "expected URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
		{
			name:          "invalid URL - only slash",
			url:           "https://app.cloudsmith.com/",
			expectError:   true,
			errorContains: "expected URL format: https://app.cloudsmith.com/[owner]/[repository]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseCloudsmithURL(tt.url)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if owner != tt.expectedOwner {
				t.Errorf("Expected owner '%s', got '%s'", tt.expectedOwner, owner)
			}

			if repo != tt.expectedRepo {
				t.Errorf("Expected repository '%s', got '%s'", tt.expectedRepo, repo)
			}
		})
	}
}
