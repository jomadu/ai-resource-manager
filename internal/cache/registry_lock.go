package cache

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

var registryLocks sync.Map // map[string]*sync.Mutex

func WithRegistryLock(registryKey string, fn func() error) error {
	// Get or create mutex for this registry
	mutex, _ := registryLocks.LoadOrStore(registryKey, &sync.Mutex{})
	mu := mutex.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	// Create file lock in dedicated .locks subdirectory
	registriesDir := GetRegistriesDir()
	locksDir := filepath.Join(registriesDir, ".locks")
	lockFile := filepath.Join(locksDir, registryKey+".lock")
	if err := os.MkdirAll(locksDir, 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}

	return fn()
}
