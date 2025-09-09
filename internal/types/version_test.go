package types

import "testing"

func TestContentSelectorMatches(t *testing.T) {
	tests := []struct {
		name     string
		selector ContentSelector
		path     string
		want     bool
	}{
		{
			name:     "empty include matches all",
			selector: ContentSelector{},
			path:     "any/path.txt",
			want:     true,
		},
		{
			name:     "simple include match",
			selector: ContentSelector{Include: []string{"*.txt"}},
			path:     "file.txt",
			want:     true,
		},
		{
			name:     "simple include no match",
			selector: ContentSelector{Include: []string{"*.txt"}},
			path:     "file.md",
			want:     false,
		},
		{
			name:     "glob pattern match",
			selector: ContentSelector{Include: []string{"**/*.md"}},
			path:     "docs/readme.md",
			want:     true,
		},
		{
			name: "exclude overrides include",
			selector: ContentSelector{
				Include: []string{"**/*.md"},
				Exclude: []string{"**/private.md"},
			},
			path: "docs/private.md",
			want: false,
		},
		{
			name: "include with exclude no match",
			selector: ContentSelector{
				Include: []string{"**/*.md"},
				Exclude: []string{"**/private.md"},
			},
			path: "docs/readme.md",
			want: true,
		},
		{
			name:     "multiple includes",
			selector: ContentSelector{Include: []string{"*.txt", "*.md"}},
			path:     "file.md",
			want:     true,
		},
		{
			name: "multiple excludes",
			selector: ContentSelector{
				Include: []string{"**/*"},
				Exclude: []string{"*.tmp", "*.log"},
			},
			path: "file.tmp",
			want: false,
		},
		{
			name:     "nested path match",
			selector: ContentSelector{Include: []string{"src/**/*.go"}},
			path:     "src/internal/cache/key.go",
			want:     true,
		},
		{
			name:     "nested path no match",
			selector: ContentSelector{Include: []string{"src/**/*.go"}},
			path:     "test/cache_test.go",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.selector.Matches(tt.path)
			if got != tt.want {
				t.Errorf("ContentSelector.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
