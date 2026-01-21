package storage

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

// FileLock provides cross-process file locking
type FileLock struct {
	lockFile string
	file     *os.File
}

// NewFileLock creates lock for given base path
func NewFileLock(basePath string) *FileLock {
	return &FileLock{
		lockFile: basePath + ".lock",
	}
}

// Lock acquires exclusive lock with timeout (default 10s)
func (fl *FileLock) Lock(ctx context.Context, timeout ...time.Duration) error {
	if err := os.MkdirAll(filepath.Dir(fl.lockFile), 0755); err != nil {
		return err
	}
	
	t := 10 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}
	
	deadline := time.Now().Add(t)
	for time.Now().Before(deadline) {
		// Check if context cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		file, err := os.OpenFile(fl.lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			fl.file = file
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return os.ErrDeadlineExceeded
}

// Unlock releases lock
func (fl *FileLock) Unlock() error {
	if fl.file != nil {
		fl.file.Close()
		fl.file = nil
	}
	err := os.Remove(fl.lockFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}