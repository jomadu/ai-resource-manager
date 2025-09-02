package cache

import (
	"os"
	"path/filepath"
)

// GetCacheDir returns the ARM cache directory path.
func GetCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".arm", "cache", "registries")
}
