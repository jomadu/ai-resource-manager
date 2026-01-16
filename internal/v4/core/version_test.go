package core

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
	}{
		// Valid semver cases
		{
			input: "1.2.3",
			expected: Version{
				Major:   1,
				Minor:   2,
				Patch:   3,
				Version: "1.2.3",
			},
		},
		{
			input: "v1.2.3",
			expected: Version{
				Major:   1,
				Minor:   2,
				Patch:   3,
				Version: "v1.2.3",
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
			},
		},
		{
			input: "1.2.3+build123",
			expected: Version{
				Major:   1,
				Minor:   2,
				Patch:   3,
				Build:   "build123",
				Version: "1.2.3+build123",
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
			},
		},
		// Partial version cases - should parse what's available
		{
			input: "1.2",
			expected: Version{
				Major:   1,
				Minor:   2,
				Version: "1.2",
			},
		},
		{
			input: "v1.2",
			expected: Version{
				Major:   1,
				Minor:   2,
				Version: "v1.2",
			},
		},
		{
			input: "5",
			expected: Version{
				Major:   5,
				Version: "5",
			},
		},
		{
			input: "v3",
			expected: Version{
				Major:   3,
				Version: "v3",
			},
		},
		// Non-semver cases - should return Version field only
		{
			input: "main",
			expected: Version{
				Version: "main",
			},
		},
		{
			input: "develop",
			expected: Version{
				Version: "develop",
			},
		},
		{
			input: "feature/new-stuff",
			expected: Version{
				Version: "feature/new-stuff",
			},
		},
		{
			input: "1.2.x",
			expected: Version{
				Version: "1.2.x",
			},
		},
		{
			input: "latest",
			expected: Version{
				Version: "latest",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseVersion(test.input)
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
		})
	}
}

func TestResolveVersion_Normal(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		available := []Version{
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3"},
			{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0"},
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
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3"},
			{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0"},
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
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3"},
			{Major: 1, Minor: 5, Patch: 0, Version: "1.5.0"},
			{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0"},
		}
		got, err := ResolveVersion("1.0.0", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.5.0" {
			t.Errorf("got %s, want 1.5.0", got.Version)
		}
	})

	t.Run("minor constraint returns newest in minor", func(t *testing.T) {
		available := []Version{
			{Major: 1, Minor: 2, Patch: 0, Version: "1.2.0"},
			{Major: 1, Minor: 2, Patch: 3, Version: "1.2.3"},
			{Major: 1, Minor: 2, Patch: 5, Version: "1.2.5"},
			{Major: 1, Minor: 3, Patch: 0, Version: "1.3.0"},
		}
		got, err := ResolveVersion("1.2.0", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.2.5" {
			t.Errorf("got %s, want 1.2.5", got.Version)
		}
	})

	t.Run("branch name exact match", func(t *testing.T) {
		available := []Version{
			{Version: "main"},
			{Version: "develop"},
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
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
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
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
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0"},
		}
		_, err := ResolveVersion("3.0.0", available)
		if err == nil {
			t.Error("expected error for no matching versions")
		}
	})

	t.Run("major constraint with no matches", func(t *testing.T) {
		available := []Version{
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 5, Patch: 0, Version: "1.5.0"},
		}
		_, err := ResolveVersion("2.0.0", available)
		if err == nil {
			t.Error("expected error for no matching major version")
		}
	})

	t.Run("minor constraint with no matches", func(t *testing.T) {
		available := []Version{
			{Major: 1, Minor: 2, Patch: 0, Version: "1.2.0"},
			{Major: 1, Minor: 2, Patch: 5, Version: "1.2.5"},
		}
		_, err := ResolveVersion("1.3.0", available)
		if err == nil {
			t.Error("expected error for no matching minor version")
		}
	})

	t.Run("branch name not found", func(t *testing.T) {
		available := []Version{
			{Version: "main"},
			{Version: "develop"},
		}
		_, err := ResolveVersion("feature/test", available)
		if err == nil {
			t.Error("expected error for non-existent branch")
		}
	})

	t.Run("mixed semantic and non-semantic versions", func(t *testing.T) {
		available := []Version{
			{Version: "main"},
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 2, Patch: 0, Version: "1.2.0"},
			{Version: "develop"},
		}
		got, err := ResolveVersion("1.0.0", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.2.0" {
			t.Errorf("got %s, want 1.2.0", got.Version)
		}
	})

	t.Run("versions with same major.minor.patch", func(t *testing.T) {
		available := []Version{
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 0, Patch: 0, Version: "v1.0.0"},
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
			available[i] = Version{
				Major:   i / 10,
				Minor:   i % 10,
				Patch:   0,
				Version: "",
			}
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
			{Major: 999, Minor: 999, Patch: 999, Version: "999.999.999"},
			{Major: 1000, Minor: 0, Patch: 0, Version: "1000.0.0"},
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
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 1, Patch: 0, Version: "1.1.0"},
			{Major: 1, Minor: 2, Patch: 0, Version: "1.2.0"},
			{Major: 1, Minor: 3, Patch: 0, Version: "1.3.0"},
		}
		got, err := ResolveVersion("1.0.0", available)
		if err != nil {
			t.Fatal(err)
		}
		if got.Version != "1.3.0" {
			t.Errorf("got %s, want 1.3.0 (newest)", got.Version)
		}
	})

	t.Run("unsorted available versions", func(t *testing.T) {
		available := []Version{
			{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0"},
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 3, Minor: 0, Patch: 0, Version: "3.0.0"},
			{Major: 1, Minor: 5, Patch: 0, Version: "1.5.0"},
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
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
			{Major: 1, Minor: 0, Patch: 0, Version: "1.0.0"},
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
			{Version: "main"},
			{Version: "develop"},
			{Version: "feature/test"},
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
			{Major: 2, Minor: 0, Patch: 0, Version: "2.0.0"},
			{Major: 3, Minor: 0, Patch: 0, Version: "3.0.0"},
		}
		_, err := ResolveVersion("1.0.0", available)
		if err == nil {
			t.Error("expected error when constraint older than all versions")
		}
	})
}
