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

func TestFileManager_GetAllRulesets(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupFile func(t *testing.T) string
		want      map[string]*RulesetConfig
		wantErr   bool
	}{
		{
			name: "success - file exists with rulesets",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    1,
					Registries: make(map[string]map[string]interface{}),
					Dependencies: Dependencies{
						Rulesets: map[string]RulesetConfig{
							"registry1/ruleset1": {
								Version: "1.0.0",
								Sinks:   []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			want: map[string]*RulesetConfig{
				"registry1/ruleset1": {
					Version: "1.0.0",
					Sinks:   []string{"sink1"},
				},
			},
			wantErr: false,
		},
		{
			name: "success - file doesn't exist, returns empty map",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent.json")
			},
			want:    make(map[string]*RulesetConfig),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetAllRulesets(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRulesets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GetAllRulesets() length = %v, want %v", len(got), len(tt.want))
				}
				for key, wantConfig := range tt.want {
					gotConfig, exists := got[key]
					if !exists {
						t.Errorf("GetAllRulesets() ruleset %v not found", key)
						continue
					}
					if !reflect.DeepEqual(gotConfig, wantConfig) {
						t.Errorf("GetAllRulesets() ruleset %v = %v, want %v", key, gotConfig, wantConfig)
					}
				}
			}
		})
	}
}

func TestFileManager_GetRulesetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		packageKey  string
		want        *RulesetConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "success - ruleset exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    1,
					Registries: make(map[string]map[string]interface{}),
					Dependencies: Dependencies{
						Rulesets: map[string]RulesetConfig{
							"registry1/ruleset1": {
								Version:  "1.0.0",
								Priority: 100,
								Sinks:    []string{"sink1"},
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			packageKey: "registry1/ruleset1",
			want: &RulesetConfig{
				Version:  "1.0.0",
				Priority: 100,
				Sinks:    []string{"sink1"},
			},
			wantErr: false,
		},
		{
			name: "error - ruleset not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:      1,
					Registries:   make(map[string]map[string]interface{}),
					Dependencies: Dependencies{},
					Sinks:        make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			packageKey:  "registry1/nonexistent",
			want:        nil,
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			got, err := fm.GetRulesetConfig(ctx, tt.packageKey)

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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRulesetConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFileManager_AddRulesetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		packageKey  string
		config      *RulesetConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "success - add ruleset",
			setupFile: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "empty.json")
			},
			packageKey: "registry1/ruleset1",
			config: &RulesetConfig{
				Version:  "1.0.0",
				Priority: 100,
				Sinks:    []string{"sink1"},
			},
			wantErr: false,
		},
		{
			name: "error - ruleset already exists",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    1,
					Registries: make(map[string]map[string]interface{}),
					Dependencies: Dependencies{
						Rulesets: map[string]RulesetConfig{
							"registry1/ruleset1": {
								Version: "1.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			packageKey: "registry1/ruleset1",
			config: &RulesetConfig{
				Version: "2.0.0",
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.AddRulesetConfig(ctx, tt.packageKey, tt.config)

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

func TestFileManager_RemoveRulesetConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFile   func(t *testing.T) string
		packageKey  string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - remove ruleset",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:    1,
					Registries: make(map[string]map[string]interface{}),
					Dependencies: Dependencies{
						Rulesets: map[string]RulesetConfig{
							"registry1/ruleset1": {
								Version: "1.0.0",
							},
							"registry1/ruleset2": {
								Version: "2.0.0",
							},
						},
					},
					Sinks: make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			packageKey: "registry1/ruleset1",
			wantErr:    false,
		},
		{
			name: "error - ruleset not found",
			setupFile: func(t *testing.T) string {
				manifest := &Manifest{
					Version:      1,
					Registries:   make(map[string]map[string]interface{}),
					Dependencies: Dependencies{},
					Sinks:        make(map[string]SinkConfig),
				}
				return createTestManifest(t, manifest)
			},
			packageKey:  "registry1/nonexistent",
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := tt.setupFile(t)
			fm := NewFileManagerWithPath(manifestPath)

			err := fm.RemoveRulesetConfig(ctx, tt.packageKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveRulesetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemoveRulesetConfig() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemoveRulesetConfig() error = %v, should contain %v", err, tt.errContains)
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

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}