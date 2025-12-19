package filetype

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsRulesetFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		want     bool
	}{
		{
			name:     "valid ruleset yml",
			filename: "test.yml",
			content: `apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
  name: "Test"
spec:
  rules:
    rule1:
      name: "Rule 1"
      body: "test rule"`,
			want: true,
		},
		{
			name:     "valid ruleset yaml",
			filename: "test.yaml",
			content: `apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
  name: "Test"
spec:
  rules:
    rule1:
      name: "Rule 1"
      body: "test rule"`,
			want: true,
		},
		{
			name:     "wrong extension",
			filename: "test.txt",
			content: `apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
spec:
  rules:
    rule1:
      body: "test"`,
			want: false,
		},
		{
			name:     "wrong kind",
			filename: "test.yml",
			content: `apiVersion: v1
kind: Promptset
metadata:
  id: "test"
spec:
  prompts:
    prompt1:
      body: "test"`,
			want: false,
		},
		{
			name:     "missing required fields",
			filename: "test.yml",
			content: `apiVersion: v1
kind: Ruleset`,
			want: false,
		},
		{
			name:     "invalid yaml",
			filename: "test.yml",
			content:  `invalid: yaml: content:`,
			want:     false,
		},
		{
			name:     "file does not exist",
			filename: "nonexistent.yml",
			content:  "",
			want:     false,
		},
		{
			name:     "case insensitive extension",
			filename: "test.YML",
			content: `apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
spec:
  rules:
    rule1:
      body: "test"`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testFile string
			if tt.content != "" && tt.filename != "nonexistent.yml" {
				// Create temp file
				tmpDir := t.TempDir()
				testFile = filepath.Join(tmpDir, tt.filename)
				err := os.WriteFile(testFile, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			} else if tt.filename == "nonexistent.yml" {
				// Use non-existent path
				testFile = filepath.Join(t.TempDir(), "nonexistent.yml")
			}

			got := IsRulesetFile(testFile)
			if got != tt.want {
				t.Errorf("IsRulesetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPromptsetFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		want     bool
	}{
		{
			name:     "valid promptset yml",
			filename: "test.yml",
			content: `apiVersion: v1
kind: Promptset
metadata:
  id: "test"
  name: "Test"
spec:
  prompts:
    prompt1:
      name: "Prompt 1"
      body: "test prompt"`,
			want: true,
		},
		{
			name:     "valid promptset yaml",
			filename: "test.yaml",
			content: `apiVersion: v1
kind: Promptset
metadata:
  id: "test"
spec:
  prompts:
    prompt1:
      body: "test prompt"`,
			want: true,
		},
		{
			name:     "wrong extension",
			filename: "test.txt",
			content: `apiVersion: v1
kind: Promptset
metadata:
  id: "test"
spec:
  prompts:
    prompt1:
      body: "test"`,
			want: false,
		},
		{
			name:     "wrong kind",
			filename: "test.yml",
			content: `apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
spec:
  rules:
    rule1:
      body: "test"`,
			want: false,
		},
		{
			name:     "missing required fields",
			filename: "test.yml",
			content: `apiVersion: v1
kind: Promptset`,
			want: false,
		},
		{
			name:     "invalid yaml",
			filename: "test.yml",
			content:  `invalid: yaml: content:`,
			want:     false,
		},
		{
			name:     "file does not exist",
			filename: "nonexistent.yml",
			content:  "",
			want:     false,
		},
		{
			name:     "case insensitive extension",
			filename: "test.YAML",
			content: `apiVersion: v1
kind: Promptset
metadata:
  id: "test"
spec:
  prompts:
    prompt1:
      body: "test"`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testFile string
			if tt.content != "" && tt.filename != "nonexistent.yml" {
				// Create temp file
				tmpDir := t.TempDir()
				testFile = filepath.Join(tmpDir, tt.filename)
				err := os.WriteFile(testFile, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			} else if tt.filename == "nonexistent.yml" {
				// Use non-existent path
				testFile = filepath.Join(t.TempDir(), "nonexistent.yml")
			}

			got := IsPromptsetFile(testFile)
			if got != tt.want {
				t.Errorf("IsPromptsetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}