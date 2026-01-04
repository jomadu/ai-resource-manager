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