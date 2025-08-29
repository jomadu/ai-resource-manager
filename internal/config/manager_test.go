package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const TEST_CONFIG_FILE = ".armrc.json"

func TestFileManager_GetRegistries(t *testing.T) {
	tests := []struct {
		name       string
		configData *Config
		want       map[string]RegistryConfig
		wantErr    bool
	}{
		{
			name: "valid config with registries",
			configData: &Config{
				Registries: map[string]RegistryConfig{
					"ai-rules": {
						URL:  "https://github.com/my-user/ai-rules",
						Type: "git",
					},
					"local-rules": {
						URL:  "/path/to/local/rules",
						Type: "local",
					},
				},
				Sinks: map[string]SinkConfig{},
			},
			want: map[string]RegistryConfig{
				"ai-rules": {
					URL:  "https://github.com/my-user/ai-rules",
					Type: "git",
				},
				"local-rules": {
					URL:  "/path/to/local/rules",
					Type: "local",
				},
			},
		},
		{
			name: "empty registries",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks:      map[string]SinkConfig{},
			},
			want: map[string]RegistryConfig{},
		},
		{
			name:    "missing config file",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.configData)
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			got, err := fm.GetRegistries(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegistries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !registriesEqual(got, tt.want) {
				t.Errorf("GetRegistries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_GetSinks(t *testing.T) {
	tests := []struct {
		name       string
		configData *Config
		want       map[string]SinkConfig
		wantErr    bool
	}{
		{
			name: "valid config with sinks",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks: map[string]SinkConfig{
					"q": {
						Directories: []string{".amazonq/rules"},
						Include:     []string{"ai-rules/amazonq-*"},
						Exclude:     []string{"ai-rules/cursor-*"},
					},
					"cursor": {
						Directories: []string{".cursor/rules"},
						Include:     []string{"ai-rules/cursor-*"},
						Exclude:     []string{"ai-rules/amazonq-*"},
					},
				},
			},
			want: map[string]SinkConfig{
				"q": {
					Directories: []string{".amazonq/rules"},
					Include:     []string{"ai-rules/amazonq-*"},
					Exclude:     []string{"ai-rules/cursor-*"},
				},
				"cursor": {
					Directories: []string{".cursor/rules"},
					Include:     []string{"ai-rules/cursor-*"},
					Exclude:     []string{"ai-rules/amazonq-*"},
				},
			},
		},
		{
			name: "empty sinks",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks:      map[string]SinkConfig{},
			},
			want: map[string]SinkConfig{},
		},
		{
			name:    "missing config file",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.configData)
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			got, err := fm.GetSinks(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !sinksEqual(got, tt.want) {
				t.Errorf("GetSinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_AddRegistry(t *testing.T) {
	tests := []struct {
		name         string
		configData   *Config
		registryName string
		url          string
		registryType string
		wantErr      bool
	}{
		{
			name: "add to empty config",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks:      map[string]SinkConfig{},
			},
			registryName: "ai-rules",
			url:          "https://github.com/my-user/ai-rules",
			registryType: "git",
		},
		{
			name: "add to existing config",
			configData: &Config{
				Registries: map[string]RegistryConfig{
					"existing": {
						URL:  "https://github.com/existing/rules",
						Type: "git",
					},
				},
				Sinks: map[string]SinkConfig{},
			},
			registryName: "ai-rules",
			url:          "https://github.com/my-user/ai-rules",
			registryType: "git",
		},
		{
			name: "add duplicate registry",
			configData: &Config{
				Registries: map[string]RegistryConfig{
					"ai-rules": {
						URL:  "https://github.com/existing/ai-rules",
						Type: "git",
					},
				},
				Sinks: map[string]SinkConfig{},
			},
			registryName: "ai-rules",
			url:          "https://github.com/my-user/ai-rules",
			registryType: "git",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.configData)
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			err := fm.AddRegistry(context.Background(), tt.registryName, tt.url, tt.registryType)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				verifyRegistryExists(t, tt.registryName, tt.url, tt.registryType)
			}
		})
	}
}

func TestFileManager_AddSink(t *testing.T) {
	tests := []struct {
		name       string
		configData *Config
		sinkName   string
		dirs       []string
		include    []string
		exclude    []string
		wantErr    bool
	}{
		{
			name: "add to empty config",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks:      map[string]SinkConfig{},
			},
			sinkName: "q",
			dirs:     []string{".amazonq/rules"},
			include:  []string{"ai-rules/amazonq-*"},
			exclude:  []string{"ai-rules/cursor-*"},
		},
		{
			name: "add to existing config",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks: map[string]SinkConfig{
					"existing": {
						Directories: []string{".existing/rules"},
						Include:     []string{"existing/*"},
					},
				},
			},
			sinkName: "cursor",
			dirs:     []string{".cursor/rules"},
			include:  []string{"ai-rules/cursor-*"},
			exclude:  []string{"ai-rules/amazonq-*"},
		},
		{
			name: "add duplicate sink",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks: map[string]SinkConfig{
					"q": {
						Directories: []string{".existing/rules"},
						Include:     []string{"existing/*"},
					},
				},
			},
			sinkName: "q",
			dirs:     []string{".amazonq/rules"},
			include:  []string{"ai-rules/amazonq-*"},
			exclude:  []string{"ai-rules/cursor-*"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.configData)
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			err := fm.AddSink(context.Background(), tt.sinkName, tt.dirs, tt.include, tt.exclude)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				verifySinkExists(t, tt.sinkName, tt.dirs, tt.include, tt.exclude)
			}
		})
	}
}

