package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ConstraintType string

const (
	Exact      ConstraintType = "exact"
	Major      ConstraintType = "major"
	Minor      ConstraintType = "minor"
	BranchHead ConstraintType = "branch-head"
	Latest     ConstraintType = "latest"
)

type Constraint struct {
	Type    ConstraintType
	Version *Version
}

func (c *Constraint) IsSatisfiedBy(version Version) bool {
	if c == nil {
		return false
	}
	switch c.Type {
	case Exact:
		if c.Version == nil {
			return false
		}
		return version.Compare(*c.Version) == 0
	case Major:
		if c.Version == nil {
			return false
		}
		return version.Major == c.Version.Major && version.Compare(*c.Version) >= 0
	case Minor:
		if c.Version == nil {
			return false
		}
		return version.Major == c.Version.Major && version.Minor == c.Version.Minor && version.Compare(*c.Version) >= 0
	case BranchHead:
		return version.Version == c.Version.Version
	case Latest:
		return true
	default:
		return false
	}
}

var constraintRegex = regexp.MustCompile(`^(v)?(\d+)(?:\.(\d+))?(?:\.(\d+))?$`)

func ParseConstraint(versionStr string) (Constraint, error) {
	if versionStr == "latest" {
		return Constraint{Type: Latest}, nil
	}

	// Strip prefix
	prefix := ""
	rest := versionStr
	if strings.HasPrefix(versionStr, "^") {
		prefix = "^"
		rest = versionStr[1:]
	} else if strings.HasPrefix(versionStr, "~") {
		prefix = "~"
		rest = versionStr[1:]
	}

	// Try to match semver/abbreviated pattern
	matches := constraintRegex.FindStringSubmatch(rest)
	if matches == nil {
		// Not a semver pattern
		if prefix != "" {
			return Constraint{}, fmt.Errorf("invalid constraint: %s (prefix requires version)", versionStr)
		}
		// Branch name
		version := Version{Version: rest, IsSemver: false}
		return Constraint{Type: BranchHead, Version: &version}, nil
	}

	// Parse and expand semver
	hasV := matches[1] != ""
	major, _ := strconv.Atoi(matches[2])
	minor := 0
	if matches[3] != "" {
		minor, _ = strconv.Atoi(matches[3])
	}
	patch := 0
	if matches[4] != "" {
		patch, _ = strconv.Atoi(matches[4])
	}

	// Build version string
	versionString := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	if hasV {
		versionString = "v" + versionString
	}

	version := Version{
		Major:    major,
		Minor:    minor,
		Patch:    patch,
		Version:  versionString,
		IsSemver: true,
	}

	// Determine constraint type
	switch prefix {
	case "^":
		return Constraint{Type: Major, Version: &version}, nil
	case "~":
		return Constraint{Type: Minor, Version: &version}, nil
	default:
		if patch > 0 {
			return Constraint{Type: Exact, Version: &version}, nil
		}
		if minor > 0 {
			return Constraint{Type: Minor, Version: &version}, nil
		}
		return Constraint{Type: Major, Version: &version}, nil
	}
}
