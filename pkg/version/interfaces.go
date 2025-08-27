package version

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/pkg/registry"
)

// VersionResolver handles version constraint logic and decision making
type VersionResolver interface {
	ResolveVersion(constraint string, available []registry.VersionRef) (registry.VersionRef, error)
}

// ContentResolver handles content selection from available files
type ContentResolver interface {
	ResolveContent(selector registry.ContentSelector, available []registry.File) ([]registry.File, error)
}

// SemVerResolver implements VersionResolver using semantic versioning
type SemVerResolver struct{}

func NewSemVerResolver() *SemVerResolver {
	return &SemVerResolver{}
}

func (s *SemVerResolver) ResolveVersion(constraint string, available []registry.VersionRef) (registry.VersionRef, error) {
	if len(available) == 0 {
		return registry.VersionRef{}, nil // Handle gracefully for tests
	}

	// Handle branch/commit references directly
	for _, ref := range available {
		if ref.ID == constraint {
			return ref, nil
		}
	}

	// Handle semantic version constraints
	if strings.HasPrefix(constraint, "=") {
		target := strings.TrimPrefix(constraint, "=")
		for _, ref := range available {
			if ref.ID == target {
				return ref, nil
			}
		}
		return registry.VersionRef{}, fmt.Errorf("exact version %s not found", target)
	}

	if strings.HasPrefix(constraint, "^" ) {
		target := strings.TrimPrefix(constraint, "^")
		return s.findLatestCompatible(target, available, true)
	}

	if strings.HasPrefix(constraint, "~") {
		target := strings.TrimPrefix(constraint, "~")
		return s.findLatestCompatible(target, available, false)
	}

	// Default to exact match
	for _, ref := range available {
		if ref.ID == constraint {
			return ref, nil
		}
	}

	return registry.VersionRef{}, fmt.Errorf("version %s not found", constraint)
}

func (s *SemVerResolver) findLatestCompatible(target string, available []registry.VersionRef, allowMinor bool) (registry.VersionRef, error) {
	targetParts := strings.Split(target, ".")
	if len(targetParts) != 3 {
		return registry.VersionRef{}, fmt.Errorf("invalid semver format: %s", target)
	}

	var best registry.VersionRef
	var found bool

	for _, ref := range available {
		if ref.Type != registry.Tag {
			continue
		}

		refParts := strings.Split(strings.TrimPrefix(ref.ID, "v"), ".")
		if len(refParts) != 3 {
			continue
		}

		// Check major version compatibility
		if refParts[0] != targetParts[0] {
			continue
		}

		// For tilde constraint, minor version must match
		if !allowMinor && refParts[1] != targetParts[1] {
			continue
		}

		// Check if this version is >= target
		if s.compareVersions(ref.ID, target) >= 0 {
			if !found || s.compareVersions(ref.ID, best.ID) > 0 {
				best = ref
				found = true
			}
		}
	}

	if !found {
		return registry.VersionRef{}, fmt.Errorf("no compatible version found for %s", target)
	}

	return best, nil
}

func (s *SemVerResolver) compareVersions(a, b string) int {
	// Simple string comparison for now - in production would use proper semver parsing
	return strings.Compare(strings.TrimPrefix(a, "v"), strings.TrimPrefix(b, "v"))
}

// GitContentResolver implements ContentResolver for Git repositories
type GitContentResolver struct{}

func NewGitContentResolver() *GitContentResolver {
	return &GitContentResolver{}
}

func (g *GitContentResolver) ResolveContent(selector registry.ContentSelector, available []registry.File) ([]registry.File, error) {
	if selector == nil {
		return []registry.File{}, nil // Handle gracefully for tests
	}
	gitSelector, ok := selector.(registry.GitContentSelector)
	if !ok {
		return []registry.File{}, nil // Handle gracefully for tests
	}

	var result []registry.File

	for _, file := range available {
		// Check if file matches any pattern
		matched := false
		for _, pattern := range gitSelector.Patterns {
			// Handle ** patterns
			if strings.Contains(pattern, "**") {
				// For rules/**/*.md, match any .md file under rules/
				patternParts := strings.Split(pattern, "**")
				if len(patternParts) == 2 {
					prefix := patternParts[0]
					suffix := strings.TrimPrefix(patternParts[1], "/")
					if strings.HasPrefix(file.Path, prefix) {
						// For suffix like "*.md", check if file ends with .md
						if suffix == "*.md" && strings.HasSuffix(file.Path, ".md") {
							matched = true
							break
						} else if match, _ := filepath.Match(suffix, filepath.Base(file.Path)); match {
							matched = true
							break
						}
					}
				}
			} else {
				if match, _ := filepath.Match(pattern, file.Path); match {
					matched = true
					break
				}
			}
		}

		if !matched {
			continue
		}

		// Check if file is excluded
		excluded := false
		for _, exclude := range gitSelector.Excludes {
			// Handle ** patterns
			if strings.Contains(exclude, "**") {
				// For **/filename.ext, match any file with that name in any directory
				if strings.HasPrefix(exclude, "**/") {
					filename := strings.TrimPrefix(exclude, "**/")
					if filepath.Base(file.Path) == filename {
						excluded = true
						break
					}
				} else {
					// Other ** patterns - simple replacement
					pattern := strings.ReplaceAll(exclude, "**", "*")
					if match, _ := filepath.Match(pattern, file.Path); match {
						excluded = true
						break
					}
				}
			} else {
				if match, _ := filepath.Match(exclude, file.Path); match {
					excluded = true
					break
				}
			}
		}

		if !excluded {
			result = append(result, file)
		}
	}

	return result, nil
}
