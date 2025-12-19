package compiler

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

func TestCompileRuleset(t *testing.T) {
	ruleset := &resource.RulesetResource{
		APIVersion: "v1",
		Kind:       "Ruleset",
		Metadata: resource.ResourceMetadata{
			ID:   "test-ruleset",
			Name: "Test Ruleset",
		},
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
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

	namespace := "sample-registry"

	files, err := CompileRuleset(TargetCursor, namespace, ruleset)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	// Should generate one file per rule
	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}

	// Files should be sorted by rule ID
	if files[0].Path != "test-ruleset_rule1.mdc" {
		t.Errorf("Expected first file test-ruleset_rule1.mdc, got %s", files[0].Path)
	}
	if files[1].Path != "test-ruleset_rule2.mdc" {
		t.Errorf("Expected second file test-ruleset_rule2.mdc, got %s", files[1].Path)
	}

	// Check first file content
	content1 := string(files[0].Content)
	if !strings.Contains(content1, "namespace: "+namespace) {
		t.Error("Expected namespace in first file content")
	}
	if !strings.Contains(content1, "Rule 1 content") {
		t.Error("Expected rule 1 body in first file content")
	}
	if !strings.Contains(content1, "alwaysApply: true") {
		t.Error("Expected alwaysApply: true for 'must' enforcement")
	}

	// Check second file content
	content2 := string(files[1].Content)
	if !strings.Contains(content2, "Rule 2 content") {
		t.Error("Expected rule 2 body in second file content")
	}
	if strings.Contains(content2, "alwaysApply: true") {
		t.Error("Should not have alwaysApply: true for 'should' enforcement")
	}
}

func TestCompileRuleset_DifferentTargets(t *testing.T) {
	testCases := []struct {
		target    CompileTarget
		extension string
	}{
		{TargetCursor, ".mdc"},
		{TargetAmazonQ, ".md"},
		{TargetCopilot, ".instructions.md"},
		{TargetMarkdown, ".md"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.target), func(t *testing.T) {
			ruleset := &resource.RulesetResource{
				APIVersion: "v1",
				Kind:       "Ruleset",
				Metadata: resource.ResourceMetadata{
					ID:   "test-ruleset",
					Name: "Test Ruleset",
				},
				Spec: resource.RulesetSpec{
					Rules: map[string]resource.Rule{
						"rule1": {
							Name: "Test Rule",
							Body: "Rule content",
						},
					},
				},
			}

			files, err := CompileRuleset(tc.target, "test-namespace", ruleset)
			if err != nil {
				t.Fatalf("Compilation failed for %s: %v", tc.target, err)
			}

			if len(files) != 1 {
				t.Fatalf("Expected 1 file for %s, got %d", tc.target, len(files))
			}

			expectedFilename := "test-ruleset_rule1" + tc.extension
			if files[0].Path != expectedFilename {
				t.Errorf("Expected filename %s for %s, got %s", expectedFilename, tc.target, files[0].Path)
			}
		})
	}
}

func TestCompilePromptset(t *testing.T) {
	promptset := &resource.PromptsetResource{
		APIVersion: "v1",
		Kind:       "Promptset",
		Metadata: resource.ResourceMetadata{
			ID:   "test-promptset",
			Name: "Test Promptset",
		},
		Spec: resource.PromptsetSpec{
			Prompts: map[string]resource.Prompt{
				"prompt1": {
					Name:        "Test Prompt 1",
					Description: "First prompt",
					Body:        "Prompt 1 content",
				},
				"prompt2": {
					Name:        "Test Prompt 2",
					Description: "Second prompt", 
					Body:        "Prompt 2 content",
				},
			},
		},
	}

	namespace := "sample-registry"

	files, err := CompilePromptset(TargetCursor, namespace, promptset)
	if err != nil {
		t.Fatalf("Promptset compilation failed: %v", err)
	}

	// Should generate one file per prompt
	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}

	// Files should be sorted by prompt ID
	if files[0].Path != "test-promptset_prompt1.md" {
		t.Errorf("Expected first file test-promptset_prompt1.md, got %s", files[0].Path)
	}
	if files[1].Path != "test-promptset_prompt2.md" {
		t.Errorf("Expected second file test-promptset_prompt2.md, got %s", files[1].Path)
	}

	// Check file contents (prompts should have no metadata)
	content1 := string(files[0].Content)
	if !strings.Contains(content1, "Prompt 1 content") {
		t.Error("Expected prompt 1 body in first file content")
	}
	if strings.Contains(content1, "namespace:") {
		t.Error("Promptset should not contain metadata")
	}

	content2 := string(files[1].Content)
	if !strings.Contains(content2, "Prompt 2 content") {
		t.Error("Expected prompt 2 body in second file content")
	}
	if strings.Contains(content2, "alwaysApply:") {
		t.Error("Promptset should not contain rule metadata")
	}
}

func TestCompilePromptset_AllTargets(t *testing.T) {
	targets := []CompileTarget{TargetCursor, TargetAmazonQ, TargetMarkdown, TargetCopilot}

	for _, target := range targets {
		t.Run(string(target), func(t *testing.T) {
			promptset := &resource.PromptsetResource{
				APIVersion: "v1",
				Kind:       "Promptset",
				Metadata: resource.ResourceMetadata{
					ID:   "test-promptset",
					Name: "Test Promptset",
				},
				Spec: resource.PromptsetSpec{
					Prompts: map[string]resource.Prompt{
						"prompt1": {
							Name: "Test Prompt",
							Body: "Prompt content",
						},
					},
				},
			}

			files, err := CompilePromptset(target, "test-namespace", promptset)
			if err != nil {
				t.Fatalf("Promptset compilation failed for %s: %v", target, err)
			}

			if len(files) != 1 {
				t.Fatalf("Expected 1 file for %s, got %d", target, len(files))
			}

			// All promptset targets should use .md extension
			if !strings.HasSuffix(files[0].Path, ".md") {
				t.Errorf("Expected .md extension for %s promptset, got %s", target, files[0].Path)
			}

			// All promptset targets should produce content-only files (no metadata)
			content := string(files[0].Content)
			if strings.Contains(content, "namespace:") {
				t.Errorf("Promptset should not contain metadata for %s target", target)
			}
			if strings.Contains(content, "alwaysApply:") {
				t.Errorf("Promptset should not contain rule metadata for %s target", target)
			}
			if !strings.Contains(content, "Prompt content") {
				t.Errorf("Expected prompt body in %s output", target)
			}
		})
	}
}