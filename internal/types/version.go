package types

import "github.com/bmatcuk/doublestar/v4"

// Version represents a resolved version with display formatting.
type Version struct {
	Version string `json:"version"` // Full version (e.g., "abc1234567890abcdef", "1.2.0")
	Display string `json:"display"` // Display version (e.g., "abc1234", "1.2.0")
}

// ContentSelector defines include/exclude patterns for content filtering.
type ContentSelector struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// Matches returns true if the path matches the selector criteria.
func (s ContentSelector) Matches(path string) bool {
	if len(s.Include) == 0 {
		return true
	}

	for _, pattern := range s.Include {
		if matched, _ := doublestar.Match(pattern, path); matched {
			for _, exclude := range s.Exclude {
				if matched, _ := doublestar.Match(exclude, path); matched {
					return false
				}
			}
			return true
		}
	}

	return false
}
