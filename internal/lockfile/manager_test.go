package lockfile

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileManager_GetEntry(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	// Create test lock file
	lockContent := `{
		"rulesets": {
			"ai-rules": {
				"amazonq-rules": {
					"url": "https://github.com/my-user/ai-rules",
					"type": "git",
					"constraint": "^2.1.0",
					"resolved": "2.1.0",
					"include": ["rules/amazonq/*.md"]
				}
			}
		}
	}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	entry, err := manager.GetEntry(ctx, "ai-rules", "amazonq-rules")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if entry.URL != "https://github.com/my-user/ai-rules" {
		t.Errorf("Expected URL 'https://github.com/my-user/ai-rules', got '%s'", entry.URL)
	}
	if entry.Constraint != "^2.1.0" {
		t.Errorf("Expected constraint '^2.1.0', got '%s'", entry.Constraint)
	}
	if entry.Resolved != "2.1.0" {
		t.Errorf("Expected resolved '2.1.0', got '%s'", entry.Resolved)
	}
}

func TestFileManager_GetEntry_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	lockContent := `{"rulesets": {}}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	entry, err := manager.GetEntry(ctx, "nonexistent", "ruleset")
	if err == nil {
		t.Error("Expected error for nonexistent entry")
	}
	if entry != nil {
		t.Error("Expected nil entry for nonexistent entry")
	}
}

func TestFileManager_GetEntries(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	lockContent := `{
		"rulesets": {
			"ai-rules": {
				"amazonq-rules": {
					"url": "https://github.com/my-user/ai-rules",
					"type": "git",
					"constraint": "^2.1.0",
					"resolved": "2.1.0",
					"include": ["rules/amazonq/*.md"]
				},
				"cursor-rules": {
					"url": "https://github.com/my-user/ai-rules",
					"type": "git",
					"constraint": "^2.1.0",
					"resolved": "2.1.0",
					"include": ["rules/cursor/*.mdc"]
				}
			}
		}
	}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	entries, err := manager.GetEntries(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 registry, got %d", len(entries))
	}

	aiRules, exists := entries["ai-rules"]
	if !exists {
		t.Error("Expected 'ai-rules' registry to exist")
	}

	if len(aiRules) != 2 {
		t.Errorf("Expected 2 rulesets in ai-rules, got %d", len(aiRules))
	}

	amazonqRules, exists := aiRules["amazonq-rules"]
	if !exists {
		t.Error("Expected 'amazonq-rules' ruleset to exist")
	}
	if amazonqRules.URL != "https://github.com/my-user/ai-rules" {
		t.Errorf("Expected URL 'https://github.com/my-user/ai-rules', got '%s'", amazonqRules.URL)
	}
}

func TestFileManager_CreateEntry(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	entry := &Entry{
		URL:        "https://github.com/my-user/ai-rules",
		Type:       "git",
		Constraint: "^1.0.0",
		Resolved:   "1.0.0",
		Include:    []string{"rules/amazonq/*.md"},
	}

	err := manager.CreateEntry(ctx, "ai-rules", "amazonq-rules", entry)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify entry was created
	retrievedEntry, err := manager.GetEntry(ctx, "ai-rules", "amazonq-rules")
	if err != nil {
		t.Fatalf("Expected no error retrieving entry, got %v", err)
	}

	if retrievedEntry.URL != entry.URL {
		t.Errorf("Expected URL '%s', got '%s'", entry.URL, retrievedEntry.URL)
	}
	if retrievedEntry.Constraint != entry.Constraint {
		t.Errorf("Expected constraint '%s', got '%s'", entry.Constraint, retrievedEntry.Constraint)
	}
}

func TestFileManager_UpdateEntry(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	// Create initial lock file
	lockContent := `{
		"rulesets": {
			"ai-rules": {
				"amazonq-rules": {
					"url": "https://github.com/my-user/ai-rules",
					"type": "git",
					"constraint": "^1.0.0",
					"resolved": "1.0.0",
					"include": ["rules/amazonq/*.md"]
				}
			}
		}
	}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	updatedEntry := &Entry{
		URL:        "https://github.com/my-user/ai-rules",
		Type:       "git",
		Constraint: "^1.0.0",
		Resolved:   "1.1.0",
		Include:    []string{"rules/amazonq/*.md"},
	}

	err := manager.UpdateEntry(ctx, "ai-rules", "amazonq-rules", updatedEntry)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify entry was updated
	retrievedEntry, err := manager.GetEntry(ctx, "ai-rules", "amazonq-rules")
	if err != nil {
		t.Fatalf("Expected no error retrieving entry, got %v", err)
	}

	if retrievedEntry.Resolved != "1.1.0" {
		t.Errorf("Expected resolved '1.1.0', got '%s'", retrievedEntry.Resolved)
	}
}

func TestFileManager_RemoveEntry(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	// Create initial lock file with two entries
	lockContent := `{
		"rulesets": {
			"ai-rules": {
				"amazonq-rules": {
					"url": "https://github.com/my-user/ai-rules",
					"type": "git",
					"constraint": "^2.1.0",
					"resolved": "2.1.0",
					"include": ["rules/amazonq/*.md"]
				},
				"cursor-rules": {
					"url": "https://github.com/my-user/ai-rules",
					"type": "git",
					"constraint": "^2.1.0",
					"resolved": "2.1.0",
					"include": ["rules/cursor/*.mdc"]
				}
			}
		}
	}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	err := manager.RemoveEntry(ctx, "ai-rules", "cursor-rules")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify entry was removed
	_, err = manager.GetEntry(ctx, "ai-rules", "cursor-rules")
	if err == nil {
		t.Error("Expected error for removed entry")
	}

	// Verify other entry still exists
	_, err = manager.GetEntry(ctx, "ai-rules", "amazonq-rules")
	if err != nil {
		t.Errorf("Expected amazonq-rules to still exist, got error: %v", err)
	}
}

func TestFileManager_RemoveEntry_LastInRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	// Create initial lock file with one entry
	lockContent := `{
		"rulesets": {
			"ai-rules": {
				"amazonq-rules": {
					"url": "https://github.com/my-user/ai-rules",
					"type": "git",
					"constraint": "^2.1.0",
					"resolved": "2.1.0",
					"include": ["rules/amazonq/*.md"]
				}
			}
		}
	}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	err := manager.RemoveEntry(ctx, "ai-rules", "amazonq-rules")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify registry is removed when empty
	entries, err := manager.GetEntries(ctx)
	if err != nil {
		t.Fatalf("Expected no error getting entries, got %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected empty rulesets, got %d registries", len(entries))
	}
}

func TestFileManager_ImplementsInterface(t *testing.T) {
	var _ Manager = (*FileManager)(nil)
}
