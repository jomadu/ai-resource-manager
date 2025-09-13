package urf

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestURFIntegration_EndToEnd(t *testing.T) {
	// Create URF service
	service := NewService()

	// Sample URF content matching the design document
	urfContent := `version: "1.0"
metadata:
  id: "security-rules"
  name: "Security Rules"
  version: "1.0.0"
  description: "Critical security validation rules"
rules:
  - id: "critical-security-check"
    name: "Critical Security Check"
    description: "Validate all user inputs to prevent security vulnerabilities"
    priority: 100
    enforcement: "must"
    scope:
      - files: ["**/*.go", "**/*.js"]
    body: |
      Always validate and sanitize user inputs to prevent injection attacks.
      Use parameterized queries for database operations.
  - id: "recommended-best-practice"
    name: "Recommended Best Practice"
    description: "Follow established coding conventions for maintainability"
    priority: 80
    enforcement: "should"
    scope:
      - files: ["**/*.go"]
    body: |
      Follow Go naming conventions and use gofmt for consistent formatting.`

	// Step 1: Process files to detect and parse URF
	files := []*types.File{
		{
			Path:    "security-rules.yaml",
			Content: []byte(urfContent),
			Size:    int64(len(urfContent)),
		},
		{
			Path:    "not-urf.txt",
			Content: []byte("regular text file"),
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
	if urf.Metadata.ID != "security-rules" {
		t.Errorf("URF ID = %s, expected security-rules", urf.Metadata.ID)
	}
	if len(urf.Rules) != 2 {
		t.Errorf("Rules count = %d, expected 2", len(urf.Rules))
	}

	// Step 2: Compile to Cursor format
	cursorFiles, err := service.CompileFiles(urfFiles, TargetCursor)
	if err != nil {
		t.Fatalf("Cursor compilation failed: %v", err)
	}

	if len(cursorFiles) != 2 {
		t.Fatalf("Expected 2 Cursor files, got %d", len(cursorFiles))
	}

	// Verify Cursor file structure
	cursorFile := cursorFiles[0]
	cursorContent := string(cursorFile.Content)

	// Should have Cursor-specific frontmatter
	if !strings.Contains(cursorContent, `description: "Validate all user inputs to prevent security vulnerabilities"`) {
		t.Error("Missing Cursor frontmatter description")
	}
	if !strings.Contains(cursorContent, "alwaysApply: true") {
		t.Error("Missing alwaysApply for must enforcement")
	}
	if !strings.Contains(cursorContent, "globs: [**/*.go **/*.js]") {
		t.Error("Missing globs in Cursor frontmatter")
	}

	// Should have shared metadata block
	if !strings.Contains(cursorContent, "namespace: security-rules") {
		t.Error("Missing namespace in metadata")
	}
	if !strings.Contains(cursorContent, "enforcement: MUST") {
		t.Error("Missing enforcement in metadata")
	}

	// Should have rule content
	if !strings.Contains(cursorContent, "# Critical Security Check (MUST)") {
		t.Error("Missing rule header")
	}
	if !strings.Contains(cursorContent, "Always validate and sanitize user inputs") {
		t.Error("Missing rule body")
	}

	// Step 3: Compile to Amazon Q format
	amazonqFiles, err := service.CompileFiles(urfFiles, TargetAmazonQ)
	if err != nil {
		t.Fatalf("Amazon Q compilation failed: %v", err)
	}

	if len(amazonqFiles) != 2 {
		t.Fatalf("Expected 2 Amazon Q files, got %d", len(amazonqFiles))
	}

	// Verify Amazon Q file structure
	amazonqFile := amazonqFiles[0]
	amazonqContent := string(amazonqFile.Content)

	// Should NOT have Cursor-specific frontmatter
	if strings.Contains(amazonqContent, "alwaysApply:") || strings.Contains(amazonqContent, "globs:") {
		t.Error("Amazon Q format should not contain Cursor frontmatter")
	}

	// Should have shared metadata block
	if !strings.Contains(amazonqContent, "namespace: security-rules") {
		t.Error("Missing namespace in Amazon Q metadata")
	}

	// Should have rule content
	if !strings.Contains(amazonqContent, "# Critical Security Check (MUST)") {
		t.Error("Missing rule header in Amazon Q format")
	}

	// Verify file extensions
	if !strings.HasSuffix(cursorFile.Path, ".mdc") {
		t.Errorf("Cursor file should have .mdc extension, got %s", cursorFile.Path)
	}
	if !strings.HasSuffix(amazonqFile.Path, ".md") {
		t.Errorf("Amazon Q file should have .md extension, got %s", amazonqFile.Path)
	}
}
