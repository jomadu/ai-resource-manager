package manifest

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
)

// Test helper functions

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

func newTestManifest() *Manifest {
	return &Manifest{
		Version:      1,
		Registries:   make(map[string]map[string]interface{}),
		Sinks:        make(map[string]SinkConfig),
		Dependencies: make(map[string]map[string]interface{}),
	}
}

// Registry tests

func TestFileManager_GetAllRegistriesConfig(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Registries["test-reg"] = map[string]interface{}{
		"type": "git",
		"url":  "https://github.com/test/repo",
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	registries, err := fm.GetAllRegistriesConfig(ctx)
	if err != nil {
		t.Fatalf("GetAllRegistriesConfig() error = %v", err)
	}
	
	if len(registries) != 1 {
		t.Errorf("Expected 1 registry, got %d", len(registries))
	}
	
	if registries["test-reg"]["type"] != "git" {
		t.Errorf("Expected git registry type")
	}
}




func TestFileManager_RemoveRegistryConfig(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Registries["test-reg"] = map[string]interface{}{
		"type": "git",
		"url":  "https://github.com/test/repo",
	}
	manifest.Dependencies["test-reg/package1"] = map[string]interface{}{
		"type":    "ruleset",
		"version": "1.0.0",
		"sinks":   []string{"sink1"},
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	err := fm.RemoveRegistryConfig(ctx, "test-reg")
	if err != nil {
		t.Fatalf("RemoveRegistryConfig() error = %v", err)
	}
	
	// Verify registry and dependencies are removed
	registries, _ := fm.GetAllRegistriesConfig(ctx)
	if len(registries) != 0 {
		t.Errorf("Expected 0 registries after removal, got %d", len(registries))
	}
	
	deps, _ := fm.GetAllDependenciesConfig(ctx)
	if len(deps) != 0 {
		t.Errorf("Expected 0 dependencies after registry removal, got %d", len(deps))
	}
}

// Sink tests

func TestFileManager_UpsertSinkConfig(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "empty.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	config := SinkConfig{
		Directory: ".cursor/rules",
		Tool:      compiler.Cursor,
	}
	
	err := fm.UpsertSinkConfig(ctx, "cursor-rules", config)
	if err != nil {
		t.Fatalf("UpsertSinkConfig() error = %v", err)
	}
	
	// Verify it was added
	retrieved, err := fm.GetSinkConfig(ctx, "cursor-rules")
	if err != nil {
		t.Fatalf("GetSinkConfig() error = %v", err)
	}
	
	if retrieved.Directory != config.Directory {
		t.Errorf("Expected directory %s, got %s", config.Directory, retrieved.Directory)
	}
	
	// Test upsert (update existing)
	config.Directory = ".cursor/updated"
	config.Tool = compiler.AmazonQ
	
	err = fm.UpsertSinkConfig(ctx, "cursor-rules", config)
	if err != nil {
		t.Fatalf("UpsertSinkConfig() update error = %v", err)
	}
	
	retrieved, _ = fm.GetSinkConfig(ctx, "cursor-rules")
	if retrieved.Directory != ".cursor/updated" {
		t.Errorf("Expected updated directory")
	}
	
	if retrieved.Tool != compiler.AmazonQ {
		t.Errorf("Expected updated tool")
	}
}

func TestFileManager_RemoveSinkConfig(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Sinks["test-sink"] = SinkConfig{
		Directory: ".cursor/rules",
		Tool:      compiler.Cursor,
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	err := fm.RemoveSinkConfig(ctx, "test-sink")
	if err != nil {
		t.Fatalf("RemoveSinkConfig() error = %v", err)
	}
	
	// Verify it was removed
	_, err = fm.GetSinkConfig(ctx, "test-sink")
	if err == nil {
		t.Errorf("Expected error when getting removed sink")
	}
}

// Dependency tests

func TestFileManager_UpsertRulesetDependencyConfig(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "empty.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	config := &RulesetDependencyConfig{
		BaseDependencyConfig: BaseDependencyConfig{
			Version: "1.0.0",
			Sinks:   []string{"cursor-rules"},
		},
		Priority: 100,
	}
	
	err := fm.UpsertRulesetDependencyConfig(ctx, "test-reg", "ruleset1", *config)
	if err != nil {
		t.Fatalf("UpsertRulesetDependencyConfig() error = %v", err)
	}
	
	// Verify it was added
	retrieved, err := fm.GetRulesetDependencyConfig(ctx, "test-reg", "ruleset1")
	if err != nil {
		t.Fatalf("GetRulesetDependencyConfig() error = %v", err)
	}
	
	if retrieved.Version != config.Version {
		t.Errorf("Expected version %s, got %s", config.Version, retrieved.Version)
	}
	
	if retrieved.Priority != config.Priority {
		t.Errorf("Expected priority %d, got %d", config.Priority, retrieved.Priority)
	}
	
	if retrieved.Type != ResourceTypeRuleset {
		t.Errorf("Expected type %s, got %s", ResourceTypeRuleset, retrieved.Type)
	}
}

func TestFileManager_UpsertPromptsetDependencyConfig(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "empty.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	config := &PromptsetDependencyConfig{
		BaseDependencyConfig: BaseDependencyConfig{
			Version: "2.0.0",
			Sinks:   []string{"cursor-commands"},
			Include: []string{"*.yml"},
		},
	}
	
	err := fm.UpsertPromptsetDependencyConfig(ctx, "test-reg", "promptset1", *config)
	if err != nil {
		t.Fatalf("UpsertPromptsetDependencyConfig() error = %v", err)
	}
	
	// Verify it was added
	retrieved, err := fm.GetPromptsetDependencyConfig(ctx, "test-reg", "promptset1")
	if err != nil {
		t.Fatalf("GetPromptsetDependencyConfig() error = %v", err)
	}
	
	if retrieved.Version != config.Version {
		t.Errorf("Expected version %s, got %s", config.Version, retrieved.Version)
	}
	
	if retrieved.Type != ResourceTypePromptset {
		t.Errorf("Expected type %s, got %s", ResourceTypePromptset, retrieved.Type)
	}
}

