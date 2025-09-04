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
					"resolved": "2.1.0",
					"checksum": "sha256:abc123def456789"
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

	if entry.Checksum != "sha256:abc123def456789" {
		t.Errorf("Expected checksum 'sha256:abc123def456789', got '%s'", entry.Checksum)
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
					"resolved": "2.1.0",
					"checksum": "sha256:abc123def456789"
				},
				"cursor-rules": {
					"resolved": "2.1.0",
					"checksum": "sha256:def456abc123789"
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
	if amazonqRules.Checksum != "sha256:abc123def456789" {
		t.Errorf("Expected checksum 'sha256:abc123def456789', got '%s'", amazonqRules.Checksum)
	}
}

func TestFileManager_CreateEntry(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "arm.lock")

	manager := NewFileManagerWithPath(lockPath)
	ctx := context.Background()

	entry := &Entry{
		Resolved: "1.0.0",
		Checksum: "sha256:abc123def456789",
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

	if retrievedEntry.Checksum != entry.Checksum {
		t.Errorf("Expected checksum '%s', got '%s'", entry.Checksum, retrievedEntry.Checksum)
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
					"resolved": "1.0.0",
					"checksum": "sha256:abc123def456789"
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
		Resolved: "1.1.0",
		Checksum: "sha256:def456abc123789",
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
	if retrievedEntry.Checksum != "sha256:def456abc123789" {
		t.Errorf("Expected checksum 'sha256:def456abc123789', got '%s'", retrievedEntry.Checksum)
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
					"resolved": "2.1.0",
					"checksum": "sha256:abc123def456789"
				},
				"cursor-rules": {
					"resolved": "2.1.0",
					"checksum": "sha256:def456abc123789"
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
					"resolved": "2.1.0",
					"checksum": "sha256:abc123def456789"
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
