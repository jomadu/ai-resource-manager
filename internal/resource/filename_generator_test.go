package resource

import "testing"

func TestCursorFilenameGenerator_GenerateFilename(t *testing.T) {
	generator := &CursorFilenameGenerator{}

	result := generator.GenerateFilename("test-ruleset", "rule-1")
	expected := "test-ruleset_rule-1.mdc"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestMarkdownFilenameGenerator_GenerateFilename(t *testing.T) {
	generator := &MarkdownFilenameGenerator{}

	result := generator.GenerateFilename("test-ruleset", "rule-1")
	expected := "test-ruleset_rule-1.md"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestFilenameGeneratorFactory_NewFilenameGenerator(t *testing.T) {
	factory := NewFilenameGeneratorFactory()

	// Test Cursor target
	cursorGen, err := factory.NewFilenameGenerator(TargetCursor)
	if err != nil {
		t.Fatalf("Expected no error for Cursor target, got %v", err)
	}
	if _, ok := cursorGen.(*CursorFilenameGenerator); !ok {
		t.Error("Expected CursorFilenameGenerator for Cursor target")
	}

	// Test Amazon Q target
	amazonqGen, err := factory.NewFilenameGenerator(TargetAmazonQ)
	if err != nil {
		t.Fatalf("Expected no error for Amazon Q target, got %v", err)
	}
	if _, ok := amazonqGen.(*MarkdownFilenameGenerator); !ok {
		t.Error("Expected MarkdownFilenameGenerator for Amazon Q target")
	}

	// Test unsupported target
	_, err = factory.NewFilenameGenerator("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported target")
	}
}
