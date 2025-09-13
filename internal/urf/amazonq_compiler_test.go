package urf

import (
	"strings"
	"testing"
)

func TestAmazonQCompiler_Compile(t *testing.T) {
	compiler := NewAmazonQCompiler()

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
					{Files: []string{"**/*.py"}},
				},
				Body: "Rule 1 content",
			},
		},
	}

	files, err := compiler.Compile(urf)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	file := files[0]
	if file.Path != "test-ruleset_rule1.md" {
		t.Errorf("File path = %s, expected test-ruleset_rule1.md", file.Path)
	}

	content := string(file.Content)

	// Amazon Q format should start directly with metadata (no Cursor frontmatter)
	if strings.Contains(content, "description:") && strings.Contains(content, "globs:") {
		t.Error("Should not contain Cursor-specific frontmatter")
	}

	// Check metadata block
	if !strings.Contains(content, "namespace: test-ruleset") {
		t.Error("Missing namespace in metadata")
	}
	if !strings.Contains(content, "enforcement: MUST") {
		t.Error("Missing enforcement in metadata")
	}

	// Check rule content
	if !strings.Contains(content, "# Rule 1 (MUST)") {
		t.Error("Missing rule header")
	}
	if !strings.Contains(content, "Rule 1 content") {
		t.Error("Missing rule body")
	}
}

func TestAmazonQCompiler_Compile_MultipleRules(t *testing.T) {
	compiler := NewAmazonQCompiler()

	urf := &URFFile{
		Metadata: Metadata{
			ID:      "test-ruleset",
			Name:    "Test Ruleset",
			Version: "1.0.0",
		},
		Rules: []Rule{
			{ID: "rule1", Name: "Rule 1", Body: "Content 1"},
			{ID: "rule2", Name: "Rule 2", Body: "Content 2"},
			{ID: "rule3", Name: "Rule 3", Body: "Content 3"},
		},
	}

	files, err := compiler.Compile(urf)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(files))
	}

	expectedPaths := []string{
		"test-ruleset_rule1.md",
		"test-ruleset_rule2.md",
		"test-ruleset_rule3.md",
	}

	for i, file := range files {
		if file.Path != expectedPaths[i] {
			t.Errorf("File %d path = %s, expected %s", i, file.Path, expectedPaths[i])
		}
	}
}
