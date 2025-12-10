package new

import "runtime"

// Build-time variables set by ldflags
var (
	VersionString   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func GetBuildInfo() BuildInfo {
	// parse VersionString into Version
	version, err := ParseVersion(VersionString)
	if err != nil {
		return BuildInfo{
			Version:   VersionString,
			Commit:    Commit,
			BuildTime: BuildTime,
			Arch:      runtime.GOOS + "/" + runtime.GOARCH,
		}
	}
	return BuildInfo{
		Version:   VersionString,
		Commit:    Commit,
		BuildTime: BuildTime,
		Arch:      runtime.GOOS + "/" + runtime.GOARCH,
	}
}

func ParseVersion(versionString string) (Version, error) {
}