package urf

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestYAMLParser_IsURF(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		file     *types.File
		expected bool
	}{
		{
			name: "valid URF file",
			file: &types.File{
				Path: "test.yaml",
				Content: []byte(`version: "1.0"
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
  version: "1.0.0"
rules:
  - id: "test-rule"
    name: "Test Rule"
    body: "Test content"`),
			},
			expected: true,
		},
		{
			name: "non-yaml file",
			file: &types.File{
				Path:    "test.txt",
				Content: []byte("some content"),
			},
			expected: false,
		},
		{
			name: "yaml file but not URF",
			file: &types.File{
				Path:    "test.yaml",
				Content: []byte("key: value"),
			},
			expected: false,
		},
		{
			name: "invalid yaml",
			file: &types.File{
				Path:    "test.yaml",
				Content: []byte("invalid: yaml: content:"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsURF(tt.file)
			if result != tt.expected {
				t.Errorf("IsURF() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestYAMLParser_Parse(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name        string
		file        *types.File
		expectError bool
		validate    func(*testing.T, *URFFile)
	}{
		{
			name: "valid URF file",
			file: &types.File{
				Path: "test.yaml",
				Content: []byte(`version: "1.0"
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
  version: "1.0.0"
  description: "Test description"
rules:
  - id: "test-rule"
    name: "Test Rule"
    description: "Test rule description"
    priority: 100
    enforcement: "must"
    scope:
      - files: ["**/*.go"]
    body: "Test content"`),
			},
			expectError: false,
			validate: func(t *testing.T, urf *URFFile) {
				if urf.Version != "1.0" {
					t.Errorf("Version = %s, expected 1.0", urf.Version)
				}
				if urf.Metadata.ID != "test-ruleset" {
					t.Errorf("Metadata.ID = %s, expected test-ruleset", urf.Metadata.ID)
				}
				if len(urf.Rules) != 1 {
					t.Errorf("Rules length = %d, expected 1", len(urf.Rules))
				}
			},
		},
		{
			name: "missing version",
			file: &types.File{
				Path: "test.yaml",
				Content: []byte(`metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
  version: "1.0.0"
rules:
  - id: "test-rule"
    name: "Test Rule"
    body: "Test content"`),
			},
			expectError: true,
		},
		{
			name: "missing metadata",
			file: &types.File{
				Path: "test.yaml",
				Content: []byte(`version: "1.0"
rules:
  - id: "test-rule"
    name: "Test Rule"
    body: "Test content"`),
			},
			expectError: true,
		},
		{
			name: "duplicate rule IDs",
			file: &types.File{
				Path: "test.yaml",
				Content: []byte(`version: "1.0"
metadata:
  id: "test-ruleset"
  name: "Test Ruleset"
  version: "1.0.0"
rules:
  - id: "test-rule"
    name: "Test Rule"
    body: "Test content"
  - id: "test-rule"
    name: "Another Rule"
    body: "Another content"`),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.file)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
