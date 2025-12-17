package cache

import (
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

func WithRegistryLock(registryKey string, fn func() error) error {
	locksDir := filepath.Join(GetRegistriesDir(), ".locks")
	if err := os.MkdirAll(locksDir, 0o755); err != nil {
		return err
	}

	lockFile := filepath.Join(locksDir, registryKey+".lock")
	fileLock := flock.New(lockFile)
	if err := fileLock.Lock(); err != nil {
		return err
	}
	defer func() { _ = fileLock.Unlock() }()

	return fn()
}