func TestFileManager_RemoveRegistry(t *testing.T) {
	tests := []struct {
		name         string
		configData   *Config
		registryName string
		wantErr      bool
	}{
		{
			name: "remove existing registry",
			configData: &Config{
				Registries: map[string]RegistryConfig{
					"ai-rules": {
						URL:  "https://github.com/my-user/ai-rules",
						Type: "git",
					},
					"other-rules": {
						URL:  "https://github.com/other/rules",
						Type: "git",
					},
				},
				Sinks: map[string]SinkConfig{},
			},
			registryName: "ai-rules",
		},
		{
			name: "remove last registry",
			configData: &Config{
				Registries: map[string]RegistryConfig{
					"ai-rules": {
						URL:  "https://github.com/my-user/ai-rules",
						Type: "git",
					},
				},
				Sinks: map[string]SinkConfig{},
			},
			registryName: "ai-rules",
		},
		{
			name: "remove non-existent registry",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks:      map[string]SinkConfig{},
			},
			registryName: "nonexistent",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.configData)
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			err := fm.RemoveRegistry(context.Background(), tt.registryName)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				verifyRegistryNotExists(t, tt.registryName)
			}
		})
	}
}

func TestFileManager_RemoveSink(t *testing.T) {
	tests := []struct {
		name       string
		configData *Config
		sinkName   string
		wantErr    bool
	}{
		{
			name: "remove existing sink",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks: map[string]SinkConfig{
					"q": {
						Directories: []string{".amazonq/rules"},
						Include:     []string{"ai-rules/amazonq-*"},
						Exclude:     []string{"ai-rules/cursor-*"},
					},
					"cursor": {
						Directories: []string{".cursor/rules"},
						Include:     []string{"ai-rules/cursor-*"},
						Exclude:     []string{"ai-rules/amazonq-*"},
					},
				},
			},
			sinkName: "q",
		},
		{
			name: "remove last sink",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks: map[string]SinkConfig{
					"q": {
						Directories: []string{".amazonq/rules"},
						Include:     []string{"ai-rules/amazonq-*"},
						Exclude:     []string{"ai-rules/cursor-*"},
					},
				},
			},
			sinkName: "q",
		},
		{
			name: "remove non-existent sink",
			configData: &Config{
				Registries: map[string]RegistryConfig{},
				Sinks:      map[string]SinkConfig{},
			},
			sinkName: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.configData)
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			err := fm.RemoveSink(context.Background(), tt.sinkName)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				verifySinkNotExists(t, tt.sinkName)
			}
		})
	}
}

func TestFileManager_ImplementsInterface(t *testing.T) {
	var _ Manager = (*FileManager)(nil)
}

// Helper functions

func setupTestDir(t *testing.T, configData *Config) string {
	tempDir := t.TempDir()
	if configData != nil {
		configPath := filepath.Join(tempDir, TEST_CONFIG_FILE)
		data, err := json.MarshalIndent(configData, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}
		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}
	}
	return tempDir
}

func verifyRegistryExists(t *testing.T, name, url, registryType string) {
	data, err := os.ReadFile(TEST_CONFIG_FILE)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	registry, exists := config.Registries[name]
	if !exists {
		t.Errorf("Registry %s not found", name)
		return
	}

	if registry.URL != url {
		t.Errorf("Expected URL %s, got %s", url, registry.URL)
	}
	if registry.Type != registryType {
		t.Errorf("Expected type %s, got %s", registryType, registry.Type)
	}
}

func verifyRegistryNotExists(t *testing.T, name string) {
	data, err := os.ReadFile(TEST_CONFIG_FILE)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if _, exists := config.Registries[name]; exists {
		t.Errorf("Registry %s should not exist", name)
	}
}

func verifySinkExists(t *testing.T, name string, dirs, include, exclude []string) {
	data, err := os.ReadFile(TEST_CONFIG_FILE)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	sink, exists := config.Sinks[name]
	if !exists {
		t.Errorf("Sink %s not found", name)
		return
	}

	if !stringSlicesEqual(sink.Directories, dirs) {
		t.Errorf("Expected directories %v, got %v", dirs, sink.Directories)
	}
	if !stringSlicesEqual(sink.Include, include) {
		t.Errorf("Expected include %v, got %v", include, sink.Include)
	}
	if !stringSlicesEqual(sink.Exclude, exclude) {
		t.Errorf("Expected exclude %v, got %v", exclude, sink.Exclude)
	}
}

func verifySinkNotExists(t *testing.T, name string) {
	data, err := os.ReadFile(TEST_CONFIG_FILE)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if _, exists := config.Sinks[name]; exists {
		t.Errorf("Sink %s should not exist", name)
	}
}

func registriesEqual(a, b map[string]RegistryConfig) bool {
	if len(a) != len(b) {
		return false
	}
	for name, aRegistry := range a {
		bRegistry, exists := b[name]
		if !exists || aRegistry.URL != bRegistry.URL || aRegistry.Type != bRegistry.Type {
			return false
		}
	}
	return true
}

func sinksEqual(a, b map[string]SinkConfig) bool {
	if len(a) != len(b) {
		return false
	}
	for name, aSink := range a {
		bSink, exists := b[name]
		if !exists || !stringSlicesEqual(aSink.Directories, bSink.Directories) || !stringSlicesEqual(aSink.Include, bSink.Include) || !stringSlicesEqual(aSink.Exclude, bSink.Exclude) {
			return false
		}
	}
	return true
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
