package lockmanager

import (
	"errors"
	"log"
	"time"

	"github.com/sascha-andres/lockutil/internal/lockmanager/types"

	"github.com/sascha-andres/lockutil/internal/lockmanager/inmemory"
)

// LockManager manages named locks with optional timeout waits.
type LockManager struct {

	// locker provides methods for acquiring and releasing locks, typically used by the LockManager to manage named locks.
	locker types.Locker

	// verbose indicates whether to log detailed information about lock operations.
	verbose bool
}

// NewLockManager creates a new LockManager instance
func NewLockManager(verbose bool) *LockManager {
	return &LockManager{
		locker:  inmemory.NewInMemoryLocker(),
		verbose: verbose,
	}
}

// RequestLock attempts to acquire a lock with the given name and PID, waiting up to timeoutSeconds.
func (lm *LockManager) RequestLock(name string, pid int32, addr string, timeoutSeconds int32) error {
	if timeoutSeconds < 0 {
		return errors.New("timeoutSeconds must be greater than or equal to 0")
	}
	waitDuration := time.Duration(timeoutSeconds) * time.Second
	timeout := time.After(waitDuration)
	ticker := time.NewTicker(100 * time.Millisecond) // Poll every 100 ms
	defer ticker.Stop()

	for {
		err := lm.locker.Lock(name, pid, addr)
		if err == nil {
			if lm.verbose {
				log.Printf("Acquired lock for %s from %s-%d", name, addr, pid)
			}
			return nil
		}
		if errors.Is(err, types.ErrLockExists) && timeoutSeconds == 0 {
			if lm.verbose {
				log.Printf("no lock for %s from %s-%d: already taken", name, addr, pid)
			}
			return nil
		}

		// Wait for the lock to be released or timeout
		select {
		case <-timeout:
			if lm.verbose {
				log.Printf("timeout before acquiring lock for %s from %s-%d", name, addr, pid)
			}
			return nil
		case <-ticker.C:
			// Retry acquiring the lock
		}
	}
}

// ReleaseLock releases the lock for the given name and PID.
func (lm *LockManager) ReleaseLock(name string, pid int32, addr string) error {
	return lm.locker.Unlock(name, pid, addr)
}

// GetLocks returns a slice of LockInfo representing all the current locks and their statuses.
func (lm *LockManager) GetLocks() []types.LockInfo {
	return lm.locker.GetLocks()
}
