package manifest

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/registry"
)

const TEST_MANIFEST_FILE = "arm.json"

func TestFileManager_GetEntry(t *testing.T) {
	tests := []struct {
		name         string
		manifestData *Manifest
		registry     string
		ruleset      string
		wantEntry    *Entry
		wantErr      bool
	}{
		{
			name: "existing entry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {
						"amazonq-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/amazonq/*.md"},
						},
					},
				},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
			wantEntry: &Entry{
				Version: "^1.0.0",
				Include: []string{"rules/amazonq/*.md"},
			},
		},
		{
			name: "missing registry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{},
				Rulesets:   map[string]map[string]Entry{},
			},
			registry: "missing-registry",
			ruleset:  "some-ruleset",
			wantErr:  true,
		},
		{
			name: "missing ruleset",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {},
				},
			},
			registry: "ai-rules",
			ruleset:  "missing-ruleset",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.manifestData)
			defer func() { _ = os.RemoveAll(tempDir) }()

			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			got, err := fm.GetEntry(context.Background(), tt.registry, tt.ruleset)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !entriesEqual(got, tt.wantEntry) {
				t.Errorf("GetEntry() = %v, want %v", got, tt.wantEntry)
			}
		})
	}
}

func TestFileManager_GetEntries(t *testing.T) {
	tests := []struct {
		name         string
		manifestData *Manifest
		want         map[string]map[string]Entry
		wantErr      bool
	}{
		{
			name: "valid manifest",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {
						"amazonq-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/amazonq/*.md"},
						},
						"cursor-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/cursor/*.mdc"},
						},
					},
				},
			},
			want: map[string]map[string]Entry{
				"ai-rules": {
					"amazonq-rules": {
						Version: "^1.0.0",
						Include: []string{"rules/amazonq/*.md"},
					},
					"cursor-rules": {
						Version: "^1.0.0",
						Include: []string{"rules/cursor/*.mdc"},
					},
				},
			},
		},
		{
			name: "empty manifest",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{},
				Rulesets:   map[string]map[string]Entry{},
			},
			want: map[string]map[string]Entry{},
		},
		{
			name:    "missing manifest file",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			if tt.manifestData != nil {
				writeManifest(t, tt.manifestData)
			}

			fm := NewFileManager()
			got, err := fm.GetEntries(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !entriesMapEqual(got, tt.want) {
				t.Errorf("GetEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_CreateEntry(t *testing.T) {
	tests := []struct {
		name         string
		manifestData *Manifest
		registry     string
		ruleset      string
		entry        Entry
		wantErr      bool
	}{
		{
			name: "create in empty manifest",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{},
				Rulesets:   map[string]map[string]Entry{},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
			entry: Entry{
				Version: "^1.0.0",
				Include: []string{"rules/amazonq/*.md"},
			},
		},
		{
			name: "create in existing registry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {
						"cursor-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/cursor/*.mdc"},
						},
					},
				},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
			entry: Entry{
				Version: "^1.0.0",
				Include: []string{"rules/amazonq/*.md"},
			},
		},
		{
			name: "create duplicate entry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {
						"amazonq-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/amazonq/*.md"},
						},
					},
				},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
			entry: Entry{
				Version: "^2.0.0",
				Include: []string{"rules/amazonq/*.md"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.manifestData)
			defer func() { _ = os.RemoveAll(tempDir) }()

			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			err := fm.CreateEntry(context.Background(), tt.registry, tt.ruleset, tt.entry)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				verifyEntryExists(t, tt.registry, tt.ruleset, tt.entry)
			}
		})
	}
}

func TestFileManager_UpdateEntry(t *testing.T) {
	tests := []struct {
		name         string
		manifestData *Manifest
		registry     string
		ruleset      string
		entry        Entry
		wantErr      bool
	}{
		{
			name: "update existing entry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {
						"amazonq-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/amazonq/*.md"},
						},
					},
				},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
			entry: Entry{
				Version: "^2.0.0",
				Include: []string{"rules/amazonq/*.md"},
				Exclude: []string{"rules/amazonq/deprecated.md"},
			},
		},
		{
			name: "update non-existent entry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{},
				Rulesets:   map[string]map[string]Entry{},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
			entry: Entry{
				Version: "^1.0.0",
				Include: []string{"rules/amazonq/*.md"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.manifestData)
			defer func() { _ = os.RemoveAll(tempDir) }()

			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			err := fm.UpdateEntry(context.Background(), tt.registry, tt.ruleset, tt.entry)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				verifyEntryExists(t, tt.registry, tt.ruleset, tt.entry)
			}
		})
	}
}

func TestFileManager_RemoveEntry(t *testing.T) {
	tests := []struct {
		name         string
		manifestData *Manifest
		registry     string
		ruleset      string
		wantErr      bool
	}{
		{
			name: "remove existing entry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {
						"amazonq-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/amazonq/*.md"},
						},
						"cursor-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/cursor/*.mdc"},
						},
					},
				},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
		},
		{
			name: "remove last entry in registry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{
					"ai-rules": {
						"amazonq-rules": {
							Version: "^1.0.0",
							Include: []string{"rules/amazonq/*.md"},
						},
					},
				},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
		},
		{
			name: "remove non-existent entry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{},
				Rulesets:   map[string]map[string]Entry{},
			},
			registry: "ai-rules",
			ruleset:  "amazonq-rules",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.manifestData)
			defer func() { _ = os.RemoveAll(tempDir) }()

			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			err := fm.RemoveEntry(context.Background(), tt.registry, tt.ruleset)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				verifyEntryNotExists(t, tt.registry, tt.ruleset)
			}
		})
	}
}

