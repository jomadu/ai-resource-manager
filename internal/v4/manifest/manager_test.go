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

// Package operation tests

func TestFileManager_GetAllPackagesConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name: "success - file exists with packages",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]interface{}{
						"registry1/package1": map[string]interface{}{
							"version": "1.0.0",
							"sinks":   []string{"sink1"},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			want: map[string]interface{}{
				"registry1/package1": map[string]interface{}{
					"version": "1.0.0",
					"sinks":   []interface{}{"sink1"},
				},
			},
			wantErr: false,
		},
		{
			name: "success - file doesn't exist, returns empty map",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			want:    make(map[string]interface{}),
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
				for key, wantConfig := range tt.want {
					gotConfig, exists := got[key]
					if !exists {
						t.Errorf("GetAllPackagesConfig() package %v not found", key)
						continue
					}
					if !mapsEqual(gotConfig.(map[string]interface{}), wantConfig.(map[string]interface{})) {
						t.Errorf("GetAllPackagesConfig() package %v = %v, want %v", key, gotConfig, wantConfig)
					}
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
					Packages: map[string]interface{}{
						"registry1/package1": map[string]interface{}{
							"version": "1.0.0",
							"sinks":   []string{"sink1"},
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
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "package1",
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

func TestFileManager_UpdatePackageConfigName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name            string
		setupFile       func(t *testing.T) string
		registryName    string
		packageName     string
		newRegistryName string
		newPackageName  string
		wantErr         bool
		errContains     string
	}{
		{
			name: "success - rename package",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]interface{}{
						"registry1/package1": map[string]interface{}{
							"version": "1.0.0",
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName:    "registry1",
			packageName:     "package1",
			newRegistryName: "registry1",
			newPackageName:  "package2",
			wantErr:         false,
		},
		{
			name: "error - package not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages:   make(map[string]interface{}),
					Sinks:      make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName:    "registry1",
			packageName:     "nonexistent",
			newRegistryName: "registry1",
			newPackageName:  "package2",
			wantErr:         true,
			errContains:     "not found",
		},
		{
			name: "error - new name already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]interface{}{
						"registry1/package1": map[string]interface{}{
							"version": "1.0.0",
						},
						"registry1/package2": map[string]interface{}{
							"version": "2.0.0",
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName:    "registry1",
			packageName:     "package1",
			newRegistryName: "registry1",
			newPackageName:  "package2",
			wantErr:         true,
			errContains:     "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.UpdatePackageConfigName(ctx, tt.registryName, tt.packageName, tt.newRegistryName, tt.newPackageName)

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
					Packages: map[string]interface{}{
						"registry1/package1": map[string]interface{}{
							"version": "1.0.0",
						},
						"registry1/package2": map[string]interface{}{
							"version": "2.0.0",
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
					Packages:   make(map[string]interface{}),
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
					Packages: map[string]interface{}{
						"registry1/ruleset1": map[string]interface{}{
							"version":      "1.0.0",
							"resourceType": "ruleset",
							"priority":     100,
							"sinks":        []string{"sink1"},
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
					Packages:   make(map[string]interface{}),
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
					Packages: map[string]interface{}{
						"registry1/promptset1": map[string]interface{}{
							"version":      "1.0.0",
							"resourceType": "promptset",
							"sinks":        []string{"sink1"},
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

func TestFileManager_AddRulesetConfig(t *testing.T) {
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
			name: "success - add ruleset",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "empty.json")
			},
			registryName: "registry1",
			packageName:  "ruleset1",
			config: &RulesetConfig{
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
			name: "error - package already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    "1.0.0",
					Registries: make(map[string]map[string]interface{}),
					Packages: map[string]interface{}{
						"registry1/ruleset1": map[string]interface{}{
							"version": "1.0.0",
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			registryName: "registry1",
			packageName:  "ruleset1",
			config: &RulesetConfig{
				PackageConfig: PackageConfig{Version: "2.0.0"},
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.AddRulesetConfig(ctx, tt.registryName, tt.packageName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddRulesetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddRulesetConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("AddRulesetConfig() error = %v, should contain %v", err, tt.errContains)
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
					Packages: map[string]interface{}{
						"registry1/ruleset1": map[string]interface{}{
							"version":      "1.0.0",
							"resourceType": "ruleset",
							"priority":     100,
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
					Packages:   make(map[string]interface{}),
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
		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}