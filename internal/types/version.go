package types

import "path/filepath"

// VersionRefType defines the type of version reference.
type VersionRefType int

const (
	Tag VersionRefType = iota
	Branch
	Commit
)

// VersionRef represents a version reference in a repository.
type VersionRef struct {
	ID   string         `json:"id"`
	Type VersionRefType `json:"type"`
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
		if matched, _ := filepath.Match(pattern, path); matched {
			for _, exclude := range s.Exclude {
				if matched, _ := filepath.Match(exclude, path); matched {
					return false
				}
			}
			return true
		}
	}

	return false
}
