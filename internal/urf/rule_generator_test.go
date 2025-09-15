package urf

import "testing"

func TestRuleGeneratorFactory_NewRuleGenerator(t *testing.T) {
	factory := NewRuleGeneratorFactory()

	// Test Cursor target
	cursorGen, err := factory.NewRuleGenerator(TargetCursor)
	if err != nil {
		t.Fatalf("Expected no error for Cursor target, got %v", err)
	}
	if _, ok := cursorGen.(*CursorRuleGenerator); !ok {
		t.Error("Expected CursorRuleGenerator for Cursor target")
	}

	// Test Amazon Q target
	amazonqGen, err := factory.NewRuleGenerator(TargetAmazonQ)
	if err != nil {
		t.Fatalf("Expected no error for Amazon Q target, got %v", err)
	}
	if _, ok := amazonqGen.(*AmazonQRuleGenerator); !ok {
		t.Error("Expected AmazonQRuleGenerator for Amazon Q target")
	}

	// Test unsupported target
	_, err = factory.NewRuleGenerator("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported target")
	}
}

func TestRuleGeneratorFactory_CreatesWithMetadataGenerator(t *testing.T) {
	factory := NewRuleGeneratorFactory()

	cursorGen, err := factory.NewRuleGenerator(TargetCursor)
	if err != nil {
		t.Fatalf("Failed to create Cursor generator: %v", err)
	}

	// Verify the generator has a metadata generator
	cursorRuleGen := cursorGen.(*CursorRuleGenerator)
	if cursorRuleGen.metadataGen == nil {
		t.Error("Expected CursorRuleGenerator to have metadata generator")
	}

	amazonqGen, err := factory.NewRuleGenerator(TargetAmazonQ)
	if err != nil {
		t.Fatalf("Failed to create Amazon Q generator: %v", err)
	}

	// Verify the generator has a metadata generator
	amazonqRuleGen := amazonqGen.(*AmazonQRuleGenerator)
	if amazonqRuleGen.metadataGen == nil {
		t.Error("Expected AmazonQRuleGenerator to have metadata generator")
	}
}
