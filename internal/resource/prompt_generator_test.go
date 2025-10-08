package resource

import (
	"strings"
	"testing"
)

func TestCursorPromptGenerator_GeneratePrompt(t *testing.T) {
	generator := &DefaultPromptGenerator{}

	promptset := &Promptset{
		APIVersion: "v1",
		Kind:       "Promptset",
		Metadata: Metadata{
			ID:          "test-promptset",
			Name:        "Test Promptset",
			Description: "Test description",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]Prompt{
				"prompt1": {
					Name:        "Test Prompt 1",
					Description: "First test prompt",
					Body:        "This is the prompt body content.",
				},
			},
		},
	}

	prompt := promptset.Spec.Prompts["prompt1"]
	namespace := "ai-rules/test@1.0.0"

	result := generator.GeneratePrompt(namespace, promptset, "prompt1", &prompt)

	// Check that it's content-only (no metadata/frontmatter)
	if strings.Contains(result, "namespace:") {
		t.Error("Prompt should not contain namespace metadata")
	}
	if strings.Contains(result, "---") {
		t.Error("Prompt should not contain YAML frontmatter")
	}
	if strings.Contains(result, "alwaysApply:") {
		t.Error("Prompt should not contain rule-specific metadata")
	}

	// Check prompt content
	if !strings.Contains(result, "This is the prompt body content.") {
		t.Error("Expected prompt body content")
	}
}

func TestMarkdownPromptGenerator_GeneratePrompt(t *testing.T) {
	generator := &DefaultPromptGenerator{}

	promptset := &Promptset{
		APIVersion: "v1",
		Kind:       "Promptset",
		Metadata: Metadata{
			ID:          "test-promptset",
			Name:        "Test Promptset",
			Description: "Test description",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]Prompt{
				"prompt1": {
					Name:        "Test Prompt 1",
					Description: "First test prompt",
					Body:        "This is the prompt body content.",
				},
			},
		},
	}

	prompt := promptset.Spec.Prompts["prompt1"]
	namespace := "ai-rules/test@1.0.0"

	result := generator.GeneratePrompt(namespace, promptset, "prompt1", &prompt)

	// Check that it's content-only (no metadata/frontmatter)
	if strings.Contains(result, "namespace:") {
		t.Error("Prompt should not contain namespace metadata")
	}
	if strings.Contains(result, "---") {
		t.Error("Prompt should not contain YAML frontmatter")
	}
	if strings.Contains(result, "alwaysApply:") {
		t.Error("Prompt should not contain rule-specific metadata")
	}

	// Check prompt content
	if !strings.Contains(result, "This is the prompt body content.") {
		t.Error("Expected prompt body content")
	}
}

func TestCopilotPromptGenerator_GeneratePrompt(t *testing.T) {
	generator := &DefaultPromptGenerator{}

	promptset := &Promptset{
		APIVersion: "v1",
		Kind:       "Promptset",
		Metadata: Metadata{
			ID:          "test-promptset",
			Name:        "Test Promptset",
			Description: "Test description",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]Prompt{
				"prompt1": {
					Name:        "Test Prompt 1",
					Description: "First test prompt",
					Body:        "This is the prompt body content.",
				},
			},
		},
	}

	prompt := promptset.Spec.Prompts["prompt1"]
	namespace := "ai-rules/test@1.0.0"

	result := generator.GeneratePrompt(namespace, promptset, "prompt1", &prompt)

	// Check that it's content-only (no metadata/frontmatter)
	if strings.Contains(result, "namespace:") {
		t.Error("Prompt should not contain namespace metadata")
	}
	if strings.Contains(result, "---") {
		t.Error("Prompt should not contain YAML frontmatter")
	}
	if strings.Contains(result, "alwaysApply:") {
		t.Error("Prompt should not contain rule-specific metadata")
	}

	// Check prompt content
	if !strings.Contains(result, "This is the prompt body content.") {
		t.Error("Expected prompt body content")
	}
}

func TestPromptGeneratorFactory_NewPromptGenerator(t *testing.T) {
	factory := NewPromptGeneratorFactory()

	// Test Cursor target
	cursorGen, err := factory.NewPromptGenerator(TargetCursor)
	if err != nil {
		t.Fatalf("Expected no error for Cursor target, got %v", err)
	}
	if _, ok := cursorGen.(*DefaultPromptGenerator); !ok {
		t.Error("Expected DefaultPromptGenerator for Cursor target")
	}

	// Test Amazon Q target
	amazonqGen, err := factory.NewPromptGenerator(TargetAmazonQ)
	if err != nil {
		t.Fatalf("Expected no error for Amazon Q target, got %v", err)
	}
	if _, ok := amazonqGen.(*DefaultPromptGenerator); !ok {
		t.Error("Expected DefaultPromptGenerator for Amazon Q target")
	}

	// Test Markdown target
	markdownGen, err := factory.NewPromptGenerator(TargetMarkdown)
	if err != nil {
		t.Fatalf("Expected no error for Markdown target, got %v", err)
	}
	if _, ok := markdownGen.(*DefaultPromptGenerator); !ok {
		t.Error("Expected DefaultPromptGenerator for Markdown target")
	}

	// Test Copilot target
	copilotGen, err := factory.NewPromptGenerator(TargetCopilot)
	if err != nil {
		t.Fatalf("Expected no error for Copilot target, got %v", err)
	}
	if _, ok := copilotGen.(*DefaultPromptGenerator); !ok {
		t.Error("Expected DefaultPromptGenerator for Copilot target")
	}

	// Test unsupported target
	_, err = factory.NewPromptGenerator("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported target")
	}
}

func TestPromptGenerator_AllTargetsProduceContentOnly(t *testing.T) {
	factory := NewPromptGeneratorFactory()
	targets := []CompileTarget{TargetCursor, TargetAmazonQ, TargetMarkdown, TargetCopilot}

	promptset := &Promptset{
		APIVersion: "v1",
		Kind:       "Promptset",
		Metadata: Metadata{
			ID:   "test-promptset",
			Name: "Test Promptset",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]Prompt{
				"prompt1": {
					Name: "Test Prompt",
					Body: "This is the prompt content.",
				},
			},
		},
	}

	prompt := promptset.Spec.Prompts["prompt1"]
	namespace := "ai-rules/test@1.0.0"

	for _, target := range targets {
		t.Run(string(target), func(t *testing.T) {
			generator, err := factory.NewPromptGenerator(target)
			if err != nil {
				t.Fatalf("Failed to create %s prompt generator: %v", target, err)
			}

			result := generator.GeneratePrompt(namespace, promptset, "prompt1", &prompt)

			// All prompt generators should produce content-only output
			if strings.Contains(result, "namespace:") {
				t.Errorf("Prompt should not contain namespace metadata for %s target", target)
			}
			if strings.Contains(result, "---") {
				t.Errorf("Prompt should not contain YAML frontmatter for %s target", target)
			}
			if strings.Contains(result, "alwaysApply:") {
				t.Errorf("Prompt should not contain rule-specific metadata for %s target", target)
			}
			if !strings.Contains(result, "This is the prompt content.") {
				t.Errorf("Expected prompt body content for %s target", target)
			}
		})
	}
}
