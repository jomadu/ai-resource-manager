package version

import "runtime"

// Build-time variables set by ldflags
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// GetVersionInfo returns version information including runtime architecture
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   Version,
		Commit:    Commit,
		Timestamp: BuildTime,
		Arch:      runtime.GOOS + "/" + runtime.GOARCH,
	}
}

// VersionInfo provides version information about the ARM tool itself.
type VersionInfo struct {
	Arch      string `json:"arch"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Timestamp string `json:"timestamp"`
}
