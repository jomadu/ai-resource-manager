package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

func TestCloudsmithRegistry_loadToken(t *testing.T) {
	tests := []struct {
		name          string
		configMgr     *mockConfigManager
		existingToken string
		wantToken     string
		wantErr       bool
		errContains   string
	}{
		{
			name: "loads token from config",
			configMgr: &mockConfigManager{
				sections: map[string]map[string]string{
					"registry https://api.cloudsmith.io/myorg/myrepo": {
						"token": "test-token-123",
					},
				},
			},
			wantToken: "test-token-123",
			wantErr:   false,
		},
		{
			name: "returns error when token not found",
			configMgr: &mockConfigManager{
				sections: map[string]map[string]string{},
			},
			wantErr:     true,
			errContains: "failed to load token from .armrc",
		},
		{
			name: "skips loading if token already set",
			configMgr: &mockConfigManager{
				sections: map[string]map[string]string{},
			},
			existingToken: "existing-token",
			wantToken:     "existing-token",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &CloudsmithRegistry{
				name: "test-registry",
				config: CloudsmithRegistryConfig{
					RegistryConfig: RegistryConfig{
						URL: "https://api.cloudsmith.io",
					},
					Owner:      "myorg",
					Repository: "myrepo",
				},
				configMgr: tt.configMgr,
				client: &cloudsmithClient{
					token: tt.existingToken,
				},
			}

			err := reg.loadToken(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if reg.client.token != tt.wantToken {
				t.Errorf("token = %q, want %q", reg.client.token, tt.wantToken)
			}
		})
	}
}

func TestCloudsmithRegistry_loadToken_NoConfigManager(t *testing.T) {
	reg := &CloudsmithRegistry{
		name: "test-registry",
		config: CloudsmithRegistryConfig{
			RegistryConfig: RegistryConfig{
				URL: "https://api.cloudsmith.io",
			},
			Owner:      "myorg",
			Repository: "myrepo",
		},
		configMgr: nil,
		client: &cloudsmithClient{
			token: "",
		},
	}

	err := reg.loadToken(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "no token configured") {
		t.Errorf("error %q does not contain %q", err.Error(), "no token configured")
	}
}

