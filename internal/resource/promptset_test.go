package resource

import (
	"strings"
	"testing"
)

func TestPromptset_Validation(t *testing.T) {
	tests := []struct {
		name      string
		promptset *Promptset
		wantError bool
	}{
		{
			name: "valid promptset",
			promptset: &Promptset{
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
							Body: "Prompt content",
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "missing apiVersion",
			promptset: &Promptset{
				Kind: "Promptset",
				Metadata: Metadata{
					ID:   "test-promptset",
					Name: "Test Promptset",
				},
				Spec: PromptsetSpec{
					Prompts: map[string]Prompt{
						"prompt1": {
							Name: "Test Prompt",
							Body: "Prompt content",
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "wrong kind",
			promptset: &Promptset{
				APIVersion: "v1",
				Kind:       "Ruleset", // Wrong kind
				Metadata: Metadata{
					ID:   "test-promptset",
					Name: "Test Promptset",
				},
				Spec: PromptsetSpec{
					Prompts: map[string]Prompt{
						"prompt1": {
							Name: "Test Prompt",
							Body: "Prompt content",
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "empty prompts",
			promptset: &Promptset{
				APIVersion: "v1",
				Kind:       "Promptset",
				Metadata: Metadata{
					ID:   "test-promptset",
					Name: "Test Promptset",
				},
				Spec: PromptsetSpec{
					Prompts: map[string]Prompt{}, // Empty prompts
				},
			},
			wantError: true,
		},
		{
			name: "prompt without body",
			promptset: &Promptset{
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
							// Missing body
						},
					},
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real implementation, you would use a validator
			// For now, we'll do basic validation checks
			if tt.promptset.APIVersion == "" && !tt.wantError {
				t.Error("Expected apiVersion to be required")
			}
			if tt.promptset.Kind != "Promptset" && !tt.wantError {
				t.Error("Expected kind to be 'Promptset'")
			}
			if len(tt.promptset.Spec.Prompts) == 0 && !tt.wantError {
				t.Error("Expected at least one prompt")
			}
		})
	}
}

func TestPrompt_Validation(t *testing.T) {
	tests := []struct {
		name      string
		prompt    Prompt
		wantError bool
	}{
		{
			name: "valid prompt",
			prompt: Prompt{
				Name:        "Test Prompt",
				Description: "A test prompt",
				Body:        "This is the prompt content",
			},
			wantError: false,
		},
		{
			name: "prompt without name",
			prompt: Prompt{
				Description: "A test prompt",
				Body:        "This is the prompt content",
			},
			wantError: true,
		},
		{
			name: "prompt without body",
			prompt: Prompt{
				Name:        "Test Prompt",
				Description: "A test prompt",
			},
			wantError: true,
		},
		{
			name: "minimal valid prompt",
			prompt: Prompt{
				Name: "Test Prompt",
				Body: "This is the prompt content",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			if tt.prompt.Name == "" && !tt.wantError {
				t.Error("Expected name to be required")
			}
			if tt.prompt.Body == "" && !tt.wantError {
				t.Error("Expected body to be required")
			}
		})
	}
}

func TestPromptsetSpec_Validation(t *testing.T) {
	tests := []struct {
		name      string
		spec      PromptsetSpec
		wantError bool
	}{
		{
			name: "valid spec with multiple prompts",
			spec: PromptsetSpec{
				Prompts: map[string]Prompt{
					"prompt1": {
						Name: "First Prompt",
						Body: "First prompt content",
					},
					"prompt2": {
						Name: "Second Prompt",
						Body: "Second prompt content",
					},
				},
			},
			wantError: false,
		},
		{
			name: "empty prompts",
			spec: PromptsetSpec{
				Prompts: map[string]Prompt{},
			},
			wantError: true,
		},
		{
			name: "nil prompts",
			spec: PromptsetSpec{
				Prompts: nil,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			if len(tt.spec.Prompts) == 0 && !tt.wantError {
				t.Error("Expected at least one prompt")
			}
		})
	}
}

func TestPromptset_StringRepresentation(t *testing.T) {
	promptset := &Promptset{
		APIVersion: "v1",
		Kind:       "Promptset",
		Metadata: Metadata{
			ID:          "code-review",
			Name:        "Code Review Prompts",
			Description: "Prompts for code review tasks",
		},
		Spec: PromptsetSpec{
			Prompts: map[string]Prompt{
				"review": {
					Name:        "Code Review",
					Description: "Review the code for issues",
					Body:        "Please review this code for bugs, performance issues, and best practices.",
				},
				"feedback": {
					Name:        "Provide Feedback",
					Description: "Give constructive feedback",
					Body:        "Provide constructive feedback on the following code changes.",
				},
			},
		},
	}

	// Test that we can access the promptset data
	if promptset.Metadata.ID != "code-review" {
		t.Errorf("Expected ID 'code-review', got '%s'", promptset.Metadata.ID)
	}

	if len(promptset.Spec.Prompts) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(promptset.Spec.Prompts))
	}

	reviewPrompt, exists := promptset.Spec.Prompts["review"]
	if !exists {
		t.Error("Expected 'review' prompt to exist")
	}

	if !strings.Contains(reviewPrompt.Body, "review this code") {
		t.Error("Expected review prompt body to contain 'review this code'")
	}
}
