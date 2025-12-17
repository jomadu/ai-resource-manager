package manifest

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFileManager_GetAllRegistriesConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		want      map[string]map[string]interface{}
		wantErr   bool
	}{
		{
			name: "success - file exists with registries",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
						"registry2": {
							"type": "cloudsmith",
							"url":  "https://api.cloudsmith.io",
							"owner": "org",
							"repository": "repo",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			want: map[string]map[string]interface{}{
				"registry1": {
					"type": "git",
					"url":  "https://github.com/example/repo",
				},
				"registry2": {
					"type": "cloudsmith",
					"url":  "https://api.cloudsmith.io",
					"owner": "org",
					"repository": "repo",
				},
			},
			wantErr: false,
		},
		{
			name: "success - file doesn't exist, returns empty map",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			want:    make(map[string]map[string]interface{}),
			wantErr: false,
		},
		{
			name: "success - file exists but no registries",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: nil,
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			want:    make(map[string]map[string]interface{}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetAllRegistriesConfig(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRegistriesConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GetAllRegistriesConfig() length = %v, want %v", len(got), len(tt.want))
				}
				// Compare registry contents
				for name, wantConfig := range tt.want {
					gotConfig, exists := got[name]
					if !exists {
						t.Errorf("GetAllRegistriesConfig() registry %v not found", name)
						continue
					}
					if !mapsEqual(gotConfig, wantConfig) {
						t.Errorf("GetAllRegistriesConfig() registry %v = %v, want %v", name, gotConfig, wantConfig)
					}
				}
			}
		})
	}
}

func TestFileManager_GetRegistryConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		registryName string
		want        map[string]interface{}
		wantErr     bool
		errContains string
	}{
		{
			name: "success - registry exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			want: map[string]interface{}{
				"type": "git",
				"url":  "https://github.com/example/repo",
			},
			wantErr: false,
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "nonexistent",
			want:        nil,
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "error - file doesn't exist, registry not found",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			registryName: "registry1",
			want:        nil,
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetRegistryConfig(ctx, tt.registryName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegistryConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetRegistryConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetRegistryConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if !mapsEqual(got, tt.want) {
				t.Errorf("GetRegistryConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_AddRegistryConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		registryName string
		config      map[string]interface{}
		wantErr     bool
		errContains string
		wantFile    *Manifest
	}{
		{
			name: "success - add registry to empty manifest",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "empty.json")
			},
			registryName: "registry1",
			config: map[string]interface{}{
				"type": "git",
				"url":  "https://github.com/example/repo",
			},
			wantErr: false,
			wantFile: &Manifest{
				Version: "1.0.0",
				Registries: map[string]map[string]interface{}{
					"registry1": {
						"type": "git",
						"url":  "https://github.com/example/repo",
					},
				},
				Packages: make(map[string]map[string]interface{}),
				Sinks:    make(map[string]SinkConfig),
			},
		},
		{
			name: "success - add registry to existing manifest",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry2",
			config: map[string]interface{}{
				"type": "cloudsmith",
				"url":  "https://api.cloudsmith.io",
				"owner": "org",
				"repository": "repo",
			},
			wantErr: false,
			wantFile: &Manifest{
				Version: "1.0.0",
				Registries: map[string]map[string]interface{}{
					"registry1": {
						"type": "git",
						"url":  "https://github.com/example/repo",
					},
					"registry2": {
						"type": "cloudsmith",
						"url":  "https://api.cloudsmith.io",
						"owner": "org",
						"repository": "repo",
					},
				},
				Packages: make(map[string]map[string]interface{}),
				Sinks:    make(map[string]SinkConfig),
			},
		},
		{
			name: "error - registry already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			config: map[string]interface{}{
				"type": "git",
				"url":  "https://github.com/other/repo",
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.AddRegistryConfig(ctx, tt.registryName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddRegistryConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddRegistryConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("AddRegistryConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if tt.wantFile != nil {
				got, err := fm.loadManifest()
				if err != nil {
					t.Fatalf("loadManifest() error = %v", err)
				}
				compareManifests(t, got, tt.wantFile)
			}
		})
	}
}

