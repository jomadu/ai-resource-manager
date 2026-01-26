package core

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

// Compare returns -1 if v is older than other, 0 if equal, 1 if newer
// Implements semver precedence rules including prerelease comparison
func (v *Version) Compare(other *Version) int {
	if v.Major < other.Major {
		return -1
	}
	if v.Major > other.Major {
		return 1
	}
	if v.Minor < other.Minor {
		return -1
	}
	if v.Minor > other.Minor {
		return 1
	}
	if v.Patch < other.Patch {
		return -1
	}
	if v.Patch > other.Patch {
		return 1
	}

	// Compare prerelease versions per semver spec
	return comparePrerelease(v.Prerelease, other.Prerelease)
}

// comparePrerelease compares two prerelease strings per semver spec
// Returns -1 if a < b, 0 if equal, 1 if a > b
func comparePrerelease(a, b string) int {
	// Per semver: version without prerelease > version with prerelease
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1 // no prerelease > prerelease
	}
	if b == "" {
		return -1 // prerelease < no prerelease
	}

	// Both have prereleases - compare identifiers
	aParts := splitPrerelease(a)
	bParts := splitPrerelease(b)

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		cmp := comparePrereleaseIdentifier(aParts[i], bParts[i])
		if cmp != 0 {
			return cmp
		}
	}

	// All compared identifiers are equal - longer prerelease wins
	if len(aParts) < len(bParts) {
		return -1
	}
	if len(aParts) > len(bParts) {
		return 1
	}
	return 0
}

// splitPrerelease splits a prerelease string by dots
func splitPrerelease(s string) []string {
	if s == "" {
		return nil
	}
	parts := []string{}
	current := ""
	for _, ch := range s {
		if ch == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// comparePrereleaseIdentifier compares two prerelease identifiers
// Per semver: numeric identifiers are compared as integers,
// alphanumeric identifiers are compared lexically,
// numeric identifiers always have lower precedence than non-numeric
func comparePrereleaseIdentifier(a, b string) int {
	aNum, aIsNum := parseNumericIdentifier(a)
	bNum, bIsNum := parseNumericIdentifier(b)

	if aIsNum && bIsNum {
		// Both numeric - compare as integers
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}

	if aIsNum {
		return -1 // numeric < alphanumeric
	}
	if bIsNum {
		return 1 // alphanumeric > numeric
	}

	// Both alphanumeric - compare lexically
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// parseNumericIdentifier attempts to parse a string as a numeric identifier
// Returns the integer value and true if successful, 0 and false otherwise
func parseNumericIdentifier(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	// Per semver: numeric identifiers must not have leading zeros
	if len(s) > 1 && s[0] == '0' {
		return 0, false
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return val, true
}

// CompareTo returns -1 if v is older than other, 0 if equal, 1 if newer
// Returns error if either version is not semver
func (v *Version) CompareTo(other *Version) (int, error) {
	if !v.IsSemver {
		return 0, fmt.Errorf("cannot compare non-semver version: %s", v.Version)
	}
	if !other.IsSemver {
		return 0, fmt.Errorf("cannot compare non-semver version: %s", other.Version)
	}
	return v.Compare(other), nil
}

// IsNewerThan returns true if v is newer than other
func (v *Version) IsNewerThan(other *Version) (bool, error) {
	cmp, err := v.CompareTo(other)
	if err != nil {
		return false, err
	}
	return cmp > 0, nil
}

// IsOlderThan returns true if v is older than other
func (v *Version) IsOlderThan(other *Version) (bool, error) {
	cmp, err := v.CompareTo(other)
	if err != nil {
		return false, err
	}
	return cmp < 0, nil
}

// ToString returns the string representation of the version
func (v *Version) ToString() string {
	return v.Version
}

// Build-time variables set by ldflags
var (
	VersionString = "dev"
	Commit        = "unknown"
	BuildTime     = "unknown"
)

var semverRegex = regexp.MustCompile(`^(v)?(\d+)\.(\d+)\.(\d+)(?:-([.\w-]+))?(?:\+([\w.-]+))?$`)

// NewVersion creates a new Version from a version string
func NewVersion(versionString string) (Version, error) {
	if versionString == "" {
		return Version{}, fmt.Errorf("version string cannot be empty")
	}

	// Try to match semver pattern
	matches := semverRegex.FindStringSubmatch(versionString)
	if matches == nil {
		// Not semver, return as plain version string
		return Version{Version: versionString, IsSemver: false}, nil
	}

	// Parse semver groups
	major, _ := strconv.Atoi(matches[2])
	minor, _ := strconv.Atoi(matches[3])
	patch, _ := strconv.Atoi(matches[4])
	prerelease := matches[5]
	build := matches[6]

	return Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: prerelease,
		Build:      build,
		Version:    versionString,
		IsSemver:   true,
	}, nil
}

// ParseVersion is deprecated, use NewVersion instead
func ParseVersion(versionString string) (Version, error) {
	return NewVersion(versionString)
}

func ResolveVersion(versionStr string, availableVersions []Version) (Version, error) {
	constraint, err := ParseConstraint(versionStr)
	if err != nil {
		return Version{}, err
	}

	var candidates []Version
	for _, v := range availableVersions {
		satisfied, err := constraint.IsSatisfiedBy(&v)
		if err != nil {
			// Skip versions that don't match constraint requirements (e.g., non-semver for semver constraints)
			continue
		}
		if satisfied {
			candidates = append(candidates, v)
		}
	}

	if len(candidates) == 0 {
		return Version{}, fmt.Errorf("no version satisfies constraint: %s", versionStr)
	}

	sort.Slice(candidates, func(i, j int) bool {
		cmp, err := candidates[i].CompareTo(&candidates[j])
		if err != nil {
			return false
		}
		return cmp > 0
	})

	return candidates[0], nil
}
