package types

import "testing"

func TestContentSelector_Matches_WindowsPathNormalization(t *testing.T) {
	tests := []struct {
		name     string
		selector ContentSelector
		path     string
		want     bool
	}{
		{
			name: "forward slash pattern matches backslash path",
			selector: ContentSelector{
				Include: []string{"**/*.yml"},
			},
			path: "rules\\my-rule.yml",
			want: true,
		},
		{
			name: "backslash pattern matches forward slash path",
			selector: ContentSelector{
				Include: []string{"**\\*.yml"},
			},
			path: "rules/my-rule.yml",
			want: true,
		},
		{
			name: "forward slash pattern matches forward slash path",
			selector: ContentSelector{
				Include: []string{"**/*.yml"},
			},
			path: "rules/my-rule.yml",
			want: true,
		},
		{
			name: "backslash pattern matches backslash path",
			selector: ContentSelector{
				Include: []string{"**\\*.yml"},
			},
			path: "rules\\my-rule.yml",
			want: true,
		},
		{
			name: "complex nested path with backslashes",
			selector: ContentSelector{
				Include: []string{"**/clean-code/**/*.yml"},
			},
			path: "rules\\clean-code\\naming\\use-meaningful-names.yml",
			want: true,
		},
		{
			name: "exclude pattern with backslashes",
			selector: ContentSelector{
				Include: []string{"**/*.yml"},
				Exclude: []string{"**/tests/**"},
			},
			path: "rules\\tests\\test-rule.yml",
			want: false,
		},
		{
			name: "exclude pattern with forward slashes matches backslash path",
			selector: ContentSelector{
				Include: []string{"**/*.yml"},
				Exclude: []string{"**/tests/**"},
			},
			path: "rules\\tests\\test-rule.yml",
			want: false,
		},
		{
			name: "non-matching extension",
			selector: ContentSelector{
				Include: []string{"**/*.yml"},
			},
			path: "rules\\my-rule.md",
			want: false,
		},
		{
			name: "empty include matches all",
			selector: ContentSelector{
				Include: []string{},
			},
			path: "any\\path\\file.txt",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.selector.Matches(tt.path); got != tt.want {
				t.Errorf("ContentSelector.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentSelector_Matches_CrossPlatformConsistency(t *testing.T) {
	selector := ContentSelector{
		Include: []string{"**/*.yml", "**/*.yaml"},
		Exclude: []string{"**/node_modules/**", "**/test/**"},
	}

	// Test that the same logical path works regardless of separator
	unixPath := "src/rules/clean-code.yml"
	windowsPath := "src\\rules\\clean-code.yml"

	unixResult := selector.Matches(unixPath)
	windowsResult := selector.Matches(windowsPath)

	if unixResult != windowsResult {
		t.Errorf("Cross-platform inconsistency: Unix path = %v, Windows path = %v", unixResult, windowsResult)
	}

	if !unixResult {
		t.Error("Expected both paths to match, but they didn't")
	}
}
