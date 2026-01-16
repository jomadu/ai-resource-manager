package core

import "testing"

func TestParseConstraint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ConstraintType
		wantErr  bool
	}{
		{"latest", "latest", Latest, false},
		{"exact version", "1.2.3", Exact, false},
		{"minor version", "1.2.0", Minor, false},
		{"major version", "1.0.0", Major, false},
		{"branch name", "main", BranchHead, false},
		{"branch name with slash", "feature/test", BranchHead, false},
		{"v prefix exact", "v1.2.3", Exact, false},
		{"v prefix minor", "v1.2.0", Minor, false},
		{"v prefix major", "v1.0.0", Major, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConstraint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConstraint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Type != tt.wantType {
				t.Errorf("ParseConstraint() type = %v, want %v", got.Type, tt.wantType)
			}
		})
	}
}

func TestParseConstraint_VersionValues(t *testing.T) {
	t.Run("exact version has correct values", func(t *testing.T) {
		c, err := ParseConstraint("1.2.3")
		if err != nil {
			t.Fatal(err)
		}
		if c.Version.Major != 1 || c.Version.Minor != 2 || c.Version.Patch != 3 {
			t.Errorf("got %d.%d.%d, want 1.2.3", c.Version.Major, c.Version.Minor, c.Version.Patch)
		}
	})

	t.Run("minor version has correct values", func(t *testing.T) {
		c, err := ParseConstraint("2.5.0")
		if err != nil {
			t.Fatal(err)
		}
		if c.Version.Major != 2 || c.Version.Minor != 5 || c.Version.Patch != 0 {
			t.Errorf("got %d.%d.%d, want 2.5.0", c.Version.Major, c.Version.Minor, c.Version.Patch)
		}
	})

	t.Run("major version has correct values", func(t *testing.T) {
		c, err := ParseConstraint("3.0.0")
		if err != nil {
			t.Fatal(err)
		}
		if c.Version.Major != 3 || c.Version.Minor != 0 || c.Version.Patch != 0 {
			t.Errorf("got %d.%d.%d, want 3.0.0", c.Version.Major, c.Version.Minor, c.Version.Patch)
		}
	})

	t.Run("branch name preserves original", func(t *testing.T) {
		c, err := ParseConstraint("develop")
		if err != nil {
			t.Fatal(err)
		}
		if c.Version.Version != "develop" {
			t.Errorf("got %s, want develop", c.Version.Version)
		}
	})

	t.Run("latest has no version", func(t *testing.T) {
		c, err := ParseConstraint("latest")
		if err != nil {
			t.Fatal(err)
		}
		if c.Version != nil {
			t.Error("latest should have nil version")
		}
	})
}
