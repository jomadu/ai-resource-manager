package core

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

// Compare returns -1 if v is older than other, 0 if equal, 1 if newer
func (v Version) Compare(other Version) int {
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
	return 0
}

// CompareTo returns -1 if v is older than other, 0 if equal, 1 if newer
// Returns error if either version is not semver
func (v Version) CompareTo(other Version) (int, error) {
	if !v.IsSemver {
		return 0, fmt.Errorf("cannot compare non-semver version: %s", v.Version)
	}
	if !other.IsSemver {
		return 0, fmt.Errorf("cannot compare non-semver version: %s", other.Version)
	}
	return v.Compare(other), nil
}

// IsNewerThan returns true if v is newer than other
func (v Version) IsNewerThan(other Version) (bool, error) {
	cmp, err := v.CompareTo(other)
	if err != nil {
		return false, err
	}
	return cmp > 0, nil
}

// IsOlderThan returns true if v is older than other
func (v Version) IsOlderThan(other Version) (bool, error) {
	cmp, err := v.CompareTo(other)
	if err != nil {
		return false, err
	}
	return cmp < 0, nil
}

// ToString returns the string representation of the version
func (v Version) ToString() string {
	return v.Version
}

// Build-time variables set by ldflags
var (
	VersionString = "dev"
	Commit        = "unknown"
	BuildTime     = "unknown"
)

var semverRegex = regexp.MustCompile(`^(v)?(\d+)\.(\d+)\.(\d+)(?:-([\.\w-]+))?(?:\+([\w.-]+))?$`)

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
		satisfied, err := constraint.IsSatisfiedBy(v)
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
		cmp, err := candidates[i].CompareTo(candidates[j])
		if err != nil {
			return false
		}
		return cmp > 0
	})

	return candidates[0], nil
}