func TestFileManager_UpdateRegistryConfigName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		oldName     string
		newName     string
		wantErr     bool
		errContains string
		wantFile    *Manifest
	}{
		{
			name: "success - rename registry and move packages",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
					},
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
								"sinks":   []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			oldName: "registry1",
			newName: "registry2",
			wantErr: false,
			wantFile: &Manifest{
				Version: "1.0.0",
				Registries: map[string]map[string]interface{}{
					"registry2": {
						"type": "git",
						"url":  "https://github.com/example/repo",
					},
				},
				Packages: map[string]map[string]interface{}{
					"registry2": {
						"package1": map[string]interface{}{
							"version": "1.0.0",
							"sinks":   []string{"sink1"},
						},
					},
				},
				Sinks: make(map[string]SinkConfig),
			},
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			oldName:     "nonexistent",
			newName:     "registry2",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "error - new name already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
						"registry2": {
							"type": "cloudsmith",
							"url":  "https://api.cloudsmith.io",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			oldName:     "registry1",
			newName:     "registry2",
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdateRegistryConfigName(ctx, tt.oldName, tt.newName)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRegistryConfigName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateRegistryConfigName() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateRegistryConfigName() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if tt.wantFile != nil {
				got, err := fm.loadManifest()
				if err != nil {
					t.Fatalf("loadManifest() error = %v", err)
				}
				compareManifests(t, got, tt.wantFile)
			}
		})
	}
}

func TestFileManager_UpdateRegistryConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		registryName string
		config      map[string]interface{}
		wantErr     bool
		errContains string
		wantFile    *Manifest
	}{
		{
			name: "success - update existing registry",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			config: map[string]interface{}{
				"type": "git",
				"url":  "https://github.com/other/repo",
				"branches": []string{"main", "develop"},
			},
			wantErr: false,
			wantFile: &Manifest{
				Version: "1.0.0",
				Registries: map[string]map[string]interface{}{
					"registry1": {
						"type": "git",
						"url":  "https://github.com/other/repo",
						"branches": []string{"main", "develop"},
					},
				},
				Packages: make(map[string]map[string]interface{}),
				Sinks:    make(map[string]SinkConfig),
			},
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "nonexistent",
			config: map[string]interface{}{
				"type": "git",
				"url":  "https://github.com/example/repo",
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdateRegistryConfig(ctx, tt.registryName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRegistryConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateRegistryConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateRegistryConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if tt.wantFile != nil {
				got, err := fm.loadManifest()
				if err != nil {
					t.Fatalf("loadManifest() error = %v", err)
				}
				compareManifests(t, got, tt.wantFile)
			}
		})
	}
}

func TestFileManager_RemoveRegistryConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		registryName string
		wantErr     bool
		errContains string
		wantFile    *Manifest
	}{
		{
			name: "success - remove registry",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"registry1": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
						"registry2": {
							"type": "cloudsmith",
							"url":  "https://api.cloudsmith.io",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			wantErr:     false,
			wantFile: &Manifest{
				Version: "1.0.0",
				Registries: map[string]map[string]interface{}{
					"registry2": {
						"type": "cloudsmith",
						"url":  "https://api.cloudsmith.io",
					},
				},
				Packages: make(map[string]map[string]interface{}),
				Sinks:    make(map[string]SinkConfig),
			},
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "nonexistent",
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.RemoveRegistryConfig(ctx, tt.registryName)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveRegistryConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoveRegistryConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemoveRegistryConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if tt.wantFile != nil {
				got, err := fm.loadManifest()
				if err != nil {
					t.Fatalf("loadManifest() error = %v", err)
				}
				compareManifests(t, got, tt.wantFile)
			}
		})
	}
}

// Type-safe registry helper tests

func TestFileManager_GetGitRegistryConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		registryName string
		want        *GitRegistryConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "success - git registry exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"git-registry": {
							"type":    "git",
							"url":     "https://github.com/example/repo",
							"branches": []string{"main", "develop"},
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "git-registry",
			want: &GitRegistryConfig{
				RegistryConfig: RegistryConfig{
					Type: "git",
					URL:  "https://github.com/example/repo",
				},
				Branches: []string{"main", "develop"},
			},
			wantErr: false,
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "nonexistent",
			want:        nil,
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "error - wrong registry type",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"cloudsmith-registry": {
							"type":       "cloudsmith",
							"url":        "https://api.cloudsmith.io",
							"owner":      "org",
							"repository": "repo",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "cloudsmith-registry",
			want:        nil,
			wantErr:     true,
			errContains: "not a git registry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetGitRegistryConfig(ctx, tt.registryName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetGitRegistryConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetGitRegistryConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetGitRegistryConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got.Type != tt.want.Type || got.URL != tt.want.URL {
				t.Errorf("GetGitRegistryConfig() = %+v, want %+v", got, tt.want)
			}
			if len(got.Branches) != len(tt.want.Branches) {
				t.Errorf("GetGitRegistryConfig() branches = %v, want %v", got.Branches, tt.want.Branches)
			}
		})
	}
}

