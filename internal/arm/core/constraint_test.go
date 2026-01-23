package core

import "testing"

func TestNewConstraint(t *testing.T) {
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
		{"branch name rejected", "main", Latest, true},
		{"branch name with slash rejected", "feature/test", Latest, true},
		{"v prefix exact", "v1.2.3", Exact, false},
		{"v prefix minor", "v1.2.0", Minor, false},
		{"v prefix major", "v1.0.0", Major, false},
		// Abbreviated versions
		{"abbreviated major", "1", Major, false},
		{"abbreviated minor", "1.2", Minor, false},
		{"abbreviated major with v", "v1", Major, false},
		{"abbreviated minor with v", "v1.2", Minor, false},
		// Caret
		{"caret exact", "^1.2.3", Major, false},
		{"caret with v", "^v1.2.3", Major, false},
		// Tilde
		{"tilde exact", "~1.2.3", Minor, false},
		{"tilde with v", "~v1.2.3", Minor, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConstraint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConstraint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Type != tt.wantType {
				t.Errorf("NewConstraint() type = %v, want %v", got.Type, tt.wantType)
			}
		})
	}
}

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
		{"branch name", "main", Latest, false},
		{"branch name with slash", "feature/test", Latest, false},
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
		wantErr    bool
	}{
		// Latest
		{"latest allows semver", "latest", "1.2.3", true, false},
		{"latest allows any semver", "latest", "0.0.1", true, false},

		// Exact
		{"1.2.3 allows 1.2.3", "1.2.3", "1.2.3", true, false},
		{"1.2.3 rejects 1.2.4", "1.2.3", "1.2.4", false, false},

		// Major (abbreviated like "1" becomes Major constraint)
		{"1 allows 1.0.0", "1", "1.0.0", true, false},
		{"1 allows 1.5.0", "1", "1.5.0", true, false},
		{"1 allows 1.9.9", "1", "1.9.9", true, false},
		{"1 rejects 2.0.0", "1", "2.0.0", false, false},
		{"1 rejects 0.9.9", "1", "0.9.9", false, false},

		// Minor (abbreviated like "1.2" becomes Minor constraint)
		{"1.2 allows 1.2.0", "1.2", "1.2.0", true, false},
		{"1.2 allows 1.2.5", "1.2", "1.2.5", true, false},
		{"1.2 allows 1.2.9", "1.2", "1.2.9", true, false},
		{"1.2 rejects 1.3.0", "1.2", "1.3.0", false, false},
		{"1.2 rejects 1.1.9", "1.2", "1.1.9", false, false},

		// Caret always major
		{"^1.1.1 allows 1.1.1", "^1.1.1", "1.1.1", true, false},
		{"^1.1.1 allows 1.2.0", "^1.1.1", "1.2.0", true, false},
		{"^1.1.1 allows 1.9.9", "^1.1.1", "1.9.9", true, false},
		{"^1.1.1 rejects 2.0.0", "^1.1.1", "2.0.0", false, false},
		{"^1.1.1 rejects 1.1.0", "^1.1.1", "1.1.0", false, false},
		{"^1.1.1 rejects 1.0.0", "^1.1.1", "1.0.0", false, false},
		{"^0.1.0 allows 0.1.0", "^0.1.0", "0.1.0", true, false},
		{"^0.1.0 allows 0.9.9", "^0.1.0", "0.9.9", true, false},
		{"^0.1.0 rejects 1.0.0", "^0.1.0", "1.0.0", false, false},
		{"^0.0.1 allows 0.0.1", "^0.0.1", "0.0.1", true, false},
		{"^0.0.1 allows 0.0.9", "^0.0.1", "0.0.9", true, false},
		{"^0.0.1 allows 0.1.0", "^0.0.1", "0.1.0", true, false},
		{"^0.0.1 rejects 1.0.0", "^0.0.1", "1.0.0", false, false},

		// Tilde
		{"~1.0.1 allows 1.0.1", "~1.0.1", "1.0.1", true, false},
		{"~1.0.1 allows 1.0.9", "~1.0.1", "1.0.9", true, false},
		{"~1.0.1 rejects 1.1.0", "~1.0.1", "1.1.0", false, false},
		{"~1.0.1 rejects 1.0.0", "~1.0.1", "1.0.0", false, false},
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
			got, err := c.IsSatisfiedBy(&v)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSatisfiedBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstraint_IsSatisfiedBy_Errors(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		version    string
	}{
		{"non-semver version with exact constraint", "1.2.3", "main"},
		{"non-semver version with major constraint", "1.0.0", "develop"},
		{"non-semver version with minor constraint", "1.2.0", "feature/test"},
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
			_, err = c.IsSatisfiedBy(&v)
			if err == nil {
				t.Error("expected error for non-semver version, got nil")
			}
		})
	}
}

func TestConstraint_ToString(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		want       string
	}{
		{"latest", "latest", "latest"},
		{"exact version", "1.2.3", "1.2.3"},
		{"minor version", "1.2.0", "~1.2.0"},
		{"major version", "1.0.0", "^1.0.0"},
		{"caret explicit", "^2.3.4", "^2.3.4"},
		{"tilde explicit", "~3.4.5", "~3.4.5"},
		// Test v prefix preservation
		{"v prefix exact", "v1.2.3", "v1.2.3"},
		{"v prefix with caret", "^v2.3.4", "^v2.3.4"},
		{"v prefix with tilde", "~v3.4.5", "~v3.4.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseConstraint(tt.constraint)
			if err != nil {
				t.Fatal(err)
			}
			got := c.ToString()
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConstraint_ToString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"latest", "latest", "latest"},
		{"exact version", "1.2.3", "1.2.3"},
		{"v prefix exact", "v1.2.3", "v1.2.3"},
		{"caret", "^1.2.3", "^1.2.3"},
		{"caret with v", "^v1.2.3", "^v1.2.3"},
		{"tilde", "~1.2.3", "~1.2.3"},
		{"tilde with v", "~v1.2.3", "~v1.2.3"},
		// Abbreviated versions expand but preserve v
		{"abbreviated major", "1", "^1.0.0"},
		{"abbreviated minor", "1.2", "~1.2.0"},
		{"abbreviated major with v", "v1", "^v1.0.0"},
		{"abbreviated minor with v", "v1.2", "~v1.2.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewConstraint(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			got := c.ToString()
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
