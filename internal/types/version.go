package types

import (
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Version represents a resolved version with display formatting.
type Version struct {
	Version string // Full version (e.g., "abc1234567890abcdef", "1.2.0")
	Display string // Display version (e.g., "abc1234", "1.2.0")
}

// ContentSelector defines include/exclude patterns for content filtering.
type ContentSelector struct {
	Include []string
	Exclude []string
}

// Matches returns true if the path matches the selector criteria.
// Both path and patterns are normalized to use forward slashes for
// consistent cross-platform behavior (Windows uses backslashes).
func (s ContentSelector) Matches(path string) bool {
	// Normalize path to use forward slashes for consistent matching across platforms
	normalizedPath := strings.ReplaceAll(path, "\\", "/")

	if len(s.Include) == 0 {
		return true
	}

	for _, pattern := range s.Include {
		// Normalize pattern to use forward slashes
		normalizedInclude := strings.ReplaceAll(pattern, "\\", "/")
		if matched, _ := doublestar.Match(normalizedInclude, normalizedPath); matched {
			for _, exclude := range s.Exclude {
				normalizedExclude := strings.ReplaceAll(exclude, "\\", "/")
				if matched, _ := doublestar.Match(normalizedExclude, normalizedPath); matched {
					return false
				}
			}
			return true
		}
	}

	return false
}
