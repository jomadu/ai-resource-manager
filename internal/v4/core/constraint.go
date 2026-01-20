package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ConstraintType string

const (
	Exact  ConstraintType = "exact"
	Major  ConstraintType = "major"
	Minor  ConstraintType = "minor"
	Latest ConstraintType = "latest"
)

type Constraint struct {
	Type    ConstraintType
	Version *Version
}

// IsSatisfiedBy checks if a version satisfies this constraint
// Returns error if version is not semver (except for Latest constraint)
func (c *Constraint) IsSatisfiedBy(version Version) (bool, error) {
	if c == nil {
		return false, nil
	}
	
	// Latest accepts any version
	if c.Type == Latest {
		return true, nil
	}
	
	// All other constraint types require semver
	if !version.IsSemver {
		return false, fmt.Errorf("constraint requires semver version, got: %s", version.Version)
	}
	
	switch c.Type {
	case Exact:
		if c.Version == nil {
			return false, nil
		}
		return version.Compare(*c.Version) == 0, nil
	case Major:
		if c.Version == nil {
			return false, nil
		}
		return version.Major == c.Version.Major && version.Compare(*c.Version) >= 0, nil
	case Minor:
		if c.Version == nil {
			return false, nil
		}
		return version.Major == c.Version.Major && version.Minor == c.Version.Minor && version.Compare(*c.Version) >= 0, nil
	default:
		return false, nil
	}
}

// ToString returns the string representation of the constraint
func (c *Constraint) ToString() string {
	if c == nil {
		return ""
	}
	
	switch c.Type {
	case Latest:
		return "latest"
	case Exact:
		if c.Version != nil {
			return c.Version.ToString()
		}
		return ""
	case Major:
		if c.Version != nil {
			return fmt.Sprintf("^%s", c.Version.ToString())
		}
		return ""
	case Minor:
		if c.Version != nil {
			return fmt.Sprintf("~%s", c.Version.ToString())
		}
		return ""
	default:
		return ""
	}
}

var constraintRegex = regexp.MustCompile(`^(v)?(\d+)(?:\.(\d+))?(?:\.(\d+))?$`)

// NewConstraint creates a new Constraint from a constraint string
// Rejects non-semver inputs (except "latest")
func NewConstraint(versionStr string) (Constraint, error) {
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
		// Not a semver pattern - reject
		return Constraint{}, fmt.Errorf("invalid constraint: %s (must be semver or 'latest')", versionStr)
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

// ParseConstraint is deprecated, use NewConstraint instead
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
		// Branch name - legacy behavior
		version := Version{Version: rest, IsSemver: false}
		return Constraint{Type: Latest, Version: &version}, nil
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
