package rcfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestService_GetValue(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "armrc-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test .armrc file
	rcContent := `# Test .armrc file
[registry]
my-gitlab = ${GITLAB_TOKEN}
another-registry = direct-token

[other-section]
some-key = some-value
`

	rcPath := filepath.Join(tmpDir, ".armrc")
	if err := os.WriteFile(rcPath, []byte(rcContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Set environment variable for test
	t.Setenv("GITLAB_TOKEN", "test-token-123")

	// Create service with injected working directory
	service := NewServiceWithPaths(tmpDir, "/nonexistent")

	tests := []struct {
		name        string
		section     string
		key         string
		expected    string
		expectError bool
	}{
		{
			name:     "get value with env var expansion",
			section:  "registry",
			key:      "my-gitlab",
			expected: "test-token-123",
		},
		{
			name:     "get direct value",
			section:  "registry",
			key:      "another-registry",
			expected: "direct-token",
		},
		{
			name:        "key not found",
			section:     "registry",
			key:         "nonexistent",
			expectError: true,
		},
		{
			name:        "section not found",
			section:     "nonexistent",
			key:         "some-key",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := service.GetValue(tt.section, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, value)
			}
		})
	}
}

func TestService_GetSection(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "armrc-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test .armrc file with multi-credential section
	rcContent := `# Test .armrc file
[registry https://nexus.company.com]
username = ${NEXUS_USER}
password = ${NEXUS_PASS}
url = https://nexus.company.com

[registry https://gitlab.example.com/project/123]
token = ${GITLAB_TOKEN}
`

	rcPath := filepath.Join(tmpDir, ".armrc")
	if err := os.WriteFile(rcPath, []byte(rcContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Set environment variables for test
	t.Setenv("NEXUS_USER", "test-user")
	t.Setenv("NEXUS_PASS", "test-pass")
	t.Setenv("GITLAB_TOKEN", "gitlab-token-123")

	// Create service with injected working directory
	service := NewServiceWithPaths(tmpDir, "/nonexistent")

	tests := []struct {
		name     string
		section  string
		expected map[string]string
	}{
		{
			name:    "nexus multi-credential",
			section: "registry https://nexus.company.com",
			expected: map[string]string{
				"username": "test-user",
				"password": "test-pass",
				"url":      "https://nexus.company.com",
			},
		},
		{
			name:    "gitlab single credential",
			section: "registry https://gitlab.example.com/project/123",
			expected: map[string]string{
				"token": "gitlab-token-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := service.GetSection(tt.section)
			if err != nil {
				t.Errorf("GetSection() error = %v", err)
				return
			}

			for key, expectedValue := range tt.expected {
				if actualValue, exists := values[key]; !exists {
					t.Errorf("Expected key %s not found", key)
				} else if actualValue != expectedValue {
					t.Errorf("Key %s: expected %q, got %q", key, expectedValue, actualValue)
				}
			}

			if len(values) != len(tt.expected) {
				t.Errorf("Expected %d keys, got %d", len(tt.expected), len(values))
			}
		})
	}
}

func TestService_GetValue_FileNotFound(t *testing.T) {
	// Create temporary directory without .armrc file
	tmpDir, err := os.MkdirTemp("", "armrc-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create service with injected paths pointing to nonexistent files
	service := NewServiceWithPaths(tmpDir, "/nonexistent")
	_, err = service.GetValue("registry", "test")

	if err == nil {
		t.Error("expected error for missing .armrc file")
	}
}

func TestService_GetSection_SectionNotFound(t *testing.T) {
	// Create temporary directory with .armrc file
	tmpDir, err := os.MkdirTemp("", "armrc-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test .armrc file
	rcContent := `[registry existing]
token = test
`
	rcPath := filepath.Join(tmpDir, ".armrc")
	if err := os.WriteFile(rcPath, []byte(rcContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create service with injected working directory
	service := NewServiceWithPaths(tmpDir, "/nonexistent")
	_, err = service.GetSection("registry nonexistent")

	if err == nil {
		t.Error("expected error for missing section")
	}
}

func TestService_HierarchicalLookup(t *testing.T) {
	// Create temporary directories
	projectDir, err := os.MkdirTemp("", "project-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(projectDir) }()

	homeDir, err := os.MkdirTemp("", "home-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(homeDir) }()

	emptyDir, err := os.MkdirTemp("", "empty-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(emptyDir) }()

	// Create user .armrc with global token
	userRcContent := `[registry https://gitlab.example.com]
token = ${GLOBAL_TOKEN}
`
	userRcPath := filepath.Join(homeDir, ".armrc")
	if err := os.WriteFile(userRcPath, []byte(userRcContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create project .armrc with project-specific token
	projectRcContent := `[registry https://gitlab.example.com]
token = ${PROJECT_TOKEN}
`
	projectRcPath := filepath.Join(projectDir, ".armrc")
	if err := os.WriteFile(projectRcPath, []byte(projectRcContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Set environment variables
	t.Setenv("GLOBAL_TOKEN", "global-token-123")
	t.Setenv("PROJECT_TOKEN", "project-token-456")

	// Test 1: From project directory - should get project token
	projectService := NewServiceWithPaths(projectDir, homeDir)

	token, err := projectService.GetValue("registry https://gitlab.example.com", "token")
	if err != nil {
		t.Errorf("GetValue() error = %v", err)
	}
	if token != "project-token-456" {
		t.Errorf("Expected project token, got %s", token)
	}

	// Test 2: From directory without .armrc - should get user token
	emptyService := NewServiceWithPaths(emptyDir, homeDir)

	token, err = emptyService.GetValue("registry https://gitlab.example.com", "token")
	if err != nil {
		t.Errorf("GetValue() error = %v", err)
	}
	if token != "global-token-123" {
		t.Errorf("Expected global token, got %s", token)
	}
}
