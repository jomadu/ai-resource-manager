package registry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
)

func TestGitLabRegistry_ListPackages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		packages := []map[string]interface{}{
			{"id": 1, "name": "test-package", "version": "1.0.0", "package_type": "generic"},
			{"id": 2, "name": "other-package", "version": "2.0.0", "package_type": "generic"},
		}
		_ = json.NewEncoder(w).Encode(packages)
	}))
	defer server.Close()

	tempDir, _ := os.MkdirTemp("", "gitlab-test")
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := GitLabRegistryConfig{
		RegistryConfig: RegistryConfig{URL: server.URL, Type: "gitlab"},
		ProjectID:      "123",
	}

	configMgr := newMockConfigManager()
	configMgr.SetValue("registry "+server.URL+"/project/123", "token", "test-token")

	registry, err := NewGitLabRegistryWithPath(tempDir, "test", &config, configMgr)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	packages, err := registry.ListPackages(context.Background())
	if err != nil {
		t.Fatalf("failed to list packages: %v", err)
	}

	if len(packages) != 2 {
		t.Errorf("expected 2 packages, got %d", len(packages))
	}
}

func TestGitLabRegistry_ListPackageVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		packages := []map[string]interface{}{
			{"id": 1, "name": "test-package", "version": "1.0.0", "package_type": "generic"},
			{"id": 2, "name": "test-package", "version": "2.0.0", "package_type": "generic"},
			{"id": 3, "name": "other-package", "version": "1.0.0", "package_type": "generic"},
		}
		_ = json.NewEncoder(w).Encode(packages)
	}))
	defer server.Close()

	tempDir, _ := os.MkdirTemp("", "gitlab-test")
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := GitLabRegistryConfig{
		RegistryConfig: RegistryConfig{URL: server.URL, Type: "gitlab"},
		ProjectID:      "123",
	}

	configMgr := newMockConfigManager()
	configMgr.SetValue("registry "+server.URL+"/project/123", "token", "test-token")

	registry, err := NewGitLabRegistryWithPath(tempDir, "test", &config, configMgr)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	versions, err := registry.ListPackageVersions(context.Background(), "test-package")
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(versions))
	}
}

func TestGitLabRegistry_GetPackage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v4/projects/123/packages":
			packages := []map[string]interface{}{
				{"id": 1, "name": "test-package", "version": "1.0.0", "package_type": "generic"},
			}
			_ = json.NewEncoder(w).Encode(packages)
		case "/api/v4/projects/123/packages/1/package_files":
			files := []map[string]interface{}{
				{"id": 1, "file_name": "test.yml", "size": 12},
			}
			_ = json.NewEncoder(w).Encode(files)
		case "/api/v4/projects/123/packages/generic/test-package/1.0.0/test.yml":
			_, _ = w.Write([]byte("test content"))
		}
	}))
	defer server.Close()

	tempDir, _ := os.MkdirTemp("", "gitlab-test")
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := GitLabRegistryConfig{
		RegistryConfig: RegistryConfig{URL: server.URL, Type: "gitlab"},
		ProjectID:      "123",
	}

	configMgr := newMockConfigManager()
	configMgr.SetValue("registry "+server.URL+"/project/123", "token", "test-token")

	registry, err := NewGitLabRegistryWithPath(tempDir, "test", &config, configMgr)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	version, _ := core.ParseVersion("1.0.0")
	pkg, err := registry.GetPackage(context.Background(), "test-package", &version, nil, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(pkg.Files))
	}

	if string(pkg.Files[0].Content) != "test content" {
		t.Errorf("expected 'test content', got %s", string(pkg.Files[0].Content))
	}
}

func TestGitLabRegistry_GroupPackages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		packages := []map[string]interface{}{
			{"id": 1, "name": "group-package", "version": "1.0.0", "package_type": "generic"},
		}
		_ = json.NewEncoder(w).Encode(packages)
	}))
	defer server.Close()

	tempDir, _ := os.MkdirTemp("", "gitlab-test")
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := GitLabRegistryConfig{
		RegistryConfig: RegistryConfig{URL: server.URL, Type: "gitlab"},
		GroupID:        "456",
	}

	configMgr := newMockConfigManager()
	configMgr.SetValue("registry "+server.URL+"/group/456", "token", "test-token")

	registry, err := NewGitLabRegistryWithPath(tempDir, "test", &config, configMgr)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	packages, err := registry.ListPackages(context.Background())
	if err != nil {
		t.Fatalf("failed to list packages: %v", err)
	}

	if len(packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(packages))
	}
}

func TestGitLabRegistry_IncludeExcludePatterns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v4/projects/123/packages":
			packages := []map[string]interface{}{
				{"id": 1, "name": "test-package", "version": "1.0.0", "package_type": "generic"},
			}
			_ = json.NewEncoder(w).Encode(packages)
		case "/api/v4/projects/123/packages/1/package_files":
			files := []map[string]interface{}{
				{"id": 1, "file_name": "rule.yml", "size": 4},
				{"id": 2, "file_name": "test.md", "size": 4},
			}
			_ = json.NewEncoder(w).Encode(files)
		case "/api/v4/projects/123/packages/generic/test-package/1.0.0/rule.yml":
			_, _ = w.Write([]byte("rule"))
		case "/api/v4/projects/123/packages/generic/test-package/1.0.0/test.md":
			_, _ = w.Write([]byte("test"))
		}
	}))
	defer server.Close()

	tempDir, _ := os.MkdirTemp("", "gitlab-test")
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := GitLabRegistryConfig{
		RegistryConfig: RegistryConfig{URL: server.URL, Type: "gitlab"},
		ProjectID:      "123",
	}

	configMgr := newMockConfigManager()
	configMgr.SetValue("registry "+server.URL+"/project/123", "token", "test-token")

	registry, err := NewGitLabRegistryWithPath(tempDir, "test", &config, configMgr)
	if err != nil {
		t.Fatalf("failed to create registry: %v", err)
	}

	version, _ := core.ParseVersion("1.0.0")
	pkg, err := registry.GetPackage(context.Background(), "test-package", &version, []string{"*.yml"}, nil)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if len(pkg.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(pkg.Files))
	}

	if pkg.Files[0].Path != "rule.yml" {
		t.Errorf("expected 'rule.yml', got %s", pkg.Files[0].Path)
	}
}
