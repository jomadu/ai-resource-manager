package main

import (
	"reflect"
	"testing"
)

func TestParseRulesetArg(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    RulesetRef
		wantErr bool
	}{
		{
			name: "registry/ruleset format",
			arg:  "ai-rules/amazonq-rules",
			want: RulesetRef{
				Registry: "ai-rules",
				Name:     "amazonq-rules",
				Version:  "",
			},
			wantErr: false,
		},
		{
			name: "registry/ruleset@version format",
			arg:  "ai-rules/amazonq-rules@1.2.3",
			want: RulesetRef{
				Registry: "ai-rules",
				Name:     "amazonq-rules",
				Version:  "1.2.3",
			},
			wantErr: false,
		},
		{
			name: "registry/ruleset@latest format",
			arg:  "company/rules@latest",
			want: RulesetRef{
				Registry: "company",
				Name:     "rules",
				Version:  "latest",
			},
			wantErr: false,
		},
		{
			name:    "invalid format - no registry",
			arg:     "just-ruleset",
			wantErr: true,
		},
		{
			name:    "empty argument",
			arg:     "",
			wantErr: true,
		},
		{
			name:    "empty registry",
			arg:     "/ruleset",
			wantErr: true,
		},
		{
			name:    "empty ruleset",
			arg:     "registry/",
			wantErr: true,
		},
		{
			name:    "empty ruleset with version",
			arg:     "registry/@1.0.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRulesetArg(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRulesetArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRulesetArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRulesetArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []RulesetRef
		wantErr bool
	}{
		{
			name: "multiple valid rulesets",
			args: []string{"ai-rules/amazonq@1.0.0", "company/cursor-rules@latest"},
			want: []RulesetRef{
				{Registry: "ai-rules", Name: "amazonq", Version: "1.0.0"},
				{Registry: "company", Name: "cursor-rules", Version: "latest"},
			},
			wantErr: false,
		},
		{
			name: "single ruleset",
			args: []string{"ai-rules/test"},
			want: []RulesetRef{
				{Registry: "ai-rules", Name: "test", Version: ""},
			},
			wantErr: false,
		},
		{
			name:    "empty args",
			args:    []string{},
			want:    []RulesetRef{},
			wantErr: false,
		},
		{
			name:    "invalid ruleset in list",
			args:    []string{"ai-rules/valid", "invalid-format"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRulesetArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRulesetArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRulesetArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultIncludePatterns(t *testing.T) {
	tests := []struct {
		name    string
		include []string
		want    []string
	}{
		{
			name:    "nil include patterns",
			include: nil,
			want:    []string{"**/*"},
		},
		{
			name:    "empty include patterns",
			include: []string{},
			want:    []string{"**/*"},
		},
		{
			name:    "existing include patterns",
			include: []string{"*.md", "*.txt"},
			want:    []string{"*.md", "*.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDefaultIncludePatterns(tt.include)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDefaultIncludePatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}
