package urf

import (
	"strings"
	"testing"
)

func TestDefaultRuleMetadataGenerator_GenerateRuleMetadata(t *testing.T) {
	generator := NewRuleMetadataGenerator()

	ruleset := &Ruleset{
		Metadata: Metadata{
			ID:      "test-ruleset",
			Name:    "Test Ruleset",
			Version: "1.0.0",
		},
		Rules: []Rule{
			{ID: "rule1", Name: "Rule 1"},
			{ID: "rule2", Name: "Rule 2"},
		},
	}

	rule := &Rule{
		ID:          "rule1",
		Name:        "Test Rule",
		Enforcement: "must",
		Priority:    100,
		Scope: []Scope{
			{Files: []string{"**/*.js", "**/*.ts"}},
		},
	}

	namespace := "ai-rules/test@1.0.0"
	result := generator.GenerateRuleMetadata(namespace, ruleset, rule)

	// Check metadata structure
	if !strings.Contains(result, "---\n") {
		t.Error("Expected YAML frontmatter delimiters")
	}
	if !strings.Contains(result, "namespace: "+namespace) {
		t.Error("Expected namespace in metadata")
	}
	if !strings.Contains(result, "id: test-ruleset") {
		t.Error("Expected ruleset ID in metadata")
	}
	if !strings.Contains(result, "name: Test Ruleset") {
		t.Error("Expected ruleset name in metadata")
	}
	if !strings.Contains(result, "version: 1.0.0") {
		t.Error("Expected ruleset version in metadata")
	}

	// Check rules list
	if !strings.Contains(result, "- rule1") {
		t.Error("Expected rule1 in rules list")
	}
	if !strings.Contains(result, "- rule2") {
		t.Error("Expected rule2 in rules list")
	}

	// Check rule details
	if !strings.Contains(result, "id: rule1") {
		t.Error("Expected rule ID in metadata")
	}
	if !strings.Contains(result, "name: Test Rule") {
		t.Error("Expected rule name in metadata")
	}
	if !strings.Contains(result, "enforcement: MUST") {
		t.Error("Expected uppercase enforcement in metadata")
	}
	if !strings.Contains(result, "priority: 100") {
		t.Error("Expected priority in metadata")
	}

	// Check scope
	if !strings.Contains(result, `files: "**/*.js"`) {
		t.Error("Expected first scope file in metadata")
	}
	if !strings.Contains(result, `files: "**/*.ts"`) {
		t.Error("Expected second scope file in metadata")
	}

	// Should end with frontmatter delimiter and newlines
	if !strings.HasSuffix(result, "---\n\n") {
		t.Error("Expected metadata to end with '---\\n\\n'")
	}
}

func TestDefaultRuleMetadataGenerator_GenerateRuleMetadata_NoScope(t *testing.T) {
	generator := NewRuleMetadataGenerator()

	ruleset := &Ruleset{
		Metadata: Metadata{ID: "test", Name: "Test", Version: "1.0.0"},
		Rules:    []Rule{{ID: "rule1"}},
	}

	rule := &Rule{
		ID:          "rule1",
		Name:        "Test Rule",
		Enforcement: "should",
		Priority:    50,
		Scope:       []Scope{}, // Empty scope
	}

	result := generator.GenerateRuleMetadata("test", ruleset, rule)

	// Should not contain scope section when no scope files
	if strings.Contains(result, "scope:") {
		t.Error("Should not contain scope section when no scope files")
	}
	if !strings.Contains(result, "enforcement: SHOULD") {
		t.Error("Expected SHOULD enforcement")
	}
}
