package resolver

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestNewGitConstraintResolver(t *testing.T) {
	resolver := NewGitConstraintResolver()
	if resolver == nil {
		t.Error("Expected non-nil resolver")
	}
}

func TestGitConstraintResolver_ParseConstraint(t *testing.T) {
	resolver := NewGitConstraintResolver()

	tests := []struct {
		name       string
		input      string
		expected   Constraint
		shouldFail bool
	}{
		{"Exact constraint", "1.0.0", Constraint{Type: Exact, Version: "1.0.0", Major: 1, Minor: 0, Patch: 0}, false},
		{"Major constraint", "^1.0.0", Constraint{Type: Major, Version: "1.0.0", Major: 1, Minor: 0, Patch: 0}, false},
		{"Minor constraint", "~1.2.3", Constraint{Type: Minor, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3}, false},
		{"Branch constraint", "main", Constraint{Type: BranchHead, Version: "main"}, false},
		{"Latest constraint explicit", "latest", Constraint{Type: Latest}, false},
		{"Major version shorthand", "1", Constraint{Type: Major, Version: "1.0.0", Major: 1, Minor: 0, Patch: 0}, false},
		{"Minor version shorthand", "1.2", Constraint{Type: Major, Version: "1.2.0", Major: 1, Minor: 2, Patch: 0}, false},
		{"Empty constraint", "", Constraint{}, true},
		{"Invalid constraint", "invalid", Constraint{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolver.ParseConstraint(tt.input)
			if tt.shouldFail {
				if err == nil {
					t.Error("Expected error for invalid constraint")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %+v, got %+v", tt.expected, result)
				}
			}
		})
	}
}

func TestGitConstraintResolver_SatisfiesConstraint(t *testing.T) {
	resolver := NewGitConstraintResolver()

	tests := []struct {
		name       string
		version    string
		constraint Constraint
		expected   bool
	}{
		{"Exact match", "1.0.0", Constraint{Type: Exact, Version: "1.0.0"}, true},
		{"Exact no match", "1.0.1", Constraint{Type: Exact, Version: "1.0.0"}, false},
		{"Major compatible", "1.0.1", Constraint{Type: Major, Version: "1.0.0"}, true},
		{"Major incompatible", "2.0.0", Constraint{Type: Major, Version: "1.0.0"}, false},
		{"Minor compatible", "1.2.4", Constraint{Type: Minor, Version: "1.2.3"}, true},
		{"Minor incompatible", "1.3.0", Constraint{Type: Minor, Version: "1.2.3"}, false},
		{"Branch match", "main", Constraint{Type: BranchHead, Version: "main"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.SatisfiesConstraint(tt.version, tt.constraint)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGitConstraintResolver_FindBestMatch(t *testing.T) {
	resolver := NewGitConstraintResolver()

	versions := []types.Version{
		{Version: "1.0.0", Display: "1.0.0"},
		{Version: "1.0.1", Display: "1.0.1"},
		{Version: "1.1.0", Display: "1.1.0"},
		{Version: "2.0.0", Display: "2.0.0"},
		{Version: "main", Display: "main"},
	}

	tests := []struct {
		name       string
		constraint Constraint
		expected   string
		shouldFail bool
	}{
		{"Exact constraint", Constraint{Type: Exact, Version: "1.0.1"}, "1.0.1", false},
		{"Major constraint", Constraint{Type: Major, Version: "1.0.0"}, "1.1.0", false},
		{"Minor constraint", Constraint{Type: Minor, Version: "1.0.0"}, "1.0.1", false},
		{"Branch constraint", Constraint{Type: BranchHead, Version: "main"}, "main", false},
		{"Latest constraint", Constraint{Type: Latest}, "1.0.0", false},
		{"No match", Constraint{Type: Exact, Version: "3.0.0"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolver.FindBestMatch(tt.constraint, versions)
			if tt.shouldFail {
				if err == nil {
					t.Error("Expected error for no match")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil || result.Version != tt.expected {
					t.Errorf("Expected %s, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestGitConstraintResolver_ImplementsInterface(t *testing.T) {
	var _ ConstraintResolver = (*GitConstraintResolver)(nil)
}
