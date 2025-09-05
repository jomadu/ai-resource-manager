package registry

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/resolver"
)

func TestGitRegistry_isBranchConstraint(t *testing.T) {
	registry := &GitRegistry{
		config: GitRegistryConfig{
			Branches: []string{"main", "develop"},
		},
	}

	tests := []struct {
		name       string
		constraint string
		expected   bool
	}{
		// Should be treated as branches (permissive)
		{"main branch", "main", true},
		{"develop branch", "develop", true},
		{"feature branch", "feature/new-feature", true},
		{"any branch name", "some-random-branch", true},
		{"branch with numbers", "release-2024", true},

		// Should NOT be treated as branches
		{"latest keyword", "latest", false},
		{"semantic version", "1.0.0", false},
		{"semantic version with v", "v1.0.0", false},
		{"caret constraint", "^1.0.0", false},
		{"tilde constraint", "~1.2.3", false},
		{"invalid dot pattern", "1.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := registry.isBranchConstraint(tt.constraint)
			if result != tt.expected {
				t.Errorf("isBranchConstraint(%s) = %v, expected %v", tt.constraint, result, tt.expected)
			}
		})
	}
}

func TestGitRegistry_BranchConstraintParsing(t *testing.T) {
	resolver := resolver.NewGitConstraintResolver()

	tests := []struct {
		name             string
		constraint       string
		expectedIsBranch bool
	}{
		{"main branch", "main", true},
		{"develop branch", "develop", true},
		{"feature branch", "feature/auth", true},
		{"release branch", "release-v1", true},
		{"semantic version", "1.0.0", false},
		{"caret constraint", "^1.0.0", false},
		{"tilde constraint", "~1.2.3", false},
		{"latest", "latest", false},
	}

	registry := &GitRegistry{
		config: GitRegistryConfig{
			Branches: []string{"main", "develop"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test constraint parsing
			constraint, err := resolver.ParseConstraint(tt.constraint)
			if err != nil {
				t.Errorf("ParseConstraint(%s) failed: %v", tt.constraint, err)
				return
			}

			// Test branch detection
			isBranch := registry.isBranchConstraint(tt.constraint)
			if isBranch != tt.expectedIsBranch {
				t.Errorf("isBranchConstraint(%s) = %v, expected %v", tt.constraint, isBranch, tt.expectedIsBranch)
			}

			// Verify constraint was parsed (basic validation)
			if constraint.Version == "" && tt.constraint != "latest" {
				t.Errorf("ParseConstraint(%s) resulted in empty version", tt.constraint)
			}
		})
	}
}
