package urf

import (
	"strings"
	"testing"
)

func TestMarkdownRuleGenerator_GenerateRule(t *testing.T) {
	generator := &MarkdownRuleGenerator{
		metadataGen: NewRuleMetadataGenerator(),
	}

	ruleset := &Ruleset{
		Metadata: Metadata{
			ID:      "test-ruleset",
			Name:    "Test Ruleset",
			Version: "1.0.0",
		},
		Rules: map[string]Rule{
			"rule1": {
				Name:        "Test Rule 1",
				Description: "First test rule",
				Priority:    80,
				Enforcement: "should",
				Scope: []Scope{
					{Files: []string{"**/*.py"}},
				},
				Body: "This is the rule body content.",
			},
		},
	}

	rule := ruleset.Rules["rule1"]
	namespace := "ai-rules/test@1.0.0"

	result := generator.GenerateRule(namespace, ruleset, "rule1", &rule)

	// Markdown should NOT have tool-specific frontmatter
	if strings.Contains(result, "description:") && strings.Index(result, "description:") < strings.Index(result, "namespace:") {
		t.Error("Markdown should not have tool-specific frontmatter before URF metadata")
	}
	if strings.Contains(result, "globs:") && strings.Index(result, "globs:") < strings.Index(result, "namespace:") {
		t.Error("Markdown should not have globs frontmatter before URF metadata")
	}

	// Check URF metadata block
	if !strings.Contains(result, "namespace: "+namespace) {
		t.Error("Expected namespace in metadata block")
	}
	if !strings.Contains(result, "id: test-ruleset") {
		t.Error("Expected ruleset ID in metadata block")
	}
	if !strings.Contains(result, "enforcement: SHOULD") {
		t.Error("Expected uppercase enforcement in metadata block")
	}

	// Check rule content
	if !strings.Contains(result, "# Test Rule 1 (SHOULD)") {
		t.Error("Expected rule title with enforcement")
	}
	if !strings.Contains(result, "This is the rule body content.") {
		t.Error("Expected rule body content")
	}
}
