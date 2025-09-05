package types

import "testing"

func TestContentSelector_Matches(t *testing.T) {
	tests := []struct {
		name     string
		selector ContentSelector
		path     string
		expected bool
	}{
		{
			name:     "Empty include matches all",
			selector: ContentSelector{},
			path:     "any/path.txt",
			expected: true,
		},
		{
			name:     "Simple pattern match",
			selector: ContentSelector{Include: []string{"*.md"}},
			path:     "file.md",
			expected: true,
		},
		{
			name:     "Simple pattern no match",
			selector: ContentSelector{Include: []string{"*.md"}},
			path:     "file.txt",
			expected: false,
		},
		{
			name:     "Globstar pattern match",
			selector: ContentSelector{Include: []string{"rules/**/*.md"}},
			path:     "rules/subdir/file.md",
			expected: true,
		},
		{
			name:     "Globstar deep match",
			selector: ContentSelector{Include: []string{"rules/**/*.md"}},
			path:     "rules/deep/nested/dir/file.md",
			expected: true,
		},
		{
			name:     "Globstar no match wrong extension",
			selector: ContentSelector{Include: []string{"rules/**/*.md"}},
			path:     "rules/subdir/file.txt",
			expected: false,
		},
		{
			name:     "Globstar no match wrong directory",
			selector: ContentSelector{Include: []string{"rules/**/*.md"}},
			path:     "other/subdir/file.md",
			expected: false,
		},
		{
			name:     "Multiple include patterns",
			selector: ContentSelector{Include: []string{"*.md", "*.txt"}},
			path:     "file.txt",
			expected: true,
		},
		{
			name:     "Exclude overrides include",
			selector: ContentSelector{Include: []string{"*.md"}, Exclude: []string{"test.md"}},
			path:     "test.md",
			expected: false,
		},
		{
			name:     "Include matches exclude doesn't",
			selector: ContentSelector{Include: []string{"*.md"}, Exclude: []string{"test.md"}},
			path:     "other.md",
			expected: true,
		},
		{
			name:     "Complex globstar with exclude",
			selector: ContentSelector{Include: []string{"src/**/*.go"}, Exclude: []string{"**/*_test.go"}},
			path:     "src/main/app.go",
			expected: true,
		},
		{
			name:     "Complex globstar excluded",
			selector: ContentSelector{Include: []string{"src/**/*.go"}, Exclude: []string{"**/*_test.go"}},
			path:     "src/main/app_test.go",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.selector.Matches(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for path %s", tt.expected, result, tt.path)
			}
		})
	}
}
