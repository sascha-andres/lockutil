package lockmanager

import (
	"errors"
	"time"

	"github.com/sascha-andres/lockutil/internal/lockmanager/types"

	"github.com/sascha-andres/lockutil/internal/lockmanager/inmemory"
)

// LockManager manages named locks with optional timeout waits.
type LockManager struct {

	// locker provides methods for acquiring and releasing locks, typically used by the LockManager to manage named locks.
	locker types.Locker
}

// NewLockManager creates a new LockManager instance
func NewLockManager() *LockManager {
	return &LockManager{
		locker: inmemory.NewInMemoryLocker(),
	}
}

// RequestLock attempts to acquire a lock with the given name and PID, waiting up to timeoutSeconds.
func (lm *LockManager) RequestLock(name string, pid int32, timeoutSeconds int32) error {
	if timeoutSeconds < 0 {
		return errors.New("timeoutSeconds must be greater than or equal to 0")
	}
	waitDuration := time.Duration(timeoutSeconds) * time.Second
	timeout := time.After(waitDuration)
	ticker := time.NewTicker(100 * time.Millisecond) // Poll every 100 ms
	defer ticker.Stop()

	for {
		err := lm.locker.Lock(name, pid)
		if err == nil {
			return nil
		}
		if errors.Is(err, types.ErrLockExists) {
			return err
		}

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
	return lm.locker.Unlock(name, pid)
}
