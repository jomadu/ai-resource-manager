package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/rcfile"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestGitLabRegistry_ListVersions(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		packages := []GitLabPackage{
			{ID: 1, Name: "ai-rules", Version: "1.2.0", PackageType: "generic", CreatedAt: time.Now()},
			{ID: 2, Name: "ai-rules", Version: "1.1.0", PackageType: "generic", CreatedAt: time.Now()},
			{ID: 3, Name: "ai-rules", Version: "1.0.0", PackageType: "generic", CreatedAt: time.Now()},
			{ID: 4, Name: "other", Version: "2.0.0", PackageType: "npm", CreatedAt: time.Now()}, // Should be filtered out
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(packages)
	}))
	defer server.Close()

	// Create test registry
	config := GitLabRegistryConfig{
		RegistryConfig: RegistryConfig{URL: server.URL, Type: "gitlab"},
		ProjectID:      "123",
		APIVersion:     "v4",
	}

	mockCache := &mockRegistryRulesetCache{}
	registry := NewGitLabRegistry("test-registry", &config, mockCache)

	// Create temporary .armrc file
	tmpDir := t.TempDir()

	// Override rcfile service to use test directory
	registry.rcService = rcfile.NewServiceWithPaths(tmpDir, "/nonexistent")
	armrcPath := filepath.Join(tmpDir, ".armrc")
	// Extract host from server URL for .armrc format
	serverHost := strings.TrimPrefix(server.URL, "http://")
	serverHost = strings.TrimPrefix(serverHost, "https://")
	armrcContent := fmt.Sprintf("[registry %s/project/123]\ntoken = test-token\n", serverHost)
	err := os.WriteFile(armrcPath, []byte(armrcContent), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(tmpDir)

	versions, err := registry.ListVersions(context.Background())
	if err != nil {
		t.Fatalf("ListVersions() error = %v", err)
	}

	// Should return versions sorted by semver (descending)
	expected := []string{"1.2.0", "1.1.0", "1.0.0"}
	if len(versions) != len(expected) {
		t.Fatalf("Expected %d versions, got %d", len(expected), len(versions))
	}

	for i, expectedVersion := range expected {
		if versions[i].Version != expectedVersion {
			t.Errorf("Expected version %s at index %d, got %s", expectedVersion, i, versions[i].Version)
		}
	}
}

func TestGitLabRegistry_GetContent(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v4/projects/123/packages":
			packages := []GitLabPackage{
				{ID: 1, Name: "ai-rules", Version: "1.0.0", PackageType: "generic"},
			}
			_ = json.NewEncoder(w).Encode(packages)
		case "/api/v4/projects/123/packages/1/package_files":
			files := []GitLabPackageFile{
				{ID: 1, FileName: "clean-code.yml", Size: 150},
				{ID: 2, FileName: "security.yml", Size: 200},
				{ID: 3, FileName: "build/cursor/clean-code/rule-1.mdc", Size: 100},
			}
			_ = json.NewEncoder(w).Encode(files)
		case "/api/v4/projects/123/packages/generic/ai-rules/1.0.0/clean-code.yml":
			_, _ = w.Write([]byte("version: 1.0\nmetadata:\n  id: clean-code"))
		case "/api/v4/projects/123/packages/generic/ai-rules/1.0.0/security.yml":
			_, _ = w.Write([]byte("version: 1.0\nmetadata:\n  id: security"))
		case "/api/v4/projects/123/packages/generic/ai-rules/1.0.0/build/cursor/clean-code/rule-1.mdc":
			_, _ = w.Write([]byte("# Test Rule\nThis is a test rule."))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := GitLabRegistryConfig{
		RegistryConfig: RegistryConfig{URL: server.URL, Type: "gitlab"},
		ProjectID:      "123",
		APIVersion:     "v4",
	}

	mockCache := &mockRegistryRulesetCache{}
	registry := NewGitLabRegistry("test-registry-2", &config, mockCache)

	// Create temporary .armrc file
	tmpDir := t.TempDir()

	// Override rcfile service to use test directory
	registry.rcService = rcfile.NewServiceWithPaths(tmpDir, "/nonexistent")
	armrcPath := filepath.Join(tmpDir, ".armrc")
	// Extract host from server URL for .armrc format
	serverHost := strings.TrimPrefix(server.URL, "http://")
	serverHost = strings.TrimPrefix(serverHost, "https://")
	armrcContent := fmt.Sprintf("[registry %s/project/123]\ntoken = test-token\n", serverHost)
	err := os.WriteFile(armrcPath, []byte(armrcContent), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(tmpDir)

	version := types.Version{Version: "1.0.0", Display: "1.0.0"}
	selector := types.ContentSelector{Include: []string{"*.yml", "*.yaml"}} // URF files

	files, err := registry.GetContent(context.Background(), version, selector)
	if err != nil {
		t.Fatalf("GetContent() error = %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Expected 2 URF files, got %d", len(files))
	}

	// Should get URF files by default
	expectedPaths := map[string]bool{
		"clean-code.yml": false,
		"security.yml":   false,
	}

	for _, file := range files {
		if _, exists := expectedPaths[file.Path]; exists {
			expectedPaths[file.Path] = true
		} else {
			t.Errorf("Unexpected file path: %s", file.Path)
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf("Expected URF file %s not found", path)
		}
	}
}

func TestGitLabClient_buildURLs(t *testing.T) {
	client := &GitLabClient{
		baseURL:    "https://gitlab.example.com",
		apiVersion: "v4",
	}

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "project package list",
			method:   func() string { return client.buildProjectPackageListURL("123") },
			expected: "https://gitlab.example.com/api/v4/projects/123/packages",
		},
		{
			name:     "project package files",
			method:   func() string { return client.buildProjectPackageFilesURL("123", 456) },
			expected: "https://gitlab.example.com/api/v4/projects/123/packages/456/package_files",
		},
		{
			name:     "project package download",
			method:   func() string { return client.buildProjectPackageDownloadURL("123", "ai-rules", "1.0.0", "test.md") },
			expected: "https://gitlab.example.com/api/v4/projects/123/packages/generic/ai-rules/1.0.0/test.md",
		},
		{
			name:     "group package list",
			method:   func() string { return client.buildGroupPackageListURL("456") },
			expected: "https://gitlab.example.com/api/v4/groups/456/packages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method()
			if got != tt.expected {
				t.Errorf("URL = %s, want %s", got, tt.expected)
			}
		})
	}
}

// Mock cache implementation for testing
type mockRegistryRulesetCache struct {
	data map[string][]types.File
}

func (m *mockRegistryRulesetCache) ListVersions(ctx context.Context, keyObj interface{}) ([]string, error) {
	return []string{}, nil
}

func (m *mockRegistryRulesetCache) GetRulesetVersion(ctx context.Context, keyObj interface{}, version string) ([]types.File, error) {
	if m.data == nil {
		return nil, fmt.Errorf("not found in cache")
	}
	files, ok := m.data[version]
	if !ok {
		return nil, fmt.Errorf("not found in cache")
	}
	return files, nil
}

func (m *mockRegistryRulesetCache) SetRulesetVersion(ctx context.Context, keyObj interface{}, version string, files []types.File) error {
	if m.data == nil {
		m.data = make(map[string][]types.File)
	}
	m.data[version] = files
	return nil
}

func (m *mockRegistryRulesetCache) InvalidateRuleset(ctx context.Context, rulesetKey string) error {
	return nil
}

func (m *mockRegistryRulesetCache) InvalidateVersion(ctx context.Context, rulesetKey, version string) error {
	return nil
}
