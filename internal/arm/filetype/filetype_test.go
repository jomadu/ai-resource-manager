package filetype

import (
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
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
			file := &core.File{
				Path:    tt.filename,
				Content: []byte(tt.content),
			}

			got := IsRulesetFile(file)
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
			file := &core.File{
				Path:    tt.filename,
				Content: []byte(tt.content),
			}

			got := IsPromptsetFile(file)
			if got != tt.want {
				t.Errorf("IsPromptsetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsResourceFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		want     bool
	}{
		{
			name:     "ruleset file",
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
			name:     "promptset file",
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
			name:     "regular yaml file",
			filename: "config.yml",
			content:  `database:
  host: localhost
  port: 5432`,
			want:     false,
		},
		{
			name:     "non-yaml file",
			filename: "readme.md",
			content:  `# README
This is a readme file.`,
			want:     false,
		},
		{
			name:     "invalid yaml",
			filename: "invalid.yml",
			content:  `invalid: yaml: content:`,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &core.File{
				Path:    tt.filename,
				Content: []byte(tt.content),
			}

			got := IsResourceFile(file)
			if got != tt.want {
				t.Errorf("IsResourceFile() = %v, want %v", got, tt.want)
			}
		})
	}
}