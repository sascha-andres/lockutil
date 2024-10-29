package lockmanager

import (
	"errors"
	"sync"
	"time"
)

// LockManager manages named locks with optional timeout waits.
type LockManager struct {
	mu    sync.Mutex
	locks map[string]*lockInfo
}

type lockInfo struct {
	pid      int32
	isLocked bool
}

// NewLockManager creates a new LockManager instance
func NewLockManager() *LockManager {
	return &LockManager{
		locks: make(map[string]*lockInfo),
	}
}

// RequestLock attempts to acquire a lock with the given name and PID, waiting up to timeoutSeconds.
func (lm *LockManager) RequestLock(name string, pid int32, timeoutSeconds int32) error {
	waitDuration := time.Duration(timeoutSeconds) * time.Second
	timeout := time.After(waitDuration)
	ticker := time.NewTicker(100 * time.Millisecond) // Poll every 100 ms
	defer ticker.Stop()

	for {
		lm.mu.Lock()
		lock, exists := lm.locks[name]
		if !exists || (exists && !lock.isLocked) {
			// Acquire lock if it does not exist or is not currently locked
			lm.locks[name] = &lockInfo{pid: pid, isLocked: true}
			lm.mu.Unlock()
			return nil
		}
		lm.mu.Unlock()

		// Wait for the lock to be released or timeout
		select {
		case <-timeout:
			return errors.New("timeout exceeded while waiting to acquire lock")
		case <-ticker.C:
			// Retry acquiring the lock
		}
	}
}

// ReleaseLock releases the lock for the given name and PID.
func (lm *LockManager) ReleaseLock(name string, pid int32) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lock, exists := lm.locks[name]; exists && lock.isLocked && lock.pid == pid {
		// Only release if the PID matches the lock holder's PID
		delete(lm.locks, name)
		return nil
	}
	return errors.New("lock not held by given PID or does not exist")
}
