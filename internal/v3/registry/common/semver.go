package common

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// SemverHelper provides semantic version utilities for registries
type SemverHelper struct{}

// NewSemverHelper creates a new semantic version helper
func NewSemverHelper() *SemverHelper {
	return &SemverHelper{}
}

// IsSemverVersion checks if a version follows semantic versioning
func (s *SemverHelper) IsSemverVersion(version string) bool {
	normalized := strings.TrimPrefix(version, "v")
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+`, normalized)
	return matched
}

// IsHigherVersion compares two semantic versions
func (s *SemverHelper) IsHigherVersion(v1, v2 string) bool {
	major1, minor1, patch1, err1 := s.ParseVersion(v1)
	major2, minor2, patch2, err2 := s.ParseVersion(v2)
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

// ParseVersion parses a semantic version string
func (s *SemverHelper) ParseVersion(version string) (major, minor, patch int, err error) {
	version = strings.TrimPrefix(version, "v")
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 4 {
		return 0, 0, 0, fmt.Errorf("invalid version format")
	}

	major, _ = strconv.Atoi(matches[1])
	minor, _ = strconv.Atoi(matches[2])
	patch, _ = strconv.Atoi(matches[3])
	return
}

// SortVersionsBySemver sorts versions by semantic version in descending order
func (s *SemverHelper) SortVersionsBySemver(versions []string) []string {
	var semverVersions []string
	var otherVersions []string

	for _, version := range versions {
		if s.IsSemverVersion(version) {
			semverVersions = append(semverVersions, version)
		} else {
			otherVersions = append(otherVersions, version)
		}
	}

	// Sort semver versions by version (descending)
	sort.Slice(semverVersions, func(i, j int) bool {
		return s.IsHigherVersion(semverVersions[i], semverVersions[j])
	})

	// Combine semver versions first, then other versions
	result := make([]string, 0, len(versions))
	result = append(result, semverVersions...)
	result = append(result, otherVersions...)
	return result
}
