package v4

// IsSatisfiedBy checks if a version satisfies the constraint.
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
