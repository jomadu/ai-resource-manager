package resolver

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestParseConstraint(t *testing.T) {
	resolver := NewGitConstraintResolver()

	tests := []struct {
		name       string
		constraint string
		want       Constraint
		wantErr    bool
	}{
		{
			name:       "exact version",
			constraint: "1.2.3",
			want:       Constraint{Type: Exact, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:       "caret constraint",
			constraint: "^1.2.3",
			want:       Constraint{Type: Major, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:       "tilde constraint",
			constraint: "~1.2.3",
			want:       Constraint{Type: Minor, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:       "latest",
			constraint: "latest",
			want:       Constraint{Type: Latest},
		},
		{
			name:       "branch",
			constraint: "main",
			want:       Constraint{Type: BranchHead, Version: "main"},
		},
		{
			name:       "shorthand major",
			constraint: "1",
			want:       Constraint{Type: Major, Version: "1.0.0", Major: 1, Minor: 0, Patch: 0},
		},
		{
			name:       "shorthand minor",
			constraint: "1.2",
			want:       Constraint{Type: Major, Version: "1.2.0", Major: 1, Minor: 2, Patch: 0},
		},
		{
			name:       "empty constraint",
			constraint: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.ParseConstraint(tt.constraint)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConstraint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseConstraint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSatisfiesConstraint(t *testing.T) {
	resolver := NewGitConstraintResolver()

	tests := []struct {
		name       string
		version    string
		constraint Constraint
		want       bool
	}{
		{
			name:       "exact match",
			version:    "1.2.3",
			constraint: Constraint{Type: Exact, Version: "1.2.3"},
			want:       true,
		},
		{
			name:       "exact no match",
			version:    "1.2.4",
			constraint: Constraint{Type: Exact, Version: "1.2.3"},
			want:       false,
		},
		{
			name:       "caret compatible",
			version:    "1.3.0",
			constraint: Constraint{Type: Major, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
			want:       true,
		},
		{
			name:       "caret incompatible major",
			version:    "2.0.0",
			constraint: Constraint{Type: Major, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
			want:       false,
		},
		{
			name:       "tilde compatible",
			version:    "1.2.5",
			constraint: Constraint{Type: Minor, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
			want:       true,
		},
		{
			name:       "tilde incompatible minor",
			version:    "1.3.0",
			constraint: Constraint{Type: Minor, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
			want:       false,
		},
		{
			name:       "branch match",
			version:    "main",
			constraint: Constraint{Type: BranchHead, Version: "main"},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolver.SatisfiesConstraint(tt.version, tt.constraint)
			if got != tt.want {
				t.Errorf("SatisfiesConstraint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindBestMatch(t *testing.T) {
	resolver := NewGitConstraintResolver()
	versions := []types.Version{
		{Version: "1.3.0", Display: "1.3.0"},
		{Version: "1.2.5", Display: "1.2.5"},
		{Version: "1.2.3", Display: "1.2.3"},
		{Version: "2.0.0", Display: "2.0.0"},
	}

	tests := []struct {
		name       string
		constraint Constraint
		want       string
		wantErr    bool
	}{
		{
			name:       "latest",
			constraint: Constraint{Type: Latest},
			want:       "1.3.0",
		},
		{
			name:       "caret finds highest compatible",
			constraint: Constraint{Type: Major, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
			want:       "1.3.0",
		},
		{
			name:       "tilde finds highest patch",
			constraint: Constraint{Type: Minor, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3},
			want:       "1.2.5",
		},
		{
			name:       "no match",
			constraint: Constraint{Type: Exact, Version: "3.0.0"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.FindBestMatch(tt.constraint, versions)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBestMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Version != tt.want {
				t.Errorf("FindBestMatch() = %v, want %v", got.Version, tt.want)
			}
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"with v prefix", "v1.2.3", "1.2.3"},
		{"without v prefix", "1.2.3", "1.2.3"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeVersion(tt.version)
			if got != tt.want {
				t.Errorf("normalizeVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		wantMajor int
		wantMinor int
		wantPatch int
		wantErr   bool
	}{
		{"valid version", "1.2.3", 1, 2, 3, false},
		{"with v prefix", "v1.2.3", 1, 2, 3, false},
		{"invalid format", "1.2", 0, 0, 0, true},
		{"non-numeric", "a.b.c", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			major, minor, patch, err := parseVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if major != tt.wantMajor || minor != tt.wantMinor || patch != tt.wantPatch {
				t.Errorf("parseVersion() = (%v, %v, %v), want (%v, %v, %v)",
					major, minor, patch, tt.wantMajor, tt.wantMinor, tt.wantPatch)
			}
		})
	}
}

func TestIsHigherVersion(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want bool
	}{
		{"higher major", "2.0.0", "1.9.9", true},
		{"higher minor", "1.2.0", "1.1.9", true},
		{"higher patch", "1.1.2", "1.1.1", true},
		{"equal", "1.1.1", "1.1.1", false},
		{"lower", "1.0.0", "1.1.0", false},
		{"invalid v1", "invalid", "1.0.0", false},
		{"invalid v2", "1.0.0", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isHigherVersion(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("isHigherVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandVersionShorthand(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		want       string
	}{
		{"major only", "1", "^1.0.0"},
		{"major.minor", "1.2", "^1.2.0"},
		{"full version", "1.2.3", "1.2.3"},
		{"with prefix", "^1.2.3", "^1.2.3"},
		{"branch", "main", "main"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandVersionShorthand(tt.constraint)
			if got != tt.want {
				t.Errorf("expandVersionShorthand() = %v, want %v", got, tt.want)
			}
		})
	}
}
