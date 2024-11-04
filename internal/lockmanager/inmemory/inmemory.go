package inmemory

import (
	"sync"

	"github.com/sascha-andres/lockutil/internal/lockmanager/types"
)

// InMemoryLocker is a simple in-memory lock manager.
type InMemoryLocker struct {
	mu    sync.Mutex
	locks map[string]*lockInfo
}

// Lock attempts to acquire a lock with the given name for the specified pid.
// Returns ErrLockExists if the lock is already held.
func (i *InMemoryLocker) Lock(name string, pid int32, addr string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	lock, exists := i.locks[name]
	if !exists || (exists && !lock.isLocked) {
		// Acquire lock if it does not exist or is not currently locked
		i.locks[name] = &lockInfo{pid: pid, isLocked: true, addr: addr}
		return nil
	}
	return types.ErrLockExists
}

// Unlock attempts to release a lock identified by the name for the given pid.
// Returns ErrStrangersLock if the lock is held by a different PID or does not exist.
func (i *InMemoryLocker) Unlock(name string, pid int32, addr string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if lock, exists := i.locks[name]; exists && lock.isLocked && lock.pid == pid && lock.addr == addr {
		// Only release if the PID matches the lock holder's PID
		delete(i.locks, name)
		return nil
	}
	return types.ErrStrangersLock
}

// lockInfo represents the lock status and the process ID (pid) holding the lock.
type lockInfo struct {

	// pid represents the process ID holding the lock.
	pid int32

	// addr represents the address associated with the lock.
	addr string

	// isLocked indicates whether the lock is currently held by a process.
	isLocked bool
}

// NewInMemoryLocker creates and initializes a new InMemoryLocker instance.
func NewInMemoryLocker() *InMemoryLocker {
	return &InMemoryLocker{
		locks: make(map[string]*lockInfo),
	}
}