func TestCloudsmithClient_makeRequest(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantHeader string
	}{
		{
			name:       "sets authorization header with token",
			token:      "test-token-123",
			wantHeader: "Token test-token-123",
		},
		{
			name:       "no authorization header when no token",
			token:      "",
			wantHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader != tt.wantHeader {
					t.Errorf("Authorization header = %q, want %q", authHeader, tt.wantHeader)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := &cloudsmithClient{
				baseURL:    server.URL,
				httpClient: &http.Client{},
				token:      tt.token,
			}

			resp, err := client.makeRequest(context.Background(), "GET", "/test")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestCloudsmithRegistry_ListPackageVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/packages/myorg/myrepo/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "test-package" {
			t.Errorf("unexpected query: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"name": "test-package", "version": "1.0.0", "format": "raw", "filename": "test-package-1.0.0.tar.gz"},
			{"name": "test-package", "version": "2.0.0", "format": "raw", "filename": "test-package-2.0.0.tar.gz"},
			{"name": "test-package", "version": "1.5.0", "format": "raw", "filename": "test-package-1.5.0.tar.gz"},
			{"name": "other-package", "version": "3.0.0", "format": "raw", "filename": "other-package-3.0.0.tar.gz"}
		]`))
	}))
	defer server.Close()

	reg := &CloudsmithRegistry{
		name: "test-registry",
		config: CloudsmithRegistryConfig{
			RegistryConfig: RegistryConfig{
				URL: server.URL,
			},
			Owner:      "myorg",
			Repository: "myrepo",
		},
		configMgr: &mockConfigManager{
			sections: map[string]map[string]string{
				"registry " + server.URL + "/myorg/myrepo": {
					"token": "test-token",
				},
			},
		},
		client: &cloudsmithClient{
			baseURL:    server.URL,
			httpClient: &http.Client{},
		},
	}

	versions, err := reg.ListPackageVersions(context.Background(), "test-package")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(versions))
	}

	// Should be sorted descending by semver
	expectedVersions := []string{"2.0.0", "1.5.0", "1.0.0"}
	for i, expected := range expectedVersions {
		if versions[i].Version != expected {
			t.Errorf("version[%d] = %s, want %s", i, versions[i].Version, expected)
		}
	}
}

func TestCloudsmithRegistry_ListPackageVersions_Pagination(t *testing.T) {
	callCount := 0
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		if callCount == 1 {
			w.Header().Set("Link", fmt.Sprintf("<%s/v1/packages/myorg/myrepo/?page=2>; rel=\"next\"", server.URL))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"name": "test-package", "version": "1.0.0", "format": "raw", "filename": "test-package-1.0.0.tar.gz"}
			]`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"name": "test-package", "version": "2.0.0", "format": "raw", "filename": "test-package-2.0.0.tar.gz"}
			]`))
		}
	}))
	defer server.Close()

	reg := &CloudsmithRegistry{
		name: "test-registry",
		config: CloudsmithRegistryConfig{
			RegistryConfig: RegistryConfig{
				URL: server.URL,
			},
			Owner:      "myorg",
			Repository: "myrepo",
		},
		configMgr: &mockConfigManager{
			sections: map[string]map[string]string{
				"registry " + server.URL + "/myorg/myrepo": {
					"token": "test-token",
				},
			},
		},
		client: &cloudsmithClient{
			baseURL:    server.URL,
			httpClient: &http.Client{},
		},
	}

	versions, err := reg.ListPackageVersions(context.Background(), "test-package")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
}

func TestParseNextURLFromLinkHeader(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		want       string
	}{
		{
			name:       "extracts next URL",
			linkHeader: `<https://api.cloudsmith.io/v1/packages/org/repo/?page=2>; rel="next"`,
			want:       "/v1/packages/org/repo/?page=2",
		},
		{
			name:       "handles multiple links",
			linkHeader: `<https://api.cloudsmith.io/v1/packages/org/repo/?page=1>; rel="prev", <https://api.cloudsmith.io/v1/packages/org/repo/?page=3>; rel="next"`,
			want:       "/v1/packages/org/repo/?page=3",
		},
		{
			name:       "returns empty for no next link",
			linkHeader: `<https://api.cloudsmith.io/v1/packages/org/repo/?page=1>; rel="prev"`,
			want:       "",
		},
		{
			name:       "returns empty for empty header",
			linkHeader: "",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseNextURLFromLinkHeader(tt.linkHeader)
			if got != tt.want {
				t.Errorf("parseNextURLFromLinkHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCloudsmithRegistry_ResolveVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"name": "test-package", "version": "1.0.0", "format": "raw", "filename": "test-package-1.0.0.tar.gz"},
			{"name": "test-package", "version": "2.0.0", "format": "raw", "filename": "test-package-2.0.0.tar.gz"},
			{"name": "test-package", "version": "1.5.0", "format": "raw", "filename": "test-package-1.5.0.tar.gz"}
		]`))
	}))
	defer server.Close()

	reg := &CloudsmithRegistry{
		name: "test-registry",
		config: CloudsmithRegistryConfig{
			RegistryConfig: RegistryConfig{
				URL: server.URL,
			},
			Owner:      "myorg",
			Repository: "myrepo",
		},
		configMgr: &mockConfigManager{
			sections: map[string]map[string]string{
				"registry " + server.URL + "/myorg/myrepo": {
					"token": "test-token",
				},
			},
		},
		client: &cloudsmithClient{
			baseURL:    server.URL,
			httpClient: &http.Client{},
		},
	}

	tests := []struct {
		name        string
		constraint  string
		wantVersion string
		wantErr     bool
	}{
		{
			name:        "resolves exact version",
			constraint:  "1.5.0",
			wantVersion: "1.5.0",
			wantErr:     false,
		},
		{
			name:        "resolves caret constraint",
			constraint:  "^1.0.0",
			wantVersion: "1.5.0",
			wantErr:     false,
		},
		{
			name:        "resolves tilde constraint",
			constraint:  "~1.0.0",
			wantVersion: "1.0.0",
			wantErr:     false,
		},
		{
			name:        "resolves latest",
			constraint:  "latest",
			wantVersion: "2.0.0",
			wantErr:     false,
		},
		{
			name:       "returns error for no match",
			constraint: "^3.0.0",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := reg.ResolveVersion(context.Background(), "test-package", tt.constraint)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if version.Version != tt.wantVersion {
				t.Errorf("version = %s, want %s", version.Version, tt.wantVersion)
			}
		})
	}
}

func TestCloudsmithRegistry_GetPackage(t *testing.T) {
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/v1/packages/") {
			packages := []map[string]interface{}{
				{
					"name":     "test-package",
					"version":  "1.0.0",
					"format":   "raw",
					"filename": "test-file.txt",
					"cdn_url":  serverURL + "/download/test-file.txt",
					"size":     11,
				},
			}
			json.NewEncoder(w).Encode(packages)
		} else if strings.Contains(r.URL.Path, "/download/") {
			w.Write([]byte("test content"))
		}
	}))
	defer server.Close()
	serverURL = server.URL

	configMgr := &mockConfigManager{
		sections: map[string]map[string]string{
			"registry " + server.URL + "/testowner/testrepo": {
				"token": "test-token",
			},
		},
	}

	registry, err := NewCloudsmithRegistry("test", CloudsmithRegistryConfig{
		RegistryConfig: RegistryConfig{
			Type: "cloudsmith",
			URL:  server.URL,
		},
		Owner:      "testowner",
		Repository: "testrepo",
	}, configMgr)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	version, _ := core.NewVersion("1.0.0")
	pkg, err := registry.GetPackage(context.Background(), "test-package", version, nil, nil)
	if err != nil {
		t.Fatalf("GetContent failed: %v", err)
	}

	if pkg.Metadata.Name != "test-package" {
		t.Errorf("Expected package name test-package, got %s", pkg.Metadata.Name)
	}

	if len(pkg.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(pkg.Files))
	}

	if pkg.Files[0].Path != "test-file.txt" {
		t.Errorf("Expected file path test-file.txt, got %s", pkg.Files[0].Path)
	}

	if string(pkg.Files[0].Content) != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", string(pkg.Files[0].Content))
	}
}

func TestCloudsmithRegistry_GetPackage_WithIncludeExclude(t *testing.T) {
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/v1/packages/") {
			packages := []map[string]interface{}{
				{
					"name":     "test-package",
					"version":  "1.0.0",
					"format":   "raw",
					"filename": "file1.txt",
					"cdn_url":  serverURL + "/download/file1.txt",
					"size":     5,
				},
				{
					"name":     "test-package",
					"version":  "1.0.0",
					"format":   "raw",
					"filename": "file2.md",
					"cdn_url":  serverURL + "/download/file2.md",
					"size":     5,
				},
			}
			json.NewEncoder(w).Encode(packages)
		} else if strings.Contains(r.URL.Path, "/download/file1.txt") {
			w.Write([]byte("file1"))
		} else if strings.Contains(r.URL.Path, "/download/file2.md") {
			w.Write([]byte("file2"))
		}
	}))
	defer server.Close()
	serverURL = server.URL

	configMgr := &mockConfigManager{
		sections: map[string]map[string]string{
			"registry " + server.URL + "/testowner/testrepo": {
				"token": "test-token",
			},
		},
	}

	registry, err := NewCloudsmithRegistry("test", CloudsmithRegistryConfig{
		RegistryConfig: RegistryConfig{
			Type: "cloudsmith",
			URL:  server.URL,
		},
		Owner:      "testowner",
		Repository: "testrepo",
	}, configMgr)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	version, _ := core.NewVersion("1.0.0")
	pkg, err := registry.GetPackage(context.Background(), "test-package", version, []string{"*.txt"}, nil)
	if err != nil {
		t.Fatalf("GetContent failed: %v", err)
	}

	if len(pkg.Files) != 1 {
		t.Fatalf("Expected 1 file after filtering, got %d", len(pkg.Files))
	}

	if pkg.Files[0].Path != "file1.txt" {
		t.Errorf("Expected file1.txt, got %s", pkg.Files[0].Path)
	}
}

func TestCloudsmithRegistry_GetPackage_Cache(t *testing.T) {
	callCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if strings.Contains(r.URL.Path, "/v1/packages/") {
			packages := []map[string]interface{}{
				{
					"name":     "test-package",
					"version":  "1.0.0",
					"format":   "raw",
					"filename": "test-file.txt",
					"cdn_url":  serverURL + "/download/test-file.txt",
					"size":     11,
				},
			}
			json.NewEncoder(w).Encode(packages)
		} else if strings.Contains(r.URL.Path, "/download/") {
			w.Write([]byte("test content"))
		}
	}))
	defer server.Close()
	serverURL = server.URL

	configMgr := &mockConfigManager{
		sections: map[string]map[string]string{
			"registry " + server.URL + "/testowner/testrepo": {
				"token": "test-token",
			},
		},
	}

	registry, err := NewCloudsmithRegistry("test", CloudsmithRegistryConfig{
		RegistryConfig: RegistryConfig{
			Type: "cloudsmith",
			URL:  server.URL,
		},
		Owner:      "testowner",
		Repository: "testrepo",
	}, configMgr)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	version, _ := core.NewVersion("1.0.0")
	
	// First call - should hit server
	_, err = registry.GetPackage(context.Background(), "test-package", version, nil, nil)
	if err != nil {
		t.Fatalf("GetContent failed: %v", err)
	}
	firstCallCount := callCount

	// Second call - should use cache
	_, err = registry.GetPackage(context.Background(), "test-package", version, nil, nil)
	if err != nil {
		t.Fatalf("GetContent failed on second call: %v", err)
	}

	if callCount != firstCallCount {
		t.Errorf("Expected cache hit, but server was called again")
	}
}
