package urf

import (
	"strings"
	"testing"
)

func TestDefaultMetadataGenerator_GenerateMetadata(t *testing.T) {
	generator := NewMetadataGenerator()

	urf := &URFFile{
		Metadata: Metadata{
			ID:      "test-ruleset",
			Name:    "Test Ruleset",
			Version: "1.0.0",
		},
		Rules: []Rule{
			{ID: "rule1", Name: "Rule 1"},
			{ID: "rule2", Name: "Rule 2"},
		},
	}

	rule := Rule{
		ID:          "rule1",
		Name:        "Rule 1",
		Enforcement: "must",
		Priority:    100,
		Scope: []Scope{
			{Files: []string{"**/*.go", "**/*.js"}},
		},
	}

	result := generator.GenerateMetadata(urf, &rule, "test-namespace")

	// Verify structure
	if !strings.Contains(result, "namespace: test-namespace") {
		t.Error("Missing namespace")
	}
	if !strings.Contains(result, "id: test-ruleset") {
		t.Error("Missing ruleset ID")
	}
	if !strings.Contains(result, "name: Test Ruleset") {
		t.Error("Missing ruleset name")
	}
	if !strings.Contains(result, "version: 1.0.0") {
		t.Error("Missing ruleset version")
	}
	if !strings.Contains(result, "- rule1") {
		t.Error("Missing rule1 in rules list")
	}
	if !strings.Contains(result, "- rule2") {
		t.Error("Missing rule2 in rules list")
	}
	if !strings.Contains(result, "enforcement: MUST") {
		t.Error("Missing enforcement")
	}
	if !strings.Contains(result, "priority: 100") {
		t.Error("Missing priority")
	}
	if !strings.Contains(result, `files: "**/*.go"`) {
		t.Error("Missing scope files")
	}
}

func TestDefaultMetadataGenerator_GenerateMetadata_NoScope(t *testing.T) {
	generator := NewMetadataGenerator()

	urf := &URFFile{
		Metadata: Metadata{
			ID:      "test-ruleset",
			Name:    "Test Ruleset",
			Version: "1.0.0",
		},
		Rules: []Rule{{ID: "rule1"}},
	}

	rule := Rule{
		ID:          "rule1",
		Name:        "Rule 1",
		Enforcement: "should",
		Priority:    50,
	}

	result := generator.GenerateMetadata(urf, &rule, "test-namespace")

	// Should not contain scope section
	if strings.Contains(result, "scope:") {
		t.Error("Should not contain scope section when no scope defined")
	}
	if !strings.Contains(result, "enforcement: SHOULD") {
		t.Error("Missing enforcement")
	}
}
