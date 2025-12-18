package core

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

import (
	"runtime"
)

// Build-time variables set by ldflags
var (
	VersionString = "dev"
	Commit        = "unknown"
	BuildTime     = "unknown"
)

func ParseVersion(versionString string) (Version, error) {
	// TODO implement:
	return Version{Version: versionString}, nil
}