func TestFileManager_RemoveDependencyConfig(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Dependencies["test-reg/package1"] = map[string]interface{}{
		"type":    "ruleset",
		"version": "1.0.0",
		"sinks":   []string{"sink1"},
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	err := fm.RemoveDependencyConfig(ctx, "test-reg", "package1")
	if err != nil {
		t.Fatalf("RemoveDependencyConfig() error = %v", err)
	}
	
	// Verify it was removed
	_, err = fm.GetDependencyConfig(ctx, "test-reg", "package1")
	if err == nil {
		t.Errorf("Expected error when getting removed dependency")
	}
}

// Error cases

func TestFileManager_GetRegistryConfig_NotFound(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "empty.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	_, err := fm.GetRegistryConfig(ctx, "nonexistent")
	if err == nil {
		t.Errorf("Expected error for nonexistent registry")
	}
}

func TestFileManager_GetRulesetDependencyConfig_WrongType(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Dependencies["test-reg/package1"] = map[string]interface{}{
		"type":    "promptset", // Wrong type
		"version": "1.0.0",
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	_, err := fm.GetRulesetDependencyConfig(ctx, "test-reg", "package1")
	if err == nil {
		t.Errorf("Expected error for wrong dependency type")
	}
}

func TestFileManager_UpdateRegistryConfigName_WithDependencies(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Registries["old-reg"] = map[string]interface{}{
		"type": "git",
		"url":  "https://github.com/test/repo",
	}
	manifest.Dependencies["old-reg/package1"] = map[string]interface{}{
		"type":    "ruleset",
		"version": "1.0.0",
	}
	manifest.Dependencies["old-reg/package2"] = map[string]interface{}{
		"type":    "promptset",
		"version": "2.0.0",
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	err := fm.UpdateRegistryConfigName(ctx, "old-reg", "new-reg")
	if err != nil {
		t.Fatalf("UpdateRegistryConfigName() error = %v", err)
	}
	
	// Verify registry was renamed
	_, err = fm.GetRegistryConfig(ctx, "new-reg")
	if err != nil {
		t.Errorf("Expected new registry to exist")
	}
	
	// Verify old registry is gone
	_, err = fm.GetRegistryConfig(ctx, "old-reg")
	if err == nil {
		t.Errorf("Expected old registry to be removed")
	}
	
	// Verify dependencies were updated
	deps, _ := fm.GetAllDependenciesConfig(ctx)
	if _, exists := deps["new-reg/package1"]; !exists {
		t.Errorf("Expected dependency to be renamed to new-reg/package1")
	}
	if _, exists := deps["new-reg/package2"]; !exists {
		t.Errorf("Expected dependency to be renamed to new-reg/package2")
	}
	if _, exists := deps["old-reg/package1"]; exists {
		t.Errorf("Expected old dependency key to be removed")
	}
}

func TestFileManager_LoadManifest_FileNotExists(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "nonexistent.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	// Should return empty manifest, not error
	registries, err := fm.GetAllRegistriesConfig(ctx)
	if err != nil {
		t.Fatalf("Expected no error for nonexistent file, got %v", err)
	}
	
	if len(registries) != 0 {
		t.Errorf("Expected empty registries map, got %d entries", len(registries))
	}
}

func TestFileManager_UpsertDependency_Overwrite(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Dependencies["test-reg/package1"] = map[string]interface{}{
		"type":     "ruleset",
		"version":  "1.0.0",
		"priority": 50,
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	// Update with new config
	newConfig := &RulesetDependencyConfig{
		BaseDependencyConfig: BaseDependencyConfig{
			Version: "2.0.0",
			Sinks:   []string{"new-sink"},
		},
		Priority: 200,
	}
	
	err := fm.UpsertRulesetDependencyConfig(ctx, "test-reg", "package1", *newConfig)
	if err != nil {
		t.Fatalf("UpsertRulesetDependencyConfig() error = %v", err)
	}
	
	// Verify it was updated
	retrieved, err := fm.GetRulesetDependencyConfig(ctx, "test-reg", "package1")
	if err != nil {
		t.Fatalf("GetRulesetDependencyConfig() error = %v", err)
	}
	
	if retrieved.Version != "2.0.0" {
		t.Errorf("Expected version 2.0.0, got %s", retrieved.Version)
	}
	
	if retrieved.Priority != 200 {
		t.Errorf("Expected priority 200, got %d", retrieved.Priority)
	}
}

// Registry type validation tests




func TestFileManager_GetGitRegistryConfig(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Registries["git-reg"] = map[string]interface{}{
		"type":     "git",
		"url":      "https://github.com/test/repo",
		"branches": []interface{}{"main", "develop"},
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	config, err := fm.GetGitRegistryConfig(ctx, "git-reg")
	if err != nil {
		t.Fatalf("GetGitRegistryConfig() error = %v", err)
	}
	
	if config.Type != "git" {
		t.Errorf("Expected type git, got %s", config.Type)
	}
	if config.URL != "https://github.com/test/repo" {
		t.Errorf("Expected URL https://github.com/test/repo, got %s", config.URL)
	}
	if len(config.Branches) != 2 {
		t.Errorf("Expected 2 branches, got %d", len(config.Branches))
	}
}

func TestFileManager_GetGitRegistryConfig_WrongType(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Registries["gitlab-reg"] = map[string]interface{}{
		"type": "gitlab",
		"url":  "https://gitlab.com/test/repo",
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	_, err := fm.GetGitRegistryConfig(ctx, "gitlab-reg")
	if err == nil {
		t.Errorf("Expected error for wrong registry type")
	}
}

func TestFileManager_GetGitLabRegistryConfig(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Registries["gitlab-reg"] = map[string]interface{}{
		"type":       "gitlab",
		"url":        "https://gitlab.com",
		"projectId":  "123",
		"groupId":    "456",
		"apiVersion": "v4",
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	config, err := fm.GetGitLabRegistryConfig(ctx, "gitlab-reg")
	if err != nil {
		t.Fatalf("GetGitLabRegistryConfig() error = %v", err)
	}
	
	if config.Type != "gitlab" {
		t.Errorf("Expected type gitlab, got %s", config.Type)
	}
	if config.ProjectID != "123" {
		t.Errorf("Expected projectId 123, got %s", config.ProjectID)
	}
	if config.GroupID != "456" {
		t.Errorf("Expected groupId 456, got %s", config.GroupID)
	}
	if config.APIVersion != "v4" {
		t.Errorf("Expected apiVersion v4, got %s", config.APIVersion)
	}
}

func TestFileManager_GetCloudsmithRegistryConfig(t *testing.T) {
	ctx := context.Background()
	
	manifest := newTestManifest()
	manifest.Registries["cloudsmith-reg"] = map[string]interface{}{
		"type":       "cloudsmith",
		"url":        "https://cloudsmith.io",
		"owner":      "myorg",
		"repository": "myrepo",
	}
	
	manifestPath := createTestManifest(t, manifest)
	fm := NewFileManagerWithPath(manifestPath)
	
	config, err := fm.GetCloudsmithRegistryConfig(ctx, "cloudsmith-reg")
	if err != nil {
		t.Fatalf("GetCloudsmithRegistryConfig() error = %v", err)
	}
	
	if config.Type != "cloudsmith" {
		t.Errorf("Expected type cloudsmith, got %s", config.Type)
	}
	if config.Owner != "myorg" {
		t.Errorf("Expected owner myorg, got %s", config.Owner)
	}
	if config.Repository != "myrepo" {
		t.Errorf("Expected repository myrepo, got %s", config.Repository)
	}
}

func TestFileManager_UpsertGitRegistryConfig(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "empty.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	config := GitRegistryConfig{
		URL:      "https://github.com/test/repo",
		Branches: []string{"main", "develop"},
	}
	
	err := fm.UpsertGitRegistryConfig(ctx, "git-reg", config)
	if err != nil {
		t.Fatalf("UpsertGitRegistryConfig() error = %v", err)
	}
	
	retrieved, err := fm.GetGitRegistryConfig(ctx, "git-reg")
	if err != nil {
		t.Fatalf("GetGitRegistryConfig() error = %v", err)
	}
	
	if retrieved.Type != "git" {
		t.Errorf("Expected type git, got %s", retrieved.Type)
	}
	if retrieved.URL != config.URL {
		t.Errorf("Expected URL %s, got %s", config.URL, retrieved.URL)
	}
}

func TestFileManager_UpsertGitLabRegistryConfig(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "empty.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	config := GitLabRegistryConfig{
		URL:        "https://gitlab.com",
		ProjectID:  "123",
		GroupID:    "456",
		APIVersion: "v4",
	}
	
	err := fm.UpsertGitLabRegistryConfig(ctx, "gitlab-reg", config)
	if err != nil {
		t.Fatalf("UpsertGitLabRegistryConfig() error = %v", err)
	}
	
	retrieved, err := fm.GetGitLabRegistryConfig(ctx, "gitlab-reg")
	if err != nil {
		t.Fatalf("GetGitLabRegistryConfig() error = %v", err)
	}
	
	if retrieved.Type != "gitlab" {
		t.Errorf("Expected type gitlab, got %s", retrieved.Type)
	}
	if retrieved.ProjectID != config.ProjectID {
		t.Errorf("Expected projectId %s, got %s", config.ProjectID, retrieved.ProjectID)
	}
}

func TestFileManager_UpsertCloudsmithRegistryConfig(t *testing.T) {
	ctx := context.Background()
	manifestPath := filepath.Join(t.TempDir(), "empty.json")
	fm := NewFileManagerWithPath(manifestPath)
	
	config := CloudsmithRegistryConfig{
		URL:        "https://cloudsmith.io",
		Owner:      "myorg",
		Repository: "myrepo",
	}
	
	err := fm.UpsertCloudsmithRegistryConfig(ctx, "cloudsmith-reg", config)
	if err != nil {
		t.Fatalf("UpsertCloudsmithRegistryConfig() error = %v", err)
	}
	
	retrieved, err := fm.GetCloudsmithRegistryConfig(ctx, "cloudsmith-reg")
	if err != nil {
		t.Fatalf("GetCloudsmithRegistryConfig() error = %v", err)
	}
	
	if retrieved.Type != "cloudsmith" {
		t.Errorf("Expected type cloudsmith, got %s", retrieved.Type)
	}
	if retrieved.Owner != config.Owner {
		t.Errorf("Expected owner %s, got %s", config.Owner, retrieved.Owner)
	}
}
