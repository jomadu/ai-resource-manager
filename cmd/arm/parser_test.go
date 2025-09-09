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
			name: "registry and ruleset only",
			arg:  "registry/ruleset",
			want: RulesetRef{Registry: "registry", Name: "ruleset", Version: ""},
		},
		{
			name: "with version",
			arg:  "registry/ruleset@1.0.0",
			want: RulesetRef{Registry: "registry", Name: "ruleset", Version: "1.0.0"},
		},
		{
			name: "with branch version",
			arg:  "registry/ruleset@main",
			want: RulesetRef{Registry: "registry", Name: "ruleset", Version: "main"},
		},
		{
			name:    "empty arg",
			arg:     "",
			wantErr: true,
		},
		{
			name:    "missing registry",
			arg:     "ruleset",
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
		{
			name: "complex registry name",
			arg:  "github.com/user/ruleset@v1.0.0",
			want: RulesetRef{Registry: "github.com", Name: "user/ruleset", Version: "v1.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRulesetArg(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRulesetArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
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
			name: "single arg",
			args: []string{"registry/ruleset"},
			want: []RulesetRef{
				{Registry: "registry", Name: "ruleset", Version: ""},
			},
		},
		{
			name: "multiple args",
			args: []string{"registry1/ruleset1", "registry2/ruleset2@1.0.0"},
			want: []RulesetRef{
				{Registry: "registry1", Name: "ruleset1", Version: ""},
				{Registry: "registry2", Name: "ruleset2", Version: "1.0.0"},
			},
		},
		{
			name: "empty args",
			args: []string{},
			want: []RulesetRef{},
		},
		{
			name:    "invalid arg",
			args:    []string{"registry/ruleset", "invalid"},
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
			if !reflect.DeepEqual(got, tt.want) {
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
			name:    "empty include",
			include: []string{},
			want:    []string{"**/*"},
		},
		{
			name:    "nil include",
			include: nil,
			want:    []string{"**/*"},
		},
		{
			name:    "existing patterns",
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
