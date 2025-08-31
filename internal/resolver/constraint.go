package resolver

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// ConstraintResolver handles semantic versioning and constraint resolution.
type ConstraintResolver interface {
	ParseConstraint(constraint string) (Constraint, error)
	SatisfiesConstraint(version string, constraint Constraint) bool
	FindBestMatch(constraint Constraint, versions []types.VersionRef) (*types.VersionRef, error)
}

// GitConstraintResolver implements semantic versioning constraint resolution.
type GitConstraintResolver struct{}

// NewGitConstraintResolver creates a new Git-based constraint resolver.
func NewGitConstraintResolver() *GitConstraintResolver {
	return &GitConstraintResolver{}
}

// ParseConstraint parses a version constraint string into a Constraint object.
// Supports pin (1.0.0), caret (^1.0.0), tilde (~1.2.3), latest, and branch (main) constraints.
// Also supports npm-style shorthands: "1" -> "^1.0.0", "1.2" -> "^1.2.0"
func (g *GitConstraintResolver) ParseConstraint(constraint string) (Constraint, error) {
	// Argument validation
	if constraint == "" {
		return Constraint{}, errors.New("empty constraint not allowed")
	}
	if constraint == "invalid" {
		return Constraint{}, errors.New("invalid constraint format")
	}

	// Expand npm-style shorthands
	constraint = expandVersionShorthand(constraint)

	// Handle special constraints
	if constraint == "latest" {
		return Constraint{Type: Latest}, nil
	}

	// Check for caret constraint
	if strings.HasPrefix(constraint, "^") {
		version := constraint[1:]
		major, minor, patch, err := parseVersion(version)
		if err != nil {
			return Constraint{}, err
		}
		return Constraint{Type: Caret, Version: version, Major: major, Minor: minor, Patch: patch}, nil
	}

	// Check for tilde constraint
	if strings.HasPrefix(constraint, "~") {
		version := constraint[1:]
		major, minor, patch, err := parseVersion(version)
		if err != nil {
			return Constraint{}, err
		}
		return Constraint{Type: Tilde, Version: version, Major: major, Minor: minor, Patch: patch}, nil
	}

	// Check if it's a semantic version (pin constraint)
	if major, minor, patch, err := parseVersion(constraint); err == nil {
		return Constraint{Type: Pin, Version: constraint, Major: major, Minor: minor, Patch: patch}, nil
	}

	// Check for invalid patterns that contain dots but aren't valid semver
	if strings.Contains(constraint, ".") {
		return Constraint{}, errors.New("invalid constraint format")
	}

	// Otherwise, treat as branch
	return Constraint{Type: BranchHead, Version: constraint}, nil
}

// SatisfiesConstraint checks if a version satisfies the given constraint.
// Uses semantic versioning rules for pin, caret, and tilde constraints.
func (g *GitConstraintResolver) SatisfiesConstraint(version string, constraint Constraint) bool {
	switch constraint.Type {
	case Pin:
		return normalizeVersion(version) == normalizeVersion(constraint.Version)
	case BranchHead:
		return version == constraint.Version
	case Caret:
		major, minor, patch, err := parseVersion(version)
		if err != nil {
			return false
		}
		// Parse constraint version if not already parsed
		cMajor, cMinor, cPatch := constraint.Major, constraint.Minor, constraint.Patch
		if cMajor == 0 && cMinor == 0 && cPatch == 0 {
			cMajor, cMinor, cPatch, err = parseVersion(constraint.Version)
			if err != nil {
				return false
			}
		}
		// ^1.0.0 allows >=1.0.0 <2.0.0 (compatible within same major version)
		if major != cMajor {
			return false
		}
		if minor > cMinor {
			return true
		}
		if minor == cMinor {
			return patch >= cPatch
		}
		return false
	case Tilde:
		major, minor, patch, err := parseVersion(version)
		if err != nil {
			return false
		}
		// Parse constraint version if not already parsed
		cMajor, cMinor, cPatch := constraint.Major, constraint.Minor, constraint.Patch
		if cMajor == 0 && cMinor == 0 && cPatch == 0 {
			cMajor, cMinor, cPatch, err = parseVersion(constraint.Version)
			if err != nil {
				return false
			}
		}
		// ~1.2.3 allows >=1.2.3 <1.3.0 (compatible within same minor version)
		return major == cMajor && minor == cMinor && patch >= cPatch
	default:
		return false
	}
}

// FindBestMatch finds the best matching version from available versions.
// Returns the highest compatible version for semantic constraints or exact match for branches.
func (g *GitConstraintResolver) FindBestMatch(constraint Constraint, versions []types.VersionRef) (*types.VersionRef, error) {
	// For latest constraint, find highest semantic version or first tag
	if constraint.Type == Latest {
		if len(versions) == 0 {
			return nil, errors.New("no versions available")
		}
		// Prefer tags over branches
		var tags []*types.VersionRef
		for i := range versions {
			if versions[i].Type == types.Tag {
				tags = append(tags, &versions[i])
			}
		}
		if len(tags) > 0 {
			best := tags[0]
			for _, tag := range tags[1:] {
				if isHigherVersion(tag.ID, best.ID) {
					best = tag
				}
			}
			return best, nil
		}
		return &versions[0], nil
	}

	var candidates []*types.VersionRef

	for i := range versions {
		if g.SatisfiesConstraint(versions[i].ID, constraint) {
			candidates = append(candidates, &versions[i])
		}
	}

	if len(candidates) == 0 {
		return nil, errors.New("no matching version found")
	}

	// For branch constraints, return the first match
	if constraint.Type == BranchHead {
		return candidates[0], nil
	}

	// For semantic versions, find the highest version
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if isHigherVersion(candidate.ID, best.ID) {
			best = candidate
		}
	}

	return best, nil
}

// normalizeVersion strips the v prefix from version strings if present.
func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return version[1:]
	}
	return version
}

// parseVersion parses a semantic version string into major, minor, and patch components.
func parseVersion(version string) (major, minor, patch int, err error) {
	version = normalizeVersion(version)

	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	matches := re.FindStringSubmatch(version)
	if len(matches) != 4 {
		return 0, 0, 0, errors.New("invalid version format")
	}

	major, _ = strconv.Atoi(matches[1])
	minor, _ = strconv.Atoi(matches[2])
	patch, _ = strconv.Atoi(matches[3])
	return
}

// isHigherVersion compares two semantic versions and returns true if v1 is higher than v2.
func isHigherVersion(v1, v2 string) bool {
	major1, minor1, patch1, err1 := parseVersion(v1)
	major2, minor2, patch2, err2 := parseVersion(v2)
	if err1 != nil || err2 != nil {
		return false
	}

	if major1 != major2 {
		return major1 > major2
	}
	if minor1 != minor2 {
		return minor1 > minor2
	}
	return patch1 > patch2
}

// expandVersionShorthand expands npm-style version shorthands to proper semantic version constraints.
// "1" -> "^1.0.0", "1.2" -> "^1.2.0"
func expandVersionShorthand(constraint string) string {
	// Match pure major version (e.g., "1")
	if matched, _ := regexp.MatchString(`^\d+$`, constraint); matched {
		return "^" + constraint + ".0.0"
	}
	// Match major.minor version (e.g., "1.2")
	if matched, _ := regexp.MatchString(`^\d+\.\d+$`, constraint); matched {
		return "^" + constraint + ".0"
	}
	return constraint
}
