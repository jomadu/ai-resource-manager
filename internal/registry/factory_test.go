package registry

import (
	"testing"
)

func TestNewRegistry_GitLab(t *testing.T) {
	tests := []struct {
		name      string
		rawConfig map[string]interface{}
		wantType  string
		wantErr   bool
	}{
		{
			name: "valid gitlab project registry",
			rawConfig: map[string]interface{}{
				"type":        "gitlab",
				"url":         "https://gitlab.example.com",
				"project_id":  "123",
				"api_version": "v4",
			},
			wantType: "*registry.GitLabRegistry",
			wantErr:  false,
		},
		{
			name: "valid gitlab group registry",
			rawConfig: map[string]interface{}{
				"type":        "gitlab",
				"url":         "https://gitlab.example.com",
				"group_id":    "456",
				"api_version": "v4",
			},
			wantType: "*registry.GitLabRegistry",
			wantErr:  false,
		},
		{
			name: "missing type",
			rawConfig: map[string]interface{}{
				"url": "https://gitlab.example.com",
			},
			wantErr: true,
		},
		{
			name: "unsupported type",
			rawConfig: map[string]interface{}{
				"type": "unsupported",
				"url":  "https://example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid gitlab config",
			rawConfig: map[string]interface{}{
				"type": "gitlab",
				"url":  123, // Invalid type
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRegistry(tt.rawConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check that we got the right type
				gotType := getTypeName(got)
				if gotType != tt.wantType {
					t.Errorf("NewRegistry() type = %v, want %v", gotType, tt.wantType)
				}
			}
		})
	}
}

func TestNewGitLabRegistry(t *testing.T) {
	tests := []struct {
		name      string
		rawConfig map[string]interface{}
		wantErr   bool
	}{
		{
			name: "valid project config",
			rawConfig: map[string]interface{}{
				"type":        "gitlab",
				"url":         "https://gitlab.example.com",
				"project_id":  "123",
				"api_version": "v4",
			},
			wantErr: false,
		},
		{
			name: "valid group config",
			rawConfig: map[string]interface{}{
				"type":        "gitlab",
				"url":         "https://gitlab.example.com",
				"group_id":    "456",
				"api_version": "v4",
			},
			wantErr: false,
		},
		{
			name: "missing url",
			rawConfig: map[string]interface{}{
				"type":        "gitlab",
				"project_id":  "123",
				"api_version": "v4",
			},
			wantErr: false, // URL is optional in config parsing
		},
		{
			name: "invalid json structure",
			rawConfig: map[string]interface{}{
				"type": "gitlab",
				"url":  make(chan int), // Cannot be marshaled to JSON
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newGitLabRegistry(tt.rawConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("newGitLabRegistry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to get type name for testing
func getTypeName(v interface{}) string {
	switch v.(type) {
	case *GitRegistry:
		return "*registry.GitRegistry"
	case *GitLabRegistry:
		return "*registry.GitLabRegistry"
	default:
		return "unknown"
	}
}
