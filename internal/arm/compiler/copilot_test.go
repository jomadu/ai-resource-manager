package compiler

import (
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/resource"
)

func TestCopilotRuleGenerator_GenerateRule(t *testing.T) {
	generator := &CopilotRuleGenerator{}

	ruleset := &resource.RulesetResource{
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"test-rule": {
					Name:        "Test Rule",
					Description: "A test rule",
					Enforcement: "could",
					Priority:    50,
					Scope: []resource.Scope{
						{Files: []string{"**/*.ts", "**/*.tsx"}},
					},
					Body: "# Test Rule\n\nThis is a test rule body.",
				},
			},
		},
	}

	result, err := generator.GenerateRule("test-namespace", ruleset, "test-rule")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := `---
namespace: test-namespace
ruleset:
  id: 
  name: 
  rules:
    - test-rule
rule:
  id: test-rule
  name: Test Rule
  enforcement: COULD
  priority: 50
  scope:
    - files: ["**/*.ts", "**/*.tsx"]
---

# Test Rule

This is a test rule body.`

	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestCopilotRuleGenerator_GenerateRule_NotFound(t *testing.T) {
	generator := &CopilotRuleGenerator{}
	ruleset := &resource.RulesetResource{
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{},
		},
	}

	_, err := generator.GenerateRule("test-namespace", ruleset, "nonexistent")
	if err == nil {
		t.Fatal("Expected error for nonexistent rule")
	}
}

func TestCopilotPromptGenerator_GeneratePrompt(t *testing.T) {
	generator := &CopilotPromptGenerator{}

	promptset := &resource.PromptsetResource{
		Spec: resource.PromptsetSpec{
			Prompts: map[string]resource.Prompt{
				"test-prompt": {
					Name: "Test Prompt",
					Body: "This is a test prompt body.",
				},
			},
		},
	}

	result, err := generator.GeneratePrompt("test-namespace", promptset, "test-prompt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "This is a test prompt body."
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCopilotPromptGenerator_GeneratePrompt_NotFound(t *testing.T) {
	generator := &CopilotPromptGenerator{}
	promptset := &resource.PromptsetResource{
		Spec: resource.PromptsetSpec{
			Prompts: map[string]resource.Prompt{},
		},
	}

	_, err := generator.GeneratePrompt("test-namespace", promptset, "nonexistent")
	if err == nil {
		t.Fatal("Expected error for nonexistent prompt")
	}
}

func TestCopilotRuleFilenameGenerator_GenerateRuleFilename(t *testing.T) {
	generator := &CopilotRuleFilenameGenerator{}

	result, err := generator.GenerateRuleFilename("my-ruleset", "my-rule")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "my-ruleset_my-rule.instructions.md"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCopilotRuleFilenameGenerator_GenerateRuleFilename_EmptyInputs(t *testing.T) {
	generator := &CopilotRuleFilenameGenerator{}

	_, err := generator.GenerateRuleFilename("", "rule")
	if err == nil {
		t.Fatal("Expected error for empty rulesetID")
	}

	_, err = generator.GenerateRuleFilename("ruleset", "")
	if err == nil {
		t.Fatal("Expected error for empty ruleID")
	}
}

func TestCopilotPromptFilenameGenerator_GeneratePromptFilename(t *testing.T) {
	generator := &CopilotPromptFilenameGenerator{}

	result, err := generator.GeneratePromptFilename("my-promptset", "my-prompt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "my-promptset_my-prompt.md"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCopilotPromptFilenameGenerator_GeneratePromptFilename_EmptyInputs(t *testing.T) {
	generator := &CopilotPromptFilenameGenerator{}

	_, err := generator.GeneratePromptFilename("", "prompt")
	if err == nil {
		t.Fatal("Expected error for empty promptsetID")
	}

	_, err = generator.GeneratePromptFilename("promptset", "")
	if err == nil {
		t.Fatal("Expected error for empty promptID")
	}
}
