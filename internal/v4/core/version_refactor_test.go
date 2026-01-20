package core

import "testing"

// TestParseVersion_RefactoredBehavior tests new strict ParseVersion
func TestParseVersion_RefactoredBehavior(t *testing.T) {
	t.Run("treats constraint strings as non-semver", func(t *testing.T) {
		constraints := []string{"^1.2.3", "~1.2.3", "latest"}
		for _, input := range constraints {
			v, err := ParseVersion(input)
			if err != nil {
				t.Errorf("ParseVersion(%q) should not error, got %v", input, err)
			}
			if v.IsSemver {
				t.Errorf("ParseVersion(%q) should not be semver", input)
			}
			if v.Version != input {
				t.Errorf("ParseVersion(%q) Version=%q, want %q", input, v.Version, input)
			}
		}
	})

	t.Run("treats abbreviated versions as non-semver", func(t *testing.T) {
		abbreviated := []string{"1.2", "1", "v1.2", "v1"}
		for _, input := range abbreviated {
			v, err := ParseVersion(input)
			if err != nil {
				t.Errorf("ParseVersion(%q) should not error, got %v", input, err)
			}
			if v.IsSemver {
				t.Errorf("ParseVersion(%q) should not be semver", input)
			}
			if v.Version != input {
				t.Errorf("ParseVersion(%q) Version=%q, want %q", input, v.Version, input)
			}
		}
	})

	t.Run("accepts full semver", func(t *testing.T) {
		valid := []string{"1.2.3", "v1.2.3", "0.0.1", "v0.0.0"}
		for _, input := range valid {
			v, err := ParseVersion(input)
			if err != nil {
				t.Errorf("ParseVersion(%q) should not error, got %v", input, err)
			}
			if !v.IsSemver {
				t.Errorf("ParseVersion(%q) should be semver", input)
			}
		}
	})

	t.Run("accepts non-semver strings", func(t *testing.T) {
		valid := []string{"main", "develop", "feature/test"}
		for _, input := range valid {
			v, err := ParseVersion(input)
			if err != nil {
				t.Errorf("ParseVersion(%q) should not error, got %v", input, err)
			}
			if v.IsSemver {
				t.Errorf("ParseVersion(%q) should not be semver", input)
			}
			if v.Version != input {
				t.Errorf("ParseVersion(%q) Version=%q, want %q", input, v.Version, input)
			}
		}
	})
}

// TestParseConstraint_RefactoredBehavior tests ParseConstraint handles abbreviations
func TestParseConstraint_RefactoredBehavior(t *testing.T) {
	t.Run("expands abbreviated versions", func(t *testing.T) {
		tests := []struct {
			input     string
			wantMajor int
			wantMinor int
			wantPatch int
			wantType  ConstraintType
		}{
			{"1.2", 1, 2, 0, Minor},
			{"1", 1, 0, 0, Major},
			{"v1.2", 1, 2, 0, Minor},
			{"v1", 1, 0, 0, Major},
		}
		for _, tt := range tests {
			c, err := ParseConstraint(tt.input)
			if err != nil {
				t.Errorf("ParseConstraint(%q) error: %v", tt.input, err)
				continue
			}
			if c.Version.Major != tt.wantMajor || c.Version.Minor != tt.wantMinor || c.Version.Patch != tt.wantPatch {
				t.Errorf("ParseConstraint(%q) = %d.%d.%d, want %d.%d.%d",
					tt.input, c.Version.Major, c.Version.Minor, c.Version.Patch,
					tt.wantMajor, tt.wantMinor, tt.wantPatch)
			}
			if c.Type != tt.wantType {
				t.Errorf("ParseConstraint(%q) type = %v, want %v", tt.input, c.Type, tt.wantType)
			}
		}
	})

	t.Run("expands with caret and tilde", func(t *testing.T) {
		tests := []struct {
			input     string
			wantMajor int
			wantMinor int
			wantPatch int
		}{
			{"^1.2", 1, 2, 0},
			{"^1", 1, 0, 0},
			{"~1.2", 1, 2, 0},
			{"~1", 1, 0, 0},
		}
		for _, tt := range tests {
			c, err := ParseConstraint(tt.input)
			if err != nil {
				t.Errorf("ParseConstraint(%q) error: %v", tt.input, err)
				continue
			}
			if c.Version.Major != tt.wantMajor || c.Version.Minor != tt.wantMinor || c.Version.Patch != tt.wantPatch {
				t.Errorf("ParseConstraint(%q) = %d.%d.%d, want %d.%d.%d",
					tt.input, c.Version.Major, c.Version.Minor, c.Version.Patch,
					tt.wantMajor, tt.wantMinor, tt.wantPatch)
			}
		}
	})
}
