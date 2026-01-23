package compiler

import "testing"

func TestRuleGeneratorFactory(t *testing.T) {
	factory := NewRuleGeneratorFactory()

	validTools := []Tool{Cursor, Markdown, AmazonQ, Copilot}

	for _, tool := range validTools {
		t.Run(string(tool), func(t *testing.T) {
			generator, err := factory.NewRuleGenerator(tool)
			if err != nil {
				t.Fatalf("Expected no error for tool %s, got %v", tool, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for tool %s, got nil", tool)
			}
		})
	}

	// Test invalid tool
	_, err := factory.NewRuleGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid tool, got nil")
	}
}

func TestPromptGeneratorFactory(t *testing.T) {
	factory := NewPromptGeneratorFactory()

	validTools := []Tool{Cursor, Markdown, AmazonQ, Copilot}

	for _, tool := range validTools {
		t.Run(string(tool), func(t *testing.T) {
			generator, err := factory.NewPromptGenerator(tool)
			if err != nil {
				t.Fatalf("Expected no error for tool %s, got %v", tool, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for tool %s, got nil", tool)
			}
		})
	}

	// Test invalid tool
	_, err := factory.NewPromptGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid tool, got nil")
	}
}

func TestRuleFilenameGeneratorFactory(t *testing.T) {
	factory := NewRuleFilenameGeneratorFactory()

	validTools := []Tool{Cursor, Markdown, AmazonQ, Copilot}

	for _, tool := range validTools {
		t.Run(string(tool), func(t *testing.T) {
			generator, err := factory.NewRuleFilenameGenerator(tool)
			if err != nil {
				t.Fatalf("Expected no error for tool %s, got %v", tool, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for tool %s, got nil", tool)
			}
		})
	}

	// Test invalid tool
	_, err := factory.NewRuleFilenameGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid tool, got nil")
	}
}

func TestPromptFilenameGeneratorFactory(t *testing.T) {
	factory := NewPromptFilenameGeneratorFactory()

	validTools := []Tool{Cursor, Markdown, AmazonQ, Copilot}

	for _, tool := range validTools {
		t.Run(string(tool), func(t *testing.T) {
			generator, err := factory.NewPromptFilenameGenerator(tool)
			if err != nil {
				t.Fatalf("Expected no error for tool %s, got %v", tool, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for tool %s, got nil", tool)
			}
		})
	}

	// Test invalid tool
	_, err := factory.NewPromptFilenameGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid tool, got nil")
	}
}
