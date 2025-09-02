package cache

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

var registryLocks sync.Map // map[string]*sync.Mutex

type RegistryLock struct {
	file *os.File
	mu   *sync.Mutex
}

func AcquireRegistryLock(registryDir string) (*RegistryLock, error) {
	// Get or create mutex for this registry
	mutex, _ := registryLocks.LoadOrStore(registryDir, &sync.Mutex{})
	mu := mutex.(*sync.Mutex)
	mu.Lock()

	// Create file lock
	lockFile := filepath.Join(registryDir, ".lock")
	if err := os.MkdirAll(registryDir, 0o755); err != nil {
		mu.Unlock()
		return nil, err
	}

	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		mu.Unlock()
		return nil, err
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		mu.Unlock()
		return nil, err
	}

	return &RegistryLock{file: file, mu: mu}, nil
}

func (rl *RegistryLock) Release() error {
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func (rl *RegistryLock) ReleaseIgnoreError() {
	defer rl.mu.Unlock()
	_ = rl.file.Close()
}
