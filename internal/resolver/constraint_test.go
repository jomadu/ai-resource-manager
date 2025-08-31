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
		{"Pin constraint", "1.0.0", Constraint{Type: Pin, Version: "1.0.0", Major: 1, Minor: 0, Patch: 0}, false},
		{"Caret constraint", "^1.0.0", Constraint{Type: Caret, Version: "1.0.0", Major: 1, Minor: 0, Patch: 0}, false},
		{"Tilde constraint", "~1.2.3", Constraint{Type: Tilde, Version: "1.2.3", Major: 1, Minor: 2, Patch: 3}, false},
		{"Branch constraint", "main", Constraint{Type: BranchHead, Version: "main"}, false},
		{"Latest constraint explicit", "latest", Constraint{Type: Latest}, false},
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
		{"Pin exact match", "1.0.0", Constraint{Type: Pin, Version: "1.0.0"}, true},
		{"Pin no match", "1.0.1", Constraint{Type: Pin, Version: "1.0.0"}, false},
		{"Caret compatible", "1.0.1", Constraint{Type: Caret, Version: "1.0.0"}, true},
		{"Caret incompatible", "2.0.0", Constraint{Type: Caret, Version: "1.0.0"}, false},
		{"Tilde compatible", "1.2.4", Constraint{Type: Tilde, Version: "1.2.3"}, true},
		{"Tilde incompatible", "1.3.0", Constraint{Type: Tilde, Version: "1.2.3"}, false},
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

	versions := []types.VersionRef{
		{ID: "1.0.0", Type: types.Tag},
		{ID: "1.0.1", Type: types.Tag},
		{ID: "1.1.0", Type: types.Tag},
		{ID: "2.0.0", Type: types.Tag},
		{ID: "main", Type: types.Branch},
	}

	tests := []struct {
		name       string
		constraint Constraint
		expected   string
		shouldFail bool
	}{
		{"Pin constraint", Constraint{Type: Pin, Version: "1.0.1"}, "1.0.1", false},
		{"Caret constraint", Constraint{Type: Caret, Version: "1.0.0"}, "1.1.0", false},
		{"Tilde constraint", Constraint{Type: Tilde, Version: "1.0.0"}, "1.0.1", false},
		{"Branch constraint", Constraint{Type: BranchHead, Version: "main"}, "main", false},
		{"Latest constraint", Constraint{Type: Latest}, "2.0.0", false},
		{"No match", Constraint{Type: Pin, Version: "3.0.0"}, "", true},
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
				if result == nil || result.ID != tt.expected {
					t.Errorf("Expected %s, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestGitConstraintResolver_ImplementsInterface(t *testing.T) {
	var _ ConstraintResolver = (*GitConstraintResolver)(nil)
}
