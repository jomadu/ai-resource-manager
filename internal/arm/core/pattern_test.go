package core

import "testing"

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		// Literal matches
		{"exact match", "file.yml", "file.yml", true},
		{"no match", "file.yml", "other.yml", false},
		{"path match", "dir/file.yml", "dir/file.yml", true},
		{"path no match", "dir/file.yml", "other/file.yml", false},

		// Single wildcard
		{"wildcard prefix", "*.yml", "file.yml", true},
		{"wildcard prefix no match", "*.yml", "file.txt", false},
		{"wildcard suffix", "file.*", "file.yml", true},
		{"wildcard suffix no match", "file.*", "other.yml", false},
		{"wildcard middle", "file*.yml", "file123.yml", true},
		{"wildcard middle no match", "file*.yml", "other123.yml", false},

		// Double star prefix (**/pattern)
		{"doublestar prefix match", "**/file.yml", "file.yml", true},
		{"doublestar prefix match deep", "**/file.yml", "dir/file.yml", true},
		{"doublestar prefix match deeper", "**/file.yml", "dir/subdir/file.yml", true},
		{"doublestar prefix no match", "**/file.yml", "other.yml", false},
		{"doublestar with wildcard", "**/*.yml", "file.yml", true},
		{"doublestar with wildcard deep", "**/*.yml", "dir/file.yml", true},
		{"doublestar with wildcard no match", "**/*.yml", "file.txt", false},

		// Double star suffix (dir/**)
		{"doublestar suffix match", "dir/**", "dir/file.yml", true},
		{"doublestar suffix match deep", "dir/**", "dir/subdir/file.yml", true},
		{"doublestar suffix match exact", "dir/**", "dir", true},
		{"doublestar suffix no match", "dir/**", "other/file.yml", false},

		// Complex patterns
		{"security prefix", "security/**/*.yml", "security/rule1.yml", true},
		{"security prefix exact match", "security/**/*.yml", "security/rule1.yml", true},
		{"security prefix rule2", "security/**/*.yml", "security/rule2.yml", true},
		{"security prefix deep", "security/**/*.yml", "security/subdir/rule1.yml", true},
		{"security prefix no match", "security/**/*.yml", "general/rule1.yml", false},
		{"experimental exclude", "**/experimental/**", "experimental/rule.yml", true},
		{"experimental exclude deep", "**/experimental/**", "dir/experimental/rule.yml", true},
		{"experimental exclude no match", "**/experimental/**", "dir/rule.yml", false},

		// Path separator normalization
		{"backslash pattern", "dir\\file.yml", "dir/file.yml", true},
		{"backslash path", "dir/file.yml", "dir\\file.yml", true},
		{"mixed separators", "dir\\**/*.yml", "dir/subdir/file.yml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchPattern(tt.pattern, tt.path)
			if got != tt.want {
				t.Errorf("MatchPattern(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
			}
		})
	}
}