func TestFileManager_AddGitRegistryConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		registryName string
		config      *GitRegistryConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "success - add git registry",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "empty.json")
			},
			registryName: "git-registry",
			config: &GitRegistryConfig{
				RegistryConfig: RegistryConfig{
					Type: "git",
					URL:  "https://github.com/example/repo",
				},
				Branches: []string{"main", "develop"},
			},
			wantErr: false,
		},
		{
			name: "error - registry already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: map[string]map[string]interface{}{
						"git-registry": {
							"type": "git",
							"url":  "https://github.com/example/repo",
						},
					},
					Packages: make(map[string]map[string]interface{}),
					Sinks:    make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "git-registry",
			config: &GitRegistryConfig{
				RegistryConfig: RegistryConfig{
					Type: "git",
					URL:  "https://github.com/other/repo",
				},
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.AddGitRegistryConfig(ctx, tt.registryName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddGitRegistryConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddGitRegistryConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("AddGitRegistryConfig() error = %v, should contain %v", err, tt.errContains)
				}
			} else {
				// Verify it was added
				got, err := fm.GetGitRegistryConfig(ctx, tt.registryName)
				if err != nil {
					t.Fatalf("GetGitRegistryConfig() error = %v", err)
				}
				if got.Type != tt.config.Type || got.URL != tt.config.URL {
					t.Errorf("GetGitRegistryConfig() = %+v, want %+v", got, tt.config)
				}
			}
		})
	}
}

// Sink operation tests

func TestFileManager_GetAllSinksConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		want      map[string]*SinkConfig
		wantErr   bool
	}{
		{
			name: "success - file exists with sinks",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks: map[string]SinkConfig{
						"sink1": {
							Directory:    ".cursor/rules",
							Layout:       "hierarchical",
							CompileTarget: "cursor",
						},
						"sink2": {
							Directory:    ".amazonq/rules",
							CompileTarget: "amazonq",
						},
					},
				}
				return createTestManifest(t, manifest)
			},
			want: map[string]*SinkConfig{
				"sink1": {
					Directory:    ".cursor/rules",
					Layout:       "hierarchical",
					CompileTarget: "cursor",
				},
				"sink2": {
					Directory:    ".amazonq/rules",
					CompileTarget: "amazonq",
				},
			},
			wantErr: false,
		},
		{
			name: "success - file doesn't exist, returns empty map",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			want:    make(map[string]*SinkConfig),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetAllSinksConfig(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllSinksConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GetAllSinksConfig() length = %v, want %v", len(got), len(tt.want))
				}
				for name, wantSink := range tt.want {
					gotSink, exists := got[name]
					if !exists {
						t.Errorf("GetAllSinksConfig() sink %v not found", name)
						continue
					}
					if gotSink.Directory != wantSink.Directory ||
						gotSink.Layout != wantSink.Layout ||
						gotSink.CompileTarget != wantSink.CompileTarget {
						t.Errorf("GetAllSinksConfig() sink %v = %+v, want %+v", name, gotSink, wantSink)
					}
				}
			}
		})
	}
}

