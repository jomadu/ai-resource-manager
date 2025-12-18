package core

import "runtime"

func GetBuildInfo() BuildInfo {
	// parse VersionString into Version
	version, err := ParseVersion(VersionString)
	if err != nil {
		// Return a default version if parsing fails
		return BuildInfo{
			Version:   Version{Version: VersionString},
			Commit:    Commit,
			BuildTime: BuildTime,
			Arch:      runtime.GOOS + "/" + runtime.GOARCH,
		}
	}
	return BuildInfo{
		Version:   version,
		Commit:    Commit,
		BuildTime: BuildTime,
		Arch:      runtime.GOOS + "/" + runtime.GOARCH,
	}
}
