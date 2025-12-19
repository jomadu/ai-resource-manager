package compiler

import (
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

func TestCursorRuleGenerator_GenerateRule(t *testing.T) {
	generator := &CursorRuleGenerator{}

	ruleset := &resource.RulesetResource{
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"test-rule": {
					Name:        "Test Rule",
					Description: "A test rule",
					Enforcement: "must",
					Priority:    100,
					Scope: []resource.Scope{
						{Files: []string{"**/*.go", "**/*.js"}},
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
description: "A test rule"
globs: **/*.go, **/*.js
alwaysApply: true
---

---
namespace: test-namespace
ruleset:
  id: 
  name: 
  rules:
    - test-rule
rule:
  id: test-rule
  name: Test Rule
  enforcement: MUST
  priority: 100
  scope:
    - files: ["**/*.go", "**/*.js"]
---

# Test Rule

This is a test rule body.`

	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestCursorRuleGenerator_GenerateRule_NotFound(t *testing.T) {
	generator := &CursorRuleGenerator{}
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

func TestCursorPromptGenerator_GeneratePrompt(t *testing.T) {
	generator := &CursorPromptGenerator{}

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

func TestCursorPromptGenerator_GeneratePrompt_NotFound(t *testing.T) {
	generator := &CursorPromptGenerator{}
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

func TestCursorRuleFilenameGenerator_GenerateRuleFilename(t *testing.T) {
	generator := &CursorRuleFilenameGenerator{}

	result, err := generator.GenerateRuleFilename("my-ruleset", "my-rule")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "my-ruleset_my-rule.mdc"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCursorRuleFilenameGenerator_GenerateRuleFilename_EmptyInputs(t *testing.T) {
	generator := &CursorRuleFilenameGenerator{}

	_, err := generator.GenerateRuleFilename("", "rule")
	if err == nil {
		t.Fatal("Expected error for empty rulesetID")
	}

	_, err = generator.GenerateRuleFilename("ruleset", "")
	if err == nil {
		t.Fatal("Expected error for empty ruleID")
	}
}

func TestCursorPromptFilenameGenerator_GeneratePromptFilename(t *testing.T) {
	generator := &CursorPromptFilenameGenerator{}

	result, err := generator.GeneratePromptFilename("my-promptset", "my-prompt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "my-promptset_my-prompt.md"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCursorPromptFilenameGenerator_GeneratePromptFilename_EmptyInputs(t *testing.T) {
	generator := &CursorPromptFilenameGenerator{}

	_, err := generator.GeneratePromptFilename("", "prompt")
	if err == nil {
		t.Fatal("Expected error for empty promptsetID")
	}

	_, err = generator.GeneratePromptFilename("promptset", "")
	if err == nil {
		t.Fatal("Expected error for empty promptID")
	}
}