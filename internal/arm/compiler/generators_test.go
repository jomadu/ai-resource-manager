package compiler

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/resource"
)

func TestGenerateRuleMetadata_Basic(t *testing.T) {
	ruleset := &resource.RulesetResource{
		Metadata: resource.ResourceMetadata{
			ID:   "clean-code",
			Name: "Clean Code Rules",
		},
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"meaningful-names": {
					Name:        "Use Meaningful Names",
					Enforcement: "must",
				},
			},
		},
	}

	rule := &resource.Rule{
		Name:        "Use Meaningful Names",
		Enforcement: "must",
	}

	metadata := GenerateRuleMetadata("sample-registry", ruleset, "meaningful-names", rule)

	// Check basic structure
	if !strings.Contains(metadata, "---") {
		t.Error("Expected metadata to contain YAML frontmatter markers")
	}

	// Check namespace
	if !strings.Contains(metadata, "namespace: sample-registry") {
		t.Error("Expected namespace in metadata")
	}

	// Check ruleset info
	if !strings.Contains(metadata, "id: clean-code") {
		t.Error("Expected ruleset ID in metadata")
	}
	if !strings.Contains(metadata, "name: Clean Code Rules") {
		t.Error("Expected ruleset name in metadata")
	}

	// Check rule info
	if !strings.Contains(metadata, "id: meaningful-names") {
		t.Error("Expected rule ID in metadata")
	}
	if !strings.Contains(metadata, "name: Use Meaningful Names") {
		t.Error("Expected rule name in metadata")
	}
	if !strings.Contains(metadata, "enforcement: MUST") {
		t.Error("Expected enforcement in uppercase")
	}
}

func TestGenerateRuleMetadata_WithPriority(t *testing.T) {
	ruleset := &resource.RulesetResource{
		Metadata: resource.ResourceMetadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"rule1": {
					Name:        "Test Rule",
					Enforcement: "should",
					Priority:    100,
				},
			},
		},
	}

	rule := &resource.Rule{
		Name:        "Test Rule",
		Enforcement: "should",
		Priority:    100,
	}

	metadata := GenerateRuleMetadata("test-namespace", ruleset, "rule1", rule)

	if !strings.Contains(metadata, "priority: 100") {
		t.Error("Expected priority in metadata")
	}
}

func TestGenerateRuleMetadata_WithoutPriority(t *testing.T) {
	ruleset := &resource.RulesetResource{
		Metadata: resource.ResourceMetadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"rule1": {
					Name:        "Test Rule",
					Enforcement: "should",
					Priority:    0,
				},
			},
		},
	}

	rule := &resource.Rule{
		Name:        "Test Rule",
		Enforcement: "should",
		Priority:    0,
	}

	metadata := GenerateRuleMetadata("test-namespace", ruleset, "rule1", rule)

	if strings.Contains(metadata, "priority:") {
		t.Error("Expected no priority in metadata when priority is 0")
	}
}

func TestGenerateRuleMetadata_WithScope(t *testing.T) {
	ruleset := &resource.RulesetResource{
		Metadata: resource.ResourceMetadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"rule1": {
					Name:        "Test Rule",
					Enforcement: "must",
					Scope: []resource.Scope{
						{
							Files: []string{"**/*.js", "**/*.ts"},
						},
					},
				},
			},
		},
	}

	rule := &resource.Rule{
		Name:        "Test Rule",
		Enforcement: "must",
		Scope: []resource.Scope{
			{
				Files: []string{"**/*.js", "**/*.ts"},
			},
		},
	}

	metadata := GenerateRuleMetadata("test-namespace", ruleset, "rule1", rule)

	if !strings.Contains(metadata, "scope:") {
		t.Error("Expected scope in metadata")
	}
	if !strings.Contains(metadata, `"**/*.js"`) {
		t.Error("Expected first file pattern in metadata")
	}
	if !strings.Contains(metadata, `"**/*.ts"`) {
		t.Error("Expected second file pattern in metadata")
	}
}

func TestGenerateRuleMetadata_MultipleRules(t *testing.T) {
	ruleset := &resource.RulesetResource{
		Metadata: resource.ResourceMetadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"rule1": {
					Name:        "Rule 1",
					Enforcement: "must",
				},
				"rule2": {
					Name:        "Rule 2",
					Enforcement: "should",
				},
				"rule3": {
					Name:        "Rule 3",
					Enforcement: "may",
				},
			},
		},
	}

	rule := &resource.Rule{
		Name:        "Rule 1",
		Enforcement: "must",
	}

	metadata := GenerateRuleMetadata("test-namespace", ruleset, "rule1", rule)

	// Check that all rules are listed
	if !strings.Contains(metadata, "- rule1") {
		t.Error("Expected rule1 in rules list")
	}
	if !strings.Contains(metadata, "- rule2") {
		t.Error("Expected rule2 in rules list")
	}
	if !strings.Contains(metadata, "- rule3") {
		t.Error("Expected rule3 in rules list")
	}

	// Check that rules are sorted (rule1, rule2, rule3)
	rule1Idx := strings.Index(metadata, "- rule1")
	rule2Idx := strings.Index(metadata, "- rule2")
	rule3Idx := strings.Index(metadata, "- rule3")

	if rule1Idx > rule2Idx || rule2Idx > rule3Idx {
		t.Error("Expected rules to be sorted alphabetically")
	}
}

func TestGenerateRuleMetadata_EnforcementUppercase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"must", "MUST"},
		{"should", "SHOULD"},
		{"may", "MAY"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			ruleset := &resource.RulesetResource{
				Metadata: resource.ResourceMetadata{
					ID:   "test-ruleset",
					Name: "Test Ruleset",
				},
				Spec: resource.RulesetSpec{
					Rules: map[string]resource.Rule{
						"rule1": {
							Name:        "Test Rule",
							Enforcement: tc.input,
						},
					},
				},
			}

			rule := &resource.Rule{
				Name:        "Test Rule",
				Enforcement: tc.input,
			}

			metadata := GenerateRuleMetadata("test-namespace", ruleset, "rule1", rule)

			if !strings.Contains(metadata, "enforcement: "+tc.expected) {
				t.Errorf("Expected enforcement: %s, but not found in metadata", tc.expected)
			}
		})
	}
}