func TestFileManager_GetSinkConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		sinkName  string
		want      *SinkConfig
		wantErr   bool
		errContains string
	}{
		{
			name: "success - sink exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks: map[string]SinkConfig{
						"sink1": {
							Directory:    ".cursor/rules",
							Layout:       "hierarchical",
							CompileTarget: "cursor",
						},
					},
				}
				return createTestManifest(t, manifest)
			},
			sinkName: "sink1",
			want: &SinkConfig{
				Directory:    ".cursor/rules",
				Layout:       "hierarchical",
				CompileTarget: "cursor",
			},
			wantErr: false,
		},
		{
			name: "error - sink not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			sinkName: "nonexistent",
			want:     nil,
			wantErr:  true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetSinkConfig(ctx, tt.sinkName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSinkConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetSinkConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetSinkConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got.Directory != tt.want.Directory ||
				got.Layout != tt.want.Layout ||
				got.CompileTarget != tt.want.CompileTarget {
				t.Errorf("GetSinkConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFileManager_AddSinkConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		sinkName    string
		config      *SinkConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "success - add sink to empty manifest",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "empty.json")
			},
			sinkName: "sink1",
			config: &SinkConfig{
				Directory:    ".cursor/rules",
				Layout:       "hierarchical",
				CompileTarget: "cursor",
			},
			wantErr: false,
		},
		{
			name: "error - sink already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version: "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks: map[string]SinkConfig{
						"sink1": {
							Directory:    ".cursor/rules",
							CompileTarget: "cursor",
						},
					},
				}
				return createTestManifest(t, manifest)
			},
			sinkName: "sink1",
			config: &SinkConfig{
				Directory:    ".amazonq/rules",
				CompileTarget: "amazonq",
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.AddSinkConfig(ctx, tt.sinkName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddSinkConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddSinkConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("AddSinkConfig() error = %v, should contain %v", err, tt.errContains)
				}
			} else {
				// Verify it was added
				got, err := fm.GetSinkConfig(ctx, tt.sinkName)
				if err != nil {
					t.Fatalf("GetSinkConfig() error = %v", err)
				}
				if got.Directory != tt.config.Directory ||
					got.CompileTarget != tt.config.CompileTarget {
					t.Errorf("GetSinkConfig() = %+v, want %+v", got, tt.config)
				}
			}
		})
	}
}

// Package operation tests

