package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const TEST_CONFIG_FILE = ".armrc.json"

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
				Sinks: map[string]SinkConfig{},
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
				Sinks: map[string]SinkConfig{},
			},
			sinkName: "q",
			dirs:     []string{".amazonq/rules"},
			include:  []string{"ai-rules/amazonq-*"},
			exclude:  []string{"ai-rules/cursor-*"},
		},
		{
			name: "add to existing config",
			configData: &Config{
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
			err := fm.AddSink(context.Background(), tt.sinkName, tt.dirs, tt.include, tt.exclude, "hierarchical", false)

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
				Sinks: map[string]SinkConfig{},
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
