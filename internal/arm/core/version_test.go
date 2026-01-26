package core

import (
	"fmt"
	"testing"
)

func TestCompareTo(t *testing.T) {
	tests := []struct {
		name    string
		v1      Version
		v2      Version
		want    int
		wantErr bool
	}{
		{
			name:    "equal versions",
			v1:      Version{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3", IsSemver: true},
			v2:      Version{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3", IsSemver: true},
			want:    0,
			wantErr: false,
		},
		{
			name:    "v1 newer than v2",
			v1:      Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true},
			v2:      Version{Major: 1, Minor: 9, Patch: 9, Version: "1.9.9", IsSemver: true},
			want:    1,
			wantErr: false,
		},
		{
			name:    "v1 older than v2",
			v1:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			v2:      Version{Major: 1, Minor: 0, Patch: 1, Version: "1.0.1", IsSemver: true},
			want:    -1,
			wantErr: false,
		},
		{
			name:    "non-semver v1",
			v1:      Version{Version: "main", IsSemver: false},
			v2:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-semver v2",
			v1:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			v2:      Version{Version: "main", IsSemver: false},
			want:    0,
			wantErr: true,
		},
		{
			name:    "both non-semver",
			v1:      Version{Version: "main", IsSemver: false},
			v2:      Version{Version: "develop", IsSemver: false},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v1.CompareTo(&tt.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNewerThan(t *testing.T) {
	tests := []struct {
		name    string
		v1      Version
		v2      Version
		want    bool
		wantErr bool
	}{
		{
			name:    "v1 newer than v2",
			v1:      Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true},
			v2:      Version{Major: 1, Minor: 9, Patch: 9, Version: "1.9.9", IsSemver: true},
			want:    true,
			wantErr: false,
		},
		{
			name:    "v1 older than v2",
			v1:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			v2:      Version{Major: 1, Minor: 0, Patch: 1, Version: "1.0.1", IsSemver: true},
			want:    false,
			wantErr: false,
		},
		{
			name:    "equal versions",
			v1:      Version{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3", IsSemver: true},
			v2:      Version{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3", IsSemver: true},
			want:    false,
			wantErr: false,
		},
		{
			name:    "non-semver v1",
			v1:      Version{Version: "main", IsSemver: false},
			v2:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			want:    false,
			wantErr: true,
		},
		{
			name:    "non-semver v2",
			v1:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			v2:      Version{Version: "main", IsSemver: false},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v1.IsNewerThan(&tt.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsNewerThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsNewerThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOlderThan(t *testing.T) {
	tests := []struct {
		name    string
		v1      Version
		v2      Version
		want    bool
		wantErr bool
	}{
		{
			name:    "v1 older than v2",
			v1:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			v2:      Version{Major: 1, Minor: 0, Patch: 1, Version: "1.0.1", IsSemver: true},
			want:    true,
			wantErr: false,
		},
		{
			name:    "v1 newer than v2",
			v1:      Version{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0", IsSemver: true},
			v2:      Version{Major: 1, Minor: 9, Patch: 9, Version: "1.9.9", IsSemver: true},
			want:    false,
			wantErr: false,
		},
		{
			name:    "equal versions",
			v1:      Version{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3", IsSemver: true},
			v2:      Version{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3", IsSemver: true},
			want:    false,
			wantErr: false,
		},
		{
			name:    "non-semver v1",
			v1:      Version{Version: "main", IsSemver: false},
			v2:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			want:    false,
			wantErr: true,
		},
		{
			name:    "non-semver v2",
			v1:      Version{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0", IsSemver: true},
			v2:      Version{Version: "main", IsSemver: false},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v1.IsOlderThan(&tt.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsOlderThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsOlderThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		want    string
	}{
		{
			name:    "semver version",
			version: Version{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3", IsSemver: true},
			want:    "1.2.3",
		},
		{
			name:    "semver with v prefix",
			version: Version{Major: 1, Minor: 2, Patch: 3, Version: "v1.2.3", IsSemver: true},
			want:    "v1.2.3",
		},
		{
			name:    "non-semver version",
			version: Version{Version: "main", IsSemver: false},
			want:    "main",
		},
		{
			name:    "branch name",
			version: Version{Version: "feature/test", IsSemver: false},
			want:    "feature/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
		wantErr  bool
	}{
		// Empty string case
		{
			input:   "",
			wantErr: true,
		},
		// Valid semver cases
		{
			input: "1.2.3",
			expected: Version{
				Major:    1,
				Minor:    2,
				Patch:    3,
				Version:  "1.2.3",
				IsSemver: true,
			},
		},
		{
			input: "v1.2.3",
			expected: Version{
				Major:    1,
				Minor:    2,
				Patch:    3,
				Version:  "v1.2.3",
				IsSemver: true,
			},
		},
		{
			input: "1.2.3-alpha",
			expected: Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				Prerelease: "alpha",
				Version:    "1.2.3-alpha",
				IsSemver:   true,
			},
		},
		{
			input: "1.2.3+build123",
			expected: Version{
				Major:    1,
				Minor:    2,
				Patch:    3,
				Build:    "build123",
				Version:  "1.2.3+build123",
				IsSemver: true,
			},
		},
		{
			input: "v2.0.0-beta.1+exp.sha.5114f85",
			expected: Version{
				Major:      2,
				Minor:      0,
				Patch:      0,
				Prerelease: "beta.1",
				Build:      "exp.sha.5114f85",
				Version:    "v2.0.0-beta.1+exp.sha.5114f85",
				IsSemver:   true,
			},
		},
		// Abbreviated versions - now treated as non-semver
		{
			input: "1.2",
			expected: Version{
				Version:  "1.2",
				IsSemver: false,
			},
		},
		{
			input: "v1.2",
			expected: Version{
				Version:  "v1.2",
				IsSemver: false,
			},
		},
		{
			input: "5",
			expected: Version{
				Version:  "5",
				IsSemver: false,
			},
		},
		{
			input: "v3",
			expected: Version{
				Version:  "v3",
				IsSemver: false,
			},
		},
		// Non-semver cases - should return Version field only
		{
			input: "main",
			expected: Version{
				Version:  "main",
				IsSemver: false,
			},
		},
		{
			input: "develop",
			expected: Version{
				Version:  "develop",
				IsSemver: false,
			},
		},
		{
			input: "feature/new-stuff",
			expected: Version{
				Version:  "feature/new-stuff",
				IsSemver: false,
			},
		},
		{
			input: "1.2.x",
			expected: Version{
				Version:  "1.2.x",
				IsSemver: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseVersion(test.input)
			if test.wantErr {
				if err == nil {
					t.Errorf("ParseVersion(%q) expected error but got none", test.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseVersion(%q) returned error: %v", test.input, err)
				return
			}

			if result.Major != test.expected.Major {
				t.Errorf("Major: got %d, want %d", result.Major, test.expected.Major)
			}
			if result.Minor != test.expected.Minor {
				t.Errorf("Minor: got %d, want %d", result.Minor, test.expected.Minor)
			}
			if result.Patch != test.expected.Patch {
				t.Errorf("Patch: got %d, want %d", result.Patch, test.expected.Patch)
			}
			if result.Prerelease != test.expected.Prerelease {
				t.Errorf("Prerelease: got %q, want %q", result.Prerelease, test.expected.Prerelease)
			}
			if result.Build != test.expected.Build {
				t.Errorf("Build: got %q, want %q", result.Build, test.expected.Build)
			}
			if result.Version != test.expected.Version {
				t.Errorf("Version: got %q, want %q", result.Version, test.expected.Version)
			}
			if result.IsSemver != test.expected.IsSemver {
				t.Errorf("IsSemver: got %v, want %v", result.IsSemver, test.expected.IsSemver)
			}
		})
	}
}

// mustVersion is a test helper that creates a Version or panics
func mustVersion(s string) Version {
	v, err := NewVersion(s)
	if err != nil {
		panic(fmt.Sprintf("mustVersion(%q): %v", s, err))
	}
	return v
}

func TestResolveVersion_Normal(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("1.2.3"),
			mustVersion("2.0.0"),
		}
		got, err := ResolveVersion("1.2.3", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.2.3" {
			t.Errorf("got %s, want 1.2.3", got.Version)
		}
	})

	t.Run("latest returns newest", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("1.2.3"),
			mustVersion("2.0.0"),
		}
		got, err := ResolveVersion("latest", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "2.0.0" {
			t.Errorf("got %s, want 2.0.0", got.Version)
		}
	})

	t.Run("major constraint returns newest in major", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("1.2.3"),
			mustVersion("1.5.0"),
			mustVersion("2.0.0"),
		}
		got, err := ResolveVersion("1", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.5.0" {
			t.Errorf("got %s, want 1.5.0", got.Version)
		}
	})

	t.Run("minor constraint returns newest in minor", func(t *testing.T) {
		available := []Version{
			mustVersion("1.2.0"),
			mustVersion("1.2.3"),
			mustVersion("1.2.5"),
			mustVersion("1.3.0"),
		}
		got, err := ResolveVersion("1.2", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.2.5" {
			t.Errorf("got %s, want 1.2.5", got.Version)
		}
	})

	t.Run("branch name exact match", func(t *testing.T) {
		available := []Version{
			mustVersion("main"),
			mustVersion("develop"),
			mustVersion("1.0.0"),
		}
		got, err := ResolveVersion("main", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "main" {
			t.Errorf("got %s, want main", got.Version)
		}
	})
}

func TestResolveVersion_Edge(t *testing.T) {
	t.Run("single version available", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
		}
		got, err := ResolveVersion("latest", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.0.0" {
			t.Errorf("got %s, want 1.0.0", got.Version)
		}
	})

	t.Run("no matching versions", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("2.0.0"),
		}
		_, err := ResolveVersion("3.0.0", available)
		if err == nil {
			t.Error("expected error for no matching versions")
		}
	})

	t.Run("major constraint with no matches", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("1.5.0"),
		}
		_, err := ResolveVersion("2.0.0", available)
		if err == nil {
			t.Error("expected error for no matching major version")
		}
	})

	t.Run("minor constraint with no matches", func(t *testing.T) {
		available := []Version{
			mustVersion("1.2.0"),
			mustVersion("1.2.5"),
		}
		_, err := ResolveVersion("1.3.0", available)
		if err == nil {
			t.Error("expected error for no matching minor version")
		}
	})

	t.Run("branch name not found", func(t *testing.T) {
		available := []Version{
			mustVersion("main"),
			mustVersion("develop"),
		}
		_, err := ResolveVersion("feature/test", available)
		if err == nil {
			t.Error("expected error for non-existent branch")
		}
	})

	t.Run("mixed semantic and non-semantic versions", func(t *testing.T) {
		available := []Version{
			mustVersion("main"),
			mustVersion("1.0.0"),
			mustVersion("1.2.0"),
			mustVersion("develop"),
		}
		got, err := ResolveVersion("1", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.2.0" {
			t.Errorf("got %s, want 1.2.0", got.Version)
		}
	})

	t.Run("versions with same major.minor.patch", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("v1.0.0"),
		}
		got, err := ResolveVersion("1.0.0", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Major != 1 || got.Minor != 0 || got.Patch != 0 {
			t.Errorf("got %d.%d.%d, want 1.0.0", got.Major, got.Minor, got.Patch)
		}
	})
}

func TestResolveVersion_Extreme(t *testing.T) {
	t.Run("empty available versions", func(t *testing.T) {
		available := []Version{}
		_, err := ResolveVersion("latest", available)
		if err == nil {
			t.Error("expected error for empty versions list")
		}
	})

	t.Run("many versions", func(t *testing.T) {
		available := make([]Version, 100)
		for i := 0; i < 100; i++ {
			v := fmt.Sprintf("%d.%d.0", i/10, i%10)
			available[i] = mustVersion(v)
		}
		got, err := ResolveVersion("latest", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Major != 9 || got.Minor != 9 {
			t.Errorf("got %d.%d.0, want 9.9.0", got.Major, got.Minor)
		}
	})

	t.Run("large version numbers", func(t *testing.T) {
		available := []Version{
			mustVersion("999.999.999"),
			mustVersion("1000.0.0"),
		}
		got, err := ResolveVersion("latest", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Major != 1000 {
			t.Errorf("got %d.0.0, want 1000.0.0", got.Major)
		}
	})

	t.Run("all versions match constraint", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("1.1.0"),
			mustVersion("1.2.0"),
			mustVersion("1.3.0"),
		}
		got, err := ResolveVersion("1", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.3.0" {
			t.Errorf("got %s, want 1.3.0 (newest)", got.Version)
		}
	})

	t.Run("unsorted available versions", func(t *testing.T) {
		available := []Version{
			mustVersion("2.0.0"),
			mustVersion("1.0.0"),
			mustVersion("3.0.0"),
			mustVersion("1.5.0"),
		}
		got, err := ResolveVersion("latest", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "3.0.0" {
			t.Errorf("got %s, want 3.0.0", got.Version)
		}
	})

	t.Run("duplicate versions", func(t *testing.T) {
		available := []Version{
			mustVersion("1.0.0"),
			mustVersion("1.0.0"),
			mustVersion("1.0.0"),
		}
		got, err := ResolveVersion("1.0.0", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.0.0" {
			t.Errorf("got %s, want 1.0.0", got.Version)
		}
	})

	t.Run("only non-semantic versions", func(t *testing.T) {
		available := []Version{
			mustVersion("main"),
			mustVersion("develop"),
			mustVersion("feature/test"),
		}
		got, err := ResolveVersion("main", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "main" {
			t.Errorf("got %s, want main", got.Version)
		}
	})

	t.Run("constraint older than all available", func(t *testing.T) {
		available := []Version{
			mustVersion("2.0.0"),
			mustVersion("3.0.0"),
		}
		_, err := ResolveVersion("1.0.0", available)
		if err == nil {
			t.Error("expected error when constraint older than all versions")
		}
	})
}

func TestPrereleaseComparison(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int // -1 if v1 < v2, 0 if equal, 1 if v1 > v2
	}{
		// Basic prerelease vs release
		{name: "prerelease < release", v1: "1.0.0-alpha", v2: "1.0.0", want: -1},
		{name: "release > prerelease", v1: "1.0.0", v2: "1.0.0-alpha", want: 1},

		// Numeric identifier comparison
		{name: "numeric identifiers", v1: "1.0.0-1", v2: "1.0.0-2", want: -1},
		{name: "numeric identifiers equal", v1: "1.0.0-5", v2: "1.0.0-5", want: 0},
		{name: "numeric identifiers reverse", v1: "1.0.0-10", v2: "1.0.0-2", want: 1},

		// Alphanumeric identifier comparison
		{name: "alpha < beta", v1: "1.0.0-alpha", v2: "1.0.0-beta", want: -1},
		{name: "beta > alpha", v1: "1.0.0-beta", v2: "1.0.0-alpha", want: 1},
		{name: "alpha equal", v1: "1.0.0-alpha", v2: "1.0.0-alpha", want: 0},

		// Numeric < alphanumeric
		{name: "numeric < alphanumeric", v1: "1.0.0-1", v2: "1.0.0-alpha", want: -1},
		{name: "alphanumeric > numeric", v1: "1.0.0-alpha", v2: "1.0.0-1", want: 1},

		// Multi-part prerelease
		{name: "alpha.1 < alpha.2", v1: "1.0.0-alpha.1", v2: "1.0.0-alpha.2", want: -1},
		{name: "alpha.1 < alpha.beta", v1: "1.0.0-alpha.1", v2: "1.0.0-alpha.beta", want: -1},
		{name: "alpha < alpha.1", v1: "1.0.0-alpha", v2: "1.0.0-alpha.1", want: -1},
		{name: "alpha.1 > alpha", v1: "1.0.0-alpha.1", v2: "1.0.0-alpha", want: 1},

		// Complex semver examples from spec
		{name: "spec example 1", v1: "1.0.0-alpha", v2: "1.0.0-alpha.1", want: -1},
		{name: "spec example 2", v1: "1.0.0-alpha.1", v2: "1.0.0-alpha.beta", want: -1},
		{name: "spec example 3", v1: "1.0.0-alpha.beta", v2: "1.0.0-beta", want: -1},
		{name: "spec example 4", v1: "1.0.0-beta", v2: "1.0.0-beta.2", want: -1},
		{name: "spec example 5", v1: "1.0.0-beta.2", v2: "1.0.0-beta.11", want: -1},
		{name: "spec example 6", v1: "1.0.0-beta.11", v2: "1.0.0-rc.1", want: -1},
		{name: "spec example 7", v1: "1.0.0-rc.1", v2: "1.0.0", want: -1},

		// Edge cases
		{name: "empty prerelease equal", v1: "1.0.0", v2: "1.0.0", want: 0},
		{name: "leading zeros not numeric", v1: "1.0.0-01", v2: "1.0.0-1", want: 1}, // "01" is alphanumeric
		{name: "complex identifiers", v1: "1.0.0-x.7.z.92", v2: "1.0.0-x.7.z.93", want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("failed to parse v1 %s: %v", tt.v1, err)
			}
			v2, err := NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("failed to parse v2 %s: %v", tt.v2, err)
			}

			got := v1.Compare(&v2)
			if got != tt.want {
				t.Errorf("Compare(%s, %s) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}

			// Verify symmetry: if v1 < v2, then v2 > v1
			reverse := v2.Compare(&v1)
			if reverse != -got {
				t.Errorf("Symmetry broken: Compare(%s, %s) = %d, but Compare(%s, %s) = %d (expected %d)",
					tt.v1, tt.v2, got, tt.v2, tt.v1, reverse, -got)
			}
		})
	}
}

func TestSplitPrerelease(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: "", want: nil},
		{input: "alpha", want: []string{"alpha"}},
		{input: "alpha.1", want: []string{"alpha", "1"}},
		{input: "alpha.beta.1", want: []string{"alpha", "beta", "1"}},
		{input: "1.2.3", want: []string{"1", "2", "3"}},
		{input: "x.7.z.92", want: []string{"x", "7", "z", "92"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitPrerelease(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("splitPrerelease(%q) length = %d, want %d", tt.input, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("splitPrerelease(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseNumericIdentifier(t *testing.T) {
	tests := []struct {
		input     string
		wantVal   int
		wantIsNum bool
	}{
		{input: "0", wantVal: 0, wantIsNum: true},
		{input: "1", wantVal: 1, wantIsNum: true},
		{input: "123", wantVal: 123, wantIsNum: true},
		{input: "01", wantVal: 0, wantIsNum: false},  // leading zero
		{input: "001", wantVal: 0, wantIsNum: false}, // leading zeros
		{input: "alpha", wantVal: 0, wantIsNum: false},
		{input: "1alpha", wantVal: 0, wantIsNum: false},
		{input: "", wantVal: 0, wantIsNum: false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotVal, gotIsNum := parseNumericIdentifier(tt.input)
			if gotVal != tt.wantVal || gotIsNum != tt.wantIsNum {
				t.Errorf("parseNumericIdentifier(%q) = (%d, %v), want (%d, %v)",
					tt.input, gotVal, gotIsNum, tt.wantVal, tt.wantIsNum)
			}
		})
	}
}