func TestFileManager_GetAllPackagesConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		want      map[string]map[string]interface{}
		wantErr   bool
	}{
		{
			name: "success - file exists with packages",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
								"sinks":   []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			want: map[string]map[string]interface{}{
				"registry1": {
					"package1": map[string]interface{}{
						"version": "1.0.0",
						"sinks":   []string{"sink1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success - file doesn't exist, returns empty map",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			want:    make(map[string]map[string]interface{}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetAllPackagesConfig(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPackagesConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GetAllPackagesConfig() length = %v, want %v", len(got), len(tt.want))
				}
			}
		})
	}
}

func TestFileManager_GetPackageConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		want         map[string]interface{}
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - package exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
								"sinks":   []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "package1",
			want: map[string]interface{}{
				"version": "1.0.0",
				"sinks":   []interface{}{"sink1"},
			},
			wantErr: false,
		},
		{
			name: "error - registry not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "nonexistent",
			packageName:  "package1",
			want:        nil,
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "nonexistent",
			want:        nil,
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetPackageConfig(ctx, tt.registryName, tt.packageName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackageConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPackageConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetPackageConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if !mapsEqual(got, tt.want) {
				t.Errorf("GetPackageConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_AddPackageConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		config       map[string]interface{}
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - add package to empty manifest",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "empty.json")
			},
			registryName: "registry1",
			packageName:  "package1",
			config: map[string]interface{}{
				"version": "1.0.0",
				"sinks":  []string{"sink1"},
			},
			wantErr: false,
		},
		{
			name: "error - package already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "package1",
			config: map[string]interface{}{
				"version": "2.0.0",
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.AddPackageConfig(ctx, tt.registryName, tt.packageName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddPackageConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddPackageConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("AddPackageConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_GetRulesetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		want         *RulesetConfig
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - ruleset exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"ruleset1": map[string]interface{}{
								"version":      "1.0.0",
								"resourceType": "ruleset",
								"priority":     100,
								"sinks":        []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "ruleset1",
			want: &RulesetConfig{
				PackageConfig: PackageConfig{
					Version:      "1.0.0",
					ResourceType: ResourceTypeRuleset,
					Sinks:        []string{"sink1"},
				},
				Priority: 100,
			},
			wantErr: false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "nonexistent",
			want:         nil,
			wantErr:      true,
			errContains:  "not found",
		},
		{
			name: "error - wrong resource type",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"promptset1": map[string]interface{}{
								"version":      "1.0.0",
								"resourceType": "promptset",
								"sinks":        []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "promptset1",
			want:         nil,
			wantErr:      true,
			errContains:  "not a ruleset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetRulesetConfig(ctx, tt.registryName, tt.packageName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRulesetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetRulesetConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetRulesetConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got.Version != tt.want.Version ||
				got.Priority != tt.want.Priority ||
				got.ResourceType != tt.want.ResourceType {
				t.Errorf("GetRulesetConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFileManager_UpdateSinkConfigName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		oldName     string
		newName     string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - rename sink",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks: map[string]SinkConfig{
						"sink1": {
							Directory:    ".cursor/rules",
							CompileTarget: "cursor",
						},
					},
				}
				return createTestManifest(t, manifest)
			},
			oldName: "sink1",
			newName: "sink2",
			wantErr: false,
		},
		{
			name: "error - sink not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			oldName:     "nonexistent",
			newName:     "sink2",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "error - new name already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks: map[string]SinkConfig{
						"sink1": {
							Directory:    ".cursor/rules",
							CompileTarget: "cursor",
						},
						"sink2": {
							Directory:    ".amazonq/rules",
							CompileTarget: "amazonq",
						},
					},
				}
				return createTestManifest(t, manifest)
			},
			oldName:     "sink1",
			newName:     "sink2",
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdateSinkConfigName(ctx, tt.oldName, tt.newName)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSinkConfigName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateSinkConfigName() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateSinkConfigName() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_UpdateSinkConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		sinkName    string
		config      *SinkConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "success - update existing sink",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks: map[string]SinkConfig{
						"sink1": {
							Directory:    ".cursor/rules",
							CompileTarget: "cursor",
						},
					},
				}
				return createTestManifest(t, manifest)
			},
			sinkName: "sink1",
			config: &SinkConfig{
				Directory:    ".cursor/new-rules",
				Layout:       "flat",
				CompileTarget: "cursor",
			},
			wantErr: false,
		},
		{
			name: "error - sink not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			sinkName:    "nonexistent",
			config:      &SinkConfig{Directory: ".cursor/rules", CompileTarget: "cursor"},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdateSinkConfig(ctx, tt.sinkName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSinkConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateSinkConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateSinkConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_RemoveSinkConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		sinkName    string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - remove sink",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks: map[string]SinkConfig{
						"sink1": {
							Directory:    ".cursor/rules",
							CompileTarget: "cursor",
						},
						"sink2": {
							Directory:    ".amazonq/rules",
							CompileTarget: "amazonq",
						},
					},
				}
				return createTestManifest(t, manifest)
			},
			sinkName: "sink1",
			wantErr: false,
		},
		{
			name: "error - sink not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			sinkName:    "nonexistent",
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.RemoveSinkConfig(ctx, tt.sinkName)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveSinkConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoveSinkConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemoveSinkConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_UpdatePackageConfigName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		newName      string
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - rename package",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "package1",
			newName:      "package2",
			wantErr:      false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "nonexistent",
			newName:      "package2",
			wantErr:      true,
			errContains:  "not found",
		},
		{
			name: "error - new name already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
							},
							"package2": map[string]interface{}{
								"version": "2.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "package1",
			newName:      "package2",
			wantErr:      true,
			errContains:  "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdatePackageConfigName(ctx, tt.registryName, tt.packageName, tt.newName)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePackageConfigName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdatePackageConfigName() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdatePackageConfigName() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_UpdatePackageConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		config       map[string]interface{}
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - update existing package",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "package1",
			config: map[string]interface{}{
				"version": "2.0.0",
				"sinks":  []string{"sink1", "sink2"},
			},
			wantErr: false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "nonexistent",
			config:       map[string]interface{}{"version": "1.0.0"},
			wantErr:      true,
			errContains:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdatePackageConfig(ctx, tt.registryName, tt.packageName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePackageConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdatePackageConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdatePackageConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_RemovePackageConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - remove package",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"package1": map[string]interface{}{
								"version": "1.0.0",
							},
							"package2": map[string]interface{}{
								"version": "2.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "package1",
			wantErr:      false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "nonexistent",
			wantErr:      true,
			errContains:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.RemovePackageConfig(ctx, tt.registryName, tt.packageName)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemovePackageConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemovePackageConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemovePackageConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_GetPromptsetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		want         *PromptsetConfig
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - promptset exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"promptset1": map[string]interface{}{
								"version":      "1.0.0",
								"resourceType": "promptset",
								"sinks":        []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "promptset1",
			want: &PromptsetConfig{
				PackageConfig: PackageConfig{
					Version:      "1.0.0",
					ResourceType: ResourceTypePromptset,
					Sinks:        []string{"sink1"},
				},
			},
			wantErr: false,
		},
		{
			name: "error - wrong resource type",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"ruleset1": map[string]interface{}{
								"version":      "1.0.0",
								"resourceType": "ruleset",
								"priority":     100,
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "ruleset1",
			want:         nil,
			wantErr:      true,
			errContains:  "not a promptset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetPromptsetConfig(ctx, tt.registryName, tt.packageName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPromptsetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPromptsetConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetPromptsetConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got.Version != tt.want.Version ||
				got.ResourceType != tt.want.ResourceType {
				t.Errorf("GetPromptsetConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFileManager_AddPromptsetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		config       *PromptsetConfig
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - add promptset",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "empty.json")
			},
			registryName: "registry1",
			packageName:  "promptset1",
			config: &PromptsetConfig{
				PackageConfig: PackageConfig{
					Version:      "1.0.0",
					ResourceType: ResourceTypePromptset,
					Sinks:        []string{"sink1"},
				},
			},
			wantErr: false,
		},
		{
			name: "error - package already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"promptset1": map[string]interface{}{
								"version": "1.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "promptset1",
			config: &PromptsetConfig{
				PackageConfig: PackageConfig{
					Version: "2.0.0",
				},
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.AddPromptsetConfig(ctx, tt.registryName, tt.packageName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddPromptsetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddPromptsetConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("AddPromptsetConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_UpdateRulesetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		config       *RulesetConfig
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - update existing ruleset",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"ruleset1": map[string]interface{}{
								"version":      "1.0.0",
								"resourceType": "ruleset",
								"priority":     100,
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "ruleset1",
			config: &RulesetConfig{
				PackageConfig: PackageConfig{
					Version:      "2.0.0",
					ResourceType: ResourceTypeRuleset,
					Sinks:        []string{"sink1"},
				},
				Priority: 200,
			},
			wantErr: false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "nonexistent",
			config: &RulesetConfig{
				PackageConfig: PackageConfig{Version: "1.0.0"},
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdateRulesetConfig(ctx, tt.registryName, tt.packageName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRulesetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateRulesetConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateRulesetConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestFileManager_UpdatePromptsetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		setupFile    func(t *testing.T) string
		registryName string
		packageName  string
		config       *PromptsetConfig
		wantErr      bool
		errContains  string
	}{
		{
			name: "success - update existing promptset",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]map[string]interface{}{
						"registry1": {
							"promptset1": map[string]interface{}{
								"version":      "1.0.0",
								"resourceType": "promptset",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "promptset1",
			config: &PromptsetConfig{
				PackageConfig: PackageConfig{
					Version:      "2.0.0",
					ResourceType: ResourceTypePromptset,
					Sinks:        []string{"sink1"},
				},
			},
			wantErr: false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "nonexistent",
			config: &PromptsetConfig{
				PackageConfig: PackageConfig{Version: "1.0.0"},
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdatePromptsetConfig(ctx, tt.registryName, tt.packageName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePromptsetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdatePromptsetConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdatePromptsetConfig() error = %v, should contain %v", err, tt.errContains)
				}
			}
		})
	}
}

// Helper functions

func createTestManifest(t *testing.T, manifest *Manifest) string {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "test-manifest.json")

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal manifest: %v", err)
	}

	err = os.WriteFile(manifestPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	return manifestPath
}

func mapsEqual(m1, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			return false
		}
		// Use reflect.DeepEqual for proper comparison of slices and nested structures
		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	return true
}

func compareManifests(t *testing.T, got, want *Manifest) {
	if got.Version != want.Version {
		t.Errorf("Manifest Version = %v, want %v", got.Version, want.Version)
	}
	// Compare registries
	if len(got.Registries) != len(want.Registries) {
		t.Errorf("Manifest Registries length = %v, want %v", len(got.Registries), len(want.Registries))
	}
	// Compare packages
	if len(got.Packages) != len(want.Packages) {
		t.Errorf("Manifest Packages length = %v, want %v", len(got.Packages), len(want.Packages))
	}
	// Compare sinks
	if len(got.Sinks) != len(want.Sinks) {
		t.Errorf("Manifest Sinks length = %v, want %v", len(got.Sinks), len(want.Sinks))
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