func setupTestDir(t *testing.T, manifestData *Manifest) string {
	tempDir := t.TempDir()
	if manifestData != nil {
		manifestPath := filepath.Join(tempDir, TEST_MANIFEST_FILE)
		data, err := json.MarshalIndent(manifestData, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal manifest: %v", err)
		}
		if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
			t.Fatalf("Failed to write manifest: %v", err)
		}
	}
	return tempDir
}

func writeManifest(t *testing.T, manifest *Manifest) {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal manifest: %v", err)
	}
	if err := os.WriteFile(TEST_MANIFEST_FILE, data, 0o644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}
}

func verifyEntryExists(t *testing.T, registry, ruleset string, expected Entry) {
	data, err := os.ReadFile(TEST_MANIFEST_FILE)
	if err != nil {
		t.Fatalf("Failed to read manifest: %v", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	registryMap, exists := manifest.Rulesets[registry]
	if !exists {
		t.Errorf("Registry %s not found", registry)
		return
	}

	entry, exists := registryMap[ruleset]
	if !exists {
		t.Errorf("Ruleset %s not found in registry %s", ruleset, registry)
		return
	}

	if !entriesEqual(&entry, &expected) {
		t.Errorf("Entry mismatch: got %v, want %v", entry, expected)
	}
}

func verifyEntryNotExists(t *testing.T, registry, ruleset string) {
	data, err := os.ReadFile(TEST_MANIFEST_FILE)
	if err != nil {
		t.Fatalf("Failed to read manifest: %v", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	if registryMap, exists := manifest.Rulesets[registry]; exists {
		if _, exists := registryMap[ruleset]; exists {
			t.Errorf("Entry %s/%s should not exist", registry, ruleset)
		}
	}
}

func entriesEqual(a, b *Entry) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Version == b.Version &&
		stringSlicesEqual(a.Include, b.Include) &&
		stringSlicesEqual(a.Exclude, b.Exclude)
}

func entriesMapEqual(a, b map[string]map[string]Entry) bool {
	if len(a) != len(b) {
		return false
	}
	for registry, aRulesets := range a {
		bRulesets, exists := b[registry]
		if !exists || len(aRulesets) != len(bRulesets) {
			return false
		}
		for ruleset, aEntry := range aRulesets {
			bEntry, exists := bRulesets[ruleset]
			if !exists || !entriesEqual(&aEntry, &bEntry) {
				return false
			}
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

func TestFileManager_AddRegistry(t *testing.T) {
	tests := []struct {
		name         string
		manifestData *Manifest
		registryName string
		url          string
		registryType string
		wantErr      bool
	}{
		{
			name: "add to empty manifest",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{},
				Rulesets:   map[string]map[string]Entry{},
			},
			registryName: "ai-rules",
			url:          "https://github.com/test/repo",
			registryType: "git",
		},
		{
			name: "add duplicate registry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/existing/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{},
			},
			registryName: "ai-rules",
			url:          "https://github.com/test/repo",
			registryType: "git",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.manifestData)
			defer func() { _ = os.RemoveAll(tempDir) }()

			oldWd, _ := os.Getwd()
			_ = os.Chdir(tempDir)
			defer func() { _ = os.Chdir(oldWd) }()

			fm := NewFileManager()
			gitConfig := registry.GitRegistryConfig{
				RegistryConfig: registry.RegistryConfig{URL: tt.url, Type: tt.registryType},
			}
			err := fm.AddGitRegistry(context.Background(), tt.registryName, gitConfig)

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

func TestFileManager_RemoveRegistry(t *testing.T) {
	tests := []struct {
		name         string
		manifestData *Manifest
		registryName string
		wantErr      bool
	}{
		{
			name: "remove existing registry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{
					"ai-rules": {"url": "https://github.com/test/repo", "type": "git"},
				},
				Rulesets: map[string]map[string]Entry{},
			},
			registryName: "ai-rules",
		},
		{
			name: "remove non-existent registry",
			manifestData: &Manifest{
				Registries: map[string]map[string]interface{}{},
				Rulesets:   map[string]map[string]Entry{},
			},
			registryName: "nonexistent",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := setupTestDir(t, tt.manifestData)
			defer func() { _ = os.RemoveAll(tempDir) }()

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

func verifyRegistryExists(t *testing.T, name, url, registryType string) {
	data, err := os.ReadFile(TEST_MANIFEST_FILE)
	if err != nil {
		t.Fatalf("Failed to read manifest: %v", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	rawRegistry, exists := manifest.Registries[name]
	if !exists {
		t.Errorf("Registry %s not found", name)
		return
	}

	if rawRegistry["url"] != url {
		t.Errorf("Expected URL %s, got %s", url, rawRegistry["url"])
	}
	if rawRegistry["type"] != registryType {
		t.Errorf("Expected type %s, got %s", registryType, rawRegistry["type"])
	}
}

func verifyRegistryNotExists(t *testing.T, name string) {
	data, err := os.ReadFile(TEST_MANIFEST_FILE)
	if err != nil {
		t.Fatalf("Failed to read manifest: %v", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("Failed to unmarshal manifest: %v", err)
	}

	if _, exists := manifest.Registries[name]; exists {
		t.Errorf("Registry %s should not exist", name)
	}
}
