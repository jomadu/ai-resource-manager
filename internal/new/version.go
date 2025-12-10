package new

// Compare returns -1 if v is older than other, 0 if equal, 1 if newer
func (v Version) Compare(other Version) int {
    // Compare Major, Minor, Patch
    // Handle prerelease/build metadata if needed
}

// IsNewerThan returns true if v is newer than other
func (v Version) IsNewerThan(other Version) bool {
    return v.Compare(other) > 0
}

// IsOlderThan returns true if v is older than other
func (v Version) IsOlderThan(other Version) bool {
    return v.Compare(other) < 0
}