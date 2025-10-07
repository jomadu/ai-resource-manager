package resource

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
	"gopkg.in/yaml.v3"
)

func TestDefaultCompiler_Compile(t *testing.T) {
	compiler, err := NewCompiler(TargetCursor)
	if err != nil {
		t.Fatalf("Failed to create compiler: %v", err)
	}

	ruleset := &Ruleset{
		APIVersion: "v1",
		Kind:       "Ruleset",
		Metadata: Metadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: RulesetSpec{
			Rules: map[string]Rule{
				"rule1": {
					Name:        "Test Rule 1",
					Description: "First rule",
					Enforcement: "must",
					Body:        "Rule 1 content",
				},
				"rule2": {
					Name:        "Test Rule 2",
					Description: "Second rule",
					Enforcement: "should",
					Body:        "Rule 2 content",
				},
			},
		},
	}

	namespace := "ai-rules/test@1.0.0"
	// Convert ruleset to resource file format
	resourceContent, err := yaml.Marshal(ruleset)
	if err != nil {
		t.Fatalf("Failed to marshal ruleset: %v", err)
	}

	resourceFile := &types.File{
		Path:    "test-ruleset.yml",
		Content: resourceContent,
		Size:    int64(len(resourceContent)),
	}

	files, err := compiler.CompileRuleset(namespace, resourceFile)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	// Should generate one file per rule
	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}

	// Check first file
	file1 := files[0]
	if file1.Path != "test-ruleset_rule1.mdc" {
		t.Errorf("Expected filename test-ruleset_rule1.mdc, got %s", file1.Path)
	}

	content1 := string(file1.Content)
	if !strings.Contains(content1, "namespace: "+namespace) {
		t.Error("Expected namespace in first file content")
	}
	if !strings.Contains(content1, "Rule 1 content") {
		t.Error("Expected rule 1 body in first file content")
	}
	if !strings.Contains(content1, "alwaysApply: true") {
		t.Error("Expected alwaysApply: true for 'must' enforcement")
	}

	// Check second file
	file2 := files[1]
	if file2.Path != "test-ruleset_rule2.mdc" {
		t.Errorf("Expected filename test-ruleset_rule2.mdc, got %s", file2.Path)
	}

	content2 := string(file2.Content)
	if !strings.Contains(content2, "Rule 2 content") {
		t.Error("Expected rule 2 body in second file content")
	}
	if strings.Contains(content2, "alwaysApply: true") {
		t.Error("Should not have alwaysApply: true for 'should' enforcement")
	}
}

func TestNewCompiler_UnsupportedTarget(t *testing.T) {
	_, err := NewCompiler("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported target")
	}
}

func TestDefaultCompiler_AmazonQTarget(t *testing.T) {
	compiler, err := NewCompiler(TargetAmazonQ)
	if err != nil {
		t.Fatalf("Failed to create Amazon Q compiler: %v", err)
	}

	ruleset := &Ruleset{
		APIVersion: "v1",
		Kind:       "Ruleset",
		Metadata: Metadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: RulesetSpec{
			Rules: map[string]Rule{
				"rule1": {
					Name: "Test Rule",
					Body: "Rule content",
				},
			},
		},
	}

	// Convert ruleset to resource file format
	resourceContent, err := yaml.Marshal(ruleset)
	if err != nil {
		t.Fatalf("Failed to marshal ruleset: %v", err)
	}

	resourceFile := &types.File{
		Path:    "test-ruleset.yml",
		Content: resourceContent,
		Size:    int64(len(resourceContent)),
	}

	files, err := compiler.CompileRuleset("test-namespace", resourceFile)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	// Amazon Q should use .md extension
	if files[0].Path != "test-ruleset_rule1.md" {
		t.Errorf("Expected .md extension for Amazon Q, got %s", files[0].Path)
	}
}
