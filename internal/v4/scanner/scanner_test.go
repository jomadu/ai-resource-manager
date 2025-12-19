package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanResources_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	config := Config{
		Dirs:      []string{tempDir},
		Recursive: false,
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 0 {
		t.Errorf("expected empty rulesets, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 0 {
		t.Errorf("expected empty promptsets, got %d", len(result.Promptsets))
	}
}

func TestScanResources_EmptyDirsArray(t *testing.T) {
	config := Config{
		Dirs: []string{},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 0 {
		t.Errorf("expected empty rulesets, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 0 {
		t.Errorf("expected empty promptsets, got %d", len(result.Promptsets))
	}
}

func TestScanResources_NonExistentDirectory(t *testing.T) {
	config := Config{
		Dirs: []string{"/does/not/exist"},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 0 {
		t.Errorf("expected empty rulesets, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 0 {
		t.Errorf("expected empty promptsets, got %d", len(result.Promptsets))
	}
}

func TestScanResources_OnlyRulesets(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create valid ruleset file
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset
  name: Test Ruleset
spec:
  rules:
    rule1:
      name: Test Rule
      body: "This is a test rule"
`
	
	rulesetPath := filepath.Join(tempDir, "ruleset.yml")
	if err := os.WriteFile(rulesetPath, []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	config := Config{
		Dirs: []string{tempDir},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 0 {
		t.Errorf("expected empty promptsets, got %d", len(result.Promptsets))
	}
	if len(result.Rulesets) > 0 && result.Rulesets[0].Metadata.ID != "test-ruleset" {
		t.Errorf("expected ruleset ID 'test-ruleset', got '%s'", result.Rulesets[0].Metadata.ID)
	}
}

func TestScanResources_OnlyPromptsets(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create valid promptset file
	promptsetContent := `apiVersion: v1
kind: Promptset
metadata:
  id: test-promptset
  name: Test Promptset
spec:
  prompts:
    prompt1:
      name: Test Prompt
      body: "This is a test prompt"
`
	
	promptsetPath := filepath.Join(tempDir, "promptset.yml")
	if err := os.WriteFile(promptsetPath, []byte(promptsetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	config := Config{
		Dirs: []string{tempDir},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 0 {
		t.Errorf("expected empty rulesets, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 1 {
		t.Errorf("expected 1 promptset, got %d", len(result.Promptsets))
	}
	if len(result.Promptsets) > 0 && result.Promptsets[0].Metadata.ID != "test-promptset" {
		t.Errorf("expected promptset ID 'test-promptset', got '%s'", result.Promptsets[0].Metadata.ID)
	}
}

func TestScanResources_BothTypes(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create ruleset
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset
spec:
  rules:
    rule1:
      body: "Test rule"
`
	
	// Create promptset
	promptsetContent := `apiVersion: v1
kind: Promptset
metadata:
  id: test-promptset
spec:
  prompts:
    prompt1:
      body: "Test prompt"
`
	
	if err := os.WriteFile(filepath.Join(tempDir, "ruleset.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write ruleset file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "promptset.yaml"), []byte(promptsetContent), 0644); err != nil {
		t.Fatalf("failed to write promptset file: %v", err)
	}
	
	config := Config{
		Dirs: []string{tempDir},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 1 {
		t.Errorf("expected 1 promptset, got %d", len(result.Promptsets))
	}
}

func TestScanResources_IgnoreNonResourceFiles(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create non-resource files
	if err := os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("not a resource"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "config.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "invalid.yml"), []byte("not: valid"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	config := Config{
		Dirs: []string{tempDir},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 0 {
		t.Errorf("expected empty rulesets, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 0 {
		t.Errorf("expected empty promptsets, got %d", len(result.Promptsets))
	}
}

func TestScanResources_Recursive(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	
	// Create ruleset in subdirectory
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: nested-ruleset
spec:
  rules:
    rule1:
      body: "Nested rule"
`
	
	if err := os.WriteFile(filepath.Join(subDir, "ruleset.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	// Test recursive=true
	config := Config{
		Dirs:      []string{tempDir},
		Recursive: true,
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(result.Rulesets))
	}
	if len(result.Rulesets) > 0 && result.Rulesets[0].Metadata.ID != "nested-ruleset" {
		t.Errorf("expected ruleset ID 'nested-ruleset', got '%s'", result.Rulesets[0].Metadata.ID)
	}
	
	// Test recursive=false
	config.Recursive = false
	result, err = ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 0 {
		t.Errorf("expected empty rulesets, got %d", len(result.Rulesets))
	}
}

func TestScanResources_MultipleDirs(t *testing.T) {
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()
	
	// Create ruleset in first dir
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: ruleset-1
spec:
  rules:
    rule1:
      body: "Rule 1"
`
	
	// Create promptset in second dir
	promptsetContent := `apiVersion: v1
kind: Promptset
metadata:
  id: promptset-2
spec:
  prompts:
    prompt1:
      body: "Prompt 2"
`
	
	if err := os.WriteFile(filepath.Join(tempDir1, "ruleset.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write ruleset file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir2, "promptset.yml"), []byte(promptsetContent), 0644); err != nil {
		t.Fatalf("failed to write promptset file: %v", err)
	}
	
	config := Config{
		Dirs: []string{tempDir1, tempDir2},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(result.Rulesets))
	}
	if len(result.Promptsets) != 1 {
		t.Errorf("expected 1 promptset, got %d", len(result.Promptsets))
	}
}

func TestScanResources_IncludePatterns(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create multiple files
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset
spec:
  rules:
    rule1:
      body: "Test rule"
`
	
	if err := os.WriteFile(filepath.Join(tempDir, "security-rules.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "style-rules.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "other.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	config := Config{
		Dirs:    []string{tempDir},
		Include: []string{"security-*.yml"},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(result.Rulesets))
	}
}

func TestScanResources_ExcludePatterns(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create multiple files
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset
spec:
  rules:
    rule1:
      body: "Test rule"
`
	
	if err := os.WriteFile(filepath.Join(tempDir, "rules.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "experimental-rules.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	config := Config{
		Dirs:    []string{tempDir},
		Exclude: []string{"experimental-*"},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset, got %d", len(result.Rulesets))
	}
}

func TestScanResources_IncludeAndExclude(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create multiple files
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: test-ruleset
spec:
  rules:
    rule1:
      body: "Test rule"
`
	
	if err := os.WriteFile(filepath.Join(tempDir, "security-rules.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "security-experimental.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "style-rules.yml"), []byte(rulesetContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	config := Config{
		Dirs:    []string{tempDir},
		Include: []string{"security-*"},
		Exclude: []string{"*experimental*"},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset (only security-rules.yml should match), got %d", len(result.Rulesets))
	}
}

func TestScanResources_SkipInvalidFiles(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create valid and invalid files
	validRuleset := `apiVersion: v1
kind: Ruleset
metadata:
  id: valid-ruleset
spec:
  rules:
    rule1:
      body: "Valid rule"
`
	
	invalidYaml := `invalid: yaml: content: [unclosed`
	
	if err := os.WriteFile(filepath.Join(tempDir, "valid.yml"), []byte(validRuleset), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "invalid.yml"), []byte(invalidYaml), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	config := Config{
		Dirs: []string{tempDir},
	}
	
	result, err := ScanResources(config)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rulesets) != 1 {
		t.Errorf("expected 1 ruleset (should skip invalid file), got %d", len(result.Rulesets))
	}
}