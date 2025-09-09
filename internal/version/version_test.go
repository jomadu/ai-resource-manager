package version

import (
	"runtime"
	"testing"
)

func TestGetVersionInfo(t *testing.T) {
	info := GetVersionInfo()

	if info.Version == "" {
		t.Error("GetVersionInfo() Version should not be empty")
	}

	if info.Commit == "" {
		t.Error("GetVersionInfo() Commit should not be empty")
	}

	if info.Timestamp == "" {
		t.Error("GetVersionInfo() Timestamp should not be empty")
	}

	expectedArch := runtime.GOOS + "/" + runtime.GOARCH
	if info.Arch != expectedArch {
		t.Errorf("GetVersionInfo() Arch = %v, want %v", info.Arch, expectedArch)
	}
}

func TestVersionInfoStruct(t *testing.T) {
	info := VersionInfo{
		Version:   "1.0.0",
		Commit:    "abc123",
		Timestamp: "2023-01-01T00:00:00Z",
		Arch:      "linux/amd64",
	}

	if info.Version != "1.0.0" {
		t.Errorf("VersionInfo.Version = %v, want 1.0.0", info.Version)
	}

	if info.Commit != "abc123" {
		t.Errorf("VersionInfo.Commit = %v, want abc123", info.Commit)
	}

	if info.Timestamp != "2023-01-01T00:00:00Z" {
		t.Errorf("VersionInfo.Timestamp = %v, want 2023-01-01T00:00:00Z", info.Timestamp)
	}

	if info.Arch != "linux/amd64" {
		t.Errorf("VersionInfo.Arch = %v, want linux/amd64", info.Arch)
	}
}
