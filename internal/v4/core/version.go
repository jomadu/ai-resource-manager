package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// IsSemanticVersion returns true if the version was successfully parsed as semantic version
func (v Version) IsSemanticVersion() bool {
	return v.Major > 0 || v.Minor > 0 || v.Patch > 0
}

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

// IsNewerThan returns true if v is newer than other
func (v Version) IsNewerThan(other Version) bool {
	return v.Compare(other) > 0
}

// IsOlderThan returns true if v is older than other
func (v Version) IsOlderThan(other Version) bool {
	return v.Compare(other) < 0
}

// Build-time variables set by ldflags
var (
	VersionString = "dev"
	Commit        = "unknown"
	BuildTime     = "unknown"
)

func ParseVersion(versionString string) (Version, error) {
	original := versionString
	
	// Remove 'v' prefix if present
	if strings.HasPrefix(versionString, "v") {
		versionString = versionString[1:]
	}
	
	// Split on '+' to separate build metadata
	parts := strings.Split(versionString, "+")
	versionPart := parts[0]
	build := ""
	if len(parts) > 1 {
		build = parts[1]
	}
	
	// Split on '-' to separate prerelease
	parts = strings.Split(versionPart, "-")
	corePart := parts[0]
	prerelease := ""
	if len(parts) > 1 {
		prerelease = strings.Join(parts[1:], "-")
	}
	
	// Parse major.minor.patch
	versionParts := strings.Split(corePart, ".")
	if len(versionParts) < 1 {
		// No version parts, return as-is
		return Version{Version: original}, nil
	}
	
	// Parse major (required)
	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return Version{Version: original}, nil
	}
	
	// Parse minor (optional)
	minor := 0
	if len(versionParts) > 1 {
		minor, err = strconv.Atoi(versionParts[1])
		if err != nil {
			return Version{Version: original}, nil
		}
	}
	
	// Parse patch (optional)
	patch := 0
	if len(versionParts) > 2 {
		patch, err = strconv.Atoi(versionParts[2])
		if err != nil {
			return Version{Version: original}, nil
		}
	}
	
	return Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: prerelease,
		Build:      build,
		Version:    original,
	}, nil
}

func ResolveVersion(versionStr string, availableVersions []Version) (Version, error) {
	constraint, err := ParseConstraint(versionStr)
	if err != nil {
		return Version{}, err
	}

	var candidates []Version
	for _, v := range availableVersions {
		if constraint.IsSatisfiedBy(v) {
			candidates = append(candidates, v)
		}
	}

	if len(candidates) == 0 {
		return Version{}, fmt.Errorf("no version satisfies constraint: %s", versionStr)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Compare(candidates[j]) > 0
	})

	return candidates[0], nil
}
