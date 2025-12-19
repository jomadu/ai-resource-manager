package compiler

import "testing"

func TestRuleGeneratorFactory(t *testing.T) {
	factory := NewRuleGeneratorFactory()
	
	validTargets := []CompileTarget{TargetCursor, TargetMarkdown, TargetAmazonQ, TargetCopilot}
	
	for _, target := range validTargets {
		t.Run(string(target), func(t *testing.T) {
			generator, err := factory.NewRuleGenerator(target)
			if err != nil {
				t.Fatalf("Expected no error for target %s, got %v", target, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for target %s, got nil", target)
			}
		})
	}
	
	// Test invalid target
	_, err := factory.NewRuleGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid target, got nil")
	}
}

func TestPromptGeneratorFactory(t *testing.T) {
	factory := NewPromptGeneratorFactory()
	
	validTargets := []CompileTarget{TargetCursor, TargetMarkdown, TargetAmazonQ, TargetCopilot}
	
	for _, target := range validTargets {
		t.Run(string(target), func(t *testing.T) {
			generator, err := factory.NewPromptGenerator(target)
			if err != nil {
				t.Fatalf("Expected no error for target %s, got %v", target, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for target %s, got nil", target)
			}
		})
	}
	
	// Test invalid target
	_, err := factory.NewPromptGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid target, got nil")
	}
}

func TestRuleFilenameGeneratorFactory(t *testing.T) {
	factory := NewRuleFilenameGeneratorFactory()
	
	validTargets := []CompileTarget{TargetCursor, TargetMarkdown, TargetAmazonQ, TargetCopilot}
	
	for _, target := range validTargets {
		t.Run(string(target), func(t *testing.T) {
			generator, err := factory.NewRuleFilenameGenerator(target)
			if err != nil {
				t.Fatalf("Expected no error for target %s, got %v", target, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for target %s, got nil", target)
			}
		})
	}
	
	// Test invalid target
	_, err := factory.NewRuleFilenameGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid target, got nil")
	}
}

func TestPromptFilenameGeneratorFactory(t *testing.T) {
	factory := NewPromptFilenameGeneratorFactory()
	
	validTargets := []CompileTarget{TargetCursor, TargetMarkdown, TargetAmazonQ, TargetCopilot}
	
	for _, target := range validTargets {
		t.Run(string(target), func(t *testing.T) {
			generator, err := factory.NewPromptFilenameGenerator(target)
			if err != nil {
				t.Fatalf("Expected no error for target %s, got %v", target, err)
			}
			if generator == nil {
				t.Fatalf("Expected generator for target %s, got nil", target)
			}
		})
	}
	
	// Test invalid target
	_, err := factory.NewPromptFilenameGenerator("invalid")
	if err == nil {
		t.Error("Expected error for invalid target, got nil")
	}
}