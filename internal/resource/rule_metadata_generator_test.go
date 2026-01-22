package resource

import (
	"strings"
	"testing"
)

func TestDefaultRuleMetadataGenerator_GenerateRuleMetadata(t *testing.T) {
	generator := NewRuleMetadataGenerator()

	ruleset := &Ruleset{
		APIVersion: "v1",
		Kind:       "Ruleset",
		Metadata: Metadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: RulesetSpec{
			Rules: map[string]Rule{
				"rule1": {Name: "Rule 1"},
				"rule2": {Name: "Rule 2"},
			},
		},
	}

	rule := &Rule{
		Name:        "Test Rule",
		Enforcement: "must",
		Priority:    100,
		Scope: []Scope{
			{Files: []string{"**/*.js", "**/*.ts"}},
		},
	}

	namespace := "ai-rules/test@1.0.0"
	result := generator.GenerateRuleMetadata(namespace, ruleset, "rule1", rule)

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

	// Check scope (should be in array format)
	if !strings.Contains(result, `files: ["**/*.js", "**/*.ts"]`) {
		t.Error("Expected scope files in array format")
	}

	// Should end with frontmatter delimiter and newlines
	if !strings.HasSuffix(result, "---\n\n") {
		t.Error("Expected metadata to end with '---\\n\\n'")
	}
}

func TestDefaultRuleMetadataGenerator_GenerateRuleMetadata_NoScope(t *testing.T) {
	generator := NewRuleMetadataGenerator()

	ruleset := &Ruleset{
		APIVersion: "v1",
		Kind:       "Ruleset",
		Metadata:   Metadata{ID: "test", Name: "Test"},
		Spec: RulesetSpec{
			Rules: map[string]Rule{"rule1": {}},
		},
	}

	rule := &Rule{
		Name:        "Test Rule",
		Enforcement: "should",
		Priority:    50,
		Scope:       []Scope{}, // Empty scope
	}

	result := generator.GenerateRuleMetadata("test", ruleset, "rule1", rule)

	// Should not contain scope section when no scope files
	if strings.Contains(result, "scope:") {
		t.Error("Should not contain scope section when no scope files")
	}
	if !strings.Contains(result, "enforcement: SHOULD") {
		t.Error("Expected SHOULD enforcement")
	}
}
