package resource

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestDefaultParser_IsRuleset(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		file     *types.File
		expected bool
	}{
		{
			name: "valid ruleset file",
			file: &types.File{
				Path: "test-ruleset.yml",
				Content: []byte(`apiVersion: v1
kind: Ruleset
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
spec:
  rules:
    rule1:
      name: "Test Rule"
      body: "Rule content"`),
			},
			expected: true,
		},
		{
			name: "valid promptset file",
			file: &types.File{
				Path: "test-promptset.yml",
				Content: []byte(`apiVersion: v1
kind: Promptset
metadata:
  id: "test-promptset"
  name: "Test Promptset"
spec:
  prompts:
    prompt1:
      name: "Test Prompt"
      body: "Prompt content"`),
			},
			expected: false,
		},
		{
			name: "invalid yaml",
			file: &types.File{
				Path:    "invalid.yml",
				Content: []byte(`invalid: yaml: content`),
			},
			expected: false,
		},
		{
			name: "missing kind",
			file: &types.File{
				Path: "missing-kind.yml",
				Content: []byte(`apiVersion: v1
metadata:
  id: "test"
spec:
  rules:
    rule1:
      name: "Test Rule"`),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsRuleset(tt.file)
			if result != tt.expected {
				t.Errorf("IsRuleset() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultParser_IsPromptset(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		file     *types.File
		expected bool
	}{
		{
			name: "valid promptset file",
			file: &types.File{
				Path: "test-promptset.yml",
				Content: []byte(`apiVersion: v1
kind: Promptset
metadata:
  id: "test-promptset"
  name: "Test Promptset"
spec:
  prompts:
    prompt1:
      name: "Test Prompt"
      body: "Prompt content"`),
			},
			expected: true,
		},
		{
			name: "valid ruleset file",
			file: &types.File{
				Path: "test-ruleset.yml",
				Content: []byte(`apiVersion: v1
kind: Ruleset
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
spec:
  rules:
    rule1:
      name: "Test Rule"
      body: "Rule content"`),
			},
			expected: false,
		},
		{
			name: "invalid yaml",
			file: &types.File{
				Path:    "invalid.yml",
				Content: []byte(`invalid: yaml: content`),
			},
			expected: false,
		},
		{
			name: "missing kind",
			file: &types.File{
				Path: "missing-kind.yml",
				Content: []byte(`apiVersion: v1
metadata:
  id: "test"
spec:
  prompts:
    prompt1:
      name: "Test Prompt"`),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsPromptset(tt.file)
			if result != tt.expected {
				t.Errorf("IsPromptset() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultParser_IsRulesetFile(t *testing.T) {
	parser := NewParser()

	tempDir := t.TempDir()

	// Create a valid ruleset file
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
spec:
  rules:
    rule1:
      name: "Test Rule"
      body: "Rule content"`

	rulesetPath := filepath.Join(tempDir, "test-ruleset.yml")
	err := os.WriteFile(rulesetPath, []byte(rulesetContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a valid promptset file
	promptsetContent := `apiVersion: v1
kind: Promptset
metadata:
  id: "test-promptset"
  name: "Test Promptset"
spec:
  prompts:
    prompt1:
      name: "Test Prompt"
      body: "Prompt content"`

	promptsetPath := filepath.Join(tempDir, "test-promptset.yml")
	err = os.WriteFile(promptsetPath, []byte(promptsetContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create an invalid file
	invalidPath := filepath.Join(tempDir, "invalid.yml")
	err = os.WriteFile(invalidPath, []byte("invalid: yaml: content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "valid ruleset file",
			path:     rulesetPath,
			expected: true,
		},
		{
			name:     "valid promptset file",
			path:     promptsetPath,
			expected: false,
		},
		{
			name:     "invalid file",
			path:     invalidPath,
			expected: false,
		},
		{
			name:     "non-existent file",
			path:     filepath.Join(tempDir, "nonexistent.yml"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsRulesetFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsRulesetFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultParser_IsPromptsetFile(t *testing.T) {
	parser := NewParser()

	tempDir := t.TempDir()

	// Create a valid ruleset file
	rulesetContent := `apiVersion: v1
kind: Ruleset
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
spec:
  rules:
    rule1:
      name: "Test Rule"
      body: "Rule content"`

	rulesetPath := filepath.Join(tempDir, "test-ruleset.yml")
	err := os.WriteFile(rulesetPath, []byte(rulesetContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a valid promptset file
	promptsetContent := `apiVersion: v1
kind: Promptset
metadata:
  id: "test-promptset"
  name: "Test Promptset"
spec:
  prompts:
    prompt1:
      name: "Test Prompt"
      body: "Prompt content"`

	promptsetPath := filepath.Join(tempDir, "test-promptset.yml")
	err = os.WriteFile(promptsetPath, []byte(promptsetContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create an invalid file
	invalidPath := filepath.Join(tempDir, "invalid.yml")
	err = os.WriteFile(invalidPath, []byte("invalid: yaml: content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "valid promptset file",
			path:     promptsetPath,
			expected: true,
		},
		{
			name:     "valid ruleset file",
			path:     rulesetPath,
			expected: false,
		},
		{
			name:     "invalid file",
			path:     invalidPath,
			expected: false,
		},
		{
			name:     "non-existent file",
			path:     filepath.Join(tempDir, "nonexistent.yml"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsPromptsetFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsPromptsetFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultParser_ParseRuleset(t *testing.T) {
	parser := NewParser()

	validRulesetFile := &types.File{
		Path: "test-ruleset.yml",
		Content: []byte(`apiVersion: v1
kind: Ruleset
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
  description: "A test ruleset"
spec:
  rules:
    rule1:
      name: "Test Rule 1"
      description: "First rule"
      enforcement: "must"
      body: "Rule 1 content"
    rule2:
      name: "Test Rule 2"
      description: "Second rule"
      enforcement: "should"
      body: "Rule 2 content"`),
	}

	ruleset, err := parser.ParseRuleset(validRulesetFile)
	if err != nil {
		t.Fatalf("ParseRuleset() error = %v", err)
	}

	if ruleset.APIVersion != "v1" {
		t.Errorf("Expected APIVersion 'v1', got '%s'", ruleset.APIVersion)
	}

	if ruleset.Kind != "Ruleset" {
		t.Errorf("Expected Kind 'Ruleset', got '%s'", ruleset.Kind)
	}

	if ruleset.Metadata.ID != "test-ruleset" {
		t.Errorf("Expected ID 'test-ruleset', got '%s'", ruleset.Metadata.ID)
	}

	if len(ruleset.Spec.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(ruleset.Spec.Rules))
	}

	rule1, exists := ruleset.Spec.Rules["rule1"]
	if !exists {
		t.Error("Expected 'rule1' to exist")
	}

	if rule1.Enforcement != "must" {
		t.Errorf("Expected rule1 enforcement 'must', got '%s'", rule1.Enforcement)
	}
}

func TestDefaultParser_ParsePromptset(t *testing.T) {
	parser := NewParser()

	validPromptsetFile := &types.File{
		Path: "test-promptset.yml",
		Content: []byte(`apiVersion: v1
kind: Promptset
metadata:
  id: "test-promptset"
  name: "Test Promptset"
  description: "A test promptset"
spec:
  prompts:
    prompt1:
      name: "Test Prompt 1"
      description: "First prompt"
      body: "Prompt 1 content"
    prompt2:
      name: "Test Prompt 2"
      description: "Second prompt"
      body: "Prompt 2 content"`),
	}

	promptset, err := parser.ParsePromptset(validPromptsetFile)
	if err != nil {
		t.Fatalf("ParsePromptset() error = %v", err)
	}

	if promptset.APIVersion != "v1" {
		t.Errorf("Expected APIVersion 'v1', got '%s'", promptset.APIVersion)
	}

	if promptset.Kind != "Promptset" {
		t.Errorf("Expected Kind 'Promptset', got '%s'", promptset.Kind)
	}

	if promptset.Metadata.ID != "test-promptset" {
		t.Errorf("Expected ID 'test-promptset', got '%s'", promptset.Metadata.ID)
	}

	if len(promptset.Spec.Prompts) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(promptset.Spec.Prompts))
	}

	prompt1, exists := promptset.Spec.Prompts["prompt1"]
	if !exists {
		t.Error("Expected 'prompt1' to exist")
	}

	if prompt1.Body != "Prompt 1 content" {
		t.Errorf("Expected prompt1 body 'Prompt 1 content', got '%s'", prompt1.Body)
	}
}
