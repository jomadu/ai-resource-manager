package urf

import (
	"strings"
	"testing"
)

func TestCursorCompiler_Compile(t *testing.T) {
	compiler := NewCursorCompiler()

	urf := &URFFile{
		Metadata: Metadata{
			ID:          "test-ruleset",
			Name:        "Test Ruleset",
			Version:     "1.0.0",
			Description: "Test description",
		},
		Rules: []Rule{
			{
				ID:          "rule1",
				Name:        "Rule 1",
				Description: "First rule",
				Priority:    100,
				Enforcement: "must",
				Scope: []Scope{
					{Files: []string{"**/*.go"}},
				},
				Body: "Rule 1 content",
			},
			{
				ID:          "rule2",
				Name:        "Rule 2",
				Description: "Second rule",
				Priority:    80,
				Enforcement: "should",
				Body:        "Rule 2 content",
			},
		},
	}

	files, err := compiler.Compile(urf)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}

	// Test first file
	file1 := files[0]
	if file1.Path != "test-ruleset_rule1.mdc" {
		t.Errorf("File1 path = %s, expected test-ruleset_rule1.mdc", file1.Path)
	}

	content1 := string(file1.Content)

	// Check Cursor-specific frontmatter
	if !strings.Contains(content1, `description: "First rule"`) {
		t.Error("Missing description in frontmatter")
	}
	if !strings.Contains(content1, `globs: [**/*.go]`) {
		t.Error("Missing globs in frontmatter")
	}
	if !strings.Contains(content1, "alwaysApply: true") {
		t.Error("Missing alwaysApply for must enforcement")
	}

	// Check metadata block
	if !strings.Contains(content1, "namespace: test-ruleset") {
		t.Error("Missing namespace in metadata")
	}

	// Check rule content
	if !strings.Contains(content1, "# Rule 1 (MUST)") {
		t.Error("Missing rule header")
	}
	if !strings.Contains(content1, "Rule 1 content") {
		t.Error("Missing rule body")
	}

	// Test second file (should enforcement)
	file2 := files[1]
	content2 := string(file2.Content)

	if strings.Contains(content2, "alwaysApply: true") {
		t.Error("Should not have alwaysApply for should enforcement")
	}
	if !strings.Contains(content2, "# Rule 2 (SHOULD)") {
		t.Error("Missing rule header with SHOULD enforcement")
	}
}

func TestCursorCompiler_Compile_EmptyRules(t *testing.T) {
	compiler := NewCursorCompiler()

	urf := &URFFile{
		Metadata: Metadata{
			ID:      "test-ruleset",
			Name:    "Test Ruleset",
			Version: "1.0.0",
		},
		Rules: []Rule{},
	}

	files, err := compiler.Compile(urf)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files for empty rules, got %d", len(files))
	}
}
