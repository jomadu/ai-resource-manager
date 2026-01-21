package storage

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestFileLock_BasicLockUnlock(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "test")
	
	lock := NewFileLock(basePath)
	
	// Should be able to acquire lock
	err := lock.Lock(context.Background())
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}
	
	// Lock file should exist
	lockFile := basePath + ".lock"
	if _, err := os.Stat(lockFile); os.IsNotExist(err) {
		t.Fatal("Lock file should exist after Lock()")
	}
	
	// Should be able to unlock
	err = lock.Unlock()
	if err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}
	
	// Lock file should be removed
	if _, err := os.Stat(lockFile); !os.IsNotExist(err) {
		t.Fatal("Lock file should be removed after Unlock()")
	}
}

func TestFileLock_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "concurrent")
	
	var wg sync.WaitGroup
	results := make(chan bool, 2)
	
	// Start two goroutines trying to acquire same lock
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock := NewFileLock(basePath)
			err := lock.Lock(context.Background())
			if err != nil {
				results <- false
				return
			}
			
			// Hold lock briefly
			time.Sleep(50 * time.Millisecond)
			
			lock.Unlock()
			results <- true
		}()
	}
	
	wg.Wait()
	close(results)
	
	// Both should succeed (one after the other)
	successCount := 0
	for success := range results {
		if success {
			successCount++
		}
	}
	
	if successCount != 2 {
		t.Fatalf("Expected 2 successful locks, got %d", successCount)
	}
}

func TestFileLock_Timeout(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "timeout")
	
	// First lock holds for longer than timeout
	lock1 := NewFileLock(basePath)
	err := lock1.Lock(context.Background())
	if err != nil {
		t.Fatalf("First lock failed: %v", err)
	}
	
	// Second lock should timeout quickly
	lock2 := NewFileLock(basePath)
	start := time.Now()
	err = lock2.Lock(context.Background(), 200 * time.Millisecond) // Fast timeout for test
	duration := time.Since(start)
	
	// Should fail with timeout
	if err == nil {
		t.Fatal("Second lock should have failed with timeout")
	}
	
	// Should timeout in reasonable time (less than 500ms)
	if duration > 500*time.Millisecond {
		t.Fatalf("Timeout took too long: %v", duration)
	}
	
	lock1.Unlock()
}

func TestFileLock_MultipleUnlock(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "multiple")
	
	lock := NewFileLock(basePath)
	
	// Lock and unlock
	lock.Lock(context.Background())
	err := lock.Unlock()
	if err != nil {
		t.Fatalf("First unlock failed: %v", err)
	}
	
	// Second unlock should not fail
	err = lock.Unlock()
	if err != nil {
		t.Fatalf("Second unlock should not fail: %v", err)
	}
}

func TestFileLock_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "cancel")
	
	// First lock holds
	lock1 := NewFileLock(basePath)
	err := lock1.Lock(context.Background())
	if err != nil {
		t.Fatalf("First lock failed: %v", err)
	}
	
	// Second lock with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	lock2 := NewFileLock(basePath)
	err = lock2.Lock(ctx)
	
	// Should fail with context error
	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled, got %v", err)
	}
	
	lock1.Unlock()
}