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

func TestParseConstraint_Caret(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ConstraintType
	}{
		{"^1.1.1 is major", "^1.1.1", Major},
		{"^1.0.0 is major", "^1.0.0", Major},
		{"^2.3.4 is major", "^2.3.4", Major},
		{"^0.1.0 is major", "^0.1.0", Major},
		{"^0.2.5 is major", "^0.2.5", Major},
		{"^0.0.1 is major", "^0.0.1", Major},
		{"^0.0.0 is major", "^0.0.0", Major},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseConstraint(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if c.Type != tt.wantType {
				t.Errorf("got type %v, want %v", c.Type, tt.wantType)
			}
		})
	}
}

func TestParseConstraint_Tilde(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"~1.0.1", "~1.0.1"},
		{"~1.2.3", "~1.2.3"},
		{"~0.0.1", "~0.0.1"},
		{"~2.5.0", "~2.5.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseConstraint(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if c.Type != Minor {
				t.Errorf("tilde should always be Minor, got %v", c.Type)
			}
		})
	}
}

func TestConstraint_IsSatisfiedBy(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		version    string
		want       bool
	}{
		// Caret always major
		{"^1.1.1 allows 1.1.1", "^1.1.1", "1.1.1", true},
		{"^1.1.1 allows 1.2.0", "^1.1.1", "1.2.0", true},
		{"^1.1.1 allows 1.9.9", "^1.1.1", "1.9.9", true},
		{"^1.1.1 rejects 2.0.0", "^1.1.1", "2.0.0", false},
		{"^1.1.1 rejects 1.1.0", "^1.1.1", "1.1.0", false},
		{"^1.1.1 rejects 1.0.0", "^1.1.1", "1.0.0", false},
		{"^0.1.0 allows 0.1.0", "^0.1.0", "0.1.0", true},
		{"^0.1.0 allows 0.9.9", "^0.1.0", "0.9.9", true},
		{"^0.1.0 rejects 1.0.0", "^0.1.0", "1.0.0", false},
		{"^0.0.1 allows 0.0.1", "^0.0.1", "0.0.1", true},
		{"^0.0.1 allows 0.0.9", "^0.0.1", "0.0.9", true},
		{"^0.0.1 allows 0.1.0", "^0.0.1", "0.1.0", true},
		{"^0.0.1 rejects 1.0.0", "^0.0.1", "1.0.0", false},

		// Tilde
		{"~1.0.1 allows 1.0.1", "~1.0.1", "1.0.1", true},
		{"~1.0.1 allows 1.0.9", "~1.0.1", "1.0.9", true},
		{"~1.0.1 rejects 1.1.0", "~1.0.1", "1.1.0", false},
		{"~1.0.1 rejects 1.0.0", "~1.0.1", "1.0.0", false},

		// Exact
		{"1.2.3 allows 1.2.3", "1.2.3", "1.2.3", true},
		{"1.2.3 rejects 1.2.4", "1.2.3", "1.2.4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseConstraint(tt.constraint)
			if err != nil {
				t.Fatal(err)
			}
			v, err := ParseVersion(tt.version)
			if err != nil {
				t.Fatal(err)
			}
			got := c.IsSatisfiedBy(v)
			if got != tt.want {
				t.Errorf("IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
