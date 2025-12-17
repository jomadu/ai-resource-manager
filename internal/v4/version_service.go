package v4

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

// Build-time variables set by ldflags
var (
	VersionString = "dev"
	Commit        = "unknown"
	BuildTime     = "unknown"
)

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

func ParseVersion(versionString string) (Version, error) {
	versionString = strings.TrimPrefix(versionString, "v")
	parts := strings.Split(versionString, ".")
	if len(parts) < 3 {
		return Version{Version: versionString}, errors.New("invalid version format")
	}
	
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{Version: versionString}, err
	}
	
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{Version: versionString}, err
	}
	
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{Version: versionString}, err
	}
	
	return Version{
		Major:   major,
		Minor:   minor,
		Patch:   patch,
		Version: versionString,
	}, nil
}