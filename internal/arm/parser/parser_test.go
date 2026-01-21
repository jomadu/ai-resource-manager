package parser

import (
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
)

func TestParseRuleset(t *testing.T) {
	tests := []struct {
		name    string
		file    *core.File
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid ruleset",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
  name: "Test"
spec:
  rules:
    rule1:
      name: "Rule 1"
      body: "test rule"`),
			},
			wantErr: false,
		},
		{
			name: "invalid yaml",
			file: &core.File{
				Path:    "test.yml",
				Content: []byte(`invalid: yaml: content:`),
			},
			wantErr: true,
			errMsg:  "failed to parse YAML",
		},
		{
			name: "wrong kind",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Promptset
metadata:
  id: "test"
spec:
  prompts:
    prompt1:
      body: "test"`),
			},
			wantErr: true,
			errMsg:  "invalid ruleset",
		},
		{
			name: "missing required fields",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Ruleset`),
			},
			wantErr: true,
			errMsg:  "invalid ruleset",
		},
		{
			name: "empty rules",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
spec:
  rules: {}`),
			},
			wantErr: true,
			errMsg:  "invalid ruleset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRuleset(tt.file)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseRuleset() expected error but got none")
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseRuleset() error = %v, want error containing %v", err, tt.errMsg)
				}
				if result != nil {
					t.Errorf("ParseRuleset() expected nil result on error, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("ParseRuleset() unexpected error = %v", err)
					return
				}
				if result == nil {
					t.Errorf("ParseRuleset() expected result, got nil")
					return
				}
				if result.Kind != "Ruleset" {
					t.Errorf("ParseRuleset() kind = %v, want Ruleset", result.Kind)
				}
			}
		})
	}
}

func TestParsePromptset(t *testing.T) {
	tests := []struct {
		name    string
		file    *core.File
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid promptset",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Promptset
metadata:
  id: "test"
  name: "Test"
spec:
  prompts:
    prompt1:
      name: "Prompt 1"
      body: "test prompt"`),
			},
			wantErr: false,
		},
		{
			name: "invalid yaml",
			file: &core.File{
				Path:    "test.yml",
				Content: []byte(`invalid: yaml: content:`),
			},
			wantErr: true,
			errMsg:  "failed to parse YAML",
		},
		{
			name: "wrong kind",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Ruleset
metadata:
  id: "test"
spec:
  rules:
    rule1:
      body: "test"`),
			},
			wantErr: true,
			errMsg:  "invalid promptset",
		},
		{
			name: "missing required fields",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Promptset`),
			},
			wantErr: true,
			errMsg:  "invalid promptset",
		},
		{
			name: "empty prompts",
			file: &core.File{
				Path: "test.yml",
				Content: []byte(`apiVersion: v1
kind: Promptset
metadata:
  id: "test"
spec:
  prompts: {}`),
			},
			wantErr: true,
			errMsg:  "invalid promptset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePromptset(tt.file)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParsePromptset() expected error but got none")
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParsePromptset() error = %v, want error containing %v", err, tt.errMsg)
				}
				if result != nil {
					t.Errorf("ParsePromptset() expected nil result on error, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("ParsePromptset() unexpected error = %v", err)
					return
				}
				if result == nil {
					t.Errorf("ParsePromptset() expected result, got nil")
					return
				}
				if result.Kind != "Promptset" {
					t.Errorf("ParsePromptset() kind = %v, want Promptset", result.Kind)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   (len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		   (len(substr) < len(s) && findInString(s, substr))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}