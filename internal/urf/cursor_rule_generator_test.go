package urf

import (
	"strings"
	"testing"
)

func TestCursorRuleGenerator_GenerateRule(t *testing.T) {
	generator := &CursorRuleGenerator{
		metadataGen: NewRuleMetadataGenerator(),
	}

	ruleset := &Ruleset{
		Metadata: Metadata{
			ID:          "test-ruleset",
			Name:        "Test Ruleset",
			Description: "Test description",
		},
		Rules: map[string]Rule{
			"rule1": {
				Name:        "Test Rule 1",
				Description: "First test rule",
				Priority:    100,
				Enforcement: "must",
				Scope: []Scope{
					{Files: []string{"**/*.js", "**/*.ts"}},
				},
				Body: "This is the rule body content.",
			},
		},
	}

	rule := ruleset.Rules["rule1"]
	namespace := "ai-rules/test@1.0.0"

	result := generator.GenerateRule(namespace, ruleset, "rule1", &rule)

	// Check Cursor-specific frontmatter
	if !strings.Contains(result, `description: "First test rule"`) {
		t.Error("Expected description in frontmatter")
	}
	if !strings.Contains(result, `globs: ["**/*.js", "**/*.ts"]`) {
		t.Error("Expected globs array in frontmatter")
	}
	if !strings.Contains(result, "alwaysApply: true") {
		t.Error("Expected alwaysApply: true for 'must' enforcement")
	}

	// Check URF metadata block
	if !strings.Contains(result, "namespace: "+namespace) {
		t.Error("Expected namespace in metadata block")
	}
	if !strings.Contains(result, "id: test-ruleset") {
		t.Error("Expected ruleset ID in metadata block")
	}
	if !strings.Contains(result, "enforcement: MUST") {
		t.Error("Expected uppercase enforcement in metadata block")
	}

	// Check rule content
	if !strings.Contains(result, "# Test Rule 1 (MUST)") {
		t.Error("Expected rule title with enforcement")
	}
	if !strings.Contains(result, "This is the rule body content.") {
		t.Error("Expected rule body content")
	}
}

func TestCursorRuleGenerator_GenerateRule_ShouldEnforcement(t *testing.T) {
	generator := &CursorRuleGenerator{
		metadataGen: NewRuleMetadataGenerator(),
	}

	ruleset := &Ruleset{
		Metadata: Metadata{ID: "test", Name: "Test"},
		Rules: map[string]Rule{
			"rule1": {
				Name:        "Test Rule",
				Enforcement: "should",
				Body:        "Rule content",
			},
		},
	}

	rule := ruleset.Rules["rule1"]
	result := generator.GenerateRule("test", ruleset, "rule1", &rule)

	// Should NOT have alwaysApply for 'should' enforcement
	if strings.Contains(result, "alwaysApply: true") {
		t.Error("Should not have alwaysApply: true for 'should' enforcement")
	}
	if !strings.Contains(result, "# Test Rule (SHOULD)") {
		t.Error("Expected SHOULD enforcement in title")
	}
}
