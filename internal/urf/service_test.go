package urf

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestService_ProcessFiles(t *testing.T) {
	service := NewService()

	files := []*types.File{
		{
			Path: "ruleset.yaml",
			Content: []byte(`version: "1.0"
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
  version: "1.0.0"
rules:
  - id: "rule1"
    name: "Rule 1"
    body: "Content 1"`),
		},
		{
			Path:    "not-urf.txt",
			Content: []byte("regular text file"),
		},
		{
			Path:    "invalid.yaml",
			Content: []byte("key: value"),
		},
	}

	urfFiles, err := service.ProcessFiles(files)
	if err != nil {
		t.Fatalf("ProcessFiles failed: %v", err)
	}

	if len(urfFiles) != 1 {
		t.Fatalf("Expected 1 URF file, got %d", len(urfFiles))
	}

	urf := urfFiles[0]
	if urf.Metadata.ID != "test-ruleset" {
		t.Errorf("URF ID = %s, expected test-ruleset", urf.Metadata.ID)
	}
	if len(urf.Rules) != 1 {
		t.Errorf("Rules count = %d, expected 1", len(urf.Rules))
	}
}

func TestService_ProcessFiles_SkipsInvalidURF(t *testing.T) {
	service := NewService()

	// This file looks like URF but is invalid, so it should be skipped
	files := []*types.File{
		{
			Path: "invalid.yaml",
			Content: []byte(`version: "1.0"
metadata:
  id: "test-ruleset"
rules:
  - id: "rule1"
    name: "Rule 1"`), // Missing required metadata.name and metadata.version
		},
	}

	urfFiles, err := service.ProcessFiles(files)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should skip invalid files and return empty result
	if len(urfFiles) != 0 {
		t.Errorf("Expected 0 URF files (invalid should be skipped), got %d", len(urfFiles))
	}
}

func TestService_CompileFiles(t *testing.T) {
	service := NewService()

	urfFiles := []*URFFile{
		{
			Metadata: Metadata{
				ID:      "ruleset1",
				Name:    "Ruleset 1",
				Version: "1.0.0",
			},
			Rules: []Rule{
				{ID: "rule1", Name: "Rule 1", Body: "Content 1"},
			},
		},
		{
			Metadata: Metadata{
				ID:      "ruleset2",
				Name:    "Ruleset 2",
				Version: "2.0.0",
			},
			Rules: []Rule{
				{ID: "rule2", Name: "Rule 2", Body: "Content 2"},
				{ID: "rule3", Name: "Rule 3", Body: "Content 3"},
			},
		},
	}

	// Test Cursor compilation
	cursorFiles, err := service.CompileFiles(urfFiles, TargetCursor)
	if err != nil {
		t.Fatalf("CompileFiles for Cursor failed: %v", err)
	}

	if len(cursorFiles) != 3 { // 1 + 2 rules
		t.Errorf("Expected 3 Cursor files, got %d", len(cursorFiles))
	}

	// Verify file extensions
	for _, file := range cursorFiles {
		if !strings.HasSuffix(file.Path, ".mdc") {
			t.Errorf("Cursor file should have .mdc extension, got %s", file.Path)
		}
	}

	// Test Amazon Q compilation
	amazonqFiles, err := service.CompileFiles(urfFiles, TargetAmazonQ)
	if err != nil {
		t.Fatalf("CompileFiles for Amazon Q failed: %v", err)
	}

	if len(amazonqFiles) != 3 {
		t.Errorf("Expected 3 Amazon Q files, got %d", len(amazonqFiles))
	}

	// Verify file extensions
	for _, file := range amazonqFiles {
		if !strings.HasSuffix(file.Path, ".md") {
			t.Errorf("Amazon Q file should have .md extension, got %s", file.Path)
		}
	}
}

func TestService_CompileFiles_UnsupportedTarget(t *testing.T) {
	service := NewService()

	urfFiles := []*URFFile{
		{
			Metadata: Metadata{ID: "test", Name: "Test", Version: "1.0.0"},
			Rules:    []Rule{{ID: "rule1", Name: "Rule 1"}},
		},
	}

	_, err := service.CompileFiles(urfFiles, CompileTarget("unsupported"))
	if err == nil {
		t.Error("Expected error for unsupported target")
	}
}

func TestService_GetSupportedTargets(t *testing.T) {
	service := NewService()
	targets := service.GetSupportedTargets()

	if len(targets) != 2 {
		t.Errorf("Expected 2 supported targets, got %d", len(targets))
	}

	targetMap := make(map[CompileTarget]bool)
	for _, target := range targets {
		targetMap[target] = true
	}

	if !targetMap[TargetCursor] {
		t.Error("Missing Cursor target")
	}
	if !targetMap[TargetAmazonQ] {
		t.Error("Missing Amazon Q target")
	}
}
