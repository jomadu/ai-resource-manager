package cache

import (
	"os"
	"path/filepath"
)

// GetCacheDir returns the base ARM cache directory path.
func GetCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".arm", "cache")
}

// GetRegistriesDir returns the ARM registries cache directory path.
func GetRegistriesDir() string {
	return filepath.Join(GetCacheDir(), "registries")
}
