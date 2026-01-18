package core

import "strings"

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

func ParseConstraint(versionStr string) (Constraint, error) {
	if versionStr == "latest" {
		return Constraint{Type: Latest}, nil
	}

	// Handle ^ (caret) - allows changes within same major version
	if strings.HasPrefix(versionStr, "^") {
		version, err := ParseVersion(versionStr[1:])
		if err != nil {
			return Constraint{}, err
		}
		if !version.IsSemver {
			return Constraint{Type: BranchHead, Version: &version}, nil
		}
		return Constraint{Type: Major, Version: &version}, nil
	}

	// Handle ~ (tilde) - allows patch-level changes
	if strings.HasPrefix(versionStr, "~") {
		version, err := ParseVersion(versionStr[1:])
		if err != nil {
			return Constraint{}, err
		}
		if !version.IsSemver {
			return Constraint{Type: BranchHead, Version: &version}, nil
		}
		return Constraint{Type: Minor, Version: &version}, nil
	}

	version, err := ParseVersion(versionStr)
	if err != nil {
		return Constraint{}, err
	}

	if !version.IsSemver {
		return Constraint{Type: BranchHead, Version: &version}, nil
	}

	if version.Patch > 0 {
		return Constraint{Type: Exact, Version: &version}, nil
	}
	if version.Minor > 0 {
		return Constraint{Type: Minor, Version: &version}, nil
	}
	return Constraint{Type: Major, Version: &version}, nil
}
