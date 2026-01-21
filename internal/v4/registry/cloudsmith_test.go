package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
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
